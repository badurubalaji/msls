-- Migration: 000029_student_enrollments.down.sql
-- Description: Drop student_enrollments and enrollment_status_changes tables

-- Drop trigger
DROP TRIGGER IF EXISTS trigger_student_enrollments_updated_at ON student_enrollments;

-- Drop function
DROP FUNCTION IF EXISTS update_student_enrollments_updated_at();

-- Drop permissions
DELETE FROM permissions WHERE code IN (
    'enrollments:read',
    'enrollments:create',
    'enrollments:update',
    'enrollments:delete'
);

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS enrollment_status_changes;
DROP TABLE IF EXISTS student_enrollments;

-- Drop enum type
DROP TYPE IF EXISTS enrollment_status;
