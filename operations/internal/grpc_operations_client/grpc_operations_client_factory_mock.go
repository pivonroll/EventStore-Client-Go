// Code generated by MockGen. DO NOT EDIT.
// Source: grpc_operations_client_factory.go

// Package grpc_operations_client is a generated GoMock package.
package grpc_operations_client

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	operations "github.com/pivonroll/EventStore-Client-Go/protos/operations"
	grpc "google.golang.org/grpc"
)

// MockFactory is a mock of Factory interface.
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory.
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance.
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockFactory) Create(arg0 *grpc.ClientConn) operations.OperationsClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(operations.OperationsClient)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockFactoryMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockFactory)(nil).Create), arg0)
}
