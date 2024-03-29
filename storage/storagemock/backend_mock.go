// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-utils/storage (interfaces: Backend)

// Package storagemock is a generated GoMock package.
package storagemock

import (
	context "context"
	io "io"
	reflect "reflect"

	types "github.com/Scalingo/go-utils/storage/types"
	gomock "github.com/golang/mock/gomock"
)

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockBackend) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockBackendMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBackend)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockBackend) Get(arg0 context.Context, arg1 string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockBackendMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockBackend)(nil).Get), arg0, arg1)
}

// Info mocks base method.
func (m *MockBackend) Info(arg0 context.Context, arg1 string) (types.Info, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info", arg0, arg1)
	ret0, _ := ret[0].(types.Info)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Info indicates an expected call of Info.
func (mr *MockBackendMockRecorder) Info(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockBackend)(nil).Info), arg0, arg1)
}

// List mocks base method.
func (m *MockBackend) List(arg0 context.Context, arg1 string, arg2 types.ListOpts) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockBackendMockRecorder) List(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockBackend)(nil).List), arg0, arg1, arg2)
}

// Move mocks base method.
func (m *MockBackend) Move(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Move", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Move indicates an expected call of Move.
func (mr *MockBackendMockRecorder) Move(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Move", reflect.TypeOf((*MockBackend)(nil).Move), arg0, arg1, arg2)
}

// Size mocks base method.
func (m *MockBackend) Size(arg0 context.Context, arg1 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Size indicates an expected call of Size.
func (mr *MockBackendMockRecorder) Size(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockBackend)(nil).Size), arg0, arg1)
}

// Upload mocks base method.
func (m *MockBackend) Upload(arg0 context.Context, arg1 io.Reader, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upload indicates an expected call of Upload.
func (mr *MockBackendMockRecorder) Upload(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockBackend)(nil).Upload), arg0, arg1, arg2)
}
