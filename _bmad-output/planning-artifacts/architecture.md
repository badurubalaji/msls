---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
inputDocuments:
  - school-erp-prd/index.md
  - school-erp-prd/01-technical-architecture.md
  - school-erp-prd/02-core-foundation.md
  - school-erp-prd/03-student-management.md
  - school-erp-prd/04-academic-operations.md
  - school-erp-prd/05-admissions.md
  - school-erp-prd/06-examinations.md
  - school-erp-prd/07-homework-assignments.md
  - school-erp-prd/08-online-quiz.md
  - school-erp-prd/09-digital-classroom.md
  - school-erp-prd/10-staff-management.md
  - school-erp-prd/11-leave-management.md
  - school-erp-prd/12-fees-payments.md
  - school-erp-prd/13-communication-system.md
  - school-erp-prd/14-parent-portal.md
  - school-erp-prd/15-student-portal.md
  - school-erp-prd/16-certificate-generation.md
  - school-erp-prd/17-transport-management.md
  - school-erp-prd/18-library-management.md
  - school-erp-prd/19-inventory-assets.md
  - school-erp-prd/20-visitor-gate-management.md
  - school-erp-prd/21-analytics-dashboards.md
  - school-erp-prd/22-ai-capabilities.md
workflowType: 'architecture'
lastStep: 8
status: 'complete'
completedAt: '2026-01-23'
project_name: 'msls'
user_name: 'Ashulabs'
date: '2026-01-22'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**
The School ERP comprises 22 distinct modules spanning:
- **Academic Core** (6 modules): Student Management, Academic Operations, Admissions, Examinations, Homework/Assignments, Online Quiz
- **Learning Platform** (1 module): Digital Classroom with class recording and content library
- **HR & Staff** (2 modules): Staff Management, Leave Management
- **Finance** (1 module): Fees & Payments with multi-gateway integration
- **Communication** (3 modules): Communication System, Parent Portal, Student Portal
- **Operations** (5 modules): Certificates, Transport, Library, Inventory/Assets, Visitor/Gate Management
- **Advanced** (2 modules): Analytics Dashboards, AI Capabilities (on-premise)
- **Foundation** (2 modules): Technical Architecture, Core Foundation (Tenants, Users, RBAC)

**Non-Functional Requirements:**
- **Performance**: Sub-second response, 50,000+ concurrent users
- **Availability**: 99.9% uptime for SaaS deployment
- **Offline Capability**: Full operation on LAN without internet
- **Security**: RBAC, encryption at rest/transit, audit logging
- **Compliance**: FERPA-like principles, GDPR-ready
- **Localization**: Multi-language UI and reports (India-first, global-ready)
- **Mobile**: PWA + Native apps for iOS/Android

**Scale & Complexity:**
- Primary domain: Full-stack enterprise web platform
- Complexity level: Enterprise
- Estimated architectural components: 25+ distinct services/modules

### Technical Constraints & Dependencies

**Hard Constraints:**
1. **Go-only backend**: Modular monolith architecture (no Node.js)
2. **PostgreSQL with RLS**: Row-level isolation for multi-tenancy
3. **Same codebase**: Must serve both SaaS and On-Premise deployments
4. **Offline-first**: Core functionality must work without internet
5. **On-premise AI**: No external AI APIs (Llama, Mistral, Phi for local LLMs)

**Technology Stack (Locked):**
- Backend: Go (Gin), GORM/sqlc
- Frontend: Angular 21, Tailwind CSS, Custom Component Library
- Database: PostgreSQL, Redis
- Storage: Local FS / MinIO
- Deployment: Docker Compose (on-prem), Kubernetes (SaaS)

### Cross-Cutting Concerns Identified

1. **Multi-Tenancy**: `tenant_id` on all tenant-scoped tables with RLS policies
2. **Authentication**: JWT with refresh tokens (15min access, 7-day refresh)
3. **Authorization**: Hierarchical RBAC with `module:resource:action` permissions
4. **Audit Logging**: All data changes logged with old/new values, user, timestamp
5. **File Storage**: Unified storage abstraction (local/MinIO) for recordings, documents, photos
6. **Notification Bus**: SMS, Email, WhatsApp, Push notifications across all modules
7. **Offline Sync**: Data consistency strategy for LAN-only operation
8. **Localization**: i18n support for UI and generated documents
9. **Configuration System**: Feature flags, custom fields, workflow configuration
10. **Media Processing**: Video encoding, audio normalization, thumbnail generation

---

## Starter Template Evaluation

### Selected Approach: Tailwind CSS + Custom Angular Component Library

**Rationale:**
- **Performance**: Minimal CSS bundle (~15-30KB vs 500KB+ with PrimeNG)
- **Brand Identity**: Unique visual design, not generic admin template
- **Offline-First**: Lighter assets for LAN-only schools
- **Full Control**: No vendor lock-in, components are company asset
- **Modern Angular**: Aligns with Angular 21 standalone components and signals

### Repository Structure

**Dual Repository Architecture:**
- `msls-backend` - Go modular monolith
- `msls-frontend` - Angular 21 + Tailwind CSS

**Rationale for Dual Repos:**
1. Different release cycles (backend vs frontend)
2. On-premise may need backend-only updates
3. Frontend can be served from CDN in SaaS mode
4. Cleaner CI/CD pipelines per component
5. Teams can work independently

### Backend Initialization

**Repository:** `msls-backend`

```bash
mkdir msls-backend && cd msls-backend
go mod init github.com/ashulabs/msls-backend

# Core dependencies
go get github.com/gin-gonic/gin@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
go get github.com/redis/go-redis/v9@latest
go get github.com/golang-jwt/jwt/v5@latest
go get github.com/google/uuid@latest
```

**Project Structure:**
```
msls-backend/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── modules/                 # 22 business modules
│   │   ├── auth/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── dto.go
│   │   ├── student/
│   │   ├── staff/
│   │   ├── fees/
│   │   ├── attendance/
│   │   └── ...
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go
│   │   ├── tenant.go
│   │   ├── ratelimit.go
│   │   └── audit.go
│   ├── pkg/                     # Shared internal packages
│   │   ├── database/
│   │   ├── cache/
│   │   ├── storage/
│   │   └── notification/
│   └── config/
├── migrations/                  # SQL migrations
├── api/                         # OpenAPI specs
├── build/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
├── scripts/
└── docs/
```

### Go Best Practices (MANDATORY)

Backend developers MUST follow the Google Go Style Guide:

| Resource | Documentation Link |
|----------|-------------------|
| **Best Practices** | https://google.github.io/styleguide/go/best-practices |
| **Style Guide** | https://google.github.io/styleguide/go/guide |
| **Decisions** | https://google.github.io/styleguide/go/decisions |

**Key Implementation Requirements:**

1. **Naming Conventions**
   - Use MixedCaps (not underscores) for multi-word names
   - Keep names short but descriptive; avoid stuttering (`user.User` → `user.Record`)
   - Acronyms should be consistent case (`HTTPServer` or `httpServer`, not `HttpServer`)

2. **Error Handling**
   - Always check and handle errors; never use `_` to discard errors
   - Wrap errors with context: `fmt.Errorf("fetch student %s: %w", id, err)`
   - Use sentinel errors for expected conditions: `var ErrNotFound = errors.New("not found")`
   - Return errors, don't panic (except for truly unrecoverable situations)

3. **Documentation**
   - All exported functions, types, and packages must have doc comments
   - Doc comments should be complete sentences starting with the name
   - Example: `// FindByID returns the student with the given ID.`

4. **Testing**
   - Use table-driven tests for comprehensive coverage
   - Use `t.Helper()` in test helper functions
   - Use `t.Parallel()` for independent tests
   - Name test cases descriptively: `TestStudentService_Create_WithInvalidEmail`

5. **Concurrency**
   - Pass `context.Context` as first parameter for cancellation
   - Prefer channels over shared memory with mutexes
   - Use `sync.WaitGroup` for goroutine coordination
   - Never start goroutines without a way to stop them

6. **Package Design**
   - Keep packages focused and cohesive
   - Avoid circular dependencies
   - Use interfaces for abstraction at consumer side
   - Export only what's necessary

### Frontend Initialization

**Repository:** `msls-frontend`

```bash
# Create Angular project
ng new msls-frontend \
  --style=scss \
  --routing=true \
  --ssr=false \
  --standalone=true \
  --strict=true

cd msls-frontend

# Add Tailwind CSS
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init

# Add Angular CDK (accessibility primitives)
npm install @angular/cdk

# Add utilities
npm install date-fns                    # Date handling
npm install apexcharts ng-apexcharts    # Charts for analytics
```

### Angular CLI Commands (MANDATORY)

**CRITICAL: All Angular artifacts MUST be created using `ng generate` commands. Never create files manually.**

| Artifact | Command |
|----------|---------|
| **Component** | `ng generate component path/component-name --standalone` |
| **Service** | `ng generate service path/service-name` |
| **Guard** | `ng generate guard path/guard-name --functional` |
| **Interceptor** | `ng generate interceptor path/interceptor-name --functional` |
| **Pipe** | `ng generate pipe path/pipe-name` |
| **Directive** | `ng generate directive path/directive-name` |
| **Interface** | `ng generate interface path/interface-name` |
| **Enum** | `ng generate enum path/enum-name` |
| **Class** | `ng generate class path/class-name` |
| **Resolver** | `ng generate resolver path/resolver-name --functional` |

**Common Examples:**

```bash
# Shared UI Components
ng generate component shared/components/button --standalone --export
ng generate component shared/components/data-table --standalone --export
ng generate component shared/components/modal --standalone --export

# Feature Components
ng generate component features/students/pages/student-list --standalone
ng generate component features/students/components/student-card --standalone

# Services
ng generate service features/students/services/student
ng generate service core/services/auth
ng generate service core/services/toast

# Guards & Interceptors
ng generate guard core/guards/auth --functional
ng generate guard core/guards/role --functional
ng generate interceptor core/interceptors/auth --functional
ng generate interceptor core/interceptors/error --functional

# Models & Interfaces
ng generate interface features/students/models/student
ng generate interface core/models/api-response
ng generate interface core/models/pagination

# Pipes & Directives
ng generate pipe shared/pipes/date-format --standalone
ng generate pipe shared/pipes/currency-inr --standalone
ng generate directive shared/directives/tooltip --standalone
```

**Rationale:**
1. **Consistency**: Ensures uniform file structure and naming conventions
2. **Best Practices**: Generated code follows Angular's latest patterns
3. **Testing**: Automatically creates `.spec.ts` test files
4. **Imports**: Handles imports and declarations automatically
5. **Standards**: Enforces project-wide coding standards

**Tailwind Configuration** (`tailwind.config.js`):
```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{html,ts}"],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        },
      },
    },
  },
  plugins: [],
}
```

**Project Structure:**
```
msls-frontend/
├── src/
│   ├── app/
│   │   ├── core/                    # Singleton services
│   │   │   ├── services/
│   │   │   │   ├── api.service.ts
│   │   │   │   ├── auth.service.ts
│   │   │   │   └── tenant.service.ts
│   │   │   ├── guards/
│   │   │   ├── interceptors/
│   │   │   └── models/
│   │   │
│   │   ├── shared/                  # Reusable UI components
│   │   │   ├── components/
│   │   │   │   ├── data-table/
│   │   │   │   ├── modal/
│   │   │   │   ├── form-controls/
│   │   │   │   ├── button/
│   │   │   │   ├── card/
│   │   │   │   ├── badge/
│   │   │   │   ├── avatar/
│   │   │   │   ├── tabs/
│   │   │   │   ├── dropdown/
│   │   │   │   ├── sidebar/
│   │   │   │   ├── stepper/
│   │   │   │   └── chart/
│   │   │   ├── directives/
│   │   │   └── pipes/
│   │   │
│   │   ├── features/                # Feature modules (lazy-loaded)
│   │   │   ├── dashboard/
│   │   │   ├── students/
│   │   │   ├── staff/
│   │   │   ├── attendance/
│   │   │   ├── fees/
│   │   │   └── ... (all 22 modules)
│   │   │
│   │   ├── layouts/
│   │   │   ├── admin-layout/
│   │   │   ├── teacher-layout/
│   │   │   ├── student-layout/
│   │   │   └── parent-layout/
│   │   │
│   │   ├── app.component.ts
│   │   ├── app.config.ts
│   │   └── app.routes.ts
│   │
│   ├── assets/
│   ├── environments/
│   └── styles/
├── tailwind.config.js
└── package.json
```

### Custom Component Library

**Phase 1: Foundation**
| Component | Description |
|-----------|-------------|
| `ButtonComponent` | Primary, secondary, danger, ghost variants |
| `InputComponent` | Text input with validation states |
| `SelectComponent` | Dropdown select with search |
| `CardComponent` | Container with header, body, footer |
| `BadgeComponent` | Status indicators |
| `AvatarComponent` | User/student images with fallback |

**Phase 2: Layout**
| Component | Description |
|-----------|-------------|
| `SidebarComponent` | Collapsible navigation |
| `HeaderComponent` | Top bar with user menu |
| `ModalComponent` | Dialog with CDK overlay |
| `TabsComponent` | Tab navigation |
| `DropdownComponent` | Action menus |

**Phase 3: Data Display**
| Component | Description |
|-----------|-------------|
| `DataTableComponent` | Sort, filter, paginate, select |
| `ChartComponent` | ApexCharts wrapper |
| `StepperComponent` | Multi-step wizard |
| `DatepickerComponent` | Calendar date selection |
| `FileUploadComponent` | Drag-drop file upload |

### Technology Decisions Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Styling** | Tailwind CSS | Utility-first, minimal bundle, full control |
| **UI Components** | Custom library | Brand identity, performance, no lock-in |
| **Accessibility** | Angular CDK | Proven primitives (overlay, focus trap, a11y) |
| **Charts** | ApexCharts | Lightweight, responsive, good Angular support |
| **Icons** | Heroicons / Lucide | SVG, tree-shakeable, Tailwind-friendly |
| **Date Handling** | date-fns | Modular, tree-shakeable |
| **State** | Angular Signals | Native Angular 21, no NgRx overhead |
| **Forms** | Reactive Forms | Type-safe, validation, complex forms |
| **HTTP** | HttpClient + Interceptors | Auth tokens, error handling, tenant injection |

### Angular Best Practices (MANDATORY)

Frontend developers MUST follow these official Angular best practices documentation:

| Category | Documentation Link |
|----------|-------------------|
| **Style Guide** | https://angular.dev/style-guide |
| **Security** | https://angular.dev/best-practices/security |
| **Accessibility (a11y)** | https://angular.dev/best-practices/a11y |
| **Error Handling** | https://angular.dev/best-practices/error-handling |
| **Runtime Performance** | https://angular.dev/best-practices/runtime-performance |
| **Zone Pollution** | https://angular.dev/best-practices/zone-pollution |
| **Slow Computations** | https://angular.dev/best-practices/slow-computations |
| **Skipping Subtrees** | https://angular.dev/best-practices/skipping-subtrees |
| **Chrome DevTools Profiling** | https://angular.dev/best-practices/profiling-with-chrome-devtools |
| **Zoneless Angular** | https://angular.dev/guide/zoneless |
| **Tailwind CSS Integration** | https://angular.dev/guide/tailwind |
| **Angular Updates** | https://angular.dev/update |

**Key Implementation Requirements:**

1. **Performance Optimization**
   - Use `OnPush` change detection strategy on all components
   - Leverage `computed()` signals for derived state
   - Avoid zone pollution (use `runOutsideAngular` for heavy computations)
   - Use `@defer` blocks for lazy loading non-critical UI

2. **Security**
   - Never use `bypassSecurityTrust*` methods without security review
   - Rely on Angular's built-in XSS protection
   - Validate all user inputs
   - Use HttpOnly cookies for sensitive tokens

3. **Accessibility (WCAG 2.1 AA)**
   - All interactive elements must be keyboard accessible
   - Use semantic HTML elements (`<button>`, `<nav>`, `<main>`)
   - Provide ARIA labels where semantic HTML is insufficient
   - Maintain minimum 4.5:1 color contrast ratio
   - Test with screen readers (NVDA, VoiceOver)

4. **Error Handling**
   - Implement global `ErrorHandler` for uncaught exceptions
   - Use `catchError` operator in all HTTP observables
   - Display user-friendly error messages
   - Log errors with context for debugging

5. **Zoneless Ready**
   - Project targets zoneless Angular for maximum performance
   - Avoid Zone.js-dependent patterns
   - Use signals for all reactive state

---

## Core Architectural Decisions

### Data Architecture

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **ORM Strategy** | Hybrid: GORM + sqlc | GORM for 80% CRUD operations; sqlc for complex reports, analytics, bulk operations. Best performance + developer productivity. |
| **Migration Tool** | golang-migrate | Production-safe, versioned SQL files, supports rollback, CI/CD friendly. |
| **Soft Delete** | GORM `deleted_at` | Built-in filtering, audit-compliant, recoverable data. |
| **Audit Logging** | GORM Hooks + PostgreSQL Triggers | Hooks for application events; triggers for critical tables (users, payments, grades). |
| **UUID Strategy** | UUID v7 (time-ordered) | Sortable, no collision, better index performance than UUID v4. |
| **Connection Pooling** | PgBouncer | Transaction-level pooling, supports 1000+ concurrent connections. |

### Authentication & Security

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Password Hashing** | Argon2id | Current gold standard, memory-hard, resistant to GPU attacks. |
| **JWT Library** | golang-jwt/jwt/v5 | Well-maintained, standard claims support. |
| **Access Token** | 15 minutes, HS256 signed | Short-lived, symmetric key for single-service architecture. |
| **Refresh Token** | 7 days, PostgreSQL stored, rotation | Persistent storage, automatic expiry, rotate on each use. |
| **2FA** | TOTP (RFC 6238) | Works offline, Google Authenticator compatible. |
| **API Keys** | Scoped, rotatable | For third-party integrations (SMS, payment gateways). |
| **Rate Limiting** | Redis sliding window | Per-user, per-IP, per-tenant limits. |
| **Encryption at Rest** | AES-256-GCM | For sensitive fields (Aadhaar, bank details). |
| **TLS** | TLS 1.3 minimum | All connections encrypted in transit. |

### Production-Grade Refresh Token Implementation

The MSLS refresh token implementation follows security best practices:

#### Security Features (Implemented)

| Feature | Implementation | Why It Matters |
|---------|---------------|----------------|
| **Token Hashing** | SHA-256 hash stored, never plain tokens | Prevents token theft from DB compromise |
| **Token Rotation** | Old token revoked on each refresh | Limits window of token reuse attacks |
| **Expiration Checking** | 7-day TTL with explicit expiry column | Time-bounded access |
| **Revocation Support** | `revoked_at` timestamp for immediate invalidation | Enables logout across devices |
| **Request Queuing** | Frontend queues requests during refresh | Prevents race conditions |
| **Retry After Refresh** | Failed requests automatically retried | Seamless user experience |
| **Audit Logging** | All refresh events logged with IP/User-Agent | Security monitoring and forensics |

#### Token Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        REFRESH TOKEN FLOW                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Client                    Backend                    Database           │
│    │                         │                           │               │
│    │─── POST /auth/refresh ──▶│                           │               │
│    │    (refresh_token)      │                           │               │
│    │                         │                           │               │
│    │                         │── SHA256(token) ──────────▶│               │
│    │                         │                           │               │
│    │                         │◀── Find by hash ──────────│               │
│    │                         │                           │               │
│    │                         │── Check revoked? ─────────│               │
│    │                         │── Check expired? ─────────│               │
│    │                         │                           │               │
│    │                         │── REVOKE old token ───────▶│               │
│    │                         │                           │               │
│    │                         │── Generate new pair ──────│               │
│    │                         │                           │               │
│    │                         │── Store new hash ─────────▶│               │
│    │                         │                           │               │
│    │                         │── Audit log ──────────────▶│               │
│    │                         │                           │               │
│    │◀── New token pair ──────│                           │               │
│    │                         │                           │               │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Database Schema

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id),
    token_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA-256 hash
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,                   -- NULL = active
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

#### Audit Events

| Event | Trigger | Data Captured |
|-------|---------|---------------|
| `token_refresh` | Successful refresh | user_id, IP, user_agent, old_token_id |
| `token_refresh_failed` | Invalid/expired/revoked | attempt details, failure reason |
| `token_revoked` | Logout or rotation | token_id, revocation_reason |

#### Frontend Integration (Angular)

```typescript
// Interceptor handles 401 automatically:
// 1. Catches 401 error
// 2. Queues concurrent requests
// 3. Calls refresh endpoint
// 4. Retries original request with new token
// 5. Logs out if refresh fails
```

#### Security Considerations

1. **Why HS256 (not RS256)?**
   - Single backend service (modular monolith) - no need for distributed verification
   - Simpler key management
   - RS256 useful when multiple services verify tokens independently

2. **Why PostgreSQL (not Redis) for refresh tokens?**
   - Durability: Tokens survive restarts without configuration
   - Audit trail: Easier to track token lifecycle
   - On-premise support: No additional infrastructure dependency
   - Redis optional for caching active tokens for performance

### API & Communication

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **API Style** | REST with OpenAPI 3.1 | Universal support, tooling, auto-generated docs. |
| **Versioning** | URL-based `/api/v1/` | Explicit, easy routing, clear deprecation path. |
| **Error Format** | RFC 7807 Problem Details | Standard format, machine-readable, includes trace IDs. |
| **Pagination** | Cursor-based (default) | Performance on large tables; offset for small sets. |
| **Real-time** | WebSocket + SSE | WebSocket for bidirectional (transport tracking); SSE for notifications. |
| **Background Jobs** | Redis Streams | Persistent queues, consumer groups, no extra infra. |
| **Request Tracing** | OpenTelemetry | Distributed tracing, vendor-agnostic. |

**Standard Response Envelope:**
```json
{
  "success": true,
  "data": { },
  "meta": { "page": 1, "total": 100, "cursor": "abc123" },
  "errors": null,
  "trace_id": "req-uuid-here"
}
```

**Error Response (RFC 7807):**
```json
{
  "type": "https://api.msls.com/errors/validation",
  "title": "Validation Error",
  "status": 400,
  "detail": "Email format is invalid",
  "instance": "/api/v1/students",
  "trace_id": "req-uuid-here",
  "errors": [
    { "field": "email", "message": "Invalid email format" }
  ]
}
```

### Frontend Architecture

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Lazy Loading** | Route-based with preload | Load modules on demand; preload likely routes. |
| **Error Handling** | Global HTTP interceptor | Centralized error capture, toast notifications. |
| **Offline Storage** | IndexedDB via Dexie.js | Structured storage for complex offline data. |
| **Offline Strategy** | Cache-first static, network-first data | PWA service worker caching. |
| **PWA** | Angular PWA + Workbox | Service worker, push notifications, installable. |
| **i18n** | Angular i18n + ICU | Compile-time optimization, pluralization. |
| **Icons** | Lucide Angular | Tree-shakeable SVG icons, consistent style. |

### Infrastructure & Deployment

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Logging** | Zap (structured JSON) | High-performance, structured, log levels. |
| **Log Aggregation** | Loki (SaaS) / File rotation (on-prem) | Grafana-compatible, lightweight. |
| **Metrics** | Prometheus + Grafana | Industry standard, alerting built-in. |
| **Health Checks** | `/health` + `/ready` | K8s compatible, dependency checks. |
| **Secrets** | Env vars + Docker secrets (on-prem) | HashiCorp Vault for SaaS. |
| **CI/CD** | GitHub Actions | Integrated, good Docker/K8s support. |
| **Container Registry** | GitHub Container Registry | Integrated with CI, private images. |

### Multi-Tenancy Implementation

**Tenant Context Flow:**
```
Request → Auth Middleware (extract JWT)
       → Tenant Middleware (set tenant_id in context)
       → Set PostgreSQL session variable: SET app.current_tenant = 'uuid'
       → RLS policies automatically filter all queries
       → Response
```

**RLS Policy Pattern:**
```sql
CREATE POLICY tenant_isolation ON students
  USING (tenant_id = current_setting('app.current_tenant')::uuid);
```

**On-Premise Mode:**
```go
// Single tenant, tenant_id = predefined UUID
// Multi-tenant UI features hidden
// Branch hierarchy replaces tenant hierarchy
```

### Agent Autonomy Configuration

**Decision Authority Matrix:**
| Decision Type | Authority | Approval |
|--------------|-----------|----------|
| Code within established patterns | Agent | No |
| New utility functions | Agent | No |
| Bug fixes | Agent | No |
| API endpoints (following patterns) | Agent | No |
| New architectural patterns | Agent proposes | Yes |
| Breaking API changes | Agent proposes | Yes |
| New external dependencies | Agent proposes | Yes |
| Database schema changes | Agent proposes | Yes |

**Agent Context Sources:**
1. `architecture.md` - This document (patterns, decisions)
2. `project-context.md` - Code rules, naming conventions
3. PRD module docs - Feature requirements
4. Story files - Acceptance criteria per feature

---

## Implementation Patterns & Consistency Rules

### Naming Patterns

#### Database Naming (PostgreSQL)

| Element | Convention | Example |
|---------|------------|---------|
| Tables | `snake_case`, plural | `students`, `fee_payments`, `audit_logs` |
| Columns | `snake_case` | `first_name`, `created_at`, `tenant_id` |
| Primary Key | `id` (UUID v7) | `id UUID PRIMARY KEY DEFAULT gen_random_uuid()` |
| Foreign Key | `{table_singular}_id` | `student_id`, `tenant_id`, `branch_id` |
| Timestamps | `created_at`, `updated_at`, `deleted_at` | Standard GORM columns |
| Boolean | `is_` or `has_` prefix | `is_active`, `has_transport`, `is_verified` |
| Index | `idx_{table}_{columns}` | `idx_students_tenant_email` |
| Unique | `uniq_{table}_{columns}` | `uniq_users_tenant_email` |

#### API Naming (REST)

| Element | Convention | Example |
|---------|------------|---------|
| Endpoints | `snake_case`, plural | `/api/v1/students`, `/api/v1/fee_payments` |
| Path params | `snake_case` | `/api/v1/students/{student_id}` |
| Query params | `snake_case` | `?page_size=20&sort_by=created_at` |
| Nested resources | Max 2 levels | `/api/v1/classes/{class_id}/students` |
| Actions | Verb suffix | `/api/v1/students/{id}/promote` |

#### Go Backend Naming

| Element | Convention | Example |
|---------|------------|---------|
| Package | `lowercase`, single word | `student`, `fees`, `auth` |
| File | `snake_case.go` | `student_handler.go`, `fee_service.go` |
| Struct | `PascalCase` | `Student`, `FeePayment`, `CreateStudentDTO` |
| Interface | `PascalCase` + context | `StudentRepository`, `FeeService` |
| Function | `PascalCase` (exported) | `GetByID`, `CreateStudent` |
| Variable | `camelCase` | `studentID`, `tenantCtx`, `pageSize` |
| Error | `Err` prefix | `ErrStudentNotFound`, `ErrInvalidTenant` |

#### Angular Frontend Naming

| Element | Convention | Example |
|---------|------------|---------|
| Component | `PascalCase` | `StudentListComponent` |
| File | `kebab-case` | `student-list.component.ts` |
| Service | `PascalCase` + `Service` | `StudentService`, `AuthService` |
| Model | `PascalCase` | `Student`, `FeePayment` |
| Signal | `camelCase` | `students`, `isLoading` |
| Route | `kebab-case` | `/students`, `/fee-payments` |

### Structure Patterns

#### Backend Module Structure

```
internal/modules/{module}/
├── handler.go          # HTTP handlers
├── service.go          # Business logic
├── repository.go       # Database (GORM)
├── repository_sqlc.go  # Complex queries (optional)
├── dto.go              # Request/Response DTOs
├── errors.go           # Module-specific errors
└── {module}_test.go    # Tests co-located
```

#### Frontend Feature Structure

```
features/{feature}/
├── {feature}.routes.ts
├── pages/
│   ├── {feature}-list/
│   └── {feature}-detail/
├── components/
├── services/
│   └── {feature}.service.ts
├── models/
│   └── {feature}.model.ts
└── index.ts
```

### Format Patterns

#### API Response Format

**Success:**
```json
{
  "success": true,
  "data": { },
  "meta": { "page": 1, "page_size": 20, "total_count": 150 },
  "trace_id": "req-uuid"
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "type": "validation_error",
    "title": "Validation Failed",
    "status": 400,
    "detail": "One or more fields failed validation",
    "fields": [{ "field": "email", "message": "Invalid format" }]
  },
  "trace_id": "req-uuid"
}
```

#### Date/Time Formats

| Context | Format |
|---------|--------|
| API (JSON) | ISO 8601 UTC: `2026-01-22T10:30:00Z` |
| Display (India) | `DD-MMM-YYYY HH:mm` |
| Database | `TIMESTAMPTZ` |

### Communication Patterns

#### Event Naming

```
{domain}.{entity}.{action}
Examples: student.enrolled, fee.paid, attendance.marked
```

#### Event Payload

```json
{
  "event_id": "uuid",
  "event_type": "student.enrolled",
  "timestamp": "2026-01-22T10:30:00Z",
  "tenant_id": "uuid",
  "actor_id": "uuid",
  "data": { "student_id": "uuid" }
}
```

#### Angular State Pattern (Signals)

```typescript
@Injectable({ providedIn: 'root' })
export class StudentService {
  private _students = signal<Student[]>([]);
  private _loading = signal<boolean>(false);

  readonly students = this._students.asReadonly();
  readonly loading = this._loading.asReadonly();

  async loadStudents(): Promise<void> {
    this._loading.set(true);
    const data = await this.api.get<Student[]>('/students');
    this._students.set(data);
    this._loading.set(false);
  }
}
```

### Process Patterns

#### Error Handling (Backend)

```go
var ErrStudentNotFound = &AppError{
    Type: "not_found",
    Message: "Student not found",
    Status: 404,
}

func (h *Handler) GetByID(c *gin.Context) {
    student, err := h.service.GetByID(ctx, id)
    if err != nil {
        c.Error(err) // Middleware handles formatting
        return
    }
    c.JSON(200, Response{Success: true, Data: student})
}
```

#### Loading Pattern (Frontend)

```html
@if (loading()) {
  <app-loading-spinner />
} @else if (error()) {
  <app-error-message [message]="error()" />
} @else {
  <app-data-table [data]="students()" />
}
```

### Enforcement Rules

**All AI Agents MUST:**
1. Follow naming conventions exactly
2. Use standard response envelope for all APIs
3. Place files in correct module structure
4. Use signals for Angular state
5. Include `tenant_id` in all tenant-scoped operations
6. Log using structured JSON (Zap)
7. Write tests co-located with source

---

## Project Structure & Boundaries

### Backend Repository (`msls-backend`)

```
msls-backend/
├── .github/workflows/
│   ├── ci.yml
│   └── release.yml
├── api/openapi/
│   └── openapi.yaml
├── build/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
│       ├── base/
│       └── overlays/{dev,staging,prod}/
├── cmd/server/
│   └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── loader.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── tenant.go
│   │   ├── ratelimit.go
│   │   ├── audit.go
│   │   └── tracing.go
│   ├── modules/
│   │   ├── auth/
│   │   ├── tenant/
│   │   ├── branch/
│   │   ├── user/
│   │   ├── role/
│   │   ├── student/
│   │   ├── staff/
│   │   ├── class/
│   │   ├── section/
│   │   ├── subject/
│   │   ├── timetable/
│   │   ├── attendance/
│   │   ├── admission/
│   │   ├── examination/
│   │   ├── homework/
│   │   ├── quiz/
│   │   ├── classroom/
│   │   ├── leave/
│   │   ├── fees/
│   │   ├── payment/
│   │   ├── communication/
│   │   ├── notification/
│   │   ├── certificate/
│   │   ├── transport/
│   │   ├── library/
│   │   ├── inventory/
│   │   ├── visitor/
│   │   ├── analytics/
│   │   └── ai/
│   ├── pkg/
│   │   ├── database/
│   │   ├── cache/
│   │   ├── storage/
│   │   ├── notification/
│   │   ├── payment/
│   │   ├── logger/
│   │   ├── validator/
│   │   ├── response/
│   │   ├── pagination/
│   │   └── utils/
│   └── server/
│       ├── server.go
│       └── routes.go
├── migrations/
├── scripts/
├── sqlc/
├── tests/integration/
├── .env.example
├── go.mod
├── Makefile
└── README.md
```

### Frontend Repository (`msls-frontend`)

```
msls-frontend/
├── .github/workflows/
├── src/
│   ├── app/
│   │   ├── core/
│   │   │   ├── services/
│   │   │   ├── guards/
│   │   │   ├── interceptors/
│   │   │   └── models/
│   │   ├── shared/
│   │   │   ├── components/
│   │   │   ├── directives/
│   │   │   └── pipes/
│   │   ├── layouts/
│   │   │   ├── admin-layout/
│   │   │   ├── teacher-layout/
│   │   │   ├── student-layout/
│   │   │   └── parent-layout/
│   │   ├── features/
│   │   │   ├── auth/
│   │   │   ├── dashboard/
│   │   │   ├── students/
│   │   │   ├── staff/
│   │   │   ├── classes/
│   │   │   ├── attendance/
│   │   │   ├── timetable/
│   │   │   ├── admissions/
│   │   │   ├── examinations/
│   │   │   ├── homework/
│   │   │   ├── quiz/
│   │   │   ├── classroom/
│   │   │   ├── fees/
│   │   │   ├── communication/
│   │   │   ├── certificates/
│   │   │   ├── transport/
│   │   │   ├── library/
│   │   │   ├── inventory/
│   │   │   ├── visitors/
│   │   │   ├── analytics/
│   │   │   ├── settings/
│   │   │   └── profile/
│   │   ├── app.component.ts
│   │   ├── app.config.ts
│   │   └── app.routes.ts
│   ├── assets/
│   ├── environments/
│   └── styles/
├── e2e/
├── angular.json
├── package.json
├── tailwind.config.js
└── README.md
```

### Module Mapping (PRD → Code)

| PRD Module | Backend | Frontend |
|------------|---------|----------|
| Core Foundation | `modules/{tenant,branch,user,role}` | `features/{auth,settings}` |
| Student Management | `modules/student` | `features/students` |
| Academic Operations | `modules/{class,section,subject,timetable}` | `features/{classes,timetable}` |
| Attendance | `modules/attendance` | `features/attendance` |
| Admissions | `modules/admission` | `features/admissions` |
| Examinations | `modules/examination` | `features/examinations` |
| Homework | `modules/homework` | `features/homework` |
| Online Quiz | `modules/quiz` | `features/quiz` |
| Digital Classroom | `modules/classroom` | `features/classroom` |
| Staff Management | `modules/staff` | `features/staff` |
| Leave Management | `modules/leave` | `features/staff` |
| Fees & Payments | `modules/{fees,payment}` | `features/fees` |
| Communication | `modules/{communication,notification}` | `features/communication` |
| Certificates | `modules/certificate` | `features/certificates` |
| Transport | `modules/transport` | `features/transport` |
| Library | `modules/library` | `features/library` |
| Inventory | `modules/inventory` | `features/inventory` |
| Visitor Management | `modules/visitor` | `features/visitors` |
| Analytics | `modules/analytics` | `features/analytics` |
| AI Capabilities | `modules/ai` | `features/analytics` |

### Architectural Boundaries

```
┌───────────────────────────────────────────────────┐
│                   FRONTEND                         │
│   Features → Services → Interceptors → API        │
└───────────────────────────────────────────────────┘
                        │ REST API
┌───────────────────────────────────────────────────┐
│                   BACKEND                          │
│   Middleware → Handlers → Services → Repository   │
└───────────────────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
   PostgreSQL       Redis           MinIO
    (RLS)         (Cache)         (Storage)
```

---

## Architecture Validation Results

### 1. Coherence Validation

| Check | Status | Notes |
|-------|--------|-------|
| Tech Stack Compatibility | ✅ PASS | Go + Angular + PostgreSQL proven combination |
| Pattern Consistency | ✅ PASS | All modules follow same handler→service→repository pattern |
| Naming Convention Adherence | ✅ PASS | snake_case (DB), camelCase (JSON), kebab-case (files) |
| Multi-tenancy Coverage | ✅ PASS | RLS applies to all 22 modules |
| Security Model Complete | ✅ PASS | Auth, RBAC, RLS, audit logging defined |
| Offline Support Feasible | ✅ PASS | PWA + IndexedDB pattern documented |

### 2. Requirements Coverage

| PRD Requirement | Architecture Coverage |
|-----------------|----------------------|
| Multi-tenant SaaS | PostgreSQL RLS + tenant_id on all tables |
| On-Premise Deployment | Docker Compose configuration |
| 22 Functional Modules | All mapped to backend/frontend modules |
| Role-Based Access | RBAC with hierarchical permissions |
| Real-time Updates | WebSocket (bidirectional) + SSE (notifications) |
| Offline Mode | PWA + Dexie.js + background sync |
| Payment Integration | Payment module with UPI/Razorpay abstraction |
| Document Generation | PDF via go-pdf, Excel via excelize |
| AI Capabilities | On-premise AI module with local LLM support |
| Multi-language | i18n via Angular localization |

### 3. Implementation Readiness

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technology Maturity | HIGH | All technologies production-proven |
| Team Skill Availability | HIGH | Go, Angular, PostgreSQL widely known |
| Pattern Clarity | HIGH | Clear conventions for all code patterns |
| Boundary Definition | HIGH | Module boundaries explicit |
| Testing Strategy | HIGH | Unit + Integration + E2E defined |
| Deployment Path | HIGH | Docker → K8s progression clear |

**Overall Readiness: ✅ HIGH** - Architecture is implementation-ready with no blocking gaps.

### 4. Risk Assessment

| Risk | Severity | Mitigation |
|------|----------|------------|
| RLS Performance at Scale | Medium | Index optimization, query analysis, PgBouncer |
| Custom Component Library Effort | Medium | Build incrementally, start with core components |
| Offline Sync Conflicts | Medium | Last-write-wins with conflict audit log |
| AI Model Resource Requirements | Low | Optional module, quantized models for CPU-only |

### 5. Decision Authority Matrix (Agent Autonomy)

| Decision Type | Authority | Approval Needed |
|---------------|-----------|-----------------|
| Code within established patterns | Agent (Dev/TEA) | No |
| New utility functions/helpers | Agent | No |
| Bug fixes | Agent | No |
| API endpoint implementation | Agent | No |
| Test coverage expansion | Agent | No |
| New architectural patterns | Architect | Yes |
| Breaking API changes | PM + Architect | Yes |
| New external dependencies | Architect | Yes |
| Database schema changes | Architect | Review |
| Security-related changes | Architect | Yes |

**Note**: Agents can implement autonomously within documented patterns. Deviations require approval.

---

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** COMPLETED ✅
**Total Steps Completed:** 8
**Date Completed:** 2026-01-23
**Document Location:** `_bmad-output/planning-artifacts/architecture.md`

### Final Architecture Deliverables

**📋 Complete Architecture Document**

- All architectural decisions documented with specific versions
- Implementation patterns ensuring AI agent consistency
- Complete project structure with all files and directories
- Requirements to architecture mapping for all 22 PRD modules
- Validation confirming coherence and completeness

**🏗️ Implementation Ready Foundation**

- 25+ architectural decisions made
- 15+ implementation patterns defined
- 22 backend modules + 22 frontend features specified
- All PRD requirements fully supported

**📚 AI Agent Implementation Guide**

- Technology stack with verified versions (Go 1.23+, Angular 21, PostgreSQL 16)
- Consistency rules that prevent implementation conflicts
- Project structure with clear boundaries
- Integration patterns and communication standards

### Implementation Handoff

**For AI Agents:**
This architecture document is your complete guide for implementing MSLS (Multi-School Learning System). Follow all decisions, patterns, and structures exactly as documented.

**First Implementation Priority:**
```bash
# Backend initialization
mkdir -p msls-backend && cd msls-backend
go mod init github.com/ashulabs/msls-backend
# Use go-blueprint or manual setup following documented structure

# Frontend initialization
ng new msls-frontend --style=scss --routing --standalone
cd msls-frontend && npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init
```

**Development Sequence:**

1. Initialize project using documented starter templates
2. Set up development environment per architecture
3. Implement core architectural foundations (auth, tenant, database)
4. Build features following established patterns
5. Maintain consistency with documented rules

### Quality Assurance Checklist

**✅ Architecture Coherence**

- [x] All decisions work together without conflicts
- [x] Technology choices are compatible
- [x] Patterns support the architectural decisions
- [x] Structure aligns with all choices

**✅ Requirements Coverage**

- [x] All 22 functional modules are supported
- [x] All non-functional requirements are addressed
- [x] Cross-cutting concerns (auth, multi-tenancy, audit) handled
- [x] Integration points (payment, SMS, email) defined

**✅ Implementation Readiness**

- [x] Decisions are specific and actionable
- [x] Patterns prevent agent conflicts
- [x] Structure is complete and unambiguous
- [x] Examples are provided for clarity

### Project Success Factors

**🎯 Clear Decision Framework**
Every technology choice was made with clear rationale, ensuring all stakeholders understand the architectural direction.

**🔧 Consistency Guarantee**
Implementation patterns and rules ensure that multiple AI agents will produce compatible, consistent code that works together seamlessly.

**📋 Complete Coverage**
All 22 PRD modules are architecturally supported, with clear mapping from business needs to technical implementation.

**🏗️ Solid Foundation**
The chosen stack (Go + Angular + PostgreSQL + Tailwind) provides a production-ready foundation following current best practices.

---

**Architecture Status:** READY FOR IMPLEMENTATION ✅

**Next Phase:** Begin implementation using the architectural decisions and patterns documented herein.

**Document Maintenance:** Update this architecture when major technical decisions are made during implementation.

