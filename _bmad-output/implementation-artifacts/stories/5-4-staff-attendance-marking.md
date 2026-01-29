# Story 5.4: Staff Attendance Marking

**Epic:** 5 - Staff Management
**Status:** completed
**Priority:** P1 (Core functionality)
**Story Points:** 8

---

## User Story

**As a** staff member,
**I want** to mark my attendance,
**So that** my presence is recorded daily.

---

## Acceptance Criteria

### AC1: Self Attendance Marking

**Given** a staff member is logged in
**When** marking attendance
**Then** they can mark: present, half-day (with reason)
**And** check-in time is recorded automatically
**And** check-out time can be recorded later
**And** late arrival is flagged if after threshold

### AC2: Attendance Dashboard (Staff View)

**Given** a staff member views their attendance
**When** on the attendance page
**Then** they see today's status (checked-in/not checked-in)
**And** they see this month's summary (present, absent, half-day, late)
**And** they can view historical attendance calendar

### AC3: HR Attendance Management

**Given** HR is on staff attendance management
**When** viewing daily attendance
**Then** they see all staff with status (present, absent, leave, half-day)
**And** they can filter by department, designation, date
**And** they can mark attendance on behalf of staff

### AC4: Attendance Regularization

**Given** attendance regularization is needed
**When** staff submits regularization request
**Then** they can select: date, reason, supporting document
**And** request goes to reporting manager/HR for approval
**And** approved regularization updates attendance record

---

## Technical Implementation

### Database Schema

```sql
-- Migration: 000037_staff_attendance.up.sql

-- Staff attendance records
CREATE TABLE staff_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id),
    attendance_date DATE NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'present',  -- present, absent, half_day, on_leave, holiday
    check_in_time TIMESTAMPTZ,
    check_out_time TIMESTAMPTZ,

    is_late BOOLEAN NOT NULL DEFAULT false,
    late_minutes INTEGER DEFAULT 0,

    half_day_type VARCHAR(20),  -- first_half, second_half
    remarks TEXT,

    marked_by UUID REFERENCES users(id),  -- Who marked (self or HR)
    marked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_staff_attendance UNIQUE (tenant_id, staff_id, attendance_date)
);

-- Enable RLS
ALTER TABLE staff_attendance ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_staff_attendance ON staff_attendance
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_staff_attendance_tenant ON staff_attendance(tenant_id);
CREATE INDEX idx_staff_attendance_staff ON staff_attendance(staff_id);
CREATE INDEX idx_staff_attendance_date ON staff_attendance(tenant_id, attendance_date);
CREATE INDEX idx_staff_attendance_status ON staff_attendance(tenant_id, status);

-- Attendance regularization requests
CREATE TABLE staff_attendance_regularization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id),
    attendance_id UUID REFERENCES staff_attendance(id),

    request_date DATE NOT NULL,
    requested_status VARCHAR(20) NOT NULL,  -- present, half_day
    reason TEXT NOT NULL,
    supporting_document_url TEXT,

    status VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending, approved, rejected
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE staff_attendance_regularization ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_regularization ON staff_attendance_regularization
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_regularization_tenant ON staff_attendance_regularization(tenant_id);
CREATE INDEX idx_regularization_staff ON staff_attendance_regularization(staff_id);
CREATE INDEX idx_regularization_status ON staff_attendance_regularization(tenant_id, status);

-- Attendance settings (per branch)
CREATE TABLE staff_attendance_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    work_start_time TIME NOT NULL DEFAULT '09:00',
    work_end_time TIME NOT NULL DEFAULT '17:00',
    late_threshold_minutes INTEGER NOT NULL DEFAULT 15,  -- Minutes after start time to mark late
    half_day_threshold_hours DECIMAL(4,2) NOT NULL DEFAULT 4.0,  -- Hours to consider half day

    allow_self_checkout BOOLEAN NOT NULL DEFAULT true,
    require_regularization_approval BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_attendance_settings UNIQUE (tenant_id, branch_id)
);

-- Enable RLS
ALTER TABLE staff_attendance_settings ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_attendance_settings ON staff_attendance_settings
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('attendance:mark_self', 'Mark Own Attendance', 'Permission to mark own attendance', 'attendance', NOW(), NOW()),
    ('attendance:mark_others', 'Mark Others Attendance', 'Permission to mark attendance for other staff', 'attendance', NOW(), NOW()),
    ('attendance:view_self', 'View Own Attendance', 'Permission to view own attendance records', 'attendance', NOW(), NOW()),
    ('attendance:view_all', 'View All Attendance', 'Permission to view all staff attendance', 'attendance', NOW(), NOW()),
    ('attendance:regularize', 'Request Regularization', 'Permission to request attendance regularization', 'attendance', NOW(), NOW()),
    ('attendance:approve_regularization', 'Approve Regularization', 'Permission to approve regularization requests', 'attendance', NOW(), NOW()),
    ('attendance:settings', 'Manage Attendance Settings', 'Permission to manage attendance settings', 'attendance', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT r.id, p.id, NOW(), NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN (
    'attendance:mark_self',
    'attendance:mark_others',
    'attendance:view_self',
    'attendance:view_all',
    'attendance:regularize',
    'attendance:approve_regularization',
    'attendance:settings'
)
ON CONFLICT DO NOTHING;
```

### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| POST | `/api/v1/attendance/check-in` | Mark check-in | `attendance:mark_self` |
| POST | `/api/v1/attendance/check-out` | Mark check-out | `attendance:mark_self` |
| GET | `/api/v1/attendance/my` | Get own attendance records | `attendance:view_self` |
| GET | `/api/v1/attendance/my/today` | Get today's attendance | `attendance:view_self` |
| GET | `/api/v1/attendance/my/summary` | Get monthly summary | `attendance:view_self` |
| GET | `/api/v1/attendance` | List all attendance (HR) | `attendance:view_all` |
| POST | `/api/v1/attendance/mark` | Mark attendance for staff (HR) | `attendance:mark_others` |
| POST | `/api/v1/attendance/regularization` | Submit regularization request | `attendance:regularize` |
| GET | `/api/v1/attendance/regularization` | List regularization requests | `attendance:approve_regularization` |
| PUT | `/api/v1/attendance/regularization/{id}/approve` | Approve regularization | `attendance:approve_regularization` |
| PUT | `/api/v1/attendance/regularization/{id}/reject` | Reject regularization | `attendance:approve_regularization` |

---

## Tasks

### Backend Tasks

- [x] **BE-5.4.1**: Create database migration for staff_attendance tables
- [x] **BE-5.4.2**: Create attendance module (entity, dto, repository)
- [x] **BE-5.4.3**: Create attendance service with business logic
- [x] **BE-5.4.4**: Create attendance handler with all endpoints
- [x] **BE-5.4.5**: Create regularization service and handler
- [x] **BE-5.4.6**: Register routes in main.go

### Frontend Tasks

- [x] **FE-5.4.1**: Create attendance feature structure
- [x] **FE-5.4.2**: Create attendance service
- [x] **FE-5.4.3**: Create staff attendance dashboard (self-view)
- [x] **FE-5.4.4**: Create check-in/check-out component
- [x] **FE-5.4.5**: Create attendance calendar view
- [x] **FE-5.4.6**: Create regularization request form
- [x] **FE-5.4.7**: Create HR attendance management page
- [x] **FE-5.4.8**: Add routes and navigation

---

## Definition of Done

- [x] Staff can check-in and check-out
- [x] Late arrival is automatically flagged
- [x] Staff can view their attendance history
- [x] Staff can request regularization
- [x] HR can view all staff attendance
- [x] HR can mark attendance on behalf of staff
- [x] HR can approve/reject regularization requests

---

## Dependencies

- Story 5.1 (Staff Profile Management) - Complete
- Story 5.2 (Department & Designation) - Complete

---

## Dev Notes

### Backend Structure

```
internal/modules/attendance/
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
src/app/features/attendance/
├── attendance.routes.ts
├── models/
│   └── attendance.model.ts
├── services/
│   └── attendance.service.ts
├── pages/
│   ├── my-attendance/           # Staff's own attendance view
│   ├── attendance-management/   # HR view for all staff
│   └── regularization/          # Regularization requests
└── components/
    ├── check-in-out/           # Check-in/out button component
    ├── attendance-calendar/    # Calendar view
    └── attendance-summary/     # Summary stats component
```
