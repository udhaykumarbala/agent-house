# APPROVED PLAN: SubTrack - Subscription Manager

Reviewed by: CEO
Date: January 18, 2025
Status: APPROVED FOR DEVELOPMENT

---

## Summary

SubTrack is a privacy-first, mobile-first subscription management web application. It helps users track their recurring payments without requiring bank connections, differentiating it from competitors like Truebill/Rocket Money. The app focuses on simplicity, beautiful UX, and user trust.

**Core Value Proposition**: "See all your subscriptions in one beautiful place. No bank connection needed."

---

## Approved Specifications

| Spec | File | Status |
|------|------|--------|
| Product | `.plans/specs/product-spec.md` | Approved |
| UX | `.plans/specs/ux-spec.md` | Approved |
| UI | `.plans/specs/ui-spec.md` | Approved |
| Security | `.plans/specs/security-spec.md` | Approved |
| Architecture | `.plans/specs/architecture-spec.md` | Approved |

---

## Key Decisions

1. **Template**: Fullstack (Next.js 14 + Go API + PostgreSQL)
2. **Authentication**: Session-based with secure cookies (not JWT for MVP)
3. **Brand Color**: Coral (#FF7A5C) - warm, memorable, unique in fintech
4. **Typography**: Geist font for clean number rendering
5. **Bottom Navigation**: 3 tabs (Home, Upcoming, Settings)
6. **Offline Strategy**: PWA with cache-first for app shell, network-first for API data
7. **No Bank Linking**: Core privacy differentiator - all data entered manually

---

## MVP Feature Scope (P0)

- Dashboard with total monthly/yearly spend
- Add subscription manually (< 30 seconds target)
- Subscription list with cards
- Popular service templates (20+)
- Edit/Delete subscriptions
- Local storage + cloud sync via account
- Dark mode
- Mobile-responsive design

---

## Implementation Order

### Phase 1: Foundation
1. Database schema and migrations (PostgreSQL)
2. User authentication (register, login, logout)
3. Subscription CRUD API endpoints
4. Security middleware (rate limiting, validation, headers)

### Phase 2: Core UI
1. Dashboard page with spend summary
2. Subscription list component
3. Add subscription form/sheet
4. Subscription detail/edit view
5. Bottom navigation
6. Empty state

### Phase 3: Polish
1. Dark mode implementation
2. Service templates integration
3. Due soon indicators
4. Swipe-to-delete gesture
5. Toast notifications
6. Loading skeletons

---

## Development Team Assignments

### Senior Developer (senior_dev)
- Backend API (Go + Chi)
- Database setup and migrations
- Authentication system with Argon2id
- Input validation middleware
- Rate limiting
- Security headers
- API testing

### UI Developer (assigned from team)
- Component library setup (shadcn/ui customization)
- Design system tokens (colors, typography, spacing)
- Dashboard components
- Subscription card/form components
- Bottom navigation
- Responsive layout
- Dark mode toggle

---

## Critical Success Metrics

| Metric | Target |
|--------|--------|
| Time to add first subscription | < 30 seconds |
| Activation (3+ subs in first session) | 60% |
| Add flow completion rate | 80% |
| 30-day retention | 40% |
| Lighthouse Performance Score | 90+ |
| API response time (p95) | < 100ms |

---

## Security Requirements (MVP - Non-negotiable)

- [x] Argon2id password hashing
- [x] HttpOnly, Secure, SameSite=Strict cookies
- [x] Rate limiting on auth endpoints (10 req/min)
- [x] Parameterized queries (sqlc)
- [x] IDOR prevention (user_id in all queries)
- [x] Security headers (HSTS, CSP, X-Frame-Options)
- [x] Generic error messages

---

## Notes for Development Team

1. **Mobile-First**: All development should start with mobile viewport. Desktop is enhancement.

2. **Touch Targets**: Minimum 44x44px for all interactive elements.

3. **Performance Budget**: Initial JS bundle < 150KB gzipped.

4. **Color System**: Use CSS custom properties defined in UI spec for easy theming.

5. **Font Note**: Geist font is from fonts.vercel.com, not Google Fonts. Use @fontsource/geist for self-hosting.

6. **Templates Data**: Hardcode popular service templates as JSON for MVP. Can move to remote config later.

7. **Currency**: Schema supports multiple currencies, but MVP only implements USD. Currency selector disabled in UI.

8. **No Over-Engineering**: Keep it simple. Don't add features not in scope. We can iterate after launch.

---

## Out of Scope (Do Not Build in MVP)

- Bank/card linking
- Automatic subscription detection
- Bill negotiation services
- Social login (Google, Apple)
- Native mobile apps
- Multi-currency display
- Team/family accounts
- API for third parties
- Admin dashboard

---

## Launch Readiness Checklist

Before declaring MVP complete:

- [ ] All P0 features functional
- [ ] Security audit passed (see security-spec.md acceptance criteria)
- [ ] Lighthouse score > 90 on mobile
- [ ] Works offline (read mode)
- [ ] Dark mode functional
- [ ] 50+ subscriptions render smoothly
- [ ] Password reset flow works
- [ ] Data export works
- [ ] Privacy policy page exists

---

## Approved By

**CEO Agent**
January 18, 2025

The team has produced excellent, well-coordinated specifications. The vision is clear, the architecture is sound, and the security posture is appropriate for handling personal financial data.

Proceed with development.

---

BEGIN DEVELOPMENT
