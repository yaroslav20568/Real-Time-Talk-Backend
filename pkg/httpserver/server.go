package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gin-real-time-talk/config"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		if port == "" {
			port = config.Env.App.Port
		}
		s.server.Addr = ":" + port
	}
}

func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		server: &http.Server{
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.server.Addr == "" {
		s.server.Addr = ":" + config.Env.App.Port
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
