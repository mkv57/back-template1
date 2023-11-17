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

	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	session_pb "github.com/ZergsLaw/back-template/api/session/v1"
	pb "github.com/ZergsLaw/back-template/api/user/v1"
	"github.com/ZergsLaw/back-template/internal/grpchelper"
	"github.com/ZergsLaw/back-template/internal/logger"
	"github.com/ZergsLaw/back-template/internal/metrics"
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

const (
	username1  = `username`
	fullName   = `full name`
	email      = `email@email.com`
	pass1      = `11111111`
	username2  = `username2`
	email2     = `email2@email.com`
	pass2      = `22222222`
	avatarPath = `testdata/test.jpg`
)

func initService(t *testing.T, ctx context.Context) (*require.Assertions, pb.UserExternalAPIClient, *queue.Queue, *config) {
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

	sessionSrv := newSessionServiceMock(ctx, t, reg, namespace)
	sessionSrvAddr := testhelper.Addr(t, assert)
	ln, err := net.Listen("tcp", sessionSrvAddr.String())
	assert.NoError(err)
	go func() { assert.NoError(sessionSrv.Serve(ln)) }()
	t.Cleanup(sessionSrv.GracefulStop)

	cockroachCfg := testhelper.CockroachDB(
		ctx,
		t,
		assert,
		filepath.Join(pwd, caCrtPath), filepath.Join(pwd, caKeyPath),
		filepath.Join(pwd, nodeCrtPath), filepath.Join(pwd, nodeKeyPath),
		filepath.Join(pwd, clientCrtPath), filepath.Join(pwd, clientKeyPath),
	)

	username := "test_svc"
	pass := "test_pass"
	natsURL := testhelper.NATS(
		ctx,
		t,
		assert,
		filepath.Join(pwd, queueCfgPath),
		username, pass,
	)

	endpoint := testhelper.Minio(
		ctx,
		t,
		assert,
		username, pass, "", false,
		"local-1",
	)

	cfg := config{
		Server: server{
			Host: testhelper.Host,
			Port: ports{
				GRPC:   testhelper.UnusedTCPPort(t, assert, testhelper.Host),
				Metric: testhelper.UnusedTCPPort(t, assert, testhelper.Host),
				GW:     testhelper.UnusedTCPPort(t, assert, testhelper.Host),
				Files:  testhelper.UnusedTCPPort(t, assert, testhelper.Host),
			},
		},
		Queue: queueConfig{
			URLs:     []string{natsURL},
			Username: queueUsername,
			Password: queuePassword,
		},
		DB: dbConfig{
			MigrateDir: migrateDir,
			Driver:     "postgres",
			Cockroach:  *cockroachCfg,
		},
		Clients: clients{
			Session: sessionSrvAddr.String(),
		},
		FileStore: fileStoreConfig{
			Secure:    false,
			Endpoint:  endpoint,
			AccessKey: username,
			SecretKey: pass,
			Region:    "local-1",
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

	return assert, pb.NewUserExternalAPIClient(conn), n, &cfg
}

var _ session_pb.SessionInternalAPIServer = &sessionServiceMock{}

type sessionServiceMock struct {
	values map[string]struct {
		sessionID string
		req       *session_pb.SaveRequest
	}
}

func newSessionServiceMock(
	ctx context.Context,
	t *testing.T,
	reg *prometheus.Registry,
	namespace string,
) *grpc.Server {
	t.Helper()
	log := logger.FromContext(ctx)
	subsystem := "session_api_mock"
	namespace = namespace + "_session_mock"

	grpcMetrics := grpchelper.NewServerMetrics(reg, namespace, subsystem)
	srv, _ := grpchelper.NewServer(metrics.New(prometheus.NewPedanticRegistry(), namespace), log, grpcMetrics, func(err error) *status.Status { return status.New(codes.Internal, err.Error()) },
		[]grpc.UnaryServerInterceptor{},
		[]grpc.StreamServerInterceptor{},
	)

	session_pb.RegisterSessionInternalAPIServer(srv, &sessionServiceMock{values: make(map[string]struct {
		sessionID string
		req       *session_pb.SaveRequest
	})})
	return srv
}

func (s *sessionServiceMock) Save(_ context.Context, request *session_pb.SaveRequest) (*session_pb.SaveResponse, error) {
	token := uuid.Must(uuid.NewV4()).String()
	sessionID := uuid.Must(uuid.NewV4()).String()

	s.values[token] = struct {
		sessionID string
		req       *session_pb.SaveRequest
	}{sessionID: sessionID, req: request}

	return &session_pb.SaveResponse{Token: token}, nil
}

func (s *sessionServiceMock) Get(_ context.Context, request *session_pb.GetRequest) (*session_pb.GetResponse, error) {
	val := s.values[request.Token]
	return &session_pb.GetResponse{
		SessionId: val.sessionID,
		UserId:    val.req.UserId,
		Kind:      val.req.Kind,
	}, nil
}

func (s *sessionServiceMock) Delete(_ context.Context, request *session_pb.DeleteRequest) (*session_pb.DeleteResponse, error) {
	for token, val := range s.values {
		if request.SessionId == val.sessionID {
			delete(s.values, token)
		}
	}

	return &session_pb.DeleteResponse{}, nil
}

func auth(ctx context.Context, token string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD{
		"authorization": {fmt.Sprintf("Bearer %s", token)},
	})
}
