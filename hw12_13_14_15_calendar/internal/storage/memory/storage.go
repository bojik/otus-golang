package memorystorage

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	events map[string]*storage.Event
	mu     sync.RWMutex
}

func (s *Storage) DeleteOldEvents() error {
	return nil
}

func (s *Storage) UpdateSentFlag(id string) error {
	return nil
}

func (s *Storage) SelectToNotify() ([]*storage.Event, error) {
	return nil, nil
}

func New() *Storage {
	return &Storage{
		events: map[string]*storage.Event{},
	}
}

func (s *Storage) InsertEvent(evt *storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if evt.ID == "" {
		evt.ID = uuid.New().String()
	}
	if _, ok := s.events[evt.ID]; ok {
		return "", fmt.Errorf("%w: id = %s", ErrEventAlreadyInserted, evt.ID)
	}
	s.events[evt.ID] = evt
	return evt.ID, nil
}

func (s *Storage) UpdateEvent(evt *storage.Event) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[evt.ID]; !ok {
		return nil, fmt.Errorf("%w: id = %s", ErrEventNotFound, evt.ID)
	}
	s.events[evt.ID] = evt
	return evt, nil
}

func (s *Storage) FindByID(id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret, ok := s.events[id]
	if !ok {
		return nil, fmt.Errorf("%w: tut id = %s", ErrEventNotFound, id)
	}
	return ret, nil
}

func (s *Storage) DeleteEventByID(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return fmt.Errorf("%w: id = %s", ErrEventNotFound, id)
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) DeleteEvent(evt *storage.Event) error {
	if err := s.DeleteEventByID(evt.ID); err != nil {
		return err
	}
	return nil
}

func (s *Storage) SelectAll() ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]*storage.Event, len(s.events))
	i := 0
	for _, evt := range s.events {
		events[i] = evt
		i++
	}
	return events, nil
}

func (s *Storage) SelectInterval(startTime, endTime time.Time) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := []*storage.Event{}
	begin := startTime.Add(-time.Second)
	end := endTime.Add(time.Second)
	for _, event := range s.events {
		if event.StartedAt.After(begin) && event.StartedAt.Before(end) {
			events = append(events, event)
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartedAt.Before(events[j].StartedAt)
	})
	return events, nil
}

func (s *Storage) SelectDay(date time.Time) ([]*storage.Event, error) {
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())
	return s.SelectInterval(startTime, endTime)
}

func (s *Storage) SelectWeek(date time.Time) ([]*storage.Event, error) {
	weekday := date.Weekday()
	startTime := time.Date(
		date.Year(),
		date.Month(),
		date.Day()-int(weekday),
		0,
		0,
		0,
		0,
		date.Location(),
	)
	endTime := startTime.Add(7 * time.Hour * 24).Add(-time.Second)
	return s.SelectInterval(startTime, endTime)
}

func (s *Storage) SelectMonth(date time.Time) ([]*storage.Event, error) {
	startTime := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endTime := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location()).
		Add(-time.Second)
	return s.SelectInterval(startTime, endTime)
}

var _ storage.DataKeeper = (*Storage)(nil)
