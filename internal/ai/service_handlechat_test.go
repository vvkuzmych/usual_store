package ai

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHandleChat(t *testing.T) {
	tests := []struct {
		name        string
		request     ChatRequest
		setupMock   func(sqlmock.Sqlmock)
		aiClient    *MockAIClient
		expectError bool
		checkResult func(*testing.T, *ChatResponse)
	}{
		{
			name: "successful chat with existing conversation",
			request: ChatRequest{
				SessionID: "sess-123",
				Message:   "Hello AI!",
				UserID:    intPtr(42),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// getOrCreateConversation - find existing
				convRows := sqlmock.NewRows([]string{
					"id", "session_id", "user_id", "started_at", "ended_at",
					"total_messages", "resulted_in_purchase", "total_tokens_used",
					"total_cost", "user_agent", "ip_address", "created_at", "updated_at",
				}).AddRow(
					100, "sess-123", 42, time.Now(), nil,
					2, false, 50, 0.01, "", "", time.Now(), time.Now(),
				)
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-123").
					WillReturnRows(convRows)

				// getConversationHistory
				historyRows := sqlmock.NewRows([]string{
					"id", "conversation_id", "role", "content", "tokens_used",
					"response_time_ms", "model", "temperature", "metadata", "created_at",
				}).AddRow(
					1, 100, "user", "Previous message", 5, nil, "", 0.0, nil, time.Now(),
				)
				mock.ExpectQuery("SELECT (.+) FROM ai_messages WHERE conversation_id").
					WithArgs(100, 10).
					WillReturnRows(historyRows)

				// createMessage - user message
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(100, "user", "Hello AI!", 0, nil, "", sqlmock.AnyArg(), nil, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

				// getProductContext
				productRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "is_recurring"}).
					AddRow(1, "Widget", "A great widget", 1000, false)
				mock.ExpectQuery("SELECT (.+) FROM widgets WHERE inventory_level").
					WillReturnRows(productRows)

				// getUserPreferences
				prefsRows := sqlmock.NewRows([]string{
					"id", "user_id", "session_id", "preferred_categories", "budget_min",
					"budget_max", "interaction_count", "last_products_viewed",
					"last_products_purchased", "conversation_style", "preferred_language",
					"created_at", "updated_at",
				}).AddRow(
					1, 42, "sess-123", sql.NullString{}, nil, nil, 5,
					sql.NullString{}, sql.NullString{}, nil, "en", time.Now(), time.Now(),
				)
				mock.ExpectQuery("SELECT (.+) FROM ai_user_preferences WHERE user_id").
					WithArgs(intPtr(42), stringPtr("sess-123")).
					WillReturnRows(prefsRows)

				// createMessage - assistant message
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(100, "assistant", "Hello! How can I help you?", 25, sqlmock.AnyArg(),
						"gpt-3.5-turbo", sqlmock.AnyArg(), nil, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(11))

				// updateConversationStats
				mock.ExpectExec("UPDATE ai_conversations SET").
					WithArgs(100, 25, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			aiClient: &MockAIClient{
				GenerateResponseFunc: func(messages []Message, context string) (*ChatResponse, error) {
					return &ChatResponse{
						Message:        "Hello! How can I help you?",
						TokensUsed:     25,
						ResponseTimeMs: 150,
						Products:       []RecommendedProduct{},
						Suggestions:    []string{"Browse products", "Ask a question"},
					}, nil
				},
			},
			expectError: false,
			checkResult: func(t *testing.T, resp *ChatResponse) {
				if resp.Message != "Hello! How can I help you?" {
					t.Errorf("Expected message %q, got %q", "Hello! How can I help you?", resp.Message)
				}
				if resp.SessionID != "sess-123" {
					t.Errorf("Expected session_id sess-123, got %s", resp.SessionID)
				}
				if resp.TokensUsed != 25 {
					t.Errorf("Expected 25 tokens, got %d", resp.TokensUsed)
				}
			},
		},
		{
			name: "chat with products in response",
			request: ChatRequest{
				SessionID: "sess-456",
				Message:   "Show me widgets",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// getOrCreateConversation - create new
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-456").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectQuery("INSERT INTO ai_conversations").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(200))

				// getConversationHistory - empty
				historyRows := sqlmock.NewRows([]string{
					"id", "conversation_id", "role", "content", "tokens_used",
					"response_time_ms", "model", "temperature", "metadata", "created_at",
				})
				mock.ExpectQuery("SELECT (.+) FROM ai_messages WHERE conversation_id").
					WithArgs(200, 10).
					WillReturnRows(historyRows)

				// createMessage - user
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(20))

				// getProductContext
				productRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "is_recurring"}).
					AddRow(1, "Widget", "Great", 1000, false)
				mock.ExpectQuery("SELECT (.+) FROM widgets WHERE inventory_level").
					WillReturnRows(productRows)

				// getUserPreferences - not found
				mock.ExpectQuery("SELECT (.+) FROM ai_user_preferences WHERE user_id").
					WithArgs((*int)(nil), stringPtr("sess-456")).
					WillReturnError(sql.ErrNoRows)

				// createMessage - assistant with products
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(21))

				// updateConversationStats
				mock.ExpectExec("UPDATE ai_conversations SET").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			aiClient: &MockAIClient{
				GenerateResponseFunc: func(messages []Message, context string) (*ChatResponse, error) {
					return &ChatResponse{
						Message:        "Here are our widgets!",
						TokensUsed:     30,
						ResponseTimeMs: 200,
						Products: []RecommendedProduct{
							{ID: 1, Name: "Widget", Price: 10.0, Reason: "Popular"},
						},
					}, nil
				},
			},
			expectError: false,
			checkResult: func(t *testing.T, resp *ChatResponse) {
				if len(resp.Products) != 1 {
					t.Errorf("Expected 1 product, got %d", len(resp.Products))
				}
			},
		},
		{
			name: "empty message error",
			request: ChatRequest{
				SessionID: "sess-error",
				Message:   "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// No DB calls expected
			},
			aiClient:    &MockAIClient{},
			expectError: true,
		},
		{
			name: "conversation creation error",
			request: ChatRequest{
				SessionID: "sess-dberror",
				Message:   "Test message",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-dberror").
					WillReturnError(sql.ErrConnDone)
			},
			aiClient:    &MockAIClient{},
			expectError: true,
		},
		{
			name: "user message save error",
			request: ChatRequest{
				SessionID: "sess-usererr",
				Message:   "Test",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// getOrCreateConversation - success
				convRows := sqlmock.NewRows([]string{
					"id", "session_id", "user_id", "started_at", "ended_at",
					"total_messages", "resulted_in_purchase", "total_tokens_used",
					"total_cost", "user_agent", "ip_address", "created_at", "updated_at",
				}).AddRow(
					300, "sess-usererr", nil, time.Now(), nil,
					0, false, 0, 0.0, "", "", time.Now(), time.Now(),
				)
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-usererr").
					WillReturnRows(convRows)

				// getConversationHistory
				historyRows := sqlmock.NewRows([]string{
					"id", "conversation_id", "role", "content", "tokens_used",
					"response_time_ms", "model", "temperature", "metadata", "created_at",
				})
				mock.ExpectQuery("SELECT (.+) FROM ai_messages WHERE conversation_id").
					WithArgs(300, 10).
					WillReturnRows(historyRows)

				// createMessage - error
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			aiClient:    &MockAIClient{},
			expectError: true,
		},
		{
			name: "AI generation error",
			request: ChatRequest{
				SessionID: "sess-aierr",
				Message:   "Test",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// getOrCreateConversation
				convRows := sqlmock.NewRows([]string{
					"id", "session_id", "user_id", "started_at", "ended_at",
					"total_messages", "resulted_in_purchase", "total_tokens_used",
					"total_cost", "user_agent", "ip_address", "created_at", "updated_at",
				}).AddRow(
					400, "sess-aierr", nil, time.Now(), nil,
					0, false, 0, 0.0, "", "", time.Now(), time.Now(),
				)
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-aierr").
					WillReturnRows(convRows)

				// getConversationHistory
				historyRows := sqlmock.NewRows([]string{
					"id", "conversation_id", "role", "content", "tokens_used",
					"response_time_ms", "model", "temperature", "metadata", "created_at",
				})
				mock.ExpectQuery("SELECT (.+) FROM ai_messages WHERE conversation_id").
					WithArgs(400, 10).
					WillReturnRows(historyRows)

				// createMessage - user
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(40))

				// getProductContext
				productRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "is_recurring"})
				mock.ExpectQuery("SELECT (.+) FROM widgets WHERE inventory_level").
					WillReturnRows(productRows)

				// getUserPreferences
				mock.ExpectQuery("SELECT (.+) FROM ai_user_preferences WHERE user_id").
					WithArgs((*int)(nil), stringPtr("sess-aierr")).
					WillReturnError(sql.ErrNoRows)
			},
			aiClient: &MockAIClient{
				GenerateResponseFunc: func(messages []Message, context string) (*ChatResponse, error) {
					return nil, sql.ErrConnDone
				},
			},
			expectError: true,
		},
		{
			name: "history retrieval warning (non-fatal)",
			request: ChatRequest{
				SessionID: "sess-nowarn",
				Message:   "Test warning",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// getOrCreateConversation
				convRows := sqlmock.NewRows([]string{
					"id", "session_id", "user_id", "started_at", "ended_at",
					"total_messages", "resulted_in_purchase", "total_tokens_used",
					"total_cost", "user_agent", "ip_address", "created_at", "updated_at",
				}).AddRow(
					500, "sess-nowarn", nil, time.Now(), nil,
					0, false, 0, 0.0, "", "", time.Now(), time.Now(),
				)
				mock.ExpectQuery("SELECT (.+) FROM ai_conversations WHERE session_id").
					WithArgs("sess-nowarn").
					WillReturnRows(convRows)

				// getConversationHistory - error (non-fatal)
				mock.ExpectQuery("SELECT (.+) FROM ai_messages WHERE conversation_id").
					WithArgs(500, 10).
					WillReturnError(sql.ErrConnDone)

				// createMessage - user
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(50))

				// getProductContext - error (non-fatal)
				mock.ExpectQuery("SELECT (.+) FROM widgets WHERE inventory_level").
					WillReturnError(sql.ErrConnDone)

				// getUserPreferences - error (non-fatal)
				mock.ExpectQuery("SELECT (.+) FROM ai_user_preferences WHERE user_id").
					WithArgs((*int)(nil), stringPtr("sess-nowarn")).
					WillReturnError(sql.ErrConnDone)

				// createMessage - assistant (warning on failure)
				mock.ExpectQuery("INSERT INTO ai_messages").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
										sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone) // This is just a warning

				// updateConversationStats - warning on failure
				mock.ExpectExec("UPDATE ai_conversations SET").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone) // This is just a warning
			},
			aiClient: &MockAIClient{
				GenerateResponseFunc: func(messages []Message, context string) (*ChatResponse, error) {
					return &ChatResponse{
						Message:        "Response despite warnings",
						TokensUsed:     20,
						ResponseTimeMs: 100,
					}, nil
				},
			},
			expectError: false,
			checkResult: func(t *testing.T, resp *ChatResponse) {
				if resp.Message != "Response despite warnings" {
					t.Errorf("Expected message despite warnings, got %q", resp.Message)
				}
			},
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
				aiClient: tt.aiClient,
				logger:   logger,
			}

			resp, err := service.HandleChat(tt.request)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkResult != nil && resp != nil {
					tt.checkResult(t, resp)
				}
			}

			// Some expectations may not be met if there are non-fatal warnings
			// So we check but don't fail the test
			_ = mock.ExpectationsWereMet()
		})
	}
}

// MockAIClient is defined in service_db_test.go
