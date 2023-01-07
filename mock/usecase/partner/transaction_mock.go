// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package mock_partner is a generated GoMock package.
package partner

import (
	reflect "reflect"

	model "github.com/adinandradrs/cezbek-engine/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockTransactionProvider is a mock of TransactionProvider interface.
type MockTransactionProvider struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionProviderMockRecorder
}

// MockTransactionProviderMockRecorder is the mock recorder for MockTransactionProvider.
type MockTransactionProviderMockRecorder struct {
	mock *MockTransactionProvider
}

// NewMockTransactionProvider creates a new mock instance.
func NewMockTransactionProvider(ctrl *gomock.Controller) *MockTransactionProvider {
	mock := &MockTransactionProvider{ctrl: ctrl}
	mock.recorder = &MockTransactionProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionProvider) EXPECT() *MockTransactionProviderMockRecorder {
	return m.recorder
}

// Detail mocks base method.
func (m *MockTransactionProvider) Detail(inp *model.FindByIdRequest) (*model.PartnerTransactionProjection, *model.BusinessError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Detail", inp)
	ret0, _ := ret[0].(*model.PartnerTransactionProjection)
	ret1, _ := ret[1].(*model.BusinessError)
	return ret0, ret1
}

// Detail indicates an expected call of Detail.
func (mr *MockTransactionProviderMockRecorder) Detail(inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Detail", reflect.TypeOf((*MockTransactionProvider)(nil).Detail), inp)
}

// Search mocks base method.
func (m *MockTransactionProvider) Search(inp *model.SearchRequest) (*model.PartnerTransactionSearchResponse, *model.BusinessError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", inp)
	ret0, _ := ret[0].(*model.PartnerTransactionSearchResponse)
	ret1, _ := ret[1].(*model.BusinessError)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockTransactionProviderMockRecorder) Search(inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockTransactionProvider)(nil).Search), inp)
}
