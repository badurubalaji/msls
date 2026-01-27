-- Migration: 000002_uuid_v7.down.sql
-- Description: Remove UUID v7 generation function

DROP FUNCTION IF EXISTS uuid_generate_v7();
