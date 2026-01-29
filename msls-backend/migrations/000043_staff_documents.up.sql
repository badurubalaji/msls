-- Migration: 000043_staff_documents.up.sql
-- Story 5.8: Staff Document Management

-- Document Types Configuration
CREATE TABLE IF NOT EXISTS staff_document_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL, -- identity, education, employment, compliance, other
    description TEXT,

    is_mandatory BOOLEAN NOT NULL DEFAULT false,
    has_expiry BOOLEAN NOT NULL DEFAULT false,
    default_validity_months INTEGER, -- Default validity period if has_expiry

    applicable_to VARCHAR(20)[] DEFAULT ARRAY['teaching', 'non_teaching'], -- Staff types
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_document_type_code UNIQUE (tenant_id, code)
);

-- Enable RLS
ALTER TABLE staff_document_types ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_document_types ON staff_document_types
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes for document types
CREATE INDEX idx_document_types_tenant ON staff_document_types(tenant_id);
CREATE INDEX idx_document_types_active ON staff_document_types(tenant_id, is_active) WHERE is_active = true;

-- Staff Documents
CREATE TABLE IF NOT EXISTS staff_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    document_type_id UUID NOT NULL REFERENCES staff_document_types(id),

    document_number VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,

    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL, -- S3/MinIO path
    file_size INTEGER NOT NULL, -- in bytes
    mime_type VARCHAR(100) NOT NULL,

    -- Verification
    verification_status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (verification_status IN ('pending', 'verified', 'rejected')),
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    verification_notes TEXT,
    rejection_reason TEXT,

    -- Metadata
    remarks TEXT,
    is_current BOOLEAN NOT NULL DEFAULT true, -- For document versions

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Enable RLS
ALTER TABLE staff_documents ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_staff_documents ON staff_documents
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes for staff documents
CREATE INDEX idx_staff_documents_tenant ON staff_documents(tenant_id);
CREATE INDEX idx_staff_documents_staff ON staff_documents(staff_id);
CREATE INDEX idx_staff_documents_type ON staff_documents(document_type_id);
CREATE INDEX idx_staff_documents_expiry ON staff_documents(expiry_date) WHERE expiry_date IS NOT NULL;
CREATE INDEX idx_staff_documents_status ON staff_documents(verification_status);
CREATE INDEX idx_staff_documents_current ON staff_documents(staff_id, document_type_id, is_current) WHERE is_current = true;

-- Document Expiry Notifications (audit trail)
CREATE TABLE IF NOT EXISTS staff_document_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    document_id UUID NOT NULL REFERENCES staff_documents(id) ON DELETE CASCADE,

    notification_type VARCHAR(50) NOT NULL, -- expiry_30_days, expiry_7_days, expired
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_to UUID[] NOT NULL, -- User IDs notified

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE staff_document_notifications ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_document_notifications ON staff_document_notifications
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

CREATE INDEX idx_document_notifications_tenant ON staff_document_notifications(tenant_id);
CREATE INDEX idx_document_notifications_document ON staff_document_notifications(document_id);

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('staff_document.view', 'View Staff Documents', 'Permission to view staff documents', 'staff_document', NOW(), NOW()),
    ('staff_document.upload', 'Upload Staff Documents', 'Permission to upload documents', 'staff_document', NOW(), NOW()),
    ('staff_document.verify', 'Verify Staff Documents', 'Permission to verify/reject documents', 'staff_document', NOW(), NOW()),
    ('staff_document.delete', 'Delete Staff Documents', 'Permission to delete documents', 'staff_document', NOW(), NOW()),
    ('staff_document.download', 'Download Staff Documents', 'Permission to download documents', 'staff_document', NOW(), NOW()),
    ('staff_document.manage_types', 'Manage Document Types', 'Permission to manage document type configuration', 'staff_document', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;
