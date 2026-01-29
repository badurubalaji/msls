-- Migration: 000041_academic_structure.down.sql
-- Description: Drop academic structure tables

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_class_subjects_updated_at ON class_subjects;
DROP TRIGGER IF EXISTS trigger_subjects_updated_at ON subjects;
DROP TRIGGER IF EXISTS trigger_sections_updated_at ON sections;
DROP TRIGGER IF EXISTS trigger_classes_updated_at ON classes;

DROP FUNCTION IF EXISTS update_classes_updated_at();

-- Drop permissions
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'class:view', 'class:create', 'class:update', 'class:delete',
        'section:view', 'section:create', 'section:update', 'section:delete',
        'subject:view', 'subject:create', 'subject:update', 'subject:delete'
    )
);

DELETE FROM permissions WHERE code IN (
    'class:view', 'class:create', 'class:update', 'class:delete',
    'section:view', 'section:create', 'section:update', 'section:delete',
    'subject:view', 'subject:create', 'subject:update', 'subject:delete'
);

-- Drop tables
DROP TABLE IF EXISTS class_subjects;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS sections;
DROP TABLE IF EXISTS classes;
