package event_streams

import (
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/shared"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/streams"
)

type deleteRequest struct {
	streamId string
	// Types that are assignable to expectedStreamRevision:
	// WriteStreamRevision
	// WriteStreamRevisionNoStream
	// WriteStreamRevisionAny
	// WriteStreamRevisionStreamExists
	expectedStreamRevision stream_revision.IsWriteStreamRevision
}

func (this deleteRequest) build() *streams.DeleteReq {
	result := &streams.DeleteReq{
		Options: &streams.DeleteReq_Options{
			StreamIdentifier: &shared.StreamIdentifier{
				StreamName: []byte(this.streamId),
			},
			ExpectedStreamRevision: nil,
		},
	}

	switch this.expectedStreamRevision.(type) {
	case stream_revision.WriteStreamRevision:
		revision := this.expectedStreamRevision.(stream_revision.WriteStreamRevision)
		result.Options.ExpectedStreamRevision = &streams.DeleteReq_Options_Revision{
			Revision: revision.Revision,
		}

	case stream_revision.WriteStreamRevisionNoStream:
		result.Options.ExpectedStreamRevision = &streams.DeleteReq_Options_NoStream{
			NoStream: &shared.Empty{},
		}
	case stream_revision.WriteStreamRevisionAny:
		result.Options.ExpectedStreamRevision = &streams.DeleteReq_Options_Any{
			Any: &shared.Empty{},
		}
	case stream_revision.WriteStreamRevisionStreamExists:
		result.Options.ExpectedStreamRevision = &streams.DeleteReq_Options_StreamExists{
			StreamExists: &shared.Empty{},
		}
	}

	return result
}
