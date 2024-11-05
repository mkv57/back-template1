package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/ZergsLaw/back-template1/cmd/user/internal/app"
	"github.com/ZergsLaw/back-template1/internal/adapters/session"
	"github.com/ZergsLaw/back-template1/internal/logger"
)

const maxAvatarSize = 25 << 20

func (a *api) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	userSession := session.FromContext(r.Context())
	if userSession == nil {
		errorHandler(w, r, http.StatusUnauthorized, ErrUserUnauthorized)

		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err)

		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.FromContext(r.Context()).Error("couldn't not close file", slog.String(logger.Error.String(), err.Error()))
		}
	}()

	if handler.Size >= maxAvatarSize {
		errorHandler(w, r, http.StatusBadRequest, ErrMaxAvatarSize)

		return
	}

	var buf [512]byte
	_, err = file.Read(buf[:])
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, fmt.Errorf("file.Read: %w", err))

		return
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, fmt.Errorf("file.Seek: %w", err))

		return
	}
	contentType := http.DetectContentType(buf[:])

	f := app.Avatar{
		Name:           handler.Filename,
		ContentType:    contentType,
		Size:           handler.Size,
		ReadSeekCloser: file,
	}

	id, err := a.app.SaveAvatar(r.Context(), *userSession, f)
	switch {
	case err == nil:
		resp := &UploadFilesResponse{
			FileID: id,
		}
		responseHandler(w, r, http.StatusCreated, resp)

		return
	case errors.Is(err, app.ErrInvalidImageFormat):
		errorHandler(w, r, http.StatusBadRequest, err)

		return
	default:
		errorHandler(w, r, http.StatusInternalServerError, err)

		return
	}
}

func (a *api) downloadAvatar(w http.ResponseWriter, r *http.Request) {
	userSession := session.FromContext(r.Context())
	if userSession == nil {
		errorHandler(w, r, http.StatusUnauthorized, ErrUserUnauthorized)

		return
	}

	v := mux.Vars(r)
	fileID := uuid.FromStringOrNil(v["id"])
	if fileID == uuid.Nil {
		errorHandler(w, r, http.StatusBadRequest, ErrInvalidArgument)

		return
	}

	file, err := a.app.GetFile(r.Context(), *userSession, fileID)
	switch {
	case err == nil:
		http.ServeContent(w, r, file.Name, file.ModTime, file.ReadSeekCloser)
		defer func() {
			err := file.Close()
			if err != nil {
				logger.FromContext(r.Context()).Error("couldn't not close file", slog.String(logger.Error.String(), err.Error()))
			}
		}()

		return
	case errors.Is(err, app.ErrNotFound):
		errorHandler(w, r, http.StatusNotFound, err)

		return
	default:
		errorHandler(w, r, http.StatusInternalServerError, err)

		return
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.WriteHeader(code)
	erR := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if erR != nil {
		logger.FromContext(r.Context()).Error("couldn't send error msg", slog.String(logger.Error.String(), err.Error()))

		return
	}
}

func responseHandler(w http.ResponseWriter, r *http.Request, code int, resp interface{}) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err)

		return
	}
}
