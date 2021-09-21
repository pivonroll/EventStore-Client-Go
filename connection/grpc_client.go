package connection

import (
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ConnectionHandle interface {
	Id() uuid.UUID
	Connection() *grpc.ClientConn
}

type GrpcClient interface {
	HandleError(handle ConnectionHandle, headers metadata.MD, trailers metadata.MD, err error) error
	GetConnectionHandle() (ConnectionHandle, error)
	Close()
}
