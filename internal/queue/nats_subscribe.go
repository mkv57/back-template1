package queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"

	"github.com/ZergsLaw/back-template1/internal/logger"
)

// Subscribe starts subscription by args.
func (c *Queue) Subscribe(
	ctx context.Context,
	subj, consumerName string,
	handler func(context.Context, Message) error,
) error {
	log := logger.FromContext(ctx)

	sub, err := c.jetStream.PullSubscribe(subj, consumerName, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("c.nats.jetStream.QueueSubscribe: %w", err)
	}
	defer func() {
		err := sub.Drain()
		if err != nil {
			log.Error("couldn't drain sub", slog.String(logger.Error.String(), err.Error()))
		}
	}()

	for {
		msgs, err := sub.Fetch(1, nats.Context(ctx))
		switch {
		case ctx.Err() != nil:
			return nil
		case errors.Is(err, context.DeadlineExceeded):
			continue // Because fetcher can return context error by default timeout.
		case err != nil:
			return fmt.Errorf("sub.Fetch: %w", err)
		}

		for i := range msgs {
			if ctx.Err() != nil {
				return nil
			}

			err = handler(ctx, &natsMessage{decoder: c.decoder, msg: msgs[i]})
			if err != nil {
				log.Error("couldn't handle message", slog.String(logger.Error.String(), err.Error()))
			}
		}
	}
}
