// Code generated by MockGen. DO NOT EDIT.
// Source: grpc_subscription_client_factory.go

// Package persistent is a generated GoMock package.
package persistent

import (
	reflect "reflect"

	persistent "github.com/EventStore/EventStore-Client-Go/protos/persistent"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockgrpcSubscriptionClientFactory is a mock of grpcSubscriptionClientFactory interface.
type MockgrpcSubscriptionClientFactory struct {
	ctrl     *gomock.Controller
	recorder *MockgrpcSubscriptionClientFactoryMockRecorder
}

// MockgrpcSubscriptionClientFactoryMockRecorder is the mock recorder for MockgrpcSubscriptionClientFactory.
type MockgrpcSubscriptionClientFactoryMockRecorder struct {
	mock *MockgrpcSubscriptionClientFactory
}

// NewMockgrpcSubscriptionClientFactory creates a new mock instance.
func NewMockgrpcSubscriptionClientFactory(ctrl *gomock.Controller) *MockgrpcSubscriptionClientFactory {
	mock := &MockgrpcSubscriptionClientFactory{ctrl: ctrl}
	mock.recorder = &MockgrpcSubscriptionClientFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgrpcSubscriptionClientFactory) EXPECT() *MockgrpcSubscriptionClientFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockgrpcSubscriptionClientFactory) Create(arg0 *grpc.ClientConn) persistent.PersistentSubscriptionsClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(persistent.PersistentSubscriptionsClient)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockgrpcSubscriptionClientFactoryMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockgrpcSubscriptionClientFactory)(nil).Create), arg0)
}
