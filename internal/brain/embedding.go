package brain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// Embedder turns text into vectors. The interface is small on purpose —
// new providers (Voyage, Cohere, local models) implement Embed and slot
// in without touching ConversationStore.
type Embedder interface {
	// Embed must return one vector per input, in the same order, with a
	// stable Dim() across calls.
	Embed(ctx context.Context, texts []string) ([][]float32, error)
	Dim() int
	Name() string
}

// ── OpenAIEmbedder ─────────────────────────────────────────────────

// OpenAIEmbedder calls OpenAI's /v1/embeddings endpoint. text-embedding-3-
// small is the cheap, well-balanced choice for short chat messages —
// 1536-dimensional, $0.02 / 1M tokens.
type OpenAIEmbedder struct {
	apiKey string
	model  string
	client *http.Client
}

func NewOpenAIEmbedder(apiKey, model string) *OpenAIEmbedder {
	if model == "" {
		model = "text-embedding-3-small"
	}
	return &OpenAIEmbedder{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (e *OpenAIEmbedder) Name() string { return "openai:" + e.model }

func (e *OpenAIEmbedder) Dim() int {
	switch e.model {
	case "text-embedding-3-small":
		return 1536
	case "text-embedding-3-large":
		return 3072
	}
	return 1536
}

type openAIEmbeddingsReq struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

type openAIEmbeddingsResp struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if e == nil {
		return nil, errors.New("nil embedder")
	}
	if len(texts) == 0 {
		return nil, nil
	}
	// Clean the inputs — OpenAI rejects empty strings and chokes on very
	// long lone inputs. Truncate aggressively; embeddings of short text
	// are what we want anyway.
	clean := make([]string, len(texts))
	for i, t := range texts {
		if len(t) > 8000 {
			t = t[:8000]
		}
		if t == "" {
			t = " " // placeholder so vectors line up with caller's order
		}
		clean[i] = t
	}

	body, err := json.Marshal(openAIEmbeddingsReq{Input: clean, Model: e.model})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.openai.com/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	res, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai embed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("openai embed http %d: %s", res.StatusCode, string(b))
	}
	var out openAIEmbeddingsResp
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("openai embed parse: %w", err)
	}
	if out.Error != nil {
		return nil, fmt.Errorf("openai embed error: %s", out.Error.Message)
	}
	if len(out.Data) != len(texts) {
		return nil, fmt.Errorf("openai embed: returned %d vectors for %d inputs", len(out.Data), len(texts))
	}
	// Reorder by Index just in case the API returns them out of order.
	vecs := make([][]float32, len(texts))
	for _, d := range out.Data {
		if d.Index < 0 || d.Index >= len(vecs) {
			return nil, fmt.Errorf("openai embed: index %d out of range", d.Index)
		}
		vecs[d.Index] = d.Embedding
	}
	return vecs, nil
}

// ── Math helpers ───────────────────────────────────────────────────

// Cosine returns the cosine similarity of two equal-length vectors, in
// [-1, 1]. Returns 0 if either vector is zero-norm or sizes mismatch —
// callers treat 0 as "no signal" rather than special-casing errors.
func Cosine(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, na, nb float64
	for i, va := range a {
		vb := float64(b[i])
		fa := float64(va)
		dot += fa * vb
		na += fa * fa
		nb += vb * vb
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return float32(dot / (math.Sqrt(na) * math.Sqrt(nb)))
}
