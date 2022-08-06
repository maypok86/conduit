package logger_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/maypok86/conduit/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	l := logger.New(&b, logger.WithLevel("debug"))
	msg := "Logger init"
	l.Info(msg)

	require.True(t, strings.Contains(b.String(), msg))
}
