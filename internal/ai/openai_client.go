package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	OpenAIAPIURL       = "https://api.openai.com/v1/chat/completions"
	OpenAIEmbeddingURL = "https://api.openai.com/v1/embeddings"
)

// OpenAIClient implements AIClient interface
type OpenAIClient struct {
	APIKey      string
	Model       string // "gpt-4", "gpt-3.5-turbo", etc.
	Temperature float64
	MaxTokens   int
	HTTPClient  *http.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey, model string, temperature float64) *OpenAIClient {
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	if temperature == 0 {
		temperature = 0.7
	}

	return &OpenAIClient{
		APIKey:      apiKey,
		Model:       model,
		Temperature: temperature,
		MaxTokens:   500,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// OpenAI API request/response structures
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type openAIError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// GenerateResponse generates a response using OpenAI API
func (c *OpenAIClient) GenerateResponse(messages []Message, context string) (*ChatResponse, error) {
	startTime := time.Now()

	// Convert our messages to OpenAI format
	openAIMessages := make([]openAIMessage, 0, len(messages)+2)

	// Add system message with context
	systemPrompt := c.buildSystemPrompt(context)
	openAIMessages = append(openAIMessages, openAIMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	// Add conversation history
	for _, msg := range messages {
		openAIMessages = append(openAIMessages, openAIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Create request
	reqBody := openAIRequest{
		Model:       c.Model,
		Messages:    openAIMessages,
		Temperature: c.Temperature,
		MaxTokens:   c.MaxTokens,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make API call
	req, err := http.NewRequest("POST", OpenAIAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		var apiErr openAIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			return nil, fmt.Errorf("OpenAI API error: %s", apiErr.Error.Message)
		}
		return nil, fmt.Errorf("OpenAI API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResp openAIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	responseTime := int(time.Since(startTime).Milliseconds())

	// Extract products mentioned (simple keyword extraction)
	products := c.extractProductsFromResponse(apiResp.Choices[0].Message.Content)

	return &ChatResponse{
		Message:        apiResp.Choices[0].Message.Content,
		TokensUsed:     apiResp.Usage.TotalTokens,
		ResponseTimeMs: responseTime,
		Products:       products,
		Suggestions:    c.generateSuggestions(apiResp.Choices[0].Message.Content),
	}, nil
}

// GetEmbedding generates embeddings for text (for semantic search)
func (c *OpenAIClient) GetEmbedding(text string) ([]float64, error) {
	reqBody := map[string]interface{}{
		"model": "text-embedding-ada-002",
		"input": text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", OpenAIEmbeddingURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embeddings in response")
	}

	return result.Data[0].Embedding, nil
}

// buildSystemPrompt creates the system prompt with product context
func (c *OpenAIClient) buildSystemPrompt(productContext string) string {
	return fmt.Sprintf(`You are a helpful shopping assistant for Usual Store, an online store specializing in widgets and subscription plans.

Your role:
- Help customers find the perfect products for their needs
- Answer questions about products, pricing, and features
- Provide personalized recommendations
- Be friendly, concise, and helpful
- If you don't know something, be honest and offer to help in another way

Current products available:
%s

Guidelines:
- Keep responses under 150 words unless more detail is specifically requested
- Always include product names and prices when recommending
- Ask clarifying questions when user needs are unclear
- Be enthusiastic but not pushy
- Focus on value and benefits, not just features

Remember: You're here to help customers make informed decisions and have a great shopping experience!`, productContext)
}

// extractProductsFromResponse extracts product mentions from the response
func (c *OpenAIClient) extractProductsFromResponse(content string) []RecommendedProduct {
	// This is a simple implementation
	// In production, you'd parse the content more intelligently or have the AI return structured data
	products := []RecommendedProduct{}

	// TODO: Implement proper product extraction or use function calling feature of GPT-4

	return products
}

// generateSuggestions generates follow-up suggestions for the user
func (c *OpenAIClient) generateSuggestions(content string) []string {
	// Simple suggestions based on common patterns
	suggestions := []string{
		"Tell me more about this product",
		"What other options do you have?",
		"What's your best deal?",
	}

	// TODO: Make this smarter based on conversation context

	return suggestions
}

// CalculateCost estimates the cost of the API call
func (c *OpenAIClient) CalculateCost(tokens int) float64 {
	// Pricing as of 2024 (adjust as needed)
	// GPT-3.5-turbo: $0.0015 per 1K input tokens, $0.002 per 1K output tokens
	// GPT-4: $0.03 per 1K input tokens, $0.06 per 1K output tokens

	var costPer1K float64
	switch c.Model {
	case "gpt-4", "gpt-4-0613":
		costPer1K = 0.045 // average of input and output
	case "gpt-3.5-turbo", "gpt-3.5-turbo-0613":
		costPer1K = 0.00175 // average of input and output
	default:
		costPer1K = 0.00175
	}

	return float64(tokens) / 1000.0 * costPer1K
}
