-- Migration: 000008_roles_permissions.down.sql
-- Description: Drop roles and permissions tables

-- Drop RLS policies
DROP POLICY IF EXISTS bypass_rls_roles ON roles;
DROP POLICY IF EXISTS tenant_isolation_roles ON roles;

-- Disable RLS
ALTER TABLE roles DISABLE ROW LEVEL SECURITY;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
