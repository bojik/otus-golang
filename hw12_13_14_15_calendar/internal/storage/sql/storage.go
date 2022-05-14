package sqlstorage

import (
	"context"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
	migrate "github.com/golang-migrate/migrate/v4"
	// postgresql driver.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// file driver.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	SQLInsertEvent = `insert into event(id, title, started_at, description, user_id, finished_at, notify_interval) 
		values(:id, :title, :started_at, :description, :user_id, :finished_at, :notify_interval)`
	SQLUpdateEvent = `update event set title = :title, started_at = :started_at, description = :description, 
                 user_id = :user_id, finished_at = :finished_at, notify_interval = :notify_interval where id = :id`
	SQLSelectByID = `select id, title, started_at, description, user_id, finished_at, notify_interval from event
				where id = :id`
	SQLDeleteByID     = `delete from event where id = :id`
	SQLSelectAll      = `select id, title, started_at, description, user_id, finished_at, notify_interval from event`
	SQLSelectInterval = `select id, title, started_at, description, user_id, finished_at, notify_interval from event 
				where started_at between :start_date and :end_date order by started_at`
	SQLSelectNotify = `select id, title, description, started_at, finished_at from event
		where started_at - (notify_interval::text||' milliseconds')::interval <= current_timestamp and not sent`
	SQLUpdateSent      = `update event set sent = true where id = :id`
	SQLDeleteOldEvents = `delete from event where finished_at < current_timestamp - '1 year'::interval`
)

type Storage struct {
	dsn             string
	db              *sqlx.DB
	maxOpenConnects int
	maxIdleConnects int
}

func New(dsn string, maxIdleConnects, maxOpenConnects int) *Storage {
	return &Storage{
		dsn:             dsn,
		maxOpenConnects: maxOpenConnects,
		maxIdleConnects: maxIdleConnects,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", s.dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(s.maxOpenConnects)
	db.SetMaxIdleConns(s.maxIdleConnects)
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	var err error
	done := make(chan struct{})
	go func() {
		err = s.db.Close()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) Migrate(migrationSource string) error {
	m, err := migrate.New(migrationSource, s.dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) MigrateDown(migrationSource string, steps int) error {
	m, err := migrate.New(migrationSource, s.dsn)
	if err != nil {
		return err
	}
	if steps > 0 {
		steps *= -1
	}
	if err := m.Steps(steps); err != nil {
		return err
	}
	return nil
}

func (s *Storage) FixAndForce(migrationSource string, version int) error {
	m, err := migrate.New(migrationSource, s.dsn)
	if err != nil {
		return err
	}
	if err := m.Force(version); err != nil {
		return err
	}
	return nil
}

func (s *Storage) InsertEvent(evt *storage.Event) (string, error) {
	id := uuid.New().String()
	_, err := s.db.NamedExec(
		SQLInsertEvent,
		map[string]interface{}{
			"id":              id,
			"title":           evt.Title,
			"user_id":         evt.UserID,
			"started_at":      evt.StartedAt,
			"finished_at":     evt.FinishedAt,
			"description":     evt.Description,
			"notify_interval": evt.NotifyInterval,
		},
	)
	if err != nil {
		return "", err
	}
	evt.ID = id
	return id, nil
}

func (s *Storage) UpdateEvent(evt *storage.Event) (*storage.Event, error) {
	_, err := s.db.NamedExec(
		SQLUpdateEvent,
		map[string]interface{}{
			"id":              evt.ID,
			"title":           evt.Title,
			"user_id":         evt.UserID,
			"started_at":      evt.StartedAt,
			"finished_at":     evt.FinishedAt,
			"description":     evt.Description,
			"notify_interval": evt.NotifyInterval,
		},
	)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func (s *Storage) FindByID(id string) (*storage.Event, error) {
	rows, err := s.db.NamedQuery(SQLSelectByID, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	evt := &storage.Event{}
	for rows.Next() {
		if err := rows.Scan(
			&evt.ID,
			&evt.Title,
			&evt.StartedAt,
			&evt.Description,
			&evt.UserID,
			&evt.FinishedAt,
			&evt.NotifyInterval,
		); err != nil {
			return nil, err
		}
	}
	return evt, nil
}

func (s *Storage) DeleteEventByID(id string) error {
	_, err := s.db.NamedExec(SQLDeleteByID, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(evt *storage.Event) error {
	if err := s.DeleteEventByID(evt.ID); err != nil {
		return err
	}
	return nil
}

func (s *Storage) SelectAll() ([]*storage.Event, error) {
	rows, err := s.db.Query(SQLSelectAll)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer func() {
		_ = rows.Close()
	}()
	var events []*storage.Event
	for rows.Next() {
		evt := &storage.Event{}
		if err := rows.Scan(
			&evt.ID,
			&evt.Title,
			&evt.StartedAt,
			&evt.Description,
			&evt.UserID,
			&evt.FinishedAt,
			&evt.NotifyInterval,
		); err != nil {
			return nil, err
		}
		events = append(events, evt)
	}
	return events, nil
}

func (s *Storage) SelectInterval(startTime, endTime time.Time) ([]*storage.Event, error) {
	rows, err := s.db.NamedQuery(
		SQLSelectInterval,
		map[string]interface{}{"start_date": startTime, "end_date": endTime},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var events []*storage.Event
	for rows.Next() {
		evt := &storage.Event{}
		if err := rows.Scan(
			&evt.ID,
			&evt.Title,
			&evt.StartedAt,
			&evt.Description,
			&evt.UserID,
			&evt.FinishedAt,
			&evt.NotifyInterval,
		); err != nil {
			return nil, err
		}
		events = append(events, evt)
	}
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

func (s *Storage) SelectToNotify() ([]*storage.Event, error) {
	rows, err := s.db.Queryx(SQLSelectNotify)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var events []*storage.Event
	for rows.Next() {
		evt := &storage.Event{}
		if err := rows.Scan(
			&evt.ID,
			&evt.Title,
			&evt.Description,
			&evt.StartedAt,
			&evt.FinishedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, evt)
	}
	return events, nil
}

func (s *Storage) UpdateSentFlag(id string) error {
	_, err := s.db.NamedExec(SQLUpdateSent, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteOldEvents() error {
	_, err := s.db.Exec(SQLDeleteOldEvents)
	if err != nil {
		return err
	}
	return nil
}

var _ storage.DataKeeper = (*Storage)(nil)
