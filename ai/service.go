package ai

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service handles all AI assistant operations
type Service struct {
	db       *sql.DB
	aiClient AIClient
	logger   *log.Logger
}

// NewService creates a new AI service
func NewService(db *sql.DB, aiClient AIClient, logger *log.Logger) *Service {
	return &Service{
		db:       db,
		aiClient: aiClient,
		logger:   logger,
	}
}

// HandleChat processes a chat message and returns a response
func (s *Service) HandleChat(req ChatRequest) (*ChatResponse, error) {
	// Validate request
	if req.Message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	// Create or get conversation
	conv, err := s.getOrCreateConversation(req.SessionID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Get conversation history
	history, err := s.getConversationHistory(conv.ID, 10)
	if err != nil {
		s.logger.Printf("Warning: failed to get history: %v", err)
		history = []Message{}
	}

	// Save user message
	userMsg := Message{
		ConversationID: conv.ID,
		Role:           "user",
		Content:        req.Message,
		Model:          "",
		CreatedAt:      time.Now(),
	}
	if err := s.createMessage(&userMsg); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Add user message to history
	history = append(history, userMsg)

	// Get product context
	productContext, err := s.getProductContext()
	if err != nil {
		s.logger.Printf("Warning: failed to get product context: %v", err)
		productContext = "No products available at the moment."
	}

	// Get user preferences if available
	prefs, err := s.getUserPreferences(req.UserID, &req.SessionID)
	if err != nil {
		s.logger.Printf("Warning: failed to get user preferences: %v", err)
	}

	// Add preferences to context if available
	if prefs != nil {
		productContext = s.enrichContextWithPreferences(productContext, prefs)
	}

	// Generate AI response
	aiResp, err := s.aiClient.GenerateResponse(history, productContext)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Save assistant message
	assistantMsg := Message{
		ConversationID: conv.ID,
		Role:           "assistant",
		Content:        aiResp.Message,
		TokensUsed:     aiResp.TokensUsed,
		ResponseTimeMs: &aiResp.ResponseTimeMs,
		Model:          "gpt-3.5-turbo", // TODO: get from client
		CreatedAt:      time.Now(),
	}

	if len(aiResp.Products) > 0 {
		productsJSON, _ := json.Marshal(aiResp.Products)
		productsStr := string(productsJSON)
		assistantMsg.Metadata = &productsStr
	}

	if err := s.createMessage(&assistantMsg); err != nil {
		s.logger.Printf("Warning: failed to save assistant message: %v", err)
	}

	// Update conversation stats
	if err := s.updateConversationStats(conv.ID, aiResp.TokensUsed); err != nil {
		s.logger.Printf("Warning: failed to update conversation stats: %v", err)
	}

	// Update user preferences based on conversation
	if req.UserID != nil || req.SessionID != "" {
		go s.updatePreferencesFromMessage(req.UserID, &req.SessionID, req.Message)
	}

	// Add session ID to response
	aiResp.SessionID = req.SessionID

	return aiResp, nil
}

// getOrCreateConversation gets an existing conversation or creates a new one
func (s *Service) getOrCreateConversation(sessionID string, userID *int) (*Conversation, error) {
	// If no session ID, create one
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// Try to get existing conversation
	query := `SELECT id, session_id, user_id, started_at, ended_at, total_messages, 
	                 resulted_in_purchase, total_tokens_used, total_cost, user_agent, 
	                 ip_address, created_at, updated_at
	          FROM ai_conversations WHERE session_id = $1`
	
	var conv Conversation
	err := s.db.QueryRow(query, sessionID).Scan(
		&conv.ID, &conv.SessionID, &conv.UserID, &conv.StartedAt, &conv.EndedAt,
		&conv.TotalMessages, &conv.ResultedInPurchase, &conv.TotalTokensUsed,
		&conv.TotalCost, &conv.UserAgent, &conv.IPAddress, &conv.CreatedAt, &conv.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new conversation
		conv = Conversation{
			SessionID:     sessionID,
			UserID:        userID,
			StartedAt:     time.Now(),
			TotalMessages: 0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		query = `INSERT INTO ai_conversations (session_id, user_id, started_at, created_at, updated_at)
		         VALUES ($1, $2, $3, $4, $5) RETURNING id`
		err = s.db.QueryRow(query, conv.SessionID, conv.UserID, conv.StartedAt, conv.CreatedAt, conv.UpdatedAt).Scan(&conv.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create conversation: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to query conversation: %w", err)
	}

	return &conv, nil
}

// getConversationHistory retrieves recent messages
func (s *Service) getConversationHistory(conversationID int, limit int) ([]Message, error) {
	query := `SELECT id, conversation_id, role, content, tokens_used, response_time_ms,
	                 model, temperature, metadata, created_at
	          FROM ai_messages
	          WHERE conversation_id = $1
	          ORDER BY created_at DESC
	          LIMIT $2`

	rows, err := s.db.Query(query, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content,
			&msg.TokensUsed, &msg.ResponseTimeMs, &msg.Model, &msg.Temperature,
			&msg.Metadata, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// createMessage saves a message to the database
func (s *Service) createMessage(msg *Message) error {
	query := `INSERT INTO ai_messages (conversation_id, role, content, tokens_used,
	                                    response_time_ms, model, temperature, metadata, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	return s.db.QueryRow(query, msg.ConversationID, msg.Role, msg.Content, msg.TokensUsed,
		msg.ResponseTimeMs, msg.Model, msg.Temperature, msg.Metadata, msg.CreatedAt).Scan(&msg.ID)
}

// getProductContext retrieves product information for AI context
func (s *Service) getProductContext() (string, error) {
	query := `SELECT id, name, description, price, is_recurring 
	          FROM widgets 
	          WHERE inventory_level > 0 
	          ORDER BY price ASC 
	          LIMIT 20`

	rows, err := s.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var context strings.Builder
	context.WriteString("Available Products:\n\n")

	for rows.Next() {
		var id int
		var name, description string
		var price float64
		var isRecurring bool

		if err := rows.Scan(&id, &name, &description, &price, &isRecurring); err != nil {
			continue
		}

		productType := "one-time"
		if isRecurring {
			productType = "subscription"
		}

		context.WriteString(fmt.Sprintf("- %s ($%.2f, %s): %s\n", name, price/100.0, productType, description))
	}

	return context.String(), nil
}

// getUserPreferences retrieves user preferences
func (s *Service) getUserPreferences(userID *int, sessionID *string) (*UserPreferences, error) {
	query := `SELECT id, user_id, session_id, preferred_categories, budget_min, budget_max,
	                 interaction_count, last_products_viewed, last_products_purchased,
	                 conversation_style, preferred_language, created_at, updated_at
	          FROM ai_user_preferences
	          WHERE user_id = $1 OR session_id = $2
	          LIMIT 1`

	var prefs UserPreferences
	var categories, viewed, purchased sql.NullString

	err := s.db.QueryRow(query, userID, sessionID).Scan(
		&prefs.ID, &prefs.UserID, &prefs.SessionID, &categories, &prefs.BudgetMin,
		&prefs.BudgetMax, &prefs.InteractionCount, &viewed, &purchased,
		&prefs.ConversationStyle, &prefs.PreferredLanguage, &prefs.CreatedAt, &prefs.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse arrays (simplified)
	if categories.Valid {
		// TODO: Parse PostgreSQL array format
	}

	return &prefs, nil
}

// enrichContextWithPreferences adds user preferences to the product context
func (s *Service) enrichContextWithPreferences(context string, prefs *UserPreferences) string {
	var additions strings.Builder
	additions.WriteString("\n\nUser Preferences:\n")

	if len(prefs.PreferredCategories) > 0 {
		additions.WriteString(fmt.Sprintf("- Interested in: %s\n", strings.Join(prefs.PreferredCategories, ", ")))
	}

	if prefs.BudgetMin != nil && prefs.BudgetMax != nil {
		additions.WriteString(fmt.Sprintf("- Budget range: $%.2f - $%.2f\n", *prefs.BudgetMin, *prefs.BudgetMax))
	}

	if prefs.ConversationStyle != nil {
		additions.WriteString(fmt.Sprintf("- Communication style: %s\n", *prefs.ConversationStyle))
	}

	return context + additions.String()
}

// updateConversationStats updates conversation statistics
func (s *Service) updateConversationStats(convID int, tokensUsed int) error {
	query := `UPDATE ai_conversations 
	          SET total_messages = total_messages + 1,
	              total_tokens_used = total_tokens_used + $2,
	              total_cost = total_cost + $3,
	              updated_at = NOW()
	          WHERE id = $1`

	cost := float64(tokensUsed) / 1000.0 * 0.00175 // Approximate cost for GPT-3.5

	_, err := s.db.Exec(query, convID, tokensUsed, cost)
	return err
}

// updatePreferencesFromMessage learns from user messages
func (s *Service) updatePreferencesFromMessage(userID *int, sessionID *string, message string) {
	// Simple keyword extraction (in production, use NLP)
	message = strings.ToLower(message)

	var categories []string
	var budgetMin, budgetMax *float64

	// Detect categories
	if strings.Contains(message, "widget") {
		categories = append(categories, "widgets")
	}
	if strings.Contains(message, "subscription") || strings.Contains(message, "plan") {
		categories = append(categories, "subscriptions")
	}

	// Detect budget (very simple)
	if strings.Contains(message, "$30") || strings.Contains(message, "30") {
		val := 30.0
		budgetMax = &val
	}

	// TODO: Implement proper NLP-based preference extraction

	if len(categories) == 0 && budgetMin == nil && budgetMax == nil {
		return // Nothing to update
	}

	// Update preferences in database
	// (Implementation simplified for brevity)
	s.logger.Printf("Would update preferences: categories=%v, budget=%v-%v", categories, budgetMin, budgetMax)
}

// SubmitFeedback allows users to rate AI responses
func (s *Service) SubmitFeedback(feedback Feedback) error {
	query := `INSERT INTO ai_feedback (message_id, conversation_id, helpful, rating,
	                                    feedback_text, feedback_type, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.Exec(query, feedback.MessageID, feedback.ConversationID, feedback.Helpful,
		feedback.Rating, feedback.FeedbackText, feedback.FeedbackType, time.Now())

	return err
}

// GetConversationStats returns analytics about conversations
func (s *Service) GetConversationStats(days int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total conversations
	var totalConversations int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM ai_conversations WHERE started_at > NOW() - INTERVAL '1 day' * $1`, days).Scan(&totalConversations)
	if err != nil {
		return nil, err
	}
	stats["total_conversations"] = totalConversations

	// Total messages
	var totalMessages int
	err = s.db.QueryRow(`SELECT COALESCE(SUM(total_messages), 0) FROM ai_conversations WHERE started_at > NOW() - INTERVAL '1 day' * $1`, days).Scan(&totalMessages)
	if err != nil {
		return nil, err
	}
	stats["total_messages"] = totalMessages

	// Total cost
	var totalCost float64
	err = s.db.QueryRow(`SELECT COALESCE(SUM(total_cost), 0) FROM ai_conversations WHERE started_at > NOW() - INTERVAL '1 day' * $1`, days).Scan(&totalCost)
	if err != nil {
		return nil, err
	}
	stats["total_cost"] = totalCost

	// Conversion rate
	var purchases int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM ai_conversations WHERE resulted_in_purchase = true AND started_at > NOW() - INTERVAL '1 day' * $1`, days).Scan(&purchases)
	if err != nil {
		return nil, err
	}
	stats["purchases"] = purchases

	if totalConversations > 0 {
		stats["conversion_rate"] = float64(purchases) / float64(totalConversations) * 100.0
	} else {
		stats["conversion_rate"] = 0.0
	}

	return stats, nil
}

