// Code generated by MockGen. DO NOT EDIT.
// Source: consul.go

// Package mock_adaptor is a generated GoMock package.
package adaptor

import (
	reflect "reflect"

	model "github.com/adinandradrs/cezbek-engine/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockConsulWatcher is a mock of ConsulWatcher interface.
type MockConsulWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockConsulWatcherMockRecorder
}

// MockConsulWatcherMockRecorder is the mock recorder for MockConsulWatcher.
type MockConsulWatcherMockRecorder struct {
	mock *MockConsulWatcher
}

// NewMockConsulWatcher creates a new mock instance.
func NewMockConsulWatcher(ctrl *gomock.Controller) *MockConsulWatcher {
	mock := &MockConsulWatcher{ctrl: ctrl}
	mock.recorder = &MockConsulWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsulWatcher) EXPECT() *MockConsulWatcherMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockConsulWatcher) Register() *model.TechnicalError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register")
	ret0, _ := ret[0].(*model.TechnicalError)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockConsulWatcherMockRecorder) Register() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockConsulWatcher)(nil).Register))
}
