// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ledgerwatch/erigon/polygon/p2p (interfaces: Service)
//
// Generated by this command:
//
//	mockgen -destination=./service_mock.go -package=p2p . Service
//

// Package p2p is a generated GoMock package.
package p2p

import (
	context "context"
	reflect "reflect"

	types "github.com/ledgerwatch/erigon/core/types"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// FetchHeaders mocks base method.
func (m *MockService) FetchHeaders(arg0 context.Context, arg1, arg2 uint64, arg3 *PeerId) ([]*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchHeaders", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchHeaders indicates an expected call of FetchHeaders.
func (mr *MockServiceMockRecorder) FetchHeaders(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchHeaders", reflect.TypeOf((*MockService)(nil).FetchHeaders), arg0, arg1, arg2, arg3)
}

// GetMessageListener mocks base method.
func (m *MockService) GetMessageListener() MessageListener {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageListener")
	ret0, _ := ret[0].(MessageListener)
	return ret0
}

// GetMessageListener indicates an expected call of GetMessageListener.
func (mr *MockServiceMockRecorder) GetMessageListener() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageListener", reflect.TypeOf((*MockService)(nil).GetMessageListener))
}

// ListPeersMayHaveBlockNum mocks base method.
func (m *MockService) ListPeersMayHaveBlockNum(arg0 uint64) []*PeerId {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPeersMayHaveBlockNum", arg0)
	ret0, _ := ret[0].([]*PeerId)
	return ret0
}

// ListPeersMayHaveBlockNum indicates an expected call of ListPeersMayHaveBlockNum.
func (mr *MockServiceMockRecorder) ListPeersMayHaveBlockNum(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPeersMayHaveBlockNum", reflect.TypeOf((*MockService)(nil).ListPeersMayHaveBlockNum), arg0)
}

// MaxPeers mocks base method.
func (m *MockService) MaxPeers() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MaxPeers")
	ret0, _ := ret[0].(int)
	return ret0
}

// MaxPeers indicates an expected call of MaxPeers.
func (mr *MockServiceMockRecorder) MaxPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MaxPeers", reflect.TypeOf((*MockService)(nil).MaxPeers))
}

// Penalize mocks base method.
func (m *MockService) Penalize(arg0 context.Context, arg1 *PeerId) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Penalize", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Penalize indicates an expected call of Penalize.
func (mr *MockServiceMockRecorder) Penalize(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Penalize", reflect.TypeOf((*MockService)(nil).Penalize), arg0, arg1)
}

// Start mocks base method.
func (m *MockService) Start(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", arg0)
}

// Start indicates an expected call of Start.
func (mr *MockServiceMockRecorder) Start(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockService)(nil).Start), arg0)
}

// Stop mocks base method.
func (m *MockService) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockServiceMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockService)(nil).Stop))
}
