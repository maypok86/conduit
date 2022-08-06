package logger_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/maypok86/conduit/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLoggerContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	log := logger.FromContext(ctx)
	require.Equal(t, zap.L(), log)

	var buffer bytes.Buffer
	log = logger.New(&buffer)
	ctx = logger.ContextWithLogger(ctx, log)
	ctxLogger := logger.FromContext(ctx)
	require.Equal(t, log, ctxLogger)
	require.NotEqual(t, zap.L(), ctxLogger)

	msg := "Success FromContext"
	ctxLogger.Info(msg)
	require.True(t, strings.Contains(buffer.String(), msg))
}
