# Security Specification: SubTrack - Subscription Manager

## Document Overview

This specification defines the security architecture, controls, and implementation requirements for SubTrack, a privacy-first subscription management web application. It is based on the security research conducted and aligned with the architecture decisions made by the team.

**Key Security Advantage**: By avoiding bank connections (unlike Truebill/Rocket Money), we significantly reduce attack surface and liability. This is a core security design decision.

---

## 1. Threat Model Summary

### 1.1 Assets to Protect

| Asset | Classification | Protection Priority |
|-------|---------------|---------------------|
| User credentials (email, password hash) | Sensitive | Critical |
| Session tokens | Sensitive | Critical |
| Subscription data (names, prices, dates) | Personal/Financial | High |
| User preferences | Personal | Medium |
| Aggregated analytics | Operational | Low |

### 1.2 Primary Threat Actors

| Actor | Risk Level | Primary Concerns |
|-------|-----------|------------------|
| Opportunistic Hackers | High | Credential stuffing, data theft |
| Script Kiddies | Medium | Automated vulnerability scanning |
| Malicious Insiders | Low | Data exfiltration |

### 1.3 Attack Surface

```
Frontend (PWA)
├── XSS via user-generated content (subscription names, notes)
├── CSRF on state-changing operations
├── Local storage data exposure
└── Service worker cache poisoning

API Layer
├── Authentication bypass
├── IDOR (Insecure Direct Object Reference)
├── Rate limiting bypass
├── Injection attacks
└── CORS misconfiguration

Database
├── SQL injection
├── Data exposure via backup leaks
└── Unencrypted sensitive fields
```

---

## 2. Authentication Security

### 2.1 Session-Based Authentication (MVP)

**Decision**: Use session-based authentication for MVP per architecture research. Simpler, more secure for initial deployment.

```
Authentication Flow:
┌────────┐    credentials     ┌────────┐    verify/create    ┌──────────┐
│ Client │ ────────────────►  │  API   │ ──────────────────► │ Database │
└────────┘                    └────────┘                     └──────────┘
     ▲                             │
     │    Set-Cookie: session_id   │
     │    HttpOnly; Secure;        │
     │    SameSite=Strict          │
     └─────────────────────────────┘
```

### 2.2 Password Requirements

| Requirement | Value | Rationale |
|-------------|-------|-----------|
| Minimum length | 12 characters | NIST 800-63B guideline |
| Maximum length | 128 characters | Allow passphrases |
| Complexity | None required | Length > complexity per NIST |
| Character set | All Unicode allowed | Don't restrict creativity |
| Breach check | Required | Check against HaveIBeenPwned API (k-anonymity) |

### 2.3 Password Hashing

```go
// Use Argon2id - OWASP recommended
Config:
  - Memory: 64 MB (65536 KB)
  - Iterations: 3
  - Parallelism: 4
  - Salt length: 16 bytes (random)
  - Key length: 32 bytes
```

### 2.4 Session Management

| Setting | Value | Rationale |
|---------|-------|-----------|
| Token generation | 32 bytes crypto random, hex encoded | Sufficient entropy |
| Storage | Hash of token in DB, not plaintext | Defense in depth |
| Cookie flags | `HttpOnly`, `Secure`, `SameSite=Strict`, `Path=/` | XSS/CSRF protection |
| Session duration | 7 days sliding window | UX balance |
| Concurrent sessions | Allow multiple | Multi-device support |
| Session rotation | On privilege change (password reset, etc.) | Prevent session fixation |

### 2.5 Login Protection

| Control | Implementation |
|---------|----------------|
| Rate limiting | 10 attempts per email per minute |
| Account lockout | Temporary 15-minute lockout after 5 failures |
| IP blocking | After 50 failed attempts across accounts in 10 minutes |
| CAPTCHA | Trigger after 3 failed attempts on same email |
| Generic errors | "Invalid email or password" - no enumeration |

### 2.6 Password Reset Flow

```
1. User requests reset → Generate 32-byte random token
2. Hash token, store in DB with 15-minute expiry
3. Email unhashed token in reset link
4. User clicks link → Validate hash(token) exists and not expired
5. User sets new password → Invalidate token, invalidate all sessions
6. Redirect to login
```

**Security Controls**:
- Token expires in 15 minutes
- Single-use tokens only
- Rate limit: 3 reset requests per email per hour
- Email does not confirm if account exists

---

## 3. Authorization Security

### 3.1 Access Control Model

**Principle**: Deny by default, explicit allow.

```
Every API request MUST:
1. Validate session token exists and is valid
2. Extract user_id from session
3. Verify resource belongs to user_id
4. Execute only if all checks pass
```

### 3.2 IDOR Prevention

| Resource | Protection Method |
|----------|-------------------|
| Subscriptions | Verify `subscription.user_id == session.user_id` |
| User data | Only access own user record |
| Sessions | Users can only view/revoke own sessions |

**Implementation Pattern**:
```sql
-- ALWAYS include user_id in WHERE clause
SELECT * FROM subscriptions
WHERE id = $1 AND user_id = $2;
-- Never: WHERE id = $1 (allows IDOR)
```

### 3.3 Resource Identifiers

| Resource | ID Type | Rationale |
|----------|---------|-----------|
| Users | UUIDv4 | Non-enumerable |
| Subscriptions | UUIDv4 | Non-enumerable |
| Sessions | UUIDv4 | Non-enumerable |

**Never use**: Sequential integers, predictable patterns.

---

## 4. Input Validation & Sanitization

### 4.1 Validation Rules

| Field | Type | Constraints | Sanitization |
|-------|------|-------------|--------------|
| `email` | string | Valid email format, max 255 chars | Lowercase, trim |
| `password` | string | 12-128 chars | None (preserve as-is) |
| `subscription.name` | string | 1-100 chars, printable | HTML escape, trim |
| `subscription.price` | decimal | 0 - 999999.99 | Round to 2 decimals |
| `subscription.billing_cycle` | enum | `weekly`, `monthly`, `yearly`, `custom` | Strict enum validation |
| `subscription.category` | enum | Predefined list | Strict enum validation |
| `subscription.notes` | string | 0-500 chars | HTML escape, trim |
| `subscription.next_billing_date` | date | Valid ISO 8601, not > 10 years future | Parse strictly |
| `subscription.reminder_days` | integer | 0-30 | Clamp to range |
| `subscription.color` | string | Valid hex color (#RRGGBB) | Regex validation |

### 4.2 Validation Implementation

```go
// Server-side validation is MANDATORY
// Never trust client-side validation alone

func ValidateSubscription(input SubscriptionInput) error {
    if len(input.Name) < 1 || len(input.Name) > 100 {
        return ErrInvalidName
    }
    if input.Price < 0 || input.Price > 999999.99 {
        return ErrInvalidPrice
    }
    if !isValidBillingCycle(input.BillingCycle) {
        return ErrInvalidBillingCycle
    }
    // ... continue for all fields
}
```

### 4.3 SQL Injection Prevention

**Required**: All database queries MUST use parameterized queries via sqlc.

```sql
-- queries/subscriptions.sql (sqlc)
-- name: GetSubscription :one
SELECT * FROM subscriptions
WHERE id = $1 AND user_id = $2;

-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, name, price, ...)
VALUES ($1, $2, $3, ...)
RETURNING *;
```

**Never**: String concatenation in queries.

### 4.4 XSS Prevention

| Context | Prevention Method |
|---------|-------------------|
| HTML output | Automatic escaping via React JSX |
| Custom HTML rendering | DOMPurify sanitization |
| URL parameters | URL encoding |
| JSON responses | Content-Type: application/json |

**Content Security Policy** (see Section 6.1).

---

## 5. API Security

### 5.1 Rate Limiting

| Endpoint Category | Limit | Window |
|-------------------|-------|--------|
| Authentication (`/auth/*`) | 10 requests | per minute |
| General API (`/api/v1/*`) | 100 requests | per minute |
| Data export (`/user/export`) | 5 requests | per hour |
| Password reset | 3 requests per email | per hour |

**Implementation**: Use sliding window algorithm with Redis or in-memory store.

### 5.2 CORS Configuration

```go
// Strict CORS for API
corsConfig := cors.Config{
    AllowOrigins:     []string{"https://subtrack.app"}, // Production domain only
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}

// Development: Allow localhost
if env == "development" {
    corsConfig.AllowOrigins = []string{"http://localhost:3000"}
}
```

### 5.3 Request Validation

```
Every API request MUST:
1. Validate Content-Type header (application/json for POST/PUT)
2. Validate request body size (max 1MB)
3. Validate JSON structure
4. Validate individual field types and constraints
5. Return 400 for validation errors (generic messages)
```

### 5.4 Response Security

| Response Type | Headers Required |
|---------------|------------------|
| All responses | `Cache-Control: no-store` |
| API responses | `Content-Type: application/json` |
| Error responses | Generic messages, no stack traces |

**Error Response Format**:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid subscription data"
  }
}
```

**Never expose**: Stack traces, SQL errors, internal paths, server versions.

---

## 6. Security Headers

### 6.1 Required Headers

```http
# Enforce HTTPS
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload

# Prevent XSS (primary defense is CSP)
Content-Security-Policy:
  default-src 'self';
  script-src 'self';
  style-src 'self' 'unsafe-inline';
  img-src 'self' data: https:;
  connect-src 'self' https://api.subtrack.app;
  font-src 'self';
  object-src 'none';
  frame-ancestors 'none';
  form-action 'self';
  base-uri 'self';

# Prevent MIME sniffing
X-Content-Type-Options: nosniff

# Prevent framing (clickjacking)
X-Frame-Options: DENY

# Control referrer
Referrer-Policy: strict-origin-when-cross-origin

# Disable dangerous browser features
Permissions-Policy: geolocation=(), microphone=(), camera=(), payment=()
```

### 6.2 API Response Headers

```http
Cache-Control: no-store
Pragma: no-cache
X-Content-Type-Options: nosniff
```

---

## 7. Data Protection

### 7.1 Encryption at Rest

| Data | Encryption | Method |
|------|------------|--------|
| Passwords | Yes | Argon2id hash (not reversible encryption) |
| Session tokens | Yes | SHA-256 hash stored |
| Subscription data | Database-level | PostgreSQL TDE (if available) or application-level for sensitive notes |
| User preferences | Database-level | Standard DB encryption |

### 7.2 Encryption in Transit

| Connection | Protocol | Minimum Version |
|------------|----------|-----------------|
| Client ↔ Frontend | HTTPS | TLS 1.2+ |
| Frontend ↔ API | HTTPS | TLS 1.2+ |
| API ↔ Database | TLS | TLS 1.2+ |

**Certificate**: Use automated certificates (Let's Encrypt via hosting provider).

### 7.3 Local Storage Security (PWA)

| Data Type | Storage Method | Protection |
|-----------|----------------|------------|
| Auth tokens | Secure HTTP-only cookies | Not in localStorage |
| Cached subscriptions | IndexedDB | Consider encryption with user-derived key |
| User preferences | IndexedDB/localStorage | Non-sensitive only |
| Pending mutations | IndexedDB | Clear on logout |

**Logout Behavior**:
```javascript
// Clear ALL local data on logout
async function logout() {
    await clearIndexedDB();
    localStorage.clear();
    sessionStorage.clear();
    await caches.delete('subtrack-cache');
    // Session cookie cleared by server
}
```

---

## 8. Privacy Compliance

### 8.1 GDPR/CCPA Requirements

| Requirement | Implementation |
|-------------|----------------|
| Data access (Article 15) | `/user/export` endpoint - JSON/CSV export |
| Data deletion (Article 17) | `/user/account` DELETE - cascading delete |
| Data portability (Article 20) | Export includes all user data |
| Consent | Clear terms acceptance on registration |
| Privacy policy | Accessible at `/privacy`, linked during signup |

### 8.2 Data Retention

| Data Type | Retention Period | Deletion Method |
|-----------|------------------|-----------------|
| User account | Until deletion requested | Hard delete |
| Subscriptions | With user account | Cascade delete |
| Sessions | Auto-expire after 7 days | Cron job cleanup |
| Audit logs | 90 days | Auto-purge |
| Backups | 30 days | Auto-rotation |

### 8.3 Data Minimization

**Collect only**:
- Email (for account identification)
- Password (for authentication)
- Subscription details (core functionality)
- User preferences (optional, for UX)

**Never collect**:
- Real names (optional only)
- Bank account information
- Credit card numbers
- Social security numbers
- Location data

---

## 9. Logging & Monitoring

### 9.1 Security Events to Log

| Event | Log Level | Data Captured |
|-------|-----------|---------------|
| Successful login | INFO | user_id, ip, timestamp, user_agent |
| Failed login | WARN | email (hashed), ip, timestamp |
| Password change | INFO | user_id, timestamp |
| Password reset request | INFO | email (hashed), ip, timestamp |
| Session created/revoked | INFO | user_id, session_id (prefix only) |
| Account deletion | INFO | user_id, timestamp |
| Rate limit triggered | WARN | ip, endpoint, timestamp |
| Authorization failure | WARN | user_id, resource, timestamp |

### 9.2 Log Security

| Control | Implementation |
|---------|----------------|
| PII handling | Hash or truncate sensitive data |
| Log storage | Append-only, separate from app data |
| Access | Restricted to ops team |
| Retention | 90 days default |
| Integrity | Immutable log storage if possible |

### 9.3 Alerting Triggers

| Condition | Alert Priority |
|-----------|---------------|
| 10+ failed logins from same IP in 5 min | High |
| Rate limiting triggered 100+ times in 10 min | High |
| Authorization failures > 5/min for a user | Medium |
| New IP accessing sensitive endpoint | Low (informational) |

---

## 10. Secure Development Requirements

### 10.1 Dependency Management

| Practice | Implementation |
|----------|----------------|
| Vulnerability scanning | `npm audit` and `govulncheck` in CI |
| Automated updates | Dependabot enabled |
| Lock files | `package-lock.json`, `go.sum` committed |
| Audit frequency | Weekly automated scan |

### 10.2 Code Review Security Checklist

Before approving PRs, verify:

- [ ] No hardcoded secrets/credentials
- [ ] All user input validated server-side
- [ ] Parameterized queries used (no string concatenation)
- [ ] Authorization checks on all protected endpoints
- [ ] Error messages don't leak internal details
- [ ] New endpoints have rate limiting
- [ ] Sensitive data not logged
- [ ] Tests cover security-relevant paths

### 10.3 Security Testing Requirements

| Test Type | When | Tool |
|-----------|------|------|
| Static analysis | Every PR | SonarQube, gosec |
| Dependency scan | Every PR | npm audit, govulncheck |
| Dynamic scanning | Weekly | OWASP ZAP baseline |
| SSL/TLS check | Monthly | SSL Labs |
| Penetration test | Pre-launch, quarterly | External vendor |

---

## 11. Incident Response

### 11.1 Severity Classification

| Level | Definition | Response Time | Example |
|-------|------------|---------------|---------|
| Critical | Active breach, data exposed | Immediate | Database leak |
| High | Exploitable vulnerability | 4 hours | Auth bypass found |
| Medium | Vulnerability, not exploited | 24 hours | XSS discovered |
| Low | Security improvement | 1 week | Missing header |

### 11.2 Response Procedure

```
1. IDENTIFY: Confirm the security incident
2. CONTAIN: Stop active exploitation
   - Revoke compromised sessions
   - Block attacking IPs
   - Take service offline if needed
3. ERADICATE: Fix root cause
4. RECOVER: Restore service
5. NOTIFY: Inform affected users (within 72 hours per GDPR)
6. REVIEW: Post-mortem and improvements
```

### 11.3 User Notification Template

```
Subject: Security Notice from SubTrack

What happened:
[Clear, non-technical description]

What we're doing:
[Actions taken]

What you should do:
[User actions - e.g., change password]

Questions?
Contact: security@subtrack.app
```

---

## 12. Implementation Priorities

### Phase 1: MVP (Must Have)

| Item | Priority | Owner |
|------|----------|-------|
| HTTPS everywhere with HSTS | Critical | DevOps |
| Session-based auth with secure cookies | Critical | senior_dev |
| Argon2id password hashing | Critical | senior_dev |
| Input validation middleware | Critical | senior_dev |
| Parameterized queries (sqlc) | Critical | senior_dev |
| IDOR prevention (user_id checks) | Critical | senior_dev |
| Rate limiting on auth endpoints | Critical | senior_dev |
| Basic security headers | High | senior_dev |
| Generic error messages | High | senior_dev |
| Account lockout | High | senior_dev |

### Phase 2: Post-Launch (Should Have)

| Item | Priority | Owner |
|------|----------|-------|
| Full CSP implementation | High | senior_dev |
| Encrypted local storage (PWA) | Medium | senior_dev |
| HaveIBeenPwned password check | Medium | senior_dev |
| Audit logging | Medium | senior_dev |
| CAPTCHA after failed attempts | Medium | senior_dev |
| Security event alerting | Medium | DevOps |

### Phase 3: Scale (Nice to Have)

| Item | Priority | Owner |
|------|----------|-------|
| MFA (TOTP) | Medium | senior_dev |
| WAF (Web Application Firewall) | Low | DevOps |
| Bug bounty program | Low | Security |
| SOC 2 compliance | Low | Security |
| Regular penetration testing | Low | Security |

---

## 13. Security Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              SECURITY LAYERS                                  │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                         EDGE / CDN (Vercel)                              │ │
│  │  • TLS 1.2+ termination                                                  │ │
│  │  • DDoS protection                                                       │ │
│  │  • Geographic distribution                                               │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│                                    │                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                         FRONTEND (Next.js PWA)                           │ │
│  │  • React automatic XSS escaping                                          │ │
│  │  • CSP meta tags                                                         │ │
│  │  • HTTPS-only cookies                                                    │ │
│  │  • Secure service worker                                                 │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│                                    │                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                         API GATEWAY (Chi Router)                         │ │
│  │  • Rate limiting middleware                                              │ │
│  │  • CORS enforcement                                                      │ │
│  │  • Security headers middleware                                           │ │
│  │  • Request validation                                                    │ │
│  │  • Authentication middleware                                             │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│                                    │                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                         APPLICATION LAYER (Go)                           │ │
│  │  • Session validation                                                    │ │
│  │  • Authorization checks                                                  │ │
│  │  • Input sanitization                                                    │ │
│  │  • Business logic                                                        │ │
│  │  • Secure error handling                                                 │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│                                    │                                          │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                         DATA LAYER (PostgreSQL)                          │ │
│  │  • Parameterized queries (sqlc)                                          │ │
│  │  • TLS encrypted connections                                             │ │
│  │  • Encrypted at rest (provider)                                          │ │
│  │  • Hashed passwords (Argon2id)                                           │ │
│  │  • Hashed session tokens                                                 │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│                                                                               │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 14. Delegation to Development Team

**DELEGATE:**
- **architect**: Review and finalize session management architecture; confirm encrypted local storage approach for PWA offline mode
- **senior_dev**: Implement all Phase 1 security controls including:
  - Argon2id password hashing with specified parameters
  - Session-based auth with secure cookie configuration
  - Rate limiting middleware (sliding window algorithm)
  - Input validation middleware following the validation rules table
  - Authorization middleware with IDOR prevention
  - Security headers middleware
  - Secure error handling (generic messages, internal logging)
  - Parameterized queries via sqlc

---

## 15. Security Acceptance Criteria

Before MVP launch, the following must be verified:

### Authentication
- [ ] Passwords hashed with Argon2id (verified parameters)
- [ ] Session cookies have HttpOnly, Secure, SameSite=Strict flags
- [ ] Login rate limiting functional (tested with 11+ attempts)
- [ ] Account lockout triggers after 5 failed attempts
- [ ] Password reset tokens expire in 15 minutes
- [ ] Generic error messages on auth failures

### Authorization
- [ ] All API endpoints require valid session
- [ ] IDOR testing passed (cannot access other users' data)
- [ ] UUIDs used for all resource identifiers

### Input Validation
- [ ] All inputs validated server-side
- [ ] SQL injection testing passed (sqlmap or manual)
- [ ] XSS testing passed (reflected and stored)
- [ ] Invalid inputs return 400, not 500

### Headers & Transport
- [ ] SSL Labs grade A or higher
- [ ] All security headers present
- [ ] CORS only allows production origin

### Monitoring
- [ ] Failed auth attempts logged
- [ ] Security events trigger alerts

---

## Summary

This security specification provides comprehensive protection for SubTrack while maintaining the privacy-first approach that differentiates the product. By avoiding bank connections, we reduce our attack surface significantly, but we must still protect user credentials and subscription data with industry-standard controls.

**Key Security Principles Applied**:
1. **Defense in depth** - Multiple security layers
2. **Least privilege** - Access only what's needed
3. **Fail securely** - Generic errors, secure defaults
4. **Don't trust client input** - Validate everything server-side
5. **Privacy by design** - Minimal data collection, user control

---

PHASE_COMPLETE: planning
