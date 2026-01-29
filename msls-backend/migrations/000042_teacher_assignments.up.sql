-- Migration: 000042_teacher_assignments.up.sql
-- Description: Create teacher subject assignment tables

-- Teacher Subject Assignments
CREATE TABLE teacher_subject_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    section_id UUID REFERENCES sections(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),

    periods_per_week INTEGER NOT NULL DEFAULT 0,
    is_class_teacher BOOLEAN NOT NULL DEFAULT false,

    effective_from DATE NOT NULL,
    effective_to DATE, -- NULL means current

    status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'inactive')),
    remarks TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
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
CREATE INDEX idx_teacher_assignments_section ON teacher_subject_assignments(section_id);
CREATE INDEX idx_teacher_assignments_year ON teacher_subject_assignments(academic_year_id);
CREATE INDEX idx_teacher_assignments_status ON teacher_subject_assignments(tenant_id, status);
CREATE INDEX idx_teacher_assignments_class_teacher ON teacher_subject_assignments(tenant_id, class_id, section_id, is_class_teacher)
    WHERE is_class_teacher = true;

-- Unique constraint: One active assignment per teacher-subject-class-section-year
CREATE UNIQUE INDEX uniq_teacher_assignment_active
    ON teacher_subject_assignments(tenant_id, staff_id, subject_id, class_id, COALESCE(section_id, '00000000-0000-0000-0000-000000000000'::uuid), academic_year_id)
    WHERE status = 'active';

-- Teacher Workload Settings (per branch)
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

-- Indexes
CREATE INDEX idx_workload_settings_tenant ON teacher_workload_settings(tenant_id);
CREATE INDEX idx_workload_settings_branch ON teacher_workload_settings(branch_id);

-- Add updated_at triggers
CREATE OR REPLACE FUNCTION update_teacher_assignments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_teacher_assignments_updated_at
    BEFORE UPDATE ON teacher_subject_assignments
    FOR EACH ROW EXECUTE FUNCTION update_teacher_assignments_updated_at();

CREATE TRIGGER trigger_workload_settings_updated_at
    BEFORE UPDATE ON teacher_workload_settings
    FOR EACH ROW EXECUTE FUNCTION update_teacher_assignments_updated_at();

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('assignment:view', 'View Teacher Assignments', 'Permission to view teacher subject assignments', 'assignment', NOW(), NOW()),
    ('assignment:create', 'Create Teacher Assignment', 'Permission to assign subjects to teachers', 'assignment', NOW(), NOW()),
    ('assignment:update', 'Update Teacher Assignment', 'Permission to modify teacher assignments', 'assignment', NOW(), NOW()),
    ('assignment:delete', 'Delete Teacher Assignment', 'Permission to remove teacher assignments', 'assignment', NOW(), NOW()),
    ('assignment:workload', 'View Workload Report', 'Permission to view teacher workload reports', 'assignment', NOW(), NOW()),
    ('assignment:class_teacher', 'Assign Class Teacher', 'Permission to assign class teachers', 'assignment', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT r.id, p.id, NOW(), NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN (
    'assignment:view',
    'assignment:create',
    'assignment:update',
    'assignment:delete',
    'assignment:workload',
    'assignment:class_teacher'
)
ON CONFLICT DO NOTHING;
