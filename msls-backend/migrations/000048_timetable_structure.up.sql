-- Migration: 000048_timetable_structure.up.sql
-- Description: Create timetable structure tables (shifts, day patterns, period slots)

-- Shifts (morning, afternoon, etc.)
CREATE TABLE shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    name VARCHAR(50) NOT NULL,              -- e.g., "Morning Shift", "Afternoon Shift"
    code VARCHAR(20) NOT NULL,              -- e.g., "MORN", "AFT"
    start_time TIME NOT NULL,               -- e.g., 08:00
    end_time TIME NOT NULL,                 -- e.g., 14:00
    description TEXT,
    display_order INTEGER NOT NULL DEFAULT 0,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    CONSTRAINT uniq_shift_code UNIQUE (tenant_id, branch_id, code)
);

-- Enable RLS
ALTER TABLE shifts ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_shifts ON shifts
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_shifts_tenant ON shifts(tenant_id);
CREATE INDEX idx_shifts_branch ON shifts(branch_id);
CREATE INDEX idx_shifts_active ON shifts(tenant_id, is_active);

-- Day Patterns (regular, half-day, etc.)
CREATE TABLE day_patterns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(50) NOT NULL,              -- e.g., "Regular Day", "Half Day", "Saturday"
    code VARCHAR(20) NOT NULL,              -- e.g., "REG", "HALF", "SAT"
    description TEXT,
    total_periods INTEGER NOT NULL DEFAULT 8,
    display_order INTEGER NOT NULL DEFAULT 0,

    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    CONSTRAINT uniq_day_pattern_code UNIQUE (tenant_id, code)
);

-- Enable RLS
ALTER TABLE day_patterns ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_day_patterns ON day_patterns
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_day_patterns_tenant ON day_patterns(tenant_id);
CREATE INDEX idx_day_patterns_active ON day_patterns(tenant_id, is_active);

-- Day Pattern Assignments (which pattern applies to which day of week)
CREATE TABLE day_pattern_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),  -- 0=Sunday, 1=Monday, ..., 6=Saturday
    day_pattern_id UUID REFERENCES day_patterns(id) ON DELETE SET NULL,
    is_working_day BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_day_assignment UNIQUE (tenant_id, branch_id, day_of_week)
);

-- Enable RLS
ALTER TABLE day_pattern_assignments ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_day_pattern_assignments ON day_pattern_assignments
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_day_pattern_assignments_tenant ON day_pattern_assignments(tenant_id);
CREATE INDEX idx_day_pattern_assignments_branch ON day_pattern_assignments(branch_id);

-- Period Slot Types
CREATE TYPE period_slot_type AS ENUM (
    'regular',      -- Regular teaching period
    'short',        -- Shortened period
    'assembly',     -- Morning assembly
    'break',        -- Short break
    'lunch',        -- Lunch break
    'activity',     -- Extra-curricular activity
    'zero_period'   -- Before regular school hours
);

-- Period Slots
CREATE TABLE period_slots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    name VARCHAR(50) NOT NULL,              -- e.g., "Period 1", "Assembly", "Lunch Break"
    period_number INTEGER,                   -- Null for breaks, 1-based for teaching periods
    slot_type period_slot_type NOT NULL DEFAULT 'regular',
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    duration_minutes INTEGER NOT NULL,

    -- Link to day pattern (slots can vary by day pattern)
    day_pattern_id UUID REFERENCES day_patterns(id) ON DELETE CASCADE,

    -- Link to shift (for multi-shift schools)
    shift_id UUID REFERENCES shifts(id) ON DELETE SET NULL,

    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    CONSTRAINT chk_period_times CHECK (end_time > start_time)
);

-- Enable RLS
ALTER TABLE period_slots ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_period_slots ON period_slots
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_period_slots_tenant ON period_slots(tenant_id);
CREATE INDEX idx_period_slots_branch ON period_slots(branch_id);
CREATE INDEX idx_period_slots_day_pattern ON period_slots(day_pattern_id);
CREATE INDEX idx_period_slots_shift ON period_slots(shift_id);
CREATE INDEX idx_period_slots_active ON period_slots(tenant_id, is_active);

-- Add updated_at triggers
CREATE TRIGGER trigger_shifts_updated_at
    BEFORE UPDATE ON shifts
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

CREATE TRIGGER trigger_day_patterns_updated_at
    BEFORE UPDATE ON day_patterns
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

CREATE TRIGGER trigger_day_pattern_assignments_updated_at
    BEFORE UPDATE ON day_pattern_assignments
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

CREATE TRIGGER trigger_period_slots_updated_at
    BEFORE UPDATE ON period_slots
    FOR EACH ROW EXECUTE FUNCTION update_classes_updated_at();

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('timetable:view', 'View Timetable', 'Permission to view timetables', 'academic', NOW(), NOW()),
    ('timetable:create', 'Create Timetable', 'Permission to create/edit timetables', 'academic', NOW(), NOW()),
    ('timetable:manage', 'Manage Timetable Settings', 'Permission to manage period slots, shifts, patterns', 'academic', NOW(), NOW()),
    ('shift:view', 'View Shifts', 'Permission to view shifts', 'academic', NOW(), NOW()),
    ('shift:manage', 'Manage Shifts', 'Permission to manage shifts', 'academic', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT r.id, p.id, NOW(), NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN ('timetable:view', 'timetable:create', 'timetable:manage', 'shift:view', 'shift:manage')
ON CONFLICT DO NOTHING;

-- Add shift_id to sections for multi-shift schools
ALTER TABLE sections ADD COLUMN IF NOT EXISTS shift_id UUID REFERENCES shifts(id);
CREATE INDEX IF NOT EXISTS idx_sections_shift ON sections(shift_id);
