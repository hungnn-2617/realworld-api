package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"realworld-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("sets CORS headers", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.CORS())
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		// Check CORS headers
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", w.Header().Get("Access-Control-Allow-Origin"), "*")
		}
		if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
			t.Errorf("Access-Control-Allow-Credentials = %q, want %q", w.Header().Get("Access-Control-Allow-Credentials"), "true")
		}
		if w.Header().Get("Access-Control-Allow-Headers") == "" {
			t.Error("Access-Control-Allow-Headers should not be empty")
		}
		if w.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("Access-Control-Allow-Methods should not be empty")
		}

		if w.Code != http.StatusOK {
			t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("handles OPTIONS request", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.CORS())
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		r.ServeHTTP(w, req)

		// OPTIONS request should return 204 No Content
		if w.Code != http.StatusNoContent {
			t.Errorf("OPTIONS Status = %d, want %d", w.Code, http.StatusNoContent)
		}

		// Should still have CORS headers
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", w.Header().Get("Access-Control-Allow-Origin"), "*")
		}
	})

	t.Run("allows different HTTP methods", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.CORS())

		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"status": "created"})
		})
		r.PUT("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "updated"})
		})
		r.DELETE("/test", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		// Test POST
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("POST Status = %d, want %d", w.Code, http.StatusCreated)
		}

		// Test PUT
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("PUT", "/test", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("PUT Status = %d, want %d", w.Code, http.StatusOK)
		}

		// Test DELETE
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/test", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Errorf("DELETE Status = %d, want %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("CORS headers contain required values", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.CORS())
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		allowMethods := w.Header().Get("Access-Control-Allow-Methods")
		expectedMethods := []string{"POST", "OPTIONS", "GET", "PUT", "DELETE", "PATCH"}
		for _, method := range expectedMethods {
			found := false
			if containsString(allowMethods, method) {
				found = true
			}
			if !found {
				t.Errorf("Access-Control-Allow-Methods should contain %q, got %q", method, allowMethods)
			}
		}

		allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
		expectedHeaders := []string{"Content-Type", "Authorization"}
		for _, header := range expectedHeaders {
			if !containsString(allowHeaders, header) {
				t.Errorf("Access-Control-Allow-Headers should contain %q, got %q", header, allowHeaders)
			}
		}
	})
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
