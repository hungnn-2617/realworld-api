package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Errors ErrorBody `json:"errors"`
}

type ErrorBody struct {
	Body []string `json:"body"`
}

func NewErrorResponse(messages ...string) ErrorResponse {
	return ErrorResponse{
		Errors: ErrorBody{
			Body: messages,
		},
	}
}

func RespondWithError(c *gin.Context, statusCode int, messages ...string) {
	c.JSON(statusCode, NewErrorResponse(messages...))
}

// BadRequest responds with a 400 Bad Request error
func BadRequest(c *gin.Context, messages ...string) {
	RespondWithError(c, http.StatusBadRequest, messages...)
}

// Unauthorized responds with a 401 Unauthorized error
func Unauthorized(c *gin.Context, messages ...string) {
	if len(messages) == 0 {
		messages = []string{"unauthorized"}
	}
	RespondWithError(c, http.StatusUnauthorized, messages...)
}

// Forbidden responds with a 403 Forbidden error
func Forbidden(c *gin.Context, messages ...string) {
	if len(messages) == 0 {
		messages = []string{"forbidden"}
	}
	RespondWithError(c, http.StatusForbidden, messages...)
}

// NotFound responds with a 404 Not Found error
func NotFound(c *gin.Context, messages ...string) {
	if len(messages) == 0 {
		messages = []string{"not found"}
	}
	RespondWithError(c, http.StatusNotFound, messages...)
}

// Conflict responds with a 409 Conflict error (e.g., user already exists)
func Conflict(c *gin.Context, messages ...string) {
	RespondWithError(c, http.StatusConflict, messages...)
}

// UnprocessableEntity responds with a 422 Unprocessable Entity error (validation errors)
func UnprocessableEntity(c *gin.Context, messages ...string) {
	RespondWithError(c, http.StatusUnprocessableEntity, messages...)
}

// InternalServerError responds with a 500 Internal Server Error
func InternalServerError(c *gin.Context, messages ...string) {
	if len(messages) == 0 {
		messages = []string{"internal server error"}
	}
	RespondWithError(c, http.StatusInternalServerError, messages...)
}

// ValidationError responds with a 422 error for validation failures
func ValidationError(c *gin.Context, err error) {
	UnprocessableEntity(c, err.Error())
}

// Success sends a successful JSON response
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// OK responds with a 200 OK and data
func OK(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, data)
}

// Created responds with a 201 Created and data
func Created(c *gin.Context, data interface{}) {
	Success(c, http.StatusCreated, data)
}

// NoContent responds with a 204 No Content (for delete operations)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
