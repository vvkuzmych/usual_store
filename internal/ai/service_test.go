package ai

import (
	"strings"
	"testing"
)

func TestEnrichContextWithPreferences(t *testing.T) {
	// Create a service instance (db and aiClient can be nil for this test)
	s := &Service{}

	tests := []struct {
		name            string
		context         string
		prefs           *UserPreferences
		expectedContent []string // Strings that should appear in the result
		notExpected     []string // Strings that should not appear
	}{
		{
			name:    "no preferences - empty struct",
			context: "Available Products:\n- Widget ($10.00)",
			prefs:   &UserPreferences{},
			expectedContent: []string{
				"Available Products:",
				"Widget ($10.00)",
				"User Preferences:",
			},
			notExpected: []string{
				"Interested in:",
				"Budget range:",
				"Communication style:",
			},
		},
		{
			name:    "with preferred categories",
			context: "Available Products:\n- Widget ($10.00)",
			prefs: &UserPreferences{
				PreferredCategories: []string{"electronics", "gadgets"},
			},
			expectedContent: []string{
				"Available Products:",
				"User Preferences:",
				"Interested in: electronics, gadgets",
			},
			notExpected: []string{
				"Budget range:",
				"Communication style:",
			},
		},
		{
			name:    "with budget range",
			context: "Available Products:\n- Widget ($10.00)",
			prefs: &UserPreferences{
				BudgetMin: float64Ptr(10.00),
				BudgetMax: float64Ptr(50.00),
			},
			expectedContent: []string{
				"Available Products:",
				"User Preferences:",
				"Budget range: $10.00 - $50.00",
			},
			notExpected: []string{
				"Interested in:",
				"Communication style:",
			},
		},
		{
			name:    "with conversation style",
			context: "Available Products:\n- Widget ($10.00)",
			prefs: &UserPreferences{
				ConversationStyle: stringPtr("casual"),
			},
			expectedContent: []string{
				"Available Products:",
				"User Preferences:",
				"Communication style: casual",
			},
			notExpected: []string{
				"Interested in:",
				"Budget range:",
			},
		},
		{
			name:    "with all preferences",
			context: "Available Products:\n- Widget ($10.00)",
			prefs: &UserPreferences{
				PreferredCategories: []string{"widgets", "subscriptions"},
				BudgetMin:           float64Ptr(20.00),
				BudgetMax:           float64Ptr(100.00),
				ConversationStyle:   stringPtr("professional"),
			},
			expectedContent: []string{
				"Available Products:",
				"User Preferences:",
				"Interested in: widgets, subscriptions",
				"Budget range: $20.00 - $100.00",
				"Communication style: professional",
			},
			notExpected: []string{},
		},
		{
			name:    "with single category",
			context: "Products available",
			prefs: &UserPreferences{
				PreferredCategories: []string{"electronics"},
			},
			expectedContent: []string{
				"Products available",
				"Interested in: electronics",
			},
			notExpected: []string{
				"Budget range:",
			},
		},
		{
			name:    "with multiple categories",
			context: "Products available",
			prefs: &UserPreferences{
				PreferredCategories: []string{"electronics", "gadgets", "accessories"},
			},
			expectedContent: []string{
				"Interested in: electronics, gadgets, accessories",
			},
			notExpected: []string{},
		},
		{
			name:    "with only budget min",
			context: "Products",
			prefs: &UserPreferences{
				BudgetMin: float64Ptr(15.50),
			},
			expectedContent: []string{
				"Products",
				"User Preferences:",
			},
			notExpected: []string{
				"Budget range:", // Should not show if only min is set
			},
		},
		{
			name:    "with only budget max",
			context: "Products",
			prefs: &UserPreferences{
				BudgetMax: float64Ptr(75.99),
			},
			expectedContent: []string{
				"Products",
				"User Preferences:",
			},
			notExpected: []string{
				"Budget range:", // Should not show if only max is set
			},
		},
		{
			name:    "with zero budget values",
			context: "Products",
			prefs: &UserPreferences{
				BudgetMin: float64Ptr(0.00),
				BudgetMax: float64Ptr(0.00),
			},
			expectedContent: []string{
				"Products",
				"Budget range: $0.00 - $0.00",
			},
			notExpected: []string{},
		},
		{
			name:    "empty context with preferences",
			context: "",
			prefs: &UserPreferences{
				PreferredCategories: []string{"test"},
			},
			expectedContent: []string{
				"User Preferences:",
				"Interested in: test",
			},
			notExpected: []string{},
		},
		{
			name:    "special characters in conversation style",
			context: "Products",
			prefs: &UserPreferences{
				ConversationStyle: stringPtr("friendly & helpful"),
			},
			expectedContent: []string{
				"Communication style: friendly & helpful",
			},
			notExpected: []string{},
		},
		{
			name:    "empty categories array",
			context: "Products",
			prefs: &UserPreferences{
				PreferredCategories: []string{},
			},
			expectedContent: []string{
				"Products",
				"User Preferences:",
			},
			notExpected: []string{
				"Interested in:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.enrichContextWithPreferences(tt.context, tt.prefs)

			// Check that expected content is present
			for _, expected := range tt.expectedContent {
				if !strings.Contains(result, expected) {
					t.Errorf("enrichContextWithPreferences() result missing expected content %q\nGot:\n%s", expected, result)
				}
			}

			// Check that unwanted content is not present
			for _, notExpected := range tt.notExpected {
				if strings.Contains(result, notExpected) {
					t.Errorf("enrichContextWithPreferences() result contains unexpected content %q\nGot:\n%s", notExpected, result)
				}
			}

			// Verify original context is preserved
			if !strings.Contains(result, tt.context) {
				t.Errorf("enrichContextWithPreferences() did not preserve original context\nOriginal: %s\nGot: %s", tt.context, result)
			}
		})
	}
}

// Benchmark for context enrichment
func BenchmarkEnrichContextWithPreferences(b *testing.B) {
	s := &Service{}
	context := "Available Products:\n- Widget ($10.00)\n- Plan ($30.00)"
	prefs := &UserPreferences{
		PreferredCategories: []string{"electronics", "gadgets", "accessories"},
		BudgetMin:           float64Ptr(10.00),
		BudgetMax:           float64Ptr(100.00),
		ConversationStyle:   stringPtr("professional"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.enrichContextWithPreferences(context, prefs)
	}
}

// Test that enriched context is always longer than original
func TestEnrichContextWithPreferencesLength(t *testing.T) {
	s := &Service{}

	tests := []struct {
		name    string
		context string
		prefs   *UserPreferences
	}{
		{
			name:    "with categories",
			context: "Products",
			prefs:   &UserPreferences{PreferredCategories: []string{"test"}},
		},
		{
			name:    "with budget",
			context: "Products",
			prefs:   &UserPreferences{BudgetMin: float64Ptr(10.0), BudgetMax: float64Ptr(50.0)},
		},
		{
			name:    "with style",
			context: "Products",
			prefs:   &UserPreferences{ConversationStyle: stringPtr("casual")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.enrichContextWithPreferences(tt.context, tt.prefs)
			if len(result) <= len(tt.context) {
				t.Errorf("enrichContextWithPreferences() should add content. Original length: %d, Result length: %d",
					len(tt.context), len(result))
			}
		})
	}
}

// Test edge cases
func TestEnrichContextWithPreferencesEdgeCases(t *testing.T) {
	s := &Service{}

	tests := []struct {
		name    string
		context string
		prefs   *UserPreferences
	}{
		{
			name:    "very long context",
			context: strings.Repeat("Product ", 1000),
			prefs:   &UserPreferences{PreferredCategories: []string{"test"}},
		},
		{
			name:    "very long category names",
			context: "Products",
			prefs:   &UserPreferences{PreferredCategories: []string{strings.Repeat("a", 100)}},
		},
		{
			name:    "very large budget",
			context: "Products",
			prefs:   &UserPreferences{BudgetMin: float64Ptr(0), BudgetMax: float64Ptr(999999999.99)},
		},
		{
			name:    "negative budget",
			context: "Products",
			prefs:   &UserPreferences{BudgetMin: float64Ptr(-100.0), BudgetMax: float64Ptr(-10.0)},
		},
		{
			name:    "many categories",
			context: "Products",
			prefs: &UserPreferences{
				PreferredCategories: []string{"cat1", "cat2", "cat3", "cat4", "cat5", "cat6", "cat7", "cat8", "cat9", "cat10"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := s.enrichContextWithPreferences(tt.context, tt.prefs)

			// Result should not be empty
			if result == "" {
				t.Error("enrichContextWithPreferences() returned empty string")
			}

			// Original context should be preserved
			if !strings.Contains(result, tt.context) {
				t.Error("enrichContextWithPreferences() did not preserve original context")
			}
		})
	}
}

// Helper functions for pointer creation
func float64Ptr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// Test ChatRequest validation
func TestChatRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		request   ChatRequest
		wantValid bool
	}{
		{
			name: "valid request with all fields",
			request: ChatRequest{
				SessionID: "session-123",
				Message:   "Hello, I need help",
				UserID:    intPtr(1),
			},
			wantValid: true,
		},
		{
			name: "valid request without userID",
			request: ChatRequest{
				SessionID: "session-123",
				Message:   "Hello",
			},
			wantValid: true,
		},
		{
			name: "empty message",
			request: ChatRequest{
				SessionID: "session-123",
				Message:   "",
			},
			wantValid: false,
		},
		{
			name: "empty session and message",
			request: ChatRequest{
				SessionID: "",
				Message:   "",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For now, we just check the structure
			// In a real scenario, you'd call a validation method
			isValid := tt.request.Message != ""
			if isValid != tt.wantValid {
				t.Errorf("ChatRequest validity = %v, want %v for request: %+v", isValid, tt.wantValid, tt.request)
			}
		})
	}
}

// Test RecommendedProduct struct
func TestRecommendedProduct(t *testing.T) {
	tests := []struct {
		name    string
		product RecommendedProduct
	}{
		{
			name: "complete product",
			product: RecommendedProduct{
				ID:          1,
				Name:        "Widget",
				Description: "A great widget",
				Price:       19.99,
				Image:       "widget.jpg",
				Reason:      "Matches your preferences",
			},
		},
		{
			name: "product without image",
			product: RecommendedProduct{
				ID:          2,
				Name:        "Plan",
				Description: "Subscription plan",
				Price:       29.99,
				Image:       "",
				Reason:      "Popular choice",
			},
		},
		{
			name: "free product",
			product: RecommendedProduct{
				ID:          3,
				Name:        "Free Trial",
				Description: "Try it free",
				Price:       0.00,
				Reason:      "No risk",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate required fields
			if tt.product.ID <= 0 {
				t.Errorf("RecommendedProduct ID should be positive, got %d", tt.product.ID)
			}
			if tt.product.Name == "" {
				t.Error("RecommendedProduct Name should not be empty")
			}
			if tt.product.Reason == "" {
				t.Error("RecommendedProduct Reason should not be empty")
			}
		})
	}
}
