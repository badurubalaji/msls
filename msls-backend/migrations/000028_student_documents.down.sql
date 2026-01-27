-- Rollback: Student Documents
-- Story: 4.5 - Document Management

-- Remove triggers
DROP TRIGGER IF EXISTS trg_student_documents_updated_at ON student_documents;
DROP TRIGGER IF EXISTS trg_document_types_updated_at ON document_types;

-- Remove functions
DROP FUNCTION IF EXISTS update_student_documents_updated_at();
DROP FUNCTION IF EXISTS update_document_types_updated_at();

-- Remove role permissions
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'document:read', 'document:create', 'document:update',
        'document:delete', 'document:verify', 'document_type:manage'
    )
);

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'document:read', 'document:create', 'document:update',
    'document:delete', 'document:verify', 'document_type:manage'
);

-- Drop tables
DROP TABLE IF EXISTS class_required_documents;
DROP TABLE IF EXISTS student_documents;
DROP TABLE IF EXISTS document_types;

-- Drop enum
DROP TYPE IF EXISTS document_status;
