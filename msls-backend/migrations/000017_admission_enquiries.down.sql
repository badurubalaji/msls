-- Migration: 000015_admission_enquiries.down.sql
-- Description: Drop admission enquiries and follow-ups tables

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_enquiries_updated_at ON admission_enquiries;
DROP FUNCTION IF EXISTS update_enquiry_updated_at();
DROP FUNCTION IF EXISTS get_next_enquiry_number(UUID);

-- Drop RLS policies
DROP POLICY IF EXISTS tenant_isolation_seq_update ON enquiry_number_sequence;
DROP POLICY IF EXISTS tenant_isolation_seq_insert ON enquiry_number_sequence;
DROP POLICY IF EXISTS tenant_isolation_seq ON enquiry_number_sequence;

DROP POLICY IF EXISTS tenant_isolation_follow_ups_delete ON enquiry_follow_ups;
DROP POLICY IF EXISTS tenant_isolation_follow_ups_update ON enquiry_follow_ups;
DROP POLICY IF EXISTS tenant_isolation_follow_ups_insert ON enquiry_follow_ups;
DROP POLICY IF EXISTS tenant_isolation_follow_ups ON enquiry_follow_ups;

DROP POLICY IF EXISTS tenant_isolation_enquiries_delete ON admission_enquiries;
DROP POLICY IF EXISTS tenant_isolation_enquiries_update ON admission_enquiries;
DROP POLICY IF EXISTS tenant_isolation_enquiries_insert ON admission_enquiries;
DROP POLICY IF EXISTS tenant_isolation_enquiries ON admission_enquiries;

-- Drop indexes
DROP INDEX IF EXISTS idx_enquiry_seq_tenant;
DROP INDEX IF EXISTS idx_follow_ups_created_at;
DROP INDEX IF EXISTS idx_follow_ups_date;
DROP INDEX IF EXISTS idx_follow_ups_enquiry;
DROP INDEX IF EXISTS idx_follow_ups_tenant;
DROP INDEX IF EXISTS idx_enquiries_deleted_at;
DROP INDEX IF EXISTS idx_enquiries_phone;
DROP INDEX IF EXISTS idx_enquiries_assigned_to;
DROP INDEX IF EXISTS idx_enquiries_follow_up_date;
DROP INDEX IF EXISTS idx_enquiries_created_at;
DROP INDEX IF EXISTS idx_enquiries_class;
DROP INDEX IF EXISTS idx_enquiries_status;
DROP INDEX IF EXISTS idx_enquiries_branch;
DROP INDEX IF EXISTS idx_enquiries_tenant;

-- Drop tables
DROP TABLE IF EXISTS enquiry_number_sequence;
DROP TABLE IF EXISTS enquiry_follow_ups;
DROP TABLE IF EXISTS admission_enquiries;

-- Drop enums
DROP TYPE IF EXISTS follow_up_outcome;
DROP TYPE IF EXISTS contact_mode;
DROP TYPE IF EXISTS enquiry_status;
DROP TYPE IF EXISTS enquiry_source;
