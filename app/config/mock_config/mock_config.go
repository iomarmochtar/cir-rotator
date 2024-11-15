// Code generated by MockGen. DO NOT EDIT.
// Source: config.go
//
// Generated by this command:
//
//	mockgen -destination mock_config/mock_config.go -source config.go IConfig
//

// Package mock_config is a generated GoMock package.
package mock_config

import (
	reflect "reflect"

	filter "github.com/iomarmochtar/cir-rotator/pkg/filter"
	http "github.com/iomarmochtar/cir-rotator/pkg/http"
	registry "github.com/iomarmochtar/cir-rotator/pkg/registry"
	gomock "go.uber.org/mock/gomock"
)

// MockIConfig is a mock of IConfig interface.
type MockIConfig struct {
	ctrl     *gomock.Controller
	recorder *MockIConfigMockRecorder
	isgomock struct{}
}

// MockIConfigMockRecorder is the mock recorder for MockIConfig.
type MockIConfigMockRecorder struct {
	mock *MockIConfig
}

// NewMockIConfig creates a new mock instance.
func NewMockIConfig(ctrl *gomock.Controller) *MockIConfig {
	mock := &MockIConfig{ctrl: ctrl}
	mock.recorder = &MockIConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIConfig) EXPECT() *MockIConfigMockRecorder {
	return m.recorder
}

// ExcludeEngine mocks base method.
func (m *MockIConfig) ExcludeEngine() filter.IFilterEngine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExcludeEngine")
	ret0, _ := ret[0].(filter.IFilterEngine)
	return ret0
}

// ExcludeEngine indicates an expected call of ExcludeEngine.
func (mr *MockIConfigMockRecorder) ExcludeEngine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExcludeEngine", reflect.TypeOf((*MockIConfig)(nil).ExcludeEngine))
}

// HTTPClient mocks base method.
func (m *MockIConfig) HTTPClient() http.IHttpClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HTTPClient")
	ret0, _ := ret[0].(http.IHttpClient)
	return ret0
}

// HTTPClient indicates an expected call of HTTPClient.
func (mr *MockIConfigMockRecorder) HTTPClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPClient", reflect.TypeOf((*MockIConfig)(nil).HTTPClient))
}

// HTTPWorkerCount mocks base method.
func (m *MockIConfig) HTTPWorkerCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HTTPWorkerCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// HTTPWorkerCount indicates an expected call of HTTPWorkerCount.
func (mr *MockIConfigMockRecorder) HTTPWorkerCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPWorkerCount", reflect.TypeOf((*MockIConfig)(nil).HTTPWorkerCount))
}

// Host mocks base method.
func (m *MockIConfig) Host() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Host")
	ret0, _ := ret[0].(string)
	return ret0
}

// Host indicates an expected call of Host.
func (mr *MockIConfigMockRecorder) Host() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Host", reflect.TypeOf((*MockIConfig)(nil).Host))
}

// ImageRegistry mocks base method.
func (m *MockIConfig) ImageRegistry() registry.ImageRegistry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageRegistry")
	ret0, _ := ret[0].(registry.ImageRegistry)
	return ret0
}

// ImageRegistry indicates an expected call of ImageRegistry.
func (mr *MockIConfigMockRecorder) ImageRegistry() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageRegistry", reflect.TypeOf((*MockIConfig)(nil).ImageRegistry))
}

// IncludeEngine mocks base method.
func (m *MockIConfig) IncludeEngine() filter.IFilterEngine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncludeEngine")
	ret0, _ := ret[0].(filter.IFilterEngine)
	return ret0
}

// IncludeEngine indicates an expected call of IncludeEngine.
func (mr *MockIConfigMockRecorder) IncludeEngine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncludeEngine", reflect.TypeOf((*MockIConfig)(nil).IncludeEngine))
}

// Init mocks base method.
func (m *MockIConfig) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockIConfigMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockIConfig)(nil).Init))
}

// IsDryRun mocks base method.
func (m *MockIConfig) IsDryRun() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDryRun")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDryRun indicates an expected call of IsDryRun.
func (mr *MockIConfigMockRecorder) IsDryRun() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDryRun", reflect.TypeOf((*MockIConfig)(nil).IsDryRun))
}

// Password mocks base method.
func (m *MockIConfig) Password() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Password")
	ret0, _ := ret[0].(string)
	return ret0
}

// Password indicates an expected call of Password.
func (mr *MockIConfigMockRecorder) Password() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Password", reflect.TypeOf((*MockIConfig)(nil).Password))
}

// RepositoryList mocks base method.
func (m *MockIConfig) RepositoryList() []registry.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RepositoryList")
	ret0, _ := ret[0].([]registry.Repository)
	return ret0
}

// RepositoryList indicates an expected call of RepositoryList.
func (mr *MockIConfigMockRecorder) RepositoryList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RepositoryList", reflect.TypeOf((*MockIConfig)(nil).RepositoryList))
}

// SkipDeletionErr mocks base method.
func (m *MockIConfig) SkipDeletionErr() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SkipDeletionErr")
	ret0, _ := ret[0].(bool)
	return ret0
}

// SkipDeletionErr indicates an expected call of SkipDeletionErr.
func (mr *MockIConfigMockRecorder) SkipDeletionErr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SkipDeletionErr", reflect.TypeOf((*MockIConfig)(nil).SkipDeletionErr))
}

// SkipList mocks base method.
func (m *MockIConfig) SkipList() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SkipList")
	ret0, _ := ret[0].([]string)
	return ret0
}

// SkipList indicates an expected call of SkipList.
func (mr *MockIConfigMockRecorder) SkipList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SkipList", reflect.TypeOf((*MockIConfig)(nil).SkipList))
}

// Username mocks base method.
func (m *MockIConfig) Username() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Username")
	ret0, _ := ret[0].(string)
	return ret0
}

// Username indicates an expected call of Username.
func (mr *MockIConfigMockRecorder) Username() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Username", reflect.TypeOf((*MockIConfig)(nil).Username))
}
