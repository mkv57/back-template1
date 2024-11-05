package session_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template1/cmd/session/client"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

func TestClient_Save(t *testing.T) {
	t.Parallel()

	var (
		token = &dom.Token{
			Value: "token",
		}

		userID = uuid.Must(uuid.NewV4())
	)

	testCases := map[string]struct {
		clientErr error
		want      *dom.Token
		wantErr   error
	}{
		"success":        {nil, token, nil},
		"c.session.Save": {client.ErrInternal, nil, errCustom},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, svc, mock, assert := start(t)

			var wantReturn *client.Token
			if tc.want != nil {
				wantReturn = &client.Token{
					Value: tc.want.Value,
				}
			}
			mock.EXPECT().Save(ctx, userID, origin, dom.UserStatusDefault).Return(wantReturn, tc.clientErr)

			res, err := svc.Save(ctx, userID, origin, dom.UserStatusDefault)
			assert.Equal(tc.want, res)
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	var (
		session = &dom.Session{
			ID:     uuid.Must(uuid.NewV4()),
			UserID: uuid.Must(uuid.NewV4()),
			Status: dom.UserStatusDefault,
		}

		token = "token"
	)

	testCases := map[string]struct {
		clientErr error
		want      *dom.Session
		wantErr   error
	}{
		"success":       {nil, session, nil},
		"c.session.Get": {client.ErrInternal, nil, errCustom},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, svc, mock, assert := start(t)
			var wantReturn *client.Session
			if tc.want != nil {
				wantReturn = &client.Session{
					ID:     tc.want.ID,
					UserID: tc.want.UserID,
					Status: tc.want.Status,
				}
			}

			mock.EXPECT().Get(ctx, token).Return(wantReturn, tc.clientErr)

			res, err := svc.Get(ctx, token)
			assert.Equal(tc.want, res)
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestClient_Delete(t *testing.T) {
	t.Parallel()

	sessionID := uuid.Must(uuid.NewV4())

	testCases := map[string]struct {
		clientErr error
		want      error
	}{
		"success":          {nil, nil},
		"c.session.Delete": {client.ErrInternal, errCustom},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, svc, mock, assert := start(t)

			mock.EXPECT().Delete(ctx, sessionID).Return(tc.clientErr)

			err := svc.Delete(ctx, sessionID)
			assert.ErrorIs(err, tc.want)
		})
	}
}
