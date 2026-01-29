-- Migration: 000048_timetable_structure.down.sql
-- Description: Remove timetable structure tables

-- Remove shift_id from sections
ALTER TABLE sections DROP COLUMN IF EXISTS shift_id;

-- Remove permissions
DELETE FROM role_permissions
WHERE permission_id IN (SELECT id FROM permissions WHERE code IN (
    'timetable:view', 'timetable:create', 'timetable:manage', 'shift:view', 'shift:manage'
));

DELETE FROM permissions WHERE code IN (
    'timetable:view', 'timetable:create', 'timetable:manage', 'shift:view', 'shift:manage'
);

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_period_slots_updated_at ON period_slots;
DROP TRIGGER IF EXISTS trigger_day_pattern_assignments_updated_at ON day_pattern_assignments;
DROP TRIGGER IF EXISTS trigger_day_patterns_updated_at ON day_patterns;
DROP TRIGGER IF EXISTS trigger_shifts_updated_at ON shifts;

-- Drop tables
DROP TABLE IF EXISTS period_slots;
DROP TABLE IF EXISTS day_pattern_assignments;
DROP TABLE IF EXISTS day_patterns;
DROP TABLE IF EXISTS shifts;

-- Drop types
DROP TYPE IF EXISTS period_slot_type;
