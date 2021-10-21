package statistics

import (
	"github.com/pivonroll/EventStore-Client-Go/projections/statistics"
	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type ClientSyncFactoryImpl struct{}

func (factory ClientSyncFactoryImpl) Create(
	statisticsClient projections.Projections_StatisticsClient) statistics.ClientSync {
	return newStatisticsClientSyncImpl(statisticsClient)
}
