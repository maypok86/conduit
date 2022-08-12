package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/controller/http/middleware"
)

// UserService is a user service interface.
type UserService interface{}

type userHandler struct {
	userService UserService
}

func newUserHandler(router *gin.Engine, authMiddleware middleware.Auth, userService UserService) {
	handler := userHandler{
		userService: userService,
	}

	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/", handler.createUser)
		usersGroup.POST("/login", handler.loginUser)
	}

	userGroup := router.Group("/user", authMiddleware.Handle)
	{
		userGroup.GET("/", handler.getCurrentUser)
		userGroup.PUT("/", handler.updateCurrentUser)
	}
}

func (h userHandler) createUser(c *gin.Context) {
}

func (h userHandler) loginUser(c *gin.Context) {
}

func (h userHandler) getCurrentUser(c *gin.Context) {
}

func (h userHandler) updateCurrentUser(c *gin.Context) {
}
