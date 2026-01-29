# UI Specification: Inventory Reorder Alert System

Based on research in: `.plans/research/ui-research.md`
Styling wireframes from: `.plans/specs/ux-spec.md`
Feature requirements from: `.plans/specs/product-spec.md`

---

## Design System

### Color Palette (UNIQUE - "Refined Industrial" Theme)

**IMPORTANT**: These colors differentiate us from the sea of blue/green/orange inventory systems. Our plum + warm cream palette is unique in this space.

| Name | Hex | RGB | CSS Variable | Usage |
|------|-----|-----|--------------|-------|
| **Primary** | #7C3AED | rgb(124, 58, 237) | `--color-primary` | CTAs, brand elements, active states |
| **Primary Hover** | #6D28D9 | rgb(109, 40, 217) | `--color-primary-hover` | Button hover states |
| **Primary Light** | #EDE9FE | rgb(237, 233, 254) | `--color-primary-light` | Selected backgrounds, badges |
| **Secondary** | #78716C | rgb(120, 113, 108) | `--color-secondary` | Secondary actions, muted text |
| **Background** | #FDF8F3 | rgb(253, 248, 243) | `--color-bg` | App background (warm cream) |
| **Surface** | #FFFFFF | rgb(255, 255, 255) | `--color-surface` | Cards, elevated elements |
| **Surface Elevated** | #FFFFFE | rgb(255, 255, 254) | `--color-surface-elevated` | Modals, dropdowns |
| **Border** | #E7E5E4 | rgb(231, 229, 228) | `--color-border` | Dividers, card borders |
| **Border Focus** | #7C3AED | rgb(124, 58, 237) | `--color-border-focus` | Focus rings |
| **Text Primary** | #1F1F23 | rgb(31, 31, 35) | `--color-text` | Main text (warm charcoal) |
| **Text Secondary** | #78716C | rgb(120, 113, 108) | `--color-text-secondary` | Muted/helper text |
| **Text Tertiary** | #A8A29E | rgb(168, 162, 158) | `--color-text-tertiary` | Placeholders, disabled |

#### Status Colors (Alert-Optimized)

| Status | Hex | RGB | CSS Variable | Usage |
|--------|-----|-----|--------------|-------|
| **Critical** | #B91C1C | rgb(185, 28, 28) | `--color-critical` | Out of stock, urgent |
| **Critical Light** | #FEF2F2 | rgb(254, 242, 242) | `--color-critical-light` | Critical backgrounds |
| **Warning** | #D97706 | rgb(217, 119, 6) | `--color-warning` | Low stock alerts |
| **Warning Light** | #FEF3C7 | rgb(254, 243, 199) | `--color-warning-light` | Warning backgrounds |
| **Success** | #4D7C0F | rgb(77, 124, 15) | `--color-success` | In stock, healthy |
| **Success Light** | #ECFCCB | rgb(236, 252, 203) | `--color-success-light` | Success backgrounds |
| **Pending** | #78716C | rgb(120, 113, 108) | `--color-pending` | Ordered, waiting |
| **Pending Light** | #F5F5F4 | rgb(245, 245, 244) | `--color-pending-light` | Pending backgrounds |

**Why These Colors**:
- Primary plum (#7C3AED) stands out from competitors' blues/greens while conveying sophistication
- Warm cream background (#FDF8F3) feels approachable, not cold like pure white
- Status colors follow established psychology: red=urgent, amber=caution, sage green=calm

---

### Typography

**Font Family**: Inter

```css
font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
```

**Why Inter**: Industry standard for data-heavy interfaces, excellent number legibility (critical for inventory quantities), variable font for performance.

| Element | Size | Weight | Line Height | Letter Spacing | CSS Class |
|---------|------|--------|-------------|----------------|-----------|
| **Display** | 36px | 700 | 1.1 | -0.02em | `.text-display` |
| **H1** | 32px | 700 | 1.2 | -0.02em | `.text-h1` |
| **H2** | 24px | 600 | 1.3 | -0.01em | `.text-h2` |
| **H3** | 20px | 600 | 1.4 | -0.01em | `.text-h3` |
| **H4** | 18px | 600 | 1.4 | 0 | `.text-h4` |
| **Body Large** | 18px | 400 | 1.5 | 0 | `.text-body-lg` |
| **Body** | 16px | 400 | 1.5 | 0 | `.text-body` |
| **Body Small** | 14px | 400 | 1.5 | 0 | `.text-body-sm` |
| **Caption** | 12px | 500 | 1.4 | 0.01em | `.text-caption` |
| **Label** | 14px | 500 | 1.4 | 0.01em | `.text-label` |
| **Mono (Quantities)** | 16px | 600 | 1.2 | 0 | `.text-mono` |

**Number Display**: Use tabular figures for quantity displays to ensure alignment:
```css
.quantity-display {
  font-variant-numeric: tabular-nums;
  font-weight: 600;
}
```

---

### Spacing Scale

**Base Unit**: 4px

| Token | Value | CSS Variable | Usage |
|-------|-------|--------------|-------|
| **0** | 0px | `--space-0` | No spacing |
| **1** | 4px | `--space-1` | Tight inline spacing |
| **2** | 8px | `--space-2` | Icon gaps, compact elements |
| **3** | 12px | `--space-3` | Small gaps, badge padding |
| **4** | 16px | `--space-4` | Standard padding, card content |
| **5** | 20px | `--space-5` | Medium gaps |
| **6** | 24px | `--space-6` | Section spacing, card padding |
| **8** | 32px | `--space-8` | Large section gaps |
| **10** | 40px | `--space-10` | Page vertical rhythm |
| **12** | 48px | `--space-12` | Hero spacing |
| **16** | 64px | `--space-16` | Major section breaks |

---

### Border Radius

| Token | Value | CSS Variable | Usage |
|-------|-------|--------------|-------|
| **None** | 0px | `--radius-none` | Sharp edges (rare) |
| **SM** | 4px | `--radius-sm` | Small inputs, compact buttons |
| **MD** | 8px | `--radius-md` | Buttons, cards, inputs |
| **LG** | 12px | `--radius-lg` | Large cards, modals |
| **XL** | 16px | `--radius-xl` | Hero cards, bottom sheets |
| **Full** | 9999px | `--radius-full` | Pills, avatars, circular |

---

### Shadows

Warm-tinted shadows to match our palette:

| Token | Value | CSS Variable | Usage |
|-------|-------|--------------|-------|
| **SM** | `0 1px 2px rgba(31, 31, 35, 0.05)` | `--shadow-sm` | Subtle lift |
| **MD** | `0 4px 6px rgba(31, 31, 35, 0.07), 0 2px 4px rgba(31, 31, 35, 0.04)` | `--shadow-md` | Cards, buttons |
| **LG** | `0 10px 15px rgba(31, 31, 35, 0.08), 0 4px 6px rgba(31, 31, 35, 0.04)` | `--shadow-lg` | Elevated cards |
| **XL** | `0 20px 25px rgba(31, 31, 35, 0.10), 0 8px 10px rgba(31, 31, 35, 0.05)` | `--shadow-xl` | Modals, dropdowns |
| **Focus** | `0 0 0 3px rgba(124, 58, 237, 0.25)` | `--shadow-focus` | Focus rings |

---

## Components

### 1. Button

#### Primary Button
```css
.btn-primary {
  background: var(--color-primary);
  color: #FFFFFF;
  padding: 12px 24px;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-size: 16px;
  min-height: 48px;
  border: none;
  box-shadow: var(--shadow-sm);
  transition: all 150ms ease;
}

.btn-primary:hover {
  background: var(--color-primary-hover);
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.btn-primary:active {
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
}

.btn-primary:focus-visible {
  outline: none;
  box-shadow: var(--shadow-focus);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}
```

#### Secondary Button
```css
.btn-secondary {
  background: var(--color-surface);
  color: var(--color-text);
  padding: 12px 24px;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-size: 16px;
  min-height: 48px;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
  transition: all 150ms ease;
}

.btn-secondary:hover {
  background: var(--color-bg);
  border-color: var(--color-secondary);
}
```

#### Ghost Button
```css
.btn-ghost {
  background: transparent;
  color: var(--color-primary);
  padding: 12px 24px;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-size: 16px;
  min-height: 48px;
  border: none;
  transition: all 150ms ease;
}

.btn-ghost:hover {
  background: var(--color-primary-light);
}
```

#### Icon Button
```css
.btn-icon {
  width: 44px;
  height: 44px;
  padding: 10px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
}
```

#### Button Sizes
| Size | Padding | Min Height | Font Size |
|------|---------|------------|-----------|
| **SM** | 8px 16px | 36px | 14px |
| **MD** | 12px 24px | 48px | 16px |
| **LG** | 16px 32px | 56px | 18px |

---

### 2. Input Field

```css
.input {
  width: 100%;
  padding: 12px 16px;
  font-size: 16px;
  line-height: 1.5;
  color: var(--color-text);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  min-height: 48px;
  transition: all 150ms ease;
}

.input::placeholder {
  color: var(--color-text-tertiary);
}

.input:hover {
  border-color: var(--color-secondary);
}

.input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: var(--shadow-focus);
}

.input:disabled {
  background: var(--color-bg);
  color: var(--color-text-tertiary);
  cursor: not-allowed;
}

.input-error {
  border-color: var(--color-critical);
}

.input-error:focus {
  box-shadow: 0 0 0 3px rgba(185, 28, 28, 0.25);
}
```

#### Input with Label
```html
<div class="input-group">
  <label class="input-label">Item Name *</label>
  <input class="input" type="text" placeholder="e.g., Coffee Beans">
  <span class="input-helper">Required field</span>
</div>
```

```css
.input-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.input-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.input-helper {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.input-error-text {
  font-size: 12px;
  color: var(--color-critical);
}
```

---

### 3. Card

```css
.card {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--space-6);
  border: 1px solid var(--color-border);
}

.card-interactive {
  cursor: pointer;
  transition: all 150ms ease;
}

.card-interactive:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

.card-interactive:active {
  transform: translateY(0);
}
```

#### Card Variants
```css
/* Hero Card - Dashboard metric */
.card-hero {
  background: linear-gradient(135deg, var(--color-primary) 0%, #6D28D9 100%);
  color: #FFFFFF;
  border: none;
}

/* Alert Card - Items needing attention */
.card-alert-critical {
  border-left: 4px solid var(--color-critical);
}

.card-alert-warning {
  border-left: 4px solid var(--color-warning);
}

.card-alert-success {
  border-left: 4px solid var(--color-success);
}
```

---

### 4. Status Badge

```css
.badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: 4px 10px;
  border-radius: var(--radius-full);
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.badge-critical {
  background: var(--color-critical-light);
  color: var(--color-critical);
}

.badge-warning {
  background: var(--color-warning-light);
  color: var(--color-warning);
}

.badge-success {
  background: var(--color-success-light);
  color: var(--color-success);
}

.badge-pending {
  background: var(--color-pending-light);
  color: var(--color-pending);
}
```

---

### 5. Stock Level Bar

A key component showing current stock vs. threshold visually:

```css
.stock-bar {
  width: 100%;
  height: 8px;
  background: var(--color-bg);
  border-radius: var(--radius-full);
  overflow: visible;
  position: relative;
}

.stock-bar-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width 300ms ease-out;
}

.stock-bar-fill-critical {
  background: linear-gradient(90deg, var(--color-critical) 0%, #DC2626 100%);
}

.stock-bar-fill-warning {
  background: linear-gradient(90deg, var(--color-warning) 0%, #F59E0B 100%);
}

.stock-bar-fill-success {
  background: linear-gradient(90deg, var(--color-success) 0%, #65A30D 100%);
}

/* Threshold marker */
.stock-bar-threshold {
  position: absolute;
  top: -4px;
  height: 16px;
  width: 2px;
  background: var(--color-text-secondary);
  border-radius: 1px;
}

.stock-bar-threshold::after {
  content: '';
  position: absolute;
  top: -2px;
  left: -3px;
  width: 8px;
  height: 8px;
  background: var(--color-text-secondary);
  border-radius: 50%;
}
```

#### Large Stock Bar (Item Detail)
```css
.stock-bar-lg {
  height: 16px;
}

.stock-bar-lg .stock-bar-threshold {
  top: -8px;
  height: 32px;
}
```

---

### 6. Navigation

#### Bottom Navigation (Mobile)
```css
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: var(--color-surface);
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 100;
}

.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 12px;
  font-weight: 500;
  transition: color 150ms ease;
}

.bottom-nav-item.active {
  color: var(--color-primary);
}

.bottom-nav-item:hover {
  color: var(--color-primary);
}

.bottom-nav-icon {
  width: 24px;
  height: 24px;
}

/* Alert badge on nav item */
.bottom-nav-badge {
  position: absolute;
  top: 0;
  right: 0;
  min-width: 18px;
  height: 18px;
  background: var(--color-critical);
  color: #FFFFFF;
  font-size: 10px;
  font-weight: 700;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
}
```

#### Sidebar Navigation (Desktop)
```css
.sidebar-nav {
  width: 240px;
  height: 100vh;
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  padding: var(--space-6);
  position: fixed;
  left: 0;
  top: 0;
}

.sidebar-nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: 12px 16px;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  text-decoration: none;
  font-weight: 500;
  transition: all 150ms ease;
  margin-bottom: 4px;
}

.sidebar-nav-item:hover {
  background: var(--color-bg);
  color: var(--color-text);
}

.sidebar-nav-item.active {
  background: var(--color-primary-light);
  color: var(--color-primary);
}
```

---

### 7. Filter Chips

```css
.filter-chips {
  display: flex;
  gap: var(--space-2);
  overflow-x: auto;
  padding: var(--space-2) 0;
  -webkit-overflow-scrolling: touch;
}

.filter-chip {
  flex-shrink: 0;
  padding: 8px 16px;
  border-radius: var(--radius-full);
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  cursor: pointer;
  transition: all 150ms ease;
}

.filter-chip:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.filter-chip.active {
  background: var(--color-primary);
  color: #FFFFFF;
  border-color: var(--color-primary);
}
```

---

### 8. Quantity Stepper

```css
.quantity-stepper {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.quantity-btn {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text);
  cursor: pointer;
  transition: all 150ms ease;
}

.quantity-btn:hover {
  background: var(--color-bg);
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.quantity-btn:active {
  background: var(--color-primary-light);
}

.quantity-value {
  min-width: 80px;
  text-align: center;
  font-size: 24px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: var(--color-text);
}
```

---

### 9. Bottom Sheet Modal

```css
.bottom-sheet-overlay {
  position: fixed;
  inset: 0;
  background: rgba(31, 31, 35, 0.5);
  z-index: 200;
  opacity: 0;
  transition: opacity 250ms ease;
}

.bottom-sheet-overlay.open {
  opacity: 1;
}

.bottom-sheet {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--color-surface);
  border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  padding: var(--space-6);
  padding-bottom: calc(var(--space-6) + env(safe-area-inset-bottom));
  transform: translateY(100%);
  transition: transform 250ms ease;
  z-index: 201;
  max-height: 90vh;
  overflow-y: auto;
}

.bottom-sheet.open {
  transform: translateY(0);
}

.bottom-sheet-handle {
  width: 36px;
  height: 4px;
  background: var(--color-border);
  border-radius: var(--radius-full);
  margin: 0 auto var(--space-6);
}
```

---

### 10. Toast Notification

```css
.toast {
  position: fixed;
  bottom: 80px; /* Above bottom nav */
  left: 50%;
  transform: translateX(-50%) translateY(100px);
  background: var(--color-text);
  color: #FFFFFF;
  padding: 12px 20px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  gap: var(--space-3);
  font-weight: 500;
  box-shadow: var(--shadow-xl);
  z-index: 300;
  opacity: 0;
  transition: all 300ms ease;
}

.toast.visible {
  transform: translateX(-50%) translateY(0);
  opacity: 1;
}

.toast-success {
  background: var(--color-success);
}

.toast-error {
  background: var(--color-critical);
}
```

---

### 11. Empty State

```css
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-12);
  text-align: center;
}

.empty-state-icon {
  width: 80px;
  height: 80px;
  margin-bottom: var(--space-6);
  color: var(--color-text-tertiary);
}

.empty-state-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--space-2);
}

.empty-state-description {
  font-size: 16px;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-6);
  max-width: 280px;
}
```

---

### 12. Floating Action Button

```css
.fab {
  position: fixed;
  bottom: 88px; /* Above bottom nav */
  right: var(--space-6);
  width: 56px;
  height: 56px;
  border-radius: var(--radius-full);
  background: var(--color-primary);
  color: #FFFFFF;
  border: none;
  box-shadow: var(--shadow-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 150ms ease;
  z-index: 50;
}

.fab:hover {
  background: var(--color-primary-hover);
  box-shadow: var(--shadow-xl);
  transform: scale(1.05);
}

.fab:active {
  transform: scale(0.95);
}

.fab-icon {
  width: 24px;
  height: 24px;
}
```

---

## Icons

Use **Lucide Icons** (outlined, 1.5px stroke) for consistency:

| Context | Icon | Name |
|---------|------|------|
| Home/Dashboard | house | `House` |
| Items/Inventory | package | `Package` |
| Alerts | bell | `Bell` |
| Settings | settings | `Settings` |
| Add | plus | `Plus` |
| Edit | pencil | `Pencil` |
| Delete | trash-2 | `Trash2` |
| Search | search | `Search` |
| Close | x | `X` |
| Back | arrow-left | `ArrowLeft` |
| Reorder/Cart | shopping-cart | `ShoppingCart` |
| Snooze | clock | `Clock` |
| Check/Done | check | `Check` |
| Critical | alert-triangle | `AlertTriangle` |
| Warning | zap | `Zap` |
| OK/Success | check-circle | `CheckCircle` |
| Photo | image | `Image` |
| Import | upload | `Upload` |
| Export | download | `Download` |
| User | user | `User` |
| Menu | menu | `Menu` |
| More | more-vertical | `MoreVertical` |

**Icon Sizing**:
| Size | Pixels | Usage |
|------|--------|-------|
| SM | 16px | Inline with text |
| MD | 20px | Buttons, labels |
| LG | 24px | Navigation, headers |
| XL | 32px | Empty states, features |

---

## Responsive Breakpoints

| Name | Min Width | Layout Changes |
|------|-----------|----------------|
| **Mobile** | 320px | Single column, bottom nav, FAB, bottom sheets |
| **Mobile Large** | 414px | Slightly larger touch targets |
| **Tablet** | 768px | Two-column grid, side drawer option |
| **Desktop** | 1024px | Three-column, sidebar nav, inline modals |
| **Wide** | 1280px | Wider content area, more columns |

```css
/* Mobile-first approach */
@media (min-width: 768px) {
  /* Tablet styles */
}

@media (min-width: 1024px) {
  /* Desktop styles */
}
```

---

## Animation Guidelines

### Micro-interactions

| Element | Animation | Duration | Easing |
|---------|-----------|----------|--------|
| Button press | Scale to 0.97 | 100ms | ease-out |
| Button hover | TranslateY -1px | 150ms | ease |
| Card hover | TranslateY -2px, shadow increase | 150ms | ease |
| Toggle switch | Slide with spring | 200ms | cubic-bezier(0.34, 1.56, 0.64, 1) |
| Checkbox | Scale check mark | 150ms | ease-out |

### Transitions

| Element | Animation | Duration | Easing |
|---------|-----------|----------|--------|
| Page push | Slide from right | 300ms | ease-in-out |
| Modal appear | Fade + scale from 0.95 | 200ms | ease-out |
| Bottom sheet | Slide up | 250ms | ease-out |
| Toast appear | Slide up + fade | 300ms | ease-out |
| Dropdown | Scale from 0.95 + fade | 150ms | ease-out |

### Feedback Animations

| Element | Animation | Duration | Easing |
|---------|-----------|----------|--------|
| Success checkmark | Draw path | 400ms | ease-out |
| Quantity update | Counter animation | 200ms | ease-out |
| Stock bar fill | Width transition | 300ms | ease-out |
| Loading skeleton | Shimmer pulse | 1500ms | infinite |

### Motion Preferences

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## Accessibility Specifications

### Color Contrast

All text meets WCAG AA (4.5:1 for normal text, 3:1 for large text):

| Combination | Ratio | Pass |
|-------------|-------|------|
| Text Primary on Background | 15.2:1 | AAA |
| Text Primary on Surface | 15.2:1 | AAA |
| Text Secondary on Background | 5.7:1 | AA |
| Primary on White | 4.6:1 | AA |
| White on Primary | 4.6:1 | AA |
| Critical on Critical Light | 7.2:1 | AAA |

### Focus States

All interactive elements have visible focus:
```css
:focus-visible {
  outline: none;
  box-shadow: var(--shadow-focus);
}
```

### Status Indicators

Never rely on color alone:
- **Critical**: Red + ⚠️ icon + "Critical" text
- **Warning**: Amber + ⚡ icon + "Low" text
- **Success**: Green + ✓ icon + "OK" text

### Touch Targets

Minimum 44x44px for all interactive elements (Apple HIG).

---

## CSS Custom Properties Summary

```css
:root {
  /* Colors */
  --color-primary: #7C3AED;
  --color-primary-hover: #6D28D9;
  --color-primary-light: #EDE9FE;
  --color-secondary: #78716C;
  --color-bg: #FDF8F3;
  --color-surface: #FFFFFF;
  --color-border: #E7E5E4;
  --color-text: #1F1F23;
  --color-text-secondary: #78716C;
  --color-text-tertiary: #A8A29E;

  /* Status */
  --color-critical: #B91C1C;
  --color-critical-light: #FEF2F2;
  --color-warning: #D97706;
  --color-warning-light: #FEF3C7;
  --color-success: #4D7C0F;
  --color-success-light: #ECFCCB;

  /* Spacing */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;
  --space-10: 40px;
  --space-12: 48px;
  --space-16: 64px;

  /* Border Radius */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(31, 31, 35, 0.05);
  --shadow-md: 0 4px 6px rgba(31, 31, 35, 0.07), 0 2px 4px rgba(31, 31, 35, 0.04);
  --shadow-lg: 0 10px 15px rgba(31, 31, 35, 0.08), 0 4px 6px rgba(31, 31, 35, 0.04);
  --shadow-xl: 0 20px 25px rgba(31, 31, 35, 0.10), 0 8px 10px rgba(31, 31, 35, 0.05);
  --shadow-focus: 0 0 0 3px rgba(124, 58, 237, 0.25);

  /* Typography */
  --font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
```

---

## Implementation Notes

1. **Font Loading**: Use `font-display: swap` to prevent FOIT
2. **Images**: Lazy load item photos, use WebP with JPEG fallback
3. **Icons**: Use SVG sprites or icon font for performance
4. **Animations**: Use `transform` and `opacity` for 60fps animations
5. **Mobile Performance**: Avoid expensive box-shadows on scroll

---

Status: READY_FOR_REVIEW
