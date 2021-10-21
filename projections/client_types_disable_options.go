package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type disableOptionsRequest string

func (disableOptionsRequest disableOptionsRequest) build() *projections.DisableReq {
	if strings.TrimSpace(string(disableOptionsRequest)) == "" {
		panic("Failed to build DisableOptionsRequest. Trimmed name is an empty string")
	}

	result := &projections.DisableReq{
		Options: &projections.DisableReq_Options{
			Name:            string(disableOptionsRequest),
			WriteCheckpoint: false,
		},
	}

	return result
}
