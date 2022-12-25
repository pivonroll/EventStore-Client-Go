package event_streams

import (
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/shared"
	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/streams2"
)

type tombstoneRequest struct {
	streamId string
	// Types that are assignable to expectedStreamRevision:
	// WriteStreamRevision
	// WriteStreamRevisionNoStream
	// WriteStreamRevisionAny
	// WriteStreamRevisionStreamExists
	expectedStreamRevision stream_revision.IsWriteStreamRevision
}

func (this tombstoneRequest) build() *streams2.TombstoneReq {
	result := &streams2.TombstoneReq{
		Options: &streams2.TombstoneReq_Options{
			StreamIdentifier: &shared.StreamIdentifier{
				StreamName: []byte(this.streamId),
			},
			ExpectedStreamRevision: nil,
		},
	}

	switch this.expectedStreamRevision.(type) {
	case stream_revision.WriteStreamRevision:
		revision := this.expectedStreamRevision.(stream_revision.WriteStreamRevision)
		result.Options.ExpectedStreamRevision = &streams2.TombstoneReq_Options_Revision{
			Revision: revision.Revision,
		}
	case stream_revision.WriteStreamRevisionNoStream:
		result.Options.ExpectedStreamRevision = &streams2.TombstoneReq_Options_NoStream{
			NoStream: &shared.Empty{},
		}
	case stream_revision.WriteStreamRevisionAny:
		result.Options.ExpectedStreamRevision = &streams2.TombstoneReq_Options_Any{
			Any: &shared.Empty{},
		}
	case stream_revision.WriteStreamRevisionStreamExists:
		result.Options.ExpectedStreamRevision = &streams2.TombstoneReq_Options_StreamExists{
			StreamExists: &shared.Empty{},
		}
	}

	return result
}
