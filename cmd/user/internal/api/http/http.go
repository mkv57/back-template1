// Package http contains all methods for working http server.
package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
	"github.com/ZergsLaw/back-template1/internal/logger"
)

var (
	ErrUserUnauthorized           = errors.New("user unauthorized")
	ErrMissingAuthorizationHeader = errors.New("missing authorization header")
	ErrBadAuthorizationString     = errors.New("bad authorization string")
	ErrInvalidArgument            = errors.New("invalid argument")
	ErrMaxAvatarSize              = errors.New("max file size 25 mb")
)

type application interface {
	SaveAvatar(ctx context.Context, session dom.Session, file app.Avatar) (uuid.UUID, error)
	GetFile(ctx context.Context, session dom.Session, fileID uuid.UUID) (*app.Avatar, error)
	Auth(ctx context.Context, token string) (*dom.Session, error)
}

type api struct {
	app application
}

// New build and return http.Handler.
func New(ctx context.Context, applications application) http.Handler {
	log := logger.FromContext(ctx)

	api := api{
		app: applications,
	}

	router := mux.NewRouter()

	router.Use(
		LogMiddleware(log),
		Recoverer(log),
		SetSessionToCtx(applications),
	)

	router.HandleFunc("/user/api/v1/file/avatar", api.uploadAvatar).Methods(http.MethodPost)
	router.HandleFunc("/user/api/v1/file/avatar/{id}", api.downloadAvatar).Methods(http.MethodGet)

	return router
}
