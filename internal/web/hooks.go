package web

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/orchestrator"
)

// WebhookConfig defines a registered webhook endpoint.
type WebhookConfig struct {
	ID        string `json:"id"`
	Secret    string `json:"secret"`     // HMAC-SHA256 secret (optional)
	AgentRole string `json:"agent_role"` // default target agent
	ProjectID string `json:"project_id"` // target project
	Template  string `json:"template"`   // task template with {{.payload}} placeholder
}

// handleWebhook handles POST /api/hooks/{hook_id}
// Converts incoming webhook payloads into injected tasks.
func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hookID := strings.TrimPrefix(r.URL.Path, "/api/hooks/")
	if hookID == "" {
		http.Error(w, "Hook ID required", http.StatusBadRequest)
		return
	}

	// Read body
	body, err := io.ReadAll(io.LimitReader(r.Body, 1*1024*1024)) // 1MB limit
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Load hook config
	hook := loadHookConfig(hookID, s.projectDir)
	if hook == nil {
		// No config — create a generic task from the webhook
		hook = &WebhookConfig{
			ID:        hookID,
			AgentRole: "senior_dev",
			Template:  "Webhook received ({{.hook_id}}):\n\n{{.payload}}",
		}
	}

	// Verify HMAC if secret is configured
	if hook.Secret != "" {
		sig := r.Header.Get("X-Signature-256")
		if sig == "" {
			sig = r.Header.Get("X-Hub-Signature-256")
		}
		if !verifyHMAC(body, hook.Secret, sig) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Build task from template
	var payload interface{}
	json.Unmarshal(body, &payload)

	taskText := hook.Template
	taskText = strings.ReplaceAll(taskText, "{{.payload}}", string(body))
	taskText = strings.ReplaceAll(taskText, "{{.hook_id}}", hookID)

	// Extract common webhook fields
	if m, ok := payload.(map[string]interface{}); ok {
		if action, ok := m["action"].(string); ok {
			taskText = strings.ReplaceAll(taskText, "{{.action}}", action)
		}
	}

	projectID := hook.ProjectID
	if projectID == "" {
		projectID = r.URL.Query().Get("project")
	}
	if projectID == "" {
		projectID = "default"
	}

	injected := &orchestrator.InjectedTask{
		ID:        fmt.Sprintf("hook_%s_%d", hookID, time.Now().UnixMilli()),
		Task:      taskText,
		AgentRole: agent.Role(hook.AgentRole),
		ProjectID: projectID,
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	s.orchestrator.InjectTask(injected)

	log.Printf("[WEBHOOK] %s → injected task for %s in %s", hookID, hook.AgentRole, projectID)

	writeJSON(w, map[string]interface{}{
		"success":   true,
		"hook_id":   hookID,
		"inject_id": injected.ID,
		"agent":     hook.AgentRole,
		"project":   projectID,
	})
}

// loadHookConfig loads webhook config from .agent-house/hooks/{id}.json
func loadHookConfig(hookID, projectDir string) *WebhookConfig {
	// Try project-level config first
	path := projectDir + "/.agent-house/hooks/" + hookID + ".json"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var config WebhookConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil
	}
	config.ID = hookID
	return &config
}

func verifyHMAC(body []byte, secret, signature string) bool {
	if signature == "" {
		return false
	}
	// Remove "sha256=" prefix if present
	signature = strings.TrimPrefix(signature, "sha256=")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
