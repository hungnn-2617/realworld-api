package dtos_test

import (
	"encoding/json"
	"testing"

	"realworld-api/internal/dtos"
)

func TestProfileResponseJSON(t *testing.T) {
	response := dtos.ProfileResponse{
		Profile: dtos.ProfileData{
			Username:  "testuser",
			Bio:       "Test bio",
			Image:     "https://example.com/avatar.jpg",
			Following: true,
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

	profile, ok := result["profile"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'profile' key in response")
	}

	if profile["username"] != "testuser" {
		t.Errorf("Username = %v, want %v", profile["username"], "testuser")
	}
	if profile["bio"] != "Test bio" {
		t.Errorf("Bio = %v, want %v", profile["bio"], "Test bio")
	}
	if profile["image"] != "https://example.com/avatar.jpg" {
		t.Errorf("Image = %v, want %v", profile["image"], "https://example.com/avatar.jpg")
	}
	if profile["following"] != true {
		t.Errorf("Following = %v, want %v", profile["following"], true)
	}
}

func TestProfileDataStruct(t *testing.T) {
	tests := []struct {
		name       string
		profile    dtos.ProfileData
		wantUser   string
		wantBio    string
		wantImage  string
		wantFollow bool
	}{
		{
			name: "following user",
			profile: dtos.ProfileData{
				Username:  "john",
				Bio:       "John's bio",
				Image:     "https://example.com/john.jpg",
				Following: true,
			},
			wantUser:   "john",
			wantBio:    "John's bio",
			wantImage:  "https://example.com/john.jpg",
			wantFollow: true,
		},
		{
			name: "not following user",
			profile: dtos.ProfileData{
				Username:  "jane",
				Bio:       "",
				Image:     "",
				Following: false,
			},
			wantUser:   "jane",
			wantBio:    "",
			wantImage:  "",
			wantFollow: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.profile.Username != tt.wantUser {
				t.Errorf("Username = %q, want %q", tt.profile.Username, tt.wantUser)
			}
			if tt.profile.Bio != tt.wantBio {
				t.Errorf("Bio = %q, want %q", tt.profile.Bio, tt.wantBio)
			}
			if tt.profile.Image != tt.wantImage {
				t.Errorf("Image = %q, want %q", tt.profile.Image, tt.wantImage)
			}
			if tt.profile.Following != tt.wantFollow {
				t.Errorf("Following = %v, want %v", tt.profile.Following, tt.wantFollow)
			}
		})
	}
}

func TestProfileDataEmptyFields(t *testing.T) {
	profile := dtos.ProfileData{}

	if profile.Username != "" {
		t.Errorf("Empty ProfileData.Username = %q, want empty string", profile.Username)
	}
	if profile.Bio != "" {
		t.Errorf("Empty ProfileData.Bio = %q, want empty string", profile.Bio)
	}
	if profile.Image != "" {
		t.Errorf("Empty ProfileData.Image = %q, want empty string", profile.Image)
	}
	if profile.Following != false {
		t.Errorf("Empty ProfileData.Following = %v, want false", profile.Following)
	}
}
