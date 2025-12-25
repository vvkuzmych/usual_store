package ai

import (
	"encoding/json"
	"testing"
	"time"
)

// Test Conversation struct
func TestConversation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		conv Conversation
	}{
		{
			name: "new conversation",
			conv: Conversation{
				ID:                 1,
				SessionID:          "sess-123",
				UserID:             intPtr(42),
				StartedAt:          now,
				EndedAt:            nil,
				TotalMessages:      0,
				ResultedInPurchase: false,
				TotalTokensUsed:    0,
				TotalCost:          0.0,
				UserAgent:          "Mozilla/5.0",
				IPAddress:          "192.168.1.1",
				CreatedAt:          now,
				UpdatedAt:          now,
			},
		},
		{
			name: "completed conversation",
			conv: Conversation{
				ID:                 2,
				SessionID:          "sess-456",
				UserID:             intPtr(99),
				StartedAt:          now.Add(-1 * time.Hour),
				EndedAt:            &now,
				TotalMessages:      15,
				ResultedInPurchase: true,
				TotalTokensUsed:    500,
				TotalCost:          0.875,
				UserAgent:          "Chrome/120.0",
				IPAddress:          "10.0.0.1",
				CreatedAt:          now.Add(-1 * time.Hour),
				UpdatedAt:          now,
			},
		},
		{
			name: "anonymous conversation",
			conv: Conversation{
				ID:                 3,
				SessionID:          "sess-anon",
				UserID:             nil,
				StartedAt:          now,
				TotalMessages:      3,
				ResultedInPurchase: false,
				TotalTokensUsed:    100,
				TotalCost:          0.175,
				CreatedAt:          now,
				UpdatedAt:          now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.conv)
			if err != nil {
				t.Fatalf("failed to marshal Conversation: %v", err)
			}

			// Test JSON unmarshaling
			var decoded Conversation
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal Conversation: %v", err)
			}

			// Verify key fields
			if decoded.ID != tt.conv.ID {
				t.Errorf("ID mismatch: got %d, want %d", decoded.ID, tt.conv.ID)
			}
			if decoded.SessionID != tt.conv.SessionID {
				t.Errorf("SessionID mismatch: got %s, want %s", decoded.SessionID, tt.conv.SessionID)
			}
			if decoded.TotalMessages != tt.conv.TotalMessages {
				t.Errorf("TotalMessages mismatch: got %d, want %d", decoded.TotalMessages, tt.conv.TotalMessages)
			}
		})
	}
}

// Test Message struct
func TestMessage(t *testing.T) {
	now := time.Now()
	responseTime := 250

	tests := []struct {
		name string
		msg  Message
	}{
		{
			name: "user message",
			msg: Message{
				ID:             1,
				ConversationID: 100,
				Role:           "user",
				Content:        "Hello, I need help finding a product",
				TokensUsed:     10,
				ResponseTimeMs: nil,
				Model:          "",
				Temperature:    0.0,
				Metadata:       nil,
				CreatedAt:      now,
			},
		},
		{
			name: "assistant message",
			msg: Message{
				ID:             2,
				ConversationID: 100,
				Role:           "assistant",
				Content:        "I'd be happy to help! What are you looking for?",
				TokensUsed:     15,
				ResponseTimeMs: &responseTime,
				Model:          "gpt-3.5-turbo",
				Temperature:    0.7,
				Metadata:       stringPtr(`{"products": [1,2,3]}`),
				CreatedAt:      now,
			},
		},
		{
			name: "system message",
			msg: Message{
				ID:             3,
				ConversationID: 100,
				Role:           "system",
				Content:        "You are a helpful shopping assistant",
				TokensUsed:     20,
				ResponseTimeMs: nil,
				Model:          "gpt-3.5-turbo",
				Temperature:    0.0,
				Metadata:       nil,
				CreatedAt:      now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.msg)
			if err != nil {
				t.Fatalf("failed to marshal Message: %v", err)
			}

			// Test JSON unmarshaling
			var decoded Message
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal Message: %v", err)
			}

			// Verify key fields
			if decoded.ID != tt.msg.ID {
				t.Errorf("ID mismatch: got %d, want %d", decoded.ID, tt.msg.ID)
			}
			if decoded.Role != tt.msg.Role {
				t.Errorf("Role mismatch: got %s, want %s", decoded.Role, tt.msg.Role)
			}
			if decoded.Content != tt.msg.Content {
				t.Errorf("Content mismatch: got %s, want %s", decoded.Content, tt.msg.Content)
			}
		})
	}
}

// Test UserPreferences struct
func TestUserPreferences(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		prefs UserPreferences
	}{
		{
			name: "complete preferences",
			prefs: UserPreferences{
				ID:                    1,
				UserID:                intPtr(42),
				SessionID:             stringPtr("sess-123"),
				PreferredCategories:   []string{"electronics", "gadgets"},
				BudgetMin:             float64Ptr(10.00),
				BudgetMax:             float64Ptr(100.00),
				InteractionCount:      5,
				LastProductsViewed:    []int{1, 2, 3},
				LastProductsPurchased: []int{1},
				ConversationStyle:     stringPtr("casual"),
				PreferredLanguage:     "en",
				CreatedAt:             now,
				UpdatedAt:             now,
			},
		},
		{
			name: "minimal preferences",
			prefs: UserPreferences{
				ID:                1,
				InteractionCount:  1,
				PreferredLanguage: "en",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		{
			name: "session-only preferences",
			prefs: UserPreferences{
				ID:                1,
				SessionID:         stringPtr("sess-anon"),
				InteractionCount:  3,
				PreferredLanguage: "en",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.prefs)
			if err != nil {
				t.Fatalf("failed to marshal UserPreferences: %v", err)
			}

			// Test JSON unmarshaling
			var decoded UserPreferences
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal UserPreferences: %v", err)
			}

			// Verify key fields
			if decoded.ID != tt.prefs.ID {
				t.Errorf("ID mismatch: got %d, want %d", decoded.ID, tt.prefs.ID)
			}
			if decoded.InteractionCount != tt.prefs.InteractionCount {
				t.Errorf("InteractionCount mismatch: got %d, want %d", decoded.InteractionCount, tt.prefs.InteractionCount)
			}
		})
	}
}

// Test Feedback struct
func TestFeedback(t *testing.T) {
	now := time.Now()
	helpful := true
	rating := 5

	tests := []struct {
		name     string
		feedback Feedback
	}{
		{
			name: "positive feedback",
			feedback: Feedback{
				ID:             1,
				MessageID:      10,
				ConversationID: 100,
				Helpful:        &helpful,
				Rating:         &rating,
				FeedbackText:   stringPtr("Very helpful response!"),
				FeedbackType:   stringPtr("praise"),
				CreatedAt:      now,
			},
		},
		{
			name: "negative feedback",
			feedback: Feedback{
				ID:             2,
				MessageID:      11,
				ConversationID: 100,
				Helpful:        boolPtr(false),
				Rating:         intPtr(2),
				FeedbackText:   stringPtr("Not quite what I was looking for"),
				FeedbackType:   stringPtr("critique"),
				CreatedAt:      now,
			},
		},
		{
			name: "minimal feedback",
			feedback: Feedback{
				ID:             3,
				MessageID:      12,
				ConversationID: 100,
				CreatedAt:      now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.feedback)
			if err != nil {
				t.Fatalf("failed to marshal Feedback: %v", err)
			}

			// Test JSON unmarshaling
			var decoded Feedback
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal Feedback: %v", err)
			}

			// Verify key fields
			if decoded.ID != tt.feedback.ID {
				t.Errorf("ID mismatch: got %d, want %d", decoded.ID, tt.feedback.ID)
			}
			if decoded.MessageID != tt.feedback.MessageID {
				t.Errorf("MessageID mismatch: got %d, want %d", decoded.MessageID, tt.feedback.MessageID)
			}
		})
	}
}

// Test ProductCache struct
func TestProductCache(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		cache ProductCache
	}{
		{
			name: "popular product",
			cache: ProductCache{
				ID:              1,
				ProductID:       42,
				DescriptionText: "Amazing wireless widget with Bluetooth connectivity",
				SearchKeywords:  []string{"wireless", "bluetooth", "widget", "gadget"},
				Category:        stringPtr("electronics"),
				PriceTier:       stringPtr("mid"),
				PopularityScore: 85,
				LastMentionedAt: &now,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},
		{
			name: "budget product",
			cache: ProductCache{
				ID:              2,
				ProductID:       10,
				DescriptionText: "Basic widget for everyday use",
				SearchKeywords:  []string{"basic", "widget", "affordable"},
				Category:        stringPtr("basics"),
				PriceTier:       stringPtr("budget"),
				PopularityScore: 45,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},
		{
			name: "premium product",
			cache: ProductCache{
				ID:              3,
				ProductID:       99,
				DescriptionText: "Premium subscription plan with all features",
				SearchKeywords:  []string{"premium", "subscription", "unlimited", "pro"},
				Category:        stringPtr("subscriptions"),
				PriceTier:       stringPtr("premium"),
				PopularityScore: 95,
				LastMentionedAt: &now,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.cache)
			if err != nil {
				t.Fatalf("failed to marshal ProductCache: %v", err)
			}

			// Test JSON unmarshaling
			var decoded ProductCache
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("failed to unmarshal ProductCache: %v", err)
			}

			// Verify key fields
			if decoded.ID != tt.cache.ID {
				t.Errorf("ID mismatch: got %d, want %d", decoded.ID, tt.cache.ID)
			}
			if decoded.ProductID != tt.cache.ProductID {
				t.Errorf("ProductID mismatch: got %d, want %d", decoded.ProductID, tt.cache.ProductID)
			}
			if decoded.PopularityScore != tt.cache.PopularityScore {
				t.Errorf("PopularityScore mismatch: got %d, want %d", decoded.PopularityScore, tt.cache.PopularityScore)
			}
		})
	}
}

// Test ChatRequest and ChatResponse
func TestChatRequestResponse(t *testing.T) {
	tests := []struct {
		name     string
		request  ChatRequest
		response ChatResponse
	}{
		{
			name: "simple chat exchange",
			request: ChatRequest{
				SessionID: "sess-123",
				Message:   "Show me widgets under $50",
				UserID:    intPtr(1),
			},
			response: ChatResponse{
				SessionID:      "sess-123",
				Message:        "Here are widgets under $50",
				TokensUsed:     25,
				ResponseTimeMs: 150,
				Suggestions:    []string{"See all products", "Filter by category"},
				Products: []RecommendedProduct{
					{
						ID:          1,
						Name:        "Basic Widget",
						Description: "Affordable widget",
						Price:       29.99,
						Image:       "widget.jpg",
						Reason:      "Under your budget",
					},
				},
				Metadata: map[string]interface{}{
					"query_type": "price_filter",
					"max_price":  50.00,
				},
			},
		},
		{
			name: "anonymous chat",
			request: ChatRequest{
				SessionID: "sess-anon",
				Message:   "What products do you have?",
			},
			response: ChatResponse{
				SessionID:      "sess-anon",
				Message:        "We have widgets and subscription plans",
				TokensUsed:     20,
				ResponseTimeMs: 120,
				Suggestions:    []string{"Browse widgets", "View plans"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test ChatRequest JSON
			reqData, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("failed to marshal ChatRequest: %v", err)
			}

			var decodedReq ChatRequest
			err = json.Unmarshal(reqData, &decodedReq)
			if err != nil {
				t.Fatalf("failed to unmarshal ChatRequest: %v", err)
			}

			if decodedReq.SessionID != tt.request.SessionID {
				t.Errorf("Request SessionID mismatch: got %s, want %s", decodedReq.SessionID, tt.request.SessionID)
			}

			// Test ChatResponse JSON
			respData, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("failed to marshal ChatResponse: %v", err)
			}

			var decodedResp ChatResponse
			err = json.Unmarshal(respData, &decodedResp)
			if err != nil {
				t.Fatalf("failed to unmarshal ChatResponse: %v", err)
			}

			if decodedResp.SessionID != tt.response.SessionID {
				t.Errorf("Response SessionID mismatch: got %s, want %s", decodedResp.SessionID, tt.response.SessionID)
			}
		})
	}
}

// Benchmark JSON marshaling
func BenchmarkConversationMarshal(b *testing.B) {
	conv := Conversation{
		ID:                 1,
		SessionID:          "sess-123",
		UserID:             intPtr(42),
		StartedAt:          time.Now(),
		TotalMessages:      10,
		ResultedInPurchase: true,
		TotalTokensUsed:    500,
		TotalCost:          0.875,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(conv)
	}
}

func BenchmarkMessageMarshal(b *testing.B) {
	msg := Message{
		ID:             1,
		ConversationID: 100,
		Role:           "assistant",
		Content:        "Here are some recommended products for you",
		TokensUsed:     25,
		Model:          "gpt-3.5-turbo",
		CreatedAt:      time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(msg)
	}
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}
