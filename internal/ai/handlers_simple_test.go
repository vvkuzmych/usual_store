package ai

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Test HTTP method validation
func TestHandleChatRequest_HTTPMethods(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{"GET not allowed", http.MethodGet, http.StatusMethodNotAllowed},
		{"PUT not allowed", http.MethodPut, http.StatusMethodNotAllowed},
		{"DELETE not allowed", http.MethodDelete, http.StatusMethodNotAllowed},
		{"PATCH not allowed", http.MethodPatch, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{logger: logger}
			handler := NewHandler(service, logger)

			req := httptest.NewRequest(tt.method, "/api/ai/chat", nil)
			w := httptest.NewRecorder()

			handler.HandleChatRequest(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test request validation
func TestHandleChatRequest_Validation(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "invalid JSON",
			body:           "not json",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Invalid request body",
		},
		{
			name:           "missing message",
			body:           map[string]interface{}{"session_id": "test"},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Message is required",
		},
		{
			name:           "missing session ID",
			body:           map[string]interface{}{"message": "test"},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Session ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{logger: logger}
			handler := NewHandler(service, logger)

			var body []byte
			if str, ok := tt.body.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/ai/chat", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.HandleChatRequest(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expectedMsg) {
				t.Errorf("Expected body to contain %q, got %q", tt.expectedMsg, w.Body.String())
			}
		})
	}
}

// Test feedback HTTP methods
func TestHandleFeedback_HTTPMethods(t *testing.T) {
	tests := []struct {
		method         string
		expectedStatus int
	}{
		{http.MethodGet, http.StatusMethodNotAllowed},
		{http.MethodPut, http.StatusMethodNotAllowed},
		{http.MethodDelete, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{logger: logger}
			handler := NewHandler(service, logger)

			req := httptest.NewRequest(tt.method, "/api/ai/feedback", nil)
			w := httptest.NewRecorder()

			handler.HandleFeedback(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test feedback request validation
func TestHandleFeedback_InvalidJSON(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := &Service{logger: logger}
	handler := NewHandler(service, logger)

	req := httptest.NewRequest(http.MethodPost, "/api/ai/feedback", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.HandleFeedback(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Test stats HTTP methods
func TestHandleStats_HTTPMethods(t *testing.T) {
	tests := []struct {
		method         string
		expectedStatus int
	}{
		{http.MethodPost, http.StatusMethodNotAllowed},
		{http.MethodPut, http.StatusMethodNotAllowed},
		{http.MethodDelete, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{logger: logger}
			handler := NewHandler(service, logger)

			req := httptest.NewRequest(tt.method, "/api/ai/stats", nil)
			w := httptest.NewRecorder()

			handler.HandleStats(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test CORS middleware
func TestEnableCORS_Headers(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := &Service{logger: logger}
	handler := NewHandler(service, logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := handler.EnableCORS(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	corsHandler(w, req)

	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected Access-Control-Allow-Origin header")
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, OPTIONS" {
		t.Error("Expected Access-Control-Allow-Methods header")
	}
}

// Test OPTIONS preflight
func TestEnableCORS_Preflight(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := &Service{logger: logger}
	handler := NewHandler(service, logger)

	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	corsHandler := handler.EnableCORS(nextHandler)

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()

	corsHandler(w, req)

	if nextCalled {
		t.Error("Next handler should not be called for OPTIONS request")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// Test contains helper
func TestContainsHelper(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"Hello World", "World", true},
		{"Hello World", "Universe", false},
		{"", "", true},
		{"test", "", true},
		{"", "test", false},
	}

	for _, tt := range tests {
		result := contains(tt.s, tt.substr)
		if result != tt.expected {
			t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
		}
	}
}

// Test NewHandler creation
func TestNewHandlerCreation(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := &Service{logger: logger}

	handler := NewHandler(service, logger)

	if handler == nil {
		t.Fatal("NewHandler returned nil")
	}

	if handler.service == nil {
		t.Error("Handler service is nil")
	}

	if handler.logger == nil {
		t.Error("Handler logger is nil")
	}
}

// Test route registration
func TestRouteRegistration(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := &Service{logger: logger}
	handler := NewHandler(service, logger)

	mux := http.NewServeMux()

	// Just call RegisterRoutes - it registers the routes
	handler.RegisterRoutes(mux)

	// Verify routes are registered by checking they don't panic
	// We can't easily test the actual routing without calling handlers
	// but the RegisterRoutes method itself is simple and tested indirectly
	// by other tests that use the actual routes

	// This test mainly ensures RegisterRoutes doesn't panic
	t.Log("RegisterRoutes completed without panic")
}
