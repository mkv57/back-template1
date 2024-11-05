package app

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template1/internal/dom"
)

type (
	// Repo interface for session data repository.
	Repo interface {
		// Save saves the new user session in a database.
		// Errors: unknown.
		Save(context.Context, Session) error
		// ByID returns user session by session id.
		// Errors: ErrNotFound, unknown.
		ByID(context.Context, uuid.UUID) (*Session, error)
		// Delete removes user session.
		// Errors: ErrNotFound, unknown.
		Delete(context.Context, uuid.UUID) error
		// UpdateStatus change user session status.
		// Errors: unknown.
		UpdateStatus(ctx context.Context, reqID, userID uuid.UUID, status dom.UserStatus) error
	}

	// Auth interface for generate access and refresh token by subject.
	Auth interface {
		// Token generate tokens by subject with expire time.
		// Errors: unknown.
		Token(uuid.UUID) (*Token, error)
		// Subject unwrap Subject info from token.
		// Errors: ErrInvalidToken, ErrExpiredToken, unknown.
		Subject(token string) (uuid.UUID, error)
	}

	// ID generator for session.
	ID interface {
		// New generate new ID for session.
		New() uuid.UUID
	}

	// Queue module for getting events from queue.
	Queue interface {
		// UpSessionStatus returns channel for getting new events.
		UpSessionStatus() <-chan dom.Event[UpdateStatus]
	}
)
