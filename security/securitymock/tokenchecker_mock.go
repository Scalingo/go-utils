// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Scalingo/go-utils/security (interfaces: TokenChecker)

// Package securitymock is a generated GoMock package.
package securitymock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTokenChecker is a mock of TokenChecker interface.
type MockTokenChecker struct {
	ctrl     *gomock.Controller
	recorder *MockTokenCheckerMockRecorder
}

// MockTokenCheckerMockRecorder is the mock recorder for MockTokenChecker.
type MockTokenCheckerMockRecorder struct {
	mock *MockTokenChecker
}

// NewMockTokenChecker creates a new mock instance.
func NewMockTokenChecker(ctrl *gomock.Controller) *MockTokenChecker {
	mock := &MockTokenChecker{ctrl: ctrl}
	mock.recorder = &MockTokenCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenChecker) EXPECT() *MockTokenCheckerMockRecorder {
	return m.recorder
}

// CheckToken mocks base method.
func (m *MockTokenChecker) CheckToken(arg0 context.Context, arg1, arg2, arg3 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckToken", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckToken indicates an expected call of CheckToken.
func (mr *MockTokenCheckerMockRecorder) CheckToken(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckToken", reflect.TypeOf((*MockTokenChecker)(nil).CheckToken), arg0, arg1, arg2, arg3)
}
