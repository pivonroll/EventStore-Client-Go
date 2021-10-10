package persistent

//go:generate mockgen -source=event_reader.go -destination=event_reader_mock.go -package=persistent

import (
	"github.com/pivonroll/EventStore-Client-Go/errors"
)

type Nack_Action int32

const (
	Nack_Unknown Nack_Action = 0
	Nack_Park    Nack_Action = 1
	Nack_Retry   Nack_Action = 2
	Nack_Skip    Nack_Action = 3
	Nack_Stop    Nack_Action = 4
)

// EventReader provides a reader for a persistent subscription.
// Clients must acknowledge (or not acknowledge) messages in the competing consumer model.
// If you enable auto-ack the subscription will automatically acknowledge messages once
// your handler completes them.
type EventReader interface {
	// ReadOne reads one message from persistent subscription.
	// Reads will block one the buffer size for persistent subscription is reached.
	// Buffer size is specified when we are connecting to a persistent subscription.
	ReadOne() (ReadResponseEvent, errors.Error) // this call must block
	// Ack sends Ack signal for a message.
	Ack(msgs ...ReadResponseEvent) errors.Error // max 2000 messages can be acknowledged
	// Nack sends Nack signal for a message.
	// Client must also specify a reason why message was nack-ed.
	Nack(reason string, action Nack_Action, msgs ...ReadResponseEvent) error
	// Close closes the connection to a persistent subscription.
	Close()
}
