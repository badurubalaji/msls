-- Migration: 000039_payroll.down.sql
-- Rollback Story 5.6: Payroll Processing

-- Remove role permissions
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code LIKE 'payroll.%'
);

-- Remove permissions
DELETE FROM permissions WHERE code LIKE 'payroll.%';

-- Drop indexes
DROP INDEX IF EXISTS idx_payslip_components_component_id;
DROP INDEX IF EXISTS idx_payslip_components_payslip_id;
DROP INDEX IF EXISTS idx_payslips_status;
DROP INDEX IF EXISTS idx_payslips_staff_id;
DROP INDEX IF EXISTS idx_payslips_pay_run_id;
DROP INDEX IF EXISTS idx_payslips_tenant_id;
DROP INDEX IF EXISTS idx_pay_runs_branch_id;
DROP INDEX IF EXISTS idx_pay_runs_status;
DROP INDEX IF EXISTS idx_pay_runs_period;
DROP INDEX IF EXISTS idx_pay_runs_tenant_id;
DROP INDEX IF EXISTS idx_pay_runs_unique_period_no_branch;
DROP INDEX IF EXISTS idx_pay_runs_unique_period_with_branch;

-- Drop RLS policies
DROP POLICY IF EXISTS tenant_isolation_payslip_components ON payslip_components;
DROP POLICY IF EXISTS tenant_isolation_payslips ON payslips;
DROP POLICY IF EXISTS tenant_isolation_pay_runs ON pay_runs;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS payslip_components;
DROP TABLE IF EXISTS payslips;
DROP TABLE IF EXISTS pay_runs;
