-- Migration: 000035_staff.down.sql
-- Description: Drop staff table and related tables

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_staff_employee_sequences_updated_at ON staff_employee_sequences;
DROP TRIGGER IF EXISTS trigger_staff_updated_at ON staff;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_staff_updated_at();

-- Drop FK from departments
ALTER TABLE departments DROP CONSTRAINT IF EXISTS fk_department_head;

-- Drop tables
DROP TABLE IF EXISTS staff_status_history;
DROP TABLE IF EXISTS staff_employee_sequences;
DROP TABLE IF EXISTS staff;

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'staff:read', 'staff:create', 'staff:update', 'staff:delete', 'staff:export'
);
