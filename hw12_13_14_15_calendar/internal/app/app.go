package app

import (
	"context"
	"fmt"
	"time"

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

func (a *App) CreateEvent(ctx context.Context, event Event) (string, error) {
	if err := event.validate(); err != nil {
		return "", err
	}
	events, err := a.storage.SelectInterval(event.StartedAt, event.FinishedAt)
	if err != nil {
		return "", err
	}
	if len(events) > 0 {
		return "", fmt.Errorf("%w: existed id = %s", ErrDateBusy, events[0].ID)
	}
	type result struct {
		id  string
		err error
	}
	resCh := make(chan result)
	defer close(resCh)
	go func() {
		id, err := a.storage.InsertEvent(event.ConvertToStorageEvent())
		if err != nil {
			resCh <- result{id: "", err: err}
			return
		}
		resCh <- result{id: id, err: nil}
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-resCh:
		if res.err != nil {
			return "", err
		}
		a.logg.Debug("inserted event", logger.NewStringParam("id", res.id))
		return res.id, nil
	}
}

func (a *App) UpdateEvent(ctx context.Context, event Event) error {
	if err := event.validate(); err != nil {
		return err
	}
	events, err := a.storage.SelectInterval(event.StartedAt, event.FinishedAt)
	if err != nil {
		return err
	}
	if len(events) > 0 && event.ID != events[0].ID {
		return fmt.Errorf("%w: existed id = %s", ErrDateBusy, events[0].ID)
	}
	resCh := make(chan error)
	defer close(resCh)
	go func() {
		_, err := a.storage.UpdateEvent(event.ConvertToStorageEvent())
		resCh <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resCh:
		if err != nil {
			return err
		}
	}
	a.logg.Debug("updated event", logger.NewStringParam("id", event.ID))
	return nil
}

func (a *App) FindByID(ctx context.Context, id string) (*Event, error) {
	type res struct {
		evt *storage.Event
		err error
	}
	resCh := make(chan res)
	defer close(resCh)
	go func() {
		evt, err := a.storage.FindByID(id)
		resCh <- res{
			evt: evt,
			err: err,
		}
	}()
	var r res
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r = <-resCh:
		if r.err != nil {
			return nil, r.err
		}
	}
	return newFromStorageEvent(r.evt), nil
}

func (a *App) DeleteByID(ctx context.Context, id string) (*Event, error) {
	type res struct {
		evt *storage.Event
		err error
	}
	resCh := make(chan res)
	defer close(resCh)
	go func() {
		evt, err := a.storage.FindByID(id)
		if err != nil {
			resCh <- res{evt: nil, err: err}
			return
		}
		err = a.storage.DeleteEventByID(id)
		if err != nil {
			resCh <- res{evt: nil, err: err}
			return
		}
		resCh <- res{evt: evt, err: nil}
	}()
	var r res
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r = <-resCh:
		if r.err != nil {
			return nil, r.err
		}
	}
	return newFromStorageEvent(r.evt), nil
}

func (a *App) FindEventsByInterval(ctx context.Context, startDate, endDate time.Time) ([]*Event, error) {
	return a.wrapSelect(ctx, func() ([]*storage.Event, error) {
		return a.storage.SelectInterval(startDate, endDate)
	})
}

func (a *App) FindToSend(ctx context.Context) ([]*Event, error) {
	return a.wrapSelect(ctx, func() ([]*storage.Event, error) {
		return a.storage.SelectToNotify()
	})
}

func (a *App) MarkAsSent(ctx context.Context, id string) error {
	return a.wrapExec(ctx, func() error {
		return a.storage.UpdateSentFlag(id)
	})
}

func (a *App) DeleteOldEvents(ctx context.Context) error {
	return a.wrapExec(ctx, func() error {
		return a.storage.DeleteOldEvents()
	})
}

func (a *App) wrapExec(ctx context.Context, fn func() error) error {
	resCh := make(chan error)
	defer close(resCh)
	go func() {
		resCh <- fn()
	}()
	select {
	case <-ctx.Done():
		return nil
	case err := <-resCh:
		return err
	}
}

func (a *App) wrapSelect(ctx context.Context, fn func() ([]*storage.Event, error)) ([]*Event, error) {
	type result struct {
		events []*storage.Event
		err    error
	}
	resCh := make(chan result)
	defer close(resCh)
	go func() {
		events, err := fn()
		if err != nil {
			resCh <- result{events: nil, err: err}
			return
		}
		resCh <- result{events: events, err: nil}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-resCh:
		if r.err != nil {
			return nil, r.err
		}
		return newFromStorageEventCollection(r.events), nil
	}
}

func (a *App) FindDayEvents(ctx context.Context, date time.Time) ([]*Event, error) {
	return a.findHandler(ctx, date, a.storage.SelectDay)
}

func (a *App) FindWeekEvents(ctx context.Context, date time.Time) ([]*Event, error) {
	return a.findHandler(ctx, date, a.storage.SelectWeek)
}

func (a *App) FindMonthEvents(ctx context.Context, date time.Time) ([]*Event, error) {
	return a.findHandler(ctx, date, a.storage.SelectMonth)
}

type handler func(date time.Time) ([]*storage.Event, error)

func (a *App) findHandler(ctx context.Context, date time.Time, handlerFunc handler) ([]*Event, error) {
	type result struct {
		events []*storage.Event
		err    error
	}
	resCh := make(chan result)
	defer close(resCh)
	go func() {
		events, err := handlerFunc(date)
		if err != nil {
			resCh <- result{events: nil, err: err}
			return
		}
		resCh <- result{events: events, err: nil}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-resCh:
		if r.err != nil {
			return nil, r.err
		}
		return newFromStorageEventCollection(r.events), nil
	}
}
