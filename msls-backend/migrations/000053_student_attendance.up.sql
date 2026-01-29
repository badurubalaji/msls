-- Student Attendance Tables
-- Story 7.1: Daily Attendance Marking Interface

-- Student Attendance Table
CREATE TABLE student_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    section_id UUID NOT NULL REFERENCES sections(id),
    attendance_date DATE NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'present',
    late_arrival_time TIME,
    remarks TEXT,

    marked_by UUID NOT NULL REFERENCES users(id),
    marked_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_student_daily_attendance
        UNIQUE (tenant_id, student_id, attendance_date)
);

-- Enable RLS
ALTER TABLE student_attendance ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_student_attendance ON student_attendance
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_student_attendance_section_date
    ON student_attendance(section_id, attendance_date);
CREATE INDEX idx_student_attendance_student_date
    ON student_attendance(student_id, attendance_date DESC);
CREATE INDEX idx_student_attendance_tenant
    ON student_attendance(tenant_id);
CREATE INDEX idx_student_attendance_date
    ON student_attendance(attendance_date);

-- Student Attendance Settings Table (per branch)
CREATE TABLE student_attendance_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    edit_window_minutes INTEGER NOT NULL DEFAULT 120,
    late_threshold_minutes INTEGER NOT NULL DEFAULT 15,
    sms_on_absent BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_student_attendance_settings_branch
        UNIQUE (tenant_id, branch_id)
);

-- Enable RLS
ALTER TABLE student_attendance_settings ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_student_attendance_settings ON student_attendance_settings
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Index
CREATE INDEX idx_student_attendance_settings_tenant
    ON student_attendance_settings(tenant_id);

-- Permissions for student attendance
INSERT INTO permissions (id, code, name, description, module, created_at, updated_at)
VALUES
    (uuid_generate_v7(), 'student_attendance:mark_class', 'Mark Class Attendance', 'Permission to mark attendance for assigned classes', 'student_attendance', NOW(), NOW()),
    (uuid_generate_v7(), 'student_attendance:view_class', 'View Class Attendance', 'Permission to view attendance for assigned classes', 'student_attendance', NOW(), NOW()),
    (uuid_generate_v7(), 'student_attendance:view_all', 'View All Attendance', 'Permission to view all student attendance records', 'student_attendance', NOW(), NOW()),
    (uuid_generate_v7(), 'student_attendance:manage_settings', 'Manage Attendance Settings', 'Permission to manage student attendance settings', 'student_attendance', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Add updated_at trigger using existing function
CREATE OR REPLACE FUNCTION update_student_attendance_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_student_attendance
    BEFORE UPDATE ON student_attendance
    FOR EACH ROW
    EXECUTE FUNCTION update_student_attendance_updated_at();

CREATE OR REPLACE FUNCTION update_student_attendance_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_student_attendance_settings
    BEFORE UPDATE ON student_attendance_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_student_attendance_settings_updated_at();
