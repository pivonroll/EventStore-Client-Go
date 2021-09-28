package persistent_integration_test

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/event_streams"
	"github.com/pivonroll/EventStore-Client-Go/persistent"
	"github.com/pivonroll/EventStore-Client-Go/test_container"
)

func initializeContainerAndClient(t *testing.T,
	envVariableOverrides map[string]string) (persistent.Client,
	event_streams.Client,
	test_container.CloseFunc) {
	grpcClient, closeFunc := test_container.InitializeContainerAndGrpcClient(t, envVariableOverrides)

	client := persistent.ClientFactoryImpl{}.CreateClient(grpcClient)

	eventStreamsClient := event_streams.ClientFactoryImpl{}.CreateClient(grpcClient)

	return client, eventStreamsClient, closeFunc
}

func initializeWithPrePopulatedDatabase(t *testing.T) (persistent.Client, test_container.CloseFunc) {
	grpcClient, closeFunc := test_container.InitializeGrpcClientWithPrePopulatedDatabase(t)
	client := persistent.ClientFactoryImpl{}.CreateClient(grpcClient)
	return client, closeFunc
}
