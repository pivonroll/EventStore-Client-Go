package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"

	"github.com/stretchr/testify/require"
)

func TestResultOptionsRequest_Build(t *testing.T) {
	t.Run("Non empty name", func(t *testing.T) {
		options := ResultOptionsRequest{
			ProjectionName: "name",
		}
		result := options.build()

		expectedResult := &projections.ResultReq{
			Options: &projections.ResultReq_Options{
				Name: "name",
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Non empty name with trailing spaces", func(t *testing.T) {
		options := ResultOptionsRequest{
			ProjectionName: " name ",
		}
		result := options.build()

		expectedResult := &projections.ResultReq{
			Options: &projections.ResultReq_Options{
				Name: " name ",
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Non empty name and partition", func(t *testing.T) {
		options := ResultOptionsRequest{
			ProjectionName: "name",
			Partition:      "partition",
		}
		result := options.build()

		expectedResult := &projections.ResultReq{
			Options: &projections.ResultReq_Options{
				Name:      "name",
				Partition: "partition",
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Panics for empty name", func(t *testing.T) {
		options := ResultOptionsRequest{
			ProjectionName: "",
		}

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics for name consisting of spaces only", func(t *testing.T) {
		options := ResultOptionsRequest{
			ProjectionName: "     ",
		}

		require.Panics(t, func() {
			options.build()
		})
	})
}
