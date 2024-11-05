package client_test

import (
	"context"
	"net"
	"testing"

	"github.com/gofrs/uuid"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	session_pb "github.com/ZergsLaw/back-template1/api/session/v1"
	user_status_pb "github.com/ZergsLaw/back-template1/api/user_status/v1"
	"github.com/ZergsLaw/back-template1/cmd/session/client"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

var (
	_ gomock.Matcher = &traceIDMatcher{}
)

type traceIDMatcher struct {
	expect string
}

// Matches for implements gomock.Matcher.
func (r traceIDMatcher) Matches(x interface{}) bool {
	//ctx, ok := x.(context.Context)
	//if !ok {
	//	return false
	//}

	//md, ok := metadata.FromIncomingContext(ctx)
	//if !ok {
	//	return false
	//}

	//traceID := strings.Join(md.Get(trace.Metadata), "")

	//return r.expect == "traceID"
	return true
}

// String for implements gomock.Matcher.
func (r traceIDMatcher) String() string {
	return r.expect
}

func TestClient_Save(t *testing.T) {
	t.Parallel()

	var (
		srvErrDeadline        = status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error())
		srvErrCanceled        = status.Error(codes.Canceled, context.Canceled.Error())
		srvErrInternal        = status.Error(codes.Internal, errAny.Error())
		srvErrInvalidArgument = status.Error(codes.InvalidArgument, errAny.Error())

		userID    = uuid.Must(uuid.NewV4())
		ip        = net.ParseIP("192.100.10.4")
		userAgent = "userAgent"
		token     = "token"
	)

	testCases := map[string]struct {
		appResponse *session_pb.SaveResponse
		appError    error
		want        *client.Token
		wantErr     error
	}{
		"success":              {&session_pb.SaveResponse{Token: token}, nil, &client.Token{Value: token}, nil},
		"c.conn.Save.":         {nil, srvErrDeadline, nil, context.DeadlineExceeded},
		"err_canceled":         {nil, srvErrCanceled, nil, context.Canceled},
		"err_internal":         {nil, srvErrInternal, nil, client.ErrInternal},
		"err_invalid_argument": {nil, srvErrInvalidArgument, nil, client.ErrInvalidArgument},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, conn, mock, assert := start(t)

			mock.EXPECT().Save(traceIDMatcher{expect: traceID.String()}, &session_pb.SaveRequest{
				UserId:    userID.String(),
				Ip:        ip.String(),
				UserAgent: userAgent,
				Kind:      user_status_pb.StatusKind_STATUS_KIND_DEFAULT,
			}).Return(tc.appResponse, tc.appError)

			token, err := conn.Save(ctx, userID, dom.Origin{IP: ip, UserAgent: userAgent}, dom.UserStatusDefault)
			assert.ErrorIs(err, tc.wantErr)
			assert.Equal(tc.want, token)
		})
	}
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	var (
		srvErrDeadline        = status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error())
		srvErrCanceled        = status.Error(codes.Canceled, context.Canceled.Error())
		srvErrNotFound        = status.Error(codes.NotFound, client.ErrNotFound.Error())
		srvErrInternal        = status.Error(codes.Internal, errAny.Error())
		srvErrInvalidArgument = status.Error(codes.InvalidArgument, errAny.Error())

		token = "token"

		want = &client.Session{
			ID:     uuid.Must(uuid.NewV4()),
			UserID: uuid.Must(uuid.NewV4()),
			Status: dom.UserStatusDefault,
		}
		pbResponse = &session_pb.GetResponse{
			SessionId: want.ID.String(),
			UserId:    want.UserID.String(),
			Kind:      user_status_pb.StatusKind_STATUS_KIND_DEFAULT,
		}
	)

	testCases := map[string]struct {
		appResponse *session_pb.GetResponse
		appError    error
		want        *client.Session
		wantErr     error
	}{
		"success":              {pbResponse, nil, want, nil},
		"err_deadline":         {nil, srvErrDeadline, nil, context.DeadlineExceeded},
		"err_canceled":         {nil, srvErrCanceled, nil, context.Canceled},
		"err_not_found":        {nil, srvErrNotFound, nil, client.ErrNotFound},
		"err_internal":         {nil, srvErrInternal, nil, client.ErrInternal},
		"err_invalid_argument": {nil, srvErrInvalidArgument, nil, client.ErrInvalidArgument},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, conn, mock, assert := start(t)

			mock.EXPECT().Get(traceIDMatcher{expect: traceID.String()}, &session_pb.GetRequest{
				Token: token,
			}).Return(tc.appResponse, tc.appError)

			session, err := conn.Get(ctx, token)
			assert.ErrorIs(err, tc.wantErr)
			assert.Equal(tc.want, session)
		})
	}
}

func TestClient_Delete(t *testing.T) {
	t.Parallel()

	var (
		srvErrDeadline        = status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error())
		srvErrCanceled        = status.Error(codes.Canceled, context.Canceled.Error())
		srvErrNotFound        = status.Error(codes.NotFound, client.ErrNotFound.Error())
		srvErrInternal        = status.Error(codes.Internal, errAny.Error())
		srvErrInvalidArgument = status.Error(codes.InvalidArgument, errAny.Error())

		pbResponse = &session_pb.DeleteResponse{}
		sessionID  = uuid.Must(uuid.NewV4())
	)

	testCases := map[string]struct {
		appResponse *session_pb.DeleteResponse
		appError    error
		want        error
	}{
		"success":              {pbResponse, nil, nil},
		"err_deadline":         {nil, srvErrDeadline, context.DeadlineExceeded},
		"err_canceled":         {nil, srvErrCanceled, context.Canceled},
		"err_not_found":        {nil, srvErrNotFound, client.ErrNotFound},
		"err_internal":         {nil, srvErrInternal, client.ErrInternal},
		"err_invalid_argument": {nil, srvErrInvalidArgument, client.ErrInvalidArgument},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, conn, mock, assert := start(t)

			mock.EXPECT().Delete(traceIDMatcher{expect: traceID.String()}, &session_pb.DeleteRequest{
				SessionId: sessionID.String(),
			}).Return(tc.appResponse, tc.appError)

			err := conn.Delete(ctx, sessionID)
			assert.ErrorIs(err, tc.want)
		})
	}
}
