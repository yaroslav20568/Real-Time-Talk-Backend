package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	v1 "gin-real-time-talk/internal/controller/http/v1"
	"gin-real-time-talk/pkg/httpserver"
	"gin-real-time-talk/pkg/logger"
	"gin-real-time-talk/pkg/postgres"
	"gin-real-time-talk/pkg/validator"
)

func Run() error {
	validator.Init()
	logger := logger.New()

	db, err := postgres.New()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize database: %v", err))
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := Migrate(db); err != nil {
		logger.Error(fmt.Sprintf("Failed to run migrations: %v", err))
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	handler := v1.NewRouter(db, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	httpServer := httpserver.New(handler, httpserver.Port(port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info(fmt.Sprintf("signal: %s", s.String()))
	case err = <-httpServer.Notify():
		logger.Error(fmt.Sprintf("httpServer.Notify: %v", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Sprintf("httpServer.Shutdown: %v", err))
		return err
	}

	return nil
}
