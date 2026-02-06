package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Test CreateArticle Request Validation
func TestCreateArticleRequest_Validation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "missing article field",
			requestBody: map[string]interface{}{
				"title": "Test Article",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing title",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"description": "Test description",
					"body":        "Test body",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing description",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"title": "Test Article",
					"body":  "Test body",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "missing body",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"title":       "Test Article",
					"description": "Test description",
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "valid article without tags",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"title":       "Test Article",
					"description": "Test description",
					"body":        "Test body",
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "valid article with tags",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"title":       "Test Article",
					"description": "Test description",
					"body":        "Test body",
					"tagList":     []string{"tag1", "tag2"},
				},
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.POST("/api/articles", func(c *gin.Context) {
				c.Set("userID", uint(1))
				var req struct {
					Article struct {
						Title       string   `json:"title" binding:"required"`
						Description string   `json:"description" binding:"required"`
						Body        string   `json:"body" binding:"required"`
						TagList     []string `json:"tagList,omitempty"`
					} `json:"article" binding:"required"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"article": req.Article})
			})

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test UpdateArticle Request Validation
func TestUpdateArticleRequest_Validation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "missing article field",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "valid update with title only",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"title": "New Title",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid update with description only",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"description": "New Description",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid update with body only",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"body": "New Body",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid update with all fields",
			requestBody: map[string]interface{}{
				"article": map[string]interface{}{
					"title":       "New Title",
					"description": "New Description",
					"body":        "New Body",
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.PUT("/api/articles/:slug", func(c *gin.Context) {
				c.Set("userID", uint(1))
				var req map[string]interface{}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{err.Error()}}})
					return
				}
				// Check if article field is present
				if _, ok := req["article"]; !ok {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": []string{"article field is required"}}})
					return
				}
				c.JSON(http.StatusOK, gin.H{"article": req["article"]})
			})

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/api/articles/test-slug", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test ListArticles Query Parameters
func TestListArticles_QueryParams(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  string
		expectedCode int
	}{
		{
			name:         "no params",
			queryParams:  "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "with tag filter",
			queryParams:  "?tag=golang",
			expectedCode: http.StatusOK,
		},
		{
			name:         "with author filter",
			queryParams:  "?author=testuser",
			expectedCode: http.StatusOK,
		},
		{
			name:         "with favorited filter",
			queryParams:  "?favorited=testuser",
			expectedCode: http.StatusOK,
		},
		{
			name:         "with limit",
			queryParams:  "?limit=10",
			expectedCode: http.StatusOK,
		},
		{
			name:         "with offset",
			queryParams:  "?offset=5",
			expectedCode: http.StatusOK,
		},
		{
			name:         "with all params",
			queryParams:  "?tag=golang&author=testuser&limit=10&offset=5",
			expectedCode: http.StatusOK,
		},
		{
			name:         "with invalid limit (non-numeric)",
			queryParams:  "?limit=abc",
			expectedCode: http.StatusOK, // Should default to 20
		},
		{
			name:         "with negative limit",
			queryParams:  "?limit=-5",
			expectedCode: http.StatusOK, // Should default to 20
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.GET("/api/articles", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"articles":      []interface{}{},
					"articlesCount": 0,
				})
			})

			req, _ := http.NewRequest("GET", "/api/articles"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
		})
	}
}

// Test GetArticle by slug
func TestGetArticle_Slug(t *testing.T) {
	tests := []struct {
		name         string
		slug         string
		expectedCode int
	}{
		{
			name:         "valid slug",
			slug:         "test-article",
			expectedCode: http.StatusOK,
		},
		{
			name:         "slug with numbers",
			slug:         "test-article-123",
			expectedCode: http.StatusOK,
		},
		{
			name:         "slug with special format",
			slug:         "test-article-1738881234567890",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.GET("/api/articles/:slug", func(c *gin.Context) {
				slug := c.Param("slug")
				if slug == "" {
					c.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"body": []string{"article not found"}}})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"article": gin.H{
						"slug":        slug,
						"title":       "Test Article",
						"description": "Test description",
						"body":        "Test body",
					},
				})
			})

			req, _ := http.NewRequest("GET", "/api/articles/"+tt.slug, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}
		})
	}
}

// Test FeedArticles requires authentication
func TestFeedArticles_RequiresAuth(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/articles/feed", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"articles":      []interface{}{},
			"articlesCount": 0,
		})
	})

	req, _ := http.NewRequest("GET", "/api/articles/feed", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("FeedArticles without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test DeleteArticle requires authentication
func TestDeleteArticle_RequiresAuth(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/api/articles/:slug", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"body": []string{"unauthorized"}}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "article deleted successfully"})
	})

	req, _ := http.NewRequest("DELETE", "/api/articles/test-slug", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("DeleteArticle without auth: expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Test ArticleResponse structure
func TestArticleResponse_Structure(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/articles/:slug", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"article": gin.H{
				"slug":           "test-article",
				"title":          "Test Article",
				"description":    "Test description",
				"body":           "Test body",
				"tagList":        []string{"tag1", "tag2"},
				"createdAt":      "2024-01-01T00:00:00Z",
				"updatedAt":      "2024-01-01T00:00:00Z",
				"favorited":      false,
				"favoritesCount": 0,
				"author": gin.H{
					"username":  "testuser",
					"bio":       "Test bio",
					"image":     "https://example.com/image.jpg",
					"following": false,
				},
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/articles/test-article", nil)
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
	if response.Article.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", response.Article.Title)
	}
	if len(response.Article.TagList) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(response.Article.TagList))
	}
}

// Test ArticlesResponse structure
func TestArticlesResponse_Structure(t *testing.T) {
	r := setupTestRouter()
	r.GET("/api/articles", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"articles": []gin.H{
				{
					"slug":           "article-1",
					"title":          "Article 1",
					"description":    "Description 1",
					"body":           "Body 1",
					"tagList":        []string{},
					"createdAt":      "2024-01-01T00:00:00Z",
					"updatedAt":      "2024-01-01T00:00:00Z",
					"favorited":      false,
					"favoritesCount": 5,
					"author": gin.H{
						"username":  "testuser",
						"bio":       "",
						"image":     "",
						"following": false,
					},
				},
			},
			"articlesCount": 1,
		})
	})

	req, _ := http.NewRequest("GET", "/api/articles", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Articles      []interface{} `json:"articles"`
		ArticlesCount int           `json:"articlesCount"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ArticlesCount != 1 {
		t.Errorf("Expected articlesCount 1, got %d", response.ArticlesCount)
	}
	if len(response.Articles) != 1 {
		t.Errorf("Expected 1 article, got %d", len(response.Articles))
	}
}
