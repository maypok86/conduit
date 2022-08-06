// Package logger provides a wrapper for the zap logger.
package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultLogLevel = "info"

// New creates a new logger.
func New(console io.Writer, opts ...Option) *zap.Logger {
	prodCfg := zap.NewProductionEncoderConfig()

	prodCfg.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(caller.TrimmedPath())
		encoder.AppendString("\t|")
	}

	prodCfg.EncodeTime = zapcore.TimeEncoderOfLayout("02/01 15:04:05") // "02/01/2006 15:04:05 |"
	prodCfg.ConsoleSeparator = " "
	prodCfg.EncodeName = func(n string, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(n)
		enc.AppendString("|")
	}

	prodCfg.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("|")
		enc.AppendString(l.CapitalString())
		enc.AppendString("|")
	}

	consoleEncoder := zapcore.NewConsoleEncoder(prodCfg)

	logLevel := &level{lvl: defaultLogLevel}
	for _, opt := range opts {
		opt(logLevel)
	}

	core := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(console),
		logLevel.zapLevel(),
	)

	return zap.New(
		core,
		zap.AddCaller(),
	)
}
