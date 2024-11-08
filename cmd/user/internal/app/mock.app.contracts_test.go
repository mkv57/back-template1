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

	app "github.com/ZergsLaw/back-template1/cmd/user/internal/app"
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

// ByEmail mocks base method.
func (m *MockRepo) ByEmail(arg0 context.Context, arg1 string) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByEmail", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByEmail indicates an expected call of ByEmail.
func (mr *MockRepoMockRecorder) ByEmail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByEmail", reflect.TypeOf((*MockRepo)(nil).ByEmail), arg0, arg1)
}

// ByID mocks base method.
func (m *MockRepo) ByID(arg0 context.Context, arg1 uuid.UUID) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByID", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByID indicates an expected call of ByID.
func (mr *MockRepoMockRecorder) ByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByID", reflect.TypeOf((*MockRepo)(nil).ByID), arg0, arg1)
}

// ByUsername mocks base method.
func (m *MockRepo) ByUsername(arg0 context.Context, arg1 string) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByUsername", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByUsername indicates an expected call of ByUsername.
func (mr *MockRepoMockRecorder) ByUsername(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByUsername", reflect.TypeOf((*MockRepo)(nil).ByUsername), arg0, arg1)
}

// DeleteAvatar mocks base method.
func (m *MockRepo) DeleteAvatar(ctx context.Context, userID, fileID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAvatar", ctx, userID, fileID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAvatar indicates an expected call of DeleteAvatar.
func (mr *MockRepoMockRecorder) DeleteAvatar(ctx, userID, fileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAvatar", reflect.TypeOf((*MockRepo)(nil).DeleteAvatar), ctx, userID, fileID)
}

// FinishTask mocks base method.
func (m *MockRepo) FinishTask(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinishTask indicates an expected call of FinishTask.
func (mr *MockRepoMockRecorder) FinishTask(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishTask", reflect.TypeOf((*MockRepo)(nil).FinishTask), arg0, arg1)
}

// GetAvatar mocks base method.
func (m *MockRepo) GetAvatar(ctx context.Context, fileID uuid.UUID) (*app.AvatarInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvatar", ctx, fileID)
	ret0, _ := ret[0].(*app.AvatarInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvatar indicates an expected call of GetAvatar.
func (mr *MockRepoMockRecorder) GetAvatar(ctx, fileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvatar", reflect.TypeOf((*MockRepo)(nil).GetAvatar), ctx, fileID)
}

// GetCountAvatars mocks base method.
func (m *MockRepo) GetCountAvatars(ctx context.Context, ownerID uuid.UUID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountAvatars", ctx, ownerID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountAvatars indicates an expected call of GetCountAvatars.
func (mr *MockRepoMockRecorder) GetCountAvatars(ctx, ownerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountAvatars", reflect.TypeOf((*MockRepo)(nil).GetCountAvatars), ctx, ownerID)
}

// ListActualTask mocks base method.
func (m *MockRepo) ListActualTask(arg0 context.Context, arg1 int) ([]app.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListActualTask", arg0, arg1)
	ret0, _ := ret[0].([]app.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListActualTask indicates an expected call of ListActualTask.
func (mr *MockRepoMockRecorder) ListActualTask(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListActualTask", reflect.TypeOf((*MockRepo)(nil).ListActualTask), arg0, arg1)
}

// ListAvatarByUserID mocks base method.
func (m *MockRepo) ListAvatarByUserID(ctx context.Context, userID uuid.UUID) ([]app.AvatarInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAvatarByUserID", ctx, userID)
	ret0, _ := ret[0].([]app.AvatarInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAvatarByUserID indicates an expected call of ListAvatarByUserID.
func (mr *MockRepoMockRecorder) ListAvatarByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAvatarByUserID", reflect.TypeOf((*MockRepo)(nil).ListAvatarByUserID), ctx, userID)
}

// Save mocks base method.
func (m *MockRepo) Save(arg0 context.Context, arg1 app.User) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockRepoMockRecorder) Save(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepo)(nil).Save), arg0, arg1)
}

// SaveAvatar mocks base method.
func (m *MockRepo) SaveAvatar(ctx context.Context, fileCache app.AvatarInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAvatar", ctx, fileCache)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAvatar indicates an expected call of SaveAvatar.
func (mr *MockRepoMockRecorder) SaveAvatar(ctx, fileCache any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAvatar", reflect.TypeOf((*MockRepo)(nil).SaveAvatar), ctx, fileCache)
}

// SaveTask mocks base method.
func (m *MockRepo) SaveTask(arg0 context.Context, arg1 app.Task) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTask", arg0, arg1)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveTask indicates an expected call of SaveTask.
func (mr *MockRepoMockRecorder) SaveTask(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTask", reflect.TypeOf((*MockRepo)(nil).SaveTask), arg0, arg1)
}

// SearchUsers mocks base method.
func (m *MockRepo) SearchUsers(arg0 context.Context, arg1 app.SearchParams) ([]app.User, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUsers", arg0, arg1)
	ret0, _ := ret[0].([]app.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SearchUsers indicates an expected call of SearchUsers.
func (mr *MockRepoMockRecorder) SearchUsers(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUsers", reflect.TypeOf((*MockRepo)(nil).SearchUsers), arg0, arg1)
}

// Tx mocks base method.
func (m *MockRepo) Tx(ctx context.Context, f func(app.Repo) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// Tx indicates an expected call of Tx.
func (mr *MockRepoMockRecorder) Tx(ctx, f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockRepo)(nil).Tx), ctx, f)
}

// Update mocks base method.
func (m *MockRepo) Update(arg0 context.Context, arg1 app.User) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRepoMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepo)(nil).Update), arg0, arg1)
}

// UsersByIDs mocks base method.
func (m *MockRepo) UsersByIDs(ctx context.Context, ids []uuid.UUID) ([]app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersByIDs", ctx, ids)
	ret0, _ := ret[0].([]app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UsersByIDs indicates an expected call of UsersByIDs.
func (mr *MockRepoMockRecorder) UsersByIDs(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersByIDs", reflect.TypeOf((*MockRepo)(nil).UsersByIDs), ctx, ids)
}

// MockFileInfoRepo is a mock of FileInfoRepo interface.
type MockFileInfoRepo struct {
	ctrl     *gomock.Controller
	recorder *MockFileInfoRepoMockRecorder
}

// MockFileInfoRepoMockRecorder is the mock recorder for MockFileInfoRepo.
type MockFileInfoRepoMockRecorder struct {
	mock *MockFileInfoRepo
}

// NewMockFileInfoRepo creates a new mock instance.
func NewMockFileInfoRepo(ctrl *gomock.Controller) *MockFileInfoRepo {
	mock := &MockFileInfoRepo{ctrl: ctrl}
	mock.recorder = &MockFileInfoRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileInfoRepo) EXPECT() *MockFileInfoRepoMockRecorder {
	return m.recorder
}

// DeleteAvatar mocks base method.
func (m *MockFileInfoRepo) DeleteAvatar(ctx context.Context, userID, fileID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAvatar", ctx, userID, fileID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAvatar indicates an expected call of DeleteAvatar.
func (mr *MockFileInfoRepoMockRecorder) DeleteAvatar(ctx, userID, fileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAvatar", reflect.TypeOf((*MockFileInfoRepo)(nil).DeleteAvatar), ctx, userID, fileID)
}

// GetAvatar mocks base method.
func (m *MockFileInfoRepo) GetAvatar(ctx context.Context, fileID uuid.UUID) (*app.AvatarInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvatar", ctx, fileID)
	ret0, _ := ret[0].(*app.AvatarInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvatar indicates an expected call of GetAvatar.
func (mr *MockFileInfoRepoMockRecorder) GetAvatar(ctx, fileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvatar", reflect.TypeOf((*MockFileInfoRepo)(nil).GetAvatar), ctx, fileID)
}

// GetCountAvatars mocks base method.
func (m *MockFileInfoRepo) GetCountAvatars(ctx context.Context, ownerID uuid.UUID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountAvatars", ctx, ownerID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountAvatars indicates an expected call of GetCountAvatars.
func (mr *MockFileInfoRepoMockRecorder) GetCountAvatars(ctx, ownerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountAvatars", reflect.TypeOf((*MockFileInfoRepo)(nil).GetCountAvatars), ctx, ownerID)
}

// ListAvatarByUserID mocks base method.
func (m *MockFileInfoRepo) ListAvatarByUserID(ctx context.Context, userID uuid.UUID) ([]app.AvatarInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAvatarByUserID", ctx, userID)
	ret0, _ := ret[0].([]app.AvatarInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAvatarByUserID indicates an expected call of ListAvatarByUserID.
func (mr *MockFileInfoRepoMockRecorder) ListAvatarByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAvatarByUserID", reflect.TypeOf((*MockFileInfoRepo)(nil).ListAvatarByUserID), ctx, userID)
}

// SaveAvatar mocks base method.
func (m *MockFileInfoRepo) SaveAvatar(ctx context.Context, fileCache app.AvatarInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAvatar", ctx, fileCache)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAvatar indicates an expected call of SaveAvatar.
func (mr *MockFileInfoRepoMockRecorder) SaveAvatar(ctx, fileCache any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAvatar", reflect.TypeOf((*MockFileInfoRepo)(nil).SaveAvatar), ctx, fileCache)
}

// MockFileStore is a mock of FileStore interface.
type MockFileStore struct {
	ctrl     *gomock.Controller
	recorder *MockFileStoreMockRecorder
}

// MockFileStoreMockRecorder is the mock recorder for MockFileStore.
type MockFileStoreMockRecorder struct {
	mock *MockFileStore
}

// NewMockFileStore creates a new mock instance.
func NewMockFileStore(ctrl *gomock.Controller) *MockFileStore {
	mock := &MockFileStore{ctrl: ctrl}
	mock.recorder = &MockFileStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileStore) EXPECT() *MockFileStoreMockRecorder {
	return m.recorder
}

// DeleteFile mocks base method.
func (m *MockFileStore) DeleteFile(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockFileStoreMockRecorder) DeleteFile(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockFileStore)(nil).DeleteFile), ctx, id)
}

// DownloadFile mocks base method.
func (m *MockFileStore) DownloadFile(ctx context.Context, id uuid.UUID) (*app.Avatar, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFile", ctx, id)
	ret0, _ := ret[0].(*app.Avatar)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadFile indicates an expected call of DownloadFile.
func (mr *MockFileStoreMockRecorder) DownloadFile(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockFileStore)(nil).DownloadFile), ctx, id)
}

// UploadFile mocks base method.
func (m *MockFileStore) UploadFile(ctx context.Context, f app.Avatar) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", ctx, f)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockFileStoreMockRecorder) UploadFile(ctx, f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockFileStore)(nil).UploadFile), ctx, f)
}

// MockPasswordHash is a mock of PasswordHash interface.
type MockPasswordHash struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordHashMockRecorder
}

// MockPasswordHashMockRecorder is the mock recorder for MockPasswordHash.
type MockPasswordHashMockRecorder struct {
	mock *MockPasswordHash
}

// NewMockPasswordHash creates a new mock instance.
func NewMockPasswordHash(ctrl *gomock.Controller) *MockPasswordHash {
	mock := &MockPasswordHash{ctrl: ctrl}
	mock.recorder = &MockPasswordHashMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPasswordHash) EXPECT() *MockPasswordHashMockRecorder {
	return m.recorder
}

// Compare mocks base method.
func (m *MockPasswordHash) Compare(hashedPassword, password []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compare", hashedPassword, password)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Compare indicates an expected call of Compare.
func (mr *MockPasswordHashMockRecorder) Compare(hashedPassword, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compare", reflect.TypeOf((*MockPasswordHash)(nil).Compare), hashedPassword, password)
}

// Hashing mocks base method.
func (m *MockPasswordHash) Hashing(password string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hashing", password)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hashing indicates an expected call of Hashing.
func (mr *MockPasswordHashMockRecorder) Hashing(password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hashing", reflect.TypeOf((*MockPasswordHash)(nil).Hashing), password)
}

// MockTaskRepo is a mock of TaskRepo interface.
type MockTaskRepo struct {
	ctrl     *gomock.Controller
	recorder *MockTaskRepoMockRecorder
}

// MockTaskRepoMockRecorder is the mock recorder for MockTaskRepo.
type MockTaskRepoMockRecorder struct {
	mock *MockTaskRepo
}

// NewMockTaskRepo creates a new mock instance.
func NewMockTaskRepo(ctrl *gomock.Controller) *MockTaskRepo {
	mock := &MockTaskRepo{ctrl: ctrl}
	mock.recorder = &MockTaskRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskRepo) EXPECT() *MockTaskRepoMockRecorder {
	return m.recorder
}

// FinishTask mocks base method.
func (m *MockTaskRepo) FinishTask(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinishTask indicates an expected call of FinishTask.
func (mr *MockTaskRepoMockRecorder) FinishTask(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishTask", reflect.TypeOf((*MockTaskRepo)(nil).FinishTask), arg0, arg1)
}

// ListActualTask mocks base method.
func (m *MockTaskRepo) ListActualTask(arg0 context.Context, arg1 int) ([]app.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListActualTask", arg0, arg1)
	ret0, _ := ret[0].([]app.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListActualTask indicates an expected call of ListActualTask.
func (mr *MockTaskRepoMockRecorder) ListActualTask(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListActualTask", reflect.TypeOf((*MockTaskRepo)(nil).ListActualTask), arg0, arg1)
}

// SaveTask mocks base method.
func (m *MockTaskRepo) SaveTask(arg0 context.Context, arg1 app.Task) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTask", arg0, arg1)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveTask indicates an expected call of SaveTask.
func (mr *MockTaskRepoMockRecorder) SaveTask(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTask", reflect.TypeOf((*MockTaskRepo)(nil).SaveTask), arg0, arg1)
}

// MockSessions is a mock of Sessions interface.
type MockSessions struct {
	ctrl     *gomock.Controller
	recorder *MockSessionsMockRecorder
}

// MockSessionsMockRecorder is the mock recorder for MockSessions.
type MockSessionsMockRecorder struct {
	mock *MockSessions
}

// NewMockSessions creates a new mock instance.
func NewMockSessions(ctrl *gomock.Controller) *MockSessions {
	mock := &MockSessions{ctrl: ctrl}
	mock.recorder = &MockSessionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessions) EXPECT() *MockSessionsMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockSessions) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionsMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessions)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockSessions) Get(arg0 context.Context, arg1 string) (*dom.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*dom.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSessionsMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSessions)(nil).Get), arg0, arg1)
}

// Save mocks base method.
func (m *MockSessions) Save(arg0 context.Context, arg1 uuid.UUID, arg2 dom.Origin, arg3 dom.UserStatus) (*dom.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*dom.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockSessionsMockRecorder) Save(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSessions)(nil).Save), arg0, arg1, arg2, arg3)
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

// AddUser mocks base method.
func (m *MockQueue) AddUser(arg0 context.Context, arg1 uuid.UUID, arg2 app.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockQueueMockRecorder) AddUser(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockQueue)(nil).AddUser), arg0, arg1, arg2)
}

// DeleteUser mocks base method.
func (m *MockQueue) DeleteUser(arg0 context.Context, arg1 uuid.UUID, arg2 app.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockQueueMockRecorder) DeleteUser(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockQueue)(nil).DeleteUser), arg0, arg1, arg2)
}

// UpdateUser mocks base method.
func (m *MockQueue) UpdateUser(arg0 context.Context, arg1 uuid.UUID, arg2 app.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockQueueMockRecorder) UpdateUser(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockQueue)(nil).UpdateUser), arg0, arg1, arg2)
}
