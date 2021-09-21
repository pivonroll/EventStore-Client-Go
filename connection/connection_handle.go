package connection

import (
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
)

type connectionHandle struct {
	id         uuid.UUID
	connection *grpc.ClientConn
	err        error
}

func (handle connectionHandle) Id() uuid.UUID {
	return handle.id
}

func (handle connectionHandle) Connection() *grpc.ClientConn {
	return handle.connection
}

func newErroredConnectionHandle(err error) connectionHandle {
	return connectionHandle{
		id:         uuid.Nil,
		connection: nil,
		err:        err,
	}
}

func newConnectionHandle(id uuid.UUID, connection *grpc.ClientConn) connectionHandle {
	return connectionHandle{
		id:         id,
		connection: connection,
		err:        nil,
	}
}
