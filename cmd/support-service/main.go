package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"usual_store/internal/driver"
	"usual_store/internal/support"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	hub      *support.Hub
	db       *sql.DB
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins in development
			return true
		},
	}
)

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Support Chat Service...")

	// Load configuration from environment variables
	port := getEnv("SUPPORT_SERVICE_PORT", "5001")
	databaseDSN := getEnv("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")

	// Connect to database
	var err error
	db, err = driver.OpenDB(databaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Initialize WebSocket hub
	hub = support.NewHub()
	go hub.Run()
	log.Println("WebSocket hub started")

	// Setup router
	router := chi.NewRouter()

	// Middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	router.Post("/api/support/ticket", handleCreateTicket)
	router.Get("/api/support/tickets", handleGetOpenTickets)
	router.Get("/api/support/ticket/{sessionID}", handleGetTicket)
	router.Get("/api/support/ticket/{sessionID}/messages", handleGetMessages)
	router.Post("/api/support/ticket/{ticketID}/assign", handleAssignTicket)
	router.Get("/ws/support/user/{sessionID}", handleUserWebSocket)
	router.Get("/ws/support/supporter/{sessionID}", handleSupporterWebSocket)
	router.Get("/api/support/health", handleHealthCheck)

	// Start server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down Support Chat Service...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Support Chat Service listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Support Chat Service stopped")
}

// handleCreateTicket creates a new support ticket
func handleCreateTicket(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    *int   `json:"user_id,omitempty"`
		Subject   string `json:"subject"`
		UserEmail string `json:"user_email"`
		UserName  string `json:"user_name"`
		Priority  string `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate session ID
	sessionID := uuid.New().String()

	// Set default priority if not provided
	if req.Priority == "" {
		req.Priority = "medium"
	}

	// Create ticket
	ticket := &support.Ticket{
		UserID:    req.UserID,
		Subject:   req.Subject,
		Status:    "open",
		Priority:  req.Priority,
		SessionID: sessionID,
		UserEmail: req.UserEmail,
		UserName:  req.UserName,
	}

	if err := support.CreateTicket(db, ticket); err != nil {
		log.Printf("Error creating ticket: %v", err)
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticket); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// handleGetOpenTickets retrieves all open tickets
func handleGetOpenTickets(w http.ResponseWriter, r *http.Request) {
	tickets, err := support.GetAllOpenTickets(db)
	if err != nil {
		log.Printf("Error getting tickets: %v", err)
		http.Error(w, "Failed to get tickets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tickets); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// handleGetTicket retrieves a ticket by session ID
func handleGetTicket(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	ticket, err := support.GetTicketBySessionID(db, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Ticket not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting ticket: %v", err)
			http.Error(w, "Failed to get ticket", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticket); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// handleGetMessages retrieves all messages for a ticket
func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	ticket, err := support.GetTicketBySessionID(db, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Ticket not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting ticket: %v", err)
			http.Error(w, "Failed to get ticket", http.StatusInternalServerError)
		}
		return
	}

	messages, err := support.GetMessagesByTicketID(db, ticket.ID)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// handleAssignTicket assigns a ticket to a supporter
func handleAssignTicket(w http.ResponseWriter, r *http.Request) {
	ticketIDStr := chi.URLParam(r, "ticketID")
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	var req struct {
		SupporterID int `json:"supporter_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := support.UpdateTicketStatus(db, ticketID, "assigned", &req.SupporterID); err != nil {
		log.Printf("Error assigning ticket: %v", err)
		http.Error(w, "Failed to assign ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// handleUserWebSocket handles WebSocket connections for users
func handleUserWebSocket(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	// Get user info from query parameters
	userName := r.URL.Query().Get("name")
	userEmail := r.URL.Query().Get("email")

	if userName == "" {
		userName = "Guest"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &support.Client{
		Hub:        hub,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		SessionID:  sessionID,
		ClientType: "user",
		UserID:     nil,
		UserName:   userName,
		DB:         db,
	}

	client.Hub.Register <- client

	// Send welcome message
	welcomeMsg := &support.Message{
		SenderType: "system",
		SenderName: "System",
		Message:    fmt.Sprintf("Welcome %s! A support agent will be with you shortly.", userName),
		CreatedAt:  time.Now(),
	}
	hub.BroadcastToSession(sessionID, welcomeMsg, "system")

	log.Printf("User connected: SessionID=%s, Name=%s, Email=%s", sessionID, userName, userEmail)

	// Start client pumps
	go client.WritePump()
	go client.ReadPump()
}

// handleSupporterWebSocket handles WebSocket connections for supporters
func handleSupporterWebSocket(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	// Get supporter info from query parameters (in production, this should come from JWT)
	supporterIDStr := r.URL.Query().Get("supporter_id")
	supporterName := r.URL.Query().Get("name")

	if supporterName == "" {
		supporterName = "Support Agent"
	}

	var supporterID *int
	if supporterIDStr != "" {
		if id, err := strconv.Atoi(supporterIDStr); err == nil {
			supporterID = &id
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &support.Client{
		Hub:        hub,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		SessionID:  sessionID,
		ClientType: "supporter",
		UserID:     supporterID,
		UserName:   supporterName,
		DB:         db,
	}

	client.Hub.Register <- client

	// Send notification to user
	joinMsg := &support.Message{
		SenderType: "system",
		SenderName: "System",
		Message:    fmt.Sprintf("%s has joined the conversation.", supporterName),
		CreatedAt:  time.Now(),
	}
	hub.BroadcastToSession(sessionID, joinMsg, "supporter_joined")

	log.Printf("Supporter connected: SessionID=%s, Name=%s", sessionID, supporterName)

	// Start client pumps
	go client.WritePump()
	go client.ReadPump()
}

// handleHealthCheck returns the health status of the service
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "healthy",
		"active_sessions": len(hub.GetAllSessions()),
		"timestamp":       time.Now().Unix(),
	}); err != nil {
		log.Printf("Error encoding health check response: %v", err)
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
