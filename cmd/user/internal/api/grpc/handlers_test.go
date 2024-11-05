package grpc_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	user_pb "github.com/ZergsLaw/back-template1/api/user/v1"
	user_status_pb "github.com/ZergsLaw/back-template1/api/user_status/v1"
	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

func TestApi_VerificationEmail(t *testing.T) {
	t.Parallel()

	var (
		email       = "email@mail.com"
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.VerificationEmail: %s", errAny))
	)

	testCases := map[string]struct {
		email   string
		appErr  error
		wantErr error
	}{
		"success":                 {email, nil, nil},
		"a.app.VerificationEmail": {email, errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			if tc.email == email {
				mockApp.EXPECT().VerificationEmail(gomock.Any(), tc.email).Return(tc.appErr)
			}

			_, err := c.VerificationEmail(ctx, &user_pb.VerificationEmailRequest{
				Email: tc.email,
			})
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestApi_VerificationUsername(t *testing.T) {
	t.Parallel()

	var (
		username    = "username"
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.VerificationUsername: %s", errAny))
	)

	testCases := map[string]struct {
		username string
		appErr   error
		wantErr  error
	}{
		"success":                    {username, nil, nil},
		"a.app.VerificationUsername": {username, errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			if tc.username == username {
				mockApp.EXPECT().VerificationUsername(gomock.Any(), tc.username).Return(tc.appErr)
			}

			_, err := c.VerificationUsername(ctx, &user_pb.VerificationUsernameRequest{
				Username: tc.username,
			})
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestApi_CreateUser(t *testing.T) {
	t.Parallel()

	var (
		want        = &user_pb.CreateUserResponse{Id: userID.String()}
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.CreateUser: %s", errAny))
	)

	testCases := map[string]struct {
		username string
		fullName string
		email    string
		password string
		want     *user_pb.CreateUserResponse
		appRes   uuid.UUID
		appErr   error
		wantErr  error
	}{
		"success":          {username, fullName, email, password, want, userID, nil, nil},
		"a.app.CreateUser": {username, fullName, email, password, nil, uuid.Nil, errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			if tc.appRes != uuid.Nil || tc.appErr != nil {
				mockApp.EXPECT().CreateUser(gomock.Any(), tc.email, tc.username, tc.fullName, tc.password).Return(tc.appRes, tc.appErr)
			}

			res, err := c.CreateUser(ctx, &user_pb.CreateUserRequest{
				Username: tc.username,
				Email:    tc.email,
				FullName: tc.fullName,
				Password: tc.password,
			})
			assert.ErrorIs(err, tc.wantErr)
			assert.True(proto.Equal(res, tc.want))
		})
	}
}

func TestApi_Login(t *testing.T) {
	t.Parallel()

	var (
		token       = &dom.Token{Value: "token"}
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.Login: %s", errAny))
	)

	testCases := map[string]struct {
		email     string
		password  string
		want      string
		appUserID uuid.UUID
		appToken  *dom.Token
		appErr    error
		wantResp  *user_pb.LoginResponse
		wantErr   error
	}{
		"success":     {email, password, token.Value, userID, token, nil, &user_pb.LoginResponse{UserId: userID.String()}, nil},
		"a.app.Login": {email, password, "", uuid.Nil, nil, errAny, nil, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)
			if tc.appToken != nil || tc.appErr != nil {
				mockApp.EXPECT().Login(gomock.Any(), tc.email, tc.password, origin).Return(tc.appUserID, tc.appToken, tc.appErr)
			}

			md := metadata.MD{}

			res, err := c.Login(ctx, &user_pb.LoginRequest{
				Email:    tc.email,
				Password: tc.password,
			}, grpc.Header(&md))
			assert.ErrorIs(err, tc.wantErr)
			assert.Equal(tc.want, strings.Join(md.Get("authorization"), ""))
			assert.True(proto.Equal(res, tc.wantResp))
		})
	}
}

func TestApi_Logout(t *testing.T) {
	t.Parallel()

	var (
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.Logout: %s", errAny))
	)

	testCases := map[string]struct {
		appErr error
		want   error
	}{
		"success":      {nil, nil},
		"a.app.Logout": {errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)
			mockApp.EXPECT().Logout(gomock.Any(), session).Return(tc.appErr)

			_, err := c.Logout(auth(ctx), &user_pb.LogoutRequest{})
			assert.ErrorIs(err, tc.want)
		})
	}
}

func TestApi_GetUser(t *testing.T) {
	t.Parallel()

	var (
		want = &user_pb.GetUserResponse{
			User: &user_pb.User{
				Id:       user.ID.String(),
				Username: user.Name,
				Email:    user.Email,
				AvatarId: user.AvatarID.String(),
				Kind:     user_status_pb.StatusKind_STATUS_KIND_DEFAULT,
			},
		}
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.UserByID: %s", errAny))
	)

	testCases := map[string]struct {
		userID  string
		want    *user_pb.GetUserResponse
		appRes  *app.User
		appErr  error
		wantErr error
	}{
		"success":          {userID.String(), want, &user, nil, nil},
		"success_empty_id": {"", want, &user, nil, nil},
		"a.app.GetUser":    {userID.String(), nil, nil, errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			if tc.appRes != nil || tc.appErr != nil {
				id := uuid.Nil
				if tc.userID != "" {
					id = uuid.Must(uuid.FromString(tc.userID))
				}
				mockApp.EXPECT().UserByID(gomock.Any(), session, id).Return(tc.appRes, tc.appErr)
			}

			res, err := c.GetUser(auth(ctx), &user_pb.GetUserRequest{
				Id: tc.userID,
			})
			assert.ErrorIs(err, tc.wantErr)
			assert.True(proto.Equal(res, tc.want))
		})
	}
}

func TestApi_SearchUser(t *testing.T) {
	t.Parallel()

	var (
		want = &user_pb.SearchUsersResponse{
			Users: []*user_pb.User{
				{Id: user.ID.String(), Username: user.Name, Email: user.Email, AvatarId: user.AvatarID.String(), Kind: user_status_pb.StatusKind_STATUS_KIND_DEFAULT},
			},
			Total: 1,
		}
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.ListUserByUsername: %s", errAny))
	)

	testCases := map[string]struct {
		name        string
		limit       int
		offset      int
		want        *user_pb.SearchUsersResponse
		appRes      []app.User
		appResTotal int
		appErr      error
		wantErr     error
	}{
		"success_pagination_min": {user.Name, 1, 0, want, []app.User{user}, 1, nil, nil},
		"success_pagination_max": {user.Name, 500, 100, want, []app.User{user}, 1, nil, nil},
		"a.app.SearchUsers":      {user.Name, 1, 1, nil, nil, 0, errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			if tc.appRes != nil || tc.appErr != nil {
				mockApp.EXPECT().ListUserByFilters(gomock.Any(), session, app.SearchParams{
					OwnerID:  session.UserID,
					Username: tc.name,
					FullName: tc.name,
					Limit:    uint64(tc.limit),
					Offset:   uint64(tc.offset),
				}).Return(tc.appRes, tc.appResTotal, tc.appErr)
			}

			res, err := c.SearchUsers(auth(ctx), &user_pb.SearchUsersRequest{
				Name:   tc.name,
				Limit:  int32(int64(tc.limit)),
				Offset: int32(tc.offset),
			})
			assert.ErrorIs(err, tc.wantErr)
			assert.True(proto.Equal(res, tc.want))
		})
	}
}

func TestApi_UpdatePassword(t *testing.T) {
	t.Parallel()

	var (
		oldPassword = string(user.PassHash)
		newPassword = string(user.PassHash) + "_new"
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.UpdatePassword: %s", errAny))
	)

	testCases := map[string]struct {
		oldPass string
		newPass string
		appErr  error
		wantErr error
	}{
		"success_pass_len_min":     {oldPassword, newPassword, nil, nil},
		"success_password_len_max": {oldPassword, newPassword, nil, nil},
		"a.app.UpdatePassword":     {oldPassword, newPassword, errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			if tc.newPass == newPassword && tc.oldPass == oldPassword {
				mockApp.EXPECT().
					UpdatePassword(gomock.Any(), session, tc.oldPass, tc.newPass).
					Return(tc.appErr)
			}

			_, err := c.UpdatePassword(auth(ctx), &user_pb.UpdatePasswordRequest{
				Old: tc.oldPass,
				New: tc.newPass,
			})
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}

func TestAPI_UpdateUser(t *testing.T) {
	t.Parallel()

	var (
		newUsername = "new username"
		newAvatarID = uuid.Must(uuid.NewV4())
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.UpdateUser: %s", errAny))
		errNotFound = status.Error(codes.NotFound, fmt.Sprintf("a.app.UpdateUser: %s", app.ErrNotFound))
	)

	testCases := map[string]struct {
		username string
		avatarID uuid.UUID
		appErr   error
		want     *user_pb.UpdateUserResponse
		wantErr  error
	}{
		"success":       {newUsername, newAvatarID, nil, &user_pb.UpdateUserResponse{}, nil},
		"err_not_found": {newUsername, newAvatarID, app.ErrNotFound, nil, errNotFound},
		"err_any":       {newUsername, newAvatarID, errAny, nil, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			if tc.username == newUsername && tc.avatarID == newAvatarID {
				mockApp.EXPECT().UpdateUser(gomock.Any(), session, tc.username, tc.avatarID).Return(tc.appErr)
			}

			res, err := c.UpdateUser(auth(ctx), &user_pb.UpdateUserRequest{
				Username: tc.username,
				AvatarId: tc.avatarID.String(),
			})
			assert.ErrorIs(err, tc.wantErr)
			assert.True(proto.Equal(res, tc.want))
		})
	}
}

func TestAPI_RemoveAvatar(t *testing.T) {
	t.Parallel()

	var (
		newAvatarID = uuid.Must(uuid.NewV4())
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.RemoveAvatar: %s", errAny))
		errNotFound = status.Error(codes.NotFound, fmt.Sprintf("a.app.RemoveAvatar: %s", app.ErrNotFound))
	)

	testCases := map[string]struct {
		avatarID uuid.UUID
		appErr   error
		want     *user_pb.RemoveAvatarResponse
		wantErr  error
	}{
		"success":       {newAvatarID, nil, &user_pb.RemoveAvatarResponse{}, nil},
		"err_not_found": {newAvatarID, app.ErrNotFound, nil, errNotFound},
		"err_any":       {newAvatarID, errAny, nil, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			if tc.avatarID == newAvatarID {
				mockApp.EXPECT().RemoveAvatar(gomock.Any(), session, tc.avatarID).Return(tc.appErr)
			}

			res, err := c.RemoveAvatar(auth(ctx), &user_pb.RemoveAvatarRequest{
				FileId: tc.avatarID.String(),
			})
			assert.ErrorIs(err, tc.wantErr)
			assert.True(proto.Equal(res, tc.want))
		})
	}
}

func TestApi_GetUsersByIDs(t *testing.T) {
	t.Parallel()

	var (
		user1 = &app.User{
			ID:        uuid.Must(uuid.NewV4()),
			Email:     email,
			Name:      username,
			PassHash:  []byte(password),
			AvatarID:  uuid.Must(uuid.NewV4()),
			Status:    dom.UserStatusDefault,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		user2 = &app.User{
			ID:        uuid.Must(uuid.NewV4()),
			Email:     email,
			Name:      username,
			PassHash:  []byte(password),
			AvatarID:  uuid.Must(uuid.NewV4()),
			Status:    dom.UserStatusDefault,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		user3 = &app.User{
			ID:        uuid.Must(uuid.NewV4()),
			Email:     email,
			Name:      username,
			PassHash:  []byte(password),
			AvatarID:  uuid.Must(uuid.NewV4()),
			Status:    dom.UserStatusDefault,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		ids = []string{
			user1.ID.String(),
			user2.ID.String(),
			user3.ID.String(),
		}
		res = []*user_pb.User{
			{Id: user1.ID.String(), Username: user1.Name, Email: user1.Email, AvatarId: user1.AvatarID.String(), Kind: user_status_pb.StatusKind_STATUS_KIND_DEFAULT, FullName: user1.FullName},
			{Id: user2.ID.String(), Username: user2.Name, Email: user2.Email, AvatarId: user2.AvatarID.String(), Kind: user_status_pb.StatusKind_STATUS_KIND_DEFAULT, FullName: user2.FullName},
			{Id: user3.ID.String(), Username: user3.Name, Email: user3.Email, AvatarId: user3.AvatarID.String(), Kind: user_status_pb.StatusKind_STATUS_KIND_DEFAULT, FullName: user3.FullName},
		}

		want = &user_pb.GetUsersByIDsResponse{
			Result: res,
		}
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.GetUsersByIDs: %s", errAny))
	)

	testCases := map[string]struct {
		ids     []string
		appRes  []app.User
		appErr  error
		want    *user_pb.GetUsersByIDsResponse
		wantErr error
	}{
		"success":             {ids, []app.User{*user1, *user2, *user3}, nil, want, nil},
		"a.app.GetUsersByIDs": {ids, []app.User{*user1, *user2, *user3}, errAny, nil, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t, dom.UserStatusDefault)

			iDS := make([]uuid.UUID, len(tc.ids))

			for i := range tc.ids {
				id, err := uuid.FromString(tc.ids[i])
				if err != nil {
					panic(err)
				}

				iDS[i] = id
			}

			mockApp.EXPECT().GetUsersByIDs(gomock.Any(), session, iDS).Return(tc.appRes, tc.appErr)

			resp, err := c.GetUsersByIDs(auth(ctx), &user_pb.GetUsersByIDsRequest{
				Ids: tc.ids,
			})
			assert.ErrorIs(err, tc.wantErr)
			assert.True(proto.Equal(resp, tc.want))
		})
	}
}
