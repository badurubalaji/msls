-- Reverse Student Attendance Tables

DROP TRIGGER IF EXISTS set_updated_at_student_attendance_settings ON student_attendance_settings;
DROP TRIGGER IF EXISTS set_updated_at_student_attendance ON student_attendance;

DROP FUNCTION IF EXISTS update_student_attendance_settings_updated_at();
DROP FUNCTION IF EXISTS update_student_attendance_updated_at();

DROP POLICY IF EXISTS tenant_isolation_student_attendance_settings ON student_attendance_settings;
DROP POLICY IF EXISTS tenant_isolation_student_attendance ON student_attendance;

DROP TABLE IF EXISTS student_attendance_settings;
DROP TABLE IF EXISTS student_attendance;

DELETE FROM permissions WHERE module = 'student_attendance';
