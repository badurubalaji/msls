-- Migration: 000031_bulk_operations.up.sql
-- Creates bulk operations and bulk operation items tables for batch processing

-- Bulk operations table
CREATE TABLE IF NOT EXISTS bulk_operations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    operation_type VARCHAR(50) NOT NULL, -- send_sms, send_email, update_status, export
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    total_count INTEGER NOT NULL DEFAULT 0,
    processed_count INTEGER NOT NULL DEFAULT 0,
    success_count INTEGER NOT NULL DEFAULT 0,
    failure_count INTEGER NOT NULL DEFAULT 0,
    parameters JSONB, -- Operation-specific params
    result_url VARCHAR(500), -- For export operations
    error_message TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id)
);

-- Enable RLS
ALTER TABLE bulk_operations ENABLE ROW LEVEL SECURITY;

-- Create RLS policy
CREATE POLICY tenant_isolation ON bulk_operations
    USING (tenant_id = current_setting('app.current_tenant', true)::UUID);

-- Indexes for bulk_operations
CREATE INDEX idx_bulk_operations_tenant ON bulk_operations(tenant_id);
CREATE INDEX idx_bulk_operations_user ON bulk_operations(created_by);
CREATE INDEX idx_bulk_operations_status ON bulk_operations(status) WHERE status IN ('pending', 'processing');
CREATE INDEX idx_bulk_operations_type ON bulk_operations(operation_type);

-- Bulk operation items table
CREATE TABLE IF NOT EXISTS bulk_operation_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    operation_id UUID NOT NULL REFERENCES bulk_operations(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'skipped')),
    error_message TEXT,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE bulk_operation_items ENABLE ROW LEVEL SECURITY;

-- Create RLS policy
CREATE POLICY tenant_isolation ON bulk_operation_items
    USING (tenant_id = current_setting('app.current_tenant', true)::UUID);

-- Indexes for bulk_operation_items
CREATE INDEX idx_bulk_items_operation ON bulk_operation_items(operation_id);
CREATE INDEX idx_bulk_items_student ON bulk_operation_items(student_id);
CREATE INDEX idx_bulk_items_status ON bulk_operation_items(operation_id, status);

-- Add comment
COMMENT ON TABLE bulk_operations IS 'Tracks bulk operations (SMS, email, status update, export) for audit and progress tracking';
COMMENT ON TABLE bulk_operation_items IS 'Individual items in a bulk operation with their processing status';
