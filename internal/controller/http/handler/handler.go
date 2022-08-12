// Package handler provides a http handler.
package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/internal/controller/http/middleware"
	"github.com/maypok86/conduit/pkg/token"
	"go.uber.org/zap"
)

type TokenMaker interface {
	CreateToken(string, time.Duration) (string, error)
	VerifyToken(token string) (*token.Payload, error)
}

// Deps is a http handler dependencies.
type Deps struct {
	TokenMaker  TokenMaker
	Logger      *zap.Logger
	UserService UserService
}

// NewRouter returns a new http router.
func NewRouter(deps Deps) *gin.Engine {
	router := gin.New()

	if config.Get().IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	middleware.ApplyMiddlewares(router, deps.Logger)

	api := router.Group("/api")
	{
		authMiddleware := middleware.NewAuth(deps.TokenMaker)

		newUserHandler(api, authMiddleware, deps.UserService, deps.TokenMaker)
	}

	return router
}
