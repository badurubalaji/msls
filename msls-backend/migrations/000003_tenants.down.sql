-- Migration: 000003_tenants.down.sql
-- Description: Remove tenants table

DROP TABLE IF EXISTS tenants CASCADE;
