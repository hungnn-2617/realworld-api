package utils_test

import (
	"testing"

	"realworld-api/internal/utils"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt allows empty passwords
		},
		{
			name:     "long password",
			password: "verylongpassword1234567890abcdefghijklmnop",
			wantErr:  false,
		},
		{
			name:     "special characters",
			password: "p@$$w0rd!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := utils.HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify that the hash is not empty
				if hashedPassword == "" {
					t.Error("HashPassword() returned empty string")
				}

				// Verify that the hash is different from the original password
				if hashedPassword == tt.password {
					t.Error("HashPassword() returned the original password")
				}

				// Verify hash starts with bcrypt identifier
				if len(hashedPassword) < 4 || hashedPassword[:4] != "$2a$" {
					t.Error("HashPassword() did not return a valid bcrypt hash")
				}
			}
		})
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "testpassword123"

	hash1, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	hash2, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	// Same password should produce different hashes (due to random salt)
	if hash1 == hash2 {
		t.Error("HashPassword() should produce different hashes for the same password")
	}
}

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		inputPassword  string
		expectedResult bool
	}{
		{
			name:           "correct password",
			password:       "password123",
			inputPassword:  "password123",
			expectedResult: true,
		},
		{
			name:           "incorrect password",
			password:       "password123",
			inputPassword:  "wrongpassword",
			expectedResult: false,
		},
		{
			name:           "empty password check",
			password:       "",
			inputPassword:  "",
			expectedResult: true,
		},
		{
			name:           "empty input against non-empty password",
			password:       "password123",
			inputPassword:  "",
			expectedResult: false,
		},
		{
			name:           "case sensitive password",
			password:       "Password123",
			inputPassword:  "password123",
			expectedResult: false,
		},
		{
			name:           "special characters",
			password:       "p@$$w0rd!",
			inputPassword:  "p@$$w0rd!",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := utils.HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			result := utils.CheckPassword(tt.inputPassword, hashedPassword)

			if result != tt.expectedResult {
				t.Errorf("CheckPassword() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestCheckPasswordWithInvalidHash(t *testing.T) {
	result := utils.CheckPassword("password123", "invalid_hash")
	if result {
		t.Error("CheckPassword() should return false for invalid hash")
	}
}

func TestCheckPasswordWithEmptyHash(t *testing.T) {
	result := utils.CheckPassword("password123", "")
	if result {
		t.Error("CheckPassword() should return false for empty hash")
	}
}
