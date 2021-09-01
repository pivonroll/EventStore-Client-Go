// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package connection is a generated GoMock package.
package connection

import (
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockConnectionHandle is a mock of ConnectionHandle interface.
type MockConnectionHandle struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionHandleMockRecorder
}

// MockConnectionHandleMockRecorder is the mock recorder for MockConnectionHandle.
type MockConnectionHandleMockRecorder struct {
	mock *MockConnectionHandle
}

// NewMockConnectionHandle creates a new mock instance.
func NewMockConnectionHandle(ctrl *gomock.Controller) *MockConnectionHandle {
	mock := &MockConnectionHandle{ctrl: ctrl}
	mock.recorder = &MockConnectionHandleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnectionHandle) EXPECT() *MockConnectionHandleMockRecorder {
	return m.recorder
}

// Connection mocks base method.
func (m *MockConnectionHandle) Connection() *grpc.ClientConn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connection")
	ret0, _ := ret[0].(*grpc.ClientConn)
	return ret0
}

// Connection indicates an expected call of Connection.
func (mr *MockConnectionHandleMockRecorder) Connection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connection", reflect.TypeOf((*MockConnectionHandle)(nil).Connection))
}

// Id mocks base method.
func (m *MockConnectionHandle) Id() uuid.UUID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(uuid.UUID)
	return ret0
}

// Id indicates an expected call of Id.
func (mr *MockConnectionHandleMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockConnectionHandle)(nil).Id))
}

// MockGrpcClient is a mock of GrpcClient interface.
type MockGrpcClient struct {
	ctrl     *gomock.Controller
	recorder *MockGrpcClientMockRecorder
}

// MockGrpcClientMockRecorder is the mock recorder for MockGrpcClient.
type MockGrpcClientMockRecorder struct {
	mock *MockGrpcClient
}

// NewMockGrpcClient creates a new mock instance.
func NewMockGrpcClient(ctrl *gomock.Controller) *MockGrpcClient {
	mock := &MockGrpcClient{ctrl: ctrl}
	mock.recorder = &MockGrpcClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGrpcClient) EXPECT() *MockGrpcClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockGrpcClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockGrpcClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockGrpcClient)(nil).Close))
}

// GetConnectionHandle mocks base method.
func (m *MockGrpcClient) GetConnectionHandle() (ConnectionHandle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionHandle")
	ret0, _ := ret[0].(ConnectionHandle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionHandle indicates an expected call of GetConnectionHandle.
func (mr *MockGrpcClientMockRecorder) GetConnectionHandle() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionHandle", reflect.TypeOf((*MockGrpcClient)(nil).GetConnectionHandle))
}

// HandleError mocks base method.
func (m *MockGrpcClient) HandleError(handle ConnectionHandle, headers, trailers metadata.MD, err error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleError", handle, headers, trailers, err)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleError indicates an expected call of HandleError.
func (mr *MockGrpcClientMockRecorder) HandleError(handle, headers, trailers, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleError", reflect.TypeOf((*MockGrpcClient)(nil).HandleError), handle, headers, trailers, err)
}
