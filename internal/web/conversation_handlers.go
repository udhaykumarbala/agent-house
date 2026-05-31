package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// handleConversationsRouter dispatches GET (list) vs POST (create) on the
// /api/conversations exact path.
func (s *Server) handleConversationsRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleConversations(w, r)
	case http.MethodPost:
		s.handleConversationCreate(w, r)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

// handleConversations lists every conversation with lightweight metadata
// for the command-palette / history picker. GET only.
func (s *Server) handleConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	if s.brainHandler == nil || s.brainHandler.convos == nil {
		writeJSON(w, map[string]any{"conversations": []any{}})
		return
	}
	writeJSON(w, map[string]any{
		"conversations": s.brainHandler.convos.List(),
	})
}

// handleConversationCreate accepts an empty POST and returns a fresh,
// stable conversation id the frontend can use immediately. We don't
// persist anything until the first chat lands.
func (s *Server) handleConversationCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	id := fmt.Sprintf("conv_%d", time.Now().UnixMilli())
	writeJSON(w, map[string]any{
		"id":   id,
		"name": "New conversation",
	})
}

// handleConversationSearch performs semantic search across every stored
// conversation when an Embedder is attached (OPENAI_API_KEY set), and
// falls back to substring matching otherwise. The response surfaces
// `mode` so the UI can honestly say which engine answered.
func (s *Server) handleConversationSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		writeJSON(w, map[string]any{"mode": "lexical", "hits": []any{}})
		return
	}
	if s.brainHandler == nil || s.brainHandler.convos == nil {
		writeJSON(w, map[string]any{"mode": "lexical", "hits": []any{}})
		return
	}
	result := s.brainHandler.convos.SearchWithMode(r.Context(), q, 30)
	writeJSON(w, map[string]any{
		"q":    q,
		"mode": result.Mode,
		"hits": result.Hits,
	})
}

// handleConversationOne returns the full message thread for a single id —
// alias for the existing GET /api/chat?user=<id> but cleaner verb.
func (s *Server) handleConversationOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/conversations/")
	if id == "" || strings.Contains(id, "/") {
		http.Error(w, "conversation id required: GET /api/conversations/{id}", 400)
		return
	}
	if s.brainHandler == nil || s.brainHandler.convos == nil {
		writeJSON(w, map[string]any{"messages": []any{}})
		return
	}
	writeJSON(w, map[string]any{
		"id":       id,
		"messages": s.brainHandler.convos.GetAll(id),
	})
}

// _ keeps imports used even if a handler is later commented out.
var _ = json.Marshal
