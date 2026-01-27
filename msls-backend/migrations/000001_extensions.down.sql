-- Migration: 000001_extensions.down.sql
-- Description: Remove PostgreSQL extensions

-- Note: Dropping extensions may fail if objects depend on them
-- This is intentional to prevent accidental data loss
DROP EXTENSION IF EXISTS "pgcrypto";
