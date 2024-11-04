package app_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template1/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

func TestApp_NewSession(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		sessionSaveErr error
		authTokenErr   error
		generateToken  *app.Token
		want           *app.Token
		wantErr        error
	}{
		"success":        {nil, nil, &app.Token{Value: "token"}, &app.Token{Value: "token"}, nil},
		"m.session.Save": {errAny, nil, &app.Token{Value: "token"}, nil, errAny},
		"m.auth.Token":   {nil, errAny, nil, nil, errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)

			userID := uuid.Must(uuid.NewV4())
			sessionID := uuid.Must(uuid.NewV4())

			mocks.id.EXPECT().New().Return(sessionID)
			mocks.auth.EXPECT().Token(sessionID).Return(tc.generateToken, tc.authTokenErr)
			if tc.authTokenErr == nil {
				mocks.repo.EXPECT().Save(ctx, app.Session{
					ID:     sessionID,
					Origin: origin,
					Token: app.Token{
						Value: tc.generateToken.Value,
					},
					UserID: userID,
					Status: dom.UserStatusDefault,
				}).Return(tc.sessionSaveErr)
			}

			resToken, err := module.NewSession(ctx, userID, dom.UserStatusDefault, origin)
			assert.Equal(tc.want, resToken)
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestApp_Session(t *testing.T) {
	t.Parallel()

	var (
		session = &app.Session{
			ID:     uuid.Must(uuid.NewV4()),
			Origin: origin,
			Token: app.Token{
				Value: uuid.Must(uuid.NewV4()).String(),
			},
			UserID:    uuid.Must(uuid.NewV4()),
			Status:    dom.UserStatusAdmin,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	testCases := map[string]struct {
		sessionByIDErr error
		authSubjectErr error
		want           *app.Session
		wantErr        error
	}{
		"success":        {nil, nil, session, nil},
		"m.session.ByID": {errAny, nil, nil, errAny},
		"m.auth.Subject": {nil, errAny, nil, errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)

			subject := uuid.Nil
			if tc.authSubjectErr == nil {
				subject = uuid.Must(uuid.NewV4())
			}

			mocks.auth.EXPECT().Subject("token").Return(subject, tc.authSubjectErr)
			if tc.authSubjectErr == nil {
				mocks.repo.EXPECT().ByID(ctx, subject).Return(tc.want, tc.sessionByIDErr)
			}

			session, err := module.Session(ctx, "token")
			assert.Equal(tc.want, session)
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestApp_RemoveSession(t *testing.T) {
	t.Parallel()

	var (
		session = &app.Session{
			ID:     uuid.Must(uuid.NewV4()),
			Origin: origin,
			Token: app.Token{
				Value: uuid.Must(uuid.NewV4()).String(),
			},
			UserID:    uuid.Must(uuid.NewV4()),
			Status:    dom.UserStatusFreeze,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	testCases := map[string]struct {
		sessionByIDErr   error
		sessionDeleteErr error
		session          *app.Session
		want             error
	}{
		"success":        {nil, nil, session, nil},
		"m.auth.Subject": {nil, errAny, session, errAny},
		"m.session.ByID": {errAny, nil, nil, errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)

			id := uuid.Must(uuid.NewV4())
			mocks.repo.EXPECT().ByID(ctx, id).Return(tc.session, tc.sessionByIDErr)
			if tc.sessionByIDErr == nil {
				mocks.repo.EXPECT().Delete(ctx, tc.session.ID).Return(tc.sessionDeleteErr)
			}

			err := module.RemoveSession(ctx, id)
			assert.ErrorIs(err, tc.want)
		})
	}
}
