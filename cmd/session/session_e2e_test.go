//go:build integration

package main

import (
	"net"
	"testing"

	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ZergsLaw/back-template/api/session/v1"
	user_pb "github.com/ZergsLaw/back-template/api/user/v1"
	user_status_pb "github.com/ZergsLaw/back-template/api/user_status/v1"
	"github.com/ZergsLaw/back-template/internal/dom"
	"github.com/ZergsLaw/back-template/internal/testhelper"
)

func TestE2E(t *testing.T) {
	t.Parallel()
	ctx := testhelper.Context(t)
	assert, grpcClient, consumer := initService(t, ctx)
	t.Cleanup(func() {
		err := consumer.Drain()
		assert.NoError(err)
	})

	saveReq := &pb.SaveRequest{
		UserId:    uuid.Must(uuid.NewV4()).String(),
		Ip:        net.ParseIP("127.0.0.1").String(),
		UserAgent: "UserAgent",
		Kind:      user_status_pb.StatusKind_STATUS_KIND_DEFAULT,
	}

	saveResp, err := grpcClient.Save(ctx, saveReq)
	assert.NoError(err)
	assert.NotNil(saveResp)

	getReq := &pb.GetRequest{
		Token: saveResp.Token,
	}

	getResp, err := grpcClient.Get(ctx, getReq)
	assert.NoError(err)
	assert.Equal(saveReq.UserId, getResp.UserId)
	assert.Equal(saveReq.Kind, user_status_pb.StatusKind_STATUS_KIND_DEFAULT)

	userID, err := uuid.FromString(saveReq.UserId)
	assert.NoError(err)

	eventAddUserUpStatusTaskID := uuid.Must(uuid.NewV4())
	eventUpdateStatus := &user_pb.Event{
		Body: &user_pb.Event_Update{
			Update: &user_pb.Update{
				User: &user_pb.User{
					Id:        userID.String(),
					Username:  "zergslaw",
					Email:     "edgar@google.com",
					AvatarId:  "",
					Kind:      dom.UserStatusToAPI(dom.UserStatusJedi),
					FullName:  "Edgar Sipki",
					CreatedAt: timestamppb.Now(),
					UpdatedAt: nil,
				},
			},
		},
	}

	err = consumer.Publish(ctx, user_pb.TopicUpdate, eventAddUserUpStatusTaskID, eventUpdateStatus)
	assert.NoError(err)

	for {
		getResp, err = grpcClient.Get(ctx, getReq)
		assert.NoError(err)
		assert.Equal(saveReq.UserId, getResp.UserId)
		if getResp.Kind == user_status_pb.StatusKind_STATUS_KIND_DEFAULT {
			continue
		}

		break
	}

	assert.Equal(getResp.Kind, user_status_pb.StatusKind_STATUS_KIND_JEDI)

	delReq := &pb.DeleteRequest{
		SessionId: getResp.SessionId,
	}

	_, err = grpcClient.Delete(ctx, delReq)
	assert.NoError(err)

	getResp, err = grpcClient.Get(ctx, getReq)
	assert.Nil(getResp)
}
