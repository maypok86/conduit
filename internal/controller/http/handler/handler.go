// Package http provides a http handler.
package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/internal/controller/http/middleware"
	"github.com/maypok86/conduit/pkg/token"
	"go.uber.org/zap"
)

type tokenMaker interface {
	CreateToken(string, time.Duration) (string, error)
	VerifyToken(token string) (*token.Payload, error)
}

// Handler is a http handler.
type Handler struct {
	middlewareManager middleware.Manager
}

// New creates a new Handler.
func New(tokenMaker tokenMaker, tokenExpired time.Duration, logger *zap.Logger) Handler {
	return Handler{
		middlewareManager: middleware.NewManager(tokenMaker, tokenExpired, logger),
	}
}

// Init initializes the http routes.
func (h Handler) Init() http.Handler {
	router := gin.New()

	if config.Get().IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	h.middlewareManager.ApplyMiddlewares(router)

	return router
}
