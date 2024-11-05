// Package client provide to internal method of service session.
package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ZergsLaw/back-template1/api/session/v1"
	"github.com/ZergsLaw/back-template1/internal/grpchelper"

	"github.com/ZergsLaw/back-template1/internal/dom"
)

// Errors.
var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInternal        = errors.New("internal error")
)

// Client to session microservice.
type Client struct {
	conn pb.SessionInternalAPIClient
}

// New build and returns new client to microservice session.
func New(ctx context.Context, logger *slog.Logger, reg *prometheus.Registry, namespace, addr string) (*Client, error) {
	const subsystem = "session_client"
	clientMetric := grpchelper.NewClientMetrics(reg, namespace, subsystem)

	conn, err := grpchelper.Dial(ctx, addr, logger, clientMetric,
		[]grpc.UnaryClientInterceptor{},
		[]grpc.StreamClientInterceptor{},
		[]grpc.DialOption{},
	)
	if err != nil {
		return nil, fmt.Errorf("grpc_helper.Dial: %w", err)
	}

	return &Client{conn: pb.NewSessionInternalAPIClient(conn)}, nil
}

type (
	// Session contains main session info.
	Session struct {
		ID     uuid.UUID
		UserID uuid.UUID
		Status dom.UserStatus
	}

	// Token contains user's authorization token.
	Token struct {
		Value string
	}
)

// Save user's session.
func (c *Client) Save(
	ctx context.Context,
	userID uuid.UUID,
	origin dom.Origin,
	status dom.UserStatus,
) (*Token, error) {
	res, err := c.conn.Save(ctx, &pb.SaveRequest{
		UserId:    userID.String(),
		Ip:        origin.IP.String(),
		UserAgent: origin.UserAgent,
		Kind:      dom.UserStatusToAPI(status),
	})
	if err != nil {
		return nil, convertError(err)
	}

	return &Token{Value: res.Token}, nil
}

// Get user's session by his auth token.
func (c *Client) Get(ctx context.Context, token string) (*Session, error) {
	res, err := c.conn.Get(ctx, &pb.GetRequest{
		Token: token,
	})
	if err != nil {
		return nil, convertError(err)
	}

	userUID, err := uuid.FromString(res.UserId)
	if err != nil {
		return nil, fmt.Errorf("uuid.FromString: %w", err)
	}

	sessionUID, err := uuid.FromString(res.SessionId)
	if err != nil {
		return nil, fmt.Errorf("uuid.FromString: %w", err)
	}

	return &Session{
		ID:     sessionUID,
		UserID: userUID,
		Status: dom.UserStatusFromAPI(res.Kind),
	}, nil
}

// Delete remove user session by session ID.
func (c *Client) Delete(ctx context.Context, sessionID uuid.UUID) error {
	_, err := c.conn.Delete(ctx, &pb.DeleteRequest{
		SessionId: sessionID.String(),
	})
	if err != nil {
		return convertError(err)
	}

	return nil
}

func convertError(err error) error {
	switch {
	case status.Code(err) == codes.NotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, err)
	case status.Code(err) == codes.DeadlineExceeded:
		return fmt.Errorf("%w: %s", context.DeadlineExceeded, err)
	case status.Code(err) == codes.Canceled:
		return fmt.Errorf("%w: %s", context.Canceled, err)
	case status.Code(err) == codes.InvalidArgument:
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err)
	default:
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		return fmt.Errorf("%w: %s", ErrInternal, st.Message())
	}
}
