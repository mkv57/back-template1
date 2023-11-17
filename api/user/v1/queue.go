package pb

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// Topics.
const (
	description          = "Events from user service for notifying about new registration, deleting or updating account."
	Stream               = "user"
	events               = Stream + ".events"
	version              = events + ".v1."
	TopicAdd             = version + "add"
	TopicDel             = version + "del"
	TopicUpdate          = version + "update"
	SubscribeToAllEvents = version + "*"
)

const (
	maxMsgReplicas  = 1
	duplicateWindow = time.Second * 30
)

// Migrate for init streams.
func Migrate(js nats.JetStreamManager) error {
	replicas := maxMsgReplicas
	eventStream := &nats.StreamConfig{
		Name:        Stream,
		Description: description,
		Subjects:    []string{TopicAdd, TopicDel, TopicUpdate},
		Retention:   nats.LimitsPolicy,
		Storage:     nats.FileStorage,
		Replicas:    replicas,
		NoAck:       false,
		Duplicates:  duplicateWindow,
	}

	_, err := js.AddStream(eventStream)
	switch {
	case errors.Is(err, nats.ErrStreamNameAlreadyInUse):
		_, err = js.UpdateStream(eventStream)
		if err != nil {
			return fmt.Errorf("js.UpdateStream: %w", err)
		}

		return nil
	case err != nil:
		return fmt.Errorf("js.AddStream: %w", err)
	}

	return nil
}
