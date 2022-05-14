package api

import (
	"context"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/bojik/otus-golang/hw12_13_14_15_calendar/pkg/calendarpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//nolint:lll
//go:generate protoc --proto_path=../../api/ --go_out=../../pkg/calendarpb --go-grpc_out=../../pkg/calendarpb ../../api/EventService.proto
type CalendarAPI struct {
	app  *app.App
	logg logger.Logger
	pb.UnimplementedCalendarServer
}

func NewCalendarAPI(app *app.App, logg logger.Logger) *CalendarAPI {
	return &CalendarAPI{
		app:  app,
		logg: logg,
	}
}

func (s *CalendarAPI) InsertEvent(ctx context.Context, evt *pb.Event) (*pb.Event, error) {
	newEvent := app.Event{
		Title:          evt.Title,
		UserID:         int(evt.UserId),
		Description:    evt.Description,
		StartedAt:      evt.StartedAt.AsTime(),
		FinishedAt:     evt.FinishedAt.AsTime(),
		NotifyInterval: evt.NotifyInterval.AsDuration(),
	}
	id, err := s.app.CreateEvent(ctx, newEvent)
	if err != nil {
		return nil, err
	}
	dbEvent, err := s.app.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.fillEvent(dbEvent), nil
}

func (s *CalendarAPI) UpdateEvent(ctx context.Context, evt *pb.Event) (*pb.Event, error) {
	event := app.Event{
		ID:             evt.Id,
		Title:          evt.Title,
		UserID:         int(evt.UserId),
		Description:    evt.Description,
		StartedAt:      evt.StartedAt.AsTime(),
		FinishedAt:     evt.FinishedAt.AsTime(),
		NotifyInterval: evt.NotifyInterval.AsDuration(),
	}
	if err := s.app.UpdateEvent(ctx, event); err != nil {
		return nil, err
	}
	dbEvent, err := s.app.FindByID(ctx, event.ID)
	if err != nil {
		return nil, err
	}
	return s.fillEvent(dbEvent), nil
}

func (s *CalendarAPI) FindEventByID(ctx context.Context, id *pb.Id) (*pb.Event, error) {
	evt, err := s.app.FindByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	event := s.fillEvent(evt)
	return event, nil
}

func (s *CalendarAPI) DeleteEvent(ctx context.Context, id *pb.Id) (*pb.Event, error) {
	evt, err := s.app.DeleteByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	event := s.fillEvent(evt)
	return event, nil
}

func (s *CalendarAPI) FindEventsByInterval(ctx context.Context, interval *pb.Interval) (*pb.Events, error) {
	ets, err := s.app.FindEventsByInterval(ctx, interval.StartedAt.AsTime(), interval.FinishedAt.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.Events{Events: s.fillEvents(ets)}, nil
}

func (s *CalendarAPI) FindDayEvents(ctx context.Context, date *timestamppb.Timestamp) (*pb.Events, error) {
	ets, err := s.app.FindDayEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.Events{Events: s.fillEvents(ets)}, nil
}

func (s *CalendarAPI) FindWeekEvents(ctx context.Context, date *timestamppb.Timestamp) (*pb.Events, error) {
	ets, err := s.app.FindDayEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.Events{Events: s.fillEvents(ets)}, nil
}

func (s *CalendarAPI) FindMonthEvents(ctx context.Context, date *timestamppb.Timestamp) (*pb.Events, error) {
	ets, err := s.app.FindDayEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.Events{Events: s.fillEvents(ets)}, nil
}

func (s *CalendarAPI) fillEvent(evt *app.Event) *pb.Event {
	event := &pb.Event{
		Id:             evt.ID,
		Title:          evt.Title,
		StartedAt:      timestamppb.New(evt.StartedAt),
		FinishedAt:     timestamppb.New(evt.FinishedAt),
		Description:    evt.Description,
		UserId:         int32(evt.UserID),
		NotifyInterval: durationpb.New(evt.NotifyInterval),
	}
	return event
}

func (s *CalendarAPI) fillEvents(evts []*app.Event) []*pb.Event {
	if len(evts) == 0 {
		return []*pb.Event{}
	}
	events := make([]*pb.Event, len(evts))
	for i, e := range evts {
		events[i] = s.fillEvent(e)
	}
	return events
}
