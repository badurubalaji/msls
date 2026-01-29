-- Migration: 000043_staff_documents.down.sql
-- Story 5.8: Staff Document Management - Rollback

-- Drop policies
DROP POLICY IF EXISTS tenant_isolation_document_notifications ON staff_document_notifications;
DROP POLICY IF EXISTS tenant_isolation_staff_documents ON staff_documents;
DROP POLICY IF EXISTS tenant_isolation_document_types ON staff_document_types;

-- Drop tables in reverse order
DROP TABLE IF EXISTS staff_document_notifications;
DROP TABLE IF EXISTS staff_documents;
DROP TABLE IF EXISTS staff_document_types;

-- Remove permissions
DELETE FROM permissions WHERE code IN (
    'staff_document.view',
    'staff_document.upload',
    'staff_document.verify',
    'staff_document.delete',
    'staff_document.download',
    'staff_document.manage_types'
);
