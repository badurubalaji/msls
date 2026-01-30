-- Migration: 000055_attendance_audit.down.sql
-- Rollback: Story 7.3 Attendance Audit Trail

-- Drop trigger first
DROP TRIGGER IF EXISTS trg_audit_student_attendance ON student_attendance;
DROP FUNCTION IF EXISTS fn_audit_student_attendance();

-- Drop audit table
DROP TABLE IF EXISTS student_attendance_audit;
