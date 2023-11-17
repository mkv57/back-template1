package queue

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/ZergsLaw/back-template/internal/logger"
)

const (
	msgIDHeader   = `Nats-Msg-Id`
	drainTimeout  = 3 * time.Second // Should be less than main.shutdownDelay.
	maxReconnects = 5
	pingInterval  = time.Second // Default 2 min isn't useful because TCP keepalive is faster.
)

// AsyncErrMsg is error from async publish handler.
type AsyncErrMsg struct {
	msg *nats.Msg
	Err error
}

// Queue is queue connection.
type Queue struct {
	conn            *nats.Conn
	jetStream       nats.JetStreamContext
	closed          chan struct{}
	asyncErrHandler chan AsyncErrMsg // Non-blocking on send, closes by Queue.Close.
	encoder         Encoder
	decoder         Decoder
}

// Connect adds ctx support and reasonable defaults to nats.Connect.
func Connect(ctx context.Context, urls, namespace, username, password string) (*Queue, error) {
	log := logger.FromContext(ctx)

	c := &Queue{
		closed: make(chan struct{}),
		encoder: &encoderProto{
			MarshalOptions: &proto.MarshalOptions{},
		},
		decoder: &decoderProto{
			UnmarshalOptions: &proto.UnmarshalOptions{},
		},
	}

	var err error
	for !(c.conn != nil && err == nil) {
		err := c.connect(ctx, urls, namespace, username, password)
		switch {
		case err != nil:
			log.Error("couldn't connect to Queue", slog.String(logger.Error.String(), err.Error()))
		case ctx.Err() != nil:
			if err == nil {
				err = ctx.Err()
			}

			return nil, err
		}
	}

	log.Info("Queue connected", slog.String(logger.URL.String(), c.conn.ConnectedUrl()))

	return c, nil
}

func (c *Queue) connect(ctx context.Context, urls, namespace, username, password string) (err error) {
	log := logger.FromContext(ctx)

	c.conn, err = nats.Connect(urls,
		nats.Name(namespace),
		nats.UserInfo(username, password),
		nats.MaxReconnects(maxReconnects),
		nats.DrainTimeout(drainTimeout),
		nats.PingInterval(pingInterval),

		nats.NoCallbacksAfterClientClose(),
		nats.ClosedHandler(func(_ *nats.Conn) {
			close(c.closed)
		}),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err == nil {
				log.Info("Queue disconnected")
			} else {
				log.Warn("Queue disconnected", slog.String(logger.Error.String(), err.Error()))
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info("Queue reconnected", slog.String(logger.URL.String(), nc.ConnectedUrl()))
		}),
		nats.ErrorHandler(func(_ *nats.Conn, sub *nats.Subscription, err error) {
			if sub == nil {
				log.Warn("Queue connection failed", slog.String(logger.Error.String(), err.Error()))
			} else {
				log.Warn("Queue connection failed", slog.String(logger.Reason.String(), sub.Subject), slog.String(logger.Error.String(), err.Error()))
			}
		}),
	)
	if err != nil {
		return fmt.Errorf("nats.Connect: %w", err)
	}

	c.jetStream, err = c.conn.JetStream(
		nats.Context(ctx),
		nats.PublishAsyncErrHandler(func(_ nats.JetStream, msg *nats.Msg, err error) {
			c.asyncErrHandler <- AsyncErrMsg{
				msg: msg,
				Err: err,
			}
		}),
	)
	if err != nil {
		return fmt.Errorf("c.conn.jetStream: %w", err)
	}

	return nil
}

// Monitor waits until ctx.Done or failure reconnecting Queue.
func (c *Queue) Monitor(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case <-c.closed:
		return nats.ErrConnectionClosed
	}
}

// Err returns channel for handling async error.
func (c *Queue) Err() <-chan AsyncErrMsg {
	return c.asyncErrHandler
}

// Drain starts to close process.
func (c *Queue) Drain() error {
	return c.conn.Drain()
}

// Migrate starts callback with jetStream connection for making streams/consumers.
func (c *Queue) Migrate(f func(manager nats.JetStreamManager) error) error {
	return f(c.jetStream)
}
