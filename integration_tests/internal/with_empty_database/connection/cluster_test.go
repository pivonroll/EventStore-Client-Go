package connection_integration_test

import (
	"context"
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/persistent"
	"github.com/stretchr/testify/require"
)

func Test_NotLeaderExceptionButWorkAfterRetry(t *testing.T) {
	// We purposely connect to a follower node so we can trigger on not leader exception.
	clientUri := "esdb://admin:changeit@localhost:2111,localhost:2112,localhost:2113?nodepreference=follower&tlsverifycert=false"
	config, err := connection.ParseConnectionString(clientUri)
	require.NoError(t, err)
	config.GossipTimeout = 100
	grpcClient := connection.NewGrpcClient(*config)

	persistentSubscriptionClient := persistent.NewClient(grpcClient)

	persistentCreateConfig := persistent.SubscriptionGroupForStreamRequest{
		StreamId:  "myfoobar_123456",
		GroupName: "a_group",
		Revision: stream_revision.ReadStreamRevision{
			Revision: 0,
		},
		Settings: persistent.DefaultRequestSettings,
	}

	err = persistentSubscriptionClient.CreateSubscriptionGroupForStream(context.Background(),
		persistentCreateConfig)
	require.Error(t, err)

	// It should work now as the client automatically reconnected to the leader node.
	err = persistentSubscriptionClient.CreateSubscriptionGroupForStream(context.Background(), persistentCreateConfig)
	require.NoError(t, err)
}
