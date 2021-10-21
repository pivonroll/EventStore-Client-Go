package user_management

import (
	"github.com/pivonroll/EventStore-Client-Go/protos/users"
	"google.golang.org/grpc"
)

type grpcUserClientFactory interface {
	Create(cc grpc.ClientConnInterface) users.UsersClient
}

type grpcUserClientFactoryImpl struct{}

func (factory grpcUserClientFactoryImpl) Create(cc grpc.ClientConnInterface) users.UsersClient {
	return users.NewUsersClient(cc)
}
