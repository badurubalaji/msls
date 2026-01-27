-- Migration: 000004_branches.down.sql
-- Description: Remove branches table and related objects

DROP TRIGGER IF EXISTS trigger_single_primary_branch ON branches;
DROP FUNCTION IF EXISTS ensure_single_primary_branch();
DROP TABLE IF EXISTS branches CASCADE;
