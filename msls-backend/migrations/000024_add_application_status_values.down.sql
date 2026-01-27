-- Migration: 000024_add_application_status_values.down.sql
-- Description: Cannot remove enum values in PostgreSQL without recreating the type
-- This is a no-op migration - enum values cannot be easily removed

-- Note: To truly rollback, you would need to:
-- 1. Create a new enum type without the values
-- 2. Update all references to use the new type
-- 3. Drop the old type
-- 4. Rename the new type

SELECT 1;
