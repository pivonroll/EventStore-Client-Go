package event_streams_with_prepopulated_database

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/core/test_utils"
	"github.com/pivonroll/EventStore-Client-Go/event_streams"
)

func initializeWithPrePopulatedDatabase(t *testing.T) (*event_streams.Client, test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClientWithPrePopulatedDatabase(t)
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
