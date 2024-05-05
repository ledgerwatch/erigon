// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ledgerwatch/erigon/consensus (interfaces: ChainHeaderReader)
//
// Generated by this command:
//
//	mockgen -typed=true -destination=./chain_header_reader_mock.go -package=consensus . ChainHeaderReader
//

// Package consensus is a generated GoMock package.
package consensus

import (
	big "math/big"
	reflect "reflect"

	chain "github.com/ledgerwatch/erigon-lib/chain"
	common "github.com/ledgerwatch/erigon-lib/common"
	types "github.com/ledgerwatch/erigon/core/types"
	gomock "go.uber.org/mock/gomock"
)

// MockChainHeaderReader is a mock of ChainHeaderReader interface.
type MockChainHeaderReader struct {
	ctrl     *gomock.Controller
	recorder *MockChainHeaderReaderMockRecorder
}

// MockChainHeaderReaderMockRecorder is the mock recorder for MockChainHeaderReader.
type MockChainHeaderReaderMockRecorder struct {
	mock *MockChainHeaderReader
}

// NewMockChainHeaderReader creates a new mock instance.
func NewMockChainHeaderReader(ctrl *gomock.Controller) *MockChainHeaderReader {
	mock := &MockChainHeaderReader{ctrl: ctrl}
	mock.recorder = &MockChainHeaderReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainHeaderReader) EXPECT() *MockChainHeaderReaderMockRecorder {
	return m.recorder
}

// BorSpan mocks base method.
func (m *MockChainHeaderReader) BorSpan(arg0 uint64) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BorSpan", arg0)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// BorSpan indicates an expected call of BorSpan.
func (mr *MockChainHeaderReaderMockRecorder) BorSpan(arg0 any) *MockChainHeaderReaderBorSpanCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BorSpan", reflect.TypeOf((*MockChainHeaderReader)(nil).BorSpan), arg0)
	return &MockChainHeaderReaderBorSpanCall{Call: call}
}

// MockChainHeaderReaderBorSpanCall wrap *gomock.Call
type MockChainHeaderReaderBorSpanCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockChainHeaderReaderBorSpanCall) Return(arg0 []byte) *MockChainHeaderReaderBorSpanCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockChainHeaderReaderBorSpanCall) Do(f func(uint64) []byte) *MockChainHeaderReaderBorSpanCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockChainHeaderReaderBorSpanCall) DoAndReturn(f func(uint64) []byte) *MockChainHeaderReaderBorSpanCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Config mocks base method.
func (m *MockChainHeaderReader) Config() *chain.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*chain.Config)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockChainHeaderReaderMockRecorder) Config() *MockChainHeaderReaderConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockChainHeaderReader)(nil).Config))
	return &MockChainHeaderReaderConfigCall{Call: call}
}

// MockChainHeaderReaderConfigCall wrap *gomock.Call
type MockChainHeaderReaderConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockChainHeaderReaderConfigCall) Return(arg0 *chain.Config) *MockChainHeaderReaderConfigCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockChainHeaderReaderConfigCall) Do(f func() *chain.Config) *MockChainHeaderReaderConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockChainHeaderReaderConfigCall) DoAndReturn(f func() *chain.Config) *MockChainHeaderReaderConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CurrentHeader mocks base method.
func (m *MockChainHeaderReader) CurrentHeader() *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentHeader")
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// CurrentHeader indicates an expected call of CurrentHeader.
func (mr *MockChainHeaderReaderMockRecorder) CurrentHeader() *MockChainHeaderReaderCurrentHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentHeader", reflect.TypeOf((*MockChainHeaderReader)(nil).CurrentHeader))
	return &MockChainHeaderReaderCurrentHeaderCall{Call: call}
}

// MockChainHeaderReaderCurrentHeaderCall wrap *gomock.Call
type MockChainHeaderReaderCurrentHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockChainHeaderReaderCurrentHeaderCall) Return(arg0 *types.Header) *MockChainHeaderReaderCurrentHeaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockChainHeaderReaderCurrentHeaderCall) Do(f func() *types.Header) *MockChainHeaderReaderCurrentHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockChainHeaderReaderCurrentHeaderCall) DoAndReturn(f func() *types.Header) *MockChainHeaderReaderCurrentHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// FrozenBlocks mocks base method.
func (m *MockChainHeaderReader) FrozenBlocks() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FrozenBlocks")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// FrozenBlocks indicates an expected call of FrozenBlocks.
func (mr *MockChainHeaderReaderMockRecorder) FrozenBlocks() *MockChainHeaderReaderFrozenBlocksCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FrozenBlocks", reflect.TypeOf((*MockChainHeaderReader)(nil).FrozenBlocks))
	return &MockChainHeaderReaderFrozenBlocksCall{Call: call}
}

// MockChainHeaderReaderFrozenBlocksCall wrap *gomock.Call
type MockChainHeaderReaderFrozenBlocksCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockChainHeaderReaderFrozenBlocksCall) Return(arg0 uint64) *MockChainHeaderReaderFrozenBlocksCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockChainHeaderReaderFrozenBlocksCall) Do(f func() uint64) *MockChainHeaderReaderFrozenBlocksCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockChainHeaderReaderFrozenBlocksCall) DoAndReturn(f func() uint64) *MockChainHeaderReaderFrozenBlocksCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetHeader mocks base method.
func (m *MockChainHeaderReader) GetHeader(arg0 common.Hash, arg1 uint64) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeader", arg0, arg1)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeader indicates an expected call of GetHeader.
func (mr *MockChainHeaderReaderMockRecorder) GetHeader(arg0, arg1 any) *MockChainHeaderReaderGetHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeader", reflect.TypeOf((*MockChainHeaderReader)(nil).GetHeader), arg0, arg1)
	return &MockChainHeaderReaderGetHeaderCall{Call: call}
}

// MockChainHeaderReaderGetHeaderCall wrap *gomock.Call
type MockChainHeaderReaderGetHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockChainHeaderReaderGetHeaderCall) Return(arg0 *types.Header) *MockChainHeaderReaderGetHeaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockChainHeaderReaderGetHeaderCall) Do(f func(common.Hash, uint64) *types.Header) *MockChainHeaderReaderGetHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockChainHeaderReaderGetHeaderCall) DoAndReturn(f func(common.Hash, uint64) *types.Header) *MockChainHeaderReaderGetHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetHeaderByHash mocks base method.
func (m *MockChainHeaderReader) GetHeaderByHash(arg0 common.Hash) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeaderByHash", arg0)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeaderByHash indicates an expected call of GetHeaderByHash.
func (mr *MockChainHeaderReaderMockRecorder) GetHeaderByHash(arg0 any) *MockChainHeaderReaderGetHeaderByHashCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeaderByHash", reflect.TypeOf((*MockChainHeaderReader)(nil).GetHeaderByHash), arg0)
	return &MockChainHeaderReaderGetHeaderByHashCall{Call: call}
}

// MockChainHeaderReaderGetHeaderByHashCall wrap *gomock.Call
type MockChainHeaderReaderGetHeaderByHashCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockChainHeaderReaderGetHeaderByHashCall) Return(arg0 *types.Header) *MockChainHeaderReaderGetHeaderByHashCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockChainHeaderReaderGetHeaderByHashCall) Do(f func(common.Hash) *types.Header) *MockChainHeaderReaderGetHeaderByHashCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockChainHeaderReaderGetHeaderByHashCall) DoAndReturn(f func(common.Hash) *types.Header) *MockChainHeaderReaderGetHeaderByHashCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetHeaderByNumber mocks base method.
func (m *MockChainHeaderReader) GetHeaderByNumber(arg0 uint64) *types.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeaderByNumber", arg0)
	ret0, _ := ret[0].(*types.Header)
	return ret0
}

// GetHeaderByNumber indicates an expected call of GetHeaderByNumber.
func (mr *MockChainHeaderReaderMockRecorder) GetHeaderByNumber(arg0 any) *MockChainHeaderReaderGetHeaderByNumberCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeaderByNumber", reflect.TypeOf((*MockChainHeaderReader)(nil).GetHeaderByNumber), arg0)
	return &MockChainHeaderReaderGetHeaderByNumberCall{Call: call}
}

// MockChainHeaderReaderGetHeaderByNumberCall wrap *gomock.Call
type MockChainHeaderReaderGetHeaderByNumberCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockChainHeaderReaderGetHeaderByNumberCall) Return(arg0 *types.Header) *MockChainHeaderReaderGetHeaderByNumberCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockChainHeaderReaderGetHeaderByNumberCall) Do(f func(uint64) *types.Header) *MockChainHeaderReaderGetHeaderByNumberCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockChainHeaderReaderGetHeaderByNumberCall) DoAndReturn(f func(uint64) *types.Header) *MockChainHeaderReaderGetHeaderByNumberCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetTd mocks base method.
func (m *MockChainHeaderReader) GetTd(arg0 common.Hash, arg1 uint64) *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTd", arg0, arg1)
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// GetTd indicates an expected call of GetTd.
func (mr *MockChainHeaderReaderMockRecorder) GetTd(arg0, arg1 any) *MockChainHeaderReaderGetTdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTd", reflect.TypeOf((*MockChainHeaderReader)(nil).GetTd), arg0, arg1)
	return &MockChainHeaderReaderGetTdCall{Call: call}
}

// MockChainHeaderReaderGetTdCall wrap *gomock.Call
type MockChainHeaderReaderGetTdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockChainHeaderReaderGetTdCall) Return(arg0 *big.Int) *MockChainHeaderReaderGetTdCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockChainHeaderReaderGetTdCall) Do(f func(common.Hash, uint64) *big.Int) *MockChainHeaderReaderGetTdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockChainHeaderReaderGetTdCall) DoAndReturn(f func(common.Hash, uint64) *big.Int) *MockChainHeaderReaderGetTdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
