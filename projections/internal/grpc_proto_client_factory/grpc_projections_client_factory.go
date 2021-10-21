package grpc_proto_client_factory

//go:generate mockgen -source=grpc_projections_client_factory.go -destination=grpc_projections_client_factory_mock.go -package=grpc_proto_client_factory

import (
	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"google.golang.org/grpc"
)

type Factory interface {
	Create(connection *grpc.ClientConn) projections.ProjectionsClient
}

type FactoryImpl struct{}

func (factory FactoryImpl) Create(
	connection *grpc.ClientConn) projections.ProjectionsClient {
	return projections.NewProjectionsClient(connection)
}
