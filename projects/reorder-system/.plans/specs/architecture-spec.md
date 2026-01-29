# Architecture Specification: Inventory Reorder Alert System

## Document Overview

| Field | Value |
|-------|-------|
| Version | 1.0 |
| Status | APPROVED |
| Author | Architect Agent |
| Based On | product-research.md, architecture-research.md, security-research.md, ux-research.md, ui-research.md |

---

## 1. Technical Analysis

### 1.1 Core Requirements

Based on research, the system must support:

| Requirement | Source | Technical Implication |
|-------------|--------|----------------------|
| 5-minute setup | Product Research | Minimal form fields, sensible defaults, no wizard |
| Multi-channel alerts | Product Research | Queue-based notification pipeline, provider abstraction |
| One-click reorder | Product Research | Optimistic UI, fast API responses |
| Multi-tenant isolation | Security Research | PostgreSQL RLS, tenant context middleware |
| Real-time dashboard | Architecture Research | WebSocket connections, Redis pub/sub |
| Mobile-first experience | UX Research | PWA, responsive design, touch-optimized |
| SSRF protection | Security Research | URL validation, IP blocklist for webhooks |

### 1.2 Non-Functional Requirements

| Metric | Target | Justification |
|--------|--------|---------------|
| API response time (p95) | < 200ms | UX research: quick update flow critical |
| Dashboard load | < 1s | Competitive parity |
| Notification delivery | < 30s | Alert urgency |
| Concurrent users | 1000+ per instance | Growth projection |
| Items per tenant | 100,000+ | Enterprise readiness |
| Uptime | 99.9% | Business-critical alerts |

### 1.3 Constraints

- **Budget**: Freemium model requires cost-efficient infrastructure
- **Team Size**: Assumed small team, favor simplicity over complexity
- **Timeline**: MVP focus, avoid over-engineering
- **Security**: GDPR compliance from day one

---

## 2. Architecture Decision

### 2.1 Pattern: Modular Monolith

**Decision**: Build as a modular monolith with clear module boundaries.

**Rationale**:
- Faster initial development vs. microservices
- Simple deployment and debugging
- Clear path to extract services if needed
- Matches product's "simple" philosophy

**Modules**:
```
┌─────────────────────────────────────────────────────────────────┐
│                     MODULAR MONOLITH                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐     │
│  │   Auth    │  │   Items   │  │  Alerts   │  │   Notify  │     │
│  │  Module   │  │  Module   │  │  Module   │  │  Module   │     │
│  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘     │
│        │              │              │              │            │
│        └──────────────┴──────────────┴──────────────┘            │
│                              │                                    │
│                    ┌─────────┴─────────┐                         │
│                    │   Shared Core     │                         │
│                    │  (DB, Cache, Bus) │                         │
│                    └───────────────────┘                         │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Communication Pattern: Hybrid

- **Synchronous**: REST API for CRUD operations
- **Asynchronous**: Event-driven for alert processing and notifications
- **Real-time**: WebSocket for dashboard updates

---

## 3. Technology Stack

### 3.1 Final Stack Selection

| Layer | Technology | Justification |
|-------|------------|---------------|
| **Frontend** | Next.js 14 (App Router) | SSR/SSG, TypeScript, React ecosystem |
| **Backend** | NestJS | Modular architecture, TypeScript, built-in validation |
| **Database** | PostgreSQL 15 | RLS for multi-tenancy, JSONB flexibility, reliability |
| **Cache/Queue** | Redis + BullMQ | TypeScript native, job queue for notifications |
| **ORM** | Prisma | Type-safe queries, migrations, schema-first |
| **Auth** | JWT + Refresh Tokens | Stateless, scalable, secure |
| **Email** | SendGrid | Reliable, good free tier |
| **SMS** | Twilio | Industry standard |
| **Real-time** | Socket.IO | Redis adapter for scaling |
| **Validation** | Zod | Runtime + compile-time validation |

### 3.2 Development Tools

| Tool | Purpose |
|------|---------|
| pnpm | Package manager (monorepo) |
| Turborepo | Build system |
| Vitest | Unit testing |
| Playwright | E2E testing |
| Swagger/OpenAPI | API documentation |
| ESLint + Prettier | Code quality |

### 3.3 Infrastructure (MVP)

| Component | Service | Cost Tier |
|-----------|---------|-----------|
| Frontend Hosting | Vercel | Free tier |
| Backend Hosting | Railway | Hobby ($5/mo) |
| Database | Supabase PostgreSQL | Free tier |
| Redis | Upstash | Free tier |
| Email | SendGrid | Free tier (100/day) |

---

## 4. Data Model

### 4.1 Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              DATA MODEL                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  ┌──────────────┐         ┌──────────────┐         ┌──────────────┐     │
│  │    Tenant    │────────▶│     User     │         │   Category   │     │
│  │              │   1:N   │              │         │              │     │
│  │ id (uuid)    │         │ id (uuid)    │         │ id (uuid)    │     │
│  │ name         │         │ tenant_id    │         │ tenant_id    │     │
│  │ plan         │         │ email        │         │ name         │     │
│  │ settings     │         │ password_hash│         │ parent_id    │     │
│  │ created_at   │         │ role         │         │ created_at   │     │
│  └──────────────┘         │ created_at   │         └──────────────┘     │
│         │                 └──────────────┘                │              │
│         │ 1:N                    │                        │ 1:N          │
│         ▼                        │                        ▼              │
│  ┌──────────────┐                │                 ┌──────────────┐     │
│  │     Item     │◀───────────────┘                 │     Alert    │     │
│  │              │                                   │              │     │
│  │ id (uuid)    │─────────────────────────────────▶│ id (uuid)    │     │
│  │ tenant_id    │   1:N                             │ tenant_id    │     │
│  │ name         │                                   │ item_id      │     │
│  │ sku          │                                   │ type         │     │
│  │ quantity     │                                   │ status       │     │
│  │ threshold    │                                   │ triggered_at │     │
│  │ max_stock    │                                   │ resolved_at  │     │
│  │ unit         │                                   │ resolved_by  │     │
│  │ category_id  │                                   └──────────────┘     │
│  │ image_url    │                                          │             │
│  │ attributes   │                                          │ 1:N         │
│  │ created_at   │                                          ▼             │
│  └──────────────┘                                   ┌──────────────┐     │
│         │                                           │ Notification │     │
│         │ 1:N                                       │              │     │
│         ▼                                           │ id (uuid)    │     │
│  ┌──────────────┐                                   │ alert_id     │     │
│  │ StockUpdate  │                                   │ channel      │     │
│  │              │                                   │ recipient    │     │
│  │ id (uuid)    │                                   │ status       │     │
│  │ item_id      │                                   │ sent_at      │     │
│  │ prev_qty     │                                   │ error        │     │
│  │ new_qty      │                                   └──────────────┘     │
│  │ reason       │                                                        │
│  │ user_id      │                                                        │
│  │ created_at   │                                                        │
│  └──────────────┘                                                        │
│                                                                           │
│  ┌──────────────┐         ┌──────────────┐                               │
│  │   Supplier   │◀───────▶│ ItemSupplier │ (junction)                    │
│  │              │   N:M   │              │                               │
│  │ id (uuid)    │         │ item_id      │                               │
│  │ tenant_id    │         │ supplier_id  │                               │
│  │ name         │         │ is_preferred │                               │
│  │ email        │         └──────────────┘                               │
│  │ phone        │                                                        │
│  │ webhook_url  │                                                        │
│  └──────────────┘                                                        │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Core Entities Detail

#### Tenant
```typescript
interface Tenant {
  id: string;           // UUID
  name: string;         // Company name
  plan: 'free' | 'pro' | 'team';
  settings: {
    timezone: string;
    digestTime?: string;       // Daily digest hour
    defaultThreshold?: number;
    businessHours?: {
      enabled: boolean;
      start: string;
      end: string;
    };
  };
  createdAt: Date;
  updatedAt: Date;
}
```

#### User
```typescript
interface User {
  id: string;           // UUID
  tenantId: string;     // FK to Tenant
  email: string;        // Unique per tenant
  passwordHash: string; // Argon2id
  role: 'admin' | 'manager' | 'staff';
  notificationPrefs: {
    email: boolean;
    sms: boolean;
    slack: boolean;
    pushEnabled: boolean;
  };
  phone?: string;       // For SMS
  createdAt: Date;
  lastLoginAt?: Date;
}
```

#### Item
```typescript
interface Item {
  id: string;           // UUID
  tenantId: string;     // FK to Tenant
  name: string;         // Display name
  sku?: string;         // Optional SKU/barcode
  quantity: number;     // Current stock level
  threshold: number;    // Reorder point
  maxStock?: number;    // Maximum capacity
  unit: string;         // pieces, boxes, kg, etc.
  categoryId?: string;  // FK to Category
  imageUrl?: string;    // Item photo
  attributes: Record<string, unknown>; // Flexible metadata
  createdAt: Date;
  updatedAt: Date;
}
```

#### Alert
```typescript
interface Alert {
  id: string;           // UUID
  tenantId: string;     // FK to Tenant
  itemId: string;       // FK to Item
  type: 'low_stock' | 'out_of_stock' | 'critical';
  status: 'active' | 'snoozed' | 'resolved' | 'ordered';
  triggeredAt: Date;
  snoozeUntil?: Date;
  resolvedAt?: Date;
  resolvedBy?: string;  // User ID
  resolvedAction?: 'reordered' | 'dismissed' | 'restocked';
}
```

### 4.3 Multi-Tenancy Implementation

**Strategy**: Shared database with tenant column + PostgreSQL Row-Level Security

```sql
-- Enable RLS on all tables
ALTER TABLE items ENABLE ROW LEVEL SECURITY;

-- Create policy for tenant isolation
CREATE POLICY tenant_isolation ON items
  FOR ALL
  USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Application sets context before queries
SET app.tenant_id = 'uuid-of-tenant';
```

**Middleware Implementation**:
```typescript
// NestJS Guard
@Injectable()
export class TenantGuard implements CanActivate {
  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest();
    const tenantId = request.user?.tenantId;

    if (!tenantId) {
      throw new UnauthorizedException('Tenant context required');
    }

    // Set database context for RLS
    this.prisma.$executeRaw`SET app.tenant_id = ${tenantId}`;
    return true;
  }
}
```

---

## 5. Project Structure

```
inventory-alerts/
├── apps/
│   ├── web/                          # Next.js frontend
│   │   ├── app/
│   │   │   ├── (auth)/               # Auth pages
│   │   │   │   ├── login/
│   │   │   │   ├── register/
│   │   │   │   └── forgot-password/
│   │   │   ├── (dashboard)/          # Protected pages
│   │   │   │   ├── layout.tsx        # Dashboard layout with nav
│   │   │   │   ├── page.tsx          # Dashboard home
│   │   │   │   ├── items/
│   │   │   │   ├── alerts/
│   │   │   │   └── settings/
│   │   │   ├── api/                  # BFF API routes
│   │   │   └── layout.tsx
│   │   ├── components/
│   │   │   ├── ui/                   # Base UI components
│   │   │   ├── items/                # Item-related components
│   │   │   ├── alerts/               # Alert-related components
│   │   │   └── dashboard/            # Dashboard widgets
│   │   ├── lib/
│   │   │   ├── api.ts                # API client
│   │   │   ├── auth.ts               # Auth utilities
│   │   │   └── utils.ts
│   │   ├── hooks/
│   │   └── public/
│   │
│   └── api/                          # NestJS backend
│       ├── src/
│       │   ├── modules/
│       │   │   ├── auth/
│       │   │   │   ├── auth.controller.ts
│       │   │   │   ├── auth.service.ts
│       │   │   │   ├── auth.guard.ts
│       │   │   │   ├── jwt.strategy.ts
│       │   │   │   └── dto/
│       │   │   ├── items/
│       │   │   │   ├── items.controller.ts
│       │   │   │   ├── items.service.ts
│       │   │   │   ├── items.gateway.ts    # WebSocket
│       │   │   │   └── dto/
│       │   │   ├── alerts/
│       │   │   │   ├── alerts.controller.ts
│       │   │   │   ├── alerts.service.ts
│       │   │   │   ├── threshold.checker.ts
│       │   │   │   └── dto/
│       │   │   ├── notifications/
│       │   │   │   ├── notifications.service.ts
│       │   │   │   ├── notification.processor.ts  # BullMQ
│       │   │   │   ├── providers/
│       │   │   │   │   ├── email.provider.ts
│       │   │   │   │   ├── sms.provider.ts
│       │   │   │   │   ├── slack.provider.ts
│       │   │   │   │   └── webhook.provider.ts
│       │   │   │   └── dto/
│       │   │   ├── users/
│       │   │   │   ├── users.controller.ts
│       │   │   │   ├── users.service.ts
│       │   │   │   └── dto/
│       │   │   └── tenants/
│       │   │       ├── tenants.service.ts
│       │   │       └── dto/
│       │   ├── common/
│       │   │   ├── guards/
│       │   │   │   ├── tenant.guard.ts
│       │   │   │   └── roles.guard.ts
│       │   │   ├── filters/
│       │   │   │   └── http-exception.filter.ts
│       │   │   ├── interceptors/
│       │   │   │   └── audit.interceptor.ts
│       │   │   ├── pipes/
│       │   │   │   └── validation.pipe.ts
│       │   │   └── decorators/
│       │   ├── config/
│       │   │   ├── database.config.ts
│       │   │   ├── redis.config.ts
│       │   │   └── notification.config.ts
│       │   ├── app.module.ts
│       │   └── main.ts
│       ├── prisma/
│       │   ├── schema.prisma
│       │   ├── migrations/
│       │   └── seed.ts
│       └── test/
│
├── packages/
│   ├── shared/                       # Shared types and utilities
│   │   ├── types/
│   │   │   ├── item.ts
│   │   │   ├── alert.ts
│   │   │   └── user.ts
│   │   ├── validators/
│   │   │   └── schemas.ts            # Zod schemas
│   │   └── utils/
│   │       └── format.ts
│   ├── ui/                           # Shared UI component library
│   │   ├── components/
│   │   │   ├── Button/
│   │   │   ├── Card/
│   │   │   ├── Input/
│   │   │   ├── StockBar/
│   │   │   └── AlertBadge/
│   │   └── styles/
│   │       └── tokens.css            # Design tokens
│   └── config/                       # Shared configs
│       ├── eslint/
│       └── tsconfig/
│
├── docker/
│   ├── Dockerfile.api
│   ├── Dockerfile.web
│   └── docker-compose.yml            # Local development
│
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── deploy.yml
│
├── turbo.json
├── pnpm-workspace.yaml
├── package.json
└── README.md
```

---

## 6. API Design

### 6.1 RESTful Endpoints

```yaml
Base URL: /api/v1

# Authentication
POST   /auth/register          # Create tenant + admin user
POST   /auth/login             # Login, returns JWT + refresh token
POST   /auth/refresh           # Refresh access token
POST   /auth/logout            # Invalidate refresh token
POST   /auth/forgot-password   # Request password reset
POST   /auth/reset-password    # Complete password reset

# Items
GET    /items                  # List items (paginated, filterable)
POST   /items                  # Create item
GET    /items/:id              # Get item detail
PUT    /items/:id              # Update item
DELETE /items/:id              # Delete item (soft delete)
PATCH  /items/:id/quantity     # Quick quantity update
POST   /items/bulk-import      # Import from CSV
GET    /items/export           # Export to CSV

# Alerts
GET    /alerts                 # List active alerts
GET    /alerts/history         # Historical alerts (paginated)
PATCH  /alerts/:id/resolve     # Mark as resolved
PATCH  /alerts/:id/snooze      # Snooze alert
PATCH  /alerts/:id/order       # Mark as ordered

# Categories
GET    /categories             # List categories
POST   /categories             # Create category
PUT    /categories/:id         # Update category
DELETE /categories/:id         # Delete category

# Users (admin only)
GET    /users                  # List users
POST   /users                  # Invite user
PUT    /users/:id              # Update user
DELETE /users/:id              # Remove user

# Settings
GET    /settings/notifications # Get notification preferences
PUT    /settings/notifications # Update preferences
GET    /settings/integrations  # List integrations
POST   /settings/integrations  # Add integration

# Dashboard
GET    /dashboard/stats        # Dashboard statistics
GET    /dashboard/trends       # Consumption trends
```

### 6.2 WebSocket Events

```typescript
// Server → Client
'alert:new'        // New alert triggered
'alert:resolved'   // Alert resolved
'item:updated'     // Stock level changed
'item:critical'    // Critical threshold hit

// Client → Server
'subscribe:dashboard'   // Subscribe to real-time updates
'unsubscribe:dashboard' // Unsubscribe
```

### 6.3 Request/Response Examples

**Create Item**:
```http
POST /api/v1/items
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "Printer Paper A4",
  "sku": "PP-A4-500",
  "quantity": 50,
  "threshold": 20,
  "maxStock": 200,
  "unit": "boxes",
  "categoryId": "uuid-optional"
}
```

**Response**:
```json
{
  "id": "uuid",
  "name": "Printer Paper A4",
  "sku": "PP-A4-500",
  "quantity": 50,
  "threshold": 20,
  "maxStock": 200,
  "unit": "boxes",
  "status": "healthy",
  "createdAt": "2024-01-15T10:00:00Z"
}
```

**Quick Quantity Update**:
```http
PATCH /api/v1/items/uuid/quantity
Content-Type: application/json
Authorization: Bearer <token>

{
  "quantity": 45,
  "reason": "daily_usage"
}
```

---

## 7. Notification Pipeline

### 7.1 Flow Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        NOTIFICATION PIPELINE                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  [Stock Update API]                                                       │
│         │                                                                 │
│         ▼                                                                 │
│  ┌─────────────────┐                                                      │
│  │ ItemsService    │                                                      │
│  │ updateQuantity()│                                                      │
│  └────────┬────────┘                                                      │
│           │                                                               │
│           ▼                                                               │
│  ┌─────────────────┐     ┌─────────────────┐                             │
│  │ ThresholdChecker│────▶│   No Alert      │ (quantity > threshold)      │
│  │ checkThreshold()│     └─────────────────┘                             │
│  └────────┬────────┘                                                      │
│           │ (quantity <= threshold)                                       │
│           ▼                                                               │
│  ┌─────────────────┐     ┌─────────────────┐                             │
│  │ CooldownChecker │────▶│ Skip (in        │                             │
│  │ isInCooldown()  │     │ cooldown)       │                             │
│  └────────┬────────┘     └─────────────────┘                             │
│           │ (not in cooldown)                                             │
│           ▼                                                               │
│  ┌─────────────────┐                                                      │
│  │ AlertsService   │                                                      │
│  │ createAlert()   │                                                      │
│  └────────┬────────┘                                                      │
│           │                                                               │
│           ▼                                                               │
│  ┌─────────────────┐                                                      │
│  │  BullMQ Queue   │                                                      │
│  │ notification-   │                                                      │
│  │ queue           │                                                      │
│  └────────┬────────┘                                                      │
│           │                                                               │
│           ▼                                                               │
│  ┌─────────────────────────────────────────────────────────────┐         │
│  │              NotificationProcessor                           │         │
│  │                                                               │         │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────┐   │         │
│  │  │   Email   │ │    SMS    │ │   Slack   │ │  Webhook  │   │         │
│  │  │ Provider  │ │ Provider  │ │ Provider  │ │ Provider  │   │         │
│  │  └───────────┘ └───────────┘ └───────────┘ └───────────┘   │         │
│  │                                                               │         │
│  └──────────────────────────┬────────────────────────────────────┘         │
│                             │                                              │
│                             ▼                                              │
│                   ┌─────────────────┐                                      │
│                   │  Retry Logic    │                                      │
│                   │  - 3 attempts   │                                      │
│                   │  - Exp backoff  │                                      │
│                   │  - Dead letter  │                                      │
│                   └─────────────────┘                                      │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Alert Coalescing Rules

```typescript
const alertRules = {
  // Don't re-alert same item within period
  cooldownPeriod: 4 * 60 * 60 * 1000, // 4 hours

  // Critical items bypass cooldown
  criticalBypass: true,

  // Rate limit per tenant
  maxAlertsPerDay: 50,

  // Digest option
  digestMode: 'realtime' | 'daily',
  digestTime: '08:00',

  // Business hours (optional)
  businessHours: {
    enabled: false,
    start: '08:00',
    end: '18:00',
    timezone: 'UTC'
  }
};
```

### 7.3 Webhook Security (SSRF Prevention)

```typescript
const validateWebhookUrl = async (url: string): Promise<boolean> => {
  const parsed = new URL(url);

  // HTTPS only in production
  if (process.env.NODE_ENV === 'production' && parsed.protocol !== 'https:') {
    throw new Error('HTTPS required');
  }

  // Resolve DNS
  const addresses = await dns.resolve4(parsed.hostname);

  // Block private/internal IPs
  const blockedRanges = [
    '10.0.0.0/8',
    '172.16.0.0/12',
    '192.168.0.0/16',
    '127.0.0.0/8',
    '169.254.0.0/16',  // Cloud metadata
    '0.0.0.0/8',
  ];

  for (const ip of addresses) {
    if (isInBlockedRange(ip, blockedRanges)) {
      throw new Error('Internal IP not allowed');
    }
  }

  return true;
};
```

---

## 8. Security Implementation

### 8.1 Authentication Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         AUTHENTICATION FLOW                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  [Login Request]                                                          │
│       │                                                                   │
│       ▼                                                                   │
│  ┌─────────────┐     ┌─────────────┐                                     │
│  │  Validate   │────▶│ Rate Limit  │────▶ 429 Too Many Requests          │
│  │  Input      │     │ Check       │                                     │
│  └─────────────┘     └──────┬──────┘                                     │
│                             │                                             │
│                             ▼                                             │
│                      ┌─────────────┐                                     │
│                      │ Find User   │────▶ 401 Invalid credentials        │
│                      │ by Email    │                                     │
│                      └──────┬──────┘                                     │
│                             │                                             │
│                             ▼                                             │
│                      ┌─────────────┐                                     │
│                      │ Verify      │────▶ 401 Invalid credentials        │
│                      │ Password    │     (Argon2id)                       │
│                      └──────┬──────┘                                     │
│                             │                                             │
│                             ▼                                             │
│                      ┌─────────────┐                                     │
│                      │ Generate    │                                     │
│                      │ JWT + Refresh│                                     │
│                      └──────┬──────┘                                     │
│                             │                                             │
│                             ▼                                             │
│                      ┌─────────────┐                                     │
│                      │ Store Refresh│                                    │
│                      │ Token (Redis)│                                    │
│                      └──────┬──────┘                                     │
│                             │                                             │
│                             ▼                                             │
│                      [Return Tokens]                                      │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.2 JWT Structure

```typescript
// Access Token Payload (15 min expiry)
interface AccessTokenPayload {
  sub: string;       // User ID
  tenantId: string;  // Tenant ID
  email: string;
  role: 'admin' | 'manager' | 'staff';
  iat: number;
  exp: number;
}

// Refresh Token (7 day expiry, stored in Redis)
interface RefreshToken {
  token: string;     // Random string
  userId: string;
  createdAt: Date;
  expiresAt: Date;
}
```

### 8.3 Rate Limiting Configuration

| Endpoint Type | Limit | Window | Scope |
|---------------|-------|--------|-------|
| Login | 5 | 5 min | IP + Email |
| Password Reset | 3 | 1 hour | Email |
| API Read | 1000 | 1 min | Tenant |
| API Write | 100 | 1 min | Tenant |
| Bulk Import | 10 | 1 hour | Tenant |

### 8.4 Security Headers

```typescript
// Applied to all responses
const securityHeaders = {
  'Content-Security-Policy': "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'",
  'X-Content-Type-Options': 'nosniff',
  'X-Frame-Options': 'DENY',
  'X-XSS-Protection': '1; mode=block',
  'Strict-Transport-Security': 'max-age=31536000; includeSubDomains',
  'Referrer-Policy': 'strict-origin-when-cross-origin',
  'Permissions-Policy': 'geolocation=(), microphone=(), camera=(self)'
};
```

---

## 9. Trade-offs Summary

| Decision | Optimizing For | Sacrificing |
|----------|----------------|-------------|
| Modular Monolith | Development speed, operational simplicity | Independent service scaling |
| PostgreSQL | Data integrity, multi-tenant security (RLS) | Horizontal scaling complexity |
| PWA over Native | Faster launch, single codebase | Some native mobile features |
| JWT + Refresh | Stateless scaling, simplicity | Cannot revoke individual tokens instantly |
| BullMQ over SQS | TypeScript integration, development experience | AWS ecosystem benefits |
| TypeScript everywhere | Type safety, maintainability | Initial development overhead |
| SendGrid/Twilio | Reliability, simplicity | Cost at scale vs. self-hosted |

---

## 10. Scaling Path

### Phase 1: MVP (0-1,000 tenants)
```
┌─────────┐    ┌─────────┐    ┌───────────────┐
│ Vercel  │───▶│ Railway │───▶│ Supabase      │
│ (Next)  │    │ (Nest)  │    │ (PostgreSQL)  │
└─────────┘    └─────────┘    └───────────────┘
                    │
                    ▼
              ┌───────────┐
              │  Upstash  │
              │  (Redis)  │
              └───────────┘
```

### Phase 2: Growth (1,000-10,000 tenants)
- Migrate to managed PostgreSQL (Supabase Pro or RDS)
- Add read replicas for dashboard queries
- Dedicated Redis cluster for queue
- CDN for static assets
- APM tooling (Datadog/New Relic)

### Phase 3: Scale (10,000+ tenants)
- Consider tenant sharding
- Extract notification service
- Kubernetes for container orchestration
- Multi-region deployment

---

## 11. Implementation Phases

### MVP (Phase 1)
1. Auth module (register, login, JWT)
2. Items CRUD with quantity updates
3. Threshold checker and alert creation
4. Email notifications via SendGrid
5. Dashboard with alert counts
6. Basic item list and detail views

### v1.1 (Phase 2)
1. Slack integration
2. SMS via Twilio
3. Barcode scanning (browser API)
4. CSV import/export
5. Categories
6. Consumption trends

### v1.2 (Phase 3)
1. Webhook notifications
2. Supplier management
3. One-click reorder (PO generation)
4. Advanced analytics
5. Team invitations
6. Two-factor authentication

---

## 12. Next Steps

**DELEGATE:**

- **senior_dev**: Implement NestJS backend modules (auth, items, alerts, notifications) with Prisma ORM, following the API design and security patterns defined in this spec
- **ui**: Create component library based on the design tokens from UI research (StockBar, AlertBadge, Dashboard cards) using the Plum/Cream color palette
- **security**: Validate the multi-tenant isolation implementation, review JWT configuration, and test SSRF prevention in webhook handler

---

**Status**: PHASE_COMPLETE: planning
