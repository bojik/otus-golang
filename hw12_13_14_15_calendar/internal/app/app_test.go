package app

import (
	"testing"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/docker/distribution/context"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	logg := logger.New()
	logg.ResetWriters()
	storage := memorystorage.New()
	app := New(logg, storage)
	evt := Event{
		Title:      "Title",
		StartedAt:  time.Now(),
		FinishedAt: time.Now().Add(time.Hour),
		UserId:     10,
	}
	id, err := app.CreateEvent(context.Background(), evt)
	require.Nil(t, err)
	require.NotEmpty(t, id)
	e, err := app.FindById(context.Background(), id)
	require.Nil(t, err)
	require.Equal(t, id, e.Id)
}
