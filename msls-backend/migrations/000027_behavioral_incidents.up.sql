-- Migration: 000027_behavioral_incidents.up.sql
-- Description: Create behavioral incidents and follow-ups tables

-- Behavioral incidents table
CREATE TABLE student_behavioral_incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    incident_type VARCHAR(30) NOT NULL CHECK (incident_type IN ('positive_recognition', 'minor_infraction', 'major_violation')),
    severity VARCHAR(20) NOT NULL DEFAULT 'medium' CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    incident_date DATE NOT NULL,
    incident_time TIME NOT NULL,
    location VARCHAR(200),
    description TEXT NOT NULL,
    witnesses JSONB, -- Array of witness names
    student_response TEXT,
    action_taken TEXT NOT NULL,
    parent_meeting_required BOOLEAN NOT NULL DEFAULT FALSE,
    parent_notified BOOLEAN NOT NULL DEFAULT FALSE,
    parent_notified_at TIMESTAMPTZ,
    reported_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for behavioral incidents
CREATE INDEX idx_behavioral_incidents_tenant ON student_behavioral_incidents(tenant_id);
CREATE INDEX idx_behavioral_incidents_student ON student_behavioral_incidents(student_id);
CREATE INDEX idx_behavioral_incidents_date ON student_behavioral_incidents(incident_date DESC);
CREATE INDEX idx_behavioral_incidents_type ON student_behavioral_incidents(incident_type);
CREATE INDEX idx_behavioral_incidents_severity ON student_behavioral_incidents(severity);

-- RLS for behavioral incidents
ALTER TABLE student_behavioral_incidents ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_behavioral_incidents ON student_behavioral_incidents
    USING (tenant_id = current_setting('app.current_tenant', true)::UUID);

-- Incident follow-ups table
CREATE TABLE incident_follow_ups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    incident_id UUID NOT NULL REFERENCES student_behavioral_incidents(id) ON DELETE CASCADE,
    scheduled_date DATE NOT NULL,
    scheduled_time TIME,
    participants JSONB, -- Array of {name, role} objects
    expected_outcomes TEXT,
    meeting_notes TEXT,
    actual_outcomes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'cancelled')),
    completed_at TIMESTAMPTZ,
    completed_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Indexes for follow-ups
CREATE INDEX idx_incident_follow_ups_tenant ON incident_follow_ups(tenant_id);
CREATE INDEX idx_incident_follow_ups_incident ON incident_follow_ups(incident_id);
CREATE INDEX idx_incident_follow_ups_scheduled ON incident_follow_ups(scheduled_date) WHERE status = 'pending';
CREATE INDEX idx_incident_follow_ups_status ON incident_follow_ups(status);

-- RLS for follow-ups
ALTER TABLE incident_follow_ups ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_follow_ups ON incident_follow_ups
    USING (tenant_id = current_setting('app.current_tenant', true)::UUID);

-- Add behavior permissions
INSERT INTO permissions (id, code, name, description, module, created_at, updated_at)
VALUES
    (uuid_generate_v7(), 'behavior:read', 'View Behavioral Records', 'Permission to view student behavioral incident records', 'students', NOW(), NOW()),
    (uuid_generate_v7(), 'behavior:write', 'Manage Behavioral Records', 'Permission to create, update behavioral incident records', 'students', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign behavior permissions to admin roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('super_admin', 'admin')
AND p.code IN ('behavior:read', 'behavior:write')
ON CONFLICT DO NOTHING;
