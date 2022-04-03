package storage

import "time"

type Event struct {
	ID             string
	StartedAt      time.Time
	Title          string
	Description    string
	UserId         int
	FinishedAt     time.Time
	NotifyInterval time.Duration
}
