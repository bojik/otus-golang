package storage

import "time"

type Event struct {
	ID             string
	StartedAt      time.Time
	Title          string
	Description    string
	UserID         int
	FinishedAt     time.Time
	NotifyInterval time.Duration
	Sent           bool
}
