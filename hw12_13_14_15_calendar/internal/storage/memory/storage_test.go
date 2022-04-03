package memorystorage

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		stor := New()
		c := 10
		ids := []string{}
		{
			wg := sync.WaitGroup{}
			wg.Add(c)
			mu := sync.Mutex{}
			for i := 0; i < c; i++ {
				i := i
				go func() {
					defer wg.Done()
					e := generateEvent()
					e.Title = fmt.Sprintf("%s #%d", e.Title, i)
					err := stor.InsertEvent(e)
					require.Nil(t, err)
					mu.Lock()
					defer mu.Unlock()
					ids = append(ids, e.ID)
				}()
			}
			wg.Wait()
		}
		events, err := stor.SelectAll()
		require.Nil(t, err)
		require.Equal(t, c, len(events))
		{
			wg := sync.WaitGroup{}
			wg.Add(len(ids))
			for i := 0; i < len(ids); i++ {
				i := i
				go func() {
					defer wg.Done()
					e := generateEvent()
					e.ID = ids[i]
					err := stor.UpdateEvent(e)
					require.Nil(t, err)
				}()
			}
			wg.Wait()
			for i := 0; i < len(ids); i++ {
				e, err := stor.FindById(ids[i])
				require.Nil(t, err)
				require.Equal(t, "title", e.Title)
			}
		}
		{
			wg := sync.WaitGroup{}
			wg.Add(len(ids))
			for i := 0; i < len(ids); i++ {
				i := i
				go func() {
					defer wg.Done()
					e, err := stor.FindById(ids[i])
					require.Nil(t, err)
					err = stor.DeleteEvent(e)
					require.Nil(t, err)
				}()
			}
			wg.Wait()
			evts, err := stor.SelectAll()
			require.Nil(t, err)
			require.Equal(t, 0, len(evts))
		}
	})
}

func TestStorage_SelectInterval(t *testing.T) {
	store := New()
	dates := []time.Time{
		time.Date(2022, 1, 1, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 2, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 3, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 4, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 5, 22, 22, 22, 0, time.Local),
	}
	for _, date := range dates {
		e := generateEvent()
		e.StartedAt = date
		err := store.InsertEvent(e)
		require.Nil(t, err)
	}
	ets, err := store.SelectInterval(dates[1], dates[3])
	require.Nil(t, err)
	require.Len(t, ets, 3)
	expected := dates[1:4]
	for i, e := range ets {
		require.True(t, expected[i].Equal(e.StartedAt))
	}
}

func TestStorage_SelectDay(t *testing.T) {
	store := New()
	dates := []time.Time{
		time.Date(2022, 1, 1, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
		time.Date(2022, 1, 2, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 2, 23, 59, 59, 0, time.Local),
		time.Date(2022, 1, 5, 22, 22, 22, 0, time.Local),
	}
	for _, date := range dates {
		e := generateEvent()
		e.StartedAt = date
		err := store.InsertEvent(e)
		require.Nil(t, err)
	}
	ets, err := store.SelectDay(dates[2])
	require.Nil(t, err)
	require.Len(t, ets, 3)
	expected := dates[1:4]
	for i, e := range ets {
		require.True(t, expected[i].Equal(e.StartedAt))
	}
}

func TestStorage_SelectWeek(t *testing.T) {
	store := New()
	dates := []time.Time{
		time.Date(2022, 3, 26, 22, 22, 22, 0, time.Local),
		time.Date(2022, 3, 27, 22, 22, 22, 0, time.Local),
		time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
		time.Date(2022, 4, 2, 22, 22, 22, 0, time.Local),
		time.Date(2022, 4, 3, 23, 59, 59, 0, time.Local),
		time.Date(2022, 4, 4, 22, 22, 22, 0, time.Local),
	}
	for _, date := range dates {
		e := generateEvent()
		e.StartedAt = date
		err := store.InsertEvent(e)
		require.Nil(t, err)
	}
	ets, err := store.SelectWeek(dates[3])
	require.Nil(t, err)
	require.Len(t, ets, 3)
	expected := dates[1:4]
	for i, e := range ets {
		require.True(t, expected[i].Equal(e.StartedAt))
	}
}

func TestStorage_SelectMonth(t *testing.T) {
	store := New()
	dates := []time.Time{
		time.Date(2022, 3, 26, 22, 22, 22, 0, time.Local),
		time.Date(2022, 3, 27, 22, 22, 22, 0, time.Local),
		time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
		time.Date(2022, 4, 2, 22, 22, 22, 0, time.Local),
		time.Date(2022, 4, 3, 23, 59, 59, 0, time.Local),
		time.Date(2022, 4, 4, 22, 22, 22, 0, time.Local),
		time.Date(2022, 5, 4, 22, 22, 22, 0, time.Local),
	}
	for _, date := range dates {
		e := generateEvent()
		e.StartedAt = date
		err := store.InsertEvent(e)
		require.Nil(t, err)
	}
	ets, err := store.SelectMonth(dates[3])
	require.Nil(t, err)
	require.Len(t, ets, 3)
	expected := dates[3:6]
	for i, e := range ets {
		require.True(t, expected[i].Equal(e.StartedAt))
	}
}

func TestError(t *testing.T) {
	stor := New()
	e := generateEvent()
	err := stor.InsertEvent(e)
	require.Nil(t, err)
	err = stor.InsertEvent(e)
	require.ErrorIs(t, err, ErrEventAlreadyInserted)
	e = generateEvent()
	err = stor.UpdateEvent(e)
	require.ErrorIs(t, err, ErrEventNotFound)
	_, err = stor.FindById("not_found")
	require.ErrorIs(t, err, ErrEventNotFound)
	err = stor.DeleteEvent(e)
	require.ErrorIs(t, err, ErrEventNotFound)
}

func generateEvent() *storage.Event {
	return &storage.Event{
		Title:          "title",
		Description:    "description",
		StartedAt:      time.Now(),
		FinishedAt:     time.Now(),
		UserId:         1,
		NotifyInterval: time.Hour,
	}
}
