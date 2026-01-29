# UX Research: Subscription Management Web App

## Pattern Analysis

Referenced apps and their successful patterns:

### Truebill (Rocket Money)
- **Pattern Used**: Card-based subscription list with monthly cost prominently displayed, color-coded categories, one-tap cancel requests
- **Why It Works**: Users can scan their spending at a glance. The visual hierarchy prioritizes cost (the most important data point). Categorization helps users understand where money goes.
- **What We Can Learn**: Lead with the dollar amount. Use cards for each subscription. Group by category or billing cycle.

### Bobby (Subscription Tracker)
- **Pattern Used**: Minimal dashboard showing total monthly/yearly spend at top, horizontal scrolling category pills, individual subscription cards with logo + name + price + billing date
- **Why It Works**: The total spend creates urgency and awareness. Horizontal pills allow quick filtering without taking vertical space. Clean cards with essential info only.
- **What We Can Learn**: Show aggregate totals prominently. Use iconography for quick recognition. Keep individual cards minimal (logo, name, price, date).

### Subscript
- **Pattern Used**: Calendar view for upcoming renewals, push notification reminders before charges, spending trends over time
- **Why It Works**: Calendar mental model matches how users think about bills. Reminders prevent unwanted charges. Trends show progress in reducing spend.
- **What We Can Learn**: Include a calendar/timeline view. Notification reminders are essential. Show spending history/trends.

### Mint (Bill Tracking)
- **Pattern Used**: Bottom navigation with 4-5 tabs, transaction list with search/filter, bill due date alerts, progress bars for budgets
- **Why It Works**: Bottom nav is thumb-friendly on mobile. Search handles users with many subscriptions. Visual progress indicators motivate action.
- **What We Can Learn**: Use bottom navigation for mobile. Include robust search/filter. Visual indicators for budget goals.

### YNAB (You Need A Budget)
- **Pattern Used**: Age of money concept, category-based budgeting, swipe actions on transactions, empty states with helpful prompts
- **Why It Works**: Unique value proposition creates stickiness. Categories match mental model. Swipe gestures speed up common actions.
- **What We Can Learn**: Good empty states guide new users. Swipe gestures for quick actions (archive, edit, delete).

## User Journey Analysis

### Entry Point
- User downloads app after receiving an unexpected subscription charge
- Or user wants to audit their recurring expenses
- First-time: Onboarding should be fast (< 30 seconds to value)

### Core Loop
1. Open app → See dashboard with total spend
2. Review upcoming renewals
3. Get notified before charges
4. Cancel/pause unused subscriptions
5. Track savings over time

### Success State
- User knows exactly what they're paying for
- No surprise charges
- Reduced subscription spend
- Feeling of financial control

## Mental Model

Users think of this task as: **"My subscription budget" or "Bills I need to track"**

They mentally categorize subscriptions as:
- Essential (Netflix, Spotify) vs. Nice-to-have (random apps)
- Shared (family plans) vs. Personal
- Active use vs. Forgotten

Common terminology:
- "Subscription" / "Recurring charge"
- "Cancel" / "Pause"
- "Billing date" / "Renewal date"
- "Monthly" / "Yearly" / "Weekly"
- "Free trial ending"

## Anti-Patterns to Avoid

- **Don't**: Require manual entry for every subscription - **Why**: High friction leads to abandonment. Use bank linking or smart detection where possible.
- **Don't**: Hide the total monthly cost - **Why**: This is the #1 insight users want. It should be above the fold, always visible.
- **Don't**: Use complex navigation or deep hierarchies - **Why**: Users want quick glances, not deep dives. Keep it flat.
- **Don't**: Show too much data per subscription card - **Why**: Cognitive overload. Show name, price, date. Details on tap.
- **Don't**: Make adding a subscription tedious - **Why**: If it takes more than 3 taps, users won't do it consistently.
- **Don't**: Neglect empty states - **Why**: First-time users with no subscriptions added will feel lost.
- **Don't**: Use unfamiliar icons without labels - **Why**: Mobile-first doesn't mean icon-only. Text labels improve clarity.
- **Don't**: Forget about dark mode - **Why**: Finance apps are often checked at night. Eye strain is real.

## Competitive Landscape Summary

| App | Strengths | Weaknesses |
|-----|-----------|------------|
| Truebill/Rocket Money | Bank linking, cancel service | Subscription required for full features |
| Bobby | Beautiful UI, one-time purchase | Manual entry only |
| Subscript | Calendar view, reminders | Cluttered interface |
| Mint | Comprehensive, free | Too broad, not subscription-focused |
| TrackMySubs | Web-based, simple | Dated UI |

## Key Insights for Our App

### Must-Have Features (Based on Research)
1. **Dashboard with total spend** - Above the fold, monthly/yearly toggle
2. **Subscription cards** - Logo, name, price, next billing date
3. **Quick add** - Pre-populated templates for common services
4. **Reminders** - Notify before renewal (configurable days)
5. **Categories** - Entertainment, Productivity, Lifestyle, etc.

### Differentiators to Consider
1. **Minimal by default** - No feature bloat, just essentials
2. **Offline-first** - Works without account/bank linking
3. **Privacy-focused** - No bank linking required, local data option
4. **Smart suggestions** - "You haven't used X in 30 days"

## Recommended Patterns for Our App

Based on research, we should use:

1. **Bottom navigation bar (3-4 tabs max)** - because thumb-friendly, mobile standard, prevents hamburger menu hunting
2. **Card-based subscription list** - because scannable, tappable, matches mental model of "items I own"
3. **Large total spend header** - because primary insight, creates urgency, proven pattern in all competitors
4. **Floating action button (FAB) for add** - because adding is the primary action, always accessible, mobile convention
5. **Swipe gestures for quick actions** - because speeds up common tasks (delete, edit), reduces tap count
6. **Category pills with horizontal scroll** - because filtering without leaving view, space-efficient on mobile
7. **Empty state with illustration + CTA** - because guides new users, reduces confusion, increases activation

## Color & Visual Patterns

### Color Associations in Finance Apps
- **Green**: Savings, positive trends, money saved
- **Red**: Alerts, overspending, cancellation
- **Blue**: Trust, stability (used by banks)
- **Purple/Gradient**: Modern, premium feel (Rocket Money uses this)

### Recommended Palette Direction
- Neutral base (white/dark gray for dark mode)
- Accent color for CTAs and highlights
- Semantic colors for states (green=saved, red=alert)

## Typography Patterns

- Large, bold numbers for prices/totals
- Medium weight for subscription names
- Light/small for secondary info (dates, categories)
- Monospace optional for prices (alignment)

## Mobile-First Considerations

1. **Touch targets**: 44x44px minimum
2. **One-handed use**: Important actions in bottom half
3. **Thumb zone**: Primary actions in easy reach
4. **Viewport**: Design for 375px width (iPhone SE/mini baseline)
5. **Scroll direction**: Vertical primary, horizontal for tabs/filters only
6. **Input minimization**: Use pickers, toggles, pre-filled options

---
Status: READY_FOR_PLANNING
