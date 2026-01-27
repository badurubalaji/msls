-- Migration: 000021_epic3_permissions.up.sql
-- Description: Add all missing permissions for Epic 3 features (School Setup & Admissions)

-- Add admission session permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('admissions:read', 'View Admissions', 'admissions', 'View admission sessions and related data'),
    ('admissions:create', 'Create Admissions', 'admissions', 'Create admission sessions'),
    ('admissions:update', 'Update Admissions', 'admissions', 'Update admission sessions and make decisions'),
    ('admissions:delete', 'Delete Admissions', 'admissions', 'Delete admission sessions')
ON CONFLICT (code) DO NOTHING;

-- Add application permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('applications:read', 'View Applications', 'admissions', 'View admission applications'),
    ('applications:create', 'Create Applications', 'admissions', 'Create admission applications'),
    ('applications:update', 'Update Applications', 'admissions', 'Update admission applications'),
    ('applications:delete', 'Delete Applications', 'admissions', 'Delete admission applications')
ON CONFLICT (code) DO NOTHING;

-- Add academic year permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('academic-years:read', 'View Academic Years', 'academic-years', 'View academic years, terms, and holidays'),
    ('academic-years:create', 'Create Academic Years', 'academic-years', 'Create academic years'),
    ('academic-years:update', 'Update Academic Years', 'academic-years', 'Update academic years, terms, and holidays'),
    ('academic-years:delete', 'Delete Academic Years', 'academic-years', 'Delete academic years')
ON CONFLICT (code) DO NOTHING;

-- Assign all new permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN (
    'admissions:read', 'admissions:create', 'admissions:update', 'admissions:delete',
    'applications:read', 'applications:create', 'applications:update', 'applications:delete',
    'academic-years:read', 'academic-years:create', 'academic-years:update', 'academic-years:delete'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign all new permissions to admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.code IN (
    'admissions:read', 'admissions:create', 'admissions:update', 'admissions:delete',
    'applications:read', 'applications:create', 'applications:update', 'applications:delete',
    'academic-years:read', 'academic-years:create', 'academic-years:update', 'academic-years:delete'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign relevant permissions to principal
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'principal'
AND p.code IN (
    'admissions:read', 'admissions:create', 'admissions:update',
    'applications:read', 'applications:create', 'applications:update',
    'academic-years:read', 'academic-years:update'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to teacher (to view academic year info)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher'
AND p.code IN (
    'academic-years:read'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to accountant (to view academic year and admission fee info)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'accountant'
AND p.code IN (
    'academic-years:read',
    'admissions:read',
    'applications:read'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;
