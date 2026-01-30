-- Migration: 000050_timetables.down.sql
-- Description: Drop timetables and timetable_entries tables

DROP TRIGGER IF EXISTS trigger_timetable_entries_updated_at ON timetable_entries;
DROP TRIGGER IF EXISTS trigger_timetables_updated_at ON timetables;

DROP TABLE IF EXISTS timetable_entries;
DROP TABLE IF EXISTS timetables;

DROP TYPE IF EXISTS timetable_status;

DELETE FROM permissions WHERE code IN (
    'timetables:read',
    'timetables:create',
    'timetables:update',
    'timetables:delete',
    'timetables:publish'
);
