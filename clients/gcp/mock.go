// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package gcp is a generated GoMock package.
package gcp

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	compute "google.golang.org/api/compute/v1"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetProjectByLabels mocks base method
func (m *MockClient) GetProjectByLabels(ctx context.Context, filters []string) ([]*cloudresourcemanager.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectByLabels", ctx, filters)
	ret0, _ := ret[0].([]*cloudresourcemanager.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectByLabels indicates an expected call of GetProjectByLabels
func (mr *MockClientMockRecorder) GetProjectByLabels(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectByLabels", reflect.TypeOf((*MockClient)(nil).GetProjectByLabels), ctx, filters)
}

// GetProjectNetworks mocks base method
func (m *MockClient) GetProjectNetworks(ctx context.Context, projects []*cloudresourcemanager.Project) ([]*compute.Network, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectNetworks", ctx, projects)
	ret0, _ := ret[0].([]*compute.Network)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectNetworks indicates an expected call of GetProjectNetworks
func (mr *MockClientMockRecorder) GetProjectNetworks(ctx, projects interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectNetworks", reflect.TypeOf((*MockClient)(nil).GetProjectNetworks), ctx, projects)
}

// GetProjectSubnetworks mocks base method
func (m *MockClient) GetProjectSubnetworks(ctx context.Context, projects []*cloudresourcemanager.Project) ([]*compute.Subnetwork, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectSubnetworks", ctx, projects)
	ret0, _ := ret[0].([]*compute.Subnetwork)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectSubnetworks indicates an expected call of GetProjectSubnetworks
func (mr *MockClientMockRecorder) GetProjectSubnetworks(ctx, projects interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectSubnetworks", reflect.TypeOf((*MockClient)(nil).GetProjectSubnetworks), ctx, projects)
}