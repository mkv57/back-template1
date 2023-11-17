//go:build integration

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	user_pb "github.com/ZergsLaw/back-template/api/user/v1"
	"github.com/ZergsLaw/back-template/internal/queue"
	"github.com/ZergsLaw/back-template/internal/testhelper"
)

func TestE2E(t *testing.T) {
	t.Parallel()
	ctx := testhelper.Context(t)
	assert, grpcClient, subscriber, cfg := initService(t, ctx)

	consumerName := t.Name()

	err := subscriber.Migrate(func(manager nats.JetStreamManager) error {
		err := subscriber.Migrate(user_pb.Migrate)
		if err != nil {
			return fmt.Errorf("subscriber.JetStream.AddConsumer: %w", err)
		}

		return nil
	})
	assert.NoError(err)

	verificationEmailReq := &user_pb.VerificationEmailRequest{Email: email}
	_, err = grpcClient.VerificationEmail(ctx, verificationEmailReq)
	assert.NoError(err)

	verificationUsernameReq := &user_pb.VerificationUsernameRequest{Username: username1}
	_, err = grpcClient.VerificationUsername(ctx, verificationUsernameReq)
	assert.NoError(err)

	createUserReq1 := &user_pb.CreateUserRequest{
		Username: username1,
		FullName: fullName,
		Email:    email,
		Password: pass1,
	}
	createUserResp, err := grpcClient.CreateUser(ctx, createUserReq1)
	assert.NoError(err)

	md := metadata.MD{}
	loginReq := &user_pb.LoginRequest{
		Email:    email,
		Password: pass1,
	}
	_, err = grpcClient.Login(ctx, loginReq, grpc.Header(&md))
	assert.NoError(err)

	ctx = auth(ctx, strings.Join(md.Get("authorization"), ""))

	getSelfInfo := &user_pb.GetUserRequest{}
	selfInfo, err := grpcClient.GetUser(ctx, getSelfInfo)
	assert.NoError(err)
	assert.Equal(createUserResp.Id, selfInfo.User.Id)
	assert.Equal(createUserReq1.Email, selfInfo.User.Email)
	assert.Equal(createUserReq1.Username, selfInfo.User.Username)

	subscribeCtx, subscribeCtxCancel := context.WithTimeout(ctx, time.Second*2)
	t.Cleanup(subscribeCtxCancel)

	err = subscriber.Subscribe(subscribeCtx, user_pb.SubscribeToAllEvents, consumerName, func(ctx context.Context, message queue.Message) error {
		err := message.Ack(ctx)
		if err != nil {
			return fmt.Errorf("message.Ack: %w", err)
		}

		assert.Equal(user_pb.TopicAdd, message.Subject())
		eventAdd := user_pb.Event{}
		err = message.Unmarshal(&eventAdd)
		assert.NoError(err)
		assert.Equal(createUserResp.Id, eventAdd.GetAdd().User.Id)
		assert.Equal(selfInfo.User.Email, eventAdd.GetAdd().User.Email)

		return nil
	})
	assert.NoError(err)

	createUserReq2 := &user_pb.CreateUserRequest{
		Username: username2,
		Email:    email2,
		FullName: fullName,
		Password: pass2,
	}
	createUserResp2, err := grpcClient.CreateUser(ctx, createUserReq2)
	assert.NoError(err)

	getUserInfoReq := &user_pb.GetUserRequest{
		Id: createUserResp2.Id,
	}
	user2Info, err := grpcClient.GetUser(ctx, getUserInfoReq)
	assert.NoError(err)
	assert.Equal(createUserResp2.Id, user2Info.User.Id)
	assert.Equal(createUserReq2.Email, user2Info.User.Email)
	assert.Equal(createUserReq2.Username, user2Info.User.Username)

	searchUserReq := &user_pb.SearchUsersRequest{
		Name:   user2Info.User.Username,
		Limit:  5,
		Offset: 0,
	}
	searchUserInfo, err := grpcClient.SearchUsers(ctx, searchUserReq)
	assert.NoError(err)
	assert.Len(searchUserInfo.Users, 1)
	assert.Equal(searchUserInfo.Total, int32(1))
	assert.Equal(createUserResp2.Id, searchUserInfo.Users[0].Id)
	assert.Equal(createUserReq2.Email, searchUserInfo.Users[0].Email)
	assert.Equal(createUserReq2.Username, searchUserInfo.Users[0].Username)

	updatePassReq := &user_pb.UpdatePasswordRequest{
		Old: pass1,
		New: pass2,
	}
	_, err = grpcClient.UpdatePassword(ctx, updatePassReq)
	assert.NoError(err)

	_, err = grpcClient.Logout(ctx, &user_pb.LogoutRequest{})
	assert.NoError(err)

	md = metadata.MD{}
	loginReq = &user_pb.LoginRequest{
		Email:    email,
		Password: pass2,
	}
	_, err = grpcClient.Login(ctx, loginReq, grpc.Header(&md))
	assert.NoError(err)

	// Upload Avatar
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	f, err := os.Open(avatarPath)
	assert.NoError(err)
	t.Cleanup(func() { assert.NoError(f.Close()) })

	fileWriter, err := bodyWriter.CreateFormFile("file", path.Base(avatarPath))
	assert.NoError(err)

	_, err = io.Copy(fileWriter, f)
	assert.NoError(err)

	err = bodyWriter.Close()
	assert.NoError(err)

	h := http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", strings.Join(md.Get("authorization"), ""))},
	}

	addr := fmt.Sprintf("http://%s:%d/user/api/v1/file/avatar", cfg.Server.Host, cfg.Server.Port.Files)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, addr, bodyBuf)
	assert.NoError(err)
	request.Header = h
	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	testClient := &http.Client{}
	resp, err := testClient.Do(request)
	assert.NoError(err)
	t.Cleanup(func() {
		assert.NoError(resp.Body.Close())
	})
	assert.Equal(http.StatusCreated, resp.StatusCode)

	getFile := struct {
		FileID uuid.UUID `json:"file_id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&getFile)
	assert.NoError(err)

	// Update UserProfile
	ctx = auth(ctx, strings.Join(md.Get("authorization"), ""))

	updateUserReq := &user_pb.UpdateUserRequest{
		Username: "update_test",
		FullName: "full_name",
		AvatarId: getFile.FileID.String(),
	}
	_, err = grpcClient.UpdateUser(ctx, updateUserReq)
	assert.NoError(err)

	updateUser, err := grpcClient.GetUser(ctx, getSelfInfo)
	assert.NoError(err)
	assert.Equal(updateUserReq.Username, updateUser.User.Username)
	assert.Equal(updateUserReq.AvatarId, updateUser.User.AvatarId)

	// get file by id
	getFileAddr := fmt.Sprintf("http://%s:%d/user/api/v1/file/avatar", cfg.Server.Host, cfg.Server.Port.Files)
	addr = fmt.Sprintf("%s/%s", getFileAddr, getFile.FileID.String())
	request, err = http.NewRequestWithContext(ctx, http.MethodGet, addr, http.NoBody)
	assert.NoError(err)
	request.Header = h
	resp, err = testClient.Do(request)
	assert.NoError(err)
	t.Cleanup(func() {
		assert.NoError(resp.Body.Close())
	})
	assert.Equal(http.StatusOK, resp.StatusCode)

	f2, err := os.Open(avatarPath)
	assert.NoError(err)
	t.Cleanup(func() { assert.NoError(f2.Close()) })

	bufFromSrv, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	expBytes, err := io.ReadAll(f2)
	assert.NoError(err)
	assert.Equal(expBytes, bufFromSrv)

	removeAvatarReq := &user_pb.RemoveAvatarRequest{
		FileId: getFile.FileID.String(),
	}
	_, err = grpcClient.RemoveAvatar(ctx, removeAvatarReq)
	assert.NoError(err)

	resp, err = testClient.Do(request)
	assert.NoError(err)
	t.Cleanup(func() {
		assert.NoError(resp.Body.Close())
	})
}
