# Security Specification: Brutalist Design System Reskin

**Role:** Security Expert
**Date:** 2026-03-16
**Phase:** Planning (Brutalist Reskin)
**Risk Level:** LOW (visual-only change, no functional modifications)

---

## 1. Change Scope & Threat Impact Assessment

This task is a **CSS/visual-only reskin** — transforming the "Warm Dark Craft" aesthetic to a brutalist design language. The JavaScript logic, data model, DOM structure, and all functional behavior remain **identical**.

### What Changes

| Layer | Changes? | Security Impact |
|-------|----------|-----------------|
| CSS (styles.css) | Yes — full restyle | None (CSS is not an attack vector in this context) |
| HTML structure | Minimal — possible class/attribute tweaks | Low — review for innerHTML/DOM changes |
| Font loading | Yes — switching to monospace system fonts or different CDN fonts | Low — CSP policy may need updating |
| JavaScript logic | No | None |
| Data model / state | No | None |
| Clipboard handling | No | None |
| Input validation | No | None |
| Export features | No | None |

### What Does NOT Change

All security controls from the original spec (R1–R10) remain in full effect. This reskin does not alter:

- Color value validation (`isValidHex()`)
- DOM manipulation patterns (`textContent` only, no `innerHTML`)
- CSP meta tag (may need font-src adjustment only)
- Clipboard API usage (copy from data model)
- Keyboard event handling
- URL parameter sanitization (if implemented)
- localStorage validation (if implemented)

---

## 2. Threat Model for the Reskin

### 2.1 Risk: CSP Policy Drift During Font Changes

**Severity: Medium | Likelihood: Medium | Priority: HIGH**

The brutalist design direction likely changes typography from Google Fonts (Inter, JetBrains Mono) to system monospace fonts or a different font stack. This is a **security opportunity** if system fonts replace CDN-loaded fonts.

**Scenarios:**

| Font Strategy | CSP Impact | Security Effect |
|---------------|------------|-----------------|
| System fonts only (`monospace`, `system-ui`) | Remove `fonts.googleapis.com` and `fonts.gstatic.com` from CSP | **Positive** — eliminates CDN dependency entirely |
| Self-hosted fonts | Remove external font domains from CSP | **Positive** — no third-party font loading |
| Different Google Fonts | No CSP change needed | **Neutral** — same CDN trust level |
| New font CDN (non-Google) | Must add new domain to `font-src` | **Negative** — new third-party dependency |

**Recommendation:**
Brutalist design naturally favors system monospace fonts (`Courier New`, `monospace`, `ui-monospace`). If the UI spec uses system fonts, **tighten the CSP** by removing Google Fonts domains:

```html
<!-- BEFORE (current) -->
<meta http-equiv="Content-Security-Policy"
      content="default-src 'self';
               style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;
               font-src 'self' https://fonts.gstatic.com;
               script-src 'self';
               img-src 'self' data:;
               connect-src 'none';">

<!-- AFTER (if system fonts only — preferred) -->
<meta http-equiv="Content-Security-Policy"
      content="default-src 'self';
               style-src 'self' 'unsafe-inline';
               font-src 'self';
               script-src 'self';
               img-src 'self' data:;
               connect-src 'none';">
```

This is a **net security improvement** — fewer external domains in CSP = smaller attack surface.

**Action required:** After UI spec finalizes the font stack, update the CSP `style-src` and `font-src` directives accordingly. If Google Fonts are no longer loaded, remove those CDN `<link>` tags AND `<link rel="preconnect">` tags from `index.html`.

### 2.2 Risk: Accidental innerHTML Introduction During Refactoring

**Severity: High | Likelihood: Low | Priority: HIGH**

When developers restyle components, there's a risk of refactoring DOM manipulation code alongside CSS changes — especially if HTML structure changes to support the brutalist layout (e.g., adding wrapper divs, changing swatch markup).

**Mitigation:**
- **Rule: Do not modify `app.js` or `colors.js` for visual changes.** The reskin should be achievable through `styles.css` and minor `index.html` structural changes only.
- If HTML structure changes are needed (new elements, changed classes), the rendering code in `app.js` `render()` function must continue to use `textContent` exclusively.
- **Code review gate:** Any PR touching `.js` files during this reskin must be explicitly justified and reviewed for `innerHTML`, `outerHTML`, `document.write()`, or `eval()` usage.

**Verification:**
```bash
# Run after reskin implementation — must return zero results
grep -rn 'innerHTML\|outerHTML\|document\.write\|eval(' js/
```

### 2.3 Risk: SVG/Data URI Injection via CSS

**Severity: Low | Likelihood: Very Low | Priority: LOW**

Brutalist design may use CSS `background-image` with data URIs or inline SVGs for decorative elements (borders, patterns, background textures). The current codebase already uses a data URI for the harmony select dropdown arrow.

**Mitigation:**
- All CSS data URIs must be static, hardcoded values — never constructed from user input or dynamic data.
- The existing pattern in `styles.css` (line ~441: `background-image: url("data:image/svg+xml,...")` for the dropdown arrow) is safe because it's a static SVG string.
- **Rule:** No JavaScript should ever construct CSS `url()` values from external input. All `background-image` data URIs must be hardcoded in the stylesheet.
- This is already the case and should remain so.

### 2.4 Risk: Reduced Visual Feedback Masking Security States

**Severity: Low | Likelihood: Low | Priority: MEDIUM**

Brutalist design uses minimal visual decoration. If visual feedback for security-relevant states is stripped (e.g., focus indicators, locked-state indicators), it could degrade the user experience without being a direct vulnerability — but poor focus visibility can affect keyboard-only users who rely on visual focus cues.

**Mitigation:**
- **Focus indicators must remain visible.** Brutalist focus states should use thick, high-contrast borders (e.g., `outline: 3px solid #000` on white, or `3px solid #FFF` on black) rather than the current subtle amber ring. This is actually a brutalist UX *improvement* — bold, raw focus indicators.
- **Lock state must remain clearly distinguishable.** The locked swatch indicator should be visually obvious — brutalist design supports this well (e.g., thick border, bold "LOCKED" text overlay, cross-hatching pattern).
- **Toast/copy feedback must remain.** The copy confirmation toast is a functional requirement, not decorative. It must persist through the reskin even if restyled.

### 2.5 Opportunity: Elimination of External Dependencies (CDN)

**Severity: N/A (positive) | Priority: MEDIUM**

If the brutalist reskin removes Google Fonts in favor of system fonts:

1. **Remove `<link rel="preconnect" href="https://fonts.googleapis.com">`** from `index.html`
2. **Remove `<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>`** from `index.html`
3. **Remove `<link href="https://fonts.googleapis.com/css2?family=..." rel="stylesheet">`** from `index.html`
4. **Tighten CSP** as described in 2.1
5. **Result:** Zero external network requests. The app becomes fully self-contained — no CDN supply chain risk whatsoever.

This makes the application **more secure** than its current state. Zero external dependencies is the gold standard for client-side security.

---

## 3. Security Requirements for Brutalist Reskin

### 3.1 CRITICAL — Non-Negotiable (Carried Forward)

These requirements are unchanged from the original security spec. They must NOT be regressed during the reskin:

| # | Requirement | Verification |
|---|-------------|-------------|
| R1 | All color values validated with `isValidHex()` before DOM insertion | `grep -n 'isValidHex' js/colors.js` — must exist |
| R2 | `textContent` only — no `innerHTML` with dynamic data | `grep -rn 'innerHTML' js/` — must return zero results |
| R3 | CSP meta tag present in `<head>` | Inspect `index.html` line 6-12 |
| R4 | No `eval()`, `document.write()`, `new Function()` | `grep -rn 'eval(\|document\.write\|new Function' js/` — zero results |
| R5 | Clipboard copies from `state.palette[i]`, not DOM | Review `app.js` `handleSwatchClick()` and `handleCopyAll()` |

### 3.2 HIGH — Must Address During Reskin

| # | Requirement | Action |
|---|-------------|--------|
| R6 | Update CSP `font-src` and `style-src` if font loading changes | Match CSP to actual font sources used |
| R7 | Remove unused CDN `<link>` tags if switching to system fonts | Delete preconnect and font stylesheet links |
| R8 | No new JavaScript files or `<script>` tags introduced | HTML review |
| R9 | If HTML structure changes, verify `app.js` selectors still work correctly | Functional test — generate, copy, lock all still work |

### 3.3 MEDIUM — Should Verify

| # | Requirement | Action |
|---|-------------|--------|
| R10 | Focus indicators remain visible and high-contrast | Visual/keyboard test: Tab through all elements |
| R11 | Lock state remains visually distinguishable | Visual test: lock a color, verify it's obvious |
| R12 | Toast notifications still function | Test: click a swatch, verify toast appears |
| R13 | `prefers-reduced-motion` media query still respected | Test with motion preference set |

---

## 4. Security Testing Checklist for Reskin

Before merging the brutalist reskin, verify:

### Automated / Code Review Checks
- [ ] `grep -rn 'innerHTML' js/` returns zero results
- [ ] `grep -rn 'eval(' js/` returns zero results
- [ ] `grep -rn 'document\.write' js/` returns zero results
- [ ] CSP meta tag present and `font-src` matches actual font sources
- [ ] No new `<script>` tags added to `index.html`
- [ ] No new external CDN `<link>` or `<script>` references added
- [ ] If Google Fonts removed, preconnect links also removed

### Functional Security Tests
- [ ] Generate palette → colors display correctly (validation working)
- [ ] Click swatch → clipboard contains only hex/RGB/HSL value (no HTML)
- [ ] Keyboard shortcuts (Space, L, C, 1-5) all still work
- [ ] Lock/unlock persists across generation cycles
- [ ] Open DevTools Console → zero CSP violation errors during normal use
- [ ] Tab through all interactive elements → visible focus indicator on each

### Regression Tests (Same as Original)
- [ ] Navigate to `?colors=<script>alert(1)</script>` → no script execution (if URL sharing exists)
- [ ] Manually set `localStorage.palettes` to malicious payload → reload → no execution (if localStorage used)

---

## 5. Risk Summary

| # | Issue | Severity | Likelihood | Priority | Status |
|---|-------|----------|------------|----------|--------|
| 2.1 | CSP policy drift during font changes | Medium | Medium | **High** | Must update CSP when fonts change |
| 2.2 | Accidental innerHTML during refactoring | High | Low | **High** | Code review gate on JS changes |
| 2.3 | SVG/data URI in CSS | Low | Very Low | **Low** | Keep data URIs static/hardcoded |
| 2.4 | Reduced visual feedback for states | Low | Low | **Medium** | Verify focus + lock + toast visibility |
| 2.5 | CDN elimination opportunity | N/A | N/A | **Medium** | Net positive — pursue if system fonts used |

---

## 6. Overall Security Assessment

**Risk Level: LOW.** This is the safest type of change — a CSS-only visual reskin with no functional modifications. The attack surface does not grow; it may actually shrink if CDN font dependencies are eliminated.

**Top 3 priorities for this reskin:**
1. **Update CSP to match new font sources** — if fonts change, CSP must change with them
2. **Verify no innerHTML introduced** — refactoring is the #1 risk vector for regressions
3. **Preserve all functional security controls** — validation, textContent, clipboard safety

**Security bonus:** If the brutalist design uses system monospace fonts (which it should — brutalism favors raw, native rendering), this reskin eliminates the last external dependency (Google Fonts), making the application a **zero-network-request, fully self-contained** client-side tool. This is the strongest possible security posture for a static web application.

---

## DELEGATE:
- **architect**: Confirm whether the brutalist font stack uses system fonts or requires external font loading. This determines whether CSP can be tightened (preferred) or must add new font-src domains. Ensure no new external dependencies are introduced.
- **senior_dev**: When implementing the reskin: (1) Do not modify `app.js` or `colors.js` unless absolutely necessary for structural HTML changes. (2) If JS changes are needed, verify zero `innerHTML` usage. (3) Update the CSP meta tag `font-src` and `style-src` to match the final font loading strategy. (4) Remove unused Google Fonts `<link>` tags if switching to system fonts.

---

PHASE_COMPLETE: planning
