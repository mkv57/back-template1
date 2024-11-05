package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"go.uber.org/mock/gomock"

	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

func TestApp_VerificationEmail(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		repoError error
		want      error
	}{
		"success":                 {app.ErrNotFound, nil},
		"m.user.ByEmail_exist":    {nil, app.ErrEmailExist},
		"m.user.ByEmail_internal": {errAny, errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)
			mocks.repo.EXPECT().ByEmail(ctx, "email@email.com").Return(nil, tc.repoError)

			err := module.VerificationEmail(ctx, "email@email.com")
			assert.ErrorIs(err, tc.want)
		})
	}
}

func TestApp_VerificationUsername(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		repoError error
		want      error
	}{
		"success":                    {app.ErrNotFound, nil},
		"m.user.ByUsername_exist":    {nil, app.ErrUsernameExist},
		"m.user.ByUsername_internal": {errAny, errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)
			mocks.repo.EXPECT().ByUsername(ctx, "name").Return(nil, tc.repoError)

			err := module.VerificationUsername(ctx, "name")
			assert.ErrorIs(err, tc.want)
		})
	}
}

func TestApp_CreateUser(t *testing.T) {
	t.Parallel()

	var ()

	var (
		pass     = `pass`
		email    = `email@email.com`
		fullname = `Andrey_Maslov`
		username = `Andrey`
		passHash = []byte(pass)
		wantID   = uuid.Must(uuid.NewV4())
		status   = dom.UserStatusDefault
		user     = &app.User{
			ID:        wantID,
			Email:     email,
			FullName:  fullname,
			Name:      username,
			PassHash:  passHash,
			AvatarID:  uuid.UUID{},
			Status:    status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	testCases := map[string]struct {
		hasherRes       []byte
		hasherErr       error
		repoSaveRes     uuid.UUID
		repoSaveErr     error
		repoGetRes      *app.User
		repGetErr       error
		repoSaveTaskRes uuid.UUID
		repoSaveTaskErr error
		want            uuid.UUID
		wantErr         error
	}{
		"success":         {passHash, nil, wantID, nil, user, nil, wantID, nil, wantID, nil},
		"m.hash.Hashing":  {nil, errAny, wantID, nil, nil, nil, uuid.Nil, nil, uuid.Nil, errAny},
		"m.user.Save":     {passHash, nil, uuid.Nil, app.ErrUsernameExist, nil, nil, uuid.Nil, nil, uuid.Nil, app.ErrUsernameExist},
		"m.user.SaveTask": {passHash, nil, wantID, nil, user, nil, uuid.Nil, errAny, uuid.Nil, errAny},
		"m.user.ByID":     {passHash, nil, wantID, nil, nil, errAny, uuid.Nil, nil, uuid.Nil, errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)

			mocks.hasher.EXPECT().Hashing(pass).Return(tc.hasherRes, tc.hasherErr)
			if tc.hasherErr == nil {
				mocks.repo.EXPECT().Tx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(r app.Repo) error) error {
					return fn(mocks.repo)
				})

				mocks.repo.EXPECT().Save(ctx, app.User{
					Email:    email,
					Name:     username,
					FullName: fullname,
					PassHash: tc.hasherRes,
					Status:   status,
				}).Return(tc.repoSaveRes, tc.repoSaveErr)

				if tc.repoSaveErr == nil {
					mocks.repo.EXPECT().ByID(ctx, tc.repoSaveRes).
						Return(tc.repoGetRes, tc.repGetErr)
				}

				if tc.repoSaveErr == nil && tc.repGetErr == nil {
					mocks.repo.EXPECT().SaveTask(ctx, app.Task{
						User: *tc.repoGetRes,
						Kind: app.TaskKindEventAdd,
					}).Return(tc.repoSaveTaskRes, tc.repoSaveTaskErr)
				}
			}

			res, err := module.CreateUser(ctx, email, username, fullname, pass)
			assert.Equal(tc.want, res)
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestApp_Login(t *testing.T) {
	t.Parallel()

	var (
		pass  = `pass`
		email = `email@email.com`
		user  = &app.User{
			ID:        uuid.Must(uuid.NewV4()),
			Email:     email,
			Name:      "name",
			FullName:  "Full name",
			PassHash:  []byte(pass),
			Status:    dom.UserStatusDefault,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		token = &dom.Token{
			Value: "token",
		}
	)

	testCases := map[string]struct {
		repoRes          *app.User
		repoErr          error
		authRes          *dom.Token
		authErr          error
		hasherCompareRes bool
		wantUserID       uuid.UUID
		wantToke         *dom.Token
		wantErr          error
	}{
		"success":        {user, nil, token, nil, true, user.ID, token, nil},
		"m.hash.Compare": {user, nil, nil, nil, false, uuid.Nil, nil, app.ErrInvalidPassword},
		"m.user.ByEmail": {nil, app.ErrNotFound, nil, nil, false, uuid.Nil, nil, app.ErrNotFound},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)

			mocks.repo.EXPECT().ByEmail(ctx, email).Return(tc.repoRes, tc.repoErr)
			if tc.repoErr == nil {
				mocks.hasher.EXPECT().Compare(tc.repoRes.PassHash, []byte(pass)).Return(tc.hasherCompareRes)
			}

			if tc.hasherCompareRes {
				mocks.sessions.EXPECT().Save(ctx, user.ID, origin, dom.UserStatusDefault).Return(tc.authRes, tc.authErr)
			}

			userID, token, err := module.Login(ctx, email, pass, origin)
			assert.Equal(tc.wantToke, token)
			assert.Equal(tc.wantUserID, userID)
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestApp_UserByID(t *testing.T) {
	t.Parallel()

	var (
		pass  = `pass`
		email = `email@email.com`
		user  = &app.User{
			ID:        uuid.Must(uuid.NewV4()),
			Email:     email,
			Name:      "name",
			FullName:  "Full name",
			Status:    dom.UserStatusDefault,
			PassHash:  []byte(pass),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		session = &dom.Session{
			ID:     uuid.Must(uuid.NewV4()),
			UserID: user.ID,
		}
		session2 = &dom.Session{
			ID:     uuid.Must(uuid.NewV4()),
			UserID: uuid.Must(uuid.NewV4()),
		}
	)

	testCases := map[string]struct {
		session      *dom.Session
		userIDArg    uuid.UUID
		userIDSearch uuid.UUID
		repoRes      *app.User
		repoErr      error
		want         *app.User
		wantErr      error
	}{
		"success":      {session, user.ID, user.ID, user, nil, user, nil},
		"success_self": {session2, uuid.Nil, session2.UserID, user, nil, user, nil},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)

			mocks.repo.EXPECT().ByID(ctx, tc.userIDSearch).Return(tc.repoRes, tc.repoErr)

			res, err := module.UserByID(ctx, *tc.session, tc.userIDArg)
			assert.Equal(tc.want, res)
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestApp_UpdatePassword(t *testing.T) {
	t.Parallel()

	var (
		user = app.User{
			ID:        uuid.Must(uuid.NewV4()),
			Email:     "email@mail.com",
			Name:      "name",
			FullName:  "Full name",
			PassHash:  []byte("pass"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		session = dom.Session{
			ID:     uuid.Must(uuid.NewV4()),
			UserID: user.ID,
		}
	)

	testCases := map[string]struct {
		repoByIDRes          *app.User
		repoByIDErr          error
		hashCompareResFirst  bool
		hashCompareResSecond bool
		hashHashingRes       []byte
		hashHashingErr       error
		oldPass              string
		newPass              string
		updateRes            *app.User
		updateErr            error
		want                 error
	}{
		"success":               {lo.ToPtr(user), nil, true, false, []byte("password"), nil, "pass", "password", &app.User{}, nil, nil},
		"m.user.ByID":           {nil, app.ErrNotFound, false, true, nil, nil, "pass", "password", &app.User{}, nil, app.ErrNotFound},
		"m.hash.Hashing":        {lo.ToPtr(user), nil, true, false, nil, errAny, "pass", "password", &app.User{}, nil, errAny},
		"m.hash.Compare_second": {lo.ToPtr(user), nil, true, true, nil, nil, "pass", "pass", &app.User{}, nil, app.ErrNotDifferent},
		"m.hash.Compare_first":  {lo.ToPtr(user), nil, false, true, nil, nil, "pass", "password", &app.User{}, nil, app.ErrInvalidPassword},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)

			mocks.repo.EXPECT().ByID(ctx, session.UserID).Return(tc.repoByIDRes, tc.repoByIDErr)

			if tc.repoByIDErr == nil {
				mocks.hasher.EXPECT().Compare(gomock.Any(), []byte(tc.oldPass)).Return(tc.hashCompareResFirst)
			}

			if tc.hashCompareResFirst {
				mocks.hasher.EXPECT().Compare(gomock.Any(), []byte(tc.newPass)).Return(tc.hashCompareResSecond)
			}

			if !tc.hashCompareResSecond {
				mocks.hasher.EXPECT().Hashing(tc.newPass).Return(tc.hashHashingRes, tc.hashHashingErr)
			}

			if tc.hashHashingErr == nil && tc.hashHashingRes != nil {
				tc.repoByIDRes.PassHash = tc.hashHashingRes
				mocks.repo.EXPECT().Update(ctx, *tc.repoByIDRes).Return(tc.updateRes, tc.updateErr)
			}

			err := module.UpdatePassword(ctx, session, tc.oldPass, tc.newPass)
			assert.ErrorIs(err, tc.want)
		})
	}
}

func TestApp_UpdateUser(t *testing.T) {
	t.Parallel()

	var (
		user = &app.User{
			ID:        uuid.Must(uuid.NewV4()),
			Email:     "email@mail.com",
			Name:      "name",
			FullName:  "Full name",
			PassHash:  []byte("pass"),
			Status:    dom.UserStatusDefault,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		session = dom.Session{
			ID:     uuid.Must(uuid.NewV4()),
			UserID: user.ID,
		}
		sessionAnotherUser = dom.Session{
			ID:     uuid.Must(uuid.NewV4()),
			UserID: uuid.Must(uuid.NewV4()),
		}
		newUserName = "new name"
		newAvatarID = uuid.Must(uuid.NewV4())
	)

	testCases := map[string]struct {
		session             dom.Session
		newUserName         string
		newAvatarID         uuid.UUID
		repoByIDRes         *app.User
		repoByIDErr         error
		repoGetFileCacheErr error
		repoUpdateRes       *app.User
		repoUpdateErr       error
		want                error
	}{
		"success":                      {session, newUserName, newAvatarID, user, nil, nil, &app.User{}, nil, nil},
		"success_avatar_id_is_empty":   {session, newUserName, uuid.Nil, user, nil, nil, &app.User{}, nil, nil},
		"err_not_found_by_id":          {sessionAnotherUser, newUserName, newAvatarID, nil, app.ErrNotFound, nil, &app.User{}, nil, app.ErrNotFound},
		"err_not_found_get_file_cache": {sessionAnotherUser, newUserName, newAvatarID, nil, nil, app.ErrNotFound, &app.User{}, nil, app.ErrNotFound},
		"err_any_by_id":                {session, newUserName, newAvatarID, nil, errAny, nil, &app.User{}, nil, errAny},
		"err_any_get_file_cache":       {session, newUserName, newAvatarID, nil, nil, errAny, &app.User{}, nil, errAny},
		"err_any_update":               {session, newUserName, newAvatarID, user, nil, nil, &app.User{}, errAny, errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, module, mocks, assert := start(t)

			mocks.repo.EXPECT().ByID(ctx, tc.session.UserID).Return(tc.repoByIDRes, tc.repoByIDErr)

			if tc.newAvatarID == uuid.Nil {
				tc.newAvatarID = tc.repoByIDRes.AvatarID
			}

			if tc.repoByIDErr == nil && tc.newAvatarID != uuid.Nil {
				mocks.repo.EXPECT().GetAvatar(ctx, tc.newAvatarID).Return(nil, tc.repoGetFileCacheErr)
			}

			if tc.repoByIDErr == nil && tc.repoGetFileCacheErr == nil {
				updateUser := app.User{
					ID:        tc.session.UserID,
					Email:     tc.repoByIDRes.Email,
					FullName:  user.FullName,
					Name:      tc.newUserName,
					PassHash:  tc.repoByIDRes.PassHash,
					AvatarID:  tc.newAvatarID,
					Status:    tc.repoByIDRes.Status,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}
				mocks.repo.EXPECT().Update(ctx, updateUser).Return(tc.repoUpdateRes, tc.repoUpdateErr)
			}

			err := module.UpdateUser(ctx, tc.session, tc.newUserName, tc.newAvatarID)
			assert.ErrorIs(err, tc.want)
		})
	}
}
