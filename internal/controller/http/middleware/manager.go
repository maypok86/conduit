// Package middleware provides a http middleware.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/pkg/token"
	"go.uber.org/zap"
)

// Manager is a http middleware manager.
type Manager struct {
	log  logMiddleware
	auth authMiddleware
}

// NewManager creates a new http middleware manager.
func NewManager(tokenMaker tokenMaker, logger *zap.Logger) Manager {
	return Manager{
		log:  newLogMiddleware(logger),
		auth: newAuthMiddleware(tokenMaker),
	}
}

// ApplyMiddlewares applies middlewares to the given router.
func (m Manager) ApplyMiddlewares(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(m.log.Handle)
	router.Use(m.auth.Handle)
}

// GetPayload returns the authorization payload from the context.
func (m Manager) GetPayload(c *gin.Context) *token.Payload {
	return m.auth.getPayload(c)
}
