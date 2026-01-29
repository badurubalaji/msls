-- Migration: 000049_timetable_seed_data.up.sql
-- Description: Seed data for timetable structure (shifts, day patterns, period slots)

-- Create default shifts for each branch
INSERT INTO shifts (tenant_id, branch_id, name, code, start_time, end_time, description, display_order, created_at, updated_at)
SELECT
    b.tenant_id,
    b.id as branch_id,
    s.name,
    s.code,
    s.start_time::TIME,
    s.end_time::TIME,
    s.description,
    s.display_order,
    NOW(),
    NOW()
FROM branches b
CROSS JOIN (
    VALUES
        ('Morning Shift', 'MORN', '08:00', '14:00', 'Regular morning shift', 1),
        ('Afternoon Shift', 'AFT', '12:30', '18:00', 'Afternoon shift for double shift schools', 2)
) AS s(name, code, start_time, end_time, description, display_order)
ON CONFLICT (tenant_id, branch_id, code) DO NOTHING;

-- Create default day patterns for each tenant
INSERT INTO day_patterns (tenant_id, name, code, description, total_periods, display_order, created_at, updated_at)
SELECT
    t.id as tenant_id,
    p.name,
    p.code,
    p.description,
    p.total_periods,
    p.display_order,
    NOW(),
    NOW()
FROM tenants t
CROSS JOIN (
    VALUES
        ('Regular Day', 'REG', 'Standard school day with 8 periods', 8, 1),
        ('Half Day', 'HALF', 'Half day with 4 periods (Saturday/special days)', 4, 2),
        ('Assembly Day', 'ASSM', 'Day with extended assembly (Monday)', 7, 3)
) AS p(name, code, description, total_periods, display_order)
ON CONFLICT (tenant_id, code) DO NOTHING;

-- Create day pattern assignments (Monday-Friday = Regular, Saturday = Half Day, Sunday = Off)
INSERT INTO day_pattern_assignments (tenant_id, branch_id, day_of_week, day_pattern_id, is_working_day, created_at, updated_at)
SELECT
    b.tenant_id,
    b.id as branch_id,
    d.day_of_week,
    CASE
        WHEN d.day_of_week = 1 THEN (SELECT id FROM day_patterns WHERE tenant_id = b.tenant_id AND code = 'ASSM')
        WHEN d.day_of_week IN (2, 3, 4, 5) THEN (SELECT id FROM day_patterns WHERE tenant_id = b.tenant_id AND code = 'REG')
        WHEN d.day_of_week = 6 THEN (SELECT id FROM day_patterns WHERE tenant_id = b.tenant_id AND code = 'HALF')
        ELSE NULL
    END as day_pattern_id,
    d.is_working_day,
    NOW(),
    NOW()
FROM branches b
CROSS JOIN (
    VALUES
        (0, false),  -- Sunday - Off
        (1, true),   -- Monday - Working (Assembly Day)
        (2, true),   -- Tuesday - Working
        (3, true),   -- Wednesday - Working
        (4, true),   -- Thursday - Working
        (5, true),   -- Friday - Working
        (6, true)    -- Saturday - Half Day
) AS d(day_of_week, is_working_day)
ON CONFLICT (tenant_id, branch_id, day_of_week) DO NOTHING;

-- Create period slots for Regular Day pattern
INSERT INTO period_slots (tenant_id, branch_id, name, period_number, slot_type, start_time, end_time, duration_minutes, day_pattern_id, shift_id, display_order, created_at, updated_at)
SELECT
    dp.tenant_id,
    b.id as branch_id,
    ps.name,
    ps.period_number,
    ps.slot_type::period_slot_type,
    ps.start_time::TIME,
    ps.end_time::TIME,
    ps.duration_minutes,
    dp.id as day_pattern_id,
    (SELECT id FROM shifts WHERE branch_id = b.id AND code = 'MORN' LIMIT 1) as shift_id,
    ps.display_order,
    NOW(),
    NOW()
FROM day_patterns dp
JOIN branches b ON b.tenant_id = dp.tenant_id
CROSS JOIN (
    VALUES
        ('Assembly', NULL, 'assembly', '08:00', '08:30', 30, 1),
        ('Period 1', 1, 'regular', '08:30', '09:15', 45, 2),
        ('Period 2', 2, 'regular', '09:15', '10:00', 45, 3),
        ('Short Break', NULL, 'break', '10:00', '10:15', 15, 4),
        ('Period 3', 3, 'regular', '10:15', '11:00', 45, 5),
        ('Period 4', 4, 'regular', '11:00', '11:45', 45, 6),
        ('Lunch Break', NULL, 'lunch', '11:45', '12:30', 45, 7),
        ('Period 5', 5, 'regular', '12:30', '13:15', 45, 8),
        ('Period 6', 6, 'regular', '13:15', '14:00', 45, 9),
        ('Short Break', NULL, 'break', '14:00', '14:10', 10, 10),
        ('Period 7', 7, 'regular', '14:10', '14:55', 45, 11),
        ('Period 8', 8, 'regular', '14:55', '15:40', 45, 12)
) AS ps(name, period_number, slot_type, start_time, end_time, duration_minutes, display_order)
WHERE dp.code = 'REG';

-- Create period slots for Half Day pattern (Saturday)
INSERT INTO period_slots (tenant_id, branch_id, name, period_number, slot_type, start_time, end_time, duration_minutes, day_pattern_id, shift_id, display_order, created_at, updated_at)
SELECT
    dp.tenant_id,
    b.id as branch_id,
    ps.name,
    ps.period_number,
    ps.slot_type::period_slot_type,
    ps.start_time::TIME,
    ps.end_time::TIME,
    ps.duration_minutes,
    dp.id as day_pattern_id,
    (SELECT id FROM shifts WHERE branch_id = b.id AND code = 'MORN' LIMIT 1) as shift_id,
    ps.display_order,
    NOW(),
    NOW()
FROM day_patterns dp
JOIN branches b ON b.tenant_id = dp.tenant_id
CROSS JOIN (
    VALUES
        ('Assembly', NULL, 'assembly', '08:00', '08:20', 20, 1),
        ('Period 1', 1, 'regular', '08:20', '09:00', 40, 2),
        ('Period 2', 2, 'regular', '09:00', '09:40', 40, 3),
        ('Short Break', NULL, 'break', '09:40', '09:50', 10, 4),
        ('Period 3', 3, 'regular', '09:50', '10:30', 40, 5),
        ('Period 4', 4, 'regular', '10:30', '11:10', 40, 6)
) AS ps(name, period_number, slot_type, start_time, end_time, duration_minutes, display_order)
WHERE dp.code = 'HALF';

-- Create period slots for Assembly Day (Monday - one less period for extended assembly)
INSERT INTO period_slots (tenant_id, branch_id, name, period_number, slot_type, start_time, end_time, duration_minutes, day_pattern_id, shift_id, display_order, created_at, updated_at)
SELECT
    dp.tenant_id,
    b.id as branch_id,
    ps.name,
    ps.period_number,
    ps.slot_type::period_slot_type,
    ps.start_time::TIME,
    ps.end_time::TIME,
    ps.duration_minutes,
    dp.id as day_pattern_id,
    (SELECT id FROM shifts WHERE branch_id = b.id AND code = 'MORN' LIMIT 1) as shift_id,
    ps.display_order,
    NOW(),
    NOW()
FROM day_patterns dp
JOIN branches b ON b.tenant_id = dp.tenant_id
CROSS JOIN (
    VALUES
        ('Morning Assembly', NULL, 'assembly', '08:00', '08:45', 45, 1),
        ('Period 1', 1, 'regular', '08:45', '09:30', 45, 2),
        ('Period 2', 2, 'regular', '09:30', '10:15', 45, 3),
        ('Short Break', NULL, 'break', '10:15', '10:30', 15, 4),
        ('Period 3', 3, 'regular', '10:30', '11:15', 45, 5),
        ('Period 4', 4, 'regular', '11:15', '12:00', 45, 6),
        ('Lunch Break', NULL, 'lunch', '12:00', '12:45', 45, 7),
        ('Period 5', 5, 'regular', '12:45', '13:30', 45, 8),
        ('Period 6', 6, 'regular', '13:30', '14:15', 45, 9),
        ('Period 7', 7, 'regular', '14:15', '15:00', 45, 10)
) AS ps(name, period_number, slot_type, start_time, end_time, duration_minutes, display_order)
WHERE dp.code = 'ASSM';
