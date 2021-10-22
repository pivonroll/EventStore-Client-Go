package persistent_integration_test

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/core/test_utils"
	"github.com/pivonroll/EventStore-Client-Go/event_streams"
	"github.com/pivonroll/EventStore-Client-Go/persistent"
)

func initializeContainerAndClient(t *testing.T,
	envVariableOverrides map[string]string) (*persistent.Client,
	*event_streams.Client,
	test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClient(t, envVariableOverrides)

	client := persistent.NewClient(grpcClient)

	eventStreamsClient := event_streams.NewClient(grpcClient)

	return client, eventStreamsClient, closeFunc
}

func initializeWithPrePopulatedDatabase(t *testing.T) (*persistent.Client, test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClientWithPrePopulatedDatabase(t)
	client := persistent.NewClient(grpcClient)
	return client, closeFunc
}
