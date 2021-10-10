package persistent

import (
	"github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
	"github.com/pivonroll/EventStore-Client-Go/stream_revision"
)

type SubscriptionGroupForStreamRequest struct {
	StreamId  string
	GroupName string
	Revision  stream_revision.IsReadStreamRevision
	Settings  CreateOrUpdateRequestSettings
}

func (request SubscriptionGroupForStreamRequest) BuildCreateStreamRequest() *persistent.CreateReq {
	streamOption := &persistent.CreateReq_StreamOptions{
		StreamIdentifier: &shared.StreamIdentifier{StreamName: []byte(request.StreamId)},
		RevisionOption:   nil,
	}

	buildCreateRequestRevision(request.Revision, streamOption)

	result := &persistent.CreateReq{
		Options: &persistent.CreateReq_Options{
			StreamOption: &persistent.CreateReq_Options_Stream{
				Stream: streamOption,
			},
			StreamIdentifier: &shared.StreamIdentifier{StreamName: []byte(request.StreamId)},
			GroupName:        request.GroupName,
			Settings:         request.Settings.buildCreateRequestSettings(),
		},
	}

	return result
}

func buildCreateRequestRevision(
	revision stream_revision.IsReadStreamRevision,
	streamOptions *persistent.CreateReq_StreamOptions) {
	switch revision.(type) {
	case stream_revision.ReadStreamRevision:
		streamOptions.RevisionOption = &persistent.CreateReq_StreamOptions_Revision{
			Revision: revision.(stream_revision.ReadStreamRevision).Revision,
		}
	case stream_revision.ReadStreamRevisionStart:
		streamOptions.RevisionOption = &persistent.CreateReq_StreamOptions_Start{
			Start: &shared.Empty{},
		}
	case stream_revision.ReadStreamRevisionEnd:
		streamOptions.RevisionOption = &persistent.CreateReq_StreamOptions_End{
			End: &shared.Empty{},
		}
	}
}

func (request SubscriptionGroupForStreamRequest) BuildUpdateStreamRequest() *persistent.UpdateReq {
	streamOption := &persistent.UpdateReq_StreamOptions{
		StreamIdentifier: &shared.StreamIdentifier{StreamName: []byte(request.StreamId)},
		RevisionOption:   nil,
	}

	buildUpdateRequestRevision(request.Revision, streamOption)

	result := &persistent.UpdateReq{
		Options: &persistent.UpdateReq_Options{
			StreamOption: &persistent.UpdateReq_Options_Stream{
				Stream: streamOption,
			},
			StreamIdentifier: &shared.StreamIdentifier{StreamName: []byte(request.StreamId)},
			GroupName:        request.GroupName,
			Settings:         request.Settings.buildUpdateRequestSettings(),
		},
	}

	return result
}

func buildUpdateRequestRevision(
	revision stream_revision.IsReadStreamRevision,
	streamOptions *persistent.UpdateReq_StreamOptions) {
	switch revision.(type) {
	case stream_revision.ReadStreamRevisionStart:
		streamOptions.RevisionOption = &persistent.UpdateReq_StreamOptions_Start{
			Start: &shared.Empty{},
		}
	case stream_revision.ReadStreamRevision:
		streamOptions.RevisionOption = &persistent.UpdateReq_StreamOptions_Revision{
			Revision: revision.(stream_revision.ReadStreamRevision).Revision,
		}
	case stream_revision.ReadStreamRevisionEnd:
		streamOptions.RevisionOption = &persistent.UpdateReq_StreamOptions_End{
			End: &shared.Empty{},
		}
	}
}
