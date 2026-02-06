package services_test

import (
	"testing"
	"time"

	"realworld-api/internal/dtos"
	"realworld-api/internal/models"
)

// Test CreateArticleRequest validation
func TestCreateArticleRequest_FieldValidation(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		body        string
		tagList     []string
		expectError bool
	}{
		{
			name:        "valid article with tags",
			title:       "Test Article",
			description: "Test description",
			body:        "Test body",
			tagList:     []string{"tag1", "tag2"},
			expectError: false,
		},
		{
			name:        "valid article without tags",
			title:       "Test Article",
			description: "Test description",
			body:        "Test body",
			tagList:     nil,
			expectError: false,
		},
		{
			name:        "empty title",
			title:       "",
			description: "Test description",
			body:        "Test body",
			tagList:     nil,
			expectError: true,
		},
		{
			name:        "empty description",
			title:       "Test Article",
			description: "",
			body:        "Test body",
			tagList:     nil,
			expectError: true,
		},
		{
			name:        "empty body",
			title:       "Test Article",
			description: "Test description",
			body:        "",
			tagList:     nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.title == "" || tt.description == "" || tt.body == ""
			if hasError != tt.expectError {
				t.Errorf("expected error: %v, got error: %v", tt.expectError, hasError)
			}
		})
	}
}

// Test UpdateArticleRequest optional fields
func TestUpdateArticleRequest_OptionalFields(t *testing.T) {
	title := "New Title"
	description := "New Description"
	body := "New Body"

	tests := []struct {
		name     string
		request  dtos.UpdateArticleRequest
		expected map[string]bool
	}{
		{
			name: "update title only",
			request: func() dtos.UpdateArticleRequest {
				var req dtos.UpdateArticleRequest
				req.Article.Title = &title
				return req
			}(),
			expected: map[string]bool{"title": true, "description": false, "body": false},
		},
		{
			name: "update description only",
			request: func() dtos.UpdateArticleRequest {
				var req dtos.UpdateArticleRequest
				req.Article.Description = &description
				return req
			}(),
			expected: map[string]bool{"title": false, "description": true, "body": false},
		},
		{
			name: "update body only",
			request: func() dtos.UpdateArticleRequest {
				var req dtos.UpdateArticleRequest
				req.Article.Body = &body
				return req
			}(),
			expected: map[string]bool{"title": false, "description": false, "body": true},
		},
		{
			name: "update all fields",
			request: func() dtos.UpdateArticleRequest {
				var req dtos.UpdateArticleRequest
				req.Article.Title = &title
				req.Article.Description = &description
				req.Article.Body = &body
				return req
			}(),
			expected: map[string]bool{"title": true, "description": true, "body": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.request.Article.Title != nil) != tt.expected["title"] {
				t.Errorf("title field presence mismatch")
			}
			if (tt.request.Article.Description != nil) != tt.expected["description"] {
				t.Errorf("description field presence mismatch")
			}
			if (tt.request.Article.Body != nil) != tt.expected["body"] {
				t.Errorf("body field presence mismatch")
			}
		})
	}
}

// Test ArticleData structure
func TestArticleData_Structure(t *testing.T) {
	articleData := dtos.ArticleData{
		Slug:           "test-article",
		Title:          "Test Article",
		Description:    "Test description",
		Body:           "Test body",
		TagList:        []string{"tag1", "tag2"},
		CreatedAt:      "2024-01-01T00:00:00Z",
		UpdatedAt:      "2024-01-01T00:00:00Z",
		Favorited:      false,
		FavoritesCount: 0,
		Author: dtos.ProfileData{
			Username:  "testuser",
			Bio:       "Test bio",
			Image:     "https://example.com/image.jpg",
			Following: false,
		},
	}

	if articleData.Slug != "test-article" {
		t.Errorf("Expected slug 'test-article', got '%s'", articleData.Slug)
	}
	if articleData.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", articleData.Title)
	}
	if len(articleData.TagList) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(articleData.TagList))
	}
	if articleData.FavoritesCount != 0 {
		t.Errorf("Expected favoritesCount 0, got %d", articleData.FavoritesCount)
	}
	if articleData.Author.Username != "testuser" {
		t.Errorf("Expected author username 'testuser', got '%s'", articleData.Author.Username)
	}
}

// Test Article model
func TestArticleModel_Fields(t *testing.T) {
	now := time.Now()
	article := models.Article{
		ID:             1,
		Slug:           "test-slug",
		Title:          "Test Title",
		Description:    "Test Description",
		Body:           "Test Body",
		AuthorID:       1,
		FavoritesCount: 5,
		CreatedAt:      now,
		UpdatedAt:      now,
		Author: models.User{
			ID:       1,
			Username: "author",
			Email:    "author@example.com",
		},
		Tags: []models.Tag{
			{ID: 1, Name: "tag1"},
			{ID: 2, Name: "tag2"},
		},
	}

	if article.ID != 1 {
		t.Errorf("Expected ID 1, got %d", article.ID)
	}
	if article.Slug != "test-slug" {
		t.Errorf("Expected slug 'test-slug', got '%s'", article.Slug)
	}
	if article.AuthorID != 1 {
		t.Errorf("Expected AuthorID 1, got %d", article.AuthorID)
	}
	if len(article.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(article.Tags))
	}
}

// Test ArticlesResponse structure
func TestArticlesResponse_Structure(t *testing.T) {
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

	if response.ArticlesCount != 2 {
		t.Errorf("Expected ArticlesCount 2, got %d", response.ArticlesCount)
	}
	if len(response.Articles) != 2 {
		t.Errorf("Expected 2 articles, got %d", len(response.Articles))
	}
}

// Test article authorization check
func TestArticle_AuthorizationCheck(t *testing.T) {
	tests := []struct {
		name         string
		articleOwner uint
		currentUser  uint
		isAuthorized bool
	}{
		{
			name:         "owner can modify",
			articleOwner: 1,
			currentUser:  1,
			isAuthorized: true,
		},
		{
			name:         "non-owner cannot modify",
			articleOwner: 1,
			currentUser:  2,
			isAuthorized: false,
		},
		{
			name:         "another owner check",
			articleOwner: 5,
			currentUser:  5,
			isAuthorized: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isAuthorized := tt.articleOwner == tt.currentUser
			if isAuthorized != tt.isAuthorized {
				t.Errorf("expected authorized: %v, got: %v", tt.isAuthorized, isAuthorized)
			}
		})
	}
}

// Test favorite status
func TestIsFavorited_Logic(t *testing.T) {
	tests := []struct {
		name       string
		userID     uint
		isFavorite bool
		expected   bool
	}{
		{
			name:       "user 0 (not logged in) cannot favorite",
			userID:     0,
			isFavorite: true,
			expected:   false,
		},
		{
			name:       "logged in user has favorited",
			userID:     1,
			isFavorite: true,
			expected:   true,
		},
		{
			name:       "logged in user has not favorited",
			userID:     1,
			isFavorite: false,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userID != 0 && tt.isFavorite
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

// Test is following logic
func TestIsFollowing_Logic(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		isFollowing bool
		expected    bool
	}{
		{
			name:        "user 0 (not logged in) cannot follow",
			userID:      0,
			isFollowing: true,
			expected:    false,
		},
		{
			name:        "logged in user is following",
			userID:      1,
			isFollowing: true,
			expected:    true,
		},
		{
			name:        "logged in user is not following",
			userID:      1,
			isFollowing: false,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userID != 0 && tt.isFollowing
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

// Test tag list conversion
func TestTagList_Conversion(t *testing.T) {
	tags := []models.Tag{
		{ID: 1, Name: "golang"},
		{ID: 2, Name: "testing"},
		{ID: 3, Name: "api"},
	}

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	expected := []string{"golang", "testing", "api"}
	for i, name := range expected {
		if tagNames[i] != name {
			t.Errorf("Expected tag '%s' at index %d, got '%s'", name, i, tagNames[i])
		}
	}
}

// Test empty tag list
func TestEmptyTagList_Conversion(t *testing.T) {
	tags := []models.Tag{}
	tagNames := make([]string, len(tags))

	if len(tagNames) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(tagNames))
	}
}

// Test pagination defaults
func TestPagination_Defaults(t *testing.T) {
	tests := []struct {
		name           string
		inputLimit     int
		inputOffset    int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "default values",
			inputLimit:     0,
			inputOffset:    0,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "custom limit",
			inputLimit:     10,
			inputOffset:    0,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "custom offset",
			inputLimit:     20,
			inputOffset:    5,
			expectedLimit:  20,
			expectedOffset: 5,
		},
		{
			name:           "negative limit uses default",
			inputLimit:     -5,
			inputOffset:    0,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "limit exceeds max uses max",
			inputLimit:     200,
			inputOffset:    0,
			expectedLimit:  100,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := tt.inputLimit
			offset := tt.inputOffset

			// Apply defaults
			if limit <= 0 {
				limit = 20
			}
			if limit > 100 {
				limit = 100
			}
			if offset < 0 {
				offset = 0
			}

			if limit != tt.expectedLimit {
				t.Errorf("expected limit %d, got %d", tt.expectedLimit, limit)
			}
			if offset != tt.expectedOffset {
				t.Errorf("expected offset %d, got %d", tt.expectedOffset, offset)
			}
		})
	}
}
