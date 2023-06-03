package grpc_operations_client

//go:generate mockgen -source=grpc_operations_client_factory.go -destination=grpc_operations_client_factory_mock.go -package=grpc_operations_client

import (
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/operations"
	"google.golang.org/grpc"
)

type Factory interface {
	Create(*grpc.ClientConn) operations.OperationsClient
}

type FactoryImpl struct{}

func (factory FactoryImpl) Create(
	conn *grpc.ClientConn,
) operations.OperationsClient {
	return operations.NewOperationsClient(conn)
}
