package api_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	session_pb "github.com/ZergsLaw/back-template1/api/session/v1"
	user_status_pb "github.com/ZergsLaw/back-template1/api/user_status/v1"
	"github.com/ZergsLaw/back-template1/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

func TestApi_Save(t *testing.T) {
	t.Parallel()

	var (
		token  = "token"
		userID = uuid.Must(uuid.NewV4()).String()
		st     = dom.UserStatusDefault

		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.NewSession: %s", errAny))
	)

	testCases := map[string]struct {
		userID   string
		ip       string
		appToken *app.Token
		want     *session_pb.SaveResponse
		appErr   error
		wantErr  error
	}{
		"success":          {userID, origin.IP.String(), &app.Token{Value: token}, &session_pb.SaveResponse{Token: token}, nil, nil},
		"a.app.NewSession": {userID, origin.IP.String(), nil, nil, errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t)

			mockApp.EXPECT().NewSession(gomock.Any(), uuid.Must(uuid.FromString(tc.userID)), st, origin).Return(tc.appToken, tc.appErr)

			res, err := c.Save(ctx, &session_pb.SaveRequest{
				UserId:    tc.userID,
				Ip:        tc.ip,
				UserAgent: origin.UserAgent,
				Kind:      user_status_pb.StatusKind_STATUS_KIND_DEFAULT,
			})
			assert.ErrorIs(err, tc.wantErr)
			assert.True(proto.Equal(tc.want, res))
		})
	}
}

func TestApi_Get(t *testing.T) {
	t.Parallel()

	const token = `token`

	var (
		session = &app.Session{
			ID:     uuid.Must(uuid.NewV4()),
			Origin: origin,
			Token: app.Token{
				Value: token,
			},
			UserID:    uuid.Must(uuid.NewV4()),
			Status:    dom.UserStatusAdmin,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.Session: %s", errAny))
	)

	testCases := map[string]struct {
		token      string
		appSession *app.Session
		want       *session_pb.GetResponse
		appErr     error
		wantErr    error
	}{
		"success":       {token, session, &session_pb.GetResponse{SessionId: session.ID.String(), UserId: session.UserID.String(), Kind: user_status_pb.StatusKind_STATUS_KIND_ADMIN}, nil, nil},
		"a.app.Session": {token, nil, nil, errAny, errInternal},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t)

			mockApp.EXPECT().Session(gomock.Any(), tc.token).Return(tc.appSession, tc.appErr)

			res, err := c.Get(ctx, &session_pb.GetRequest{
				Token: tc.token,
			})
			assert.ErrorIs(err, tc.wantErr)
			assert.True(proto.Equal(tc.want, res))
		})
	}
}

func TestApi_Delete(t *testing.T) {
	t.Parallel()

	var (
		sessionID   = uuid.Must(uuid.NewV4()).String()
		errInternal = status.Error(codes.Internal, fmt.Sprintf("a.app.RemoveSession: %s", errAny))
		errNotFound = status.Error(codes.NotFound, fmt.Sprintf("a.app.RemoveSession: %s", app.ErrNotFound))
	)

	testCases := map[string]struct {
		sessionID string
		appErr    error
		wantErr   error
	}{
		"success":              {sessionID, nil, nil},
		"a.app.RemoveSession":  {sessionID, errAny, errInternal},
		"validation_not_found": {sessionID, app.ErrNotFound, errNotFound},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, c, mockApp, assert := start(t)

			mockApp.EXPECT().RemoveSession(gomock.Any(), uuid.Must(uuid.FromString(tc.sessionID))).Return(tc.appErr)

			_, err := c.Delete(ctx, &session_pb.DeleteRequest{
				SessionId: tc.sessionID,
			})
			assert.ErrorIs(err, tc.wantErr)
		})
	}
}
