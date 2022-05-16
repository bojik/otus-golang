package api

import (
	"context"
	"testing"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/mocks"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
	pb "github.com/bojik/otus-golang/hw12_13_14_15_calendar/pkg/calendarpb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCalendarApi_InsertEvent(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	api := NewCalendarAPI(a, logg)
	cases := []struct {
		e *app.Event
	}{{
		&app.Event{
			Title:          "title",
			StartedAt:      time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
			FinishedAt:     time.Date(2022, 3, 2, 0, 0, 0, 0, time.UTC),
			Description:    "",
			UserID:         10,
			NotifyInterval: time.Hour,
		},
	}}
	for _, tc := range cases {
		tc := tc
		t.Run("insert event", func(t *testing.T) {
			logg.EXPECT().Info(gomock.Any()).AnyTimes()
			logg.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
			dataKeeper.EXPECT().SelectInterval(tc.e.StartedAt, tc.e.FinishedAt).Return(nil, nil)
			dataKeeper.EXPECT().InsertEvent(tc.e.ConvertToStorageEvent()).Return("123-123", nil)
			dataKeeper.EXPECT().FindById("123-123").Return(tc.e.ConvertToStorageEvent(), nil)
			e, err := api.InsertEvent(context.Background(), api.fillEvent(tc.e))
			require.Nil(t, err)
			require.Equal(t, tc.e.Title, e.Title)
		})
	}
}

func TestCalendarApi_UpdateEvent(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	api := NewCalendarAPI(a, logg)
	cases := []struct {
		e *app.Event
	}{{
		&app.Event{
			Title:          "title",
			StartedAt:      time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
			FinishedAt:     time.Date(2022, 3, 2, 0, 0, 0, 0, time.UTC),
			Description:    "",
			UserID:         10,
			NotifyInterval: time.Hour,
		},
	}}
	for _, tc := range cases {
		tc := tc
		t.Run("update event", func(t *testing.T) {
			logg.EXPECT().Info(gomock.Any()).AnyTimes()
			logg.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
			dataKeeper.EXPECT().SelectInterval(tc.e.StartedAt, tc.e.FinishedAt).Return(nil, nil)
			dataKeeper.EXPECT().UpdateEvent(tc.e.ConvertToStorageEvent()).Return(tc.e.ConvertToStorageEvent(), nil)
			dataKeeper.EXPECT().FindById("").Return(tc.e.ConvertToStorageEvent(), nil)
			e, err := api.UpdateEvent(context.Background(), api.fillEvent(tc.e))
			require.Nil(t, err)
			require.Equal(t, tc.e.Title, e.Title)
		})
	}
}

func TestCalendarApi_DeleteEvent(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	api := NewCalendarAPI(a, logg)
	cases := []struct {
		e *app.Event
	}{{
		&app.Event{
			ID:             "123-123",
			Title:          "title",
			StartedAt:      time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
			FinishedAt:     time.Date(2022, 3, 2, 0, 0, 0, 0, time.UTC),
			Description:    "",
			UserID:         10,
			NotifyInterval: time.Hour,
		},
	}}
	for _, tc := range cases {
		tc := tc
		t.Run("delete event", func(t *testing.T) {
			logg.EXPECT().Info(gomock.Any()).AnyTimes()
			logg.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
			dataKeeper.EXPECT().FindById(tc.e.ID).Return(tc.e.ConvertToStorageEvent(), nil)
			dataKeeper.EXPECT().DeleteEventById(tc.e.ID).Return(nil)
			e, err := api.DeleteEvent(context.Background(), &pb.Id{Id: tc.e.ID})
			require.Nil(t, err)
			require.Equal(t, tc.e.Title, e.Title)
		})
	}
}

func TestCalendarApi_FindEventsByInterval(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	api := NewCalendarAPI(a, logg)
	cases := []struct {
		e []*storage.Event
		f time.Time
		t time.Time
	}{{
		[]*storage.Event{{
			ID:             "123-123",
			Title:          "title",
			StartedAt:      time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
			FinishedAt:     time.Date(2022, 3, 2, 0, 0, 0, 0, time.UTC),
			Description:    "",
			UserID:         10,
			NotifyInterval: time.Hour,
		}},
		time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 3, 2, 0, 0, 0, 0, time.UTC),
	}}
	for _, tc := range cases {
		tc := tc
		t.Run("find events", func(t *testing.T) {
			logg.EXPECT().Info(gomock.Any()).AnyTimes()
			logg.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
			dataKeeper.EXPECT().SelectInterval(tc.f, tc.t).Return(tc.e, nil)
			e, err := api.FindEventsByInterval(
				context.Background(),
				&pb.Interval{StartedAt: timestamppb.New(tc.f), FinishedAt: timestamppb.New(tc.t)},
			)
			require.Nil(t, err)
			require.Equal(t, tc.e[0].Title, e.Events[0].Title)
		})
	}
}
