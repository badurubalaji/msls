-- Migration: 000037_staff_attendance.up.sql
-- Description: Create staff attendance tables for Story 5.4

-- Staff attendance records
CREATE TABLE staff_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
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

    CONSTRAINT uniq_staff_attendance UNIQUE (tenant_id, staff_id, attendance_date),
    CONSTRAINT chk_attendance_status CHECK (status IN ('present', 'absent', 'half_day', 'on_leave', 'holiday')),
    CONSTRAINT chk_half_day_type CHECK (half_day_type IS NULL OR half_day_type IN ('first_half', 'second_half'))
);

-- Enable RLS
ALTER TABLE staff_attendance ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_staff_attendance ON staff_attendance
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_staff_attendance_tenant ON staff_attendance(tenant_id);
CREATE INDEX idx_staff_attendance_staff ON staff_attendance(staff_id);
CREATE INDEX idx_staff_attendance_date ON staff_attendance(tenant_id, attendance_date);
CREATE INDEX idx_staff_attendance_staff_date ON staff_attendance(staff_id, attendance_date);
CREATE INDEX idx_staff_attendance_status ON staff_attendance(tenant_id, status);

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_staff_attendance_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_staff_attendance_updated_at
    BEFORE UPDATE ON staff_attendance
    FOR EACH ROW
    EXECUTE FUNCTION update_staff_attendance_updated_at();

-- Attendance regularization requests
CREATE TABLE staff_attendance_regularization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    attendance_id UUID REFERENCES staff_attendance(id) ON DELETE SET NULL,

    request_date DATE NOT NULL,
    requested_status VARCHAR(20) NOT NULL,  -- present, half_day
    reason TEXT NOT NULL,
    supporting_document_url TEXT,

    status VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending, approved, rejected
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMPTZ,
    rejection_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_regularization_requested_status CHECK (requested_status IN ('present', 'half_day')),
    CONSTRAINT chk_regularization_status CHECK (status IN ('pending', 'approved', 'rejected'))
);

-- Enable RLS
ALTER TABLE staff_attendance_regularization ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_regularization ON staff_attendance_regularization
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_regularization_tenant ON staff_attendance_regularization(tenant_id);
CREATE INDEX idx_regularization_staff ON staff_attendance_regularization(staff_id);
CREATE INDEX idx_regularization_status ON staff_attendance_regularization(tenant_id, status);
CREATE INDEX idx_regularization_date ON staff_attendance_regularization(tenant_id, request_date);

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_regularization_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_regularization_updated_at
    BEFORE UPDATE ON staff_attendance_regularization
    FOR EACH ROW
    EXECUTE FUNCTION update_regularization_updated_at();

-- Attendance settings (per branch)
CREATE TABLE staff_attendance_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,

    work_start_time TIME NOT NULL DEFAULT '09:00',
    work_end_time TIME NOT NULL DEFAULT '17:00',
    late_threshold_minutes INTEGER NOT NULL DEFAULT 15,
    half_day_threshold_hours DECIMAL(4,2) NOT NULL DEFAULT 4.0,

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

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_attendance_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_attendance_settings_updated_at
    BEFORE UPDATE ON staff_attendance_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_attendance_settings_updated_at();

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
