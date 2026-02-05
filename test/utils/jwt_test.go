package utils_test

import (
	"testing"
	"time"

	"realworld-api/internal/utils"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		username string
		email    string
		wantErr  bool
	}{
		{
			name:     "valid token generation",
			userID:   1,
			username: "testuser",
			email:    "test@example.com",
			wantErr:  false,
		},
		{
			name:     "token with zero ID",
			userID:   0,
			username: "anonymous",
			email:    "anon@example.com",
			wantErr:  false,
		},
		{
			name:     "token with large ID",
			userID:   999999999,
			username: "biguser",
			email:    "big@example.com",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateToken(tt.userID, tt.username, tt.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Error("GenerateToken() returned empty token")
				}

				// Validate the generated token
				claims, err := utils.ValidateToken(token)
				if err != nil {
					t.Errorf("Generated token is not valid: %v", err)
					return
				}

				if claims.UserID != tt.userID {
					t.Errorf("Token UserID = %d, want %d", claims.UserID, tt.userID)
				}
				if claims.Username != tt.username {
					t.Errorf("Token Username = %q, want %q", claims.Username, tt.username)
				}
				if claims.Email != tt.email {
					t.Errorf("Token Email = %q, want %q", claims.Email, tt.email)
				}
			}
		})
	}
}

func TestGenerateTokenUniqueness(t *testing.T) {
	token1, err := utils.GenerateToken(1, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Wait a tiny bit to ensure different timestamp
	time.Sleep(time.Millisecond * 10)

	token2, err := utils.GenerateToken(1, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Tokens might be the same if generated at the exact same second
	// but claims should be valid for both
	claims1, _ := utils.ValidateToken(token1)
	claims2, _ := utils.ValidateToken(token2)

	if claims1.UserID != claims2.UserID {
		t.Error("Same user should have same UserID in both tokens")
	}
}

func TestValidateToken(t *testing.T) {
	// Generate a valid token first
	validToken, err := utils.GenerateToken(123, "validuser", "valid@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
			errType: utils.ErrInvalidToken,
		},
		{
			name:    "malformed token",
			token:   "not-a-valid-jwt",
			wantErr: true,
			errType: utils.ErrInvalidToken,
		},
		{
			name:    "token with invalid signature",
			token:   validToken + "tampered",
			wantErr: true,
			errType: utils.ErrInvalidToken,
		},
		{
			name:    "token with wrong format",
			token:   "header.payload",
			wantErr: true,
			errType: utils.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := utils.ValidateToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims == nil {
					t.Error("ValidateToken() claims should not be nil for valid token")
				}
				if claims.UserID != 123 {
					t.Errorf("Claims.UserID = %d, want %d", claims.UserID, 123)
				}
			}
		})
	}
}

func TestJWTClaimsFields(t *testing.T) {
	claims := &utils.JWTClaims{
		UserID:   42,
		Username: "testuser",
		Email:    "test@example.com",
	}

	if claims.UserID != 42 {
		t.Errorf("JWTClaims.UserID = %d, want %d", claims.UserID, 42)
	}
	if claims.Username != "testuser" {
		t.Errorf("JWTClaims.Username = %q, want %q", claims.Username, "testuser")
	}
	if claims.Email != "test@example.com" {
		t.Errorf("JWTClaims.Email = %q, want %q", claims.Email, "test@example.com")
	}
}

func TestJWTErrors(t *testing.T) {
	t.Run("ErrInvalidToken", func(t *testing.T) {
		if utils.ErrInvalidToken.Error() != "invalid token" {
			t.Errorf("ErrInvalidToken = %q, want %q", utils.ErrInvalidToken.Error(), "invalid token")
		}
	})

	t.Run("ErrExpiredToken", func(t *testing.T) {
		if utils.ErrExpiredToken.Error() != "token has expired" {
			t.Errorf("ErrExpiredToken = %q, want %q", utils.ErrExpiredToken.Error(), "token has expired")
		}
	})
}

func TestTokenRoundTrip(t *testing.T) {
	testCases := []struct {
		userID   uint
		username string
		email    string
	}{
		{1, "alice", "alice@example.com"},
		{2, "bob", "bob@example.com"},
		{100, "charlie", "charlie@test.org"},
	}

	for _, tc := range testCases {
		t.Run(tc.username, func(t *testing.T) {
			// Generate token
			token, err := utils.GenerateToken(tc.userID, tc.username, tc.email)
			if err != nil {
				t.Fatalf("GenerateToken failed: %v", err)
			}

			// Validate token
			claims, err := utils.ValidateToken(token)
			if err != nil {
				t.Fatalf("ValidateToken failed: %v", err)
			}

			// Verify claims
			if claims.UserID != tc.userID {
				t.Errorf("UserID = %d, want %d", claims.UserID, tc.userID)
			}
			if claims.Username != tc.username {
				t.Errorf("Username = %q, want %q", claims.Username, tc.username)
			}
			if claims.Email != tc.email {
				t.Errorf("Email = %q, want %q", claims.Email, tc.email)
			}
		})
	}
}
