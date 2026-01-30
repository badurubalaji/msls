-- Migration: 000050_timetables.up.sql
-- Description: Create timetables and timetable_entries tables for Story 6.4

-- Create timetable status enum
DO $$ BEGIN
    CREATE TYPE timetable_status AS ENUM ('draft', 'published', 'archived');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Timetables table - master record for section timetables
CREATE TABLE timetables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),
    section_id UUID NOT NULL REFERENCES sections(id),
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status timetable_status NOT NULL DEFAULT 'draft',
    effective_from DATE,
    effective_to DATE,
    published_at TIMESTAMPTZ,
    published_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    version INTEGER NOT NULL DEFAULT 1
);

-- Enable Row Level Security
ALTER TABLE timetables ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_timetables ON timetables
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for timetables
CREATE INDEX idx_timetables_tenant ON timetables(tenant_id);
CREATE INDEX idx_timetables_branch ON timetables(branch_id);
CREATE INDEX idx_timetables_section ON timetables(section_id);
CREATE INDEX idx_timetables_academic_year ON timetables(academic_year_id);
CREATE INDEX idx_timetables_status ON timetables(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_timetables_deleted_at ON timetables(deleted_at) WHERE deleted_at IS NOT NULL;

-- Only one published timetable per section per academic year
CREATE UNIQUE INDEX uniq_published_timetable ON timetables(section_id, academic_year_id)
    WHERE status = 'published' AND deleted_at IS NULL;

-- Timetable entries table - individual period assignments
CREATE TABLE timetable_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    timetable_id UUID NOT NULL REFERENCES timetables(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    period_slot_id UUID NOT NULL REFERENCES period_slots(id),
    subject_id UUID REFERENCES subjects(id),
    staff_id UUID REFERENCES staff(id),
    room_number VARCHAR(50),
    notes TEXT,
    is_free_period BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_timetable_entry UNIQUE (timetable_id, day_of_week, period_slot_id)
);

-- Enable Row Level Security
ALTER TABLE timetable_entries ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_timetable_entries ON timetable_entries
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for timetable_entries
CREATE INDEX idx_timetable_entries_tenant ON timetable_entries(tenant_id);
CREATE INDEX idx_timetable_entries_timetable ON timetable_entries(timetable_id);
CREATE INDEX idx_timetable_entries_day ON timetable_entries(day_of_week);
CREATE INDEX idx_timetable_entries_period ON timetable_entries(period_slot_id);
CREATE INDEX idx_timetable_entries_subject ON timetable_entries(subject_id);
CREATE INDEX idx_timetable_entries_staff ON timetable_entries(staff_id);

-- Index for finding teacher conflicts
CREATE INDEX idx_timetable_entries_staff_day_period ON timetable_entries(staff_id, day_of_week, period_slot_id);

-- Add timetable permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('timetables:read', 'View Timetables', 'Permission to view timetables', 'academics', NOW(), NOW()),
    ('timetables:create', 'Create Timetables', 'Permission to create timetables', 'academics', NOW(), NOW()),
    ('timetables:update', 'Update Timetables', 'Permission to update timetables', 'academics', NOW(), NOW()),
    ('timetables:delete', 'Delete Timetables', 'Permission to delete timetables', 'academics', NOW(), NOW()),
    ('timetables:publish', 'Publish Timetables', 'Permission to publish timetables', 'academics', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Trigger for updated_at
CREATE TRIGGER trigger_timetables_updated_at
    BEFORE UPDATE ON timetables
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_timetable_entries_updated_at
    BEFORE UPDATE ON timetable_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
