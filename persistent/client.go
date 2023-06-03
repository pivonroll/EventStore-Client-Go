// Empty header

// Package persistent provides a set of available interactions with persistent subscription groups in EventStoreDB.
//
// Unlink regular stream subscriptions, persistent subscription groups save the current position
// of the subscription to a stream on EvenStoreDB server.
// This save occurs regularly, depending on what are the settings for a subscription group.
//
// Note that if you set that save of the position should occur too often, for example after each message,
// it might occur that if messages are consumed very fast that the save will be postponed until
// there is some 'breathing space' for the server to save the position.
//
// Read more at https://developers.eventstore.com/server/v22.10/persistent-subscriptions.
package persistent

import (
	"context"

	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/pivonroll/EventStore-Client-Go/core/errors"
	"github.com/pivonroll/EventStore-Client-Go/persistent/event_reader"
	"github.com/pivonroll/EventStore-Client-Go/persistent/internal/event_reader_factory"
	"github.com/pivonroll/EventStore-Client-Go/persistent/internal/grpc_subscription_client"
	"github.com/pivonroll/EventStore-Client-Go/persistent/internal/message_adapter"
	persistentProto "github.com/pivonroll/EventStore-Client-Go/protos/v22.10/persistent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewClient(grpcClient connection.GrpcClient) *Client {
	return &Client{
		grpcClient:                    grpcClient,
		syncReadConnectionFactory:     event_reader_factory.FactoryImpl{},
		messageAdapterProvider:        message_adapter.MessageAdapterProviderImpl{},
		grpcSubscriptionClientFactory: grpc_subscription_client.FactoryImpl{},
	}
}

// Client provides interface to interact with EventStoreDB streams.
type Client struct {
	grpcClient                    connection.GrpcClient
	syncReadConnectionFactory     event_reader_factory.Factory
	messageAdapterProvider        message_adapter.MessageAdapterProvider
	grpcSubscriptionClientFactory grpc_subscription_client.Factory
}

const (
	// SubscribeToStreamSync_FailedToInitPersistentSubscriptionClientErr is returned if unknown
	// error was received when client tried to initialize protobuf client for persistent subscription.
	SubscribeToStreamSync_FailedToInitPersistentSubscriptionClientErr errors.ErrorCode = "SubscribeToStreamSync_FailedToInitPersistentSubscriptionClientErr"
	// SubscribeToStreamSync_FailedToSendStreamInitializationErr is returned if unknown error was
	// received when client failed to send settings for protobuf stream.
	SubscribeToStreamSync_FailedToSendStreamInitializationErr errors.ErrorCode = "SubscribeToStreamSync_FailedToSendStreamInitializationErr"
	// SubscribeToStreamSync_FailedToReceiveStreamInitializationErr is returned if unknown error was
	// received when client failed to receive a response for settings sent for protobuf stream.
	SubscribeToStreamSync_FailedToReceiveStreamInitializationErr errors.ErrorCode = "SubscribeToStreamSync_FailedToReceiveStreamInitializationErr"
	// SubscribeToStreamSync_NoSubscriptionConfirmationErr is returned if unknown error was
	// received instead a subscription confirmation.
	SubscribeToStreamSync_NoSubscriptionConfirmationErr errors.ErrorCode = "SubscribeToStreamSync_NoSubscriptionConfirmationErr"
)

// SubscribeToStreamSync connects to an existing persistent subscription group.
// Buffer size determines how many messages can be received util some of them must be
// acknowledged (ack) or not-acknowledged (nack).
// Group name and stream name determine to which persistent subscription we are connecting to.
// Beside a group name a stream ID must be provided.
// Two different persistent subscriptions to same stream can exist.
// They just have to be on different groups.
func (client Client) SubscribeToStreamSync(
	ctx context.Context,
	bufferSize int32,
	groupName string,
	streamName string,
) (event_reader.EventReader, errors.Error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	ctx, cancel := context.WithCancel(ctx)
	readClient, protoErr := persistentSubscriptionClient.Read(ctx,
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		defer cancel()
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			SubscribeToStreamSync_FailedToInitPersistentSubscriptionClientErr)
		return nil, err
	}

	protoErr = readClient.Send(toPersistentReadRequest(bufferSize, groupName, streamName))
	if protoErr != nil {
		defer cancel()
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			SubscribeToStreamSync_FailedToSendStreamInitializationErr)
		return nil, err
	}

	readResult, protoErr := readClient.Recv()
	if protoErr != nil {
		defer cancel()
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			SubscribeToStreamSync_FailedToReceiveStreamInitializationErr)
		return nil, err
	}
	switch readResult.Content.(type) {
	case *persistentProto.ReadResp_SubscriptionConfirmation_:
		{
			asyncConnection := client.syncReadConnectionFactory.Create(
				readClient,
				readResult.GetSubscriptionConfirmation().SubscriptionId,
				client.messageAdapterProvider.GetMessageAdapter(),
				cancel)

			return asyncConnection, nil
		}
	}

	defer cancel()
	return nil, errors.NewErrorCode(SubscribeToStreamSync_NoSubscriptionConfirmationErr)
}

// CreateStreamSubscription_FailedToCreateErr is returned if unknown error is received while client tried
// to update a persistent subscription for a stream.
const CreateStreamSubscription_FailedToCreateErr errors.ErrorCode = "CreateStreamSubscription_FailedToCreateErr"

// CreateSubscriptionGroupForStream creates a persistent subscription group on a stream.
// Persistent subscription is identified by group name.
//
// You must have admin permissions to create a persistent subscription group.
func (client Client) CreateSubscriptionGroupForStream(
	ctx context.Context,
	request SubscriptionGroupForStreamRequest,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Create(ctx, request.buildCreateRequest(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			CreateStreamSubscription_FailedToCreateErr)
		return err
	}

	return nil
}

// CreateAllSubscription_FailedToCreateErr is returned if unknown error is received while client tried
// to create a persistent subscription for stream $all.
const CreateAllSubscription_FailedToCreateErr errors.ErrorCode = "CreateAllSubscription_FailedToCreateErr"

// CreateSubscriptionGroupForStreamAll creates a persistent subscription group to stream $all.
// Persistent subscription to stream $all is identified by group name.
//
// You must have admin permissions to create a persistent subscription group.
func (client Client) CreateSubscriptionGroupForStreamAll(
	ctx context.Context,
	request SubscriptionGroupForStreamAllRequest,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Create(ctx, request.build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			CreateAllSubscription_FailedToCreateErr)
		return err
	}

	return nil
}

// UpdateStreamSubscription_FailedToUpdateErr is returned if unknown error is received while client tried
// to update a persistent subscription for a stream.
const UpdateStreamSubscription_FailedToUpdateErr errors.ErrorCode = "UpdateStreamSubscription_FailedToUpdateErr"

// UpdateSubscriptionGroupForStream updates a persistent subscription group.
// Once settings for a persistent subscription group is updated,
// all existing connections will be dropped, and clients must reconnect.
//
// You must have admin permissions to update a persistent subscription group.
func (client Client) UpdateSubscriptionGroupForStream(
	ctx context.Context,
	request SubscriptionGroupForStreamRequest,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Update(ctx, request.buildUpdateRequest(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			UpdateStreamSubscription_FailedToUpdateErr)
		return err
	}

	return nil
}

// UpdateAllSubscription_FailedToUpdateErr is returned if unknown error is received while client tried
// to update a persistent subscription for stream $all.
const UpdateAllSubscription_FailedToUpdateErr errors.ErrorCode = "UpdateAllSubscription_FailedToUpdateErr"

// UpdateSubscriptionGroupForStreamAll updates a persistent subscription group for stream $all.
// Once settings for a persistent subscription group is updated,
// all existing connections will be dropped, and clients must reconnect.
//
// You must have admin permissions to update a persistent subscription group.
func (client Client) UpdateSubscriptionGroupForStreamAll(
	ctx context.Context,
	request UpdateSubscriptionGroupForStreamAllRequest,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Update(ctx, request.build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			UpdateAllSubscription_FailedToUpdateErr)
		return err
	}

	return nil
}

// DeleteStreamSubscription_FailedToDeleteErr is returned if unknown error is received while client tried
// to delete a persistent subscription for a stream.
const DeleteStreamSubscription_FailedToDeleteErr errors.ErrorCode = "DeleteStreamSubscription_FailedToDeleteErr"

// DeleteSubscriptionGroupForStream deletes a persistent subscription group for stream.
// Once persistent subscription group is deleted, all existing connections will be dropped.
//
// You must have admin permissions to delete a persistent subscription group.
func (client Client) DeleteSubscriptionGroupForStream(
	ctx context.Context,
	streamId string,
	groupName string,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	request := deleteSubscriptionGroupForStreamRequest{
		StreamId:  streamId,
		GroupName: groupName,
	}
	_, protoErr := persistentSubscriptionClient.Delete(ctx, request.build(),
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			DeleteStreamSubscription_FailedToDeleteErr)
		return err
	}

	return nil
}

// DeleteAllSubscription_FailedToDeleteErr is returned if unknown error is received while client tried
// to delete a persistent subscription for stream $all.
const DeleteAllSubscription_FailedToDeleteErr errors.ErrorCode = "DeleteAllSubscription_FailedToDeleteErr"

// DeleteSubscriptionGroupForStreamAll deletes a persistent subscription group for stream $all.
// Once persistent subscription group is deleted, all existing connections will be dropped.
//
// You must have admin permissions to delete a persistent subscription group.
func (client Client) DeleteSubscriptionGroupForStreamAll(
	ctx context.Context,
	groupName string,
) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	protoRequest := deleteRequestAllOptionsProto(groupName)
	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Delete(ctx, protoRequest,
		grpc.Header(&headers), grpc.Trailer(&trailers))
	if protoErr != nil {
		err := client.grpcClient.HandleError(handle, headers, trailers, protoErr,
			DeleteAllSubscription_FailedToDeleteErr)
		return err
	}

	return nil
}
