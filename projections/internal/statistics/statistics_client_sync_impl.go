package statistics

import (
	"io"

	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/pivonroll/EventStore-Client-Go/core/errors"
	"github.com/pivonroll/EventStore-Client-Go/projections/statistics"
	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type ClientSyncImpl struct {
	client             projections.Projections_StatisticsClient
	readRequestChannel chan chan statisticsReadResult
}

type statisticsReadResult struct {
	statisticsClientResponse statistics.Response
	err                      errors.Error
}

func (this *ClientSyncImpl) Read() (statistics.Response, errors.Error) {
	channel := make(chan statisticsReadResult)

	this.readRequestChannel <- channel
	resp := <-channel

	return resp.statisticsClientResponse, resp.err
}

func (statisticsSync *ClientSyncImpl) readOne() (statistics.Response, errors.Error) {
	result, protoErr := statisticsSync.client.Recv()
	if protoErr != nil {
		if protoErr == io.EOF {
			return statistics.Response{}, errors.NewError(errors.EndOfStream, protoErr)
		}
		trailer := statisticsSync.client.Trailer()
		err := connection.GetErrorFromProtoException(trailer, protoErr)
		if err != nil {
			return statistics.Response{}, err
		}
		return statistics.Response{}, errors.NewError(errors.FatalError, protoErr)
	}

	return statistics.Response{
		CoreProcessingTime:                 result.Details.CoreProcessingTime,
		Version:                            result.Details.Version,
		Epoch:                              result.Details.Epoch,
		EffectiveName:                      result.Details.EffectiveName,
		WritesInProgress:                   result.Details.WritesInProgress,
		ReadsInProgress:                    result.Details.ReadsInProgress,
		PartitionsCached:                   result.Details.PartitionsCached,
		Status:                             result.Details.Status,
		StateReason:                        result.Details.StateReason,
		Name:                               result.Details.Name,
		Mode:                               result.Details.Mode,
		Position:                           result.Details.Position,
		Progress:                           result.Details.Progress,
		LastCheckpoint:                     result.Details.LastCheckpoint,
		EventsProcessedAfterRestart:        result.Details.EventsProcessedAfterRestart,
		CheckpointStatus:                   result.Details.CheckpointStatus,
		BufferedEvents:                     result.Details.BufferedEvents,
		WritePendingEventsBeforeCheckpoint: result.Details.WritePendingEventsBeforeCheckpoint,
		WritePendingEventsAfterCheckpoint:  result.Details.WritePendingEventsAfterCheckpoint,
	}, nil
}

func (statisticsSync *ClientSyncImpl) readLoop() {
	for {
		responseChannel := <-statisticsSync.readRequestChannel
		result, err := statisticsSync.readOne()

		response := statisticsReadResult{
			statisticsClientResponse: result,
			err:                      err,
		}

		responseChannel <- response
	}
}

func newStatisticsClientSyncImpl(client projections.Projections_StatisticsClient) *ClientSyncImpl {
	statisticsReadClient := &ClientSyncImpl{
		client:             client,
		readRequestChannel: make(chan chan statisticsReadResult),
	}

	go statisticsReadClient.readLoop()

	return statisticsReadClient
}
