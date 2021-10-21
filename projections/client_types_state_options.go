package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

// StateRequest represents input required to fetch state of the projection.
type StateRequest struct {
	ProjectionName string
	Partition      string
}

func (stateOptionsRequest *StateRequest) build() *projections.StateReq {
	if strings.TrimSpace(stateOptionsRequest.ProjectionName) == "" {
		panic("Failed to build StateRequest. Trimmed name is an empty string")
	}

	result := &projections.StateReq{
		Options: &projections.StateReq_Options{
			Name:      stateOptionsRequest.ProjectionName,
			Partition: stateOptionsRequest.Partition,
		},
	}

	return result
}
