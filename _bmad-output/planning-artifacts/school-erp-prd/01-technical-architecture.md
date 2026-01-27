# 01 - Technical Architecture

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: Foundation

---

## 1. Overview

This document defines the technical architecture for the School ERP platform, covering technology stack, deployment models, multi-tenancy strategy, and infrastructure design.

---

## 2. Technology Stack

### 2.1 Backend

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Language** | Go (Golang) | High performance, single binary deployment, excellent concurrency |
| **Framework** | Gin / Echo | Lightweight, fast HTTP framework |
| **ORM** | GORM / sqlc | Type-safe database operations |
| **API Style** | REST | Widely supported, simple integration |

### 2.2 Frontend

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Web Framework** | Angular | Enterprise-grade, TypeScript-first, strong ecosystem |
| **UI Library** | Angular Material / PrimeNG | Rich component library |
| **State Management** | NgRx / Signals | Predictable state management |
| **Mobile (Hybrid)** | Ionic | Code sharing with web, PWA support |
| **Desktop** | Electron | Cross-platform desktop app |

### 2.3 Database

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Primary DB** | PostgreSQL | ACID compliance, RLS support, JSON capabilities |
| **Caching** | Redis | Session management, rate limiting, queues |
| **Search** | PostgreSQL Full-Text / Meilisearch | Fast search for students, staff |
| **File Storage** | Local FS / MinIO | On-premise friendly object storage |

### 2.4 Infrastructure

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Containerization** | Docker | Consistent deployment |
| **Orchestration** | Docker Compose (on-prem) / K8s (SaaS) | Scalable deployment |
| **Reverse Proxy** | Nginx / Caddy | SSL termination, load balancing |
| **Message Queue** | Redis Streams / NATS | Async processing |

---

## 3. Multi-Tenancy Architecture

### 3.1 Strategy: Row-Level Isolation with PostgreSQL RLS

**Chosen Approach**: Row-Level Security (RLS) with `tenant_id` column on all tenant-scoped tables.

```
┌─────────────────────────────────────────────────────────────┐
│                    APPLICATION LAYER                         │
│   ┌─────────────────────────────────────────────────────┐   │
│   │  Every request extracts tenant_id from JWT/session  │   │
│   │  Sets: SET app.current_tenant = 'tenant_uuid'       │   │
│   └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    DATABASE LAYER                            │
│   ┌─────────────────────────────────────────────────────┐   │
│   │  PostgreSQL Row-Level Security (RLS)                │   │
│   │                                                      │   │
│   │  CREATE POLICY tenant_isolation ON students         │   │
│   │    USING (tenant_id = current_setting('app.tenant'))│   │
│   └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 Why RLS Over Schema-per-Tenant

| Factor | Schema-per-Tenant | Row-Level (RLS) | Decision |
|--------|-------------------|-----------------|----------|
| Code complexity | High | Low | **RLS** |
| Migrations | Run N times | Run once | **RLS** |
| On-prem compatibility | Awkward | Natural | **RLS** |
| Cross-tenant queries | Hard | Easy | **RLS** |
| Connection pooling | Complex | Simple | **RLS** |
| Data isolation | Physical | Logical (RLS enforced) | Schema |

### 3.3 Tenant-Scoped vs Global Tables

**Tenant-Scoped Tables** (have `tenant_id`):
- students, staff, classes, sections
- attendance, fees, payments
- admissions, examinations
- All operational data

**Global Tables** (no `tenant_id`):
- tenants (tenant registry)
- subscription_plans
- system_configurations
- audit_logs (with tenant reference)

### 3.4 On-Premise Mode

Same codebase, single tenant:
```go
// On-prem: tenant_id is always "default" or UUID "1"
// Multi-tenant features hidden in UI
// Branch hierarchy replaces tenant hierarchy
```

---

## 4. Service Architecture

### 4.1 Modular Monolith (Phase 1)

Start with a modular monolith, designed for future extraction:

```
┌─────────────────────────────────────────────────────────────┐
│                     API GATEWAY                              │
│           (Authentication, Rate Limiting, Routing)          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    CORE APPLICATION                          │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │  Auth   │ │ Student │ │  Staff  │ │  Fees   │           │
│  │ Module  │ │ Module  │ │ Module  │ │ Module  │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │  Exam   │ │Attendance│ │Timetable│ │  Comms  │           │
│  │ Module  │ │ Module  │ │ Module  │ │ Module  │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    DATA LAYER                                │
│    ┌──────────────┐    ┌─────────┐    ┌─────────────┐       │
│    │  PostgreSQL  │    │  Redis  │    │  File Store │       │
│    │   (Primary)  │    │ (Cache) │    │   (MinIO)   │       │
│    └──────────────┘    └─────────┘    └─────────────┘       │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 Module Boundaries

Each module is a Go package with:
- **Handlers**: HTTP handlers (routes)
- **Services**: Business logic
- **Repositories**: Database operations
- **Models**: Domain entities
- **DTOs**: Request/Response structures

```
/internal
  /modules
    /student
      handlers.go
      service.go
      repository.go
      models.go
      dto.go
    /staff
      ...
    /fees
      ...
```

---

## 5. API Design

### 5.1 REST API Standards

**Base URL**: `/api/v1`

**Naming Conventions**:
- Resources: plural nouns (`/students`, `/classes`)
- Actions: HTTP verbs (GET, POST, PUT, DELETE)
- Nested resources: `/classes/{id}/students`

**Response Format**:
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150
  },
  "errors": null
}
```

### 5.2 API Versioning

- URL-based versioning: `/api/v1/`, `/api/v2/`
- Major version changes only for breaking changes
- Deprecation notice 6 months before removal

### 5.3 Rate Limiting

| Endpoint Type | Limit |
|---------------|-------|
| Public (login) | 10/min per IP |
| Authenticated | 100/min per user |
| Bulk operations | 10/min per tenant |
| Reports | 5/min per user |

---

## 6. Authentication & Authorization

### 6.1 Authentication

**Strategy**: JWT with Refresh Tokens

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│   Login     │────▶│   Issue     │
│             │     │   Endpoint  │     │   JWT +     │
│             │◀────│             │◀────│   Refresh   │
└─────────────┘     └─────────────┘     └─────────────┘

JWT Payload:
{
  "sub": "user_uuid",
  "tenant_id": "tenant_uuid",
  "branch_id": "branch_uuid",
  "roles": ["teacher", "class_teacher"],
  "exp": 1234567890
}
```

**Token Lifecycle**:
- Access Token: 15 minutes
- Refresh Token: 7 days (stored in Redis)
- Refresh rotation on each use

### 6.2 Authorization (RBAC)

**Role Hierarchy**:
```
super_admin
  └── school_admin
        └── branch_admin
              ├── principal
              ├── accountant
              ├── teacher
              │     └── class_teacher
              ├── librarian
              └── transport_incharge
```

**Permission Structure**:
```
module:resource:action
Examples:
- student:profile:read
- student:profile:write
- fees:payment:collect
- exam:result:publish
```

---

## 7. Configuration System

### 7.1 Three-Layer Configuration

**Layer 1: Feature Flags**
```yaml
modules:
  admissions:
    enabled: true
  library:
    enabled: false
  hostel:
    enabled: false
  transport:
    enabled: true
```

**Layer 2: Custom Fields**
```yaml
custom_fields:
  student:
    - name: blood_group
      type: dropdown
      options: [A+, A-, B+, B-, O+, O-, AB+, AB-]
      required: false
    - name: bus_route
      type: reference
      reference: transport_routes
      required: false
```

**Layer 3: Workflow Configuration**
```yaml
admission_workflow:
  stages:
    - name: enquiry
      next: [document_submission, rejected]
    - name: document_submission
      next: [entrance_test, direct_admission]
    - name: entrance_test
      next: [approved, waitlist, rejected]
```

---

## 8. Deployment Architecture

### 8.1 SaaS Deployment

```
┌─────────────────────────────────────────────────────────────┐
│                      LOAD BALANCER                           │
│                    (AWS ALB / Nginx)                        │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│   API Pod 1   │     │   API Pod 2   │     │   API Pod N   │
│   (Go App)    │     │   (Go App)    │     │   (Go App)    │
└───────────────┘     └───────────────┘     └───────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    DATA LAYER                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ PostgreSQL  │  │    Redis    │  │       MinIO         │  │
│  │  (Primary   │  │   Cluster   │  │   (File Storage)    │  │
│  │  + Replica) │  │             │  │                     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 8.2 On-Premise Deployment

```
┌─────────────────────────────────────────────────────────────┐
│                   SCHOOL SERVER                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                 Docker Compose                       │    │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐   │    │
│  │  │  Nginx  │ │ Go App  │ │PostgreSQL│ │  Redis  │   │    │
│  │  │ :80/443 │ │  :8080  │ │  :5432  │ │  :6379  │   │    │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘   │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Local Storage: /data/files, /data/backups                  │
└─────────────────────────────────────────────────────────────┘
        │
        │ LAN Access
        ▼
┌─────────────────┐
│  School Network │
│  (No Internet)  │
└─────────────────┘
```

---

## 9. Security Architecture

### 9.1 Data Security

| Layer | Protection |
|-------|------------|
| Transit | TLS 1.3 for all connections |
| Rest | AES-256 encryption for sensitive fields |
| Database | PostgreSQL RLS, encrypted backups |
| Files | Encrypted storage, signed URLs |

### 9.2 Application Security

- Input validation on all endpoints
- SQL injection prevention (parameterized queries)
- XSS prevention (Content Security Policy)
- CSRF tokens for state-changing operations
- Rate limiting and DDoS protection

### 9.3 Audit Logging

```sql
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  user_id UUID NOT NULL,
  action VARCHAR(50) NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  resource_id UUID,
  old_values JSONB,
  new_values JSONB,
  ip_address INET,
  user_agent TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 10. Scalability Considerations

| Metric | Target | Strategy |
|--------|--------|----------|
| Concurrent users | 50,000+ | Horizontal pod scaling |
| Database connections | 1,000+ | PgBouncer connection pooling |
| File storage | 10TB+ | MinIO with sharding |
| API latency | < 200ms p95 | Redis caching, query optimization |

---

## 11. Related Documents

- [02-core-foundation.md](./02-core-foundation.md) - Institution setup, Users, RBAC
- [index.md](./index.md) - Main PRD index

---

**Next**: [02-core-foundation.md](./02-core-foundation.md)
