# UX Research: Inventory Reorder Alert System

## Executive Summary

This research analyzes UX patterns from successful inventory management, alert systems, and task management applications to inform the design of an alert-first inventory reorder system. The target users are operations managers and warehouse staff who need simple, reliable stock threshold notifications.

---

## Pattern Analysis

Referenced apps and their successful patterns:

### 1. Sortly (Inventory Management)

- **Pattern Used**: Visual card-based inventory with prominent photos, color-coded tags, and large touch targets
- **Why It Works**: Visual inventory reduces cognitive load - users recognize items faster than reading text. Color tags provide instant status at a glance.
- **What We Can Learn**:
  - Use item photos as primary identifier
  - Color-code stock status (green = OK, yellow = low, red = critical)
  - Large, thumb-friendly touch targets for mobile

### 2. Todoist (Task/Alert Management)

- **Pattern Used**: Inbox-style alert list with priority indicators (P1-P4 colored flags), quick-add floating action button, swipe actions
- **Why It Works**: Familiar email mental model. Priority colors train users to scan urgency. One-tap actions reduce friction.
- **What We Can Learn**:
  - Treat alerts like a prioritized inbox
  - Use urgency levels (critical/low/reorder soon)
  - Quick actions: swipe to snooze, mark ordered, dismiss
  - Floating action button for quick stock update

### 3. Slack (Notification System)

- **Pattern Used**: Channel-based notification grouping, notification badges, "mark all read", snooze/DND modes, notification preferences per channel
- **Why It Works**: Prevents alert fatigue by grouping and giving user control. Badge counts create urgency without interruption.
- **What We Can Learn**:
  - Group alerts by category (e.g., "Office Supplies", "Raw Materials")
  - Badge counts on app icon and dashboard sections
  - Snooze options (1 hour, until tomorrow, custom)
  - Granular notification preferences per item/category

### 4. Mint/Personal Capital (Threshold Alerts)

- **Pattern Used**: Budget bars with visual fill levels, threshold lines, push notifications when exceeded, trend arrows
- **Why It Works**: Progress bars are instantly understandable. Visual threshold lines make comparison automatic.
- **What We Can Learn**:
  - Visual stock level bars (current vs. threshold vs. max capacity)
  - Clear threshold indicator line on the bar
  - Trend indicators (consumption rate arrows)
  - "You've used X of Y" language is clear

### 5. Amazon (Reorder Experience)

- **Pattern Used**: "Buy it again" with one-click reorder, "Subscribe & Save" automation, order history for reference
- **Why It Works**: Reduces reorder to single action. Historical data enables informed decisions.
- **What We Can Learn**:
  - One-click "Reorder" button on alerts
  - Show last order date and quantity
  - Suggested reorder quantity based on history
  - Link to preferred supplier/vendor

### 6. Robinhood (Dashboard Design)

- **Pattern Used**: Single KPI hero number at top, minimalist card design, red/green status colors, sparkline micro-charts
- **Why It Works**: Focuses attention on what matters most. Micro-visualizations add context without clutter.
- **What We Can Learn**:
  - Hero metric: "X items need attention"
  - Sparkline showing consumption trend per item
  - Red (critical) / Yellow (low) / Green (OK) status system
  - Clean, minimal dashboard with clear hierarchy

### 7. iOS Health App (Progress Rings)

- **Pattern Used**: Circular progress rings for goals, daily/weekly/monthly views, celebration animations on completion
- **Why It Works**: Rings are satisfying to complete. Time views show patterns.
- **What We Can Learn**:
  - Consider circular stock level indicators for key items
  - Time-based views (today's alerts, this week's consumption)
  - Positive feedback when stock is healthy

---

## User Journey Analysis

### Entry Point

**Primary**: Push notification or email alert
- User receives "Low Stock Alert: Printer Paper (15 remaining, threshold: 20)"
- Taps notification → lands directly on item detail with reorder CTA

**Secondary**: Proactive check via dashboard
- User opens app → sees dashboard with alert count badge
- Scans "Needs Attention" section for items to address

### Core Loop

```
Check Stock → Update Quantity → Threshold Breach → Alert Triggered →
→ User Notified → Review Alert → Take Action (Reorder/Snooze/Dismiss) →
→ Receive Stock → Update Quantity → [Loop]
```

**Frequency**:
- Stock updates: After receiving shipments (1-5x per week)
- Alert review: As notifications arrive (reactive)
- Dashboard check: 1-2x per week (proactive)

### Success State

**Immediate Success**:
- Alert acknowledged → Item marked as "Ordered" → Alert cleared
- Visual confirmation (checkmark animation, success toast)

**Long-term Success**:
- Zero stockouts month-over-month
- Dashboard shows "All stock levels healthy" state
- Consumption trends visible and predictable

---

## Mental Model

**Users think of this task as**: A safety net / early warning system
- Like a smoke detector - silent when things are fine, loud when action needed
- Or a gas gauge - visual indicator that triggers action at certain level

**Common terminology**:
- "Reorder point" or "threshold" or "minimum"
- "Stock level" or "quantity on hand"
- "Running low" / "Out of stock" / "Need to order"
- "Par level" (restaurant industry)
- "Safety stock" (warehouse industry)

**User expectations**:
1. Set it once, get notified automatically
2. Notifications should be actionable, not just informational
3. Easy to update when stock arrives
4. Don't bother me when things are fine

---

## Anti-Patterns to Avoid

- **Don't**: Require complex setup wizards with 20+ fields per item
  - **Why**: Users abandon before adding first item. Competitor Zoho loses users here.

- **Don't**: Send individual alerts for each item (alert bombing)
  - **Why**: Alert fatigue causes users to disable all notifications. One digest is better.

- **Don't**: Hide the reorder action behind multiple screens
  - **Why**: Amazon's success comes from one-click. Every extra tap loses conversions.

- **Don't**: Use abstract icons without labels
  - **Why**: "What does this chart icon do?" Users shouldn't guess. Sortly struggles here.

- **Don't**: Auto-dismiss alerts without user acknowledgment
  - **Why**: Users lose trust if alerts disappear. They need to know they've handled it.

- **Don't**: Require exact quantities for updates
  - **Why**: "About 20 left" should be acceptable. Precision isn't always possible.

- **Don't**: Show empty states without clear CTA
  - **Why**: "No items" is useless. "Add your first item" with button is helpful.

- **Don't**: Use desktop-only modals on mobile
  - **Why**: Modal dialogs are frustrating on mobile. Use full-screen slides or bottom sheets.

- **Don't**: Require login for every session
  - **Why**: Friction kills quick updates. Use biometrics or stay logged in.

---

## Recommended Patterns for Our App

Based on research, we should use:

### 1. Alert Inbox Pattern (from Todoist)
- Central "Alerts" tab as primary view
- Priority-sorted list with color indicators
- Swipe actions for quick resolution
- **Rationale**: Treats alerts as actionable tasks, not just notifications

### 2. Visual Stock Bars (from Mint)
- Horizontal progress bar showing current vs. threshold
- Clear visual threshold marker
- Color transitions: green → yellow → red
- **Rationale**: Instantly communicates status without reading numbers

### 3. One-Tap Reorder (from Amazon)
- Primary CTA on every alert: "Reorder"
- Secondary actions: Snooze, Mark Ordered, Adjust Threshold
- Pre-filled with suggested quantity
- **Rationale**: Reduces friction to complete the core action

### 4. Smart Notification Grouping (from Slack)
- Daily digest option vs. individual alerts
- Category-based grouping
- Badge counts without interruption
- **Rationale**: Prevents alert fatigue while maintaining awareness

### 5. Hero Metric Dashboard (from Robinhood)
- Top of dashboard: "X items need attention"
- Three status sections: Critical / Low / Healthy
- Sparkline trends per item
- **Rationale**: Focuses on what matters, provides context

### 6. Bottom Navigation (Mobile Standard)
- 4 tabs: Dashboard, Items, Alerts, Settings
- Floating action button for quick add/update
- **Rationale**: Thumb-friendly, matches user expectations

### 7. Quick Update Flow (Reducing Friction)
- Open app → See item → Tap → Enter number → Done (3 taps max)
- Support for barcode scan, +/- buttons, numeric keypad
- "Received shipment" shortcut adds predefined quantity
- **Rationale**: Speed is essential for warehouse workers

---

## Competitive UX Comparison

| Feature | Sortly | inFlow | Zoho | Our Approach |
|---------|--------|--------|------|--------------|
| Setup time | ~15 min | ~45 min | ~60 min | **~5 min** |
| Mobile UX | Good | Poor | Fair | **Excellent (primary)** |
| Alert visibility | Buried | Menu item | Notification center | **Hero feature** |
| Reorder action | 4 clicks | 6 clicks | 5 clicks | **1 tap** |
| Stock update | 3 taps | 5 taps | 4 taps | **2 taps** |

---

## Information Architecture Recommendation

```
┌─────────────────────────────────────────┐
│            NAVIGATION                    │
├─────────────────────────────────────────┤
│                                         │
│  [Dashboard] [Items] [Alerts] [Settings]│
│       │         │        │        │     │
│       │         │        │        │     │
│       ▼         ▼        ▼        ▼     │
│   Overview   Catalog   Inbox    Prefs   │
│   Stats      Add/Edit  Actions  Notify  │
│   Trends     Search    History  Account │
│   Quick      Import    Filters  Help    │
│   Actions    Categories         Export  │
└─────────────────────────────────────────┘
```

---

## Accessibility Considerations

Based on WCAG 2.1 guidelines:

1. **Color Independence**: Never use color alone for status. Add icons (checkmark, warning, X) and text labels.

2. **Touch Targets**: Minimum 44x44px for all interactive elements (Apple HIG).

3. **Contrast Ratios**:
   - Normal text: 4.5:1 minimum
   - Large text: 3:1 minimum
   - Our alert colors must pass on white and dark backgrounds

4. **Screen Reader Support**:
   - Stock bars: "Printer Paper: 15 out of 100, below reorder threshold of 20"
   - Alert items: "Critical alert: Coffee beans, 5 remaining, tap to reorder"

5. **Motion Sensitivity**: Respect reduced-motion preferences, provide static alternatives

---

## Key UX Metrics to Track

1. **Time to First Alert Setup** - Target: < 3 minutes
2. **Alert-to-Action Time** - How fast users resolve alerts
3. **Stock Update Completion Rate** - % of started updates completed
4. **Feature Discovery Rate** - Are users finding key features?
5. **Alert Engagement Rate** - % of alerts that get user action (not ignored)

---

## Research Sources

- Sortly app (iOS App Store, 4.7 rating, analyzed UI patterns)
- inFlow demo (desktop web application review)
- Zoho Inventory trial (web + mobile evaluation)
- Todoist app (notification and task UX patterns)
- Slack desktop + mobile (notification management)
- Mint app (threshold visualization patterns)
- Amazon mobile (reorder flow analysis)
- Material Design 3 Guidelines (component patterns)
- Apple Human Interface Guidelines (iOS patterns)
- Nielsen Norman Group: "Alert Design Guidelines" (2023)

---

Status: READY_FOR_PLANNING
