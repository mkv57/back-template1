// Code generated by MockGen. DO NOT EDIT.
// Source: contracts.go
//
// Generated by this command:
//
//	mockgen -source=contracts.go -destination mock.app.contracts_test.go -package app_test
//
// Package app_test is a generated GoMock package.
package app_test

import (
	context "context"
	reflect "reflect"

	app "github.com/ZergsLaw/back-template1/cmd/session/internal/app"
	dom "github.com/ZergsLaw/back-template1/internal/dom"
	uuid "github.com/gofrs/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// ByID mocks base method.
func (m *MockRepo) ByID(arg0 context.Context, arg1 uuid.UUID) (*app.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByID", arg0, arg1)
	ret0, _ := ret[0].(*app.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByID indicates an expected call of ByID.
func (mr *MockRepoMockRecorder) ByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByID", reflect.TypeOf((*MockRepo)(nil).ByID), arg0, arg1)
}

// Delete mocks base method.
func (m *MockRepo) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepoMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepo)(nil).Delete), arg0, arg1)
}

// Save mocks base method.
func (m *MockRepo) Save(arg0 context.Context, arg1 app.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockRepoMockRecorder) Save(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepo)(nil).Save), arg0, arg1)
}

// UpdateStatus mocks base method.
func (m *MockRepo) UpdateStatus(ctx context.Context, reqID, userID uuid.UUID, status dom.UserStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, reqID, userID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockRepoMockRecorder) UpdateStatus(ctx, reqID, userID, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockRepo)(nil).UpdateStatus), ctx, reqID, userID, status)
}

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// Subject mocks base method.
func (m *MockAuth) Subject(token string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subject", token)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subject indicates an expected call of Subject.
func (mr *MockAuthMockRecorder) Subject(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subject", reflect.TypeOf((*MockAuth)(nil).Subject), token)
}

// Token mocks base method.
func (m *MockAuth) Token(arg0 uuid.UUID) (*app.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Token", arg0)
	ret0, _ := ret[0].(*app.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Token indicates an expected call of Token.
func (mr *MockAuthMockRecorder) Token(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Token", reflect.TypeOf((*MockAuth)(nil).Token), arg0)
}

// MockID is a mock of ID interface.
type MockID struct {
	ctrl     *gomock.Controller
	recorder *MockIDMockRecorder
}

// MockIDMockRecorder is the mock recorder for MockID.
type MockIDMockRecorder struct {
	mock *MockID
}

// NewMockID creates a new mock instance.
func NewMockID(ctrl *gomock.Controller) *MockID {
	mock := &MockID{ctrl: ctrl}
	mock.recorder = &MockIDMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockID) EXPECT() *MockIDMockRecorder {
	return m.recorder
}

// New mocks base method.
func (m *MockID) New() uuid.UUID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New")
	ret0, _ := ret[0].(uuid.UUID)
	return ret0
}

// New indicates an expected call of New.
func (mr *MockIDMockRecorder) New() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockID)(nil).New))
}

// MockQueue is a mock of Queue interface.
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue.
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance.
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// UpSessionStatus mocks base method.
func (m *MockQueue) UpSessionStatus() <-chan dom.Event[app.UpdateStatus] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpSessionStatus")
	ret0, _ := ret[0].(<-chan dom.Event[app.UpdateStatus])
	return ret0
}

// UpSessionStatus indicates an expected call of UpSessionStatus.
func (mr *MockQueueMockRecorder) UpSessionStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpSessionStatus", reflect.TypeOf((*MockQueue)(nil).UpSessionStatus))
}
