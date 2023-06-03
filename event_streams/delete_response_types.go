package event_streams

import "github.com/pivonroll/EventStore-Client-Go/protos/v22.10/streams"

// DeleteResponse is response received when stream is soft-deleted by using Client.DeleteStream.
type DeleteResponse struct {
	position isDeleteResponsePosition
}

// GetPosition returns a position at which stream was soft-deleted.
// If position was received it will also return a true as a second return value.
// If position does not exist a zero initialized Position and a false will be returned.
// Position may not exist if an empty stream was soft-deleted.
func (response DeleteResponse) GetPosition() (Position, bool) {
	if response.position != nil {
		if position, isPosition := response.position.(deleteResponsePosition); isPosition {
			return Position{
				CommitPosition:  position.CommitPosition,
				PreparePosition: position.PreparePosition,
			}, true
		}
	}

	return Position{}, false
}

type isDeleteResponsePosition interface {
	isDeleteResponsePositionOption()
}

type deleteResponsePosition struct {
	CommitPosition  uint64
	PreparePosition uint64
}

func (this deleteResponsePosition) isDeleteResponsePositionOption() {
}

type deleteResponseAdapter interface {
	Create(resp *streams.DeleteResp) DeleteResponse
}

type deleteResponseAdapterImpl struct{}

func (this deleteResponseAdapterImpl) Create(protoResponse *streams.DeleteResp) DeleteResponse {
	result := DeleteResponse{}

	switch protoResponse.PositionOption.(type) {
	case *streams.DeleteResp_Position_:
		protoPosition := protoResponse.PositionOption.(*streams.DeleteResp_Position_)
		result.position = deleteResponsePosition{
			CommitPosition:  protoPosition.Position.CommitPosition,
			PreparePosition: protoPosition.Position.PreparePosition,
		}
	case *streams.DeleteResp_NoPosition:
		result.position = nil
	}

	return result
}
