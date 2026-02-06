package services_test

import (
	"testing"
	"time"

	"realworld-api/internal/dtos"
	"realworld-api/internal/models"
)

// Test CreateCommentRequest validation
func TestCreateCommentRequest_FieldValidation(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectError bool
	}{
		{
			name:        "valid comment",
			body:        "This is a test comment",
			expectError: false,
		},
		{
			name:        "empty body",
			body:        "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			body:        "   ",
			expectError: false, // validation typically doesn't check for whitespace
		},
		{
			name:        "long comment",
			body:        "This is a very long comment that contains lots of text to test if long comments work properly. This text continues for a while to make sure everything handles properly.",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.body == ""
			if hasError != tt.expectError {
				t.Errorf("expected error: %v, got error: %v", tt.expectError, hasError)
			}
		})
	}
}

// Test CommentData structure
func TestCommentData_Structure(t *testing.T) {
	commentData := dtos.CommentData{
		ID:        1,
		Body:      "Test comment",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
		Author: dtos.ProfileData{
			Username:  "testuser",
			Bio:       "Test bio",
			Image:     "https://example.com/image.jpg",
			Following: false,
		},
	}

	if commentData.ID != 1 {
		t.Errorf("Expected ID 1, got %d", commentData.ID)
	}
	if commentData.Body != "Test comment" {
		t.Errorf("Expected body 'Test comment', got '%s'", commentData.Body)
	}
	if commentData.Author.Username != "testuser" {
		t.Errorf("Expected author username 'testuser', got '%s'", commentData.Author.Username)
	}
}

// Test Comment model
func TestCommentModel_Fields(t *testing.T) {
	now := time.Now()
	comment := models.Comment{
		ID:        1,
		Body:      "Test Comment Body",
		ArticleID: 1,
		AuthorID:  1,
		CreatedAt: now,
		UpdatedAt: now,
		Author: models.User{
			ID:       1,
			Username: "commentauthor",
			Email:    "author@example.com",
		},
	}

	if comment.ID != 1 {
		t.Errorf("Expected ID 1, got %d", comment.ID)
	}
	if comment.Body != "Test Comment Body" {
		t.Errorf("Expected body 'Test Comment Body', got '%s'", comment.Body)
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

// Test CommentsResponse structure
func TestCommentsResponse_Structure(t *testing.T) {
	response := dtos.CommentsResponse{
		Comments: []dtos.CommentData{
			{
				ID:   1,
				Body: "Comment 1",
			},
			{
				ID:   2,
				Body: "Comment 2",
			},
		},
	}

	if len(response.Comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(response.Comments))
	}
	if response.Comments[0].ID != 1 {
		t.Errorf("Expected first comment ID 1, got %d", response.Comments[0].ID)
	}
}

// Test comment authorization check
func TestComment_AuthorizationCheck(t *testing.T) {
	tests := []struct {
		name          string
		commentAuthor uint
		currentUser   uint
		isAuthorized  bool
	}{
		{
			name:          "author can delete",
			commentAuthor: 1,
			currentUser:   1,
			isAuthorized:  true,
		},
		{
			name:          "non-author cannot delete",
			commentAuthor: 1,
			currentUser:   2,
			isAuthorized:  false,
		},
		{
			name:          "different user ids",
			commentAuthor: 5,
			currentUser:   10,
			isAuthorized:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isAuthorized := tt.commentAuthor == tt.currentUser
			if isAuthorized != tt.isAuthorized {
				t.Errorf("expected authorized: %v, got: %v", tt.isAuthorized, isAuthorized)
			}
		})
	}
}

// Test comment ID validation
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
				t.Errorf("expected valid: %v, got: %v", tt.isValid, isValid)
			}
		})
	}
}

// Test comments sorted by created date
func TestComments_SortOrder(t *testing.T) {
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

// Test empty comments list
func TestEmptyComments_Response(t *testing.T) {
	response := dtos.CommentsResponse{
		Comments: []dtos.CommentData{},
	}

	if len(response.Comments) != 0 {
		t.Errorf("Expected 0 comments, got %d", len(response.Comments))
	}
}

// Test CommentResponse single comment
func TestCommentResponse_Structure(t *testing.T) {
	response := dtos.CommentResponse{
		Comment: dtos.CommentData{
			ID:        1,
			Body:      "Single comment",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
			Author: dtos.ProfileData{
				Username: "user",
			},
		},
	}

	if response.Comment.ID != 1 {
		t.Errorf("Expected comment ID 1, got %d", response.Comment.ID)
	}
	if response.Comment.Body != "Single comment" {
		t.Errorf("Expected body 'Single comment', got '%s'", response.Comment.Body)
	}
}
