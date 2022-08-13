package logger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLoggerRequest(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	want := zap.L().Named("test")

	logger.RequestWithLogger(c, want)

	got := logger.FromRequest(c)

	require.Equal(t, want, got)
	require.NotEqual(t, zap.L(), got)
}

func TestLoggerRequestDefault(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	got := logger.FromRequest(c)

	require.Equal(t, zap.L(), got)
}

func TestLoggerRequestToContext(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	want := zap.L().Named("test")
	logger.RequestWithLogger(c, want)

	ctx := logger.FromRequestToContext(c)
	got := logger.FromContext(ctx)

	require.Equal(t, want, got)
	require.NotEqual(t, zap.L(), got)
}
