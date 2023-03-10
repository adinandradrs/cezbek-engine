// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package mock_repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	model "github.com/adinandradrs/cezbek-engine/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockTransactionPersister is a mock of TransactionPersister interface.
type MockTransactionPersister struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionPersisterMockRecorder
}

// MockTransactionPersisterMockRecorder is the mock recorder for MockTransactionPersister.
type MockTransactionPersisterMockRecorder struct {
	mock *MockTransactionPersister
}

// NewMockTransactionPersister creates a new mock instance.
func NewMockTransactionPersister(ctrl *gomock.Controller) *MockTransactionPersister {
	mock := &MockTransactionPersister{ctrl: ctrl}
	mock.recorder = &MockTransactionPersisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionPersister) EXPECT() *MockTransactionPersisterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockTransactionPersister) Add(trx model.Transaction) (*int64, *model.TechnicalError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", trx)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(*model.TechnicalError)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockTransactionPersisterMockRecorder) Add(trx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockTransactionPersister)(nil).Add), trx)
}

// CountByPartner mocks base method.
func (m *MockTransactionPersister) CountByPartner(inp *model.SearchRequest) (*int, *model.TechnicalError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByPartner", inp)
	ret0, _ := ret[0].(*int)
	ret1, _ := ret[1].(*model.TechnicalError)
	return ret0, ret1
}

// CountByPartner indicates an expected call of CountByPartner.
func (mr *MockTransactionPersisterMockRecorder) CountByPartner(inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByPartner", reflect.TypeOf((*MockTransactionPersister)(nil).CountByPartner), inp)
}

// DetailByPartner mocks base method.
func (m *MockTransactionPersister) DetailByPartner(inp *model.FindByIdRequest) (*model.PartnerTransactionProjection, *model.TechnicalError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetailByPartner", inp)
	ret0, _ := ret[0].(*model.PartnerTransactionProjection)
	ret1, _ := ret[1].(*model.TechnicalError)
	return ret0, ret1
}

// DetailByPartner indicates an expected call of DetailByPartner.
func (mr *MockTransactionPersisterMockRecorder) DetailByPartner(inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetailByPartner", reflect.TypeOf((*MockTransactionPersister)(nil).DetailByPartner), inp)
}

// SearchByPartner mocks base method.
func (m *MockTransactionPersister) SearchByPartner(inp *model.SearchRequest) ([]model.PartnerTransactionProjection, *model.TechnicalError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByPartner", inp)
	ret0, _ := ret[0].([]model.PartnerTransactionProjection)
	ret1, _ := ret[1].(*model.TechnicalError)
	return ret0, ret1
}

// SearchByPartner indicates an expected call of SearchByPartner.
func (mr *MockTransactionPersisterMockRecorder) SearchByPartner(inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByPartner", reflect.TypeOf((*MockTransactionPersister)(nil).SearchByPartner), inp)
}
