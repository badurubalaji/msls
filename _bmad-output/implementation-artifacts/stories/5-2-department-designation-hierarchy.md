# Story 5.2: Department & Designation Hierarchy

**Epic:** 5 - Staff Management
**Status:** completed
**Priority:** P0 (Foundation for staff management)
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
**And** they can assign a department head (staff member)
**And** departments can be set active/inactive
**And** department code must be unique within tenant

### AC2: Designation Management

**Given** designations are configured
**When** creating a designation
**Then** they can enter: designation name, level (1-10)
**And** they can optionally link to a department
**And** they can link a salary grade (for payroll integration)
**And** reporting hierarchy level determines seniority

### AC3: Department Listing

**Given** HR is on departments page
**When** viewing the list
**Then** they see: name, code, head, staff count, status
**And** they can filter by status (active/inactive)
**And** they can search by name or code

### AC4: Designation Listing

**Given** HR is on designations page
**When** viewing the list
**Then** they see: name, level, department, salary grade, status
**And** they can filter by department and status
**And** levels are sorted in hierarchy order

### AC5: Organization Chart View

**Given** staff members exist with departments
**When** viewing org chart
**Then** hierarchical view shows departments
**And** department heads are shown at top of each department
**And** staff can be filtered by department
**And** chart is collapsible/expandable

---

## Technical Implementation

### Database Schema

```sql
-- Migration: 000034_departments_designations.up.sql

-- Departments table
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID REFERENCES branches(id),  -- NULL means applies to all branches

    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    head_id UUID,  -- Will reference staff(id) once staff table exists

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- RLS Policy
ALTER TABLE departments ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON departments
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

-- Indexes
CREATE INDEX idx_departments_tenant ON departments(tenant_id);
CREATE INDEX idx_departments_branch ON departments(branch_id);
CREATE UNIQUE INDEX uniq_departments_code ON departments(tenant_id, code);

-- Salary Grades table (for payroll integration)
CREATE TABLE salary_grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(50) NOT NULL,           -- e.g., "Grade A", "Level 1"
    code VARCHAR(20) NOT NULL,
    min_salary DECIMAL(12,2),
    max_salary DECIMAL(12,2),

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE salary_grades ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON salary_grades
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

CREATE UNIQUE INDEX uniq_salary_grades_code ON salary_grades(tenant_id, code);

-- Designations table
CREATE TABLE designations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    level INT NOT NULL DEFAULT 1 CHECK (level >= 1 AND level <= 10),  -- 1 = highest (CEO), 10 = lowest
    department_id UUID REFERENCES departments(id),
    salary_grade_id UUID REFERENCES salary_grades(id),

    reports_to_id UUID REFERENCES designations(id),  -- Parent designation for hierarchy

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- RLS Policy
ALTER TABLE designations ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON designations
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

-- Indexes
CREATE INDEX idx_designations_tenant ON designations(tenant_id);
CREATE INDEX idx_designations_department ON designations(department_id);
CREATE INDEX idx_designations_level ON designations(tenant_id, level);

-- Add permissions for department and designation management
INSERT INTO permissions (name, description, module, action)
SELECT * FROM (VALUES
    ('department:create', 'Create departments', 'department', 'create'),
    ('department:read', 'View departments', 'department', 'read'),
    ('department:update', 'Update departments', 'department', 'update'),
    ('department:delete', 'Delete departments', 'department', 'delete'),
    ('designation:create', 'Create designations', 'designation', 'create'),
    ('designation:read', 'View designations', 'designation', 'read'),
    ('designation:update', 'Update designations', 'designation', 'update'),
    ('designation:delete', 'Delete designations', 'designation', 'delete'),
    ('salary_grade:create', 'Create salary grades', 'salary_grade', 'create'),
    ('salary_grade:read', 'View salary grades', 'salary_grade', 'read'),
    ('salary_grade:update', 'Update salary grades', 'salary_grade', 'update'),
    ('salary_grade:delete', 'Delete salary grades', 'salary_grade', 'delete')
) AS v(name, description, module, action)
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE permissions.name = v.name);

-- Grant permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.name IN (
    'department:create', 'department:read', 'department:update', 'department:delete',
    'designation:create', 'designation:read', 'designation:update', 'designation:delete',
    'salary_grade:create', 'salary_grade:read', 'salary_grade:update', 'salary_grade:delete'
)
AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
);
```

### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/departments` | List departments | `department:read` |
| GET | `/api/v1/departments/{id}` | Get department | `department:read` |
| POST | `/api/v1/departments` | Create department | `department:create` |
| PUT | `/api/v1/departments/{id}` | Update department | `department:update` |
| DELETE | `/api/v1/departments/{id}` | Delete department | `department:delete` |
| GET | `/api/v1/designations` | List designations | `designation:read` |
| GET | `/api/v1/designations/{id}` | Get designation | `designation:read` |
| POST | `/api/v1/designations` | Create designation | `designation:create` |
| PUT | `/api/v1/designations/{id}` | Update designation | `designation:update` |
| DELETE | `/api/v1/designations/{id}` | Delete designation | `designation:delete` |
| GET | `/api/v1/salary-grades` | List salary grades | `salary_grade:read` |
| POST | `/api/v1/salary-grades` | Create salary grade | `salary_grade:create` |
| PUT | `/api/v1/salary-grades/{id}` | Update salary grade | `salary_grade:update` |
| GET | `/api/v1/organization/chart` | Get org chart data | `department:read` |

---

## Tasks

### Backend Tasks

- [x] **BE-5.2.1**: Create database migration for departments, designations, salary_grades
- [x] **BE-5.2.2**: Create department module (handler, service, repository, entity, dto)
- [x] **BE-5.2.3**: Create designation module (handler, service, repository, entity, dto)
- [ ] **BE-5.2.4**: Create salary grade module (CRUD only, simple) - Deferred to payroll integration
- [ ] **BE-5.2.5**: Create org chart endpoint - Nice-to-have, deferred
- [ ] **BE-5.2.6**: Write unit tests - Deferred

### Frontend Tasks

- [x] **FE-5.2.1**: Create organization settings feature structure
- [x] **FE-5.2.2**: Create department service
- [x] **FE-5.2.3**: Create department list page
- [x] **FE-5.2.4**: Create department form dialog/page
- [x] **FE-5.2.5**: Create designation service
- [x] **FE-5.2.6**: Create designation list page
- [x] **FE-5.2.7**: Create designation form dialog/page
- [ ] **FE-5.2.8**: Create org chart component (tree view) - Nice-to-have, deferred
- [x] **FE-5.2.9**: Add routes and navigation
- [ ] **FE-5.2.10**: Write component tests - Deferred

---

## Definition of Done

- [x] All acceptance criteria verified (AC1-AC4 complete, AC5 org chart deferred)
- [x] Department CRUD working
- [x] Designation CRUD working
- [ ] Org chart displays correctly - Deferred to future enhancement
- [ ] Backend tests passing - Deferred
- [ ] Frontend tests passing - Deferred

---

## Dependencies

- Story 1.2 (Database Foundation) - Complete
- Story 2.5 (Role/Permission Management) - Complete

---

## Dev Notes

### Frontend Structure

```
src/app/features/staff/
├── staff.routes.ts              # Updated with org settings routes
├── pages/
│   ├── staff-list/              # Existing (to be enhanced in 5.1)
│   ├── department-list/         # NEW
│   ├── designation-list/        # NEW
│   └── org-chart/               # NEW
├── components/
│   ├── department-form/         # NEW
│   ├── designation-form/        # NEW
│   └── org-tree/                # NEW
├── services/
│   ├── staff.service.ts         # Existing (enhanced)
│   ├── department.service.ts    # NEW
│   └── designation.service.ts   # NEW
└── models/
    ├── staff.model.ts           # Existing
    ├── department.model.ts      # NEW
    └── designation.model.ts     # NEW
```

### Backend Structure

```
internal/modules/department/
├── handler.go
├── service.go
├── repository.go
├── entity.go
└── dto.go

internal/modules/designation/
├── handler.go
├── service.go
├── repository.go
├── entity.go
└── dto.go
```
