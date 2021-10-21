package statistics

import (
	"github.com/pivonroll/EventStore-Client-Go/errors"
	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

//go:generate mockgen -source=statistics.go -destination=../internal/statistics/statistics_client_sync_mock.go -package=statistics

type ClientSync interface {
	Read() (Response, errors.Error)
}

type ClientSyncFactory interface {
	Create(client projections.Projections_StatisticsClient) ClientSync
}

const (
	StatusAborted = "Stopped"
	StatusStopped = "Aborted/Stopped"
	StatusRunning = "Running"
)

const ModeOneTime = "OneTime"

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
