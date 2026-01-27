-- Migration: 000011_otp_codes.down.sql
-- Description: Drop OTP codes table and related objects

-- Drop functions
DROP FUNCTION IF EXISTS cleanup_expired_otp_codes();
DROP FUNCTION IF EXISTS cleanup_old_rate_limits();

-- Drop indexes
DROP INDEX IF EXISTS idx_otp_rate_limits_window_start;
DROP INDEX IF EXISTS idx_otp_rate_limits_identifier;
DROP INDEX IF EXISTS idx_otp_codes_identifier_channel;
DROP INDEX IF EXISTS idx_otp_codes_created_at;
DROP INDEX IF EXISTS idx_otp_codes_expires_at;
DROP INDEX IF EXISTS idx_otp_codes_identifier;
DROP INDEX IF EXISTS idx_otp_codes_user_id;

-- Drop tables
DROP TABLE IF EXISTS otp_rate_limits;
DROP TABLE IF EXISTS otp_codes;

-- Drop enums
DROP TYPE IF EXISTS otp_channel;
DROP TYPE IF EXISTS otp_type;
