package session

import (
	"context"

	"github.com/ZergsLaw/back-template/internal/dom"
)

type ctxMarker struct{}

// NewContext returns context with slog.Logger.
func NewContext(ctx context.Context, session *dom.Session) context.Context {
	return context.WithValue(ctx, ctxMarker{}, session)
}

// FromContext returns slog.Logger from context.
func FromContext(ctx context.Context) *dom.Session {
	l, ok := ctx.Value(ctxMarker{}).(*dom.Session)
	if !ok {
		return nil
	}

	return l
}
