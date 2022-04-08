package app

import (
	"context"
	"fmt"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logg    logger.Logger
	storage storage.DataKeeper
}

func New(logg logger.Logger, storage storage.DataKeeper) *App {
	return &App{
		logg:    logg,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	events, err := a.storage.SelectInterval(event.StartedAt, event.FinishedAt)
	if err != nil {
		return err
	}
	if len(events) > 0 {
		return fmt.Errorf("%v: existed id = %s", ErrDateBusy, events[0].ID)
	}
	resCh := make(chan error)
	defer close(resCh)
	go func() {
		resCh <- a.storage.InsertEvent(event)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resCh:
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	events, err := a.storage.SelectInterval(event.StartedAt, event.FinishedAt)
	if err != nil {
		return err
	}
	if len(events) > 0 {
		return fmt.Errorf("%v: existed id = %s", ErrDateBusy, events[0].ID)
	}
	resCh := make(chan error)
	defer close(resCh)
	go func() {
		resCh <- a.storage.UpdateEvent(event)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resCh:
		if err != nil {
			return err
		}
	}
	return nil
}
