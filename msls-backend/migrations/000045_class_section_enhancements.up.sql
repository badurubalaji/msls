-- Migration: 000045_class_section_enhancements.up.sql
-- Description: Enhance classes and sections with streams, class teacher, and academic year support

-- Streams (Science, Commerce, Arts for senior classes)
CREATE TABLE IF NOT EXISTS streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,             -- e.g., "Science", "Commerce", "Arts"
    code VARCHAR(20) NOT NULL,              -- e.g., "SCI", "COM", "ARTS"
    description TEXT,
    display_order INTEGER NOT NULL DEFAULT 0,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    CONSTRAINT uniq_stream_code UNIQUE (tenant_id, code)
);

-- Enable RLS
ALTER TABLE streams ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_streams ON streams
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_streams_tenant ON streams(tenant_id);
CREATE INDEX idx_streams_active ON streams(tenant_id, is_active);

-- Add has_streams flag to classes (for senior classes like 11th, 12th)
ALTER TABLE classes ADD COLUMN IF NOT EXISTS has_streams BOOLEAN NOT NULL DEFAULT false;

-- Class-Stream mapping (which streams are available for which class)
CREATE TABLE IF NOT EXISTS class_streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_class_stream UNIQUE (class_id, stream_id)
);

-- Enable RLS
ALTER TABLE class_streams ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_class_streams ON class_streams
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_class_streams_class ON class_streams(class_id);
CREATE INDEX idx_class_streams_stream ON class_streams(stream_id);

-- Add class_teacher_id, academic_year_id, and stream_id to sections
ALTER TABLE sections ADD COLUMN IF NOT EXISTS class_teacher_id UUID REFERENCES staff(id);
ALTER TABLE sections ADD COLUMN IF NOT EXISTS academic_year_id UUID REFERENCES academic_years(id);
ALTER TABLE sections ADD COLUMN IF NOT EXISTS stream_id UUID REFERENCES streams(id);
ALTER TABLE sections ADD COLUMN IF NOT EXISTS room_number VARCHAR(50);

-- Add index for class teacher
CREATE INDEX IF NOT EXISTS idx_sections_class_teacher ON sections(class_teacher_id);
CREATE INDEX IF NOT EXISTS idx_sections_academic_year ON sections(academic_year_id);
CREATE INDEX IF NOT EXISTS idx_sections_stream ON sections(stream_id);

-- Trigger for streams updated_at
CREATE TRIGGER trigger_streams_updated_at
    BEFORE UPDATE ON streams
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

-- Add stream permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('stream:view', 'View Streams', 'Permission to view streams list', 'academic', NOW(), NOW()),
    ('stream:create', 'Create Stream', 'Permission to create streams', 'academic', NOW(), NOW()),
    ('stream:update', 'Update Stream', 'Permission to update streams', 'academic', NOW(), NOW()),
    ('stream:delete', 'Delete Stream', 'Permission to delete streams', 'academic', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign stream permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT r.id, p.id, NOW(), NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN ('stream:view', 'stream:create', 'stream:update', 'stream:delete')
ON CONFLICT DO NOTHING;

-- Seed default streams
INSERT INTO streams (tenant_id, name, code, description, display_order, created_at, updated_at)
SELECT
    t.id,
    s.name,
    s.code,
    s.description,
    s.display_order,
    NOW(),
    NOW()
FROM tenants t
CROSS JOIN (
    VALUES
        ('Science', 'SCI', 'Science stream with Physics, Chemistry, Mathematics, Biology', 1),
        ('Commerce', 'COM', 'Commerce stream with Accountancy, Business Studies, Economics', 2),
        ('Arts', 'ARTS', 'Arts/Humanities stream with History, Geography, Political Science', 3)
) AS s(name, code, description, display_order)
ON CONFLICT (tenant_id, code) DO NOTHING;
