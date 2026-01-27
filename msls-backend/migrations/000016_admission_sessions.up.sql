-- Migration: 000015_admission_sessions.up.sql
-- Description: Create admission sessions and seats tables for admissions module

-- Create admission_sessions table
CREATE TABLE admission_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    academic_year_id UUID,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'upcoming',
    application_fee DECIMAL(10,2) DEFAULT 0,
    required_documents JSONB DEFAULT '[]',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT admission_sessions_status_check CHECK (status IN ('upcoming', 'open', 'closed')),
    CONSTRAINT admission_sessions_date_check CHECK (end_date >= start_date),
    CONSTRAINT admission_sessions_tenant_name_unique UNIQUE(tenant_id, name, academic_year_id)
);

-- Create admission_seats table (per class)
CREATE TABLE admission_seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id) ON DELETE CASCADE,
    class_name VARCHAR(100) NOT NULL,
    total_seats INT NOT NULL DEFAULT 0,
    filled_seats INT NOT NULL DEFAULT 0,
    waitlist_limit INT DEFAULT 10,
    reserved_seats JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT admission_seats_total_positive CHECK (total_seats >= 0),
    CONSTRAINT admission_seats_filled_positive CHECK (filled_seats >= 0),
    CONSTRAINT admission_seats_filled_check CHECK (filled_seats <= total_seats),
    CONSTRAINT admission_seats_session_class_unique UNIQUE(session_id, class_name)
);

-- Create indexes for admission_sessions
CREATE INDEX idx_admission_sessions_tenant_id ON admission_sessions(tenant_id);
CREATE INDEX idx_admission_sessions_branch_id ON admission_sessions(branch_id);
CREATE INDEX idx_admission_sessions_academic_year_id ON admission_sessions(academic_year_id);
CREATE INDEX idx_admission_sessions_status ON admission_sessions(status);
CREATE INDEX idx_admission_sessions_start_date ON admission_sessions(start_date);
CREATE INDEX idx_admission_sessions_end_date ON admission_sessions(end_date);

-- Create indexes for admission_seats
CREATE INDEX idx_admission_seats_tenant_id ON admission_seats(tenant_id);
CREATE INDEX idx_admission_seats_session_id ON admission_seats(session_id);
CREATE INDEX idx_admission_seats_class_name ON admission_seats(class_name);

-- Enable RLS on admission_sessions
ALTER TABLE admission_sessions ENABLE ROW LEVEL SECURITY;

-- Create RLS policy for admission_sessions
CREATE POLICY tenant_isolation_admission_sessions ON admission_sessions
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_admission_sessions ON admission_sessions
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- Enable RLS on admission_seats
ALTER TABLE admission_seats ENABLE ROW LEVEL SECURITY;

-- Create RLS policy for admission_seats
CREATE POLICY tenant_isolation_admission_seats ON admission_seats
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_admission_seats ON admission_seats
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- Add comments
COMMENT ON TABLE admission_sessions IS 'Admission sessions for organizing admission cycles';
COMMENT ON COLUMN admission_sessions.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN admission_sessions.tenant_id IS 'Reference to the parent tenant';
COMMENT ON COLUMN admission_sessions.branch_id IS 'Optional reference to a specific branch';
COMMENT ON COLUMN admission_sessions.academic_year_id IS 'Reference to the academic year';
COMMENT ON COLUMN admission_sessions.name IS 'Name of the admission session';
COMMENT ON COLUMN admission_sessions.description IS 'Description of the admission session';
COMMENT ON COLUMN admission_sessions.start_date IS 'Start date of the admission session';
COMMENT ON COLUMN admission_sessions.end_date IS 'End date of the admission session';
COMMENT ON COLUMN admission_sessions.status IS 'Session status: upcoming, open, or closed';
COMMENT ON COLUMN admission_sessions.application_fee IS 'Application fee amount';
COMMENT ON COLUMN admission_sessions.required_documents IS 'JSON array of required document types';
COMMENT ON COLUMN admission_sessions.settings IS 'Additional session settings as JSON';

COMMENT ON TABLE admission_seats IS 'Seat configuration per class for admission sessions';
COMMENT ON COLUMN admission_seats.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN admission_seats.tenant_id IS 'Reference to the parent tenant';
COMMENT ON COLUMN admission_seats.session_id IS 'Reference to the admission session';
COMMENT ON COLUMN admission_seats.class_name IS 'Name of the class (e.g., Grade 1, Class 10)';
COMMENT ON COLUMN admission_seats.total_seats IS 'Total available seats for this class';
COMMENT ON COLUMN admission_seats.filled_seats IS 'Number of seats already filled';
COMMENT ON COLUMN admission_seats.waitlist_limit IS 'Maximum number of waitlist entries';
COMMENT ON COLUMN admission_seats.reserved_seats IS 'JSON object for seat reservations by category';

-- Add admission permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('admissions:read', 'View Admissions', 'admissions', 'View admission sessions and applications'),
    ('admissions:create', 'Create Admissions', 'admissions', 'Create admission sessions'),
    ('admissions:update', 'Update Admissions', 'admissions', 'Update admission sessions and seat configurations'),
    ('admissions:delete', 'Delete Admissions', 'admissions', 'Delete admission sessions');

-- Assign admission permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.module = 'admissions';

-- Assign admission permissions to principal role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'principal'
AND p.code IN ('admissions:read', 'admissions:create', 'admissions:update');

-- Assign super_admin all admission permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.module = 'admissions';

COMMENT ON POLICY tenant_isolation_admission_sessions ON admission_sessions IS 'Restricts admission session access to the current tenant';
COMMENT ON POLICY bypass_rls_admission_sessions ON admission_sessions IS 'Allows bypass of RLS for admin operations';
COMMENT ON POLICY tenant_isolation_admission_seats ON admission_seats IS 'Restricts admission seat access to the current tenant';
COMMENT ON POLICY bypass_rls_admission_seats ON admission_seats IS 'Allows bypass of RLS for admin operations';
