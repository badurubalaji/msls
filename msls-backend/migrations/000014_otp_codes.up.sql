-- Migration: 000011_otp_codes.up.sql
-- Description: Create OTP codes table for passwordless authentication

-- Create OTP type enum
CREATE TYPE otp_type AS ENUM ('login', 'verify', 'phone_verify');

-- Create OTP channel enum
CREATE TYPE otp_channel AS ENUM ('sms', 'email');

-- Create otp_codes table
CREATE TABLE otp_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    identifier VARCHAR(255) NOT NULL,
    code_hash VARCHAR(255) NOT NULL,
    type otp_type NOT NULL,
    channel otp_channel NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ,
    attempts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT otp_codes_max_attempts CHECK (attempts <= 5)
);

-- Create otp_rate_limits table for tracking OTP requests per identifier
CREATE TABLE otp_rate_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    identifier VARCHAR(255) NOT NULL,
    channel otp_channel NOT NULL,
    request_count INTEGER NOT NULL DEFAULT 1,
    window_start TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_request_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT otp_rate_limits_identifier_channel_unique UNIQUE(identifier, channel)
);

-- Create indexes
CREATE INDEX idx_otp_codes_user_id ON otp_codes(user_id);
CREATE INDEX idx_otp_codes_identifier ON otp_codes(identifier);
CREATE INDEX idx_otp_codes_expires_at ON otp_codes(expires_at);
CREATE INDEX idx_otp_codes_created_at ON otp_codes(created_at);
CREATE INDEX idx_otp_codes_identifier_channel ON otp_codes(identifier, channel);

CREATE INDEX idx_otp_rate_limits_identifier ON otp_rate_limits(identifier);
CREATE INDEX idx_otp_rate_limits_window_start ON otp_rate_limits(window_start);

-- Add comments
COMMENT ON TABLE otp_codes IS 'Stores OTP codes for passwordless authentication';
COMMENT ON COLUMN otp_codes.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN otp_codes.user_id IS 'Reference to the user (null if user not found yet)';
COMMENT ON COLUMN otp_codes.identifier IS 'Phone number or email address';
COMMENT ON COLUMN otp_codes.code_hash IS 'SHA-256 hash of the OTP code';
COMMENT ON COLUMN otp_codes.type IS 'Type of OTP: login, verify, phone_verify';
COMMENT ON COLUMN otp_codes.channel IS 'Delivery channel: sms or email';
COMMENT ON COLUMN otp_codes.expires_at IS 'When this OTP expires (5 minutes from creation)';
COMMENT ON COLUMN otp_codes.verified_at IS 'When this OTP was successfully verified';
COMMENT ON COLUMN otp_codes.attempts IS 'Number of validation attempts (max 5)';

COMMENT ON TABLE otp_rate_limits IS 'Tracks OTP request rate limits per identifier';
COMMENT ON COLUMN otp_rate_limits.identifier IS 'Phone number or email address';
COMMENT ON COLUMN otp_rate_limits.channel IS 'Delivery channel: sms or email';
COMMENT ON COLUMN otp_rate_limits.request_count IS 'Number of OTP requests in current window';
COMMENT ON COLUMN otp_rate_limits.window_start IS 'Start of the current rate limit window';
COMMENT ON COLUMN otp_rate_limits.last_request_at IS 'Timestamp of last OTP request';

-- Create function to clean up expired OTP codes
CREATE OR REPLACE FUNCTION cleanup_expired_otp_codes()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM otp_codes
    WHERE expires_at < NOW() - INTERVAL '1 hour'
       OR verified_at IS NOT NULL;
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_expired_otp_codes() IS 'Removes expired and used OTP codes older than 1 hour';

-- Create function to reset rate limit windows older than 1 hour
CREATE OR REPLACE FUNCTION cleanup_old_rate_limits()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM otp_rate_limits
    WHERE window_start < NOW() - INTERVAL '2 hours';
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_old_rate_limits() IS 'Removes rate limit records older than 2 hours';
