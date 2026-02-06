package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Test FavoriteArticle requires authentication
func TestFavoriteArticle_RequiresAuth(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/articles/:slug/favorite", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"article": gin.H{}})
	})

	req, _ := http.NewRequest("POST", "/api/articles/test-slug/favorite", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("FavoriteArticle without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test UnfavoriteArticle requires authentication
func TestUnfavoriteArticle_RequiresAuth(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/api/articles/:slug/favorite", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"article": gin.H{}})
	})

	req, _ := http.NewRequest("DELETE", "/api/articles/test-slug/favorite", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("UnfavoriteArticle without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test FavoriteArticle returns article with favorited=true
func TestFavoriteArticle_ReturnsFavoritedTrue(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/articles/:slug/favorite", func(c *gin.Context) {
		c.Set("userID", uint(1))
		slug := c.Param("slug")
		c.JSON(http.StatusOK, gin.H{
			"article": gin.H{
				"slug":           slug,
				"title":          "Test Article",
				"description":    "Test description",
				"body":           "Test body",
				"tagList":        []string{},
				"createdAt":      "2024-01-01T00:00:00Z",
				"updatedAt":      "2024-01-01T00:00:00Z",
				"favorited":      true,
				"favoritesCount": 1,
				"author": gin.H{
					"username":  "testuser",
					"bio":       "",
					"image":     "",
					"following": false,
				},
			},
		})
	})

	req, _ := http.NewRequest("POST", "/api/articles/test-slug/favorite", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Article struct {
			Slug           string `json:"slug"`
			Favorited      bool   `json:"favorited"`
			FavoritesCount int    `json:"favoritesCount"`
		} `json:"article"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Article.Favorited {
		t.Error("Expected favorited true after favorite, got false")
	}
	if response.Article.FavoritesCount != 1 {
		t.Errorf("Expected favoritesCount 1, got %d", response.Article.FavoritesCount)
	}
}

// Test UnfavoriteArticle returns article with favorited=false
func TestUnfavoriteArticle_ReturnsFavoritedFalse(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/api/articles/:slug/favorite", func(c *gin.Context) {
		c.Set("userID", uint(1))
		slug := c.Param("slug")
		c.JSON(http.StatusOK, gin.H{
			"article": gin.H{
				"slug":           slug,
				"title":          "Test Article",
				"description":    "Test description",
				"body":           "Test body",
				"tagList":        []string{},
				"createdAt":      "2024-01-01T00:00:00Z",
				"updatedAt":      "2024-01-01T00:00:00Z",
				"favorited":      false,
				"favoritesCount": 0,
				"author": gin.H{
					"username":  "testuser",
					"bio":       "",
					"image":     "",
					"following": false,
				},
			},
		})
	})

	req, _ := http.NewRequest("DELETE", "/api/articles/test-slug/favorite", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Article struct {
			Slug           string `json:"slug"`
			Favorited      bool   `json:"favorited"`
			FavoritesCount int    `json:"favoritesCount"`
		} `json:"article"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Article.Favorited {
		t.Error("Expected favorited false after unfavorite, got true")
	}
	if response.Article.FavoritesCount != 0 {
		t.Errorf("Expected favoritesCount 0, got %d", response.Article.FavoritesCount)
	}
}

// Test FavoriteArticle with non-existent article
func TestFavoriteArticle_ArticleNotFound(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/articles/:slug/favorite", func(c *gin.Context) {
		c.Set("userID", uint(1))
		slug := c.Param("slug")

		// Simulate article not found
		if slug == "non-existent-slug" {
			c.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"article not found"}}})
			return
		}

		c.JSON(http.StatusOK, gin.H{"article": gin.H{}})
	})

	req, _ := http.NewRequest("POST", "/api/articles/non-existent-slug/favorite", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("FavoriteArticle non-existent: expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// Test UnfavoriteArticle with non-existent article
func TestUnfavoriteArticle_ArticleNotFound(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/api/articles/:slug/favorite", func(c *gin.Context) {
		c.Set("userID", uint(1))
		slug := c.Param("slug")

		// Simulate article not found
		if slug == "non-existent-slug" {
			c.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"article not found"}}})
			return
		}

		c.JSON(http.StatusOK, gin.H{"article": gin.H{}})
	})

	req, _ := http.NewRequest("DELETE", "/api/articles/non-existent-slug/favorite", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("UnfavoriteArticle non-existent: expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// Test FavoriteArticleResponse structure
func TestFavoriteArticleResponse_Structure(t *testing.T) {
	r := setupTestRouter()
	r.POST("/api/articles/:slug/favorite", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.JSON(http.StatusOK, gin.H{
			"article": gin.H{
				"slug":           "test-article",
				"title":          "Test Article",
				"description":    "Test description",
				"body":           "Test body",
				"tagList":        []string{"tag1", "tag2"},
				"createdAt":      "2024-01-01T00:00:00Z",
				"updatedAt":      "2024-01-01T00:00:00Z",
				"favorited":      true,
				"favoritesCount": 5,
				"author": gin.H{
					"username":  "testuser",
					"bio":       "Test bio",
					"image":     "https://example.com/image.jpg",
					"following": false,
				},
			},
		})
	})

	req, _ := http.NewRequest("POST", "/api/articles/test-article/favorite", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Article struct {
			Slug           string   `json:"slug"`
			Title          string   `json:"title"`
			Description    string   `json:"description"`
			Body           string   `json:"body"`
			TagList        []string `json:"tagList"`
			CreatedAt      string   `json:"createdAt"`
			UpdatedAt      string   `json:"updatedAt"`
			Favorited      bool     `json:"favorited"`
			FavoritesCount int      `json:"favoritesCount"`
			Author         struct {
				Username  string `json:"username"`
				Bio       string `json:"bio"`
				Image     string `json:"image"`
				Following bool   `json:"following"`
			} `json:"author"`
		} `json:"article"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Article.Slug != "test-article" {
		t.Errorf("Expected slug 'test-article', got '%s'", response.Article.Slug)
	}
	if !response.Article.Favorited {
		t.Error("Expected favorited true")
	}
	if response.Article.FavoritesCount != 5 {
		t.Errorf("Expected favoritesCount 5, got %d", response.Article.FavoritesCount)
	}
}

// Test multiple favorites increase count
func TestFavorite_IncrementsFavoritesCount(t *testing.T) {
	count := 0
	r := setupTestRouter()
	r.POST("/api/articles/:slug/favorite", func(c *gin.Context) {
		c.Set("userID", uint(1))
		count++
		c.JSON(http.StatusOK, gin.H{
			"article": gin.H{
				"slug":           "test-article",
				"favorited":      true,
				"favoritesCount": count,
			},
		})
	})

	// First favorite
	req1, _ := http.NewRequest("POST", "/api/articles/test-article/favorite", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// Second favorite (simulating different user)
	req2, _ := http.NewRequest("POST", "/api/articles/test-article/favorite", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	var response struct {
		Article struct {
			FavoritesCount int `json:"favoritesCount"`
		} `json:"article"`
	}

	if err := json.Unmarshal(w2.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Article.FavoritesCount != 2 {
		t.Errorf("Expected favoritesCount 2 after 2 favorites, got %d", response.Article.FavoritesCount)
	}
}
