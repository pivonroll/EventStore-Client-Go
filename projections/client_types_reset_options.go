package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"
)

// ResetRequest contains projection's name and whether
// to write a checkpoint or not when projection is reset.
type ResetRequest struct {
	ProjectionName  string
	WriteCheckpoint bool
}

func (resetOptionsRequest *ResetRequest) build() *projections.ResetReq {
	if strings.TrimSpace(resetOptionsRequest.ProjectionName) == "" {
		panic("Failed to build ResetRequest. Trimmed projection name is an empty string")
	}

	result := &projections.ResetReq{
		Options: &projections.ResetReq_Options{
			Name:            resetOptionsRequest.ProjectionName,
			WriteCheckpoint: resetOptionsRequest.WriteCheckpoint,
		},
	}

	return result
}
