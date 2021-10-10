package persistent_test

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/pivonroll/EventStore-Client-Go/connection"
	"github.com/pivonroll/EventStore-Client-Go/event_streams"
	"github.com/pivonroll/EventStore-Client-Go/persistent"
	"github.com/pivonroll/EventStore-Client-Go/stream_revision"
)

func ExampleClient_CreateSubscriptionGroupForStream() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	persistentClient := persistent.NewClient(grpcClient)
	eventStreamsClient := event_streams.NewClient(grpcClient)

	streamID := "WithDefaultSettings"

	// create a stream (not necessary)
	testEvent := event_streams.ProposedEvent{
		EventId:      uuid.Must(uuid.NewRandom()),
		EventType:    "TestEvent",
		ContentType:  "application/octet-stream",
		UserMetadata: []byte(strings.Repeat("$", 5)),
		Data:         []byte{0xb, 0xe, 0xe, 0xf},
	}

	_, err = eventStreamsClient.AppendToStream(
		context.Background(),
		streamID,
		stream_revision.WriteStreamRevisionNoStream{},
		event_streams.ProposedEventList{testEvent})

	if err != nil {
		log.Fatalln("Events not appended to stream")
	}

	err = persistentClient.CreateSubscriptionGroupForStream(
		context.Background(),
		persistent.SubscriptionGroupForStreamRequest{
			StreamId:  streamID,
			GroupName: "Group WithDefaultSettings",
			Revision:  stream_revision.ReadStreamRevisionStart{},
			Settings:  persistent.DefaultRequestSettings,
		},
	)

	if err != nil {
		log.Fatalln("Failed to create persistent subscription group")
	}
}

func ExampleClient_UpdateSubscriptionGroupForStream() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	persistentClient := persistent.NewClient(grpcClient)
	eventStreamsClient := event_streams.NewClient(grpcClient)

	streamID := "WithDefaultSettings"

	// create a stream (not necessary)
	testEvent := event_streams.ProposedEvent{
		EventId:      uuid.Must(uuid.NewRandom()),
		EventType:    "TestEvent",
		ContentType:  "application/octet-stream",
		UserMetadata: []byte(strings.Repeat("$", 5)),
		Data:         []byte{0xb, 0xe, 0xe, 0xf},
	}

	_, err = eventStreamsClient.AppendToStream(
		context.Background(),
		streamID,
		stream_revision.WriteStreamRevisionNoStream{},
		event_streams.ProposedEventList{testEvent})

	if err != nil {
		log.Fatalln("Events not appended to stream")
	}

	streamConfig := persistent.SubscriptionGroupForStreamRequest{
		StreamId:  streamID,
		GroupName: "Group DefaultToNewSettings",
		Revision:  stream_revision.ReadStreamRevisionStart{},
		Settings:  persistent.DefaultRequestSettings,
	}

	err = persistentClient.CreateSubscriptionGroupForStream(
		context.Background(),
		streamConfig,
	)

	if err != nil {
		log.Fatalln("Failed to create persistent subscription group")
	}

	streamConfig.Settings.MaxSubscriberCount = 2
	err = persistentClient.UpdateSubscriptionGroupForStream(context.Background(), streamConfig)

	if err != nil {
		log.Fatalln("Failed to create persistent subscription group")
	}
}

func ExampleClient_DeleteSubscriptionGroupForStream() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	persistentClient := persistent.NewClient(grpcClient)
	eventStreamsClient := event_streams.NewClient(grpcClient)

	streamID := "WithDefaultSettings"

	// create a stream (not necessary)
	testEvent := event_streams.ProposedEvent{
		EventId:      uuid.Must(uuid.NewRandom()),
		EventType:    "TestEvent",
		ContentType:  "application/octet-stream",
		UserMetadata: []byte(strings.Repeat("$", 5)),
		Data:         []byte{0xb, 0xe, 0xe, 0xf},
	}

	_, err = eventStreamsClient.AppendToStream(
		context.Background(),
		streamID,
		stream_revision.WriteStreamRevisionNoStream{},
		event_streams.ProposedEventList{testEvent})

	if err != nil {
		log.Fatalln("Events not appended to stream")
	}

	err = persistentClient.CreateSubscriptionGroupForStream(
		context.Background(),
		persistent.SubscriptionGroupForStreamRequest{
			StreamId:  streamID,
			GroupName: "Group WithDefaultSettings",
			Revision:  stream_revision.ReadStreamRevisionStart{},
			Settings:  persistent.DefaultRequestSettings,
		},
	)

	if err != nil {
		log.Fatalln("Failed to create persistent subscription group")
	}

	err = persistentClient.DeleteSubscriptionGroupForStream(context.Background(), streamID, "Group WithDefaultSettings")

	if err != nil {
		log.Fatalln("Failed to create persistent subscription group")
	}
}

func ExampleClient_SubscribeToStreamSync() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	persistentClient := persistent.NewClient(grpcClient)
	eventStreamsClient := event_streams.NewClient(grpcClient)

	streamID := "some_stream"

	createEvents := func(eventCount uint32) event_streams.ProposedEventList {
		result := make(event_streams.ProposedEventList, eventCount)

		for i := uint32(0); i < eventCount; i++ {
			result = append(result, event_streams.ProposedEvent{
				EventId:      uuid.Must(uuid.NewRandom()),
				EventType:    "TestEvent",
				ContentType:  "application/octet-stream",
				UserMetadata: []byte{},
				Data:         []byte("some event data"),
			})
		}
		return result
	}

	events := createEvents(10)

	// append events to a stream
	_, err = eventStreamsClient.AppendToStream(context.Background(),
		streamID,
		stream_revision.WriteStreamRevisionNoStream{},
		events)

	if err != nil {
		log.Fatalln(err)
	}

	// create persistent stream subscription group to stream
	groupName := "Group StartFromBeginning_AndEventsInIt"
	request := persistent.SubscriptionGroupForStreamRequest{
		StreamId:  streamID,
		GroupName: groupName,
		Revision:  stream_revision.ReadStreamRevisionStart{},
		Settings:  persistent.DefaultRequestSettings,
	}
	err = persistentClient.CreateSubscriptionGroupForStream(
		context.Background(),
		request,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// subscribe to persistent subscription group
	readConnectionClient, err := persistentClient.
		SubscribeToStreamSync(context.Background(), 10, groupName, streamID)
	if err != nil {
		log.Fatalln(err)
	}

	// read one event
	readEvent, err := readConnectionClient.ReadOne()
	if err != nil {
		log.Fatalln(err)
	}

	if readEvent.GetOriginalEvent() == nil {
		log.Fatalln("Event not read")
	}

	fmt.Println(readEvent.GetOriginalEvent().EventNumber)
	fmt.Println(readEvent.GetOriginalEvent().EventId)
	fmt.Println(string(readEvent.GetOriginalEvent().Data))
}

func ExampleClient_CreateSubscriptionGroupForStreamAll() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	persistentClient := persistent.NewClient(grpcClient)

	err = persistentClient.CreateSubscriptionGroupForStreamAll(context.Background(),
		persistent.SubscriptionGroupForStreamAllRequest{
			GroupName: "some group",
			Position:  stream_revision.ReadPositionAllStart{},
			Filter: persistent.Filter{
				FilterBy: persistent.FilterByEventType{
					Matcher: persistent.PrefixFilterMatcher{
						PrefixList: []string{"event_type_prefix_1"},
					},
				},
				Window: persistent.FilterWindowMax{
					Max: 10, // maximum 10 events per read window
				},
				CheckpointIntervalMultiplier: 1,
			},
			Settings: persistent.DefaultRequestSettings,
		})

	if err != nil {
		log.Fatalln("Failed to create persistent subscription group to stream $all")
	}
}

func ExampleClient_UpdateSubscriptionGroupForStreamAll() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	persistentClient := persistent.NewClient(grpcClient)

	groupName := "some group"

	// Create a persistent subscription group
	err = persistentClient.CreateSubscriptionGroupForStreamAll(context.Background(),
		persistent.SubscriptionGroupForStreamAllRequest{
			GroupName: groupName,
			Position:  stream_revision.ReadPositionAllStart{},
			Filter: persistent.Filter{
				FilterBy: persistent.FilterByEventType{
					Matcher: persistent.PrefixFilterMatcher{
						PrefixList: []string{"event_type_prefix_1"},
					},
				},
				Window: persistent.FilterWindowMax{
					Max: 10, // maximum 10 events per read window
				},
				CheckpointIntervalMultiplier: 1,
			},
			Settings: persistent.DefaultRequestSettings,
		})

	if err != nil {
		log.Fatalln("Failed to create persistent subscription group to stream $all")
	}

	// Update a persistent subscription group
	err = persistentClient.UpdateSubscriptionGroupForStreamAll(context.Background(),
		persistent.UpdateSubscriptionGroupForStreamAllRequest{
			GroupName: groupName,
			Position:  stream_revision.ReadPositionAllEnd{},
			Settings:  persistent.DefaultRequestSettings,
		})

	if err != nil {
		log.Fatalln("Failed to update persistent subscription group")
	}
}

func ExampleClient_DeleteSubscriptionGroupForStreamAll() {
	username := "admin"
	password := "changeit"
	eventStoreEndpoint := "localhost:2113" // assuming that EventStoreDB is running on port 2113
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, eventStoreEndpoint)
	config, err := connection.ParseConnectionString(clientURI)
	if err != nil {
		log.Fatalln(err)
	}
	grpcClient := connection.NewGrpcClient(*config)
	persistentClient := persistent.NewClient(grpcClient)

	groupName := "some group"

	// Create a persistent subscription group
	err = persistentClient.CreateSubscriptionGroupForStreamAll(context.Background(),
		persistent.SubscriptionGroupForStreamAllRequest{
			GroupName: groupName,
			Position:  stream_revision.ReadPositionAllStart{},
			Filter: persistent.Filter{
				FilterBy: persistent.FilterByEventType{
					Matcher: persistent.PrefixFilterMatcher{
						PrefixList: []string{"event_type_prefix_1"},
					},
				},
				Window: persistent.FilterWindowMax{
					Max: 10, // maximum 10 events per read window
				},
				CheckpointIntervalMultiplier: 1,
			},
			Settings: persistent.DefaultRequestSettings,
		})

	if err != nil {
		log.Fatalln("Failed to create persistent subscription group to stream $all")
	}

	// Update a persistent subscription group
	err = persistentClient.DeleteSubscriptionGroupForStreamAll(context.Background(),
		groupName)

	if err != nil {
		log.Fatalln("Failed to delete persistent subscription group")
	}
}
