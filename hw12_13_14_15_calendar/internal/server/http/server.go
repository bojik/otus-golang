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

func (s *Server) newHandler(ctx context.Context) http.Handler {
	serverMux := http.NewServeMux()
	serverMux.Handle("/", loggingMiddleware(s.logg, NewDefaultController(s.logg)))
	return serverMux
}

func (s *Server) Start(ctx context.Context) error {
	s.srv = &http.Server{
		Addr:           s.addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.newHandler(ctx),
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
