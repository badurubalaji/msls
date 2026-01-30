-- Migration: 000054_period_attendance.up.sql
-- Description: Add period-wise attendance support (Story 7.2)
-- Extends student_attendance table to support period-specific attendance tracking

-- Step 1: Add period_id and timetable_entry_id columns to student_attendance
ALTER TABLE student_attendance
    ADD COLUMN period_id UUID REFERENCES period_slots(id),
    ADD COLUMN timetable_entry_id UUID REFERENCES timetable_entries(id);

-- Step 2: Drop the old unique constraint (only allowed one record per student per day)
ALTER TABLE student_attendance
    DROP CONSTRAINT IF EXISTS uniq_student_daily_attendance;

-- Step 3: Add new unique constraint that allows multiple records per day (one per period)
-- Using COALESCE to handle NULL period_id (daily attendance) and specific period_id (period-wise)
-- NULL period_id records are treated as a special case (uuid_nil) for uniqueness
CREATE UNIQUE INDEX uniq_student_period_attendance
    ON student_attendance(tenant_id, student_id, attendance_date, COALESCE(period_id, '00000000-0000-0000-0000-000000000000'::UUID));

-- Step 4: Add index for period-wise attendance queries
CREATE INDEX idx_student_attendance_period
    ON student_attendance(section_id, attendance_date, period_id)
    WHERE period_id IS NOT NULL;

-- Step 5: Add index for timetable entry lookups
CREATE INDEX idx_student_attendance_timetable_entry
    ON student_attendance(timetable_entry_id)
    WHERE timetable_entry_id IS NOT NULL;

-- Step 6: Add period_attendance_enabled setting to student_attendance_settings
ALTER TABLE student_attendance_settings
    ADD COLUMN period_attendance_enabled BOOLEAN NOT NULL DEFAULT false;

-- Step 7: Add comment for documentation
COMMENT ON COLUMN student_attendance.period_id IS 'Reference to period_slots for period-wise attendance. NULL for daily attendance (Story 7.1 compatibility).';
COMMENT ON COLUMN student_attendance.timetable_entry_id IS 'Reference to timetable_entries for subject context. Links attendance to specific subject/teacher.';
COMMENT ON COLUMN student_attendance_settings.period_attendance_enabled IS 'When true, enables period-wise attendance marking for the branch.';
