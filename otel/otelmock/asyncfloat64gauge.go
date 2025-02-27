// Code generated by MockGen. DO NOT EDIT.
// Source: go.opentelemetry.io/otel/metric (interfaces: Float64ObservableGauge)

// Package otelmock is a generated GoMock package.
package otelmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFloat64ObservableGauge is a mock of Float64ObservableGauge interface.
type MockFloat64ObservableGauge struct {
	ctrl     *gomock.Controller
	recorder *MockFloat64ObservableGaugeMockRecorder
}

// MockFloat64ObservableGaugeMockRecorder is the mock recorder for MockFloat64ObservableGauge.
type MockFloat64ObservableGaugeMockRecorder struct {
	mock *MockFloat64ObservableGauge
}

// NewMockFloat64ObservableGauge creates a new mock instance.
func NewMockFloat64ObservableGauge(ctrl *gomock.Controller) *MockFloat64ObservableGauge {
	mock := &MockFloat64ObservableGauge{ctrl: ctrl}
	mock.recorder = &MockFloat64ObservableGaugeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFloat64ObservableGauge) EXPECT() *MockFloat64ObservableGaugeMockRecorder {
	return m.recorder
}

// float64Observable mocks base method.
func (m *MockFloat64ObservableGauge) float64Observable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "float64Observable")
}

// float64Observable indicates an expected call of float64Observable.
func (mr *MockFloat64ObservableGaugeMockRecorder) float64Observable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "float64Observable", reflect.TypeOf((*MockFloat64ObservableGauge)(nil).float64Observable))
}

// float64ObservableGauge mocks base method.
func (m *MockFloat64ObservableGauge) float64ObservableGauge() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "float64ObservableGauge")
}

// float64ObservableGauge indicates an expected call of float64ObservableGauge.
func (mr *MockFloat64ObservableGaugeMockRecorder) float64ObservableGauge() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "float64ObservableGauge", reflect.TypeOf((*MockFloat64ObservableGauge)(nil).float64ObservableGauge))
}

// observable mocks base method.
func (m *MockFloat64ObservableGauge) observable() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "observable")
}

// observable indicates an expected call of observable.
func (mr *MockFloat64ObservableGaugeMockRecorder) observable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "observable", reflect.TypeOf((*MockFloat64ObservableGauge)(nil).observable))
}
