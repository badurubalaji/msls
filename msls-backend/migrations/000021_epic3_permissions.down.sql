-- Migration: 000021_epic3_permissions.down.sql
-- Description: Remove Epic 3 permissions

-- Remove role_permissions first (foreign key constraints)
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'admissions:read', 'admissions:create', 'admissions:update', 'admissions:delete',
        'applications:read', 'applications:create', 'applications:update', 'applications:delete',
        'academic-years:read', 'academic-years:create', 'academic-years:update', 'academic-years:delete'
    )
);

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'admissions:read', 'admissions:create', 'admissions:update', 'admissions:delete',
    'applications:read', 'applications:create', 'applications:update', 'applications:delete',
    'academic-years:read', 'academic-years:create', 'academic-years:update', 'academic-years:delete'
);
