// Package api contains all methods for working grpc server.
package api

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	pb "github.com/ZergsLaw/back-template/api/session/v1"
	"github.com/ZergsLaw/back-template/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template/internal/dom"
	"github.com/ZergsLaw/back-template/internal/grpchelper"
	"github.com/ZergsLaw/back-template/internal/logger"
	"github.com/ZergsLaw/back-template/internal/metrics"
)

// For convenient testing.
type application interface {
	NewSession(ctx context.Context, userID uuid.UUID, status dom.UserStatus, origin app.Origin) (*app.Token, error)
	Session(ctx context.Context, token string) (*app.Session, error)
	RemoveSession(ctx context.Context, sessionID uuid.UUID) error
}

type api struct {
	app application
}

// New creates and returns gRPC server.
func New(ctx context.Context, m metrics.Metrics, app application, reg *prometheus.Registry, namespace string) *grpc.Server {
	log := logger.FromContext(ctx)
	subsystem := "api"

	grpcMetrics := grpchelper.NewServerMetrics(reg, namespace, subsystem)

	srv, health := grpchelper.NewServer(m, log, grpcMetrics, apiError,
		[]grpc.UnaryServerInterceptor{},
		[]grpc.StreamServerInterceptor{},
	)
	health.SetServingStatus(pb.SessionInternalAPI_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	pb.RegisterSessionInternalAPIServer(srv, &api{app: app})

	return srv
}

func apiError(err error) *status.Status {
	if err == nil {
		return nil
	}

	code := codes.Internal
	switch {
	case errors.Is(err, app.ErrNotFound):
		code = codes.NotFound
	case errors.Is(err, app.ErrInvalidToken):
		code = codes.InvalidArgument
	case errors.Is(err, context.DeadlineExceeded):
		code = codes.DeadlineExceeded
	case errors.Is(err, context.Canceled):
		code = codes.Canceled
	}

	return status.New(code, err.Error())
}
