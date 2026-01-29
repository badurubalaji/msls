-- Migration: 000041_academic_structure.up.sql
-- Description: Create academic structure tables (classes, sections, subjects)

-- Classes (e.g., Class 1, Class 2, ... Class 12)
CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    name VARCHAR(50) NOT NULL,          -- e.g., "Class 1", "Grade 10"
    code VARCHAR(20) NOT NULL,          -- e.g., "C1", "G10"
    display_order INTEGER NOT NULL DEFAULT 0,
    description TEXT,

    -- Academic configuration
    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    CONSTRAINT uniq_class_code UNIQUE (tenant_id, branch_id, code)
);

-- Enable RLS
ALTER TABLE classes ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_classes ON classes
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_classes_tenant ON classes(tenant_id);
CREATE INDEX idx_classes_branch ON classes(branch_id);
CREATE INDEX idx_classes_active ON classes(tenant_id, is_active);

-- Sections (e.g., A, B, C within a class)
CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,

    name VARCHAR(20) NOT NULL,          -- e.g., "A", "B", "Science"
    code VARCHAR(20) NOT NULL,          -- e.g., "A", "SCI"
    capacity INTEGER DEFAULT 40,        -- Max students in section
    display_order INTEGER NOT NULL DEFAULT 0,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    CONSTRAINT uniq_section_code UNIQUE (tenant_id, class_id, code)
);

-- Enable RLS
ALTER TABLE sections ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_sections ON sections
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_sections_tenant ON sections(tenant_id);
CREATE INDEX idx_sections_class ON sections(class_id);

-- Subjects
CREATE TABLE subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,         -- e.g., "Mathematics", "English"
    code VARCHAR(20) NOT NULL,          -- e.g., "MATH", "ENG"
    short_name VARCHAR(20),             -- e.g., "Math", "Eng"
    description TEXT,

    subject_type VARCHAR(30) NOT NULL DEFAULT 'core'
        CHECK (subject_type IN ('core', 'elective', 'language', 'co_curricular', 'vocational')),

    -- Credit/Marks configuration
    max_marks INTEGER DEFAULT 100,
    passing_marks INTEGER DEFAULT 35,
    credit_hours DECIMAL(4,2) DEFAULT 0,

    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    CONSTRAINT uniq_subject_code UNIQUE (tenant_id, code)
);

-- Enable RLS
ALTER TABLE subjects ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_subjects ON subjects
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_subjects_tenant ON subjects(tenant_id);
CREATE INDEX idx_subjects_type ON subjects(tenant_id, subject_type);
CREATE INDEX idx_subjects_active ON subjects(tenant_id, is_active);

-- Class-Subject mapping (which subjects are taught in which class)
CREATE TABLE class_subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,

    is_mandatory BOOLEAN NOT NULL DEFAULT true,
    periods_per_week INTEGER NOT NULL DEFAULT 5,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_class_subject UNIQUE (tenant_id, class_id, subject_id)
);

-- Enable RLS
ALTER TABLE class_subjects ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_class_subjects ON class_subjects
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_class_subjects_tenant ON class_subjects(tenant_id);
CREATE INDEX idx_class_subjects_class ON class_subjects(class_id);
CREATE INDEX idx_class_subjects_subject ON class_subjects(subject_id);

-- Add updated_at triggers
CREATE OR REPLACE FUNCTION update_classes_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_classes_updated_at
    BEFORE UPDATE ON classes
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

CREATE TRIGGER trigger_sections_updated_at
    BEFORE UPDATE ON sections
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

CREATE TRIGGER trigger_subjects_updated_at
    BEFORE UPDATE ON subjects
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

CREATE TRIGGER trigger_class_subjects_updated_at
    BEFORE UPDATE ON class_subjects
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('class:view', 'View Classes', 'Permission to view class list', 'academic', NOW(), NOW()),
    ('class:create', 'Create Class', 'Permission to create classes', 'academic', NOW(), NOW()),
    ('class:update', 'Update Class', 'Permission to update classes', 'academic', NOW(), NOW()),
    ('class:delete', 'Delete Class', 'Permission to delete classes', 'academic', NOW(), NOW()),
    ('section:view', 'View Sections', 'Permission to view sections', 'academic', NOW(), NOW()),
    ('section:create', 'Create Section', 'Permission to create sections', 'academic', NOW(), NOW()),
    ('section:update', 'Update Section', 'Permission to update sections', 'academic', NOW(), NOW()),
    ('section:delete', 'Delete Section', 'Permission to delete sections', 'academic', NOW(), NOW()),
    ('subject:view', 'View Subjects', 'Permission to view subjects', 'academic', NOW(), NOW()),
    ('subject:create', 'Create Subject', 'Permission to create subjects', 'academic', NOW(), NOW()),
    ('subject:update', 'Update Subject', 'Permission to update subjects', 'academic', NOW(), NOW()),
    ('subject:delete', 'Delete Subject', 'Permission to delete subjects', 'academic', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT r.id, p.id, NOW(), NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN (
    'class:view', 'class:create', 'class:update', 'class:delete',
    'section:view', 'section:create', 'section:update', 'section:delete',
    'subject:view', 'subject:create', 'subject:update', 'subject:delete'
)
ON CONFLICT DO NOTHING;
