package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/projections"
)

type ResultOptionsRequest struct {
	ProjectionName string
	Partition      string
}

func (resultOptionsRequest *ResultOptionsRequest) build() *projections.ResultReq {
	if strings.TrimSpace(resultOptionsRequest.ProjectionName) == "" {
		panic("Failed to build ResultOptionsRequest. Trimmed projection name is an empty string")
	}

	result := &projections.ResultReq{
		Options: &projections.ResultReq_Options{
			Name:      resultOptionsRequest.ProjectionName,
			Partition: resultOptionsRequest.Partition,
		},
	}

	return result
}
