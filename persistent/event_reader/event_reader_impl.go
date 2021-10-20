package event_reader

import (
	"context"
	"io"
	"sync"

	"github.com/google/uuid"
	"github.com/pivonroll/EventStore-Client-Go/connection"
	"github.com/pivonroll/EventStore-Client-Go/errors"
	"github.com/pivonroll/EventStore-Client-Go/persistent/internal/message_adapter"
	"github.com/pivonroll/EventStore-Client-Go/persistent/persistent_action"
	"github.com/pivonroll/EventStore-Client-Go/persistent/persistent_event"
	"github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	"github.com/pivonroll/EventStore-Client-Go/protos/shared"
)

// MAX_ACK_COUNT is maximum number of messages which can be acknowledged.
const MAX_ACK_COUNT = 2000

type eventReaderImpl struct {
	protoClient        persistent.PersistentSubscriptions_ReadClient
	subscriptionId     string
	messageAdapter     message_adapter.MessageAdapter
	readRequestChannel chan chan readResponse
	cancel             context.CancelFunc
	once               sync.Once
}

// ReadOne reads one message from persistent subscription.
// Reads will block one the buffer size for persistent subscription is reached.
// Buffer size is specified when we are connecting to a persistent subscription.
func (reader *eventReaderImpl) ReadOne() (persistent_event.ReadResponseEvent, errors.Error) {
	channel := make(chan readResponse)

	reader.readRequestChannel <- channel
	resp := <-channel

	return resp.ReadResponseEvent, resp.Error
}

func (reader *eventReaderImpl) readOne() (persistent_event.ReadResponseEvent, errors.Error) {
	protoResponse, protoErr := reader.protoClient.Recv()
	if protoErr != nil {
		if protoErr == io.EOF {
			return persistent_event.ReadResponseEvent{}, errors.NewError(errors.EndOfStream, protoErr)
		}
		trailer := reader.protoClient.Trailer()
		err := connection.GetErrorFromProtoException(trailer, protoErr)
		if err != nil {
			return persistent_event.ReadResponseEvent{}, err
		}
		return persistent_event.ReadResponseEvent{}, errors.NewError(errors.FatalError, protoErr)
	}

	result := reader.messageAdapter.FromProtoResponse(protoResponse.GetEvent())
	return result, nil
}

// Close closes the connection to a persistent subscription.
func (reader *eventReaderImpl) Close() {
	reader.once.Do(reader.cancel)
}

const Exceeds_Max_Message_Count_Err errors.ErrorCode = "Exceeds_Max_Message_Count_Err"

// Ack sends Ack signal for a message.
// Maximum of 2000 messages can be acknowledged at once.
// If number of messages exceeds the maximum error Exceeds_Max_Message_Count_Err will be returned.
func (reader *eventReaderImpl) Ack(messages ...persistent_event.ReadResponseEvent) errors.Error {
	if len(messages) == 0 {
		return nil
	}

	if len(messages) > MAX_ACK_COUNT {
		return errors.NewErrorCode(Exceeds_Max_Message_Count_Err)
	}

	var ids []uuid.UUID
	for _, event := range messages {
		ids = append(ids, event.GetOriginalEvent().EventId)
	}

	protoErr := reader.protoClient.Send(&persistent.ReadReq{
		Content: &persistent.ReadReq_Ack_{
			Ack: &persistent.ReadReq_Ack{
				Id:  []byte(reader.subscriptionId),
				Ids: messageIdSliceToProto(ids...),
			},
		},
	})

	if protoErr != nil {
		trailers := reader.protoClient.Trailer()
		err := connection.GetErrorFromProtoException(trailers, protoErr)
		if err != nil {
			return err
		}

		return errors.NewErrorCode(errors.FatalError)
	}

	return nil
}

// Nack sends Nack signal for a message.
// Client must also specify a reason why message was nack-ed.
func (reader *eventReaderImpl) Nack(reason string,
	action persistent_action.Nack_Action, messages ...persistent_event.ReadResponseEvent) error {
	if len(messages) == 0 {
		return nil
	}

	ids := []uuid.UUID{}
	for _, event := range messages {
		ids = append(ids, event.GetOriginalEvent().EventId)
	}

	err := reader.protoClient.Send(&persistent.ReadReq{
		Content: &persistent.ReadReq_Nack_{
			Nack: &persistent.ReadReq_Nack{
				Id:     []byte(reader.subscriptionId),
				Ids:    messageIdSliceToProto(ids...),
				Action: persistent.ReadReq_Nack_Action(action),
				Reason: reason,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (reader *eventReaderImpl) readLoopWithRequest() {
	for {
		responseChannel := <-reader.readRequestChannel
		result, err := reader.readOne()

		responseChannel <- readResponse{
			ReadResponseEvent: result,
			Error:             err,
		}
	}
}

func messageIdSliceToProto(messageIds ...uuid.UUID) []*shared.UUID {
	result := make([]*shared.UUID, len(messageIds))

	for index, messageId := range messageIds {
		result[index] = toProtoUUID(messageId)
	}

	return result
}

type readResponse struct {
	ReadResponseEvent persistent_event.ReadResponseEvent
	Error             errors.Error
}

// NewEventReader creates a new event reader
func NewEventReader(
	client persistent.PersistentSubscriptions_ReadClient,
	subscriptionId string,
	messageAdapter message_adapter.MessageAdapter,
	cancel context.CancelFunc,
) EventReader {
	channel := make(chan chan readResponse)

	reader := &eventReaderImpl{
		protoClient:        client,
		subscriptionId:     subscriptionId,
		messageAdapter:     messageAdapter,
		readRequestChannel: channel,
		cancel:             cancel,
	}

	// It is not safe to consume a stream in different goroutines. This is why we only consume
	// the stream in a dedicated goroutine.
	//
	// Current implementation doesn't terminate the goroutine. When a subscription is dropped,
	// we keep user requests coming but will always send back a subscription dropped event.
	// This implementation is simple to maintain while letting the user sharing their subscription
	// among as many goroutines as they want.
	go reader.readLoopWithRequest()

	return reader
}
