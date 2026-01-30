-- Remove attendance seed data
-- Note: This will remove all attendance records - use carefully in production

DELETE FROM student_attendance WHERE remarks IN ('Sick leave', 'Traffic delay') OR remarks IS NULL;
DELETE FROM student_attendance_settings;
DELETE FROM student_enrollments WHERE enrollment_date > CURRENT_DATE - INTERVAL '1 year';
DELETE FROM students WHERE admission_number LIKE 'STU-2026-%';
