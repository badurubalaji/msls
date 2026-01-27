-- Migration: 000016_admission_applications.up.sql
-- Description: Create admission applications and related tables for Story 3.6

-- Create application status type
CREATE TYPE application_status AS ENUM (
    'draft',
    'submitted',
    'under_review',
    'documents_pending',
    'test_scheduled',
    'test_completed',
    'shortlisted',
    'approved',
    'rejected',
    'waitlisted',
    'enrolled'
);

-- Create document verification status type
CREATE TYPE verification_status AS ENUM (
    'pending',
    'verified',
    'rejected',
    'resubmit_required'
);

-- Create test registration status type
CREATE TYPE test_registration_status AS ENUM (
    'registered',
    'hall_ticket_generated',
    'appeared',
    'absent',
    'cancelled'
);

-- Create test result type
CREATE TYPE test_result AS ENUM (
    'pass',
    'fail',
    'merit',
    'distinction'
);

-- Create review type enum
CREATE TYPE review_type AS ENUM (
    'initial_screening',
    'document_verification',
    'academic_review',
    'interview',
    'final_decision'
);

-- Create review status enum
CREATE TYPE review_status AS ENUM (
    'approved',
    'rejected',
    'pending_info',
    'escalated'
);

-- ============================================================================
-- Admission Applications Table
-- ============================================================================
CREATE TABLE admission_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    enquiry_id UUID REFERENCES admission_enquiries(id) ON DELETE SET NULL,
    application_number VARCHAR(50) NOT NULL,

    -- Student Information
    student_name VARCHAR(200) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(20) NOT NULL,
    blood_group VARCHAR(10),
    nationality VARCHAR(100) DEFAULT 'Indian',
    religion VARCHAR(100),
    caste VARCHAR(100),
    category VARCHAR(50),
    mother_tongue VARCHAR(100),
    aadhar_number VARCHAR(12),

    -- Academic Information
    class_applying VARCHAR(50) NOT NULL,
    previous_school VARCHAR(200),
    previous_class VARCHAR(50),
    previous_percentage DECIMAL(5,2),
    medium_of_instruction VARCHAR(50),

    -- Contact Information
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'India',

    -- Parent/Guardian Information
    father_name VARCHAR(200),
    father_phone VARCHAR(20),
    father_email VARCHAR(255),
    father_occupation VARCHAR(200),
    father_qualification VARCHAR(100),
    mother_name VARCHAR(200),
    mother_phone VARCHAR(20),
    mother_email VARCHAR(255),
    mother_occupation VARCHAR(200),
    mother_qualification VARCHAR(100),
    guardian_name VARCHAR(200),
    guardian_phone VARCHAR(20),
    guardian_email VARCHAR(255),
    guardian_relation VARCHAR(50),

    -- Application Details
    status application_status NOT NULL DEFAULT 'draft',
    submitted_at TIMESTAMPTZ,
    remarks TEXT,
    internal_notes TEXT,
    priority INT DEFAULT 0,

    -- Payment Information
    fee_paid BOOLEAN DEFAULT false,
    payment_reference VARCHAR(100),
    payment_date TIMESTAMPTZ,

    -- Additional Data
    extra_data JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT uniq_applications_tenant_number UNIQUE(tenant_id, application_number)
);

-- ============================================================================
-- Application Documents Table
-- ============================================================================
CREATE TABLE application_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    document_type VARCHAR(100) NOT NULL,
    document_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size INT,
    mime_type VARCHAR(100),
    verification_status verification_status NOT NULL DEFAULT 'pending',
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    verification_remarks TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Entrance Tests Table
-- ============================================================================
CREATE TABLE entrance_tests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id) ON DELETE CASCADE,
    test_name VARCHAR(200) NOT NULL,
    test_date DATE NOT NULL,
    start_time TIME NOT NULL,
    duration_minutes INT NOT NULL DEFAULT 60,
    venue VARCHAR(200),
    class_names JSONB DEFAULT '[]',
    max_candidates INT DEFAULT 100,
    status VARCHAR(20) DEFAULT 'scheduled',
    subjects JSONB DEFAULT '[]',
    instructions TEXT,
    passing_percentage DECIMAL(5,2) DEFAULT 33.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    CONSTRAINT entrance_tests_status_check CHECK (status IN ('scheduled', 'in_progress', 'completed', 'cancelled')),
    CONSTRAINT entrance_tests_duration_check CHECK (duration_minutes > 0)
);

-- ============================================================================
-- Test Registrations Table
-- ============================================================================
CREATE TABLE test_registrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES entrance_tests(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    roll_number VARCHAR(20),
    status test_registration_status NOT NULL DEFAULT 'registered',
    marks JSONB DEFAULT '{}',
    total_marks DECIMAL(6,2),
    max_marks DECIMAL(6,2),
    percentage DECIMAL(5,2),
    result test_result,
    remarks TEXT,
    hall_ticket_generated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_test_registration UNIQUE(test_id, application_id),
    CONSTRAINT uniq_test_roll_number UNIQUE(test_id, roll_number)
);

-- ============================================================================
-- Application Reviews Table
-- ============================================================================
CREATE TABLE application_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(id),
    review_type review_type NOT NULL,
    status review_status NOT NULL,
    comments TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Application Number Sequence Table
-- ============================================================================
CREATE TABLE application_number_sequence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id) ON DELETE CASCADE,
    last_sequence INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_app_seq_tenant_session UNIQUE(tenant_id, session_id)
);

-- ============================================================================
-- Roll Number Sequence Table (for entrance tests)
-- ============================================================================
CREATE TABLE roll_number_sequence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES entrance_tests(id) ON DELETE CASCADE,
    last_sequence INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_roll_seq_tenant_test UNIQUE(tenant_id, test_id)
);

-- ============================================================================
-- Indexes
-- ============================================================================

-- Admission Applications Indexes
CREATE INDEX idx_applications_tenant ON admission_applications(tenant_id);
CREATE INDEX idx_applications_session ON admission_applications(session_id);
CREATE INDEX idx_applications_branch ON admission_applications(branch_id);
CREATE INDEX idx_applications_enquiry ON admission_applications(enquiry_id);
CREATE INDEX idx_applications_status ON admission_applications(tenant_id, status);
CREATE INDEX idx_applications_class ON admission_applications(tenant_id, class_applying);
CREATE INDEX idx_applications_submitted_at ON admission_applications(tenant_id, submitted_at);
CREATE INDEX idx_applications_deleted_at ON admission_applications(deleted_at) WHERE deleted_at IS NULL;

-- Application Documents Indexes
CREATE INDEX idx_app_docs_tenant ON application_documents(tenant_id);
CREATE INDEX idx_app_docs_application ON application_documents(application_id);
CREATE INDEX idx_app_docs_status ON application_documents(verification_status);

-- Entrance Tests Indexes
CREATE INDEX idx_entrance_tests_tenant ON entrance_tests(tenant_id);
CREATE INDEX idx_entrance_tests_session ON entrance_tests(session_id);
CREATE INDEX idx_entrance_tests_date ON entrance_tests(test_date);
CREATE INDEX idx_entrance_tests_status ON entrance_tests(status);

-- Test Registrations Indexes
CREATE INDEX idx_test_registrations_tenant ON test_registrations(tenant_id);
CREATE INDEX idx_test_registrations_test ON test_registrations(test_id);
CREATE INDEX idx_test_registrations_application ON test_registrations(application_id);
CREATE INDEX idx_test_registrations_status ON test_registrations(status);

-- Application Reviews Indexes
CREATE INDEX idx_app_reviews_tenant ON application_reviews(tenant_id);
CREATE INDEX idx_app_reviews_application ON application_reviews(application_id);
CREATE INDEX idx_app_reviews_reviewer ON application_reviews(reviewer_id);
CREATE INDEX idx_app_reviews_type ON application_reviews(review_type);

-- ============================================================================
-- Row Level Security
-- ============================================================================

-- Enable RLS
ALTER TABLE admission_applications ENABLE ROW LEVEL SECURITY;
ALTER TABLE application_documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE entrance_tests ENABLE ROW LEVEL SECURITY;
ALTER TABLE test_registrations ENABLE ROW LEVEL SECURITY;
ALTER TABLE application_reviews ENABLE ROW LEVEL SECURITY;
ALTER TABLE application_number_sequence ENABLE ROW LEVEL SECURITY;
ALTER TABLE roll_number_sequence ENABLE ROW LEVEL SECURITY;

-- RLS Policies for admission_applications
CREATE POLICY tenant_isolation_applications ON admission_applications
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_applications ON admission_applications
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- RLS Policies for application_documents
CREATE POLICY tenant_isolation_app_docs ON application_documents
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_app_docs ON application_documents
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- RLS Policies for entrance_tests
CREATE POLICY tenant_isolation_entrance_tests ON entrance_tests
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_entrance_tests ON entrance_tests
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- RLS Policies for test_registrations
CREATE POLICY tenant_isolation_test_registrations ON test_registrations
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_test_registrations ON test_registrations
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- RLS Policies for application_reviews
CREATE POLICY tenant_isolation_app_reviews ON application_reviews
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_app_reviews ON application_reviews
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- RLS Policies for application_number_sequence
CREATE POLICY tenant_isolation_app_seq ON application_number_sequence
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_app_seq ON application_number_sequence
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- RLS Policies for roll_number_sequence
CREATE POLICY tenant_isolation_roll_seq ON roll_number_sequence
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_roll_seq ON roll_number_sequence
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- ============================================================================
-- Functions
-- ============================================================================

-- Function to generate next application number
CREATE OR REPLACE FUNCTION get_next_application_number(p_tenant_id UUID, p_session_id UUID)
RETURNS VARCHAR(50) AS $$
DECLARE
    v_next_seq INTEGER;
    v_year VARCHAR(4);
    v_application_number VARCHAR(50);
BEGIN
    -- Get current year
    v_year := TO_CHAR(CURRENT_DATE, 'YYYY');

    -- Insert or update the sequence
    INSERT INTO application_number_sequence (tenant_id, session_id, last_sequence)
    VALUES (p_tenant_id, p_session_id, 1)
    ON CONFLICT (tenant_id, session_id)
    DO UPDATE SET
        last_sequence = application_number_sequence.last_sequence + 1,
        updated_at = NOW()
    RETURNING last_sequence INTO v_next_seq;

    -- Format: APP-YYYY-XXXXX
    v_application_number := 'APP-' || v_year || '-' || LPAD(v_next_seq::TEXT, 5, '0');

    RETURN v_application_number;
END;
$$ LANGUAGE plpgsql;

-- Function to generate next roll number
CREATE OR REPLACE FUNCTION get_next_roll_number(p_tenant_id UUID, p_test_id UUID)
RETURNS VARCHAR(20) AS $$
DECLARE
    v_next_seq INTEGER;
    v_roll_number VARCHAR(20);
BEGIN
    -- Insert or update the sequence
    INSERT INTO roll_number_sequence (tenant_id, test_id, last_sequence)
    VALUES (p_tenant_id, p_test_id, 1)
    ON CONFLICT (tenant_id, test_id)
    DO UPDATE SET
        last_sequence = roll_number_sequence.last_sequence + 1,
        updated_at = NOW()
    RETURNING last_sequence INTO v_next_seq;

    -- Format: ROLL-XXXXX
    v_roll_number := 'ROLL-' || LPAD(v_next_seq::TEXT, 5, '0');

    RETURN v_roll_number;
END;
$$ LANGUAGE plpgsql;

-- Trigger function to update application status when test results are submitted
CREATE OR REPLACE FUNCTION update_application_on_test_result()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.result IS NOT NULL AND OLD.result IS NULL THEN
        UPDATE admission_applications
        SET status = 'test_completed', updated_at = NOW()
        WHERE id = NEW.application_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_app_on_test_result
    AFTER UPDATE ON test_registrations
    FOR EACH ROW
    EXECUTE FUNCTION update_application_on_test_result();

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_applications_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_applications_updated_at
    BEFORE UPDATE ON admission_applications
    FOR EACH ROW
    EXECUTE FUNCTION update_applications_updated_at();

CREATE TRIGGER trigger_app_docs_updated_at
    BEFORE UPDATE ON application_documents
    FOR EACH ROW
    EXECUTE FUNCTION update_applications_updated_at();

CREATE TRIGGER trigger_entrance_tests_updated_at
    BEFORE UPDATE ON entrance_tests
    FOR EACH ROW
    EXECUTE FUNCTION update_applications_updated_at();

CREATE TRIGGER trigger_test_registrations_updated_at
    BEFORE UPDATE ON test_registrations
    FOR EACH ROW
    EXECUTE FUNCTION update_applications_updated_at();

-- ============================================================================
-- Permissions
-- ============================================================================

-- Add application-specific permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('applications:read', 'View Applications', 'admissions', 'View admission applications'),
    ('applications:create', 'Create Applications', 'admissions', 'Create admission applications'),
    ('applications:update', 'Update Applications', 'admissions', 'Update admission applications and verify documents'),
    ('applications:delete', 'Delete Applications', 'admissions', 'Delete admission applications');

-- Assign to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.code IN ('applications:read', 'applications:create', 'applications:update', 'applications:delete');

-- Assign to principal role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'principal'
AND p.code IN ('applications:read', 'applications:create', 'applications:update');

-- Assign to super_admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.code IN ('applications:read', 'applications:create', 'applications:update', 'applications:delete');

-- ============================================================================
-- Comments
-- ============================================================================

COMMENT ON TABLE admission_applications IS 'Stores admission applications from prospective students';
COMMENT ON COLUMN admission_applications.application_number IS 'Auto-generated unique application number (APP-YYYY-XXXXX)';
COMMENT ON COLUMN admission_applications.status IS 'Current status of the application in the admission workflow';
COMMENT ON COLUMN admission_applications.extra_data IS 'Additional custom fields as JSON';

COMMENT ON TABLE application_documents IS 'Documents submitted with admission applications';
COMMENT ON COLUMN application_documents.verification_status IS 'Document verification status';

COMMENT ON TABLE entrance_tests IS 'Entrance tests scheduled for admission sessions';
COMMENT ON COLUMN entrance_tests.subjects IS 'JSON array of subjects with max marks';
COMMENT ON COLUMN entrance_tests.class_names IS 'JSON array of class names this test is for';

COMMENT ON TABLE test_registrations IS 'Student registrations for entrance tests';
COMMENT ON COLUMN test_registrations.marks IS 'JSON object with marks per subject';
COMMENT ON COLUMN test_registrations.roll_number IS 'Auto-generated roll number for the test';

COMMENT ON TABLE application_reviews IS 'Review history for admission applications';
COMMENT ON COLUMN application_reviews.review_type IS 'Type of review conducted';
COMMENT ON COLUMN application_reviews.status IS 'Outcome of the review';

COMMENT ON FUNCTION get_next_application_number(UUID, UUID) IS 'Generates next sequential application number for a session';
COMMENT ON FUNCTION get_next_roll_number(UUID, UUID) IS 'Generates next sequential roll number for a test';
