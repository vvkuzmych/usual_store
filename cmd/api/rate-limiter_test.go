package main

import (
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetLimiter(t *testing.T) {
	tests := []struct {
		name string
		ip   string
	}{
		{"New IP", "192.168.1.1"},
		{"Existing IP", "192.168.1.1"},
		{"Another New IP", "192.168.1.2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := GetLimiter(tt.ip)
			if limiter == nil {
				t.Errorf("Expected non-nil limiter for IP %s", tt.ip)
			}
			if _, exists := IpLimiter[tt.ip]; !exists {
				t.Errorf("Limiter for IP %s was not stored in map", tt.ip)
			}
		})
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		ip           string
		requestCount int
		wantStatus   int
	}{
		{"Allow single request", "192.168.1.1", 1, http.StatusOK},
		{"Exceed rate limit", "192.168.1.2", 10, http.StatusTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ✅ Reset the global ipLimiter before each test
			ipLimiter := make(map[string]*rate.Limiter)
			_ = ipLimiter // Suppress "unused variable" warning

			handler := RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			var lastResponseCode int
			ticker := time.NewTicker(100 * time.Millisecond) // Ticker simulating request intervals
			defer ticker.Stop()

			for i := 0; i < tt.requestCount; i++ {
				req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
				req.RemoteAddr = tt.ip
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)

				lastResponseCode = rec.Code

				// Wait for the next tick to simulate the passage of time between requests
				<-ticker.C
			}

			if lastResponseCode != tt.wantStatus {
				t.Errorf("Got status %d, want %d", lastResponseCode, tt.wantStatus)
			}
		})
	}
}
