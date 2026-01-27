-- Migration: 000035_staff.up.sql
-- Description: Create staff table and staff_status_history for Epic 5 Story 5.1

-- Staff table (core entity)
CREATE TABLE staff (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),

    -- Employee identification
    employee_id VARCHAR(50) NOT NULL,
    employee_id_prefix VARCHAR(10) NOT NULL DEFAULT 'EMP',

    -- Personal details
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(20) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    blood_group VARCHAR(10),
    nationality VARCHAR(50) DEFAULT 'Indian',
    religion VARCHAR(50),
    marital_status VARCHAR(20),

    -- Contact details
    personal_email VARCHAR(255),
    work_email VARCHAR(255) NOT NULL,
    personal_phone VARCHAR(20),
    work_phone VARCHAR(20) NOT NULL,
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(20),
    emergency_contact_relation VARCHAR(50),

    -- Current Address
    current_address_line1 VARCHAR(255),
    current_address_line2 VARCHAR(255),
    current_city VARCHAR(100),
    current_state VARCHAR(100),
    current_pincode VARCHAR(10),
    current_country VARCHAR(100) DEFAULT 'India',

    -- Permanent Address
    permanent_address_line1 VARCHAR(255),
    permanent_address_line2 VARCHAR(255),
    permanent_city VARCHAR(100),
    permanent_state VARCHAR(100),
    permanent_pincode VARCHAR(10),
    permanent_country VARCHAR(100) DEFAULT 'India',
    same_as_current BOOLEAN DEFAULT false,

    -- Employment details
    staff_type VARCHAR(20) NOT NULL CHECK (staff_type IN ('teaching', 'non_teaching')),
    department_id UUID REFERENCES departments(id),
    designation_id UUID REFERENCES designations(id),
    reporting_manager_id UUID REFERENCES staff(id),
    join_date DATE NOT NULL,
    confirmation_date DATE,
    probation_end_date DATE,

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'terminated', 'on_leave')),
    status_reason TEXT,
    termination_date DATE,

    -- Profile
    photo_url VARCHAR(500),
    bio TEXT,

    -- Audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    deleted_at TIMESTAMPTZ,
    version INTEGER NOT NULL DEFAULT 1,

    CONSTRAINT uniq_staff_employee_id UNIQUE (tenant_id, employee_id)
);

-- Enable RLS
ALTER TABLE staff ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_staff ON staff
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_staff_tenant ON staff(tenant_id);
CREATE INDEX idx_staff_branch ON staff(branch_id);
CREATE INDEX idx_staff_department ON staff(department_id);
CREATE INDEX idx_staff_designation ON staff(designation_id);
CREATE INDEX idx_staff_employee_id ON staff(tenant_id, employee_id);
CREATE INDEX idx_staff_status ON staff(tenant_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_name ON staff(tenant_id, last_name, first_name);
CREATE INDEX idx_staff_type ON staff(tenant_id, staff_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_staff_deleted_at ON staff(deleted_at) WHERE deleted_at IS NOT NULL;

-- Employee ID sequence table for tracking sequences per tenant
CREATE TABLE staff_employee_sequences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    prefix VARCHAR(10) NOT NULL DEFAULT 'EMP',
    last_sequence INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_employee_sequence UNIQUE (tenant_id, prefix)
);

ALTER TABLE staff_employee_sequences ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_employee_sequences ON staff_employee_sequences
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE INDEX idx_employee_sequences_tenant ON staff_employee_sequences(tenant_id);

-- Staff status history table
CREATE TABLE staff_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    old_status VARCHAR(20),
    new_status VARCHAR(20) NOT NULL,
    reason TEXT,
    effective_date DATE NOT NULL,
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE staff_status_history ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_staff_status_history ON staff_status_history
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE INDEX idx_staff_status_history_staff ON staff_status_history(staff_id);
CREATE INDEX idx_staff_status_history_tenant ON staff_status_history(tenant_id);

-- Add FK from departments.head_id to staff.id now that staff table exists
ALTER TABLE departments ADD CONSTRAINT fk_department_head FOREIGN KEY (head_id) REFERENCES staff(id);

-- Add staff permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('staff:read', 'View Staff', 'Permission to view staff profiles', 'staff', NOW(), NOW()),
    ('staff:create', 'Create Staff', 'Permission to create staff profiles', 'staff', NOW(), NOW()),
    ('staff:update', 'Update Staff', 'Permission to update staff profiles', 'staff', NOW(), NOW()),
    ('staff:delete', 'Delete Staff', 'Permission to delete staff profiles', 'staff', NOW(), NOW()),
    ('staff:export', 'Export Staff', 'Permission to export staff data', 'staff', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_staff_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_staff_updated_at
    BEFORE UPDATE ON staff
    FOR EACH ROW
    EXECUTE FUNCTION update_staff_updated_at();

CREATE TRIGGER trigger_staff_employee_sequences_updated_at
    BEFORE UPDATE ON staff_employee_sequences
    FOR EACH ROW
    EXECUTE FUNCTION update_staff_updated_at();
