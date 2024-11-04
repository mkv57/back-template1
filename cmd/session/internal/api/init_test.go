package api_test

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/ZergsLaw/back-template1/api/session/v1"
	"github.com/ZergsLaw/back-template1/cmd/session/internal/api"
	"github.com/ZergsLaw/back-template1/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template1/internal/metrics"
	"github.com/ZergsLaw/back-template1/internal/testhelper"
)

var (
	errAny = errors.New("any err")
	origin = app.Origin{
		IP:        net.ParseIP("192.100.10.4"),
		UserAgent: "UserAgent",
	}
)

func start(t *testing.T) (context.Context, pb.SessionInternalAPIClient, *Mockapplication, *require.Assertions) {
	t.Helper()
	assert := require.New(t)
	ctx := testhelper.Context(t)

	ctrl := gomock.NewController(t)
	mockApp := NewMockapplication(ctrl)

	reg := prometheus.NewPedanticRegistry()
	namespace := testhelper.Namespace(t)

	m := metrics.New(reg, namespace)

	server := api.New(ctx, m, mockApp, reg, namespace)

	addr := testhelper.Addr(t, assert)
	ln, err := net.Listen("tcp", addr.String())
	assert.NoError(err)

	go func() {
		err := server.Serve(ln)
		assert.NoError(err)
	}()

	conn, err := grpc.DialContext(ctx, addr.String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO Add TLS and remove this.
		grpc.WithBlock(),
	)
	assert.NoError(err)

	t.Cleanup(func() {
		err := conn.Close()
		assert.NoError(err)
		server.GracefulStop()
	})

	return ctx, pb.NewSessionInternalAPIClient(conn), mockApp, assert
}

func randString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
