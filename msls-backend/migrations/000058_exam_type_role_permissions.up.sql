-- Fix: Assign exam type permissions to roles
-- This migration adds role_permissions that were missing from 000057

-- Assign exam type permissions to admin roles (super_admin, admin, principal)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('super_admin', 'admin', 'principal')
AND p.code IN ('exam:type:view', 'exam:type:create', 'exam:type:update', 'exam:type:delete')
ON CONFLICT DO NOTHING;

-- Teachers can view exam types
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher'
AND p.code = 'exam:type:view'
ON CONFLICT DO NOTHING;
