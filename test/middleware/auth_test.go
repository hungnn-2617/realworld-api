package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"realworld-api/internal/middleware"
	"realworld-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter() *gin.Engine {
	r := gin.New()
	return r
}

func TestAuthRequired_NoHeader(t *testing.T) {
	r := setupRouter()
	r.GET("/test", middleware.AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("AuthRequired() without header status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthRequired_InvalidFormat(t *testing.T) {
	r := setupRouter()
	r.GET("/test", middleware.AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	tests := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "missing bearer prefix",
			authHeader: "some-token",
		},
		{
			name:       "wrong prefix",
			authHeader: "Basic some-token",
		},
		{
			name:       "only bearer",
			authHeader: "Bearer",
		},
		{
			name:       "bearer with empty token",
			authHeader: "Bearer ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("AuthRequired() with %s status = %d, want %d", tt.name, w.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	r := setupRouter()
	r.GET("/test", middleware.AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("AuthRequired() with invalid token status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthOptional_NoHeader(t *testing.T) {
	r := setupRouter()
	r.GET("/test", middleware.AuthOptional(), func(c *gin.Context) {
		_, exists := middleware.GetUserID(c)
		c.JSON(http.StatusOK, gin.H{"authenticated": exists})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// Should pass through without authentication
	if w.Code != http.StatusOK {
		t.Errorf("AuthOptional() without header status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAuthOptional_InvalidToken(t *testing.T) {
	r := setupRouter()
	r.GET("/test", middleware.AuthOptional(), func(c *gin.Context) {
		_, exists := middleware.GetUserID(c)
		c.JSON(http.StatusOK, gin.H{"authenticated": exists})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	r.ServeHTTP(w, req)

	// Should pass through even with invalid token
	if w.Code != http.StatusOK {
		t.Errorf("AuthOptional() with invalid token status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetUserID(t *testing.T) {
	t.Run("user ID exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(middleware.ContextUserID, uint(123))

		userID, exists := middleware.GetUserID(c)

		if !exists {
			t.Error("GetUserID() exists = false, want true")
		}
		if userID != 123 {
			t.Errorf("GetUserID() = %d, want %d", userID, 123)
		}
	})

	t.Run("user ID does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID, exists := middleware.GetUserID(c)

		if exists {
			t.Error("GetUserID() exists = true, want false")
		}
		if userID != 0 {
			t.Errorf("GetUserID() = %d, want %d", userID, 0)
		}
	})

	t.Run("wrong type in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(middleware.ContextUserID, "not-a-uint")

		userID, exists := middleware.GetUserID(c)

		if exists {
			t.Error("GetUserID() exists = true for wrong type, want false")
		}
		if userID != 0 {
			t.Errorf("GetUserID() = %d, want %d", userID, 0)
		}
	})
}

func TestGetUsername(t *testing.T) {
	t.Run("username exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(middleware.ContextUsername, "testuser")

		username, exists := middleware.GetUsername(c)

		if !exists {
			t.Error("GetUsername() exists = false, want true")
		}
		if username != "testuser" {
			t.Errorf("GetUsername() = %q, want %q", username, "testuser")
		}
	})

	t.Run("username does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		username, exists := middleware.GetUsername(c)

		if exists {
			t.Error("GetUsername() exists = true, want false")
		}
		if username != "" {
			t.Errorf("GetUsername() = %q, want %q", username, "")
		}
	})
}

func TestGetEmail(t *testing.T) {
	t.Run("email exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(middleware.ContextEmail, "test@example.com")

		email, exists := middleware.GetEmail(c)

		if !exists {
			t.Error("GetEmail() exists = false, want true")
		}
		if email != "test@example.com" {
			t.Errorf("GetEmail() = %q, want %q", email, "test@example.com")
		}
	})

	t.Run("email does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		email, exists := middleware.GetEmail(c)

		if exists {
			t.Error("GetEmail() exists = true, want false")
		}
		if email != "" {
			t.Errorf("GetEmail() = %q, want %q", email, "")
		}
	})
}

func TestIsAuthenticated(t *testing.T) {
	t.Run("is authenticated", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(middleware.ContextUserID, uint(123))

		if !middleware.IsAuthenticated(c) {
			t.Error("IsAuthenticated() = false, want true")
		}
	})

	t.Run("is not authenticated", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		if middleware.IsAuthenticated(c) {
			t.Error("IsAuthenticated() = true, want false")
		}
	})
}

func TestMiddlewareErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "missing auth header",
			err:      middleware.ErrMissingAuthHeader,
			expected: "authorization header is required",
		},
		{
			name:     "invalid auth format",
			err:      middleware.ErrInvalidAuthFormat,
			expected: "authorization header format must be 'Bearer {token}'",
		},
		{
			name:     "empty token",
			err:      middleware.ErrEmptyToken,
			expected: "token is empty",
		},
		{
			name:     "invalid token",
			err:      middleware.ErrInvalidToken,
			expected: "invalid or malformed token",
		},
		{
			name:     "token expired",
			err:      middleware.ErrTokenExpired,
			expected: "token has expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Error message = %q, want %q", tt.err.Error(), tt.expected)
			}
		})
	}
}

// Test JWTClaims struct
func TestJWTClaimsStruct(t *testing.T) {
	claims := &utils.JWTClaims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	if claims.UserID != 1 {
		t.Errorf("JWTClaims.UserID = %d, want %d", claims.UserID, 1)
	}
	if claims.Username != "testuser" {
		t.Errorf("JWTClaims.Username = %q, want %q", claims.Username, "testuser")
	}
	if claims.Email != "test@example.com" {
		t.Errorf("JWTClaims.Email = %q, want %q", claims.Email, "test@example.com")
	}
}
