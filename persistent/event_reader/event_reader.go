package event_reader

//go:generate mockgen -source=event_reader.go -destination=../internal/event_reader_mock/event_reader_mock.go -package=event_reader_mock

import (
	"github.com/pivonroll/EventStore-Client-Go/errors"
	"github.com/pivonroll/EventStore-Client-Go/persistent/persistent_action"
	"github.com/pivonroll/EventStore-Client-Go/persistent/persistent_event"
)

// EventReader provides a reader for a persistent subscription.
// Clients must acknowledge (or not acknowledge) messages in the competing consumer model.
// If you enable auto-ack the subscription will automatically acknowledge messages once
// your handler completes them.
type EventReader interface {
	// ReadOne reads one message from persistent subscription.
	// Reads will block one the buffer size for persistent subscription is reached.
	// Buffer size is specified when we are connecting to a persistent subscription.
	ReadOne() (persistent_event.ReadResponseEvent, errors.Error) // this call must block
	// Ack sends Ack signal for a message.
	// Maximum of 2000 messages can be acknowledged at once.
	Ack(msgs ...persistent_event.ReadResponseEvent) errors.Error
	// Nack sends Nack signal for a message.
	// Client must also specify a reason why message was nack-ed.
	Nack(reason string, action persistent_action.Nack_Action, msgs ...persistent_event.ReadResponseEvent) error
	// Close closes the connection to a persistent subscription.
	Close()
}
