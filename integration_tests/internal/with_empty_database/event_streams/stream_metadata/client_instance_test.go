package stream_metadata

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/core/test_utils"
	"github.com/pivonroll/EventStore-Client-Go/event_streams"
)

func initializeContainerAndClient(t *testing.T,
	envVariableOverrides map[string]string) (*event_streams.Client, test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClient(t, envVariableOverrides)

	client := event_streams.NewClient(grpcClient)
	return client, closeFunc
}

func initializeContainerAndClientWithCredentials(t *testing.T,
	username string,
	password string, envVariableOverrides map[string]string) (*event_streams.Client, test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClientWithCredentials(t, username, password, envVariableOverrides)

	client := event_streams.NewClient(grpcClient)
	return client, closeFunc
}
