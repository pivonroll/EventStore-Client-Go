package persistent

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/cenkalti/backoff/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_AsyncReadConnection_CallHandlerOnce(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	message := &messages.RecordedEvent{
		EventType: "some event type",
	}
	syncReadConnection := NewMockSyncReadConnection(ctrl)
	firstCall := syncReadConnection.EXPECT().Read().Return(message, nil) // first call should return valid message

	// we block the second call because we want to validate the message read from the first call
	syncReadConnection.EXPECT().Read().DoAndReturn(func() (*messages.RecordedEvent, error) {
		wg.Wait()
		return nil, nil
	}).After(firstCall)
	asyncConnection := newAsyncConnection(syncReadConnection, &backoff.ZeroBackOff{})

	asyncConnection.Start()                        // non-blocking call
	receivedMessage := <-asyncConnection.Updates() // blocks
	err := asyncConnection.Stop()                  // blocks
	require.NoError(t, err)
	wg.Done()
	require.Equal(t, *message, receivedMessage)
}

func Test_AsyncReadConnection_CallHandlerTwice(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	message := &messages.RecordedEvent{
		EventType: "some event type",
	}

	message2 := &messages.RecordedEvent{
		EventType: "some event type 2",
	}
	syncReadConnection := NewMockSyncReadConnection(ctrl)

	gomock.InOrder(
		syncReadConnection.EXPECT().Read().Return(message, nil),
		syncReadConnection.EXPECT().Read().Return(message2, nil),
		// we block the second call because we want to validate the message read from the first call
		syncReadConnection.EXPECT().Read().DoAndReturn(func() (*messages.RecordedEvent, error) {
			wg.Wait()
			return nil, nil
		}))

	asyncConnection := newAsyncConnection(syncReadConnection, &backoff.ZeroBackOff{})
	asyncConnection.Start() // non-blocking call
	messageChannel := asyncConnection.Updates()
	receivedMessage := <-messageChannel // blocks
	receivedMessage = <-messageChannel  // blocks

	err := asyncConnection.Stop() // blocks
	require.NoError(t, err)
	wg.Done()
	require.Equal(t, *message2, receivedMessage)
}

func Test_AsyncReadConnection_ErrorRead(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	message := &messages.RecordedEvent{
		EventType: "some event type",
	}

	message2 := &messages.RecordedEvent{
		EventType: "some event type 2",
	}
	syncReadConnection := NewMockSyncReadConnection(ctrl)

	gomock.InOrder(
		syncReadConnection.EXPECT().Read().Return(message, nil),
		syncReadConnection.EXPECT().Read().Return(nil, errors.New("some error")),
		syncReadConnection.EXPECT().Read().Return(message2, nil),
		syncReadConnection.EXPECT().Read().DoAndReturn(func() (*messages.RecordedEvent, error) {
			wg.Wait()
			return nil, nil
		}))

	asyncConnection := newAsyncConnection(syncReadConnection, &backoff.ConstantBackOff{
		Interval: time.Millisecond,
	})
	asyncConnection.Start() // non-blocking call
	messageChannel := asyncConnection.Updates()
	receivedMessage := <-messageChannel // blocks
	receivedMessage = <-messageChannel  // blocks

	err := asyncConnection.Stop() // blocks
	require.NoError(t, err)
	wg.Done()
	require.Equal(t, *message2, receivedMessage)
}

func Test_AsyncReadConnection_StartIsIdempotent(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	message := &messages.RecordedEvent{
		EventType: "some event type",
	}

	syncReadConnection := NewMockSyncReadConnection(ctrl)

	gomock.InOrder(
		syncReadConnection.EXPECT().Read().Return(message, nil),
		syncReadConnection.EXPECT().Read().DoAndReturn(func() (*messages.RecordedEvent, error) {
			wg.Wait()
			return nil, nil
		}))

	asyncConnection := newAsyncConnection(syncReadConnection, &backoff.ConstantBackOff{
		Interval: time.Millisecond,
	})
	asyncConnection.Start() // non-blocking call
	messageChannel := asyncConnection.Updates()
	asyncConnection.Start() // non-blocking call
	messageChannel2 := asyncConnection.Updates()

	require.Equal(t, messageChannel, messageChannel2)
	receivedMessage := <-messageChannel // blocks

	err := asyncConnection.Stop() // blocks
	require.NoError(t, err)
	wg.Done()
	require.Equal(t, *message, receivedMessage)
}

func Test_AsyncReadConnection_MessageChannelIsClosedAfterStop(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	message := &messages.RecordedEvent{
		EventType: "some event type",
	}

	syncReadConnection := NewMockSyncReadConnection(ctrl)

	gomock.InOrder(
		syncReadConnection.EXPECT().Read().Return(message, nil).AnyTimes(),
	)

	asyncConnection := newAsyncConnection(syncReadConnection, &backoff.ZeroBackOff{})
	asyncConnection.Start() // non-blocking call
	messageChannel := asyncConnection.Updates()
	asyncConnection.Start() // non-blocking call

	err := asyncConnection.Stop() // blocks
	require.NoError(t, err)

	_, open := <-messageChannel
	require.False(t, open)
}
