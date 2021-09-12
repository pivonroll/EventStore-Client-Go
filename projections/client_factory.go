package projections

//go:generate mockgen -source=client_factory.go -destination=client_factory_mock.go -package=projections

import (
	"github.com/pivonroll/EventStore-Client-Go/connection"
	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type ClientFactory interface {
	CreateClient(grpcClient connection.GrpcClient, projectionsClient projections.ProjectionsClient) Client
}

type ClientFactoryImpl struct{}

func (clientFactory ClientFactoryImpl) CreateClient(
	grpcClient connection.GrpcClient,
	projectionsClient projections.ProjectionsClient) Client {
	return newClientImpl(grpcClient, projectionsClient)
}
