// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ledgerwatch/erigon-lib/gointerfaces/remoteproto (interfaces: KVClient)
//
// Generated by this command:
//
//	mockgen -typed=true -destination=./kv_client_mock.go -package=remoteproto . KVClient
//

// Package remoteproto is a generated GoMock package.
package remoteproto

import (
	context "context"
	reflect "reflect"

	typesproto "github.com/ledgerwatch/erigon-lib/gointerfaces/typesproto"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockKVClient is a mock of KVClient interface.
type MockKVClient struct {
	ctrl     *gomock.Controller
	recorder *MockKVClientMockRecorder
}

// MockKVClientMockRecorder is the mock recorder for MockKVClient.
type MockKVClientMockRecorder struct {
	mock *MockKVClient
}

// NewMockKVClient creates a new mock instance.
func NewMockKVClient(ctrl *gomock.Controller) *MockKVClient {
	mock := &MockKVClient{ctrl: ctrl}
	mock.recorder = &MockKVClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKVClient) EXPECT() *MockKVClientMockRecorder {
	return m.recorder
}

// DomainGet mocks base method.
func (m *MockKVClient) DomainGet(arg0 context.Context, arg1 *DomainGetReq, arg2 ...grpc.CallOption) (*DomainGetReply, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DomainGet", varargs...)
	ret0, _ := ret[0].(*DomainGetReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainGet indicates an expected call of DomainGet.
func (mr *MockKVClientMockRecorder) DomainGet(arg0, arg1 any, arg2 ...any) *MockKVClientDomainGetCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainGet", reflect.TypeOf((*MockKVClient)(nil).DomainGet), varargs...)
	return &MockKVClientDomainGetCall{Call: call}
}

// MockKVClientDomainGetCall wrap *gomock.Call
type MockKVClientDomainGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientDomainGetCall) Return(arg0 *DomainGetReply, arg1 error) *MockKVClientDomainGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientDomainGetCall) Do(f func(context.Context, *DomainGetReq, ...grpc.CallOption) (*DomainGetReply, error)) *MockKVClientDomainGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientDomainGetCall) DoAndReturn(f func(context.Context, *DomainGetReq, ...grpc.CallOption) (*DomainGetReply, error)) *MockKVClientDomainGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DomainRange mocks base method.
func (m *MockKVClient) DomainRange(arg0 context.Context, arg1 *DomainRangeReq, arg2 ...grpc.CallOption) (*Pairs, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DomainRange", varargs...)
	ret0, _ := ret[0].(*Pairs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DomainRange indicates an expected call of DomainRange.
func (mr *MockKVClientMockRecorder) DomainRange(arg0, arg1 any, arg2 ...any) *MockKVClientDomainRangeCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DomainRange", reflect.TypeOf((*MockKVClient)(nil).DomainRange), varargs...)
	return &MockKVClientDomainRangeCall{Call: call}
}

// MockKVClientDomainRangeCall wrap *gomock.Call
type MockKVClientDomainRangeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientDomainRangeCall) Return(arg0 *Pairs, arg1 error) *MockKVClientDomainRangeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientDomainRangeCall) Do(f func(context.Context, *DomainRangeReq, ...grpc.CallOption) (*Pairs, error)) *MockKVClientDomainRangeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientDomainRangeCall) DoAndReturn(f func(context.Context, *DomainRangeReq, ...grpc.CallOption) (*Pairs, error)) *MockKVClientDomainRangeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HistoryGet mocks base method.
func (m *MockKVClient) HistoryGet(arg0 context.Context, arg1 *HistoryGetReq, arg2 ...grpc.CallOption) (*HistoryGetReply, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HistorySeek", varargs...)
	ret0, _ := ret[0].(*HistoryGetReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HistoryGet indicates an expected call of HistoryGet.
func (mr *MockKVClientMockRecorder) HistoryGet(arg0, arg1 any, arg2 ...any) *MockKVClientHistoryGetCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HistorySeek", reflect.TypeOf((*MockKVClient)(nil).HistoryGet), varargs...)
	return &MockKVClientHistoryGetCall{Call: call}
}

// MockKVClientHistoryGetCall wrap *gomock.Call
type MockKVClientHistoryGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientHistoryGetCall) Return(arg0 *HistoryGetReply, arg1 error) *MockKVClientHistoryGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientHistoryGetCall) Do(f func(context.Context, *HistoryGetReq, ...grpc.CallOption) (*HistoryGetReply, error)) *MockKVClientHistoryGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientHistoryGetCall) DoAndReturn(f func(context.Context, *HistoryGetReq, ...grpc.CallOption) (*HistoryGetReply, error)) *MockKVClientHistoryGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HistoryRange mocks base method.
func (m *MockKVClient) HistoryRange(arg0 context.Context, arg1 *HistoryRangeReq, arg2 ...grpc.CallOption) (*Pairs, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HistoryRange", varargs...)
	ret0, _ := ret[0].(*Pairs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HistoryRange indicates an expected call of HistoryRange.
func (mr *MockKVClientMockRecorder) HistoryRange(arg0, arg1 any, arg2 ...any) *MockKVClientHistoryRangeCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HistoryRange", reflect.TypeOf((*MockKVClient)(nil).HistoryRange), varargs...)
	return &MockKVClientHistoryRangeCall{Call: call}
}

// MockKVClientHistoryRangeCall wrap *gomock.Call
type MockKVClientHistoryRangeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientHistoryRangeCall) Return(arg0 *Pairs, arg1 error) *MockKVClientHistoryRangeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientHistoryRangeCall) Do(f func(context.Context, *HistoryRangeReq, ...grpc.CallOption) (*Pairs, error)) *MockKVClientHistoryRangeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientHistoryRangeCall) DoAndReturn(f func(context.Context, *HistoryRangeReq, ...grpc.CallOption) (*Pairs, error)) *MockKVClientHistoryRangeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IndexRange mocks base method.
func (m *MockKVClient) IndexRange(arg0 context.Context, arg1 *IndexRangeReq, arg2 ...grpc.CallOption) (*IndexRangeReply, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IndexRange", varargs...)
	ret0, _ := ret[0].(*IndexRangeReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexRange indicates an expected call of IndexRange.
func (mr *MockKVClientMockRecorder) IndexRange(arg0, arg1 any, arg2 ...any) *MockKVClientIndexRangeCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexRange", reflect.TypeOf((*MockKVClient)(nil).IndexRange), varargs...)
	return &MockKVClientIndexRangeCall{Call: call}
}

// MockKVClientIndexRangeCall wrap *gomock.Call
type MockKVClientIndexRangeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientIndexRangeCall) Return(arg0 *IndexRangeReply, arg1 error) *MockKVClientIndexRangeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientIndexRangeCall) Do(f func(context.Context, *IndexRangeReq, ...grpc.CallOption) (*IndexRangeReply, error)) *MockKVClientIndexRangeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientIndexRangeCall) DoAndReturn(f func(context.Context, *IndexRangeReq, ...grpc.CallOption) (*IndexRangeReply, error)) *MockKVClientIndexRangeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Range mocks base method.
func (m *MockKVClient) Range(arg0 context.Context, arg1 *RangeReq, arg2 ...grpc.CallOption) (*Pairs, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Range", varargs...)
	ret0, _ := ret[0].(*Pairs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Range indicates an expected call of Range.
func (mr *MockKVClientMockRecorder) Range(arg0, arg1 any, arg2 ...any) *MockKVClientRangeCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Range", reflect.TypeOf((*MockKVClient)(nil).Range), varargs...)
	return &MockKVClientRangeCall{Call: call}
}

// MockKVClientRangeCall wrap *gomock.Call
type MockKVClientRangeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientRangeCall) Return(arg0 *Pairs, arg1 error) *MockKVClientRangeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientRangeCall) Do(f func(context.Context, *RangeReq, ...grpc.CallOption) (*Pairs, error)) *MockKVClientRangeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientRangeCall) DoAndReturn(f func(context.Context, *RangeReq, ...grpc.CallOption) (*Pairs, error)) *MockKVClientRangeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Snapshots mocks base method.
func (m *MockKVClient) Snapshots(arg0 context.Context, arg1 *SnapshotsRequest, arg2 ...grpc.CallOption) (*SnapshotsReply, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Snapshots", varargs...)
	ret0, _ := ret[0].(*SnapshotsReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Snapshots indicates an expected call of Snapshots.
func (mr *MockKVClientMockRecorder) Snapshots(arg0, arg1 any, arg2 ...any) *MockKVClientSnapshotsCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Snapshots", reflect.TypeOf((*MockKVClient)(nil).Snapshots), varargs...)
	return &MockKVClientSnapshotsCall{Call: call}
}

// MockKVClientSnapshotsCall wrap *gomock.Call
type MockKVClientSnapshotsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientSnapshotsCall) Return(arg0 *SnapshotsReply, arg1 error) *MockKVClientSnapshotsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientSnapshotsCall) Do(f func(context.Context, *SnapshotsRequest, ...grpc.CallOption) (*SnapshotsReply, error)) *MockKVClientSnapshotsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientSnapshotsCall) DoAndReturn(f func(context.Context, *SnapshotsRequest, ...grpc.CallOption) (*SnapshotsReply, error)) *MockKVClientSnapshotsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StateChanges mocks base method.
func (m *MockKVClient) StateChanges(arg0 context.Context, arg1 *StateChangeRequest, arg2 ...grpc.CallOption) (KV_StateChangesClient, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StateChanges", varargs...)
	ret0, _ := ret[0].(KV_StateChangesClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StateChanges indicates an expected call of StateChanges.
func (mr *MockKVClientMockRecorder) StateChanges(arg0, arg1 any, arg2 ...any) *MockKVClientStateChangesCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateChanges", reflect.TypeOf((*MockKVClient)(nil).StateChanges), varargs...)
	return &MockKVClientStateChangesCall{Call: call}
}

// MockKVClientStateChangesCall wrap *gomock.Call
type MockKVClientStateChangesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientStateChangesCall) Return(arg0 KV_StateChangesClient, arg1 error) *MockKVClientStateChangesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientStateChangesCall) Do(f func(context.Context, *StateChangeRequest, ...grpc.CallOption) (KV_StateChangesClient, error)) *MockKVClientStateChangesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientStateChangesCall) DoAndReturn(f func(context.Context, *StateChangeRequest, ...grpc.CallOption) (KV_StateChangesClient, error)) *MockKVClientStateChangesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Tx mocks base method.
func (m *MockKVClient) Tx(arg0 context.Context, arg1 ...grpc.CallOption) (KV_TxClient, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Tx", varargs...)
	ret0, _ := ret[0].(KV_TxClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Tx indicates an expected call of Tx.
func (mr *MockKVClientMockRecorder) Tx(arg0 any, arg1 ...any) *MockKVClientTxCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockKVClient)(nil).Tx), varargs...)
	return &MockKVClientTxCall{Call: call}
}

// MockKVClientTxCall wrap *gomock.Call
type MockKVClientTxCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientTxCall) Return(arg0 KV_TxClient, arg1 error) *MockKVClientTxCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientTxCall) Do(f func(context.Context, ...grpc.CallOption) (KV_TxClient, error)) *MockKVClientTxCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientTxCall) DoAndReturn(f func(context.Context, ...grpc.CallOption) (KV_TxClient, error)) *MockKVClientTxCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Version mocks base method.
func (m *MockKVClient) Version(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*typesproto.VersionReply, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Version", varargs...)
	ret0, _ := ret[0].(*typesproto.VersionReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Version indicates an expected call of Version.
func (mr *MockKVClientMockRecorder) Version(arg0, arg1 any, arg2 ...any) *MockKVClientVersionCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockKVClient)(nil).Version), varargs...)
	return &MockKVClientVersionCall{Call: call}
}

// MockKVClientVersionCall wrap *gomock.Call
type MockKVClientVersionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockKVClientVersionCall) Return(arg0 *typesproto.VersionReply, arg1 error) *MockKVClientVersionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockKVClientVersionCall) Do(f func(context.Context, *emptypb.Empty, ...grpc.CallOption) (*typesproto.VersionReply, error)) *MockKVClientVersionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockKVClientVersionCall) DoAndReturn(f func(context.Context, *emptypb.Empty, ...grpc.CallOption) (*typesproto.VersionReply, error)) *MockKVClientVersionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
