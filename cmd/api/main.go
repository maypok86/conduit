package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/maypok86/conduit/internal/app/api"
	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/pkg/logger"
	"go.uber.org/zap"
)

var (
	version   string
	buildDate string
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg := config.Get()
	l := logger.New(os.Stdout, logger.WithLevel(cfg.Logger.Level))
	l.Info("conduit", zap.String("version", version), zap.String("build_date", buildDate))

	app, err := api.New(ctx, l)
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}
