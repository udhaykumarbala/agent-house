# Product Specification: Inventory Reorder Alert System

Based on research in:
- `.plans/research/product-research.md` (Market & Competitor Analysis)
- `.plans/research/ux-research.md` (User Experience Patterns)
- `.plans/research/ui-research.md` (Visual Design Direction)

---

## Target User

**Primary Persona: "Operations Manager Omar"**
- Age: 28-45
- Role: Office manager, ops coordinator, small business owner
- Company size: 5-50 employees
- Industries: Dental/medical offices, restaurants, retail shops, warehouses, manufacturing workshops

**Secondary Persona: "Warehouse Worker Wendy"**
- Age: 22-55
- Role: Stock clerk, warehouse staff, inventory associate
- Needs: Fast mobile updates, clear visibility, minimal data entry

---

## Unique Value Proposition

**"Never run out of supplies again - get smart alerts before stock runs low."**

**Core Differentiators:**
1. Alert-first design (not a buried feature)
2. 5-minute setup (no complex onboarding)
3. Multi-channel notifications (email, SMS, Slack, Teams, webhooks)
4. One-tap reorder action
5. Mobile-first experience with desktop dashboard

---

## Features (Prioritized)

### P0 - Must Have (MVP)

- [ ] **Item Catalog Management**: Add/edit items with name, description, current quantity, reorder threshold, optional photo
  - *Why essential*: Core data model required for all functionality

- [ ] **Stock Level Tracking**: Display current quantity with visual progress bar showing threshold position
  - *Why essential*: Users must see at-a-glance status (research: Mint-style visual bars)

- [ ] **Configurable Thresholds**: Set per-item reorder point (minimum quantity trigger)
  - *Why essential*: Core alert trigger mechanism

- [ ] **Alert Dashboard**: Central "Needs Attention" view showing items at/below threshold
  - *Why essential*: Primary interaction point (research: Robinhood hero metric pattern)

- [ ] **Email Notifications**: Automatic email when item crosses threshold
  - *Why essential*: Primary alert delivery for MVP

- [ ] **Quick Quantity Update**: Simple +/- buttons or direct number entry to update stock
  - *Why essential*: Must be < 3 taps (research: warehouse workers need speed)

- [ ] **Mobile-Responsive UI**: Works on phone/tablet for on-floor updates
  - *Why essential*: Target users update stock while receiving shipments

- [ ] **User Authentication**: Secure login with email/password
  - *Why essential*: Multi-user accounts require identity

### P1 - Should Have (Post-MVP v1.1)

- [ ] **Barcode/QR Scanning**: Scan item to quickly update quantity
  - *Why valuable*: Reduces friction for high-volume updates

- [ ] **Category Organization**: Group items by type (Office Supplies, Raw Materials, etc.)
  - *Why valuable*: Alert grouping prevents fatigue (research: Slack pattern)

- [ ] **Slack/Teams Integration**: Push alerts to team channels
  - *Why valuable*: Meet users where they work

- [ ] **Alert Digest Mode**: Daily/weekly summary instead of individual alerts
  - *Why valuable*: Prevents notification fatigue

- [ ] **Consumption Trend Tracking**: Show usage rate and days-until-stockout
  - *Why valuable*: Proactive planning, smarter thresholds

- [ ] **Bulk CSV Import**: Upload existing inventory from spreadsheet
  - *Why valuable*: Critical for migration from Excel trackers

- [ ] **Mark as Ordered**: Track reorder status without external system
  - *Why valuable*: Closes the alert loop

- [ ] **Alert Snooze**: Temporarily dismiss alert (1 hour, 1 day, custom)
  - *Why valuable*: Prevents repeat alerts when action is pending

### P2 - Nice to Have (Future)

- [ ] **Purchase Order Generation**: Create PO from alert with one click
  - *Future consideration*: Bridges to ordering workflow

- [ ] **Supplier Management**: Link items to preferred vendors
  - *Future consideration*: Enables direct reorder actions

- [ ] **Multi-Location Support**: Track inventory across locations
  - *Future consideration*: Scale to larger organizations

- [ ] **Predictive Reorder Suggestions**: ML-based threshold recommendations
  - *Future consideration*: Smart automation based on patterns

- [ ] **API & Webhooks**: Enable external integrations
  - *Future consideration*: Connect to ERP, accounting systems

- [ ] **SMS Notifications**: Critical alerts via text message
  - *Future consideration*: Urgent notification channel

---

## User Stories

### MVP User Stories

1. **As an operations manager**, I want to add new items to my inventory catalog so that I can track their stock levels.

2. **As an operations manager**, I want to set a reorder threshold for each item so that I'm alerted when stock is low.

3. **As an operations manager**, I want to see a dashboard of all items needing attention so that I can prioritize reorders.

4. **As an operations manager**, I want to receive email alerts when an item falls below threshold so that I don't have to constantly check.

5. **As a warehouse worker**, I want to quickly update stock quantities after receiving a shipment so that levels are always accurate.

6. **As a warehouse worker**, I want to use my phone to update inventory so that I can do it while on the floor.

7. **As an operations manager**, I want to see visual indicators (green/yellow/red) for stock status so that I can scan inventory health at a glance.

8. **As a new user**, I want to add my first item in under 5 minutes so that I can see value immediately.

### Post-MVP User Stories

9. **As a warehouse worker**, I want to scan a barcode to find and update an item so that I don't have to search manually.

10. **As an operations manager**, I want to mark items as "ordered" so that I know which alerts are being addressed.

11. **As an operations manager**, I want to receive a daily digest of all low-stock items so that I'm not overwhelmed with individual notifications.

12. **As an operations manager**, I want to see consumption trends so that I can set smarter thresholds.

---

## Acceptance Criteria

### Feature: Item Catalog Management
- [ ] User can create item with name (required), quantity (required), threshold (required)
- [ ] User can optionally add description and photo
- [ ] User can edit existing items
- [ ] User can delete items (with confirmation)
- [ ] Items persist across sessions
- [ ] Duplicate item names show warning

### Feature: Stock Level Tracking
- [ ] Each item displays current quantity numerically
- [ ] Each item displays visual progress bar showing current vs. max capacity
- [ ] Threshold position is marked on progress bar
- [ ] Bar color changes: Green (>150% of threshold), Yellow (100-150% of threshold), Red (below threshold)
- [ ] Status badge shows "OK", "Low", or "Critical"

### Feature: Threshold Configuration
- [ ] Each item has a configurable numeric threshold
- [ ] Threshold must be >= 0
- [ ] Threshold can be updated at any time
- [ ] Changes to threshold immediately update item status

### Feature: Alert Dashboard
- [ ] Dashboard shows hero metric: "X items need attention"
- [ ] Items below threshold appear in "Needs Attention" section
- [ ] Items are sorted by urgency (lowest stock % first)
- [ ] Each alert shows: item name, current qty, threshold, days since triggered
- [ ] Dashboard refreshes automatically (or on pull-to-refresh on mobile)
- [ ] Empty state shows "All stock levels healthy" with positive messaging

### Feature: Email Notifications
- [ ] Email sent within 5 minutes of threshold breach
- [ ] Email contains: item name, current quantity, threshold, link to dashboard
- [ ] User can configure email frequency (immediate, daily digest)
- [ ] User can disable email notifications per item or globally
- [ ] Email is mobile-friendly

### Feature: Quick Quantity Update
- [ ] Update accessible in 2 taps from dashboard
- [ ] Supports +/- increment buttons
- [ ] Supports direct numeric input
- [ ] Supports "Set to" and "Add/Remove" modes
- [ ] Shows confirmation of update
- [ ] Updates reflect immediately in UI

### Feature: Mobile-Responsive UI
- [ ] All features usable on 320px+ width screens
- [ ] Touch targets minimum 44x44px
- [ ] Bottom navigation on mobile (Dashboard, Items, Alerts, Settings)
- [ ] Floating action button for quick add
- [ ] No horizontal scrolling required

### Feature: User Authentication
- [ ] User can register with email and password
- [ ] User can log in with email and password
- [ ] Password reset via email
- [ ] Sessions persist (no login required every visit)
- [ ] Logout clears session

---

## Success Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| **Time to First Alert** | Time from signup to receiving first threshold alert | < 10 minutes |
| **Stockout Prevention Rate** | % of items reordered before reaching zero | > 85% |
| **Alert Action Rate** | % of alerts that result in reorder or snooze action | > 70% |
| **Weekly Active Usage** | % of users updating stock at least once per week | > 60% |
| **Stock Update Speed** | Time to complete a quantity update | < 10 seconds |
| **Setup Completion Rate** | % of new users who add at least 5 items | > 50% |

---

## Out of Scope (for MVP)

- **Purchase order creation** - Users create POs in their existing systems
- **Supplier management** - No vendor database in v1
- **Multi-location inventory** - Single location only for MVP
- **Barcode scanning** - Manual search/select for MVP
- **SMS notifications** - Email only for MVP
- **Integrations** (QuickBooks, Shopify, etc.) - Standalone system
- **Role-based permissions** - All users have full access
- **Historical reporting** - No analytics dashboard
- **Predictive/ML features** - Manual thresholds only
- **Inventory counting workflows** - No audit/count features
- **Units of measure conversion** - Single unit per item

---

## Technical Boundaries (from Architecture Research)

- Web-based application (PWA for mobile experience)
- Real-time dashboard updates via WebSocket
- RESTful API for all operations
- Queue-based notification system with retry logic
- Mobile-first responsive design

---

## Design Direction (from UI Research)

**Theme: "Refined Industrial"**
- Background: Warm cream #FDF8F3
- Primary: Deep Plum #7C3AED (distinctive, not generic blue/green)
- Status colors: Crimson (critical), Amber (warning), Sage (healthy)
- Typography: Inter (excellent number legibility)
- Style: Rounded corners (8px), soft shadows, warm palette

---

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Users don't update quantities regularly | Make updates < 3 taps, add barcode scan in v1.1 |
| Alert fatigue from too many notifications | Daily digest option, snooze capability |
| Data entry burden during onboarding | Bulk CSV import in v1.1, minimal required fields |
| Users don't set appropriate thresholds | Provide suggestions based on similar items |
| Email delivery issues | Queue with retry, show in-app alert history |

---

## Open Questions for Team Review

1. Should we include a "low stock" prediction (days until stockout) in MVP or defer to v1.1?
2. Do we need multi-user support in MVP, or is single-user sufficient for launch?
3. Should categories be part of MVP for better alert grouping?
4. What's the default threshold suggestion when users don't know what to set?

---

Status: READY_FOR_REVIEW
