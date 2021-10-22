package persistent

import (
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

// SubscriptionGroupForStreamAllRequest is a struct with all data necessary to create a persistent subscription group.
type SubscriptionGroupForStreamAllRequest struct {
	GroupName string                            // name of the persistent subscription group
	Position  stream_revision.IsReadPositionAll // position from which we want to start to receive events from a stream $all
	// CreateRequestAllNoFilter
	// CreateRequestAllFilter
	Filter   isFilter                  // filter for messages from stream $all
	Settings SubscriptionGroupSettings // setting for a persistent subscription group
}

func (request SubscriptionGroupForStreamAllRequest) build() *persistent.CreateReq {
	streamOption := &persistent.CreateReq_Options_All{
		All: &persistent.CreateReq_AllOptions{},
	}

	switch request.Filter.(type) {
	case Filter:
		streamOption.All.FilterOption = buildProtoFilter(request.Filter.(Filter))
	case NoFilter:
		streamOption.All.FilterOption = buildProtoNoFilter()
	}

	buildCreateRequestPosition(request.Position, streamOption.All)

	protoSettings := request.Settings.buildCreateRequestSettings()

	result := &persistent.CreateReq{
		Options: &persistent.CreateReq_Options{
			StreamOption: streamOption,
			GroupName:    request.GroupName,
			Settings:     protoSettings,
		},
	}

	return result
}

func buildCreateRequestPosition(
	position stream_revision.IsReadPositionAll,
	protoOptions *persistent.CreateReq_AllOptions) {
	switch position.(type) {
	case stream_revision.ReadPositionAllStart:
		protoOptions.AllOption = &persistent.CreateReq_AllOptions_Start{
			Start: &shared.Empty{},
		}
	case stream_revision.ReadPositionAll:
		protoOptions.AllOption = &persistent.CreateReq_AllOptions_Position{
			Position: &persistent.CreateReq_Position{
				CommitPosition:  position.(stream_revision.ReadPositionAll).CommitPosition,
				PreparePosition: position.(stream_revision.ReadPositionAll).PreparePosition,
			},
		}
	case stream_revision.ReadPositionAllEnd:
		protoOptions.AllOption = &persistent.CreateReq_AllOptions_End{
			End: &shared.Empty{},
		}
	}
}
