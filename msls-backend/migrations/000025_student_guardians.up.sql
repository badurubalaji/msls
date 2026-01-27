-- Migration: 000025_student_guardians.up.sql
-- Description: Create student_guardians and student_emergency_contacts tables for Epic 4 Story 4.2

-- Guardian relation types
-- father, mother, grandfather, grandmother, uncle, aunt, sibling, other

-- Student guardians table
CREATE TABLE student_guardians (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    relation VARCHAR(20) NOT NULL CHECK (relation IN ('father', 'mother', 'grandfather', 'grandmother', 'uncle', 'aunt', 'sibling', 'guardian', 'other')),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    email VARCHAR(255),
    occupation VARCHAR(100),
    annual_income DECIMAL(15, 2),
    education VARCHAR(100),
    is_primary BOOLEAN NOT NULL DEFAULT false,
    has_portal_access BOOLEAN NOT NULL DEFAULT false,
    user_id UUID REFERENCES users(id),
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(10),
    country VARCHAR(100) DEFAULT 'India',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Enable Row Level Security
ALTER TABLE student_guardians ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_student_guardians ON student_guardians
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for student_guardians
CREATE INDEX idx_student_guardians_tenant ON student_guardians(tenant_id);
CREATE INDEX idx_student_guardians_student ON student_guardians(student_id);
CREATE INDEX idx_student_guardians_phone ON student_guardians(tenant_id, phone);
CREATE INDEX idx_student_guardians_email ON student_guardians(tenant_id, email) WHERE email IS NOT NULL;
CREATE INDEX idx_student_guardians_user ON student_guardians(user_id) WHERE user_id IS NOT NULL;

-- Ensure only one primary guardian per student
CREATE UNIQUE INDEX uniq_student_primary_guardian ON student_guardians(student_id) WHERE is_primary = true;

-- Student emergency contacts table
CREATE TABLE student_emergency_contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    relation VARCHAR(50) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    alternate_phone VARCHAR(15),
    priority INTEGER NOT NULL DEFAULT 1 CHECK (priority >= 1 AND priority <= 5),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Enable Row Level Security
ALTER TABLE student_emergency_contacts ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_student_emergency_contacts ON student_emergency_contacts
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for student_emergency_contacts
CREATE INDEX idx_student_emergency_contacts_tenant ON student_emergency_contacts(tenant_id);
CREATE INDEX idx_student_emergency_contacts_student ON student_emergency_contacts(student_id);
CREATE INDEX idx_student_emergency_contacts_priority ON student_emergency_contacts(student_id, priority);

-- Ensure unique priority per student
CREATE UNIQUE INDEX uniq_student_emergency_priority ON student_emergency_contacts(student_id, priority);

-- Add guardian permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('guardians:read', 'View Guardians', 'Permission to view student guardian information', 'students', NOW(), NOW()),
    ('guardians:write', 'Manage Guardians', 'Permission to create, update, and delete student guardians', 'students', NOW(), NOW()),
    ('emergency_contacts:read', 'View Emergency Contacts', 'Permission to view student emergency contacts', 'students', NOW(), NOW()),
    ('emergency_contacts:write', 'Manage Emergency Contacts', 'Permission to create, update, and delete emergency contacts', 'students', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Add triggers to update updated_at timestamp
CREATE TRIGGER trigger_student_guardians_updated_at
    BEFORE UPDATE ON student_guardians
    FOR EACH ROW
    EXECUTE FUNCTION update_students_updated_at();

CREATE TRIGGER trigger_student_emergency_contacts_updated_at
    BEFORE UPDATE ON student_emergency_contacts
    FOR EACH ROW
    EXECUTE FUNCTION update_students_updated_at();
