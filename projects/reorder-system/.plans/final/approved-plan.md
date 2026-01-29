# APPROVED PLAN: Inventory Reorder Alert System

Reviewed by: CEO
Date: 2025-12-27
Status: APPROVED FOR DEVELOPMENT

---

## Summary

The Inventory Reorder Alert System is a mobile-first web application that helps small-to-medium businesses track stock levels and receive smart alerts when quantities fall below configured thresholds. The system differentiates from competitors through its alert-first design, 5-minute setup experience, and multi-channel notification support.

**Core Value Proposition**: "Never run out of supplies again - get smart alerts before stock runs low."

---

## Approved Specifications

| Document | Path | Status |
|----------|------|--------|
| Product Spec | `.plans/specs/product-spec.md` | Approved |
| UX Spec | `.plans/specs/ux-spec.md` | Approved |
| UI Spec | `.plans/specs/ui-spec.md` | Approved |
| Security Spec | `.plans/specs/security-spec.md` | Approved |
| Architecture Spec | `.plans/specs/architecture-spec.md` | Approved |

---

## Key Decisions

1. **Architecture**: Modular monolith pattern for MVP, with clear module boundaries enabling future service extraction if needed.

2. **Technology Stack**:
   - Frontend: Next.js 14 (App Router) + TypeScript
   - Backend: NestJS + Prisma ORM
   - Database: PostgreSQL 15 with Row-Level Security
   - Cache/Queue: Redis + BullMQ
   - Real-time: Socket.IO

3. **Multi-Tenancy**: Shared database with tenant column + PostgreSQL RLS for data isolation.

4. **Design System**: "Refined Industrial" theme with Deep Plum (#7C3AED) primary and warm cream (#FDF8F3) background - distinctive positioning vs. competitors.

5. **Security**: Argon2id password hashing, JWT with refresh tokens, SSRF prevention for webhooks, comprehensive audit logging.

6. **Notifications**: Queue-based pipeline with email in MVP, Slack/SMS/webhooks in later phases.

---

## MVP Feature Scope (P0)

| Feature | Description |
|---------|-------------|
| Item Catalog Management | Add/edit items with name, quantity, threshold, optional photo |
| Stock Level Tracking | Visual progress bars with threshold markers |
| Configurable Thresholds | Per-item reorder point configuration |
| Alert Dashboard | Hero metric + items needing attention list |
| Email Notifications | Automatic alerts when thresholds crossed |
| Quick Quantity Update | +/- buttons, < 3 taps to update |
| Mobile-Responsive UI | PWA, works on all screen sizes |
| User Authentication | Email/password, JWT sessions |

---

## Implementation Order

### Phase 1: Foundation
1. Project scaffolding (monorepo with Turborepo)
2. Database schema and Prisma setup
3. Authentication module (register, login, JWT)
4. Multi-tenant middleware and RLS policies

### Phase 2: Core Features
5. Items CRUD API and frontend
6. Stock level visualization (StockBar component)
7. Threshold checker service
8. Alert creation and dashboard

### Phase 3: Notifications
9. BullMQ notification queue
10. Email provider (SendGrid integration)
11. Real-time updates (WebSocket)

### Phase 4: Polish
12. Settings and preferences
13. Empty states and error handling
14. Mobile optimizations
15. Testing and security audit

---

## Design Tokens Summary

```css
/* Primary Colors */
--color-primary: #7C3AED;
--color-bg: #FDF8F3;
--color-surface: #FFFFFF;

/* Status Colors */
--color-critical: #B91C1C;
--color-warning: #D97706;
--color-success: #4D7C0F;

/* Typography */
font-family: 'Inter', sans-serif;

/* Spacing Base */
4px unit system
```

---

## Notes for Development Team

1. **Security First**: All P0 security controls from security-spec.md must be implemented before MVP launch, especially tenant isolation and input validation.

2. **Mobile First**: Build responsive layouts mobile-first, then enhance for tablet/desktop. All touch targets minimum 44x44px.

3. **Stock Update Speed**: The < 10 second stock update flow is critical. Optimize for this path with optimistic UI updates.

4. **Categories**: Include optional category field in Add Item form for MVP, but defer category management UI to v1.1.

5. **Mark as Ordered**: Include basic button in Item Detail screen for MVP to enable alert resolution. Full order tracking is P1.

6. **Accessibility**: Maintain WCAG AA compliance. Never use color alone for status - always include icons and text labels.

7. **Error Handling**: Use toast notifications for success/error feedback. Keep error messages generic per security requirements.

8. **Testing Requirements**:
   - Unit tests for all services (Vitest)
   - E2E tests for critical flows (Playwright)
   - Security test cases from security-spec.md

---

## Success Metrics to Track

| Metric | Target |
|--------|--------|
| Time to First Alert | < 10 minutes |
| Stockout Prevention Rate | > 85% |
| Alert Action Rate | > 70% |
| Stock Update Speed | < 10 seconds |
| Setup Completion Rate | > 50% (5+ items) |

---

## Out of Scope for MVP

- Purchase order creation/generation
- Supplier management
- Multi-location support
- Barcode scanning
- SMS notifications
- External integrations (Shopify, QuickBooks)
- Role-based permissions (all users have full access)
- Predictive/ML features
- Historical reporting/analytics

---

## Post-MVP Roadmap

### v1.1 - Enhanced Usability
- Barcode/QR scanning
- Category management
- Slack/Teams integration
- Alert digest mode
- Bulk CSV import
- Consumption trend tracking

### v1.2 - Integrations
- Webhook notifications
- Supplier management
- One-click reorder (PO generation)
- SMS notifications
- Two-factor authentication

---

## Review Sign-Off

| Role | Status | Date |
|------|--------|------|
| CEO | APPROVED | 2025-12-27 |
| Product Manager | Spec Complete | 2025-12-27 |
| UX Designer | Spec Complete | 2025-12-27 |
| UI Designer | Spec Complete | 2025-12-27 |
| Security Expert | Spec Complete | 2025-12-27 |
| Architect | Spec Complete | 2025-12-27 |

---

BEGIN DEVELOPMENT
