# UX Specification: Inventory Reorder Alert System

Based on research in: `.plans/research/ux-research.md`
Implementing features from: `.plans/specs/product-spec.md`
Visual direction from: `.plans/research/ui-research.md`

---

## Screen List

1. **Dashboard** - Hero metrics overview with items needing attention
2. **Items List** - Catalog of all tracked inventory items
3. **Item Detail** - Individual item view with stock bar and quick update
4. **Add/Edit Item** - Form for creating or modifying items
5. **Alerts Inbox** - Action-oriented list of all active threshold breaches
6. **Quick Update Modal** - Minimal overlay for fast quantity changes
7. **Settings** - Notification preferences and account management
8. **Empty States** - First-time user guidance screens
9. **Login/Register** - Authentication screens

---

## Wireframes

### Screen 1: Dashboard (Home)

```
┌─────────────────────────────────────┐
│  ☰  Inventory Alerts        [👤]   │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────┐   │
│  │    🔔 3 items need          │   │
│  │       attention             │   │
│  │    ─────────────────────    │   │
│  │    1 Critical · 2 Low       │   │
│  └─────────────────────────────┘   │
│                                     │
│  NEEDS ATTENTION                    │
│  ┌─────────────────────────────┐   │
│  │ ⚠ Coffee Beans        [→]  │   │
│  │ ▓▓░░░░░░░░░░░░░░░░ 5/50    │   │
│  │ Critical · Reorder now      │   │
│  └─────────────────────────────┘   │
│  ┌─────────────────────────────┐   │
│  │ ⚡ Printer Paper       [→]  │   │
│  │ ▓▓▓▓░░░░░░░░░░░░░░ 15/100  │   │
│  │ Low · Below threshold       │   │
│  └─────────────────────────────┘   │
│  ┌─────────────────────────────┐   │
│  │ ⚡ Toner Cartridge     [→]  │   │
│  │ ▓▓▓░░░░░░░░░░░░░░░ 8/40    │   │
│  │ Low · Below threshold       │   │
│  └─────────────────────────────┘   │
│                                     │
│  STOCK OVERVIEW                     │
│  ┌──────┐ ┌──────┐ ┌──────┐       │
│  │  12  │ │   3  │ │  24  │       │
│  │ Total│ │Alerts│ │  OK  │       │
│  └──────┘ └──────┘ └──────┘       │
│                                     │
│              [+ Add Item]           │
│                                     │
├─────────────────────────────────────┤
│ [🏠]    [📦]    [🔔]    [⚙️]     │
│ Home    Items   Alerts  Settings   │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Hero Metric Card**: Prominent "X items need attention" count with breakdown
- **Needs Attention Section**: Priority-sorted alerts with visual stock bars
- **Stock Bar**: Horizontal progress showing current vs. threshold (Mint pattern)
- **Status Badge**: Color-coded Critical/Low/OK indicators
- **Quick Overview Stats**: Total items, active alerts, healthy items
- **Floating Add Button**: Quick access to add new item
- **Bottom Navigation**: 4-tab structure for primary navigation

---

### Screen 2: Items List (Catalog)

```
┌─────────────────────────────────────┐
│  ←  All Items              🔍 ⊕   │
├─────────────────────────────────────┤
│  [All ▼]  [Critical] [Low] [OK]    │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────┐   │
│  │ [📷] Coffee Beans     ⚠️    │   │
│  │      5 units · Threshold: 20│   │
│  │      ▓▓░░░░░░░░░░░░░░░░░░░ │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ [📷] Printer Paper    ⚡    │   │
│  │      15 units · Threshold:25│   │
│  │      ▓▓▓▓░░░░░░░░░░░░░░░░░ │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ [📷] Desk Pens        ✓    │   │
│  │      48 units · Threshold:10│   │
│  │      ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░ │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ [📷] Hand Sanitizer   ✓    │   │
│  │      30 units · Threshold:15│   │
│  │      ▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░ │   │
│  └─────────────────────────────┘   │
│                                     │
│           [ Load More ]             │
│                                     │
├─────────────────────────────────────┤
│ [🏠]    [📦]    [🔔]    [⚙️]     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Search Icon**: Quick item lookup
- **Filter Chips**: Status-based filtering (All/Critical/Low/OK)
- **Item Cards**: Photo thumbnail, name, quantity, threshold, visual bar
- **Status Indicators**: ⚠️ Critical, ⚡ Low, ✓ OK
- **Add Button**: Quick access in header

---

### Screen 3: Item Detail

```
┌─────────────────────────────────────┐
│  ←  Coffee Beans          [✏️][🗑] │
├─────────────────────────────────────┤
│                                     │
│         ┌───────────────┐           │
│         │               │           │
│         │    [Photo]    │           │
│         │               │           │
│         └───────────────┘           │
│                                     │
│  STOCK LEVEL                        │
│  ┌─────────────────────────────┐   │
│  │                             │   │
│  │  ▓▓▓░░░░░░░░░░░░░░░░░░░░  │   │
│  │  ↑                          │   │
│  │  5 units                    │   │
│  │                    ┆        │   │
│  │              Threshold: 20  │   │
│  │                             │   │
│  └─────────────────────────────┘   │
│                                     │
│  Status: ⚠️ CRITICAL                │
│  Triggered: 2 days ago              │
│                                     │
│  ┌──────────────────────────────┐  │
│  │         UPDATE STOCK          │  │
│  │  ┌─────┐ ┌────────┐ ┌─────┐  │  │
│  │  │  -  │ │   5    │ │  +  │  │  │
│  │  └─────┘ └────────┘ └─────┘  │  │
│  │                              │  │
│  │  [ Set to ] [ Add ] [ Remove ]│  │
│  │                              │  │
│  │      [ Update Quantity ]     │  │
│  └──────────────────────────────┘  │
│                                     │
│  DETAILS                            │
│  Description: Whole bean coffee     │
│               for break room        │
│  Category: Break Room Supplies      │
│  Last updated: Dec 26, 2024         │
│                                     │
│  ┌──────────────────────────────┐  │
│  │    [ 🛒 Mark as Ordered ]    │  │
│  └──────────────────────────────┘  │
│                                     │
├─────────────────────────────────────┤
│ [🏠]    [📦]    [🔔]    [⚙️]     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Item Photo**: Visual identifier at top
- **Stock Level Bar**: Large visual with threshold marker line
- **Status Badge**: Prominent status with days since triggered
- **Quick Update Controls**: +/- buttons with direct input (Sortly pattern)
- **Mode Toggle**: Set to / Add / Remove options
- **Details Section**: Description, category, last update
- **Mark as Ordered CTA**: Primary action for items needing reorder
- **Edit/Delete Actions**: Header icons for management

---

### Screen 4: Add/Edit Item Form

```
┌─────────────────────────────────────┐
│  ×  Add New Item            [Save] │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────┐   │
│  │                             │   │
│  │    [  + Add Photo  ]        │   │
│  │        (optional)           │   │
│  │                             │   │
│  └─────────────────────────────┘   │
│                                     │
│  Item Name *                        │
│  ┌─────────────────────────────┐   │
│  │ e.g., Coffee Beans          │   │
│  └─────────────────────────────┘   │
│                                     │
│  Current Quantity *                 │
│  ┌─────────────────────────────┐   │
│  │ 0                           │   │
│  └─────────────────────────────┘   │
│                                     │
│  Reorder Threshold *                │
│  ┌─────────────────────────────┐   │
│  │ e.g., 20                    │   │
│  └─────────────────────────────┘   │
│  ℹ️ Alert when stock falls below   │
│     this number                     │
│                                     │
│  Description (optional)             │
│  ┌─────────────────────────────┐   │
│  │                             │   │
│  │                             │   │
│  └─────────────────────────────┘   │
│                                     │
│  Category (optional)                │
│  ┌─────────────────────────────┐   │
│  │ Select or create...       ▼ │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌──────────────────────────────┐  │
│  │         [ Save Item ]        │  │
│  └──────────────────────────────┘  │
│                                     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Minimal Required Fields**: Only 3 required (name, quantity, threshold)
- **Photo Upload**: Optional but encouraged (Sortly pattern)
- **Helper Text**: Clarifies threshold purpose
- **Category Dropdown**: Optional organization
- **Cancel Button**: X to dismiss without saving
- **Save CTA**: Primary action button

---

### Screen 5: Alerts Inbox

```
┌─────────────────────────────────────┐
│  Alerts                    [Mark ✓]│
├─────────────────────────────────────┤
│  [Active (3)]  [Resolved]  [All]   │
├─────────────────────────────────────┤
│                                     │
│  TODAY                              │
│  ┌─────────────────────────────┐   │
│  │ ⚠️ CRITICAL                  │   │
│  │ Coffee Beans                │   │
│  │ 5 remaining · Threshold: 20 │   │
│  │ ─────────────────────────── │   │
│  │ [🛒 Reorder] [⏰ Snooze] [✓]│   │
│  └─────────────────────────────┘   │
│                           ← swipe → │
│                                     │
│  YESTERDAY                          │
│  ┌─────────────────────────────┐   │
│  │ ⚡ LOW STOCK                 │   │
│  │ Printer Paper               │   │
│  │ 15 remaining · Threshold: 25│   │
│  │ ─────────────────────────── │   │
│  │ [🛒 Reorder] [⏰ Snooze] [✓]│   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ ⚡ LOW STOCK                 │   │
│  │ Toner Cartridge             │   │
│  │ 8 remaining · Threshold: 15 │   │
│  │ ─────────────────────────── │   │
│  │ [🛒 Reorder] [⏰ Snooze] [✓]│   │
│  └─────────────────────────────┘   │
│                                     │
├─────────────────────────────────────┤
│ [🏠]    [📦]    [🔔•]   [⚙️]     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Tab Filters**: Active / Resolved / All views
- **Grouped by Date**: Today, Yesterday, This Week
- **Alert Cards**: Status badge, item name, stock info
- **Quick Actions**: Reorder, Snooze, Mark Resolved (Todoist pattern)
- **Swipe Gestures**: Swipe left = snooze, right = resolve
- **Badge on Nav**: Shows active alert count

---

### Screen 6: Quick Update Modal (Bottom Sheet)

```
┌─────────────────────────────────────┐
│                                     │
│        ┌───────────────────┐        │
│        │    ▔▔▔▔▔▔▔▔▔    │        │
│        │                   │        │
│        │   Coffee Beans    │        │
│        │   Current: 5      │        │
│        │                   │        │
│        │  ┌───┐ ┌───┐ ┌───┐│        │
│        │  │ - │ │ 5 │ │ + │ │        │
│        │  └───┘ └───┘ └───┘│        │
│        │                   │        │
│        │ ○ Set to  ● Add   │        │
│        │                   │        │
│        │  Received shipment?│        │
│        │  [ +50 ] [ +100 ] │        │
│        │                   │        │
│        │ ┌────────────────┐│        │
│        │ │    Update      ││        │
│        │ └────────────────┘│        │
│        └───────────────────┘        │
│                                     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Bottom Sheet**: Mobile-friendly slide-up modal
- **Current Value Display**: Shows existing quantity
- **+/- Controls**: Large touch targets (44px+)
- **Mode Toggle**: Set to vs. Add mode
- **Quick Presets**: Common shipment quantities
- **Update Button**: Primary action

---

### Screen 7: Settings

```
┌─────────────────────────────────────┐
│  ←  Settings                        │
├─────────────────────────────────────┤
│                                     │
│  NOTIFICATIONS                      │
│  ┌─────────────────────────────┐   │
│  │ Email Alerts            [●] │   │
│  │ Get notified when stock     │   │
│  │ falls below threshold       │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Alert Frequency         [▼]│   │
│  │ Immediate                   │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌─────────────────────────────┐   │
│  │ Daily Digest           [○]  │   │
│  │ Receive one summary email   │   │
│  │ at 8:00 AM                  │   │
│  └─────────────────────────────┘   │
│                                     │
│  ACCOUNT                            │
│  ┌─────────────────────────────┐   │
│  │ Email                    →  │   │
│  │ omar@example.com            │   │
│  └─────────────────────────────┘   │
│  ┌─────────────────────────────┐   │
│  │ Change Password          →  │   │
│  └─────────────────────────────┘   │
│                                     │
│  DATA                               │
│  ┌─────────────────────────────┐   │
│  │ Import Items (CSV)       →  │   │
│  └─────────────────────────────┘   │
│  ┌─────────────────────────────┐   │
│  │ Export Data              →  │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌──────────────────────────────┐  │
│  │         [ Log Out ]          │  │
│  └──────────────────────────────┘  │
│                                     │
├─────────────────────────────────────┤
│ [🏠]    [📦]    [🔔]    [⚙️]     │
└─────────────────────────────────────┘
```

**Key Elements**:
- **Notification Controls**: Toggle and frequency options
- **Digest Option**: Alternative to immediate alerts (Slack pattern)
- **Account Management**: Email, password change
- **Data Options**: Import/Export for migration
- **Logout**: Clear session action

---

### Screen 8: Empty States

#### First-Time Dashboard
```
┌─────────────────────────────────────┐
│  ☰  Inventory Alerts        [👤]   │
├─────────────────────────────────────┤
│                                     │
│                                     │
│            ┌───────────┐            │
│            │   📦      │            │
│            │   ___     │            │
│            │  |___|    │            │
│            └───────────┘            │
│                                     │
│       Start tracking your           │
│           inventory                 │
│                                     │
│   Add your first item to get        │
│   alerts when stock runs low.       │
│                                     │
│    ┌─────────────────────────┐     │
│    │    + Add Your First     │     │
│    │         Item            │     │
│    └─────────────────────────┘     │
│                                     │
│   Or import from CSV →              │
│                                     │
├─────────────────────────────────────┤
│ [🏠]    [📦]    [🔔]    [⚙️]     │
└─────────────────────────────────────┘
```

#### All Stock Healthy
```
┌─────────────────────────────────────┐
│                                     │
│            ┌───────────┐            │
│            │    ✓      │            │
│            │  (ring)   │            │
│            └───────────┘            │
│                                     │
│      All stock levels healthy       │
│                                     │
│   You have 24 items tracked.        │
│   We'll notify you when any         │
│   fall below threshold.             │
│                                     │
│    [ View All Items ]               │
│                                     │
└─────────────────────────────────────┘
```

---

### Screen 9: Login / Register

```
┌─────────────────────────────────────┐
│                                     │
│                                     │
│            ┌───────────┐            │
│            │   LOGO    │            │
│            │  Reorder  │            │
│            │   Alert   │            │
│            └───────────┘            │
│                                     │
│   Never run out of supplies again   │
│                                     │
│  Email                              │
│  ┌─────────────────────────────┐   │
│  │                             │   │
│  └─────────────────────────────┘   │
│                                     │
│  Password                           │
│  ┌─────────────────────────────┐   │
│  │                         👁  │   │
│  └─────────────────────────────┘   │
│                                     │
│  ┌──────────────────────────────┐  │
│  │         [ Log In ]           │  │
│  └──────────────────────────────┘  │
│                                     │
│        Forgot password?             │
│                                     │
│  ──────────── or ────────────       │
│                                     │
│  ┌──────────────────────────────┐  │
│  │     [ Create Account ]       │  │
│  └──────────────────────────────┘  │
│                                     │
└─────────────────────────────────────┘
```

---

## User Flows

### Primary Flow: Respond to Reorder Alert

```
1. User receives push/email notification
   → "Low Stock Alert: Coffee Beans (5 remaining)"

2. User taps notification
   → Lands on Item Detail screen (deep link)

3. User sees stock level and status
   → Visual bar shows critical state

4. User taps "Mark as Ordered"
   → Confirmation appears
   → Alert moves to "Resolved" state

5. (Later) Shipment arrives
   → User opens app, finds item
   → Uses quick update to add quantity

6. Stock restored above threshold
   → Status changes to green/OK
   → No more alerts for this item
```

### Secondary Flow: Proactive Stock Check

```
1. User opens app
   → Dashboard loads with hero metric

2. User sees "3 items need attention"
   → Reviews list in Needs Attention section

3. User taps item card
   → Views detail with full stock bar

4. User decides to update quantity
   → Taps +/- to adjust
   → Confirms update

5. Return to dashboard
   → Metrics update, item status changes
```

### Flow: Add First Item

```
1. New user sees empty dashboard
   → Clear CTA: "Add Your First Item"

2. User taps add button
   → Add Item form opens

3. User enters:
   → Name: "Printer Paper"
   → Quantity: 50
   → Threshold: 20

4. User taps Save
   → Returns to dashboard
   → Item appears in list

5. Success feedback
   → "Item added! You'll be alerted when stock drops below 20"
```

### Flow: Quick Stock Update (After Shipment)

```
1. User receives shipment
   → Opens app from floor

2. Searches or scrolls to item
   → Taps item card

3. Taps "+ Add" mode
   → Enters quantity received: 50

4. Taps Update
   → Stock updates instantly
   → Visual bar fills

5. Repeat for other items
   → 2-3 taps per item
```

---

## Interaction Patterns

| Action | Element | Response |
|--------|---------|----------|
| Tap | Alert card | Navigate to item detail |
| Tap | +/- buttons | Increment/decrement by 1 |
| Long press | +/- buttons | Fast increment (hold to repeat) |
| Swipe left | Alert in inbox | Reveal "Snooze" action |
| Swipe right | Alert in inbox | Mark as resolved |
| Pull down | Dashboard | Refresh data |
| Tap | Filter chip | Filter list by status |
| Tap | FAB (+ button) | Open add item form |
| Tap | Stock bar | Open quick update modal |

---

## States

### Empty State: No Items
- **When**: User has no items in catalog
- **Show**: Illustration + "Start tracking your inventory" message
- **Action**: Primary CTA to add first item, secondary link to CSV import

### Empty State: No Alerts
- **When**: All items are above threshold
- **Show**: Checkmark illustration + "All stock levels healthy"
- **Action**: Link to view all items

### Loading State
- **Show**: Skeleton placeholders for cards (not spinner)
- **Duration**: Max 2 seconds before timeout message

### Error State: Network
- **Show**: "Unable to connect" message with retry option
- **Action**: "Try Again" button

### Error State: Form Validation
- **Show**: Inline red text below field
- **Examples**: "Item name is required", "Threshold must be a number"

### Success State: Item Updated
- **Show**: Green toast notification "Quantity updated"
- **Duration**: 3 seconds, auto-dismiss

### Success State: Marked as Ordered
- **Show**: Checkmark animation + "Marked as ordered" toast
- **Visual**: Alert card transitions to resolved state

---

## Accessibility Considerations

### Keyboard Navigation
- Tab order: Header → Hero card → Alert list → Nav
- Enter/Space activates buttons and links
- Arrow keys navigate within lists
- Escape closes modals

### Screen Reader Announcements
- Stock bars: "Coffee Beans: 5 out of 50 units, below reorder threshold of 20, critical status"
- Alert cards: "Critical alert: Coffee Beans, 5 remaining, actions available"
- Status changes: "Quantity updated to 55 units"
- Modal open: "Quick update modal opened for Coffee Beans"

### Touch Targets
- All buttons: Minimum 44x44px (Apple HIG)
- +/- controls: 48x48px for easier tapping
- List items: Full width tap area
- Close buttons: 44x44px with adequate padding

### Color Contrast
- Text on backgrounds: 4.5:1 minimum
- Status indicators: Never color alone (always icon + text)
  - Critical: Red + ⚠️ icon + "Critical" label
  - Low: Yellow + ⚡ icon + "Low" label
  - OK: Green + ✓ icon + "OK" label

### Motion
- Respect `prefers-reduced-motion`
- Provide static alternatives for animations
- No auto-playing animations

---

## Responsive Breakpoints

### Mobile (320px - 767px)
- Single column layout
- Bottom navigation bar
- Full-width cards
- Bottom sheet modals
- Floating action button

### Tablet (768px - 1023px)
- Two-column item grid
- Side drawer navigation option
- Cards with more horizontal space
- Dialog modals (centered)

### Desktop (1024px+)
- Three-column layout: Nav sidebar + Main + Detail panel
- Persistent sidebar navigation
- Inline quick actions (no modals needed)
- Hover states for cards

---

## Gesture Support (Mobile)

| Gesture | Location | Action |
|---------|----------|--------|
| Pull to refresh | Dashboard, Items list | Refresh data |
| Swipe left | Alert card | Show snooze action |
| Swipe right | Alert card | Mark resolved |
| Long press | Item card | Quick actions menu |
| Pinch | Photo in detail | Zoom in/out |

---

## Animation Guidelines

### Micro-interactions
- Button press: 100ms scale to 0.97
- Toggle switch: 200ms slide with spring easing
- Card tap: 150ms subtle lift effect

### Transitions
- Screen push: 300ms slide from right
- Modal appear: 250ms slide up + fade
- Modal dismiss: 200ms slide down

### Feedback
- Success checkmark: 400ms draw animation
- Quantity update: 200ms number counter
- Stock bar fill: 300ms ease-out transition

---

## Component Library Reference

Based on research, implement using:

1. **Cards**: Rounded 8px, soft shadow, warm tint
2. **Buttons**: 8px radius, 44px min height
3. **Progress Bars**: 8px height, rounded, with threshold marker
4. **Status Pills**: Rounded full, 24px height, icon + text
5. **Input Fields**: 8px radius, 48px height
6. **Bottom Sheet**: Rounded top 16px, drag handle
7. **Toast Notifications**: Rounded 8px, icon + text

---

## UX Metrics to Track

| Metric | Target | Measurement |
|--------|--------|-------------|
| Time to first item | < 3 min | Signup to first item saved |
| Stock update speed | < 10 sec | Open app to quantity updated |
| Alert resolution rate | > 70% | Alerts acted on vs ignored |
| Update completion | > 90% | Started updates that complete |
| Feature discovery | > 50% | Users who use snooze in first week |

---

Status: READY_FOR_REVIEW
