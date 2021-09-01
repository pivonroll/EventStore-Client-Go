package persistent

//go:generate mockgen -source=sync_read_connection_factory.go -destination=sync_read_connection_factory_mock.go -package=persistent

type SyncReadConnectionFactory interface {
	NewSyncReadConnection(client protoClient,
		subscriptionId string,
		messageAdapter messageAdapter,
	) SyncReadConnection
}

type SyncReadConnectionFactoryImpl struct{}

func (factory SyncReadConnectionFactoryImpl) NewSyncReadConnection(
	client protoClient,
	subscriptionId string,
	messageAdapter messageAdapter) SyncReadConnection {
	return &syncReadConnectionImpl{
		client:         client,
		subscriptionId: subscriptionId,
		messageAdapter: messageAdapter,
	}
}
