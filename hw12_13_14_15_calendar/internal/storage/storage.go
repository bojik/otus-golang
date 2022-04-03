package storage

import "time"

type DataKeeper interface {
	InsertEvent(evt *Event) error
	UpdateEvent(evt *Event) error
	FindById(id string) (*Event, error)
	DeleteEvent(evt *Event) error
	SelectAll() ([]*Event, error)
	SelectInterval(startTime, endTime time.Time) ([]*Event, error)
	SelectDay(date time.Time) ([]*Event, error)
	SelectWeek(date time.Time) ([]*Event, error)
	SelectMonth(date time.Time) ([]*Event, error)
}
