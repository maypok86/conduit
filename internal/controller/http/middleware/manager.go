// Package middleware provides a http middleware.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Manager is a http middleware manager.
type Manager struct {
	log  logMiddleware
	auth authMiddleware
}

// NewManager creates a new http middleware manager.
func NewManager(tokenMaker tokenMaker, tokenExpired time.Duration, logger *zap.Logger) Manager {
	return Manager{
		log:  newLogMiddleware(logger),
		auth: newAuthMiddleware(tokenMaker, tokenExpired),
	}
}

// ApplyMiddlewares applies middlewares to the given router.
func (m Manager) ApplyMiddlewares(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(m.log.Handle)
	router.Use(m.auth.Handle)
}
