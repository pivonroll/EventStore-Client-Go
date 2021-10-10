package persistent

import (
	"github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
	"github.com/pivonroll/EventStore-Client-Go/stream_revision"
)

type UpdateSubscriptionGroupForStreamAllRequest struct {
	GroupName string
	// AllPosition
	// AllPositionStart
	// AllPositionEnd
	Position stream_revision.IsReadPositionAll
	Settings CreateOrUpdateRequestSettings
}

func (request UpdateSubscriptionGroupForStreamAllRequest) Build() *persistent.UpdateReq {
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
