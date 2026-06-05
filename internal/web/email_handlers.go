package web

import (
	"encoding/json"
	"net/http"
	"strings"

	emailpkg "pty-claude-test/internal/email"
)

// EmailHandlers manages email-related HTTP endpoints.
type EmailHandlers struct {
	engine *emailpkg.Engine
	resend *emailpkg.ResendClient
}

// NewEmailHandlers creates email handlers.
func NewEmailHandlers(dataDir string) *EmailHandlers {
	return &EmailHandlers{
		engine: emailpkg.NewEngine(dataDir),
		resend: emailpkg.NewResendClient(),
	}
}

// GetEngine returns the email engine (for wiring to Brain).
func (eh *EmailHandlers) GetEngine() *emailpkg.Engine {
	return eh.engine
}

// HandleInbox handles GET /api/email/inbox
func (eh *EmailHandlers) HandleInbox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for specific email ID in path
	path := strings.TrimPrefix(r.URL.Path, "/api/email/inbox")
	if len(path) > 1 {
		id := strings.TrimPrefix(path, "/")
		email := eh.engine.GetEmail(id)
		if email == nil {
			http.Error(w, "Email not found", http.StatusNotFound)
			return
		}
		writeJSON(w, email)
		return
	}

	emails := eh.engine.GetInbox()
	writeJSON(w, map[string]interface{}{
		"emails": emails,
		"count":  len(emails),
	})
}

// HandleReceive handles POST /api/email/receive
func (eh *EmailHandlers) HandleReceive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		From     string `json:"from"`
		FromName string `json:"from_name"`
		To       string `json:"to"`
		Subject  string `json:"subject"`
		Body     string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if req.From == "" || req.Subject == "" {
		http.Error(w, "from and subject required", http.StatusBadRequest)
		return
	}
	if req.To == "" {
		req.To = "projects@epc-demo.com"
	}

	email, trust := eh.engine.ReceiveEmail(req.From, req.FromName, req.To, req.Subject, req.Body)

	// Gateway dropped the message — return 403 with the reason so the
	// caller (simulator, UI, or real intake) can surface the alert
	// instead of pretending the email was accepted.
	if trust != nil && trust.Status == "blocked" {
		writeJSON(w, map[string]interface{}{
			"success": false,
			"blocked": true,
			"trust":   trust,
			"error":   trust.Reason,
		})
		return
	}

	result := map[string]interface{}{
		"success":  true,
		"email":    email,
		"trust":    trust,
	}

	// Add alert if trust issue
	if trust.Status == "impersonation" {
		result["alert"] = map[string]interface{}{
			"level":        "critical",
			"message":      trust.Reason,
			"risk_factors": trust.RiskFactors,
			"action":       "Verify vendor identity before processing",
		}
	} else if trust.Status == "new_contact" {
		result["alert"] = map[string]interface{}{
			"level":   "warning",
			"message": trust.Reason,
			"action":  "Review and add to trusted list if legitimate",
		}
	}

	writeJSON(w, result)
}

// HandleSend handles POST /api/email/send
func (eh *EmailHandlers) HandleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if req.To == "" || req.Subject == "" {
		http.Error(w, "to and subject required", http.StatusBadRequest)
		return
	}
	if req.From == "" {
		req.From = "projects@epc-demo.com"
	}

	err := eh.resend.Send(req.From, req.To, req.Subject, req.Body)
	if err != nil {
		writeJSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"note":    "Email saved as draft. Set RESEND_API_KEY to send real emails.",
		})
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
		"message": "Email sent via Resend",
	})
}

// HandleApplicants handles GET /api/email/applicants
func (eh *EmailHandlers) HandleApplicants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	applicants := eh.engine.GetApplicants()
	writeJSON(w, map[string]interface{}{
		"applicants": applicants,
		"count":      len(applicants),
	})
}

// HandleTrust handles POST /api/email/trust — add email to vendor trusted list
func (eh *EmailHandlers) HandleTrust(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		VendorID string `json:"vendor_id"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if eh.engine.AddTrustedEmail(req.VendorID, req.Email) {
		writeJSON(w, map[string]string{"status": "added"})
	} else {
		writeJSON(w, map[string]string{"status": "vendor not found"})
	}
}

// handleEmail is the Server method that routes to EmailHandlers.
func (s *Server) handleEmail(w http.ResponseWriter, r *http.Request) {
	if s.emailHandlers == nil {
		http.Error(w, "Email not initialized", http.StatusServiceUnavailable)
		return
	}
	s.emailHandlers.HandleEmailRouting(w, r)
}

// HandleEmailRouting routes /api/email/* requests.
func (eh *EmailHandlers) HandleEmailRouting(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case strings.HasPrefix(path, "/api/email/inbox"):
		eh.HandleInbox(w, r)
	case path == "/api/email/receive":
		eh.HandleReceive(w, r)
	case path == "/api/email/send":
		eh.HandleSend(w, r)
	case path == "/api/email/applicants":
		eh.HandleApplicants(w, r)
	case path == "/api/email/trust":
		eh.HandleTrust(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
