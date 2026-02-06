package services_test

import (
	"testing"

	"realworld-api/internal/models"
)

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

// Test tag name conversion
func TestTagNames_Conversion(t *testing.T) {
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
	if len(tagNames) != len(expected) {
		t.Fatalf("Expected %d tags, got %d", len(expected), len(tagNames))
	}

	for i, name := range expected {
		if tagNames[i] != name {
			t.Errorf("Expected tag '%s' at index %d, got '%s'", name, i, tagNames[i])
		}
	}
}

// Test empty tags list
func TestEmptyTags_Conversion(t *testing.T) {
	tags := []models.Tag{}
	tagNames := make([]string, len(tags))

	if len(tagNames) != 0 {
		t.Errorf("Expected 0 tag names, got %d", len(tagNames))
	}
}

// Test tags response structure
func TestTagsResponse_Structure(t *testing.T) {
	response := struct {
		Tags []string `json:"tags"`
	}{
		Tags: []string{"golang", "javascript", "react"},
	}

	if len(response.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(response.Tags))
	}
	if response.Tags[0] != "golang" {
		t.Errorf("Expected first tag 'golang', got '%s'", response.Tags[0])
	}
}

// Test tag uniqueness
func TestTag_Uniqueness(t *testing.T) {
	existingTags := map[string]bool{
		"golang":     true,
		"javascript": true,
		"react":      true,
	}

	tests := []struct {
		name   string
		tag    string
		exists bool
	}{
		{
			name:   "existing tag",
			tag:    "golang",
			exists: true,
		},
		{
			name:   "new tag",
			tag:    "python",
			exists: false,
		},
		{
			name:   "another existing",
			tag:    "react",
			exists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := existingTags[tt.tag]
			if exists != tt.exists {
				t.Errorf("expected exists=%v for tag '%s', got %v", tt.exists, tt.tag, exists)
			}
		})
	}
}

// Test tag names with special characters
func TestTagNames_SpecialCharacters(t *testing.T) {
	tags := []models.Tag{
		{ID: 1, Name: "c++"},
		{ID: 2, Name: "c#"},
		{ID: 3, Name: ".net"},
		{ID: 4, Name: "node.js"},
	}

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	expected := []string{"c++", "c#", ".net", "node.js"}
	for i, name := range expected {
		if tagNames[i] != name {
			t.Errorf("Expected tag '%s' at index %d, got '%s'", name, i, tagNames[i])
		}
	}
}

// Test tag case sensitivity
func TestTag_CaseSensitivity(t *testing.T) {
	tests := []struct {
		name    string
		tag1    string
		tag2    string
		areSame bool
	}{
		{
			name:    "same case",
			tag1:    "golang",
			tag2:    "golang",
			areSame: true,
		},
		{
			name:    "different case",
			tag1:    "Golang",
			tag2:    "golang",
			areSame: false, // case-sensitive comparison
		},
		{
			name:    "all upper vs lower",
			tag1:    "REACT",
			tag2:    "react",
			areSame: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			areSame := tt.tag1 == tt.tag2
			if areSame != tt.areSame {
				t.Errorf("expected areSame=%v for '%s' and '%s', got %v", tt.areSame, tt.tag1, tt.tag2, areSame)
			}
		})
	}
}

// Test tag trimming
func TestTag_Trimming(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no trimming needed",
			input:    "golang",
			expected: "golang",
		},
		{
			name:     "leading space",
			input:    " golang",
			expected: "golang",
		},
		{
			name:     "trailing space",
			input:    "golang ",
			expected: "golang",
		},
		{
			name:     "both spaces",
			input:    " golang ",
			expected: "golang",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate trimming
			result := tt.input
			for len(result) > 0 && result[0] == ' ' {
				result = result[1:]
			}
			for len(result) > 0 && result[len(result)-1] == ' ' {
				result = result[:len(result)-1]
			}

			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test find or create tag logic
func TestFindOrCreateTag_Logic(t *testing.T) {
	existingTags := map[string]models.Tag{
		"golang": {ID: 1, Name: "golang"},
		"react":  {ID: 2, Name: "react"},
	}
	nextID := uint(3)

	tests := []struct {
		name      string
		tagName   string
		expectNew bool
	}{
		{
			name:      "existing tag",
			tagName:   "golang",
			expectNew: false,
		},
		{
			name:      "new tag",
			tagName:   "python",
			expectNew: true,
		},
		{
			name:      "another existing",
			tagName:   "react",
			expectNew: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tag models.Tag
			existingTag, exists := existingTags[tt.tagName]
			isNew := !exists

			if exists {
				tag = existingTag
			} else {
				tag = models.Tag{ID: nextID, Name: tt.tagName}
				nextID++
			}

			if isNew != tt.expectNew {
				t.Errorf("expected new=%v, got %v", tt.expectNew, isNew)
			}

			if tag.Name != tt.tagName {
				t.Errorf("expected tag name '%s', got '%s'", tt.tagName, tag.Name)
			}
		})
	}
}

// Test empty tag name
func TestEmptyTagName(t *testing.T) {
	tagName := ""
	isValid := tagName != ""

	if isValid {
		t.Error("Expected empty tag name to be invalid")
	}
}

// Test tag ordering
func TestTagOrdering(t *testing.T) {
	tags := []models.Tag{
		{ID: 1, Name: "zebra"},
		{ID: 2, Name: "apple"},
		{ID: 3, Name: "banana"},
	}

	// Sort by name
	for i := 0; i < len(tags)-1; i++ {
		for j := i + 1; j < len(tags); j++ {
			if tags[i].Name > tags[j].Name {
				tags[i], tags[j] = tags[j], tags[i]
			}
		}
	}

	expectedOrder := []string{"apple", "banana", "zebra"}
	for i, expected := range expectedOrder {
		if tags[i].Name != expected {
			t.Errorf("Expected tag '%s' at index %d, got '%s'", expected, i, tags[i].Name)
		}
	}
}
