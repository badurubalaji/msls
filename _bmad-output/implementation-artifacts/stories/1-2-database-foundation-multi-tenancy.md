# Story 1.2: Database Foundation with Multi-Tenancy

**Epic:** 1 - Project Foundation & Design System
**Status:** ready-for-dev
**Priority:** Critical

## User Story

As a **system**,
I want **PostgreSQL database with Row-Level Security policies and multi-tenant support**,
So that **tenant data is completely isolated at the database level**.

## Acceptance Criteria

**Given** the database is initialized
**When** migrations are applied
**Then** the following foundation tables exist:
- `tenants` (id, name, slug, settings, status, created_at, updated_at)
- `branches` (id, tenant_id, name, code, address, settings, is_primary, status)
- `users` (id, tenant_id, email, phone, password_hash, status, 2fa_enabled)
**And** UUID v7 generation function `uuid_generate_v7()` is available
**And** RLS is enabled on all tables with tenant_id
**And** RLS policy uses `current_setting('app.tenant_id')::UUID`
**And** Indexes exist on tenant_id for all tables
**And** Audit columns (created_at, updated_at, created_by, updated_by) exist on all tables

## Technical Requirements

### Database Migrations Location
`msls-backend/migrations/`

### Migration Naming Convention
`YYYYMMDDHHMMSS_description.up.sql` and `YYYYMMDDHHMMSS_description.down.sql`

### Tables to Create

#### 1. Extension Setup (000001)
```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

#### 2. UUID v7 Function (000002)
```sql
CREATE OR REPLACE FUNCTION uuid_generate_v7()
RETURNS uuid AS $$
BEGIN
  -- UUID v7 implementation with timestamp + random
END;
$$ LANGUAGE plpgsql;
```

#### 3. Tenants Table (000003)
```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    settings JSONB DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### 4. Branches Table (000004)
```sql
CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    address JSONB DEFAULT '{}',
    settings JSONB DEFAULT '{}',
    is_primary BOOLEAN DEFAULT false,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_by UUID,
    UNIQUE(tenant_id, code)
);
```

#### 5. Users Table (000005)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    email VARCHAR(255),
    phone VARCHAR(20),
    password_hash VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    two_factor_enabled BOOLEAN DEFAULT false,
    two_factor_secret VARCHAR(255),
    email_verified_at TIMESTAMPTZ,
    phone_verified_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_by UUID,
    UNIQUE(tenant_id, email),
    UNIQUE(tenant_id, phone)
);
```

#### 6. RLS Setup (000006)
```sql
-- Enable RLS on tables
ALTER TABLE branches ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
CREATE POLICY tenant_isolation_branches ON branches
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_users ON users
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);
```

#### 7. Indexes (000007)
```sql
CREATE INDEX idx_branches_tenant_id ON branches(tenant_id);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
```

### GORM Models Location
`msls-backend/internal/pkg/database/models/`

### Model Files to Create
- `tenant.go`
- `branch.go`
- `user.go`
- `base.go` (common audit fields)

## Tasks

1. [ ] Create migrations directory structure
2. [ ] Create extension setup migration
3. [ ] Create UUID v7 function migration
4. [ ] Create tenants table migration
5. [ ] Create branches table migration
6. [ ] Create users table migration
7. [ ] Create RLS policies migration
8. [ ] Create indexes migration
9. [ ] Create GORM base model with audit fields
10. [ ] Create Tenant GORM model
11. [ ] Create Branch GORM model
12. [ ] Create User GORM model
13. [ ] Create database connection package
14. [ ] Add migration runner to Makefile
15. [ ] Test migrations up/down

## Definition of Done

- [ ] All migrations run successfully
- [ ] RLS policies work correctly (tested with different tenant contexts)
- [ ] GORM models have proper tags and relationships
- [ ] Indexes verified in database
- [ ] Down migrations properly rollback
