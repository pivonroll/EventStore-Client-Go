package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

type CreateConfigModeType string

const (
	CreateConfigModeOneTimeOptionType    CreateConfigModeType = "CreateConfigModeOneTimeOptionType"
	CreateConfigModeTransientOptionType  CreateConfigModeType = "CreateConfigModeTransientOptionType"
	CreateConfigModeContinuousOptionType CreateConfigModeType = "CreateConfigModeContinuousOptionType"
)

type CreateConfigMode interface {
	GetType() CreateConfigModeType
}

type CreateConfigModeOneTimeOption struct{}

func (mode CreateConfigModeOneTimeOption) GetType() CreateConfigModeType {
	return CreateConfigModeOneTimeOptionType
}

type CreateConfigModeTransientOption struct {
	ProjectionName string
}

func (mode CreateConfigModeTransientOption) GetType() CreateConfigModeType {
	return CreateConfigModeTransientOptionType
}

type CreateConfigModeContinuousOption struct {
	ProjectionName      string
	TrackEmittedStreams bool
}

func (mode CreateConfigModeContinuousOption) GetType() CreateConfigModeType {
	return CreateConfigModeContinuousOptionType
}

type CreateOptionsRequest struct {
	Mode  CreateConfigMode
	Query string
}

func (createConfig *CreateOptionsRequest) build() *projections.CreateReq {
	if strings.TrimSpace(createConfig.Query) == "" {
		panic("Failed to build CreateOptionsRequest. Trimmed query is an empty string")
	}

	result := &projections.CreateReq{
		Options: &projections.CreateReq_Options{
			Mode:  nil,
			Query: createConfig.Query,
		},
	}

	if createConfig.Mode.GetType() == CreateConfigModeOneTimeOptionType {
		result.Options.Mode = &projections.CreateReq_Options_OneTime{
			OneTime: &shared.Empty{},
		}
	} else if createConfig.Mode.GetType() == CreateConfigModeTransientOptionType {
		transientOption := createConfig.Mode.(CreateConfigModeTransientOption)
		result.Options.Mode = &projections.CreateReq_Options_Transient_{
			Transient: &projections.CreateReq_Options_Transient{
				Name: transientOption.ProjectionName,
			},
		}
	} else if createConfig.Mode.GetType() == CreateConfigModeContinuousOptionType {
		continuousOption := createConfig.Mode.(CreateConfigModeContinuousOption)
		result.Options.Mode = &projections.CreateReq_Options_Continuous_{
			Continuous: &projections.CreateReq_Options_Continuous{
				Name:                continuousOption.ProjectionName,
				TrackEmittedStreams: continuousOption.TrackEmittedStreams,
			},
		}
	}

	return result
}
