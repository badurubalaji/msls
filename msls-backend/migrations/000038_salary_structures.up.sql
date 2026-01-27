-- Migration: 000038_salary_structures.up.sql
-- Description: Create salary components, structures, and staff salary tables

-- Salary Components (earning/deduction types)
CREATE TABLE salary_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    component_type VARCHAR(20) NOT NULL CHECK (component_type IN ('earning', 'deduction')),
    calculation_type VARCHAR(20) NOT NULL CHECK (calculation_type IN ('fixed', 'percentage')),
    percentage_of_id UUID REFERENCES salary_components(id),
    is_taxable BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_salary_component_code UNIQUE (tenant_id, code)
);

-- Enable RLS
ALTER TABLE salary_components ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_salary_components ON salary_components
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_salary_components_tenant ON salary_components(tenant_id);
CREATE INDEX idx_salary_components_type ON salary_components(tenant_id, component_type);
CREATE INDEX idx_salary_components_active ON salary_components(tenant_id, is_active) WHERE is_active = true;

-- Salary Structures (templates)
CREATE TABLE salary_structures (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    designation_id UUID REFERENCES designations(id),
    is_active BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_salary_structure_code UNIQUE (tenant_id, code)
);

-- Enable RLS
ALTER TABLE salary_structures ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_salary_structures ON salary_structures
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_salary_structures_tenant ON salary_structures(tenant_id);
CREATE INDEX idx_salary_structures_designation ON salary_structures(designation_id);
CREATE INDEX idx_salary_structures_active ON salary_structures(tenant_id, is_active) WHERE is_active = true;

-- Structure Components (components in a structure with default values)
CREATE TABLE salary_structure_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    structure_id UUID NOT NULL REFERENCES salary_structures(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES salary_components(id),

    amount DECIMAL(12,2),
    percentage DECIMAL(5,2),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_structure_component UNIQUE (structure_id, component_id)
);

-- Index
CREATE INDEX idx_structure_components_structure ON salary_structure_components(structure_id);

-- Staff Salary Assignment
CREATE TABLE staff_salaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id),
    structure_id UUID REFERENCES salary_structures(id),

    effective_from DATE NOT NULL,
    effective_to DATE,

    gross_salary DECIMAL(12,2) NOT NULL,
    net_salary DECIMAL(12,2) NOT NULL,
    ctc DECIMAL(12,2),

    revision_reason TEXT,
    is_current BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Enable RLS
ALTER TABLE staff_salaries ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_staff_salaries ON staff_salaries
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_staff_salaries_tenant ON staff_salaries(tenant_id);
CREATE INDEX idx_staff_salaries_staff ON staff_salaries(staff_id);
CREATE INDEX idx_staff_salaries_current ON staff_salaries(staff_id, is_current) WHERE is_current = true;

-- Staff Salary Components (actual values for staff)
CREATE TABLE staff_salary_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    staff_salary_id UUID NOT NULL REFERENCES staff_salaries(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES salary_components(id),

    amount DECIMAL(12,2) NOT NULL,
    is_overridden BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_staff_salary_component UNIQUE (staff_salary_id, component_id)
);

-- Index
CREATE INDEX idx_staff_salary_components_salary ON staff_salary_components(staff_salary_id);

-- Updated_at triggers
CREATE OR REPLACE FUNCTION update_salary_components_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_salary_components_updated_at
    BEFORE UPDATE ON salary_components
    FOR EACH ROW
    EXECUTE FUNCTION update_salary_components_updated_at();

CREATE OR REPLACE FUNCTION update_salary_structures_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_salary_structures_updated_at
    BEFORE UPDATE ON salary_structures
    FOR EACH ROW
    EXECUTE FUNCTION update_salary_structures_updated_at();

CREATE OR REPLACE FUNCTION update_staff_salaries_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_staff_salaries_updated_at
    BEFORE UPDATE ON staff_salaries
    FOR EACH ROW
    EXECUTE FUNCTION update_staff_salaries_updated_at();

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('salary:read', 'View Salary', 'Permission to view salary structures and staff salaries', 'staff', NOW(), NOW()),
    ('salary:create', 'Create Salary', 'Permission to create salary structures', 'staff', NOW(), NOW()),
    ('salary:update', 'Update Salary', 'Permission to update salary structures', 'staff', NOW(), NOW()),
    ('salary:delete', 'Delete Salary', 'Permission to delete salary structures', 'staff', NOW(), NOW()),
    ('salary:assign', 'Assign Salary', 'Permission to assign salaries to staff', 'staff', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign salary permissions to super_admin role
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN ('salary:read', 'salary:create', 'salary:update', 'salary:delete', 'salary:assign')
ON CONFLICT DO NOTHING;
