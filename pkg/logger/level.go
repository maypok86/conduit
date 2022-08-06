package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type level struct {
	lvl string
}

func (l *level) zapLevel() zapcore.Level {
	switch l.lvl {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		panic("invalid log level")
	}
}

// Option is a function that can be passed to New to customize the logger.
type Option func(*level)

// WithLevel option sets the log level.
func WithLevel(lvl string) Option {
	return func(l *level) {
		l.lvl = lvl
	}
}
