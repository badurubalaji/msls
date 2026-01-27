-- Migration: 000032_fix_updated_at_columns.down.sql
-- Removes updated_at columns added in the up migration

ALTER TABLE user_roles DROP COLUMN IF EXISTS updated_at;
ALTER TABLE role_permissions DROP COLUMN IF EXISTS updated_at;
ALTER TABLE login_attempts DROP COLUMN IF EXISTS updated_at;
ALTER TABLE refresh_tokens DROP COLUMN IF EXISTS updated_at;
