// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ledgerwatch/erigon/cl/phase1/network/services (interfaces: AttestationService)
//
// Generated by this command:
//
//	mockgen -destination=./mock_services/attestation_service_mock.go -package=mock_services . AttestationService
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	solid "github.com/ledgerwatch/erigon/cl/cltypes/solid"
	gomock "go.uber.org/mock/gomock"
)

// MockAttestationService is a mock of AttestationService interface.
type MockAttestationService struct {
	ctrl     *gomock.Controller
	recorder *MockAttestationServiceMockRecorder
}

// MockAttestationServiceMockRecorder is the mock recorder for MockAttestationService.
type MockAttestationServiceMockRecorder struct {
	mock *MockAttestationService
}

// NewMockAttestationService creates a new mock instance.
func NewMockAttestationService(ctrl *gomock.Controller) *MockAttestationService {
	mock := &MockAttestationService{ctrl: ctrl}
	mock.recorder = &MockAttestationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAttestationService) EXPECT() *MockAttestationServiceMockRecorder {
	return m.recorder
}

// ProcessMessage mocks base method.
func (m *MockAttestationService) ProcessMessage(arg0 context.Context, arg1 *uint64, arg2 *solid.Attestation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessMessage", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessMessage indicates an expected call of ProcessMessage.
func (mr *MockAttestationServiceMockRecorder) ProcessMessage(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessMessage", reflect.TypeOf((*MockAttestationService)(nil).ProcessMessage), arg0, arg1, arg2)
}
