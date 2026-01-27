-- Migration: 000016_merit_list_admission_decision.down.sql
-- Description: Drop merit_lists and admission_decisions tables

-- Drop trigger
DROP TRIGGER IF EXISTS trigger_admission_decisions_updated_at ON admission_decisions;

-- Drop function
DROP FUNCTION IF EXISTS update_admission_decision_updated_at();

-- Drop RLS policies for merit_lists
DROP POLICY IF EXISTS tenant_isolation_merit_lists ON merit_lists;
DROP POLICY IF EXISTS bypass_rls_merit_lists ON merit_lists;

-- Drop RLS policies for admission_decisions
DROP POLICY IF EXISTS tenant_isolation_admission_decisions ON admission_decisions;
DROP POLICY IF EXISTS bypass_rls_admission_decisions ON admission_decisions;

-- Drop indexes for merit_lists
DROP INDEX IF EXISTS idx_merit_lists_tenant;
DROP INDEX IF EXISTS idx_merit_lists_session;
DROP INDEX IF EXISTS idx_merit_lists_class;
DROP INDEX IF EXISTS idx_merit_lists_test;
DROP INDEX IF EXISTS idx_merit_lists_generated_at;
DROP INDEX IF EXISTS idx_merit_lists_is_final;

-- Drop indexes for admission_decisions
DROP INDEX IF EXISTS idx_admission_decisions_tenant;
DROP INDEX IF EXISTS idx_admission_decisions_application;
DROP INDEX IF EXISTS idx_admission_decisions_decision;
DROP INDEX IF EXISTS idx_admission_decisions_decided_by;
DROP INDEX IF EXISTS idx_admission_decisions_decision_date;
DROP INDEX IF EXISTS idx_admission_decisions_waitlist;

-- Drop tables
DROP TABLE IF EXISTS merit_lists;
DROP TABLE IF EXISTS admission_decisions;

-- Drop enum
DROP TYPE IF EXISTS decision_type;
