package dtos_test

import (
	"encoding/json"
	"testing"

	"realworld-api/internal/dtos"
)

func TestCreateCommentRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{
			name:     "valid comment",
			json:     `{"comment":{"body":"This is a comment"}}`,
			expected: "This is a comment",
		},
		{
			name:     "empty body",
			json:     `{"comment":{"body":""}}`,
			expected: "",
		},
		{
			name:     "long comment",
			json:     `{"comment":{"body":"This is a very long comment with multiple sentences. It contains a lot of text to test the parsing."}}`,
			expected: "This is a very long comment with multiple sentences. It contains a lot of text to test the parsing.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result dtos.CreateCommentRequest
			err := json.Unmarshal([]byte(tt.json), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if result.Comment.Body != tt.expected {
				t.Errorf("Body = %q, want %q", result.Comment.Body, tt.expected)
			}
		})
	}
}

func TestCommentResponseJSON(t *testing.T) {
	response := dtos.CommentResponse{
		Comment: dtos.CommentData{
			ID:        1,
			Body:      "Test comment body",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
			Author: dtos.ProfileData{
				Username:  "commenter",
				Bio:       "Commenter bio",
				Image:     "https://example.com/avatar.jpg",
				Following: true,
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

	comment, ok := result["comment"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'comment' key in response")
	}

	if comment["id"] != float64(1) {
		t.Errorf("ID = %v, want %v", comment["id"], 1)
	}
	if comment["body"] != "Test comment body" {
		t.Errorf("Body = %v, want %v", comment["body"], "Test comment body")
	}

	author, ok := comment["author"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'author' key in comment")
	}
	if author["username"] != "commenter" {
		t.Errorf("Author.Username = %v, want %v", author["username"], "commenter")
	}
}

func TestCommentsResponseJSON(t *testing.T) {
	response := dtos.CommentsResponse{
		Comments: []dtos.CommentData{
			{
				ID:   1,
				Body: "First comment",
				Author: dtos.ProfileData{
					Username: "user1",
				},
			},
			{
				ID:   2,
				Body: "Second comment",
				Author: dtos.ProfileData{
					Username: "user2",
				},
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

	comments, ok := result["comments"].([]interface{})
	if !ok {
		t.Fatal("Expected 'comments' key as array in response")
	}

	if len(comments) != 2 {
		t.Errorf("Comments length = %d, want %d", len(comments), 2)
	}
}

func TestCommentDataStruct(t *testing.T) {
	commentData := dtos.CommentData{
		ID:        42,
		Body:      "Comment content",
		CreatedAt: "2024-01-01T12:00:00Z",
		UpdatedAt: "2024-01-02T12:00:00Z",
		Author: dtos.ProfileData{
			Username:  "testauthor",
			Bio:       "Author bio",
			Image:     "https://example.com/img.png",
			Following: false,
		},
	}

	if commentData.ID != 42 {
		t.Errorf("ID = %d, want %d", commentData.ID, 42)
	}
	if commentData.Body != "Comment content" {
		t.Errorf("Body = %q, want %q", commentData.Body, "Comment content")
	}
	if commentData.CreatedAt != "2024-01-01T12:00:00Z" {
		t.Errorf("CreatedAt = %q, want %q", commentData.CreatedAt, "2024-01-01T12:00:00Z")
	}
	if commentData.Author.Username != "testauthor" {
		t.Errorf("Author.Username = %q, want %q", commentData.Author.Username, "testauthor")
	}
}
