package message_adapter

import (
	"log"
	"strconv"
	"time"

	"github.com/pivonroll/EventStore-Client-Go/persistent/persistent_event"
	"github.com/pivonroll/EventStore-Client-Go/position"
	"github.com/pivonroll/EventStore-Client-Go/protobuf_uuid"
	"github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	"github.com/pivonroll/EventStore-Client-Go/ptr"
	system_metadata "github.com/pivonroll/EventStore-Client-Go/systemmetadata"
)

type MessageAdapter interface {
	FromProtoResponse(resp *persistent.ReadResp_ReadEvent) persistent_event.ReadResponseEvent
}

type MessageAdapterImpl struct{}

func (adapter MessageAdapterImpl) FromProtoResponse(
	readEvent *persistent.ReadResp_ReadEvent) persistent_event.ReadResponseEvent {
	positionWire := readEvent.GetPosition()
	eventWire := readEvent.GetEvent()
	linkWire := readEvent.GetLink()
	retryCount := readEvent.GetCount()

	result := persistent_event.ReadResponseEvent{
		Event:          nil,
		Link:           nil,
		CommitPosition: nil,
		RetryCount:     nil,
	}

	if eventWire != nil {
		event := newRecordedEventFromProto(eventWire)
		result.Event = &event
	}

	if linkWire != nil {
		link := newRecordedEventFromProto(linkWire)
		result.Link = &link
	}

	if retryCount != nil {
		result.RetryCount = ptr.Int32(readEvent.GetRetryCount())
	}

	if positionWire != nil {
		result.CommitPosition = ptr.UInt64(readEvent.GetCommitPosition())
	}

	return result
}

func newRecordedEventFromProto(
	recordedEvent *persistent.ReadResp_ReadEvent_RecordedEvent) persistent_event.RecordedEvent {
	streamIdentifier := recordedEvent.GetStreamIdentifier()

	return persistent_event.RecordedEvent{
		EventId:        protobuf_uuid.GetUUID(recordedEvent.GetId()),
		EventType:      recordedEvent.Metadata[system_metadata.SystemMetadataKeysType],
		ContentType:    getContentTypeFromProto(recordedEvent),
		StreamId:       string(streamIdentifier.StreamName),
		EventNumber:    recordedEvent.GetStreamRevision(),
		CreatedDate:    createdFromProto(recordedEvent),
		Position:       positionFromProto(recordedEvent),
		Data:           recordedEvent.GetData(),
		SystemMetadata: recordedEvent.GetMetadata(),
		UserMetadata:   recordedEvent.GetCustomMetadata(),
	}
}

func getContentTypeFromProto(recordedEvent *persistent.ReadResp_ReadEvent_RecordedEvent) string {
	return recordedEvent.Metadata[system_metadata.SystemMetadataKeysContentType]
}

func createdFromProto(recordedEvent *persistent.ReadResp_ReadEvent_RecordedEvent) time.Time {
	timeSinceEpoch, err := strconv.ParseInt(
		recordedEvent.Metadata[system_metadata.SystemMetadataKeysCreated], 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse created date as int from %+v",
			recordedEvent.Metadata[system_metadata.SystemMetadataKeysCreated])
	}
	// The metadata contains the number of .NET "ticks" (100ns increments) since the UNIX epoch
	return time.Unix(0, timeSinceEpoch*100).UTC()
}

func positionFromProto(recordedEvent *persistent.ReadResp_ReadEvent_RecordedEvent) position.Position {
	return position.Position{
		Commit:  recordedEvent.GetCommitPosition(),
		Prepare: recordedEvent.GetPreparePosition(),
	}
}