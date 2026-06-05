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
	"strings"
	"time"
)

// APIClient makes direct one-shot calls to a chat-completions API.
//
// It speaks two wire formats:
//   - "anthropic" (default): POST {base}/v1/messages, x-api-key header — works
//     with Anthropic, MiniMax, and any Anthropic-compatible endpoint.
//   - "openai": POST {base}/chat/completions, Bearer auth — works with OpenAI,
//     Featherless, Moonshot, and any OpenAI-compatible endpoint.
//
// The session-mode agents (the `claude` CLI) only speak Anthropic, so they keep
// using ANTHROPIC_*. This client (the Brain + one-shot agents) can be pointed at
// an OpenAI-compatible provider independently via ONESHOT_*.
type APIClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	Format     string // "anthropic" | "openai"
	JSONMode   bool   // force structured JSON output + low temperature (for routers)
	HTTPClient *http.Client
}

// APIUsage holds token usage from an API call.
type APIUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// APIResponse holds the parsed response from the API.
type APIResponse struct {
	Text       string
	Usage      APIUsage
	Model      string
	StopReason string
	DurationMS int64
}

// NewAPIClient creates a one-shot client from environment variables.
//
// Brain / one-shot path (this client), in precedence order:
//
//	ONESHOT_FORMAT    — "openai" or "anthropic" (default: "anthropic")
//	ONESHOT_BASE_URL  — base URL; falls back to ANTHROPIC_BASE_URL, then the
//	                    Anthropic default. For OpenAI format point this at the
//	                    provider root incl. /v1 (e.g. https://api.featherless.ai/v1).
//	ONESHOT_API_KEY   — key; for openai format falls back to KIMI_KEY, then
//	                    ANTHROPIC_API_KEY.
//	ONESHOT_MODEL     — model id (default: claude-sonnet-4-6)
func NewAPIClient() *APIClient {
	format := strings.ToLower(strings.TrimSpace(os.Getenv("ONESHOT_FORMAT")))
	if format == "" {
		format = "anthropic"
	}

	baseURL := os.Getenv("ONESHOT_BASE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("ANTHROPIC_BASE_URL")
	}
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	// Key selection: an explicit ONESHOT_API_KEY always wins; otherwise pick the
	// per-provider key that matches the base URL host. This lets every provider's
	// key live in .env and switch providers by swapping ONESHOT_BASE_URL/MODEL only.
	apiKey := os.Getenv("ONESHOT_API_KEY")
	if apiKey == "" {
		switch {
		case strings.Contains(baseURL, "featherless"):
			apiKey = os.Getenv("KIMI_KEY")
		case strings.Contains(baseURL, "googleapis"):
			apiKey = os.Getenv("GEMINI_KEY")
		case strings.Contains(baseURL, "openai.com"):
			apiKey = os.Getenv("OPENAI_KEY")
		}
	}
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	model := os.Getenv("ONESHOT_MODEL")
	if model == "" {
		model = "claude-sonnet-4-6"
	}

	client := &APIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
		Format:  format,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}

	if apiKey != "" {
		hint := apiKey
		if len(hint) > 6 {
			hint = hint[:6]
		}
		log.Printf("[API-CLIENT] Initialized: format=%s base=%s model=%s key=%s…(len %d)", format, baseURL, model, hint, len(apiKey))
	} else {
		log.Printf("[API-CLIENT] No oneshot API key set — oneshot mode unavailable, will fall back to session")
	}

	return client
}

// IsAvailable returns true if the API client has credentials configured.
func (c *APIClient) IsAvailable() bool {
	return c.APIKey != ""
}

// WithModel returns a shallow copy of the client that uses a different model
// (same endpoint, key, format, and JSON mode). Used to run a fast "router" tier
// (e.g. gemini-3.1-flash-lite) alongside the main response model.
func (c *APIClient) WithModel(model string) *APIClient {
	clone := *c
	clone.Model = model
	return &clone
}

// SendMessage sends a one-shot message and returns the text response, routing to
// the configured wire format.
func (c *APIClient) SendMessage(ctx context.Context, systemPrompt, userMessage string) (*APIResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("API client not configured (no oneshot API key)")
	}
	switch c.Format {
	case "openai":
		return c.sendOpenAI(ctx, systemPrompt, userMessage)
	case "gemini":
		return c.sendGemini(ctx, systemPrompt, userMessage)
	default:
		return c.sendAnthropic(ctx, systemPrompt, userMessage)
	}
}

// sendAnthropic uses the Anthropic Messages format (POST /v1/messages).
func (c *APIClient) sendAnthropic(ctx context.Context, systemPrompt, userMessage string) (*APIResponse, error) {
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

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/v1/messages", bytes.NewReader(jsonBody))
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

// sendOpenAI uses the OpenAI chat-completions format (POST /chat/completions).
func (c *APIClient) sendOpenAI(ctx context.Context, systemPrompt, userMessage string) (*APIResponse, error) {
	messages := []map[string]string{}
	if systemPrompt != "" {
		messages = append(messages, map[string]string{"role": "system", "content": systemPrompt})
	}
	messages = append(messages, map[string]string{"role": "user", "content": userMessage})

	body := map[string]interface{}{
		"model":       c.Model,
		"max_tokens":  8192,
		"temperature": 0.6,
		"messages":    messages,
	}
	if c.JSONMode {
		body["response_format"] = map[string]string{"type": "json_object"}
		body["temperature"] = 0.3
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

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

	var apiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Model string `json:"model"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	var text, stop string
	if len(apiResp.Choices) > 0 {
		text = apiResp.Choices[0].Message.Content
		stop = apiResp.Choices[0].FinishReason
	}

	return &APIResponse{
		Text: text,
		Usage: APIUsage{
			InputTokens:  apiResp.Usage.PromptTokens,
			OutputTokens: apiResp.Usage.CompletionTokens,
		},
		Model:      apiResp.Model,
		StopReason: stop,
		DurationMS: duration.Milliseconds(),
	}, nil
}

// sendGemini uses Google's native Generative Language API
// (POST {base}/models/{model}:generateContent, x-goog-api-key header). The
// OpenAI-compat layer (/v1beta/openai) is in beta and rejects valid keys with
// "API key not valid"; the native endpoint is stable and this key is confirmed
// to work against it. Set ONESHOT_BASE_URL to .../v1beta for this format.
func (c *APIClient) sendGemini(ctx context.Context, systemPrompt, userMessage string) (*APIResponse, error) {
	genConfig := map[string]interface{}{
		"maxOutputTokens": 8192,
		"temperature":     0.6,
	}
	if c.JSONMode {
		// Force a single raw JSON object (no markdown/prose wrapping) and lower
		// temperature for deterministic routing — fixes the prose-wrapped-JSON
		// failure mode in the Brain router.
		genConfig["responseMimeType"] = "application/json"
		genConfig["temperature"] = 0.3
	}
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"role": "user", "parts": []map[string]string{{"text": userMessage}}},
		},
		"generationConfig": genConfig,
	}
	if systemPrompt != "" {
		body["systemInstruction"] = map[string]interface{}{
			"parts": []map[string]string{{"text": systemPrompt}},
		}
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.BaseURL + "/models/" + c.Model + ":generateContent"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", c.APIKey)

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

	var apiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
		} `json:"usageMetadata"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	var sb strings.Builder
	var stop string
	if len(apiResp.Candidates) > 0 {
		for _, p := range apiResp.Candidates[0].Content.Parts {
			sb.WriteString(p.Text)
		}
		stop = apiResp.Candidates[0].FinishReason
	}

	return &APIResponse{
		Text: sb.String(),
		Usage: APIUsage{
			InputTokens:  apiResp.UsageMetadata.PromptTokenCount,
			OutputTokens: apiResp.UsageMetadata.CandidatesTokenCount,
		},
		Model:      c.Model,
		StopReason: stop,
		DurationMS: duration.Milliseconds(),
	}, nil
}
