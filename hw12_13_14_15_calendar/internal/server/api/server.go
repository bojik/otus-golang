package internalapi

import (
	"context"
	"net"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/api"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/pkg/calendarpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	addr string
	logg logger.Logger
	srv  *grpc.Server
}

func NewServer(addr string, api *api.CalendarApi, logg logger.Logger) *Server {
	s := &Server{
		addr: addr,
		logg: logg,
	}
	s.srv = grpc.NewServer(
		grpc.UnaryInterceptor(loggingMiddleware(logg)),
	)
	reflection.Register(s.srv)
	calendarpb.RegisterCalendarServer(s.srv, api)
	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("starting api server: " + s.addr)
	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", s.addr)
	if err != nil {
		return err
	}
	if err := s.srv.Serve(listener); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	doneCh := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		close(doneCh)
	}()
	select {
	case <-doneCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
