package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
)

const (
	duplEmail            = "users_email_key"
	duplUsername         = "users_name_key"
	duplOwnerIDAndFileID = "files_owner_id_file_id_key"
	fkUserID             = "fk_user_id_ref_users"
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
	case strings.HasSuffix(pqErr.Message, fmt.Sprintf("unique constraint \"%s\"", duplEmail)):
		return app.ErrEmailExist
	case strings.HasSuffix(pqErr.Message, fmt.Sprintf("unique constraint \"%s\"", duplUsername)):
		return app.ErrUsernameExist
	case strings.HasSuffix(pqErr.Message, fmt.Sprintf("unique constraint \"%s\"", duplOwnerIDAndFileID)):
		return app.ErrUserIDAndFileIDExist
	case strings.HasSuffix(pqErr.Message, fmt.Sprintf("violates foreign key constraint \"%s\"", fkUserID)):
		return app.ErrNotFound
	default:
		return pqErr
	}
}
