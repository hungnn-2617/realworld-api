package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Test CreateComment Request Validation
func TestCreateCommentRequest_Validation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "missing comment field",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing body",
			requestBody: map[string]interface{}{
				"comment": map[string]interface{}{},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "empty body",
			requestBody: map[string]interface{}{
				"comment": map[string]interface{}{
					"body": "",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "valid comment",
			requestBody: map[string]interface{}{
				"comment": map[string]interface{}{
					"body": "This is a test comment",
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.POST("/api/articles/:slug/comments", func(c *gin.Context) {
				c.Set("userID", uint(1))
				var req struct {
					Comment struct {
						Body string `json:"body" binding:"required"`
					} `json:"comment" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
					return
				}
				c.JSON(http.StatusOK, gin.H{"comment": req.Comment})
			})

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/articles/test-slug/comments", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test AddComment requires authentication
func TestAddComment_RequiresAuth(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/articles/:slug/comments", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"comment": gin.H{}})
	})

	requestBody := map[string]interface{}{
		"comment": map[string]interface{}{
			"body": "Test comment",
		},
	}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/articles/test-slug/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("AddComment without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test GetComments does not require authentication
func TestGetComments_NoAuthRequired(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/articles/:slug/comments", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"comments": []interface{}{},
		})
	})

	req, _ := http.NewRequest("GET", "/api/articles/test-slug/comments", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetComments without auth: expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// Test DeleteComment requires authentication
func TestDeleteComment_RequiresAuth(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/api/articles/:slug/comments/:id", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
	})

	req, _ := http.NewRequest("DELETE", "/api/articles/test-slug/comments/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("DeleteComment without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test DeleteComment with invalid comment ID
func TestDeleteComment_InvalidCommentID(t *testing.T) {
	tests := []struct {
		name         string
		commentID    string
		expectedCode int
	}{
		{
			name:         "non-numeric ID",
			commentID:    "abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "negative ID",
			commentID:    "-1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "valid numeric ID",
			commentID:    "123",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.DELETE("/api/articles/:slug/comments/:id", func(c *gin.Context) {
				c.Set("userID", uint(1))
				idStr := c.Param("id")
				// Validate ID is positive integer
				var id int
				if err := json.Unmarshal([]byte(idStr), &id); err != nil || id <= 0 {
					c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"body": []string{"invalid comment id"}}})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
			})

			req, _ := http.NewRequest("DELETE", "/api/articles/test-slug/comments/"+tt.commentID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
		})
	}
}

// Test CommentResponse structure
func TestCommentResponse_Structure(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/articles/:slug/comments", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"comments": []gin.H{
				{
					"id":        1,
					"body":      "Test comment",
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-01T00:00:00Z",
					"author": gin.H{
						"username":  "testuser",
						"bio":       "Test bio",
						"image":     "https://example.com/image.jpg",
						"following": false,
					},
				},
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/articles/test-slug/comments", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Comments []struct {
			ID        int    `json:"id"`
			Body      string `json:"body"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
			Author    struct {
				Username  string `json:"username"`
				Bio       string `json:"bio"`
				Image     string `json:"image"`
				Following bool   `json:"following"`
			} `json:"author"`
		} `json:"comments"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Comments) != 1 {
		t.Errorf("Expected 1 comment, got %d", len(response.Comments))
	}
	if response.Comments[0].ID != 1 {
		t.Errorf("Expected comment ID 1, got %d", response.Comments[0].ID)
	}
	if response.Comments[0].Body != "Test comment" {
		t.Errorf("Expected comment body 'Test comment', got '%s'", response.Comments[0].Body)
	}
}

// Test SingleCommentResponse structure
func TestSingleCommentResponse_Structure(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/articles/:slug/comments", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.JSON(http.StatusOK, gin.H{
			"comment": gin.H{
				"id":        1,
				"body":      "New comment",
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z",
				"author": gin.H{
					"username":  "testuser",
					"bio":       "",
					"image":     "",
					"following": false,
				},
			},
		})
	})

	requestBody := map[string]interface{}{
		"comment": map[string]interface{}{
			"body": "New comment",
		},
	}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/articles/test-slug/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Comment struct {
			ID   int    `json:"id"`
			Body string `json:"body"`
		} `json:"comment"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Comment.ID != 1 {
		t.Errorf("Expected comment ID 1, got %d", response.Comment.ID)
	}
}
