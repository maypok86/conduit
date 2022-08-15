// Package handler provides a http handler.
package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/internal/controller/http/middleware"
	"github.com/maypok86/conduit/internal/domain"
	"github.com/maypok86/conduit/pkg/token"
	"go.uber.org/zap"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler_test.go -package=handler_test

// TokenMaker is a token maker.
type TokenMaker interface {
	CreateToken(email string, duration time.Duration) (string, error)
	VerifyToken(accessToken string) (*token.Payload, error)
}

// Deps is a http handler dependencies.
type Deps struct {
	TokenMaker TokenMaker
	Logger     *zap.Logger
	Services   domain.Services
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

		newUserHandler(userDeps{
			router:         api,
			authMiddleware: authMiddleware,
			userService:    deps.Services.User,
			tokenMaker:     deps.TokenMaker,
		})

		newProfileHandler(profileDeps{
			router:         api,
			authMiddleware: authMiddleware,
			profileService: deps.Services.Profile,
		})
	}

	return router
}
