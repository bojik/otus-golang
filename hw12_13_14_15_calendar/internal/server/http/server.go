package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	logg logger.Logger
	app  *app.App
	addr string
	srv  *http.Server
}

func NewServer(logger logger.Logger, app *app.App, addr string) *Server {
	return &Server{
		logg: logger,
		app:  app,
		addr: addr,
	}
}

func (s *Server) newHandler() http.Handler {
	loggingMiddleware := newLoggingMiddleware(s.logg)
	serverMux := http.NewServeMux()
	serverMux.Handle("/", loggingMiddleware(NewIndexHandler(s.logg)))
	serverMux.Handle("/events/create/", loggingMiddleware(NewCreateEventHandler(s.app, s.logg)))
	serverMux.Handle("/events/edit/", loggingMiddleware(NewUpdateEventHandler(s.app, s.logg)))
	serverMux.Handle("/events/get/", loggingMiddleware(NewGetEventHandler(s.app, s.logg)))
	serverMux.Handle("/events/delete/", loggingMiddleware(NewDeleteEventHandler(s.app, s.logg)))
	serverMux.Handle("/events/list/", loggingMiddleware(NewFindEventsHandler(s.app, s.logg)))
	serverMux.Handle("/events/list/day/", loggingMiddleware(NewFindEventsDayHandler(s.app, s.logg)))
	serverMux.Handle("/events/list/week/", loggingMiddleware(NewFindEventsWeekHandler(s.app, s.logg)))
	serverMux.Handle("/events/list/month/", loggingMiddleware(NewFindEventsMonthHandler(s.app, s.logg)))
	return serverMux
}

func (s *Server) Start(ctx context.Context) error {
	s.srv = &http.Server{
		Addr:           s.addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.newHandler(),
	}
	s.logg.Info("Starting HTTP at " + s.addr)
	err := s.srv.ListenAndServe()
	if err != nil {
		s.logg.Error(err.Error())
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
