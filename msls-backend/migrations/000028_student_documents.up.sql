-- Migration: Student Documents
-- Story: 4.5 - Document Management

-- Document status enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_status') THEN
        CREATE TYPE document_status AS ENUM ('pending_verification', 'verified', 'rejected');
    END IF;
END$$;

-- Document types table (configurable per tenant)
CREATE TABLE IF NOT EXISTS document_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_mandatory BOOLEAN NOT NULL DEFAULT FALSE,
    has_expiry BOOLEAN NOT NULL DEFAULT FALSE,
    allowed_extensions VARCHAR(100) NOT NULL DEFAULT 'pdf,jpg,jpeg,png',
    max_size_mb INTEGER NOT NULL DEFAULT 5,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_doc_type_code UNIQUE (tenant_id, code)
);

-- Enable RLS on document_types
ALTER TABLE document_types ENABLE ROW LEVEL SECURITY;

-- RLS policy for document_types
DROP POLICY IF EXISTS tenant_isolation_document_types ON document_types;
CREATE POLICY tenant_isolation_document_types ON document_types
    USING (tenant_id = COALESCE(
        NULLIF(current_setting('app.current_tenant', true), '')::UUID,
        '00000000-0000-0000-0000-000000000000'::UUID
    ));

-- Student documents table
CREATE TABLE IF NOT EXISTS student_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    document_type_id UUID NOT NULL REFERENCES document_types(id),
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size_bytes INTEGER NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    document_number VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,
    status document_status NOT NULL DEFAULT 'pending_verification',
    rejection_reason TEXT,
    verified_at TIMESTAMPTZ,
    verified_by UUID REFERENCES users(id),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    uploaded_by UUID NOT NULL REFERENCES users(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT uniq_student_doc_type UNIQUE (student_id, document_type_id)
);

-- Enable RLS on student_documents
ALTER TABLE student_documents ENABLE ROW LEVEL SECURITY;

-- RLS policy for student_documents
DROP POLICY IF EXISTS tenant_isolation_student_documents ON student_documents;
CREATE POLICY tenant_isolation_student_documents ON student_documents
    USING (tenant_id = COALESCE(
        NULLIF(current_setting('app.current_tenant', true), '')::UUID,
        '00000000-0000-0000-0000-000000000000'::UUID
    ));

-- Indexes
CREATE INDEX IF NOT EXISTS idx_student_documents_student ON student_documents(student_id);
CREATE INDEX IF NOT EXISTS idx_student_documents_status ON student_documents(status);
CREATE INDEX IF NOT EXISTS idx_student_documents_type ON student_documents(document_type_id);
CREATE INDEX IF NOT EXISTS idx_document_types_tenant ON document_types(tenant_id);
CREATE INDEX IF NOT EXISTS idx_document_types_active ON document_types(tenant_id, is_active);

-- Required documents per class (optional configuration)
CREATE TABLE IF NOT EXISTS class_required_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    class_id UUID, -- NULL means applies to all classes
    document_type_id UUID NOT NULL REFERENCES document_types(id) ON DELETE CASCADE,
    is_required BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_class_doc UNIQUE (tenant_id, class_id, document_type_id)
);

-- Enable RLS on class_required_documents
ALTER TABLE class_required_documents ENABLE ROW LEVEL SECURITY;

-- RLS policy for class_required_documents
DROP POLICY IF EXISTS tenant_isolation_class_required_documents ON class_required_documents;
CREATE POLICY tenant_isolation_class_required_documents ON class_required_documents
    USING (tenant_id = COALESCE(
        NULLIF(current_setting('app.current_tenant', true), '')::UUID,
        '00000000-0000-0000-0000-000000000000'::UUID
    ));

-- Updated timestamp triggers
CREATE OR REPLACE FUNCTION update_document_types_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_student_documents_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_document_types_updated_at ON document_types;
CREATE TRIGGER trg_document_types_updated_at
    BEFORE UPDATE ON document_types
    FOR EACH ROW
    EXECUTE FUNCTION update_document_types_updated_at();

DROP TRIGGER IF EXISTS trg_student_documents_updated_at ON student_documents;
CREATE TRIGGER trg_student_documents_updated_at
    BEFORE UPDATE ON student_documents
    FOR EACH ROW
    EXECUTE FUNCTION update_student_documents_updated_at();

-- Add document permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('document:read', 'View Documents', 'Can view student documents', 'documents', NOW(), NOW()),
    ('document:create', 'Upload Documents', 'Can upload student documents', 'documents', NOW(), NOW()),
    ('document:update', 'Update Documents', 'Can update document metadata', 'documents', NOW(), NOW()),
    ('document:delete', 'Delete Documents', 'Can delete student documents', 'documents', NOW(), NOW()),
    ('document:verify', 'Verify Documents', 'Can verify or reject documents', 'documents', NOW(), NOW()),
    ('document_type:manage', 'Manage Document Types', 'Can create and manage document types', 'documents', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Grant document permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Admin' AND p.code IN ('document:read', 'document:create', 'document:update', 'document:delete', 'document:verify', 'document_type:manage')
ON CONFLICT DO NOTHING;

-- Grant basic document permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Staff' AND p.code IN ('document:read', 'document:create', 'document:update')
ON CONFLICT DO NOTHING;
