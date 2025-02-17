// Code generated by MockGen. DO NOT EDIT.
// Source: grpc_subscription_client_factory.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	persistent "github.com/pivonroll/EventStore-Client-Go/protos/persistent"
	grpc "google.golang.org/grpc"
)

// GrpcClientFactory is a mock of Factory interface.
type GrpcClientFactory struct {
	ctrl     *gomock.Controller
	recorder *GrpcClientFactoryMockRecorder
}

// GrpcClientFactoryMockRecorder is the mock recorder for GrpcClientFactory.
type GrpcClientFactoryMockRecorder struct {
	mock *GrpcClientFactory
}

// NewGrpcClientFactory creates a new mock instance.
func NewGrpcClientFactory(ctrl *gomock.Controller) *GrpcClientFactory {
	mock := &GrpcClientFactory{ctrl: ctrl}
	mock.recorder = &GrpcClientFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *GrpcClientFactory) EXPECT() *GrpcClientFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *GrpcClientFactory) Create(arg0 *grpc.ClientConn) persistent.PersistentSubscriptionsClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(persistent.PersistentSubscriptionsClient)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *GrpcClientFactoryMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*GrpcClientFactory)(nil).Create), arg0)
}
