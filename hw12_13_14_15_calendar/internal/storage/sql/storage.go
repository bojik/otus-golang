package sqlstorage

import (
	"context"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	SqlInsertEvent = `insert into event(id, title, started_at, description, user_id, finished_at, notify_interval) 
		values(:id, :title, :started_at, :description, :user_id, :finished_at, :notify_interval)`
	SqlUpdateEvent = `update event set title = :title, started_at = :started_at, description = :description, 
                 user_id = :user_id, finished_at = :finished_at, notify_interval = :notify_interval where id = :id`
	SqlSelectById = `select id, title, started_at, description, user_id, finished_at, notify_interval from event
				where id = :id`
	SqlDeleteById     = `delete from event where id = :id`
	SqlSelectAll      = `select id, title, started_at, description, user_id, finished_at, notify_interval from event`
	SqlSelectInterval = `select id, title, started_at, description, user_id, finished_at, notify_interval from event 
				where started_at between :start_date and :end_date order by started_at`
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

func (s *Storage) InsertEvent(evt *storage.Event) error {
	id := uuid.New().String()
	_, err := s.db.NamedExec(
		SqlInsertEvent,
		map[string]interface{}{
			"id":              id,
			"title":           evt.Title,
			"user_id":         evt.UserId,
			"started_at":      evt.StartedAt,
			"finished_at":     evt.FinishedAt,
			"description":     evt.Description,
			"notify_interval": evt.NotifyInterval,
		},
	)
	if err != nil {
		return err
	}
	evt.ID = id
	return nil
}

func (s *Storage) UpdateEvent(evt *storage.Event) error {
	_, err := s.db.NamedExec(
		SqlUpdateEvent,
		map[string]interface{}{
			"id":              evt.ID,
			"title":           evt.Title,
			"user_id":         evt.UserId,
			"started_at":      evt.StartedAt,
			"finished_at":     evt.FinishedAt,
			"description":     evt.Description,
			"notify_interval": evt.NotifyInterval,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) FindById(id string) (*storage.Event, error) {
	rows, err := s.db.NamedQuery(SqlSelectById, map[string]interface{}{"id": id})
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
			&evt.UserId,
			&evt.FinishedAt,
			&evt.NotifyInterval,
		); err != nil {
			return nil, err
		}
	}
	return evt, nil
}

func (s *Storage) DeleteEvent(evt *storage.Event) error {
	_, err := s.db.NamedExec(SqlDeleteById, map[string]interface{}{"id": evt.ID})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) SelectAll() ([]*storage.Event, error) {
	rows, err := s.db.Query(SqlSelectAll)
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
			&evt.UserId,
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
		SqlSelectInterval,
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
			&evt.UserId,
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

func (s *Storage) getDb() *sqlx.DB {
	return s.db
}

var _ storage.DataKeeper = (*Storage)(nil)
