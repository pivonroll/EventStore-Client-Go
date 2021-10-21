package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

type ProjectionType string

const (
	OneTimeProjectionType     ProjectionType = "OneTimeProjectionType"
	TransientProjectionType   ProjectionType = "TransientProjectionType"
	ContinuousProjectionsType ProjectionType = "ContinuousProjectionsType"
)

type isProjectionType interface {
	GetType() ProjectionType
}

// OneTimeProjection mode instructs EventStoreDB to create one-time projection.
type OneTimeProjection struct{}

// GetType returns a mode type. In this case it will be OneTimeProjectionType.
// See constant OneTimeProjectionType.
func (mode OneTimeProjection) GetType() ProjectionType {
	return OneTimeProjectionType
}

// TransientProjection mode is used to create transient projection.
type TransientProjection struct {
	ProjectionName string
}

// GetType returns a mode type. In this case it will be TransientProjectionType.
// See constant TransientProjectionType.
func (mode TransientProjection) GetType() ProjectionType {
	return TransientProjectionType
}

// ContinuousProjection mode is used to create a continuous projection.
type ContinuousProjection struct {
	ProjectionName      string
	TrackEmittedStreams bool
}

// GetType returns a mode type. In this case it will be ContinuousProjectionsType.
// See constant ContinuousProjectionsType.
func (mode ContinuousProjection) GetType() ProjectionType {
	return ContinuousProjectionsType
}

// CreateRequest represents data required to create a projection at EventStoreDB.
type CreateRequest struct {
	Mode  isProjectionType
	Query string
}

func (createConfig *CreateRequest) build() *projections.CreateReq {
	if strings.TrimSpace(createConfig.Query) == "" {
		panic("Failed to build CreateRequest. Trimmed query is an empty string")
	}

	result := &projections.CreateReq{
		Options: &projections.CreateReq_Options{
			Mode:  nil,
			Query: createConfig.Query,
		},
	}

	if createConfig.Mode.GetType() == OneTimeProjectionType {
		result.Options.Mode = &projections.CreateReq_Options_OneTime{
			OneTime: &shared.Empty{},
		}
	} else if createConfig.Mode.GetType() == TransientProjectionType {
		transientOption := createConfig.Mode.(TransientProjection)
		result.Options.Mode = &projections.CreateReq_Options_Transient_{
			Transient: &projections.CreateReq_Options_Transient{
				Name: transientOption.ProjectionName,
			},
		}
	} else if createConfig.Mode.GetType() == ContinuousProjectionsType {
		continuousOption := createConfig.Mode.(ContinuousProjection)
		result.Options.Mode = &projections.CreateReq_Options_Continuous_{
			Continuous: &projections.CreateReq_Options_Continuous{
				Name:                continuousOption.ProjectionName,
				TrackEmittedStreams: continuousOption.TrackEmittedStreams,
			},
		}
	}

	return result
}
