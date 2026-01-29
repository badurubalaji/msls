# Story 5.7: Teacher Subject Assignment

**Epic:** 5 - Staff Management
**Status:** review
**Priority:** P2 (Required for timetable and academic operations)
**Story Points:** 8

---

## User Story

**As an** academic administrator,
**I want** to assign subjects to teachers,
**So that** teaching responsibilities are clearly defined and workload is balanced.

---

## Acceptance Criteria

### AC1: Subject Assignment to Teachers

**Given** a teacher profile exists
**When** assigning subjects
**Then** they can select: subject, class, section
**And** they can set: is class teacher (yes/no)
**And** multiple assignments can be made for the same teacher
**And** assignment is linked to academic year

### AC2: Workload Summary View

**Given** subject assignments exist
**When** viewing workload summary
**Then** they see: total periods per week
**And** they see: classes assigned, subjects taught
**And** workload comparison across teachers is available
**And** can filter by department, subject, class

### AC3: Assignment Conflict Detection

**Given** assignment conflicts may exist
**When** creating or modifying assignments
**Then** system warns if teacher is over-assigned (exceeds max periods)
**And** system warns if subject has no teacher assigned
**And** system warns if same subject-class-section has multiple teachers
**And** workload balance recommendations are shown

### AC4: Class Teacher Assignment

**Given** class-section combinations exist
**When** assigning class teacher
**Then** only one teacher can be class teacher per class-section
**And** class teacher gets additional responsibilities flag
**And** previous class teacher assignment is ended if changed

### AC5: Assignment History

**Given** teacher assignments change over time
**When** viewing assignment history
**Then** they see: previous assignments with date ranges
**And** they see: reason for change
**And** can compare year-over-year assignments

---

## Technical Implementation

### Database Schema

```sql
-- Migration: 000045_teacher_subject_assignment.up.sql

-- Teacher Subject Assignments
CREATE TABLE teacher_subject_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    staff_id UUID NOT NULL REFERENCES staff(id),
    subject_id UUID NOT NULL REFERENCES subjects(id),
    class_id UUID NOT NULL REFERENCES classes(id),
    section_id UUID REFERENCES sections(id),
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),

    periods_per_week INTEGER NOT NULL DEFAULT 0,
    is_class_teacher BOOLEAN NOT NULL DEFAULT false,

    effective_from DATE NOT NULL,
    effective_to DATE, -- NULL means current

    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    remarks TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    -- Only one active assignment per teacher-subject-class-section
    CONSTRAINT uniq_teacher_assignment UNIQUE (tenant_id, staff_id, subject_id, class_id, section_id, academic_year_id)
);

-- Enable RLS
ALTER TABLE teacher_subject_assignments ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_teacher_assignments ON teacher_subject_assignments
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_teacher_assignments_tenant ON teacher_subject_assignments(tenant_id);
CREATE INDEX idx_teacher_assignments_staff ON teacher_subject_assignments(staff_id);
CREATE INDEX idx_teacher_assignments_subject ON teacher_subject_assignments(subject_id);
CREATE INDEX idx_teacher_assignments_class ON teacher_subject_assignments(class_id);
CREATE INDEX idx_teacher_assignments_year ON teacher_subject_assignments(academic_year_id);
CREATE INDEX idx_teacher_assignments_class_teacher ON teacher_subject_assignments(tenant_id, class_id, section_id, is_class_teacher) WHERE is_class_teacher = true;

-- Teacher Workload Settings (per branch/department)
CREATE TABLE teacher_workload_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    min_periods_per_week INTEGER NOT NULL DEFAULT 20,
    max_periods_per_week INTEGER NOT NULL DEFAULT 35,
    max_subjects_per_teacher INTEGER DEFAULT 5,
    max_classes_per_teacher INTEGER DEFAULT 8,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_workload_settings UNIQUE (tenant_id, branch_id)
);

-- Enable RLS
ALTER TABLE teacher_workload_settings ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_workload_settings ON teacher_workload_settings
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('assignment:view', 'View Teacher Assignments', 'Permission to view teacher subject assignments', 'assignment', NOW(), NOW()),
    ('assignment:create', 'Create Teacher Assignment', 'Permission to assign subjects to teachers', 'assignment', NOW(), NOW()),
    ('assignment:update', 'Update Teacher Assignment', 'Permission to modify teacher assignments', 'assignment', NOW(), NOW()),
    ('assignment:delete', 'Delete Teacher Assignment', 'Permission to remove teacher assignments', 'assignment', NOW(), NOW()),
    ('assignment:workload', 'View Workload Report', 'Permission to view teacher workload reports', 'assignment', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;
```

### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/teacher-assignments` | List all assignments with filters | `assignment:view` |
| POST | `/api/v1/teacher-assignments` | Create new assignment | `assignment:create` |
| GET | `/api/v1/teacher-assignments/:id` | Get assignment details | `assignment:view` |
| PUT | `/api/v1/teacher-assignments/:id` | Update assignment | `assignment:update` |
| DELETE | `/api/v1/teacher-assignments/:id` | Remove assignment | `assignment:delete` |
| GET | `/api/v1/staff/:id/assignments` | Get teacher's assignments | `assignment:view` |
| GET | `/api/v1/teacher-assignments/workload` | Get workload summary report | `assignment:workload` |
| GET | `/api/v1/teacher-assignments/unassigned` | Get subjects without teachers | `assignment:view` |
| POST | `/api/v1/teacher-assignments/bulk` | Bulk assign subjects | `assignment:create` |
| GET | `/api/v1/classes/:id/class-teacher` | Get class teacher | `assignment:view` |
| PUT | `/api/v1/classes/:id/class-teacher` | Set class teacher | `assignment:update` |

---

## Tasks

### Backend Tasks

- [x] **BE-5.7.1**: Create database migration for teacher assignment tables
- [x] **BE-5.7.2**: Create assignment module (entity, dto, repository)
- [x] **BE-5.7.3**: Create assignment service with workload validation
- [x] **BE-5.7.4**: Create assignment handler with all endpoints
- [x] **BE-5.7.5**: Add workload report generation
- [x] **BE-5.7.6**: Register routes and add permissions

### Frontend Tasks

- [x] **FE-5.7.1**: Create assignment models and service
- [x] **FE-5.7.2**: Create teacher assignment list page
- [x] **FE-5.7.3**: Create assignment form component
- [x] **FE-5.7.4**: Create workload summary component
- [x] **FE-5.7.5**: Create class teacher assignment component
- [x] **FE-5.7.6**: Add validation warnings for conflicts
- [x] **FE-5.7.7**: Add routes and navigation

---

## Definition of Done

- [x] Teachers can be assigned to subjects and classes
- [x] Class teacher can be assigned per class-section
- [x] Workload summary shows periods per teacher
- [x] System warns on over-assignment
- [x] Unassigned subjects are highlighted
- [x] Assignment history is maintained
- [ ] Unit tests pass

---

## Dependencies

- Story 5.1: Staff Profile Management (completed)
- Epic 6: Academic Structure (classes, sections, subjects) - Required

---

## Dev Notes

### Backend Structure

```
internal/modules/assignment/
├── entity.go
├── dto.go
├── repository.go
├── service.go
├── handler.go
├── errors.go
└── routes.go
```

### Frontend Structure

```
src/app/features/academic/assignment/
├── assignment.routes.ts
├── models/
│   └── assignment.model.ts
├── services/
│   └── assignment.service.ts
├── pages/
│   ├── assignment-list/
│   ├── assignment-form/
│   └── workload-report/
└── components/
    ├── teacher-workload-card/
    └── unassigned-subjects/
```

### Workload Calculation

- Total periods = Sum of periods_per_week for all active assignments
- Over-assigned = Total periods > max_periods_per_week from settings
- Under-assigned = Total periods < min_periods_per_week from settings
