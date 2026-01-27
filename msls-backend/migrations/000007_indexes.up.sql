-- Migration: 000007_indexes.up.sql
-- Description: Create performance indexes for common query patterns

-- Branches table indexes
CREATE INDEX idx_branches_tenant_id ON branches(tenant_id);
CREATE INDEX idx_branches_status ON branches(status);
CREATE INDEX idx_branches_is_primary ON branches(tenant_id, is_primary) WHERE is_primary = true;

-- Users table indexes
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_tenant_email ON users(tenant_id, email) WHERE email IS NOT NULL;
CREATE INDEX idx_users_tenant_phone ON users(tenant_id, phone) WHERE phone IS NOT NULL;

-- Add comments
COMMENT ON INDEX idx_branches_tenant_id IS 'Speeds up tenant-based branch queries';
COMMENT ON INDEX idx_users_tenant_id IS 'Speeds up tenant-based user queries';
COMMENT ON INDEX idx_users_email IS 'Speeds up email lookup for authentication';
COMMENT ON INDEX idx_users_phone IS 'Speeds up phone lookup for authentication';
