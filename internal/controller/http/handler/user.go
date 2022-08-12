package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/internal/controller/http/httperr"
	"github.com/maypok86/conduit/internal/controller/http/middleware"
	"github.com/maypok86/conduit/internal/domain/user"
	"github.com/maypok86/conduit/pkg/logger"
)

// UserService is a user service interface.
type UserService interface {
	CreateUser(ctx context.Context, dto user.CreateDTO) (user.User, error)
}

type userHandler struct {
	userService UserService
	tokenMaker  TokenMaker
}

type userDeps struct {
	router         *gin.RouterGroup
	authMiddleware middleware.Auth
	userService    UserService
	tokenMaker     TokenMaker
}

func newUserHandler(deps userDeps) {
	handler := userHandler{
		userService: deps.userService,
		tokenMaker:  deps.tokenMaker,
	}

	usersGroup := deps.router.Group("/users")
	{
		usersGroup.POST("/", handler.createUser)
		usersGroup.POST("/login", handler.loginUser)
	}

	userGroup := deps.router.Group("/user", deps.authMiddleware.Handle)
	{
		userGroup.GET("/", handler.getCurrentUser)
		userGroup.PUT("/", handler.updateCurrentUser)
	}
}

type userRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type createUserRequest struct {
	User userRequest `json:"user"`
}

type createUserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func (h userHandler) createUser(c *gin.Context) {
	var request createUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperr.BadRequest(c, "invalid-request", err)
		return
	}

	userEntity, err := h.userService.CreateUser(logger.FromRequestToContext(c), user.CreateDTO{
		Email:    request.User.Email,
		Username: request.User.Username,
		Password: request.User.Password,
	})
	if err != nil {
		httperr.RespondWithSlugError(c, err)
		return
	}

	token, err := h.tokenMaker.CreateToken(userEntity.Email, config.Get().Token.Expired)
	if err != nil {
		httperr.RespondWithSlugError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": createUserResponse{
			Email:    userEntity.Email,
			Username: userEntity.Username,
			Bio:      userEntity.GetBio(),
			Image:    userEntity.GetImage(),
			Token:    token,
		},
	})
}

func (h userHandler) loginUser(c *gin.Context) {
}

func (h userHandler) getCurrentUser(c *gin.Context) {
}

func (h userHandler) updateCurrentUser(c *gin.Context) {
}
