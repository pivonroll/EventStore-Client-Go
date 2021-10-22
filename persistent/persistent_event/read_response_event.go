package persistent_event

import (
	"time"

	"github.com/google/uuid"
	"github.com/pivonroll/EventStore-Client-Go/core/position"
)

// ReadResponseEvent is an event received from a stream.
// Each event has either Event or Link set.
// If event has no commit position CommitPosition will be nil.
type ReadResponseEvent struct {
	Event          *RecordedEvent
	Link           *RecordedEvent
	CommitPosition *uint64 // nil if NoCommit Position is received
	RetryCount     *int32  // nil if NoRetryCount is received
}

// GetOriginalEvent returns an original event.
// It chooses between Link and Event fields.
// Link field has precedence over Event field.
func (responseEvent ReadResponseEvent) GetOriginalEvent() *RecordedEvent {
	if responseEvent.Link != nil {
		return responseEvent.Link
	}

	return responseEvent.Event
}

// RecordedEvent represents an event recorded in the EventStoreDB.
type RecordedEvent struct {
	EventId        uuid.UUID         // ID of an event. Event's ID is provided by user when event is appended to a stream
	EventType      string            // user defined event type
	ContentType    string            // content type used to store event in EventStoreDB. Supported types are application/json and application/octet-stream
	StreamId       string            // stream identifier of a stream on which this event is stored
	EventNumber    uint64            // index number of an event in a stream
	Position       position.Position // event's position in stream $all
	CreatedDate    time.Time         // a date and time when event was stored in a stream
	Data           []byte            // user data stored in an event
	SystemMetadata map[string]string // EventStoreDB's metadata set for an event
	UserMetadata   []byte            // user defined metadata
}
