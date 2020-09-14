// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package planner is a generated GoMock package.
package planner

import (
	context "context"
	network "github.com/estafette/estafette-gcp-network-planner/api/network/v1"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// LoadConfig mocks base method
func (m *MockService) LoadConfig(ctx context.Context) (*network.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadConfig", ctx)
	ret0, _ := ret[0].(*network.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadConfig indicates an expected call of LoadConfig
func (mr *MockServiceMockRecorder) LoadConfig(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfig", reflect.TypeOf((*MockService)(nil).LoadConfig), ctx)
}

// Suggest mocks base method
func (m *MockService) Suggest(ctx context.Context, filter string) ([]network.RangeConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Suggest", ctx, filter)
	ret0, _ := ret[0].([]network.RangeConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Suggest indicates an expected call of Suggest
func (mr *MockServiceMockRecorder) Suggest(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Suggest", reflect.TypeOf((*MockService)(nil).Suggest), ctx, filter)
}
