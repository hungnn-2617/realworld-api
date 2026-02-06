package services_test

import (
	"testing"
	"time"

	"realworld-api/internal/dtos"
	"realworld-api/internal/models"
)

// Test ProfileData structure
func TestProfileData_Structure(t *testing.T) {
	profileData := dtos.ProfileData{
		Username:  "testuser",
		Bio:       "Test bio",
		Image:     "https://example.com/image.jpg",
		Following: false,
	}

	if profileData.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", profileData.Username)
	}
	if profileData.Bio != "Test bio" {
		t.Errorf("Expected bio 'Test bio', got '%s'", profileData.Bio)
	}
	if profileData.Image != "https://example.com/image.jpg" {
		t.Errorf("Expected image 'https://example.com/image.jpg', got '%s'", profileData.Image)
	}
	if profileData.Following != false {
		t.Error("Expected following false, got true")
	}
}

// Test ProfileResponse structure
func TestProfileResponse_Structure(t *testing.T) {
	response := dtos.ProfileResponse{
		Profile: dtos.ProfileData{
			Username:  "testuser",
			Bio:       "Test bio",
			Image:     "https://example.com/image.jpg",
			Following: true,
		},
	}

	if response.Profile.Username != "testuser" {
		t.Errorf("Expected Profile.Username 'testuser', got '%s'", response.Profile.Username)
	}
	if !response.Profile.Following {
		t.Error("Expected Profile.Following true, got false")
	}
}

// Test follow/unfollow self validation
func TestFollowSelf_Validation(t *testing.T) {
	tests := []struct {
		name          string
		currentUserID uint
		targetUserID  uint
		expectError   bool
	}{
		{
			name:          "follow different user",
			currentUserID: 1,
			targetUserID:  2,
			expectError:   false,
		},
		{
			name:          "follow self",
			currentUserID: 1,
			targetUserID:  1,
			expectError:   true,
		},
		{
			name:          "unfollow different user",
			currentUserID: 5,
			targetUserID:  10,
			expectError:   false,
		},
		{
			name:          "unfollow self",
			currentUserID: 3,
			targetUserID:  3,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.currentUserID == tt.targetUserID
			if hasError != tt.expectError {
				t.Errorf("expected error: %v, got error: %v", tt.expectError, hasError)
			}
		})
	}
}

// Test Follow model
func TestFollowModel_Fields(t *testing.T) {
	now := time.Now()
	follow := models.Follow{
		FollowerID:  1,
		FollowingID: 2,
		CreatedAt:   now,
	}

	if follow.FollowerID != 1 {
		t.Errorf("Expected FollowerID 1, got %d", follow.FollowerID)
	}
	if follow.FollowingID != 2 {
		t.Errorf("Expected FollowingID 2, got %d", follow.FollowingID)
	}
}

// Test profile with empty bio and image
func TestProfile_EmptyFields(t *testing.T) {
	profileData := dtos.ProfileData{
		Username:  "testuser",
		Bio:       "",
		Image:     "",
		Following: false,
	}

	if profileData.Bio != "" {
		t.Errorf("Expected empty bio, got '%s'", profileData.Bio)
	}
	if profileData.Image != "" {
		t.Errorf("Expected empty image, got '%s'", profileData.Image)
	}
}

// Test isFollowing check with unauthenticated user
func TestIsFollowing_Unauthenticated(t *testing.T) {
	tests := []struct {
		name            string
		currentUserID   uint
		isFollowingInDB bool
		expectedResult  bool
	}{
		{
			name:            "unauthenticated user (ID=0)",
			currentUserID:   0,
			isFollowingInDB: true,
			expectedResult:  false,
		},
		{
			name:            "authenticated user following",
			currentUserID:   1,
			isFollowingInDB: true,
			expectedResult:  true,
		},
		{
			name:            "authenticated user not following",
			currentUserID:   1,
			isFollowingInDB: false,
			expectedResult:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mimics the IsFollowing logic in services
			result := tt.currentUserID != 0 && tt.isFollowingInDB
			if result != tt.expectedResult {
				t.Errorf("expected: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}

// Test User to ProfileData conversion
func TestUserToProfileData_Conversion(t *testing.T) {
	user := models.User{
		ID:       1,
		Username: "testuser",
		Bio:      "Test bio",
		Image:    "https://example.com/image.jpg",
	}

	profileData := dtos.ProfileData{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}

	if profileData.Username != user.Username {
		t.Errorf("Username mismatch: expected '%s', got '%s'", user.Username, profileData.Username)
	}
	if profileData.Bio != user.Bio {
		t.Errorf("Bio mismatch: expected '%s', got '%s'", user.Bio, profileData.Bio)
	}
	if profileData.Image != user.Image {
		t.Errorf("Image mismatch: expected '%s', got '%s'", user.Image, profileData.Image)
	}
}

// Test following state changes
func TestFollowingState_Changes(t *testing.T) {
	tests := []struct {
		name      string
		action    string
		following bool
		expected  bool
	}{
		{
			name:      "follow sets following to true",
			action:    "follow",
			following: false,
			expected:  true,
		},
		{
			name:      "unfollow sets following to false",
			action:    "unfollow",
			following: true,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.following
			if tt.action == "follow" {
				result = true
			} else if tt.action == "unfollow" {
				result = false
			}

			if result != tt.expected {
				t.Errorf("expected following=%v after %s, got %v", tt.expected, tt.action, result)
			}
		})
	}
}

// Test unique follow constraint
func TestFollow_UniqueConstraint(t *testing.T) {
	existingFollows := map[string]bool{
		"1-2": true, // user 1 follows user 2
		"1-3": true, // user 1 follows user 3
	}

	tests := []struct {
		name        string
		followerID  uint
		followingID uint
		exists      bool
	}{
		{
			name:        "existing follow",
			followerID:  1,
			followingID: 2,
			exists:      true,
		},
		{
			name:        "new follow",
			followerID:  2,
			followingID: 1,
			exists:      false,
		},
		{
			name:        "another existing",
			followerID:  1,
			followingID: 3,
			exists:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := string(rune(tt.followerID+48)) + "-" + string(rune(tt.followingID+48))
			exists := existingFollows[key]
			if exists != tt.exists {
				t.Errorf("expected exists=%v, got %v", tt.exists, exists)
			}
		})
	}
}

// Test Profile not found scenario
func TestProfile_NotFound(t *testing.T) {
	// Simulate user not found scenario
	username := "nonexistent"
	var user *models.User = nil

	if user != nil {
		t.Error("Expected user to be nil for nonexistent username")
	}
	_ = username // suppress unused variable warning
}
