package dtos_test

import (
	"encoding/json"
	"testing"

	"realworld-api/internal/dtos"
)

func TestErrorResponseJSON(t *testing.T) {
	tests := []struct {
		name     string
		response dtos.ErrorResponse
		expected []string
	}{
		{
			name: "single error",
			response: dtos.ErrorResponse{
				Errors: dtos.ErrorBody{
					Body: []string{"error message"},
				},
			},
			expected: []string{"error message"},
		},
		{
			name: "multiple errors",
			response: dtos.ErrorResponse{
				Errors: dtos.ErrorBody{
					Body: []string{"error 1", "error 2", "error 3"},
				},
			},
			expected: []string{"error 1", "error 2", "error 3"},
		},
		{
			name: "empty errors",
			response: dtos.ErrorResponse{
				Errors: dtos.ErrorBody{
					Body: []string{},
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			var result map[string]interface{}
			err = json.Unmarshal(jsonBytes, &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			errors, ok := result["errors"].(map[string]interface{})
			if !ok {
				t.Fatal("Expected 'errors' key in response")
			}

			body, ok := errors["body"].([]interface{})
			if !ok {
				t.Fatal("Expected 'body' key as array in errors")
			}

			if len(body) != len(tt.expected) {
				t.Errorf("Body length = %d, want %d", len(body), len(tt.expected))
			}

			for i, msg := range body {
				if msg != tt.expected[i] {
					t.Errorf("Body[%d] = %v, want %v", i, msg, tt.expected[i])
				}
			}
		})
	}
}

func TestErrorBodyStruct(t *testing.T) {
	t.Run("with messages", func(t *testing.T) {
		errorBody := dtos.ErrorBody{
			Body: []string{"validation failed", "email is invalid"},
		}

		if len(errorBody.Body) != 2 {
			t.Errorf("Body length = %d, want %d", len(errorBody.Body), 2)
		}
		if errorBody.Body[0] != "validation failed" {
			t.Errorf("Body[0] = %q, want %q", errorBody.Body[0], "validation failed")
		}
		if errorBody.Body[1] != "email is invalid" {
			t.Errorf("Body[1] = %q, want %q", errorBody.Body[1], "email is invalid")
		}
	})

	t.Run("empty body", func(t *testing.T) {
		errorBody := dtos.ErrorBody{
			Body: []string{},
		}

		if len(errorBody.Body) != 0 {
			t.Errorf("Body length = %d, want %d", len(errorBody.Body), 0)
		}
	})

	t.Run("nil body", func(t *testing.T) {
		errorBody := dtos.ErrorBody{}

		if errorBody.Body != nil {
			t.Errorf("Body should be nil")
		}
	})
}

func TestErrorResponseUnmarshal(t *testing.T) {
	jsonStr := `{"errors":{"body":["unauthorized","token expired"]}}`

	var response dtos.ErrorResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(response.Errors.Body) != 2 {
		t.Errorf("Errors.Body length = %d, want %d", len(response.Errors.Body), 2)
	}
	if response.Errors.Body[0] != "unauthorized" {
		t.Errorf("Errors.Body[0] = %q, want %q", response.Errors.Body[0], "unauthorized")
	}
	if response.Errors.Body[1] != "token expired" {
		t.Errorf("Errors.Body[1] = %q, want %q", response.Errors.Body[1], "token expired")
	}
}
