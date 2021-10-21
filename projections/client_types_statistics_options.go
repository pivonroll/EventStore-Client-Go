package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

// IsStatisticsProjection is general interface type which is used to represent all statistics modes
// through which we can select to fetch statistics of specific projections.
type IsStatisticsProjection interface {
	isStatisticsProjectionType()
}

// StatisticsForAllProjections fetch statistics for all projections.
type StatisticsForAllProjections struct{}

func (s StatisticsForAllProjections) isStatisticsProjectionType() {
}

// StatisticsForProjectionByName fetch statistics for a specific projection.
type StatisticsForProjectionByName struct {
	ProjectionName string
}

func (s StatisticsForProjectionByName) isStatisticsProjectionType() {
}

// StatisticsForTransientProjections fetch statistics for all transient projections.
type StatisticsForTransientProjections struct{}

func (s StatisticsForTransientProjections) isStatisticsProjectionType() {
}

// StatisticsForContinuousProjections fetch statistics for all continuous projections.
type StatisticsForContinuousProjections struct{}

func (s StatisticsForContinuousProjections) isStatisticsProjectionType() {
}

// StatisticsForOneTimeProjections fetch statistics for all one-time projections.
type StatisticsForOneTimeProjections struct{}

func (s StatisticsForOneTimeProjections) isStatisticsProjectionType() {
}

func buildStatisticsRequest(mode IsStatisticsProjection) *projections.StatisticsReq {
	result := &projections.StatisticsReq{
		Options: &projections.StatisticsReq_Options{
			Mode: nil,
		},
	}

	switch mode.(type) {
	case StatisticsForAllProjections:
		result.Options.Mode = &projections.StatisticsReq_Options_All{
			All: &shared.Empty{},
		}
	case StatisticsForTransientProjections:
		result.Options.Mode = &projections.StatisticsReq_Options_Transient{
			Transient: &shared.Empty{},
		}
	case StatisticsForContinuousProjections:
		result.Options.Mode = &projections.StatisticsReq_Options_Continuous{
			Continuous: &shared.Empty{},
		}
	case StatisticsForOneTimeProjections:
		result.Options.Mode = &projections.StatisticsReq_Options_OneTime{
			OneTime: &shared.Empty{},
		}
	case StatisticsForProjectionByName:
		mode := mode.(StatisticsForProjectionByName)
		if strings.TrimSpace(mode.ProjectionName) == "" {
			panic("Failed to build StatisticsOptionsRequest. Trimmed projection name is an empty string")
		}

		result.Options.Mode = &projections.StatisticsReq_Options_Name{
			Name: mode.ProjectionName,
		}
	}

	return result
}
