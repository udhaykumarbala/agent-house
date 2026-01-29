# Architecture Specification: SubTrack

## Executive Summary

This document defines the technical architecture for SubTrack, a privacy-first, mobile-first subscription management web application. Based on comprehensive research from all team members, this spec establishes the system design, tech stack, data models, and implementation patterns.

**Architecture Philosophy**: Simple, secure, and scalable - in that order.

---

## 1. Architecture Decision Records (ADRs)

### ADR-001: Fullstack Template Selection

**Decision**: Use `fullstack` template (Next.js 14 + Go API + PostgreSQL)

**Context**: The app requires user authentication, persistent data storage, CRUD operations, and calculations that cannot be achieved with static sites.

**Alternatives Considered**:
- `nextjs-frontend` with localStorage: Rejected - no cross-device sync, data loss risk
- `static-enhanced`: Rejected - insufficient for user accounts and persistent data
- Firebase/Supabase BaaS: Rejected - less control, vendor lock-in

**Consequences**:
- More complex deployment (2 services)
- Clear separation of concerns
- Full control over data and security
- Scalable architecture foundation

---

### ADR-002: Session-Based Authentication

**Decision**: Use cookie-based sessions for MVP (not JWT)

**Context**: Need secure authentication for personal financial data.

**Rationale**:
- Simpler to implement correctly
- Automatic CSRF protection with `SameSite=Strict`
- Easy revocation (delete session record)
- HttpOnly cookies prevent XSS token theft
- Sufficient for single-domain web app

**Future Path**: Add JWT support if native mobile app or third-party API access needed.

---

### ADR-003: PostgreSQL Database

**Decision**: PostgreSQL 16 for data persistence

**Context**: Need reliable storage for subscription and user data.

**Rationale**:
- ACID compliance for financial data integrity
- JSON support for flexible preferences storage
- Excellent indexing for query performance
- Mature tooling and ecosystem
- Easy to host (Railway, Supabase, RDS)

---

### ADR-004: Go for Backend API

**Decision**: Go with Chi router for REST API

**Context**: Need a performant, maintainable backend.

**Rationale**:
- Excellent performance with minimal resources
- Strong typing catches errors at compile time
- Simple deployment (single binary)
- Chi is idiomatic and middleware-friendly
- Good match for team expertise

---

### ADR-005: State Management Strategy

**Decision**: Zustand for client state, React Query for server state

**Context**: Need predictable state management without boilerplate.

**Rationale**:
- **Zustand**: Minimal API, no boilerplate, TypeScript-native
- **React Query**: Caching, background refetch, offline support built-in
- Clear separation: UI state (Zustand) vs server data (React Query)

---

## 2. System Architecture

### High-Level Component Diagram

```
                                 INTERNET
                                    |
                    +---------------+---------------+
                    |                               |
            +-------v-------+               +-------v-------+
            |    Vercel     |               |   Railway     |
            |   (Frontend)  |               |   (Backend)   |
            +-------+-------+               +-------+-------+
                    |                               |
            +-------v-------+               +-------v-------+
            |   Next.js 14  |    REST API   |    Go API     |
            |   App Router  | ------------> |  Chi Router   |
            |   TypeScript  |   /api/v1/*   |               |
            |   Tailwind    |               |               |
            |   shadcn/ui   |               +-------+-------+
            +---------------+                       |
                                            +-------v-------+
                                            |  PostgreSQL   |
                                            |   Database    |
                                            +---------------+
```

### Component Responsibilities

| Component | Technology | Responsibility |
|-----------|------------|----------------|
| **Frontend** | Next.js 14 | UI rendering, routing, PWA shell, client state |
| **API Gateway** | Chi Router | Request routing, middleware, rate limiting |
| **Auth Service** | Go + Argon2 | User registration, login, session management |
| **Subscription Service** | Go | CRUD operations, business logic, validation |
| **Analytics Service** | Go | Aggregation, calculations, reporting |
| **Database** | PostgreSQL | Data persistence, transactions, indexing |
| **Service Worker** | Workbox | Offline caching, push notifications |

---

## 3. Data Model

### Entity Relationship Diagram

```
+------------------+          +----------------------+
|      users       |          |    subscriptions     |
+------------------+          +----------------------+
| id (uuid) PK     |<---------| id (uuid) PK         |
| email (unique)   |    1:N   | user_id (uuid) FK    |
| password_hash    |          | name                 |
| email_verified   |          | price                |
| preferences      |          | currency             |
| created_at       |          | billing_cycle        |
| updated_at       |          | next_billing_date    |
+------------------+          | category             |
        |                     | icon                 |
        |                     | color                |
        |                     | is_active            |
        | 1:N                 | reminder_days        |
        v                     | created_at           |
+------------------+          | updated_at           |
|    sessions      |          +----------------------+
+------------------+
| id (uuid) PK     |
| user_id (uuid)   |
| token_hash       |
| expires_at       |
| ip_address       |
| user_agent       |
| created_at       |
+------------------+
```

### Table Schemas

#### users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    preferences JSONB DEFAULT '{
        "currency": "USD",
        "theme": "system",
        "reminderDays": 3,
        "emailNotifications": true
    }',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

#### sessions
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sessions_token ON sessions(token_hash);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expiry ON sessions(expires_at) WHERE expires_at > NOW();
```

#### subscriptions
```sql
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(3) DEFAULT 'USD',
    billing_cycle VARCHAR(20) NOT NULL
        CHECK (billing_cycle IN ('weekly', 'monthly', 'yearly', 'custom')),
    custom_days INTEGER CHECK (custom_days > 0),
    next_billing_date DATE NOT NULL,
    category VARCHAR(50) DEFAULT 'other'
        CHECK (category IN ('entertainment', 'productivity', 'lifestyle',
               'utilities', 'news', 'education', 'finance', 'health', 'other')),
    icon VARCHAR(50),
    color VARCHAR(7),
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    reminder_days INTEGER DEFAULT 3 CHECK (reminder_days >= 0 AND reminder_days <= 30),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_next_billing ON subscriptions(next_billing_date)
    WHERE is_active = TRUE;
CREATE INDEX idx_subscriptions_user_active ON subscriptions(user_id, is_active);
```

### Type Definitions (TypeScript)

```typescript
// types/subscription.ts
export type BillingCycle = 'weekly' | 'monthly' | 'yearly' | 'custom';

export type Category =
  | 'entertainment'
  | 'productivity'
  | 'lifestyle'
  | 'utilities'
  | 'news'
  | 'education'
  | 'finance'
  | 'health'
  | 'other';

export type Currency = 'USD' | 'EUR' | 'GBP' | 'CAD' | 'AUD' | 'INR';

export interface Subscription {
  id: string;
  name: string;
  description?: string;
  price: number;
  currency: Currency;
  billingCycle: BillingCycle;
  customDays?: number;
  nextBillingDate: string; // ISO date
  category: Category;
  icon?: string;
  color?: string;
  notes?: string;
  isActive: boolean;
  reminderDays: number;
  createdAt: string;
  updatedAt: string;
}

export interface SubscriptionInput {
  name: string;
  description?: string;
  price: number;
  currency?: Currency;
  billingCycle: BillingCycle;
  customDays?: number;
  nextBillingDate: string;
  category?: Category;
  icon?: string;
  color?: string;
  notes?: string;
  reminderDays?: number;
}
```

---

## 4. API Specification

### Base URL
- Development: `http://localhost:8080/api/v1`
- Production: `https://api.subtrack.app/v1`

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Create new account |
| POST | `/auth/login` | Login, create session |
| POST | `/auth/logout` | Invalidate session |
| POST | `/auth/forgot-password` | Request password reset |
| POST | `/auth/reset-password` | Reset with token |
| GET | `/auth/me` | Get current user |

### Subscription Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/subscriptions` | List all subscriptions |
| POST | `/subscriptions` | Create subscription |
| GET | `/subscriptions/:id` | Get subscription |
| PUT | `/subscriptions/:id` | Update subscription |
| DELETE | `/subscriptions/:id` | Delete subscription |

### Analytics Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/analytics/summary` | Total spend, by category |
| GET | `/analytics/upcoming` | Next 30 days renewals |
| GET | `/analytics/history` | Spending over time |

### User Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/user/preferences` | Update preferences |
| GET | `/user/export` | Export all data |
| DELETE | `/user/account` | Delete account |

### Request/Response Standards

#### Success Response
```json
{
  "data": { ... },
  "meta": {
    "total": 42,
    "page": 1,
    "perPage": 20
  }
}
```

#### Error Response
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human readable message",
    "details": [
      { "field": "price", "message": "Must be positive" }
    ]
  }
}
```

#### Error Codes
| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Invalid input data |
| `UNAUTHORIZED` | 401 | Missing or invalid session |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `RATE_LIMITED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

---

## 5. Project Structure

### Frontend (Next.js 14)

```
frontend/
├── app/
│   ├── (auth)/
│   │   ├── login/page.tsx
│   │   ├── register/page.tsx
│   │   ├── forgot-password/page.tsx
│   │   └── layout.tsx
│   ├── (dashboard)/
│   │   ├── page.tsx                    # Dashboard home
│   │   ├── subscriptions/
│   │   │   ├── page.tsx                # List view
│   │   │   ├── new/page.tsx            # Add new
│   │   │   └── [id]/page.tsx           # Edit/detail
│   │   ├── analytics/page.tsx
│   │   ├── settings/page.tsx
│   │   └── layout.tsx
│   ├── layout.tsx
│   ├── globals.css
│   └── manifest.ts                     # PWA manifest
├── components/
│   ├── ui/                             # shadcn/ui primitives
│   ├── subscription/
│   │   ├── subscription-card.tsx
│   │   ├── subscription-form.tsx
│   │   ├── subscription-list.tsx
│   │   └── category-filter.tsx
│   ├── dashboard/
│   │   ├── spend-summary.tsx
│   │   ├── upcoming-renewals.tsx
│   │   └── category-breakdown.tsx
│   └── layout/
│       ├── header.tsx
│       ├── bottom-nav.tsx
│       └── mobile-drawer.tsx
├── lib/
│   ├── api.ts                          # API client
│   ├── auth.ts                         # Auth utilities
│   ├── utils.ts                        # Helpers
│   └── constants.ts                    # App constants
├── hooks/
│   ├── use-subscriptions.ts
│   ├── use-analytics.ts
│   └── use-auth.ts
├── stores/
│   └── app-store.ts                    # Zustand store
├── types/
│   ├── subscription.ts
│   ├── user.ts
│   └── api.ts
└── public/
    ├── icons/                          # PWA icons
    └── sw.js                           # Service worker
```

### Backend (Go)

```
backend/
├── cmd/
│   └── api/
│       └── main.go                     # Entry point
├── internal/
│   ├── config/
│   │   └── config.go                   # Environment config
│   ├── handler/
│   │   ├── auth.go
│   │   ├── subscription.go
│   │   ├── analytics.go
│   │   ├── user.go
│   │   └── middleware.go
│   ├── service/
│   │   ├── auth.go
│   │   ├── subscription.go
│   │   └── analytics.go
│   ├── repository/
│   │   ├── user.go
│   │   ├── session.go
│   │   └── subscription.go
│   ├── model/
│   │   ├── user.go
│   │   ├── session.go
│   │   └── subscription.go
│   ├── database/
│   │   ├── postgres.go
│   │   └── migrations/
│   │       ├── 001_init.up.sql
│   │       └── 001_init.down.sql
│   └── validator/
│       └── validator.go
├── pkg/
│   ├── password/
│   │   └── argon2.go
│   └── response/
│       └── json.go
├── Dockerfile
├── go.mod
└── Makefile
```

---

## 6. Security Architecture

### Authentication Flow

```
[Client]                    [API]                      [Database]
   |                          |                            |
   |  POST /auth/login        |                            |
   |  {email, password}       |                            |
   | -----------------------> |                            |
   |                          |  SELECT user BY email      |
   |                          | -------------------------> |
   |                          |                            |
   |                          |  Verify password (Argon2)  |
   |                          |                            |
   |                          |  INSERT session            |
   |                          | -------------------------> |
   |                          |                            |
   |  Set-Cookie: session=xxx |                            |
   |  HttpOnly; Secure;       |                            |
   |  SameSite=Strict         |                            |
   | <----------------------- |                            |
```

### Security Controls

| Control | Implementation |
|---------|----------------|
| **Password Hashing** | Argon2id with 64MB memory, 3 iterations |
| **Session Tokens** | 32 bytes crypto random, SHA-256 hashed for storage |
| **Cookie Settings** | HttpOnly, Secure, SameSite=Strict, 7-day expiry |
| **Rate Limiting** | 10 req/min on auth, 100 req/min general |
| **Input Validation** | Server-side validation on all endpoints |
| **CORS** | Strict origin allowlist |
| **CSP** | Restrictive Content-Security-Policy headers |

### Security Headers

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

---

## 7. Performance Specifications

### Frontend Targets

| Metric | Target | Strategy |
|--------|--------|----------|
| **LCP** | < 2.5s | SSR critical content, optimize fonts |
| **FID** | < 100ms | Code splitting, lazy load routes |
| **CLS** | < 0.1 | Reserve space for dynamic content |
| **Bundle Size** | < 150KB (gzipped) | Tree shaking, no heavy deps |

### Backend Targets

| Metric | Target | Strategy |
|--------|--------|----------|
| **API Response** | < 100ms p95 | Connection pooling, optimized queries |
| **Database Query** | < 50ms p95 | Proper indexing, prepared statements |
| **Cold Start** | < 500ms | Single binary, minimal dependencies |

### Caching Strategy

| Resource | Cache Strategy | TTL |
|----------|----------------|-----|
| Static assets | Cache-first | 1 year (immutable) |
| App shell | Cache-first | 1 day |
| API: subscriptions | Network-first | Cache fallback |
| API: analytics | Network-first | 5 min stale-while-revalidate |

---

## 8. Offline Support

### Service Worker Strategy

```javascript
// Cache strategies by route
const routes = {
  // App shell - always cache
  '/': 'CacheFirst',
  '/dashboard': 'CacheFirst',
  '/subscriptions': 'CacheFirst',

  // API - network first with cache fallback
  '/api/v1/subscriptions': 'NetworkFirst',
  '/api/v1/analytics/*': 'StaleWhileRevalidate',

  // Static assets - cache first
  '/_next/static/*': 'CacheFirst',
  '/icons/*': 'CacheFirst',
};
```

### Offline Data Model

```typescript
// IndexedDB schema for offline
interface OfflineStore {
  subscriptions: Subscription[];
  pendingMutations: Array<{
    id: string;
    type: 'create' | 'update' | 'delete';
    data: Partial<Subscription>;
    timestamp: number;
  }>;
  lastSyncAt: number;
  user: { id: string; email: string };
}
```

### Sync Strategy

1. **Read**: Serve from cache, refresh in background
2. **Write**: Queue mutation, sync when online
3. **Conflict**: Last-write-wins with timestamp
4. **Notification**: Toast when back online with pending changes

---

## 9. Development Environment

### Prerequisites

- Node.js 20+
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 16 (via Docker)

### Local Setup Commands

```bash
# Clone and setup
git clone <repo>
cd subscription-manager

# Start all services
docker-compose up -d

# Frontend development
cd frontend
npm install
npm run dev  # http://localhost:3000

# Backend development
cd backend
go mod download
make run     # http://localhost:8080

# Database migrations
make migrate-up
```

### Environment Variables

```bash
# Backend (.env)
DATABASE_URL=postgres://user:pass@localhost:5432/subtrack
SESSION_SECRET=<32-byte-random-hex>
CORS_ORIGINS=http://localhost:3000
PORT=8080

# Frontend (.env.local)
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## 10. Deployment Strategy

### MVP Deployment

| Component | Platform | Rationale |
|-----------|----------|-----------|
| **Frontend** | Vercel | Next.js native, edge network, preview deployments |
| **Backend** | Railway | Simple Go deployment, PostgreSQL included |
| **Database** | Railway PostgreSQL | Managed, automatic backups |

### Deployment Pipeline

```
[Push to main]
       |
       v
[GitHub Actions]
   |         |
   v         v
[Build]   [Test]
   |         |
   +----+----+
        |
        v
   [Deploy]
   /       \
  v         v
Vercel   Railway
```

### Health Checks

- **Frontend**: Vercel automatic
- **Backend**: `GET /health` returns `{"status": "ok"}`
- **Database**: Connection pool health check

---

## 11. Implementation Phases

### Phase 1: Foundation (MVP Core)
1. Database schema and migrations
2. User authentication (register, login, logout)
3. Subscription CRUD API
4. Dashboard page with total spend
5. Subscription list with add/edit/delete
6. Mobile-responsive UI
7. Basic security headers

### Phase 2: Enhanced Experience
1. Category filtering and sorting
2. Analytics dashboard (charts, trends)
3. Dark mode support
4. PWA with offline read
5. Push notification setup
6. Data export (CSV/JSON)

### Phase 3: Polish & Scale
1. Offline write with sync
2. Email reminders before renewal
3. Subscription templates (popular services)
4. Multi-currency support
5. Performance optimization
6. Error tracking (Sentry)

---

## 12. Team Delegation

### Senior Developer
- Backend API implementation
- Database setup and migrations
- Authentication system
- Input validation middleware
- API testing

### UI Developer
- Component library (shadcn/ui customization)
- Design system implementation
- Dashboard components
- Subscription card/form components
- Responsive layout

### Security Review
- Authentication flow audit
- Rate limiting implementation
- Security headers configuration
- Input sanitization review
- Penetration testing coordination

---

## 13. Technical Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Session scaling limits | Low | Medium | Redis session store if needed |
| Offline sync conflicts | Medium | Low | Last-write-wins, user notification |
| Mobile performance | Medium | Medium | Performance budgets, testing on low-end devices |
| Database connection exhaustion | Low | High | Connection pooling, monitoring |
| Third-party API changes | Low | Low | No external dependencies for core features |

---

## 14. Questions Resolved

Based on research, these decisions are made:

1. **Social login in MVP?** No - defer to Phase 3 for simplicity
2. **PWA offline in MVP?** Offline read in Phase 2, offline write in Phase 3
3. **Email verification?** Recommended for Phase 1, but soft-required (allow usage)
4. **Admin dashboard?** No - not needed for MVP

---

## Summary

This architecture provides:

- **Clear separation**: Frontend (UI/UX) and Backend (Logic/Data) are independent
- **Type safety**: TypeScript + Go provide compile-time guarantees
- **Security-first**: Session auth, input validation, security headers from day 1
- **Mobile-first**: PWA capabilities, responsive design, offline support
- **Simple deployment**: Two services, managed database, automated CI/CD

The fullstack template is the right choice. We have the flexibility to add features while keeping the core simple and maintainable.

---

Status: PHASE_COMPLETE: planning
