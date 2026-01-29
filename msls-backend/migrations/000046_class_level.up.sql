-- Migration: 000046_class_level.up.sql
-- Description: Add level field to classes as per PRD (nursery, primary, middle, secondary, senior_secondary)

-- Create class_level enum type
DO $$ BEGIN
    CREATE TYPE class_level AS ENUM ('nursery', 'primary', 'middle', 'secondary', 'senior_secondary');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add level column to classes table
ALTER TABLE classes ADD COLUMN IF NOT EXISTS level class_level;

-- Set default levels based on display_order for existing classes
-- This is a reasonable default mapping:
-- display_order 1-2: nursery (LKG, UKG)
-- display_order 3-7: primary (Class 1-5)
-- display_order 8-10: middle (Class 6-8)
-- display_order 11-12: secondary (Class 9-10)
-- display_order 13+: senior_secondary (Class 11-12)
UPDATE classes SET level =
    CASE
        WHEN display_order <= 2 THEN 'nursery'::class_level
        WHEN display_order <= 7 THEN 'primary'::class_level
        WHEN display_order <= 10 THEN 'middle'::class_level
        WHEN display_order <= 12 THEN 'secondary'::class_level
        ELSE 'senior_secondary'::class_level
    END
WHERE level IS NULL;

-- Create index for level filtering
CREATE INDEX IF NOT EXISTS idx_classes_level ON classes(tenant_id, level);

-- Comment on the column
COMMENT ON COLUMN classes.level IS 'Class level: nursery (LKG/UKG), primary (1-5), middle (6-8), secondary (9-10), senior_secondary (11-12)';
