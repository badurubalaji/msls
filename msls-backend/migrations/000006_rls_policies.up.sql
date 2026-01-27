-- Migration: 000006_rls_policies.up.sql
-- Description: Enable Row-Level Security policies for multi-tenant data isolation

-- Enable RLS on branches table
ALTER TABLE branches ENABLE ROW LEVEL SECURITY;

-- Enable RLS on users table
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Create RLS policy for branches table
-- Uses app.tenant_id session variable set by the application
CREATE POLICY tenant_isolation_branches ON branches
    FOR ALL
    USING (tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID));

-- Create RLS policy for users table
CREATE POLICY tenant_isolation_users ON users
    FOR ALL
    USING (tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID));

-- Create bypass policy for superusers/admin operations
-- This allows the application to perform cross-tenant operations when needed
CREATE POLICY bypass_rls_branches ON branches
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

CREATE POLICY bypass_rls_users ON users
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- Add comments
COMMENT ON POLICY tenant_isolation_branches ON branches IS 'Restricts branch access to current tenant only';
COMMENT ON POLICY tenant_isolation_users ON users IS 'Restricts user access to current tenant only';
COMMENT ON POLICY bypass_rls_branches ON branches IS 'Allows bypass of RLS for admin operations';
COMMENT ON POLICY bypass_rls_users ON users IS 'Allows bypass of RLS for admin operations';
