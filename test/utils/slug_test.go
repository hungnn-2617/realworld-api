package utils_test

import (
	"testing"

	"realworld-api/internal/utils"
)

func TestGenerateSlug(t *testing.T) {
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
			name:     "title with special characters",
			title:    "Hello, World! How are you?",
			expected: "hello-world-how-are-you",
		},
		{
			name:     "title with numbers",
			title:    "Top 10 Tips for 2024",
			expected: "top-10-tips-for-2024",
		},
		{
			name:     "title with multiple spaces",
			title:    "Hello    World",
			expected: "hello-world",
		},
		{
			name:     "title with leading and trailing spaces",
			title:    "  Hello World  ",
			expected: "hello-world",
		},
		{
			name:     "title with uppercase",
			title:    "HELLO WORLD",
			expected: "hello-world",
		},
		{
			name:     "empty title",
			title:    "",
			expected: "",
		},
		{
			name:     "title with only special characters",
			title:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "title with unicode characters",
			title:    "Café Résumé",
			expected: "cafe-resume",
		},
		{
			name:     "title with hyphens",
			title:    "hello-world-test",
			expected: "hello-world-test",
		},
		{
			name:     "title with multiple hyphens",
			title:    "hello---world",
			expected: "hello-world",
		},
		{
			name:     "title with mixed case and special chars",
			title:    "The Quick Brown Fox! Jumps Over",
			expected: "the-quick-brown-fox-jumps-over",
		},
		{
			name:     "title with apostrophe",
			title:    "It's a Test",
			expected: "its-a-test",
		},
		{
			name:     "title with ampersand",
			title:    "Rock & Roll",
			expected: "rock-roll",
		},
		{
			name:     "title starting with number",
			title:    "123 Test",
			expected: "123-test",
		},
		{
			name:     "Vietnamese title",
			title:    "Xin chào thế giới",
			expected: "xin-chao-the-gioi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GenerateSlug(tt.title)
			if result != tt.expected {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.title, result, tt.expected)
			}
		})
	}
}

func TestGenerateUniqueSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		suffix   string
		expected string
	}{
		{
			name:     "with suffix",
			title:    "Hello World",
			suffix:   "abc123",
			expected: "hello-world-abc123",
		},
		{
			name:     "without suffix",
			title:    "Hello World",
			suffix:   "",
			expected: "hello-world",
		},
		{
			name:     "empty title with suffix",
			title:    "",
			suffix:   "abc123",
			expected: "-abc123",
		},
		{
			name:     "special title with suffix",
			title:    "Hello, World!",
			suffix:   "123",
			expected: "hello-world-123",
		},
		{
			name:     "numeric suffix",
			title:    "Test Article",
			suffix:   "1234567890",
			expected: "test-article-1234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GenerateUniqueSlug(tt.title, tt.suffix)
			if result != tt.expected {
				t.Errorf("GenerateUniqueSlug(%q, %q) = %q, want %q", tt.title, tt.suffix, result, tt.expected)
			}
		})
	}
}

func TestGenerateSlugIdempotent(t *testing.T) {
	title := "Hello World Test"

	result1 := utils.GenerateSlug(title)
	result2 := utils.GenerateSlug(title)

	if result1 != result2 {
		t.Errorf("GenerateSlug() should be idempotent: got %q and %q", result1, result2)
	}
}
