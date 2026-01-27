-- Migration: 000008_roles_permissions.up.sql
-- Description: Create roles and permissions tables for RBAC

-- Create permissions table (global, not tenant-specific)
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    module VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create roles table (tenant-specific)
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT roles_tenant_name_unique UNIQUE(tenant_id, name)
);

-- Create role_permissions junction table
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT role_permissions_unique UNIQUE(role_id, permission_id)
);

-- Create indexes
CREATE INDEX idx_permissions_module ON permissions(module);
CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_is_system ON roles(is_system);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Add comments
COMMENT ON TABLE permissions IS 'System-wide permissions for access control';
COMMENT ON COLUMN permissions.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN permissions.code IS 'Unique permission code (e.g., users:read, students:write)';
COMMENT ON COLUMN permissions.name IS 'Human-readable permission name';
COMMENT ON COLUMN permissions.module IS 'Module this permission belongs to';
COMMENT ON COLUMN permissions.description IS 'Description of what this permission allows';

COMMENT ON TABLE roles IS 'Roles that can be assigned to users within a tenant';
COMMENT ON COLUMN roles.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN roles.tenant_id IS 'Reference to the parent tenant (NULL for system roles)';
COMMENT ON COLUMN roles.name IS 'Role name (unique per tenant)';
COMMENT ON COLUMN roles.description IS 'Description of the role';
COMMENT ON COLUMN roles.is_system IS 'Whether this is a system-defined role';

COMMENT ON TABLE role_permissions IS 'Junction table linking roles to permissions';
COMMENT ON COLUMN role_permissions.role_id IS 'Reference to the role';
COMMENT ON COLUMN role_permissions.permission_id IS 'Reference to the permission';

-- Insert default permissions
-- Users module
INSERT INTO permissions (code, name, module, description) VALUES
    ('users:read', 'View Users', 'users', 'View user accounts and profiles'),
    ('users:write', 'Manage Users', 'users', 'Create and update user accounts'),
    ('users:delete', 'Delete Users', 'users', 'Delete user accounts'),
    ('users:assign_roles', 'Assign Roles', 'users', 'Assign roles to users');

-- Students module
INSERT INTO permissions (code, name, module, description) VALUES
    ('students:read', 'View Students', 'students', 'View student records and profiles'),
    ('students:write', 'Manage Students', 'students', 'Create and update student records'),
    ('students:delete', 'Delete Students', 'students', 'Delete student records'),
    ('students:enroll', 'Enroll Students', 'students', 'Enroll students in classes and programs');

-- Staff module
INSERT INTO permissions (code, name, module, description) VALUES
    ('staff:read', 'View Staff', 'staff', 'View staff records and profiles'),
    ('staff:write', 'Manage Staff', 'staff', 'Create and update staff records'),
    ('staff:delete', 'Delete Staff', 'staff', 'Delete staff records');

-- Academics module
INSERT INTO permissions (code, name, module, description) VALUES
    ('academics:read', 'View Academics', 'academics', 'View academic programs, classes, and schedules'),
    ('academics:write', 'Manage Academics', 'academics', 'Create and update academic programs and classes'),
    ('academics:delete', 'Delete Academics', 'academics', 'Delete academic records'),
    ('academics:grades', 'Manage Grades', 'academics', 'Enter and modify student grades');

-- Finance module
INSERT INTO permissions (code, name, module, description) VALUES
    ('finance:read', 'View Finance', 'finance', 'View financial records and reports'),
    ('finance:write', 'Manage Finance', 'finance', 'Create and update financial records'),
    ('finance:delete', 'Delete Finance', 'finance', 'Delete financial records'),
    ('finance:collect', 'Collect Payments', 'finance', 'Collect fees and process payments');

-- Attendance module
INSERT INTO permissions (code, name, module, description) VALUES
    ('attendance:read', 'View Attendance', 'attendance', 'View attendance records'),
    ('attendance:write', 'Manage Attendance', 'attendance', 'Mark and update attendance'),
    ('attendance:reports', 'Attendance Reports', 'attendance', 'Generate attendance reports');

-- Reports module
INSERT INTO permissions (code, name, module, description) VALUES
    ('reports:read', 'View Reports', 'reports', 'View system reports'),
    ('reports:export', 'Export Reports', 'reports', 'Export reports to various formats');

-- Settings module
INSERT INTO permissions (code, name, module, description) VALUES
    ('settings:read', 'View Settings', 'settings', 'View system settings'),
    ('settings:write', 'Manage Settings', 'settings', 'Update system settings');

-- Audit module
INSERT INTO permissions (code, name, module, description) VALUES
    ('audit:read', 'View Audit Logs', 'audit', 'View system audit logs');

-- Branches module
INSERT INTO permissions (code, name, module, description) VALUES
    ('branches:read', 'View Branches', 'branches', 'View branch information'),
    ('branches:write', 'Manage Branches', 'branches', 'Create and update branches'),
    ('branches:delete', 'Delete Branches', 'branches', 'Delete branches');

-- Roles module
INSERT INTO permissions (code, name, module, description) VALUES
    ('roles:read', 'View Roles', 'roles', 'View roles and permissions'),
    ('roles:write', 'Manage Roles', 'roles', 'Create and update roles'),
    ('roles:delete', 'Delete Roles', 'roles', 'Delete roles');

-- Insert system roles (tenant_id NULL for system-wide roles)
INSERT INTO roles (tenant_id, name, description, is_system) VALUES
    (NULL, 'super_admin', 'System administrator with full access to all features', true),
    (NULL, 'admin', 'Tenant administrator with full access to tenant features', true),
    (NULL, 'principal', 'School principal with access to most features', true),
    (NULL, 'teacher', 'Teacher with access to academic and attendance features', true),
    (NULL, 'accountant', 'Accountant with access to finance features', true),
    (NULL, 'parent', 'Parent with read-only access to child information', true),
    (NULL, 'student', 'Student with limited access to own information', true);

-- Assign all permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin';

-- Assign permissions to admin (all except audit)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.module != 'audit';

-- Assign permissions to principal
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'principal'
AND p.code IN (
    'users:read', 'users:write',
    'students:read', 'students:write', 'students:enroll',
    'staff:read', 'staff:write',
    'academics:read', 'academics:write', 'academics:grades',
    'finance:read',
    'attendance:read', 'attendance:write', 'attendance:reports',
    'reports:read', 'reports:export',
    'settings:read',
    'branches:read',
    'roles:read'
);

-- Assign permissions to teacher
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher'
AND p.code IN (
    'students:read',
    'academics:read', 'academics:grades',
    'attendance:read', 'attendance:write',
    'reports:read'
);

-- Assign permissions to accountant
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'accountant'
AND p.code IN (
    'students:read',
    'finance:read', 'finance:write', 'finance:collect',
    'reports:read', 'reports:export'
);

-- Assign permissions to parent
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'parent'
AND p.code IN (
    'students:read',
    'academics:read',
    'finance:read',
    'attendance:read'
);

-- Assign permissions to student
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'student'
AND p.code IN (
    'academics:read',
    'attendance:read'
);

-- Enable RLS on roles table
ALTER TABLE roles ENABLE ROW LEVEL SECURITY;

-- Create RLS policy for roles table (system roles visible to all, tenant roles only to tenant)
CREATE POLICY tenant_isolation_roles ON roles
    FOR ALL
    USING (
        tenant_id IS NULL -- System roles visible to all
        OR tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_roles ON roles
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

COMMENT ON POLICY tenant_isolation_roles ON roles IS 'Restricts role access - system roles visible to all, tenant roles to tenant only';
COMMENT ON POLICY bypass_rls_roles ON roles IS 'Allows bypass of RLS for admin operations';
