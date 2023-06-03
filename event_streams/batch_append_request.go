package event_streams

import (
	"time"

	"github.com/google/uuid"
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/shared"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/streams"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type batchAppendRequest struct {
	correlationId    uuid.UUID
	options          batchAppendRequestOptions
	proposedMessages []batchAppendRequestProposedMessage
	isFinal          bool
}

func (this batchAppendRequest) build() *streams.BatchAppendReq {
	result := &streams.BatchAppendReq{
		CorrelationId: &shared.UUID{
			Value: &shared.UUID_String_{
				String_: this.correlationId.String(),
			},
		},
		Options: &streams.BatchAppendReq_Options{
			StreamIdentifier: &shared.StreamIdentifier{
				StreamName: []byte(this.options.streamId),
			},
			DeadlineOption: &streams.BatchAppendReq_Options_Deadline_21_10_0{
				Deadline_21_10_0: timestamppb.New(this.options.deadline),
			},
		},
		ProposedMessages: this.buildProposedMessages(),
		IsFinal:          this.isFinal,
	}

	this.buildExpectedStreamPosition(
		result,
		this.options.expectedStreamRevision)

	return result
}

func (this batchAppendRequest) buildExpectedStreamPosition(
	protoResult *streams.BatchAppendReq,
	position stream_revision.IsWriteStreamRevision,
) {
	switch position.(type) {
	case stream_revision.WriteStreamRevision:
		streamPosition := position.(stream_revision.WriteStreamRevision)
		protoResult.Options.ExpectedStreamPosition = &streams.BatchAppendReq_Options_StreamPosition{
			StreamPosition: streamPosition.Revision,
		}
	case stream_revision.WriteStreamRevisionNoStream:
		protoResult.Options.ExpectedStreamPosition = &streams.BatchAppendReq_Options_NoStream{
			NoStream: &emptypb.Empty{},
		}

	case stream_revision.WriteStreamRevisionAny:
		protoResult.Options.ExpectedStreamPosition = &streams.BatchAppendReq_Options_Any{
			Any: &emptypb.Empty{},
		}

	case stream_revision.WriteStreamRevisionStreamExists:
		protoResult.Options.ExpectedStreamPosition = &streams.BatchAppendReq_Options_StreamExists{
			StreamExists: &emptypb.Empty{},
		}
	}
}

func (this batchAppendRequest) buildProposedMessages() []*streams.BatchAppendReq_ProposedMessage {
	result := make([]*streams.BatchAppendReq_ProposedMessage, len(this.proposedMessages))

	for index, value := range this.proposedMessages {
		result[index] = &streams.BatchAppendReq_ProposedMessage{
			Id: &shared.UUID{
				Value: &shared.UUID_String_{
					String_: value.Id.String(),
				},
			},
			Metadata:       value.Metadata,
			CustomMetadata: value.CustomMetadata,
			Data:           value.Data,
		}
	}

	return result
}

type batchAppendRequestProposedMessage struct {
	Id             uuid.UUID
	Metadata       map[string]string
	CustomMetadata []byte
	Data           []byte
}

type batchAppendRequestOptions struct {
	streamId string
	// WriteStreamRevision
	// WriteStreamRevisionNoStream
	// WriteStreamRevisionAny
	// WriteStreamRevisionStreamExists
	expectedStreamRevision stream_revision.IsWriteStreamRevision
	deadline               time.Time
}
