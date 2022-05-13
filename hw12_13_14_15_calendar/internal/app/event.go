package app

import (
	"fmt"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Event struct {
	ID             string
	Title          string
	StartedAt      time.Time
	FinishedAt     time.Time
	Description    string
	UserID         int
	NotifyInterval time.Duration
}

func newFromStorageEvent(evt *storage.Event) *Event {
	return &Event{
		ID:             evt.ID,
		Title:          evt.Title,
		StartedAt:      evt.StartedAt,
		FinishedAt:     evt.FinishedAt,
		Description:    evt.Description,
		UserID:         evt.UserID,
		NotifyInterval: evt.NotifyInterval,
	}
}

func newFromStorageEventCollection(events []*storage.Event) []*Event {
	if len(events) == 0 {
		return nil
	}
	appEvents := make([]*Event, len(events))
	for i, e := range events {
		appEvents[i] = newFromStorageEvent(e)
	}
	return appEvents
}

func (e Event) ConvertToStorageEvent() *storage.Event {
	return &storage.Event{
		ID:             e.ID,
		Title:          e.Title,
		StartedAt:      e.StartedAt,
		FinishedAt:     e.FinishedAt,
		Description:    e.Description,
		UserID:         e.UserID,
		NotifyInterval: e.NotifyInterval,
	}
}

func (e Event) validate() error {
	if e.Title == "" {
		return fmt.Errorf("title %w", ErrRequiredField)
	}
	if e.UserID == 0 {
		return fmt.Errorf("userId %w", ErrRequiredField)
	}
	if e.FinishedAt.Before(e.StartedAt) {
		return fmt.Errorf("%w: finishedAt < startedAt", ErrInvalidDate)
	}
	return nil
}
