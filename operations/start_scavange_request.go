package operations

import "github.com/pivonroll/EventStore-Client-Go/protos/v21.6/operations"

// StartScavengeRequest are set of options which can be set for a scavenge operation.
type StartScavengeRequest struct {
	ThreadCount    int32
	StartFromChunk int32
}

func (request StartScavengeRequest) build() *operations.StartScavengeReq {
	return &operations.StartScavengeReq{
		Options: &operations.StartScavengeReq_Options{
			ThreadCount:    request.ThreadCount,
			StartFromChunk: request.StartFromChunk,
		},
	}
}
