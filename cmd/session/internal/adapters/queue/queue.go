package queue

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"

	user_pb "github.com/ZergsLaw/back-template/api/user/v1"
	"github.com/ZergsLaw/back-template/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template/internal/dom"
	"github.com/ZergsLaw/back-template/internal/queue"
)

var _ app.Queue = &Client{}

type (
	// Config provide connection info for message broker.
	Config struct {
		URLs     []string
		Username string
		Password string
	}
	// Client provided data from and to message broker.
	Client struct {
		consumerName string
		queue        *queue.Queue
		m            Metrics
		chUpStatus   chan dom.Event[app.UpdateStatus]
	}
)

// New build and returns new queue instance.
func New(ctx context.Context, reg *prometheus.Registry, namespace string, cfg Config) (*Client, error) {
	const subsystem = "queue"
	m := NewMetrics(reg, namespace, subsystem, []string{})

	client, err := queue.Connect(ctx, strings.Join(cfg.URLs, ","), namespace, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("queue.Connect: %w", err)
	}

	err = client.Migrate(migrate(namespace))
	if err != nil {
		return nil, fmt.Errorf("client.Migrate: %w", err)
	}

	return &Client{
		consumerName: namespace,
		queue:        client,
		m:            m,
		chUpStatus:   make(chan dom.Event[app.UpdateStatus]),
	}, nil
}

// UpSessionStatus implements app.Queue.
func (c *Client) UpSessionStatus() <-chan dom.Event[app.UpdateStatus] {
	return c.chUpStatus
}

// Process starts worker for collecting events from queue.
func (c *Client) Process(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	subjects := []string{
		user_pb.TopicUpdate,
	}

	for i := range subjects {
		i := i
		group.Go(func() error {
			return c.queue.Subscribe(ctx, subjects[i], c.consumerName, c.handleEvent)
		})
	}

	return group.Wait()
}

func (c *Client) handleEvent(ctx context.Context, msg queue.Message) error {
	ack := make(chan dom.AcknowledgeKind)

	var err error
	switch {
	case ctx.Err() != nil:
		return nil
	case msg.Subject() == user_pb.TopicUpdate:
		err = c.handleUpStatus(ctx, ack, msg.ID(), msg)
	default:
		err = fmt.Errorf("%w: unknown topic %s", app.ErrInvalidArgument, msg.Subject())
	}
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return nil
	case ackKind := <-ack:
		switch ackKind {
		case dom.AcknowledgeKindAck:
			err = msg.Ack(ctx)
		case dom.AcknowledgeKindNack:
			err = msg.Nack(ctx)
		}
		if err != nil {
			return fmt.Errorf("msg.Ack|Nack: %w", err)
		}
	}

	return nil
}

func (c *Client) handleUpStatus(ctx context.Context, ack chan dom.AcknowledgeKind, msgID uuid.UUID, msg queue.Message) error {
	event := &user_pb.Event{}
	err := msg.Unmarshal(event)
	if err != nil {
		return fmt.Errorf("msg.Unmarshal: %w", err)
	}

	updateEvent := event.GetUpdate()
	if err != nil {
		return fmt.Errorf("%w: event.GetUpdate: %+v", queue.ErrIncorrectMessage, event.GetBody())
	}

	arg := dom.NewEvent(msgID, ack, app.UpdateStatus{
		UserID: uuid.Must(uuid.FromString(updateEvent.User.Id)),
		Status: dom.UserStatusFromAPI(updateEvent.User.Kind),
	})

	select {
	case <-ctx.Done():
		return nil
	case c.chUpStatus <- *arg:
	}

	return nil
}

// Close implements io.Closer.
func (c *Client) Close() error {
	return c.queue.Drain()
}

// Monitor for starting monitor connection background logic.
func (c *Client) Monitor(ctx context.Context) error {
	return c.queue.Monitor(ctx)
}

const (
	ackWait       = time.Second * 5
	maxDeliver    = 5
	maxAckPending = 1
)

func migrate(namespace string) func(manager nats.JetStreamManager) error {
	return func(manager nats.JetStreamManager) error {
		err := user_pb.Migrate(manager)
		if err != nil {
			return fmt.Errorf("user.Migrate: %w", err)
		}

		_, err = manager.AddConsumer(user_pb.Stream, &nats.ConsumerConfig{
			Durable:       namespace,
			Description:   "Consumer for updating user's session status.",
			DeliverPolicy: nats.DeliverAllPolicy,
			AckPolicy:     nats.AckExplicitPolicy,
			AckWait:       ackWait,
			MaxDeliver:    maxDeliver,
			BackOff:       []time.Duration{time.Second / 10, time.Second / 2, time.Second},
			MaxAckPending: maxAckPending,
		})
		if err != nil && !errors.Is(err, nats.ErrConsumerNameAlreadyInUse) {
			return fmt.Errorf("manager.AddConsumer: %w", err)
		}

		return nil
	}
}
