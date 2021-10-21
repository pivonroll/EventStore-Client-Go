package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type DeleteOptionsRequest struct {
	ProjectionName         string
	DeleteEmittedStreams   bool
	DeleteStateStream      bool
	DeleteCheckpointStream bool
}

func (deleteOptions *DeleteOptionsRequest) build() *projections.DeleteReq {
	if strings.TrimSpace(deleteOptions.ProjectionName) == "" {
		panic("Failed to build DeleteOptionsRequest. Trimmed name is an empty string")
	}

	result := &projections.DeleteReq{
		Options: &projections.DeleteReq_Options{
			Name:                   deleteOptions.ProjectionName,
			DeleteEmittedStreams:   deleteOptions.DeleteEmittedStreams,
			DeleteStateStream:      deleteOptions.DeleteStateStream,
			DeleteCheckpointStream: deleteOptions.DeleteCheckpointStream,
		},
	}

	return result
}
