//go:build integration
// +build integration

package sqlstorage

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

const (
	DbDsn           = "postgres://postgres:postgres@localhost:5432/postgres_test?sslmode=disable"
	MigrationSource = "file://../../../migrations"
	MaxOpenConnects = 10
	MaxIdleConnects = 10
)

var (
	dsn, migrationSource       string
	openConnects, idleConnects int
)

// go test -tags=integration ./internal/storage/sql/* -args -db-dsn='postgres://postgres:postgres@localhost:5432/postgres_test?sslmode=disable' -migration-source='file://../../../migrations' -open-connects=10 -idle-connects=10
func init() {
	flag.StringVar(&dsn, "db-dsn", DbDsn, "")
	flag.StringVar(&migrationSource, "migration-source", MigrationSource, "")
	flag.IntVar(&openConnects, "open-connects", MaxOpenConnects, "")
	flag.IntVar(&idleConnects, "idle-connects", MaxIdleConnects, "")
}

type StorageTestSuite struct {
	suite.Suite
	storage *Storage
}

func (s *StorageTestSuite) SetupSuite() {
	s.storage = New(dsn, idleConnects, openConnects)
	err := s.storage.Connect(context.Background())
	require.Nil(s.T(), err)
	err = s.storage.Migrate(migrationSource)
	if err != nil && !xerrors.Is(err, migrate.ErrNoChange) {
		require.Nil(s.T(), err)
	}
	_, err = s.storage.getDb().Exec("truncate event")
	require.Nil(s.T(), err)
}

func (s *StorageTestSuite) TearDownSuite() {
	err := s.storage.Close(context.Background())
	require.Nil(s.T(), err)
}

func (s *StorageTestSuite) TestCRUD() {
	e := &storage.Event{
		Title:          "Title",
		StartedAt:      time.Now().Add(24 * time.Hour),
		Description:    "Description",
		UserID:         1,
		FinishedAt:     time.Now().Add(2 * 24 * time.Hour),
		NotifyInterval: time.Hour,
	}
	// Insert
	{
		id, err := s.storage.InsertEvent(e)
		require.Nil(s.T(), err)
		require.NotEmpty(s.T(), id)
		require.NotEqual(s.T(), "", e.ID)
	}
	// Select
	{
		events, err := s.storage.SelectAll()
		require.Nil(s.T(), err)
		require.Len(s.T(), events, 1)
	}
	// Update
	{
		e2, err := s.storage.FindById(e.ID)
		require.Nil(s.T(), err)
		require.Equal(s.T(), e.Title, e2.Title)
		require.True(s.T(), e2.StartedAt.Equal(e.StartedAt))
		require.Equal(s.T(), e.NotifyInterval, e2.NotifyInterval)
		e2.Title = "Title2"
		e4, err := s.storage.UpdateEvent(e2)
		require.Nil(s.T(), err)
		e3, err := s.storage.FindById(e.ID)
		require.Equal(s.T(), e2.Title, e3.Title)
		require.Equal(s.T(), e2.Title, e4.Title)
	}
	// Delete
	{
		err := s.storage.DeleteEvent(e)
		events, err := s.storage.SelectAll()
		require.Nil(s.T(), err)
		require.Len(s.T(), events, 0)
	}
}

func (s *StorageTestSuite) TestSelectInterval() {
	dates := []time.Time{
		time.Date(2022, 1, 1, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 2, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 3, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 4, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 5, 22, 22, 22, 0, time.Local),
	}
	for _, date := range dates {
		e := &storage.Event{StartedAt: date}
		_, err := s.storage.InsertEvent(e)
		require.Nil(s.T(), err)
	}
	ets, err := s.storage.SelectInterval(dates[1], dates[3])
	require.Nil(s.T(), err)
	require.Len(s.T(), ets, 3)
	expected := dates[1:4]
	for i, e := range ets {
		require.True(s.T(), expected[i].Equal(e.StartedAt))
	}
	_, err = s.storage.getDb().Exec("truncate event")
	require.Nil(s.T(), err)
}

func (s *StorageTestSuite) TestSelectDay() {
	dates := []time.Time{
		time.Date(2022, 1, 1, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
		time.Date(2022, 1, 2, 22, 22, 22, 0, time.Local),
		time.Date(2022, 1, 2, 23, 59, 59, 0, time.Local),
		time.Date(2022, 1, 5, 22, 22, 22, 0, time.Local),
	}
	for _, date := range dates {
		e := &storage.Event{StartedAt: date}
		_, err := s.storage.InsertEvent(e)
		require.Nil(s.T(), err)
	}
	ets, err := s.storage.SelectDay(dates[2])
	require.Nil(s.T(), err)
	require.Len(s.T(), ets, 3)
	expected := dates[1:4]
	for i, e := range ets {
		require.True(s.T(), expected[i].Equal(e.StartedAt))
	}
	_, err = s.storage.getDb().Exec("truncate event")
	require.Nil(s.T(), err)
}

func (s *StorageTestSuite) TestSelectWeek() {
	dates := []time.Time{
		time.Date(2022, 3, 26, 22, 22, 22, 0, time.Local),
		time.Date(2022, 3, 27, 22, 22, 22, 0, time.Local),
		time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
		time.Date(2022, 4, 2, 22, 22, 22, 0, time.Local),
		time.Date(2022, 4, 3, 23, 59, 59, 0, time.Local),
		time.Date(2022, 4, 4, 22, 22, 22, 0, time.Local),
	}
	for _, date := range dates {
		e := &storage.Event{StartedAt: date}
		_, err := s.storage.InsertEvent(e)
		require.Nil(s.T(), err)
	}
	ets, err := s.storage.SelectWeek(dates[3])
	require.Nil(s.T(), err)
	require.Len(s.T(), ets, 3)
	expected := dates[1:4]
	for i, e := range ets {
		require.True(s.T(), expected[i].Equal(e.StartedAt))
	}
	_, err = s.storage.getDb().Exec("truncate event")
	require.Nil(s.T(), err)
}

func (s *StorageTestSuite) TestSelectMonth() {
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
		e := &storage.Event{StartedAt: date}
		_, err := s.storage.InsertEvent(e)
		require.Nil(s.T(), err)
	}
	ets, err := s.storage.SelectMonth(dates[3])
	require.Nil(s.T(), err)
	require.Len(s.T(), ets, 3)
	expected := dates[3:6]
	for i, e := range ets {
		require.True(s.T(), expected[i].Equal(e.StartedAt))
	}
	_, err = s.storage.getDb().Exec("truncate event")
	require.Nil(s.T(), err)
}

func TestStorage(t *testing.T) {
	flag.Parse()
	suite.Run(t, new(StorageTestSuite))
}
