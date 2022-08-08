// Package http provides a http handler.
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/maypok86/conduit/pkg/logger"
	"go.uber.org/zap"
)

// Handler is a http handler.
type Handler struct{}

// NewHandler creates a new Handler.
func NewHandler() Handler {
	return Handler{}
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
