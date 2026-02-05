package errors_test

import (
	"errors"
	"testing"

	appErrors "realworld-api/internal/errors"
)

func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		// Authentication & Authorization errors
		{name: "ErrUnauthorized", constant: appErrors.ErrUnauthorized, expected: "unauthorized"},
		{name: "ErrForbidden", constant: appErrors.ErrForbidden, expected: "forbidden"},
		{name: "ErrInvalidCredentials", constant: appErrors.ErrInvalidCredentials, expected: "invalid email or password"},
		{name: "ErrInvalidToken", constant: appErrors.ErrInvalidToken, expected: "invalid token"},
		{name: "ErrTokenExpired", constant: appErrors.ErrTokenExpired, expected: "token expired"},
		{name: "ErrTokenRequired", constant: appErrors.ErrTokenRequired, expected: "authorization token required"},

		// User errors
		{name: "ErrUserNotFound", constant: appErrors.ErrUserNotFound, expected: "user not found"},
		{name: "ErrUserAlreadyExists", constant: appErrors.ErrUserAlreadyExists, expected: "user already exists"},
		{name: "ErrEmailAlreadyTaken", constant: appErrors.ErrEmailAlreadyTaken, expected: "email already taken"},
		{name: "ErrUsernameAlreadyTaken", constant: appErrors.ErrUsernameAlreadyTaken, expected: "username already taken"},

		// Profile errors
		{name: "ErrCannotFollowYourself", constant: appErrors.ErrCannotFollowYourself, expected: "cannot follow yourself"},

		// Article errors
		{name: "ErrArticleNotFound", constant: appErrors.ErrArticleNotFound, expected: "article not found"},
		{name: "ErrNotAuthorizedUpdate", constant: appErrors.ErrNotAuthorizedUpdate, expected: "not authorized to update this article"},
		{name: "ErrNotAuthorizedDelete", constant: appErrors.ErrNotAuthorizedDelete, expected: "not authorized to delete this article"},

		// Comment errors
		{name: "ErrCommentNotFound", constant: appErrors.ErrCommentNotFound, expected: "comment not found"},
		{name: "ErrNotAuthorizedDeleteComment", constant: appErrors.ErrNotAuthorizedDeleteComment, expected: "not authorized to delete this comment"},
		{name: "ErrInvalidCommentID", constant: appErrors.ErrInvalidCommentID, expected: "invalid comment id"},

		// Favorite errors
		{name: "ErrAlreadyFavorited", constant: appErrors.ErrAlreadyFavorited, expected: "article already favorited"},
		{name: "ErrNotFavorited", constant: appErrors.ErrNotFavorited, expected: "article not favorited"},

		// Validation errors
		{name: "ErrValidationFailed", constant: appErrors.ErrValidationFailed, expected: "validation failed"},
		{name: "ErrInvalidInput", constant: appErrors.ErrInvalidInput, expected: "invalid input"},

		// Server errors
		{name: "ErrInternalServer", constant: appErrors.ErrInternalServer, expected: "internal server error"},
		{name: "ErrFailedToFetch", constant: appErrors.ErrFailedToFetch, expected: "failed to fetch data"},
		{name: "ErrFailedToCreate", constant: appErrors.ErrFailedToCreate, expected: "failed to create"},
		{name: "ErrFailedToUpdate", constant: appErrors.ErrFailedToUpdate, expected: "failed to update"},
		{name: "ErrFailedToDelete", constant: appErrors.ErrFailedToDelete, expected: "failed to delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestIsArticleNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "article not found error",
			err:      errors.New("article not found"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appErrors.IsArticleNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsArticleNotFound() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsUserNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "user not found error",
			err:      errors.New("user not found"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appErrors.IsUserNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsUserNotFound() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsCommentNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "comment not found error",
			err:      errors.New("comment not found"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appErrors.IsCommentNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsCommentNotFound() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNotAuthorizedUpdate(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "not authorized update error",
			err:      errors.New("not authorized to update this article"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appErrors.IsNotAuthorizedUpdate(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotAuthorizedUpdate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNotAuthorizedDelete(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "not authorized delete error",
			err:      errors.New("not authorized to delete this article"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appErrors.IsNotAuthorizedDelete(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotAuthorizedDelete() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNotAuthorizedDeleteComment(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "not authorized delete comment error",
			err:      errors.New("not authorized to delete this comment"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appErrors.IsNotAuthorizedDeleteComment(tt.err)
			if result != tt.expected {
				t.Errorf("IsNotAuthorizedDeleteComment() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsCannotFollowYourself(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "cannot follow yourself error",
			err:      errors.New("cannot follow yourself"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appErrors.IsCannotFollowYourself(tt.err)
			if result != tt.expected {
				t.Errorf("IsCannotFollowYourself() = %v, want %v", result, tt.expected)
			}
		})
	}
}
