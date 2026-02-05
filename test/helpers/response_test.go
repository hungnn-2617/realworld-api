package helpers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"realworld-api/internal/helpers"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		messages []string
		expected helpers.ErrorResponse
	}{
		{
			name:     "single message",
			messages: []string{"error message"},
			expected: helpers.ErrorResponse{
				Errors: helpers.ErrorBody{
					Body: []string{"error message"},
				},
			},
		},
		{
			name:     "multiple messages",
			messages: []string{"error 1", "error 2"},
			expected: helpers.ErrorResponse{
				Errors: helpers.ErrorBody{
					Body: []string{"error 1", "error 2"},
				},
			},
		},
		{
			name:     "empty messages",
			messages: []string{},
			expected: helpers.ErrorResponse{
				Errors: helpers.ErrorBody{
					Body: []string{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.NewErrorResponse(tt.messages...)

			if len(result.Errors.Body) != len(tt.expected.Errors.Body) {
				t.Errorf("NewErrorResponse() body length = %d, want %d", len(result.Errors.Body), len(tt.expected.Errors.Body))
			}

			for i, msg := range result.Errors.Body {
				if msg != tt.expected.Errors.Body[i] {
					t.Errorf("NewErrorResponse() body[%d] = %q, want %q", i, msg, tt.expected.Errors.Body[i])
				}
			}
		})
	}
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		messages       []string
		expectedStatus int
	}{
		{
			name:           "bad request",
			statusCode:     http.StatusBadRequest,
			messages:       []string{"bad request"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "internal server error",
			statusCode:     http.StatusInternalServerError,
			messages:       []string{"internal error"},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			helpers.RespondWithError(c, tt.statusCode, tt.messages...)

			if w.Code != tt.expectedStatus {
				t.Errorf("RespondWithError() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var response helpers.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if len(response.Errors.Body) != len(tt.messages) {
				t.Errorf("RespondWithError() body length = %d, want %d", len(response.Errors.Body), len(tt.messages))
			}
		})
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	helpers.BadRequest(c, "invalid input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("BadRequest() status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var response helpers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Errors.Body) != 1 || response.Errors.Body[0] != "invalid input" {
		t.Errorf("BadRequest() unexpected body: %v", response.Errors.Body)
	}
}

func TestUnauthorized(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helpers.Unauthorized(c, "custom message")

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Unauthorized() status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		var response helpers.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response.Errors.Body) != 1 || response.Errors.Body[0] != "custom message" {
			t.Errorf("Unauthorized() unexpected body: %v", response.Errors.Body)
		}
	})

	t.Run("without message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helpers.Unauthorized(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Unauthorized() status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		var response helpers.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response.Errors.Body) != 1 || response.Errors.Body[0] != "unauthorized" {
			t.Errorf("Unauthorized() expected default message 'unauthorized', got: %v", response.Errors.Body)
		}
	})
}

func TestForbidden(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helpers.Forbidden(c, "access denied")

		if w.Code != http.StatusForbidden {
			t.Errorf("Forbidden() status = %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("without message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helpers.Forbidden(c)

		var response helpers.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response.Errors.Body) != 1 || response.Errors.Body[0] != "forbidden" {
			t.Errorf("Forbidden() expected default message 'forbidden', got: %v", response.Errors.Body)
		}
	})
}

func TestNotFound(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helpers.NotFound(c, "user not found")

		if w.Code != http.StatusNotFound {
			t.Errorf("NotFound() status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("without message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helpers.NotFound(c)

		var response helpers.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response.Errors.Body) != 1 || response.Errors.Body[0] != "not found" {
			t.Errorf("NotFound() expected default message 'not found', got: %v", response.Errors.Body)
		}
	})
}

func TestConflict(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	helpers.Conflict(c, "resource already exists")

	if w.Code != http.StatusConflict {
		t.Errorf("Conflict() status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestUnprocessableEntity(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	helpers.UnprocessableEntity(c, "validation failed")

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("UnprocessableEntity() status = %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
}

func TestInternalServerError(t *testing.T) {
	t.Run("with message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helpers.InternalServerError(c, "database error")

		if w.Code != http.StatusInternalServerError {
			t.Errorf("InternalServerError() status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("without message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		helpers.InternalServerError(c)

		var response helpers.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response.Errors.Body) != 1 || response.Errors.Body[0] != "internal server error" {
			t.Errorf("InternalServerError() expected default message, got: %v", response.Errors.Body)
		}
	})
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"message": "success"}
	helpers.Success(c, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("Success() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "success" {
		t.Errorf("Success() unexpected body: %v", response)
	}
}

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]int{"count": 10}
	helpers.OK(c, data)

	if w.Code != http.StatusOK {
		t.Errorf("OK() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"id": "123"}
	helpers.Created(c, data)

	if w.Code != http.StatusCreated {
		t.Errorf("Created() status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestNoContent(t *testing.T) {
	r := gin.New()
	r.DELETE("/test", func(c *gin.Context) {
		helpers.NoContent(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("NoContent() status = %d, want %d", w.Code, http.StatusNoContent)
	}

	if w.Body.Len() != 0 {
		t.Errorf("NoContent() should have empty body, got: %s", w.Body.String())
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	helpers.Error(c, http.StatusBadRequest, "error message")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Error() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"status": "ok"}
	helpers.SuccessResponse(c, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("SuccessResponse() status = %d, want %d", w.Code, http.StatusOK)
	}
}
