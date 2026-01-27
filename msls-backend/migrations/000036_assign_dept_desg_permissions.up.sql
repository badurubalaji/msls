-- Migration: 000036_assign_dept_desg_permissions.up.sql
-- Description: Assign department and designation permissions to super_admin role

-- Assign department/designation permissions to super_admin role
INSERT INTO role_permissions (role_id, permission_id, created_at, updated_at)
SELECT r.id, p.id, NOW(), NOW()
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN (
    'department:read',
    'department:create',
    'department:update',
    'department:delete',
    'designation:read',
    'designation:create',
    'designation:update',
    'designation:delete'
)
ON CONFLICT DO NOTHING;
