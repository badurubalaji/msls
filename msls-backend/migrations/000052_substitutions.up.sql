-- Migration: 000052_substitutions.up.sql
-- Description: Create substitutions table for teacher replacement management

-- Create substitution status enum
DO $$ BEGIN
    CREATE TYPE substitution_status AS ENUM ('pending', 'confirmed', 'completed', 'cancelled');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create substitutions table
CREATE TABLE IF NOT EXISTS substitutions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    original_staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    substitute_staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    substitution_date DATE NOT NULL,
    reason VARCHAR(255),
    status substitution_status NOT NULL DEFAULT 'pending',
    notes TEXT,
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT chk_different_staff CHECK (original_staff_id != substitute_staff_id)
);

-- Create substitution periods table (links substitution to specific timetable entries)
CREATE TABLE IF NOT EXISTS substitution_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    substitution_id UUID NOT NULL REFERENCES substitutions(id) ON DELETE CASCADE,
    timetable_entry_id UUID REFERENCES timetable_entries(id) ON DELETE SET NULL,
    period_slot_id UUID NOT NULL REFERENCES period_slots(id) ON DELETE CASCADE,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,
    section_id UUID REFERENCES sections(id) ON DELETE SET NULL,
    room_number VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_substitutions_tenant ON substitutions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_substitutions_branch ON substitutions(branch_id);
CREATE INDEX IF NOT EXISTS idx_substitutions_original_staff ON substitutions(original_staff_id);
CREATE INDEX IF NOT EXISTS idx_substitutions_substitute_staff ON substitutions(substitute_staff_id);
CREATE INDEX IF NOT EXISTS idx_substitutions_date ON substitutions(substitution_date);
CREATE INDEX IF NOT EXISTS idx_substitutions_status ON substitutions(status);
CREATE INDEX IF NOT EXISTS idx_substitutions_date_range ON substitutions(tenant_id, substitution_date, status) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_substitution_periods_substitution ON substitution_periods(substitution_id);
CREATE INDEX IF NOT EXISTS idx_substitution_periods_entry ON substitution_periods(timetable_entry_id);

-- Unique constraint: One substitution per original teacher per date per period
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_substitution_period
ON substitution_periods(substitution_id, period_slot_id);

-- Add updated_at trigger
DROP TRIGGER IF EXISTS trigger_substitutions_updated_at ON substitutions;
CREATE TRIGGER trigger_substitutions_updated_at
    BEFORE UPDATE ON substitutions
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Enable RLS
ALTER TABLE substitutions ENABLE ROW LEVEL SECURITY;
ALTER TABLE substitution_periods ENABLE ROW LEVEL SECURITY;

-- RLS Policies for substitutions
DROP POLICY IF EXISTS tenant_isolation_substitutions ON substitutions;
CREATE POLICY tenant_isolation_substitutions ON substitutions
    USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid);

DROP POLICY IF EXISTS tenant_isolation_substitution_periods ON substitution_periods;
CREATE POLICY tenant_isolation_substitution_periods ON substitution_periods
    USING (substitution_id IN (
        SELECT id FROM substitutions WHERE tenant_id = current_setting('app.current_tenant_id', true)::uuid
    ));

-- Grant permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON substitutions TO msls;
GRANT SELECT, INSERT, UPDATE, DELETE ON substitution_periods TO msls;

-- Add substitution permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('substitution:view', 'View Substitutions', 'View teacher substitutions', 'timetable', NOW(), NOW()),
    ('substitution:create', 'Create Substitutions', 'Create teacher substitutions', 'timetable', NOW(), NOW()),
    ('substitution:update', 'Update Substitutions', 'Update teacher substitutions', 'timetable', NOW(), NOW()),
    ('substitution:delete', 'Delete Substitutions', 'Delete teacher substitutions', 'timetable', NOW(), NOW()),
    ('substitution:approve', 'Approve Substitutions', 'Approve teacher substitutions', 'timetable', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Add permissions to roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('super_admin', 'admin', 'principal')
AND p.code IN ('substitution:view', 'substitution:create', 'substitution:update', 'substitution:delete', 'substitution:approve')
ON CONFLICT DO NOTHING;

-- Teachers can view substitutions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher'
AND p.code = 'substitution:view'
ON CONFLICT DO NOTHING;
