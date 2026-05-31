# UX Specification: UnitShift — Premium Minimal Unit Converter

Based on research in: .plans/research/ux-research.md, .plans/research/ui-research.md
Implementing features from: .plans/research/product-research.md
Architecture context: .plans/research/architecture-research.md
Security context: .plans/research/security_expert.md

---

## Design Philosophy

**"Typography as interface, restraint as luxury."**

Every element earns its place. The converted number is the hero — everything else supports it. Inspired by Elk's gesture-driven minimalism, Apple Calculator's clarity, and Linear's polish. The app should feel like a precision instrument: a Braun calculator reimagined for the web.

---

## Screen List

1. **Main Converter** — Single-screen layout housing category tabs, input/output fields, unit selectors, and swap control. This IS the app — no other pages.
2. **Empty/Welcome State** — First-launch variant of the main screen with a warm prompt instead of zeros.
3. **Error State** — Inline variant when edge-case inputs occur (e.g., below absolute zero).

There is intentionally no settings page, no history page, no about page. One screen, one purpose.

---

## Wireframes

### Screen 1: Main Converter (Primary State)

```
┌──────────────────────────────────────────┐
│                                          │
│     [ Length ]  [ Weight ]  [ Temp ]     │  ← Category tabs (segmented pill)
│     ─────────                            │     Active tab has sliding pill bg
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  FROM                              │  │
│  │                                    │  │
│  │  ┌──────┐ ┌──────┐ ┌──────┐ ···   │  │  ← Unit chips (scrollable row)
│  │  │  km  │ │  mi  │ │  m   │       │  │     Active chip = filled accent
│  │  └──────┘ └──────┘ └──────┘       │  │
│  │                                    │  │
│  │         42.5                       │  │  ← Input value (editable, 48px)
│  │         Kilometers                 │  │  ← Unit label (muted, 14px)
│  │                                    │  │
│  └────────────────────────────────────┘  │
│                                          │
│              ┌──────┐                    │
│              │  ⇅   │                    │  ← Swap button (rotates 180° on tap)
│              └──────┘                    │
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  TO                                │  │
│  │                                    │  │
│  │  ┌──────┐ ┌──────┐ ┌──────┐ ···   │  │  ← Unit chips (scrollable row)
│  │  │  km  │ │  mi  │ │  m   │       │  │     Active chip = filled accent
│  │  └──────┘ └──────┘ └──────┘       │  │
│  │                                    │  │
│  │         26.4097                    │  │  ← Result value (live, 56px bold)
│  │         Miles                      │  │  ← Unit label (muted, 14px)
│  │                                    │  │
│  └────────────────────────────────────┘  │
│                                          │
│     1 km = 0.6214 mi                    │  ← Quick reference (12px, muted)
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  7  │  8  │  9  │                  │  │
│  ├─────┼─────┼─────┤                  │  │  ← Custom numeric keypad
│  │  4  │  5  │  6  │                  │  │     (mobile only, hidden on desktop)
│  ├─────┼─────┼─────┤                  │  │     Large 48px touch targets
│  │  1  │  2  │  3  │                  │  │     Clean, minimal styling
│  ├─────┼─────┼─────┤                  │  │
│  │  .  │  0  │  ⌫  │                  │  │
│  └────────────────────────────────────┘  │
│                                          │
└──────────────────────────────────────────┘
```

**Key Elements**:
- **Category Tabs**: Segmented control with sliding pill indicator. Three tabs only: Length, Weight, Temperature.
- **FROM Card**: Muted surface card containing source unit chips and the editable input value.
- **Swap Button**: Circular button centered between cards. Rotates 180° with spring animation on tap. Swaps both units and values.
- **TO Card**: Slightly elevated surface card containing target unit chips and the live result value. Result is visually dominant (larger, bolder).
- **Quick Reference**: Subtle one-liner showing the base conversion rate (e.g., "1 km = 0.6214 mi"). Updates with unit changes.
- **Numeric Keypad**: Mobile-only custom keypad. Large touch targets, decimal point, backspace. Disappears on desktop (keyboard input replaces it).

### Screen 1b: Main Converter (Desktop Variant)

```
┌──────────────────────────────────────────┐
│                                          │
│     [ Length ]  [ Weight ]  [ Temp ]     │
│     ─────────                            │
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  FROM                              │  │
│  │  [ km ] [ mi ] [ m ] [ ft ] [in]  │  │  ← All chips visible (wider space)
│  │                                    │  │
│  │         42.5|                      │  │  ← Editable inline (cursor visible)
│  │         Kilometers                 │  │
│  └────────────────────────────────────┘  │
│                                          │
│              [ ⇅ ]                       │
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  TO                                │  │
│  │  [ km ] [ mi ] [ m ] [ ft ] [in]  │  │
│  │                                    │  │
│  │         26.4097                    │  │  ← Also editable (bidirectional)
│  │         Miles                      │  │
│  └────────────────────────────────────┘  │
│                                          │
│     1 km = 0.6214 mi                    │
│                                          │
└──────────────────────────────────────────┘
```

**Desktop Differences**:
- No numeric keypad — standard keyboard input
- Both fields editable (bidirectional conversion, per Google Converter pattern)
- All unit chips visible without scrolling (wider viewport)
- Max-width 480px container, centered on screen
- Hover states on all interactive elements

### Screen 2: Empty/Welcome State

```
┌──────────────────────────────────────────┐
│                                          │
│     [ Length ]  [ Weight ]  [ Temp ]     │
│     ─────────                            │
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  FROM                              │  │
│  │  [ km ] [ mi ] [ m ] [ ft ] ···   │  │
│  │                                    │  │
│  │         Type a value...            │  │  ← Placeholder text (muted, pulsing
│  │         Kilometers                 │  │     cursor animation)
│  │                                    │  │
│  └────────────────────────────────────┘  │
│                                          │
│              [ ⇅ ]                       │
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  TO                                │  │
│  │  [ km ] [ mi ] [ m ] [ ft ] ···   │  │
│  │                                    │  │
│  │         —                          │  │  ← Em dash placeholder (not "0")
│  │         Miles                      │  │
│  │                                    │  │
│  └────────────────────────────────────┘  │
│                                          │
└──────────────────────────────────────────┘
```

**Key Decisions**:
- No zeros on launch — feels broken (per UX research anti-pattern)
- "Type a value..." placeholder with a soft pulsing cursor gently invites input
- Result shows an em dash (—) not "0" to indicate "awaiting input"
- Smart defaults per category: km ↔ mi (length), kg ↔ lb (weight), °C ↔ °F (temperature)

### Screen 3: Error / Edge Case State

```
┌────────────────────────────────────────┐
│  ...                                   │
│                                        │
│  ┌──────────────────────────────────┐  │
│  │  TO                              │  │
│  │  [ °C ] [ °F ] [ K ]            │  │
│  │                                  │  │
│  │         -280                     │  │  ← Invalid result (red accent)
│  │         ⚠ Below absolute zero    │  │  ← Inline warning (no modal/alert)
│  │                                  │  │
│  └──────────────────────────────────┘  │
│                                        │
└────────────────────────────────────────┘
```

**Error Handling**:
- Inline messages, never modal alerts or toast popups
- Gentle warning icon + text below the value
- Value still shown (not hidden) so user understands the issue
- Only applies to temperature edge cases (below absolute zero)
- Non-numeric input simply ignored (input field filters to digits, decimal, minus)

---

## User Flows

### Primary Flow: Quick Conversion (3-8 seconds)

1. **User opens app** → sees Main Converter with last-used category and smart defaults (or Welcome State on first visit)
2. **User taps a category tab** (if needed) → sliding pill animates to selected tab, content cross-fades with vertical parallax
3. **User taps FROM input field** → keypad slides up (mobile), field gets subtle focus glow
4. **User types a number** → result updates live in the TO field with rolling digit animation, no delay
5. **User reads the result** → task complete. Result is large, bold, unmissable
6. **User leaves** → app remembers category and unit selections for next visit (localStorage)

### Secondary Flow: Swap Direction

1. User has a conversion showing (e.g., 42.5 km → 26.41 mi)
2. **User taps swap button (⇅)** → button rotates 180° with spring animation
3. FROM and TO units swap positions with cross-fade
4. Values recalculate: 26.41 mi → 42.5 km (or user's current input re-converts)
5. The swap feels physical and satisfying — the signature micro-interaction

### Tertiary Flow: Change Units

1. User sees the unit chip row (e.g., [km] [mi] [m] [ft] [in] [cm] [yd] [mm])
2. **User taps a different unit chip** → chip fills with accent color, previous chip unfills
3. Result recalculates instantly with digit-roll animation
4. Quick reference line updates to reflect new unit pair

### Tertiary Flow: Bidirectional Input (Desktop)

1. User clicks the TO field (result field) instead of FROM
2. TO field becomes editable with cursor
3. **User types in TO field** → FROM value updates live (reverse conversion)
4. Either field drives conversion — whichever was last edited is "source"

---

## Interaction Patterns

| Action | Element | Response | Animation |
|--------|---------|----------|-----------|
| Tap | Category tab | Switch category, load default units | Pill slides to new tab (200ms spring). Content cross-fades (150ms). |
| Tap | Unit chip (FROM or TO) | Select unit, recalculate | Chip fills with accent color (100ms). Result digits roll to new value (300ms). |
| Tap | Swap button (⇅) | Swap FROM/TO units and values | Button rotates 180° (400ms spring with overshoot). Values cross-fade (200ms). |
| Type | Numeric keypad / keyboard | Update input, live convert | Input digit scales up briefly (50ms spring). Result digits roll/morph (100ms delay, 200ms animation). |
| Tap | Backspace (⌫) | Delete last digit, recalculate | Digit shrinks out (100ms). Result updates. |
| Tap | Clear (long-press ⌫) | Clear entire input | All digits shrink out staggered (150ms). Returns to empty state. |
| Tap | Result value | Copy to clipboard (mobile) | Brief flash + checkmark morph on value (300ms). Subtle "Copied" text fades in/out. |
| Hover | Unit chip (desktop) | Visual hover state | Background opacity increases (80ms). |
| Hover | Swap button (desktop) | Highlight | Scale up to 1.05 (100ms). |
| Focus | Input field | Focus indicator | Soft warm glow around card border (200ms fade in). |

---

## States

### Empty State (First Launch)
- **When**: App loads for the first time, or input is cleared
- **Show**: "Type a value..." placeholder in FROM field, em dash (—) in TO field, default unit pair pre-selected
- **Action**: Tapping FROM field (or just typing on desktop) begins input immediately
- **Feel**: Welcoming, not sterile. The placeholder has a soft pulsing cursor animation.

### Active/Converting State
- **When**: User has entered a numeric value
- **Show**: Input value in FROM, live result in TO, quick reference line visible
- **Feel**: Responsive, alive. Every keystroke triggers an immediate result update.

### Loading State
- **Not applicable**: All conversions are local O(1) math — zero loading time. No spinners, no skeleton screens. This is a key UX advantage.

### Error State
- **When**: Temperature result below absolute zero (-273.15°C / -459.67°F / 0K), or invalid input edge case
- **Show**: Result still displayed but with warning icon and inline message below the value. Warning text in a muted warm tone (not alarming red).
- **Action**: User can continue typing to correct the input. Warning disappears when input becomes valid.

### Success State
- **When**: A valid conversion is displayed
- **Show**: The result IS the success state — large, bold, confident. No confirmation needed.
- **Optional**: Tapping the result copies it to clipboard with a brief checkmark flash (subtle confirmation).

---

## Category-Specific Details

### Length
- **Units**: mm, cm, m, km, in, ft, yd, mi
- **Default pair**: km ↔ mi
- **Accent color**: Warm amber (#C7882A at 15% opacity for chip backgrounds, full opacity for active indicators)
- **Quick reference example**: "1 km = 0.6214 mi"

### Weight
- **Units**: mg, g, kg, oz, lb, t
- **Default pair**: kg ↔ lb
- **Accent color**: Soft sage green (#6B8F71 at 15% opacity, full for active)
- **Quick reference example**: "1 kg = 2.2046 lb"

### Temperature
- **Units**: °C, °F, K
- **Default pair**: °C ↔ °F
- **Accent color**: Cool slate blue (#6B8FA3 at 15% opacity, full for active)
- **Quick reference example**: "0°C = 32°F"
- **Special**: Fewer units means chips are centered, not scrollable. Warning for below-absolute-zero values.

---

## Typography Hierarchy

Based on research consensus (ux-research.md, ui-research.md, architecture-research.md):

| Element | Font | Size | Weight | Color |
|---------|------|------|--------|-------|
| Result value (TO) | Space Grotesk | 56px / 3.5rem | 700 (Bold) | #FFFFFF |
| Input value (FROM) | Space Grotesk | 48px / 3rem | 500 (Medium) | #E0E0E0 |
| Unit label below value | Inter | 14px / 0.875rem | 400 (Regular) | #888888 |
| Category tab label | Inter | 14px / 0.875rem | 600 (Semi-bold) | Active: #FFFFFF, Inactive: #666666 |
| Unit chip text | Inter | 13px / 0.8125rem | 500 (Medium) | Active: #FFFFFF, Inactive: #999999 |
| Quick reference | Inter | 12px / 0.75rem | 400 (Regular) | #555555 |
| FROM/TO label | Inter | 11px / 0.6875rem | 600 (Semi-bold) | #555555, uppercase, letter-spacing 0.08em |
| Keypad digits | Space Grotesk | 24px / 1.5rem | 500 (Medium) | #CCCCCC |

**Key typographic decisions**:
- Space Grotesk for all numbers — geometric, distinctive digit forms, immediate visual identity
- Inter for all labels/UI text — neutral, reliable, excellent at small sizes
- Tabular lining figures enabled on Space Grotesk to prevent layout shifts during number animation
- Result value intentionally 8px larger than input value to establish clear visual hierarchy (output > input)

---

## Layout & Spacing System

- **Grid**: 8px base unit
- **Container max-width**: 480px (centered on desktop/tablet)
- **Horizontal padding**: 24px (mobile), 32px (desktop)
- **Card padding**: 24px vertical, 20px horizontal
- **Gap between FROM and TO cards**: 16px (swap button overlaps this gap, centered)
- **Card corner radius**: 16px
- **Chip corner radius**: 20px (pill shape)
- **Swap button**: 48px diameter circle, 16px icon inside
- **Keypad grid**: 3 columns, each key min 64px tall, 8px gap between keys
- **Category tabs height**: 40px
- **Minimum touch target**: 44x44px for all interactive elements

---

## Responsive Behavior

### Mobile (< 640px) — Primary
- Full-viewport layout, no horizontal margins beyond 16px padding
- Custom numeric keypad visible at bottom third of screen
- Unit chips in scrollable horizontal row (overflow-x: auto, no scrollbar visible)
- FROM field auto-focused on load (keypad appears)
- Copy result on tap

### Tablet (640px – 1024px)
- Centered container at 480px max-width
- Keypad still shown but with larger touch targets
- All unit chips may fit without scrolling

### Desktop (> 1024px)
- Centered container at 480px max-width, generous vertical centering
- No numeric keypad — standard keyboard input
- Both FROM and TO fields editable (bidirectional)
- Hover states on chips, swap button, and result (copy affordance)
- Tab key navigates: FROM input → FROM chips → Swap → TO chips → TO input
- Enter key on FROM input focuses TO input (or vice versa)

---

## Accessibility Considerations

### Keyboard Navigation
- **Tab order**: Category tabs → FROM unit chips → FROM input → Swap button → TO unit chips → TO input
- **Arrow keys**: Left/Right to navigate between unit chips within a row
- **Enter/Space**: Activate focused chip or swap button
- **Escape**: Blur/unfocus active input

### Screen Reader Support
- Category tabs: `role="tablist"`, each tab `role="tab"` with `aria-selected`
- Unit chips: `role="radiogroup"` per row, each chip `role="radio"` with `aria-checked`
- FROM input: `aria-label="Value to convert from"`, `aria-describedby` linked to unit label
- TO result: `aria-live="polite"` so screen reader announces new results after conversion
- Swap button: `aria-label="Swap conversion direction"`
- Quick reference: `aria-hidden="true"` (supplementary, not essential)

### Motion & Preferences
- **`prefers-reduced-motion: reduce`**: Disable all spring animations, digit rolling, swap rotation. Use instant state changes instead.
- **`prefers-color-scheme: light`**: Switch to light theme variant (warm off-white background, dark text)
- **`prefers-contrast: more`**: Increase border visibility on cards, use solid chip backgrounds instead of translucent

### Color Contrast
- All text meets WCAG AA (4.5:1 for normal text, 3:1 for large text)
- Active chip text on accent background: verified ≥ 4.5:1
- Muted text (#888 on #1A1714 background): 5.2:1 ratio — passes AA
- Focus indicators: 3px solid ring in accent color (not just color change)

### Touch Targets
- Minimum 44x44px for all interactive elements (per WCAG 2.5.5)
- Keypad keys: 64px tall for comfortable thumb tapping
- Unit chips: min 36px height with 8px margin (effective 44px target)

---

## Persistence (localStorage)

Minimal state stored to improve return visits:

| Key | Value | Purpose |
|-----|-------|---------|
| `unitshift-category` | `"length"` / `"weight"` / `"temperature"` | Remember last-used category |
| `unitshift-length-from` | `"km"` | Last-used FROM unit for length |
| `unitshift-length-to` | `"mi"` | Last-used TO unit for length |
| `unitshift-weight-from` | `"kg"` | Same for weight |
| `unitshift-weight-to` | `"lb"` | Same for weight |
| `unitshift-temp-from` | `"celsius"` | Same for temperature |
| `unitshift-temp-to` | `"fahrenheit"` | Same for temperature |

- No input values persisted (fresh start each session)
- No theme preference stored (follow system `prefers-color-scheme`)
- Per security research: no sensitive data — this is safe for localStorage

---

## Micro-Interaction Specification

### 1. Category Tab Switch
- **Trigger**: Tap on inactive tab
- **Animation**: Pill background slides from current tab to new tab (200ms, spring ease with slight overshoot)
- **Content**: Cross-fade (old content fades out 100ms, new fades in 100ms, slight 4px vertical shift)
- **Reduced motion**: Instant switch, no animation

### 2. Unit Chip Selection
- **Trigger**: Tap on inactive chip
- **Animation**: Previous chip: background fades to transparent (100ms). New chip: background fills with accent color (100ms). Result value rolls to new number (300ms).
- **Reduced motion**: Instant color swap, instant value update

### 3. Swap Button
- **Trigger**: Tap swap button
- **Animation**: Button icon rotates 180° (400ms spring, slight overshoot to 190° then settles). FROM and TO values cross-fade through midpoint (200ms out, 200ms in). Unit chips swap with horizontal slide (300ms).
- **Reduced motion**: Instant swap, no rotation

### 4. Value Input (Digit Entry)
- **Trigger**: Keypad tap or keyboard press
- **Animation**: New digit scales from 0.8 → 1.0 (50ms spring). Result digits roll/morph to new value (100ms delay to feel "computed", then 200ms roll animation).
- **Reduced motion**: Instant digit appearance, instant result update

### 5. Copy Result
- **Trigger**: Tap on result value
- **Animation**: Brief flash (opacity 1 → 0.6 → 1 in 150ms). Small checkmark fades in next to value (200ms), then fades out (after 1.5s, 300ms fade).
- **Clipboard**: `navigator.clipboard.writeText(resultValue)`
- **Reduced motion**: Same (flash is subtle enough to keep)

### 6. Focus Glow
- **Trigger**: Input field receives focus
- **Animation**: Card border transitions to accent color at 30% opacity (200ms ease). Subtle box-shadow glow appears (200ms).
- **Reduced motion**: Border color change only, no glow animation

---

## Visual Theme Summary

Per combined research from UI and UX research:

### Dark Theme (Default)
- **Background**: Warm obsidian `#1A1714` (slight brown undertone, not cold black)
- **Card surface**: `#242019` with 1px border `rgba(255,255,255,0.06)`
- **Elevated card (TO)**: `#2A2520` — slightly lighter to indicate prominence
- **Swap button**: `#2E2924` background, accent icon
- **Keypad surface**: `#1E1A15`
- **Keypad key**: `#2A2520`, active state `#353025`
- **Text primary**: `#F5F5F0` (warm white)
- **Text secondary**: `#888880`
- **Text muted**: `#555550`

### Light Theme (System Preference)
- **Background**: Warm stone `#F5F2ED`
- **Card surface**: `#FFFFFF` with 1px border `rgba(0,0,0,0.06)`
- **Text primary**: `#1A1714`
- **Text secondary**: `#666660`
- **Accent colors**: Same hues, adjusted for contrast on light backgrounds

### Category Accent Colors
- **Length**: Warm copper `#C17F59`
- **Weight**: Sage green `#6B8F71`
- **Temperature**: Slate blue `#6B8FA3`

These are used sparingly: active tab indicator, active chip background, focus glow, swap button icon. The rest of the UI remains monochromatic.

---

## Edge Cases & Special Handling

| Scenario | Behavior |
|----------|----------|
| User enters `0` | Show `0` in result (valid conversion), not empty state |
| User enters very large number (>1e15) | Switch to scientific notation in result (e.g., `1.23e+18`) |
| User enters many decimal places | Allow up to 10 digits of input, result shows up to 6 significant figures |
| Temperature below absolute zero | Show result with inline warning: "⚠ Below absolute zero" |
| User pastes non-numeric text | Strip to numeric characters only (digits, one decimal, leading minus) |
| User rapidly switches categories | Debounce animation (cancel in-flight animation, start new one) |
| Network offline | No impact — app is fully client-side, works offline after first load |
| Screen reader user converts | Announce: "[result value] [unit name]" via `aria-live` region |

---

Status: READY_FOR_REVIEW
