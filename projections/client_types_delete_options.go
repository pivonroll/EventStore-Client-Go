package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type DeleteOptionsRequest struct {
	Name                   string
	DeleteEmittedStreams   bool
	DeleteStateStream      bool
	DeleteCheckpointStream bool
}

func (deleteOptions *DeleteOptionsRequest) build() *projections.DeleteReq {
	if strings.TrimSpace(deleteOptions.Name) == "" {
		panic("Failed to build DeleteOptionsRequest. Trimmed name is an empty string")
	}

	result := &projections.DeleteReq{
		Options: &projections.DeleteReq_Options{
			Name:                   deleteOptions.Name,
			DeleteEmittedStreams:   deleteOptions.DeleteEmittedStreams,
			DeleteStateStream:      deleteOptions.DeleteStateStream,
			DeleteCheckpointStream: deleteOptions.DeleteCheckpointStream,
		},
	}

	return result
}
