package client_test

import (
	"context"
	"testing"

	"github.com/EventStore/EventStore-Client-Go/errors"

	"github.com/EventStore/EventStore-Client-Go/event_streams"
	"github.com/EventStore/EventStore-Client-Go/systemmetadata"
	uuid "github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

type Created struct {
	Seconds int64 `json:"seconds"`
	Nanos   int   `json:"nanos"`
}

type StreamRevision struct {
	Value uint64 `json:"value"`
}

type TestEvent struct {
	Event Event `json:"event"`
}

type Event struct {
	StreamID       string         `json:"streamId"`
	StreamRevision StreamRevision `json:"streamRevision"`
	EventID        uuid.UUID      `json:"eventId"`
	EventType      string         `json:"eventType"`
	EventData      []byte         `json:"eventData"`
	UserMetadata   []byte         `json:"userMetadata"`
	ContentType    string         `json:"contentType"`
	Position       Position       `json:"position"`
	Created        Created        `json:"created"`
}

func Test_Read_Forwards_Linked_Stream_Big_Count(t *testing.T) {
	container := GetEmptyDatabase()
	defer container.Close()
	client := CreateTestClient(container, t)
	defer func() {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}()

	deletedStream := "deleted_stream"
	linkedStream := "linked_stream"

	_, err := client.AppendToStream(context.Background(),
		deletedStream,
		event_streams.AppendRequestExpectedStreamRevisionNoStream{},
		testCreateEvents(1))
	require.NoError(t, err)

	two := 2
	streamMetaData := event_streams.StreamMetadata{
		MaxCount: &two,
	}

	_, err = client.SetStreamMetadata(context.Background(),
		deletedStream,
		event_streams.AppendRequestExpectedStreamRevisionAny{},
		streamMetaData,
	)
	require.NoError(t, err)

	_, err = client.AppendToStream(context.Background(),
		deletedStream,
		event_streams.AppendRequestExpectedStreamRevisionNoStream{},
		testCreateEvents(1))
	require.NoError(t, err)

	_, err = client.AppendToStream(context.Background(),
		deletedStream,
		event_streams.AppendRequestExpectedStreamRevisionNoStream{},
		testCreateEvents(1))
	require.NoError(t, err)

	newUUID, _ := uuid.NewV4()
	linkedEvent := event_streams.ProposedEvent{
		EventID:      newUUID,
		EventType:    string(systemmetadata.SystemEventLinkTo),
		ContentType:  event_streams.ContentTypeOctetStream,
		Data:         []byte("0@" + deletedStream),
		UserMetadata: []byte{},
	}
	_, err = client.AppendToStream(context.Background(),
		linkedStream,
		event_streams.AppendRequestExpectedStreamRevisionNoStream{},
		[]event_streams.ProposedEvent{linkedEvent})
	require.NoError(t, err)

	readEvents, err := client.ReadStreamEvents(context.Background(),
		linkedStream,
		event_streams.ReadRequestDirectionForward,
		event_streams.ReadRequestOptionsStreamRevisionStart{},
		232323,
		true)
	require.NoError(t, err)
	require.Len(t, readEvents, 1)
}

func Test_Read_Forwards_Linked_Stream(t *testing.T) {
	container := GetEmptyDatabase()
	defer container.Close()
	client := CreateTestClient(container, t)
	defer func() {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}()

	deletedStream := "deleted_stream"
	linkedStream := "linked_stream"

	_, err := client.AppendToStream(context.Background(),
		deletedStream,
		event_streams.AppendRequestExpectedStreamRevisionAny{},
		testCreateEvents(1))
	require.NoError(t, err)

	newUUID, _ := uuid.NewV4()
	linkedEvent := event_streams.ProposedEvent{
		EventID:      newUUID,
		EventType:    string(systemmetadata.SystemEventLinkTo),
		ContentType:  event_streams.ContentTypeOctetStream,
		Data:         []byte("0@" + deletedStream),
		UserMetadata: []byte{},
	}
	_, err = client.AppendToStream(context.Background(),
		linkedStream,
		event_streams.AppendRequestExpectedStreamRevisionNoStream{},
		[]event_streams.ProposedEvent{linkedEvent})
	require.NoError(t, err)

	_, err = client.DeleteStream(
		context.Background(),
		deletedStream, event_streams.DeleteRequestExpectedStreamRevisionAny{})
	require.NoError(t, err)

	readEvents, err := client.ReadStreamEvents(context.Background(),
		linkedStream,
		event_streams.ReadRequestDirectionForward,
		event_streams.ReadRequestOptionsStreamRevisionStart{},
		1,
		true)

	require.NoError(t, err)

	t.Run("One Event Is Read", func(t *testing.T) {
		require.Len(t, readEvents, 1)
	})

	t.Run("Event Is Not Resolved", func(t *testing.T) {
		require.Nil(t, readEvents[0].Event)
	})

	t.Run("Link Event Is Included", func(t *testing.T) {
		require.NotNil(t, readEvents[0].Link)
	})
}

func Test_Read_Backwards_Linked_Stream(t *testing.T) {
	container := GetEmptyDatabase()
	defer container.Close()
	client := CreateTestClient(container, t)
	defer func() {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}()

	deletedStream := "deleted_stream"
	linkedStream := "linked_stream"

	_, err := client.AppendToStream(context.Background(),
		deletedStream,
		event_streams.AppendRequestExpectedStreamRevisionAny{},
		testCreateEvents(1))
	require.NoError(t, err)

	newUUID, _ := uuid.NewV4()
	linkedEvent := event_streams.ProposedEvent{
		EventID:      newUUID,
		EventType:    string(systemmetadata.SystemEventLinkTo),
		ContentType:  event_streams.ContentTypeOctetStream,
		Data:         []byte("0@" + deletedStream),
		UserMetadata: []byte{},
	}
	_, err = client.AppendToStream(context.Background(),
		linkedStream,
		event_streams.AppendRequestExpectedStreamRevisionNoStream{},
		[]event_streams.ProposedEvent{linkedEvent})
	require.NoError(t, err)

	_, err = client.DeleteStream(
		context.Background(),
		deletedStream, event_streams.DeleteRequestExpectedStreamRevisionAny{})
	require.NoError(t, err)

	readEvents, err := client.ReadStreamEvents(context.Background(),
		linkedStream,
		event_streams.ReadRequestDirectionBackward,
		event_streams.ReadRequestOptionsStreamRevisionStart{},
		1,
		true)

	require.NoError(t, err)

	t.Run("One Event Is Read", func(t *testing.T) {
		require.Len(t, readEvents, 1)
	})

	t.Run("Event Is Not Resolved", func(t *testing.T) {
		require.Nil(t, readEvents[0].Event)
	})

	t.Run("Link Event Is Included", func(t *testing.T) {
		require.NotNil(t, readEvents[0].Link)
	})
}

func Test_Read_Backwards(t *testing.T) {
	container := GetEmptyDatabase()
	defer container.Close()
	client := CreateTestClient(container, t)
	defer func() {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}()

	t.Run("Stream Does Not Exist Throws", func(t *testing.T) {
		streamId := "stream_does_not_exist_throws"
		_, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevisionEnd{},
			1,
			false)
		require.Equal(t, event_streams.StreamNotFoundErr, err.Code())
	})

	t.Run("Stream Deleted Throws", func(t *testing.T) {
		streamId := "stream_deleted_throws"

		_, err := client.TombstoneStream(context.Background(),
			streamId,
			event_streams.TombstoneRequestExpectedStreamRevisionNoStream{})
		require.NoError(t, err)

		_, err = client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevisionEnd{},
			1,
			false)
		require.Equal(t, errors.StreamDeletedErr, err.Code())
	})

	t.Run("Returns Events In Reversed Order Small Events", func(t *testing.T) {
		streamId := "returns_events_in_reversed_order_small_events"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevisionEnd{},
			10,
			false)

		readEventsToProposed := readEvents.Reverse().ToProposedEvents()
		require.Equal(t, events, readEventsToProposed)
	})

	t.Run("Read Single Event From Arbitrary Position", func(t *testing.T) {
		streamId := "read_single_event_from_arbitrary_position"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevision{Revision: 7},
			1,
			false)

		require.Equal(t, events[7], readEvents[0].ToProposedEvent())
	})

	t.Run("Read Many Events From Arbitrary Position", func(t *testing.T) {
		streamId := "read_many_events_from_arbitrary_position"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevision{Revision: 3},
			2,
			false)

		require.Equal(t, events[2:4], readEvents.Reverse().ToProposedEvents())
	})

	t.Run("Read First Event", func(t *testing.T) {
		streamId := "read_first_event"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			1,
			false)

		require.Equal(t, events[0], readEvents[0].ToProposedEvent())
	})

	t.Run("Read Last Event", func(t *testing.T) {
		streamId := "read_last_event"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevisionEnd{},
			1,
			false)

		require.Equal(t, events[9], readEvents[0].ToProposedEvent())
	})

	t.Run("Max Count Is Respected", func(t *testing.T) {
		streamId := "max_count_is_respected"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevisionEnd{},
			5,
			false)

		require.Len(t, readEvents, 5)
	})

	t.Run("Deadline Exceeded if Context Times Out", func(t *testing.T) {
		streamId := "deadline_exceeded"

		ctx := context.Background()
		timeoutCtx, cancelFunc := context.WithTimeout(ctx, 0)
		defer cancelFunc()

		_, err := client.ReadStreamEvents(timeoutCtx,
			streamId,
			event_streams.ReadRequestDirectionBackward,
			event_streams.ReadRequestOptionsStreamRevisionEnd{},
			1,
			false)

		require.Equal(t, errors.DeadlineExceededErr, err.Code())
	})
}

func Test_Read_Forwards(t *testing.T) {
	container := GetEmptyDatabase()
	defer container.Close()
	client := CreateTestClient(container, t)
	defer func() {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}()

	t.Run("Stream Does Not Exist Throws", func(t *testing.T) {
		streamId := "stream_does_not_exist_throws"
		_, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			1,
			false)
		require.Equal(t, event_streams.StreamNotFoundErr, err.Code())
	})

	t.Run("Stream Deleted Throws", func(t *testing.T) {
		streamId := "stream_deleted_throws"

		_, err := client.TombstoneStream(context.Background(),
			streamId,
			event_streams.TombstoneRequestExpectedStreamRevisionNoStream{})
		require.NoError(t, err)

		_, err = client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			1,
			false)
		require.Equal(t, errors.StreamDeletedErr, err.Code())
	})

	t.Run("Returns Events In Order Small Events", func(t *testing.T) {
		streamId := "returns_events_in_reversed_order_small_events"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			10,
			false)

		readEventsToProposed := readEvents.ToProposedEvents()
		require.Equal(t, events, readEventsToProposed)
	})

	t.Run("Returns Events In Order Large Events", func(t *testing.T) {
		streamId := "returns_events_in_reversed_order_small_events"
		events := testCreateEventsWithMetadata(2, 1_000_000)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			2,
			false)

		readEventsToProposed := readEvents.ToProposedEvents()
		require.Equal(t, events, readEventsToProposed)
	})

	t.Run("Read Single Event From Arbitrary Position", func(t *testing.T) {
		streamId := "read_single_event_from_arbitrary_position"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevision{Revision: 7},
			1,
			false)

		require.Equal(t, events[7], readEvents[0].ToProposedEvent())
	})

	t.Run("Read Many Events From Arbitrary Position", func(t *testing.T) {
		streamId := "read_many_events_from_arbitrary_position"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevision{Revision: 3},
			2,
			false)

		require.Equal(t, events[3:5], readEvents.ToProposedEvents())
	})

	t.Run("Read First Event", func(t *testing.T) {
		streamId := "read_first_event"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			1,
			false)

		require.Equal(t, events[0], readEvents[0].ToProposedEvent())
	})

	t.Run("Max Count Is Respected", func(t *testing.T) {
		streamId := "max_count_is_respected"
		events := testCreateEvents(10)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			5,
			false)

		require.Len(t, readEvents, 5)
	})

	t.Run("Read All Events By Default", func(t *testing.T) {
		streamId := "reads_all_events_by_default"
		events := testCreateEvents(200)

		_, err := client.AppendToStream(context.Background(),
			streamId,
			event_streams.AppendRequestExpectedStreamRevisionNoStream{},
			events)
		require.NoError(t, err)

		readEvents, err := client.ReadStreamEvents(context.Background(),
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			event_streams.ReadCountMax,
			false)

		require.Len(t, readEvents, 200)
	})

	t.Run("Deadline Exceeded if Context Times Out", func(t *testing.T) {
		streamId := "deadline_exceeded"

		ctx := context.Background()
		timeoutCtx, cancelFunc := context.WithTimeout(ctx, 0)
		defer cancelFunc()

		_, err := client.ReadStreamEvents(timeoutCtx,
			streamId,
			event_streams.ReadRequestDirectionForward,
			event_streams.ReadRequestOptionsStreamRevisionStart{},
			1,
			false)

		require.Equal(t, errors.DeadlineExceededErr, err.Code())
	})
}

//func TestReadStreamEventsForwardsFromZeroPosition(t *testing.T) {
//	eventsContent, err := ioutil.ReadFile("../resources/test/dataset20M-1800-e0-e10.json")
//	require.NoError(t, err)
//
//	var testEvents []TestEvent
//	err = json.Unmarshal(eventsContent, &testEvents)
//	require.NoError(t, err)
//
//	container := GetPrePopulatedDatabase()
//	defer container.Close()
//
//	client := CreateTestClient(container, t)
//	defer client.Close()
//
//	context, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
//	defer cancel()
//
//	numberOfEventsToRead := 10
//	numberOfEvents := uint64(numberOfEventsToRead)
//
//	streamId := "dataset20M-1800"
//
//	stream, err := client.ReadStreamEvents_OLD(context, direction.Forwards, streamId, stream_position.Start{}, numberOfEvents, true)
//	if err != nil {
//		t.Fatalf("Unexpected failure %+v", err)
//	}
//
//	defer stream.Close()
//
//	events, err := collectStreamEvents(stream)
//	if err != nil {
//		t.Fatalf("Unexpected failure %+v", err)
//	}
//
//	assert.Equal(t, numberOfEvents, uint64(len(events)), "Expected the correct number of messages to be returned")
//
//	for i := 0; i < numberOfEventsToRead; i++ {
//		assert.Equal(t, testEvents[i].Event.EventID, events[i].GetOriginalEvent().EventID)
//		assert.Equal(t, testEvents[i].Event.EventType, events[i].GetOriginalEvent().EventType)
//		assert.Equal(t, testEvents[i].Event.StreamID, events[i].GetOriginalEvent().StreamID)
//		assert.Equal(t, testEvents[i].Event.StreamRevision.Value, events[i].GetOriginalEvent().EventNumber)
//		assert.Equal(t, testEvents[i].Event.Created.Nanos, events[i].GetOriginalEvent().CreatedDate.Nanosecond())
//		assert.Equal(t, testEvents[i].Event.Created.Seconds, events[i].GetOriginalEvent().CreatedDate.Unix())
//		assert.Equal(t, testEvents[i].Event.Position.Commit, events[i].GetOriginalEvent().Position.Commit)
//		assert.Equal(t, testEvents[i].Event.Position.Prepare, events[i].GetOriginalEvent().Position.Prepare)
//		assert.Equal(t, testEvents[i].Event.ContentType, events[i].GetOriginalEvent().ContentType)
//	}
//}

//func TestReadStreamEventsBackwardsFromEndPosition(t *testing.T) {
//	eventsContent, err := ioutil.ReadFile("../resources/test/dataset20M-1800-e1999-e1990.json")
//	require.NoError(t, err)
//
//	var testEvents []TestEvent
//	err = json.Unmarshal(eventsContent, &testEvents)
//
//	require.NoError(t, err)
//	container := GetPrePopulatedDatabase()
//	defer container.Close()
//
//	client := CreateTestClient(container, t)
//	defer client.Close()
//
//	context, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
//	defer cancel()
//
//	numberOfEventsToRead := 10
//	numberOfEvents := uint64(numberOfEventsToRead)
//
//	streamId := "dataset20M-1800"
//
//	stream, err := client.ReadStreamEvents_OLD(context, direction.Backwards, streamId, stream_position.End{}, numberOfEvents, true)
//	if err != nil {
//		t.Fatalf("Unexpected failure %+v", err)
//	}
//
//	defer stream.Close()
//
//	events, err := collectStreamEvents(stream)
//	if err != nil {
//		t.Fatalf("Unexpected failure %+v", err)
//	}
//
//	assert.Equal(t, numberOfEvents, uint64(len(events)), "Expected the correct number of messages to be returned")
//
//	for i := 0; i < numberOfEventsToRead; i++ {
//		assert.Equal(t, testEvents[i].Event.EventID, events[i].GetOriginalEvent().EventID)
//		assert.Equal(t, testEvents[i].Event.EventType, events[i].GetOriginalEvent().EventType)
//		assert.Equal(t, testEvents[i].Event.StreamID, events[i].GetOriginalEvent().StreamID)
//		assert.Equal(t, testEvents[i].Event.StreamRevision.Value, events[i].GetOriginalEvent().EventNumber)
//		assert.Equal(t, testEvents[i].Event.Created.Nanos, events[i].GetOriginalEvent().CreatedDate.Nanosecond())
//		assert.Equal(t, testEvents[i].Event.Created.Seconds, events[i].GetOriginalEvent().CreatedDate.Unix())
//		assert.Equal(t, testEvents[i].Event.Position.Commit, events[i].GetOriginalEvent().Position.Commit)
//		assert.Equal(t, testEvents[i].Event.Position.Prepare, events[i].GetOriginalEvent().Position.Prepare)
//		assert.Equal(t, testEvents[i].Event.ContentType, events[i].GetOriginalEvent().ContentType)
//	}
//}

//func TestReadStreamReturnsEOFAfterCompletion(t *testing.T) {
//	container := GetEmptyDatabase()
//	defer container.Close()
//	client := CreateTestClient(container, t)
//	defer client.Close()
//
//	var waitingForError sync.WaitGroup
//
//	proposedEvents := []messages.ProposedEvent{}
//
//	for i := 1; i <= 10; i++ {
//		proposedEvents = append(proposedEvents, createTestEvent())
//	}
//
//	_, err := client.AppendToStream_OLD(context.Background(), "testing-closing", streamrevision.StreamRevisionNoStream, proposedEvents)
//	require.NoError(t, err)
//
//	stream, err := client.ReadStreamEvents_OLD(context.Background(), direction.Forwards, "testing-closing", stream_position.Start{}, 1_024, false)
//
//	require.NoError(t, err)
//	_, err = collectStreamEvents(stream)
//	require.NoError(t, err)
//
//	go func() {
//		_, err := stream.Recv()
//		require.Error(t, err)
//		require.True(t, err == io.EOF)
//		waitingForError.Done()
//	}()
//
//	require.NoError(t, err)
//	waitingForError.Add(1)
//	timedOut := waitWithTimeout(&waitingForError, time.Duration(5)*time.Second)
//	require.False(t, timedOut, "Timed out waiting for read stream to return io.EOF on completion")
//}
