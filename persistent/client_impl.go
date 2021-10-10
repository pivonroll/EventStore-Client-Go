package persistent

import (
	"context"

	"github.com/pivonroll/EventStore-Client-Go/connection"
	"github.com/pivonroll/EventStore-Client-Go/errors"
	persistentProto "github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type clientImpl struct {
	grpcClient                    connection.GrpcClient
	syncReadConnectionFactory     eventReaderFactory
	messageAdapterProvider        messageAdapterProvider
	grpcSubscriptionClientFactory grpcSubscriptionClientFactory
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
func (client clientImpl) SubscribeToStreamSync(
	ctx context.Context,
	bufferSize int32,
	groupName string,
	streamName string,
) (EventReader, errors.Error) {
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
func (client clientImpl) CreateSubscriptionGroupForStream(
	ctx context.Context,
	request CreateOrUpdateStreamRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Create(ctx, request.BuildCreateStreamRequest(),
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
func (client clientImpl) CreateSubscriptionGroupForStreamAll(
	ctx context.Context,
	request CreateAllRequest) errors.Error {
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
func (client clientImpl) UpdateSubscriptionGroupForStream(
	ctx context.Context,
	request CreateOrUpdateStreamRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Update(ctx, request.BuildUpdateStreamRequest(),
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
func (client clientImpl) UpdateSubscriptionGroupForStreamAll(
	ctx context.Context,
	request UpdateAllRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Update(ctx, request.Build(),
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
func (client clientImpl) DeleteSubscriptionGroupForStream(
	ctx context.Context,
	request DeleteRequest) errors.Error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	persistentSubscriptionClient := client.grpcSubscriptionClientFactory.Create(handle.Connection())

	var headers, trailers metadata.MD
	_, protoErr := persistentSubscriptionClient.Delete(ctx, request.Build(),
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
func (client clientImpl) DeleteSubscriptionGroupForStreamAll(
	ctx context.Context,
	groupName string) errors.Error {
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

func newClientImpl(grpcClient connection.GrpcClient) clientImpl {
	return clientImpl{
		grpcClient:                    grpcClient,
		syncReadConnectionFactory:     eventReaderFactoryImpl{},
		messageAdapterProvider:        messageAdapterProviderImpl{},
		grpcSubscriptionClientFactory: grpcSubscriptionClientFactoryImpl{},
	}
}
