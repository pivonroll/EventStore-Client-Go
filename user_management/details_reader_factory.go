package user_management

import (
	"context"
	"sync"

	"github.com/pivonroll/EventStore-Client-Go/protos/users"
)

type detailsReaderFactory interface {
	Create(protoStreamReader users.Users_DetailsClient,
		cancelFunc context.CancelFunc) detailsReader
}

type detailsReaderFactoryImpl struct{}

func (factory detailsReaderFactoryImpl) Create(
	protoStreamReader users.Users_DetailsClient,
	cancelFunc context.CancelFunc) detailsReader {
	return &DetailsReaderImpl{
		protoStreamReader:      protoStreamReader,
		detailsResponseAdapter: detailsResponseAdapterImpl{},
		once:                   sync.Once{},
		cancelFunc:             cancelFunc,
	}
}
