package logger_test

import (
	"testing"

	"github.com/maypok86/conduit/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLevel_ZapLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lvl      string
		zapLevel zapcore.Level
	}{
		{
			name:     "debug level",
			lvl:      "debug",
			zapLevel: zapcore.DebugLevel,
		},
		{
			name:     "info level",
			lvl:      "info",
			zapLevel: zapcore.InfoLevel,
		},
		{
			name:     "warn level",
			lvl:      "warn",
			zapLevel: zapcore.WarnLevel,
		},
		{
			name:     "error level",
			lvl:      "error",
			zapLevel: zapcore.ErrorLevel,
		},
		{
			name:     "dpanic level",
			lvl:      "dpanic",
			zapLevel: zapcore.DPanicLevel,
		},
		{
			name:     "panic level",
			lvl:      "panic",
			zapLevel: zapcore.PanicLevel,
		},
		{
			name:     "fatal level",
			lvl:      "fatal",
			zapLevel: zapcore.FatalLevel,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			level := logger.NewLevel(tt.lvl)
			require.Equal(t, tt.zapLevel, level.ZapLevel())
		})
	}
}

func TestLevel_ZapLevelPanic(t *testing.T) {
	t.Parallel()

	level := logger.NewLevel("invalid")

	require.Panics(t, func() {
		level.ZapLevel()
	}, "invalid log level")
}
