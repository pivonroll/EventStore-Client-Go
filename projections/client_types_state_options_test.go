package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/projections"

	"github.com/stretchr/testify/require"
)

func TestStateOptionsRequest_Build(t *testing.T) {
	t.Run("Non empty name", func(t *testing.T) {
		options := StateRequest{
			ProjectionName: "name",
		}
		result := options.build()

		expectedState := &projections.StateReq{
			Options: &projections.StateReq_Options{
				Name: "name",
			},
		}

		require.Equal(t, expectedState, result)
	})

	t.Run("Non empty name with trailing spaces", func(t *testing.T) {
		options := StateRequest{
			ProjectionName: " name ",
		}
		result := options.build()

		expectedState := &projections.StateReq{
			Options: &projections.StateReq_Options{
				Name: " name ",
			},
		}

		require.Equal(t, expectedState, result)
	})

	t.Run("Non empty name and partition", func(t *testing.T) {
		options := StateRequest{
			ProjectionName: "name",
			Partition:      "partition",
		}
		result := options.build()

		expectedState := &projections.StateReq{
			Options: &projections.StateReq_Options{
				Name:      "name",
				Partition: "partition",
			},
		}

		require.Equal(t, expectedState, result)
	})

	t.Run("Panics for empty name", func(t *testing.T) {
		options := StateRequest{
			ProjectionName: "",
		}

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics for name consisting of spaces only", func(t *testing.T) {
		options := StateRequest{
			ProjectionName: "     ",
		}

		require.Panics(t, func() {
			options.build()
		})
	})
}
