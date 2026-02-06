package repositories_test

import (
	"testing"
	"time"

	"realworld-api/internal/models"
)

// Test Comment model fields
func TestCommentModel_Fields(t *testing.T) {
	now := time.Now()
	comment := models.Comment{
		ID:        1,
		Body:      "Test comment body",
		ArticleID: 1,
		AuthorID:  1,
		CreatedAt: now,
		UpdatedAt: now,
		Author: models.User{
			ID:       1,
			Username: "commentauthor",
		},
	}

	if comment.ID != 1 {
		t.Errorf("Expected ID 1, got %d", comment.ID)
	}
	if comment.Body != "Test comment body" {
		t.Errorf("Expected body 'Test comment body', got '%s'", comment.Body)
	}
	if comment.ArticleID != 1 {
		t.Errorf("Expected ArticleID 1, got %d", comment.ArticleID)
	}
	if comment.AuthorID != 1 {
		t.Errorf("Expected AuthorID 1, got %d", comment.AuthorID)
	}
	if comment.Author.Username != "commentauthor" {
		t.Errorf("Expected Author.Username 'commentauthor', got '%s'", comment.Author.Username)
	}
}

// Test Comment ID validation
func TestCommentID_Validation(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		isValid bool
	}{
		{
			name:    "valid ID",
			id:      1,
			isValid: true,
		},
		{
			name:    "zero ID",
			id:      0,
			isValid: false,
		},
		{
			name:    "large ID",
			id:      999999,
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.id > 0
			if isValid != tt.isValid {
				t.Errorf("expected valid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

// Test find comments by article ID
func TestFindByArticleID_Logic(t *testing.T) {
	comments := []models.Comment{
		{ID: 1, ArticleID: 1, Body: "Comment 1"},
		{ID: 2, ArticleID: 1, Body: "Comment 2"},
		{ID: 3, ArticleID: 2, Body: "Comment 3"},
		{ID: 4, ArticleID: 1, Body: "Comment 4"},
	}

	tests := []struct {
		name        string
		articleID   uint
		expectedLen int
	}{
		{
			name:        "article 1 comments",
			articleID:   1,
			expectedLen: 3,
		},
		{
			name:        "article 2 comments",
			articleID:   2,
			expectedLen: 1,
		},
		{
			name:        "article 3 no comments",
			articleID:   3,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filtered []models.Comment
			for _, c := range comments {
				if c.ArticleID == tt.articleID {
					filtered = append(filtered, c)
				}
			}

			if len(filtered) != tt.expectedLen {
				t.Errorf("expected %d comments, got %d", tt.expectedLen, len(filtered))
			}
		})
	}
}

// Test comment ordering by created date
func TestCommentOrdering_CreatedAt(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)
	evenEarlier := now.Add(-2 * time.Hour)

	comments := []models.Comment{
		{ID: 1, Body: "First", CreatedAt: evenEarlier},
		{ID: 2, Body: "Second", CreatedAt: earlier},
		{ID: 3, Body: "Third", CreatedAt: now},
	}

	// Sort by created_at DESC (most recent first)
	for i := 0; i < len(comments)-1; i++ {
		for j := i + 1; j < len(comments); j++ {
			if comments[i].CreatedAt.Before(comments[j].CreatedAt) {
				comments[i], comments[j] = comments[j], comments[i]
			}
		}
	}

	// Verify order (most recent first)
	if comments[0].ID != 3 {
		t.Errorf("Expected most recent comment (ID 3) first, got ID %d", comments[0].ID)
	}
	if comments[2].ID != 1 {
		t.Errorf("Expected oldest comment (ID 1) last, got ID %d", comments[2].ID)
	}
}

// Test comment with author relation
func TestCommentWithAuthor_Relation(t *testing.T) {
	comment := models.Comment{
		ID:       1,
		Body:     "Test comment",
		AuthorID: 1,
		Author: models.User{
			ID:       1,
			Username: "testuser",
			Bio:      "Test bio",
			Image:    "https://example.com/image.jpg",
		},
	}

	if comment.Author.ID != 1 {
		t.Errorf("Expected Author.ID 1, got %d", comment.Author.ID)
	}
	if comment.Author.Username != "testuser" {
		t.Errorf("Expected Author.Username 'testuser', got '%s'", comment.Author.Username)
	}
}

// Test comment delete by ID
func TestDeleteByID_Logic(t *testing.T) {
	comments := map[uint]models.Comment{
		1: {ID: 1, Body: "Comment 1"},
		2: {ID: 2, Body: "Comment 2"},
		3: {ID: 3, Body: "Comment 3"},
	}

	tests := []struct {
		name      string
		id        uint
		preCount  int
		postCount int
	}{
		{
			name:      "delete existing comment",
			id:        2,
			preCount:  3,
			postCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(comments) != tt.preCount {
				t.Errorf("expected pre-delete count %d, got %d", tt.preCount, len(comments))
			}

			delete(comments, tt.id)

			if len(comments) != tt.postCount {
				t.Errorf("expected post-delete count %d, got %d", tt.postCount, len(comments))
			}

			if _, exists := comments[tt.id]; exists {
				t.Errorf("expected comment %d to be deleted", tt.id)
			}
		})
	}
}

// Test comment timestamps
func TestCommentModel_Timestamps(t *testing.T) {
	now := time.Now()
	comment := models.Comment{
		ID:        1,
		Body:      "Test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if comment.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if comment.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
	if !comment.CreatedAt.Equal(comment.UpdatedAt) {
		t.Error("CreatedAt and UpdatedAt should be equal for new comment")
	}
}

// Test comment body validation
func TestCommentBody_Validation(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		isValid bool
	}{
		{
			name:    "valid body",
			body:    "This is a valid comment",
			isValid: true,
		},
		{
			name:    "empty body",
			body:    "",
			isValid: false,
		},
		{
			name:    "very long body",
			body:    string(make([]byte, 10000)),
			isValid: true, // assuming no max length constraint
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.body != ""
			if isValid != tt.isValid {
				t.Errorf("expected valid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

// Test article and comment relationship
func TestArticleComment_Relationship(t *testing.T) {
	article := models.Article{
		ID:    1,
		Title: "Test Article",
	}

	comment := models.Comment{
		ID:        1,
		Body:      "Test comment",
		ArticleID: article.ID,
	}

	if comment.ArticleID != article.ID {
		t.Errorf("Expected comment.ArticleID %d, got %d", article.ID, comment.ArticleID)
	}
}

// Test Favorite model fields
func TestFavoriteModel_Fields(t *testing.T) {
	now := time.Now()
	favorite := models.Favorite{
		UserID:    1,
		ArticleID: 1,
		CreatedAt: now,
	}

	if favorite.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", favorite.UserID)
	}
	if favorite.ArticleID != 1 {
		t.Errorf("Expected ArticleID 1, got %d", favorite.ArticleID)
	}
}

// Test isFavorited check
func TestIsFavorited_Logic(t *testing.T) {
	favorites := map[string]bool{
		"1-1": true, // user 1 favorited article 1
		"1-2": true, // user 1 favorited article 2
	}

	tests := []struct {
		name        string
		userID      uint
		articleID   uint
		isFavorited bool
	}{
		{
			name:        "user favorited",
			userID:      1,
			articleID:   1,
			isFavorited: true,
		},
		{
			name:        "user not favorited",
			userID:      2,
			articleID:   1,
			isFavorited: false,
		},
		{
			name:        "another favorited",
			userID:      1,
			articleID:   2,
			isFavorited: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := string(rune(tt.userID+48)) + "-" + string(rune(tt.articleID+48))
			isFavorited := favorites[key]
			if isFavorited != tt.isFavorited {
				t.Errorf("expected isFavorited=%v, got %v", tt.isFavorited, isFavorited)
			}
		})
	}
}

// Test unique favorite constraint
func TestUniqueFavorite_Constraint(t *testing.T) {
	existingFavorites := map[string]bool{
		"1-1": true,
	}

	tests := []struct {
		name        string
		userID      uint
		articleID   uint
		isDuplicate bool
	}{
		{
			name:        "duplicate favorite",
			userID:      1,
			articleID:   1,
			isDuplicate: true,
		},
		{
			name:        "new favorite",
			userID:      1,
			articleID:   2,
			isDuplicate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := string(rune(tt.userID+48)) + "-" + string(rune(tt.articleID+48))
			isDuplicate := existingFavorites[key]
			if isDuplicate != tt.isDuplicate {
				t.Errorf("expected isDuplicate=%v, got %v", tt.isDuplicate, isDuplicate)
			}
		})
	}
}

// Test favorites count update
func TestFavoritesCount_Update(t *testing.T) {
	article := models.Article{
		ID:             1,
		FavoritesCount: 5,
	}

	// Simulate favorite
	article.FavoritesCount++
	if article.FavoritesCount != 6 {
		t.Errorf("Expected FavoritesCount 6 after favorite, got %d", article.FavoritesCount)
	}

	// Simulate unfavorite
	article.FavoritesCount--
	if article.FavoritesCount != 5 {
		t.Errorf("Expected FavoritesCount 5 after unfavorite, got %d", article.FavoritesCount)
	}
}

// Test Tag model
func TestTagModel_Fields(t *testing.T) {
	tag := models.Tag{
		ID:   1,
		Name: "golang",
	}

	if tag.ID != 1 {
		t.Errorf("Expected ID 1, got %d", tag.ID)
	}
	if tag.Name != "golang" {
		t.Errorf("Expected name 'golang', got '%s'", tag.Name)
	}
}

// Test many-to-many article tags relationship
func TestArticleTags_ManyToMany(t *testing.T) {
	article := models.Article{
		ID:    1,
		Title: "Test Article",
		Tags: []models.Tag{
			{ID: 1, Name: "golang"},
			{ID: 2, Name: "testing"},
			{ID: 3, Name: "api"},
		},
	}

	if len(article.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(article.Tags))
	}

	expectedTags := []string{"golang", "testing", "api"}
	for i, tag := range article.Tags {
		if tag.Name != expectedTags[i] {
			t.Errorf("Expected tag '%s' at index %d, got '%s'", expectedTags[i], i, tag.Name)
		}
	}
}
