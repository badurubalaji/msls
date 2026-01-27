-- Migration: 000016_admission_applications.down.sql
-- Description: Drop admission applications and related tables

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_test_registrations_updated_at ON test_registrations;
DROP TRIGGER IF EXISTS trigger_entrance_tests_updated_at ON entrance_tests;
DROP TRIGGER IF EXISTS trigger_app_docs_updated_at ON application_documents;
DROP TRIGGER IF EXISTS trigger_applications_updated_at ON admission_applications;
DROP TRIGGER IF EXISTS trigger_update_app_on_test_result ON test_registrations;

-- Drop functions
DROP FUNCTION IF EXISTS update_applications_updated_at();
DROP FUNCTION IF EXISTS update_application_on_test_result();
DROP FUNCTION IF EXISTS get_next_roll_number(UUID, UUID);
DROP FUNCTION IF EXISTS get_next_application_number(UUID, UUID);

-- Remove permissions
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'applications:read', 'applications:create', 'applications:update', 'applications:delete'
    )
);

DELETE FROM permissions WHERE code IN (
    'applications:read', 'applications:create', 'applications:update', 'applications:delete'
);

-- Drop tables
DROP TABLE IF EXISTS roll_number_sequence;
DROP TABLE IF EXISTS application_number_sequence;
DROP TABLE IF EXISTS application_reviews;
DROP TABLE IF EXISTS test_registrations;
DROP TABLE IF EXISTS entrance_tests;
DROP TABLE IF EXISTS application_documents;
DROP TABLE IF EXISTS admission_applications;

-- Drop types
DROP TYPE IF EXISTS review_status;
DROP TYPE IF EXISTS review_type;
DROP TYPE IF EXISTS test_result;
DROP TYPE IF EXISTS test_registration_status;
DROP TYPE IF EXISTS verification_status;
DROP TYPE IF EXISTS application_status;
