-- Migration: 000006_rls_policies.down.sql
-- Description: Remove Row-Level Security policies

-- Drop policies
DROP POLICY IF EXISTS tenant_isolation_branches ON branches;
DROP POLICY IF EXISTS tenant_isolation_users ON users;
DROP POLICY IF EXISTS bypass_rls_branches ON branches;
DROP POLICY IF EXISTS bypass_rls_users ON users;

-- Disable RLS on tables
ALTER TABLE branches DISABLE ROW LEVEL SECURITY;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
