-- Migration: 000046_class_level.down.sql
-- Description: Remove level field from classes

-- Drop index
DROP INDEX IF EXISTS idx_classes_level;

-- Remove column
ALTER TABLE classes DROP COLUMN IF EXISTS level;

-- Drop enum type
DROP TYPE IF EXISTS class_level;
