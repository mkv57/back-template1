package api

import (
	"context"
	"fmt"
	"net"

	"github.com/gofrs/uuid"

	pb "github.com/ZergsLaw/back-template/api/session/v1"
	"github.com/ZergsLaw/back-template/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template/internal/dom"
)

// Save implements pb.SessionAPIServer.
func (a *api) Save(ctx context.Context, request *pb.SaveRequest) (*pb.SaveResponse, error) {
	userID, err := uuid.FromString(request.UserId)
	if err != nil {
		return nil, fmt.Errorf("uuid.FromString: %w", err)
	}

	token, err := a.app.NewSession(ctx, userID, dom.UserStatusFromAPI(request.Kind), app.Origin{
		IP:        net.ParseIP(request.Ip),
		UserAgent: request.UserAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("a.app.NewSession: %w", err)
	}

	return &pb.SaveResponse{Token: token.Value}, nil
}

// Get implements pb.SessionAPIServer.
func (a *api) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	session, err := a.app.Session(ctx, request.Token)
	if err != nil {
		return nil, fmt.Errorf("a.app.Session: %w", err)
	}

	return &pb.GetResponse{
		SessionId: session.ID.String(),
		UserId:    session.UserID.String(),
		Kind:      dom.UserStatusToAPI(session.Status),
	}, nil
}

// Delete implements pb.SessionAPIServer.
func (a *api) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	uid, err := uuid.FromString(request.SessionId)
	if err != nil {
		return nil, fmt.Errorf("uuid.FromString: %w", err)
	}

	err = a.app.RemoveSession(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("a.app.RemoveSession: %w", err)
	}

	return &pb.DeleteResponse{}, nil
}
