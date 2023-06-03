// Header

// Package operations provides interaction with EventStoreDb admin operations.
// Before accessing streams a grpc connection needs to be established with EventStore through
// github.com/pivonroll/EventStore-Client-Go/core/connection package.
package operations

import (
	"context"

	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/pivonroll/EventStore-Client-Go/core/errors"
	"github.com/pivonroll/EventStore-Client-Go/operations/internal/grpc_operations_client"
	protoOperations "github.com/pivonroll/EventStore-Client-Go/protos/v22.10/operations"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client provides interaction with admin operations at EvenStoreDB.
type Client struct {
	grpcClient                  connection.GrpcClient
	scavengeResponseAdapter     scavengeResponseAdapter
	grpcOperationsClientFactory grpc_operations_client.Factory
}

// NewClient creates a new client which interacts with admin operations at EventStoreDB.
func NewClient(grpcClient connection.GrpcClient) *Client {
	return &Client{
		grpcClient:                  grpcClient,
		scavengeResponseAdapter:     scavengeResponseAdapterImpl{},
		grpcOperationsClientFactory: grpc_operations_client.FactoryImpl{},
	}
}

// Shutdown shuts down the EventStoreDB node.
func (client *Client) Shutdown(ctx context.Context) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcOperationsClient := client.grpcOperationsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcOperationsClient.Shutdown(ctx, &shared.Empty{},
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// MergeIndexes initiates an index merge operation.
func (client *Client) MergeIndexes(ctx context.Context) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcOperationsClient := client.grpcOperationsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcOperationsClient.MergeIndexes(ctx, &shared.Empty{},
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// ResignNode resigns a node.
func (client *Client) ResignNode(ctx context.Context) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcOperationsClient := client.grpcOperationsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcOperationsClient.ResignNode(ctx, &shared.Empty{},
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// SetNodePriority sets the node priority.
func (client *Client) SetNodePriority(ctx context.Context, priority int32) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcOperationsClient := client.grpcOperationsClientFactory.Create(handle.Connection())

	request := &protoOperations.SetNodePriorityReq{
		Priority: priority,
	}

	var headers, trailers metadata.MD
	_, protoErr := grpcOperationsClient.SetNodePriority(ctx, request,
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr, errors.FatalError)
		return err
	}

	return nil
}

// RestartPersistentSubscriptions restart persistent subscriptions
func (client *Client) RestartPersistentSubscriptions(ctx context.Context) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	grpcOperationsClient := client.grpcOperationsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := grpcOperationsClient.RestartPersistentSubscriptions(ctx, &shared.Empty{},
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return err
	}

	return nil
}

const (
	StartScavenge_ThreadCountLessOrEqualZeroErr errors.ErrorCode = "StartScavenge_ThreadCountLessOrEqualZeroErr"
	StartScavenge_StartFromChunkLessThanZeroErr errors.ErrorCode = "StartScavenge_StartFromChunkLessThanZeroErr"
)

// StartScavenge starts a scavenge operation.
func (client *Client) StartScavenge(ctx context.Context,
	request StartScavengeRequest,
) (ScavengeResponse, errors.Error) {
	if request.ThreadCount <= 0 {
		return ScavengeResponse{}, errors.NewErrorCode(StartScavenge_ThreadCountLessOrEqualZeroErr)
	}

	if request.StartFromChunk < 0 {
		return ScavengeResponse{}, errors.NewErrorCode(StartScavenge_StartFromChunkLessThanZeroErr)
	}

	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return ScavengeResponse{}, err
	}

	grpcOperationsClient := client.grpcOperationsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	protoResponse, protoErr := grpcOperationsClient.StartScavenge(ctx, request.build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return ScavengeResponse{}, err
	}

	return client.scavengeResponseAdapter.create(protoResponse), nil
}

// StopScavenge stops a scavenge operation.
func (client *Client) StopScavenge(ctx context.Context,
	scavengeId string,
) (ScavengeResponse, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return ScavengeResponse{}, err
	}

	request := &protoOperations.StopScavengeReq{
		Options: &protoOperations.StopScavengeReq_Options{
			ScavengeId: scavengeId,
		},
	}

	grpcOperationsClient := client.grpcOperationsClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	protoResponse, protoErr := grpcOperationsClient.StopScavenge(ctx, request,
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			errors.FatalError)
		return ScavengeResponse{}, err
	}

	return client.scavengeResponseAdapter.create(protoResponse), nil
}
