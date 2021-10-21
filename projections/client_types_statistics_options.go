package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

type StatisticsOptionsRequestModeType string

const (
	StatisticsOptionsRequestModeAllType        StatisticsOptionsRequestModeType = "StatisticsOptionsRequestModeAllType"
	StatisticsOptionsRequestModeNameType       StatisticsOptionsRequestModeType = "StatisticsOptionsRequestModeNameType"
	StatisticsOptionsRequestModeTransientType  StatisticsOptionsRequestModeType = "StatisticsOptionsRequestModeTransientType"
	StatisticsOptionsRequestModeContinuousType StatisticsOptionsRequestModeType = "StatisticsOptionsRequestModeContinuousType"
	StatisticsOptionsRequestModeOneTimeType    StatisticsOptionsRequestModeType = "StatisticsOptionsRequestModeOneTimeType"
)

type StatisticsOptionsRequestMode interface {
	GetType() StatisticsOptionsRequestModeType
}

type StatisticsOptionsRequestModeAll struct{}

func (s StatisticsOptionsRequestModeAll) GetType() StatisticsOptionsRequestModeType {
	return StatisticsOptionsRequestModeAllType
}

type StatisticsOptionsRequestModeName struct {
	Name string
}

func (s StatisticsOptionsRequestModeName) GetType() StatisticsOptionsRequestModeType {
	return StatisticsOptionsRequestModeNameType
}

type StatisticsOptionsRequestModeTransient struct{}

func (s StatisticsOptionsRequestModeTransient) GetType() StatisticsOptionsRequestModeType {
	return StatisticsOptionsRequestModeTransientType
}

type StatisticsOptionsRequestModeContinuous struct{}

func (s StatisticsOptionsRequestModeContinuous) GetType() StatisticsOptionsRequestModeType {
	return StatisticsOptionsRequestModeContinuousType
}

type StatisticsOptionsRequestModeOneTime struct{}

func (s StatisticsOptionsRequestModeOneTime) GetType() StatisticsOptionsRequestModeType {
	return StatisticsOptionsRequestModeOneTimeType
}

func buildStatisticsRequest(mode StatisticsOptionsRequestMode) *projections.StatisticsReq {
	result := &projections.StatisticsReq{
		Options: &projections.StatisticsReq_Options{
			Mode: nil,
		},
	}

	if mode.GetType() == StatisticsOptionsRequestModeAllType {
		result.Options.Mode = &projections.StatisticsReq_Options_All{
			All: &shared.Empty{},
		}
	} else if mode.GetType() == StatisticsOptionsRequestModeTransientType {
		result.Options.Mode = &projections.StatisticsReq_Options_Transient{
			Transient: &shared.Empty{},
		}
	} else if mode.GetType() == StatisticsOptionsRequestModeContinuousType {
		result.Options.Mode = &projections.StatisticsReq_Options_Continuous{
			Continuous: &shared.Empty{},
		}
	} else if mode.GetType() == StatisticsOptionsRequestModeOneTimeType {
		result.Options.Mode = &projections.StatisticsReq_Options_OneTime{
			OneTime: &shared.Empty{},
		}
	} else if mode.GetType() == StatisticsOptionsRequestModeNameType {
		mode := mode.(StatisticsOptionsRequestModeName)
		if strings.TrimSpace(mode.Name) == "" {
			panic("Failed to build StatisticsOptionsRequest. Trimmed name is an empty string")
		}

		result.Options.Mode = &projections.StatisticsReq_Options_Name{
			Name: mode.Name,
		}
	}

	return result
}
