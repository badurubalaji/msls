-- Migration: 000007_indexes.down.sql
-- Description: Remove performance indexes

-- Drop branches indexes
DROP INDEX IF EXISTS idx_branches_tenant_id;
DROP INDEX IF EXISTS idx_branches_status;
DROP INDEX IF EXISTS idx_branches_is_primary;

-- Drop users indexes
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_tenant_email;
DROP INDEX IF EXISTS idx_users_tenant_phone;
