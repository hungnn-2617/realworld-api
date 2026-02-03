package handlers

import (
	"net/http"

	"realworld-api/internal/dtos"
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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.RegisterRequest true "User registration data"
// @Success 201 {object} dtos.UserResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 422 {object} dtos.ErrorResponse
// @Router /api/users [post]
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

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.LoginRequest true "User login credentials"
// @Success 200 {object} dtos.UserResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 422 {object} dtos.ErrorResponse
// @Router /api/users/login [post]
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

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get currently logged-in user
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dtos.UserResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Router /api/user [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		helpers.Error(c, http.StatusNotFound, "user not found")
		return
	}

	token, _ := h.userService.GenerateToken(user)
	helpers.SuccessResponse(c, http.StatusOK, h.toUserResponse(user, token))
}

// UpdateUser godoc
// @Summary Update current user
// @Description Update currently logged-in user's profile
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body dtos.UpdateUserRequest true "User update data"
// @Success 200 {object} dtos.UserResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 422 {object} dtos.ErrorResponse
// @Router /api/user [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dtos.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		helpers.Error(c, http.StatusNotFound, "user not found")
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
