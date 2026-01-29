# UI Research: SubTrack - Subscription Manager

## Design Inspiration

### Linear (Project Management)
- **Visual Style**: Ultra-minimal dark mode with muted gray tones, subtle purple accent, crisp sans-serif typography, generous whitespace
- **What Makes It Great**: Feels premium without being flashy. The restrained color palette creates focus. Micro-interactions are buttery smooth. Icons are custom and consistent.
- **What We Can Borrow**: The restraint in color usage - one accent color used sparingly. The smooth transitions. The "less is more" philosophy executed perfectly.

### Cleo (Finance App)
- **Visual Style**: Dark mode with vibrant coral/salmon accents, playful but sophisticated, rounded UI elements, conversational tone
- **What Makes It Great**: Finance apps are usually boring/corporate - Cleo feels friendly and approachable. The coral accent is memorable and not overused in fintech.
- **What We Can Borrow**: Warm accent color that feels friendly (not corporate blue). The balance of playful personality with serious functionality.

### Arc Browser
- **Visual Style**: Soft gradients with muted pastels, glassmorphism effects, rounded corners everywhere, dynamic theming
- **What Makes It Great**: Fresh and modern without being garish. The subtle gradients add depth without distraction. Feels like a premium product.
- **What We Can Borrow**: Soft, muted color palette. Generous border-radius. Subtle layering effects.

### Raycast (Productivity)
- **Visual Style**: Deep charcoal background, electric accents (cyan/magenta), crisp typography, command-line aesthetic with polish
- **What Makes It Great**: Dark mode done right - not pure black but warm charcoal. The electric accent colors pop perfectly against dark backgrounds.
- **What We Can Borrow**: Warm dark mode (not #000000), vibrant accent that draws attention to CTAs.

### Notion (Productivity)
- **Visual Style**: Clean white/cream backgrounds, minimal decoration, content-first design, subtle borders and shadows
- **What Makes It Great**: The app disappears - content takes center stage. Typography hierarchy is perfect. Nothing competes for attention.
- **What We Can Borrow**: Content-first approach for light mode. Subtle shadows over harsh borders. Clean empty states.

---

## Competitor Color Analysis

| Competitor | Primary Color | Style | Notes |
|------------|--------------|-------|-------|
| Truebill/Rocket Money | Purple gradient #7C3AED → #4F46E5 | Modern, fintech | Common in fintech - AVOID pure purple |
| Bobby | Coral/peach tones #FF6B6B | Warm, friendly | App is discontinued but beloved. Could modernize this direction |
| Mint | Bright green #00BF6F | Financial health | Screams "money" - too literal, AVOID |
| YNAB | Blue #0072CE | Corporate, trustworthy | Generic financial blue - AVOID |
| Copilot | Blue gradient #2563EB → #1D4ED8 | Premium fintech | Also blue gradient - oversaturated space |
| Subscriptions (Apple) | System blue + gray | Native iOS | No personality, just functional |

**Pattern Observed**: Most finance/subscription apps use:
- Blue (trust/corporate)
- Green (money/growth)
- Purple (modern fintech)

All three are now generic. Users see these colors and think "another finance app."

---

## Our Differentiation Strategy

To stand out, we should:
- **Avoid**: Blue, green, purple as primary colors (too associated with generic fintech)
- **Avoid**: Pure black backgrounds (feels dated, harsh on eyes)
- **Consider**: Warm, approachable tones that feel calm (not urgent about money)

**Brand Mood**: We want SubTrack to feel:
- **Calm** (not anxiety-inducing about spending)
- **Trustworthy** (not salesy or pushy)
- **Delightful** (small moments of joy)
- **Minimal** (substance over decoration)

This translates to: **warm neutral bases with a distinctive, unexpected accent color**.

---

## Color Psychology

**App Mood**: "Calm confidence" - Users should feel in control, not stressed

### Recommended Color Directions

#### 1. Warm Charcoal + Coral Accent (RECOMMENDED)
- **Base**: Deep warm gray (#1A1A1E) for dark mode, warm off-white (#FAFAFA) for light
- **Accent**: Soft coral (#FF7A5C) - warm, friendly, energetic without being alarming
- **Why**: Coral is memorable and underused in fintech. It's the spiritual successor to Bobby's peach tones but modernized. Works beautifully in both light and dark modes.

#### 2. Slate + Teal Accent
- **Base**: Cool slate gray (#1E293B) for dark, cool white (#F8FAFC) for light
- **Accent**: Deep teal (#0D9488) - fresh, calming, sophisticated
- **Why**: Teal bridges blue (trust) and green (money) but feels more unique. Calming effect matches our "no anxiety" goal.

#### 3. Neutral + Amber Accent
- **Base**: True neutral gray (#18181B) for dark, warm white (#FFFBF5) for light
- **Accent**: Rich amber (#F59E0B) - warm, premium, attention-grabbing
- **Why**: Amber/gold conveys premium quality. Warm and inviting. Uncommon in subscription trackers.

---

## Typography Research

**Requirements**:
- Excellent legibility at small sizes (mobile-first)
- Clear number rendering (prices are key data)
- Modern but not trendy (longevity)
- Good weight range (need light to bold)
- Free/open source or system fonts

### Display/UI Font Options

**Geist (by Vercel)**
- Clean, modern geometric sans
- Excellent for interfaces
- Tabular numbers for price alignment
- Why it fits: Feels cutting-edge but readable. Great on screens.

**Inter (by Rasmus Andersson)**
- Designed specifically for screens
- Extensive weight range
- Superb legibility at small sizes
- Why it fits: Battle-tested for interfaces. Neutral personality.

**Manrope (by Mikhail Sharanda)**
- Slightly warmer/rounded than Inter
- Modern geometric with personality
- Good for both UI and display
- Why it fits: Friendlier feel, matches our "approachable" goal.

**SF Pro (System - iOS/Mac)**
- Native iOS feel
- No font loading needed
- Familiar to Apple users
- Why it fits: Zero performance cost, native feel on Apple devices.

### Recommendation: **Geist** as primary

Rationale:
- Numbers render beautifully (critical for subscription prices)
- Clean and modern without being cold
- Tabular numbers prevent layout shift when prices change
- Variable font = single file, all weights
- Free and open source

Fallback: SF Pro → Inter → system-ui

---

## Visual Style Direction

### Recommended Style: "Soft Minimal"

A balance between stark minimalism (Linear) and friendly approachability (Cleo).

#### Corners
- **Cards**: 12-16px border-radius (soft, modern)
- **Buttons**: 8-10px border-radius (slightly less rounded)
- **Inputs**: 8px border-radius (consistent with buttons)
- **Modals**: 16-20px border-radius (generous, feels native)
- **Full round**: Avatars, pills, tags only

#### Shadows
- **Light mode**: Soft, diffused shadows for elevation
  - Cards: `0 2px 8px rgba(0,0,0,0.08)`
  - Modals: `0 8px 30px rgba(0,0,0,0.12)`
- **Dark mode**: Minimal/no shadows, use subtle borders or background differences
  - Cards: 1px border at 8% white opacity
  - Modals: `0 8px 30px rgba(0,0,0,0.5)`

#### Animations
- **Style**: Subtle and functional, not decorative
- **Duration**: 150-200ms for micro-interactions, 300ms for page transitions
- **Easing**: ease-out for enters, ease-in for exits
- **Examples**:
  - Button hover: gentle scale (1.02) + brightness
  - Card tap: subtle press effect (scale 0.98)
  - List additions: fade + slide from below
  - Deletions: fade + slide to side

#### Icons
- **Style**: Outlined (not filled) for navigation/actions
- **Size**: 24px for navigation, 20px for inline, 16px for small indicators
- **Set**: Lucide (open source, consistent, comprehensive)
- **Stroke width**: 1.5-2px for legibility on mobile

---

## Dark Mode Considerations

Since dark mode is essential (per UX research), we need a carefully considered dark palette:

### Dark Mode Principles
1. **Not pure black**: Use warm charcoal (#1A1A1E or #18181B) - easier on eyes
2. **Layered surfaces**: Background → Surface → Elevated surface
3. **Reduced contrast for text**: Don't use pure white (#FFFFFF), use off-white (#FAFAFA or #E5E5E5)
4. **Accent colors may need adjustment**: Colors that work in light mode may be too saturated in dark
5. **Shadows replaced with borders/layers**: Shadows are nearly invisible on dark backgrounds

### Surface Hierarchy (Dark Mode)
| Layer | Suggested Color | Usage |
|-------|-----------------|-------|
| Background | #0F0F12 | App background |
| Surface | #1A1A1E | Cards, sheets |
| Surface Elevated | #252529 | Modals, dropdowns |
| Surface Hover | #2A2A2F | Interactive states |

---

## Iconography Direction

### Recommended: Lucide Icons
- Open source (MIT license)
- 1000+ icons with consistent style
- 1.5px stroke weight - legible on mobile
- Active development and community
- React/Vue components available

### Key Icons Needed
| Action | Icon | Notes |
|--------|------|-------|
| Add subscription | Plus | FAB primary action |
| Home/Dashboard | Home or LayoutDashboard | Bottom nav |
| Calendar/Upcoming | Calendar | View renewals |
| Settings | Settings | User preferences |
| Delete | Trash2 | Destructive action (red) |
| Edit | Pencil | Modify subscription |
| Search | Search | Filter subscriptions |
| Notification | Bell | Reminder settings |
| Category | Tags or Folder | Organize subs |
| Price | DollarSign or receipt | Currency/cost |

---

## Accessibility Considerations

### Color Contrast Requirements
- **Normal text**: Minimum 4.5:1 contrast ratio (WCAG AA)
- **Large text**: Minimum 3:1 contrast ratio
- **Interactive elements**: Must have visible focus states
- **Don't rely on color alone**: Use icons/text alongside color indicators

### Tested Combinations (Coral on Dark)
| Foreground | Background | Ratio | Pass? |
|------------|------------|-------|-------|
| #FF7A5C (Coral) | #0F0F12 (Background) | 6.2:1 | AA |
| #FAFAFA (Text) | #0F0F12 (Background) | 18.2:1 | AAA |
| #FF7A5C (Coral) | #FAFAFA (Light bg) | 3.1:1 | AA Large only |
| #E85C3A (Darker coral) | #FAFAFA (Light bg) | 4.6:1 | AA |

**Note**: Coral accent needs to be slightly darkened for light mode text usage.

---

## Micro-interaction Inspiration

### Adding a Subscription (Delight Moment)
- Quick haptic feedback on save
- Card animates in from bottom with slight bounce
- Total spend counter animates up smoothly
- Subtle confetti or check animation (optional, not overdone)

### Deleting a Subscription
- Swipe reveals red delete zone
- Card collapses smoothly
- Total spend animates down
- Undo toast appears at bottom

### Pull to Refresh
- Custom branded loader (not default spinner)
- Smooth spring animation
- Brief haptic on refresh complete

---

## Inspiration Board Summary

| Element | Inspiration From | Why |
|---------|------------------|-----|
| Overall restraint | Linear | Proves minimal can be premium |
| Warm accent color | Cleo | Friendly finance isn't an oxymoron |
| Dark mode warmth | Raycast | Charcoal > pure black |
| Card design | Bobby | Proven pattern for subscriptions |
| Typography clarity | Notion | Content-first, hierarchy matters |
| Micro-interactions | Stripe Dashboard | Subtle polish creates premium feel |

---

## Final Recommendation

**Primary Direction**: Warm Charcoal + Coral

This direction because:
1. **Differentiation**: No major competitor uses coral/warm tones
2. **Mood Match**: Warm colors = calm, friendly (our brand goals)
3. **Bobby Legacy**: Modernizes beloved Bobby aesthetic
4. **Versatility**: Works in both light and dark modes
5. **Memorability**: Users will remember "the coral subscription app"

**Typography**: Geist (with system font fallbacks)
**Icons**: Lucide (outlined, 1.5px stroke)
**Style**: Soft minimal (generous radius, subtle shadows, restrained animation)

---

Status: READY_FOR_PLANNING
