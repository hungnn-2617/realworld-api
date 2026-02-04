package handlers

import (
	"net/http"

	"realworld-api/internal/dtos"
	appErrors "realworld-api/internal/errors"
	"realworld-api/internal/helpers"
	"realworld-api/internal/models"
	"realworld-api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, token, err := h.userService.Register(&req)
	if err != nil {
		helpers.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusCreated, h.toUserResponse(user, token))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, token, err := h.userService.Login(&req)
	if err != nil {
		helpers.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, h.toUserResponse(user, token))
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.Error(c, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	user, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		helpers.Error(c, http.StatusNotFound, appErrors.ErrUserNotFound)
		return
	}

	token, _ := h.userService.GenerateToken(user)
	helpers.SuccessResponse(c, http.StatusOK, h.toUserResponse(user, token))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.Error(c, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	var req dtos.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		helpers.Error(c, http.StatusNotFound, appErrors.ErrUserNotFound)
		return
	}

	if err := h.userService.Update(user, &req); err != nil {
		helpers.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	token, _ := h.userService.GenerateToken(user)
	helpers.SuccessResponse(c, http.StatusOK, h.toUserResponse(user, token))
}

func (h *UserHandler) toUserResponse(user *models.User, token string) dtos.UserResponse {
	return dtos.UserResponse{
		User: dtos.UserData{
			Email:    user.Email,
			Token:    token,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}
}
