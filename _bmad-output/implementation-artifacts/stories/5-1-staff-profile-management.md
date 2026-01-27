# Story 5.1: Staff Profile Management

**Epic:** 5 - Staff Management
**Status:** dev-complete
**Priority:** P0 (Foundation for all Epic 5 stories)
**Story Points:** 8

---

## User Story

**As an** HR administrator,
**I want** to create and manage staff profiles,
**So that** all employee information is centrally maintained.

---

## Acceptance Criteria

### AC1: Staff Profile Creation

**Given** HR is creating a new staff member
**When** they fill the staff form
**Then** they can enter personal details (first name, last name, DOB, gender, blood group)
**And** they can enter contact details (phone, email, address)
**And** they can enter employment details (employee ID, join date, designation)
**And** they can select staff type (teaching/non-teaching)
**And** they can select department and reporting manager
**And** employee ID is auto-generated with configurable prefix

### AC2: Staff Profile View

**Given** a staff profile exists
**When** viewing the profile
**Then** all information is displayed in organized tabs
**And** photo is displayed with fallback avatar
**And** quick actions are available (edit, view attendance, view salary)

### AC3: Staff List & Search

**Given** HR is on staff management page
**When** viewing the staff list
**Then** they see all staff with basic info (name, department, designation, status)
**And** they can search by name, employee ID, phone
**And** they can filter by department, designation, staff type, status
**And** pagination is cursor-based

### AC4: Staff Profile Update

**Given** a staff profile exists
**When** HR edits the profile
**Then** all fields are editable except employee ID
**And** changes are audited with old/new values
**And** updated_by and updated_at are recorded

### AC5: Staff Status Management

**Given** a staff profile exists
**When** managing staff status
**Then** status can be: active, inactive, terminated, on_leave
**And** status change reason is captured
**And** status history is maintained

---

## Technical Implementation

### Database Schema

```sql
-- Staff table (core entity)
CREATE TABLE staff (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    -- Employee identification
    employee_id VARCHAR(50) NOT NULL,
    employee_id_prefix VARCHAR(10) DEFAULT 'EMP',

    -- Personal details
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(20) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    blood_group VARCHAR(10),
    nationality VARCHAR(50) DEFAULT 'Indian',
    religion VARCHAR(50),
    marital_status VARCHAR(20),

    -- Contact details
    personal_email VARCHAR(255),
    work_email VARCHAR(255) NOT NULL,
    personal_phone VARCHAR(20),
    work_phone VARCHAR(20) NOT NULL,
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(20),
    emergency_contact_relation VARCHAR(50),

    -- Address
    current_address_line1 VARCHAR(255),
    current_address_line2 VARCHAR(255),
    current_city VARCHAR(100),
    current_state VARCHAR(100),
    current_pincode VARCHAR(10),
    current_country VARCHAR(100) DEFAULT 'India',
    permanent_address_line1 VARCHAR(255),
    permanent_address_line2 VARCHAR(255),
    permanent_city VARCHAR(100),
    permanent_state VARCHAR(100),
    permanent_pincode VARCHAR(10),
    permanent_country VARCHAR(100) DEFAULT 'India',
    same_as_current BOOLEAN DEFAULT false,

    -- Employment details
    staff_type VARCHAR(20) NOT NULL CHECK (staff_type IN ('teaching', 'non_teaching')),
    department_id UUID REFERENCES departments(id),
    designation_id UUID REFERENCES designations(id),
    reporting_manager_id UUID REFERENCES staff(id),
    join_date DATE NOT NULL,
    confirmation_date DATE,
    probation_end_date DATE,

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'terminated', 'on_leave')),
    status_reason TEXT,
    termination_date DATE,

    -- Profile
    photo_url VARCHAR(500),
    bio TEXT,

    -- Audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    deleted_at TIMESTAMPTZ
);

-- RLS Policy
ALTER TABLE staff ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON staff
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

-- Indexes
CREATE INDEX idx_staff_tenant ON staff(tenant_id);
CREATE INDEX idx_staff_branch ON staff(branch_id);
CREATE INDEX idx_staff_department ON staff(department_id);
CREATE INDEX idx_staff_employee_id ON staff(tenant_id, employee_id);
CREATE UNIQUE INDEX uniq_staff_employee_id ON staff(tenant_id, employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_status ON staff(tenant_id, status);
CREATE INDEX idx_staff_name ON staff(tenant_id, first_name, last_name);

-- Staff status history
CREATE TABLE staff_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    old_status VARCHAR(20),
    new_status VARCHAR(20) NOT NULL,
    reason TEXT,
    effective_date DATE NOT NULL,
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_staff_status_history ON staff_status_history(staff_id);

-- Departments table (dependency)
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    head_id UUID REFERENCES staff(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE departments ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON departments
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

-- Designations table (dependency)
CREATE TABLE designations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    level INT NOT NULL DEFAULT 1,
    department_id UUID REFERENCES departments(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE designations ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON designations
    USING (tenant_id = current_setting('app.tenant_id')::UUID);
```

### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/staff` | List staff with filters | `staff:read` |
| GET | `/api/v1/staff/{id}` | Get staff profile | `staff:read` |
| POST | `/api/v1/staff` | Create staff | `staff:create` |
| PUT | `/api/v1/staff/{id}` | Update staff | `staff:update` |
| PATCH | `/api/v1/staff/{id}/status` | Update status | `staff:update` |
| POST | `/api/v1/staff/{id}/photo` | Upload photo | `staff:update` |
| GET | `/api/v1/staff/employee-id/preview` | Preview next employee ID | `staff:read` |
| GET | `/api/v1/departments` | List departments | `department:read` |
| GET | `/api/v1/designations` | List designations | `designation:read` |

### Backend Tasks

- [x] **BE-5.1.1**: Create departments and designations tables with migrations
- [ ] **BE-5.1.2**: Create staff table migration with all columns
- [ ] **BE-5.1.3**: Create staff module (handler, service, repository, entity, dto)
- [ ] **BE-5.1.4**: Implement auto-generated employee ID with prefix
- [ ] **BE-5.1.5**: Implement staff CRUD operations
- [ ] **BE-5.1.6**: Implement staff list with search/filter/pagination
- [ ] **BE-5.1.7**: Implement status change with history tracking
- [ ] **BE-5.1.8**: Implement photo upload endpoint
- [ ] **BE-5.1.9**: Write backend unit tests

### Frontend Tasks

- [ ] **FE-5.1.1**: Create staff feature module structure
- [ ] **FE-5.1.2**: Create staff model interfaces
- [ ] **FE-5.1.3**: Create staff service with all API calls
- [ ] **FE-5.1.4**: Create staff list page with data table
- [ ] **FE-5.1.5**: Implement search and filter components
- [ ] **FE-5.1.6**: Create staff form component (multi-step or tabbed)
- [ ] **FE-5.1.7**: Create staff detail/view page with tabs
- [ ] **FE-5.1.8**: Implement photo upload with preview
- [ ] **FE-5.1.9**: Write frontend component tests

---

## Dev Notes

### Architecture Patterns (MUST FOLLOW)

**Backend Module Structure:**
```
internal/modules/staff/
├── handler.go      # HTTP handlers
├── service.go      # Business logic
├── repository.go   # Database operations
├── entity.go       # Staff, Department, Designation entities
├── dto.go          # CreateStaffDTO, UpdateStaffDTO, StaffResponse
└── staff_test.go   # Tests
```

**Frontend Feature Structure:**
```
src/app/features/staff/
├── staff.routes.ts
├── pages/
│   ├── staff-list/
│   └── staff-detail/
├── components/
│   ├── staff-form/
│   ├── staff-card/
│   └── staff-filters/
├── services/
│   └── staff.service.ts
└── models/
    └── staff.model.ts
```

### Similar Patterns from Epic 4 (Student Module)

Reference the student module implementation for patterns:
- `internal/modules/student/` - Backend patterns
- `src/app/features/students/` - Frontend patterns
- Photo upload: Use same storage abstraction as student photos
- List pagination: Use cursor-based pagination like student list
- Form validation: Similar multi-step form pattern

### Employee ID Auto-Generation

```go
// Service implementation
func (s *StaffService) GenerateEmployeeID(ctx context.Context, tenantID uuid.UUID) (string, error) {
    // Get tenant's employee ID prefix from config (default: "EMP")
    prefix := s.getEmployeeIDPrefix(tenantID)

    // Get next sequence number
    var maxNum int
    err := s.db.WithContext(ctx).Raw(`
        SELECT COALESCE(MAX(CAST(SUBSTRING(employee_id FROM '\d+$') AS INT)), 0)
        FROM staff WHERE tenant_id = ? AND employee_id LIKE ?
    `, tenantID, prefix+"%").Scan(&maxNum).Error

    return fmt.Sprintf("%s%05d", prefix, maxNum+1), nil
}
```

### Status Change Pattern

```go
// Status change with history tracking
func (s *StaffService) UpdateStatus(ctx context.Context, id uuid.UUID, req StatusUpdateDTO) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Get current status
        var staff Staff
        if err := tx.First(&staff, id).Error; err != nil {
            return err
        }

        // Record history
        history := StaffStatusHistory{
            StaffID:       id,
            OldStatus:     staff.Status,
            NewStatus:     req.Status,
            Reason:        req.Reason,
            EffectiveDate: req.EffectiveDate,
            ChangedBy:     ctx.Value("user_id").(uuid.UUID),
        }
        if err := tx.Create(&history).Error; err != nil {
            return err
        }

        // Update staff status
        return tx.Model(&staff).Updates(map[string]interface{}{
            "status":       req.Status,
            "status_reason": req.Reason,
        }).Error
    })
}
```

### Project Structure Notes

- Follow existing patterns from `internal/modules/student/` for backend
- Follow existing patterns from `src/app/features/students/` for frontend
- Use shared components from `src/app/shared/components/`
- Staff photos stored in MinIO/local storage like student photos

### References

- [Source: architecture.md#Backend-Module-Structure] - Module patterns
- [Source: architecture.md#Angular-Frontend-Rules] - Angular patterns
- [Source: architecture.md#Database-Rules] - Schema patterns
- [Source: project-context.md#Testing-Requirements] - Test patterns
- [Source: epic-05-staff-management.md#Story-5.1] - Requirements

---

## Definition of Done

- [ ] All acceptance criteria verified
- [ ] Staff CRUD operations working
- [ ] Employee ID auto-generation working
- [ ] Search and filters working
- [ ] Status change with history working
- [ ] Photo upload working
- [ ] Backend tests passing
- [ ] Frontend tests passing

---

## Dependencies

- Story 1.2 (Database Foundation) - ✅ Complete
- Story 2.5 (Role/Permission Management) - ✅ Complete
- Branches table exists - ✅ Complete

---

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

