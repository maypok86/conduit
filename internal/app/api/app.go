// Package api provides a application interface.
package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/internal/controller/http"
	"github.com/maypok86/conduit/pkg/httpserver"
	"go.uber.org/zap"
)

// App is a application interface.
type App struct {
	logger     *zap.Logger
	httpServer httpserver.Server
}

// New creates a new App.
func New(ctx context.Context, logger *zap.Logger) (App, error) {
	cfg := config.Get()
	handler := http.NewHandler()

	return App{
		logger: logger,
		httpServer: httpserver.New(
			handler.Init(logger),
			httpserver.WithHost(cfg.HTTP.Host),
			httpserver.WithPort(cfg.HTTP.Port),
			httpserver.WithMaxHeaderBytes(cfg.HTTP.MaxHeaderBytes),
			httpserver.WithReadTimeout(cfg.HTTP.ReadTimeout),
			httpserver.WithWriteTimeout(cfg.HTTP.WriteTimeout),
		),
	}, nil
}

// Run runs the application.
func (a App) Run(ctx context.Context) error {
	eChan := make(chan error)
	interrupt := make(chan os.Signal, 1)

	a.logger.Info("Http server is starting")

	go func() {
		if err := a.httpServer.Start(); err != nil {
			eChan <- fmt.Errorf("failed to listen and serve: %w", err)
		}
	}()

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-eChan:
		return fmt.Errorf("conduit started failed: %w", err)
	case <-interrupt:
	}

	const httpShutdownTimeout = 5 * time.Second
	if err := a.httpServer.Stop(ctx, httpShutdownTimeout); err != nil {
		return fmt.Errorf("failed to stop http server: %w", err)
	}

	return nil
}
