package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ZergsLaw/back-template1/internal/dom"
	"github.com/ZergsLaw/back-template1/internal/logger"
)

func (a *App) Process(ctx context.Context) error {
	// TODO: Add support for master-slaves nodes.
	// TODO: Add tests.
	var (
		log   = logger.FromContext(ctx)
		wg    = &sync.WaitGroup{}
		tasks = make(chan *dom.Event[Task])
	)
	defer wg.Wait()

	wg.Add(1)
	go a.collectingTasks(ctx, wg, tasks)

	for {
		var err error
		select {
		case <-ctx.Done():
			return nil
		case task := <-tasks:
			switch task.Body().Kind {
			case TaskKindEventAdd:
				err = a.handleTaskKindEventAdd(ctx, task.Body())
			case TaskKindEventDel:
				err = a.handleTaskKindEventDel(ctx, task.Body())
			case TaskKindEventUpdate:
				err = a.handleTaskKindEventUpdate(ctx, task.Body())
			default:
				log.Error("unknown task",
					slog.String(logger.TaskID.String(), task.Body().ID.String()),
					slog.String(logger.TaskKind.String(), task.Body().Kind.String()),
				)

				continue
			}
			if err != nil {
				task.Nack(ctx)
				log.Error("couldn't handle event", slog.String(logger.Error.String(), err.Error()))

				continue
			}

			task.Ack(ctx)
		}
	}
}

func (a *App) handleTaskKindEventAdd(ctx context.Context, task Task) error {
	err := a.queue.AddUser(ctx, task.ID, task.User)
	if err != nil {
		return fmt.Errorf("a.queue.AddUser: %w", err)
	}

	err = a.repo.FinishTask(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("a.repo.FinishTask: %w", err)
	}

	return nil
}

func (a *App) handleTaskKindEventDel(ctx context.Context, task Task) error {
	err := a.queue.DeleteUser(ctx, task.ID, task.User)
	if err != nil {
		return fmt.Errorf("a.queue.DeleteUser: %w", err)
	}

	err = a.repo.FinishTask(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("a.repo.FinishTask: %w", err)
	}

	return nil
}

func (a *App) handleTaskKindEventUpdate(ctx context.Context, task Task) error {
	err := a.queue.UpdateUser(ctx, task.ID, task.User)
	if err != nil {
		return fmt.Errorf("a.queue.UpdateUser: %w", err)
	}

	err = a.repo.FinishTask(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("a.subscribe.FinishTask: %w", err)
	}

	return nil
}

func (a *App) collectingTasks(ctx context.Context, wg *sync.WaitGroup, out chan *dom.Event[Task]) {
	defer wg.Done()

	const (
		taskTickerTimeout = time.Second / 10
		taskLimit         = 100
	)
	var (
		log    = logger.FromContext(ctx)
		ticker = time.NewTicker(taskTickerTimeout)
	)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tasks, err := a.repo.ListActualTask(ctx, taskLimit)
			if err != nil {
				log.Error("couldn't get tasks", slog.String(logger.Error.String(), err.Error()))

				continue
			}

			handle := func(event *dom.Event[Task], ackCh chan dom.AcknowledgeKind) {
				select {
				case <-ctx.Done():
					return
				case out <- event:
				}

				for {
					select {
					case <-ctx.Done():
						return
					case ack := <-ackCh:
						if ack == dom.AcknowledgeKindAck {
							return
						}

						select {
						case <-ctx.Done():
							return
						case out <- event:
						}
					}
				}
			}

			ackCh := make(chan dom.AcknowledgeKind)
			for i := range tasks {
				event := dom.NewEvent(tasks[i].ID, ackCh, tasks[i])

				handle(event, ackCh)
			}
		}
	}
}
