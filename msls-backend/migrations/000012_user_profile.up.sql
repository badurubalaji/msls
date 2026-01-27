-- Migration: 000010_user_profile.up.sql
-- Description: Add profile fields to users table and create user_preferences table

-- Add new columns to users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS first_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS last_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500),
    ADD COLUMN IF NOT EXISTS bio TEXT,
    ADD COLUMN IF NOT EXISTS timezone VARCHAR(50) DEFAULT 'UTC',
    ADD COLUMN IF NOT EXISTS locale VARCHAR(10) DEFAULT 'en',
    ADD COLUMN IF NOT EXISTS notification_preferences JSONB DEFAULT '{"email": true, "push": true, "sms": false}'::jsonb,
    ADD COLUMN IF NOT EXISTS last_login_ip INET,
    ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS failed_login_attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS account_deletion_requested_at TIMESTAMPTZ;

-- Create user_preferences table for extended settings
CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category VARCHAR(50) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Unique constraint for user + category + key
    CONSTRAINT user_preferences_unique UNIQUE (user_id, category, key)
);

-- Create index on user_preferences
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_user_preferences_category ON user_preferences(user_id, category);

-- Add RLS policy for user_preferences
ALTER TABLE user_preferences ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only access their own preferences
CREATE POLICY user_preferences_isolation_policy ON user_preferences
    USING (user_id IN (
        SELECT id FROM users WHERE tenant_id = current_setting('app.current_tenant_id', true)::uuid
    ));

-- Add comments
COMMENT ON COLUMN users.first_name IS 'User first name';
COMMENT ON COLUMN users.last_name IS 'User last name';
COMMENT ON COLUMN users.avatar_url IS 'URL to user avatar image';
COMMENT ON COLUMN users.bio IS 'User biography or description';
COMMENT ON COLUMN users.timezone IS 'User preferred timezone (IANA format)';
COMMENT ON COLUMN users.locale IS 'User preferred locale/language code';
COMMENT ON COLUMN users.notification_preferences IS 'JSONB object with notification preferences (email, push, sms)';
COMMENT ON COLUMN users.last_login_ip IS 'IP address of last login';
COMMENT ON COLUMN users.locked_until IS 'Account locked until this timestamp';
COMMENT ON COLUMN users.failed_login_attempts IS 'Count of consecutive failed login attempts';
COMMENT ON COLUMN users.account_deletion_requested_at IS 'Timestamp when account deletion was requested';

COMMENT ON TABLE user_preferences IS 'Extended user preferences stored as key-value pairs';
COMMENT ON COLUMN user_preferences.user_id IS 'Reference to the user';
COMMENT ON COLUMN user_preferences.category IS 'Preference category (e.g., appearance, accessibility)';
COMMENT ON COLUMN user_preferences.key IS 'Preference key within the category';
COMMENT ON COLUMN user_preferences.value IS 'Preference value as JSONB';
