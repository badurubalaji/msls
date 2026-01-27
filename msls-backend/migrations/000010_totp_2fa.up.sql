-- Migration: 000010_totp_2fa.up.sql
-- Description: Add TOTP-based two-factor authentication support

-- Add TOTP verification timestamp to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS totp_verified_at TIMESTAMPTZ;

COMMENT ON COLUMN users.totp_verified_at IS 'Timestamp when TOTP 2FA was verified and enabled';

-- Create backup codes table for 2FA account recovery
CREATE TABLE backup_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash VARCHAR(255) NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT backup_codes_hash_unique UNIQUE(code_hash)
);

-- Create indexes
CREATE INDEX idx_backup_codes_user_id ON backup_codes(user_id);
CREATE INDEX idx_backup_codes_code_hash ON backup_codes(code_hash);

-- Add comments
COMMENT ON TABLE backup_codes IS 'Backup codes for 2FA account recovery';
COMMENT ON COLUMN backup_codes.user_id IS 'Reference to the user who owns these backup codes';
COMMENT ON COLUMN backup_codes.code_hash IS 'SHA-256 hash of the backup code';
COMMENT ON COLUMN backup_codes.used_at IS 'Timestamp when this code was used (null if unused)';

-- Create 2FA attempts table for rate limiting
CREATE TABLE totp_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    ip_address INET NOT NULL,
    success BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for TOTP attempts
CREATE INDEX idx_totp_attempts_user_id ON totp_attempts(user_id);
CREATE INDEX idx_totp_attempts_ip_address ON totp_attempts(ip_address);
CREATE INDEX idx_totp_attempts_created_at ON totp_attempts(created_at);

-- Add comments for TOTP attempts
COMMENT ON TABLE totp_attempts IS 'Tracks 2FA validation attempts for rate limiting and security monitoring';
COMMENT ON COLUMN totp_attempts.user_id IS 'Reference to the user attempting 2FA validation';
COMMENT ON COLUMN totp_attempts.ip_address IS 'IP address of the request';
COMMENT ON COLUMN totp_attempts.success IS 'Whether the 2FA attempt was successful';

-- Create index for users with 2FA enabled for efficient lookups
CREATE INDEX idx_users_two_factor_enabled ON users(two_factor_enabled) WHERE two_factor_enabled = true;
