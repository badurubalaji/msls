-- Migration: 000033_fix_epic4_role_permissions.up.sql
-- Description: Add missing health permissions and assign all Epic 4 permissions to admin roles

-- First, add any missing permissions (health permissions were missing)
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('health:read', 'View Health Records', 'Permission to view student health information', 'students', NOW(), NOW()),
    ('health:write', 'Manage Health Records', 'Permission to manage student health information', 'students', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign all Epic 4 permissions to super_admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
  AND p.code IN (
    -- Guardian permissions
    'guardians:read', 'guardians:write',
    -- Emergency contact permissions
    'emergency_contacts:read', 'emergency_contacts:write',
    -- Health permissions
    'health:read', 'health:write',
    -- Behavioral permissions
    'behavior:read', 'behavior:write',
    -- Document permissions
    'document:read', 'document:create', 'document:update', 'document:delete', 'document:verify', 'document_type:manage',
    -- Enrollment permissions
    'enrollments:read', 'enrollments:create', 'enrollments:update', 'enrollments:delete',
    -- Promotion permissions
    'promotion:read', 'promotion:create', 'promotion:update', 'promotion:process', 'promotion:cancel',
    'promotion_rules:read', 'promotion_rules:manage',
    -- Student permissions
    'students:create', 'students:update', 'students:export'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign all Epic 4 permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
  AND p.code IN (
    -- Guardian permissions
    'guardians:read', 'guardians:write',
    -- Emergency contact permissions
    'emergency_contacts:read', 'emergency_contacts:write',
    -- Health permissions
    'health:read', 'health:write',
    -- Behavioral permissions
    'behavior:read', 'behavior:write',
    -- Document permissions
    'document:read', 'document:create', 'document:update', 'document:delete', 'document:verify', 'document_type:manage',
    -- Enrollment permissions
    'enrollments:read', 'enrollments:create', 'enrollments:update', 'enrollments:delete',
    -- Promotion permissions
    'promotion:read', 'promotion:create', 'promotion:update', 'promotion:process', 'promotion:cancel',
    'promotion_rules:read', 'promotion_rules:manage',
    -- Student permissions
    'students:create', 'students:update', 'students:export'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;
