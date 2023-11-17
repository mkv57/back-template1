package client_test

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	pb "github.com/ZergsLaw/back-template/api/session/v1"
	"github.com/ZergsLaw/back-template/cmd/session/client"
	"github.com/ZergsLaw/back-template/internal/testhelper"
)

var (
	traceID = xid.New()
	errAny  = errors.New("any err")
)

func start(t *testing.T) (context.Context, *client.Client, *MockSessionInternalAPIServer, *require.Assertions) {
	t.Helper()
	assert := require.New(t)
	ctx := testhelper.Context(t)
	ctrl := gomock.NewController(t)
	mock := NewMockSessionInternalAPIServer(ctrl)

	srv := grpc.NewServer()
	pb.RegisterSessionInternalAPIServer(srv, mock)
	addr := testhelper.Addr(t, assert)
	ln, err := net.Listen("tcp", addr.String())
	assert.NoError(err)
	go func() { assert.NoError(srv.Serve(ln)) }()
	t.Cleanup(srv.Stop)
	reg := prometheus.NewPedanticRegistry()
	namespace := testhelper.Namespace(t)

	log := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{ //nolint:exhaustruct
				AddSource: true,
				Level:     slog.LevelDebug,
			},
		),
	)

	svc, err := client.New(ctx, log, reg, namespace, addr.String())
	assert.NoError(err)

	return ctx, svc, mock, assert
}
