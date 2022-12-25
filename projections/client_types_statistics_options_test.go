package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/shared"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"

	"github.com/stretchr/testify/require"
)

func TestStatisticsOptionsRequest_Build(t *testing.T) {
	t.Run("With Mode OneTime", func(t *testing.T) {
		result := buildStatisticsRequest(StatisticsForOneTimeProjections{})

		expectedResult := &projections.StatisticsReq{
			Options: &projections.StatisticsReq_Options{
				Mode: &projections.StatisticsReq_Options_OneTime{
					OneTime: &shared.Empty{},
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("With Mode All", func(t *testing.T) {
		result := buildStatisticsRequest(StatisticsForAllProjections{})

		expectedResult := &projections.StatisticsReq{
			Options: &projections.StatisticsReq_Options{
				Mode: &projections.StatisticsReq_Options_All{
					All: &shared.Empty{},
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("With Mode Continuous", func(t *testing.T) {
		result := buildStatisticsRequest(StatisticsForContinuousProjections{})

		expectedResult := &projections.StatisticsReq{
			Options: &projections.StatisticsReq_Options{
				Mode: &projections.StatisticsReq_Options_Continuous{
					Continuous: &shared.Empty{},
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("With Mode Transient", func(t *testing.T) {
		result := buildStatisticsRequest(StatisticsForTransientProjections{})

		expectedResult := &projections.StatisticsReq{
			Options: &projections.StatisticsReq_Options{
				Mode: &projections.StatisticsReq_Options_Transient{
					Transient: &shared.Empty{},
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("With Mode Name", func(t *testing.T) {
		result := buildStatisticsRequest(StatisticsForProjectionByName{
			ProjectionName: "name",
		})

		expectedResult := &projections.StatisticsReq{
			Options: &projections.StatisticsReq_Options{
				Mode: &projections.StatisticsReq_Options_Name{
					Name: "name",
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("With Mode Name, non empty name with trailing spaces", func(t *testing.T) {
		result := buildStatisticsRequest(StatisticsForProjectionByName{
			ProjectionName: " name ",
		})

		expectedResult := &projections.StatisticsReq{
			Options: &projections.StatisticsReq_Options{
				Mode: &projections.StatisticsReq_Options_Name{
					Name: " name ",
				},
			},
		}

		require.Equal(t, expectedResult, result)
	})

	t.Run("With Mode Name, Panics with empty name", func(t *testing.T) {
		require.Panics(t, func() {
			buildStatisticsRequest(StatisticsForProjectionByName{
				ProjectionName: "",
			})
		})
	})

	t.Run("With Mode Name, Panics with name consisting from spaces only", func(t *testing.T) {
		require.Panics(t, func() {
			buildStatisticsRequest(
				StatisticsForProjectionByName{
					ProjectionName: "    ",
				})
		})
	})
}
