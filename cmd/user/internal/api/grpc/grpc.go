// Package grpc contains all methods for working grpc server.
package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/gofrs/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	user_pb "github.com/ZergsLaw/back-template1/api/user/v1"
	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
	"github.com/ZergsLaw/back-template1/internal/grpchelper"
	"github.com/ZergsLaw/back-template1/internal/logger"
	"github.com/ZergsLaw/back-template1/internal/metrics"
)

// For convenient testing.
type application interface {
	VerificationEmail(ctx context.Context, email string) error
	VerificationUsername(ctx context.Context, username string) error
	CreateUser(ctx context.Context, email, username, fullName, password string) (uuid.UUID, error)
	Login(ctx context.Context, email, password string, origin dom.Origin) (uuid.UUID, *dom.Token, error)
	UserByID(ctx context.Context, session dom.Session, userID uuid.UUID) (*app.User, error)
	ListUserByFilters(ctx context.Context, _ dom.Session, filters app.SearchParams) ([]app.User, int, error)
	Logout(ctx context.Context, session dom.Session) error
	UpdatePassword(ctx context.Context, session dom.Session, oldPass, newPass string) error
	Auth(ctx context.Context, token string) (*dom.Session, error)
	UpdateUser(ctx context.Context, session dom.Session, username string, avatarID uuid.UUID) error
	RemoveAvatar(ctx context.Context, session dom.Session, fileID uuid.UUID) error
	ListUserAvatars(ctx context.Context, session dom.Session) ([]app.AvatarInfo, error)
	GetUsersByIDs(ctx context.Context, session dom.Session, ids []uuid.UUID) ([]app.User, error)
}

type api struct {
	app  application
	auth map[string]bool
}

// New creates and returns gRPC server.
func New(ctx context.Context, m metrics.Metrics, applications application, reg *prometheus.Registry, namespace string) *grpc.Server {
	log := logger.FromContext(ctx)
	subsystem := "api"

	grpcMetrics := grpchelper.NewServerMetrics(reg, namespace, subsystem)

	srv, health := grpchelper.NewServer(m, log, grpcMetrics, apiError,
		[]grpc.UnaryServerInterceptor{grpc_auth.UnaryServerInterceptor(nil)},   // Nil because we are using override.
		[]grpc.StreamServerInterceptor{grpc_auth.StreamServerInterceptor(nil)}, // Nil because we are using override.
	)
	health.SetServingStatus(user_pb.UserExternalAPI_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	user_pb.RegisterUserExternalAPIServer(srv, &api{
		app: applications,
		auth: map[string]bool{
			"VerificationEmail":    false,
			"VerificationUsername": false,
			"CreateUser":           false,
			"Login":                false,
			"GetUser":              true,
			"SearchUsers":          true,
			"Logout":               true,
			"UpdatePassword":       true,
			"UpdateUser":           true,
			"RemoveAvatar":         true,
			"ListUserAvatar":       true,
			"GetUsersByIDs":        true,
		},
	})

	return srv
}

func originFromCtx(ctx context.Context) (*dom.Origin, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("peer.FromContext: %w", app.ErrNotFound)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata.FromIncomingContext: %w", app.ErrNotFound)
	}

	host, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, fmt.Errorf("net.SplitHostPort: %w", err)
	}

	clientUserAgent := strings.Join(md.Get(userAgentForward), "")
	if clientUserAgent == "" {
		clientUserAgent = strings.Join(md.Get(userAgent), "")
	}

	return &dom.Origin{
		IP:        net.ParseIP(host),
		UserAgent: clientUserAgent,
	}, nil
}

func apiError(err error) *status.Status {
	if err == nil {
		return nil
	}

	code := codes.Internal
	switch {
	case errors.Is(err, app.ErrInvalidAuth):
		code = codes.Unauthenticated
	case errors.Is(err, app.ErrEmailExist):
		code = codes.AlreadyExists
	case errors.Is(err, app.ErrUsernameExist):
		code = codes.AlreadyExists
	case errors.Is(err, app.ErrAccessDenied):
		code = codes.PermissionDenied
	case errors.Is(err, app.ErrNotFound):
		code = codes.NotFound
	case errors.Is(err, app.ErrNotDifferent):
		code = codes.InvalidArgument
	case errors.Is(err, app.ErrInvalidPassword):
		code = codes.InvalidArgument
	case errors.Is(err, context.DeadlineExceeded):
		code = codes.DeadlineExceeded
	case errors.Is(err, context.Canceled):
		code = codes.Canceled
	}

	return status.New(code, err.Error())
}
