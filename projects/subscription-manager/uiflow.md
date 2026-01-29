# SubTrack UI Flow Documentation

## Overview

SubTrack is a mobile-first subscription management web application. This document describes the user interface flows and interactions.

---

## Screen Map

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐       │
│  │   Login     │────▶│  Register   │     │  Dashboard  │       │
│  └─────────────┘     └─────────────┘     └──────┬──────┘       │
│         │                    │                   │               │
│         └────────────────────┼───────────────────┘               │
│                              │                                   │
│                              ▼                                   │
│                     ┌─────────────┐                             │
│                     │  Dashboard  │                             │
│                     │   (Home)    │                             │
│                     └──────┬──────┘                             │
│                            │                                     │
│         ┌──────────────────┼──────────────────┐                 │
│         │                  │                  │                 │
│         ▼                  ▼                  ▼                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │  Upcoming   │    │ Add Sheet   │    │  Settings   │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│         │                  │                                     │
│         ▼                  │                                     │
│  ┌─────────────┐           │                                     │
│  │Subscription │◀──────────┘                                     │
│  │   Detail    │                                                 │
│  └─────────────┘                                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## User Flows

### 1. First-Time User Flow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   App Load   │────▶│  Empty State │────▶│  Add Sheet   │
│              │     │  (No subs)   │     │   Opens      │
└──────────────┘     └──────────────┘     └──────────────┘
                                                 │
                     ┌──────────────┐            │
                     │   Dashboard  │◀───────────┘
                     │ (First sub!) │   (After save)
                     └──────────────┘
```

**Steps:**
1. User opens app → Empty state shows with illustration
2. User taps "Add Subscription" or popular template
3. Add sheet slides up from bottom
4. User fills form (name, price, date)
5. User taps "Add Subscription"
6. Sheet closes → Dashboard shows first card

---

### 2. Add Subscription Flow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Tap FAB (+) │────▶│  Add Sheet   │────▶│ Fill Form or │
│              │     │   Opens      │     │ Pick Template│
└──────────────┘     └──────────────┘     └──────────────┘
                                                 │
                     ┌──────────────┐            │
                     │   Success    │◀───────────┘
                     │  Toast + Card│   (Tap "Add")
                     └──────────────┘
```

**Add Sheet Contents:**
- Search bar for templates
- Popular templates (horizontal scroll)
- Divider: "or enter manually"
- Form fields: Name, Price, Cycle, Date, Category
- "Add Subscription" button

---

### 3. View/Edit Subscription Flow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ Tap Card on  │────▶│  Detail View │────▶│  Edit Sheet  │
│  Dashboard   │     │   Opens      │     │   Opens      │
└──────────────┘     └──────────────┘     └──────────────┘
                            │                    │
                            │                    │
                            ▼                    ▼
                     ┌──────────────┐     ┌──────────────┐
                     │   Delete     │     │    Save      │
                     │  Confirm     │     │   Changes    │
                     └──────────────┘     └──────────────┘
```

**Detail View Shows:**
- Large icon + name
- Price with billing cycle
- Next billing date (with days remaining)
- Category
- Notes (if any)
- Edit buttons on each row
- Delete button in header

---

### 4. Delete Subscription Flow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ Tap Delete   │────▶│  Confirm     │────▶│   Delete     │
│   Button     │     │  Dialog      │     │  Completed   │
└──────────────┘     └──────────────┘     └──────────────┘
                            │                    │
                            │                    ▼
                            │              ┌──────────────┐
                            └─────────────▶│  Dashboard   │
                              (Cancel)     │  Updated     │
                                           └──────────────┘
```

---

### 5. Check Upcoming Renewals Flow

```
┌──────────────┐     ┌──────────────┐
│ Tap "Upcoming"│────▶│  Upcoming    │
│ in Bottom Nav │     │   View       │
└──────────────┘     └──────────────┘
                            │
                            ▼
                     ┌──────────────────────────┐
                     │  Grouped by:             │
                     │  - This Week             │
                     │  - Next Week             │
                     │  - Later This Month      │
                     │                          │
                     │  Total at bottom         │
                     └──────────────────────────┘
```

---

## Component Interactions

### Bottom Navigation
| Tab | Icon | Action |
|-----|------|--------|
| Home | 🏠 | Show dashboard with spend summary + list |
| Upcoming | 📅 | Show upcoming renewals grouped by time |
| Settings | ⚙️ | Show settings (theme, data, account) |

### FAB (Floating Action Button)
- Position: Bottom right, above nav
- Action: Opens Add Subscription sheet
- Color: Primary coral (#FF7A5C)

### Subscription Card
- Tap: Navigate to detail view
- Shows: Icon, name, price, next date, due soon indicator

### Filter Pills
- Tap "All": Show all subscriptions
- Tap category: Filter to that category only
- Active state: Filled with primary color

### Theme Toggle
- Light → Dark → System → Light
- Updates immediately
- Persists in localStorage

---

## States

### Empty State
- **When**: No subscriptions exist
- **Shows**:
  - Illustration (📋✨)
  - "No subscriptions yet" heading
  - Helpful subtext
  - Primary "Add Subscription" button
  - Quick template buttons

### Loading State
- **When**: Fetching data
- **Shows**: Skeleton cards (pulse animation)

### Error State
- **When**: API error
- **Shows**: Error toast with message

### Due Soon State
- **When**: Subscription due within 7 days
- **Shows**:
  - Warning icon (⚠️) on card
  - Amber color for date text
  - "Due in X days" instead of date

---

## Gestures

| Gesture | Element | Action |
|---------|---------|--------|
| Tap | Card | Open detail view |
| Tap | FAB | Open add sheet |
| Tap | Nav item | Switch view |
| Swipe down | Sheet | Dismiss sheet |

---

## Responsive Behavior

### Mobile (< 640px)
- Single column layout
- Bottom navigation
- Full-width cards
- Bottom sheets for forms

### Tablet (640-1024px)
- Two-column card grid
- Wider sheets

### Desktop (> 1024px)
- Sidebar navigation (left)
- Three-column grid
- Modal dialogs instead of sheets

---

## Accessibility

- All interactive elements: min 44x44px touch targets
- Focus indicators on keyboard navigation
- Screen reader labels for icons
- Reduced motion support
- WCAG AA contrast compliance
