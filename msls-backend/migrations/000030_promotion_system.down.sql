-- Migration: 000030_promotion_system.down.sql
-- Description: Rollback promotion rules, batches, and records tables

-- Remove promotion permissions from roles
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'promotion:read',
        'promotion:create',
        'promotion:update',
        'promotion:process',
        'promotion:cancel',
        'promotion_rules:read',
        'promotion_rules:manage'
    )
);

-- Remove promotion permissions
DELETE FROM permissions WHERE code IN (
    'promotion:read',
    'promotion:create',
    'promotion:update',
    'promotion:process',
    'promotion:cancel',
    'promotion_rules:read',
    'promotion_rules:manage'
);

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_promotion_records_updated_at ON promotion_records;
DROP TRIGGER IF EXISTS trigger_promotion_batches_updated_at ON promotion_batches;
DROP TRIGGER IF EXISTS trigger_promotion_rules_updated_at ON promotion_rules;

-- Drop functions
DROP FUNCTION IF EXISTS update_promotion_records_updated_at();
DROP FUNCTION IF EXISTS update_promotion_batches_updated_at();
DROP FUNCTION IF EXISTS update_promotion_rules_updated_at();

-- Drop tables (in reverse order of creation due to foreign keys)
DROP TABLE IF EXISTS promotion_records CASCADE;
DROP TABLE IF EXISTS promotion_batches CASCADE;
DROP TABLE IF EXISTS promotion_rules CASCADE;

-- Drop enum type
DROP TYPE IF EXISTS promotion_decision;
