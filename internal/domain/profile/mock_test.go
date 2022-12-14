// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package profile_test is a generated GoMock package.
package profile_test

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	profile "github.com/maypok86/conduit/internal/domain/profile"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CheckFollowing mocks base method.
func (m *MockRepository) CheckFollowing(ctx context.Context, followeeID, followerID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckFollowing", ctx, followeeID, followerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckFollowing indicates an expected call of CheckFollowing.
func (mr *MockRepositoryMockRecorder) CheckFollowing(ctx, followeeID, followerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckFollowing", reflect.TypeOf((*MockRepository)(nil).CheckFollowing), ctx, followeeID, followerID)
}

// Follow mocks base method.
func (m *MockRepository) Follow(ctx context.Context, followeeID, followerID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Follow", ctx, followeeID, followerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Follow indicates an expected call of Follow.
func (mr *MockRepositoryMockRecorder) Follow(ctx, followeeID, followerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Follow", reflect.TypeOf((*MockRepository)(nil).Follow), ctx, followeeID, followerID)
}

// GetByEmail mocks base method.
func (m *MockRepository) GetByEmail(ctx context.Context, email string) (profile.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	ret0, _ := ret[0].(profile.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockRepositoryMockRecorder) GetByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockRepository)(nil).GetByEmail), ctx, email)
}

// GetByUsername mocks base method.
func (m *MockRepository) GetByUsername(ctx context.Context, username string) (profile.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", ctx, username)
	ret0, _ := ret[0].(profile.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsername indicates an expected call of GetByUsername.
func (mr *MockRepositoryMockRecorder) GetByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockRepository)(nil).GetByUsername), ctx, username)
}

// Unfollow mocks base method.
func (m *MockRepository) Unfollow(ctx context.Context, followeeID, followerID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unfollow", ctx, followeeID, followerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unfollow indicates an expected call of Unfollow.
func (mr *MockRepositoryMockRecorder) Unfollow(ctx, followeeID, followerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unfollow", reflect.TypeOf((*MockRepository)(nil).Unfollow), ctx, followeeID, followerID)
}
