// Package auth contains methods for working with authorization tokens,
// their generation and parsing.
package auth

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/o1egl/paseto/v2"

	"github.com/ZergsLaw/back-template/cmd/session/internal/app"
)

var _ app.Auth = &Auth{}

// Auth is implements app.Auth.
// Responsible for working with authorization tokens, be it cookies or jwt.
type Auth struct {
	key []byte
}

// New creates and returns new instance auth.
func New(secretKey string) *Auth {
	return &Auth{
		key: []byte(secretKey),
	}
}

type jsonToken struct {
	SessionID uuid.UUID `json:"session_id"`
}

// Token need for implements app.Auth.
func (a *Auth) Token(subject uuid.UUID) (*app.Token, error) {
	t := jsonToken{
		SessionID: subject,
	}

	value, err := paseto.Encrypt(a.key, t, "")
	if err != nil {
		return nil, fmt.Errorf("paseto.Encrypt: %w", err)
	}

	res := &app.Token{
		Value: value,
	}

	return res, nil
}

// Subject need for implements app.Auth.
func (a *Auth) Subject(token string) (uuid.UUID, error) {
	t := jsonToken{}

	err := paseto.Decrypt(token, a.key, &t, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %s", app.ErrInvalidToken, err)
	}

	return t.SessionID, nil
}
