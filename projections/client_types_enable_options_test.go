package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"

	"github.com/stretchr/testify/require"
)

func TestEnableOptionsRequest_Build(t *testing.T) {
	t.Run("Non empty name", func(t *testing.T) {
		options := enableOptionsRequest("name")
		result := options.build()

		expectedResult := &projections.EnableReq{
			Options: &projections.EnableReq_Options{
				Name: "name",
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Non empty name with trailing spaces", func(t *testing.T) {
		options := enableOptionsRequest(" name ")
		result := options.build()

		expectedResult := &projections.EnableReq{
			Options: &projections.EnableReq_Options{
				Name: " name ",
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Panics for empty name", func(t *testing.T) {
		options := enableOptionsRequest("")

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics for name consisting of spaces only", func(t *testing.T) {
		options := enableOptionsRequest("    ")

		require.Panics(t, func() {
			options.build()
		})
	})
}
