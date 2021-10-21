// Header

// Package projections provides interaction with projections in EventStoreDB.
// Before accessing streams a grpc connection needs to be established with EventStore through
// github.com/pivonroll/EventStore-Client-Go/connection package.
package projections

import (
	"context"

	"github.com/pivonroll/EventStore-Client-Go/errors"
	"github.com/pivonroll/EventStore-Client-Go/projections/internal/grpc_proto_client_factory"
	statistics_internal "github.com/pivonroll/EventStore-Client-Go/projections/internal/statistics"
	"github.com/pivonroll/EventStore-Client-Go/projections/statistics"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/pivonroll/EventStore-Client-Go/connection"
)

type Client struct {
	grpcClient                   connection.GrpcClient
	grpcProjectionsClientFactory grpc_proto_client_factory.Factory
	statisticsClientSyncFactory  statistics.ClientSyncFactory
}

const FailedToCreateProjectionErr errors.ErrorCode = "FailedToCreateProjectionErr"

func (client *Client) CreateProjection(
	ctx context.Context,
	options CreateOptionsRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Create(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToCreateProjectionErr)
		return err
	}

	return nil
}

const FailedToUpdateProjectionErr errors.ErrorCode = "FailedToUpdateProjectionErr"

func (client *Client) UpdateProjection(
	ctx context.Context,
	options UpdateOptionsRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Update(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToUpdateProjectionErr)
		return err
	}

	return nil
}

const FailedToDeleteProjectionErr errors.ErrorCode = "FailedToDeleteProjectionErr"

func (client *Client) DeleteProjection(
	ctx context.Context,
	options DeleteOptionsRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Delete(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToDeleteProjectionErr)
		return err
	}

	return nil
}

const FailedToFetchProjectionStatisticsErr errors.ErrorCode = "FailedToFetchProjectionStatisticsErr"

func (client *Client) GetProjectionStatistics(
	ctx context.Context,
	options StatisticsOptionsRequest) (statistics.ClientSync, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD

	statisticsClient, protoErr := projectionsClient.Statistics(ctx, options.Build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToFetchProjectionStatisticsErr)
		return nil, err
	}

	return client.statisticsClientSyncFactory.Create(statisticsClient), nil
}

const FailedToDisableProjectionErr errors.ErrorCode = "FailedToDisableProjectionErr"

func (client *Client) DisableProjection(
	ctx context.Context,
	options DisableOptionsRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD

	_, protoErr := projectionsClient.Disable(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToDisableProjectionErr)
		return err
	}

	return nil
}

const FailedToAbortProjectionErr errors.ErrorCode = "FailedToAbortProjectionErr"

func (client *Client) AbortProjection(
	ctx context.Context,
	options AbortOptionsRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Disable(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToAbortProjectionErr)
		return err
	}

	return nil
}

const FailedToEnableProjectionErr errors.ErrorCode = "FailedToEnableProjectionErr"

func (client *Client) EnableProjection(
	ctx context.Context,
	options EnableOptionsRequest) errors.Error {

	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Enable(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToEnableProjectionErr)
		return err
	}

	return nil
}

const FailedToResetProjectionErr errors.ErrorCode = "FailedToResetProjectionErr"

func (client *Client) ResetProjection(
	ctx context.Context,
	options ResetOptionsRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Reset(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToResetProjectionErr)
		return err
	}

	return nil
}

const FailedToGetProjectionStateErr errors.ErrorCode = "FailedToGetProjectionStateErr"

func (client *Client) GetProjectionState(
	ctx context.Context,
	options StateOptionsRequest) (StateResponse, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	result, protoErr := projectionsClient.State(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToGetProjectionStateErr)
		return nil, err
	}

	return newStateResponse(result), nil
}

const FailedToGetProjectionResultErr errors.ErrorCode = "FailedToGetProjectionResultErr"

func (client *Client) GetProjectionResult(
	ctx context.Context,
	options ResultOptionsRequest) (ResultResponse, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	result, protoErr := projectionsClient.Result(ctx, options.Build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, FailedToGetProjectionResultErr)
		return nil, err
	}

	return newResultResponse(result), nil
}

const FailedToRestartProjectionsSubsystemErr errors.ErrorCode = "FailedToRestartProjectionsSubsystemErr"

func (client *Client) RestartProjectionsSubsystem(ctx context.Context) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.RestartSubsystem(ctx, &shared.Empty{},
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			FailedToRestartProjectionsSubsystemErr)
		return err
	}

	return nil
}

const (
	FailedToReadStatistics errors.ErrorCode = "FailedToReadStatistics"
)

func (client *Client) ListAllProjections(
	ctx context.Context) ([]statistics.Response, errors.Error) {

	options := StatisticsOptionsRequest{}
	options.SetMode(StatisticsOptionsRequestModeAll{})

	statisticsClient, err := client.GetProjectionStatistics(ctx, options)
	if err != nil {
		return nil, err
	}

	var result []statistics.Response

	for {
		statisticsResult, err := statisticsClient.Read()
		if err != nil {
			if err.Code() == errors.EndOfStream {
				break
			}

			return nil, err
		}

		result = append(result, statisticsResult)
	}

	return result, nil
}

func (client *Client) ListContinuousProjections(
	ctx context.Context) ([]statistics.Response, errors.Error) {
	options := StatisticsOptionsRequest{}
	options.SetMode(StatisticsOptionsRequestModeContinuous{})

	statisticsClient, err := client.GetProjectionStatistics(ctx, options)
	if err != nil {
		return nil, err
	}

	var result []statistics.Response

	for {
		statisticsResult, err := statisticsClient.Read()
		if err != nil {
			if err.Code() == errors.EndOfStream {
				break
			}
			return nil, err
		}

		result = append(result, statisticsResult)
	}

	return result, nil
}

func (client *Client) ListOneTimeProjections(
	ctx context.Context) ([]statistics.Response, errors.Error) {
	options := StatisticsOptionsRequest{}
	options.SetMode(StatisticsOptionsRequestModeOneTime{})

	statisticsClient, err := client.GetProjectionStatistics(ctx, options)
	if err != nil {
		return nil, err
	}

	var result []statistics.Response

	for {
		statisticsResult, err := statisticsClient.Read()
		if err != nil {
			if err.Code() == errors.EndOfStream {
				break
			}
			return nil, err
		}

		result = append(result, statisticsResult)
	}

	return result, nil
}

func NewClient(
	grpcClient connection.GrpcClient) *Client {
	return &Client{
		grpcProjectionsClientFactory: grpc_proto_client_factory.FactoryImpl{},
		grpcClient:                   grpcClient,
		statisticsClientSyncFactory:  statistics_internal.ClientSyncFactoryImpl{},
	}
}
