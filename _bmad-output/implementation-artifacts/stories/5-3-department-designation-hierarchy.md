# Story 5.3: Department & Designation Hierarchy

**Epic:** 5 - Staff Management
**Status:** dev-complete
**Priority:** P1 (Required for organizational structure)
**Story Points:** 5

---

## User Story

**As an** HR administrator,
**I want** to configure departments and designations,
**So that** organizational structure is properly defined.

---

## Acceptance Criteria

### AC1: Department Management

**Given** HR is on organization settings
**When** creating a department
**Then** they can enter: name, code, description
**And** they can select a department head (staff member)
**And** departments can be active/inactive
**And** department code must be unique within tenant

### AC2: Designation Management

**Given** HR is configuring designations
**When** creating a designation
**Then** they can enter: designation name, code, level (numeric for hierarchy)
**And** they can optionally link to a department
**And** they can set reporting designation (who this role reports to)
**And** designations can be active/inactive

### AC3: Staff Department Assignment

**Given** a staff member exists
**When** assigning department/designation
**Then** they can select department and designation
**And** assignment date is recorded
**And** previous assignments are tracked in history

### AC4: Department Listing & Filtering

**Given** HR is on department list
**When** viewing departments
**Then** they see: name, code, head, staff count, status
**And** they can filter by status (active/inactive)
**And** they can search by name or code

### AC5: Designation Listing & Hierarchy

**Given** HR is on designation list
**When** viewing designations
**Then** they see: name, code, level, department, reports to
**And** they can filter by department
**And** hierarchy levels are clearly displayed

---

## Technical Implementation

### Database Schema

```sql
-- Departments table
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    head_id UUID REFERENCES staff(id),

    is_active BOOLEAN DEFAULT true,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    UNIQUE(tenant_id, code)
);

-- Designations table
CREATE TABLE designations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    level INTEGER NOT NULL DEFAULT 1, -- Lower = more senior
    department_id UUID REFERENCES departments(id),
    reports_to_id UUID REFERENCES designations(id),

    is_active BOOLEAN DEFAULT true,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    UNIQUE(tenant_id, code)
);

-- Add department_id and designation_id to staff table
ALTER TABLE staff
    ADD COLUMN department_id UUID REFERENCES departments(id),
    ADD COLUMN designation_id UUID REFERENCES designations(id);

-- RLS Policies
ALTER TABLE departments ENABLE ROW LEVEL SECURITY;
ALTER TABLE designations ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_departments ON departments
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY tenant_isolation_designations ON designations
    USING (tenant_id = current_setting('app.tenant_id')::uuid);
```

### API Endpoints

```
# Departments
GET    /api/v1/departments              - List departments
POST   /api/v1/departments              - Create department
GET    /api/v1/departments/:id          - Get department
PUT    /api/v1/departments/:id          - Update department
DELETE /api/v1/departments/:id          - Delete department (soft)

# Designations
GET    /api/v1/designations             - List designations
POST   /api/v1/designations             - Create designation
GET    /api/v1/designations/:id         - Get designation
PUT    /api/v1/designations/:id         - Update designation
DELETE /api/v1/designations/:id         - Delete designation (soft)
```

### Frontend Components

```
src/app/features/organization/
├── organization.routes.ts
├── models/
│   ├── department.model.ts
│   └── designation.model.ts
├── services/
│   ├── department.service.ts
│   └── designation.service.ts
└── pages/
    ├── department-list/
    │   └── department-list.component.ts
    ├── department-form/
    │   └── department-form.component.ts
    ├── designation-list/
    │   └── designation-list.component.ts
    └── designation-form/
        └── designation-form.component.ts
```

---

## Tasks

### Backend Tasks

- [ ] BE-5.3.1: Create departments and designations migration
- [ ] BE-5.3.2: Create department module (model, repository, service, handler)
- [ ] BE-5.3.3: Create designation module (model, repository, service, handler)
- [ ] BE-5.3.4: Update staff module to include department/designation
- [ ] BE-5.3.5: Add permissions for departments and designations

### Frontend Tasks

- [ ] FE-5.3.1: Create organization feature module with routing
- [ ] FE-5.3.2: Create department models and service
- [ ] FE-5.3.3: Create department list and form components
- [ ] FE-5.3.4: Create designation models and service
- [ ] FE-5.3.5: Create designation list and form components
- [ ] FE-5.3.6: Update staff form to include department/designation selection

### Testing Tasks

- [ ] BE-5.3.6: Write unit tests for department and designation modules
- [ ] FE-5.3.7: Test all components

---

## Dependencies

- Story 5.1: Staff Profile Management (completed) - For staff head selection

## Notes

- Department head selection should show active staff members only
- Designation reports_to creates the hierarchy chain
- Level field helps with sorting and display (1 = highest, e.g., Principal)
