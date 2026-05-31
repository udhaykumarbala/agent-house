# UI Flow — Color Palette Generator (Brutalist)

## Design System: Brutalist ("Raw Machine")
- **Font**: Space Mono 400/700 (monospace everywhere)
- **Colors**: Background #0A0A0A, Text #FFFFFF, Accent #EBFF00, Surface #141414
- **Borders**: 2-3px solid white, no border-radius (0px everywhere)
- **Shadows**: None
- **Transitions**: None (instant state changes)
- **Typography**: Uppercase headings/buttons, wide letter-spacing, weight 700

## Phase 1: Core Features

### Page Load
1. Browser loads `index.html` with CSP meta tag, Space Mono font, and styles
2. Scripts load: `colors.js` → `clipboard.js` → `app.js`
3. `app.js` auto-generates initial 5-color palette via `generatePalette('random')`
4. Each swatch renders with background color and auto-contrast hex code (uppercase, 700 weight)

### Generate New Palette
- **Trigger**: Click "GENERATE" button OR press Spacebar
- **Flow**: `generate()` → `generatePalette(mode, existingPalette)` → `render()`
- **Result**: All 5 swatches update with new colors (locked colors preserved)
- **Animation**: None — instant color snap (transition: none)

### Copy Color
- **Trigger**: Click any swatch (excluding lock button)
- **Flow**: `handleSwatchClick()` → reads `state.palette[i]` → `formatColor()` → `copyToClipboard()` → `showToast()`
- **Result**: Color value copied to clipboard in current display format (from data model, never DOM)
- **Feedback**: Toast appears "Copied [value]!" — auto-dismisses after 1.5s

### Switch Color Format
- **Trigger**: Click HEX / RGB / HSL button in toolbar format toggle
- **Flow**: `handleFormatToggle()` → updates `state.colorFormat` → toggles `.active` class → `render()`
- **Result**: All 5 swatch labels update to the selected format
- **Formats**:
  - HEX: `#3A86FF`
  - RGB: `rgb(58, 134, 255)`
  - HSL: `hsl(220, 100%, 61%)`
- **Copy behavior**: Copies in the currently displayed format

## Phase 2: Interactive Features

### Lock/Unlock Colors
- **Trigger**: Click lock button on swatch OR press L key while swatch is focused/hovered
- **Flow**: `handleLockClick()` / keyboard handler → toggles `palette[i].locked` → `render()`
- **Visual**:
  - Unlocked: outline lock SVG, button hidden until hover (desktop) or always visible (mobile)
  - Locked: filled lock SVG in contrast color (black/white based on luminance), always visible, thick 3px top border
- **Behavior**: Locked colors are NOT regenerated when Generate is clicked
- **All locked**: If all 5 are locked, Generate shows toast "Unlock a color to generate new ones" and button appears muted
- **ARIA**: `aria-pressed="true/false"`, `aria-label="Lock this color" / "Unlock this color"`

### Harmony Mode Selector
- **Trigger**: Change the harmony `<select>` dropdown in toolbar
- **Flow**: `handleHarmonyChange()` → updates `state.harmonyMode` → `generate()`
- **Modes**: Random, Analogous, Complementary, Triadic, Split-Complementary, Monochromatic
- **Result**: Palette regenerates immediately using the selected harmony algorithm
- **Style**: Custom-styled native `<select>` matching brutalist theme (#141414 bg, #FFFFFF text, 2px solid border, uppercase monospace)

### Keyboard Shortcuts
| Key | Action | Context |
|-----|--------|---------|
| Space | Generate new palette | Global (not on toolbar controls or inputs) |
| L | Toggle lock on focused/hovered swatch | Swatch must be focused or hovered |
| C | Copy focused/hovered swatch color | Swatch must be focused or hovered |
| 1-5 | Focus swatch by number | Global (not on toolbar controls) |
| Tab | Navigate between swatches | Standard browser behavior (tabindex=0) |

- Shortcuts do NOT fire when user is focused on `<select>`, `<input>`, `<textarea>`, or toolbar buttons
- `getActiveSwatch()` checks `document.activeElement` first, then falls back to `:hover`

### State Model
```javascript
{
  palette: [
    { h, s, l, hex, rgb: {r,g,b}, locked: false },
    // ... 5 color objects
  ],
  harmonyMode: 'random',  // 'random' | 'analogous' | 'complementary' | 'triadic' | 'split-complementary' | 'monochromatic'
  colorFormat: 'hex'       // 'hex' | 'rgb' | 'hsl'
}
```

### Data Flow
```
User action → state update → render() → DOM (textContent only)
                                ↓
                          Lock button: aria-pressed, aria-label, CSS class .locked
                          Color code: formatColor(color, state.colorFormat)
                          Contrast: getContrastColor(hex) for text + icon colors
                                ↓
                          copyToClipboard() reads from state, not DOM
```

### Responsive Layout
| Viewport | Layout |
|----------|--------|
| < 768px (mobile) | 5 horizontal rows, lock buttons always visible, toolbar wraps |
| >= 768px (desktop) | 5 vertical columns, lock buttons appear on hover |

## Phase 3: Polish & Refinement

### Animations & Micro-interactions

#### Animations (Brutalist — All Removed)
- **Palette generation**: No stagger delays, no wave effect. All 5 swatches change simultaneously and instantly. `triggerStagger()` removed from JS. All `@keyframes` and `.stagger` CSS rules removed.
- **Copy feedback**: Toast appears/disappears instantly (opacity toggle, no slide). No copy-pulse scale animation. `addCopyPulse()` removed from JS, `@keyframes copy-pulse` removed from CSS.
- **Lock toggle**: Icon swaps instantly, thick top border appears. No lock-bounce animation. `addLockBounce()` removed from JS, `@keyframes lock-bounce` removed from CSS.
- **Hover**: No brightness overlay — `.swatch::before` pseudo-element and dark/light hover overlays removed from CSS. Lock button becomes visible on hover (desktop).
- **Swatch dividers**: 1px solid #333333 — `border-bottom` on mobile, `border-right` on desktop.
- **Reduced Motion**: Already satisfied — no transitions or animations exist.

### Accessibility Refinements

#### Screen Reader Announcements
- **Element**: `<div id="sr-announce" class="sr-only" aria-live="polite" aria-atomic="true">`
- **Generate**: `announce('New palette generated')` on every successful generation
- **Copy**: `announce('Copied [value]')` on individual copy; `announce('Palette copied to clipboard')` on Copy All
- **Lock/Unlock**: `announce('Color N locked/unlocked')` on every toggle
- **Technique**: Clear text → `requestAnimationFrame` → set new text (ensures AT picks up change)

#### ARIA Labels
| Element | aria-label |
|---------|-----------|
| Generate button | "Generate new color palette" |
| Harmony select | "Harmony mode" |
| Lock button (unlocked) | "Lock color N #XXXXXX" |
| Lock button (locked) | "Unlock color N #XXXXXX" |
| Format button HEX | "HEX color format" |
| Format button RGB | "RGB color format" |
| Format button HSL | "HSL color format" |
| Copy All button | "Copy all palette colors to clipboard" |
| Each swatch | "Color N: #XXXXXX. Locked" or "Color N: #XXXXXX. Unlocked" |

#### Focus Rings
- All focusable elements: `3px solid #EBFF00` outline, `3px` offset via `*:focus-visible`
- Swatches: `3px solid #FFFFFF` outline, `-3px` offset (white for visibility against any swatch color)
- Mouse users: no outline via `*:focus:not(:focus-visible)`

#### Touch Targets
- All interactive elements have `min-width: 44px` and `min-height: 44px` via `--touch-target`
- Lock buttons: 44x44px with 8px padding
- Format buttons: 44x44px min dimensions
- Generate button: 44px min-height via padding

#### Tab Order
1. Swatch 1 → Lock btn 1 → Swatch 2 → Lock btn 2 → ... → Swatch 5 → Lock btn 5
2. Harmony select → Generate button
3. Format toggle (HEX, RGB, HSL) → Copy All button

### Copy All Palette
- **Trigger**: Click "Copy All" button in toolbar
- **Flow**: `handleCopyAll()` → `formatColor()` for each → `copyToClipboard(joined)` → `showToast()`
- **Format**: Comma-separated values in current display format (e.g., "#A45C2F, #D4913B, #E8C97A, #5B9A6B, #2E4A3F")
- **Toast**: "Palette copied!"
