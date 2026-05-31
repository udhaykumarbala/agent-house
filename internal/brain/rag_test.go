package brain

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestCosine(t *testing.T) {
	cases := []struct {
		name string
		a, b []float32
		want float32
	}{
		{"identical unit vectors", []float32{1, 0, 0}, []float32{1, 0, 0}, 1},
		{"orthogonal", []float32{1, 0, 0}, []float32{0, 1, 0}, 0},
		{"opposite", []float32{1, 0, 0}, []float32{-1, 0, 0}, -1},
		{"empty", []float32{}, []float32{}, 0},
		{"length mismatch", []float32{1, 0}, []float32{1, 0, 0}, 0},
		{"zero norm", []float32{0, 0, 0}, []float32{1, 0, 0}, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Cosine(c.a, c.b)
			if math.Abs(float64(got-c.want)) > 1e-5 {
				t.Errorf("Cosine(%v, %v) = %v, want %v", c.a, c.b, got, c.want)
			}
		})
	}
}

// SearchWithMode must always return mode="lexical" when no embedder is
// attached, never panic, and pass through to the substring search.
func TestSearchWithMode_LexicalFallback(t *testing.T) {
	dir := t.TempDir()
	convDir := filepath.Join(dir, "conversations")
	if err := os.MkdirAll(convDir, 0o755); err != nil {
		t.Fatal(err)
	}
	store := NewConversationStore(convDir)
	store.Add("c1", "user", "What's the impersonation status?")
	store.Add("c1", "assistant", "Two impersonation emails were flagged.")

	res := store.SearchWithMode(context.Background(), "impersonation", 10)
	if res.Mode != "lexical" {
		t.Errorf("mode = %q, want %q (no embedder attached)", res.Mode, "lexical")
	}
	if len(res.Hits) == 0 {
		t.Errorf("lexical fallback returned 0 hits; expected ≥1")
	}
}

// fakeEmbedder lets us prove the semantic path is taken without hitting
// the real OpenAI API. Vectors are tiny but deterministic.
type fakeEmbedder struct {
	calls [][]string
}

func (f *fakeEmbedder) Name() string { return "fake" }
func (f *fakeEmbedder) Dim() int     { return 3 }
func (f *fakeEmbedder) Embed(_ context.Context, texts []string) ([][]float32, error) {
	f.calls = append(f.calls, texts)
	out := make([][]float32, len(texts))
	for i, t := range texts {
		// Map presence of a few tokens to one-hot-ish vectors so cosine
		// rankings are predictable.
		v := []float32{0, 0, 0}
		switch {
		case containsCI(t, "impersonation"):
			v = []float32{1, 0, 0}
		case containsCI(t, "schedule"):
			v = []float32{0, 1, 0}
		case containsCI(t, "vendor"):
			v = []float32{0.7, 0.7, 0} // close to impersonation, not identical
		default:
			v = []float32{0, 0, 1}
		}
		out[i] = v
	}
	return out, nil
}

func containsCI(s, sub string) bool {
	return len(s) >= len(sub) && (indexFold(s, sub) >= 0)
}

func indexFold(s, sub string) int {
	// Tiny case-insensitive Index. Avoids depending on package layout in this test.
	if sub == "" {
		return 0
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		match := true
		for j := 0; j < len(sub); j++ {
			a, b := s[i+j], sub[j]
			if a >= 'A' && a <= 'Z' {
				a += 32
			}
			if b >= 'A' && b <= 'Z' {
				b += 32
			}
			if a != b {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// Multiple matches inside one conversation should collapse to a single
// row, with MatchCount reflecting the total. Lexical path.
func TestSearchWithMode_LexicalDedupesByConversation(t *testing.T) {
	dir := t.TempDir()
	convDir := filepath.Join(dir, "conversations")
	if err := os.MkdirAll(convDir, 0o755); err != nil {
		t.Fatal(err)
	}
	store := NewConversationStore(convDir)
	// All four matches live in "atlas".
	store.Add("atlas", "user", "whats pending")
	store.Add("atlas", "assistant", "Here is whats pending across your workspace.")
	store.Add("atlas", "user", "anything else pending?")
	store.Add("atlas", "assistant", "Two pending impersonation flags remain.")
	// One match in a different conversation.
	store.Add("other", "user", "any pending invoices?")

	res := store.SearchWithMode(context.Background(), "pending", 10)
	if res.Mode != "lexical" {
		t.Fatalf("mode = %q", res.Mode)
	}
	if len(res.Hits) != 2 {
		t.Fatalf("expected 2 deduped hits (one per conversation), got %d: %+v",
			len(res.Hits), res.Hits)
	}
	// Find the atlas hit and check MatchCount.
	var atlas *SearchHit
	for i := range res.Hits {
		if res.Hits[i].ConversationID == "atlas" {
			atlas = &res.Hits[i]
		}
	}
	if atlas == nil {
		t.Fatalf("no hit for atlas conversation in %+v", res.Hits)
	}
	if atlas.MatchCount != 4 {
		t.Errorf("atlas MatchCount = %d, want 4", atlas.MatchCount)
	}
}

// With a fake embedder, the same impersonation query should rank the
// "impersonation" message above the "schedule" message, and the search
// should be flagged as semantic.
func TestSearchWithMode_Semantic(t *testing.T) {
	dir := t.TempDir()
	convDir := filepath.Join(dir, "conversations")
	if err := os.MkdirAll(convDir, 0o755); err != nil {
		t.Fatal(err)
	}
	store := NewConversationStore(convDir)
	store.SetEmbedder(&fakeEmbedder{})

	store.Add("ops", "user", "Anything about the schedule today?")
	store.Add("ops", "assistant", "Foundation work is 53 days late.")
	store.Add("sec", "user", "Any impersonation attempts?")
	store.Add("sec", "assistant", "Yes — two emails flagged as impersonation from gmail.")

	res := store.SearchWithMode(context.Background(), "impersonation", 5)
	if res.Mode != "semantic" {
		t.Fatalf("mode = %q, want %q", res.Mode, "semantic")
	}
	if len(res.Hits) == 0 {
		t.Fatalf("semantic search returned 0 hits")
	}
	// Top hit should come from the "sec" conversation (impersonation match).
	if res.Hits[0].ConversationID != "sec" {
		t.Errorf("top hit conv = %q, want %q",
			res.Hits[0].ConversationID, "sec")
	}
}
