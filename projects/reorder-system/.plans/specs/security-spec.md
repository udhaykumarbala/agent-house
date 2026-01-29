# Security Specification: Inventory Reorder Alert System

## Document Info

| Item | Detail |
|------|--------|
| **Author** | Security Expert Agent |
| **Date** | 2025-12-27 |
| **Status** | Planning Phase |
| **Based On** | security-research.md, architecture-research.md, product-research.md |

---

## 1. Executive Summary

This specification defines the security controls, implementation requirements, and testing criteria for the Inventory Reorder Alert System. The system manages business-critical inventory data, integrates with external notification channels (email, SMS, Slack, webhooks), and operates in a multi-tenant environment.

**Key Security Priorities**:
1. Multi-tenant data isolation (prevent cross-tenant data leakage)
2. SSRF prevention in webhook functionality
3. Secure authentication with optional MFA
4. Input validation and injection prevention
5. Audit logging for compliance and incident response

---

## 2. Threat Assessment Summary

### 2.1 Critical Assets

| Asset | Classification | Protection Level |
|-------|---------------|------------------|
| User credentials (passwords) | Secret | Argon2id hashed, never plaintext |
| API keys (user-generated) | Secret | Argon2id hashed at rest |
| Integration tokens (Slack, Twilio) | Confidential | AES-256-GCM encrypted |
| Inventory data | Business | Tenant-isolated, encrypted at rest |
| Audit logs | Internal | Append-only, integrity protected |
| Webhook URLs | Medium | SSRF-validated, HTTPS enforced |

### 2.2 Primary Threat Actors

| Actor | Risk Level | Key Mitigations |
|-------|------------|-----------------|
| Opportunistic attackers | High | Rate limiting, input validation, WAF |
| Disgruntled insiders | Medium | RBAC, audit logging, tenant isolation |
| Competitors | Medium | Data encryption, access controls |

### 2.3 Attack Surface

1. **Web Application** - Dashboard, forms, API endpoints
2. **REST API** - CRUD operations, bulk imports
3. **Notification Channels** - Email, SMS, Slack, webhooks (SSRF risk)
4. **Authentication System** - Login, registration, password reset

---

## 3. Security Requirements by OWASP Top 10

### 3.1 A01 - Broken Access Control

**Requirement**: All resources must be tenant-scoped with server-side authorization.

| Control | Implementation | Priority |
|---------|----------------|----------|
| Tenant isolation | PostgreSQL RLS + middleware validation | P0 |
| UUID for resources | No sequential/guessable IDs | P0 |
| RBAC enforcement | Guard on every protected route | P0 |
| CORS policy | Strict allowed origins list | P1 |
| Access denial logging | Log all 403 responses with context | P1 |

**Implementation Specification**:

```typescript
// Tenant isolation middleware (NestJS)
@Injectable()
export class TenantGuard implements CanActivate {
  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest();
    const user = request.user;
    const resourceTenantId = request.params.tenantId || request.body?.tenantId;

    // Always validate tenant ownership
    if (resourceTenantId && resourceTenantId !== user.tenantId) {
      throw new ForbiddenException('Access denied');
    }

    // Set tenant context for PostgreSQL RLS
    // SET app.tenant_id = user.tenantId
    return true;
  }
}
```

**PostgreSQL RLS Policy**:

```sql
-- Enable RLS on all tenant-scoped tables
ALTER TABLE items ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON items
  FOR ALL
  USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- Repeat for: users, alerts, categories, suppliers, notifications, stock_updates
```

---

### 3.2 A02 - Cryptographic Failures

**Requirement**: Strong cryptography for all sensitive data.

| Data Type | Algorithm | Implementation |
|-----------|-----------|----------------|
| Passwords | Argon2id | Memory: 64MB, Iterations: 3, Parallelism: 4 |
| API keys | Argon2id | Same as passwords (store hash only) |
| Integration tokens | AES-256-GCM | Envelope encryption via secrets manager |
| Session tokens | CSPRNG | crypto.randomBytes(32).toString('hex') |
| Webhook signatures | HMAC-SHA256 | Unique secret per integration |

**Password Policy**:

```typescript
const passwordPolicy = {
  minLength: 12,
  maxLength: 128,
  requireComplexity: false, // Per NIST 800-63B
  checkBreached: true,      // HaveIBeenPwned API
  hashAlgorithm: 'argon2id',
  argonConfig: {
    memoryCost: 65536,  // 64 MB
    timeCost: 3,
    parallelism: 4,
  },
};
```

**TLS Requirements**:

- Minimum TLS 1.2 (prefer TLS 1.3)
- Disable SSLv3, TLS 1.0, TLS 1.1
- HSTS header: `max-age=31536000; includeSubDomains`

---

### 3.3 A03 - Injection Prevention

**Requirement**: Prevent all injection attacks via parameterized queries and input validation.

| Vector | Protection | Priority |
|--------|------------|----------|
| SQL Injection | Prisma ORM (parameterized) | P0 |
| NoSQL Injection | N/A (using PostgreSQL) | N/A |
| Command Injection | No shell execution | P0 |
| Email Header Injection | Strict email validation | P1 |
| Template Injection | Sandboxed templating | P1 |
| CSV Formula Injection | Strip =, +, -, @, \| prefixes | P1 |

**Input Validation Schema (Zod)**:

```typescript
import { z } from 'zod';

export const ItemSchema = z.object({
  name: z.string()
    .min(1)
    .max(200)
    .regex(/^[a-zA-Z0-9\s\-_\.]+$/),
  sku: z.string()
    .max(50)
    .regex(/^[a-zA-Z0-9\-_]+$/)
    .optional(),
  quantity: z.number()
    .int()
    .min(0)
    .max(1_000_000),
  threshold: z.number()
    .int()
    .min(0)
    .max(1_000_000),
  unit: z.enum(['pieces', 'boxes', 'kg', 'liters', 'units']),
  categoryId: z.string().uuid().optional(),
});

export const EmailSchema = z.string()
  .email()
  .max(254)
  .transform(e => e.toLowerCase().trim());

export const PhoneSchema = z.string()
  .regex(/^\+[1-9]\d{1,14}$/); // E.164 format
```

---

### 3.4 A04 - Secure Design Principles

**Requirement**: Security built into architecture from inception.

| Principle | Implementation |
|-----------|----------------|
| Defense in Depth | Validation at API, service, and DB layers |
| Fail Secure | Generic error messages, detailed logging |
| Least Privilege | RBAC with minimal default permissions |
| Secure Defaults | MFA encouraged, alerts minimal disclosure |

**Role-Based Access Control (RBAC)**:

| Role | Permissions |
|------|-------------|
| `owner` | Full access, billing, delete tenant |
| `admin` | Manage users, settings, integrations |
| `manager` | Full item CRUD, alert management |
| `staff` | Update quantities, view items, acknowledge alerts |
| `viewer` | Read-only access to items and alerts |

```typescript
// Permission matrix
const permissions = {
  'items:create': ['owner', 'admin', 'manager'],
  'items:read': ['owner', 'admin', 'manager', 'staff', 'viewer'],
  'items:update': ['owner', 'admin', 'manager', 'staff'],
  'items:delete': ['owner', 'admin', 'manager'],
  'items:bulk_import': ['owner', 'admin', 'manager'],
  'settings:read': ['owner', 'admin'],
  'settings:update': ['owner', 'admin'],
  'users:manage': ['owner', 'admin'],
  'integrations:manage': ['owner', 'admin'],
  'analytics:view': ['owner', 'admin', 'manager'],
};
```

---

### 3.5 A05 - Security Configuration

**Requirement**: Hardened configuration with no defaults.

**Required HTTP Security Headers**:

```typescript
// NestJS helmet configuration
import helmet from 'helmet';

app.use(helmet({
  contentSecurityPolicy: {
    directives: {
      defaultSrc: ["'self'"],
      scriptSrc: ["'self'"],
      styleSrc: ["'self'", "'unsafe-inline'"], // For Tailwind
      imgSrc: ["'self'", "data:", "https:"],
      connectSrc: ["'self'", "wss:", process.env.API_URL],
      fontSrc: ["'self'"],
      objectSrc: ["'none'"],
      mediaSrc: ["'none'"],
      frameSrc: ["'none'"],
    },
  },
  crossOriginEmbedderPolicy: true,
  crossOriginOpenerPolicy: true,
  crossOriginResourcePolicy: { policy: 'same-origin' },
  hsts: {
    maxAge: 31536000,
    includeSubDomains: true,
    preload: true,
  },
  noSniff: true,
  referrerPolicy: { policy: 'strict-origin-when-cross-origin' },
  xssFilter: true,
}));
```

**Environment Configuration**:

```typescript
// Required environment variables (never hardcode)
const requiredEnvVars = [
  'DATABASE_URL',
  'REDIS_URL',
  'JWT_SECRET',           // 256-bit random
  'JWT_REFRESH_SECRET',   // 256-bit random
  'ENCRYPTION_KEY',       // 256-bit for AES-256-GCM
  'SENDGRID_API_KEY',
  'TWILIO_ACCOUNT_SID',
  'TWILIO_AUTH_TOKEN',
];
```

---

### 3.6 A06 - Vulnerable Components

**Requirement**: Continuous dependency scanning and updates.

| Tool | Integration Point | Action |
|------|-------------------|--------|
| Dependabot | GitHub | Auto-create PRs for updates |
| npm audit | CI pipeline | Fail build on high/critical |
| Snyk | CI pipeline | Block PRs with vulnerabilities |
| Trivy | Docker build | Scan container images |

**Package.json Scripts**:

```json
{
  "scripts": {
    "audit": "npm audit --audit-level=high",
    "audit:fix": "npm audit fix",
    "security:check": "snyk test"
  }
}
```

---

### 3.7 A07 - Authentication & Session Management

**Requirement**: Secure authentication with session protection.

**JWT Configuration**:

```typescript
const jwtConfig = {
  accessToken: {
    expiresIn: '15m',
    algorithm: 'HS256',
  },
  refreshToken: {
    expiresIn: '7d',
    algorithm: 'HS256',
    rotateOnUse: true,
  },
};
```

**Session Security**:

| Control | Value |
|---------|-------|
| Cookie flags | `Secure; HttpOnly; SameSite=Strict` |
| Session ID length | 256 bits (32 bytes) |
| Absolute timeout | 24 hours |
| Idle timeout | 30 minutes |
| Regenerate on login | Yes |
| Regenerate on privilege change | Yes |

**Password Reset Security**:

```typescript
const passwordResetConfig = {
  tokenExpiry: '1 hour',
  tokenLength: 32,             // bytes
  singleUse: true,
  rateLimit: {
    perEmail: { limit: 3, window: '1 hour' },
    perIP: { limit: 5, window: '1 hour' },
  },
  response: 'generic',         // Don't confirm email exists
  notifyOnChange: true,        // Email confirmation of change
};
```

**Multi-Factor Authentication (MFA)**:

| Setting | Value |
|---------|-------|
| Algorithm | TOTP (RFC 6238) |
| Period | 30 seconds |
| Digits | 6 |
| Required for | Admin/Owner roles (optional for others) |
| Backup codes | 10 codes, single-use |

---

### 3.8 A08 - Data Integrity

**Requirement**: Validate all external data and protect against tampering.

**CSV Import Security**:

```typescript
const csvImportConfig = {
  maxFileSize: '5MB',
  maxRows: 10_000,
  allowedMimeTypes: ['text/csv', 'application/vnd.ms-excel'],
  sanitizeFormulas: true,     // Strip =, +, -, @, |
  validateSchema: true,
  processAsync: true,         // Background queue
  auditLog: true,
};

// Formula injection prevention
function sanitizeCsvCell(value: string): string {
  const dangerousPrefixes = ['=', '+', '-', '@', '|', '%'];
  if (dangerousPrefixes.some(p => value.startsWith(p))) {
    return `'${value}`; // Prefix with single quote
  }
  return value;
}
```

**Webhook Signature Verification**:

```typescript
// Outgoing webhooks - sign with HMAC
function signWebhookPayload(payload: object, secret: string): string {
  const timestamp = Date.now();
  const message = `${timestamp}.${JSON.stringify(payload)}`;
  const signature = crypto
    .createHmac('sha256', secret)
    .update(message)
    .digest('hex');
  return `t=${timestamp},v1=${signature}`;
}

// Headers sent with webhook
// X-Webhook-Signature: t=1234567890,v1=abc123...
// X-Webhook-Timestamp: 1234567890
```

---

### 3.9 A09 - Security Logging & Monitoring

**Requirement**: Comprehensive audit logging with alerting.

**Audit Events**:

| Event | Data Logged | Retention |
|-------|-------------|-----------|
| Login success | userId, IP, userAgent, timestamp | 90 days |
| Login failure | email (masked), IP, reason, timestamp | 90 days |
| Password change | userId, IP, timestamp | 1 year |
| MFA enable/disable | userId, actor, timestamp | 1 year |
| Item CRUD | actor, action, itemId, changes | 90 days |
| Threshold change | actor, itemId, oldValue, newValue | 90 days |
| Bulk import | actor, recordCount, status | 90 days |
| API key create/revoke | actor, keyId (partial) | 1 year |
| Permission change | actor, targetUser, oldRole, newRole | 1 year |
| Integration add/remove | actor, integrationType | 1 year |

**Log Format**:

```typescript
interface AuditLog {
  id: string;
  timestamp: string;           // ISO 8601
  tenantId: string;
  actorId: string;
  actorType: 'user' | 'system' | 'api_key';
  action: string;
  resource: string;
  resourceId: string;
  changes?: {
    field: string;
    oldValue: any;
    newValue: any;
  }[];
  metadata: {
    ip: string;
    userAgent: string;
    requestId: string;
  };
  severity: 'info' | 'warning' | 'critical';
}
```

**Security Alerts**:

| Trigger | Threshold | Action |
|---------|-----------|--------|
| Failed logins | 5 in 5 minutes | Temporary account lock |
| New geography login | First from country | Email notification |
| Bulk data export | Any | Email to owner |
| Admin privilege granted | Any | Email to all admins |
| MFA disabled | Any | Email confirmation |
| API rate limit exceeded | 3x in hour | Temporary block |

---

### 3.10 A10 - SSRF Prevention (Critical for Webhooks)

**Requirement**: Strict validation of all user-provided URLs.

**Webhook URL Validation**:

```typescript
import { isIP } from 'net';
import { lookup } from 'dns/promises';

const BLOCKED_IP_RANGES = [
  { start: '10.0.0.0', end: '10.255.255.255' },       // 10.0.0.0/8
  { start: '172.16.0.0', end: '172.31.255.255' },    // 172.16.0.0/12
  { start: '192.168.0.0', end: '192.168.255.255' },  // 192.168.0.0/16
  { start: '127.0.0.0', end: '127.255.255.255' },    // 127.0.0.0/8
  { start: '169.254.0.0', end: '169.254.255.255' },  // Link-local
  { start: '0.0.0.0', end: '0.255.255.255' },        // 0.0.0.0/8
];

async function validateWebhookUrl(url: string): Promise<void> {
  // 1. Parse URL
  let parsed: URL;
  try {
    parsed = new URL(url);
  } catch {
    throw new BadRequestException('Invalid URL format');
  }

  // 2. Protocol check
  if (parsed.protocol !== 'https:') {
    throw new BadRequestException('HTTPS required for webhooks');
  }

  // 3. Port restriction
  const port = parsed.port || '443';
  if (!['443', '8443'].includes(port)) {
    throw new BadRequestException('Invalid port');
  }

  // 4. DNS resolution and IP validation
  const hostname = parsed.hostname;

  // Check if hostname is already an IP
  if (isIP(hostname)) {
    if (isPrivateIP(hostname)) {
      throw new BadRequestException('Private IPs not allowed');
    }
    return;
  }

  // Resolve hostname
  const addresses = await lookup(hostname, { all: true });

  for (const addr of addresses) {
    if (isPrivateIP(addr.address)) {
      throw new BadRequestException('URL resolves to private IP');
    }
  }
}

// Re-resolve before each request to prevent DNS rebinding
async function sendWebhookSecurely(url: string, payload: object): Promise<void> {
  await validateWebhookUrl(url); // Re-validate

  // Use custom DNS resolver that checks IP again
  const response = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Webhook-Signature': signPayload(payload),
    },
    body: JSON.stringify(payload),
    signal: AbortSignal.timeout(5000), // 5 second timeout
  });
}
```

---

## 4. Rate Limiting Specification

**Requirement**: Prevent abuse and DoS attacks.

| Endpoint Category | Limit | Window | Action on Exceed |
|-------------------|-------|--------|------------------|
| Login attempts | 5 | 5 minutes | Temporary block |
| Password reset | 3 | 1 hour | Block + notify |
| Registration | 5 | 1 hour | Block IP |
| API read (auth'd) | 1000 | 1 minute | 429 response |
| API write (auth'd) | 100 | 1 minute | 429 response |
| Bulk import | 10 | 1 hour | 429 response |
| Webhook triggers | 100 | 1 minute | Queue delay |
| Unauthenticated | 30 | 1 minute | Block IP |

**Implementation**:

```typescript
// Redis-based rate limiter
import { RateLimiterRedis } from 'rate-limiter-flexible';

const loginLimiter = new RateLimiterRedis({
  storeClient: redisClient,
  keyPrefix: 'rl:login',
  points: 5,              // attempts
  duration: 300,          // 5 minutes
  blockDuration: 900,     // 15 minute block
});

const apiLimiter = new RateLimiterRedis({
  storeClient: redisClient,
  keyPrefix: 'rl:api',
  points: 1000,
  duration: 60,
});
```

---

## 5. API Security Specification

### 5.1 API Key Management

```typescript
interface APIKeyConfig {
  prefix: 'iras_';                    // Inventory Reorder Alert System
  keyLength: 32;                      // bytes (256 bits)
  hashAlgorithm: 'argon2id';
  scopes: ['read', 'write', 'admin'];
  expirationOptions: ['30d', '90d', '1y', 'never'];
  maxKeysPerUser: 10;
}

// Key format: iras_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
// Only show full key once on creation
// Store Argon2id hash only
```

### 5.2 API Response Security

```typescript
// Never expose internal errors
const errorResponse = {
  statusCode: 500,
  message: 'An error occurred',
  requestId: 'req_abc123',  // For support lookup
  // NO stack traces, NO internal details
};

// Standardized error codes
enum ErrorCode {
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  UNAUTHORIZED = 'UNAUTHORIZED',
  FORBIDDEN = 'FORBIDDEN',
  NOT_FOUND = 'NOT_FOUND',
  RATE_LIMITED = 'RATE_LIMITED',
  INTERNAL_ERROR = 'INTERNAL_ERROR',
}
```

---

## 6. Notification Channel Security

### 6.1 Email Security

| Requirement | Implementation |
|-------------|----------------|
| Sender authentication | SPF, DKIM, DMARC configured |
| Recipient validation | E.164 + domain MX check |
| Rate limiting | 100 emails/tenant/hour |
| Content | No sensitive data in body |
| Unsubscribe | One-click unsubscribe link |

### 6.2 SMS Security

| Requirement | Implementation |
|-------------|----------------|
| Provider | Twilio (verified account) |
| Format | E.164 international |
| Rate limiting | 50 SMS/tenant/hour |
| Content | Short, no PII |
| Compliance | TCPA opt-in tracking |

### 6.3 Slack/Teams Security

| Requirement | Implementation |
|-------------|----------------|
| OAuth scopes | Minimal required only |
| Token storage | AES-256-GCM encrypted |
| Token refresh | Automatic with retry |
| Channel validation | Verify membership before post |

---

## 7. Data Protection & Privacy

### 7.1 Data Retention Policy

| Data Type | Retention | Deletion Method |
|-----------|-----------|-----------------|
| Active user data | While account active | Cascade delete |
| Deleted accounts | 30 days soft delete | Hard delete after |
| Audit logs | 1 year | Automated purge |
| Session data | 7 days | Redis TTL |
| Notification logs | 90 days | Automated purge |

### 7.2 GDPR Compliance

| Right | Implementation |
|-------|----------------|
| Access | Data export endpoint (JSON) |
| Rectification | Self-service profile edit |
| Erasure | Account deletion with cascade |
| Portability | Export in machine-readable format |
| Objection | Granular notification preferences |

---

## 8. Security Testing Requirements

### 8.1 Automated Tests (CI Pipeline)

| Test Type | Tool | Frequency |
|-----------|------|-----------|
| SAST | CodeQL / Semgrep | Every commit |
| Dependency scan | npm audit + Snyk | Every commit |
| Secret detection | GitLeaks | Every commit |
| Container scan | Trivy | On image build |

### 8.2 Manual Security Tests

| Test | Scope | Frequency |
|------|-------|-----------|
| OWASP ZAP scan | Full API | Weekly |
| Penetration test | Full application | Quarterly |
| Tenant isolation audit | Cross-tenant access | Monthly |
| SSRF testing | Webhook endpoints | Per release |

### 8.3 Security Test Cases

```
Authentication:
[ ] Brute force protection locks account after 5 attempts
[ ] Session regenerates on login
[ ] Session expires after idle timeout
[ ] Password reset tokens are single-use
[ ] MFA cannot be bypassed

Authorization:
[ ] User A cannot access User B's items
[ ] Staff cannot access admin settings
[ ] API key respects assigned scopes
[ ] Tenant A data isolated from Tenant B

Injection:
[ ] SQL injection blocked on search
[ ] XSS blocked in item names
[ ] CSRF tokens validated on mutations
[ ] CSV import sanitizes formulas

SSRF:
[ ] Webhook to 127.0.0.1 blocked
[ ] Webhook to 169.254.169.254 blocked
[ ] Webhook to internal DNS blocked
[ ] DNS rebinding prevented
```

---

## 9. Incident Response Requirements

### 9.1 Security Incident Classification

| Severity | Response Time | Examples |
|----------|---------------|----------|
| Critical | 1 hour | Data breach, active attack |
| High | 4 hours | Account compromise, service DoS |
| Medium | 24 hours | Vulnerability disclosure |
| Low | 48 hours | Suspicious activity |

### 9.2 Kill Switch Requirements

| Capability | Endpoint | Access |
|------------|----------|--------|
| Force logout all | Admin API | Owner only |
| Disable tenant | Admin API | Platform admin |
| Revoke all API keys | Admin API | Owner only |
| Disable integrations | Admin API | Owner/Admin |
| IP block | Platform | Platform admin |

---

## 10. Implementation Priority

| Control | Priority | Effort | Risk Mitigation |
|---------|----------|--------|-----------------|
| Input validation (Zod) | P0 | Low | Injection, XSS |
| Password hashing (Argon2id) | P0 | Low | Credential theft |
| HTTPS enforcement | P0 | Low | MitM attacks |
| Tenant isolation (RLS) | P0 | Medium | Data leakage |
| Session security | P0 | Medium | Session hijacking |
| Rate limiting | P1 | Medium | DoS, brute force |
| SSRF prevention | P1 | Medium | Internal network access |
| Security headers | P1 | Low | Various attacks |
| Audit logging | P1 | Medium | Incident response |
| MFA support | P2 | High | Account takeover |
| API key management | P2 | Medium | API abuse |
| Dependency scanning | P2 | Low | Supply chain |

---

## 11. Security Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           SECURITY ARCHITECTURE                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐                                                            │
│  │   Client    │◄──────── TLS 1.3 ────────────────────────────┐            │
│  │  (Browser)  │                                               │            │
│  └──────┬──────┘                                               │            │
│         │                                                       │            │
│         ▼                                                       │            │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐        │            │
│  │    CDN/     │───▶│    WAF      │───▶│   Next.js   │        │            │
│  │ Cloudflare  │    │  (Rules)    │    │  Frontend   │        │            │
│  └─────────────┘    └─────────────┘    └──────┬──────┘        │            │
│                                               │                │            │
│                     ┌─────────────────────────┼────────────────┘            │
│                     │                         │                              │
│                     ▼                         ▼                              │
│  ┌──────────────────────────────────────────────────────────────┐          │
│  │                      NestJS API Server                        │          │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────┐    │          │
│  │  │  Helmet   │ │   CORS    │ │Rate Limit │ │   Auth    │    │          │
│  │  │ (Headers) │ │ (Origins) │ │  (Redis)  │ │  Guard    │    │          │
│  │  └───────────┘ └───────────┘ └───────────┘ └───────────┘    │          │
│  │                         │                                     │          │
│  │                         ▼                                     │          │
│  │  ┌───────────────────────────────────────────────────────┐   │          │
│  │  │                  Tenant Guard                          │   │          │
│  │  │  • Extract tenant from JWT                            │   │          │
│  │  │  • Set PostgreSQL app.tenant_id                       │   │          │
│  │  │  • Validate resource ownership                        │   │          │
│  │  └───────────────────────────────────────────────────────┘   │          │
│  │                         │                                     │          │
│  │                         ▼                                     │          │
│  │  ┌───────────────────────────────────────────────────────┐   │          │
│  │  │                  Input Validation                      │   │          │
│  │  │  • Zod schemas for all endpoints                      │   │          │
│  │  │  • Parameterized queries (Prisma)                     │   │          │
│  │  └───────────────────────────────────────────────────────┘   │          │
│  └──────────────────────────────────────────────────────────────┘          │
│                                    │                                         │
│         ┌──────────────────────────┼────────────────────────┐               │
│         │                          │                        │               │
│         ▼                          ▼                        ▼               │
│  ┌─────────────┐           ┌─────────────┐          ┌─────────────┐        │
│  │ PostgreSQL  │           │    Redis    │          │   BullMQ    │        │
│  │             │           │             │          │   Workers   │        │
│  │ • RLS       │           │ • Sessions  │          │             │        │
│  │ • Encrypted │           │ • Rate limit│          │ • Email     │        │
│  │ • Backups   │           │ • Cache     │          │ • SMS       │        │
│  └─────────────┘           └─────────────┘          │ • Webhooks  │        │
│                                                      └──────┬──────┘        │
│                                                             │               │
│                    ┌────────────────────────────────────────┘               │
│                    │                                                         │
│                    ▼                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐       │
│  │                    SSRF Protection Layer                          │       │
│  │  • URL validation (HTTPS only)                                   │       │
│  │  • DNS resolution + IP check                                     │       │
│  │  • Private IP range blocking                                     │       │
│  │  • Timeout enforcement (5s)                                      │       │
│  └─────────────────────────────────────────────────────────────────┘       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 12. Delegation

```
DELEGATE:
- architect: Validate multi-tenant RLS implementation, ensure defense-in-depth layers are properly integrated
- senior_dev: Implement authentication module (Argon2id, JWT, MFA), rate limiting middleware, input validation with Zod, SSRF protection for webhooks, audit logging
```

---

## Summary

This security specification addresses all OWASP Top 10 concerns for the Inventory Reorder Alert System with particular emphasis on:

1. **Multi-tenant isolation** - PostgreSQL RLS + application-level guards
2. **SSRF prevention** - Critical for webhook functionality
3. **Authentication strength** - Argon2id, JWT, optional MFA
4. **Input validation** - Zod schemas on all endpoints
5. **Audit logging** - Comprehensive event tracking

Security controls are prioritized by risk, with P0 items required before MVP launch.

---

**Status**: PHASE_COMPLETE: planning
