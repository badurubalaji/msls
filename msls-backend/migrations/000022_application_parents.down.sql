-- Migration: 000022_application_parents.down.sql
-- Description: Remove application_parents table

DROP TRIGGER IF EXISTS trigger_app_parents_updated_at ON application_parents;
DROP POLICY IF EXISTS tenant_isolation_app_parents ON application_parents;
DROP POLICY IF EXISTS bypass_rls_app_parents ON application_parents;
DROP TABLE IF EXISTS application_parents;
DROP TYPE IF EXISTS parent_relation;
