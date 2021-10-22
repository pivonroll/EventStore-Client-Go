package statistics

import (
	"github.com/pivonroll/EventStore-Client-Go/core/errors"
)

//go:generate mockgen -source=statistics.go -destination=../internal/statistics/statistics_client_sync_mock.go -package=statistics

// Reader is a reader interface for projection's statistics.
type Reader interface {
	// Read reads one statistics result for a projection.
	Read() (Response, errors.Error)
}

// Status is a status of the projection.
type Status string

const (
	// StatusAborted projection is aborted.
	StatusAborted Status = "Stopped"
	// StatusAbortedOrStopped projection is either aborted or stopped.
	StatusAbortedOrStopped Status = "Aborted/Stopped"
	// StatusStopped projection is stopped.
	StatusStopped Status = "Stopped"
	// StatusRunning projection is running.
	StatusRunning Status = "Running"
)

// ModeOneTime means that mode od projection is one-time.
const ModeOneTime = "OneTime"

// Response is a result of a read operation of the statistics reader.
type Response struct {
	CoreProcessingTime                 int64
	Version                            int64
	Epoch                              int64
	EffectiveName                      string
	WritesInProgress                   int32
	ReadsInProgress                    int32
	PartitionsCached                   int32
	Status                             string
	StateReason                        string
	Name                               string
	Mode                               string
	Position                           string
	Progress                           float32
	LastCheckpoint                     string
	EventsProcessedAfterRestart        int64
	CheckpointStatus                   string
	BufferedEvents                     int64
	WritePendingEventsBeforeCheckpoint int32
	WritePendingEventsAfterCheckpoint  int32
}
