package brain

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// ConversationMeta is the lightweight summary used to populate the
// conversation list / command-palette without loading every full thread.
type ConversationMeta struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	MessageCount int    `json:"message_count"`
	LastAt       int64  `json:"last_at"`       // ms since epoch; 0 if empty
	FirstAt      int64  `json:"first_at"`      // ms since epoch; 0 if empty
	Preview      string `json:"preview"`       // first user message, truncated
}

// SearchHit is one match — but the palette wants conversation-level rows,
// so the API contract is "one hit per conversation" with MatchCount
// telling the UI how many messages matched inside it. Snippet + Role +
// Timestamp belong to the best-ranked matching message in that conv.
type SearchHit struct {
	ConversationID string `json:"conversation_id"`
	Name           string `json:"name"`
	Role           string `json:"role"`        // best-matching message's role
	Snippet        string `json:"snippet"`     // best-matching message's snippet
	Timestamp      int64  `json:"timestamp"`   // best-matching message's ts
	MatchCount     int    `json:"match_count"` // total matching messages in this conv
}

// ConversationStore manages chat history per user.
//
// When an Embedder is attached, the store persists per-message embeddings
// alongside (one file per conversation, mirroring the messages file) and
// upgrades Search() from substring matching to true semantic ranking.
// Without an Embedder it falls back to lexical search — same shape,
// honest about the mode in the response.
type ConversationStore struct {
	mu       sync.Mutex
	dir      string
	embDir   string
	messages map[string][]ChatMessage // userID → messages
	embed    map[string][][]float32   // userID → per-message vectors (len matches messages)
	embedder Embedder                 // nil = lexical-only mode
}

// NewConversationStore creates a store backed by the filesystem. Embeddings
// (when an Embedder is attached later via SetEmbedder) live in a sibling
// `embeddings/` directory under the same root.
func NewConversationStore(dir string) *ConversationStore {
	os.MkdirAll(dir, 0755)
	embDir := filepath.Join(filepath.Dir(dir), "embeddings")
	os.MkdirAll(embDir, 0755)
	return &ConversationStore{
		dir:      dir,
		embDir:   embDir,
		messages: make(map[string][]ChatMessage),
		embed:    make(map[string][][]float32),
	}
}

// SetEmbedder attaches a vector model. Called from server bootstrap when
// OPENAI_API_KEY is present.
func (cs *ConversationStore) SetEmbedder(e Embedder) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.embedder = e
	if e != nil {
		log.Printf("[brain] semantic search ENABLED (%s)", e.Name())
	}
}

// Add appends a message to a user's conversation.
func (cs *ConversationStore) Add(userID, role, content string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	msgs := cs.messages[userID]
	msgs = append(msgs, ChatMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now().UnixMilli(),
	})

	// Keep bounded
	if len(msgs) > 50 {
		msgs = msgs[len(msgs)-30:]
	}

	cs.messages[userID] = msgs
	cs.save(userID)
}

// GetRecent returns the last N messages for a user.
func (cs *ConversationStore) GetRecent(userID string, n int) []ChatMessage {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Load from disk if not in memory
	if _, ok := cs.messages[userID]; !ok {
		cs.load(userID)
	}

	msgs := cs.messages[userID]
	if len(msgs) <= n {
		return msgs
	}
	return msgs[len(msgs)-n:]
}

// GetAll returns all messages for a user.
func (cs *ConversationStore) GetAll(userID string) []ChatMessage {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, ok := cs.messages[userID]; !ok {
		cs.load(userID)
	}
	return cs.messages[userID]
}

func (cs *ConversationStore) save(userID string) {
	data, err := json.MarshalIndent(cs.messages[userID], "", "  ")
	if err != nil {
		return
	}
	path := filepath.Join(cs.dir, userID+".json")
	os.WriteFile(path, data, 0644)
}

func (cs *ConversationStore) load(userID string) {
	path := filepath.Join(cs.dir, userID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		cs.messages[userID] = []ChatMessage{}
		return
	}
	var msgs []ChatMessage
	if json.Unmarshal(data, &msgs) == nil {
		cs.messages[userID] = msgs
	}
}

// List scans the conversation directory and returns lightweight metadata
// for every conversation (one file = one conversation). Sorted most-
// recently-active first.
func (cs *ConversationStore) List() []ConversationMeta {
	entries, err := os.ReadDir(cs.dir)
	if err != nil {
		return nil
	}
	out := make([]ConversationMeta, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		out = append(out, cs.metaFor(id))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LastAt > out[j].LastAt })
	return out
}

// metaFor builds metadata for a single conversation, lazily loading it.
func (cs *ConversationStore) metaFor(id string) ConversationMeta {
	msgs := cs.GetAll(id)
	meta := ConversationMeta{
		ID:           id,
		Name:         deriveName(id, msgs),
		MessageCount: len(msgs),
	}
	if len(msgs) > 0 {
		meta.FirstAt = msgs[0].Timestamp
		meta.LastAt = msgs[len(msgs)-1].Timestamp
		for _, m := range msgs {
			if m.Role == "user" {
				meta.Preview = truncateRunes(m.Content, 140)
				break
			}
		}
	}
	return meta
}

// deriveName produces a human-readable name from the first user message,
// or falls back to the id when there's no content yet.
func deriveName(id string, msgs []ChatMessage) string {
	for _, m := range msgs {
		if m.Role != "user" {
			continue
		}
		// Use the first non-empty user message, truncated to a headline.
		head := strings.TrimSpace(m.Content)
		if head == "" {
			continue
		}
		return truncateRunes(firstLine(head), 60)
	}
	// id like "conv_<ms>" → "Conversation <last-4>"
	if strings.HasPrefix(id, "conv_") {
		return "New conversation"
	}
	return id
}

func firstLine(s string) string {
	if i := strings.IndexAny(s, "\n\r"); i > 0 {
		return s[:i]
	}
	return s
}

func truncateRunes(s string, n int) string {
	s = strings.TrimSpace(s)
	if len([]rune(s)) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n-1]) + "…"
}

// SearchResult bundles the ranked hits with the mode used so the UI can
// honestly say whether semantic or lexical answered.
type SearchResult struct {
	Mode string      // "semantic" | "lexical"
	Hits []SearchHit
}

// Search dispatches between semantic (embedder-backed cosine ranking) and
// lexical (substring) based on whether an Embedder is attached. Lexical
// is the legacy entrypoint preserved for callers/tests that don't need
// the mode tag.
func (cs *ConversationStore) Search(q string, limit int) []SearchHit {
	return cs.SearchWithMode(context.Background(), q, limit).Hits
}

// SearchWithMode is the new full-featured search. Caller gets back which
// engine actually answered.
func (cs *ConversationStore) SearchWithMode(ctx context.Context, q string, limit int) SearchResult {
	q = strings.TrimSpace(q)
	if q == "" {
		return SearchResult{Mode: "lexical", Hits: nil}
	}
	if cs.embedder != nil {
		hits, err := cs.semanticSearch(ctx, q, limit)
		if err != nil {
			log.Printf("[brain] semantic search failed, falling back to lexical: %v", err)
		} else {
			return SearchResult{Mode: "semantic", Hits: hits}
		}
	}
	return SearchResult{Mode: "lexical", Hits: cs.lexicalSearch(q, limit)}
}

// lexicalSearch is the original O(N) substring matcher, kept for callers
// without an embedder and as the safety net when embedding fails.
func (cs *ConversationStore) lexicalSearch(q string, limit int) []SearchHit {
	if limit <= 0 {
		limit = 20
	}
	q = strings.TrimSpace(q)
	if q == "" {
		return nil
	}
	low := strings.ToLower(q)
	hits := []SearchHit{}

	entries, err := os.ReadDir(cs.dir)
	if err != nil {
		return nil
	}
	// Oversample (limit * 5) so dedupe can collapse multi-hit conversations
	// without starving the result list. Bounded so a pathological corpus
	// can't blow up memory.
	hardCap := limit * 5
	if hardCap < limit+20 {
		hardCap = limit + 20
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		msgs := cs.GetAll(id)
		name := deriveName(id, msgs)
		for _, m := range msgs {
			idx := strings.Index(strings.ToLower(m.Content), low)
			if idx < 0 {
				continue
			}
			hits = append(hits, SearchHit{
				ConversationID: id,
				Name:           name,
				Role:           m.Role,
				Snippet:        snippetAround(m.Content, idx, len(q), 60),
				Timestamp:      m.Timestamp,
			})
			if len(hits) >= hardCap {
				break
			}
		}
		if len(hits) >= hardCap {
			break
		}
	}
	// Sort most-recent-first so dedupe keeps the freshest match per conv.
	sort.Slice(hits, func(i, j int) bool { return hits[i].Timestamp > hits[j].Timestamp })
	hits = dedupeByConversation(hits)
	if len(hits) > limit {
		hits = hits[:limit]
	}
	return hits
}

// dedupeByConversation collapses a ranked hit list down to one row per
// conversation, keeping the first occurrence (which both search paths
// arrange to be the best-ranked one) and aggregating MatchCount.
func dedupeByConversation(hits []SearchHit) []SearchHit {
	if len(hits) == 0 {
		return hits
	}
	idx := map[string]int{} // conv id → index in out
	out := make([]SearchHit, 0, len(hits))
	for _, h := range hits {
		if pos, seen := idx[h.ConversationID]; seen {
			out[pos].MatchCount++
			continue
		}
		h.MatchCount = 1
		out = append(out, h)
		idx[h.ConversationID] = len(out) - 1
	}
	return out
}

func snippetAround(s string, idx, matchLen, window int) string {
	start := idx - window
	if start < 0 {
		start = 0
	}
	end := idx + matchLen + window
	if end > len(s) {
		end = len(s)
	}
	out := s[start:end]
	if start > 0 {
		out = "…" + out
	}
	if end < len(s) {
		out += "…"
	}
	return strings.ReplaceAll(out, "\n", " ")
}

// ── Semantic search (embedder-backed) ───────────────────────────────

// semanticSearch ensures every stored message has an embedding (lazy
// backfill on first call per conversation, then incremental), embeds the
// query, and returns the top-`limit` messages ranked by cosine similarity.
func (cs *ConversationStore) semanticSearch(ctx context.Context, q string, limit int) ([]SearchHit, error) {
	if limit <= 0 {
		limit = 20
	}
	// Discover every conversation on disk.
	entries, err := os.ReadDir(cs.dir)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		ids = append(ids, strings.TrimSuffix(e.Name(), ".json"))
	}
	// Backfill embeddings for any conversation that's missing some.
	for _, id := range ids {
		if err := cs.ensureEmbeddings(ctx, id); err != nil {
			return nil, err
		}
	}
	// Embed the query.
	qVecs, err := cs.embedder.Embed(ctx, []string{q})
	if err != nil || len(qVecs) != 1 {
		return nil, err
	}
	qv := qVecs[0]

	type scored struct {
		hit   SearchHit
		score float32
	}
	all := []scored{}
	cs.mu.Lock()
	for _, id := range ids {
		msgs := cs.messages[id]
		vecs := cs.embed[id]
		if len(vecs) != len(msgs) {
			continue // mid-update; will catch up next call
		}
		name := deriveName(id, msgs)
		for i, m := range msgs {
			s := Cosine(qv, vecs[i])
			if s <= 0 {
				continue
			}
			all = append(all, scored{
				hit: SearchHit{
					ConversationID: id,
					Name:           name,
					Role:           m.Role,
					Snippet:        snippetForSemantic(m.Content),
					Timestamp:      m.Timestamp,
				},
				score: s,
			})
		}
	}
	cs.mu.Unlock()
	// Rank by similarity, dedupe to one row per conversation (keeps the
	// highest-scoring message because order is preserved), then truncate.
	sort.Slice(all, func(i, j int) bool { return all[i].score > all[j].score })
	hits := make([]SearchHit, len(all))
	for i, s := range all {
		hits[i] = s.hit
	}
	hits = dedupeByConversation(hits)
	if len(hits) > limit {
		hits = hits[:limit]
	}
	return hits, nil
}

// snippetForSemantic returns a leading window of the message — semantic
// hits aren't anchored to a specific phrase, so showing the head of the
// content is the most informative preview.
func snippetForSemantic(content string) string {
	content = strings.ReplaceAll(content, "\n", " ")
	if len([]rune(content)) <= 140 {
		return content
	}
	r := []rune(content)
	return string(r[:139]) + "…"
}

// ensureEmbeddings makes sure every message in convID has a vector. The
// embeddings file is sibling to the conversation file (`<embDir>/<id>.json`)
// and holds a `[][]float32` aligned with messages by index.
func (cs *ConversationStore) ensureEmbeddings(ctx context.Context, convID string) error {
	if cs.embedder == nil {
		return nil
	}
	cs.mu.Lock()
	msgs := append([]ChatMessage(nil), cs.messages[convID]...)
	cur := append([][]float32(nil), cs.embed[convID]...)
	cs.mu.Unlock()

	// Lazy first-load for both messages and vectors.
	if msgs == nil {
		cs.mu.Lock()
		if _, ok := cs.messages[convID]; !ok {
			cs.load(convID)
		}
		msgs = append([]ChatMessage(nil), cs.messages[convID]...)
		cs.mu.Unlock()
	}
	if cur == nil {
		cur = cs.loadEmbeddings(convID)
	}

	// Figure out which messages need embedding.
	missing := []int{}
	missingTexts := []string{}
	for i, m := range msgs {
		var v []float32
		if i < len(cur) {
			v = cur[i]
		}
		if len(v) == 0 {
			missing = append(missing, i)
			missingTexts = append(missingTexts, m.Content)
		}
	}
	if len(missing) == 0 {
		// Cache + return.
		cs.mu.Lock()
		cs.embed[convID] = cur
		cs.mu.Unlock()
		return nil
	}
	vecs, err := cs.embedder.Embed(ctx, missingTexts)
	if err != nil {
		return err
	}
	if len(vecs) != len(missing) {
		return nil // anomalous; skip rather than corrupt
	}
	// Grow cur to len(msgs) and fill in.
	if len(cur) < len(msgs) {
		grow := make([][]float32, len(msgs))
		copy(grow, cur)
		cur = grow
	}
	for i, idx := range missing {
		cur[idx] = vecs[i]
	}
	cs.mu.Lock()
	cs.embed[convID] = cur
	cs.mu.Unlock()
	cs.saveEmbeddings(convID, cur)
	return nil
}

func (cs *ConversationStore) loadEmbeddings(convID string) [][]float32 {
	path := filepath.Join(cs.embDir, convID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var out [][]float32
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	return out
}

func (cs *ConversationStore) saveEmbeddings(convID string, vecs [][]float32) {
	data, err := json.Marshal(vecs)
	if err != nil {
		return
	}
	path := filepath.Join(cs.embDir, convID+".json")
	os.WriteFile(path, data, 0o644)
}
