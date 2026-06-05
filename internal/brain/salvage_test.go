package brain

import "testing"

// A decision envelope whose markdown "response" carries literal (unescaped)
// newlines — exactly the shape that breaks strict JSON parsing and used to be
// dumped at the user as raw JSON. parseDecision's Try-4 must salvage it.
func TestParseDecisionSalvagesMalformedEnvelope(t *testing.T) {
	malformed := `{"action": "respond",
"response": "### Project Alpha — Mitigation Plan
The 20-day C-7 slab delay cascades into structural erection.
- Option A: fast-track RFI #18
- Option B: enforce XYZ Steel penalty",
"suggestions": ["Draft RFI #18 response"]}`

	d, err := parseDecision(malformed)
	if err != nil {
		t.Fatalf("expected salvage, got error: %v", err)
	}
	if d.Action != ActionRespond {
		t.Fatalf("action = %q, want respond", d.Action)
	}
	if got := d.Response; len(got) < 20 || got[:3] != "###" {
		t.Fatalf("response not salvaged cleanly: %.60q", got)
	}
	// the raw envelope must NOT leak into the response
	if containsSub(d.Response, `"action"`) {
		t.Fatalf("raw envelope leaked into response: %.80q", d.Response)
	}
}

func TestExtractJSONStringFieldHandlesEscapes(t *testing.T) {
	in := `{"response": "line one\nline two with \"quotes\" inside", "x": 1}`
	got := extractJSONStringField(in, "response")
	want := "line one\nline two with \"quotes\" inside"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func containsSub(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
