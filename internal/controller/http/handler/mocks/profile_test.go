// Code generated by MockGen. DO NOT EDIT.
// Source: profile.go

// Package handler_test is a generated GoMock package.
package handler_test

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	profile "github.com/maypok86/conduit/internal/domain/profile"
)

// MockProfileService is a mock of ProfileService interface.
type MockProfileService struct {
	ctrl     *gomock.Controller
	recorder *MockProfileServiceMockRecorder
}

// MockProfileServiceMockRecorder is the mock recorder for MockProfileService.
type MockProfileServiceMockRecorder struct {
	mock *MockProfileService
}

// NewMockProfileService creates a new mock instance.
func NewMockProfileService(ctrl *gomock.Controller) *MockProfileService {
	mock := &MockProfileService{ctrl: ctrl}
	mock.recorder = &MockProfileServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfileService) EXPECT() *MockProfileServiceMockRecorder {
	return m.recorder
}

// Follow mocks base method.
func (m *MockProfileService) Follow(ctx context.Context, email, username string) (profile.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Follow", ctx, email, username)
	ret0, _ := ret[0].(profile.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Follow indicates an expected call of Follow.
func (mr *MockProfileServiceMockRecorder) Follow(ctx, email, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Follow", reflect.TypeOf((*MockProfileService)(nil).Follow), ctx, email, username)
}

// GetByUsername mocks base method.
func (m *MockProfileService) GetByUsername(ctx context.Context, username string) (profile.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", ctx, username)
	ret0, _ := ret[0].(profile.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsername indicates an expected call of GetByUsername.
func (mr *MockProfileServiceMockRecorder) GetByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockProfileService)(nil).GetByUsername), ctx, username)
}

// GetWithFollow mocks base method.
func (m *MockProfileService) GetWithFollow(ctx context.Context, email, username string) (profile.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithFollow", ctx, email, username)
	ret0, _ := ret[0].(profile.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithFollow indicates an expected call of GetWithFollow.
func (mr *MockProfileServiceMockRecorder) GetWithFollow(ctx, email, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithFollow", reflect.TypeOf((*MockProfileService)(nil).GetWithFollow), ctx, email, username)
}

// Unfollow mocks base method.
func (m *MockProfileService) Unfollow(ctx context.Context, email, username string) (profile.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unfollow", ctx, email, username)
	ret0, _ := ret[0].(profile.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unfollow indicates an expected call of Unfollow.
func (mr *MockProfileServiceMockRecorder) Unfollow(ctx, email, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unfollow", reflect.TypeOf((*MockProfileService)(nil).Unfollow), ctx, email, username)
}
