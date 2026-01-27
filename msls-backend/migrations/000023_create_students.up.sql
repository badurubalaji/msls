-- Migration: 000023_create_students.up.sql
-- Description: Create students and student_addresses tables for Epic 4 Story 4.1

-- Students table
CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),
    admission_number VARCHAR(20) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    blood_group VARCHAR(5),
    aadhaar_number VARCHAR(12),
    photo_url VARCHAR(500),
    birth_certificate_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'transferred', 'graduated')),
    admission_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    version INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT uniq_students_tenant_admission UNIQUE (tenant_id, admission_number)
);

-- Enable Row Level Security
ALTER TABLE students ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_students ON students
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for students
CREATE INDEX idx_students_tenant ON students(tenant_id);
CREATE INDEX idx_students_branch ON students(branch_id);
CREATE INDEX idx_students_status ON students(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_name ON students(tenant_id, last_name, first_name);
CREATE INDEX idx_students_admission_number ON students(tenant_id, admission_number);
CREATE INDEX idx_students_deleted_at ON students(deleted_at) WHERE deleted_at IS NOT NULL;

-- Admission number sequence table for tracking sequences per branch/year
CREATE TABLE student_admission_sequences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),
    year INTEGER NOT NULL,
    last_sequence INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_admission_sequence UNIQUE (tenant_id, branch_id, year)
);

ALTER TABLE student_admission_sequences ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_admission_sequences ON student_admission_sequences
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE INDEX idx_admission_sequences_tenant ON student_admission_sequences(tenant_id);

-- Student addresses table
CREATE TABLE student_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    address_type VARCHAR(20) NOT NULL CHECK (address_type IN ('current', 'permanent')),
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'India',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_student_address_type UNIQUE (student_id, address_type)
);

-- Enable Row Level Security
ALTER TABLE student_addresses ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY tenant_isolation_student_addresses ON student_addresses
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes for student_addresses
CREATE INDEX idx_student_addresses_tenant ON student_addresses(tenant_id);
CREATE INDEX idx_student_addresses_student ON student_addresses(student_id);

-- Add student permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('students:read', 'View Students', 'Permission to view student profiles', 'students', NOW(), NOW()),
    ('students:create', 'Create Students', 'Permission to create student profiles', 'students', NOW(), NOW()),
    ('students:update', 'Update Students', 'Permission to update student profiles', 'students', NOW(), NOW()),
    ('students:delete', 'Delete Students', 'Permission to delete student profiles', 'students', NOW(), NOW()),
    ('students:export', 'Export Students', 'Permission to export student data', 'students', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_students_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_students_updated_at
    BEFORE UPDATE ON students
    FOR EACH ROW
    EXECUTE FUNCTION update_students_updated_at();

CREATE TRIGGER trigger_student_addresses_updated_at
    BEFORE UPDATE ON student_addresses
    FOR EACH ROW
    EXECUTE FUNCTION update_students_updated_at();
