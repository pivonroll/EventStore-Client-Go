package user_management_integration_test

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/core/test_utils"
	"github.com/pivonroll/EventStore-Client-Go/user_management"
)

func initializeContainerAndClient(t *testing.T,
	envVariableOverrides map[string]string) (*user_management.Client,
	test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClient(t, envVariableOverrides)

	client := user_management.NewClient(grpcClient)

	return client, closeFunc
}
