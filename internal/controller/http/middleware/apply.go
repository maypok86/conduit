// Package middleware provides a http middleware.
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ApplyMiddlewares applies middlewares to the given router.
func ApplyMiddlewares(router *gin.Engine, l *zap.Logger) {
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(logMiddleware(l))
}
