package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ChatRequest struct {
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id,omitempty"`
	UserID         int    `json:"user_id,omitempty"`
}

type ChatResponse struct {
	Response       string    `json:"response"`
	ConversationID string    `json:"conversation_id"`
	Timestamp      time.Time `json:"timestamp"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "ai-assistant",
		}); err != nil {
			log.Printf("Error encoding health response: %v", err)
		}
	})

	// Chat endpoint (placeholder implementation)
	http.HandleFunc("/api/ai/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Placeholder response - in production, this would call OpenAI API
		response := ChatResponse{
			Response:       fmt.Sprintf("AI Assistant received your message: '%s'. (This is a placeholder - OpenAI integration pending)", req.Message),
			ConversationID: req.ConversationID,
			Timestamp:      time.Now(),
		}

		if response.ConversationID == "" {
			response.ConversationID = fmt.Sprintf("conv_%d", time.Now().Unix())
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding chat response: %v", err)
		}
	})

	// Info endpoint
	http.HandleFunc("/api/ai/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"service":     "AI Shopping Assistant",
			"version":     "1.0.0",
			"status":      "placeholder",
			"description": "AI assistant service (OpenAI integration pending)",
			"endpoints": []string{
				"POST /api/ai/chat - Send chat message",
				"GET /api/ai/info - Service information",
				"GET /health - Health check",
			},
		}); err != nil {
			log.Printf("Error encoding info response: %v", err)
		}
	})

	log.Printf("Starting AI Assistant Service on port %s", port)
	log.Printf("Health check: http://localhost:%s/health", port)
	log.Printf("Note: This is a placeholder service. OpenAI integration pending.")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
