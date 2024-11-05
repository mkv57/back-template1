package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/ZergsLaw/back-template1/cmd/session/internal/app"
)

const (
	duplPKey    = "deduplication_pkey"
	duplPrimary = "primary"
)

func convertErr(err error) error {
	var pqErr *pq.Error

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return app.ErrNotFound
	case errors.As(err, &pqErr):
		return constraint(pqErr)
	default:
		return err
	}
}

func constraint(pqErr *pq.Error) error {
	switch {
	case strings.HasSuffix(pqErr.Message, fmt.Sprintf("unique constraint %q", duplPKey)):
		return app.ErrDuplicate
	case strings.HasSuffix(pqErr.Message, fmt.Sprintf("unique constraint %q", duplPrimary)):
		return app.ErrDuplicate
	default:
		return pqErr
	}
}
