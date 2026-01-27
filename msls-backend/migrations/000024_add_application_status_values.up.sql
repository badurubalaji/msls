-- Migration: 000024_add_application_status_values.up.sql
-- Description: Add missing application status enum values

-- Add new enum values to application_status type
ALTER TYPE application_status ADD VALUE IF NOT EXISTS 'documents_verified';
ALTER TYPE application_status ADD VALUE IF NOT EXISTS 'interview_scheduled';
ALTER TYPE application_status ADD VALUE IF NOT EXISTS 'interview_completed';
ALTER TYPE application_status ADD VALUE IF NOT EXISTS 'withdrawn';
