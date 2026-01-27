-- Migration: 000034_departments_designations.down.sql
-- Description: Drop departments and designations tables

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_designations_updated_at ON designations;
DROP TRIGGER IF EXISTS trigger_departments_updated_at ON departments;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_designations_updated_at();
DROP FUNCTION IF EXISTS update_departments_updated_at();

-- Drop tables
DROP TABLE IF EXISTS designations;
DROP TABLE IF EXISTS departments;

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'department:read', 'department:create', 'department:update', 'department:delete',
    'designation:read', 'designation:create', 'designation:update', 'designation:delete'
);
