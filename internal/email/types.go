package email

import "time"

// Email represents an inbound or outbound email.
type Email struct {
	ID          string    `json:"id"`
	From        string    `json:"from"`
	FromName    string    `json:"from_name"`
	To          string    `json:"to"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	Date        time.Time `json:"date"`
	Read        bool      `json:"read"`
	Category    string    `json:"category"`     // vendor, job_application, client, internal, unknown
	VendorID    string    `json:"vendor_id,omitempty"`
	TrustStatus string    `json:"trust_status"` // trusted, new_contact, untrusted, impersonation
	TrustReason string    `json:"trust_reason"`
	Direction   string    `json:"direction"`    // inbound, outbound
	ReplyTo     string    `json:"reply_to,omitempty"`
}

// Vendor represents a vendor with trusted email contacts.
type Vendor struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Domain        string   `json:"domain"`
	TrustedEmails []string `json:"trusted_emails"`
	ContactPerson string   `json:"contact_person"`
	ContractActive bool    `json:"contract_active"`
}

// Applicant represents a job applicant parsed from email.
type Applicant struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	AppliedFor     string   `json:"applied_for"`
	ExperienceYears int     `json:"experience_years,omitempty"`
	KeySkills      []string `json:"key_skills,omitempty"`
	Certifications []string `json:"certifications,omitempty"`
	ResumeEmailID  string   `json:"resume_email_id"`
	ReceivedDate   string   `json:"received_date"`
	Status         string   `json:"status"` // new, reviewed, shortlisted, rejected
}

// TrustResult holds the outcome of a vendor trust check.
type TrustResult struct {
	Status      string   `json:"status"`       // trusted, new_contact, untrusted, impersonation
	Reason      string   `json:"reason"`
	RiskFactors []string `json:"risk_factors,omitempty"`
	VendorID    string   `json:"vendor_id,omitempty"`
	VendorName  string   `json:"vendor_name,omitempty"`
}
