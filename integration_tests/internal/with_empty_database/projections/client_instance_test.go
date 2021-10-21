package projections_integration_test

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/event_streams"
	"github.com/pivonroll/EventStore-Client-Go/projections"
	"github.com/pivonroll/EventStore-Client-Go/test_utils"
)

func initializeContainerAndClient(t *testing.T) (projections.Client,
	test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClient(t,
		map[string]string{
			"EVENTSTORE_RUN_PROJECTIONS":            "All",
			"EVENTSTORE_START_STANDARD_PROJECTIONS": "true",
		})

	client := projections.NewClientImpl(grpcClient)

	return client, closeFunc
}

func initializeClientAndEventStreamsClient(t *testing.T) (projections.Client,
	*event_streams.Client,
	test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClient(t,
		map[string]string{
			"EVENTSTORE_RUN_PROJECTIONS":            "All",
			"EVENTSTORE_START_STANDARD_PROJECTIONS": "true",
		})

	client := projections.NewClientImpl(grpcClient)
	eventStreamsClient := event_streams.NewClient(grpcClient)

	return client, eventStreamsClient, closeFunc
}

func initializeContainerAndClientWithCredentials(t *testing.T,
	username string,
	password string, envVariableOverrides map[string]string) (projections.Client, test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClientWithCredentials(t, username, password, envVariableOverrides)

	client := projections.NewClientImpl(grpcClient)
	return client, closeFunc
}
