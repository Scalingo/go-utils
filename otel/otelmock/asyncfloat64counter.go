// Code generated by MockGen. DO NOT EDIT.
// Source: go.opentelemetry.io/otel/metric (interfaces: Float64ObservableCounter)

// Package otelmock is a generated GoMock package.
package otelmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"go.opentelemetry.io/otel/metric"
)

// MockFloat64ObservableCounter is a mock of Float64ObservableCounter interface.
type MockFloat64ObservableCounter struct {
	metric.Float64ObservableCounter

	ctrl     *gomock.Controller
	recorder *MockFloat64ObservableCounterMockRecorder
}

// MockFloat64ObservableCounterMockRecorder is the mock recorder for MockFloat64ObservableCounter.
type MockFloat64ObservableCounterMockRecorder struct {
	mock *MockFloat64ObservableCounter
}

// NewMockFloat64ObservableCounter creates a new mock instance.
func NewMockFloat64ObservableCounter(ctrl *gomock.Controller) *MockFloat64ObservableCounter {
	mock := &MockFloat64ObservableCounter{ctrl: ctrl}
	mock.recorder = &MockFloat64ObservableCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFloat64ObservableCounter) EXPECT() *MockFloat64ObservableCounterMockRecorder {
	return m.recorder
}

// float64Observable mocks base method.
func (m *MockFloat64ObservableCounter) float64Observable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "float64Observable")
}

// float64Observable indicates an expected call of float64Observable.
func (mr *MockFloat64ObservableCounterMockRecorder) float64Observable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "float64Observable", reflect.TypeOf((*MockFloat64ObservableCounter)(nil).float64Observable))
}

// float64ObservableCounter mocks base method.
func (m *MockFloat64ObservableCounter) float64ObservableCounter() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "float64ObservableCounter")
}

// float64ObservableCounter indicates an expected call of float64ObservableCounter.
func (mr *MockFloat64ObservableCounterMockRecorder) float64ObservableCounter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "float64ObservableCounter", reflect.TypeOf((*MockFloat64ObservableCounter)(nil).float64ObservableCounter))
}

// observable mocks base method.
func (m *MockFloat64ObservableCounter) observable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "observable")
}

// observable indicates an expected call of observable.
func (mr *MockFloat64ObservableCounterMockRecorder) observable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "observable", reflect.TypeOf((*MockFloat64ObservableCounter)(nil).observable))
}
