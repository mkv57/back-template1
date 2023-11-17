//go:build integration

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/ZergsLaw/back-template/api/session/v1"
	"github.com/ZergsLaw/back-template/internal/grpchelper"
	"github.com/ZergsLaw/back-template/internal/queue"
	"github.com/ZergsLaw/back-template/internal/testhelper"
)

const (
	queueCfgPath  = `testdata/nats.conf`
	queueUsername = `test_svc`
	queuePassword = `test_pass`
	migrateDir    = `./migrate`
	caCrtPath     = `../../certs/cockroach/ca.crt`
	caKeyPath     = `../../certs/cockroach/ca.key`
	nodeCrtPath   = `../../certs/cockroach/nodes/node1/node.crt`
	nodeKeyPath   = `../../certs/cockroach/nodes/node1/node.key`
	clientCrtPath = `../../certs/cockroach/client.root.crt`
	clientKeyPath = `../../certs/cockroach/client.root.key`
)

func initService(t *testing.T, ctx context.Context) (*require.Assertions, pb.SessionInternalAPIClient, *queue.Queue) {
	t.Helper()

	assert := require.New(t)
	pwd, err := os.Getwd()
	assert.NoError(err)

	reg := prometheus.NewPedanticRegistry()
	namespace := testhelper.Namespace(t)
	subsystem := testhelper.Namespace(t) + "_subsystem"

	devLogger := slog.New(slog.NewJSONHandler(
		os.Stderr, &slog.HandlerOptions{ //nolint:exhaustruct
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)
	grpclog.SetLoggerV2(grpchelper.NewLogger(devLogger))

	username := "test_svc"
	pass := "test_pass"
	natsURL := testhelper.NATS(
		ctx,
		t,
		assert,
		filepath.Join(pwd, queueCfgPath),
		username, pass,
	)

	cockroachCfg := testhelper.CockroachDB(
		ctx,
		t,
		assert,
		filepath.Join(pwd, caCrtPath), filepath.Join(pwd, caKeyPath),
		filepath.Join(pwd, nodeCrtPath), filepath.Join(pwd, nodeKeyPath),
		filepath.Join(pwd, clientCrtPath), filepath.Join(pwd, clientKeyPath),
	)

	cfg := config{
		AuthKey: "super-duper-secret-key-qwertyuio",
		Server: server{
			Host: testhelper.Host,
			Port: ports{
				GRPC:   testhelper.UnusedTCPPort(t, assert, testhelper.Host),
				Metric: testhelper.UnusedTCPPort(t, assert, testhelper.Host),
			},
		},
		DB: dbConfig{
			MigrateDir: migrateDir,
			Driver:     "postgres",
			Cockroach:  *cockroachCfg,
		},
		Queue: queueConfig{
			URLs:     []string{natsURL},
			Username: queueUsername,
			Password: queuePassword,
		},
	}
	addr := net.JoinHostPort(cfg.Server.Host, fmt.Sprintf("%d", cfg.Server.Port.GRPC))

	errc := make(chan error)
	ctxShutdown, shutdown := context.WithCancel(ctx)
	go func() { errc <- run(ctxShutdown, cfg, reg, namespace) }()
	t.Cleanup(func() {
		shutdown()
		assert.NoError(<-errc)
	})
	assert.NoError(testhelper.WaitTCPPort(ctx, addr))

	clientMetric := grpchelper.NewClientMetrics(reg, namespace, subsystem)
	conn, err := grpchelper.Dial(ctx, addr, devLogger, clientMetric,
		[]grpc.UnaryClientInterceptor{},
		[]grpc.StreamClientInterceptor{},
		[]grpc.DialOption{},
	)
	assert.NoError(err)

	n, err := queue.Connect(ctx, natsURL, namespace, queueUsername, queuePassword)
	assert.NoError(err)

	return assert, pb.NewSessionInternalAPIClient(conn), n
}
