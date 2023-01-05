// Code generated by MockGen. DO NOT EDIT.
// Source: josvo.go

// Package mock_adaptor is a generated GoMock package.
package adaptor

import (
	reflect "reflect"

	model "github.com/adinandradrs/cezbek-engine/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockJosvoAdapter is a mock of JosvoAdapter interface.
type MockJosvoAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockJosvoAdapterMockRecorder
}

// MockJosvoAdapterMockRecorder is the mock recorder for MockJosvoAdapter.
type MockJosvoAdapterMockRecorder struct {
	mock *MockJosvoAdapter
}

// NewMockJosvoAdapter creates a new mock instance.
func NewMockJosvoAdapter(ctrl *gomock.Controller) *MockJosvoAdapter {
	mock := &MockJosvoAdapter{ctrl: ctrl}
	mock.recorder = &MockJosvoAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJosvoAdapter) EXPECT() *MockJosvoAdapterMockRecorder {
	return m.recorder
}

// AccountTransfer mocks base method.
func (m *MockJosvoAdapter) AccountTransfer(inp *model.JosvoAccountTransferRequest) (*model.JosvoAccountTransferResponse, *model.TechnicalError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountTransfer", inp)
	ret0, _ := ret[0].(*model.JosvoAccountTransferResponse)
	ret1, _ := ret[1].(*model.TechnicalError)
	return ret0, ret1
}

// AccountTransfer indicates an expected call of AccountTransfer.
func (mr *MockJosvoAdapterMockRecorder) AccountTransfer(inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountTransfer", reflect.TypeOf((*MockJosvoAdapter)(nil).AccountTransfer), inp)
}