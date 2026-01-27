-- Migration: 000010_feature_flags.down.sql
-- Description: Drop feature flags tables

-- Drop RLS policy
DROP POLICY IF EXISTS tenant_feature_flags_tenant_isolation ON tenant_feature_flags;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_feature_flags_flag_id;
DROP INDEX IF EXISTS idx_user_feature_flags_user_id;
DROP INDEX IF EXISTS idx_tenant_feature_flags_flag_id;
DROP INDEX IF EXISTS idx_tenant_feature_flags_tenant_id;
DROP INDEX IF EXISTS idx_feature_flags_key;

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS user_feature_flags;
DROP TABLE IF EXISTS tenant_feature_flags;
DROP TABLE IF EXISTS feature_flags;
