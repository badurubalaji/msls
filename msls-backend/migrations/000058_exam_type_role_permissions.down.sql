-- Rollback: Remove exam type role permissions

DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code LIKE 'exam:type:%'
);
