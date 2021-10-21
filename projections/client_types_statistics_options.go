package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

type StatisticsForProjectionType string

const (
	StatisticsForAllProjectionsType        StatisticsForProjectionType = "StatisticsForAllProjectionsType"
	StatisticsForProjectionByNameType      StatisticsForProjectionType = "StatisticsForProjectionByNameType"
	StatisticsForTransientProjectionsType  StatisticsForProjectionType = "StatisticsForTransientProjectionsType"
	StatisticsForContinuousProjectionsType StatisticsForProjectionType = "StatisticsForContinuousProjectionsType"
	StatisticsForOneTimeProjectionsType    StatisticsForProjectionType = "StatisticsForOneTimeProjectionsType"
)

type IsStatisticsByProjection interface {
	GetType() StatisticsForProjectionType
}

type StatisticsForAllProjections struct{}

func (s StatisticsForAllProjections) GetType() StatisticsForProjectionType {
	return StatisticsForAllProjectionsType
}

type StatisticsForProjectionByName struct {
	ProjectionName string
}

func (s StatisticsForProjectionByName) GetType() StatisticsForProjectionType {
	return StatisticsForProjectionByNameType
}

type StatisticsForTransientProjections struct{}

func (s StatisticsForTransientProjections) GetType() StatisticsForProjectionType {
	return StatisticsForTransientProjectionsType
}

type StatisticsForContinuousProjections struct{}

func (s StatisticsForContinuousProjections) GetType() StatisticsForProjectionType {
	return StatisticsForContinuousProjectionsType
}

type StatisticsForOneTimeProjections struct{}

func (s StatisticsForOneTimeProjections) GetType() StatisticsForProjectionType {
	return StatisticsForOneTimeProjectionsType
}

func buildStatisticsRequest(mode IsStatisticsByProjection) *projections.StatisticsReq {
	result := &projections.StatisticsReq{
		Options: &projections.StatisticsReq_Options{
			Mode: nil,
		},
	}

	if mode.GetType() == StatisticsForAllProjectionsType {
		result.Options.Mode = &projections.StatisticsReq_Options_All{
			All: &shared.Empty{},
		}
	} else if mode.GetType() == StatisticsForTransientProjectionsType {
		result.Options.Mode = &projections.StatisticsReq_Options_Transient{
			Transient: &shared.Empty{},
		}
	} else if mode.GetType() == StatisticsForContinuousProjectionsType {
		result.Options.Mode = &projections.StatisticsReq_Options_Continuous{
			Continuous: &shared.Empty{},
		}
	} else if mode.GetType() == StatisticsForOneTimeProjectionsType {
		result.Options.Mode = &projections.StatisticsReq_Options_OneTime{
			OneTime: &shared.Empty{},
		}
	} else if mode.GetType() == StatisticsForProjectionByNameType {
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
