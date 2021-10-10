package persistent

//go:generate mockgen -source=client.go -destination=client_mock.go -package=persistent

import (
	"context"

	"github.com/pivonroll/EventStore-Client-Go/errors"
	"github.com/pivonroll/EventStore-Client-Go/stream_revision"
)

// Client interface to interact with persistent subscriptions in EventStoreDB.
type Client interface {
	// SubscribeToStreamSync connects to an existing persistent subscription group.
	// Buffer size determines how many messages can be received util some of them must be
	// acknowledged (ack) or not-acknowledged (nack).
	// Group name and stream name determine to which persistent subscription we are connecting to.
	// Beside a group name a stream ID must be provided.
	// Two different persistent subscriptions to same stream can exist.
	// They just have to be on different groups.
	SubscribeToStreamSync(ctx context.Context,
		bufferSize int32,
		groupName string,
		streamId string,
	) (EventReader, errors.Error)
	// CreateSubscriptionGroupForStream creates a persistent subscription group on a stream.
	// Persistent subscription is identified by group name.
	//
	// You must have admin permissions to create a persistent subscription group.
	CreateSubscriptionGroupForStream(
		ctx context.Context,
		request SubscriptionGroupForStreamRequest,
	) errors.Error
	// CreateSubscriptionGroupForStreamAll creates a persistent subscription group to stream $all.
	// Persistent subscription to stream $all is identified by group name.
	//
	// You must have admin permissions to create a persistent subscription group.
	CreateSubscriptionGroupForStreamAll(
		ctx context.Context,
		request SubscriptionGroupForStreamAllRequest,
	) errors.Error
	// UpdateSubscriptionGroupForStream updates a persistent subscription group.
	// Once settings for a persistent subscription group is updated,
	// all existing connections will be dropped, and clients must reconnect.
	//
	// You must have admin permissions to update a persistent subscription group.
	UpdateSubscriptionGroupForStream(
		ctx context.Context,
		request SubscriptionGroupForStreamRequest,
	) errors.Error
	// UpdateSubscriptionGroupForStreamAll updates a persistent subscription group for stream $all.
	// Once settings for a persistent subscription group is updated,
	// all existing connections will be dropped, and clients must reconnect.
	//
	// You must have admin permissions to update a persistent subscription group.
	UpdateSubscriptionGroupForStreamAll(
		ctx context.Context,
		GroupName string,
		Position stream_revision.IsReadPositionAll,
		Settings CreateOrUpdateRequestSettings,
	) errors.Error
	// DeleteSubscriptionGroupForStream deletes a persistent subscription group for stream.
	// Once persistent subscription group is deleted, all existing connections will be dropped.
	//
	// You must have admin permissions to delete a persistent subscription group.
	DeleteSubscriptionGroupForStream(
		ctx context.Context,
		streamId string,
		groupName string,
	) errors.Error
	// DeleteSubscriptionGroupForStreamAll deletes a persistent subscription group for stream $all.
	// Once persistent subscription group is deleted, all existing connections will be dropped.
	//
	// You must have admin permissions to delete a persistent subscription group.
	DeleteSubscriptionGroupForStreamAll(
		ctx context.Context,
		groupName string,
	) errors.Error
}
