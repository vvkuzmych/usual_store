package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// MockHTTPClient is a mock HTTP client for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewOpenAIClient(t *testing.T) {
	tests := []struct {
		name              string
		apiKey            string
		model             string
		temperature       float64
		expectedModel     string
		expectedTemp      float64
		expectedMaxTokens int
	}{
		{
			name:              "with all parameters",
			apiKey:            "sk-test123",
			model:             "gpt-4",
			temperature:       0.8,
			expectedModel:     "gpt-4",
			expectedTemp:      0.8,
			expectedMaxTokens: 500,
		},
		{
			name:              "with default model",
			apiKey:            "sk-test456",
			model:             "",
			temperature:       0.5,
			expectedModel:     "gpt-3.5-turbo",
			expectedTemp:      0.5,
			expectedMaxTokens: 500,
		},
		{
			name:              "with default temperature",
			apiKey:            "sk-test789",
			model:             "gpt-3.5-turbo",
			temperature:       0,
			expectedModel:     "gpt-3.5-turbo",
			expectedTemp:      0.7,
			expectedMaxTokens: 500,
		},
		{
			name:              "with all defaults",
			apiKey:            "sk-test000",
			model:             "",
			temperature:       0,
			expectedModel:     "gpt-3.5-turbo",
			expectedTemp:      0.7,
			expectedMaxTokens: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewOpenAIClient(tt.apiKey, tt.model, tt.temperature)

			if client == nil {
				t.Fatal("NewOpenAIClient returned nil")
			}

			if client.APIKey != tt.apiKey {
				t.Errorf("APIKey = %q, want %q", client.APIKey, tt.apiKey)
			}

			if client.Model != tt.expectedModel {
				t.Errorf("Model = %q, want %q", client.Model, tt.expectedModel)
			}

			if client.Temperature != tt.expectedTemp {
				t.Errorf("Temperature = %v, want %v", client.Temperature, tt.expectedTemp)
			}

			if client.MaxTokens != tt.expectedMaxTokens {
				t.Errorf("MaxTokens = %d, want %d", client.MaxTokens, tt.expectedMaxTokens)
			}

			if client.HTTPClient == nil {
				t.Error("HTTPClient is nil")
			}
		})
	}
}

func TestGenerateResponse(t *testing.T) {
	tests := []struct {
		name           string
		messages       []Message
		context        string
		mockHTTPFunc   func(req *http.Request) (*http.Response, error)
		expectedError  bool
		errorContains  string
		validateResult func(*testing.T, *ChatResponse)
	}{
		{
			name: "successful response",
			messages: []Message{
				{Role: "user", Content: "Hello"},
			},
			context: "Available Products:\n- Widget ($10.00)",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request
				if req.Header.Get("Authorization") == "" {
					t.Error("Authorization header missing")
				}

				// Create mock response
				apiResp := openAIResponse{
					ID:      "chatcmpl-123",
					Object:  "chat.completion",
					Created: time.Now().Unix(),
					Model:   "gpt-3.5-turbo",
					Choices: []struct {
						Index   int `json:"index"`
						Message struct {
							Role    string `json:"role"`
							Content string `json:"content"`
						} `json:"message"`
						FinishReason string `json:"finish_reason"`
					}{
						{
							Index: 0,
							Message: struct {
								Role    string `json:"role"`
								Content string `json:"content"`
							}{
								Role:    "assistant",
								Content: "Hello! How can I help you today?",
							},
							FinishReason: "stop",
						},
					},
					Usage: struct {
						PromptTokens     int `json:"prompt_tokens"`
						CompletionTokens int `json:"completion_tokens"`
						TotalTokens      int `json:"total_tokens"`
					}{
						PromptTokens:     10,
						CompletionTokens: 15,
						TotalTokens:      25,
					},
				}

				body, _ := json.Marshal(apiResp)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				}, nil
			},
			expectedError: false,
			validateResult: func(t *testing.T, resp *ChatResponse) {
				if resp.Message != "Hello! How can I help you today?" {
					t.Errorf("Expected message %q, got %q", "Hello! How can I help you today?", resp.Message)
				}
				if resp.TokensUsed != 25 {
					t.Errorf("Expected 25 tokens, got %d", resp.TokensUsed)
				}
				if resp.ResponseTimeMs < 0 {
					t.Error("Expected non-negative response time")
				}
			},
		},
		{
			name:     "HTTP request error",
			messages: []Message{{Role: "user", Content: "Hello"}},
			context:  "Products",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			expectedError: true,
			errorContains: "failed to make request",
		},
		{
			name:     "OpenAI API error response",
			messages: []Message{{Role: "user", Content: "Hello"}},
			context:  "Products",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				apiErr := openAIError{
					Error: struct {
						Message string `json:"message"`
						Type    string `json:"type"`
						Code    string `json:"code"`
					}{
						Message: "Incorrect API key provided",
						Type:    "invalid_request_error",
						Code:    "invalid_api_key",
					},
				}
				body, _ := json.Marshal(apiErr)
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				}, nil
			},
			expectedError: true,
			errorContains: "Incorrect API key provided",
		},
		{
			name:     "empty choices in response",
			messages: []Message{{Role: "user", Content: "Hello"}},
			context:  "Products",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				apiResp := openAIResponse{
					Choices: []struct {
						Index   int `json:"index"`
						Message struct {
							Role    string `json:"role"`
							Content string `json:"content"`
						} `json:"message"`
						FinishReason string `json:"finish_reason"`
					}{},
				}
				body, _ := json.Marshal(apiResp)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				}, nil
			},
			expectedError: true,
			errorContains: "no choices in response",
		},
		{
			name:     "invalid JSON response",
			messages: []Message{{Role: "user", Content: "Hello"}},
			context:  "Products",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
					Header:     make(http.Header),
				}, nil
			},
			expectedError: true,
			errorContains: "failed to parse response",
		},
		{
			name: "multiple messages in history",
			messages: []Message{
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
				{Role: "user", Content: "Show me products"},
			},
			context: "Available Products:\n- Widget ($10.00)\n- Plan ($30.00)",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				apiResp := openAIResponse{
					Choices: []struct {
						Index   int `json:"index"`
						Message struct {
							Role    string `json:"role"`
							Content string `json:"content"`
						} `json:"message"`
						FinishReason string `json:"finish_reason"`
					}{
						{
							Message: struct {
								Role    string `json:"role"`
								Content string `json:"content"`
							}{
								Role:    "assistant",
								Content: "Here are our products: Widget and Plan",
							},
						},
					},
					Usage: struct {
						PromptTokens     int `json:"prompt_tokens"`
						CompletionTokens int `json:"completion_tokens"`
						TotalTokens      int `json:"total_tokens"`
					}{
						TotalTokens: 50,
					},
				}
				body, _ := json.Marshal(apiResp)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				}, nil
			},
			expectedError: false,
			validateResult: func(t *testing.T, resp *ChatResponse) {
				if !strings.Contains(resp.Message, "products") {
					t.Error("Expected response to mention products")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewOpenAIClient("sk-test123", "gpt-3.5-turbo", 0.7)

			// Replace HTTP client with mock
			client.HTTPClient = &http.Client{
				Transport: &mockTransport{doFunc: tt.mockHTTPFunc},
			}

			resp, err := client.GenerateResponse(tt.messages, tt.context)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.validateResult != nil && resp != nil {
					tt.validateResult(t, resp)
				}
			}
		})
	}
}

// mockTransport is a mock HTTP transport
type mockTransport struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestGetEmbedding(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		mockHTTPFunc  func(req *http.Request) (*http.Response, error)
		expectedError bool
		errorContains string
	}{
		{
			name: "successful embedding",
			text: "test text",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				result := struct {
					Data []struct {
						Embedding []float64 `json:"embedding"`
					} `json:"data"`
				}{
					Data: []struct {
						Embedding []float64 `json:"embedding"`
					}{
						{
							Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
						},
					},
				}
				body, _ := json.Marshal(result)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				}, nil
			},
			expectedError: false,
		},
		{
			name: "API error",
			text: "test",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString("Unauthorized")),
					Header:     make(http.Header),
				}, nil
			},
			expectedError: true,
			errorContains: "OpenAI API returned status 401",
		},
		{
			name: "network error",
			text: "test",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("connection refused")
			},
			expectedError: true,
			errorContains: "failed to make request",
		},
		{
			name: "empty embeddings",
			text: "test",
			mockHTTPFunc: func(req *http.Request) (*http.Response, error) {
				result := struct {
					Data []struct {
						Embedding []float64 `json:"embedding"`
					} `json:"data"`
				}{
					Data: []struct {
						Embedding []float64 `json:"embedding"`
					}{},
				}
				body, _ := json.Marshal(result)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				}, nil
			},
			expectedError: true,
			errorContains: "no embeddings in response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewOpenAIClient("sk-test123", "gpt-3.5-turbo", 0.7)
			client.HTTPClient = &http.Client{
				Transport: &mockTransport{doFunc: tt.mockHTTPFunc},
			}

			embedding, err := client.GetEmbedding(tt.text)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(embedding) == 0 {
					t.Error("Expected non-empty embedding")
				}
			}
		})
	}
}

func TestBuildSystemPrompt(t *testing.T) {
	tests := []struct {
		name             string
		context          string
		expectedContains []string
	}{
		{
			name:    "with product context",
			context: "- Widget ($10.00)\n- Plan ($30.00)",
			expectedContains: []string{
				"shopping assistant",
				"Usual Store",
				"Widget ($10.00)",
				"Plan ($30.00)",
				"recommendations",
			},
		},
		{
			name:    "with empty context",
			context: "",
			expectedContains: []string{
				"shopping assistant",
				"helpful",
			},
		},
		{
			name:    "with detailed context",
			context: "Available Products:\n\n- Premium Widget ($50.00, one-time): High quality widget\n- Basic Plan ($10.00, subscription): Entry level",
			expectedContains: []string{
				"Premium Widget",
				"Basic Plan",
				"products available",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewOpenAIClient("sk-test123", "gpt-3.5-turbo", 0.7)
			prompt := client.buildSystemPrompt(tt.context)

			if prompt == "" {
				t.Error("buildSystemPrompt returned empty string")
			}

			for _, expected := range tt.expectedContains {
				if !strings.Contains(prompt, expected) {
					t.Errorf("Expected prompt to contain %q, but it doesn't", expected)
				}
			}
		})
	}
}

func TestExtractProductsFromResponse(t *testing.T) {
	client := NewOpenAIClient("sk-test123", "gpt-3.5-turbo", 0.7)

	tests := []struct {
		name     string
		content  string
		expected int // Expected number of products
	}{
		{
			name:     "response with products",
			content:  "I recommend the Widget ($10) and the Premium Plan ($30)",
			expected: 0, // Current implementation returns empty
		},
		{
			name:     "response without products",
			content:  "I can help you find products",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products := client.extractProductsFromResponse(tt.content)
			if len(products) != tt.expected {
				t.Errorf("Expected %d products, got %d", tt.expected, len(products))
			}
		})
	}
}

func TestGenerateSuggestions(t *testing.T) {
	client := NewOpenAIClient("sk-test123", "gpt-3.5-turbo", 0.7)

	tests := []struct {
		name     string
		content  string
		minCount int
	}{
		{
			name:     "any content generates suggestions",
			content:  "Here are some products for you",
			minCount: 1,
		},
		{
			name:     "empty content generates suggestions",
			content:  "",
			minCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := client.generateSuggestions(tt.content)
			if len(suggestions) < tt.minCount {
				t.Errorf("Expected at least %d suggestions, got %d", tt.minCount, len(suggestions))
			}
		})
	}
}

func TestCalculateCost(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		tokens   int
		expected float64
	}{
		{
			name:     "GPT-4 cost calculation",
			model:    "gpt-4",
			tokens:   1000,
			expected: 0.045,
		},
		{
			name:     "GPT-3.5 cost calculation",
			model:    "gpt-3.5-turbo",
			tokens:   1000,
			expected: 0.00175,
		},
		{
			name:     "default model cost",
			model:    "unknown-model",
			tokens:   1000,
			expected: 0.00175,
		},
		{
			name:     "zero tokens",
			model:    "gpt-3.5-turbo",
			tokens:   0,
			expected: 0.0,
		},
		{
			name:     "large token count",
			model:    "gpt-4",
			tokens:   100000,
			expected: 4.5,
		},
		{
			name:     "gpt-4-0613 specific version",
			model:    "gpt-4-0613",
			tokens:   1000,
			expected: 0.045,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewOpenAIClient("sk-test123", tt.model, 0.7)
			cost := client.CalculateCost(tt.tokens)

			if cost != tt.expected {
				t.Errorf("Expected cost %v, got %v", tt.expected, cost)
			}
		})
	}
}

// Benchmark tests
func BenchmarkGenerateResponse(b *testing.B) {
	client := NewOpenAIClient("sk-test123", "gpt-3.5-turbo", 0.7)
	client.HTTPClient = &http.Client{
		Transport: &mockTransport{
			doFunc: func(req *http.Request) (*http.Response, error) {
				apiResp := openAIResponse{
					Choices: []struct {
						Index   int `json:"index"`
						Message struct {
							Role    string `json:"role"`
							Content string `json:"content"`
						} `json:"message"`
						FinishReason string `json:"finish_reason"`
					}{
						{
							Message: struct {
								Role    string `json:"role"`
								Content string `json:"content"`
							}{
								Role:    "assistant",
								Content: "Test response",
							},
						},
					},
					Usage: struct {
						PromptTokens     int `json:"prompt_tokens"`
						CompletionTokens int `json:"completion_tokens"`
						TotalTokens      int `json:"total_tokens"`
					}{
						TotalTokens: 25,
					},
				}
				body, _ := json.Marshal(apiResp)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
				}, nil
			},
		},
	}

	messages := []Message{{Role: "user", Content: "Hello"}}
	context := "Products"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GenerateResponse(messages, context)
	}
}

func BenchmarkCalculateCost(b *testing.B) {
	client := NewOpenAIClient("sk-test123", "gpt-3.5-turbo", 0.7)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.CalculateCost(1000)
	}
}

func BenchmarkBuildSystemPrompt(b *testing.B) {
	client := NewOpenAIClient("sk-test123", "gpt-3.5-turbo", 0.7)
	context := "Available Products:\n- Widget ($10.00)\n- Plan ($30.00)"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.buildSystemPrompt(context)
	}
}
