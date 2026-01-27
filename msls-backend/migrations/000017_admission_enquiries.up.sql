-- Migration: 000015_admission_enquiries.up.sql
-- Description: Create admission enquiries and follow-ups tables

-- Create enquiry source enum
CREATE TYPE enquiry_source AS ENUM (
    'walk_in',
    'phone',
    'website',
    'referral',
    'advertisement',
    'social_media',
    'other'
);

-- Create enquiry status enum
CREATE TYPE enquiry_status AS ENUM (
    'new',
    'contacted',
    'interested',
    'converted',
    'closed'
);

-- Create follow-up contact mode enum
CREATE TYPE contact_mode AS ENUM (
    'phone',
    'email',
    'whatsapp',
    'in_person',
    'other'
);

-- Create follow-up outcome enum
CREATE TYPE follow_up_outcome AS ENUM (
    'interested',
    'not_interested',
    'follow_up_required',
    'converted',
    'no_response'
);

-- Create admission_enquiries table
CREATE TABLE admission_enquiries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    session_id UUID,  -- Will reference admission_sessions when that table is created
    enquiry_number VARCHAR(50) NOT NULL,
    student_name VARCHAR(200) NOT NULL,
    date_of_birth DATE,
    gender VARCHAR(20),
    class_applying VARCHAR(50) NOT NULL,
    parent_name VARCHAR(200) NOT NULL,
    parent_phone VARCHAR(20) NOT NULL,
    parent_email VARCHAR(255),
    source enquiry_source NOT NULL DEFAULT 'walk_in',
    referral_details TEXT,
    remarks TEXT,
    status enquiry_status NOT NULL DEFAULT 'new',
    follow_up_date DATE,
    assigned_to UUID REFERENCES users(id),
    converted_application_id UUID,  -- Will reference admission_applications when created
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    deleted_at TIMESTAMPTZ,

    -- Ensure enquiry number is unique per tenant
    CONSTRAINT uniq_enquiries_tenant_number UNIQUE(tenant_id, enquiry_number)
);

-- Create enquiry_follow_ups table
CREATE TABLE enquiry_follow_ups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    enquiry_id UUID NOT NULL REFERENCES admission_enquiries(id) ON DELETE CASCADE,
    follow_up_date DATE NOT NULL,
    contact_mode contact_mode NOT NULL DEFAULT 'phone',
    notes TEXT,
    outcome follow_up_outcome,
    next_follow_up DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Create enquiry number sequence table for auto-generation
CREATE TABLE enquiry_number_sequence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    date_key VARCHAR(8) NOT NULL,  -- Format: YYYYMMDD
    last_sequence INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_enquiry_seq_tenant_date UNIQUE(tenant_id, date_key)
);

-- Create indexes for admission_enquiries
CREATE INDEX idx_enquiries_tenant ON admission_enquiries(tenant_id);
CREATE INDEX idx_enquiries_branch ON admission_enquiries(branch_id);
CREATE INDEX idx_enquiries_status ON admission_enquiries(tenant_id, status);
CREATE INDEX idx_enquiries_class ON admission_enquiries(tenant_id, class_applying);
CREATE INDEX idx_enquiries_created_at ON admission_enquiries(tenant_id, created_at);
CREATE INDEX idx_enquiries_follow_up_date ON admission_enquiries(tenant_id, follow_up_date) WHERE follow_up_date IS NOT NULL;
CREATE INDEX idx_enquiries_assigned_to ON admission_enquiries(assigned_to) WHERE assigned_to IS NOT NULL;
CREATE INDEX idx_enquiries_phone ON admission_enquiries(tenant_id, parent_phone);
CREATE INDEX idx_enquiries_deleted_at ON admission_enquiries(deleted_at) WHERE deleted_at IS NULL;

-- Create indexes for enquiry_follow_ups
CREATE INDEX idx_follow_ups_tenant ON enquiry_follow_ups(tenant_id);
CREATE INDEX idx_follow_ups_enquiry ON enquiry_follow_ups(enquiry_id);
CREATE INDEX idx_follow_ups_date ON enquiry_follow_ups(follow_up_date);
CREATE INDEX idx_follow_ups_created_at ON enquiry_follow_ups(created_at);

-- Create index for enquiry_number_sequence
CREATE INDEX idx_enquiry_seq_tenant ON enquiry_number_sequence(tenant_id);

-- Enable Row Level Security
ALTER TABLE admission_enquiries ENABLE ROW LEVEL SECURITY;
ALTER TABLE enquiry_follow_ups ENABLE ROW LEVEL SECURITY;
ALTER TABLE enquiry_number_sequence ENABLE ROW LEVEL SECURITY;

-- Create RLS policies for admission_enquiries
CREATE POLICY tenant_isolation_enquiries ON admission_enquiries
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_enquiries_insert ON admission_enquiries
    FOR INSERT WITH CHECK (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_enquiries_update ON admission_enquiries
    FOR UPDATE USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_enquiries_delete ON admission_enquiries
    FOR DELETE USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Create RLS policies for enquiry_follow_ups
CREATE POLICY tenant_isolation_follow_ups ON enquiry_follow_ups
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_follow_ups_insert ON enquiry_follow_ups
    FOR INSERT WITH CHECK (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_follow_ups_update ON enquiry_follow_ups
    FOR UPDATE USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_follow_ups_delete ON enquiry_follow_ups
    FOR DELETE USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Create RLS policies for enquiry_number_sequence
CREATE POLICY tenant_isolation_seq ON enquiry_number_sequence
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_seq_insert ON enquiry_number_sequence
    FOR INSERT WITH CHECK (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_seq_update ON enquiry_number_sequence
    FOR UPDATE USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Add comments for documentation
COMMENT ON TABLE admission_enquiries IS 'Stores admission enquiries from prospective students/parents';
COMMENT ON COLUMN admission_enquiries.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN admission_enquiries.tenant_id IS 'Reference to the tenant (school) this enquiry belongs to';
COMMENT ON COLUMN admission_enquiries.branch_id IS 'Reference to the branch this enquiry is for';
COMMENT ON COLUMN admission_enquiries.session_id IS 'Reference to the admission session (academic year)';
COMMENT ON COLUMN admission_enquiries.enquiry_number IS 'Auto-generated unique enquiry number (ENQ-YYYYMMDD-XXXX)';
COMMENT ON COLUMN admission_enquiries.student_name IS 'Name of the prospective student';
COMMENT ON COLUMN admission_enquiries.class_applying IS 'Class/grade the student is applying for';
COMMENT ON COLUMN admission_enquiries.parent_name IS 'Name of the parent/guardian making the enquiry';
COMMENT ON COLUMN admission_enquiries.parent_phone IS 'Contact phone number of parent/guardian';
COMMENT ON COLUMN admission_enquiries.source IS 'How the enquiry was received';
COMMENT ON COLUMN admission_enquiries.status IS 'Current status of the enquiry';
COMMENT ON COLUMN admission_enquiries.follow_up_date IS 'Next scheduled follow-up date';
COMMENT ON COLUMN admission_enquiries.assigned_to IS 'Staff member assigned to handle this enquiry';
COMMENT ON COLUMN admission_enquiries.converted_application_id IS 'Reference to application if converted';

COMMENT ON TABLE enquiry_follow_ups IS 'Tracks follow-up interactions for admission enquiries';
COMMENT ON COLUMN enquiry_follow_ups.enquiry_id IS 'Reference to the parent enquiry';
COMMENT ON COLUMN enquiry_follow_ups.contact_mode IS 'How the follow-up was conducted';
COMMENT ON COLUMN enquiry_follow_ups.outcome IS 'Result of the follow-up interaction';
COMMENT ON COLUMN enquiry_follow_ups.next_follow_up IS 'Suggested date for next follow-up';

COMMENT ON TABLE enquiry_number_sequence IS 'Tracks sequence numbers for enquiry number generation per tenant per day';

-- Create function to get next enquiry number
CREATE OR REPLACE FUNCTION get_next_enquiry_number(p_tenant_id UUID)
RETURNS VARCHAR(50) AS $$
DECLARE
    v_date_key VARCHAR(8);
    v_next_seq INTEGER;
    v_enquiry_number VARCHAR(50);
BEGIN
    -- Get current date in YYYYMMDD format
    v_date_key := TO_CHAR(CURRENT_DATE, 'YYYYMMDD');

    -- Insert or update the sequence, returning the new sequence number
    INSERT INTO enquiry_number_sequence (tenant_id, date_key, last_sequence)
    VALUES (p_tenant_id, v_date_key, 1)
    ON CONFLICT (tenant_id, date_key)
    DO UPDATE SET
        last_sequence = enquiry_number_sequence.last_sequence + 1,
        updated_at = NOW()
    RETURNING last_sequence INTO v_next_seq;

    -- Format: ENQ-YYYYMMDD-XXXX
    v_enquiry_number := 'ENQ-' || v_date_key || '-' || LPAD(v_next_seq::TEXT, 4, '0');

    RETURN v_enquiry_number;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_next_enquiry_number(UUID) IS 'Generates next sequential enquiry number for a tenant';

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_enquiry_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_enquiries_updated_at
    BEFORE UPDATE ON admission_enquiries
    FOR EACH ROW
    EXECUTE FUNCTION update_enquiry_updated_at();
