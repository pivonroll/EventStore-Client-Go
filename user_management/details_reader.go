package user_management

import (
	"context"
	"io"
	"sync"

	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/pivonroll/EventStore-Client-Go/core/errors"
	"github.com/pivonroll/EventStore-Client-Go/protos/users"
)

type detailsReader interface {
	Recv() (DetailsResponse, errors.Error)
}

type DetailsReaderImpl struct {
	protoStreamReader      users.Users_DetailsClient
	detailsResponseAdapter detailsResponseAdapter
	once                   sync.Once
	cancelFunc             context.CancelFunc
}

func (reader *DetailsReaderImpl) Recv() (DetailsResponse, errors.Error) {
	protoResponse, protoErr := reader.protoStreamReader.Recv()
	if protoErr != nil {
		if protoErr == io.EOF {
			return DetailsResponse{}, errors.NewError(errors.EndOfStream, protoErr)
		}
		trailer := reader.protoStreamReader.Trailer()
		err := connection.GetErrorFromProtoException(trailer, protoErr)
		if err != nil {
			return DetailsResponse{}, err
		}
		return DetailsResponse{}, errors.NewError(errors.FatalError, protoErr)
	}

	return reader.detailsResponseAdapter.Create(protoResponse), nil
}

func (this *DetailsReaderImpl) Close() {
	this.once.Do(this.cancelFunc)
}
