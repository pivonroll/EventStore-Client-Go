package event_streams_test

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/google/uuid"
	"github.com/pivonroll/EventStore-Client-Go/connection"
	"github.com/pivonroll/EventStore-Client-Go/errors"
	"github.com/pivonroll/EventStore-Client-Go/event_streams"
	"github.com/pivonroll/EventStore-Client-Go/stream_revision"
)

// Example of reading events from the start of a stream.
func ExampleClient_ReadStreamEvents_readEventsFromStart() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	client := event_streams.NewClient(grpcClient)

	streamId := "some_stream"
	proposedEvent := event_streams.ProposedEvent{
		EventId:      uuid.Must(uuid.NewRandom()),
		EventType:    "TestEvent",
		ContentType:  "application/octet-stream",
		UserMetadata: []byte{},
		Data:         []byte("some event data"),
	}

	eventsToWrite := []event_streams.ProposedEvent{proposedEvent}

	// create a stream with one event
	_, err = client.AppendToStream(context.Background(),
		streamId,
		stream_revision.WriteStreamRevisionNoStream{},
		eventsToWrite)
	if err != nil {
		log.Fatalln(err)
	}

	// read events from existing stream
	readEvents, err := client.ReadStreamEvents(context.Background(),
		streamId,
		event_streams.ReadDirectionForward,
		stream_revision.ReadStreamRevisionStart{},
		10, // set to be bigger than current number of events in a stream
		false)
	if err != nil {
		log.Fatalln(err)
	}

	// since readEvents are of type ResolvedEvent we must convert them to slice of ProposedEvents
	if !reflect.DeepEqual(eventsToWrite, readEvents.ToProposedEvents()) {
		log.Fatalln("Events read from a stream must match")
	}
}

// Example of reading events backwards from the end of a stream.
func ExampleClient_ReadStreamEvents_readEventsBackwardsFromEnd() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	client := event_streams.NewClient(grpcClient)

	streamId := "some_stream"

	// create 10 events to write (EventId must be unique)
	eventsToWrite := make(event_streams.ProposedEventList, 10)
	for i := uint32(0); i < 10; i++ {
		eventsToWrite[i] = event_streams.ProposedEvent{
			EventId:      uuid.Must(uuid.NewRandom()),
			EventType:    "TestEvent",
			ContentType:  "application/octet-stream",
			UserMetadata: []byte{},
			Data:         []byte("some event data"),
		}
	}

	// create a stream with 10 events
	_, err = client.AppendToStream(context.Background(),
		streamId,
		stream_revision.WriteStreamRevisionNoStream{},
		eventsToWrite)
	if err != nil {
		log.Fatalln(err)
	}

	// read events from existing stream
	readEvents, err := client.ReadStreamEvents(context.Background(),
		streamId,
		event_streams.ReadDirectionBackward,
		stream_revision.ReadStreamRevisionEnd{},
		10, // set to be bigger than current number of events in a stream
		false)
	if err != nil {
		log.Fatalln(err)
	}

	// Event read must be in reversed order
	// since readEvents are of type ResolvedEvent we must convert them to slice of ProposedEvents
	if !reflect.DeepEqual(eventsToWrite, readEvents.Reverse().ToProposedEvents()) {
		log.Fatalln("Events read from a stream must match")
	}
}

// Example of trying to read events from a stream which does not exist.
func ExampleClient_ReadStreamEvents_streamDoesNotExist() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, stdErr := connection.ParseConnectionString(clientURI)
	if stdErr != nil {
		log.Fatalln(stdErr)
	}
	grpcClient := connection.NewGrpcClient(*config)
	client := event_streams.NewClient(grpcClient)

	streamId := "some_stream"

	_, err := client.ReadStreamEvents(context.Background(),
		streamId,
		event_streams.ReadDirectionBackward,
		stream_revision.ReadStreamRevisionEnd{},
		1,
		false)

	if err.Code() != errors.StreamNotFoundErr {
		log.Fatalln("Stream must not exist")
	}
}

// Example of trying to read events from a stream which is soft-deleted.
func ExampleClient_ReadStreamEvents_streamIsSoftDeleted() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, stdErr := connection.ParseConnectionString(clientURI)
	if stdErr != nil {
		log.Fatalln(stdErr)
	}
	grpcClient := connection.NewGrpcClient(*config)
	client := event_streams.NewClient(grpcClient)

	streamId := "some_stream"

	proposedEvent := event_streams.ProposedEvent{
		EventId:      uuid.Must(uuid.NewRandom()),
		EventType:    "TestEvent",
		ContentType:  "application/octet-stream",
		UserMetadata: []byte{},
		Data:         []byte("some event data"),
	}

	// create a stream with one event
	writeResult, err := client.AppendToStream(context.Background(),
		streamId,
		stream_revision.WriteStreamRevisionNoStream{},
		[]event_streams.ProposedEvent{proposedEvent})
	if err != nil {
		log.Fatalln(err)
	}

	// soft-delete a stream
	_, err = client.DeleteStream(context.Background(),
		streamId,
		stream_revision.WriteStreamRevision{Revision: writeResult.GetCurrentRevision()})
	if err != nil {
		log.Fatalln(err)
	}

	// reading a soft-deleted stream fails with error code StreamNotFoundErr
	_, err = client.ReadStreamEvents(context.Background(),
		streamId,
		event_streams.ReadDirectionBackward,
		stream_revision.ReadStreamRevisionEnd{},
		event_streams.ReadCountMax,
		false)

	if err.Code() != errors.StreamNotFoundErr {
		log.Fatalln("Stream must not exist")
	}
}
