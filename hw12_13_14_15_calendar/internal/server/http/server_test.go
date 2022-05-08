package internalhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/mocks"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestIndexServeHTTP(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	logg.EXPECT().Info(gomock.Any()).AnyTimes()
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	server := NewServer(logg, a, "")
	req, err := http.NewRequest("GET", "/", nil)
	require.Nil(t, err)
	rr := httptest.NewRecorder()
	handler := server.newHandler()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestEventCreate(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	server := NewServer(logg, a, "")
	cases := []struct {
		e storage.Event
		b string
	}{
		{
			e: storage.Event{
				Title:          "title",
				UserId:         10,
				StartedAt:      time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
				FinishedAt:     time.Date(2022, 12, 1, 23, 0, 0, 0, time.UTC),
				NotifyInterval: time.Hour,
			},
			b: `{"Success":true,"Error":"","Data":{"Event":{"Id":"","Title":"title","StartedAt":"2022-12-01T00:00:00Z","FinishedAt":"2022-12-01T23:00:00Z","Description":"","UserId":10,"NotifyInterval":3600000000000}}}`,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.e.Title, func(t *testing.T) {
			logg.EXPECT().Info(gomock.Any()).AnyTimes()
			logg.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
			dataKeeper.EXPECT().SelectInterval(tc.e.StartedAt, tc.e.FinishedAt).Return(nil, nil)
			dataKeeper.EXPECT().InsertEvent(&tc.e).Return("123-123", nil)
			dataKeeper.EXPECT().FindById("123-123").Return(&tc.e, nil)
			params := url.Values{
				"title":           []string{tc.e.Title},
				"user_id":         []string{fmt.Sprintf("%d", tc.e.UserId)},
				"started_at":      []string{tc.e.StartedAt.Format(time.RFC3339)},
				"finished_at":     []string{tc.e.FinishedAt.Format(time.RFC3339)},
				"notify_interval": []string{"1h"},
			}
			req := httptest.NewRequest("GET", "/events/create/?"+params.Encode(), nil)
			resp := httptest.NewRecorder()
			handler := server.newHandler()
			handler.ServeHTTP(resp, req)
			if status := resp.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
			require.Equal(t, tc.b, resp.Body.String())
		})
	}
}

func TestEventUpdate(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	server := NewServer(logg, a, "")
	cases := []struct {
		e storage.Event
		b string
	}{
		{
			e: storage.Event{
				ID:             "123-123",
				Title:          "title",
				UserId:         10,
				StartedAt:      time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
				FinishedAt:     time.Date(2022, 12, 1, 23, 0, 0, 0, time.UTC),
				NotifyInterval: time.Hour,
			},
			b: `{"Success":true,"Error":"","Data":null}`,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.e.Title, func(t *testing.T) {
			logg.EXPECT().Info(gomock.Any()).AnyTimes()
			logg.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
			dataKeeper.EXPECT().FindById(tc.e.ID).Return(&tc.e, nil)
			dataKeeper.EXPECT().SelectInterval(tc.e.StartedAt, tc.e.FinishedAt).Return(nil, nil)
			dataKeeper.EXPECT().UpdateEvent(&tc.e).Return(&tc.e, nil)
			params := url.Values{
				"id":              []string{tc.e.ID},
				"title":           []string{tc.e.Title},
				"user_id":         []string{fmt.Sprintf("%d", tc.e.UserId)},
				"started_at":      []string{tc.e.StartedAt.Format(time.RFC3339)},
				"finished_at":     []string{tc.e.FinishedAt.Format(time.RFC3339)},
				"notify_interval": []string{"1h"},
			}
			req := httptest.NewRequest("GET", "/events/edit/?"+params.Encode(), nil)
			resp := httptest.NewRecorder()
			handler := server.newHandler()
			handler.ServeHTTP(resp, req)
			if status := resp.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
			require.Equal(t, tc.b, resp.Body.String())
		})
	}
}

func TestEventDelete(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	server := NewServer(logg, a, "")
	cases := []struct {
		e storage.Event
		b string
	}{
		{
			e: storage.Event{
				ID:             "123-123",
				Title:          "title",
				UserId:         10,
				StartedAt:      time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
				FinishedAt:     time.Date(2022, 12, 1, 23, 0, 0, 0, time.UTC),
				NotifyInterval: time.Hour,
			},
			b: `{"Success":true,"Error":"","Data":{"Event":{"Id":"123-123","Title":"title","StartedAt":"2022-12-01T00:00:00Z","FinishedAt":"2022-12-01T23:00:00Z","Description":"","UserId":10,"NotifyInterval":3600000000000}}}`,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.e.Title, func(t *testing.T) {
			logg.EXPECT().Info(gomock.Any()).AnyTimes()
			logg.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
			dataKeeper.EXPECT().FindById(tc.e.ID).Return(&tc.e, nil)
			dataKeeper.EXPECT().DeleteEventById(tc.e.ID).Return(nil)
			params := url.Values{
				"id": []string{tc.e.ID},
			}
			req := httptest.NewRequest("GET", "/events/delete/?"+params.Encode(), nil)
			resp := httptest.NewRecorder()
			handler := server.newHandler()
			handler.ServeHTTP(resp, req)
			if status := resp.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
			require.Equal(t, tc.b, resp.Body.String())
		})
	}
}

func TestFindEvents(t *testing.T) {
	ctr := gomock.NewController(t)
	logg := mocks.NewMockLogger(ctr)
	dataKeeper := mocks.NewMockDataKeeper(ctr)
	a := app.New(logg, dataKeeper)
	server := NewServer(logg, a, "")
	cases := []struct {
		e []*storage.Event
		f time.Time
		t time.Time
		b string
	}{
		{
			e: []*storage.Event{
				&storage.Event{
					ID:             "123-123",
					Title:          "title",
					UserId:         10,
					StartedAt:      time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
					FinishedAt:     time.Date(2022, 12, 1, 23, 0, 0, 0, time.UTC),
					NotifyInterval: time.Hour,
				},
				&storage.Event{
					ID:             "321-456",
					Title:          "title",
					UserId:         10,
					StartedAt:      time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
					FinishedAt:     time.Date(2022, 12, 1, 23, 0, 0, 0, time.UTC),
					NotifyInterval: time.Hour,
				},
			},
			f: time.Date(2022, 11, 30, 0, 0, 0, 0, time.UTC),
			t: time.Date(2022, 12, 2, 23, 0, 0, 0, time.UTC),
			b: `{"Success":true,"Error":"","Data":{"Events":[{"Id":"123-123","Title":"title","StartedAt":"2022-12-01T00:00:00Z","FinishedAt":"2022-12-01T23:00:00Z","Description":"","UserId":10,"NotifyInterval":3600000000000},{"Id":"321-456","Title":"title","StartedAt":"2022-12-01T00:00:00Z","FinishedAt":"2022-12-01T23:00:00Z","Description":"","UserId":10,"NotifyInterval":3600000000000}]}}`,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run("list", func(t *testing.T) {
			logg.EXPECT().Info(gomock.Any()).AnyTimes()
			logg.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
			dataKeeper.EXPECT().SelectInterval(tc.f, tc.t).Return(tc.e, nil)
			params := url.Values{
				"from_date": []string{tc.f.Format(time.RFC3339)},
				"to_date":   []string{tc.t.Format(time.RFC3339)},
			}
			req := httptest.NewRequest("GET", "/events/list/?"+params.Encode(), nil)
			resp := httptest.NewRecorder()
			handler := server.newHandler()
			handler.ServeHTTP(resp, req)
			if status := resp.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
			require.Equal(t, tc.b, resp.Body.String())
		})
	}
}
