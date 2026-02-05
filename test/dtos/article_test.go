package dtos_test

import (
	"encoding/json"
	"testing"

	"realworld-api/internal/dtos"
)

func TestCreateArticleRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected dtos.CreateArticleRequest
	}{
		{
			name: "with tags",
			json: `{"article":{"title":"Test Title","description":"Test Description","body":"Test Body","tagList":["tag1","tag2"]}}`,
			expected: func() dtos.CreateArticleRequest {
				r := dtos.CreateArticleRequest{}
				r.Article.Title = "Test Title"
				r.Article.Description = "Test Description"
				r.Article.Body = "Test Body"
				r.Article.TagList = []string{"tag1", "tag2"}
				return r
			}(),
		},
		{
			name: "without tags",
			json: `{"article":{"title":"Test Title","description":"Test Description","body":"Test Body"}}`,
			expected: func() dtos.CreateArticleRequest {
				r := dtos.CreateArticleRequest{}
				r.Article.Title = "Test Title"
				r.Article.Description = "Test Description"
				r.Article.Body = "Test Body"
				r.Article.TagList = nil
				return r
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result dtos.CreateArticleRequest
			err := json.Unmarshal([]byte(tt.json), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if result.Article.Title != tt.expected.Article.Title {
				t.Errorf("Title = %q, want %q", result.Article.Title, tt.expected.Article.Title)
			}
			if result.Article.Description != tt.expected.Article.Description {
				t.Errorf("Description = %q, want %q", result.Article.Description, tt.expected.Article.Description)
			}
			if result.Article.Body != tt.expected.Article.Body {
				t.Errorf("Body = %q, want %q", result.Article.Body, tt.expected.Article.Body)
			}
			if len(result.Article.TagList) != len(tt.expected.Article.TagList) {
				t.Errorf("TagList length = %d, want %d", len(result.Article.TagList), len(tt.expected.Article.TagList))
			}
		})
	}
}

func TestUpdateArticleRequestJSON(t *testing.T) {
	t.Run("full update", func(t *testing.T) {
		jsonStr := `{"article":{"title":"New Title","description":"New Description","body":"New Body"}}`
		var result dtos.UpdateArticleRequest
		err := json.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if result.Article.Title == nil || *result.Article.Title != "New Title" {
			t.Errorf("Title unexpected value")
		}
		if result.Article.Description == nil || *result.Article.Description != "New Description" {
			t.Errorf("Description unexpected value")
		}
		if result.Article.Body == nil || *result.Article.Body != "New Body" {
			t.Errorf("Body unexpected value")
		}
	})

	t.Run("partial update - only title", func(t *testing.T) {
		jsonStr := `{"article":{"title":"New Title"}}`
		var result dtos.UpdateArticleRequest
		err := json.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if result.Article.Title == nil || *result.Article.Title != "New Title" {
			t.Errorf("Title unexpected value")
		}
		if result.Article.Description != nil {
			t.Errorf("Description should be nil for partial update")
		}
		if result.Article.Body != nil {
			t.Errorf("Body should be nil for partial update")
		}
	})
}

func TestArticleResponseJSON(t *testing.T) {
	response := dtos.ArticleResponse{
		Article: dtos.ArticleData{
			Slug:           "test-slug",
			Title:          "Test Title",
			Description:    "Test Description",
			Body:           "Test Body",
			TagList:        []string{"tag1", "tag2"},
			CreatedAt:      "2024-01-01T00:00:00Z",
			UpdatedAt:      "2024-01-01T00:00:00Z",
			Favorited:      true,
			FavoritesCount: 5,
			Author: dtos.ProfileData{
				Username:  "author",
				Bio:       "Author bio",
				Image:     "https://example.com/img.jpg",
				Following: false,
			},
		},
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	article, ok := result["article"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'article' key in response")
	}

	if article["slug"] != "test-slug" {
		t.Errorf("Slug = %v, want %v", article["slug"], "test-slug")
	}
	if article["title"] != "Test Title" {
		t.Errorf("Title = %v, want %v", article["title"], "Test Title")
	}
	if article["favorited"] != true {
		t.Errorf("Favorited = %v, want %v", article["favorited"], true)
	}
	if article["favoritesCount"] != float64(5) {
		t.Errorf("FavoritesCount = %v, want %v", article["favoritesCount"], 5)
	}
}

func TestArticlesResponseJSON(t *testing.T) {
	response := dtos.ArticlesResponse{
		Articles: []dtos.ArticleData{
			{
				Slug:  "article-1",
				Title: "Article 1",
			},
			{
				Slug:  "article-2",
				Title: "Article 2",
			},
		},
		ArticlesCount: 2,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	articles, ok := result["articles"].([]interface{})
	if !ok {
		t.Fatal("Expected 'articles' key as array in response")
	}

	if len(articles) != 2 {
		t.Errorf("Articles length = %d, want %d", len(articles), 2)
	}

	if result["articlesCount"] != float64(2) {
		t.Errorf("ArticlesCount = %v, want %v", result["articlesCount"], 2)
	}
}

func TestArticleDataStruct(t *testing.T) {
	articleData := dtos.ArticleData{
		Slug:           "test-article",
		Title:          "Test Article",
		Description:    "Description",
		Body:           "Body content",
		TagList:        []string{"go", "testing"},
		CreatedAt:      "2024-01-01T00:00:00Z",
		UpdatedAt:      "2024-01-02T00:00:00Z",
		Favorited:      false,
		FavoritesCount: 10,
		Author: dtos.ProfileData{
			Username: "testuser",
		},
	}

	if articleData.Slug != "test-article" {
		t.Errorf("Slug = %q, want %q", articleData.Slug, "test-article")
	}
	if articleData.FavoritesCount != 10 {
		t.Errorf("FavoritesCount = %d, want %d", articleData.FavoritesCount, 10)
	}
	if len(articleData.TagList) != 2 {
		t.Errorf("TagList length = %d, want %d", len(articleData.TagList), 2)
	}
	if articleData.Author.Username != "testuser" {
		t.Errorf("Author.Username = %q, want %q", articleData.Author.Username, "testuser")
	}
}
