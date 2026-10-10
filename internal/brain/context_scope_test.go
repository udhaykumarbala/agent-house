package brain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The Brain's chat context must respect the company mailboxes: the software
// studio never sees EPC mail in "check mails" and vice versa. (The leak:
// AssembleContext used to read the whole inbox dir unfiltered.)
func TestAssembleContextScopedFiltersMailboxes(t *testing.T) {
	dir := t.TempDir()
	inbox := filepath.Join(dir, "inbox")
	_ = os.MkdirAll(inbox, 0o755)
	write := func(id string, e map[string]any) {
		e["id"] = id
		raw, _ := json.Marshal(e)
		_ = os.WriteFile(filepath.Join(inbox, id+".json"), raw, 0o644)
	}
	write("e_sw", map[string]any{"from": "cto@fitflow.app", "subject": "FitFlow requirements", "mailbox": "software"})
	write("e_epc", map[string]any{"from": "ar@vendorxyz.com", "subject": "Steel invoice ST-0847", "mailbox": "epc"})
	write("e_legacy", map[string]any{"from": "x@y.com", "subject": "Legacy site update"}) // pre-split → epc

	sw := AssembleContextScoped(dir, "software")
	if !strings.Contains(sw, "FitFlow requirements") {
		t.Fatal("software context lost its own mail")
	}
	if strings.Contains(sw, "Steel invoice") || strings.Contains(sw, "Legacy site update") {
		t.Fatalf("software context leaked EPC mail:\n%s", sw)
	}

	epc := AssembleContextScoped(dir, "atlas-site")
	if !strings.Contains(epc, "Steel invoice") || !strings.Contains(epc, "Legacy site update") {
		t.Fatal("EPC context lost its own/legacy mail")
	}
	if strings.Contains(epc, "FitFlow requirements") {
		t.Fatalf("EPC context leaked software mail:\n%s", epc)
	}
}

// "View the FitFlow project page" on a discovery-only project used to hit
// the project.json dead end. It must now summarize the discovery run and
// point at the Documents page.
func TestProjectStatusFallsBackToDiscovery(t *testing.T) {
	dir := t.TempDir()
	disc := filepath.Join(dir, "fitflow-app", "discovery")
	_ = os.MkdirAll(disc, 0o755)
	_ = os.WriteFile(filepath.Join(disc, "manifest.json"), []byte(`{
		"project_id":"fitflow-app","status":"done","client":"cto@fitflow.app",
		"steps":[{"agent":"pm","action":"write_prd","ok":true},{"agent":"architect","action":"system_design","ok":true}]
	}`), 0o644)
	_ = os.WriteFile(filepath.Join(disc, "01-prd.md"), []byte("# FitFlow PRD"), 0o644)

	e := NewExecutor(dir)
	res := e.executeProjectStatus(&BrainDecision{
		Action: ActionProjectStatus,
		Params: paramMap{"project_id": "fitflow-app"},
	})
	if !res.Success {
		t.Fatalf("expected success, got %#v", res)
	}
	low := strings.ToLower(res.Response)
	if !strings.Contains(low, "documents") || !strings.Contains(low, "fitflow-app") {
		t.Fatalf("response should point at the Documents page, got %q", res.Response)
	}
	if !strings.Contains(low, "done") {
		t.Fatalf("response should carry the pipeline status, got %q", res.Response)
	}
}
