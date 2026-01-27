-- Migration: 000030_promotion_system.up.sql
-- Description: Create promotion rules, batches, and records tables for Epic 4 Story 4.7

-- Create promotion decision enum type
DO $$ BEGIN
    CREATE TYPE promotion_decision AS ENUM ('pending', 'promote', 'retain', 'transfer');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Promotion rules per class - configurable criteria for auto-promotion
CREATE TABLE promotion_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    class_id UUID NOT NULL, -- From class (references classes table when created in Epic 6)
    min_attendance_pct DECIMAL(5,2) DEFAULT 75.00,
    min_overall_marks_pct DECIMAL(5,2) DEFAULT 33.00,
    min_subjects_passed INTEGER DEFAULT 0,
    auto_promote_on_criteria BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    CONSTRAINT uniq_promotion_rule UNIQUE (tenant_id, class_id)
);

-- Enable Row Level Security
ALTER TABLE promotion_rules ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_promotion_rules ON promotion_rules
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for promotion_rules
CREATE INDEX idx_promotion_rules_tenant ON promotion_rules(tenant_id);
CREATE INDEX idx_promotion_rules_class ON promotion_rules(class_id);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_promotion_rules_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_promotion_rules_updated_at
    BEFORE UPDATE ON promotion_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_promotion_rules_updated_at();

-- Promotion batches - groups students being processed together
CREATE TABLE promotion_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    from_academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    to_academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    from_class_id UUID NOT NULL, -- References classes table (Epic 6)
    from_section_id UUID, -- References sections table (Epic 6), null means all sections
    to_class_id UUID, -- Target class for promotions (null if retaining in same class)
    status VARCHAR(20) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'processing', 'completed', 'cancelled')),
    total_students INTEGER NOT NULL DEFAULT 0,
    promoted_count INTEGER NOT NULL DEFAULT 0,
    retained_count INTEGER NOT NULL DEFAULT 0,
    transferred_count INTEGER NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ,
    processed_by UUID REFERENCES users(id),
    cancelled_at TIMESTAMPTZ,
    cancelled_by UUID REFERENCES users(id),
    cancellation_reason TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Enable Row Level Security
ALTER TABLE promotion_batches ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_promotion_batches ON promotion_batches
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for promotion_batches
CREATE INDEX idx_promotion_batches_tenant ON promotion_batches(tenant_id);
CREATE INDEX idx_promotion_batches_from_year ON promotion_batches(from_academic_year_id);
CREATE INDEX idx_promotion_batches_to_year ON promotion_batches(to_academic_year_id);
CREATE INDEX idx_promotion_batches_class ON promotion_batches(from_class_id);
CREATE INDEX idx_promotion_batches_status ON promotion_batches(status);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_promotion_batches_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_promotion_batches_updated_at
    BEFORE UPDATE ON promotion_batches
    FOR EACH ROW
    EXECUTE FUNCTION update_promotion_batches_updated_at();

-- Individual promotion records - one per student per batch
CREATE TABLE promotion_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    batch_id UUID NOT NULL REFERENCES promotion_batches(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id),
    from_enrollment_id UUID NOT NULL REFERENCES student_enrollments(id),
    to_enrollment_id UUID REFERENCES student_enrollments(id), -- Created after processing
    decision promotion_decision NOT NULL DEFAULT 'pending',
    to_class_id UUID, -- Target class (null for pending, same as from for retain)
    to_section_id UUID, -- Target section (assigned before processing)
    roll_number VARCHAR(20), -- New roll number (generated during processing)
    auto_decided BOOLEAN NOT NULL DEFAULT FALSE,
    decision_reason TEXT, -- Why auto-decision was made
    -- Performance metrics (for rule evaluation)
    attendance_pct DECIMAL(5,2), -- From Epic 7 (placeholder for now)
    overall_marks_pct DECIMAL(5,2), -- From Epic 8 (placeholder for now)
    subjects_passed INTEGER, -- From Epic 8 (placeholder for now)
    -- Manual override tracking
    override_by UUID REFERENCES users(id),
    override_at TIMESTAMPTZ,
    override_reason TEXT,
    -- Retention specific
    retention_reason TEXT,
    -- Transfer specific
    transfer_destination TEXT, -- School name if transferring out
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_promotion_record UNIQUE (batch_id, student_id)
);

-- Enable Row Level Security
ALTER TABLE promotion_records ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_promotion_records ON promotion_records
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for promotion_records
CREATE INDEX idx_promotion_records_tenant ON promotion_records(tenant_id);
CREATE INDEX idx_promotion_records_batch ON promotion_records(batch_id);
CREATE INDEX idx_promotion_records_student ON promotion_records(student_id);
CREATE INDEX idx_promotion_records_decision ON promotion_records(decision);
CREATE INDEX idx_promotion_records_from_enrollment ON promotion_records(from_enrollment_id);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_promotion_records_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_promotion_records_updated_at
    BEFORE UPDATE ON promotion_records
    FOR EACH ROW
    EXECUTE FUNCTION update_promotion_records_updated_at();

-- Add promotion permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('promotion:read', 'View Promotions', 'Permission to view promotion batches and records', 'students', NOW(), NOW()),
    ('promotion:create', 'Create Promotions', 'Permission to create promotion batches', 'students', NOW(), NOW()),
    ('promotion:update', 'Update Promotions', 'Permission to update promotion decisions', 'students', NOW(), NOW()),
    ('promotion:process', 'Process Promotions', 'Permission to execute promotion processing', 'students', NOW(), NOW()),
    ('promotion:cancel', 'Cancel Promotions', 'Permission to cancel promotion batches', 'students', NOW(), NOW()),
    ('promotion_rules:read', 'View Promotion Rules', 'Permission to view promotion rules', 'students', NOW(), NOW()),
    ('promotion_rules:manage', 'Manage Promotion Rules', 'Permission to create/update promotion rules', 'students', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign promotion permissions to Admin role
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Admin'
  AND p.code IN (
    'promotion:read',
    'promotion:create',
    'promotion:update',
    'promotion:process',
    'promotion:cancel',
    'promotion_rules:read',
    'promotion_rules:manage'
  )
ON CONFLICT DO NOTHING;

-- Assign read-only promotion permissions to Staff role
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Staff'
  AND p.code IN ('promotion:read', 'promotion_rules:read')
ON CONFLICT DO NOTHING;
