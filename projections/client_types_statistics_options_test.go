package projections

import (
	"testing"

	"github.com/pivonroll/EventStore-Client-Go/protos/shared"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"

	"github.com/stretchr/testify/require"
)

func TestStatisticsOptionsRequest_Build(t *testing.T) {
	t.Run("With Mode OneTime", func(t *testing.T) {
		result := buildStatisticsRequest(StatisticsOptionsRequestModeOneTime{})

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
		result := buildStatisticsRequest(StatisticsOptionsRequestModeAll{})

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
		result := buildStatisticsRequest(StatisticsOptionsRequestModeContinuous{})

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
		result := buildStatisticsRequest(StatisticsOptionsRequestModeTransient{})

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
		result := buildStatisticsRequest(StatisticsOptionsRequestModeName{
			Name: "name",
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
		result := buildStatisticsRequest(StatisticsOptionsRequestModeName{
			Name: " name ",
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
			buildStatisticsRequest(StatisticsOptionsRequestModeName{
				Name: "",
			})
		})
	})

	t.Run("With Mode Name, Panics with name consisting from spaces only", func(t *testing.T) {
		require.Panics(t, func() {
			buildStatisticsRequest(
				StatisticsOptionsRequestModeName{
					Name: "    ",
				})
		})
	})
}
