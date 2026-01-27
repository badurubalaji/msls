-- Migration: 000034_departments_designations.up.sql
-- Description: Create departments and designations tables for staff management

-- Departments table
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    head_id UUID, -- Will reference staff(id) after staff table is created
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_department_code UNIQUE (tenant_id, branch_id, code)
);

-- Enable RLS
ALTER TABLE departments ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_departments ON departments
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_departments_tenant ON departments(tenant_id);
CREATE INDEX idx_departments_branch ON departments(branch_id);
CREATE INDEX idx_departments_active ON departments(tenant_id, is_active) WHERE is_active = true;

-- Designations table
CREATE TABLE designations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    level INTEGER NOT NULL DEFAULT 1,
    department_id UUID REFERENCES departments(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_designation_name UNIQUE (tenant_id, name)
);

-- Enable RLS
ALTER TABLE designations ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_designations ON designations
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_designations_tenant ON designations(tenant_id);
CREATE INDEX idx_designations_department ON designations(department_id);
CREATE INDEX idx_designations_active ON designations(tenant_id, is_active) WHERE is_active = true;

-- Updated_at triggers
CREATE OR REPLACE FUNCTION update_departments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_departments_updated_at
    BEFORE UPDATE ON departments
    FOR EACH ROW
    EXECUTE FUNCTION update_departments_updated_at();

CREATE OR REPLACE FUNCTION update_designations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_designations_updated_at
    BEFORE UPDATE ON designations
    FOR EACH ROW
    EXECUTE FUNCTION update_designations_updated_at();

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('department:read', 'View Departments', 'Permission to view departments', 'staff', NOW(), NOW()),
    ('department:create', 'Create Departments', 'Permission to create departments', 'staff', NOW(), NOW()),
    ('department:update', 'Update Departments', 'Permission to update departments', 'staff', NOW(), NOW()),
    ('department:delete', 'Delete Departments', 'Permission to delete departments', 'staff', NOW(), NOW()),
    ('designation:read', 'View Designations', 'Permission to view designations', 'staff', NOW(), NOW()),
    ('designation:create', 'Create Designations', 'Permission to create designations', 'staff', NOW(), NOW()),
    ('designation:update', 'Update Designations', 'Permission to update designations', 'staff', NOW(), NOW()),
    ('designation:delete', 'Delete Designations', 'Permission to delete designations', 'staff', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;
