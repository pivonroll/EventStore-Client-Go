package persistent

import (
	"github.com/EventStore/EventStore-Client-Go/messages"
)

type AsyncReadConnection interface {
	Start()
	Stop() error
	Updates() <-chan messages.RecordedEvent
}
