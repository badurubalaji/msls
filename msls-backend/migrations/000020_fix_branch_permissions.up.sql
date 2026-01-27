-- Migration: 000017_fix_branch_permissions.up.sql
-- Description: Add missing permissions required by API routes

-- Add missing branch permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('branches:create', 'Create Branches', 'branches', 'Create new branches'),
    ('branches:update', 'Update Branches', 'branches', 'Update branch information')
ON CONFLICT (code) DO NOTHING;

-- Add missing application review permission
INSERT INTO permissions (code, name, module, description) VALUES
    ('applications:review', 'Review Applications', 'admissions', 'Review and process admission applications')
ON CONFLICT (code) DO NOTHING;

-- Add enquiry permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('enquiries:read', 'View Enquiries', 'admissions', 'View admission enquiries'),
    ('enquiries:create', 'Create Enquiries', 'admissions', 'Create admission enquiries'),
    ('enquiries:update', 'Update Enquiries', 'admissions', 'Update admission enquiries'),
    ('enquiries:delete', 'Delete Enquiries', 'admissions', 'Delete admission enquiries')
ON CONFLICT (code) DO NOTHING;

-- Add test permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('tests:read', 'View Tests', 'admissions', 'View entrance tests'),
    ('tests:create', 'Create Tests', 'admissions', 'Create entrance tests'),
    ('tests:update', 'Update Tests', 'admissions', 'Update entrance tests'),
    ('tests:manage', 'Manage Tests', 'admissions', 'Manage test registrations and results'),
    ('tests:delete', 'Delete Tests', 'admissions', 'Delete entrance tests')
ON CONFLICT (code) DO NOTHING;

-- Assign all new permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN (
    'branches:create', 'branches:update',
    'applications:review',
    'enquiries:read', 'enquiries:create', 'enquiries:update', 'enquiries:delete',
    'tests:read', 'tests:create', 'tests:update', 'tests:manage', 'tests:delete'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign all new permissions to admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.code IN (
    'branches:create', 'branches:update',
    'applications:review',
    'enquiries:read', 'enquiries:create', 'enquiries:update', 'enquiries:delete',
    'tests:read', 'tests:create', 'tests:update', 'tests:manage', 'tests:delete'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to principal
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'principal'
AND p.code IN (
    'branches:update',
    'applications:review',
    'enquiries:read', 'enquiries:create', 'enquiries:update',
    'tests:read', 'tests:create', 'tests:update', 'tests:manage'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;
