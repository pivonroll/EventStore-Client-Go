package event_streams

import (
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/shared"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/streams"
)

type subscribeToStreamRequest struct {
	streamOption isSubscribeRequestStreamOptions
	direction    subscribeRequestDirection
	resolveLinks bool
	// Filter
	// noFilter
	filter isFilter
}

func (this subscribeToStreamRequest) build() *streams.ReadReq {
	result := &streams.ReadReq{
		Options: &streams.ReadReq_Options{
			ResolveLinks: this.resolveLinks,
			FilterOption: nil,
			CountOption: &streams.ReadReq_Options_Subscription{
				Subscription: &streams.ReadReq_Options_SubscriptionOptions{},
			},
			UuidOption: &streams.ReadReq_Options_UUIDOption{
				Content: &streams.ReadReq_Options_UUIDOption_String_{
					String_: &shared.Empty{},
				},
			},
		},
	}

	if this.direction == subscribeRequestDirectionForward {
		result.Options.ReadDirection = streams.ReadReq_Options_Forwards
	} else {
		result.Options.ReadDirection = streams.ReadReq_Options_Backwards
	}

	this.buildStreamOption(result.Options)
	this.buildFilterOption(result.Options)

	return result
}

func (this subscribeToStreamRequest) buildStreamOption(options *streams.ReadReq_Options) {
	switch this.streamOption.(type) {
	case subscribeRequestStreamOptions:
		options.StreamOption = this.buildStreamOptions(
			this.streamOption.(subscribeRequestStreamOptions))
	case subscribeRequestStreamOptionsAll:
		options.StreamOption = this.buildStreamOptionAll(
			this.streamOption.(subscribeRequestStreamOptionsAll))
	}
}

func (this subscribeToStreamRequest) buildStreamOptionAll(all subscribeRequestStreamOptionsAll) *streams.ReadReq_Options_All {
	result := &streams.ReadReq_Options_All{
		All: &streams.ReadReq_Options_AllOptions{},
	}

	switch all.Position.(type) {
	case stream_revision.ReadPositionAll:
		allPosition := all.Position.(stream_revision.ReadPositionAll)
		result.All.AllOption = &streams.ReadReq_Options_AllOptions_Position{
			Position: &streams.ReadReq_Options_Position{
				CommitPosition:  allPosition.CommitPosition,
				PreparePosition: allPosition.PreparePosition,
			},
		}
	case stream_revision.ReadPositionAllStart:
		result.All.AllOption = &streams.ReadReq_Options_AllOptions_Start{
			Start: &shared.Empty{},
		}
	case stream_revision.ReadPositionAllEnd:
		result.All.AllOption = &streams.ReadReq_Options_AllOptions_End{
			End: &shared.Empty{},
		}
	}

	return result
}

func (this subscribeToStreamRequest) buildStreamOptions(
	streamOptions subscribeRequestStreamOptions,
) *streams.ReadReq_Options_Stream {
	result := &streams.ReadReq_Options_Stream{
		Stream: &streams.ReadReq_Options_StreamOptions{
			StreamIdentifier: &shared.StreamIdentifier{
				StreamName: []byte(streamOptions.StreamIdentifier),
			},
			RevisionOption: nil,
		},
	}

	switch streamOptions.Revision.(type) {
	case stream_revision.ReadStreamRevision:
		streamRevision := streamOptions.Revision.(stream_revision.ReadStreamRevision)
		result.Stream.RevisionOption = &streams.ReadReq_Options_StreamOptions_Revision{
			Revision: streamRevision.Revision,
		}
	case stream_revision.ReadStreamRevisionStart:
		result.Stream.RevisionOption = &streams.ReadReq_Options_StreamOptions_Start{
			Start: &shared.Empty{},
		}
	case stream_revision.ReadStreamRevisionEnd:
		result.Stream.RevisionOption = &streams.ReadReq_Options_StreamOptions_End{
			End: &shared.Empty{},
		}
	}

	return result
}

func (this subscribeToStreamRequest) buildFilterOption(options *streams.ReadReq_Options) {
	switch this.filter.(type) {
	case Filter:
		options.FilterOption = buildProtoFilter(this.filter)

	case noFilter:
		options.FilterOption = &streams.ReadReq_Options_NoFilter{NoFilter: &shared.Empty{}}
	}
}

type isSubscribeRequestStreamOptions interface {
	isSubscribeRequestStreamOptions()
}

type subscribeRequestStreamOptions struct {
	StreamIdentifier string
	// ReadStreamRevision
	// ReadStreamRevisionStart
	// ReadStreamRevisionEnd
	Revision stream_revision.IsReadStreamRevision
}

func (this subscribeRequestStreamOptions) isSubscribeRequestStreamOptions() {}

type subscribeRequestStreamOptionsAll struct {
	// ReadPositionAll
	// ReadPositionAllStart
	// ReadPositionAllEnd
	Position stream_revision.IsReadPositionAll
}

func (this subscribeRequestStreamOptionsAll) isSubscribeRequestStreamOptions() {}

type subscribeRequestDirection string

const (
	subscribeRequestDirectionForward subscribeRequestDirection = "subscribeRequestDirectionForward"
	subscribeRequestDirectionEnd     subscribeRequestDirection = "subscribeRequestDirectionEnd"
)
