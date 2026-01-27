-- Migration: 000003_tenants.up.sql
-- Description: Create tenants table for multi-tenancy foundation

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    settings JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT tenants_status_check CHECK (status IN ('active', 'inactive', 'suspended'))
);

-- Add indexes
CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_status ON tenants(status);

-- Add comments
COMMENT ON TABLE tenants IS 'Multi-tenant organizations/companies using the MSLS system';
COMMENT ON COLUMN tenants.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN tenants.name IS 'Display name of the tenant organization';
COMMENT ON COLUMN tenants.slug IS 'URL-friendly unique identifier for the tenant';
COMMENT ON COLUMN tenants.settings IS 'Tenant-specific configuration as JSON';
COMMENT ON COLUMN tenants.status IS 'Tenant status: active, inactive, or suspended';
