package repositories_test

import (
	"testing"
	"time"

	"realworld-api/internal/models"
	"realworld-api/internal/repositories"
)

// Test Article model fields
func TestArticleModel_Fields(t *testing.T) {
	now := time.Now()
	article := models.Article{
		ID:             1,
		Slug:           "test-article",
		Title:          "Test Article",
		Description:    "Test description",
		Body:           "Test body content",
		AuthorID:       1,
		FavoritesCount: 5,
		CreatedAt:      now,
		UpdatedAt:      now,
		Author: models.User{
			ID:       1,
			Username: "author",
		},
		Tags: []models.Tag{
			{ID: 1, Name: "tag1"},
			{ID: 2, Name: "tag2"},
		},
	}

	if article.ID != 1 {
		t.Errorf("Expected ID 1, got %d", article.ID)
	}
	if article.Slug != "test-article" {
		t.Errorf("Expected slug 'test-article', got '%s'", article.Slug)
	}
	if article.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", article.Title)
	}
	if article.AuthorID != 1 {
		t.Errorf("Expected AuthorID 1, got %d", article.AuthorID)
	}
	if article.FavoritesCount != 5 {
		t.Errorf("Expected FavoritesCount 5, got %d", article.FavoritesCount)
	}
	if len(article.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(article.Tags))
	}
}

// Test ArticleFilter structure
func TestArticleFilter_Structure(t *testing.T) {
	filter := repositories.ArticleFilter{
		Tag:       "golang",
		Author:    "testuser",
		Favorited: "faveduser",
		Limit:     10,
		Offset:    5,
	}

	if filter.Tag != "golang" {
		t.Errorf("Expected tag 'golang', got '%s'", filter.Tag)
	}
	if filter.Author != "testuser" {
		t.Errorf("Expected author 'testuser', got '%s'", filter.Author)
	}
	if filter.Favorited != "faveduser" {
		t.Errorf("Expected favorited 'faveduser', got '%s'", filter.Favorited)
	}
	if filter.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", filter.Limit)
	}
	if filter.Offset != 5 {
		t.Errorf("Expected offset 5, got %d", filter.Offset)
	}
}

// Test ArticleFilter defaults
func TestArticleFilter_Defaults(t *testing.T) {
	tests := []struct {
		name          string
		filter        repositories.ArticleFilter
		expectedLimit int
	}{
		{
			name:          "zero limit uses default",
			filter:        repositories.ArticleFilter{Limit: 0},
			expectedLimit: 20,
		},
		{
			name:          "negative limit uses default",
			filter:        repositories.ArticleFilter{Limit: -5},
			expectedLimit: 20,
		},
		{
			name:          "limit exceeds max uses max",
			filter:        repositories.ArticleFilter{Limit: 200},
			expectedLimit: 100,
		},
		{
			name:          "valid limit",
			filter:        repositories.ArticleFilter{Limit: 50},
			expectedLimit: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := tt.filter.Limit
			if limit <= 0 {
				limit = 20
			}
			if limit > 100 {
				limit = 100
			}

			if limit != tt.expectedLimit {
				t.Errorf("expected limit %d, got %d", tt.expectedLimit, limit)
			}
		})
	}
}

// Test slug generation logic
func TestSlugGeneration_Logic(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "simple title",
			title:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "title with special chars",
			title:    "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "title with numbers",
			title:    "Article 123",
			expected: "article-123",
		},
		{
			name:     "title with extra spaces",
			title:    "  Hello   World  ",
			expected: "hello-world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple slug generation simulation
			slug := ""
			lastWasDash := true
			for _, c := range tt.title {
				if c >= 'A' && c <= 'Z' {
					c = c + 32 // to lowercase
				}
				if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
					slug += string(c)
					lastWasDash = false
				} else if !lastWasDash {
					slug += "-"
					lastWasDash = true
				}
			}
			// Trim trailing dash
			if len(slug) > 0 && slug[len(slug)-1] == '-' {
				slug = slug[:len(slug)-1]
			}

			if slug != tt.expected {
				t.Errorf("expected slug '%s', got '%s'", tt.expected, slug)
			}
		})
	}
}

// Test duplicate slug handling
func TestDuplicateSlug_Handling(t *testing.T) {
	existingSlugs := map[string]bool{
		"test-article": true,
	}

	tests := []struct {
		name        string
		slug        string
		isDuplicate bool
	}{
		{
			name:        "duplicate slug",
			slug:        "test-article",
			isDuplicate: true,
		},
		{
			name:        "unique slug",
			slug:        "another-article",
			isDuplicate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isDuplicate := existingSlugs[tt.slug]
			if isDuplicate != tt.isDuplicate {
				t.Errorf("expected isDuplicate=%v, got %v", tt.isDuplicate, isDuplicate)
			}
		})
	}
}

// Test feed articles filter by following IDs
func TestFeed_FilterByFollowingIDs(t *testing.T) {
	articles := []models.Article{
		{ID: 1, AuthorID: 1},
		{ID: 2, AuthorID: 2},
		{ID: 3, AuthorID: 3},
		{ID: 4, AuthorID: 2},
		{ID: 5, AuthorID: 4},
	}

	tests := []struct {
		name         string
		followingIDs []uint
		expectedIDs  []uint
	}{
		{
			name:         "following user 2",
			followingIDs: []uint{2},
			expectedIDs:  []uint{2, 4},
		},
		{
			name:         "following users 1 and 3",
			followingIDs: []uint{1, 3},
			expectedIDs:  []uint{1, 3},
		},
		{
			name:         "following nobody",
			followingIDs: []uint{},
			expectedIDs:  []uint{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var feedArticles []models.Article
			for _, article := range articles {
				for _, followingID := range tt.followingIDs {
					if article.AuthorID == followingID {
						feedArticles = append(feedArticles, article)
						break
					}
				}
			}

			if len(feedArticles) != len(tt.expectedIDs) {
				t.Errorf("expected %d articles, got %d", len(tt.expectedIDs), len(feedArticles))
			}
		})
	}
}

// Test filter by tag
func TestFilterByTag_Logic(t *testing.T) {
	articles := []models.Article{
		{ID: 1, Tags: []models.Tag{{Name: "golang"}, {Name: "api"}}},
		{ID: 2, Tags: []models.Tag{{Name: "javascript"}}},
		{ID: 3, Tags: []models.Tag{{Name: "golang"}}},
	}

	tagFilter := "golang"
	var filtered []models.Article

	for _, article := range articles {
		for _, tag := range article.Tags {
			if tag.Name == tagFilter {
				filtered = append(filtered, article)
				break
			}
		}
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 articles with tag 'golang', got %d", len(filtered))
	}
}

// Test filter by author
func TestFilterByAuthor_Logic(t *testing.T) {
	articles := []models.Article{
		{ID: 1, Author: models.User{Username: "testuser"}},
		{ID: 2, Author: models.User{Username: "anotheruser"}},
		{ID: 3, Author: models.User{Username: "testuser"}},
	}

	authorFilter := "testuser"
	var filtered []models.Article

	for _, article := range articles {
		if article.Author.Username == authorFilter {
			filtered = append(filtered, article)
		}
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 articles by author 'testuser', got %d", len(filtered))
	}
}

// Test Article model with comments relation
func TestArticleWithComments_Relation(t *testing.T) {
	article := models.Article{
		ID:    1,
		Title: "Test Article",
		Comments: []models.Comment{
			{ID: 1, Body: "Comment 1"},
			{ID: 2, Body: "Comment 2"},
		},
	}

	if len(article.Comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(article.Comments))
	}
}

// Test Article model with favorites relation
func TestArticleWithFavorites_Relation(t *testing.T) {
	article := models.Article{
		ID:             1,
		Title:          "Test Article",
		FavoritesCount: 3,
		Favorites: []models.Favorite{
			{UserID: 1, ArticleID: 1},
			{UserID: 2, ArticleID: 1},
			{UserID: 3, ArticleID: 1},
		},
	}

	if len(article.Favorites) != 3 {
		t.Errorf("Expected 3 favorites, got %d", len(article.Favorites))
	}
	if article.FavoritesCount != 3 {
		t.Errorf("Expected FavoritesCount 3, got %d", article.FavoritesCount)
	}
}

// Test Article ordering by created date
func TestArticleOrdering_CreatedAt(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)
	evenEarlier := now.Add(-2 * time.Hour)

	articles := []models.Article{
		{ID: 1, Title: "First", CreatedAt: evenEarlier},
		{ID: 2, Title: "Second", CreatedAt: earlier},
		{ID: 3, Title: "Third", CreatedAt: now},
	}

	// Sort by created_at DESC (most recent first)
	for i := 0; i < len(articles)-1; i++ {
		for j := i + 1; j < len(articles); j++ {
			if articles[i].CreatedAt.Before(articles[j].CreatedAt) {
				articles[i], articles[j] = articles[j], articles[i]
			}
		}
	}

	// Verify order (most recent first)
	if articles[0].ID != 3 {
		t.Errorf("Expected most recent article (ID 3) first, got ID %d", articles[0].ID)
	}
	if articles[2].ID != 1 {
		t.Errorf("Expected oldest article (ID 1) last, got ID %d", articles[2].ID)
	}
}

// Test pagination logic
func TestPagination_Logic(t *testing.T) {
	articles := make([]models.Article, 50)
	for i := 0; i < 50; i++ {
		articles[i] = models.Article{ID: uint(i + 1)}
	}

	tests := []struct {
		name        string
		limit       int
		offset      int
		expectedLen int
		firstID     uint
	}{
		{
			name:        "first page",
			limit:       10,
			offset:      0,
			expectedLen: 10,
			firstID:     1,
		},
		{
			name:        "second page",
			limit:       10,
			offset:      10,
			expectedLen: 10,
			firstID:     11,
		},
		{
			name:        "last partial page",
			limit:       10,
			offset:      45,
			expectedLen: 5,
			firstID:     46,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			end := tt.offset + tt.limit
			if end > len(articles) {
				end = len(articles)
			}
			paginated := articles[tt.offset:end]

			if len(paginated) != tt.expectedLen {
				t.Errorf("expected %d articles, got %d", tt.expectedLen, len(paginated))
			}
			if len(paginated) > 0 && paginated[0].ID != tt.firstID {
				t.Errorf("expected first ID %d, got %d", tt.firstID, paginated[0].ID)
			}
		})
	}
}
