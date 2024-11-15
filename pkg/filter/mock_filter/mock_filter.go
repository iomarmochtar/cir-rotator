// Code generated by MockGen. DO NOT EDIT.
// Source: filter.go
//
// Generated by this command:
//
//	mockgen -destination mock_filter/mock_filter.go -source filter.go IFilterEngine
//

// Package mock_filter is a generated GoMock package.
package mock_filter

import (
	reflect "reflect"

	filter "github.com/iomarmochtar/cir-rotator/pkg/filter"
	gomock "go.uber.org/mock/gomock"
)

// MockIFilterEngine is a mock of IFilterEngine interface.
type MockIFilterEngine struct {
	ctrl     *gomock.Controller
	recorder *MockIFilterEngineMockRecorder
	isgomock struct{}
}

// MockIFilterEngineMockRecorder is the mock recorder for MockIFilterEngine.
type MockIFilterEngineMockRecorder struct {
	mock *MockIFilterEngine
}

// NewMockIFilterEngine creates a new mock instance.
func NewMockIFilterEngine(ctrl *gomock.Controller) *MockIFilterEngine {
	mock := &MockIFilterEngine{ctrl: ctrl}
	mock.recorder = &MockIFilterEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIFilterEngine) EXPECT() *MockIFilterEngineMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockIFilterEngine) Process(fields filter.Fields) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", fields)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Process indicates an expected call of Process.
func (mr *MockIFilterEngineMockRecorder) Process(fields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockIFilterEngine)(nil).Process), fields)
}
