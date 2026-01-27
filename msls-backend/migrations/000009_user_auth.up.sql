-- Migration: 000009_user_auth.up.sql
-- Description: Create authentication-related tables (user_roles, refresh_tokens, verification_tokens, login_attempts, audit_logs)

-- Create user_roles junction table
CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT user_roles_unique UNIQUE(user_id, role_id)
);

-- Create refresh_tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,

    -- Index for fast lookup
    CONSTRAINT refresh_tokens_hash_unique UNIQUE(token_hash)
);

-- Create verification_tokens table
CREATE TABLE verification_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT verification_tokens_hash_unique UNIQUE(token_hash),
    CONSTRAINT verification_tokens_type_check CHECK (type IN ('email_verify', 'password_reset', 'phone_verify'))
);

-- Create login_attempts table
CREATE TABLE login_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    email VARCHAR(255),
    ip_address INET NOT NULL,
    user_agent TEXT,
    success BOOLEAN NOT NULL DEFAULT false,
    failure_reason VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create audit_logs table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,
    old_data JSONB,
    new_data JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

CREATE INDEX idx_verification_tokens_user_id ON verification_tokens(user_id);
CREATE INDEX idx_verification_tokens_token_hash ON verification_tokens(token_hash);
CREATE INDEX idx_verification_tokens_type ON verification_tokens(type);
CREATE INDEX idx_verification_tokens_expires_at ON verification_tokens(expires_at);

CREATE INDEX idx_login_attempts_user_id ON login_attempts(user_id);
CREATE INDEX idx_login_attempts_email ON login_attempts(email);
CREATE INDEX idx_login_attempts_ip_address ON login_attempts(ip_address);
CREATE INDEX idx_login_attempts_created_at ON login_attempts(created_at);

CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Add comments
COMMENT ON TABLE user_roles IS 'Junction table linking users to roles';
COMMENT ON COLUMN user_roles.user_id IS 'Reference to the user';
COMMENT ON COLUMN user_roles.role_id IS 'Reference to the role';

COMMENT ON TABLE refresh_tokens IS 'Stores hashed refresh tokens for JWT authentication';
COMMENT ON COLUMN refresh_tokens.user_id IS 'Reference to the user who owns this token';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hash of the refresh token';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'When this refresh token expires';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'When this token was revoked (null if active)';

COMMENT ON TABLE verification_tokens IS 'Stores tokens for email verification and password reset';
COMMENT ON COLUMN verification_tokens.user_id IS 'Reference to the user';
COMMENT ON COLUMN verification_tokens.token_hash IS 'SHA-256 hash of the verification token';
COMMENT ON COLUMN verification_tokens.type IS 'Type of verification: email_verify, password_reset, phone_verify';
COMMENT ON COLUMN verification_tokens.expires_at IS 'When this token expires';
COMMENT ON COLUMN verification_tokens.used_at IS 'When this token was used (null if unused)';

COMMENT ON TABLE login_attempts IS 'Tracks login attempts for security monitoring and rate limiting';
COMMENT ON COLUMN login_attempts.user_id IS 'Reference to the user (null if user not found)';
COMMENT ON COLUMN login_attempts.email IS 'Email used in the login attempt';
COMMENT ON COLUMN login_attempts.ip_address IS 'IP address of the request';
COMMENT ON COLUMN login_attempts.user_agent IS 'User agent string of the request';
COMMENT ON COLUMN login_attempts.success IS 'Whether the login attempt was successful';
COMMENT ON COLUMN login_attempts.failure_reason IS 'Reason for failure if unsuccessful';

COMMENT ON TABLE audit_logs IS 'Audit trail for tracking changes and actions';
COMMENT ON COLUMN audit_logs.tenant_id IS 'Reference to the tenant';
COMMENT ON COLUMN audit_logs.user_id IS 'Reference to the user who performed the action';
COMMENT ON COLUMN audit_logs.action IS 'Action performed (create, update, delete, login, logout, etc.)';
COMMENT ON COLUMN audit_logs.entity_type IS 'Type of entity affected (user, student, etc.)';
COMMENT ON COLUMN audit_logs.entity_id IS 'ID of the affected entity';
COMMENT ON COLUMN audit_logs.old_data IS 'Previous state of the entity (for updates)';
COMMENT ON COLUMN audit_logs.new_data IS 'New state of the entity (for creates/updates)';
COMMENT ON COLUMN audit_logs.ip_address IS 'IP address of the request';
COMMENT ON COLUMN audit_logs.user_agent IS 'User agent string of the request';

-- Enable RLS on audit_logs
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- Create RLS policy for audit_logs table
CREATE POLICY tenant_isolation_audit_logs ON audit_logs
    FOR ALL
    USING (
        tenant_id IS NULL
        OR tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_audit_logs ON audit_logs
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

COMMENT ON POLICY tenant_isolation_audit_logs ON audit_logs IS 'Restricts audit log access to current tenant only';
COMMENT ON POLICY bypass_rls_audit_logs ON audit_logs IS 'Allows bypass of RLS for admin operations';

-- Add first_name and last_name to users table if not exists
ALTER TABLE users ADD COLUMN IF NOT EXISTS first_name VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_name VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS failed_login_attempts INTEGER NOT NULL DEFAULT 0;

COMMENT ON COLUMN users.first_name IS 'User first name';
COMMENT ON COLUMN users.last_name IS 'User last name';
COMMENT ON COLUMN users.locked_until IS 'Account locked until this time (null if not locked)';
COMMENT ON COLUMN users.failed_login_attempts IS 'Number of consecutive failed login attempts';
