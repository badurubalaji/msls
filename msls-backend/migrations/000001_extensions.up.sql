-- Migration: 000001_extensions.up.sql
-- Description: Enable required PostgreSQL extensions for MSLS

-- Enable pgcrypto extension for cryptographic functions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
