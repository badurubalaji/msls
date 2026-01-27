-- Migration: 000010_totp_2fa.down.sql
-- Description: Remove TOTP-based two-factor authentication support

-- Drop index for users with 2FA enabled
DROP INDEX IF EXISTS idx_users_two_factor_enabled;

-- Drop TOTP attempts table
DROP TABLE IF EXISTS totp_attempts;

-- Drop backup codes table
DROP TABLE IF EXISTS backup_codes;

-- Remove TOTP verification timestamp column
ALTER TABLE users DROP COLUMN IF EXISTS totp_verified_at;
