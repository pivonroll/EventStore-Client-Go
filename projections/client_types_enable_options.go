package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"
)

type enableOptionsRequest string

func (enableOptionsRequest enableOptionsRequest) build() *projections.EnableReq {
	if strings.TrimSpace(string(enableOptionsRequest)) == "" {
		panic("Failed to build EnableOptionsRequest. Trimmed name is an empty string")
	}

	result := &projections.EnableReq{
		Options: &projections.EnableReq_Options{
			Name: string(enableOptionsRequest),
		},
	}

	return result
}
