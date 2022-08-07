package main

import (
	"os"

	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/pkg/logger"
	"go.uber.org/zap"
)

var (
	version   string
	buildDate string
)

func main() {
	run()
}

func run() {
	cfg := config.Get()
	l := logger.New(os.Stdout, logger.WithLevel(cfg.Logger.Level))
	l.Info("conduit", zap.String("version", version), zap.String("build_date", buildDate))
}
