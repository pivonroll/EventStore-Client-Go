package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type StateOptionsRequest struct {
	ProjectionName string
	Partition      string
}

func (stateOptionsRequest *StateOptionsRequest) build() *projections.StateReq {
	if strings.TrimSpace(stateOptionsRequest.ProjectionName) == "" {
		panic("Failed to build StateOptionsRequest. Trimmed name is an empty string")
	}

	result := &projections.StateReq{
		Options: &projections.StateReq_Options{
			Name:      stateOptionsRequest.ProjectionName,
			Partition: stateOptionsRequest.Partition,
		},
	}

	return result
}
