-- Migration: 000004_branches.up.sql
-- Description: Create branches table for tenant branch locations

CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    address JSONB NOT NULL DEFAULT '{}',
    settings JSONB NOT NULL DEFAULT '{}',
    is_primary BOOLEAN NOT NULL DEFAULT false,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_by UUID,

    -- Constraints
    CONSTRAINT branches_tenant_code_unique UNIQUE(tenant_id, code),
    CONSTRAINT branches_status_check CHECK (status IN ('active', 'inactive', 'suspended'))
);

-- Function to ensure only one primary branch per tenant
CREATE OR REPLACE FUNCTION ensure_single_primary_branch()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_primary = true THEN
        UPDATE branches
        SET is_primary = false, updated_at = NOW()
        WHERE tenant_id = NEW.tenant_id
          AND id != NEW.id
          AND is_primary = true;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to enforce single primary branch
CREATE TRIGGER trigger_single_primary_branch
    BEFORE INSERT OR UPDATE ON branches
    FOR EACH ROW
    EXECUTE FUNCTION ensure_single_primary_branch();

-- Add comments
COMMENT ON TABLE branches IS 'Physical or logical branches/locations belonging to a tenant';
COMMENT ON COLUMN branches.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN branches.tenant_id IS 'Reference to the parent tenant';
COMMENT ON COLUMN branches.name IS 'Display name of the branch';
COMMENT ON COLUMN branches.code IS 'Unique code within tenant for the branch';
COMMENT ON COLUMN branches.address IS 'Branch address details as JSON';
COMMENT ON COLUMN branches.settings IS 'Branch-specific configuration as JSON';
COMMENT ON COLUMN branches.is_primary IS 'Whether this is the primary/main branch';
COMMENT ON COLUMN branches.status IS 'Branch status: active, inactive, or suspended';
COMMENT ON COLUMN branches.created_by IS 'User who created this record';
COMMENT ON COLUMN branches.updated_by IS 'User who last updated this record';
