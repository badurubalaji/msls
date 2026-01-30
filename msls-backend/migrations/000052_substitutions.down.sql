-- Migration: 000052_substitutions.down.sql
-- Description: Drop substitutions tables

-- Remove permissions from roles
DELETE FROM role_permissions
WHERE permission_id IN (SELECT id FROM permissions WHERE code LIKE 'substitution:%');

-- Remove permissions
DELETE FROM permissions WHERE code LIKE 'substitution:%';

-- Drop tables
DROP TABLE IF EXISTS substitution_periods;
DROP TABLE IF EXISTS substitutions;

-- Drop enum
DROP TYPE IF EXISTS substitution_status;
