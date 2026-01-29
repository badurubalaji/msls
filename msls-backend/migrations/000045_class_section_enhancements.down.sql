-- Migration: 000045_class_section_enhancements.down.sql
-- Description: Revert class section enhancements

-- Remove columns from sections
ALTER TABLE sections DROP COLUMN IF EXISTS class_teacher_id;
ALTER TABLE sections DROP COLUMN IF EXISTS academic_year_id;
ALTER TABLE sections DROP COLUMN IF EXISTS stream_id;
ALTER TABLE sections DROP COLUMN IF EXISTS room_number;

-- Remove column from classes
ALTER TABLE classes DROP COLUMN IF EXISTS has_streams;

-- Drop class_streams table
DROP TABLE IF EXISTS class_streams;

-- Drop streams table
DROP TABLE IF EXISTS streams;

-- Remove permissions
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'stream:view', 'stream:create', 'stream:update', 'stream:delete'
    )
);

DELETE FROM permissions WHERE code IN (
    'stream:view', 'stream:create', 'stream:update', 'stream:delete'
);
