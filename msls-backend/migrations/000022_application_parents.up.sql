-- Migration: 000022_application_parents.up.sql
-- Description: Create application_parents table for storing multiple parent/guardian records

-- Create parent relation type
DO $$ BEGIN
    CREATE TYPE parent_relation AS ENUM (
        'father',
        'mother',
        'guardian',
        'other'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create application_parents table
CREATE TABLE IF NOT EXISTS application_parents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    relation VARCHAR(20) NOT NULL,
    name VARCHAR(200) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    occupation VARCHAR(200),
    education VARCHAR(100),
    annual_income VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_app_parents_tenant ON application_parents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_app_parents_application ON application_parents(application_id);
CREATE INDEX IF NOT EXISTS idx_app_parents_relation ON application_parents(relation);

-- Enable RLS
ALTER TABLE application_parents ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
CREATE POLICY tenant_isolation_app_parents ON application_parents
    FOR ALL
    USING (tenant_id = COALESCE(NULLIF(current_setting('app.tenant_id', true), '')::UUID, '00000000-0000-0000-0000-000000000000'::UUID));

CREATE POLICY bypass_rls_app_parents ON application_parents
    FOR ALL
    USING (current_setting('app.bypass_rls', true) = 'true');

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_app_parents_updated_at ON application_parents;
CREATE TRIGGER trigger_app_parents_updated_at
    BEFORE UPDATE ON application_parents
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Add comments
COMMENT ON TABLE application_parents IS 'Parent/guardian records for admission applications';
COMMENT ON COLUMN application_parents.relation IS 'Relationship to student (father, mother, guardian, other)';
