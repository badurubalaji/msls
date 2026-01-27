-- Migration: 000016_merit_list_admission_decision.up.sql
-- Description: Create merit_lists and admission_decisions tables for Story 3.7

-- Create decision_type enum
CREATE TYPE decision_type AS ENUM (
    'approved',
    'waitlisted',
    'rejected'
);

-- Create admission_decisions table
CREATE TABLE admission_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL,  -- References admission_applications when created
    decision decision_type NOT NULL,
    decision_date DATE NOT NULL,
    decided_by UUID REFERENCES users(id),
    section_assigned VARCHAR(50),
    waitlist_position INT,
    rejection_reason TEXT,
    offer_letter_url VARCHAR(500),
    offer_valid_until DATE,
    offer_accepted BOOLEAN,
    offer_accepted_at TIMESTAMPTZ,
    remarks TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT admission_decisions_waitlist_check CHECK (
        (decision = 'waitlisted' AND waitlist_position IS NOT NULL AND waitlist_position > 0) OR
        (decision != 'waitlisted' AND waitlist_position IS NULL)
    ),
    CONSTRAINT admission_decisions_rejection_check CHECK (
        (decision = 'rejected' AND rejection_reason IS NOT NULL) OR
        (decision != 'rejected')
    )
);

-- Create merit_lists table (snapshot for audit)
CREATE TABLE merit_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id) ON DELETE CASCADE,
    class_name VARCHAR(50) NOT NULL,
    test_id UUID,  -- References entrance_tests when created
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    generated_by UUID REFERENCES users(id),
    cutoff_score DECIMAL(5,2),
    entries JSONB NOT NULL DEFAULT '[]',
    is_final BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure unique merit list per session, class, and test
    CONSTRAINT merit_lists_session_class_test_unique UNIQUE(session_id, class_name, test_id)
);

-- Create indexes for admission_decisions
CREATE INDEX idx_admission_decisions_tenant ON admission_decisions(tenant_id);
CREATE INDEX idx_admission_decisions_application ON admission_decisions(application_id);
CREATE INDEX idx_admission_decisions_decision ON admission_decisions(tenant_id, decision);
CREATE INDEX idx_admission_decisions_decided_by ON admission_decisions(decided_by);
CREATE INDEX idx_admission_decisions_decision_date ON admission_decisions(tenant_id, decision_date);
CREATE INDEX idx_admission_decisions_waitlist ON admission_decisions(tenant_id, waitlist_position) WHERE decision = 'waitlisted';

-- Create indexes for merit_lists
CREATE INDEX idx_merit_lists_tenant ON merit_lists(tenant_id);
CREATE INDEX idx_merit_lists_session ON merit_lists(session_id);
CREATE INDEX idx_merit_lists_class ON merit_lists(tenant_id, class_name);
CREATE INDEX idx_merit_lists_test ON merit_lists(test_id) WHERE test_id IS NOT NULL;
CREATE INDEX idx_merit_lists_generated_at ON merit_lists(generated_at);
CREATE INDEX idx_merit_lists_is_final ON merit_lists(tenant_id, is_final);

-- Enable RLS on admission_decisions
ALTER TABLE admission_decisions ENABLE ROW LEVEL SECURITY;

-- Create RLS policies for admission_decisions
CREATE POLICY tenant_isolation_admission_decisions ON admission_decisions
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_admission_decisions ON admission_decisions
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- Enable RLS on merit_lists
ALTER TABLE merit_lists ENABLE ROW LEVEL SECURITY;

-- Create RLS policies for merit_lists
CREATE POLICY tenant_isolation_merit_lists ON merit_lists
    FOR ALL
    USING (
        tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID)
    );

CREATE POLICY bypass_rls_merit_lists ON merit_lists
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_admission_decision_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_admission_decisions_updated_at
    BEFORE UPDATE ON admission_decisions
    FOR EACH ROW
    EXECUTE FUNCTION update_admission_decision_updated_at();

-- Add comments
COMMENT ON TABLE admission_decisions IS 'Stores admission decisions for applications';
COMMENT ON COLUMN admission_decisions.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN admission_decisions.tenant_id IS 'Reference to the tenant (school)';
COMMENT ON COLUMN admission_decisions.application_id IS 'Reference to the admission application';
COMMENT ON COLUMN admission_decisions.decision IS 'Decision type: approved, waitlisted, or rejected';
COMMENT ON COLUMN admission_decisions.decision_date IS 'Date when the decision was made';
COMMENT ON COLUMN admission_decisions.decided_by IS 'User who made the decision';
COMMENT ON COLUMN admission_decisions.section_assigned IS 'Section assigned for approved applications';
COMMENT ON COLUMN admission_decisions.waitlist_position IS 'Position in waitlist for waitlisted applications';
COMMENT ON COLUMN admission_decisions.rejection_reason IS 'Reason for rejection';
COMMENT ON COLUMN admission_decisions.offer_letter_url IS 'URL to the generated offer letter PDF';
COMMENT ON COLUMN admission_decisions.offer_valid_until IS 'Date until the offer is valid';
COMMENT ON COLUMN admission_decisions.offer_accepted IS 'Whether the offer was accepted';
COMMENT ON COLUMN admission_decisions.offer_accepted_at IS 'Timestamp when offer was accepted';
COMMENT ON COLUMN admission_decisions.remarks IS 'Additional remarks about the decision';

COMMENT ON TABLE merit_lists IS 'Stores merit list snapshots for admission sessions';
COMMENT ON COLUMN merit_lists.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN merit_lists.tenant_id IS 'Reference to the tenant (school)';
COMMENT ON COLUMN merit_lists.session_id IS 'Reference to the admission session';
COMMENT ON COLUMN merit_lists.class_name IS 'Class name for which merit list is generated';
COMMENT ON COLUMN merit_lists.test_id IS 'Reference to the entrance test (optional)';
COMMENT ON COLUMN merit_lists.generated_at IS 'When the merit list was generated';
COMMENT ON COLUMN merit_lists.generated_by IS 'User who generated the merit list';
COMMENT ON COLUMN merit_lists.cutoff_score IS 'Cutoff score applied to filter candidates';
COMMENT ON COLUMN merit_lists.entries IS 'JSON array of merit list entries with rank, applicant details, and scores';
COMMENT ON COLUMN merit_lists.is_final IS 'Whether this is the final merit list';

COMMENT ON POLICY tenant_isolation_admission_decisions ON admission_decisions IS 'Restricts admission decision access to the current tenant';
COMMENT ON POLICY bypass_rls_admission_decisions ON admission_decisions IS 'Allows bypass of RLS for admin operations';
COMMENT ON POLICY tenant_isolation_merit_lists ON merit_lists IS 'Restricts merit list access to the current tenant';
COMMENT ON POLICY bypass_rls_merit_lists ON merit_lists IS 'Allows bypass of RLS for admin operations';
