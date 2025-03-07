// Code generated by MockGen. DO NOT EDIT.
// Source: go.opentelemetry.io/otel/metric (interfaces: Int64Counter)

// Package otelmock is a generated GoMock package.
package otelmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	metric "go.opentelemetry.io/otel/metric"
)

// MockInt64Counter is a mock of Int64Counter interface.
type MockInt64Counter struct {
	metric.Int64Counter

	ctrl     *gomock.Controller
	recorder *MockInt64CounterMockRecorder
}

// MockInt64CounterMockRecorder is the mock recorder for MockInt64Counter.
type MockInt64CounterMockRecorder struct {
	mock *MockInt64Counter
}

// NewMockInt64Counter creates a new mock instance.
func NewMockInt64Counter(ctrl *gomock.Controller) *MockInt64Counter {
	mock := &MockInt64Counter{ctrl: ctrl}
	mock.recorder = &MockInt64CounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInt64Counter) EXPECT() *MockInt64CounterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockInt64Counter) Add(arg0 context.Context, arg1 int64, arg2 ...metric.AddOption) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Add", varargs...)
}

// Add indicates an expected call of Add.
func (mr *MockInt64CounterMockRecorder) Add(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockInt64Counter)(nil).Add), varargs...)
}

// int64Counter mocks base method.
func (m *MockInt64Counter) int64Counter() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "int64Counter")
}

// int64Counter indicates an expected call of int64Counter.
func (mr *MockInt64CounterMockRecorder) int64Counter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "int64Counter", reflect.TypeOf((*MockInt64Counter)(nil).int64Counter))
}
