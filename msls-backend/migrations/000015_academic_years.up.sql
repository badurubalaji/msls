-- Migration: 000015_academic_years.up.sql
-- Description: Create academic years, terms, and holidays tables for school calendar management

-- Academic Years Table
CREATE TABLE academic_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT academic_years_tenant_branch_name_unique UNIQUE(tenant_id, branch_id, name),
    CONSTRAINT academic_years_dates_check CHECK (end_date > start_date)
);

-- Function to ensure only one current academic year per tenant/branch
CREATE OR REPLACE FUNCTION ensure_single_current_academic_year()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_current = TRUE THEN
        -- Unset is_current for other academic years in the same tenant/branch scope
        UPDATE academic_years
        SET is_current = FALSE, updated_at = NOW()
        WHERE tenant_id = NEW.tenant_id
          AND (
              (NEW.branch_id IS NULL AND branch_id IS NULL)
              OR (NEW.branch_id IS NOT NULL AND branch_id = NEW.branch_id)
          )
          AND id != NEW.id
          AND is_current = TRUE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to enforce single current academic year
CREATE TRIGGER trigger_single_current_academic_year
    BEFORE INSERT OR UPDATE ON academic_years
    FOR EACH ROW
    EXECUTE FUNCTION ensure_single_current_academic_year();

-- Academic Terms/Semesters Table
CREATE TABLE academic_terms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    sequence INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT academic_terms_year_name_unique UNIQUE(academic_year_id, name),
    CONSTRAINT academic_terms_dates_check CHECK (end_date > start_date),
    CONSTRAINT academic_terms_sequence_positive CHECK (sequence > 0)
);

-- Holidays Table
CREATE TABLE holidays (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    date DATE NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'public',
    is_optional BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT holidays_type_check CHECK (type IN ('public', 'religious', 'national', 'school', 'other'))
);

-- Create indexes
CREATE INDEX idx_academic_years_tenant_id ON academic_years(tenant_id);
CREATE INDEX idx_academic_years_branch_id ON academic_years(branch_id);
CREATE INDEX idx_academic_years_is_current ON academic_years(is_current) WHERE is_current = TRUE;
CREATE INDEX idx_academic_years_dates ON academic_years(start_date, end_date);

CREATE INDEX idx_academic_terms_tenant_id ON academic_terms(tenant_id);
CREATE INDEX idx_academic_terms_academic_year_id ON academic_terms(academic_year_id);
CREATE INDEX idx_academic_terms_dates ON academic_terms(start_date, end_date);
CREATE INDEX idx_academic_terms_sequence ON academic_terms(academic_year_id, sequence);

CREATE INDEX idx_holidays_tenant_id ON holidays(tenant_id);
CREATE INDEX idx_holidays_academic_year_id ON holidays(academic_year_id);
CREATE INDEX idx_holidays_branch_id ON holidays(branch_id);
CREATE INDEX idx_holidays_date ON holidays(date);
CREATE INDEX idx_holidays_type ON holidays(type);

-- Enable RLS
ALTER TABLE academic_years ENABLE ROW LEVEL SECURITY;
ALTER TABLE academic_terms ENABLE ROW LEVEL SECURITY;
ALTER TABLE holidays ENABLE ROW LEVEL SECURITY;

-- RLS policies for academic_years
CREATE POLICY tenant_isolation_academic_years ON academic_years
    FOR ALL
    USING (tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID));

CREATE POLICY bypass_rls_academic_years ON academic_years
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- RLS policies for academic_terms
CREATE POLICY tenant_isolation_academic_terms ON academic_terms
    FOR ALL
    USING (tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID));

CREATE POLICY bypass_rls_academic_terms ON academic_terms
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- RLS policies for holidays
CREATE POLICY tenant_isolation_holidays ON holidays
    FOR ALL
    USING (tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID));

CREATE POLICY bypass_rls_holidays ON holidays
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- Add comments
COMMENT ON TABLE academic_years IS 'Academic years for school calendar management';
COMMENT ON COLUMN academic_years.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN academic_years.tenant_id IS 'Reference to the parent tenant';
COMMENT ON COLUMN academic_years.branch_id IS 'Optional branch-specific academic year (NULL for tenant-wide)';
COMMENT ON COLUMN academic_years.name IS 'Name of the academic year (e.g., 2025-26)';
COMMENT ON COLUMN academic_years.start_date IS 'Start date of the academic year';
COMMENT ON COLUMN academic_years.end_date IS 'End date of the academic year';
COMMENT ON COLUMN academic_years.is_current IS 'Whether this is the current active academic year';
COMMENT ON COLUMN academic_years.is_active IS 'Whether the academic year is active for operations';
COMMENT ON COLUMN academic_years.created_by IS 'User who created this record';
COMMENT ON COLUMN academic_years.updated_by IS 'User who last updated this record';

COMMENT ON TABLE academic_terms IS 'Terms/semesters within an academic year';
COMMENT ON COLUMN academic_terms.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN academic_terms.tenant_id IS 'Reference to the parent tenant';
COMMENT ON COLUMN academic_terms.academic_year_id IS 'Reference to the parent academic year';
COMMENT ON COLUMN academic_terms.name IS 'Name of the term (e.g., Term 1, Semester 1)';
COMMENT ON COLUMN academic_terms.start_date IS 'Start date of the term';
COMMENT ON COLUMN academic_terms.end_date IS 'End date of the term';
COMMENT ON COLUMN academic_terms.sequence IS 'Order of the term within the academic year';

COMMENT ON TABLE holidays IS 'Holidays within an academic year';
COMMENT ON COLUMN holidays.id IS 'Unique identifier using UUID v7 (time-ordered)';
COMMENT ON COLUMN holidays.tenant_id IS 'Reference to the parent tenant';
COMMENT ON COLUMN holidays.academic_year_id IS 'Reference to the parent academic year';
COMMENT ON COLUMN holidays.branch_id IS 'Optional branch-specific holiday (NULL for all branches)';
COMMENT ON COLUMN holidays.name IS 'Name of the holiday';
COMMENT ON COLUMN holidays.date IS 'Date of the holiday';
COMMENT ON COLUMN holidays.type IS 'Type of holiday: public, religious, national, school, other';
COMMENT ON COLUMN holidays.is_optional IS 'Whether the holiday is optional';

-- Insert academic years permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('academic-years:read', 'View Academic Years', 'academic-years', 'View academic years, terms, and holidays'),
    ('academic-years:create', 'Create Academic Years', 'academic-years', 'Create new academic years'),
    ('academic-years:update', 'Update Academic Years', 'academic-years', 'Update academic years, terms, and holidays'),
    ('academic-years:delete', 'Delete Academic Years', 'academic-years', 'Delete academic years');

-- Assign academic-years permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
AND p.module = 'academic-years';

-- Assign academic-years permissions to admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
AND p.module = 'academic-years';

-- Assign read permission to principal
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'principal'
AND p.code IN ('academic-years:read', 'academic-years:create', 'academic-years:update');

-- Assign read permission to teacher
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher'
AND p.code = 'academic-years:read';
