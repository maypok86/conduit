// Package http provides a http handler.
package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/pkg/token"
	"go.uber.org/zap"
)

type tokenMaker interface {
	CreateToken(string, time.Duration) (string, error)
	VerifyToken(token string) (*token.Payload, error)
}

// Handler is a http handler.
type Handler struct {
	tokenMaker   tokenMaker
	tokenExpired time.Duration
}

// New creates a new Handler.
func New(tokenMaker tokenMaker, tokenExpired time.Duration) Handler {
	return Handler{
		tokenMaker:   tokenMaker,
		tokenExpired: tokenExpired,
	}
}

// Init initializes the http routes.
func (h Handler) Init(_ *zap.Logger) http.Handler {
	router := gin.New()

	if config.Get().IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	setMiddlewares(router, nil)

	return router
}

func setMiddlewares(router *gin.Engine, _ *zap.Logger) {
	router.Use(gin.Recovery())
}
