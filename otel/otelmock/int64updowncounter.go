// Code generated by MockGen. DO NOT EDIT.
// Source: go.opentelemetry.io/otel/metric (interfaces: Int64UpDownCounter)

// Package otelmock is a generated GoMock package.
package otelmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	metric "go.opentelemetry.io/otel/metric"
)

// MockInt64UpDownCounter is a mock of Int64UpDownCounter interface.
type MockInt64UpDownCounter struct {
	ctrl     *gomock.Controller
	recorder *MockInt64UpDownCounterMockRecorder
}

// MockInt64UpDownCounterMockRecorder is the mock recorder for MockInt64UpDownCounter.
type MockInt64UpDownCounterMockRecorder struct {
	mock *MockInt64UpDownCounter
}

// NewMockInt64UpDownCounter creates a new mock instance.
func NewMockInt64UpDownCounter(ctrl *gomock.Controller) *MockInt64UpDownCounter {
	mock := &MockInt64UpDownCounter{ctrl: ctrl}
	mock.recorder = &MockInt64UpDownCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInt64UpDownCounter) EXPECT() *MockInt64UpDownCounterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockInt64UpDownCounter) Add(arg0 context.Context, arg1 int64, arg2 ...metric.AddOption) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Add", varargs...)
}

// Add indicates an expected call of Add.
func (mr *MockInt64UpDownCounterMockRecorder) Add(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockInt64UpDownCounter)(nil).Add), varargs...)
}

// int64UpDownCounter mocks base method.
func (m *MockInt64UpDownCounter) int64UpDownCounter() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "int64UpDownCounter")
}

// int64UpDownCounter indicates an expected call of int64UpDownCounter.
func (mr *MockInt64UpDownCounterMockRecorder) int64UpDownCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "int64UpDownCounter", reflect.TypeOf((*MockInt64UpDownCounter)(nil).int64UpDownCounter))
}
