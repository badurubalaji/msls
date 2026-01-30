-- Migration: 000051_timetable_entries_seed.down.sql
-- Description: Remove seed timetables and entries

-- Delete timetable entries for seed timetables
DELETE FROM timetable_entries
WHERE timetable_id IN (
    SELECT tt.id FROM timetables tt
    JOIN sections s ON s.id = tt.section_id
    JOIN classes c ON c.id = s.class_id
    WHERE c.code IN ('V', 'VI', 'VII') AND s.code = 'A'
);

-- Delete seed timetables
DELETE FROM timetables
WHERE section_id IN (
    SELECT s.id FROM sections s
    JOIN classes c ON c.id = s.class_id
    WHERE c.code IN ('V', 'VI', 'VII') AND s.code = 'A'
);
