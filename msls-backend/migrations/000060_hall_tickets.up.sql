-- Hall Ticket Generation
-- Story 8.3: Hall Ticket Generation

-- ============================================================
-- Hall Ticket Templates
-- ============================================================

CREATE TABLE hall_ticket_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    header_logo_url VARCHAR(500),
    school_name VARCHAR(200),
    school_address TEXT,
    instructions TEXT,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Enable RLS
ALTER TABLE hall_ticket_templates ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_hall_ticket_templates ON hall_ticket_templates
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_hall_ticket_templates_tenant ON hall_ticket_templates(tenant_id);
CREATE INDEX idx_hall_ticket_templates_default ON hall_ticket_templates(tenant_id, is_default) WHERE is_default = true;

-- Updated at trigger
CREATE TRIGGER set_updated_at_hall_ticket_templates
    BEFORE UPDATE ON hall_ticket_templates
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- Hall Tickets
-- ============================================================

CREATE TABLE hall_tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id),
    roll_number VARCHAR(50) NOT NULL,
    qr_code_data VARCHAR(500) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'generated',
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    printed_at TIMESTAMPTZ,
    downloaded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_hall_ticket_status CHECK (status IN ('generated', 'printed', 'downloaded')),
    CONSTRAINT uq_hall_ticket_exam_student UNIQUE (examination_id, student_id)
);

-- Enable RLS
ALTER TABLE hall_tickets ENABLE ROW LEVEL SECURITY;

-- RLS Policy
CREATE POLICY tenant_isolation_hall_tickets ON hall_tickets
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_hall_tickets_tenant ON hall_tickets(tenant_id);
CREATE INDEX idx_hall_tickets_examination ON hall_tickets(examination_id);
CREATE INDEX idx_hall_tickets_student ON hall_tickets(student_id);
CREATE INDEX idx_hall_tickets_roll_number ON hall_tickets(tenant_id, examination_id, roll_number);
CREATE INDEX idx_hall_tickets_status ON hall_tickets(tenant_id, examination_id, status);

-- Updated at trigger
CREATE TRIGGER set_updated_at_hall_tickets
    BEFORE UPDATE ON hall_tickets
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- Permissions
-- ============================================================

INSERT INTO permissions (id, code, name, description, module, created_at, updated_at)
VALUES
    (uuid_generate_v7(), 'hall-ticket:view', 'View Hall Tickets', 'Permission to view hall tickets', 'hall-ticket', NOW(), NOW()),
    (uuid_generate_v7(), 'hall-ticket:generate', 'Generate Hall Tickets', 'Permission to generate hall tickets', 'hall-ticket', NOW(), NOW()),
    (uuid_generate_v7(), 'hall-ticket:download', 'Download Hall Tickets', 'Permission to download hall ticket PDFs', 'hall-ticket', NOW(), NOW()),
    (uuid_generate_v7(), 'hall-ticket:template-manage', 'Manage Hall Ticket Templates', 'Permission to create and manage hall ticket templates', 'hall-ticket', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Assign permissions to roles
-- Super admin, admin, principal - full access
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name IN ('super_admin', 'admin', 'principal')
AND p.code IN ('hall-ticket:view', 'hall-ticket:generate', 'hall-ticket:download', 'hall-ticket:template-manage')
ON CONFLICT DO NOTHING;

-- Coordinators can view, generate, and download
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'coordinator'
AND p.code IN ('hall-ticket:view', 'hall-ticket:generate', 'hall-ticket:download')
ON CONFLICT DO NOTHING;

-- Teachers can view and download
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher'
AND p.code IN ('hall-ticket:view', 'hall-ticket:download')
ON CONFLICT DO NOTHING;
