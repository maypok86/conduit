package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/config"
	"github.com/maypok86/conduit/internal/controller/http/httperr"
	"github.com/maypok86/conduit/internal/controller/http/middleware"
	"github.com/maypok86/conduit/internal/domain/user"
	"github.com/maypok86/conduit/pkg/logger"
	"go.uber.org/zap"
)

//go:generate mockgen -source=user.go -destination=mocks/user_test.go -package=handler_test

// ErrAtLeastOneFieldRequired is returned when at least one field is required to update user.
var ErrAtLeastOneFieldRequired = errors.New("at least one field in update current user request must be provided")

// UserService is a user service interface.
type UserService interface {
	Create(ctx context.Context, dto user.CreateDTO) (user.User, error)
	Login(ctx context.Context, email, password string) (user.User, error)
	GetByEmail(ctx context.Context, email string) (user.User, error)
	UpdateByEmail(ctx context.Context, email string, dto user.UpdateDTO) (user.User, error)
}

type userHandler struct {
	authMiddleware middleware.Auth
	userService    UserService
	tokenMaker     TokenMaker
}

type userDeps struct {
	router         *gin.RouterGroup
	authMiddleware middleware.Auth
	userService    UserService
	tokenMaker     TokenMaker
}

func newUserHandler(deps userDeps) {
	handler := userHandler{
		userService:    deps.userService,
		tokenMaker:     deps.tokenMaker,
		authMiddleware: deps.authMiddleware,
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
	User userRequest `json:"user" binding:"required"`
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

	userEntity, err := h.userService.Create(logger.FromRequestToContext(c), user.CreateDTO{
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

type loginUserRequest struct {
	User struct {
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	} `json:"user" binding:"required"`
}

type loginUserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func (h userHandler) loginUser(c *gin.Context) {
	var request loginUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperr.BadRequest(c, "invalid-request", err)
		return
	}

	userEntity, err := h.userService.Login(logger.FromRequestToContext(c), request.User.Email, request.User.Password)
	if err != nil {
		httperr.RespondWithSlugError(c, err)
		return
	}

	accessToken, err := h.tokenMaker.CreateToken(userEntity.Email, config.Get().Token.Expired)
	if err != nil {
		httperr.RespondWithSlugError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": loginUserResponse{
			Email:    userEntity.Email,
			Username: userEntity.Username,
			Bio:      userEntity.GetBio(),
			Image:    userEntity.GetImage(),
			Token:    accessToken,
		},
	})
}

type getCurrentUserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func (h userHandler) getCurrentUser(c *gin.Context) {
	payload := h.authMiddleware.GetPayload(c)
	if payload == nil {
		return
	}

	userEntity, err := h.userService.GetByEmail(logger.FromRequestToContext(c), payload.Email)
	if err != nil {
		httperr.RespondWithSlugError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": getCurrentUserResponse{
			Email:    userEntity.Email,
			Username: userEntity.Username,
			Bio:      userEntity.GetBio(),
			Image:    userEntity.GetImage(),
			Token:    h.authMiddleware.GetToken(c),
		},
	})
}

type updateCurrentUserRequest struct {
	User struct {
		Username *string `json:"username" binding:"omitempty,alphanum"`
		Email    *string `json:"email"    binding:"omitempty,email"`
		Token    *string `json:"token"    binding:"omitempty"`
		Bio      *string `json:"bio"      binding:"omitempty,max=1024"`
		Image    *string `json:"image"    binding:"omitempty,url"`
	} `json:"user" binding:"required"`
}

type updateCurrentUserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func (ucur updateCurrentUserRequest) validate() error {
	if ucur.User.Username != nil {
		return nil
	}

	if ucur.User.Email != nil {
		return nil
	}

	if ucur.User.Bio != nil {
		return nil
	}

	if ucur.User.Image != nil {
		return nil
	}

	return ErrAtLeastOneFieldRequired
}

func (h userHandler) updateCurrentUser(c *gin.Context) {
	payload := h.authMiddleware.GetPayload(c)
	if payload == nil {
		return
	}

	var request updateCurrentUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperr.BadRequest(c, "invalid-request", err)
		return
	}

	if err := request.validate(); err != nil {
		logger.FromRequest(c).Info("validate", zap.Any("request", request))
		httperr.RespondWithSlugError(c, err)

		return
	}

	userEntity, err := h.userService.UpdateByEmail(logger.FromRequestToContext(c), payload.Email, user.UpdateDTO{
		Username: request.User.Username,
		Email:    request.User.Email,
		Bio:      request.User.Bio,
		Image:    request.User.Image,
	})
	if err != nil {
		httperr.RespondWithSlugError(c, err)
		return
	}

	token := h.authMiddleware.GetToken(c)
	if request.User.Token != nil {
		token = *request.User.Token
	}

	c.JSON(http.StatusOK, gin.H{
		"user": updateCurrentUserResponse{
			Email:    userEntity.Email,
			Username: userEntity.Username,
			Bio:      userEntity.GetBio(),
			Image:    userEntity.GetImage(),
			Token:    token,
		},
	})
}
