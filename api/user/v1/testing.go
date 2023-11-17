package pb

import (
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

var (
	_ gomock.Matcher = &VerificationEmailRequest{}
	_ gomock.Matcher = &VerificationUsernameRequest{}
	_ gomock.Matcher = &CreateUserRequest{}
	_ gomock.Matcher = &LoginRequest{}
	_ gomock.Matcher = &GetUserRequest{}
	_ gomock.Matcher = &SearchUsersRequest{}
	_ gomock.Matcher = &LogoutRequest{}
	_ gomock.Matcher = &UpdatePasswordRequest{}
	_ gomock.Matcher = &UpdateUserRequest{}
	_ gomock.Matcher = &RemoveAvatarRequest{}
	_ gomock.Matcher = &ListUserAvatarRequest{}
	_ gomock.Matcher = &GetUsersByIDsRequest{}
)

func (x *VerificationEmailRequest) Matches(y interface{}) bool    { return match(x, y) }
func (x *VerificationUsernameRequest) Matches(y interface{}) bool { return match(x, y) }
func (x *CreateUserRequest) Matches(y interface{}) bool           { return match(x, y) }
func (x *LoginRequest) Matches(y interface{}) bool                { return match(x, y) }
func (x *GetUserRequest) Matches(y interface{}) bool              { return match(x, y) }
func (x *SearchUsersRequest) Matches(y interface{}) bool          { return match(x, y) }
func (x *LogoutRequest) Matches(y interface{}) bool               { return match(x, y) }
func (x *UpdatePasswordRequest) Matches(y interface{}) bool       { return match(x, y) }
func (x *UpdateUserRequest) Matches(y interface{}) bool           { return match(x, y) }
func (x *RemoveAvatarRequest) Matches(y interface{}) bool         { return match(x, y) }
func (x *ListUserAvatarRequest) Matches(y interface{}) bool       { return match(x, y) }
func (x *GetUsersByIDsRequest) Matches(y interface{}) bool        { return match(x, y) }

func match(x proto.Message, y interface{}) bool {
	p2, ok := y.(proto.Message)
	if !ok {
		return false
	}

	return proto.Equal(x, p2)
}
