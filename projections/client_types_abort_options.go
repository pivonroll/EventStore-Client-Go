package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/v22.10/projections"
)

type abortOptionsRequest string

func (abortOptionsRequest abortOptionsRequest) build() *projections.DisableReq {
	if strings.TrimSpace(string(abortOptionsRequest)) == "" {
		panic("Failed to build AbortOptionsRequest. Trimmed name is an empty string")
	}

	result := &projections.DisableReq{
		Options: &projections.DisableReq_Options{
			Name:            string(abortOptionsRequest),
			WriteCheckpoint: true,
		},
	}

	return result
}
