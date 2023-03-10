// Code generated by MockGen. DO NOT EDIT.
// Source: cashback.go

// Package mock_repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	model "github.com/adinandradrs/cezbek-engine/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockCashbackPersister is a mock of CashbackPersister interface.
type MockCashbackPersister struct {
	ctrl     *gomock.Controller
	recorder *MockCashbackPersisterMockRecorder
}

// MockCashbackPersisterMockRecorder is the mock recorder for MockCashbackPersister.
type MockCashbackPersisterMockRecorder struct {
	mock *MockCashbackPersister
}

// NewMockCashbackPersister creates a new mock instance.
func NewMockCashbackPersister(ctrl *gomock.Controller) *MockCashbackPersister {
	mock := &MockCashbackPersister{ctrl: ctrl}
	mock.recorder = &MockCashbackPersisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCashbackPersister) EXPECT() *MockCashbackPersisterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCashbackPersister) Add(cashback model.Cashback) *model.TechnicalError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", cashback)
	ret0, _ := ret[0].(*model.TechnicalError)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockCashbackPersisterMockRecorder) Add(cashback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCashbackPersister)(nil).Add), cashback)
}
