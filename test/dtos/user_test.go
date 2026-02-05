package dtos_test

import (
	"encoding/json"
	"testing"

	"realworld-api/internal/dtos"
)

func TestRegisterRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected dtos.RegisterRequest
	}{
		{
			name: "valid register request",
			json: `{"user":{"username":"testuser","email":"test@example.com","password":"password123"}}`,
			expected: func() dtos.RegisterRequest {
				r := dtos.RegisterRequest{}
				r.User.Username = "testuser"
				r.User.Email = "test@example.com"
				r.User.Password = "password123"
				return r
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result dtos.RegisterRequest
			err := json.Unmarshal([]byte(tt.json), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if result.User.Username != tt.expected.User.Username {
				t.Errorf("Username = %q, want %q", result.User.Username, tt.expected.User.Username)
			}
			if result.User.Email != tt.expected.User.Email {
				t.Errorf("Email = %q, want %q", result.User.Email, tt.expected.User.Email)
			}
			if result.User.Password != tt.expected.User.Password {
				t.Errorf("Password = %q, want %q", result.User.Password, tt.expected.User.Password)
			}
		})
	}
}

func TestLoginRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected dtos.LoginRequest
	}{
		{
			name: "valid login request",
			json: `{"user":{"email":"test@example.com","password":"password123"}}`,
			expected: func() dtos.LoginRequest {
				r := dtos.LoginRequest{}
				r.User.Email = "test@example.com"
				r.User.Password = "password123"
				return r
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result dtos.LoginRequest
			err := json.Unmarshal([]byte(tt.json), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if result.User.Email != tt.expected.User.Email {
				t.Errorf("Email = %q, want %q", result.User.Email, tt.expected.User.Email)
			}
			if result.User.Password != tt.expected.User.Password {
				t.Errorf("Password = %q, want %q", result.User.Password, tt.expected.User.Password)
			}
		})
	}
}

func TestUpdateUserRequestJSON(t *testing.T) {
	t.Run("full update", func(t *testing.T) {
		jsonStr := `{"user":{"username":"newuser","email":"new@example.com","password":"newpass","bio":"New bio","image":"https://example.com/img.jpg"}}`
		var result dtos.UpdateUserRequest
		err := json.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if result.User.Username == nil || *result.User.Username != "newuser" {
			t.Errorf("Username unexpected value")
		}
		if result.User.Email == nil || *result.User.Email != "new@example.com" {
			t.Errorf("Email unexpected value")
		}
		if result.User.Password == nil || *result.User.Password != "newpass" {
			t.Errorf("Password unexpected value")
		}
		if result.User.Bio == nil || *result.User.Bio != "New bio" {
			t.Errorf("Bio unexpected value")
		}
		if result.User.Image == nil || *result.User.Image != "https://example.com/img.jpg" {
			t.Errorf("Image unexpected value")
		}
	})

	t.Run("partial update", func(t *testing.T) {
		jsonStr := `{"user":{"username":"newuser"}}`
		var result dtos.UpdateUserRequest
		err := json.Unmarshal([]byte(jsonStr), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if result.User.Username == nil || *result.User.Username != "newuser" {
			t.Errorf("Username unexpected value")
		}
		if result.User.Email != nil {
			t.Errorf("Email should be nil for partial update")
		}
		if result.User.Password != nil {
			t.Errorf("Password should be nil for partial update")
		}
	})
}

func TestUserResponseJSON(t *testing.T) {
	response := dtos.UserResponse{
		User: dtos.UserData{
			Email:    "test@example.com",
			Token:    "jwt-token-here",
			Username: "testuser",
			Bio:      "Test bio",
			Image:    "https://example.com/avatar.jpg",
		},
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	user, ok := result["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'user' key in response")
	}

	if user["email"] != "test@example.com" {
		t.Errorf("Email = %v, want %v", user["email"], "test@example.com")
	}
	if user["token"] != "jwt-token-here" {
		t.Errorf("Token = %v, want %v", user["token"], "jwt-token-here")
	}
	if user["username"] != "testuser" {
		t.Errorf("Username = %v, want %v", user["username"], "testuser")
	}
}

func TestUserDataStruct(t *testing.T) {
	userData := dtos.UserData{
		Email:    "user@test.com",
		Token:    "token123",
		Username: "username",
		Bio:      "User bio",
		Image:    "https://image.url",
	}

	if userData.Email != "user@test.com" {
		t.Errorf("Email = %q, want %q", userData.Email, "user@test.com")
	}
	if userData.Token != "token123" {
		t.Errorf("Token = %q, want %q", userData.Token, "token123")
	}
	if userData.Username != "username" {
		t.Errorf("Username = %q, want %q", userData.Username, "username")
	}
	if userData.Bio != "User bio" {
		t.Errorf("Bio = %q, want %q", userData.Bio, "User bio")
	}
	if userData.Image != "https://image.url" {
		t.Errorf("Image = %q, want %q", userData.Image, "https://image.url")
	}
}
