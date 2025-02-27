// Code generated by MockGen. DO NOT EDIT.
// Source: go.opentelemetry.io/otel/metric (interfaces: Int64ObservableGauge)

// Package otelmock is a generated GoMock package.
package otelmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInt64ObservableGauge is a mock of Int64ObservableGauge interface.
type MockInt64ObservableGauge struct {
	ctrl     *gomock.Controller
	recorder *MockInt64ObservableGaugeMockRecorder
}

// MockInt64ObservableGaugeMockRecorder is the mock recorder for MockInt64ObservableGauge.
type MockInt64ObservableGaugeMockRecorder struct {
	mock *MockInt64ObservableGauge
}

// NewMockInt64ObservableGauge creates a new mock instance.
func NewMockInt64ObservableGauge(ctrl *gomock.Controller) *MockInt64ObservableGauge {
	mock := &MockInt64ObservableGauge{ctrl: ctrl}
	mock.recorder = &MockInt64ObservableGaugeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInt64ObservableGauge) EXPECT() *MockInt64ObservableGaugeMockRecorder {
	return m.recorder
}

// int64Observable mocks base method.
func (m *MockInt64ObservableGauge) int64Observable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "int64Observable")
}

// int64Observable indicates an expected call of int64Observable.
func (mr *MockInt64ObservableGaugeMockRecorder) int64Observable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "int64Observable", reflect.TypeOf((*MockInt64ObservableGauge)(nil).int64Observable))
}

// int64ObservableGauge mocks base method.
func (m *MockInt64ObservableGauge) int64ObservableGauge() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "int64ObservableGauge")
}

// int64ObservableGauge indicates an expected call of int64ObservableGauge.
func (mr *MockInt64ObservableGaugeMockRecorder) int64ObservableGauge() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "int64ObservableGauge", reflect.TypeOf((*MockInt64ObservableGauge)(nil).int64ObservableGauge))
}

// observable mocks base method.
func (m *MockInt64ObservableGauge) observable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "observable")
}

// observable indicates an expected call of observable.
func (mr *MockInt64ObservableGaugeMockRecorder) observable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "observable", reflect.TypeOf((*MockInt64ObservableGauge)(nil).observable))
}
