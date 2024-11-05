package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"

	"github.com/ZergsLaw/back-template1/internal/adapters/session"
)

const headerAuthorize = "authorization"

const scheme = "Bearer"

// LogMiddleware add logger in middleware.
func LogMiddleware(log *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// some evil middlewares modify this values
			path := r.URL.EscapedPath()
			m := httpsnoop.CaptureMetrics(next, w, r)

			duration := time.Since(start)
			fields := []any{
				slog.String("request_method", r.Method),
				slog.String("url_path", path),
				slog.Duration("event_duration", duration),
				slog.String("client_address", r.RemoteAddr),
			}

			if m.Code >= http.StatusInternalServerError {
				fields = append(fields, slog.Int("status_code", m.Code))
				log.Error("server error", fields...)

				return
			}

			if m.Code >= http.StatusBadRequest {
				fields = append(fields, slog.Int("status_code", m.Code))
				log.Warn("client error", fields...)

				return
			}
			log.Info("", fields...)
		})
	}
}

// Recoverer recover service after a panic.
func Recoverer(log *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					start := time.Now()
					log.Error("panic",
						slog.String("http_method", r.Method),
						slog.String("path", r.URL.EscapedPath()),
						slog.String("trace", string(debug.Stack())),
						slog.Duration("duration", time.Since(start)),
					)

					jsonBody, err := json.Marshal(map[string]string{
						"error": "There was an internal server error",
					})
					if err != nil {
						errorHandler(w, r, http.StatusInternalServerError, err)

						return
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write(jsonBody)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func SetSessionToCtx(app application) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := r.Header.Get(headerAuthorize)
			if value == "" {
				errorHandler(w, r, http.StatusUnauthorized, ErrMissingAuthorizationHeader)

				return
			}

			const headerSize = 2
			splits := strings.SplitN(value, " ", headerSize)
			if len(splits) < headerSize {
				errorHandler(w, r, http.StatusUnauthorized, ErrBadAuthorizationString)

				return
			}
			if !strings.EqualFold(splits[0], scheme) {
				errorHandler(w, r, http.StatusInternalServerError, ErrUserUnauthorized)

				return
			}

			userSession, err := app.Auth(r.Context(), splits[1])
			if err != nil {
				errorHandler(w, r, http.StatusUnauthorized, err)

				return
			}

			next.ServeHTTP(w, r.WithContext(session.NewContext(r.Context(), userSession)))
		})
	}
}
