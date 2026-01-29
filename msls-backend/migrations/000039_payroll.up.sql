-- Migration: 000039_payroll.up.sql
-- Story 5.6: Payroll Processing
-- Creates pay_runs, payslips, and payslip_components tables

-- Pay Run (monthly payroll batch)
CREATE TABLE pay_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    pay_period_month INTEGER NOT NULL,
    pay_period_year INTEGER NOT NULL,
    branch_id UUID REFERENCES branches(id),

    status VARCHAR(20) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'processing', 'calculated', 'approved', 'finalized', 'reversed')),

    total_staff INTEGER DEFAULT 0,
    total_gross DECIMAL(14,2) DEFAULT 0,
    total_deductions DECIMAL(14,2) DEFAULT 0,
    total_net DECIMAL(14,2) DEFAULT 0,

    calculated_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES users(id),
    finalized_at TIMESTAMPTZ,
    finalized_by UUID REFERENCES users(id),

    notes TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID REFERENCES users(id),

    -- Note: Unique constraint for period handled via index below
    CONSTRAINT check_pay_period_month CHECK (pay_period_month BETWEEN 1 AND 12)
);

-- Unique index for pay runs: one pay run per tenant/month/year/branch combination
-- Two separate indexes: one for when branch_id IS NULL, one for when it IS NOT NULL
CREATE UNIQUE INDEX idx_pay_runs_unique_period_no_branch
    ON pay_runs (tenant_id, pay_period_month, pay_period_year)
    WHERE branch_id IS NULL;

CREATE UNIQUE INDEX idx_pay_runs_unique_period_with_branch
    ON pay_runs (tenant_id, pay_period_month, pay_period_year, branch_id)
    WHERE branch_id IS NOT NULL;

-- Payslips (individual staff pay records for a pay run)
CREATE TABLE payslips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    pay_run_id UUID NOT NULL REFERENCES pay_runs(id) ON DELETE CASCADE,
    staff_id UUID NOT NULL REFERENCES staff(id),
    staff_salary_id UUID REFERENCES staff_salaries(id),

    -- Attendance summary for the period
    working_days INTEGER NOT NULL DEFAULT 0,
    present_days INTEGER NOT NULL DEFAULT 0,
    leave_days INTEGER NOT NULL DEFAULT 0,
    absent_days INTEGER NOT NULL DEFAULT 0,
    lop_days INTEGER NOT NULL DEFAULT 0,

    -- Calculated amounts
    gross_salary DECIMAL(12,2) NOT NULL,
    total_earnings DECIMAL(12,2) NOT NULL,
    total_deductions DECIMAL(12,2) NOT NULL,
    net_salary DECIMAL(12,2) NOT NULL,

    -- LOP deduction if any
    lop_deduction DECIMAL(12,2) DEFAULT 0,

    status VARCHAR(20) NOT NULL DEFAULT 'calculated'
        CHECK (status IN ('calculated', 'adjusted', 'approved', 'paid')),

    payment_date DATE,
    payment_reference VARCHAR(100),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(pay_run_id, staff_id)
);

-- Payslip Components (breakdown of each payslip)
CREATE TABLE payslip_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    payslip_id UUID NOT NULL REFERENCES payslips(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES salary_components(id),

    component_name VARCHAR(100) NOT NULL,
    component_code VARCHAR(20) NOT NULL,
    component_type VARCHAR(20) NOT NULL,

    amount DECIMAL(12,2) NOT NULL,
    is_prorated BOOLEAN DEFAULT false,

    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(payslip_id, component_id)
);

-- Enable RLS
ALTER TABLE pay_runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE payslips ENABLE ROW LEVEL SECURITY;
ALTER TABLE payslip_components ENABLE ROW LEVEL SECURITY;

-- RLS Policies for pay_runs
CREATE POLICY tenant_isolation_pay_runs ON pay_runs
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- RLS Policies for payslips
CREATE POLICY tenant_isolation_payslips ON payslips
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- RLS Policy for payslip_components (through payslips join)
CREATE POLICY tenant_isolation_payslip_components ON payslip_components
    USING (EXISTS (
        SELECT 1 FROM payslips p
        WHERE p.id = payslip_components.payslip_id
        AND p.tenant_id = current_setting('app.tenant_id', true)::UUID
    ));

-- Indexes for pay_runs
CREATE INDEX idx_pay_runs_tenant_id ON pay_runs(tenant_id);
CREATE INDEX idx_pay_runs_period ON pay_runs(tenant_id, pay_period_year, pay_period_month);
CREATE INDEX idx_pay_runs_status ON pay_runs(tenant_id, status);
CREATE INDEX idx_pay_runs_branch_id ON pay_runs(branch_id) WHERE branch_id IS NOT NULL;

-- Indexes for payslips
CREATE INDEX idx_payslips_tenant_id ON payslips(tenant_id);
CREATE INDEX idx_payslips_pay_run_id ON payslips(pay_run_id);
CREATE INDEX idx_payslips_staff_id ON payslips(staff_id);
CREATE INDEX idx_payslips_status ON payslips(status);

-- Indexes for payslip_components
CREATE INDEX idx_payslip_components_payslip_id ON payslip_components(payslip_id);
CREATE INDEX idx_payslip_components_component_id ON payslip_components(component_id);

-- Add payroll permissions
INSERT INTO permissions (code, name, module, description) VALUES
    ('payroll.view', 'View Payroll', 'payroll', 'View payroll runs and payslips'),
    ('payroll.create', 'Create Payroll', 'payroll', 'Create payroll runs'),
    ('payroll.calculate', 'Calculate Payroll', 'payroll', 'Calculate payroll'),
    ('payroll.approve', 'Approve Payroll', 'payroll', 'Approve payroll runs'),
    ('payroll.finalize', 'Finalize Payroll', 'payroll', 'Finalize payroll runs'),
    ('payroll.adjust', 'Adjust Payslips', 'payroll', 'Adjust individual payslips'),
    ('payroll.export', 'Export Payroll', 'payroll', 'Export payroll reports'),
    ('payroll.delete', 'Delete Payroll', 'payroll', 'Delete draft payroll runs')
ON CONFLICT (code) DO NOTHING;

-- Assign payroll permissions to Super Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Super Admin'
AND p.code IN (
    'payroll.view', 'payroll.create', 'payroll.calculate',
    'payroll.approve', 'payroll.finalize', 'payroll.adjust',
    'payroll.export', 'payroll.delete'
)
ON CONFLICT DO NOTHING;

-- Assign payroll permissions to HR Manager role (if exists)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'HR Manager'
AND p.code IN (
    'payroll.view', 'payroll.create', 'payroll.calculate',
    'payroll.approve', 'payroll.finalize', 'payroll.adjust',
    'payroll.export'
)
ON CONFLICT DO NOTHING;

-- Assign view-only payroll permission to Accountant role (if exists)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Accountant'
AND p.code IN ('payroll.view', 'payroll.export')
ON CONFLICT DO NOTHING;
