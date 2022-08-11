package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/pkg/logger"
	"go.uber.org/zap"
)

type logMiddleware struct {
	logger *zap.Logger
}

func newLogMiddleware(logger *zap.Logger) logMiddleware {
	return logMiddleware{
		logger: logger,
	}
}

func (lm logMiddleware) Handle(c *gin.Context) {
	fields := []zap.Field{
		zap.String("endpoint", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("remote_addr", c.Request.RemoteAddr),
	}
	logger.RequestWithLogger(c, lm.logger.With(fields...))
	c.Next()
}
