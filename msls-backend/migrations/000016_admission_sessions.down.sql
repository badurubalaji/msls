-- Migration: 000015_admission_sessions.down.sql
-- Description: Drop admission sessions and seats tables

-- Remove admission permissions from roles
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE module = 'admissions'
);

-- Remove admission permissions
DELETE FROM permissions WHERE module = 'admissions';

-- Drop RLS policies for admission_seats
DROP POLICY IF EXISTS tenant_isolation_admission_seats ON admission_seats;
DROP POLICY IF EXISTS bypass_rls_admission_seats ON admission_seats;

-- Drop RLS policies for admission_sessions
DROP POLICY IF EXISTS tenant_isolation_admission_sessions ON admission_sessions;
DROP POLICY IF EXISTS bypass_rls_admission_sessions ON admission_sessions;

-- Drop indexes for admission_seats
DROP INDEX IF EXISTS idx_admission_seats_class_name;
DROP INDEX IF EXISTS idx_admission_seats_session_id;
DROP INDEX IF EXISTS idx_admission_seats_tenant_id;

-- Drop indexes for admission_sessions
DROP INDEX IF EXISTS idx_admission_sessions_end_date;
DROP INDEX IF EXISTS idx_admission_sessions_start_date;
DROP INDEX IF EXISTS idx_admission_sessions_status;
DROP INDEX IF EXISTS idx_admission_sessions_academic_year_id;
DROP INDEX IF EXISTS idx_admission_sessions_branch_id;
DROP INDEX IF EXISTS idx_admission_sessions_tenant_id;

-- Drop admission_seats table
DROP TABLE IF EXISTS admission_seats;

-- Drop admission_sessions table
DROP TABLE IF EXISTS admission_sessions;
