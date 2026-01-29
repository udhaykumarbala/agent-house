# UI Flow Documentation - Reorder Alert System

## Overview

The Reorder Alert System is a mobile-first web application for tracking inventory stock levels and receiving alerts when quantities fall below configured thresholds.

---

## Screen Flow Diagram

```
┌─────────────────┐
│   Auth Screen   │
│  (Login/Signup) │
└────────┬────────┘
         │ Login Success
         ▼
┌─────────────────┐     ┌─────────────────┐
│    Dashboard    │◄───►│   Items List    │
│  (Hero Metrics) │     │ (All Inventory) │
└────────┬────────┘     └────────┬────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐     ┌─────────────────┐
│   Alerts List   │     │  Item Detail    │
│ (Active Alerts) │     │ (Stock + Update)│
└─────────────────┘     └────────┬────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │  Add/Edit Form  │
                        │  (Item CRUD)    │
                        └─────────────────┘

┌─────────────────┐
│    Settings     │
│ (Preferences)   │
└─────────────────┘

┌─────────────────┐
│  Quick Update   │
│ (Bottom Sheet)  │
└─────────────────┘
```

---

## Navigation Structure

### Bottom Navigation (Mobile)
- **Home** → Dashboard
- **Items** → Items List
- **Alerts** → Alerts List (with badge count)
- **Settings** → Settings Page

### Floating Action Button (FAB)
- Always visible on Dashboard and Items pages
- Tapping opens Add Item form

---

## Screen Descriptions

### 1. Auth Screen
**Purpose**: User authentication (login/signup)

**Elements**:
- Logo and tagline
- Email input field
- Password input field (with visibility toggle)
- Login/Signup button
- Toggle between login and signup modes
- "Forgot password" link

**User Actions**:
- Enter credentials and submit
- Switch between login/signup
- Toggle password visibility

---

### 2. Dashboard
**Purpose**: At-a-glance overview of inventory health

**Elements**:
- Hero metric card (items needing attention count)
- "Needs Attention" section with top 3 critical items
- Stock overview stats (Total, Alerts, OK)
- FAB for adding new item

**States**:
- Empty: First-time user CTA
- All healthy: Success message
- Alerts present: Priority-sorted alert cards

**User Actions**:
- Tap item card → Item Detail
- Tap "View All" → Alerts page
- Tap FAB → Add Item form

---

### 3. Items List
**Purpose**: Browse and manage all inventory items

**Elements**:
- Filter chips (All, Critical, Low, OK)
- Scrollable list of item cards
- Each card shows: name, quantity, threshold, stock bar, status badge
- FAB for adding new item

**User Actions**:
- Tap filter chip → Filter items by status
- Tap item card → Item Detail
- Tap FAB → Add Item form

---

### 4. Item Detail
**Purpose**: View full item info and update quantity

**Elements**:
- Back button
- Edit/Delete action buttons
- Item emoji/photo
- Large stock level bar with threshold marker
- Current quantity (large number)
- Status badge
- Quantity stepper (+/- buttons)
- Mode toggle (Set to / Add / Remove)
- Update button
- Item metadata (description, category, last updated)
- "Mark as Ordered" button (for low stock items)

**User Actions**:
- Tap +/- to adjust quantity
- Select update mode
- Tap Update to save
- Tap Edit → Edit Item form
- Tap Delete → Confirm dialog
- Tap Mark as Ordered → Toast confirmation

---

### 5. Add/Edit Item Form
**Purpose**: Create new or edit existing item

**Elements**:
- Cancel button (X)
- Save button
- Emoji picker (optional)
- Name input (required)
- Quantity input (required)
- Threshold input (required, with helper text)
- Description textarea (optional)
- Category input (optional)
- Delete button (edit mode only)

**Validation**:
- Name, quantity, threshold are required
- Quantity and threshold must be >= 0

**User Actions**:
- Fill form fields
- Tap emoji area to randomize emoji
- Tap Save → Create/update item
- Tap Cancel → Return without saving
- Tap Delete → Confirm dialog

---

### 6. Alerts List
**Purpose**: Action-oriented view of items needing attention

**Elements**:
- Tab filters (Active, Resolved)
- Date grouping (Today, Yesterday, Earlier)
- Alert cards with:
  - Status badge (CRITICAL / LOW STOCK)
  - Item name
  - Stock info (current / threshold)
  - Action buttons (Update, Snooze, Order)

**User Actions**:
- Tap Update → Quick Update sheet
- Tap Snooze → Temporarily dismiss alert (1 hour)
- Tap Order → Mark as ordered (moves to Resolved tab)
- Switch to Resolved tab → View previously resolved alerts

---

### 7. Quick Update (Bottom Sheet)
**Purpose**: Fast quantity update without leaving current page

**Elements**:
- Drag handle
- Item name
- Current quantity display
- Quantity stepper
- Mode toggle (Set to / Add / Remove)
- Preset buttons (+10, +25, +50, +100)
- Update button

**User Actions**:
- Adjust quantity with stepper or presets
- Select update mode
- Tap Update → Save and close
- Tap overlay → Close without saving

---

### 8. Settings
**Purpose**: Configure notifications and account

**Sections**:
- **Notifications**:
  - Email alerts toggle
  - Daily digest toggle
- **Account**:
  - Email display
  - Change password link
- **Data**:
  - Export CSV
  - Import CSV (coming soon)
- Logout button
- Version info

**User Actions**:
- Toggle notification settings
- Tap Export → Download CSV
- Tap Logout → Confirm dialog

---

## User Flows

### Flow 1: First-Time User Onboarding
```
1. User opens app
2. See Auth screen → Create account
3. Enter email/password → Submit
4. See empty Dashboard with CTA
5. Tap "Add Your First Item"
6. Fill form: Name, Quantity, Threshold
7. Tap Save
8. Return to Dashboard with item visible
```

### Flow 2: Respond to Low Stock Alert
```
1. Open app (or tap notification)
2. Dashboard shows "3 items need attention"
3. Tap on alert item card
4. View Item Detail with low stock warning
5. Option A: Update quantity (+10 units)
   - Tap + button
   - Tap Update
6. Option B: Mark as ordered
   - Tap "Mark as Ordered"
   - See confirmation toast
```

### Flow 3: Quick Stock Update After Shipment
```
1. Open app from warehouse floor
2. Navigate to Items
3. Find item in list
4. Tap item card → Item Detail
5. Tap "Add" mode
6. Enter received quantity
7. Tap Update
8. See "Quantity updated" toast
9. Stock bar animates to new level
```

### Flow 4: Export Inventory Data
```
1. Navigate to Settings
2. Scroll to Data section
3. Tap "Export Data"
4. CSV file downloads
5. See "Data exported" toast
```

---

## State Management

### Item Status Logic
```javascript
function getItemStatus(item) {
  if (item.quantity <= item.threshold * 0.5) return 'critical';
  if (item.quantity <= item.threshold) return 'warning';
  return 'success';
}
```

### Alert Badge Count
The alert badge on the Alerts nav item shows the count of items with status 'critical' or 'warning'.

---

## Responsive Behavior

### Mobile (< 768px)
- Single column layout
- Bottom navigation visible
- FAB visible
- Bottom sheet modals

### Tablet (768px - 1023px)
- Two-column item grid
- Centered content (max-width: 600px)

### Desktop (1024px+)
- Three-column layout possible
- Bottom navigation hidden
- FAB positioned at bottom-right
- Dialog modals instead of bottom sheets

---

## Accessibility

- All touch targets minimum 44x44px
- Color + icon + text for status indicators
- Focus-visible states for keyboard navigation
- Screen reader announcements for quantity updates
- Respects prefers-reduced-motion

---

## Error States

### Network Error
- Show retry message
- Allow offline viewing of cached data

### Form Validation
- Inline error messages below fields
- Red border on invalid inputs
- Prevent submission until valid

### Empty States
- Friendly messaging
- Clear CTA to fix the situation
- Illustration/icon for visual interest
