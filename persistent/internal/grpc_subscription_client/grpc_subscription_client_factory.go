package grpc_subscription_client

//go:generate mockgen -source=grpc_subscription_client_factory.go -destination=../mocks/grpc_subscription_client_factory_mock.go -mock_names=Factory=GrpcClientFactory -package=mocks

import (
	persistentProto "github.com/pivonroll/EventStore-Client-Go/protos/v21.6/persistent"
	"google.golang.org/grpc"
)

type Factory interface {
	Create(*grpc.ClientConn) persistentProto.PersistentSubscriptionsClient
}

type FactoryImpl struct{}

func (factory FactoryImpl) Create(
	conn *grpc.ClientConn) persistentProto.PersistentSubscriptionsClient {
	return persistentProto.NewPersistentSubscriptionsClient(conn)
}
