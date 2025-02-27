// Code generated by MockGen. DO NOT EDIT.
// Source: go.opentelemetry.io/otel/metric (interfaces: Int64ObservableUpDownCounter)

// Package otelmock is a generated GoMock package.
package otelmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInt64ObservableUpDownCounter is a mock of Int64ObservableUpDownCounter interface.
type MockInt64ObservableUpDownCounter struct {
	ctrl     *gomock.Controller
	recorder *MockInt64ObservableUpDownCounterMockRecorder
}

// MockInt64ObservableUpDownCounterMockRecorder is the mock recorder for MockInt64ObservableUpDownCounter.
type MockInt64ObservableUpDownCounterMockRecorder struct {
	mock *MockInt64ObservableUpDownCounter
}

// NewMockInt64ObservableUpDownCounter creates a new mock instance.
func NewMockInt64ObservableUpDownCounter(ctrl *gomock.Controller) *MockInt64ObservableUpDownCounter {
	mock := &MockInt64ObservableUpDownCounter{ctrl: ctrl}
	mock.recorder = &MockInt64ObservableUpDownCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInt64ObservableUpDownCounter) EXPECT() *MockInt64ObservableUpDownCounterMockRecorder {
	return m.recorder
}

// int64Observable mocks base method.
func (m *MockInt64ObservableUpDownCounter) int64Observable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "int64Observable")
}

// int64Observable indicates an expected call of int64Observable.
func (mr *MockInt64ObservableUpDownCounterMockRecorder) int64Observable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "int64Observable", reflect.TypeOf((*MockInt64ObservableUpDownCounter)(nil).int64Observable))
}

// int64ObservableUpDownCounter mocks base method.
func (m *MockInt64ObservableUpDownCounter) int64ObservableUpDownCounter() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "int64ObservableUpDownCounter")
}

// int64ObservableUpDownCounter indicates an expected call of int64ObservableUpDownCounter.
func (mr *MockInt64ObservableUpDownCounterMockRecorder) int64ObservableUpDownCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "int64ObservableUpDownCounter", reflect.TypeOf((*MockInt64ObservableUpDownCounter)(nil).int64ObservableUpDownCounter))
}

// observable mocks base method.
func (m *MockInt64ObservableUpDownCounter) observable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "observable")
}

// observable indicates an expected call of observable.
func (mr *MockInt64ObservableUpDownCounterMockRecorder) observable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "observable", reflect.TypeOf((*MockInt64ObservableUpDownCounter)(nil).observable))
}
