package projections

import (
	"strings"

	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/projections"
)

// ResultRequest represents input required to fetch a result of a projection.
type ResultRequest struct {
	ProjectionName string
	Partition      string
}

func (resultOptionsRequest *ResultRequest) build() *projections.ResultReq {
	if strings.TrimSpace(resultOptionsRequest.ProjectionName) == "" {
		panic("Failed to build ResultRequest. Trimmed projection name is an empty string")
	}

	result := &projections.ResultReq{
		Options: &projections.ResultReq_Options{
			Name:      resultOptionsRequest.ProjectionName,
			Partition: resultOptionsRequest.Partition,
		},
	}

	return result
}
