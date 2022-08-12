package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/pkg/logger"
	"go.uber.org/zap"
)

func logMiddleware(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := []zap.Field{
			zap.String("endpoint", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("remote_addr", c.Request.RemoteAddr),
		}
		logger.RequestWithLogger(c, l.With(fields...))
		c.Next()
	}
}
