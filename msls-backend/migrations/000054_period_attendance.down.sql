-- Migration: 000054_period_attendance.down.sql
-- Description: Rollback period-wise attendance support (Story 7.2)

-- Step 1: Remove period_attendance_enabled from settings
ALTER TABLE student_attendance_settings
    DROP COLUMN IF EXISTS period_attendance_enabled;

-- Step 2: Remove indexes
DROP INDEX IF EXISTS idx_student_attendance_timetable_entry;
DROP INDEX IF EXISTS idx_student_attendance_period;

-- Step 3: Remove the new unique constraint
DROP INDEX IF EXISTS uniq_student_period_attendance;

-- Step 4: Restore the original unique constraint (only one record per student per day)
-- Note: This may fail if there are multiple records per student per day
-- In that case, manual cleanup is required before running this migration
ALTER TABLE student_attendance
    ADD CONSTRAINT uniq_student_daily_attendance
        UNIQUE (tenant_id, student_id, attendance_date);

-- Step 5: Remove the columns
ALTER TABLE student_attendance
    DROP COLUMN IF EXISTS timetable_entry_id,
    DROP COLUMN IF EXISTS period_id;

-- Step 6: Remove comments (no action needed - columns are dropped)
