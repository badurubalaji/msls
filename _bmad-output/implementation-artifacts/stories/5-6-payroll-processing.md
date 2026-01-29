# Story 5.6: Payroll Processing

**Epic:** 5 - Staff Management
**Status:** ready-for-review
**Priority:** P1 (Required for staff salary disbursement)
**Story Points:** 8

---

## User Story

**As an** HR/Finance administrator,
**I want** to process monthly payroll for all staff,
**So that** salaries can be calculated and disbursed on time.

---

## Acceptance Criteria

### AC1: Payroll Run Initiation

**Given** HR is on payroll processing page
**When** initiating a new pay run
**Then** they can select: month, year, branch (optional)
**And** system validates no duplicate pay run exists for the period
**And** system shows staff count and estimated total

### AC2: Payroll Calculation

**Given** a pay run is initiated
**When** calculating payroll
**Then** system calculates each staff member's salary based on:
  - Assigned salary structure/components
  - Attendance for the month (present days, leaves, absences)
  - Any LOP (Loss of Pay) deductions
**And** gross salary, deductions, and net salary are computed
**And** calculation summary shows totals and breakdowns

### AC3: Payslip Generation

**Given** payroll is calculated
**When** viewing individual payslip
**Then** they see: employee details, pay period
**And** they see: all earnings with amounts
**And** they see: all deductions with amounts
**And** they see: gross, total deductions, net pay
**And** payslip can be downloaded as PDF

### AC4: Pay Run Approval & Finalization

**Given** payroll is calculated
**When** reviewing pay run
**Then** HR can review: summary by department, exceptions list
**And** HR can modify individual calculations if needed
**And** HR can approve/finalize the pay run
**And** finalized pay run cannot be edited (only reversed)

### AC5: Payroll Reports & Export

**Given** a pay run is finalized
**When** generating reports
**Then** they can export: bank transfer file (CSV/Excel)
**And** they can export: payslip batch (PDF)
**And** they can view: department-wise summary
**And** they can view: year-to-date summary per employee

---

## Technical Implementation

### Database Schema

```sql
-- Pay Run (monthly payroll batch)
CREATE TABLE pay_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    pay_period_month INTEGER NOT NULL CHECK (pay_period_month BETWEEN 1 AND 12),
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

    UNIQUE(tenant_id, pay_period_month, pay_period_year, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'))
);

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

-- RLS Policies
ALTER TABLE pay_runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE payslips ENABLE ROW LEVEL SECURITY;
ALTER TABLE payslip_components ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_pay_runs ON pay_runs
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_payslips ON payslips
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);
```

### API Endpoints

```
# Pay Runs
GET    /api/v1/payroll/runs                - List pay runs
POST   /api/v1/payroll/runs                - Create new pay run
GET    /api/v1/payroll/runs/:id            - Get pay run details
POST   /api/v1/payroll/runs/:id/calculate  - Calculate payroll
POST   /api/v1/payroll/runs/:id/approve    - Approve pay run
POST   /api/v1/payroll/runs/:id/finalize   - Finalize pay run
DELETE /api/v1/payroll/runs/:id            - Delete draft pay run

# Payslips
GET    /api/v1/payroll/runs/:id/payslips   - List payslips in a pay run
GET    /api/v1/payroll/payslips/:id        - Get payslip details
PUT    /api/v1/payroll/payslips/:id        - Adjust payslip (before approval)
GET    /api/v1/payroll/payslips/:id/pdf    - Download payslip PDF

# Reports
GET    /api/v1/payroll/runs/:id/export     - Export bank transfer file
GET    /api/v1/payroll/runs/:id/summary    - Get department-wise summary

# Staff payslip history
GET    /api/v1/staff/:id/payslips          - Get staff payslip history
```

### Frontend Components

```
src/app/features/payroll/
├── payroll.routes.ts
├── models/
│   └── payroll.model.ts
├── services/
│   └── payroll.service.ts
└── pages/
    ├── pay-runs/
    │   ├── pay-runs.component.ts
    │   └── pay-run-form.component.ts
    ├── pay-run-detail/
    │   └── pay-run-detail.component.ts
    └── payslip-detail/
        └── payslip-detail.component.ts
```

---

## Tasks

### Backend Tasks

- [x] BE-5.6.1: Create payroll tables migration
- [x] BE-5.6.2: Create pay run module (model, repository, service, handler)
- [x] BE-5.6.3: Create payslip module (model, repository, service, handler)
- [x] BE-5.6.4: Implement payroll calculation logic
- [x] BE-5.6.5: Add payroll permissions and routes

### Frontend Tasks

- [x] FE-5.6.1: Create payroll models and service
- [x] FE-5.6.2: Create pay runs list page
- [x] FE-5.6.3: Create pay run detail page with payslips
- [x] FE-5.6.4: Create payslip detail view
- [x] FE-5.6.5: Add payroll routes and navigation

---

## Dependencies

- Story 5.5: Salary Structure Configuration (completed)
- Story 5.4: Staff Attendance Tracking (for LOP calculation)

## Notes

- LOP (Loss of Pay) = (Daily rate) × (Absent days without approved leave)
- Daily rate = Monthly gross / Working days in month
- Payslip should show YTD (Year-to-Date) totals
- Bank transfer file format: Staff Name, Bank Account, IFSC, Net Amount
