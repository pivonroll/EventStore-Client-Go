package persistent_integration_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pivonroll/EventStore-Client-Go/core/errors"
	"github.com/pivonroll/EventStore-Client-Go/core/stream_revision"
	"github.com/pivonroll/EventStore-Client-Go/core/test_utils"
	"github.com/pivonroll/EventStore-Client-Go/persistent"
	"github.com/stretchr/testify/require"
)

func Test_CreatePersistentStreamSubscription(t *testing.T) {
	client, eventStreamsClient, closeFunc := initializeContainerAndClient(t, nil)
	defer closeFunc()

	t.Run("With Default Settings", func(t *testing.T) {
		streamID := "WithDefaultSettings"
		pushEventToStream(t, eventStreamsClient, streamID)

		err := client.CreateSubscriptionGroupForStream(
			context.Background(),
			persistent.SubscriptionGroupForStreamRequest{
				StreamId:  streamID,
				GroupName: "Group WithDefaultSettings",
				Revision:  stream_revision.ReadStreamRevisionStart{},
				Settings:  persistent.DefaultRequestSettings,
			},
		)

		require.NoError(t, err)
	})

	t.Run("MessageTimeoutZero", func(t *testing.T) {
		streamID := "MessageTimeoutZero"
		pushEventToStream(t, eventStreamsClient, streamID)

		settings := persistent.DefaultRequestSettings
		settings.MessageTimeout = persistent.MessageTimeoutInMs{MilliSeconds: 0}

		err := client.CreateSubscriptionGroupForStream(
			context.Background(),
			persistent.SubscriptionGroupForStreamRequest{
				StreamId:  streamID,
				GroupName: "Group MessageTimeoutZero",
				Revision:  stream_revision.ReadStreamRevisionStart{},
				Settings:  settings,
			},
		)
		require.NoError(t, err)
	})

	t.Run("StreamNotExits", func(t *testing.T) {
		streamID := "StreamNotExits"

		err := client.CreateSubscriptionGroupForStream(
			context.Background(),
			persistent.SubscriptionGroupForStreamRequest{
				StreamId:  streamID,
				GroupName: "Group StreamNotExits",
				Revision:  stream_revision.ReadStreamRevisionStart{},
				Settings:  persistent.DefaultRequestSettings,
			},
		)

		require.NoError(t, err)
	})

	t.Run("FailsIfAlreadyExists", func(t *testing.T) {
		streamID := "FailsIfAlreadyExists"
		pushEventToStream(t, eventStreamsClient, streamID)

		err := client.CreateSubscriptionGroupForStream(
			context.Background(),
			persistent.SubscriptionGroupForStreamRequest{
				StreamId:  streamID,
				GroupName: "Group FailsIfAlreadyExists",
				Revision:  stream_revision.ReadStreamRevisionStart{},
				Settings:  persistent.DefaultRequestSettings,
			},
		)

		require.NoError(t, err)

		err = client.CreateSubscriptionGroupForStream(
			context.Background(),
			persistent.SubscriptionGroupForStreamRequest{
				StreamId:  streamID,
				GroupName: "Group FailsIfAlreadyExists",
				Revision:  stream_revision.ReadStreamRevisionStart{},
				Settings:  persistent.DefaultRequestSettings,
			},
		)

		require.Error(t, err)
	})

	t.Run("AfterDeleting", func(t *testing.T) {
		streamID := "AfterDeleting"
		pushEventToStream(t, eventStreamsClient, streamID)

		streamConfig := persistent.SubscriptionGroupForStreamRequest{
			StreamId:  streamID,
			GroupName: "Group AfterDeleting",
			Revision:  stream_revision.ReadStreamRevisionStart{},
			Settings:  persistent.DefaultRequestSettings,
		}

		err := client.CreateSubscriptionGroupForStream(context.Background(), streamConfig)

		require.NoError(t, err)

		err = client.DeleteSubscriptionGroupForStream(context.Background(),
			streamID,
			streamConfig.GroupName,
		)

		require.NoError(t, err)

		err = client.CreateSubscriptionGroupForStream(context.Background(), streamConfig)

		require.NoError(t, err)
	})
}

func Test_UpdatePersistentStreamSubscription(t *testing.T) {
	client, eventStreamsClient, closeFunc := initializeContainerAndClient(t, nil)
	defer closeFunc()

	t.Run("Default To New Settings", func(t *testing.T) {
		streamID := "DefaultToNewSettings"
		pushEventToStream(t, eventStreamsClient, streamID)

		streamConfig := persistent.SubscriptionGroupForStreamRequest{
			StreamId:  streamID,
			GroupName: "Group DefaultToNewSettings",
			Revision:  stream_revision.ReadStreamRevisionStart{},
			Settings:  persistent.DefaultRequestSettings,
		}

		err := client.CreateSubscriptionGroupForStream(context.Background(), streamConfig)

		require.NoError(t, err)

		streamConfig.Settings.HistoryBufferSize = streamConfig.Settings.HistoryBufferSize + 1
		streamConfig.Settings.NamedConsumerStrategy = persistent.ConsumerStrategy_DispatchToSingle
		streamConfig.Settings.MaxSubscriberCount = streamConfig.Settings.MaxSubscriberCount + 1
		streamConfig.Settings.ReadBatchSize = streamConfig.Settings.ReadBatchSize + 1
		streamConfig.Settings.CheckpointAfter = persistent.CheckpointAfterMs{MilliSeconds: 112}
		streamConfig.Settings.MaxCheckpointCount = streamConfig.Settings.MaxCheckpointCount + 1
		streamConfig.Settings.MinCheckpointCount = streamConfig.Settings.MinCheckpointCount + 1
		streamConfig.Settings.LiveBufferSize = streamConfig.Settings.LiveBufferSize + 1
		streamConfig.Settings.MaxRetryCount = streamConfig.Settings.MaxRetryCount + 1
		streamConfig.Settings.MessageTimeout = persistent.MessageTimeoutInMs{MilliSeconds: 222}
		streamConfig.Settings.ExtraStatistics = !streamConfig.Settings.ExtraStatistics
		streamConfig.Settings.ResolveLinks = !streamConfig.Settings.ResolveLinks

		err = client.UpdateSubscriptionGroupForStream(context.Background(), streamConfig)

		require.NoError(t, err)
	})

	t.Run("Error If Subscription Does Not Exist", func(t *testing.T) {
		streamID := "ErrorIfSubscriptionDoesNotExist"

		streamConfig := persistent.SubscriptionGroupForStreamRequest{
			StreamId:  streamID,
			GroupName: "Group ErrorIfSubscriptionDoesNotExist",
			Revision:  stream_revision.ReadStreamRevisionStart{},
			Settings:  persistent.DefaultRequestSettings,
		}

		err := client.UpdateSubscriptionGroupForStream(context.Background(), streamConfig)

		require.Error(t, err)
	})
}

func Test_DeletePersistentStreamSubscription(t *testing.T) {
	client, eventStreamsClient, closeFunc := initializeContainerAndClient(t, nil)
	defer closeFunc()

	t.Run("Delete Existing", func(t *testing.T) {
		streamID := "DeleteExisting"
		pushEventToStream(t, eventStreamsClient, streamID)

		streamConfig := persistent.SubscriptionGroupForStreamRequest{
			StreamId:  streamID,
			GroupName: "Group DeleteExisting",
			Revision:  stream_revision.ReadStreamRevisionStart{},
			Settings:  persistent.DefaultRequestSettings,
		}

		err := client.CreateSubscriptionGroupForStream(context.Background(), streamConfig)

		require.NoError(t, err)

		err = client.DeleteSubscriptionGroupForStream(context.Background(),
			streamID,
			streamConfig.GroupName,
		)

		require.NoError(t, err)
	})

	t.Run("Error If Subscription Does Not Exist", func(t *testing.T) {
		err := client.DeleteSubscriptionGroupForStream(context.Background(),
			"a",
			"a",
		)

		require.Error(t, err)
	})
}

func TestPersistentSubscriptionClosing(t *testing.T) {
	client, closeFunc := initializeWithPrePopulatedDatabase(t)
	defer closeFunc()

	streamID := "dataset20M-0"
	groupName := "Group 1"
	bufferSize := int32(2)

	streamConfig := persistent.SubscriptionGroupForStreamRequest{
		StreamId:  streamID,
		GroupName: "Group 1",
		Revision:  stream_revision.ReadStreamRevisionStart{},
		Settings:  persistent.DefaultRequestSettings,
	}

	err := client.CreateSubscriptionGroupForStream(context.Background(), streamConfig)

	require.NoError(t, err)

	receivedEvents := sync.WaitGroup{}
	droppedEvent := sync.WaitGroup{}
	waitForClose := sync.WaitGroup{}
	waitForClose.Add(1)
	receivedEvents.Add(10)
	droppedEvent.Add(1)

	subscription, err := client.SubscribeToStreamSync(
		context.Background(), bufferSize, groupName, streamID)

	require.NoError(t, err)

	go func() {
		current := 1

		for {
			if current == 11 {
				waitForClose.Wait()
			}
			result, err := subscription.ReadOne()

			if err != nil && err.Code() == errors.CanceledErr {
				droppedEvent.Done()
				break
			}

			if current <= 10 {
				receivedEvents.Done()
				current++
			}

			err = subscription.Ack(result)
			require.NoError(t, err)

			continue
		}
	}()

	timedOut := test_utils.WaitWithTimeout(&receivedEvents, time.Duration(5)*time.Second)
	require.False(t, timedOut, "Timed out waiting for initial set of events")
	subscription.Close()
	waitForClose.Done()
	timedOut = test_utils.WaitWithTimeout(&droppedEvent, time.Duration(5)*time.Second)
	require.False(t, timedOut, "Timed out waiting for dropped event")
}
