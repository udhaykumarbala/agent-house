package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// ResendClient sends emails via the Resend API.
type ResendClient struct {
	APIKey string
}

// NewResendClient creates a client from RESEND_API_KEY env var.
func NewResendClient() *ResendClient {
	key := os.Getenv("RESEND_API_KEY")
	if key != "" {
		log.Printf("[RESEND] API key configured — outbound email enabled")
	}
	return &ResendClient{APIKey: key}
}

// IsAvailable returns true if Resend is configured.
func (r *ResendClient) IsAvailable() bool {
	return r.APIKey != ""
}

// Send sends an email via Resend API.
func (r *ResendClient) Send(from, to, subject, htmlBody string) error {
	if !r.IsAvailable() {
		return fmt.Errorf("RESEND_API_KEY not set — email not sent (would send to %s)", to)
	}

	payload := map[string]string{
		"from":    from,
		"to":      to,
		"subject": subject,
		"html":    htmlBody,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend API error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend API %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[RESEND] Sent email: to=%s subject=%q", to, subject)
	return nil
}
