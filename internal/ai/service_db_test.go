package ai

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetOrCreateConversation(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		userID      *int
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		checkResult func(*testing.T, *Conversation)
	}{
		{
			name:      "existing conversation found",
			sessionID: "sess-123",
			userID:    intPtr(42),
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "session_id", "user_id", "started_at", "ended_at",
					"total_messages", "resulted_in_purchase", "total_tokens_used",
					"total_cost", "user_agent", "ip_address", "created_at", "updated_at",
				}).AddRow(
					1, "sess-123", 42, time.Now(), nil,
					5, false, 100, 0.175, "Mozilla/5.0", "192.168.1.1",
					time.Now(), time.Now(),
				)
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-123").
					WillReturnRows(rows)
			},
			expectError: false,
			checkResult: func(t *testing.T, conv *Conversation) {
				if conv.ID != 1 {
					t.Errorf("Expected ID 1, got %d", conv.ID)
				}
				if conv.SessionID != "sess-123" {
					t.Errorf("Expected session_id sess-123, got %s", conv.SessionID)
				}
				if *conv.UserID != 42 {
					t.Errorf("Expected user_id 42, got %d", *conv.UserID)
				}
			},
		},
		{
			name:      "create new conversation",
			sessionID: "sess-new",
			userID:    intPtr(99),
			setupMock: func(mock sqlmock.Sqlmock) {
				// First query returns no rows
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-new").
					WillReturnError(sql.ErrNoRows)

				// Insert new conversation
				mock.ExpectQuery("INSERT INTO ai_conversations").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
			},
			expectError: false,
			checkResult: func(t *testing.T, conv *Conversation) {
				if conv.ID != 10 {
					t.Errorf("Expected ID 10, got %d", conv.ID)
				}
				if conv.SessionID != "sess-new" {
					t.Errorf("Expected session_id sess-new, got %s", conv.SessionID)
				}
			},
		},
		{
			name:      "empty session ID generates UUID",
			sessionID: "",
			userID:    nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				// Will try to find conversation with generated UUID (won't find it)
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)

				// Insert new conversation
				mock.ExpectQuery("INSERT INTO ai_conversations").
					WithArgs(sqlmock.AnyArg(), nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(20))
			},
			expectError: false,
			checkResult: func(t *testing.T, conv *Conversation) {
				if conv.ID != 20 {
					t.Errorf("Expected ID 20, got %d", conv.ID)
				}
				if conv.SessionID == "" {
					t.Error("Expected non-empty generated session ID")
				}
			},
		},
		{
			name:      "database query error",
			sessionID: "sess-error",
			userID:    nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-error").
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{
				db:       db,
				aiClient: nil,
				logger:   logger,
			}

			conv, err := service.getOrCreateConversation(tt.sessionID, tt.userID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkResult != nil && conv != nil {
					tt.checkResult(t, conv)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetConversationHistory(t *testing.T) {
	tests := []struct {
		name           string
		conversationID int
		limit          int
		setupMock      func(sqlmock.Sqlmock)
		expectedCount  int
		expectError    bool
	}{
		{
			name:           "retrieve multiple messages",
			conversationID: 100,
			limit:          10,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "conversation_id", "role", "content", "tokens_used",
					"response_time_ms", "model", "temperature", "metadata", "created_at",
				}).
					AddRow(3, 100, "assistant", "Hello!", 10, 100, "gpt-3.5", 0.7, nil, time.Now()).
					AddRow(2, 100, "user", "Hi there", 5, nil, "", 0.0, nil, time.Now()).
					AddRow(1, 100, "system", "You are helpful", 20, nil, "gpt-3.5", 0.0, nil, time.Now())

				mock.ExpectQuery("SELECT (.+) FROM ai_messages WHERE conversation_id").
					WithArgs(100, 10).
					WillReturnRows(rows)
			},
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:           "no messages found",
			conversationID: 200,
			limit:          5,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "conversation_id", "role", "content", "tokens_used",
					"response_time_ms", "model", "temperature", "metadata", "created_at",
				})

				mock.ExpectQuery("SELECT (.+) FROM ai_messages WHERE conversation_id").
					WithArgs(200, 5).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:           "database error",
			conversationID: 300,
			limit:          10,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM ai_messages WHERE conversation_id").
					WithArgs(300, 10).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{
				db:       db,
				aiClient: nil,
				logger:   logger,
			}

			messages, err := service.getConversationHistory(tt.conversationID, tt.limit)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(messages) != tt.expectedCount {
					t.Errorf("Expected %d messages, got %d", tt.expectedCount, len(messages))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestCreateMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     *Message
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "create user message",
			message: &Message{
				ConversationID: 100,
				Role:           "user",
				Content:        "Hello AI",
				TokensUsed:     5,
				Model:          "",
				CreatedAt:      time.Now(),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(100, "user", "Hello AI", 5, nil, "", sqlmock.AnyArg(), nil, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
			},
			expectError: false,
		},
		{
			name: "create assistant message with metadata",
			message: &Message{
				ConversationID: 100,
				Role:           "assistant",
				Content:        "I can help you!",
				TokensUsed:     15,
				ResponseTimeMs: intPtr(250),
				Model:          "gpt-3.5-turbo",
				Temperature:    0.7,
				Metadata:       stringPtr(`{"products": [1,2,3]}`),
				CreatedAt:      time.Now(),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(43))
			},
			expectError: false,
		},
		{
			name: "database error",
			message: &Message{
				ConversationID: 100,
				Role:           "user",
				Content:        "Test",
				TokensUsed:     5,
				CreatedAt:      time.Now(),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{
				db:       db,
				aiClient: nil,
				logger:   logger,
			}

			err = service.createMessage(tt.message)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !tt.expectError && tt.message.ID == 0 {
					t.Error("Expected message ID to be set")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetProductContext(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func(sqlmock.Sqlmock)
		expectedContent []string
		expectError     bool
	}{
		{
			name: "retrieve multiple products",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "is_recurring"}).
					AddRow(1, "Widget", "A great widget", 1000, false).
					AddRow(2, "Premium Plan", "Monthly subscription", 3000, true)

				mock.ExpectQuery("SELECT (.+) FROM widgets WHERE inventory_level").
					WillReturnRows(rows)
			},
			expectedContent: []string{
				"Available Products:",
				"Widget",
				"$10.00",
				"one-time",
				"Premium Plan",
				"$30.00",
				"subscription",
			},
			expectError: false,
		},
		{
			name: "no products available",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "is_recurring"})

				mock.ExpectQuery("SELECT (.+) FROM widgets WHERE inventory_level").
					WillReturnRows(rows)
			},
			expectedContent: []string{
				"Available Products:",
			},
			expectError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM widgets WHERE inventory_level").
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{
				db:       db,
				aiClient: nil,
				logger:   logger,
			}

			context, err := service.getProductContext()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				for _, expected := range tt.expectedContent {
					if !contains(context, expected) {
						t.Errorf("Expected context to contain %q, but it doesn't", expected)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetUserPreferences(t *testing.T) {
	tests := []struct {
		name        string
		userID      *int
		sessionID   *string
		setupMock   func(sqlmock.Sqlmock)
		expectFound bool
		expectError bool
	}{
		{
			name:      "preferences found for user",
			userID:    intPtr(42),
			sessionID: stringPtr("sess-123"),
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "session_id", "preferred_categories", "budget_min",
					"budget_max", "interaction_count", "last_products_viewed",
					"last_products_purchased", "conversation_style", "preferred_language",
					"created_at", "updated_at",
				}).AddRow(
					1, 42, "sess-123", sql.NullString{}, float64Ptr(10.0),
					float64Ptr(100.0), 5, sql.NullString{}, sql.NullString{},
					stringPtr("casual"), "en", time.Now(), time.Now(),
				)

				mock.ExpectQuery("SELECT (.+) FROM ai_user_preferences WHERE user_id").
					WithArgs(intPtr(42), stringPtr("sess-123")).
					WillReturnRows(rows)
			},
			expectFound: true,
			expectError: false,
		},
		{
			name:      "no preferences found",
			userID:    intPtr(99),
			sessionID: stringPtr("sess-new"),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM ai_user_preferences WHERE user_id").
					WithArgs(intPtr(99), stringPtr("sess-new")).
					WillReturnError(sql.ErrNoRows)
			},
			expectFound: false,
			expectError: false,
		},
		{
			name:      "database error",
			userID:    intPtr(42),
			sessionID: stringPtr("sess-error"),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM ai_user_preferences WHERE user_id").
					WithArgs(intPtr(42), stringPtr("sess-error")).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{
				db:       db,
				aiClient: nil,
				logger:   logger,
			}

			prefs, err := service.getUserPreferences(tt.userID, tt.sessionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.expectFound && prefs == nil {
					t.Error("Expected preferences to be found, got nil")
				}
				if !tt.expectFound && prefs != nil {
					t.Error("Expected no preferences, got non-nil")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestUpdateConversationStats(t *testing.T) {
	tests := []struct {
		name        string
		convID      int
		tokensUsed  int
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:       "successful update",
			convID:     100,
			tokensUsed: 50,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE ai_conversations SET").
					WithArgs(100, 50, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name:       "database error",
			convID:     200,
			tokensUsed: 25,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE ai_conversations SET").
					WithArgs(200, 25, sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{
				db:       db,
				aiClient: nil,
				logger:   logger,
			}

			err = service.updateConversationStats(tt.convID, tt.tokensUsed)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestSubmitFeedback(t *testing.T) {
	tests := []struct {
		name        string
		feedback    Feedback
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "successful feedback submission",
			feedback: Feedback{
				MessageID:      10,
				ConversationID: 100,
				Helpful:        boolPtr(true),
				Rating:         intPtr(5),
				FeedbackText:   stringPtr("Very helpful!"),
				FeedbackType:   stringPtr("praise"),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO ai_feedback").
					WithArgs(10, 100, boolPtr(true), intPtr(5), stringPtr("Very helpful!"), stringPtr("praise"), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name: "database error",
			feedback: Feedback{
				MessageID:      20,
				ConversationID: 200,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO ai_feedback").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{
				db:       db,
				aiClient: nil,
				logger:   logger,
			}

			err = service.SubmitFeedback(tt.feedback)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetConversationStats(t *testing.T) {
	tests := []struct {
		name        string
		days        int
		setupMock   func(sqlmock.Sqlmock)
		checkResult func(*testing.T, map[string]interface{})
		expectError bool
	}{
		{
			name: "retrieve stats for 7 days",
			days: 7,
			setupMock: func(mock sqlmock.Sqlmock) {
				// Total conversations
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ai_conversations WHERE started_at").
					WithArgs(7).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))

				// Total messages
				mock.ExpectQuery("SELECT COALESCE\\(SUM\\(total_messages\\)").
					WithArgs(7).
					WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(500))

				// Total cost
				mock.ExpectQuery("SELECT COALESCE\\(SUM\\(total_cost\\)").
					WithArgs(7).
					WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(1.25))

				// Purchases
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ai_conversations WHERE resulted_in_purchase").
					WithArgs(7).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
			},
			checkResult: func(t *testing.T, stats map[string]interface{}) {
				if stats["total_conversations"] != 100 {
					t.Errorf("Expected 100 conversations, got %v", stats["total_conversations"])
				}
				if stats["total_messages"] != 500 {
					t.Errorf("Expected 500 messages, got %v", stats["total_messages"])
				}
				if stats["total_cost"] != 1.25 {
					t.Errorf("Expected cost 1.25, got %v", stats["total_cost"])
				}
				if stats["purchases"] != 10 {
					t.Errorf("Expected 10 purchases, got %v", stats["purchases"])
				}
				if stats["conversion_rate"] != 10.0 {
					t.Errorf("Expected conversion rate 10.0%%, got %v", stats["conversion_rate"])
				}
			},
			expectError: false,
		},
		{
			name: "zero conversations - zero conversion rate",
			days: 30,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ai_conversations WHERE started_at").
					WithArgs(30).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectQuery("SELECT COALESCE\\(SUM\\(total_messages\\)").
					WithArgs(30).
					WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(0))

				mock.ExpectQuery("SELECT COALESCE\\(SUM\\(total_cost\\)").
					WithArgs(30).
					WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(0))

				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ai_conversations WHERE resulted_in_purchase").
					WithArgs(30).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			checkResult: func(t *testing.T, stats map[string]interface{}) {
				if stats["conversion_rate"] != 0.0 {
					t.Errorf("Expected conversion rate 0.0%%, got %v", stats["conversion_rate"])
				}
			},
			expectError: false,
		},
		{
			name: "database error on first query",
			days: 7,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ai_conversations WHERE started_at").
					WithArgs(7).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
			service := &Service{
				db:       db,
				aiClient: nil,
				logger:   logger,
			}

			stats, err := service.GetConversationStats(tt.days)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkResult != nil {
					tt.checkResult(t, stats)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestUpdatePreferencesFromMessage(t *testing.T) {
	// This is a void function that logs, so we test it doesn't panic
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service := &Service{
		db:       nil, // Not used in this function
		aiClient: nil,
		logger:   logger,
	}

	tests := []struct {
		name      string
		userID    *int
		sessionID *string
		message   string
	}{
		{"with widget mention", intPtr(1), stringPtr("sess-1"), "I want a widget"},
		{"with subscription mention", intPtr(2), stringPtr("sess-2"), "Show me subscription plans"},
		{"with budget mention", nil, stringPtr("sess-3"), "I have $30 to spend"},
		{"generic message", nil, stringPtr("sess-4"), "Hello"},
		{"empty message", intPtr(5), stringPtr("sess-5"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure it doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Function panicked: %v", r)
				}
			}()

			service.updatePreferencesFromMessage(tt.userID, tt.sessionID, tt.message)
		})
	}
}

func TestNewService(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	// Create a mock AI client
	mockClient := &MockAIClient{}

	service := NewService(db, mockClient, logger)

	if service == nil {
		t.Fatal("NewService returned nil")
	}

	if service.db == nil {
		t.Error("Service db is nil")
	}

	if service.aiClient == nil {
		t.Error("Service aiClient is nil")
	}

	if service.logger == nil {
		t.Error("Service logger is nil")
	}
}

// MockAIClient for testing
type MockAIClient struct {
	GenerateResponseFunc func(messages []Message, context string) (*ChatResponse, error)
	GetEmbeddingFunc     func(text string) ([]float64, error)
	SpeechToTextFunc     func(audioData []byte, language string) (string, error)
	TextToSpeechFunc     func(text, voice string) ([]byte, error)
}

func (m *MockAIClient) GenerateResponse(messages []Message, context string) (*ChatResponse, error) {
	if m.GenerateResponseFunc != nil {
		return m.GenerateResponseFunc(messages, context)
	}
	return &ChatResponse{
		Message:        "Test response",
		TokensUsed:     10,
		ResponseTimeMs: 100,
	}, nil
}

func (m *MockAIClient) GetEmbedding(text string) ([]float64, error) {
	if m.GetEmbeddingFunc != nil {
		return m.GetEmbeddingFunc(text)
	}
	return []float64{0.1, 0.2, 0.3}, nil
}

func (m *MockAIClient) SpeechToText(audioData []byte, language string) (string, error) {
	if m.SpeechToTextFunc != nil {
		return m.SpeechToTextFunc(audioData, language)
	}
	return "Test transcription", nil
}

func (m *MockAIClient) TextToSpeech(text, voice string) ([]byte, error) {
	if m.TextToSpeechFunc != nil {
		return m.TextToSpeechFunc(text, voice)
	}
	return []byte("test audio data"), nil
}
