package main

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"net/http"
)

var failedAttempts = make(map[string]int)

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Context().Value(TraceIDKey).(string) // Retrieve trace ID from context

		app.infoLog.Printf("[TraceID: %s] %s %s", traceID, r.Method, r.URL)

		_, err := app.authenticateToken(r)
		if err != nil {
			app.invalidCredentials(w)
			return
		}

		ip := r.RemoteAddr
		if failedAttempts[ip] > 5 {
			http.Error(w, "Access Denied", http.StatusForbidden)
			return
		}

		// Simulate authentication
		if r.Header.Get("Authorization") != "valid-token" {
			failedAttempts[ip]++
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" || userAgent == "curl" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TraceMiddleware adds a trace ID to each request and logs it
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for an existing trace ID, or generate a new one
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Add trace ID to the context
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)

		// Add the trace ID to the response header for visibility
		w.Header().Set("X-Trace-Id", traceID)

		// Log the trace ID
		LogWithTrace(ctx, "Incoming request")

		// Pass the request with the updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
