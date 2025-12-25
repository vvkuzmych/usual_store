package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"usual_store/internal/ai"

	_ "github.com/lib/pq"
)

func main() {
	// Get configuration from environment
	dbDSN := os.Getenv("DATABASE_DSN")
	if dbDSN == "" {
		log.Fatal("DATABASE_DSN environment variable is required")
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Println("⚠️  WARNING: OPENAI_API_KEY not set - AI chat will not work!")
		log.Println("   Set OPENAI_API_KEY in .env file to enable AI features")
		openaiKey = "dummy-key-for-development" // Allow service to start
	}

	// Connect to database
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Create OpenAI client
	openaiClient := ai.NewOpenAIClient(
		openaiKey,
		"gpt-3.5-turbo", // or "gpt-4" for better quality
		0.7,             // temperature
	)

	// Create AI service
	logger := log.New(os.Stdout, "[AI-Assistant] ", log.LstdFlags)
	aiService := ai.NewService(db, openaiClient, logger)

	// Create HTTP handler
	aiHandler := ai.NewHandler(aiService, logger)

	// Set up routes
	mux := http.NewServeMux()

	// Register AI routes
	mux.HandleFunc("/api/ai/chat", aiHandler.EnableCORS(aiHandler.HandleChatRequest))
	mux.HandleFunc("/api/ai/feedback", aiHandler.EnableCORS(aiHandler.HandleFeedback))
	mux.HandleFunc("/api/ai/stats", aiHandler.EnableCORS(aiHandler.HandleStats))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("AI Assistant API server starting on port %s", port)
	log.Printf("Endpoints:")
	log.Printf("  POST   /api/ai/chat       - Send chat message")
	log.Printf("  POST   /api/ai/feedback   - Submit feedback")
	log.Printf("  GET    /api/ai/stats      - Get statistics")
	log.Printf("  GET    /health            - Health check")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
