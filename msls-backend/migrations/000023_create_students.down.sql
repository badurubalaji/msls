-- Migration: 000023_create_students.down.sql
-- Description: Drop students and student_addresses tables

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_students_updated_at ON students;
DROP TRIGGER IF EXISTS trigger_student_addresses_updated_at ON student_addresses;

-- Drop function
DROP FUNCTION IF EXISTS update_students_updated_at();

-- Remove student permissions
DELETE FROM permissions WHERE module = 'students';

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS student_addresses;
DROP TABLE IF EXISTS student_admission_sequences;
DROP TABLE IF EXISTS students;
