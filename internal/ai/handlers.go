package ai

import (
	"encoding/json"
	"io"
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
	mux.HandleFunc("/api/ai/voice", h.HandleVoiceRequest)
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

// HandleVoiceRequest handles POST /api/ai/voice
// Accepts audio file (multipart/form-data) and optionally returns audio response
func (h *Handler) HandleVoiceRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get session ID
	sessionID := r.FormValue("session_id")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// Get optional parameters
	language := r.FormValue("language") // Optional: language code for speech-to-text
	returnAudio := r.FormValue("return_audio") == "true" // Optional: return audio response
	voice := r.FormValue("voice") // Optional: voice for TTS (alloy, echo, fable, onyx, nova, shimmer)
	if voice == "" {
		voice = "alloy" // Default voice
	}

	// Get audio file
	file, header, err := r.FormFile("audio")
	if err != nil {
		h.logger.Printf("Error getting audio file: %v", err)
		http.Error(w, "Audio file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read audio data
	audioData, err := io.ReadAll(file)
	if err != nil {
		h.logger.Printf("Error reading audio file: %v", err)
		http.Error(w, "Failed to read audio file", http.StatusInternalServerError)
		return
	}

	if len(audioData) == 0 {
		http.Error(w, "Audio file is empty", http.StatusBadRequest)
		return
	}

	h.logger.Printf("Received voice request: session=%s, file=%s, size=%d bytes", sessionID, header.Filename, len(audioData))

	// Convert speech to text and process through chat
	transcription, chatResp, err := h.service.HandleVoiceInput(audioData, sessionID, language)
	if err != nil {
		h.logger.Printf("Error processing voice input: %v", err)
		
		// Check if it's an OpenAI API key error
		errMsg := err.Error()
		if contains(errMsg, "Incorrect API key") || contains(errMsg, "dummy-key") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error":   "AI assistant not configured",
				"message": "Voice assistant needs an OpenAI API key to work. Please configure OPENAI_API_KEY.",
			}); err != nil {
				h.logger.Printf("Error encoding error response: %v", err)
			}
			return
		}

		http.Error(w, "Failed to process voice input", http.StatusInternalServerError)
		return
	}

	// Get the response text from chat response
	responseText := transcription
	if chatResp != nil && chatResp.Message != "" {
		responseText = chatResp.Message
	}

	// If return_audio is true, convert response to speech
	if returnAudio {
		audioResponse, err := h.service.HandleVoiceOutput(responseText, voice)
		if err != nil {
			h.logger.Printf("Error generating audio response: %v", err)
			// Fall back to text response
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"session_id":    sessionID,
				"transcription": transcription,
				"message":       responseText,
				"error":         "Failed to generate audio response, returning text",
			}
			if chatResp != nil {
				response["tokens_used"] = chatResp.TokensUsed
				response["response_time_ms"] = chatResp.ResponseTimeMs
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				h.logger.Printf("Error encoding fallback response: %v", err)
			}
			return
		}

		// Return audio response
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Disposition", "inline; filename=response.mp3")
		if _, err := w.Write(audioResponse); err != nil {
			h.logger.Printf("Error writing audio response: %v", err)
		}
		return
	}

	// Return text response with transcription and chat response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"session_id":    sessionID,
		"transcription": transcription,
		"message":       responseText,
	}
	if chatResp != nil {
		response["tokens_used"] = chatResp.TokensUsed
		response["response_time_ms"] = chatResp.ResponseTimeMs
		response["suggestions"] = chatResp.Suggestions
		response["products"] = chatResp.Products
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
