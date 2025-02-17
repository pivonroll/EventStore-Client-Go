package connection_integration_test

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/pivonroll/EventStore-Client-Go/core/test_utils"
	"github.com/pivonroll/EventStore-Client-Go/event_streams"
)

func initializeContainerAndClient(t *testing.T,
	envVariableOverrides map[string]string) (*test_utils.Container, connection.GrpcClient, *event_streams.Client) {
	container, grpcClient := test_utils.InitializeContainerAndGrpcClient(t, envVariableOverrides)

	eventStreamsClient := event_streams.NewClient(grpcClient)
	return container, grpcClient, eventStreamsClient
}

func initializeContainerAndClientWithTLS(t *testing.T,
	envVariableOverrides map[string]string) (connection.GrpcClient, *event_streams.Client, test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClientWithTLS(t, envVariableOverrides)

	eventStreamsClient := event_streams.NewClient(grpcClient)
	return grpcClient, eventStreamsClient, closeFunc
}

func initializeEventStreamsWithGrpcClient(grpcClient connection.GrpcClient) *event_streams.Client {
	return event_streams.NewClient(grpcClient)
}
