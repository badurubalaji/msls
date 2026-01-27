-- Migration: 000032_fix_updated_at_columns.up.sql
-- Adds missing updated_at columns to tables that use BaseModel in GORM

-- Add updated_at to user_roles
ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Add updated_at to role_permissions
ALTER TABLE role_permissions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Add updated_at to login_attempts
ALTER TABLE login_attempts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- Add updated_at to refresh_tokens
ALTER TABLE refresh_tokens ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();
