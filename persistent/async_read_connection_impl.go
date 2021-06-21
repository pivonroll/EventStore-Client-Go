package persistent

import (
	"sync"
	"time"

	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/protos/persistent"
	"github.com/cenkalti/backoff/v3"
)

type asyncConnectionImpl struct {
	readConnection     SyncReadConnection
	stopRequestChannel chan chan error             // for Close
	updatesChannel     chan messages.RecordedEvent // sends items to the user
	backOff            backoff.BackOff
	once               sync.Once
}

func newAsyncConnectionImpl(
	client persistent.PersistentSubscriptions_ReadClient,
	subscriptionId string,
) AsyncReadConnection {
	return newAsyncConnection(
		newSyncReadConnection(client, subscriptionId, messageAdapterImpl{}),
		&backoff.ConstantBackOff{
			Interval: 1 * time.Second,
		},
	)
}

func newAsyncConnection(
	readConnection SyncReadConnection,
	backOff backoff.BackOff) AsyncReadConnection {
	return &asyncConnectionImpl{
		readConnection:     readConnection,
		updatesChannel:     make(chan messages.RecordedEvent), // for Updates
		stopRequestChannel: make(chan chan error),             // for Close
		backOff:            backOff,
		once:               sync.Once{},
	}
}

func (asyncConnection *asyncConnectionImpl) Start() {
	asyncConnection.once.Do(func() {
		go asyncConnection.loop()
	})
}

// Stop blocks until context cancel is acknowledged
func (asyncConnection *asyncConnectionImpl) Stop() error {
	stopResponseChannel := make(chan error)
	asyncConnection.stopRequestChannel <- stopResponseChannel
	return <-stopResponseChannel
}

func (asyncConnection *asyncConnectionImpl) Updates() <-chan messages.RecordedEvent {
	return asyncConnection.updatesChannel
}

const maxCacheSize = 10

type messageCache struct {
	messages []messages.RecordedEvent
}

func (cache *messageCache) isFull() bool {
	return len(cache.messages) >= maxCacheSize
}

func (cache *messageCache) hasMessages() bool {
	return len(cache.messages) > 0
}

func (cache *messageCache) firstMessage() messages.RecordedEvent {
	return cache.messages[0]
}

func (cache *messageCache) addMessage(newMessage messages.RecordedEvent) {
	cache.messages = append(cache.messages, newMessage)
}

func (cache *messageCache) popFirstMessage() {
	cache.messages = cache.messages[1:]
}

func (asyncConnection *asyncConnectionImpl) loop() {
	type fetchResult struct {
		fetchedMessage *messages.RecordedEvent
		err            error
	}

	var fetchDone chan fetchResult // if non-nil, Fetch is running
	messageCache := messageCache{}
	var err error
	var fetchDelay time.Duration
	for {
		// Enable new Fetch if current fetching is not in progress and cache is not full
		var doFetch <-chan time.Time
		if fetchDone == nil && !messageCache.isFull() {
			doFetch = time.After(fetchDelay) // enable fetch case
		}

		var nextMessageForConsumer messages.RecordedEvent
		var updatesChannel chan messages.RecordedEvent // block send case with nil channel
		if messageCache.hasMessages() {
			nextMessageForConsumer = messageCache.firstMessage()
			updatesChannel = asyncConnection.updatesChannel // enable send case
		}

		select {
		case <-doFetch:
			fetchDone = make(chan fetchResult, 1)
			go func() {
				fetched, err := asyncConnection.readConnection.Read()
				fetchDone <- fetchResult{fetchedMessage: fetched, err: err}
			}()
		case result := <-fetchDone:
			fetchDone = nil // when fetch is done block listening on fetchDone channel

			fetchedMessage := result.fetchedMessage
			err = result.err
			if err != nil { // if error was received from call to Read initiate backoff mechanism
				fetchDelay = asyncConnection.backOff.NextBackOff()
				break
			}

			asyncConnection.backOff.Reset()
			fetchDelay = 0
			if fetchedMessage != nil {
				messageCache.addMessage(*fetchedMessage)
			}

		case stopResponseChannel := <-asyncConnection.stopRequestChannel: // if Stop is initiated
			stopResponseChannel <- err
			close(asyncConnection.updatesChannel)
			return
		case updatesChannel <- nextMessageForConsumer:
			messageCache.popFirstMessage()
		}
	}
}
