package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"

	"github.com/stretchr/testify/require"
)

func TestDeleteOptionsRequest_Build(t *testing.T) {
	t.Run("Set non empty name without trailing spaces", func(t *testing.T) {
		options := DeleteRequest{
			ProjectionName:         "name",
			DeleteEmittedStreams:   false,
			DeleteStateStream:      false,
			DeleteCheckpointStream: false,
		}

		result := options.build()

		expectedResult := &projections.DeleteReq{
			Options: &projections.DeleteReq_Options{
				Name:                   "name",
				DeleteEmittedStreams:   false,
				DeleteStateStream:      false,
				DeleteCheckpointStream: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Set non empty name with trailing spaces", func(t *testing.T) {
		options := DeleteRequest{
			ProjectionName:         " name ",
			DeleteEmittedStreams:   false,
			DeleteStateStream:      false,
			DeleteCheckpointStream: false,
		}

		result := options.build()

		expectedResult := &projections.DeleteReq{
			Options: &projections.DeleteReq_Options{
				Name:                   " name ",
				DeleteEmittedStreams:   false,
				DeleteStateStream:      false,
				DeleteCheckpointStream: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Set non empty name and DeleteEmittedStreams true", func(t *testing.T) {
		options := DeleteRequest{
			ProjectionName:         "name",
			DeleteEmittedStreams:   true,
			DeleteStateStream:      false,
			DeleteCheckpointStream: false,
		}

		result := options.build()

		expectedResult := &projections.DeleteReq{
			Options: &projections.DeleteReq_Options{
				Name:                   "name",
				DeleteEmittedStreams:   true,
				DeleteStateStream:      false,
				DeleteCheckpointStream: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Set non empty name and SetDeleteStateStream true", func(t *testing.T) {
		options := DeleteRequest{
			ProjectionName:         "name",
			DeleteEmittedStreams:   false,
			DeleteStateStream:      true,
			DeleteCheckpointStream: false,
		}

		result := options.build()

		expectedResult := &projections.DeleteReq{
			Options: &projections.DeleteReq_Options{
				Name:                   "name",
				DeleteEmittedStreams:   false,
				DeleteStateStream:      true,
				DeleteCheckpointStream: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Set non empty name and SetDeleteStateStream true", func(t *testing.T) {
		options := DeleteRequest{
			ProjectionName:         "name",
			DeleteEmittedStreams:   false,
			DeleteStateStream:      false,
			DeleteCheckpointStream: true,
		}

		result := options.build()

		expectedResult := &projections.DeleteReq{
			Options: &projections.DeleteReq_Options{
				Name:                   "name",
				DeleteEmittedStreams:   false,
				DeleteStateStream:      false,
				DeleteCheckpointStream: true,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Panics for empty name", func(t *testing.T) {
		options := DeleteRequest{
			ProjectionName: "",
		}
		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics for name consisting of spaces only", func(t *testing.T) {
		options := DeleteRequest{
			ProjectionName: "     ",
		}
		require.Panics(t, func() {
			options.build()
		})
	})
}
