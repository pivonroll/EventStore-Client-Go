package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"

	"github.com/stretchr/testify/require"
)

func TestAbortOptionsRequest_Build(t *testing.T) {
	t.Run("Non empty name", func(t *testing.T) {
		options := abortOptionsRequest("some name")
		result := options.build()

		expectedResult := &projections.DisableReq{
			Options: &projections.DisableReq_Options{
				Name:            "some name",
				WriteCheckpoint: true,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Non empty name with trailing spaces", func(t *testing.T) {
		options := abortOptionsRequest(" some name ")
		result := options.build()

		expectedResult := &projections.DisableReq{
			Options: &projections.DisableReq_Options{
				Name:            " some name ",
				WriteCheckpoint: true,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Panics for empty name", func(t *testing.T) {
		options := abortOptionsRequest("")
		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics for name consisting of spaces only", func(t *testing.T) {
		options := abortOptionsRequest("    ")
		require.Panics(t, func() {
			options.build()
		})
	})
}
