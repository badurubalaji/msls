-- Migration: 000029_student_enrollments.up.sql
-- Description: Create student_enrollments and enrollment_status_changes tables for Epic 4 Story 4.6

-- Create enrollment status enum type
DO $$ BEGIN
    CREATE TYPE enrollment_status AS ENUM ('active', 'completed', 'transferred', 'dropout');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Student enrollments table - tracks enrollment history per academic year
CREATE TABLE student_enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    class_id UUID, -- References classes table (Epic 6), nullable until then
    section_id UUID, -- References sections table (Epic 6), nullable
    roll_number VARCHAR(20),
    class_teacher_id UUID, -- References staff (Epic 5), nullable until then
    status enrollment_status NOT NULL DEFAULT 'active',
    enrollment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    completion_date DATE,
    transfer_date DATE,
    transfer_reason TEXT,
    dropout_date DATE,
    dropout_reason TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    CONSTRAINT uniq_student_year UNIQUE (student_id, academic_year_id)
);

-- Enable Row Level Security
ALTER TABLE student_enrollments ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_student_enrollments ON student_enrollments
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for student_enrollments
CREATE INDEX idx_enrollments_tenant ON student_enrollments(tenant_id);
CREATE INDEX idx_enrollments_student ON student_enrollments(student_id);
CREATE INDEX idx_enrollments_year ON student_enrollments(academic_year_id);
CREATE INDEX idx_enrollments_class ON student_enrollments(class_id) WHERE class_id IS NOT NULL;
CREATE INDEX idx_enrollments_status ON student_enrollments(status) WHERE status = 'active';

-- Ensure only one active enrollment per student (across all academic years)
CREATE UNIQUE INDEX uniq_active_enrollment ON student_enrollments(student_id)
    WHERE status = 'active';

-- Enrollment status change log - audit trail for status changes
CREATE TABLE enrollment_status_changes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    enrollment_id UUID NOT NULL REFERENCES student_enrollments(id) ON DELETE CASCADE,
    from_status enrollment_status,
    to_status enrollment_status NOT NULL,
    change_reason TEXT,
    change_date DATE NOT NULL DEFAULT CURRENT_DATE,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    changed_by UUID NOT NULL REFERENCES users(id)
);

-- Enable Row Level Security
ALTER TABLE enrollment_status_changes ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_enrollment_status_changes ON enrollment_status_changes
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for enrollment_status_changes
CREATE INDEX idx_status_changes_tenant ON enrollment_status_changes(tenant_id);
CREATE INDEX idx_status_changes_enrollment ON enrollment_status_changes(enrollment_id);
CREATE INDEX idx_status_changes_date ON enrollment_status_changes(changed_at);

-- Add enrollment permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('enrollments:read', 'View Enrollments', 'Permission to view student enrollment history', 'students', NOW(), NOW()),
    ('enrollments:create', 'Create Enrollments', 'Permission to create student enrollments', 'students', NOW(), NOW()),
    ('enrollments:update', 'Update Enrollments', 'Permission to update student enrollments', 'students', NOW(), NOW()),
    ('enrollments:delete', 'Delete Enrollments', 'Permission to delete student enrollments', 'students', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_student_enrollments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_student_enrollments_updated_at
    BEFORE UPDATE ON student_enrollments
    FOR EACH ROW
    EXECUTE FUNCTION update_student_enrollments_updated_at();
