// Code generated by MockGen. DO NOT EDIT.
// Source: linksaja.go

// Package mock_adaptor is a generated GoMock package.
package adaptor

import (
	reflect "reflect"

	model "github.com/adinandradrs/cezbek-engine/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockLinksajaAdapter is a mock of LinksajaAdapter interface.
type MockLinksajaAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockLinksajaAdapterMockRecorder
}

// MockLinksajaAdapterMockRecorder is the mock recorder for MockLinksajaAdapter.
type MockLinksajaAdapterMockRecorder struct {
	mock *MockLinksajaAdapter
}

// NewMockLinksajaAdapter creates a new mock instance.
func NewMockLinksajaAdapter(ctrl *gomock.Controller) *MockLinksajaAdapter {
	mock := &MockLinksajaAdapter{ctrl: ctrl}
	mock.recorder = &MockLinksajaAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLinksajaAdapter) EXPECT() *MockLinksajaAdapterMockRecorder {
	return m.recorder
}

// Authorization mocks base method.
func (m *MockLinksajaAdapter) Authorization() (*model.LinksajaAuthorizationResponse, *model.TechnicalError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorization")
	ret0, _ := ret[0].(*model.LinksajaAuthorizationResponse)
	ret1, _ := ret[1].(*model.TechnicalError)
	return ret0, ret1
}

// Authorization indicates an expected call of Authorization.
func (mr *MockLinksajaAdapterMockRecorder) Authorization() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorization", reflect.TypeOf((*MockLinksajaAdapter)(nil).Authorization))
}

// FundTransfer mocks base method.
func (m *MockLinksajaAdapter) FundTransfer(inp *model.LinksajaFundTransferRequest) (*model.LinksajaFundTransferResponse, *model.TechnicalError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FundTransfer", inp)
	ret0, _ := ret[0].(*model.LinksajaFundTransferResponse)
	ret1, _ := ret[1].(*model.TechnicalError)
	return ret0, ret1
}

// FundTransfer indicates an expected call of FundTransfer.
func (mr *MockLinksajaAdapterMockRecorder) FundTransfer(inp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FundTransfer", reflect.TypeOf((*MockLinksajaAdapter)(nil).FundTransfer), inp)
}