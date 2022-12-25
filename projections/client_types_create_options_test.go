package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/shared"

	"github.com/stretchr/testify/require"
)

func TestCreateOptionsRequest_Build(t *testing.T) {
	t.Run("Non empty query and mode set", func(t *testing.T) {
		options := CreateRequest{
			Mode:  OneTimeProjection{},
			Query: "some query",
		}

		result := options.build()

		expectedResult := &projections.CreateReq{
			Options: &projections.CreateReq_Options{
				Mode: &projections.CreateReq_Options_OneTime{
					OneTime: &shared.Empty{},
				},
				Query: "some query",
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Query with trailing space is not transformed", func(t *testing.T) {
		options := CreateRequest{
			Mode:  OneTimeProjection{},
			Query: " some query ",
		}

		result := options.build()

		expectedResult := &projections.CreateReq{
			Options: &projections.CreateReq_Options{
				Mode: &projections.CreateReq_Options_OneTime{
					OneTime: &shared.Empty{},
				},
				Query: " some query ",
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("Panics for empty query", func(t *testing.T) {
		options := CreateRequest{
			Mode:  OneTimeProjection{},
			Query: "",
		}

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics for query consisting of spaces only", func(t *testing.T) {
		options := CreateRequest{
			Mode:  OneTimeProjection{},
			Query: "    ",
		}

		require.Panics(t, func() {
			options.build()
		})
	})

	t.Run("Panics if mode is not set", func(t *testing.T) {
		options := CreateRequest{
			Query: "some query",
		}

		require.Panics(t, func() {
			options.build()
		})
	})
}
