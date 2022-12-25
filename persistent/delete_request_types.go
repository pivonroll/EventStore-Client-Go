package persistent

import (
	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/persistent"
	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/shared"
)

type deleteSubscriptionGroupForStreamRequest struct {
	StreamId  string
	GroupName string
}

func (options deleteSubscriptionGroupForStreamRequest) build() *persistent.DeleteReq {
	return &persistent.DeleteReq{
		Options: &persistent.DeleteReq_Options{
			GroupName: options.GroupName,
			StreamOption: &persistent.DeleteReq_Options_StreamIdentifier{
				StreamIdentifier: &shared.StreamIdentifier{
					StreamName: []byte(options.StreamId),
				},
			},
		},
	}
}

func deleteRequestAllOptionsProto(groupName string) *persistent.DeleteReq {
	return &persistent.DeleteReq{
		Options: &persistent.DeleteReq_Options{
			GroupName: groupName,
			StreamOption: &persistent.DeleteReq_Options_All{
				All: &shared.Empty{},
			},
		},
	}
}
