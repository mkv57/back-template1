package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ZergsLaw/back-template1/internal/dom"
	"github.com/ZergsLaw/back-template1/internal/logger"
)

func (a *App) Process(ctx context.Context) error {
	log := logger.FromContext(ctx)

	for {
		var err error
		select {
		case <-ctx.Done():
			return nil
		case msg := <-a.queue.UpSessionStatus():
			err = a.handleUpdateStatus(ctx, msg)
		}
		if err != nil {
			log.Error("couldn't handle event", slog.String(logger.Error.String(), err.Error()))

			continue
		}
	}
}

func (a *App) handleUpdateStatus(ctx context.Context, event dom.Event[UpdateStatus]) error {
	err := a.session.UpdateStatus(ctx, event.ID(), event.Body().UserID, event.Body().Status)
	switch {
	case errors.Is(err, ErrDuplicate):
		// We must acknowledge this message.
	case err != nil:
		event.Nack(ctx)

		return fmt.Errorf("a.session.UpdateStatus: %w", err)
	}

	event.Ack(ctx)

	return nil
}
