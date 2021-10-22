package event_streams_test

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/core/systemmetadata"
	"github.com/pivonroll/EventStore-Client-Go/event_streams"
)

// Example demonstrates how to subscribe to stream $all without a filter.
//
// We create two streams and write events to them.
// Subscription to stream $all must catch all events written to those two streams.
func ExampleClient_SubscribeToStreamAll() {
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

	createEvents := func(count uint32) event_streams.ProposedEventList {
		result := make(event_streams.ProposedEventList, 10)
		for i := uint32(0); i < count; i++ {
			result[i] = event_streams.ProposedEvent{
				EventId:      uuid.Must(uuid.NewRandom()),
				EventType:    "TestEvent",
				ContentType:  "application/octet-stream",
				UserMetadata: []byte{},
				Data:         []byte("some event data"),
			}
		}
		return result
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	firstStream := "firstStream"
	secondStream := "secondStream"

	beforeEvents := createEvents(10)
	afterEvents := createEvents(10)

	allUserEvents := append(beforeEvents, afterEvents...)

	// create a stream with some events in it
	_, err := client.AppendToStream(context.Background(),
		firstStream,
		stream_revision.WriteStreamRevisionNoStream{},
		beforeEvents)
	if err != nil {
		log.Fatalln(err)
	}

	// create a subscription to a stream
	streamReader, err := client.SubscribeToStreamAll(context.Background(),
		stream_revision.ReadPositionAllStart{},
		false)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer wg.Done()

		var resultsRead event_streams.ProposedEventList
		readResult, err := streamReader.ReadOne()
		if err != nil {
			log.Fatalln(err)
		}

		if event, isEvent := readResult.GetEvent(); isEvent {
			if !systemmetadata.IsSystemStream(event.Event.StreamId) {
				resultsRead = append(resultsRead, event.ToProposedEvent())
			}
		}

		// if we have read all user defined event stop listening for events and return from go routine
		if reflect.DeepEqual(allUserEvents, resultsRead) {
			streamReader.Close()
			return
		}
	}()

	// append some events to a stream after a listening go routine has started
	_, err = client.AppendToStream(context.Background(),
		secondStream,
		stream_revision.WriteStreamRevisionNoStream{},
		afterEvents)
	if err != nil {
		log.Fatalln(err)
	}

	// wait for subscription to receive all events
	wg.Wait()
}
