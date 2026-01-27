-- Migration: 000025_student_guardians.down.sql
-- Description: Drop student_guardians and student_emergency_contacts tables

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_student_emergency_contacts_updated_at ON student_emergency_contacts;
DROP TRIGGER IF EXISTS trigger_student_guardians_updated_at ON student_guardians;

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'guardians:read',
    'guardians:write',
    'emergency_contacts:read',
    'emergency_contacts:write'
);

-- Drop tables
DROP TABLE IF EXISTS student_emergency_contacts;
DROP TABLE IF EXISTS student_guardians;
