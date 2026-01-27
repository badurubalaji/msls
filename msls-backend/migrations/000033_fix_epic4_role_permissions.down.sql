-- Migration: 000033_fix_epic4_role_permissions.down.sql
-- Description: Remove Epic 4 permission assignments from admin roles
-- Note: This removes role_permissions entries and the health permissions we added

-- Remove the role_permissions assignments
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'guardians:read', 'guardians:write',
        'emergency_contacts:read', 'emergency_contacts:write',
        'health:read', 'health:write',
        'behavior:read', 'behavior:write',
        'document:read', 'document:create', 'document:update', 'document:delete', 'document:verify', 'document_type:manage',
        'enrollments:read', 'enrollments:create', 'enrollments:update', 'enrollments:delete',
        'promotion:read', 'promotion:create', 'promotion:update', 'promotion:process', 'promotion:cancel',
        'promotion_rules:read', 'promotion_rules:manage',
        'students:create', 'students:update', 'students:export'
    )
)
AND role_id IN (
    SELECT id FROM roles WHERE name IN ('super_admin', 'admin')
);

-- Remove the health permissions we added
DELETE FROM permissions WHERE code IN ('health:read', 'health:write');
