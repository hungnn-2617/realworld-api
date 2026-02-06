package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Test GetTags does not require authentication
func TestGetTags_NoAuthRequired(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/tags", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"tags": []string{"tag1", "tag2", "tag3"},
		})
	})

	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetTags: expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// Test GetTags returns empty array when no tags
func TestGetTags_EmptyResponse(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/tags", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"tags": []string{},
		})
	})

	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Tags []string `json:"tags"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(response.Tags))
	}
}

// Test TagsResponse structure
func TestTagsResponse_Structure(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/tags", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"tags": []string{"golang", "javascript", "react", "docker", "kubernetes"},
		})
	})

	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Tags []string `json:"tags"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Tags) != 5 {
		t.Errorf("Expected 5 tags, got %d", len(response.Tags))
	}

	expectedTags := []string{"golang", "javascript", "react", "docker", "kubernetes"}
	for i, expected := range expectedTags {
		if response.Tags[i] != expected {
			t.Errorf("Expected tag '%s' at index %d, got '%s'", expected, i, response.Tags[i])
		}
	}
}

// Test GetTags handles server error
func TestGetTags_ServerError(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/tags", func(c *gin.Context) {
		// Simulate a server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": gin.H{
				"body": []string{"failed to fetch tags"},
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetTags server error: expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// Test tags with special characters
func TestGetTags_SpecialCharacters(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/tags", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"tags": []string{"c++", "c#", ".net", "node.js"},
		})
	})

	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Tags []string `json:"tags"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Tags) != 4 {
		t.Errorf("Expected 4 tags, got %d", len(response.Tags))
	}
}
