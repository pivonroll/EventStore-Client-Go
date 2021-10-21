package client

import (
	"github.com/pivonroll/EventStore-Client-Go/connection"
	"github.com/pivonroll/EventStore-Client-Go/event_streams"
	"github.com/pivonroll/EventStore-Client-Go/operations"
	"github.com/pivonroll/EventStore-Client-Go/persistent"
	"github.com/pivonroll/EventStore-Client-Go/projections"
	"github.com/pivonroll/EventStore-Client-Go/user_management"
)

type Configuration = connection.Configuration

func ParseConnectionString(str string) (*connection.Configuration, error) {
	return connection.ParseConnectionString(str)
}

// Client ...
type Client struct {
	grpcClient              connection.GrpcClient
	Config                  *connection.Configuration
	projectionClientFactory projections.ClientFactory
	operationsClientFactory operations.ClientFactory
}

// NewClient ...
func NewClient(configuration *connection.Configuration) (*Client, error) {
	grpcClient := connection.NewGrpcClient(*configuration)
	return &Client{
		grpcClient:              grpcClient,
		Config:                  configuration,
		projectionClientFactory: projections.ClientFactoryImpl{},
		operationsClientFactory: operations.ClientFactoryImpl{},
	}, nil
}

// Close ...
func (client *Client) Close() {
	client.grpcClient.Close()
}

func (client *Client) Projections() projections.Client {
	return client.projectionClientFactory.Create(client.grpcClient)
}

func (client *Client) UserManagement() *user_management.Client {
	return user_management.NewClient(client.grpcClient)
}

func (client *Client) EventStreams() *event_streams.Client {
	return event_streams.NewClient(client.grpcClient)
}

func (client *Client) PersistentSubscriptions() *persistent.Client {
	return persistent.NewClient(client.grpcClient)
}

func (client *Client) Operations() operations.Client {
	return client.operationsClientFactory.Create(client.grpcClient)
}
