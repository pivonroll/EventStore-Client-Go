package event_streams

import (
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/shared"
	"github.com/pivonroll/EventStore-Client-Go/protos/v22.10/streams"
)

const ReadCountMax = ^uint64(0)

type readRequest struct {
	streamOption isReadRequestStreamOptions
	direction    ReadDirection
	resolveLinks bool
	count        uint64
	// Filter
	// noFilter
	filter isFilter
}

func (this readRequest) build() *streams.ReadReq {
	result := &streams.ReadReq{
		Options: &streams.ReadReq_Options{
			ResolveLinks: this.resolveLinks,
			FilterOption: nil,
			UuidOption: &streams.ReadReq_Options_UUIDOption{
				Content: &streams.ReadReq_Options_UUIDOption_String_{
					String_: &shared.Empty{},
				},
			},
			CountOption: &streams.ReadReq_Options_Count{
				Count: this.count,
			},
		},
	}

	if this.direction == ReadDirectionForward {
		result.Options.ReadDirection = streams.ReadReq_Options_Forwards
	} else {
		result.Options.ReadDirection = streams.ReadReq_Options_Backwards
	}

	this.buildStreamOption(result.Options)
	this.buildFilterOption(result.Options)

	return result
}

func (this readRequest) buildStreamOption(options *streams.ReadReq_Options) {
	switch this.streamOption.(type) {
	case readRequestStreamOptions:
		options.StreamOption = this.buildStreamOptions(
			this.streamOption.(readRequestStreamOptions))
	case readRequestStreamOptionsAll:
		options.StreamOption = thisBuildStreamOptionAll(
			this.streamOption.(readRequestStreamOptionsAll))
	}
}

func thisBuildStreamOptionAll(all readRequestStreamOptionsAll) *streams.ReadReq_Options_All {
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

func (this readRequest) buildStreamOptions(
	streamOptions readRequestStreamOptions,
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

func (this readRequest) buildFilterOption(options *streams.ReadReq_Options) {
	switch this.filter.(type) {
	case Filter:
		options.FilterOption = buildProtoFilter(this.filter)

	case noFilter:
		options.FilterOption = &streams.ReadReq_Options_NoFilter{NoFilter: &shared.Empty{}}
	}
}

type isReadRequestStreamOptions interface {
	isReadRequestStreamOptions()
}

type readRequestStreamOptions struct {
	StreamIdentifier string
	// ReadStreamRevision
	// ReadStreamRevisionStart
	// ReadStreamRevisionEnd
	Revision stream_revision.IsReadStreamRevision
}

func (this readRequestStreamOptions) isReadRequestStreamOptions() {}

type readRequestStreamOptionsAll struct {
	// ReadPositionAll
	// ReadPositionAllStart
	// ReadPositionAllEnd
	Position stream_revision.IsReadPositionAll
}

func (this readRequestStreamOptionsAll) isReadRequestStreamOptions() {}
