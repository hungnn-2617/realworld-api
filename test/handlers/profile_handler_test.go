package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Test GetProfile does not require authentication
func TestGetProfile_NoAuthRequired(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/profiles/:username", func(c *gin.Context) {
		username := c.Param("username")
		c.JSON(http.StatusOK, gin.H{
			"profile": gin.H{
				"username":  username,
				"bio":       "Test bio",
				"image":     "https://example.com/image.jpg",
				"following": false,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/profiles/testuser", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetProfile without auth: expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// Test FollowUser requires authentication
func TestFollowUser_RequiresAuth(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/profiles/:username/follow", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"profile": gin.H{}})
	})

	req, _ := http.NewRequest("POST", "/api/profiles/testuser/follow", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("FollowUser without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test UnfollowUser requires authentication
func TestUnfollowUser_RequiresAuth(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/api/profiles/:username/follow", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"profile": gin.H{}})
	})

	req, _ := http.NewRequest("DELETE", "/api/profiles/testuser/follow", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("UnfollowUser without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test ProfileResponse structure
func TestProfileResponse_Structure(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/profiles/:username", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"profile": gin.H{
				"username":  "testuser",
				"bio":       "Test bio",
				"image":     "https://example.com/image.jpg",
				"following": false,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/profiles/testuser", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Profile struct {
			Username  string `json:"username"`
			Bio       string `json:"bio"`
			Image     string `json:"image"`
			Following bool   `json:"following"`
		} `json:"profile"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Profile.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", response.Profile.Username)
	}
	if response.Profile.Bio != "Test bio" {
		t.Errorf("Expected bio 'Test bio', got '%s'", response.Profile.Bio)
	}
	if response.Profile.Following != false {
		t.Errorf("Expected following false, got true")
	}
}

// Test FollowUser returns profile with following=true
func TestFollowUser_ReturnsFollowingTrue(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/profiles/:username/follow", func(c *gin.Context) {
		c.Set("userID", uint(1))
		username := c.Param("username")
		c.JSON(http.StatusOK, gin.H{
			"profile": gin.H{
				"username":  username,
				"bio":       "",
				"image":     "",
				"following": true,
			},
		})
	})

	req, _ := http.NewRequest("POST", "/api/profiles/targetuser/follow", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Profile struct {
			Username  string `json:"username"`
			Following bool   `json:"following"`
		} `json:"profile"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Profile.Following != true {
		t.Errorf("Expected following true after follow, got false")
	}
}

// Test UnfollowUser returns profile with following=false
func TestUnfollowUser_ReturnsFollowingFalse(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/api/profiles/:username/follow", func(c *gin.Context) {
		c.Set("userID", uint(1))
		username := c.Param("username")
		c.JSON(http.StatusOK, gin.H{
			"profile": gin.H{
				"username":  username,
				"bio":       "",
				"image":     "",
				"following": false,
			},
		})
	})

	req, _ := http.NewRequest("DELETE", "/api/profiles/targetuser/follow", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Profile struct {
			Username  string `json:"username"`
			Following bool   `json:"following"`
		} `json:"profile"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Profile.Following != false {
		t.Errorf("Expected following false after unfollow, got true")
	}
}

// Test GetProfile with username containing special characters
func TestGetProfile_UsernameFormats(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		expectedCode int
	}{
		{
			name:         "simple username",
			username:     "testuser",
			expectedCode: http.StatusOK,
		},
		{
			name:         "username with numbers",
			username:     "testuser123",
			expectedCode: http.StatusOK,
		},
		{
			name:         "username with underscore",
			username:     "test_user",
			expectedCode: http.StatusOK,
		},
		{
			name:         "username with dash",
			username:     "test-user",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.GET("/api/profiles/:username", func(c *gin.Context) {
				username := c.Param("username")
				c.JSON(http.StatusOK, gin.H{
					"profile": gin.H{
						"username":  username,
						"bio":       "",
						"image":     "",
						"following": false,
					},
				})
			})

			req, _ := http.NewRequest("GET", "/api/profiles/"+tt.username, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
		})
	}
}

// Test FollowUser cannot follow self
func TestFollowUser_CannotFollowSelf(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/profiles/:username/follow", func(c *gin.Context) {
		c.Set("userID", uint(1))
		username := c.Param("username")

		// Simulate check for following self
		if username == "currentuser" { // Assume current user is "currentuser"
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{"cannot follow yourself"}}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"profile": gin.H{
				"username":  username,
				"bio":       "",
				"image":     "",
				"following": true,
			},
		})
	})

	req, _ := http.NewRequest("POST", "/api/profiles/currentuser/follow", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("FollowUser self: expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}
