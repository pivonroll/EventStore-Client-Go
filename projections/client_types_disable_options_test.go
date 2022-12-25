package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"

	"github.com/stretchr/testify/require"
)

func TestDisableOptionsRequest_Build(t *testing.T) {
	t.Run("Non empty name", func(t *testing.T) {
		options := disableOptionsRequest("name")
		result := options.build()

		expectedResult := &projections.DisableReq{
			Options: &projections.DisableReq_Options{
				Name:            "name",
				WriteCheckpoint: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Non empty name with trailing spaces", func(t *testing.T) {
		options := disableOptionsRequest(" name ")
		result := options.build()

		expectedResult := &projections.DisableReq{
			Options: &projections.DisableReq_Options{
				Name:            " name ",
				WriteCheckpoint: false,
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Panics for empty name", func(t *testing.T) {
		options := disableOptionsRequest("")

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics for name consisting of spaces only", func(t *testing.T) {
		options := disableOptionsRequest("    ")

		require.Panics(t, func() {
			options.build()
		})
	})
}
