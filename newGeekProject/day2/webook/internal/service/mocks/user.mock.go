// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/user.go

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	domain "GeekProject/newGeekProject/day2/webook/internal/domain"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserServiceInterface is a mock of UserServiceInterface interface.
type MockUserServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceInterfaceMockRecorder
}

// MockUserServiceInterfaceMockRecorder is the mock recorder for MockUserServiceInterface.
type MockUserServiceInterfaceMockRecorder struct {
	mock *MockUserServiceInterface
}

// NewMockUserServiceInterface creates a new mock instance.
func NewMockUserServiceInterface(ctrl *gomock.Controller) *MockUserServiceInterface {
	mock := &MockUserServiceInterface{ctrl: ctrl}
	mock.recorder = &MockUserServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServiceInterface) EXPECT() *MockUserServiceInterfaceMockRecorder {
	return m.recorder
}

// FindByCreateWechat mocks base method.
func (m *MockUserServiceInterface) FindByCreateWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCreateWechat", ctx, wechatInfo)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCreateWechat indicates an expected call of FindByCreateWechat.
func (mr *MockUserServiceInterfaceMockRecorder) FindByCreateWechat(ctx, wechatInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCreateWechat", reflect.TypeOf((*MockUserServiceInterface)(nil).FindByCreateWechat), ctx, wechatInfo)
}

// FindOrCreate mocks base method.
func (m *MockUserServiceInterface) FindOrCreate(ctx context.Context, phone string) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreate", ctx, phone)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockUserServiceInterfaceMockRecorder) FindOrCreate(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockUserServiceInterface)(nil).FindOrCreate), ctx, phone)
}

// Login mocks base method.
func (m *MockUserServiceInterface) Login(ctx context.Context, domainU domain.UserDomain) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, domainU)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceInterfaceMockRecorder) Login(ctx, domainU interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserServiceInterface)(nil).Login), ctx, domainU)
}

// Profile mocks base method.
func (m *MockUserServiceInterface) Profile(ctx context.Context, id int64) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Profile", ctx, id)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Profile indicates an expected call of Profile.
func (mr *MockUserServiceInterfaceMockRecorder) Profile(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Profile", reflect.TypeOf((*MockUserServiceInterface)(nil).Profile), ctx, id)
}

// SignUp mocks base method.
func (m *MockUserServiceInterface) SignUp(ctx context.Context, domainU domain.UserDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, domainU)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserServiceInterfaceMockRecorder) SignUp(ctx, domainU interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserServiceInterface)(nil).SignUp), ctx, domainU)
}
