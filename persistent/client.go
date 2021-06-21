package persistent

import (
	"context"
)

type Client interface {
	SubscribeToStreamAsync(ctx context.Context,
		bufferSize int32,
		groupName string,
		streamName []byte,
	) (AsyncReadConnection, error)

	SubscribeToStreamSync(ctx context.Context,
		bufferSize int32,
		groupName string,
		streamName []byte,
		eventAppeared EventAppearedHandler,
		subscriptionDropped SubscriptionDroppedHandler,
	) (SyncReadConnection, error)
	SubscribeToAllAsync(ctx context.Context) (SyncReadConnection, error)
	CreateStreamSubscription(ctx context.Context, streamConfig SubscriptionStreamConfig) error
	CreateAllSubscription(ctx context.Context, allOptions SubscriptionAllOptionConfig) error
	UpdateStreamSubscription(ctx context.Context, streamConfig SubscriptionStreamConfig) error
	UpdateAllSubscription(ctx context.Context, allOptions SubscriptionUpdateAllOptionConfig) error
	DeleteStreamSubscription(ctx context.Context, deleteOptions DeleteOptions) error
	DeleteAllSubscription(ctx context.Context, groupName string) error
}
