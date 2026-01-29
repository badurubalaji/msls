-- Migration: 000040_payroll_seed_data.down.sql
-- Rollback: Remove payroll seed data

-- Remove payslip components for seed payslips
DELETE FROM payslip_components WHERE payslip_id IN (
    SELECT p.id FROM payslips p
    JOIN pay_runs pr ON p.pay_run_id = pr.id
    WHERE (pr.pay_period_month = 12 AND pr.pay_period_year = 2025)
       OR (pr.pay_period_month = 1 AND pr.pay_period_year = 2026)
);

-- Remove payslips for seed pay runs
DELETE FROM payslips WHERE pay_run_id IN (
    SELECT id FROM pay_runs
    WHERE (pay_period_month = 12 AND pay_period_year = 2025)
       OR (pay_period_month = 1 AND pay_period_year = 2026)
);

-- Remove seed pay runs
DELETE FROM pay_runs
WHERE (pay_period_month = 12 AND pay_period_year = 2025)
   OR (pay_period_month = 1 AND pay_period_year = 2026);

-- Note: Staff, salary components, structures, and staff salaries are NOT removed
-- as they may be used by other data. Remove manually if needed.
