package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"pty-claude-test/internal/capability"
)

// dataRoot is where every capability persists its per-scope JSON.
// Kept under the existing project dir so backup/restore is one folder.
const dataRoot = "data"

func scopeOf(r *http.Request) string {
	s := r.URL.Query().Get("scope")
	if s == "" {
		s = "default"
	}
	return s
}

// ── HR ────────────────────────────────────────────────────────────────

func (s *Server) handleHRApplicants(w http.ResponseWriter, r *http.Request) {
	hr, err := capability.NewHR(dataRoot, scopeOf(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	switch r.Method {
	case http.MethodGet:
		apps, err := hr.List()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"applicants": apps, "count": len(apps)})
	case http.MethodPost:
		var a capability.Applicant
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, "bad json: "+err.Error(), 400)
			return
		}
		saved, err := hr.Store(a)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, saved)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (s *Server) handleHRApplicantOne(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/cap/hr/applicants/")
	if id == "" {
		http.Error(w, "id required", 400)
		return
	}
	hr, err := capability.NewHR(dataRoot, scopeOf(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if r.Method != http.MethodPatch && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json: "+err.Error(), 400)
		return
	}
	a, err := hr.UpdateStatus(id, body.Status)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	writeJSON(w, a)
}

func (s *Server) handleHRMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	hr, err := capability.NewHR(dataRoot, scopeOf(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var jd capability.JD
	if err := json.NewDecoder(r.Body).Decode(&jd); err != nil {
		http.Error(w, "bad json: "+err.Error(), 400)
		return
	}
	matches, err := hr.Match(jd, 5)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{
		"jd":      jd,
		"matches": matches,
		"count":   len(matches),
	})
}

// ── Procurement ───────────────────────────────────────────────────────

func (s *Server) handleProcurementVendors(w http.ResponseWriter, r *http.Request) {
	p, err := capability.NewProcurement(dataRoot, scopeOf(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	switch r.Method {
	case http.MethodGet:
		vs, err := p.List()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"vendors": vs, "count": len(vs)})
	case http.MethodPost:
		var v capability.Vendor
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			http.Error(w, "bad json: "+err.Error(), 400)
			return
		}
		if v.ID == "" {
			http.Error(w, "vendor id required", 400)
			return
		}
		saved, err := p.Store(v)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, saved)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (s *Server) handleProcurementValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	p, err := capability.NewProcurement(dataRoot, scopeOf(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var body struct {
		SenderEmail string  `json:"sender_email"`
		VendorID    string  `json:"vendor_id"`
		Amount      float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json: "+err.Error(), 400)
		return
	}
	if body.SenderEmail == "" {
		http.Error(w, "sender_email required", 400)
		return
	}
	check, err := p.ValidateInvoice(body.SenderEmail, body.VendorID, body.Amount)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, check)
}

// ── Schedule ──────────────────────────────────────────────────────────

func (s *Server) handleScheduleMilestones(w http.ResponseWriter, r *http.Request) {
	sc, err := capability.NewSchedule(dataRoot, scopeOf(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	switch r.Method {
	case http.MethodGet:
		ms, err := sc.List()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"milestones": ms, "count": len(ms)})
	case http.MethodPost:
		var m capability.Milestone
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "bad json: "+err.Error(), 400)
			return
		}
		if m.Title == "" || m.DueDate == "" {
			http.Error(w, "title and due_date required", 400)
			return
		}
		saved, err := sc.Store(m)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, saved)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (s *Server) handleScheduleSlips(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	sc, err := capability.NewSchedule(dataRoot, scopeOf(r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	slips, err := sc.Slips()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"slips": slips, "count": len(slips)})
}

// ── Inbox (global; one mailbox per the production model) ──────────────

// handleInboxList dispatches GET (list) vs POST (mock-inject) on the
// /api/cap/email/inbox path so the Lab page can drop synthetic emails
// into the inbox without going through a real mail server.
func (s *Server) handleInboxList(w http.ResponseWriter, r *http.Request) {
	in := capability.NewInbox(dataRoot)
	switch r.Method {
	case http.MethodGet:
		mails, err := in.List()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"emails": mails, "count": len(mails)})
	case http.MethodPost:
		var e capability.Email
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, "bad json: "+err.Error(), 400)
			return
		}
		saved, err := in.Inject(e)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, saved)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (s *Server) handleInboxSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	in := capability.NewInbox(dataRoot)
	sum, err := in.Summary()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, sum)
}

// ── Capability index (UI + docs) ──────────────────────────────────────

func (s *Server) handleCapabilities(w http.ResponseWriter, r *http.Request) {
	// Static descriptor — what each agent declares it can do. The HTTP
	// endpoints listed here are the real, tested surface.
	writeJSON(w, map[string]any{
		"agents": []map[string]any{
			{
				"role":  "hr",
				"name":  "HR Lead",
				"description": "Job applicant pipeline + HRMS integration.",
				"capabilities": []map[string]any{
					{"name": "list_applicants", "method": "GET", "path": "/api/cap/hr/applicants"},
					{"name": "store_applicant", "method": "POST", "path": "/api/cap/hr/applicants"},
					{"name": "match_against_jd", "method": "POST", "path": "/api/cap/hr/match"},
					{"name": "update_status", "method": "PATCH", "path": "/api/cap/hr/applicants/{id}"},
				},
			},
			{
				"role":  "procurement",
				"name":  "Procurement Manager",
				"description": "Vendor registry + invoice validation against BEC patterns.",
				"capabilities": []map[string]any{
					{"name": "list_vendors", "method": "GET", "path": "/api/cap/procurement/vendors"},
					{"name": "store_vendor", "method": "POST", "path": "/api/cap/procurement/vendors"},
					{"name": "validate_invoice", "method": "POST", "path": "/api/cap/procurement/validate-invoice"},
				},
			},
			{
				"role":  "project_manager",
				"name":  "Project Manager",
				"description": "Schedule milestones + slip detection.",
				"capabilities": []map[string]any{
					{"name": "list_milestones", "method": "GET", "path": "/api/cap/schedule/milestones"},
					{"name": "store_milestone", "method": "POST", "path": "/api/cap/schedule/milestones"},
					{"name": "detect_slips", "method": "GET", "path": "/api/cap/schedule/slips"},
				},
			},
			{
				"role":  "conductor",
				"name":  "Conductor",
				"description": "Global inbox sweep + briefing composition.",
				"capabilities": []map[string]any{
					{"name": "list_inbox", "method": "GET", "path": "/api/cap/email/inbox"},
					{"name": "inbox_summary", "method": "GET", "path": "/api/cap/email/summary"},
				},
			},
		},
	})
}

// notFoundError is used to translate capability errors into 404s, but isn't
// exposed yet; placeholder for future error mapping.
var _ = errors.New
