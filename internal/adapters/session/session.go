// Package session needed for get user session by token.
package session

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template/cmd/session/client"
	"github.com/ZergsLaw/back-template/internal/dom"
)

// For easy testing.
type sessionClient interface {
	Save(ctx context.Context, userID uuid.UUID, origin dom.Origin, status dom.UserStatus) (*client.Token, error)
	Get(ctx context.Context, token string) (*client.Session, error)
	Delete(ctx context.Context, sessionID uuid.UUID) error
}

// Client wrapper for session microservice.
type Client struct {
	session      sessionClient
	errConverter func(error) error
}

// New build and returns new session Client.
func New(svc sessionClient, errConverter func(error) error) *Client {
	return &Client{
		session:      svc,
		errConverter: errConverter,
	}
}

// Save for implements app.Sessions.
func (c *Client) Save(ctx context.Context, userID uuid.UUID, origin dom.Origin, status dom.UserStatus) (*dom.Token, error) {
	res, err := c.session.Save(ctx, userID, origin, status)
	if err != nil {
		return nil, c.errConverter(err)
	}

	return &dom.Token{Value: res.Value}, nil
}

// Get for implements app.Sessions.
func (c *Client) Get(ctx context.Context, token string) (*dom.Session, error) {
	res, err := c.session.Get(ctx, token)
	if err != nil {
		return nil, c.errConverter(err)
	}

	return &dom.Session{
		ID:     res.ID,
		UserID: res.UserID,
		Status: res.Status,
	}, nil
}

// Delete for implements app.Sessions.
func (c *Client) Delete(ctx context.Context, sessionID uuid.UUID) error {
	err := c.session.Delete(ctx, sessionID)
	if err != nil {
		return c.errConverter(err)
	}

	return nil
}
