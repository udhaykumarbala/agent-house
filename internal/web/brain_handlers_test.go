package web

import (
	"reflect"
	"testing"
)

func TestDefaultSuggestionsFor_ConfirmationPrompts(t *testing.T) {
	cases := []struct {
		name     string
		action   string
		response string
		want     []string
	}{
		{
			"delete_email confirmation with count 2 — the screenshot's case",
			"delete_email",
			"Found 2 impersonation emails to delete. Confirm deletion of both impersonation emails?",
			[]string{"Yes, delete all 2 emails", "No, keep them", "Show inbox"},
		},
		{
			"delete_email confirmation with count 1 → singular",
			"delete_email",
			"Found 1 impersonation email. Confirm deletion?",
			[]string{"Yes, delete the email", "No, keep them", "Show inbox"},
		},
		{
			"delete_email confirmation without count",
			"delete_email",
			"Confirm deletion of these messages?",
			[]string{"Yes, delete", "No, keep them", "Show inbox"},
		},
		{
			"delete_email confirmation prefers count in question line",
			"delete_email",
			// 9 mentioned earlier in the body; the question references 2 — the
			// matcher should pick 2 because it's on the "?" line.
			"You have 9 unread emails in your inbox.\nConfirm deletion of 2 impersonation emails?",
			[]string{"Yes, delete all 2 emails", "No, keep them", "Show inbox"},
		},
		{
			"archive_email confirmation with count 5",
			"archive_email",
			"Are you sure you want to archive these 5 emails?",
			[]string{"Yes, archive all 5 emails", "Cancel", "Show inbox"},
		},
		{
			"send_reply confirmation, no count present",
			"send_reply",
			"Shall I send this draft to Sarah Jones?",
			[]string{"Yes, send", "Edit draft first", "Cancel"},
		},
		{
			"shortlist with 3 candidates",
			"shortlist_applicant",
			"Found 3 candidates matching the JD. Do you want to shortlist all of them?",
			[]string{"Yes, shortlist all 3 candidates", "Show their CV", "Cancel"},
		},
		{
			"shortlist_applicant confirmation, single named candidate (no number)",
			"shortlist_applicant",
			"Do you want to shortlist Raj Kumar?",
			[]string{"Yes, shortlist", "Show their CV", "Cancel"},
		},
		{
			"escalate confirmation with 2 slips → matched noun wins",
			"escalate",
			"There are 2 critical slips today. Should I escalate them?",
			[]string{"Yes, escalate all 2 slips", "Add context first", "Cancel"},
		},
		{
			"delete_email confirmation with 'messages' — matched noun preserved",
			"delete_email",
			"Found 4 suspicious messages. Confirm deletion?",
			[]string{"Yes, delete all 4 messages", "No, keep them", "Show inbox"},
		},
		{
			"unknown action with confirmation falls back to Yes/Cancel",
			"some_new_action",
			"Should I proceed?",
			[]string{"Yes, proceed", "Cancel"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := defaultSuggestionsFor(c.action, c.response)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("defaultSuggestionsFor(%q, %q) = %v, want %v",
					c.action, c.response, got, c.want)
			}
		})
	}
}

func TestDefaultSuggestionsFor_NonConfirmation(t *testing.T) {
	cases := []struct {
		name     string
		action   string
		response string
		want     []string
	}{
		{
			"list_projects",
			"list_projects",
			"Here are your projects: ...",
			[]string{"Show one in detail", "Create a new project", "Check inbox"},
		},
		{
			"shortlist_applicant after-action (no question)",
			"shortlist_applicant",
			"Marked Raj Kumar as shortlisted.",
			[]string{"Schedule interview", "Show their CV", "Find more like them"},
		},
		{
			"unknown action falls back to global generic",
			"some_unknown",
			"Some plain text response.",
			[]string{"Show inbox", "List projects", "Check status"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := defaultSuggestionsFor(c.action, c.response)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("defaultSuggestionsFor(%q, %q) = %v, want %v",
					c.action, c.response, got, c.want)
			}
		})
	}
}

func TestIsConfirmationPrompt(t *testing.T) {
	yes := []string{
		"Confirm deletion of both impersonation emails?",
		"Are you sure you want to send this?",
		"Shall I proceed?",
		"Should I escalate this to the CEO?",
		"Do you want me to continue?",
		"Ready to go ahead?",
	}
	no := []string{
		"Found 9 unread emails in your inbox.",
		"Done. Marked as shortlisted.",
		"Why did this fail?",                // a "?" but not a confirmation
		"What needs my attention today?",    // user-style question, not confirmation
	}
	for _, s := range yes {
		if !isConfirmationPrompt(s) {
			t.Errorf("expected confirmation: %q", s)
		}
	}
	for _, s := range no {
		if isConfirmationPrompt(s) {
			t.Errorf("expected NOT confirmation: %q", s)
		}
	}
}
