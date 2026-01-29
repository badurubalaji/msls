-- Migration: 000042_teacher_assignments.down.sql
-- Description: Drop teacher subject assignment tables

-- Remove permissions from role_permissions
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'assignment:view',
        'assignment:create',
        'assignment:update',
        'assignment:delete',
        'assignment:workload',
        'assignment:class_teacher'
    )
);

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'assignment:view',
    'assignment:create',
    'assignment:update',
    'assignment:delete',
    'assignment:workload',
    'assignment:class_teacher'
);

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_workload_settings_updated_at ON teacher_workload_settings;
DROP TRIGGER IF EXISTS trigger_teacher_assignments_updated_at ON teacher_subject_assignments;
DROP FUNCTION IF EXISTS update_teacher_assignments_updated_at();

-- Drop tables
DROP TABLE IF EXISTS teacher_workload_settings;
DROP TABLE IF EXISTS teacher_subject_assignments;
