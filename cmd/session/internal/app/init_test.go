package app_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ZergsLaw/back-template1/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template1/internal/testhelper"
)

var (
	errAny = errors.New("any error")
	origin = app.Origin{
		IP:        net.ParseIP("192.100.10.4"),
		UserAgent: "UserAgent",
	}
)

type mocks struct {
	repo  *MockRepo
	id    *MockID
	auth  *MockAuth
	queue *MockQueue
}

func start(t *testing.T) (context.Context, *app.App, *mocks, *require.Assertions) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockRepo := NewMockRepo(ctrl)
	mockID := NewMockID(ctrl)
	mockAuth := NewMockAuth(ctrl)
	mockQueue := NewMockQueue(ctrl)

	module := app.New(mockRepo, mockAuth, mockID, mockQueue)

	mocks := &mocks{
		repo:  mockRepo,
		id:    mockID,
		auth:  mockAuth,
		queue: mockQueue,
	}

	return testhelper.Context(t), module, mocks, require.New(t)
}
