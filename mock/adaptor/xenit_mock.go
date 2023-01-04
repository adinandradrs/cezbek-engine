// Code generated by MockGen. DO NOT EDIT.
// Source: xenit.go

// Package mock_adaptor is a generated GoMock package.
package adaptor

import (
	reflect "reflect"

	model "github.com/adinandradrs/cezbek-engine/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockXenitAdapter is a mock of XenitAdapter interface.
type MockXenitAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockXenitAdapterMockRecorder
}

// MockXenitAdapterMockRecorder is the mock recorder for MockXenitAdapter.
type MockXenitAdapterMockRecorder struct {
	mock *MockXenitAdapter
}

// NewMockXenitAdapter creates a new mock instance.
func NewMockXenitAdapter(ctrl *gomock.Controller) *MockXenitAdapter {
	mock := &MockXenitAdapter{ctrl: ctrl}
	mock.recorder = &MockXenitAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockXenitAdapter) EXPECT() *MockXenitAdapterMockRecorder {
	return m.recorder
}

// WalletTopup mocks base method.
func (m *MockXenitAdapter) WalletTopup(inp *model.XenitWalletTopupRequest) (*model.XenitWalletTopupResponse, *model.TechnicalError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WalletTopup", inp)
	ret0, _ := ret[0].(*model.XenitWalletTopupResponse)
	ret1, _ := ret[1].(*model.TechnicalError)
	return ret0, ret1
}

// WalletTopup indicates an expected call of WalletTopup.
func (mr *MockXenitAdapterMockRecorder) WalletTopup(inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WalletTopup", reflect.TypeOf((*MockXenitAdapter)(nil).WalletTopup), inp)
}
