-- Migration: 000017_fix_branch_permissions.down.sql
-- Description: Remove permissions added in this migration

-- Remove role_permissions entries
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'branches:create', 'branches:update',
        'applications:review',
        'enquiries:read', 'enquiries:create', 'enquiries:update', 'enquiries:delete',
        'tests:read', 'tests:create', 'tests:update', 'tests:manage', 'tests:delete'
    )
);

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'branches:create', 'branches:update',
    'applications:review',
    'enquiries:read', 'enquiries:create', 'enquiries:update', 'enquiries:delete',
    'tests:read', 'tests:create', 'tests:update', 'tests:manage', 'tests:delete'
);
