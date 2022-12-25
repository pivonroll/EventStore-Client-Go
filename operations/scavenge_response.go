package operations

import "github.com/pivonroll/EventStore-Client-Go/protos/v21.6/operations"

// ScavengeResponse is a result received when scavenge operation is started and stopped.
type ScavengeResponse struct {
	ScavengeId string         // ID of the scavenge operations
	Status     ScavengeStatus // Status of the operation
}

// ScavengeStatus is a status of the scavenge operation.
type ScavengeStatus int32

const (
	// ScavengeStatus_Started scavenge operation is started
	ScavengeStatus_Started ScavengeStatus = 0
	// ScavengeStatus_InProgress scavenge operation is in progress
	ScavengeStatus_InProgress ScavengeStatus = 1
	// ScavengeStatus_Stopped scavenge operation is stopped
	ScavengeStatus_Stopped ScavengeStatus = 2
)

type scavengeResponseAdapter interface {
	create(protoResponse *operations.ScavengeResp) ScavengeResponse
}

type scavengeResponseAdapterImpl struct{}

func (adapter scavengeResponseAdapterImpl) create(
	protoResponse *operations.ScavengeResp) ScavengeResponse {
	result := ScavengeResponse{
		ScavengeId: protoResponse.ScavengeId,
		Status:     ScavengeStatus(protoResponse.ScavengeResult),
	}
	return result
}
