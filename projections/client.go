// Header

// Package projections provides interaction with projections in EventStoreDB.
// Before accessing streams a grpc connection needs to be established with EventStore through
// github.com/pivonroll/EventStore-Client-Go/core/connection package.
package projections

import (
	"context"

	"github.com/pivonroll/EventStore-Client-Go/core/errors"
	"github.com/pivonroll/EventStore-Client-Go/projections/internal/grpc_proto_client_factory"
	statistics_internal "github.com/pivonroll/EventStore-Client-Go/projections/internal/statistics"
	"github.com/pivonroll/EventStore-Client-Go/projections/statistics"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/pivonroll/EventStore-Client-Go/core/connection"
)

// Client which can interact with EventStoreDB projections.
type Client struct {
	grpcClient                   connection.GrpcClient
	grpcProjectionsClientFactory grpc_proto_client_factory.Factory
	statisticsReaderFactory      statistics_internal.ReaderFactory
}

// CreateProjection creates a new projection on EventStoreDB.
func (client *Client) CreateProjection(
	ctx context.Context,
	options CreateRequest,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Create(ctx, options.build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// UpdateProjection updates an existing projection on EventStoreDB.
func (client *Client) UpdateProjection(
	ctx context.Context,
	options UpdateRequest,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Update(ctx, options.build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// DeleteProjection removes a projection from EventStoreDB.
func (client *Client) DeleteProjection(
	ctx context.Context,
	options DeleteRequest,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Delete(ctx, options.build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// GetProjectionStatistics returns a reader for projection's statistics.
func (client *Client) GetProjectionStatistics(
	ctx context.Context,
	mode IsStatisticsProjection,
) (statistics.Reader, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD

	statisticsClient, protoErr := projectionsClient.Statistics(ctx, buildStatisticsRequest(mode),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return nil, err
	}

	return client.statisticsReaderFactory.Create(statisticsClient), nil
}

// DisableProjection disables an existing projection.
func (client *Client) DisableProjection(
	ctx context.Context,
	projectionName string,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD

	_, protoErr := projectionsClient.Disable(ctx, disableOptionsRequest(projectionName).build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// AbortProjection aborts an existing projection.
func (client *Client) AbortProjection(
	ctx context.Context,
	projectionName string,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Disable(ctx, abortOptionsRequest(projectionName).build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// EnableProjection enables a disabled projection. If projection was already enabled it does noop (no operation).
func (client *Client) EnableProjection(
	ctx context.Context,
	projectionName string,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Enable(ctx, enableOptionsRequest(projectionName).build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// ResetProjection resets an existing projection.
func (client *Client) ResetProjection(
	ctx context.Context,
	options ResetRequest,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := projectionsClient.Reset(ctx, options.build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// GetProjectionState fetches a state of the projection.
func (client *Client) GetProjectionState(
	ctx context.Context,
	options StateRequest,
) (IsResult, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	result, protoErr := projectionsClient.State(ctx, options.build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return nil, err
	}

	return newStateResponse(result), nil
}

// GetProjectionResult fetches a result of a projection.
func (client *Client) GetProjectionResult(
	ctx context.Context,
	options ResultRequest,
) (IsResult, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.grpcProjectionsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	result, protoErr := projectionsClient.Result(ctx, options.build(), grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return nil, err
	}

	return newResultResponse(result), nil
}

// RestartProjectionsSubsystem restarts a projection subsystem at EventStoreDB.
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
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// ListAllProjections lists details of all projections.
func (client *Client) ListAllProjections(
	ctx context.Context,
) ([]statistics.Response, errors.Error) {
	statisticsClient, err := client.GetProjectionStatistics(ctx, StatisticsForAllProjections{})
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

// ListContinuousProjections lists details of all continuous projections.
func (client *Client) ListContinuousProjections(
	ctx context.Context,
) ([]statistics.Response, errors.Error) {
	statisticsClient, err := client.GetProjectionStatistics(ctx, StatisticsForContinuousProjections{})
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

// ListOneTimeProjections lists details of all one-time projections.
func (client *Client) ListOneTimeProjections(
	ctx context.Context,
) ([]statistics.Response, errors.Error) {
	statisticsClient, err := client.GetProjectionStatistics(ctx, StatisticsForOneTimeProjections{})
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

// NewClient creates a new client instance.
// Grpc connection must be passed to the new client instance.
func NewClient(
	grpcClient connection.GrpcClient,
) *Client {
	return &Client{
		grpcProjectionsClientFactory: grpc_proto_client_factory.FactoryImpl{},
		grpcClient:                   grpcClient,
		statisticsReaderFactory:      statistics_internal.ClientSyncFactoryImpl{},
	}
}
