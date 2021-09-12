package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/EventStore/EventStore-Client-Go/client/filtering"
	"github.com/EventStore/EventStore-Client-Go/connection"
	"github.com/EventStore/EventStore-Client-Go/direction"
	errors2 "github.com/EventStore/EventStore-Client-Go/errors"
	"github.com/EventStore/EventStore-Client-Go/event_streams"
	"github.com/EventStore/EventStore-Client-Go/internal/protoutils"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/persistent"
	"github.com/EventStore/EventStore-Client-Go/projections"
	persistentProto "github.com/EventStore/EventStore-Client-Go/protos/persistent"
	projectionsProto "github.com/EventStore/EventStore-Client-Go/protos/projections"
	api "github.com/EventStore/EventStore-Client-Go/protos/streams"
	"github.com/EventStore/EventStore-Client-Go/protos/streams2"
	"github.com/EventStore/EventStore-Client-Go/stream_position"
	stream_revision "github.com/EventStore/EventStore-Client-Go/streamrevision"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Configuration = connection.Configuration

func ParseConnectionString(str string) (*connection.Configuration, error) {
	return connection.ParseConnectionString(str)
}

// Client ...
type Client struct {
	grpcClient                connection.GrpcClient
	Config                    *connection.Configuration
	persistentClientFactory   persistent.ClientFactory
	projectionClientFactory   projections.ClientFactory
	eventStreamsClientFactory event_streams.ClientFactory
}

// NewClient ...
func NewClient(configuration *connection.Configuration) (*Client, error) {
	grpcClient := connection.NewGrpcClient(*configuration)
	return &Client{
		grpcClient:                grpcClient,
		Config:                    configuration,
		persistentClientFactory:   persistent.ClientFactoryImpl{},
		projectionClientFactory:   projections.ClientFactoryImpl{},
		eventStreamsClientFactory: event_streams.ClientFactoryImpl{},
	}, nil
}

// Close ...
func (client *Client) Close() error {
	client.grpcClient.Close()
	return nil
}

const (
	AppendToStream_WrongExpectedVersionErr         = "AppendToStream_WrongExpectedVersionErr"
	AppendToStream_WrongExpectedVersion_20_6_0_Err = "AppendToStream_WrongExpectedVersion_20_6_0_Err"
)

func (client *Client) AppendToStream(
	ctx context.Context,
	options event_streams.AppendRequestContentOptions,
	events []event_streams.ProposedEvent,
) (WriteResult, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return WriteResult{}, err
	}
	eventStreamsClient := client.eventStreamsClientFactory.CreateClient(
		client.grpcClient, streams2.NewStreamsClient(handle.Connection()))

	appender, err := eventStreamsClient.GetAppender(ctx, handle)
	if err != nil {
		return WriteResult{}, err
	}

	err = appender.Send(handle, event_streams.AppendRequest{
		Content: options,
	})

	if err != nil {
		log.Println("Could not send append request header", err)
		return WriteResult{}, err
	}

	for _, event := range events {
		err = appender.Send(handle, event_streams.AppendRequest{
			Content: event.ToProposedMessage(),
		})

		if err != nil {
			log.Println("Could not send append request", err)
			return WriteResult{}, err
		}
	}

	response, err := appender.CloseAndRecv(handle)
	if err != nil {
		log.Println("Could not close sender end", err)
		return WriteResult{}, err
	}

	switch response.Result.(type) {
	case event_streams.AppendResponseSuccess:
		successResponse := response.Result.(event_streams.AppendResponseSuccess)

		var streamRevision uint64 = 1
		revision, isCurrentRevision := successResponse.CurrentRevision.(event_streams.AppendResponseSuccessCurrentRevision)
		if isCurrentRevision {
			streamRevision = revision.CurrentRevision
		}

		var commitPosition uint64
		var preparePosition uint64
		if position, ok := successResponse.Position.(event_streams.AppendResponseSuccessPosition); ok {
			commitPosition = position.CommitPosition
			preparePosition = position.PreparePosition
		} else if !isCurrentRevision {
			streamRevision = 0
		}

		return WriteResult{
			CommitPosition:      commitPosition,
			PreparePosition:     preparePosition,
			NextExpectedVersion: streamRevision,
		}, nil
	case event_streams.AppendResponseWrongExpectedVersion:
		wrongVersion := response.Result.(event_streams.AppendResponseWrongExpectedVersion)
		type currentRevisionType struct {
			Revision uint64
			NoStream bool
		}
		var currentRevision *currentRevisionType

		if wrongVersion.CurrentRevision_20_6_0 != nil {
			switch wrongVersion.CurrentRevision_20_6_0.(type) {
			case event_streams.AppendResponseWrongCurrentRevision_20_6_0:
				revision := wrongVersion.CurrentRevision_20_6_0.(event_streams.AppendResponseWrongCurrentRevision_20_6_0)
				currentRevision = &currentRevisionType{
					Revision: revision.CurrentRevision,
					NoStream: false,
				}
			case event_streams.AppendResponseWrongCurrentRevisionNoStream_20_6_0:
				currentRevision = &currentRevisionType{
					NoStream: true,
				}
			}
		} else if wrongVersion.CurrentRevision != nil {
			switch wrongVersion.CurrentRevision.(type) {
			case event_streams.AppendResponseWrongCurrentRevision:
				revision := wrongVersion.CurrentRevision.(event_streams.AppendResponseWrongCurrentRevision)

				currentRevision = &currentRevisionType{
					Revision: revision.CurrentRevision,
					NoStream: false,
				}
			case event_streams.AppendResponseWrongCurrentRevisionNoStream:
				currentRevision = &currentRevisionType{
					NoStream: true,
				}
			}
		}

		if currentRevision != nil {
			if currentRevision.NoStream {
				log.Println("Wrong expected revision. Current revision no stream")
			} else {
				log.Println("Wrong expected revision. Current revision:", currentRevision.Revision)
			}
		}

		type expectedRevisionType struct {
			Revision     uint64
			IsAny        bool
			StreamExists bool
			NoStream     bool
		}

		var expectedRevision *expectedRevisionType

		if wrongVersion.ExpectedRevision_20_6_0 != nil {
			switch wrongVersion.ExpectedRevision_20_6_0.(type) {
			case event_streams.AppendResponseWrongExpectedRevision_20_6_0:
				revision := wrongVersion.ExpectedRevision_20_6_0.(event_streams.AppendResponseWrongExpectedRevision_20_6_0)
				expectedRevision = &expectedRevisionType{
					Revision: revision.ExpectedRevision,
				}
			case event_streams.AppendResponseWrongExpectedRevisionAny_20_6_0:
				expectedRevision = &expectedRevisionType{
					IsAny: true,
				}
			case event_streams.AppendResponseWrongExpectedRevisionStreamExists_20_6_0:
				expectedRevision = &expectedRevisionType{
					StreamExists: true,
				}
			}
		} else if wrongVersion.ExpectedRevision != nil {
			switch wrongVersion.ExpectedRevision.(type) {
			case event_streams.AppendResponseWrongExpectedRevision:
				revision := wrongVersion.ExpectedRevision.(event_streams.AppendResponseWrongExpectedRevision)
				expectedRevision = &expectedRevisionType{
					Revision: revision.ExpectedRevision,
				}
			case event_streams.AppendResponseWrongExpectedRevisionAny:
				expectedRevision = &expectedRevisionType{
					IsAny: true,
				}
			case event_streams.AppendResponseWrongExpectedRevisionStreamExists:
				expectedRevision = &expectedRevisionType{
					StreamExists: true,
				}
			case event_streams.AppendResponseWrongExpectedRevisionNoStream:
				expectedRevision = &expectedRevisionType{
					NoStream: true,
				}
			}
		}

		if expectedRevision != nil {
			if expectedRevision.StreamExists {
				log.Println("Wrong expected revision. Stream Exists!")
			} else if expectedRevision.IsAny {
				log.Println("Wrong expected revision. Any!")
			} else if expectedRevision.NoStream {
				log.Println("Wrong expected revision. No Stream!")
			} else {
				log.Println("Wrong expected revision. Expected revision: ", expectedRevision.Revision)
			}
		}

		if wrongVersion.CurrentRevision_20_6_0 != nil || wrongVersion.ExpectedRevision_20_6_0 != nil {
			return WriteResult{}, errors.New(AppendToStream_WrongExpectedVersion_20_6_0_Err)
		}
		return WriteResult{}, errors.New(AppendToStream_WrongExpectedVersionErr)
	}

	return WriteResult{
		CommitPosition:      0,
		PreparePosition:     0,
		NextExpectedVersion: 1,
	}, nil
}

// AppendToStream_OLD ...
func (client *Client) AppendToStream_OLD(
	context context.Context,
	streamID string,
	streamRevision stream_revision.StreamRevision,
	events []messages.ProposedEvent,
) (*WriteResult, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	streamsClient := api.NewStreamsClient(handle.Connection())
	var headers, trailers metadata.MD

	appendOperation, err := streamsClient.Append(context, grpc.Header(&headers), grpc.Trailer(&trailers))
	if err != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Could not construct append operation. Reason: %v", err)
	}

	header := protoutils.ToAppendHeader(streamID, streamRevision)

	if err := appendOperation.Send(header); err != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Could not send append request header. Reason: %v", err)
	}

	for _, event := range events {
		appendRequest := &api.AppendReq{
			Content: &api.AppendReq_ProposedMessage_{
				ProposedMessage: protoutils.ToProposedMessage(event),
			},
		}

		if err = appendOperation.Send(appendRequest); err != nil {
			err = client.grpcClient.HandleError(handle, headers, trailers, err)
			return nil, fmt.Errorf("Could not send append request. Reason: %v", err)
		}
	}

	response, err := appendOperation.CloseAndRecv()
	if err != nil {
		return nil, client.grpcClient.HandleError(handle, headers, trailers, err)
	}

	result := response.GetResult()
	switch result.(type) {
	case *api.AppendResp_Success_:
		{
			success := result.(*api.AppendResp_Success_)
			var streamRevision uint64
			if _, ok := success.Success.GetCurrentRevisionOption().(*api.AppendResp_Success_NoStream); ok {
				streamRevision = 1
			} else {
				streamRevision = success.Success.GetCurrentRevision()
			}

			var commitPosition uint64
			var preparePosition uint64
			if position, ok := success.Success.GetPositionOption().(*api.AppendResp_Success_Position); ok {
				commitPosition = position.Position.CommitPosition
				preparePosition = position.Position.PreparePosition
			} else {
				streamRevision = success.Success.GetCurrentRevision()
			}

			return &WriteResult{
				CommitPosition:      commitPosition,
				PreparePosition:     preparePosition,
				NextExpectedVersion: streamRevision,
			}, nil
		}
	case *api.AppendResp_WrongExpectedVersion_:
		{
			return nil, errors2.ErrWrongExpectedStreamRevision
		}
	}

	return &WriteResult{
		CommitPosition:      0,
		PreparePosition:     0,
		NextExpectedVersion: 1,
	}, nil
}

// DeleteStream_OLD ...
func (client *Client) DeleteStream_OLD(
	context context.Context,
	streamID string,
	streamRevision stream_revision.StreamRevision,
) (*DeleteResult, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	streamsClient := api.NewStreamsClient(handle.Connection())
	var headers, trailers metadata.MD
	deleteRequest := protoutils.ToDeleteRequest(streamID, streamRevision)
	deleteResponse, err := streamsClient.Delete(context, deleteRequest, grpc.Header(&headers), grpc.Trailer(&trailers))
	if err != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Failed to perform delete, details: %v", err)
	}

	return &DeleteResult{Position: protoutils.DeletePositionFromProto(deleteResponse)}, nil
}

// TombstoneStream_OLD Tombstone ...
func (client *Client) TombstoneStream_OLD(
	context context.Context,
	streamID string,
	streamRevision stream_revision.StreamRevision,
) (*DeleteResult, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	streamsClient := api.NewStreamsClient(handle.Connection())
	var headers, trailers metadata.MD
	tombstoneRequest := protoutils.ToTombstoneRequest(streamID, streamRevision)
	tombstoneResponse, err := streamsClient.Tombstone(context, tombstoneRequest, grpc.Header(&headers), grpc.Trailer(&trailers))
	if err != nil {
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Failed to perform delete, details: %v", err)
	}

	return &DeleteResult{Position: protoutils.TombstonePositionFromProto(tombstoneResponse)}, nil
}

// ReadStreamEvents_OLD ...
func (client *Client) ReadStreamEvents_OLD(
	context context.Context,
	direction direction.Direction,
	streamID string,
	from stream_position.StreamPosition,
	count uint64,
	resolveLinks bool) (*ReadStream, error) {
	readRequest := protoutils.ToReadStreamRequest(streamID, direction, from, count, resolveLinks)
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	streamsClient := api.NewStreamsClient(handle.Connection())

	return readInternal(context, client.grpcClient, handle, streamsClient, readRequest)
}

// ReadAllEvents_OLD ...
func (client *Client) ReadAllEvents_OLD(
	context context.Context,
	direction direction.Direction,
	from stream_position.AllStreamPosition,
	count uint64,
	resolveLinks bool,
) (*ReadStream, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	streamsClient := api.NewStreamsClient(handle.Connection())
	readRequest := protoutils.ToReadAllRequest(direction, from, count, resolveLinks)
	return readInternal(context, client.grpcClient, handle, streamsClient, readRequest)
}

// SubscribeToStream_OLD ...
func (client *Client) SubscribeToStream_OLD(
	ctx context.Context,
	streamID string,
	from stream_position.StreamPosition,
	resolveLinks bool,
) (*Subscription, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	var headers, trailers metadata.MD
	streamsClient := api.NewStreamsClient(handle.Connection())
	subscriptionRequest, err := protoutils.ToStreamSubscriptionRequest(streamID, from, resolveLinks, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to construct subscription. Reason: %v", err)
	}
	ctx, cancel := context.WithCancel(ctx)
	readClient, err := streamsClient.Read(ctx, subscriptionRequest, grpc.Header(&headers), grpc.Trailer(&trailers))
	if err != nil {
		defer cancel()
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Failed to construct subscription. Reason: %v", err)
	}
	readResult, err := readClient.Recv()
	if err != nil {
		defer cancel()
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Failed to perform read. Reason: %v", err)
	}
	switch readResult.Content.(type) {
	case *api.ReadResp_Confirmation:
		{
			confirmation := readResult.GetConfirmation()
			return NewSubscription(client, cancel, readClient, confirmation.SubscriptionId), nil
		}
	case *api.ReadResp_StreamNotFound_:
		{
			defer cancel()
			return nil, fmt.Errorf("Failed to initiate subscription because the stream (%s) was not found.", streamID)
		}
	}
	defer cancel()
	return nil, fmt.Errorf("Failed to initiate subscription.")
}

// SubscribeToAll_OLD ...
func (client *Client) SubscribeToAll_OLD(
	ctx context.Context,
	from stream_position.AllStreamPosition,
	resolveLinks bool,
) (*Subscription, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	streamsClient := api.NewStreamsClient(handle.Connection())
	var headers, trailers metadata.MD
	subscriptionRequest, err := protoutils.ToAllSubscriptionRequest(from, resolveLinks, nil)
	ctx, cancel := context.WithCancel(ctx)
	readClient, err := streamsClient.Read(ctx, subscriptionRequest, grpc.Header(&headers), grpc.Trailer(&trailers))
	if err != nil {
		defer cancel()
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Failed to construct subscription. Reason: %v", err)
	}
	readResult, err := readClient.Recv()
	if err != nil {
		defer cancel()
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Failed to perform read. Reason: %v", err)
	}
	switch readResult.Content.(type) {
	case *api.ReadResp_Confirmation:
		{
			confirmation := readResult.GetConfirmation()
			return NewSubscription(client, cancel, readClient, confirmation.SubscriptionId), nil
		}
	}
	defer cancel()
	return nil, fmt.Errorf("Failed to initiate subscription.")
}

// SubscribeToAllFiltered_OLD ...
func (client *Client) SubscribeToAllFiltered_OLD(
	ctx context.Context,
	from stream_position.AllStreamPosition,
	resolveLinks bool,
	filterOptions filtering.SubscriptionFilterOptions,
) (*Subscription, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	streamsClient := api.NewStreamsClient(handle.Connection())
	subscriptionRequest, err := protoutils.ToAllSubscriptionRequest(from, resolveLinks, &filterOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to construct subscription. Reason: %v", err)
	}
	var headers, trailers metadata.MD
	ctx, cancel := context.WithCancel(ctx)
	readClient, err := streamsClient.Read(ctx, subscriptionRequest, grpc.Header(&headers), grpc.Trailer(&trailers))
	if err != nil {
		defer cancel()
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Failed to initiate subscription. Reason: %v", err)
	}
	readResult, err := readClient.Recv()
	if err != nil {
		defer cancel()
		err = client.grpcClient.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("Failed to read from subscription. Reason: %v", err)
	}
	switch readResult.Content.(type) {
	case *api.ReadResp_Confirmation:
		{
			confirmation := readResult.GetConfirmation()
			return NewSubscription(client, cancel, readClient, confirmation.SubscriptionId), nil
		}
	}
	defer cancel()
	return nil, fmt.Errorf("Failed to initiate subscription.")
}

// ConnectToPersistentSubscription ...
func (client *Client) ConnectToPersistentSubscription(
	ctx context.Context,
	bufferSize int32,
	groupName string,
	streamName []byte,
) (persistent.SyncReadConnection, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}
	persistentSubscriptionClient := client.persistentClientFactory.
		CreateClient(client.grpcClient, persistentProto.NewPersistentSubscriptionsClient(handle.Connection()))

	return persistentSubscriptionClient.SubscribeToStreamSync(
		ctx,
		handle,
		bufferSize,
		groupName,
		streamName,
	)
}

func (client *Client) CreatePersistentSubscription(
	ctx context.Context,
	streamConfig persistent.SubscriptionStreamConfig,
) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}
	persistentSubscriptionClient := client.persistentClientFactory.
		CreateClient(client.grpcClient, persistentProto.NewPersistentSubscriptionsClient(handle.Connection()))

	return persistentSubscriptionClient.CreateStreamSubscription(ctx, handle, streamConfig)
}

func (client *Client) CreatePersistentSubscriptionAll(
	ctx context.Context,
	allOptions persistent.SubscriptionAllOptionConfig,
) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}
	persistentSubscriptionClient := client.persistentClientFactory.
		CreateClient(client.grpcClient, persistentProto.NewPersistentSubscriptionsClient(handle.Connection()))

	return persistentSubscriptionClient.CreateAllSubscription(ctx, handle, allOptions)
}

func (client *Client) UpdatePersistentStreamSubscription(
	ctx context.Context,
	streamConfig persistent.SubscriptionStreamConfig,
) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}
	persistentSubscriptionClient := client.persistentClientFactory.
		CreateClient(client.grpcClient, persistentProto.NewPersistentSubscriptionsClient(handle.Connection()))

	return persistentSubscriptionClient.UpdateStreamSubscription(ctx, handle, streamConfig)
}

func (client *Client) UpdatePersistentSubscriptionAll(
	ctx context.Context,
	allOptions persistent.SubscriptionUpdateAllOptionConfig,
) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}
	persistentSubscriptionClient := client.persistentClientFactory.
		CreateClient(client.grpcClient, persistentProto.NewPersistentSubscriptionsClient(handle.Connection()))

	return persistentSubscriptionClient.UpdateAllSubscription(ctx, handle, allOptions)
}

func (client *Client) DeletePersistentSubscription(
	ctx context.Context,
	deleteOptions persistent.DeleteOptions,
) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}
	persistentSubscriptionClient := client.persistentClientFactory.
		CreateClient(client.grpcClient, persistentProto.NewPersistentSubscriptionsClient(handle.Connection()))

	return persistentSubscriptionClient.DeleteStreamSubscription(ctx, handle, deleteOptions)
}

func (client *Client) DeletePersistentSubscriptionAll(
	ctx context.Context,
	groupName string,
) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}
	persistentSubscriptionClient := client.persistentClientFactory.
		CreateClient(client.grpcClient, persistentProto.NewPersistentSubscriptionsClient(handle.Connection()))

	return persistentSubscriptionClient.DeleteAllSubscription(ctx, handle, groupName)
}

func (client *Client) CreateProjection(ctx context.Context, options projections.CreateOptionsRequest) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.CreateProjection(ctx, handle, options)
}

func (client *Client) UpdateProjection(ctx context.Context, options projections.UpdateOptionsRequest) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.UpdateProjection(ctx, handle, options)
}

func (client *Client) AbortProjection(ctx context.Context, options projections.AbortOptionsRequest) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.AbortProjection(ctx, handle, options)
}

func (client *Client) DisableProjection(ctx context.Context, options projections.DisableOptionsRequest) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.DisableProjection(ctx, handle, options)
}

func (client *Client) ResetProjection(ctx context.Context, options projections.ResetOptionsRequest) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.ResetProjection(ctx, handle, options)
}

func (client *Client) DeleteProjection(ctx context.Context, options projections.DeleteOptionsRequest) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.DeleteProjection(ctx, handle, options)
}

func (client *Client) EnableProjection(ctx context.Context, options projections.EnableOptionsRequest) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.EnableProjection(ctx, handle, options)
}

func (client *Client) RestartProjectionsSubsystem(ctx context.Context) error {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.RestartProjectionsSubsystem(ctx, handle)
}

func (client *Client) GetProjectionState(
	ctx context.Context,
	options projections.StateOptionsRequest) (projections.StateResponse, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.GetProjectionState(ctx, handle, options)
}

func (client *Client) GetProjectionResult(
	ctx context.Context,
	options projections.ResultOptionsRequest) (projections.ResultResponse, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.GetProjectionResult(ctx, handle, options)
}

func (client *Client) GetProjectionStatistics(
	ctx context.Context,
	options projections.StatisticsOptionsRequest) (projections.StatisticsClientSync, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.GetProjectionStatistics(ctx, handle, options)
}

func (client *Client) ListAllProjections(
	ctx context.Context) ([]projections.StatisticsClientResponse, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.ListAllProjections(ctx, handle)
}

func (client *Client) ListContinuousProjections(
	ctx context.Context) ([]projections.StatisticsClientResponse, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.ListContinuousProjections(ctx, handle)
}

func (client *Client) ListOneTimeProjections(
	ctx context.Context) ([]projections.StatisticsClientResponse, error) {
	handle, err := client.grpcClient.GetConnectionHandle()
	if err != nil {
		return nil, err
	}

	projectionsClient := client.projectionClientFactory.CreateClient(client.grpcClient,
		projectionsProto.NewProjectionsClient(handle.Connection()))

	return projectionsClient.ListOneTimeProjections(ctx, handle)
}

func readInternal(
	ctx context.Context,
	client connection.GrpcClient,
	handle connection.ConnectionHandle,
	streamsClient api.StreamsClient,
	readRequest *api.ReadReq,
) (*ReadStream, error) {
	var headers, trailers metadata.MD
	ctx, cancel := context.WithCancel(ctx)
	result, err := streamsClient.Read(ctx, readRequest, grpc.Header(&headers), grpc.Trailer(&trailers))
	if err != nil {
		defer cancel()
		err = client.HandleError(handle, headers, trailers, err)
		return nil, fmt.Errorf("failed to construct read stream. Reason: %v", err)
	}

	params := ReadStreamParams{
		client:   client,
		handle:   handle,
		cancel:   cancel,
		inner:    result,
		headers:  headers,
		trailers: trailers,
	}

	stream := NewReadStream(params)

	return stream, nil
}
