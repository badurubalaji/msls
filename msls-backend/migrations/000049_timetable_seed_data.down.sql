-- Migration: 000049_timetable_seed_data.down.sql
-- Description: Remove timetable seed data

-- Remove period slots
DELETE FROM period_slots WHERE day_pattern_id IN (
    SELECT id FROM day_patterns WHERE code IN ('REG', 'HALF', 'ASSM')
);

-- Remove day pattern assignments
DELETE FROM day_pattern_assignments;

-- Remove day patterns
DELETE FROM day_patterns WHERE code IN ('REG', 'HALF', 'ASSM');

-- Remove shifts
DELETE FROM shifts WHERE code IN ('MORN', 'AFT');
