package connection

import (
	"fmt"

	"github.com/gofrs/uuid"
)

type getConnection struct {
	channel chan connectionHandle
}

func newGetConnectionMsg() getConnection {
	return getConnection{
		channel: make(chan connectionHandle),
	}
}

func (msg getConnection) handle(state connectionState) connectionState {
	if state.IsConnected() { // Means we are connected
		handle := connectionHandle{
			id:         state.correlation,
			connection: state.connection,
			err:        nil,
		}

		msg.channel <- handle
		return state
	}

	// Means we need to create a grpc connection.
	newConnection, err := discoverNode(state.config)
	if err != nil {
		state.lastError = err
		resp := newErroredConnectionHandle(err)
		msg.channel <- resp
		return state
	}

	id, err := uuid.NewV4()
	if err != nil {
		state.lastError = fmt.Errorf("error when trying to generate a random UUID: %v", err)
		return state
	}

	state.correlation = id
	state.connection = newConnection

	resp := newConnectionHandle(id, newConnection)
	msg.channel <- resp
	return state
}
