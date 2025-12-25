package ai

import (
	"time"
)

// Conversation represents an AI chat session
type Conversation struct {
	ID                 int        `json:"id"`
	SessionID          string     `json:"session_id"`
	UserID             *int       `json:"user_id,omitempty"`
	StartedAt          time.Time  `json:"started_at"`
	EndedAt            *time.Time `json:"ended_at,omitempty"`
	TotalMessages      int        `json:"total_messages"`
	ResultedInPurchase bool       `json:"resulted_in_purchase"`
	TotalTokensUsed    int        `json:"total_tokens_used"`
	TotalCost          float64    `json:"total_cost"`
	UserAgent          string     `json:"user_agent,omitempty"`
	IPAddress          string     `json:"ip_address,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// Message represents a single message in a conversation
type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	Role           string    `json:"role"` // "user", "assistant", "system"
	Content        string    `json:"content"`
	TokensUsed     int       `json:"tokens_used"`
	ResponseTimeMs *int      `json:"response_time_ms,omitempty"`
	Model          string    `json:"model"`
	Temperature    float64   `json:"temperature,omitempty"`
	Metadata       *string   `json:"metadata,omitempty"` // JSON string
	CreatedAt      time.Time `json:"created_at"`
}

// UserPreferences stores learned preferences from interactions
type UserPreferences struct {
	ID                    int       `json:"id"`
	UserID                *int      `json:"user_id,omitempty"`
	SessionID             *string   `json:"session_id,omitempty"`
	PreferredCategories   []string  `json:"preferred_categories,omitempty"`
	BudgetMin             *float64  `json:"budget_min,omitempty"`
	BudgetMax             *float64  `json:"budget_max,omitempty"`
	InteractionCount      int       `json:"interaction_count"`
	LastProductsViewed    []int     `json:"last_products_viewed,omitempty"`
	LastProductsPurchased []int     `json:"last_products_purchased,omitempty"`
	ConversationStyle     *string   `json:"conversation_style,omitempty"`
	PreferredLanguage     string    `json:"preferred_language"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Feedback represents user feedback on AI responses
type Feedback struct {
	ID             int       `json:"id"`
	MessageID      int       `json:"message_id"`
	ConversationID int       `json:"conversation_id"`
	Helpful        *bool     `json:"helpful,omitempty"`
	Rating         *int      `json:"rating,omitempty"` // 1-5
	FeedbackText   *string   `json:"feedback_text,omitempty"`
	FeedbackType   *string   `json:"feedback_type,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ProductCache stores product information optimized for AI
type ProductCache struct {
	ID              int        `json:"id"`
	ProductID       int        `json:"product_id"`
	DescriptionText string     `json:"description_text"`
	SearchKeywords  []string   `json:"search_keywords"`
	Category        *string    `json:"category,omitempty"`
	PriceTier       *string    `json:"price_tier,omitempty"` // "budget", "mid", "premium"
	PopularityScore int        `json:"popularity_score"`
	LastMentionedAt *time.Time `json:"last_mentioned_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ChatRequest represents an incoming chat message from the user
type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
	UserID    *int   `json:"user_id,omitempty"`
}

// ChatResponse represents the AI assistant's response
type ChatResponse struct {
	SessionID      string                 `json:"session_id"`
	Message        string                 `json:"message"`
	TokensUsed     int                    `json:"tokens_used"`
	ResponseTimeMs int                    `json:"response_time_ms"`
	Suggestions    []string               `json:"suggestions,omitempty"`
	Products       []RecommendedProduct   `json:"products,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// RecommendedProduct represents a product suggested by the AI
type RecommendedProduct struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image,omitempty"`
	Reason      string  `json:"reason"` // Why this product was recommended
}

// ConversationHistory returns the last N messages for context
type ConversationHistory struct {
	Messages []Message `json:"messages"`
}

// DB interface for database operations
type DB interface {
	CreateConversation(conv *Conversation) error
	GetConversation(sessionID string) (*Conversation, error)
	UpdateConversation(conv *Conversation) error

	CreateMessage(msg *Message) error
	GetMessages(conversationID int, limit int) ([]Message, error)

	GetUserPreferences(userID *int, sessionID *string) (*UserPreferences, error)
	UpsertUserPreferences(prefs *UserPreferences) error

	CreateFeedback(feedback *Feedback) error

	GetProductCache(productID int) (*ProductCache, error)
	GetAllProductsCache() ([]ProductCache, error)
	UpdateProductPopularity(productID int) error
}

// AIClient interface for interacting with AI services (OpenAI, etc.)
type AIClient interface {
	GenerateResponse(messages []Message, context string) (*ChatResponse, error)
	GetEmbedding(text string) ([]float64, error)
}
