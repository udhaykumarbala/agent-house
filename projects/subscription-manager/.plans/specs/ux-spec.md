# UX Specification: SubTrack - Subscription Manager

Based on research in: `.plans/research/ux-research.md`, `.plans/research/ui-research.md`
Implementing features from: `.plans/specs/product-spec.md`

---

## Screen List

1. **Dashboard (Home)** - Primary view with total spend and subscription list
2. **Add Subscription** - Form to add new subscription (sheet/modal)
3. **Service Templates** - Browse and search popular services
4. **Subscription Detail** - View/edit individual subscription
5. **Upcoming Renewals** - Calendar/timeline view of upcoming charges
6. **Settings** - App preferences (theme, reminders, data)
7. **Empty State** - First-time user experience

---

## Wireframes

### Screen 1: Dashboard (Home)

Primary screen users see upon opening the app. Optimized for one-handed mobile use.

```
┌─────────────────────────────────────┐
│  ≡  SubTrack              [☀/🌙]   │  ← Minimal header (hamburger optional)
├─────────────────────────────────────┤
│                                     │
│         $247.85 /month              │  ← Large total spend (primary insight)
│         $2,974.20 /year             │  ← Secondary: yearly total
│                                     │
│   [ Monthly ▼ ]   [ All ▼ ]         │  ← Filter pills (cycle, category)
│                                     │
├─────────────────────────────────────┤
│  ┌─────────────────────────────┐    │
│  │ 🎵  Spotify                 │    │
│  │     $9.99/mo    Due Jan 22  │    │  ← Subscription card
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 📺  Netflix                 │    │
│  │     $15.49/mo   Due Jan 25  │    │  ← Card: icon, name, price, date
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ ☁️  iCloud+                 │    │
│  │     $2.99/mo    Due Feb 1   │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │ 💼  Adobe CC         ⚠️     │    │  ← Warning icon: due soon
│  │     $54.99/mo   Due in 2d   │    │
│  └─────────────────────────────┘    │
│                                     │
│           [Load more...]            │
│                                     │
├─────────────────────────────────────┤
│                                     │
│              ( + )                  │  ← FAB: Add subscription
│                                     │
├─────────────────────────────────────┤
│  🏠        📅        ⚙️            │  ← Bottom nav: Home, Upcoming, Settings
│  Home    Upcoming   Settings        │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Total Spend Header**: Large, bold typography. Most important number above the fold.
- **Filter Pills**: Horizontal scroll for filtering by cycle (monthly/yearly) and category.
- **Subscription Cards**: Tappable to view details. Show logo, name, price, next billing date.
- **Due Soon Indicator**: Visual warning for subscriptions due within 7 days.
- **FAB**: Floating action button for adding new subscription (primary action).
- **Bottom Navigation**: 3 tabs - Home, Upcoming, Settings. Thumb-friendly zone.

---

### Screen 2: Add Subscription (Bottom Sheet)

Slides up from bottom. Can be dismissed by swiping down.

```
┌─────────────────────────────────────┐
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │  ← Drag handle
├─────────────────────────────────────┤
│                                     │
│  Add Subscription         [ × ]     │  ← Title + close button
│                                     │
├─────────────────────────────────────┤
│                                     │
│  [ 🔍 Search popular services... ]  │  ← Quick search for templates
│                                     │
│  Popular:                           │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐       │
│  │ 🎵 │ │ 📺 │ │ 🎮 │ │ ☁️  │ →    │  ← Horizontal scroll: templates
│  │Spot│ │Netf│ │Xbox│ │Adob│       │
│  └────┘ └────┘ └────┘ └────┘       │
│                                     │
│  ─── or enter manually ───          │
│                                     │
│  Service Name *                     │
│  ┌─────────────────────────────┐    │
│  │                             │    │  ← Text input
│  └─────────────────────────────┘    │
│                                     │
│  Price *                            │
│  ┌───────────┐  ┌───────────────┐   │
│  │  $ 0.00   │  │  Monthly  ▼   │   │  ← Price + billing cycle
│  └───────────┘  └───────────────┘   │
│                                     │
│  Next Billing Date *                │
│  ┌─────────────────────────────┐    │
│  │  📅  Select date...         │    │  ← Native date picker
│  └─────────────────────────────┘    │
│                                     │
│  Category (optional)                │
│  ┌─────────────────────────────┐    │
│  │  Entertainment  ▼           │    │  ← Dropdown
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │       Add Subscription       │    │  ← Primary CTA button
│  └─────────────────────────────┘    │
│                                     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Drag Handle**: Visual affordance for dismissing sheet by swiping down.
- **Popular Services**: Quick-select templates to reduce manual entry.
- **Minimal Fields**: Only 3 required fields (name, price, date) to meet <30 second goal.
- **Cycle Picker**: Dropdown with Weekly, Monthly, Yearly options.
- **Native Date Picker**: Uses OS native picker for familiarity on mobile.
- **Primary CTA**: Full-width button in coral accent color.

---

### Screen 3: Service Templates

Full-screen view for browsing all available templates.

```
┌─────────────────────────────────────┐
│  ←  Choose a Service                │  ← Back button + title
├─────────────────────────────────────┤
│                                     │
│  [ 🔍 Search services...        ]   │  ← Search input
│                                     │
├─────────────────────────────────────┤
│  Streaming                          │  ← Category header
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌──   │
│  │ 📺 │ │ 🎵 │ │ 📺 │ │ 🎵 │ │    │
│  │Netf│ │Spot│ │Hulu│ │Apple│ │    │  ← Horizontal scroll
│  │$15 │ │$10 │ │$8  │ │$11  │ │    │
│  └────┘ └────┘ └────┘ └────┘ └──   │
│                                     │
│  Productivity                       │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌──   │
│  │ 📝 │ │ ☁️ │ │ 📊 │ │ 💬 │ │    │
│  │Noti│ │Adob│ │Msft│ │Slac│ │    │
│  │$8  │ │$55 │ │$10 │ │$8  │ │    │
│  └────┘ └────┘ └────┘ └────┘ └──   │
│                                     │
│  Gaming                             │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌──   │
│  │ 🎮 │ │ 🎮 │ │ 🎮 │ │ 🎮 │ │    │
│  │Xbox│ │PS+ │ │Nint│ │EA  │ │    │
│  └────┘ └────┘ └────┘ └────┘ └──   │
│                                     │
│  Cloud Storage                      │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌──   │
│  │ ☁️ │ │ ☁️ │ │ ☁️ │ │ ☁️ │ │    │
│  │iClo│ │Goog│ │Drop│ │OnDr│ │    │
│  └────┘ └────┘ └────┘ └────┘ └──   │
│                                     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Search**: Filter templates by name.
- **Categorized Grid**: Grouped by type for easy browsing.
- **Template Cards**: Show logo, name, typical price.
- **Tap to Select**: Opens add form with pre-filled values.

---

### Screen 4: Subscription Detail

View and edit an existing subscription.

```
┌─────────────────────────────────────┐
│  ←  Subscription Details   [ 🗑️ ]   │  ← Back + delete button
├─────────────────────────────────────┤
│                                     │
│              ┌────────┐             │
│              │   🎵   │             │  ← Large service icon
│              │ Spotify│             │
│              └────────┘             │
│                                     │
│         $9.99 / month               │  ← Price prominently displayed
│                                     │
│  ─────────────────────────────────  │
│                                     │
│  Next billing                       │
│  January 22, 2025        [ Edit ]   │
│                                     │
│  Billing cycle                      │
│  Monthly                 [ Edit ]   │
│                                     │
│  Category                           │
│  Entertainment           [ Edit ]   │
│                                     │
│  ─────────────────────────────────  │
│                                     │
│  Reminder                           │
│  ┌────────────────────────────┐     │
│  │  ○ Off  ○ 1 day  ● 3 days  │     │  ← Reminder toggle
│  └────────────────────────────┘     │
│                                     │
│  ─────────────────────────────────  │
│                                     │
│  Notes                              │
│  Family plan, shared with 5 others  │  ← Optional notes
│                                     │
│                                     │
│  ─────────────────────────────────  │
│                                     │
│  Added: Dec 15, 2024                │  ← Metadata
│                                     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Large Icon**: Visual identity at top.
- **Inline Editing**: Each field has an edit button for quick changes.
- **Reminder Settings**: Configure notification timing for this specific subscription.
- **Delete Button**: In header, with confirmation prompt.
- **Notes Field**: Optional context (who shares it, what plan tier, etc.).

---

### Screen 5: Upcoming Renewals

Timeline/calendar view of upcoming charges.

```
┌─────────────────────────────────────┐
│  ≡  Upcoming                        │
├─────────────────────────────────────┤
│                                     │
│  This Week                          │  ← Section: This week
│  ───────────────────                │
│                                     │
│  Tomorrow, Jan 19                   │
│  ┌─────────────────────────────┐    │
│  │ 💼  Adobe CC        $54.99  │    │
│  └─────────────────────────────┘    │
│                                     │
│  Wednesday, Jan 22                  │
│  ┌─────────────────────────────┐    │
│  │ 🎵  Spotify          $9.99  │    │
│  └─────────────────────────────┘    │
│                                     │
│  Next Week                          │  ← Section: Next week
│  ───────────────────                │
│                                     │
│  Saturday, Jan 25                   │
│  ┌─────────────────────────────┐    │
│  │ 📺  Netflix         $15.49  │    │
│  └─────────────────────────────┘    │
│                                     │
│  Later This Month                   │  ← Section: Rest of month
│  ───────────────────                │
│                                     │
│  February 1                         │
│  ┌─────────────────────────────┐    │
│  │ ☁️  iCloud+          $2.99  │    │
│  └─────────────────────────────┘    │
│  ┌─────────────────────────────┐    │
│  │ 🎮  Xbox Game Pass  $14.99  │    │
│  └─────────────────────────────┘    │
│                                     │
│  ─────────────────────────────────  │
│  Total upcoming (30 days): $98.45   │  ← Summary of upcoming charges
│                                     │
├─────────────────────────────────────┤
│  🏠        📅        ⚙️            │
│  Home    Upcoming   Settings        │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Grouped by Time**: This week, next week, later - matches mental model.
- **Date Headers**: Clear date labeling for each renewal.
- **30-Day Summary**: Total upcoming spend at bottom.
- **Tappable Cards**: Open detail view on tap.

---

### Screen 6: Settings

App preferences and data management.

```
┌─────────────────────────────────────┐
│  ≡  Settings                        │
├─────────────────────────────────────┤
│                                     │
│  Appearance                         │
│  ───────────────────                │
│                                     │
│  Theme                              │
│  ┌────────────────────────────┐     │
│  │ ○ Light  ○ Dark  ● System  │     │  ← Theme picker
│  └────────────────────────────┘     │
│                                     │
│  Currency                           │
│  USD ($)                       >    │
│                                     │
│  Notifications                      │
│  ───────────────────                │
│                                     │
│  Renewal Reminders              🔘  │  ← Toggle: on/off
│                                     │
│  Default Reminder Time              │
│  3 days before               >      │
│                                     │
│  Data                               │
│  ───────────────────                │
│                                     │
│  Export Subscriptions          >    │  ← Future: CSV export
│                                     │
│  Clear All Data                >    │  ← Destructive, needs confirm
│                                     │
│  About                              │
│  ───────────────────                │
│                                     │
│  Version                            │
│  1.0.0                              │
│                                     │
│  Privacy Policy                >    │
│                                     │
│  Terms of Service              >    │
│                                     │
│  Made with ♥ by [Team]              │
│                                     │
├─────────────────────────────────────┤
│  🏠        📅        ⚙️            │
│  Home    Upcoming   Settings        │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Grouped Sections**: Appearance, Notifications, Data, About.
- **Theme Toggle**: Light, Dark, System (default: System).
- **Reminder Default**: Sets default for new subscriptions.
- **Data Management**: Export and clear options.

---

### Screen 7: Empty State (First-Time User)

When user has no subscriptions yet.

```
┌─────────────────────────────────────┐
│  ≡  SubTrack              [☀/🌙]   │
├─────────────────────────────────────┤
│                                     │
│                                     │
│                                     │
│           ┌─────────┐               │
│           │         │               │
│           │   📋    │               │  ← Friendly illustration
│           │   ✨    │               │
│           │         │               │
│           └─────────┘               │
│                                     │
│     No subscriptions yet            │  ← Clear headline
│                                     │
│   Add your first subscription to    │
│   start tracking your spending.     │  ← Helpful subtext
│                                     │
│                                     │
│  ┌─────────────────────────────┐    │
│  │    + Add Subscription        │    │  ← Primary CTA
│  └─────────────────────────────┘    │
│                                     │
│                                     │
│    Or choose from popular:          │
│                                     │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐       │
│  │ 📺 │ │ 🎵 │ │ ☁️ │ │ 🎮 │       │  ← Quick templates
│  │Netf│ │Spot│ │iClo│ │Xbox│       │
│  └────┘ └────┘ └────┘ └────┘       │
│                                     │
│                                     │
├─────────────────────────────────────┤
│  🏠        📅        ⚙️            │
│  Home    Upcoming   Settings        │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Illustration**: Friendly, non-empty visual.
- **Clear Message**: Explains what to do.
- **Prominent CTA**: Large button to add first subscription.
- **Quick Templates**: One-tap options for common services.

---

## User Flows

### Primary Flow: Add First Subscription

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Open App  │ ──▶ │ Empty State │ ──▶ │ Tap "Add"   │ ──▶ │  Add Sheet  │
│             │     │  Displayed  │     │   Button    │     │   Opens     │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                                                   │
                    ┌─────────────┐     ┌─────────────┐            │
                    │  Dashboard  │ ◀── │   Success!  │ ◀──────────┘
                    │ Shows Card  │     │  Card Added │     (Fill form, tap save)
                    │ + Total $   │     │             │
                    └─────────────┘     └─────────────┘
```

**Steps**:
1. User opens app → sees empty state with illustration
2. User taps "Add Subscription" button or a quick template
3. Add sheet slides up from bottom
4. User fills required fields (name, price, date) OR selects template
5. User taps "Add Subscription" button
6. Success haptic/animation → Sheet dismisses
7. Dashboard shows new card with updated total spend

**Target Time**: < 30 seconds from open to first subscription added

---

### Secondary Flow: Review and Cancel Subscription

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Dashboard  │ ──▶ │ Tap on Card │ ──▶ │   Detail    │ ──▶ │ Tap Delete  │
│             │     │             │     │    View     │     │   Button    │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                                                   │
                    ┌─────────────┐     ┌─────────────┐            │
                    │  Undo Toast │ ◀── │  Confirm    │ ◀──────────┘
                    │   Appears   │     │   Delete?   │
                    │  (5 sec)    │     │   [Cancel]  │
                    └─────────────┘     │   [Delete]  │
                           │            └─────────────┘
                           ▼
                    ┌─────────────┐
                    │  Dashboard  │
                    │  Updated    │
                    └─────────────┘
```

**Steps**:
1. User views dashboard with subscription list
2. User taps on a subscription card
3. Detail view opens with all subscription info
4. User taps delete icon in header
5. Confirmation modal appears: "Delete Spotify? This cannot be undone."
6. User confirms deletion
7. Card animates out, undo toast appears for 5 seconds
8. Total spend updates immediately

---

### Flow: Quick Add via Template

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Dashboard  │ ──▶ │  Tap FAB    │ ──▶ │  Add Sheet  │ ──▶ │ Tap Template│
│             │     │   (+)       │     │   Opens     │     │  (Netflix)  │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                                                   │
                    ┌─────────────┐     ┌─────────────┐            │
                    │  Dashboard  │ ◀── │  Pre-filled │ ◀──────────┘
                    │  Updated    │     │ Form - User │
                    │             │     │ Picks Date  │
                    └─────────────┘     └─────────────┘
```

**Steps**:
1. User taps FAB (+) on dashboard
2. Add sheet opens with popular templates visible
3. User taps "Netflix" template
4. Form auto-fills: Name=Netflix, Price=$15.49, Cycle=Monthly
5. User only needs to select next billing date
6. User taps "Add Subscription"
7. Done - total time: ~10 seconds

---

### Flow: Check Upcoming Renewals

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Dashboard  │ ──▶ │ Tap Upcoming│ ──▶ │  Upcoming   │
│             │     │  in Nav     │     │    View     │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │ See grouped │
                                        │  renewals   │
                                        │ (this week, │
                                        │ next week)  │
                                        └─────────────┘
```

---

### Flow: Edit Subscription

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Dashboard  │ ──▶ │  Tap Card   │ ──▶ │   Detail    │ ──▶ │ Tap "Edit"  │
│             │     │             │     │    View     │     │  on field   │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                                                   │
                    ┌─────────────┐     ┌─────────────┐            │
                    │   Detail    │ ◀── │  Edit Modal │ ◀──────────┘
                    │  Updated    │     │  or Inline  │
                    └─────────────┘     └─────────────┘
```

---

## Interaction Patterns

| Action | Element | Response | Notes |
|--------|---------|----------|-------|
| Tap | Subscription Card | Navigate to detail view | Subtle press effect (scale 0.98) |
| Tap | FAB (+) | Open add sheet from bottom | Spring animation |
| Tap | Bottom nav item | Switch view | Smooth crossfade |
| Swipe Left | Subscription Card | Reveal delete action | Red background slides in |
| Swipe Down | Add Sheet | Dismiss sheet | Elastic snap back or dismiss |
| Long Press | Subscription Card | Open quick action menu | Haptic feedback |
| Pull Down | Subscription List | Refresh (future: sync) | Custom branded loader |
| Tap | Filter Pill | Toggle filter on/off | Pill fills with accent color |

---

## States

### Empty State
- **When**: User has 0 subscriptions
- **Show**: Illustration + "No subscriptions yet" + CTA button + quick templates
- **Action**: Prominent "Add Subscription" button

### Loading State
- **When**: Initial app load, refreshing data
- **Show**: Skeleton cards (pulse animation) - 3 placeholder cards
- **Duration**: Should be <1 second for local data

### Error State
- **When**: Failed to save, failed to load
- **Show**: Inline error message with retry button
- **Tone**: Friendly, not technical ("Something went wrong. Try again?")

### Success State
- **When**: Subscription added/edited/deleted
- **Show**: Brief haptic feedback + toast confirmation
- **Duration**: Toast auto-dismisses after 3 seconds (except undo which lasts 5s)

### Due Soon State
- **When**: Subscription due within 7 days
- **Show**: Warning icon on card + optional badge on Upcoming nav
- **Color**: Amber/orange accent (not red - not an error)

### Free Trial State (P1)
- **When**: Subscription marked as trial
- **Show**: "TRIAL" badge on card + countdown ("Ends in 5 days")
- **Color**: Distinct color (blue or purple)

---

## Gesture Reference

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  TAP          SWIPE LEFT      SWIPE DOWN      LONG PRESS       │
│  ┌───┐        ──────▶         │               ┌───────┐        │
│  │ · │        [Delete]        │               │ · · · │        │
│  └───┘                        ▼               │ (hold)│        │
│  Open         Reveal          Dismiss         └───────┘        │
│  Detail       Delete          Sheet           Context          │
│               Action                          Menu             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Touch Target Specifications

All interactive elements follow minimum 44x44px touch targets:

| Element | Minimum Size | Actual Size | Notes |
|---------|--------------|-------------|-------|
| FAB | 44x44px | 56x56px | Larger for primary action |
| Subscription Card | 44px height | 72px height | Full-width tap area |
| Bottom Nav Item | 44x44px | 64x48px | With label |
| Filter Pill | 44px height | 32px height, 44px tap | Padding extends tap area |
| Close Button | 44x44px | 44x44px | Minimum |
| List Row | 44px height | 56px height | With padding |

---

## Thumb Zone Optimization

Mobile-first design keeps primary actions in easy thumb reach:

```
┌─────────────────────────────────────┐
│                                     │
│         Hard to reach               │  ← Secondary info only
│                                     │
├─────────────────────────────────────┤
│                                     │
│       Scrollable content            │  ← Subscription list
│                                     │
│                                     │
├─────────────────────────────────────┤
│         Easy to reach               │  ← FAB in this zone
│              ( + )                  │
├─────────────────────────────────────┤
│  🏠        📅        ⚙️            │  ← Navigation in thumb zone
│  Natural thumb zone                 │
└─────────────────────────────────────┘
```

---

## Accessibility Considerations

### Keyboard Navigation
- Tab order follows visual hierarchy: header → filters → list → FAB → nav
- Enter/Space activates buttons and cards
- Escape closes modals and sheets
- Arrow keys navigate within lists

### Screen Reader Support
- All icons have aria-labels (e.g., "Add subscription" not just "plus")
- Cards announce: "Spotify, $9.99 per month, due January 22"
- Total spend announced on page load
- Live regions announce changes (subscription added/deleted)

### Visual Accessibility
- Minimum 4.5:1 contrast for all text (WCAG AA)
- Don't rely on color alone - icons accompany status colors
- Focus indicators visible in both light/dark modes
- Reduced motion option respects prefers-reduced-motion

### Touch Accessibility
- All touch targets minimum 44x44px
- Adequate spacing between interactive elements (8px minimum)
- No time-limited interactions (except undo toast, which is generous 5s)

---

## Responsive Breakpoints

While mobile-first, the design scales gracefully:

| Breakpoint | Width | Layout Changes |
|------------|-------|----------------|
| Mobile | < 640px | Single column, bottom nav, full-width cards |
| Tablet | 640-1024px | 2-column card grid, bottom nav persists |
| Desktop | > 1024px | Sidebar nav replaces bottom nav, 3-column grid, modals instead of sheets |

---

## Animation Specifications

| Animation | Duration | Easing | Properties |
|-----------|----------|--------|------------|
| Card press | 100ms | ease-out | scale: 0.98 |
| Card add | 300ms | ease-out | opacity: 0→1, translateY: 20→0 |
| Card delete | 250ms | ease-in | opacity: 1→0, translateX: 0→-100% |
| Sheet open | 300ms | ease-out | translateY: 100%→0 |
| Sheet close | 200ms | ease-in | translateY: 0→100% |
| Total counter | 400ms | ease-out | number tween |
| Theme switch | 200ms | ease-in-out | background-color, color |
| FAB press | 150ms | ease-out | scale: 0.95 |

---

## Design Tokens Reference

Based on UI research, key tokens for implementation:

### Spacing Scale
- 4px (xs), 8px (sm), 12px (md), 16px (lg), 24px (xl), 32px (2xl)

### Border Radius
- Cards: 12px
- Buttons: 8px
- Inputs: 8px
- Modals: 16px
- Pills/Tags: 9999px (full)

### Typography Scale
- Total spend: 32px, bold
- Card title: 16px, semibold
- Card secondary: 14px, regular
- Labels: 12px, medium
- Body: 14px, regular

---

Status: READY_FOR_REVIEW

PHASE_COMPLETE: planning
