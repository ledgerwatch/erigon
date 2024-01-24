// Code generated by MockGen. DO NOT EDIT.
// Source: io.go

// Package heimdall is a generated GoMock package.
package heimdall

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/ledgerwatch/erigon-lib/common"
	kv "github.com/ledgerwatch/erigon-lib/kv"
	rlp "github.com/ledgerwatch/erigon/rlp"
)

// MockSpanReader is a mock of SpanReader interface.
type MockSpanReader struct {
	ctrl     *gomock.Controller
	recorder *MockSpanReaderMockRecorder
}

// MockSpanReaderMockRecorder is the mock recorder for MockSpanReader.
type MockSpanReaderMockRecorder struct {
	mock *MockSpanReader
}

// NewMockSpanReader creates a new mock instance.
func NewMockSpanReader(ctrl *gomock.Controller) *MockSpanReader {
	mock := &MockSpanReader{ctrl: ctrl}
	mock.recorder = &MockSpanReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpanReader) EXPECT() *MockSpanReaderMockRecorder {
	return m.recorder
}

// LastSpanId mocks base method.
func (m *MockSpanReader) LastSpanId(ctx context.Context) (SpanId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastSpanId", ctx)
	ret0, _ := ret[0].(SpanId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastSpanId indicates an expected call of LastSpanId.
func (mr *MockSpanReaderMockRecorder) LastSpanId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastSpanId", reflect.TypeOf((*MockSpanReader)(nil).LastSpanId), ctx)
}

// ReadSpan mocks base method.
func (m *MockSpanReader) ReadSpan(ctx context.Context, spanId SpanId) (*Span, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSpan", ctx, spanId)
	ret0, _ := ret[0].(*Span)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSpan indicates an expected call of ReadSpan.
func (mr *MockSpanReaderMockRecorder) ReadSpan(ctx, spanId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSpan", reflect.TypeOf((*MockSpanReader)(nil).ReadSpan), ctx, spanId)
}

// MockSpanWriter is a mock of SpanWriter interface.
type MockSpanWriter struct {
	ctrl     *gomock.Controller
	recorder *MockSpanWriterMockRecorder
}

// MockSpanWriterMockRecorder is the mock recorder for MockSpanWriter.
type MockSpanWriterMockRecorder struct {
	mock *MockSpanWriter
}

// NewMockSpanWriter creates a new mock instance.
func NewMockSpanWriter(ctrl *gomock.Controller) *MockSpanWriter {
	mock := &MockSpanWriter{ctrl: ctrl}
	mock.recorder = &MockSpanWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpanWriter) EXPECT() *MockSpanWriterMockRecorder {
	return m.recorder
}

// WriteSpan mocks base method.
func (m *MockSpanWriter) WriteSpan(ctx context.Context, span *Span) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteSpan", ctx, span)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteSpan indicates an expected call of WriteSpan.
func (mr *MockSpanWriterMockRecorder) WriteSpan(ctx, span interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteSpan", reflect.TypeOf((*MockSpanWriter)(nil).WriteSpan), ctx, span)
}

// MockSpanIO is a mock of SpanIO interface.
type MockSpanIO struct {
	ctrl     *gomock.Controller
	recorder *MockSpanIOMockRecorder
}

// MockSpanIOMockRecorder is the mock recorder for MockSpanIO.
type MockSpanIOMockRecorder struct {
	mock *MockSpanIO
}

// NewMockSpanIO creates a new mock instance.
func NewMockSpanIO(ctrl *gomock.Controller) *MockSpanIO {
	mock := &MockSpanIO{ctrl: ctrl}
	mock.recorder = &MockSpanIOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpanIO) EXPECT() *MockSpanIOMockRecorder {
	return m.recorder
}

// LastSpanId mocks base method.
func (m *MockSpanIO) LastSpanId(ctx context.Context) (SpanId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastSpanId", ctx)
	ret0, _ := ret[0].(SpanId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastSpanId indicates an expected call of LastSpanId.
func (mr *MockSpanIOMockRecorder) LastSpanId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastSpanId", reflect.TypeOf((*MockSpanIO)(nil).LastSpanId), ctx)
}

// ReadSpan mocks base method.
func (m *MockSpanIO) ReadSpan(ctx context.Context, spanId SpanId) (*Span, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSpan", ctx, spanId)
	ret0, _ := ret[0].(*Span)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSpan indicates an expected call of ReadSpan.
func (mr *MockSpanIOMockRecorder) ReadSpan(ctx, spanId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSpan", reflect.TypeOf((*MockSpanIO)(nil).ReadSpan), ctx, spanId)
}

// WriteSpan mocks base method.
func (m *MockSpanIO) WriteSpan(ctx context.Context, span *Span) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteSpan", ctx, span)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteSpan indicates an expected call of WriteSpan.
func (mr *MockSpanIOMockRecorder) WriteSpan(ctx, span interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteSpan", reflect.TypeOf((*MockSpanIO)(nil).WriteSpan), ctx, span)
}

// MockMilestoneReader is a mock of MilestoneReader interface.
type MockMilestoneReader struct {
	ctrl     *gomock.Controller
	recorder *MockMilestoneReaderMockRecorder
}

// MockMilestoneReaderMockRecorder is the mock recorder for MockMilestoneReader.
type MockMilestoneReaderMockRecorder struct {
	mock *MockMilestoneReader
}

// NewMockMilestoneReader creates a new mock instance.
func NewMockMilestoneReader(ctrl *gomock.Controller) *MockMilestoneReader {
	mock := &MockMilestoneReader{ctrl: ctrl}
	mock.recorder = &MockMilestoneReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMilestoneReader) EXPECT() *MockMilestoneReaderMockRecorder {
	return m.recorder
}

// LastMilestoneId mocks base method.
func (m *MockMilestoneReader) LastMilestoneId(ctx context.Context) (MilestoneId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastMilestoneId", ctx)
	ret0, _ := ret[0].(MilestoneId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastMilestoneId indicates an expected call of LastMilestoneId.
func (mr *MockMilestoneReaderMockRecorder) LastMilestoneId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastMilestoneId", reflect.TypeOf((*MockMilestoneReader)(nil).LastMilestoneId), ctx)
}

// ReadMilestone mocks base method.
func (m *MockMilestoneReader) ReadMilestone(ctx context.Context, milestoneId MilestoneId) (*Milestone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMilestone", ctx, milestoneId)
	ret0, _ := ret[0].(*Milestone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadMilestone indicates an expected call of ReadMilestone.
func (mr *MockMilestoneReaderMockRecorder) ReadMilestone(ctx, milestoneId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMilestone", reflect.TypeOf((*MockMilestoneReader)(nil).ReadMilestone), ctx, milestoneId)
}

// MockMilestoneWriter is a mock of MilestoneWriter interface.
type MockMilestoneWriter struct {
	ctrl     *gomock.Controller
	recorder *MockMilestoneWriterMockRecorder
}

// MockMilestoneWriterMockRecorder is the mock recorder for MockMilestoneWriter.
type MockMilestoneWriterMockRecorder struct {
	mock *MockMilestoneWriter
}

// NewMockMilestoneWriter creates a new mock instance.
func NewMockMilestoneWriter(ctrl *gomock.Controller) *MockMilestoneWriter {
	mock := &MockMilestoneWriter{ctrl: ctrl}
	mock.recorder = &MockMilestoneWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMilestoneWriter) EXPECT() *MockMilestoneWriterMockRecorder {
	return m.recorder
}

// WriteMilestone mocks base method.
func (m *MockMilestoneWriter) WriteMilestone(ctx context.Context, milestone *Milestone) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMilestone", ctx, milestone)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMilestone indicates an expected call of WriteMilestone.
func (mr *MockMilestoneWriterMockRecorder) WriteMilestone(ctx, milestone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMilestone", reflect.TypeOf((*MockMilestoneWriter)(nil).WriteMilestone), ctx, milestone)
}

// MockMilestoneIO is a mock of MilestoneIO interface.
type MockMilestoneIO struct {
	ctrl     *gomock.Controller
	recorder *MockMilestoneIOMockRecorder
}

// MockMilestoneIOMockRecorder is the mock recorder for MockMilestoneIO.
type MockMilestoneIOMockRecorder struct {
	mock *MockMilestoneIO
}

// NewMockMilestoneIO creates a new mock instance.
func NewMockMilestoneIO(ctrl *gomock.Controller) *MockMilestoneIO {
	mock := &MockMilestoneIO{ctrl: ctrl}
	mock.recorder = &MockMilestoneIOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMilestoneIO) EXPECT() *MockMilestoneIOMockRecorder {
	return m.recorder
}

// LastMilestoneId mocks base method.
func (m *MockMilestoneIO) LastMilestoneId(ctx context.Context) (MilestoneId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastMilestoneId", ctx)
	ret0, _ := ret[0].(MilestoneId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastMilestoneId indicates an expected call of LastMilestoneId.
func (mr *MockMilestoneIOMockRecorder) LastMilestoneId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastMilestoneId", reflect.TypeOf((*MockMilestoneIO)(nil).LastMilestoneId), ctx)
}

// ReadMilestone mocks base method.
func (m *MockMilestoneIO) ReadMilestone(ctx context.Context, milestoneId MilestoneId) (*Milestone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMilestone", ctx, milestoneId)
	ret0, _ := ret[0].(*Milestone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadMilestone indicates an expected call of ReadMilestone.
func (mr *MockMilestoneIOMockRecorder) ReadMilestone(ctx, milestoneId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMilestone", reflect.TypeOf((*MockMilestoneIO)(nil).ReadMilestone), ctx, milestoneId)
}

// WriteMilestone mocks base method.
func (m *MockMilestoneIO) WriteMilestone(ctx context.Context, milestone *Milestone) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMilestone", ctx, milestone)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMilestone indicates an expected call of WriteMilestone.
func (mr *MockMilestoneIOMockRecorder) WriteMilestone(ctx, milestone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMilestone", reflect.TypeOf((*MockMilestoneIO)(nil).WriteMilestone), ctx, milestone)
}

// MockCheckpointReader is a mock of CheckpointReader interface.
type MockCheckpointReader struct {
	ctrl     *gomock.Controller
	recorder *MockCheckpointReaderMockRecorder
}

// MockCheckpointReaderMockRecorder is the mock recorder for MockCheckpointReader.
type MockCheckpointReaderMockRecorder struct {
	mock *MockCheckpointReader
}

// NewMockCheckpointReader creates a new mock instance.
func NewMockCheckpointReader(ctrl *gomock.Controller) *MockCheckpointReader {
	mock := &MockCheckpointReader{ctrl: ctrl}
	mock.recorder = &MockCheckpointReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCheckpointReader) EXPECT() *MockCheckpointReaderMockRecorder {
	return m.recorder
}

// LastCheckpointId mocks base method.
func (m *MockCheckpointReader) LastCheckpointId(ctx context.Context) (CheckpointId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastCheckpointId", ctx)
	ret0, _ := ret[0].(CheckpointId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastCheckpointId indicates an expected call of LastCheckpointId.
func (mr *MockCheckpointReaderMockRecorder) LastCheckpointId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastCheckpointId", reflect.TypeOf((*MockCheckpointReader)(nil).LastCheckpointId), ctx)
}

// ReadCheckpoint mocks base method.
func (m *MockCheckpointReader) ReadCheckpoint(ctx context.Context, checkpointId CheckpointId) (*Checkpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCheckpoint", ctx, checkpointId)
	ret0, _ := ret[0].(*Checkpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCheckpoint indicates an expected call of ReadCheckpoint.
func (mr *MockCheckpointReaderMockRecorder) ReadCheckpoint(ctx, checkpointId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCheckpoint", reflect.TypeOf((*MockCheckpointReader)(nil).ReadCheckpoint), ctx, checkpointId)
}

// MockCheckpointWriter is a mock of CheckpointWriter interface.
type MockCheckpointWriter struct {
	ctrl     *gomock.Controller
	recorder *MockCheckpointWriterMockRecorder
}

// MockCheckpointWriterMockRecorder is the mock recorder for MockCheckpointWriter.
type MockCheckpointWriterMockRecorder struct {
	mock *MockCheckpointWriter
}

// NewMockCheckpointWriter creates a new mock instance.
func NewMockCheckpointWriter(ctrl *gomock.Controller) *MockCheckpointWriter {
	mock := &MockCheckpointWriter{ctrl: ctrl}
	mock.recorder = &MockCheckpointWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCheckpointWriter) EXPECT() *MockCheckpointWriterMockRecorder {
	return m.recorder
}

// WriteCheckpoint mocks base method.
func (m *MockCheckpointWriter) WriteCheckpoint(ctx context.Context, checkpoint *Checkpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteCheckpoint", ctx, checkpoint)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteCheckpoint indicates an expected call of WriteCheckpoint.
func (mr *MockCheckpointWriterMockRecorder) WriteCheckpoint(ctx, checkpoint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteCheckpoint", reflect.TypeOf((*MockCheckpointWriter)(nil).WriteCheckpoint), ctx, checkpoint)
}

// MockCheckpointIO is a mock of CheckpointIO interface.
type MockCheckpointIO struct {
	ctrl     *gomock.Controller
	recorder *MockCheckpointIOMockRecorder
}

// MockCheckpointIOMockRecorder is the mock recorder for MockCheckpointIO.
type MockCheckpointIOMockRecorder struct {
	mock *MockCheckpointIO
}

// NewMockCheckpointIO creates a new mock instance.
func NewMockCheckpointIO(ctrl *gomock.Controller) *MockCheckpointIO {
	mock := &MockCheckpointIO{ctrl: ctrl}
	mock.recorder = &MockCheckpointIOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCheckpointIO) EXPECT() *MockCheckpointIOMockRecorder {
	return m.recorder
}

// LastCheckpointId mocks base method.
func (m *MockCheckpointIO) LastCheckpointId(ctx context.Context) (CheckpointId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastCheckpointId", ctx)
	ret0, _ := ret[0].(CheckpointId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastCheckpointId indicates an expected call of LastCheckpointId.
func (mr *MockCheckpointIOMockRecorder) LastCheckpointId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastCheckpointId", reflect.TypeOf((*MockCheckpointIO)(nil).LastCheckpointId), ctx)
}

// ReadCheckpoint mocks base method.
func (m *MockCheckpointIO) ReadCheckpoint(ctx context.Context, checkpointId CheckpointId) (*Checkpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCheckpoint", ctx, checkpointId)
	ret0, _ := ret[0].(*Checkpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCheckpoint indicates an expected call of ReadCheckpoint.
func (mr *MockCheckpointIOMockRecorder) ReadCheckpoint(ctx, checkpointId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCheckpoint", reflect.TypeOf((*MockCheckpointIO)(nil).ReadCheckpoint), ctx, checkpointId)
}

// WriteCheckpoint mocks base method.
func (m *MockCheckpointIO) WriteCheckpoint(ctx context.Context, checkpoint *Checkpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteCheckpoint", ctx, checkpoint)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteCheckpoint indicates an expected call of WriteCheckpoint.
func (mr *MockCheckpointIOMockRecorder) WriteCheckpoint(ctx, checkpoint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteCheckpoint", reflect.TypeOf((*MockCheckpointIO)(nil).WriteCheckpoint), ctx, checkpoint)
}

// MockIO is a mock of IO interface.
type MockIO struct {
	ctrl     *gomock.Controller
	recorder *MockIOMockRecorder
}

// MockIOMockRecorder is the mock recorder for MockIO.
type MockIOMockRecorder struct {
	mock *MockIO
}

// NewMockIO creates a new mock instance.
func NewMockIO(ctrl *gomock.Controller) *MockIO {
	mock := &MockIO{ctrl: ctrl}
	mock.recorder = &MockIOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIO) EXPECT() *MockIOMockRecorder {
	return m.recorder
}

// LastCheckpointId mocks base method.
func (m *MockIO) LastCheckpointId(ctx context.Context) (CheckpointId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastCheckpointId", ctx)
	ret0, _ := ret[0].(CheckpointId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastCheckpointId indicates an expected call of LastCheckpointId.
func (mr *MockIOMockRecorder) LastCheckpointId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastCheckpointId", reflect.TypeOf((*MockIO)(nil).LastCheckpointId), ctx)
}

// LastMilestoneId mocks base method.
func (m *MockIO) LastMilestoneId(ctx context.Context) (MilestoneId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastMilestoneId", ctx)
	ret0, _ := ret[0].(MilestoneId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastMilestoneId indicates an expected call of LastMilestoneId.
func (mr *MockIOMockRecorder) LastMilestoneId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastMilestoneId", reflect.TypeOf((*MockIO)(nil).LastMilestoneId), ctx)
}

// LastSpanId mocks base method.
func (m *MockIO) LastSpanId(ctx context.Context) (SpanId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastSpanId", ctx)
	ret0, _ := ret[0].(SpanId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastSpanId indicates an expected call of LastSpanId.
func (mr *MockIOMockRecorder) LastSpanId(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastSpanId", reflect.TypeOf((*MockIO)(nil).LastSpanId), ctx)
}

// ReadCheckpoint mocks base method.
func (m *MockIO) ReadCheckpoint(ctx context.Context, checkpointId CheckpointId) (*Checkpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadCheckpoint", ctx, checkpointId)
	ret0, _ := ret[0].(*Checkpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadCheckpoint indicates an expected call of ReadCheckpoint.
func (mr *MockIOMockRecorder) ReadCheckpoint(ctx, checkpointId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadCheckpoint", reflect.TypeOf((*MockIO)(nil).ReadCheckpoint), ctx, checkpointId)
}

// ReadMilestone mocks base method.
func (m *MockIO) ReadMilestone(ctx context.Context, milestoneId MilestoneId) (*Milestone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMilestone", ctx, milestoneId)
	ret0, _ := ret[0].(*Milestone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadMilestone indicates an expected call of ReadMilestone.
func (mr *MockIOMockRecorder) ReadMilestone(ctx, milestoneId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMilestone", reflect.TypeOf((*MockIO)(nil).ReadMilestone), ctx, milestoneId)
}

// ReadSpan mocks base method.
func (m *MockIO) ReadSpan(ctx context.Context, spanId SpanId) (*Span, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSpan", ctx, spanId)
	ret0, _ := ret[0].(*Span)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSpan indicates an expected call of ReadSpan.
func (mr *MockIOMockRecorder) ReadSpan(ctx, spanId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSpan", reflect.TypeOf((*MockIO)(nil).ReadSpan), ctx, spanId)
}

// WriteCheckpoint mocks base method.
func (m *MockIO) WriteCheckpoint(ctx context.Context, checkpoint *Checkpoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteCheckpoint", ctx, checkpoint)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteCheckpoint indicates an expected call of WriteCheckpoint.
func (mr *MockIOMockRecorder) WriteCheckpoint(ctx, checkpoint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteCheckpoint", reflect.TypeOf((*MockIO)(nil).WriteCheckpoint), ctx, checkpoint)
}

// WriteMilestone mocks base method.
func (m *MockIO) WriteMilestone(ctx context.Context, milestone *Milestone) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMilestone", ctx, milestone)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMilestone indicates an expected call of WriteMilestone.
func (mr *MockIOMockRecorder) WriteMilestone(ctx, milestone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMilestone", reflect.TypeOf((*MockIO)(nil).WriteMilestone), ctx, milestone)
}

// WriteSpan mocks base method.
func (m *MockIO) WriteSpan(ctx context.Context, span *Span) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteSpan", ctx, span)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteSpan indicates an expected call of WriteSpan.
func (mr *MockIOMockRecorder) WriteSpan(ctx, span interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteSpan", reflect.TypeOf((*MockIO)(nil).WriteSpan), ctx, span)
}

// Mockreader is a mock of reader interface.
type Mockreader struct {
	ctrl     *gomock.Controller
	recorder *MockreaderMockRecorder
}

// MockreaderMockRecorder is the mock recorder for Mockreader.
type MockreaderMockRecorder struct {
	mock *Mockreader
}

// NewMockreader creates a new mock instance.
func NewMockreader(ctrl *gomock.Controller) *Mockreader {
	mock := &Mockreader{ctrl: ctrl}
	mock.recorder = &MockreaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockreader) EXPECT() *MockreaderMockRecorder {
	return m.recorder
}

// Checkpoint mocks base method.
func (m *Mockreader) Checkpoint(ctx context.Context, tx kv.Getter, checkpointId uint64) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checkpoint", ctx, tx, checkpointId)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Checkpoint indicates an expected call of Checkpoint.
func (mr *MockreaderMockRecorder) Checkpoint(ctx, tx, checkpointId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checkpoint", reflect.TypeOf((*Mockreader)(nil).Checkpoint), ctx, tx, checkpointId)
}

// EventLookup mocks base method.
func (m *Mockreader) EventLookup(ctx context.Context, tx kv.Getter, txnHash common.Hash) (uint64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventLookup", ctx, tx, txnHash)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EventLookup indicates an expected call of EventLookup.
func (mr *MockreaderMockRecorder) EventLookup(ctx, tx, txnHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventLookup", reflect.TypeOf((*Mockreader)(nil).EventLookup), ctx, tx, txnHash)
}

// EventsByBlock mocks base method.
func (m *Mockreader) EventsByBlock(ctx context.Context, tx kv.Tx, hash common.Hash, blockNum uint64) ([]rlp.RawValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventsByBlock", ctx, tx, hash, blockNum)
	ret0, _ := ret[0].([]rlp.RawValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EventsByBlock indicates an expected call of EventsByBlock.
func (mr *MockreaderMockRecorder) EventsByBlock(ctx, tx, hash, blockNum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventsByBlock", reflect.TypeOf((*Mockreader)(nil).EventsByBlock), ctx, tx, hash, blockNum)
}

// LastEventId mocks base method.
func (m *Mockreader) LastEventId(ctx context.Context, tx kv.Tx) (uint64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastEventId", ctx, tx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastEventId indicates an expected call of LastEventId.
func (mr *MockreaderMockRecorder) LastEventId(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastEventId", reflect.TypeOf((*Mockreader)(nil).LastEventId), ctx, tx)
}

// LastSpanId mocks base method.
func (m *Mockreader) LastSpanId(ctx context.Context, tx kv.Tx) (uint64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastSpanId", ctx, tx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LastSpanId indicates an expected call of LastSpanId.
func (mr *MockreaderMockRecorder) LastSpanId(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastSpanId", reflect.TypeOf((*Mockreader)(nil).LastSpanId), ctx, tx)
}

// Milestone mocks base method.
func (m *Mockreader) Milestone(ctx context.Context, tx kv.Getter, milestoneId uint64) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Milestone", ctx, tx, milestoneId)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Milestone indicates an expected call of Milestone.
func (mr *MockreaderMockRecorder) Milestone(ctx, tx, milestoneId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Milestone", reflect.TypeOf((*Mockreader)(nil).Milestone), ctx, tx, milestoneId)
}

// Span mocks base method.
func (m *Mockreader) Span(ctx context.Context, tx kv.Getter, spanId uint64) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Span", ctx, tx, spanId)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Span indicates an expected call of Span.
func (mr *MockreaderMockRecorder) Span(ctx, tx, spanId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Span", reflect.TypeOf((*Mockreader)(nil).Span), ctx, tx, spanId)
}
