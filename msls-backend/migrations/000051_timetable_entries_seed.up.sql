-- Migration: 000051_timetable_entries_seed.up.sql
-- Description: Seed sample timetables with entries for testing

-- Create sample timetables for Class 5, 6, 7 Section A
INSERT INTO timetables (tenant_id, branch_id, section_id, academic_year_id, name, status, created_at, updated_at)
SELECT
    s.tenant_id,
    c.branch_id,
    s.id as section_id,
    ay.id as academic_year_id,
    c.name || ' - ' || s.name || ' Timetable',
    'published',
    NOW(),
    NOW()
FROM sections s
JOIN classes c ON c.id = s.class_id
JOIN academic_years ay ON ay.tenant_id = s.tenant_id AND ay.is_current = true
WHERE c.code IN ('V', 'VI', 'VII') AND s.code = 'A'
ON CONFLICT DO NOTHING;

-- Create timetable entries for Class 5 Section A
-- Monday (day_of_week = 1) - Assembly Day pattern
INSERT INTO timetable_entries (tenant_id, timetable_id, day_of_week, period_slot_id, subject_id, staff_id, is_free_period, created_at, updated_at)
SELECT
    tt.tenant_id,
    tt.id as timetable_id,
    1 as day_of_week,
    ps.id as period_slot_id,
    CASE ps.period_number
        WHEN 1 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'ENG' LIMIT 1)
        WHEN 2 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'MATH' LIMIT 1)
        WHEN 3 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'HIN' LIMIT 1)
        WHEN 4 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'EVS' LIMIT 1)
        WHEN 5 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SCI' LIMIT 1)
        WHEN 6 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SST' LIMIT 1)
        WHEN 7 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'ART' LIMIT 1)
    END as subject_id,
    (SELECT id FROM staff WHERE tenant_id = tt.tenant_id AND staff_type = 'teaching' ORDER BY RANDOM() LIMIT 1) as staff_id,
    false as is_free_period,
    NOW(),
    NOW()
FROM timetables tt
JOIN sections s ON s.id = tt.section_id
JOIN classes c ON c.id = s.class_id
JOIN branches b ON b.id = tt.branch_id
JOIN day_patterns dp ON dp.tenant_id = tt.tenant_id AND dp.code = 'ASSM'
JOIN period_slots ps ON ps.day_pattern_id = dp.id AND ps.branch_id = b.id AND ps.slot_type = 'regular'
WHERE c.code = 'V' AND s.code = 'A' AND tt.status = 'published'
ON CONFLICT DO NOTHING;

-- Tuesday to Friday (day_of_week = 2-5) - Regular Day pattern for Class 5
INSERT INTO timetable_entries (tenant_id, timetable_id, day_of_week, period_slot_id, subject_id, staff_id, is_free_period, created_at, updated_at)
SELECT
    tt.tenant_id,
    tt.id as timetable_id,
    d.day_of_week,
    ps.id as period_slot_id,
    CASE
        WHEN ps.period_number = 1 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'ENG' LIMIT 1)
        WHEN ps.period_number = 2 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'MATH' LIMIT 1)
        WHEN ps.period_number = 3 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'HIN' LIMIT 1)
        WHEN ps.period_number = 4 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'EVS' LIMIT 1)
        WHEN ps.period_number = 5 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SCI' LIMIT 1)
        WHEN ps.period_number = 6 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SST' LIMIT 1)
        WHEN ps.period_number = 7 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'GK' LIMIT 1)
        WHEN ps.period_number = 8 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'PE' LIMIT 1)
    END as subject_id,
    (SELECT id FROM staff WHERE tenant_id = tt.tenant_id AND staff_type = 'teaching' ORDER BY RANDOM() LIMIT 1) as staff_id,
    false as is_free_period,
    NOW(),
    NOW()
FROM timetables tt
JOIN sections s ON s.id = tt.section_id
JOIN classes c ON c.id = s.class_id
JOIN branches b ON b.id = tt.branch_id
JOIN day_patterns dp ON dp.tenant_id = tt.tenant_id AND dp.code = 'REG'
JOIN period_slots ps ON ps.day_pattern_id = dp.id AND ps.branch_id = b.id AND ps.slot_type = 'regular'
CROSS JOIN (VALUES (2), (3), (4), (5)) AS d(day_of_week)
WHERE c.code = 'V' AND s.code = 'A' AND tt.status = 'published'
ON CONFLICT DO NOTHING;

-- Saturday (day_of_week = 6) - Half Day pattern for Class 5
INSERT INTO timetable_entries (tenant_id, timetable_id, day_of_week, period_slot_id, subject_id, staff_id, is_free_period, created_at, updated_at)
SELECT
    tt.tenant_id,
    tt.id as timetable_id,
    6 as day_of_week,
    ps.id as period_slot_id,
    CASE ps.period_number
        WHEN 1 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'ENG' LIMIT 1)
        WHEN 2 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'MATH' LIMIT 1)
        WHEN 3 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'MOR' LIMIT 1)
        WHEN 4 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'ART' LIMIT 1)
    END as subject_id,
    (SELECT id FROM staff WHERE tenant_id = tt.tenant_id AND staff_type = 'teaching' ORDER BY RANDOM() LIMIT 1) as staff_id,
    false as is_free_period,
    NOW(),
    NOW()
FROM timetables tt
JOIN sections s ON s.id = tt.section_id
JOIN classes c ON c.id = s.class_id
JOIN branches b ON b.id = tt.branch_id
JOIN day_patterns dp ON dp.tenant_id = tt.tenant_id AND dp.code = 'HALF'
JOIN period_slots ps ON ps.day_pattern_id = dp.id AND ps.branch_id = b.id AND ps.slot_type = 'regular'
WHERE c.code = 'V' AND s.code = 'A' AND tt.status = 'published'
ON CONFLICT DO NOTHING;

-- Similar entries for Class 6 Section A (Middle school - Science focused)
INSERT INTO timetable_entries (tenant_id, timetable_id, day_of_week, period_slot_id, subject_id, staff_id, is_free_period, created_at, updated_at)
SELECT
    tt.tenant_id,
    tt.id as timetable_id,
    d.day_of_week,
    ps.id as period_slot_id,
    CASE
        WHEN ps.period_number = 1 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'ENG' LIMIT 1)
        WHEN ps.period_number = 2 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'MATH' LIMIT 1)
        WHEN ps.period_number = 3 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SCI' LIMIT 1)
        WHEN ps.period_number = 4 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'HIN' LIMIT 1)
        WHEN ps.period_number = 5 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SST' LIMIT 1)
        WHEN ps.period_number = 6 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'CS' LIMIT 1)
        WHEN ps.period_number = 7 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SANS' LIMIT 1)
        WHEN ps.period_number = 8 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'PE' LIMIT 1)
    END as subject_id,
    (SELECT id FROM staff WHERE tenant_id = tt.tenant_id AND staff_type = 'teaching' ORDER BY RANDOM() LIMIT 1) as staff_id,
    false as is_free_period,
    NOW(),
    NOW()
FROM timetables tt
JOIN sections s ON s.id = tt.section_id
JOIN classes c ON c.id = s.class_id
JOIN branches b ON b.id = tt.branch_id
JOIN day_patterns dp ON dp.tenant_id = tt.tenant_id AND dp.code = 'REG'
JOIN period_slots ps ON ps.day_pattern_id = dp.id AND ps.branch_id = b.id AND ps.slot_type = 'regular'
CROSS JOIN (VALUES (1), (2), (3), (4), (5)) AS d(day_of_week)
WHERE c.code = 'VI' AND s.code = 'A' AND tt.status = 'published'
ON CONFLICT DO NOTHING;

-- Class 7 Section A entries
INSERT INTO timetable_entries (tenant_id, timetable_id, day_of_week, period_slot_id, subject_id, staff_id, is_free_period, created_at, updated_at)
SELECT
    tt.tenant_id,
    tt.id as timetable_id,
    d.day_of_week,
    ps.id as period_slot_id,
    CASE
        WHEN ps.period_number = 1 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'MATH' LIMIT 1)
        WHEN ps.period_number = 2 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'ENG' LIMIT 1)
        WHEN ps.period_number = 3 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SCI' LIMIT 1)
        WHEN ps.period_number = 4 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'SST' LIMIT 1)
        WHEN ps.period_number = 5 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'HIN' LIMIT 1)
        WHEN ps.period_number = 6 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'CS' LIMIT 1)
        WHEN ps.period_number = 7 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'ART' LIMIT 1)
        WHEN ps.period_number = 8 THEN (SELECT id FROM subjects WHERE tenant_id = tt.tenant_id AND code = 'MUS' LIMIT 1)
    END as subject_id,
    (SELECT id FROM staff WHERE tenant_id = tt.tenant_id AND staff_type = 'teaching' ORDER BY RANDOM() LIMIT 1) as staff_id,
    false as is_free_period,
    NOW(),
    NOW()
FROM timetables tt
JOIN sections s ON s.id = tt.section_id
JOIN classes c ON c.id = s.class_id
JOIN branches b ON b.id = tt.branch_id
JOIN day_patterns dp ON dp.tenant_id = tt.tenant_id AND dp.code = 'REG'
JOIN period_slots ps ON ps.day_pattern_id = dp.id AND ps.branch_id = b.id AND ps.slot_type = 'regular'
CROSS JOIN (VALUES (1), (2), (3), (4), (5)) AS d(day_of_week)
WHERE c.code = 'VII' AND s.code = 'A' AND tt.status = 'published'
ON CONFLICT DO NOTHING;

-- Update published_at for all published timetables
UPDATE timetables SET published_at = NOW() WHERE status = 'published' AND published_at IS NULL;
