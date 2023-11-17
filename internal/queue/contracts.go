package queue

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/proto"
)

// Contracts for encoding/decoding messages.
type (
	// Encoder marshals any values.
	Encoder interface {
		// Marshal serializes any value and returns bytes.
		Marshal(any) ([]byte, error)
	}
	// Decoder unmarshal bytes to values.
	Decoder interface {
		// Unmarshal deserialize bytes to message.
		Unmarshal([]byte, any) error
	}
)

// Validator for validation values.
type Validator interface {
	// ValidateAll validate all fields.
	ValidateAll() error
}

// Message contains any message.
type Message interface {
	// ID returns message id for deduplication.
	ID() uuid.UUID
	// Subject returns message subject.
	Subject() string
	// Unmarshal payload to message.
	// Notice:
	//  * If message has a method ValidateAll(), will call it.
	Unmarshal(any) error
	// Ack acknowledges a message.
	Ack(context.Context) error
	// Nack negatively acknowledges a message.
	Nack(context.Context) error
}

var (
	_ Encoder = &encoderProto{}
	_ Decoder = &decoderProto{}
)

type (
	encoderProto struct{ *proto.MarshalOptions }
	decoderProto struct{ *proto.UnmarshalOptions }
)

var ErrIncorrectMessage = errors.New("incorrect message")

// Marshal implements Encoder.
func (e *encoderProto) Marshal(a any) ([]byte, error) {
	msg, ok := a.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("%w [%T]: %+v ", ErrIncorrectMessage, a, a)
	}

	return e.MarshalOptions.Marshal(msg)
}

// Unmarshal implements Decoder.
func (d *decoderProto) Unmarshal(buf []byte, a any) error {
	msg, ok := a.(proto.Message)
	if !ok {
		return fmt.Errorf("%w: [%T]", ErrIncorrectMessage, a)
	}

	return d.UnmarshalOptions.Unmarshal(buf, msg)
}
