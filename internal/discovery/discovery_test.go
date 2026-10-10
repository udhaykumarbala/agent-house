package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pty-claude-test/internal/message"
)

// fakeLLM records every call and returns a canned rich-markdown doc.
type fakeLLM struct {
	calls   []struct{ system, user string }
	failFor string // substring of the user prompt that triggers an error
}

func (f *fakeLLM) Generate(_ context.Context, system, user string) (string, error) {
	f.calls = append(f.calls, struct{ system, user string }{system, user})
	if f.failFor != "" && strings.Contains(user, f.failFor) {
		return "", errors.New("llm boom")
	}
	return "# Generated Doc\n\n## Goals\n\n```mermaid\nflowchart TD\n  A-->B\n```\n\n## Scope Table\n\n| k | v |\n|---|---|\n| a | 1 |\n", nil
}

// recordingEmitter captures everything the pipeline broadcasts.
type recordingEmitter struct{ notes []string }

func (r *recordingEmitter) Emit(m *message.Message) { r.notes = append(r.notes, m.Content) }

func newRunner(t *testing.T, llm LLM) *Runner {
	t.Helper()
	return &Runner{
		ProjectsDir: t.TempDir(),
		LLM:         llm,
		Prompts:     func(role string) string { return "You are the " + role + " agent." },
	}
}

func TestDiscoveryRunWritesAllDocs(t *testing.T) {
	llm := &fakeLLM{}
	r := newRunner(t, llm)

	man, err := r.Run(context.Background(), RunInput{
		ProjectID: "fitflow-app",
		Client:    "cto@fitflow.app",
		Request:   "Build a mobile fitness tracking app with social challenges and wearable sync.",
	})
	if err != nil {
		t.Fatal(err)
	}
	if man.Status != StatusDone {
		t.Fatalf("expected done, got %q (%#v)", man.Status, man.Steps)
	}

	dir := filepath.Join(r.ProjectsDir, "fitflow-app", "discovery")
	for _, f := range []string{
		"00-request.md", "01-prd.md", "02-competitor-research.md",
		"03-compliance-research.md", "04-architecture.md", "05-client-reply.md",
		"manifest.json",
	} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Fatalf("missing artifact %s: %v", f, err)
		}
	}

	// 5 LLM calls: prd, competitor, compliance, architecture, reply.
	if len(llm.calls) != 5 {
		t.Fatalf("expected 5 agent calls, got %d", len(llm.calls))
	}
	// Doc prompts must demand rich output and carry the client request.
	for i, c := range llm.calls[:4] {
		if !strings.Contains(strings.ToLower(c.user), "mermaid") {
			t.Fatalf("call %d prompt does not demand a mermaid diagram", i)
		}
		if !strings.Contains(c.user, "fitness tracking") {
			t.Fatalf("call %d prompt lost the client request", i)
		}
	}
	// Steps attributed to the right agents, in story order.
	wantAgents := []string{"conductor", "pm", "senior_dev", "security", "architect", "pm"}
	if len(man.Steps) != len(wantAgents) {
		t.Fatalf("expected %d steps, got %d: %#v", len(wantAgents), len(man.Steps), man.Steps)
	}
	for i, w := range wantAgents {
		if man.Steps[i].Agent != w || !man.Steps[i].OK {
			t.Fatalf("step %d: want agent %s OK, got %#v", i, w, man.Steps[i])
		}
	}

	// Reply draft also lands in the outbox for the email loop.
	entries, _ := os.ReadDir(filepath.Join(r.ProjectsDir, "outbox"))
	if len(entries) != 1 {
		t.Fatalf("expected 1 outbox draft, got %d", len(entries))
	}
	raw, _ := os.ReadFile(filepath.Join(r.ProjectsDir, "outbox", entries[0].Name()))
	var draft map[string]any
	_ = json.Unmarshal(raw, &draft)
	if draft["to"] != "cto@fitflow.app" || draft["draft"] != true {
		t.Fatalf("outbox draft malformed: %s", raw)
	}
}

func TestDiscoveryPicksNewestRequirementsEmail(t *testing.T) {
	llm := &fakeLLM{}
	r := newRunner(t, llm)
	inboxDir := filepath.Join(r.ProjectsDir, "inbox")
	_ = os.MkdirAll(inboxDir, 0o755)
	for i, e := range []map[string]any{
		{"id": "email_old_vendor", "from": "ar@vendorxyz.com", "subject": "Invoice #42", "body": "please pay", "category": "vendor", "date": "2026-10-01T10:00:00Z"},
		{"id": "email_req", "from": "cto@fitflow.app", "subject": "Requirements — FitFlow fitness app", "body": "We need you to build a mobile fitness tracking app. Requirements attached.", "category": "client", "mailbox": "software", "date": "2026-10-09T10:00:00Z"},
		// NEWER but in the EPC mailbox with requirements-y words — the
		// software pipeline must never grab the other company's mail.
		{"id": "email_epc_est", "from": "projects@newco.com", "subject": "Estimation — warehouse build", "body": "please build and develop our warehouse proposal", "category": "client", "mailbox": "epc", "date": "2026-10-10T10:00:00Z"},
	} {
		raw, _ := json.Marshal(e)
		_ = os.WriteFile(filepath.Join(inboxDir, fmt.Sprintf("%v.json", e["id"])), raw, 0o644)
		_ = i
	}

	man, err := r.Run(context.Background(), RunInput{})
	if err != nil {
		t.Fatal(err)
	}
	if man.Status != StatusDone {
		t.Fatalf("expected done, got %#v", man)
	}
	if man.Client != "cto@fitflow.app" || man.EmailID != "email_req" {
		t.Fatalf("picked wrong email: %#v", man)
	}
	if man.ProjectID == "" {
		t.Fatal("expected a derived project id")
	}
	// Original email marked replied for the loop-closure story.
	raw, _ := os.ReadFile(filepath.Join(inboxDir, "email_req.json"))
	if !strings.Contains(string(raw), `"replied":true`) && !strings.Contains(string(raw), `"replied": true`) {
		t.Fatalf("source email not marked replied: %s", raw)
	}
}

func TestDiscoveryNoRequestAnywhere(t *testing.T) {
	r := newRunner(t, &fakeLLM{})
	man, err := r.Run(context.Background(), RunInput{})
	if err == nil && man.Status != StatusFailed {
		t.Fatalf("expected failure without any request, got %#v", man)
	}
}

// Live progress: a "working" beat streams BEFORE each agent's call (so the
// workforce panel pulses during generation, not after), completion notes
// carry a section summary, and the manifest exposes what is running now.
func TestDiscoveryStreamsLiveProgress(t *testing.T) {
	em := &recordingEmitter{}
	r := newRunner(t, &fakeLLM{})
	r.Emitter = em

	man, err := r.Run(context.Background(), RunInput{
		ProjectID: "p1", Client: "c@x.com", Request: "Build a fitness tracking app",
	})
	if err != nil {
		t.Fatal(err)
	}
	// Started beats: one per agent doc + one for the reply.
	started := 0
	for _, n := range em.notes {
		if strings.Contains(n, "working on") {
			started++
		}
	}
	if started < 5 {
		t.Fatalf("expected a working beat before each of the 5 agent calls, got %d in %v", started, em.notes)
	}
	// Completion notes summarize the document's sections for the user.
	var prdNote string
	for _, s := range man.Steps {
		if s.Action == "write_prd" {
			prdNote = s.Note
		}
	}
	if !strings.Contains(prdNote, "Goals") {
		t.Fatalf("completion note should name the doc's sections, got %q", prdNote)
	}
	// Nothing is "current" once the pipeline finished.
	if man.Current != "" {
		t.Fatalf("current should clear on completion, got %q", man.Current)
	}
}

func TestDiscoverySurvivesOneAgentFailure(t *testing.T) {
	llm := &fakeLLM{failFor: `"Competitor Research" document`}
	r := newRunner(t, llm)
	man, err := r.Run(context.Background(), RunInput{
		ProjectID: "p1", Client: "c@x.com", Request: "Build an inventory requirements portal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if man.Status != StatusDoneWithErrors {
		t.Fatalf("expected done_with_errors, got %q", man.Status)
	}
	var sawFail bool
	for _, s := range man.Steps {
		if s.Agent == "senior_dev" && !s.OK {
			sawFail = true
		}
	}
	if !sawFail {
		t.Fatalf("senior_dev failure not recorded: %#v", man.Steps)
	}
	// The rest of the pipeline still produced its docs.
	if _, err := os.Stat(filepath.Join(r.ProjectsDir, "p1", "discovery", "04-architecture.md")); err != nil {
		t.Fatal("architecture doc should exist despite competitor failure")
	}
	if _, err := os.Stat(filepath.Join(r.ProjectsDir, "p1", "discovery", "05-client-reply.md")); err != nil {
		t.Fatal("client reply should exist despite competitor failure")
	}
}
