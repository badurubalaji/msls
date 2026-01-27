-- Migration: 000010_feature_flags.up.sql
-- Description: Create feature flags system for gradual feature rollout and tenant-specific configurations

-- Feature flags master table
CREATE TABLE feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    key VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    default_value BOOLEAN NOT NULL DEFAULT false,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tenant-specific feature flag overrides
CREATE TABLE tenant_feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    flag_id UUID NOT NULL REFERENCES feature_flags(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL,
    custom_value JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT tenant_feature_flags_unique UNIQUE (tenant_id, flag_id)
);

-- User-specific feature flag overrides (for beta testing)
CREATE TABLE user_feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    flag_id UUID NOT NULL REFERENCES feature_flags(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT user_feature_flags_unique UNIQUE (user_id, flag_id)
);

-- Add indexes for performance
CREATE INDEX idx_feature_flags_key ON feature_flags(key);
CREATE INDEX idx_tenant_feature_flags_tenant_id ON tenant_feature_flags(tenant_id);
CREATE INDEX idx_tenant_feature_flags_flag_id ON tenant_feature_flags(flag_id);
CREATE INDEX idx_user_feature_flags_user_id ON user_feature_flags(user_id);
CREATE INDEX idx_user_feature_flags_flag_id ON user_feature_flags(flag_id);

-- Add Row-Level Security for tenant_feature_flags
ALTER TABLE tenant_feature_flags ENABLE ROW LEVEL SECURITY;

-- RLS policy for tenant_feature_flags: users can only see their tenant's flags
CREATE POLICY tenant_feature_flags_tenant_isolation ON tenant_feature_flags
    FOR ALL
    USING (tenant_id::text = current_setting('app.tenant_id', true));

-- Add comments
COMMENT ON TABLE feature_flags IS 'Master list of feature flags for the system';
COMMENT ON COLUMN feature_flags.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN feature_flags.key IS 'Unique key used to identify the flag in code';
COMMENT ON COLUMN feature_flags.name IS 'Human-readable name of the feature flag';
COMMENT ON COLUMN feature_flags.description IS 'Description of what the feature flag controls';
COMMENT ON COLUMN feature_flags.default_value IS 'Default enabled state when no override exists';
COMMENT ON COLUMN feature_flags.metadata IS 'Additional metadata for the flag (e.g., rollout percentage)';

COMMENT ON TABLE tenant_feature_flags IS 'Tenant-specific feature flag overrides';
COMMENT ON COLUMN tenant_feature_flags.tenant_id IS 'Reference to the tenant';
COMMENT ON COLUMN tenant_feature_flags.flag_id IS 'Reference to the feature flag';
COMMENT ON COLUMN tenant_feature_flags.enabled IS 'Whether the flag is enabled for this tenant';
COMMENT ON COLUMN tenant_feature_flags.custom_value IS 'Custom JSON value for the flag (for non-boolean flags)';

COMMENT ON TABLE user_feature_flags IS 'User-specific feature flag overrides for beta testing';
COMMENT ON COLUMN user_feature_flags.user_id IS 'Reference to the user';
COMMENT ON COLUMN user_feature_flags.flag_id IS 'Reference to the feature flag';
COMMENT ON COLUMN user_feature_flags.enabled IS 'Whether the flag is enabled for this user';

-- Insert predefined feature flags
INSERT INTO feature_flags (key, name, description, default_value, metadata) VALUES
    ('online_admissions', 'Online Admissions', 'Enable online student admissions and enrollment', true, '{"category": "admissions"}'),
    ('transport_tracking', 'Transport Tracking', 'Enable GPS tracking for school transport vehicles', false, '{"category": "transport", "requires_setup": true}'),
    ('ai_insights', 'AI Insights', 'Enable AI-powered analytics and insights dashboard', false, '{"category": "analytics", "beta": true}'),
    ('parent_messaging', 'Parent Messaging', 'Enable direct messaging between staff and parents', true, '{"category": "communication"}'),
    ('student_portal', 'Student Portal', 'Enable self-service student portal access', true, '{"category": "portal"}');
