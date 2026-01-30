-- Exam Types Table
-- Story 8.1: Exam Type Configuration

CREATE TABLE exam_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    weightage DECIMAL(5,2) NOT NULL DEFAULT 0,
    evaluation_type VARCHAR(20) NOT NULL DEFAULT 'marks',
    default_max_marks INTEGER NOT NULL DEFAULT 100,
    default_passing_marks INTEGER,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT uniq_exam_type_code UNIQUE (tenant_id, code),
    CONSTRAINT chk_evaluation_type CHECK (evaluation_type IN ('marks', 'grade')),
    CONSTRAINT chk_weightage CHECK (weightage >= 0 AND weightage <= 100),
    CONSTRAINT chk_max_marks CHECK (default_max_marks > 0)
);

-- Enable RLS
ALTER TABLE exam_types ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_exam_types ON exam_types
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_exam_types_tenant ON exam_types(tenant_id);
CREATE INDEX idx_exam_types_active ON exam_types(tenant_id, is_active) WHERE is_active = true;
CREATE INDEX idx_exam_types_order ON exam_types(tenant_id, display_order);

-- Updated at trigger
CREATE OR REPLACE FUNCTION update_exam_types_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_exam_types
    BEFORE UPDATE ON exam_types
    FOR EACH ROW
    EXECUTE FUNCTION update_exam_types_updated_at();

-- Permissions for exam types
INSERT INTO permissions (id, code, name, description, module, created_at, updated_at)
VALUES
    (uuid_generate_v7(), 'exam:type:view', 'View Exam Types', 'Permission to view exam type configurations', 'exam', NOW(), NOW()),
    (uuid_generate_v7(), 'exam:type:create', 'Create Exam Types', 'Permission to create new exam types', 'exam', NOW(), NOW()),
    (uuid_generate_v7(), 'exam:type:update', 'Update Exam Types', 'Permission to update exam type configurations', 'exam', NOW(), NOW()),
    (uuid_generate_v7(), 'exam:type:delete', 'Delete Exam Types', 'Permission to delete exam types', 'exam', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign exam type permissions to admin roles (super_admin, admin, principal)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('super_admin', 'admin', 'principal')
AND p.code IN ('exam:type:view', 'exam:type:create', 'exam:type:update', 'exam:type:delete')
ON CONFLICT DO NOTHING;

-- Teachers can view exam types
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher'
AND p.code = 'exam:type:view'
ON CONFLICT DO NOTHING;
