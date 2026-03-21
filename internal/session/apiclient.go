package session

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// APIClient makes direct calls to an Anthropic-compatible API.
// Supports MiniMax, Anthropic, and any API that implements the Anthropic messages format.
type APIClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

// APIUsage holds token usage from an API call.
type APIUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// APIResponse holds the parsed response from the API.
type APIResponse struct {
	Text         string
	Usage        APIUsage
	Model        string
	StopReason   string
	DurationMS   int64
}

// NewAPIClient creates a client from environment variables.
// ANTHROPIC_API_KEY — API key (required for oneshot mode)
// ANTHROPIC_BASE_URL — Base URL (default: https://api.anthropic.com)
// ONESHOT_MODEL — Model to use (default: claude-sonnet-4-6)
func NewAPIClient() *APIClient {
	baseURL := os.Getenv("ANTHROPIC_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	model := os.Getenv("ONESHOT_MODEL")
	if model == "" {
		model = "claude-sonnet-4-6"
	}

	client := &APIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}

	if apiKey != "" {
		log.Printf("[API-CLIENT] Initialized: base=%s model=%s", baseURL, model)
	} else {
		log.Printf("[API-CLIENT] No ANTHROPIC_API_KEY set — oneshot mode unavailable, will fall back to session")
	}

	return client
}

// IsAvailable returns true if the API client has credentials configured.
func (c *APIClient) IsAvailable() bool {
	return c.APIKey != ""
}

// SendMessage sends a one-shot message to the API and returns the text response.
func (c *APIClient) SendMessage(ctx context.Context, systemPrompt, userMessage string) (*APIResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("API client not configured (no ANTHROPIC_API_KEY)")
	}

	// Build request body
	body := map[string]interface{}{
		"model":      c.Model,
		"max_tokens": 8192,
		"system":     systemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": userMessage},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Build HTTP request
	url := c.BaseURL + "/v1/messages"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	start := time.Now()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	duration := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var apiResp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Model      string   `json:"model"`
		StopReason string   `json:"stop_reason"`
		Usage      APIUsage `json:"usage"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	// Extract text from content blocks
	var text string
	for _, block := range apiResp.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}

	return &APIResponse{
		Text:       text,
		Usage:      apiResp.Usage,
		Model:      apiResp.Model,
		StopReason: apiResp.StopReason,
		DurationMS: duration.Milliseconds(),
	}, nil
}
