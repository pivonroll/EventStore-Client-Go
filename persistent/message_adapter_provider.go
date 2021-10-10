package persistent

//go:generate mockgen -source=message_adapter_provider.go -destination=message_adapter_provider_mock.go -package=persistent

type messageAdapterProvider interface {
	getMessageAdapter() messageAdapter
}

type messageAdapterProviderImpl struct{}

func (provider messageAdapterProviderImpl) getMessageAdapter() messageAdapter {
	return messageAdapterImpl{}
}
