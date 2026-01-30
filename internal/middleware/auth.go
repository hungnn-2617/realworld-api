package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/realworld-api/internal/utils"
)

const (
	AuthorizationHeader = "Authorization"
	AuthorizationType   = "Token"
	ContextUserID       = "userID"
	ContextUsername     = "username"
	ContextEmail        = "email"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := extractAndValidateToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": gin.H{
					"body": []string{err.Error()},
				},
			})
			return
		}

		setUserContext(c, claims)
		c.Next()
	}
}

func AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := extractAndValidateToken(c)
		if err == nil && claims != nil {
			setUserContext(c, claims)
		}
		c.Next()
	}
}

func extractAndValidateToken(c *gin.Context) (*utils.JWTClaims, error) {
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return nil, ErrMissingAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], AuthorizationType) {
		return nil, ErrInvalidAuthFormat
	}

	tokenString := parts[1]
	if tokenString == "" {
		return nil, ErrEmptyToken
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		if err == utils.ErrExpiredToken {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func setUserContext(c *gin.Context, claims *utils.JWTClaims) {
	c.Set(ContextUserID, claims.UserID)
	c.Set(ContextUsername, claims.Username)
	c.Set(ContextEmail, claims.Email)
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(ContextUsername)
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(ContextEmail)
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(ContextUserID)
	return exists
}
