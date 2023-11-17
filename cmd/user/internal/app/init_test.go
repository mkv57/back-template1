package app_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ZergsLaw/back-template/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template/internal/dom"
	"github.com/ZergsLaw/back-template/internal/testhelper"
)

const pngFilePath = `./testdata/test.png`

var (
	errAny = errors.New("any error")
	origin = dom.Origin{
		IP:        net.ParseIP("192.100.10.4"),
		UserAgent: "UserAgent",
	}
	ownerID = uuid.Must(uuid.NewV4())
	fileID  = uuid.Must(uuid.NewV4())
)

type mocks struct {
	hasher   *MockPasswordHash
	repo     *MockRepo
	sessions *MockSessions
	file     *MockFileStore
	queue    *MockQueue
}

func start(t *testing.T) (context.Context, *app.App, *mocks, *require.Assertions) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockRepo := NewMockRepo(ctrl)
	mockHasher := NewMockPasswordHash(ctrl)
	mockSession := NewMockSessions(ctrl)
	mockFileStore := NewMockFileStore(ctrl)
	mockQueue := NewMockQueue(ctrl)

	module := app.New(mockRepo, mockHasher, mockSession, mockFileStore, mockQueue)

	mocks := &mocks{
		hasher:   mockHasher,
		repo:     mockRepo,
		sessions: mockSession,
		file:     mockFileStore,
		queue:    mockQueue,
	}

	return testhelper.Context(t), module, mocks, require.New(t)
}
