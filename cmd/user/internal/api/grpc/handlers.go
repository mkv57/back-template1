package grpc

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	user_pb "github.com/ZergsLaw/back-template/api/user/v1"
	"github.com/ZergsLaw/back-template/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template/internal/adapters/session"
)

const (
	userAgentForward = `grpcgateway-user-agent`
	userAgent        = `user-agent`
	auth             = `authorization`
	scheme           = `Bearer`
)

var ErrUnauthenticated = status.Error(codes.Unauthenticated, "unauthenticated")

// VerificationEmail implements pb.UserExternalAPIServer.
func (a *api) VerificationEmail(ctx context.Context, request *user_pb.VerificationEmailRequest) (*user_pb.VerificationEmailResponse, error) {
	err := a.app.VerificationEmail(ctx, request.Email)
	if err != nil {
		return nil, fmt.Errorf("a.app.VerificationEmail: %w", err)
	}

	return &user_pb.VerificationEmailResponse{}, nil
}

// VerificationUsername implements pb.UserExternalAPIServer.
func (a *api) VerificationUsername(ctx context.Context, request *user_pb.VerificationUsernameRequest) (*user_pb.VerificationUsernameResponse, error) {
	err := a.app.VerificationUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("a.app.VerificationUsername: %w", err)
	}

	return &user_pb.VerificationUsernameResponse{}, nil
}

// CreateUser implements pb.UserExternalAPIServer.
func (a *api) CreateUser(ctx context.Context, request *user_pb.CreateUserRequest) (*user_pb.CreateUserResponse, error) {
	id, err := a.app.CreateUser(ctx, request.Email, request.Username, request.FullName, request.Password)
	if err != nil {
		return nil, fmt.Errorf("a.app.CreateUser: %w", err)
	}

	return &user_pb.CreateUserResponse{Id: id.String()}, nil
}

// Login implements pb.UserExternalAPIServer.
func (a *api) Login(ctx context.Context, request *user_pb.LoginRequest) (*user_pb.LoginResponse, error) {
	origin, err := originFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("originFromCtx: %w", err)
	}

	userID, token, err := a.app.Login(ctx, request.Email, request.Password, *origin)
	if err != nil {
		return nil, fmt.Errorf("a.app.Login: %w", err)
	}

	err = grpc.SendHeader(ctx, metadata.MD{auth: {token.Value}})
	if err != nil {
		return nil, fmt.Errorf("grpc.SendHeader: %w", err)
	}

	return &user_pb.LoginResponse{UserId: userID.String()}, nil
}

// Logout implements pb.UserExternalAPIServer.
func (a *api) Logout(ctx context.Context, _ *user_pb.LogoutRequest) (*user_pb.LogoutResponse, error) {
	userSession := session.FromContext(ctx)
	if userSession == nil {
		return nil, ErrUnauthenticated
	}

	err := a.app.Logout(ctx, *userSession)
	if err != nil {
		return nil, fmt.Errorf("a.app.Logout: %w", err)
	}

	return &user_pb.LogoutResponse{}, nil
}

// GetUser implements pb.UserExternalAPIServer.
func (a *api) GetUser(ctx context.Context, request *user_pb.GetUserRequest) (*user_pb.GetUserResponse, error) {
	userSession := session.FromContext(ctx)
	if userSession == nil {
		return nil, ErrUnauthenticated
	}

	userID := uuid.FromStringOrNil(request.Id)
	user, err := a.app.UserByID(ctx, *userSession, userID)
	if err != nil {
		return nil, fmt.Errorf("a.app.UserByID: %w", err)
	}

	return &user_pb.GetUserResponse{User: toUser(*user)}, nil
}

// SearchUsers implements pb.UserExternalAPIServer.
func (a *api) SearchUsers(ctx context.Context, request *user_pb.SearchUsersRequest) (*user_pb.SearchUsersResponse, error) {
	userSession := session.FromContext(ctx)
	if userSession == nil {
		return nil, ErrUnauthenticated
	}

	users, total, err := a.app.ListUserByFilters(ctx, *userSession,
		app.SearchParams{
			OwnerID:  userSession.UserID,
			Username: request.Name,
			FullName: request.Name,
			Limit:    uint64(request.Limit),
			Offset:   uint64(request.Offset),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("a.app.ListUserByUsername: %w", err)
	}

	pbUsers := make([]*user_pb.User, len(users))
	for i := range users {
		pbUsers[i] = toUser(users[i])
	}

	return &user_pb.SearchUsersResponse{Users: pbUsers, Total: int32(total)}, nil
}

// UpdatePassword implements pb.UserExternalAPIServer.
func (a *api) UpdatePassword(ctx context.Context, request *user_pb.UpdatePasswordRequest) (*user_pb.UpdatePasswordResponse, error) {
	userSession := session.FromContext(ctx)
	if userSession == nil {
		return nil, ErrUnauthenticated
	}

	err := a.app.UpdatePassword(ctx, *userSession, request.Old, request.New)
	if err != nil {
		return nil, fmt.Errorf("a.app.UpdatePassword: %w", err)
	}

	return &user_pb.UpdatePasswordResponse{}, nil
}

// UpdateUser implements pb.UserExternalAPIServer.
func (a *api) UpdateUser(ctx context.Context, request *user_pb.UpdateUserRequest) (*user_pb.UpdateUserResponse, error) {
	userSession := session.FromContext(ctx)
	if userSession == nil {
		return nil, ErrUnauthenticated
	}

	err := a.app.UpdateUser(ctx, *userSession, request.Username, uuid.FromStringOrNil(request.AvatarId))
	if err != nil {
		return nil, fmt.Errorf("a.app.UpdateUser: %w", err)
	}

	return &user_pb.UpdateUserResponse{}, nil
}

// RemoveAvatar implements pb.UserExternalAPIServer.
func (a *api) RemoveAvatar(ctx context.Context, request *user_pb.RemoveAvatarRequest) (*user_pb.RemoveAvatarResponse, error) {
	userSession := session.FromContext(ctx)
	if userSession == nil {
		return nil, ErrUnauthenticated
	}

	err := a.app.RemoveAvatar(ctx, *userSession, uuid.FromStringOrNil(request.FileId))
	if err != nil {
		return nil, fmt.Errorf("a.app.RemoveAvatar: %w", err)
	}

	return &user_pb.RemoveAvatarResponse{}, nil
}

// ListUserAvatar implements pb.UserExternalAPIServer.
func (a *api) ListUserAvatar(ctx context.Context, _ *user_pb.ListUserAvatarRequest) (*user_pb.ListUserAvatarResponse, error) {
	userSession := session.FromContext(ctx)
	if userSession == nil {
		return nil, ErrUnauthenticated
	}

	filesCache, err := a.app.ListUserAvatars(ctx, *userSession)
	if err != nil {
		return nil, fmt.Errorf("a.app.ListUserAvatars: %w", err)
	}

	pbAvatars := make([]*user_pb.UserAvatar, len(filesCache))
	for i := range filesCache {
		pbAvatars[i] = toUserFile(filesCache[i])
	}

	return &user_pb.ListUserAvatarResponse{Avatars: pbAvatars}, nil
}

func (a *api) GetUsersByIDs(ctx context.Context, request *user_pb.GetUsersByIDsRequest) (*user_pb.GetUsersByIDsResponse, error) {
	userSession := session.FromContext(ctx)
	if userSession == nil {
		return nil, ErrUnauthenticated
	}

	ids := make([]uuid.UUID, len(request.Ids))

	for i := range request.Ids {
		id, err := uuid.FromString(request.Ids[i])
		if err != nil {
			return nil, fmt.Errorf("uuid.FromString: %w", err)
		}

		ids[i] = id
	}

	users, err := a.app.GetUsersByIDs(ctx, *userSession, ids)
	if err != nil {
		return nil, fmt.Errorf("a.app.GetUsersByIDs: %w", err)
	}

	pbUsers := make([]*user_pb.User, len(users))
	for i := range users {
		pbUsers[i] = toUser(users[i])
	}

	return &user_pb.GetUsersByIDsResponse{
		Result: pbUsers,
	}, nil
}
