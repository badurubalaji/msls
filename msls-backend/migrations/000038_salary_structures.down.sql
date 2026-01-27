-- Migration: 000038_salary_structures.down.sql
-- Description: Drop salary components, structures, and staff salary tables

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_staff_salaries_updated_at ON staff_salaries;
DROP TRIGGER IF EXISTS trigger_salary_structures_updated_at ON salary_structures;
DROP TRIGGER IF EXISTS trigger_salary_components_updated_at ON salary_components;

-- Drop functions
DROP FUNCTION IF EXISTS update_staff_salaries_updated_at();
DROP FUNCTION IF EXISTS update_salary_structures_updated_at();
DROP FUNCTION IF EXISTS update_salary_components_updated_at();

-- Remove permissions from roles
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'salary:read', 'salary:create', 'salary:update', 'salary:delete', 'salary:assign'
    )
);

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'salary:read', 'salary:create', 'salary:update', 'salary:delete', 'salary:assign'
);

-- Drop tables in order (respecting foreign keys)
DROP TABLE IF EXISTS staff_salary_components;
DROP TABLE IF EXISTS staff_salaries;
DROP TABLE IF EXISTS salary_structure_components;
DROP TABLE IF EXISTS salary_structures;
DROP TABLE IF EXISTS salary_components;
