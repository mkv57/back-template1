package queue

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
)

var _ Message = &natsMessage{}

type natsMessage struct {
	decoder Decoder
	msg     *nats.Msg
}

// ID implements Message.
func (m *natsMessage) ID() uuid.UUID {
	return uuid.Must(uuid.FromString(m.msg.Header.Get(msgIDHeader)))
}

// Subject implements Message.
func (m *natsMessage) Subject() string {
	return m.msg.Subject
}

// Ack implements Message.
func (m *natsMessage) Ack(ctx context.Context) error {
	return m.msg.Ack(nats.Context(ctx))
}

// Nack implements Message.
func (m *natsMessage) Nack(ctx context.Context) error {
	return m.msg.Nak(nats.Context(ctx))
}

// Unmarshal implements Message.
func (m *natsMessage) Unmarshal(a any) error {
	err := m.decoder.Unmarshal(m.msg.Data, a)
	if err != nil {
		return fmt.Errorf("m.decoder.Unmarshal: %w", err)
	}

	event, ok := a.(Validator)
	if !ok {
		return nil
	}

	err = event.ValidateAll()
	if err != nil {
		return fmt.Errorf("event.ValidateAll: %w", err)
	}

	return nil
}
