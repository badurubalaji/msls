-- Migration: 000009_user_auth.down.sql
-- Description: Drop authentication-related tables

-- Drop RLS policies
DROP POLICY IF EXISTS bypass_rls_audit_logs ON audit_logs;
DROP POLICY IF EXISTS tenant_isolation_audit_logs ON audit_logs;

-- Disable RLS
ALTER TABLE audit_logs DISABLE ROW LEVEL SECURITY;

-- Remove columns added to users table
ALTER TABLE users DROP COLUMN IF EXISTS first_name;
ALTER TABLE users DROP COLUMN IF EXISTS last_name;
ALTER TABLE users DROP COLUMN IF EXISTS locked_until;
ALTER TABLE users DROP COLUMN IF EXISTS failed_login_attempts;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS login_attempts;
DROP TABLE IF EXISTS verification_tokens;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_roles;
