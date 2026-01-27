-- Migration: 000036_assign_dept_desg_permissions.down.sql
-- Description: Remove department and designation permissions from super_admin role

DELETE FROM role_permissions
WHERE role_id = (SELECT id FROM roles WHERE name = 'super_admin')
AND permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'department:read',
        'department:create',
        'department:update',
        'department:delete',
        'designation:read',
        'designation:create',
        'designation:update',
        'designation:delete'
    )
);
