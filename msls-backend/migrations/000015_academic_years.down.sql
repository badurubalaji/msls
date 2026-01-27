-- Migration: 000015_academic_years.down.sql
-- Description: Drop academic years, terms, and holidays tables

-- Remove permissions from roles
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE module = 'academic-years'
);

-- Remove permissions
DELETE FROM permissions WHERE module = 'academic-years';

-- Drop RLS policies
DROP POLICY IF EXISTS tenant_isolation_holidays ON holidays;
DROP POLICY IF EXISTS bypass_rls_holidays ON holidays;
DROP POLICY IF EXISTS tenant_isolation_academic_terms ON academic_terms;
DROP POLICY IF EXISTS bypass_rls_academic_terms ON academic_terms;
DROP POLICY IF EXISTS tenant_isolation_academic_years ON academic_years;
DROP POLICY IF EXISTS bypass_rls_academic_years ON academic_years;

-- Drop triggers and functions
DROP TRIGGER IF EXISTS trigger_single_current_academic_year ON academic_years;
DROP FUNCTION IF EXISTS ensure_single_current_academic_year();

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS holidays;
DROP TABLE IF EXISTS academic_terms;
DROP TABLE IF EXISTS academic_years;
