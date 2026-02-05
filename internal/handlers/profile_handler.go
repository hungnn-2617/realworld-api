package handlers

import (
	"net/http"

	"realworld-api/internal/dtos"
	appErrors "realworld-api/internal/errors"
	"realworld-api/internal/helpers"
	"realworld-api/internal/middleware"
	"realworld-api/internal/services"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile handles GET /api/profiles/:username
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	username := c.Param("username")

	// Get current user ID if authenticated (optional)
	currentUserID, _ := middleware.GetUserID(c)

	profile, err := h.profileService.GetProfile(username, currentUserID)
	if err != nil {
		if appErrors.IsUserNotFound(err) {
			helpers.NotFound(c, appErrors.ErrUserNotFound)
		} else {
			helpers.InternalServerError(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dtos.ProfileResponse{
		Profile: *profile,
	})
}

// FollowUser handles POST /api/profiles/:username/follow
func (h *ProfileHandler) FollowUser(c *gin.Context) {
	currentUserID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	username := c.Param("username")

	profile, err := h.profileService.FollowUser(username, currentUserID)
	if err != nil {
		if appErrors.IsUserNotFound(err) {
			helpers.NotFound(c, appErrors.ErrUserNotFound)
		} else if appErrors.IsCannotFollowYourself(err) {
			helpers.UnprocessableEntity(c, appErrors.ErrCannotFollowYourself)
		} else {
			helpers.UnprocessableEntity(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dtos.ProfileResponse{
		Profile: *profile,
	})
}

// UnfollowUser handles DELETE /api/profiles/:username/follow
func (h *ProfileHandler) UnfollowUser(c *gin.Context) {
	currentUserID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	username := c.Param("username")

	profile, err := h.profileService.UnfollowUser(username, currentUserID)
	if err != nil {
		if appErrors.IsUserNotFound(err) {
			helpers.NotFound(c, appErrors.ErrUserNotFound)
		} else if appErrors.IsCannotFollowYourself(err) {
			helpers.UnprocessableEntity(c, appErrors.ErrCannotFollowYourself)
		} else {
			helpers.UnprocessableEntity(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dtos.ProfileResponse{
		Profile: *profile,
	})
}
