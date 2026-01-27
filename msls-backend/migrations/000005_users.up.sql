-- Migration: 000005_users.up.sql
-- Description: Create users table for tenant user accounts

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255),
    phone VARCHAR(20),
    password_hash VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    two_factor_enabled BOOLEAN NOT NULL DEFAULT false,
    two_factor_secret VARCHAR(255),
    email_verified_at TIMESTAMPTZ,
    phone_verified_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_by UUID,

    -- Constraints
    CONSTRAINT users_tenant_email_unique UNIQUE(tenant_id, email),
    CONSTRAINT users_tenant_phone_unique UNIQUE(tenant_id, phone),
    CONSTRAINT users_status_check CHECK (status IN ('active', 'inactive', 'suspended', 'pending')),
    CONSTRAINT users_email_or_phone_required CHECK (email IS NOT NULL OR phone IS NOT NULL)
);

-- Add comments
COMMENT ON TABLE users IS 'User accounts belonging to tenants';
COMMENT ON COLUMN users.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN users.tenant_id IS 'Reference to the parent tenant';
COMMENT ON COLUMN users.email IS 'User email address (unique per tenant)';
COMMENT ON COLUMN users.phone IS 'User phone number (unique per tenant)';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.status IS 'User status: active, inactive, suspended, or pending';
COMMENT ON COLUMN users.two_factor_enabled IS 'Whether 2FA is enabled for this user';
COMMENT ON COLUMN users.two_factor_secret IS 'Encrypted TOTP secret for 2FA';
COMMENT ON COLUMN users.email_verified_at IS 'Timestamp when email was verified';
COMMENT ON COLUMN users.phone_verified_at IS 'Timestamp when phone was verified';
COMMENT ON COLUMN users.last_login_at IS 'Timestamp of last successful login';
COMMENT ON COLUMN users.created_by IS 'User who created this record';
COMMENT ON COLUMN users.updated_by IS 'User who last updated this record';
