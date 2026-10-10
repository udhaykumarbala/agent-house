package capability

import (
	"os"
	"testing"
)

// The inbox is one directory, but each company gets its own mailbox view:
// mailbox "software" vs "epc", with legacy emails (no field) counting as
// EPC so nothing seeded before the split changes behavior.
func TestInboxMailboxFiltering(t *testing.T) {
	old, _ := os.Getwd()
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	in := NewInbox("")
	seed := []Email{
		{ID: "e_sw", From: "cto@fitflow.app", Subject: "Requirements — FitFlow", Mailbox: "software"},
		{ID: "e_epc", From: "ar@vendorxyz.com", Subject: "Invoice #42", Mailbox: "epc"},
		{ID: "e_legacy", From: "client@newco.com", Subject: "Estimation request"}, // pre-split email
	}
	for _, e := range seed {
		if _, err := in.Inject(e); err != nil {
			t.Fatal(err)
		}
	}

	sw, err := in.ListMailbox("software")
	if err != nil {
		t.Fatal(err)
	}
	if len(sw) != 1 || sw[0].ID != "e_sw" {
		t.Fatalf("software mailbox wrong: %#v", sw)
	}

	epc, _ := in.ListMailbox("epc")
	if len(epc) != 2 {
		t.Fatalf("epc mailbox should include the legacy email, got %d", len(epc))
	}

	all, _ := in.ListMailbox("")
	if len(all) != 3 {
		t.Fatalf("empty box should mean all, got %d", len(all))
	}

	// Unread + summary respect the same boundary.
	un, _ := in.UnreadMailbox("epc")
	if len(un) != 2 {
		t.Fatalf("epc unread: want 2, got %d", len(un))
	}
	sum, _ := in.SummaryMailbox("software")
	if sum.Total != 1 {
		t.Fatalf("software summary: want 1, got %d", sum.Total)
	}
}

// Inject infers the mailbox from the recipient when none is set, so real
// inbound mail routes itself: software-studio addresses → software.
func TestInboxInjectInfersMailboxFromRecipient(t *testing.T) {
	old, _ := os.Getwd()
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	in := NewInbox("")
	saved, err := in.Inject(Email{ID: "e1", From: "x@y.com", To: "hello@software-studio.dev", Subject: "Build us an app"})
	if err != nil {
		t.Fatal(err)
	}
	if saved.Mailbox != "software" {
		t.Fatalf("expected inferred software mailbox, got %q", saved.Mailbox)
	}
	saved2, _ := in.Inject(Email{ID: "e2", From: "x@y.com", To: "ops@alredaa.com", Subject: "Site update"})
	if saved2.Mailbox != "epc" {
		t.Fatalf("expected epc default, got %q", saved2.Mailbox)
	}
}
