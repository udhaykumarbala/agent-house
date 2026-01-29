# Product Research: Inventory Reorder Alert System

## Executive Summary

An inventory reorder alert system that tracks stock levels of commonly procured items and triggers alerts when quantities fall below configurable thresholds. Designed for small-to-medium businesses that lack enterprise ERP solutions but need reliable stock management.

---

## Competitor Analysis

| Competitor | Strengths | Weaknesses | Price |
|------------|-----------|------------|-------|
| **Sortly** | Visual inventory with photos, barcode scanning, mobile-first, QR labels | Expensive for small teams, limited integration options, alert system basic | $49-149/mo |
| **inFlow Inventory** | Full-featured (PO, sales orders, reports), multi-location, B2B focus | Complex setup, overkill for simple tracking, desktop-heavy UI | $89-439/mo |
| **Zoho Inventory** | Tight Zoho ecosystem, multichannel selling, good automations | Requires Zoho buy-in, steep learning curve, slow mobile app | $59-329/mo |
| **Stockpile** | Free tier available, simple interface, good for beginners | Limited alert customization, no API, basic reporting | Free-$25/mo |

### Key Observations

1. **Enterprise tools dominate** - Most solutions target mid-to-large businesses with complex needs
2. **Alert systems are secondary** - Reorder alerts are buried features, not primary value props
3. **Mobile experience varies** - Many are desktop-first with clunky mobile apps
4. **Integration complexity** - Small businesses struggle to connect tools

---

## Our Differentiation

**Alert-First Design** - While competitors treat reorder alerts as a checkbox feature, we make it the core experience. Simple setup, smart notifications, actionable alerts.

**Key Differentiators:**
1. **5-minute setup** - No complex onboarding, just add items and thresholds
2. **Multi-channel alerts** - Email, SMS, Slack, Teams, webhook support
3. **Smart thresholds** - Suggest reorder points based on consumption patterns
4. **One-click reorder** - Generate purchase orders or notify suppliers directly
5. **Lightweight** - Solves one problem well, no feature bloat

---

## Target User Persona

**Name**: "Operations Manager Omar"

**Demographics**:
- Age: 28-45
- Role: Office manager, ops coordinator, small business owner
- Company size: 5-50 employees
- Industries: Dental/medical offices, restaurants, retail shops, warehouses, manufacturing workshops

**Pain Points**:
- Runs out of critical supplies unexpectedly (coffee, printer paper, cleaning supplies, raw materials)
- Currently tracks inventory in spreadsheets that nobody updates
- Gets blamed when items run out
- No time to check stock levels daily
- Over-orders to compensate, tying up cash and storage space

**Current Solutions**:
- Excel/Google Sheets (manual, outdated)
- Physical clipboard counts (time-consuming)
- Memory/"we'll know when we're low" (unreliable)
- Enterprise ERP (too expensive, too complex)

**Why They'd Switch**:
- Set it and forget it - get notified only when action needed
- Avoid the embarrassment of stockouts
- Stop wasting time on manual counts
- Simple enough that the whole team will actually use it

**Context of Use**:
- Quick mobile check during stock receipt
- Desktop dashboard for planning
- Notifications push to phone/email/Slack

---

## Secondary Persona

**Name**: "Warehouse Worker Wendy"

**Demographics**:
- Age: 22-55
- Role: Stock clerk, warehouse staff, inventory associate
- Responsibility: Updates stock levels, receives shipments

**Pain Points**:
- Current system is too complex/slow
- Doesn't get timely info on what to reorder
- Multiple systems to log into

**Why They'd Switch**:
- Fast mobile updates (scan or tap)
- Clear visibility into what's running low
- Less paperwork and data entry

---

## Market Positioning

**One-Liner Pitch**: "Never run out of supplies again - get smart alerts before stock runs low."

**Alternative Pitches**:
- "The inventory alert system that's actually easy to use"
- "Set your thresholds. Get alerts. Reorder. Done."

**Price Point**: **Freemium**
- **Free tier**: Up to 50 items, 1 user, email alerts only
- **Pro tier ($19/mo)**: Unlimited items, 5 users, multi-channel alerts, consumption analytics
- **Team tier ($49/mo)**: Unlimited users, API access, multiple locations, integrations

**Competitive Positioning**:
```
                        Complex
                           |
                    inFlow | Zoho
                           |
        Expensive -------- + -------- Affordable
                           |
                    Sortly | [US]
                           |
                        Simple
```

We occupy the **simple + affordable** quadrant with a laser focus on alerts.

---

## Feature Priorities (Initial Research)

Based on competitor analysis and user needs:

### Must-Haves (MVP)
1. Item catalog with current quantity tracking
2. Configurable reorder thresholds per item
3. Email/notification alerts when threshold breached
4. Simple mobile-friendly quantity updates
5. Dashboard showing items at/below threshold

### Should-Haves (v1.1)
1. Barcode/QR scanning for quick updates
2. Slack/Teams integration
3. Consumption trend tracking
4. Bulk import from CSV/Excel
5. Purchase order generation

### Nice-to-Haves (Future)
1. Supplier management
2. Multi-location support
3. Predictive reorder suggestions (ML-based)
4. Integrations (QuickBooks, Shopify, etc.)

---

## Risks & Considerations

| Risk | Mitigation |
|------|------------|
| Users don't update quantities | Make updates dead-simple (1-tap, barcode scan) |
| Alert fatigue | Smart grouping, digest options, snooze capability |
| Data entry burden | Bulk import, barcode scan, suggest from history |
| Competition undercuts | Focus on UX and alert intelligence, not price war |

---

## Success Metrics (North Stars)

1. **Stockout Prevention Rate** - % of items that get reordered before running out
2. **Time to First Alert** - How quickly users set up and receive first alert
3. **Active Usage** - Weekly quantity updates per account
4. **Alert Action Rate** - % of alerts that result in reorder action

---

## Technical Considerations

- **Real-time updates** - WebSocket for live dashboard
- **Notification reliability** - Queue-based alert system with retry
- **Mobile-first** - PWA or native apps for stock updates
- **API-first** - Enable integrations from day one

---

Status: READY_FOR_PLANNING
