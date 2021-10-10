package message_adapter

//go:generate mockgen -source=message_adapter_provider.go -destination=../mocks/message_adapter_provider_mock.go -package=mocks

type MessageAdapterProvider interface {
	GetMessageAdapter() MessageAdapter
}

type MessageAdapterProviderImpl struct{}

func (provider MessageAdapterProviderImpl) GetMessageAdapter() MessageAdapter {
	return MessageAdapterImpl{}
}
