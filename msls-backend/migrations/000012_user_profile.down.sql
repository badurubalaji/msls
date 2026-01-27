-- Migration: 000010_user_profile.down.sql
-- Description: Remove profile fields from users table and drop user_preferences table

-- Drop RLS policy
DROP POLICY IF EXISTS user_preferences_isolation_policy ON user_preferences;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP INDEX IF EXISTS idx_user_preferences_category;

-- Drop user_preferences table
DROP TABLE IF EXISTS user_preferences;

-- Remove columns from users table
ALTER TABLE users
    DROP COLUMN IF EXISTS first_name,
    DROP COLUMN IF EXISTS last_name,
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS bio,
    DROP COLUMN IF EXISTS timezone,
    DROP COLUMN IF EXISTS locale,
    DROP COLUMN IF EXISTS notification_preferences,
    DROP COLUMN IF EXISTS last_login_ip,
    DROP COLUMN IF EXISTS locked_until,
    DROP COLUMN IF EXISTS failed_login_attempts,
    DROP COLUMN IF EXISTS account_deletion_requested_at;
