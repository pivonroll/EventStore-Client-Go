package operations_integration_test

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/core/test_utils"
	"github.com/pivonroll/EventStore-Client-Go/operations"
)

func initializeClient(t *testing.T,
	envVariableOverrides map[string]string) (operations.Client,
	test_utils.CloseFunc) {
	grpcClient, closeFunc := test_utils.InitializeGrpcClient(t, envVariableOverrides)

	client := operations.ClientFactoryImpl{}.Create(grpcClient)

	return client, closeFunc
}
