package app_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/ZergsLaw/back-template1/cmd/session/internal/app"
	"github.com/ZergsLaw/back-template1/internal/dom"
)

func TestApp_Process(t *testing.T) {
	t.Parallel()

	ctx, module, mocks, assert := start(t)

	newStatus := app.UpdateStatus{
		Status: dom.UserStatusDefault,
		UserID: uuid.Must(uuid.NewV4()),
	}

	chUpStatus := make(chan dom.Event[app.UpdateStatus])
	ack := make(chan dom.AcknowledgeKind)

	mocks.queue.EXPECT().UpSessionStatus().AnyTimes().Return(chUpStatus)

	errC := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	go func() { errC <- module.Process(ctx) }()

	eventID := uuid.Must(uuid.NewV4())
	mocks.repo.EXPECT().UpdateStatus(ctx, eventID, newStatus.UserID, newStatus.Status).Return(nil)
	chUpStatus <- *dom.NewEvent(eventID, ack, app.UpdateStatus{UserID: newStatus.UserID, Status: newStatus.Status})
	assert.Equal(dom.AcknowledgeKindAck, <-ack)

	eventID = uuid.Must(uuid.NewV4())
	mocks.repo.EXPECT().UpdateStatus(ctx, eventID, newStatus.UserID, newStatus.Status).Return(errAny)
	chUpStatus <- *dom.NewEvent(eventID, ack, app.UpdateStatus{UserID: newStatus.UserID, Status: newStatus.Status})
	assert.Equal(dom.AcknowledgeKindNack, <-ack)

	cancel()
	assert.NoError(<-errC)
}
