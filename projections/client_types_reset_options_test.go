package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"

	"github.com/stretchr/testify/require"
)

func TestResetOptionsRequest_Build(t *testing.T) {
	t.Run("Non empty name", func(t *testing.T) {
		options := ResetOptionsRequest{
			ProjectionName: "name",
		}
		result := options.build()

		expectedResult := &projections.ResetReq{
			Options: &projections.ResetReq_Options{
				Name:            "name",
				WriteCheckpoint: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Non empty name with trailing spaces", func(t *testing.T) {
		options := ResetOptionsRequest{
			ProjectionName: " name ",
		}
		result := options.build()

		expectedResult := &projections.ResetReq{
			Options: &projections.ResetReq_Options{
				Name:            " name ",
				WriteCheckpoint: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("WriteCheckpoint set to false", func(t *testing.T) {
		options := ResetOptionsRequest{
			ProjectionName:  "name",
			WriteCheckpoint: false,
		}
		result := options.build()

		expectedResult := &projections.ResetReq{
			Options: &projections.ResetReq_Options{
				Name:            "name",
				WriteCheckpoint: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("WriteCheckpoint set to true", func(t *testing.T) {
		options := ResetOptionsRequest{
			ProjectionName:  "name",
			WriteCheckpoint: true,
		}
		result := options.build()

		expectedResult := &projections.ResetReq{
			Options: &projections.ResetReq_Options{
				Name:            "name",
				WriteCheckpoint: true,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Panics for empty name", func(t *testing.T) {
		options := ResetOptionsRequest{
			ProjectionName: "",
		}

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics for name consisting of spaces only", func(t *testing.T) {
		options := ResetOptionsRequest{
			ProjectionName: "      ",
		}

		require.Panics(t, func() {
			options.build()
		})
	})
}
