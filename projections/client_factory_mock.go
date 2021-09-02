// Code generated by MockGen. DO NOT EDIT.
// Source: client_factory.go

// Package projections is a generated GoMock package.
package projections

import (
	reflect "reflect"

	connection "github.com/EventStore/EventStore-Client-Go/connection"
	projections "github.com/EventStore/EventStore-Client-Go/protos/projections"
	gomock "github.com/golang/mock/gomock"
)

// MockClientFactory is a mock of ClientFactory interface.
type MockClientFactory struct {
	ctrl     *gomock.Controller
	recorder *MockClientFactoryMockRecorder
}

// MockClientFactoryMockRecorder is the mock recorder for MockClientFactory.
type MockClientFactoryMockRecorder struct {
	mock *MockClientFactory
}

// NewMockClientFactory creates a new mock instance.
func NewMockClientFactory(ctrl *gomock.Controller) *MockClientFactory {
	mock := &MockClientFactory{ctrl: ctrl}
	mock.recorder = &MockClientFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientFactory) EXPECT() *MockClientFactoryMockRecorder {
	return m.recorder
}

// CreateClient mocks base method.
func (m *MockClientFactory) CreateClient(grpcClient connection.GrpcClient, projectionsClient projections.ProjectionsClient) Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClient", grpcClient, projectionsClient)
	ret0, _ := ret[0].(Client)
	return ret0
}

// CreateClient indicates an expected call of CreateClient.
func (mr *MockClientFactoryMockRecorder) CreateClient(grpcClient, projectionsClient interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClient", reflect.TypeOf((*MockClientFactory)(nil).CreateClient), grpcClient, projectionsClient)
}
