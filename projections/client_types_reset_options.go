package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type ResetOptionsRequest struct {
	ProjectionName  string
	WriteCheckpoint bool
}

func (resetOptionsRequest *ResetOptionsRequest) build() *projections.ResetReq {
	if strings.TrimSpace(resetOptionsRequest.ProjectionName) == "" {
		panic("Failed to build ResetOptionsRequest. Trimmed name is an empty string")
	}

	result := &projections.ResetReq{
		Options: &projections.ResetReq_Options{
			Name:            resetOptionsRequest.ProjectionName,
			WriteCheckpoint: resetOptionsRequest.WriteCheckpoint,
		},
	}

	return result
}
