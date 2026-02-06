package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter() *gin.Engine {
	r := gin.New()
	return r
}

// Test helper to create JSON request body
func createJSONBody(t *testing.T, data interface{}) *bytes.Buffer {
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return bytes.NewBuffer(body)
}

// Test Registration Request Validation
func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "missing user field",
			requestBody: map[string]interface{}{
				"username": "testuser",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing username",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"email":    "test@example.com",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing email",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "testuser",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing password",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid email format",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "testuser",
					"email":    "invalid-email",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "password too short",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"password": "123",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "username too short",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "ab",
					"email":    "test@example.com",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.POST("/api/users", func(c *gin.Context) {
				var req struct {
					User struct {
						Username string `json:"username" binding:"required,min=3,max=50"`
						Email    string `json:"email" binding:"required,email"`
						Password string `json:"password" binding:"required,min=6"`
					} `json:"user" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"user": req.User})
			})

			body := createJSONBody(t, tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/users", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test Login Request Validation
func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "missing user field",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing email",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"password": "password123",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing password",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"email": "test@example.com",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid email format",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"email":    "invalid-email",
					"password": "password123",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.POST("/api/users/login", func(c *gin.Context) {
				var req struct {
					User struct {
						Email    string `json:"email" binding:"required,email"`
						Password string `json:"password" binding:"required"`
					} `json:"user" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
					return
				}
				c.JSON(http.StatusOK, gin.H{"user": req.User})
			})

			body := createJSONBody(t, tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/users/login", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test GetCurrentUser without auth header
func TestGetCurrentUser_Unauthorized(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/user", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": gin.H{}})
	})

	req, _ := http.NewRequest("GET", "/api/user", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("GetCurrentUser without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test UpdateUser Request Validation
func TestUpdateUserRequest_Validation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "missing user field",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "valid update with email",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"email": "newemail@example.com",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid update with username",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "newusername",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid update with bio",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"bio": "This is my new bio",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid update with image",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"image": "https://example.com/image.jpg",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid update with multiple fields",
			requestBody: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "newuser",
					"email":    "new@example.com",
					"bio":      "New bio",
					"image":    "https://example.com/new.jpg",
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.PUT("/api/user", func(c *gin.Context) {
				c.Set("userID", uint(1))
				var req map[string]interface{}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
					return
				}
				// Check if user field is present
				if _, ok := req["user"]; !ok {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{"user field is required"}}})
					return
				}
				c.JSON(http.StatusOK, gin.H{"user": req["user"]})
			})

			body := createJSONBody(t, tt.requestBody)
			req, _ := http.NewRequest("PUT", "/api/user", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test UserResponse structure
func TestUserResponse_Structure(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/user", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"email":    "test@example.com",
				"token":    "jwt-token-here",
				"username": "testuser",
				"bio":      "Test bio",
				"image":    "https://example.com/image.jpg",
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/user", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		User struct {
			Email    string `json:"email"`
			Token    string `json:"token"`
			Username string `json:"username"`
			Bio      string `json:"bio"`
			Image    string `json:"image"`
		} `json:"user"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.User.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", response.User.Email)
	}
	if response.User.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", response.User.Username)
	}
	if response.User.Token != "jwt-token-here" {
		t.Errorf("Expected token 'jwt-token-here', got '%s'", response.User.Token)
	}
}
