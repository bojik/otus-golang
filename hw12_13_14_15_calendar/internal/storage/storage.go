package storage

import "time"

//nolint:lll
//go:generate mockgen -destination=../mocks/mock_data_keeper.go -package=mocks github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage DataKeeper
type DataKeeper interface {
	InsertEvent(evt *Event) (string, error)
	UpdateEvent(evt *Event) (*Event, error)
	FindByID(id string) (*Event, error)
	DeleteEventByID(id string) error
	DeleteEvent(evt *Event) error
	SelectAll() ([]*Event, error)
	SelectInterval(startTime, endTime time.Time) ([]*Event, error)
	SelectDay(date time.Time) ([]*Event, error)
	SelectWeek(date time.Time) ([]*Event, error)
	SelectMonth(date time.Time) ([]*Event, error)
	SelectToNotify() ([]*Event, error)
	UpdateSentFlag(id string) error
	DeleteOldEvents() error
}
