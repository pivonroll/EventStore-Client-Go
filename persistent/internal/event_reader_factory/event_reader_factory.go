package event_reader_factory

//go:generate mockgen -source=event_reader_factory.go -destination=../mocks/event_reader_factory_mock.go -mock_names=Factory=EventReaderFactory -package=mocks

import (
	"context"

	"github.com/pivonroll/EventStore-Client-Go/persistent/event_reader"
	"github.com/pivonroll/EventStore-Client-Go/persistent/internal/message_adapter"
	"github.com/pivonroll/EventStore-Client-Go/protos/v21.6/persistent"
)

type Factory interface {
	Create(client persistent.PersistentSubscriptions_ReadClient,
		subscriptionId string,
		messageAdapter message_adapter.MessageAdapter,
		cancel context.CancelFunc,
	) event_reader.EventReader
}

type FactoryImpl struct{}

func (factory FactoryImpl) Create(
	client persistent.PersistentSubscriptions_ReadClient,
	subscriptionId string,
	messageAdapter message_adapter.MessageAdapter,
	cancel context.CancelFunc) event_reader.EventReader {

	return event_reader.NewEventReader(client, subscriptionId, messageAdapter, cancel)
}
