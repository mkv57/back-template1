package app

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template1/internal/dom"
)

// NewSession save new user session.
func (a *App) NewSession(ctx context.Context, userID uuid.UUID, status dom.UserStatus, origin Origin) (*Token, error) {
	sessionID := a.id.New()
	token, err := a.auth.Token(sessionID)
	if err != nil {
		return nil, fmt.Errorf("a.auth.Token: %w", err)
	}

	session := Session{
		ID:     sessionID,
		Origin: origin,
		Token:  *token,
		UserID: userID,
		Status: status,
	}

	err = a.session.Save(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("a.session.Save: %w", err)
	}

	return token, nil
}

// Session get user session by access token.
func (a *App) Session(ctx context.Context, token string) (*Session, error) {
	subject, err := a.auth.Subject(token)
	if err != nil {
		return nil, fmt.Errorf("a.auth.Subject: %w", err)
	}

	session, err := a.session.ByID(ctx, subject)
	if err != nil {
		return nil, fmt.Errorf("a.session.ByID: %w", err)
	}

	return session, nil
}

// RemoveSession remove user's session by id.
func (a *App) RemoveSession(ctx context.Context, id uuid.UUID) error {
	session, err := a.session.ByID(ctx, id)
	if err != nil {
		return fmt.Errorf("a.session.ByID: %w", err)
	}

	return a.session.Delete(ctx, session.ID)
}
