// Package discovery runs the software studio's client-discovery pipeline:
// a client requirements email (or chat request) fans out to real LLM agent
// calls — PM (PRD), senior_dev (competitor research), security (compliance
// research), architect (overall design) — each writing a rich markdown
// document (tables + mermaid diagrams) under projects/<id>/discovery/, and
// ends with a client reply draft in the outbox. Unlike the deterministic
// EPC scenarios, every document here is genuinely generated per request;
// unlike the orchestrator's oneshot phases, the output is PERSISTED so the
// Documents page can render it.
package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"pty-claude-test/internal/message"
)

// LLM is the one-call generation interface (adapter over session.APIClient
// in production, a fake in tests).
type LLM interface {
	Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// Emitter mirrors scenario.Emitter so steps stream to the workforce panel.
type Emitter interface {
	Emit(m *message.Message)
}

// Manifest statuses.
const (
	StatusRunning        = "running"
	StatusDone           = "done"
	StatusDoneWithErrors = "done_with_errors"
	StatusFailed         = "failed"
)

// Step is one pipeline stage, mirroring the scenario Step shape the UI knows.
type Step struct {
	Agent      string `json:"agent"`
	Action     string `json:"action"`
	File       string `json:"file,omitempty"` // path relative to the project dir
	Note       string `json:"note,omitempty"`
	OK         bool   `json:"ok"`
	DurationMS int64  `json:"duration_ms,omitempty"`
}

// Manifest is the pipeline's durable status, written after every step so
// the UI can poll it live.
type Manifest struct {
	ProjectID  string    `json:"project_id"`
	Client     string    `json:"client,omitempty"`
	Subject    string    `json:"subject,omitempty"`
	EmailID    string    `json:"email_id,omitempty"`
	Request    string    `json:"request"`
	Status     string    `json:"status"`
	Steps      []Step    `json:"steps"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

// RunInput selects what to discover. Everything is optional: with no
// request and no email id, the newest requirements-looking inbox email is
// picked up automatically.
type RunInput struct {
	ProjectID string
	Request   string
	Client    string
	Subject   string
	EmailID   string
}

// Runner executes the pipeline. ProjectsDir is the projects root (inbox at
// <ProjectsDir>/inbox, outbox at <ProjectsDir>/outbox, artifacts under
// <ProjectsDir>/<project>/discovery).
type Runner struct {
	ProjectsDir string
	LLM         LLM
	Prompts     func(role string) string // system prompt per role
	Emitter     Emitter                  // optional
}

// docSpec is one agent-authored document in the pipeline.
type docSpec struct {
	role   string
	action string
	file   string
	title  string
	brief  string // what the document must cover
}

var pipelineDocs = []docSpec{
	{
		role: "pm", action: "write_prd", file: "01-prd.md", title: "Product Requirements Document",
		brief: "a complete PRD: product summary, goals, personas, user stories (as a table), functional scope table (MoSCoW), non-functional requirements, success metrics, and a user-flow diagram",
	},
	{
		role: "senior_dev", action: "competitor_research", file: "02-competitor-research.md", title: "Competitor Research",
		brief: "desk research on the competitive landscape: a competitor comparison matrix (table), feature gaps we can exploit, pricing patterns, and a positioning diagram. Be explicit that this is desk research from general knowledge, flagging anything that should be verified",
	},
	{
		role: "security", action: "compliance_research", file: "03-compliance-research.md", title: "Compliance & Security Research",
		brief: "the compliance and security posture this product needs: applicable regulations (GDPR, Saudi PDPL, SOC2 where relevant) as a requirements table, data classification, a data-flow diagram showing where personal data moves, authentication/authorization requirements, and a pre-launch security checklist",
	},
	{
		role: "architect", action: "system_design", file: "04-architecture.md", title: "System Architecture",
		brief: "the overall system design: architecture overview with a component diagram, a sequence diagram for the core user journey, recommended tech stack (table with rationale), data model sketch, API surface outline, and delivery phases",
	},
}

// Resolve fills the input from the inbox when needed (explicit text >
// named email > newest requirements-looking email) and derives the project
// id. It does no LLM work, so handlers can validate synchronously before
// launching the pipeline in the background.
func (r *Runner) Resolve(in RunInput) (RunInput, error) {
	if in.Request == "" {
		em := r.findRequirementsEmail(in.EmailID)
		if em == nil {
			return in, fmt.Errorf("no requirements found: pass request text or drop the client's email in the inbox")
		}
		in.Request = em.Subject + "\n\n" + em.Body
		in.Client = em.From
		in.Subject = em.Subject
		in.EmailID = em.ID
	}
	if in.ProjectID == "" {
		in.ProjectID = deriveProjectID(in.Subject, in.Request)
	}
	return in, nil
}

// Run executes the full pipeline synchronously. Callers wanting async wrap
// it in a goroutine and poll the manifest.
func (r *Runner) Run(ctx context.Context, in RunInput) (*Manifest, error) {
	man := &Manifest{Status: StatusRunning, StartedAt: time.Now().UTC()}

	var err error
	if in, err = r.Resolve(in); err != nil {
		man.Status = StatusFailed
		man.FinishedAt = time.Now().UTC()
		return man, err
	}
	man.ProjectID = in.ProjectID
	man.Client = in.Client
	man.Subject = in.Subject
	man.EmailID = in.EmailID
	man.Request = trimTo(in.Request, 4000)

	dir := filepath.Join(r.ProjectsDir, in.ProjectID, "discovery")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		man.Status = StatusFailed
		return man, err
	}

	// Step 0: conductor files the request.
	reqDoc := fmt.Sprintf("# Client Request\n\n**From:** %s\n**Subject:** %s\n**Received:** %s\n\n---\n\n%s\n",
		orDash(in.Client), orDash(in.Subject), man.StartedAt.Format(time.RFC3339), in.Request)
	_ = os.WriteFile(filepath.Join(dir, "00-request.md"), []byte(reqDoc), 0o644)
	r.recordStep(man, dir, Step{
		Agent: "conductor", Action: "file_request", OK: true,
		File: "discovery/00-request.md",
		Note: "client requirements filed — fanning out to the team",
	})

	// Steps 1-4: the agent documents. The PRD is generated first and fed
	// to the later agents as shared context.
	var prdContext string
	for _, d := range pipelineDocs {
		start := time.Now()
		prompt := r.docPrompt(d, in, prdContext)
		text, err := r.LLM.Generate(ctx, r.systemPrompt(d.role), prompt)
		step := Step{Agent: d.role, Action: d.action, File: "discovery/" + d.file, DurationMS: time.Since(start).Milliseconds()}
		if err != nil {
			step.OK = false
			step.Note = d.title + " failed: " + err.Error()
		} else {
			doc := ensureTitle(text, d.title)
			if werr := os.WriteFile(filepath.Join(dir, d.file), []byte(doc), 0o644); werr != nil {
				step.OK = false
				step.Note = d.title + " write failed: " + werr.Error()
			} else {
				step.OK = true
				step.Note = fmt.Sprintf("%s written (%d chars)", d.title, len(doc))
				if d.role == "pm" {
					prdContext = trimTo(doc, 2500)
				}
			}
		}
		r.recordStep(man, dir, step)
	}

	// Step 5: PM composes the client reply; draft goes to the outbox.
	start := time.Now()
	replyPrompt := fmt.Sprintf(
		"The client sent us this request:\n\n%s\n\nOur team has completed discovery: a PRD, competitor research, compliance research, and a system architecture are ready.\n\nPRD summary for grounding:\n%s\n\nWrite the reply EMAIL to the client (markdown, starting with 'Subject: ...' on the first line): thank them, summarize in 4-6 bullet points what our discovery found (scope, differentiators, compliance posture, proposed architecture and phases), state that detailed documents are attached, and ask for approval to proceed to the build phase. Professional, warm, concise — no placeholder brackets.",
		trimTo(in.Request, 1500), orDash(prdContext))
	replyText, rerr := r.LLM.Generate(ctx, r.systemPrompt("pm"), replyPrompt)
	replyStep := Step{Agent: "pm", Action: "draft_client_reply", File: "discovery/05-client-reply.md", DurationMS: time.Since(start).Milliseconds()}
	if rerr != nil {
		replyStep.OK = false
		replyStep.Note = "client reply failed: " + rerr.Error()
	} else {
		_ = os.WriteFile(filepath.Join(dir, "05-client-reply.md"), []byte(ensureTitle(replyText, "Client Reply Draft")), 0o644)
		r.writeOutboxDraft(in, replyText)
		if in.EmailID != "" {
			r.markReplied(in.EmailID)
		}
		replyStep.OK = true
		replyStep.Note = "reply drafted to " + orDash(in.Client) + " and saved to the outbox"
	}
	r.recordStep(man, dir, replyStep)

	man.Status = StatusDone
	for _, s := range man.Steps {
		if !s.OK {
			man.Status = StatusDoneWithErrors
			break
		}
	}
	man.FinishedAt = time.Now().UTC()
	r.saveManifest(man, dir)
	return man, nil
}

func (r *Runner) docPrompt(d docSpec, in RunInput, prdContext string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CLIENT REQUEST (verbatim):\n%s\n\n", trimTo(in.Request, 2500))
	if prdContext != "" && d.role != "pm" {
		fmt.Fprintf(&b, "PRD (already written by the PM — stay consistent with it):\n%s\n\n", prdContext)
	}
	fmt.Fprintf(&b, "Write the %q document for this engagement: %s.\n\n", d.title, d.brief)
	b.WriteString("FORMAT RULES (strict):\n")
	b.WriteString("- GitHub-flavored markdown, starting with a single `# ` title.\n")
	b.WriteString("- Use tables for any comparison or list of requirements.\n")
	b.WriteString("- Include at least one ```mermaid diagram (flowchart, sequenceDiagram, or graph) that genuinely illustrates this product — keep node labels short, 2-6 words, and the diagram under 20 nodes.\n")
	b.WriteString("- No placeholders like TBD/[x] — commit to concrete recommendations.\n")
	b.WriteString("- 400-900 words. Output ONLY the markdown document, no preamble.\n")
	return b.String()
}

func (r *Runner) systemPrompt(role string) string {
	if r.Prompts != nil {
		if p := r.Prompts(role); p != "" {
			return p
		}
	}
	return "You are the " + role + " of a software studio, producing client-facing discovery documents."
}

func (r *Runner) recordStep(man *Manifest, dir string, s Step) {
	man.Steps = append(man.Steps, s)
	r.saveManifest(man, dir)
	if r.Emitter != nil {
		content := fmt.Sprintf("[discovery:%s] %s · %s — %s", man.ProjectID, s.Agent, s.Action, s.Note)
		m := message.NewMessage(message.TypeSystem, s.Agent, "conductor", content)
		m.Metadata.ProjectID = man.ProjectID
		m.Metadata.Tags = []string{"discovery", man.ProjectID, s.Agent}
		r.Emitter.Emit(m)
	}
}

func (r *Runner) saveManifest(man *Manifest, dir string) {
	data, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, "manifest.json"), data, 0o644)
}

// inboxEmail is the minimal on-disk email shape (owned copy; the global
// inbox files are written by the email engine and seed scripts).
type inboxEmail struct {
	ID       string `json:"id"`
	From     string `json:"from"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Category string `json:"category"`
	Date     string `json:"date"`
	Replied  bool   `json:"replied"`
}

var requirementsKeywords = []string{
	"requirement", "build", "develop", "mvp", "proposal", "app", "platform", "website", "software",
}

// findRequirementsEmail returns the email with the given id, or the newest
// inbox email that reads like a software requirements request.
func (r *Runner) findRequirementsEmail(emailID string) *inboxEmail {
	dir := filepath.Join(r.ProjectsDir, "inbox")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var all []inboxEmail
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") ||
			e.Name() == "applicants.json" || e.Name() == "vendors.json" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var em inboxEmail
		if json.Unmarshal(raw, &em) != nil || em.ID == "" {
			continue
		}
		all = append(all, em)
	}
	sort.Slice(all, func(a, b int) bool { return all[a].Date > all[b].Date })
	for i := range all {
		if emailID != "" {
			if all[i].ID == emailID {
				return &all[i]
			}
			continue
		}
		if all[i].Category == "vendor" || all[i].Category == "job_application" || all[i].Category == "applicant" {
			continue
		}
		hay := strings.ToLower(all[i].Subject + " " + all[i].Body)
		for _, kw := range requirementsKeywords {
			if strings.Contains(hay, kw) {
				return &all[i]
			}
		}
	}
	return nil
}

// markReplied flips replied=true on the source email file (best effort).
func (r *Runner) markReplied(emailID string) {
	path := filepath.Join(r.ProjectsDir, "inbox", emailID+".json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var obj map[string]any
	if json.Unmarshal(raw, &obj) != nil {
		return
	}
	obj["replied"] = true
	obj["read"] = true
	if data, err := json.MarshalIndent(obj, "", "  "); err == nil {
		_ = os.WriteFile(path, data, 0o644)
	}
}

func (r *Runner) writeOutboxDraft(in RunInput, body string) {
	outDir := filepath.Join(r.ProjectsDir, "outbox")
	_ = os.MkdirAll(outDir, 0o755)
	subject := "Re: " + orDash(in.Subject)
	if in.Subject == "" {
		subject = "Your project proposal — discovery complete"
	}
	draft := map[string]any{
		"id":      fmt.Sprintf("reply_%d", time.Now().UnixMilli()),
		"to":      in.Client,
		"subject": subject,
		"body":    body,
		"draft":   true,
		"project": in.ProjectID,
		"date":    time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.MarshalIndent(draft, "", "  ")
	_ = os.WriteFile(filepath.Join(outDir, draft["id"].(string)+".json"), data, 0o644)
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

// deriveProjectID builds a kebab slug from the subject (or request head).
func deriveProjectID(subject, request string) string {
	src := subject
	if src == "" {
		src = trimTo(request, 60)
	}
	// Drop common email prefixes so slugs read like product names.
	low := strings.ToLower(src)
	for _, p := range []string{"re:", "fwd:", "requirements", "request for", "proposal"} {
		low = strings.ReplaceAll(low, p, " ")
	}
	slug := strings.Trim(slugRe.ReplaceAllString(low, "-"), "-")
	if len(slug) > 32 {
		slug = strings.Trim(slug[:32], "-")
	}
	if slug == "" {
		slug = fmt.Sprintf("discovery-%d", time.Now().Unix())
	}
	return slug
}

func ensureTitle(text, fallback string) string {
	t := strings.TrimSpace(text)
	if !strings.HasPrefix(t, "#") {
		t = "# " + fallback + "\n\n" + t
	}
	return t + "\n"
}

func trimTo(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}
