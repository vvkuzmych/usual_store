package ai

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Handler handles HTTP requests for AI assistant
type Handler struct {
	service *Service
	logger  *log.Logger
}

// NewHandler creates a new AI handler
func NewHandler(service *Service, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// HandleChatRequest handles POST /api/ai/chat
func (h *Handler) HandleChatRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Extract user ID from session/JWT token if authenticated

	// Process chat
	resp, err := h.service.HandleChat(req)
	if err != nil {
		h.logger.Printf("Error handling chat: %v", err)

		// Check if it's an OpenAI API key error
		errMsg := err.Error()
		if contains(errMsg, "Incorrect API key") || contains(errMsg, "dummy-key") {
			// Return user-friendly message for missing API key
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error":   "AI assistant not configured",
				"message": "Hi! The AI assistant is running but needs an OpenAI API key to respond. For now, I can help you browse our products! Try clicking on 'Products' in the menu above.",
			}); err != nil {
				h.logger.Printf("Error encoding error response: %v", err)
			}
			return
		}

		http.Error(w, "Failed to process message", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleFeedback handles POST /api/ai/feedback
func (h *Handler) HandleFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var feedback Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		h.logger.Printf("Error decoding feedback: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.SubmitFeedback(feedback); err != nil {
		h.logger.Printf("Error saving feedback: %v", err)
		http.Error(w, "Failed to save feedback", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Thank you for your feedback!",
	}); err != nil {
		h.logger.Printf("Error encoding feedback response: %v", err)
	}
}

// HandleStats handles GET /api/ai/stats
func (h *Handler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get days parameter (default to 7)
	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	stats, err := h.service.GetConversationStats(days)
	if err != nil {
		h.logger.Printf("Error getting stats: %v", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		h.logger.Printf("Error encoding stats: %v", err)
	}
}

// RegisterRoutes registers all AI assistant routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ai/chat", h.HandleChatRequest)
	mux.HandleFunc("/api/ai/feedback", h.HandleFeedback)
	mux.HandleFunc("/api/ai/stats", h.HandleStats)
}

// CORS middleware
func (h *Handler) EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
