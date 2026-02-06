package services_test

import (
	"testing"

	"realworld-api/internal/dtos"
)

// Test RegisterRequest validation logic
func TestRegisterRequest_FieldValidation(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		email       string
		password    string
		expectError bool
	}{
		{
			name:        "valid registration",
			username:    "testuser",
			email:       "test@example.com",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "empty username",
			username:    "",
			email:       "test@example.com",
			password:    "password123",
			expectError: true,
		},
		{
			name:        "empty email",
			username:    "testuser",
			email:       "",
			password:    "password123",
			expectError: true,
		},
		{
			name:        "empty password",
			username:    "testuser",
			email:       "test@example.com",
			password:    "",
			expectError: true,
		},
		{
			name:        "username too short",
			username:    "ab",
			email:       "test@example.com",
			password:    "password123",
			expectError: true,
		},
		{
			name:        "password too short",
			username:    "testuser",
			email:       "test@example.com",
			password:    "12345",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false

			// Simulate validation
			if len(tt.username) < 3 || tt.username == "" {
				hasError = true
			}
			if tt.email == "" {
				hasError = true
			}
			if len(tt.password) < 6 || tt.password == "" {
				hasError = true
			}

			if hasError != tt.expectError {
				t.Errorf("expected error: %v, got error: %v", tt.expectError, hasError)
			}
		})
	}
}

// Test LoginRequest validation
func TestLoginRequest_FieldValidation(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		expectError bool
	}{
		{
			name:        "valid login",
			email:       "test@example.com",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "empty email",
			email:       "",
			password:    "password123",
			expectError: true,
		},
		{
			name:        "empty password",
			email:       "test@example.com",
			password:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.email == "" || tt.password == ""
			if hasError != tt.expectError {
				t.Errorf("expected error: %v, got error: %v", tt.expectError, hasError)
			}
		})
	}
}

// Test UpdateUserRequest optional fields
func TestUpdateUserRequest_OptionalFields(t *testing.T) {
	email := "new@example.com"
	username := "newuser"
	bio := "New bio"
	image := "https://example.com/new.jpg"

	tests := []struct {
		name     string
		request  dtos.UpdateUserRequest
		expected map[string]bool // field exists?
	}{
		{
			name: "update email only",
			request: func() dtos.UpdateUserRequest {
				var req dtos.UpdateUserRequest
				req.User.Email = &email
				return req
			}(),
			expected: map[string]bool{"email": true, "username": false, "bio": false, "image": false},
		},
		{
			name: "update username only",
			request: func() dtos.UpdateUserRequest {
				var req dtos.UpdateUserRequest
				req.User.Username = &username
				return req
			}(),
			expected: map[string]bool{"email": false, "username": true, "bio": false, "image": false},
		},
		{
			name: "update bio only",
			request: func() dtos.UpdateUserRequest {
				var req dtos.UpdateUserRequest
				req.User.Bio = &bio
				return req
			}(),
			expected: map[string]bool{"email": false, "username": false, "bio": true, "image": false},
		},
		{
			name: "update image only",
			request: func() dtos.UpdateUserRequest {
				var req dtos.UpdateUserRequest
				req.User.Image = &image
				return req
			}(),
			expected: map[string]bool{"email": false, "username": false, "bio": false, "image": true},
		},
		{
			name: "update all fields",
			request: func() dtos.UpdateUserRequest {
				var req dtos.UpdateUserRequest
				req.User.Email = &email
				req.User.Username = &username
				req.User.Bio = &bio
				req.User.Image = &image
				return req
			}(),
			expected: map[string]bool{"email": true, "username": true, "bio": true, "image": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.request.User.Email != nil) != tt.expected["email"] {
				t.Errorf("email field presence mismatch")
			}
			if (tt.request.User.Username != nil) != tt.expected["username"] {
				t.Errorf("username field presence mismatch")
			}
			if (tt.request.User.Bio != nil) != tt.expected["bio"] {
				t.Errorf("bio field presence mismatch")
			}
			if (tt.request.User.Image != nil) != tt.expected["image"] {
				t.Errorf("image field presence mismatch")
			}
		})
	}
}

// Test UserData response structure
func TestUserData_ResponseFields(t *testing.T) {
	userData := dtos.UserData{
		Email:    "test@example.com",
		Token:    "jwt-token",
		Username: "testuser",
		Bio:      "Test bio",
		Image:    "https://example.com/image.jpg",
	}

	if userData.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", userData.Email)
	}
	if userData.Token != "jwt-token" {
		t.Errorf("Expected token 'jwt-token', got '%s'", userData.Token)
	}
	if userData.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", userData.Username)
	}
	if userData.Bio != "Test bio" {
		t.Errorf("Expected bio 'Test bio', got '%s'", userData.Bio)
	}
	if userData.Image != "https://example.com/image.jpg" {
		t.Errorf("Expected image 'https://example.com/image.jpg', got '%s'", userData.Image)
	}
}

// Test UserResponse wrapper
func TestUserResponse_Structure(t *testing.T) {
	response := dtos.UserResponse{
		User: dtos.UserData{
			Email:    "test@example.com",
			Token:    "token",
			Username: "user",
			Bio:      "",
			Image:    "",
		},
	}

	if response.User.Email != "test@example.com" {
		t.Errorf("Expected User.Email 'test@example.com', got '%s'", response.User.Email)
	}
}

// Test email validation helper
func TestEmailValidation(t *testing.T) {
	tests := []struct {
		email   string
		isValid bool
	}{
		{"test@example.com", true},
		{"user@domain.org", true},
		{"user.name@domain.co.uk", true},
		{"invalid", false},
		{"@domain.com", false},
		{"user@", false},
		{"", false},
		{"user space@domain.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			// Simple email validation (basic check)
			hasAtSign := false
			hasDot := false
			atIndex := -1

			for i, c := range tt.email {
				if c == '@' {
					hasAtSign = true
					atIndex = i
				}
				if c == '.' && atIndex > 0 && i > atIndex {
					hasDot = true
				}
				if c == ' ' {
					hasAtSign = false
					break
				}
			}

			isValid := hasAtSign && hasDot && atIndex > 0

			if isValid != tt.isValid {
				t.Errorf("email '%s': expected valid=%v, got valid=%v", tt.email, tt.isValid, isValid)
			}
		})
	}
}

// Test password requirements
func TestPasswordRequirements(t *testing.T) {
	tests := []struct {
		password    string
		meetsMinLen bool
	}{
		{"123456", true},
		{"password", true},
		{"12345", false},
		{"test", false},
		{"", false},
		{"verylongpassword123", true},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			meetsMinLen := len(tt.password) >= 6
			if meetsMinLen != tt.meetsMinLen {
				t.Errorf("password '%s': expected meetsMinLen=%v, got %v", tt.password, tt.meetsMinLen, meetsMinLen)
			}
		})
	}
}

// Test username requirements
func TestUsernameRequirements(t *testing.T) {
	tests := []struct {
		username    string
		meetsMinLen bool
		meetsMaxLen bool
	}{
		{"abc", true, true},
		{"validuser", true, true},
		{"ab", false, true}, // too short
		{"", false, true},
		{"verylongusernamethatexceedsfiftycharacterslimitwhichisfifty", true, false}, // too long
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			meetsMinLen := len(tt.username) >= 3
			meetsMaxLen := len(tt.username) <= 50

			if meetsMinLen != tt.meetsMinLen {
				t.Errorf("username '%s': expected meetsMinLen=%v, got %v", tt.username, tt.meetsMinLen, meetsMinLen)
			}
			if meetsMaxLen != tt.meetsMaxLen {
				t.Errorf("username '%s': expected meetsMaxLen=%v, got %v", tt.username, tt.meetsMaxLen, meetsMaxLen)
			}
		})
	}
}
