// Package http provides a http handler.
package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/maypok86/conduit/pkg/logger"
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

// NewHandler creates a new Handler.
func NewHandler(tokenMaker tokenMaker, tokenExpired time.Duration) Handler {
	return Handler{
		tokenMaker:   tokenMaker,
		tokenExpired: tokenExpired,
	}
}

// Init initializes the http routes.
func (h Handler) Init(l *zap.Logger) http.Handler {
	router := chi.NewRouter()

	setMiddlewares(router, l)

	return router
}

func setMiddlewares(router *chi.Mux, l *zap.Logger) {
	router.Use(setLoggerMiddleware(l))
}

func setLoggerMiddleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.ContextWithLogger(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
