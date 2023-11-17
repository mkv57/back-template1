package queue

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/types/known/timestamppb"

	user_pb "github.com/ZergsLaw/back-template/api/user/v1"
	"github.com/ZergsLaw/back-template/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template/internal/dom"
	"github.com/ZergsLaw/back-template/internal/queue"
)

var _ app.Queue = &Client{}

type (
	// Config provide connection info for message broker.
	Config struct {
		URLs        []string
		Username    string
		Password    string
		ClusterMode bool
	}
	// Client provided data from and to message broker.
	Client struct {
		queue *queue.Queue
		m     Metrics
	}
)

// New build and returns new queue instance.
func New(ctx context.Context, reg *prometheus.Registry, namespace string, cfg Config) (*Client, error) {
	const subsystem = "queue"
	m := NewMetrics(reg, namespace, subsystem, []string{})

	client, err := queue.Connect(ctx, strings.Join(cfg.URLs, ","), namespace, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("queue.ConnectNATS: %w", err)
	}

	err = client.Migrate(user_pb.Migrate)
	if err != nil {
		return nil, fmt.Errorf("client.Migrate: %w", err)
	}

	return &Client{
		queue: client,
		m:     m,
	}, nil
}

// AddUser implements app.Queue.
func (c *Client) AddUser(ctx context.Context, id uuid.UUID, user app.User) error {
	return c.queue.Publish(ctx,
		user_pb.TopicAdd,
		id,
		&user_pb.Event{
			Body: &user_pb.Event_Add{
				Add: &user_pb.Add{
					User: &user_pb.User{
						Id:        user.ID.String(),
						Username:  user.Name,
						Email:     user.Email,
						AvatarId:  user.AvatarID.String(),
						Kind:      dom.UserStatusToAPI(user.Status),
						FullName:  user.FullName,
						CreatedAt: timestamppb.New(user.CreatedAt),
						UpdatedAt: timestamppb.New(user.UpdatedAt),
					},
				},
			},
		},
	)
}

// DeleteUser implements app.Queue.
func (c *Client) DeleteUser(ctx context.Context, id uuid.UUID, user app.User) error {
	return c.queue.Publish(ctx,
		user_pb.TopicDel,
		id,
		&user_pb.Event{
			Body: &user_pb.Event_Delete{
				Delete: &user_pb.Delete{
					UserId: user.ID.String(),
				},
			},
		},
	)
}

// UpdateUser implements app.Queue.
func (c *Client) UpdateUser(ctx context.Context, id uuid.UUID, user app.User) error {
	return c.queue.Publish(ctx,
		user_pb.TopicUpdate,
		id,
		&user_pb.Event{
			Body: &user_pb.Event_Update{
				Update: &user_pb.Update{
					User: &user_pb.User{
						Id:        user.ID.String(),
						Username:  user.Name,
						Email:     user.Email,
						AvatarId:  user.AvatarID.String(),
						Kind:      dom.UserStatusToAPI(user.Status),
						FullName:  user.FullName,
						CreatedAt: timestamppb.New(user.CreatedAt),
						UpdatedAt: timestamppb.New(user.UpdatedAt),
					},
				},
			},
		},
	)
}

// Close implements io.Closer.
func (c *Client) Close() error {
	return c.queue.Drain()
}

// Monitor for starting monitor connection background logic.
func (c *Client) Monitor(ctx context.Context) error {
	return c.queue.Monitor(ctx)
}
