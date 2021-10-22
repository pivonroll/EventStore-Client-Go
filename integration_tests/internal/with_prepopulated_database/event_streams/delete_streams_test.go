package event_streams_with_prepopulated_database

import (
	"context"
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/stretchr/testify/require"
)

func Test_DeleteStream(t *testing.T) {
	client, closeFunc := initializeWithPrePopulatedDatabase(t)
	defer closeFunc()

	result, err := client.DeleteStream(context.Background(),
		"dataset20M-1800",
		stream_revision.WriteStreamRevision{Revision: 1999})
	require.NoError(t, err)

	position, isPosition := result.GetPosition()
	require.True(t, isPosition)
	require.True(t, position.PreparePosition > 0)
	require.True(t, position.CommitPosition > 0)
}
