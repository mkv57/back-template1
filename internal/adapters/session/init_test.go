package session_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ZergsLaw/back-template/cmd/session/client"
	"github.com/ZergsLaw/back-template/internal/adapters/session"
	"github.com/ZergsLaw/back-template/internal/dom"
	"github.com/ZergsLaw/back-template/internal/testhelper"
)

var (
	errCustom = errors.New("custom err") // For checking err converter.

	origin = dom.Origin{
		IP:        net.ParseIP("192.100.10.4"),
		UserAgent: "UserAgent",
	}
)

func convertErr(t *testing.T, assert *require.Assertions) func(error) error {
	t.Helper()

	return func(err error) error {
		t.Helper()

		assert.ErrorIs(err, client.ErrInternal)

		return errCustom
	}
}

func start(t *testing.T) (context.Context, *session.Client, *MocksessionClient, *require.Assertions) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mock := NewMocksessionClient(ctrl)
	assert := require.New(t)

	return testhelper.Context(t), session.New(mock, convertErr(t, assert)), mock, assert
}
