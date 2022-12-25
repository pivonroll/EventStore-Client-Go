package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"
)

// DeleteRequest contains projection's name and a set of options what to do
// with emit, state and checkpoint stream.
type DeleteRequest struct {
	ProjectionName         string
	DeleteEmittedStreams   bool
	DeleteStateStream      bool
	DeleteCheckpointStream bool
}

func (deleteOptions *DeleteRequest) build() *projections.DeleteReq {
	if strings.TrimSpace(deleteOptions.ProjectionName) == "" {
		panic("Failed to build DeleteRequest. Trimmed projection name is an empty string")
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
