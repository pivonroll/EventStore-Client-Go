package statistics

//go:generate mockgen -source=factory.go -destination=factory_mock.go -package=statistics

import (
	"github.com/pivonroll/EventStore-Client-Go/projections/statistics"
	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type ReaderFactory interface {
	Create(client projections.Projections_StatisticsClient) statistics.Reader
}
