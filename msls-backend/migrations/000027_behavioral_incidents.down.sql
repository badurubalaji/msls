-- Migration: 000027_behavioral_incidents.down.sql
-- Description: Drop behavioral incidents and follow-ups tables

-- Remove role permissions
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN ('behavior:read', 'behavior:write')
);

-- Remove permissions
DELETE FROM permissions WHERE code IN ('behavior:read', 'behavior:write');

-- Drop tables
DROP TABLE IF EXISTS incident_follow_ups;
DROP TABLE IF EXISTS student_behavioral_incidents;
