package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/conduit/internal/controller/http/httperr"
	"github.com/maypok86/conduit/internal/controller/http/middleware"
	"github.com/maypok86/conduit/internal/domain/profile"
	"github.com/maypok86/conduit/pkg/logger"
)

//go:generate mockgen -source=profile.go -destination=mocks/profile_test.go -package=handler_test

// ProfileService is a profile service.
type ProfileService interface {
	GetByUsername(ctx context.Context, username string) (profile.Profile, error)
	GetWithFollow(ctx context.Context, email, username string) (profile.Profile, error)
	Follow(ctx context.Context, email, username string) (profile.Profile, error)
}

type profileHandler struct {
	authMiddleware middleware.Auth
	profileService ProfileService
}

type profileDeps struct {
	router         *gin.RouterGroup
	authMiddleware middleware.Auth
	profileService ProfileService
}

func newProfileHandler(deps profileDeps) {
	handler := profileHandler{
		authMiddleware: deps.authMiddleware,
		profileService: deps.profileService,
	}

	deps.router.GET("/profiles/:username", deps.authMiddleware.OptionalHandle, handler.getProfile)

	profilesGroup := deps.router.Group("/profiles", deps.authMiddleware.Handle)
	{
		profilesGroup.POST("/:username/follow", handler.follow)
		profilesGroup.DELETE("/:username/unfollow", handler.unfollow)
	}
}

type getProfileRequest struct {
	Username string `uri:"username" binding:"required"`
}

type getProfileResponse struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func (h profileHandler) getProfile(c *gin.Context) {
	var request getProfileRequest
	if err := c.ShouldBindUri(&request); err != nil {
		httperr.BadRequest(c, "invalid-request", err)
		return
	}

	username := request.Username
	payload := h.authMiddleware.GetPayload(c)

	if payload == nil {
		profileEntity, err := h.profileService.GetByUsername(logger.FromRequestToContext(c), username)
		if err != nil {
			httperr.RespondWithSlugError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"profile": getProfileResponse{
				Username:  profileEntity.Username,
				Bio:       profileEntity.GetBio(),
				Image:     profileEntity.GetImage(),
				Following: profileEntity.Following,
			},
		})
	} else {
		profileEntity, err := h.profileService.GetWithFollow(logger.FromRequestToContext(c), payload.Email, username)
		if err != nil {
			httperr.RespondWithSlugError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"profile": getProfileResponse{
				Username:  profileEntity.Username,
				Bio:       profileEntity.GetBio(),
				Image:     profileEntity.GetImage(),
				Following: profileEntity.Following,
			},
		})
	}
}

type followRequest struct {
	Username string `uri:"username" binding:"required"`
}

type followResponse struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func (h profileHandler) follow(c *gin.Context) {
	var request followRequest
	if err := c.ShouldBindUri(&request); err != nil {
		httperr.BadRequest(c, "invalid-request", err)
		return
	}

	payload := h.authMiddleware.GetPayload(c)

	profileEntity, err := h.profileService.Follow(logger.FromRequestToContext(c), payload.Email, request.Username)
	if err != nil {
		httperr.RespondWithSlugError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": followResponse{
			Username:  profileEntity.Username,
			Bio:       profileEntity.GetBio(),
			Image:     profileEntity.GetImage(),
			Following: profileEntity.Following,
		},
	})
}

func (h profileHandler) unfollow(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "unfollow",
	})
}
