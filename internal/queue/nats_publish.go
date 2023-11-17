package queue

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/nats-io/nats.go"
)

// Publish send message to queue.
func (c *Queue) Publish(ctx context.Context, topic string, msgID uuid.UUID, event any) error {
	if event, ok := event.(Validator); ok {
		err := event.ValidateAll()
		if err != nil {
			return fmt.Errorf("event.ValidateAll: %w", err)
		}
	}

	buf, err := c.encoder.Marshal(event)
	if err != nil {
		return fmt.Errorf("c.encoder.Marshal: %w", err)
	}

	_, err = c.jetStream.Publish(
		topic,
		buf,
		nats.MsgId(msgID.String()),
		nats.Context(ctx),
	)
	if err != nil {
		return fmt.Errorf("c.jetStream.Publish: %w", err)
	}

	return nil
}
