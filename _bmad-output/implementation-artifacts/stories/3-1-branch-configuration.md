# Story 3.1: Branch Configuration

**Epic:** 3 - School Setup & Admissions
**Status:** code-complete
**Priority:** High
**Estimated Effort:** Medium

---

## User Story

As a **super admin**,
I want **to configure branches within a tenant**,
So that **multi-branch schools can manage each location separately**.

---

## Acceptance Criteria

### AC1: Create Branch
**Given** a super admin is on tenant settings
**When** they add a new branch
**Then** they can enter: name, code, address, contact details
**And** they can set one branch as primary
**And** they can configure branch-specific settings (logo, timezone)
**And** branch is created and available for assignment

### AC2: List and Manage Branches
**Given** branches exist
**When** viewing the branch list
**Then** all branches are displayed with status
**And** each branch shows student/staff count (placeholder for now)
**And** branches can be activated/deactivated

### AC3: Edit Branch
**Given** a branch exists
**When** admin edits the branch
**Then** all fields can be updated
**And** changes are saved and reflected immediately

### AC4: Primary Branch
**Given** multiple branches exist
**When** setting a branch as primary
**Then** only one branch can be primary at a time
**And** previous primary is automatically unset

---

## Technical Requirements

### Backend (Go)

#### Database Schema

```sql
-- Migration: {timestamp}_create_branches_table.up.sql
CREATE TABLE IF NOT EXISTS branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(200) NOT NULL,
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'India',
    phone VARCHAR(20),
    email VARCHAR(255),
    logo_url VARCHAR(500),
    timezone VARCHAR(50) DEFAULT 'Asia/Kolkata',
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    UNIQUE(tenant_id, code)
);

-- RLS Policy
ALTER TABLE branches ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON branches
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

-- Indexes
CREATE INDEX idx_branches_tenant ON branches(tenant_id);
CREATE INDEX idx_branches_tenant_code ON branches(tenant_id, code);
CREATE INDEX idx_branches_tenant_primary ON branches(tenant_id, is_primary) WHERE is_primary = TRUE;
```

#### Module Structure

```
internal/modules/branch/
├── handler.go      # HTTP handlers
├── service.go      # Business logic
├── repository.go   # Database operations
├── dto.go          # Request/Response structs
├── entity.go       # Branch entity
└── branch_test.go  # Tests
```

#### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | /api/v1/branches | List all branches | branches:read |
| GET | /api/v1/branches/:id | Get branch by ID | branches:read |
| POST | /api/v1/branches | Create branch | branches:create |
| PUT | /api/v1/branches/:id | Update branch | branches:update |
| PATCH | /api/v1/branches/:id/primary | Set as primary | branches:update |
| PATCH | /api/v1/branches/:id/status | Toggle active status | branches:update |
| DELETE | /api/v1/branches/:id | Delete branch | branches:delete |

#### DTO Examples

```go
// CreateBranchRequest
type CreateBranchRequest struct {
    Code        string `json:"code" validate:"required,max=20"`
    Name        string `json:"name" validate:"required,max=200"`
    AddressLine1 string `json:"addressLine1" validate:"max=255"`
    AddressLine2 string `json:"addressLine2" validate:"max=255"`
    City        string `json:"city" validate:"max=100"`
    State       string `json:"state" validate:"max=100"`
    PostalCode  string `json:"postalCode" validate:"max=20"`
    Country     string `json:"country" validate:"max=100"`
    Phone       string `json:"phone" validate:"max=20"`
    Email       string `json:"email" validate:"omitempty,email,max=255"`
    Timezone    string `json:"timezone" validate:"max=50"`
    IsPrimary   bool   `json:"isPrimary"`
}

// BranchResponse
type BranchResponse struct {
    ID          string `json:"id"`
    Code        string `json:"code"`
    Name        string `json:"name"`
    AddressLine1 string `json:"addressLine1,omitempty"`
    AddressLine2 string `json:"addressLine2,omitempty"`
    City        string `json:"city,omitempty"`
    State       string `json:"state,omitempty"`
    PostalCode  string `json:"postalCode,omitempty"`
    Country     string `json:"country"`
    Phone       string `json:"phone,omitempty"`
    Email       string `json:"email,omitempty"`
    LogoURL     string `json:"logoUrl,omitempty"`
    Timezone    string `json:"timezone"`
    IsPrimary   bool   `json:"isPrimary"`
    IsActive    bool   `json:"isActive"`
    CreatedAt   string `json:"createdAt"`
    UpdatedAt   string `json:"updatedAt"`
}
```

### Frontend (Angular)

#### Feature Structure

```
src/app/features/admin/branches/
├── branches.routes.ts
├── branches.component.ts           # Main list page
├── branches.component.html
├── branches.component.scss
├── branch-form.component.ts        # Create/Edit form modal
├── branch-form.component.html
└── services/
    └── branch.service.ts
```

#### Models

```typescript
// src/app/features/admin/branches/models/branch.model.ts
export interface Branch {
  id: string;
  code: string;
  name: string;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;
  country: string;
  phone?: string;
  email?: string;
  logoUrl?: string;
  timezone: string;
  isPrimary: boolean;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateBranchRequest {
  code: string;
  name: string;
  addressLine1?: string;
  addressLine2?: string;
  city?: string;
  state?: string;
  postalCode?: string;
  country?: string;
  phone?: string;
  email?: string;
  timezone?: string;
  isPrimary?: boolean;
}

export type UpdateBranchRequest = Partial<CreateBranchRequest>;
```

#### UI Components Required

1. **Branches List Page** (`/admin/branches`)
   - Table with columns: Name, Code, City, Status, Primary, Actions
   - Search/filter functionality
   - "Add Branch" button
   - Action buttons: Edit, Set Primary, Toggle Status, Delete

2. **Branch Form Modal**
   - Fields: Code, Name, Address (multi-line), City, State, Postal Code, Country, Phone, Email, Timezone, Is Primary
   - Validation messages
   - Save/Cancel buttons

#### Route Configuration

```typescript
// Add to admin.routes.ts
{
  path: 'branches',
  loadComponent: () => import('./branches/branches.component').then(m => m.BranchesComponent),
  data: { permission: 'branches:read' }
}
```

#### Navigation

Add to admin sidebar:
```typescript
{
  label: 'Branches',
  icon: 'building-office-2',
  route: '/admin/branches',
  permission: 'branches:read'
}
```

---

## Out of Scope

- Logo upload functionality (future story)
- Branch-specific settings configuration UI
- Student/staff count display (requires those modules)

---

## Testing Requirements

### Backend
- Unit tests for service layer
- Integration tests for repository
- API endpoint tests

### Frontend
- Component unit tests
- Service tests with mocked HTTP

---

## Definition of Done

- [x] Backend: Migration created and applied (migration 000004_branches already exists)
- [x] Backend: All CRUD endpoints implemented
- [x] Backend: Unit tests passing
- [x] Backend: API documented (Swagger annotations added)
- [x] Frontend: Branch list page implemented
- [x] Frontend: Branch form modal implemented
- [x] Frontend: Service with all API calls
- [x] Frontend: Navigation link added
- [x] Frontend: Component tests passing
- [ ] Code reviewed and approved
- [ ] No console errors or warnings

## Backend Implementation Notes (2026-01-23)

**Files Created:**
- `/internal/services/branch/errors.go` - Branch-specific error definitions
- `/internal/services/branch/service.go` - Branch service with all business logic
- `/internal/services/branch/service_test.go` - Unit tests for branch service
- `/internal/handlers/branch/dto.go` - Request/Response DTOs
- `/internal/handlers/branch/handler.go` - HTTP handlers with Swagger annotations

**Files Modified:**
- `/cmd/api/main.go` - Added branch routes and service initialization
- `/internal/services/rbac/errors.go` - Added ModuleBranches constant
- `/internal/services/rbac/permission_service.go` - Added branch permissions to seed
- `/internal/services/rbac/role_service.go` - Added branch permissions to system roles

**API Endpoints Implemented:**
| Method | Endpoint | Permission |
|--------|----------|------------|
| GET | /api/v1/branches | branches:read |
| GET | /api/v1/branches/:id | branches:read |
| POST | /api/v1/branches | branches:create |
| PUT | /api/v1/branches/:id | branches:update |
| PATCH | /api/v1/branches/:id/primary | branches:update |
| PATCH | /api/v1/branches/:id/status | branches:update |
| DELETE | /api/v1/branches/:id | branches:delete |

**Permissions Added:**
- `branches:read` - View branch information
- `branches:create` - Create new branches
- `branches:update` - Update branch information
- `branches:delete` - Delete branches

## Frontend Implementation Notes (2026-01-24)

**Files Created:**
- `/msls-frontend/src/app/features/admin/branches/branch.model.ts` - Branch interface and request types
- `/msls-frontend/src/app/features/admin/branches/branch.service.ts` - BranchService with all API calls
- `/msls-frontend/src/app/features/admin/branches/branch.service.spec.ts` - Service unit tests
- `/msls-frontend/src/app/features/admin/branches/branches.component.ts` - Main list page with inline template
- `/msls-frontend/src/app/features/admin/branches/branches.component.spec.ts` - Component unit tests
- `/msls-frontend/src/app/features/admin/branches/branch-form.component.ts` - Create/Edit form modal component
- `/msls-frontend/src/app/features/admin/branches/branch-form.component.spec.ts` - Form component unit tests

**Files Modified:**
- `/msls-frontend/src/app/features/admin/admin.routes.ts` - Added branches route
- `/msls-frontend/src/app/layouts/nav-config.ts` - Added Branches navigation item under Admin

**Features Implemented:**
- Branch list page with search/filter functionality
- Create branch modal with all required fields
- Edit branch functionality
- Set branch as primary functionality
- Toggle branch active/inactive status
- Delete branch with confirmation modal
- Form validation with error messages
- Loading states and error handling
- Responsive table design with Tailwind CSS
- Angular Signals for reactive state management
