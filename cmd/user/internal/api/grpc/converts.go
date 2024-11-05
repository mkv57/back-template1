package grpc

import (
	user_pb "github.com/ZergsLaw/back-template1/api/user/v1"
	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

func toUser(u app.User) *user_pb.User {
	return &user_pb.User{
		Id:       u.ID.String(),
		Username: u.Name,
		Email:    u.Email,
		FullName: u.FullName,
		AvatarId: u.AvatarID.String(),
		Kind:     dom.UserStatusToAPI(u.Status),
	}
}

func toUserFile(f app.AvatarInfo) *user_pb.UserAvatar {
	return &user_pb.UserAvatar{
		UserId: f.OwnerID.String(),
		FileId: f.FileID.String(),
	}
}
