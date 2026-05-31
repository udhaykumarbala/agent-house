# Security Specification — WorkTime Tracker

**Role:** Security Expert
**Date:** 2026-03-18
**Phase:** Planning
**Based on:** Research from `security.md`, `architecture-research.md`, `product-research.md`

---

## 1. Threat Assessment Summary

This is a **client-side only, single-user application** with localStorage persistence. No backend, no authentication, no network requests. The attack surface is narrower than a full-stack app but **client-side security discipline must be rigorous** — there is no server-side safety net.

### Threat Model

| Threat Actor | Motivation | Vector | Impact | Likelihood |
|---|---|---|---|---|
| XSS via crafted input | Data theft / manipulation | Employee names, project names, notes, task descriptions | **High** — full data access | Medium |
| Malicious browser extension | Data exfiltration | Read localStorage, inject scripts | High | Low |
| Physical access attacker | View employee data | Open browser, read localStorage via DevTools | Medium | Medium (shared computers) |
| Social engineering / self-XSS | Trick user into running JS | Console paste attacks, crafted CSV imports | Medium | Low |
| CSV formula injection | Execute commands on user's machine | Crafted field values exported to CSV, opened in Excel | **High** | Medium |

### Data Classification

| Data Type | Sensitivity | Protection Required |
|---|---|---|
| Employee names | **Medium** (PII) | Output encoding, clear-data option |
| Work hours / clock in-out | **Medium** | Validation, integrity checks |
| Project/task names | Low-Medium | Output encoding |
| Report exports (CSV) | **Medium** | Formula injection protection |
| App configuration | Low | Schema validation |

---

## 2. Security Requirements

### 2.1 XSS Prevention — CRITICAL

**Requirement:** All user-supplied data MUST be safely rendered to the DOM. No user input may ever be interpreted as HTML or JavaScript.

**Implementation Rules:**

| Rule | Detail |
|---|---|
| **R-XSS-01** | Use `element.textContent` for ALL user-supplied data. NEVER use `innerHTML`, `outerHTML`, or `insertAdjacentHTML` with unsanitized user data |
| **R-XSS-02** | Create a centralized `escapeHTML()` utility in `src/utils/sanitize.js` for any case where template literal HTML rendering is unavoidable |
| **R-XSS-03** | NEVER use `eval()`, `Function()`, `document.write()`, or `setTimeout/setInterval` with string arguments containing user data |
| **R-XSS-04** | For dynamic element creation, use `document.createElement()` + set properties via `textContent`, `setAttribute()`, etc. |
| **R-XSS-05** | If a sanitization library is ever needed, use DOMPurify (~7KB). Do not build a custom HTML sanitizer |

**`escapeHTML()` specification:**

```
Input:  any string
Output: string with &, <, >, ", ' replaced with HTML entities
Chars:  & → &amp;  < → &lt;  > → &gt;  " → &quot;  ' → &#x27;
```

This utility is a **fallback**, not the primary defense. `textContent` is the primary defense.

**Affected fields (all user-editable text):**
- Employee: `name`, `email`, `role`, `department`
- Project: `name`, `description`
- Task: `name`, `description`
- TimeEntry: `notes`
- Any future text fields

### 2.2 Input Validation — HIGH

**Requirement:** All user inputs MUST be validated before storage. Validation happens at the service layer (data services), not the view layer.

**Validation Rules:**

| Field | Validation | Max Length |
|---|---|---|
| Employee name | Required, non-empty after trim, no control characters | 100 chars |
| Employee email | Required, valid email format (regex), unique across employees | 254 chars |
| Employee role | Optional, alphanumeric + spaces | 50 chars |
| Employee department | Optional, alphanumeric + spaces | 50 chars |
| Project name | Required, non-empty after trim, unique across active projects | 100 chars |
| Project description | Optional | 500 chars |
| Task name | Required, non-empty after trim, unique within project | 100 chars |
| Task description | Optional | 500 chars |
| TimeEntry notes | Optional | 500 chars |
| TimeEntry clockIn | Required, valid ISO 8601 timestamp, not in the future beyond current day | — |
| TimeEntry clockOut | Valid ISO 8601 timestamp, MUST be after clockIn, not in the future beyond current day | — |
| TimeEntry duration | Computed value, must be 0–1440 minutes (0–24 hours) | — |
| Hourly rate | Optional, number, 0–9999.99 | — |
| Color (hex) | Must match `/^#[0-9A-Fa-f]{6}$/` | 7 chars |
| UUID fields | Must match UUID v4 format | 36 chars |

**Validation utility location:** `src/utils/validation.js`

**Validation pattern:**
```
User Input → trim() → check required → check length → check format → check business rules → store
```

**Business rule validations:**
- TimeEntry: `clockOut > clockIn` (no zero or negative duration)
- TimeEntry: No overlapping entries for the same employee on the same day
- TimeEntry: Duration between clockIn and clockOut must not exceed 24 hours
- Employee: Email must be unique among active employees
- Project: Name must be unique among non-archived projects

### 2.3 CSV Export Security — CRITICAL

**Requirement:** All CSV exports MUST be protected against formula injection (CSV injection / DDE injection).

**Implementation Rules:**

| Rule | Detail |
|---|---|
| **R-CSV-01** | Every cell value starting with `=`, `+`, `-`, `@`, `\t`, `\r`, `\n` MUST be prefixed with a single quote character (`'`) |
| **R-CSV-02** | All cell values MUST be wrapped in double quotes |
| **R-CSV-03** | Double quotes within values MUST be escaped as `""` |
| **R-CSV-04** | Use UTF-8 BOM (`\uFEFF`) at the start of the CSV file for proper Excel encoding |
| **R-CSV-05** | Sanitization MUST happen in the export utility, not in the display layer — export-time sanitization is mandatory regardless of display sanitization |

**CSV sanitize function specification (in `src/utils/export.js`):**

```
function sanitizeCSVCell(value):
  1. Convert to string
  2. If first char is =, +, -, @, \t, \r, \n → prepend with '
  3. Replace all " with ""
  4. Wrap in double quotes
  5. Return result
```

### 2.4 localStorage Security — HIGH

**Requirement:** The localStorage wrapper MUST handle all edge cases securely and gracefully.

**Implementation Rules:**

| Rule | Detail |
|---|---|
| **R-LS-01** | All `JSON.parse()` calls MUST be wrapped in try/catch. On parse failure: log warning, return default empty state, do NOT crash |
| **R-LS-02** | All `localStorage.setItem()` calls MUST catch `QuotaExceededError`. On quota exceeded: show user-visible warning, suggest data export |
| **R-LS-03** | On app load, validate stored data against expected schema. Reject entries with unexpected types or missing required fields |
| **R-LS-04** | Use `wtt_` prefix for all localStorage keys to avoid collisions |
| **R-LS-05** | Provide a "Clear All Data" button in Settings that calls `localStorage` removal for all `wtt_*` keys with user confirmation |
| **R-LS-06** | Monitor storage usage: warn user when approaching 80% of estimated 5MB limit |
| **R-LS-07** | On data load, validate referential integrity: TimeEntry.employeeId must reference an existing Employee, TimeEntry.projectId must reference an existing Project |

**Schema validation on load (in `src/data/store.js`):**

Each entity array read from localStorage must be validated:
- Is it a valid JSON array?
- Does each item have the required fields with correct types?
- Are all ID references valid?
- Are date strings valid ISO 8601?
- Are numeric values within expected ranges?

Invalid entries should be quarantined (logged, excluded from active data) — not silently dropped or allowed to corrupt the UI.

### 2.5 Data Integrity — MEDIUM

**Requirement:** Detect and handle data tampering or corruption gracefully.

**Implementation Rules:**

| Rule | Detail |
|---|---|
| **R-DI-01** | Store a lightweight checksum alongside each collection (e.g., entry count + last-modified timestamp) to detect external modifications |
| **R-DI-02** | Log all create/update/delete operations with timestamps in a `wtt_audit_log` key (capped at last 500 entries to control size) |
| **R-DI-03** | On data load, if checksum mismatch is detected, show a non-blocking warning: "Data may have been modified externally" |
| **R-DI-04** | Validate that all time entries have logically consistent data (clockOut > clockIn, duration matches timestamps, dates are real calendar dates) |

### 2.6 Content Security Policy — MEDIUM

**Requirement:** If the app is served from a web server (even a local dev server), set appropriate CSP headers.

**Recommended CSP:**

```
default-src 'self';
script-src 'self';
style-src 'self' 'unsafe-inline';
img-src 'self' data:;
font-src 'self' https://fonts.gstatic.com;
connect-src 'none';
object-src 'none';
base-uri 'self';
form-action 'none';
```

**Notes:**
- `'unsafe-inline'` for styles is needed for Tailwind's runtime styles. If Tailwind is fully compiled at build time, this can be removed.
- `connect-src 'none'` — the app makes no network requests. This prevents any XSS payload from exfiltrating data.
- `form-action 'none'` — no forms submit to a server.
- CSP should be configured in `index.html` via `<meta>` tag for static hosting compatibility.

### 2.7 Dependency Security — LOW

**Requirement:** Minimize supply chain risk from third-party packages.

| Rule | Detail |
|---|---|
| **R-DEP-01** | Commit `package-lock.json` to version control |
| **R-DEP-02** | Run `npm audit` during development. Address critical/high vulnerabilities before any deployment |
| **R-DEP-03** | Minimize runtime dependencies — prefer vanilla JS for CSV export, date formatting, UUID generation |
| **R-DEP-04** | Build-time-only dependencies (Vite, Tailwind, PostCSS) are lower risk but should still be version-pinned |
| **R-DEP-05** | If DOMPurify or any runtime dependency is added, document the justification |

---

## 3. Secure Coding Standards

### DOM Rendering Checklist

Every view/component that renders user data MUST follow this pattern:

```
SAFE:
  element.textContent = userData
  element.setAttribute('title', userData)   // safe for non-event attributes
  document.createElement('div').textContent = userData

UNSAFE (NEVER with user data):
  element.innerHTML = userData
  element.outerHTML = userData
  element.insertAdjacentHTML('beforeend', userData)
  document.write(userData)
  eval(userData)
  new Function(userData)
```

### localStorage Access Checklist

Every localStorage read/write MUST follow this pattern:

```
READ:
  1. try { JSON.parse(localStorage.getItem(key)) } catch { return defaultValue }
  2. Validate schema of parsed data
  3. Return validated data or default

WRITE:
  1. Validate data to be stored
  2. try { localStorage.setItem(key, JSON.stringify(data)) }
  3. catch (QuotaExceededError) { notify user, suggest export }
  4. Update checksum
```

### Input Handling Pipeline

Every user input goes through this pipeline before storage:

```
1. TRIM      → Remove leading/trailing whitespace
2. REQUIRED  → Check non-empty for required fields
3. LENGTH    → Enforce max length limits
4. FORMAT    → Validate pattern (email, date, hex color, etc.)
5. BUSINESS  → Check business rules (uniqueness, range, referential integrity)
6. STORE     → Write to localStorage via service layer
```

Validation errors MUST be shown inline next to the form field, with clear error messages. Do NOT use `alert()`.

---

## 4. Security Testing Plan

### Manual Tests

| # | Test | Steps | Expected Result | Priority |
|---|---|---|---|---|
| T-01 | XSS in employee name | Enter `<script>alert(1)</script>` as name, save, reload | Text displayed literally, no script execution | **Critical** |
| T-02 | XSS in project name | Enter `<img src=x onerror=alert(1)>` as project name, save | Text displayed literally, no image tag rendered | **Critical** |
| T-03 | XSS in time entry notes | Enter `"><svg onload=alert(1)>` in notes, save | Text displayed literally | **Critical** |
| T-04 | XSS in all remaining fields | Test task name, task description, project description, employee role, department with XSS payloads | No execution in any field | **Critical** |
| T-05 | Stored XSS via localStorage | Manually set `wtt_employees` in DevTools with XSS payload, reload app | Payload does not execute, data is escaped or rejected | **Critical** |
| T-06 | CSV formula injection | Set employee name to `=CMD("calc")`, export CSV, open in Excel | Formula does NOT execute, prefixed with `'` | **Critical** |
| T-07 | CSV injection variants | Test `+`, `-`, `@`, `\t` prefixed values in exports | All dangerous prefixes neutralized | **Critical** |
| T-08 | Input max length | Enter 10,000+ character string in name field | Truncated or rejected at 100 chars | High |
| T-09 | Special characters | Enter `<>"'&\`` in all fields | Stored and displayed correctly, no rendering issues | High |
| T-10 | Unicode edge cases | Enter emoji, RTL text, zero-width chars, null bytes | Handled gracefully, no crashes | High |
| T-11 | localStorage corruption | Set `wtt_employees` to `{invalid json`, reload | App loads with empty state, shows warning | High |
| T-12 | localStorage quota | Fill localStorage near limit, attempt save | Graceful error message, no crash, suggest export | Medium |
| T-13 | Overlapping time entries | Create two entries for same employee with overlapping times | Rejected with clear error message | High |
| T-14 | Negative duration | Set clockOut before clockIn via manual entry | Rejected with validation error | High |
| T-15 | Future date entry | Enter a time entry for next week | Rejected (entries limited to current day and past) | Medium |
| T-16 | Clear all data | Use "Clear All Data" in settings | All `wtt_*` keys removed, app resets to empty state | Medium |
| T-17 | Referential integrity | Delete an employee, check their time entries | Entries handled gracefully (orphaned entries shown with "Unknown Employee" or filtered) | Medium |
| T-18 | Console output | Check browser console in production build | No sensitive data (employee names, hours) logged | Low |

### Automated Test Cases (unit tests)

| Module | Test | Assertion |
|---|---|---|
| `sanitize.js` | `escapeHTML('<script>')` | Returns `&lt;script&gt;` |
| `sanitize.js` | `escapeHTML('normal text')` | Returns unchanged |
| `sanitize.js` | `escapeHTML('"quotes" & \'apostrophe\'')` | Returns properly escaped |
| `export.js` | `sanitizeCSVCell('=SUM(A1)')` | Returns `"'=SUM(A1)"` |
| `export.js` | `sanitizeCSVCell('+cmd')` | Returns `"'+cmd"` |
| `export.js` | `sanitizeCSVCell('normal')` | Returns `"normal"` |
| `export.js` | `sanitizeCSVCell('value with "quotes"')` | Returns `"value with ""quotes"""` |
| `validation.js` | Validate employee name > 100 chars | Returns error |
| `validation.js` | Validate empty required field | Returns error |
| `validation.js` | Validate email format | Accepts valid, rejects invalid |
| `validation.js` | Validate clockOut before clockIn | Returns error |
| `validation.js` | Validate hex color format | Accepts `#FF0000`, rejects `red`, `#GGG` |
| `store.js` | Parse corrupted JSON | Returns default, no throw |
| `store.js` | Handle QuotaExceededError | Returns error status, no throw |
| `store.js` | Schema validation rejects bad data | Malformed entries excluded |

---

## 5. Risk Register

| # | Issue | Severity | Likelihood | Impact | Mitigation | Status |
|---|---|---|---|---|---|---|
| SEC-01 | XSS via user input in DOM | **Critical** | Medium | Full data access/manipulation | Use `textContent`, `escapeHTML()` utility, never `innerHTML` with user data | Spec'd |
| SEC-02 | CSV formula injection | **Critical** | Medium | Command execution on user machine | Prefix dangerous chars, quote all cells | Spec'd |
| SEC-03 | Missing input validation | **High** | High | Data corruption, XSS enablement | Centralized validation in service layer | Spec'd |
| SEC-04 | localStorage parse failure | **High** | Medium | App crash on load | Try/catch all JSON.parse, return defaults | Spec'd |
| SEC-05 | localStorage quota exceeded | **Medium** | Low | Data loss, save failure | Catch error, warn user, suggest export | Spec'd |
| SEC-06 | Shared computer data exposure | **Medium** | Medium | PII visible to next user | "Clear All Data" button, warn on shared use | Spec'd |
| SEC-07 | Data tampering via DevTools | **Low-Medium** | Low | Inflated/deflated hours | Checksum validation, audit log | Spec'd |
| SEC-08 | Dependency supply chain | **Low** | Low | Build-time compromise | Lock files, npm audit, minimize deps | Spec'd |
| SEC-09 | No CSP headers | **Medium** | Low | XSS payload can exfiltrate data | Meta tag CSP in index.html | Spec'd |

---

## 6. Implementation Priority

The development team should implement security features in this order:

### Phase 1 — Build into foundation (before any UI work)
1. **`src/utils/sanitize.js`** — `escapeHTML()` utility
2. **`src/utils/validation.js`** — All input validation functions
3. **`src/data/store.js`** — Secure localStorage wrapper with try/catch, schema validation, quota handling
4. **`src/utils/export.js`** — CSV cell sanitization function

### Phase 2 — Integrate during view development
5. Every view uses `textContent` for user data rendering (code review checkpoint)
6. Every form uses validation pipeline before calling service layer
7. Export functions use CSV sanitization

### Phase 3 — Hardening
8. CSP meta tag in `index.html`
9. "Clear All Data" button in settings
10. Storage usage monitoring and warning
11. Audit log implementation
12. Data integrity checksums

---

## 7. OWASP Top 10 Compliance Matrix

| OWASP Category | Relevant? | Our Mitigation | Status |
|---|---|---|---|
| A01: Broken Access Control | Low (single user) | "Clear All Data" for shared computers | Spec'd |
| A02: Cryptographic Failures | Low (no crypto needed) | N/A — no passwords, no secrets | N/A |
| A03: Injection | **HIGH** | XSS prevention via textContent, input validation, CSV sanitization | Spec'd |
| A04: Insecure Design | **Medium** | Security built into data layer by design, validation pipeline | Spec'd |
| A05: Security Misconfiguration | Medium | CSP headers, no debug output in production | Spec'd |
| A06: Vulnerable Components | Low | Minimal deps, lock files, npm audit | Spec'd |
| A07: Auth Failures | N/A | No authentication in scope | N/A |
| A08: Data Integrity Failures | **Medium** | Schema validation, checksums, audit log | Spec'd |
| A09: Logging & Monitoring | Low | Client-side audit trail for data changes | Spec'd |
| A10: SSRF | N/A | No server-side, no outbound requests | N/A |

---

## 8. Delegation

```
DELEGATE:
- architect: Ensure store.js localStorage wrapper includes try/catch, schema validation, and quota handling as foundational architecture. CSP meta tag in index.html. Data service layer must enforce validation before any write.
- senior_dev: Implement sanitize.js (escapeHTML), validation.js (all input validators), CSV sanitization in export.js. These are security-critical utilities — they must exist before any view code is written. Every PR that renders user data must be reviewed for innerHTML usage.
```

---

**Summary:** This application's #1 security risk is XSS through user input rendered to the DOM, followed by CSV formula injection in exports. Both are **cheap to prevent if built in from the start** — the `escapeHTML()` utility, `textContent`-only rendering rule, and CSV sanitization function are small pieces of code that prevent the most damaging attacks. The localStorage wrapper with proper error handling prevents crashes and data corruption. These four utilities (`sanitize.js`, `validation.js`, `export.js`, `store.js`) form the security foundation and MUST be implemented before any UI views.

---

PHASE_COMPLETE: planning
