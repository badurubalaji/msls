-- Rollback: Examinations

-- Remove role permissions
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code LIKE 'exam:view' OR code LIKE 'exam:create' OR code LIKE 'exam:update' OR code LIKE 'exam:delete' OR code LIKE 'exam:publish'
);

-- Remove permissions
DELETE FROM permissions WHERE code IN ('exam:view', 'exam:create', 'exam:update', 'exam:delete', 'exam:publish');

-- Drop tables (in correct order due to foreign keys)
DROP TABLE IF EXISTS exam_schedules;
DROP TABLE IF EXISTS examination_classes;
DROP TABLE IF EXISTS examinations;
