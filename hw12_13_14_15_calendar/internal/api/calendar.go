package api

import (
	"context"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/bojik/otus-golang/hw12_13_14_15_calendar/pkg/calendarpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate protoc --proto_path=../../api/ --go_out=../../pkg/calendarpb --go-grpc_out=../../pkg/calendarpb ../../api/EventService.proto
type CalendarApi struct {
	app  *app.App
	logg logger.Logger
	pb.UnimplementedCalendarServer
}

func NewCalendarApi(app *app.App, logg logger.Logger) *CalendarApi {
	return &CalendarApi{
		app:  app,
		logg: logg,
	}
}

func (s *CalendarApi) InsertEvent(ctx context.Context, evt *pb.Event) (*pb.Event, error) {
	newEvent := app.Event{
		Title:          evt.Title,
		UserId:         int(evt.UserId),
		Description:    evt.Description,
		StartedAt:      evt.StartedAt.AsTime(),
		FinishedAt:     evt.FinishedAt.AsTime(),
		NotifyInterval: evt.NotifyInterval.AsDuration(),
	}
	id, err := s.app.CreateEvent(ctx, newEvent)
	if err != nil {
		return nil, err
	}
	dbEvent, err := s.app.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.fillEvent(dbEvent), nil
}

func (s *CalendarApi) UpdateEvent(ctx context.Context, evt *pb.Event) (*pb.Event, error) {
	event := app.Event{
		Id:             evt.Id,
		Title:          evt.Title,
		UserId:         int(evt.UserId),
		Description:    evt.Description,
		StartedAt:      evt.StartedAt.AsTime(),
		FinishedAt:     evt.FinishedAt.AsTime(),
		NotifyInterval: evt.NotifyInterval.AsDuration(),
	}
	if err := s.app.UpdateEvent(ctx, event); err != nil {
		return nil, err
	}
	dbEvent, err := s.app.FindById(ctx, event.Id)
	if err != nil {
		return nil, err
	}
	return s.fillEvent(dbEvent), nil
}

func (s *CalendarApi) FindEventById(ctx context.Context, id *pb.Id) (*pb.Event, error) {
	evt, err := s.app.FindById(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	event := s.fillEvent(evt)
	return event, nil
}

func (s *CalendarApi) DeleteEvent(ctx context.Context, id *pb.Id) (*pb.Event, error) {
	evt, err := s.app.DeleteById(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	event := s.fillEvent(evt)
	return event, nil
}

func (s *CalendarApi) FindEventsByInterval(ctx context.Context, interval *pb.Interval) (*pb.Events, error) {
	ets, err := s.app.FindEventsByInterval(ctx, interval.StartedAt.AsTime(), interval.FinishedAt.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.Events{Events: s.fillEvents(ets)}, nil
}

func (s *CalendarApi) FindDayEvents(ctx context.Context, date *timestamppb.Timestamp) (*pb.Events, error) {
	ets, err := s.app.FindDayEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.Events{Events: s.fillEvents(ets)}, nil
}
func (s *CalendarApi) FindWeekEvents(ctx context.Context, date *timestamppb.Timestamp) (*pb.Events, error) {
	ets, err := s.app.FindDayEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.Events{Events: s.fillEvents(ets)}, nil
}

func (s *CalendarApi) FindMonthEvents(ctx context.Context, date *timestamppb.Timestamp) (*pb.Events, error) {
	ets, err := s.app.FindDayEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.Events{Events: s.fillEvents(ets)}, nil
}

func (s *CalendarApi) fillEvent(evt *app.Event) *pb.Event {
	event := &pb.Event{
		Id:             evt.Id,
		Title:          evt.Title,
		StartedAt:      timestamppb.New(evt.StartedAt),
		FinishedAt:     timestamppb.New(evt.FinishedAt),
		Description:    evt.Description,
		UserId:         int32(evt.UserId),
		NotifyInterval: durationpb.New(evt.NotifyInterval),
	}
	return event
}

func (s *CalendarApi) fillEvents(evts []*app.Event) []*pb.Event {
	if len(evts) == 0 {
		return []*pb.Event{}
	}
	events := make([]*pb.Event, len(evts))
	for i, e := range evts {
		events[i] = s.fillEvent(e)
	}
	return events
}
