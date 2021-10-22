package persistent

import (
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

// UpdateSubscriptionGroupForStreamAllRequest is a set of data necessary to update a subscription to stream $all.
type UpdateSubscriptionGroupForStreamAllRequest struct {
	GroupName string                            // name of the persistent subscription group
	Position  stream_revision.IsReadPositionAll // position from which we want to start to receive events from a stream $all
	Settings  SubscriptionGroupSettings         // setting for a persistent subscription group
}

func (request UpdateSubscriptionGroupForStreamAllRequest) build() *persistent.UpdateReq {
	streamOption := &persistent.UpdateReq_Options_All{
		All: &persistent.UpdateReq_AllOptions{},
	}

	buildUpdateRequestPosition(request.Position, streamOption.All)
	protoSettings := request.Settings.buildUpdateRequestSettings()

	result := &persistent.UpdateReq{
		Options: &persistent.UpdateReq_Options{
			StreamOption: streamOption,
			GroupName:    request.GroupName,
			Settings:     protoSettings,
		},
	}

	return result
}

func buildUpdateRequestPosition(
	position stream_revision.IsReadPositionAll,
	protoOptions *persistent.UpdateReq_AllOptions) {
	switch position.(type) {
	case stream_revision.ReadPositionAllStart:
		protoOptions.AllOption = &persistent.UpdateReq_AllOptions_Start{
			Start: &shared.Empty{},
		}
	case stream_revision.ReadPositionAll:
		protoOptions.AllOption = &persistent.UpdateReq_AllOptions_Position{
			Position: &persistent.UpdateReq_Position{
				CommitPosition:  position.(stream_revision.ReadPositionAll).CommitPosition,
				PreparePosition: position.(stream_revision.ReadPositionAll).PreparePosition,
			},
		}
	case stream_revision.ReadPositionAllEnd:
		protoOptions.AllOption = &persistent.UpdateReq_AllOptions_End{
			End: &shared.Empty{},
		}
	}
}
