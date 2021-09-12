package projections

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pivonroll/EventStore-Client-Go/connection"
	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/stretchr/testify/require"
)

func TestClientFactoryImpl_CreateClient(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	grpcClient := connection.NewMockGrpcClient(ctrl)
	projectionsProtoClient := projections.NewMockProjectionsClient(ctrl)

	factory := ClientFactoryImpl{}
	result := factory.CreateClient(grpcClient, projectionsProtoClient)

	expectedResult := newClientImpl(grpcClient, projectionsProtoClient)

	require.Equal(t, expectedResult, result)
}
