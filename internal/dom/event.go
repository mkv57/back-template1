package dom

import (
	"context"

	"github.com/gofrs/uuid"
)

// AcknowledgeKind represents kind of acknowledgment.
type AcknowledgeKind uint8

//go:generate stringer -output=stringer.AcknowledgeKind.go -type=AcknowledgeKind -trimprefix=AcknowledgeKind
const (
	_ AcknowledgeKind = iota
	AcknowledgeKindAck
	AcknowledgeKindNack
)

// Event contains event message information.
type Event[T any] struct {
	id   uuid.UUID
	ack  chan AcknowledgeKind
	body T
}

// NewEvent build and returns new event from message broker.
func NewEvent[T any](
	id uuid.UUID,
	ack chan AcknowledgeKind,
	body T,
) *Event[T] {
	return &Event[T]{
		id:   id,
		ack:  ack,
		body: body,
	}
}

func (e *Event[T]) ID() uuid.UUID {
	return e.id
}

func (e *Event[T]) Ack(ctx context.Context) {
	select {
	case <-ctx.Done():
	case e.ack <- AcknowledgeKindAck:
	}
}

func (e *Event[T]) Nack(ctx context.Context) {
	select {
	case <-ctx.Done():
	case e.ack <- AcknowledgeKindNack:
	}
}

func (e *Event[T]) Body() T {
	return e.body
}
