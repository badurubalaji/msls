# Story 5.5: Salary Structure Configuration

**Epic:** 5 - Staff Management
**Status:** complete
**Priority:** P1 (Foundation for payroll processing)
**Story Points:** 8

---

## User Story

**As an** HR administrator,
**I want** to configure salary structures with components,
**So that** payroll can be calculated correctly for staff.

---

## Acceptance Criteria

### AC1: Salary Component Management

**Given** HR is on salary settings
**When** managing salary components
**Then** they can create earning components (Basic, HRA, DA, TA, Special Allowance)
**And** they can create deduction components (PF, ESI, Professional Tax, TDS)
**And** each component has: name, code, type (earning/deduction), calculation type (fixed/percentage)
**And** percentage-based components reference a base component

### AC2: Salary Structure Creation

**Given** HR is creating a salary structure
**When** defining the structure
**Then** they can enter: structure name, code, description
**And** they can add multiple earning components with amounts/percentages
**And** they can add multiple deduction components with amounts/percentages
**And** they can link structure to designation (optional)
**And** gross and net salary are auto-calculated for preview

### AC3: Staff Salary Assignment

**Given** a staff member exists
**When** assigning salary structure
**Then** they can select a predefined salary structure
**And** they can override individual component values
**And** effective date is recorded
**And** CTC (Cost to Company) is calculated and displayed

### AC4: Salary Revision History

**Given** a staff member has assigned salary
**When** salary is revised
**Then** new effective date is set
**And** previous salary is archived with end date
**And** revision history is maintained with old/new values
**And** revision reason is captured

### AC5: Salary Structure Listing

**Given** HR is on salary structures page
**When** viewing structures
**Then** they see: name, code, linked designation, component count, staff count
**And** they can filter by status (active/inactive)
**And** they can search by name or code

---

## Technical Implementation

### Database Schema

```sql
-- Salary Components (earning/deduction types)
CREATE TABLE salary_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    component_type VARCHAR(20) NOT NULL CHECK (component_type IN ('earning', 'deduction')),
    calculation_type VARCHAR(20) NOT NULL CHECK (calculation_type IN ('fixed', 'percentage')),
    percentage_of_id UUID REFERENCES salary_components(id), -- For percentage-based components
    is_taxable BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(tenant_id, code)
);

-- Salary Structures (templates)
CREATE TABLE salary_structures (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    designation_id UUID REFERENCES designations(id),
    is_active BOOLEAN DEFAULT true,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(tenant_id, code)
);

-- Structure Components (components in a structure with default values)
CREATE TABLE salary_structure_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    structure_id UUID NOT NULL REFERENCES salary_structures(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES salary_components(id),

    amount DECIMAL(12,2), -- For fixed amounts
    percentage DECIMAL(5,2), -- For percentage-based

    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(structure_id, component_id)
);

-- Staff Salary Assignment
CREATE TABLE staff_salaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id),
    structure_id UUID REFERENCES salary_structures(id),

    effective_from DATE NOT NULL,
    effective_to DATE, -- NULL means current

    gross_salary DECIMAL(12,2) NOT NULL,
    net_salary DECIMAL(12,2) NOT NULL,
    ctc DECIMAL(12,2), -- Cost to Company

    revision_reason TEXT,
    is_current BOOLEAN DEFAULT true,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Staff Salary Components (actual values for staff)
CREATE TABLE staff_salary_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    staff_salary_id UUID NOT NULL REFERENCES staff_salaries(id) ON DELETE CASCADE,
    component_id UUID NOT NULL REFERENCES salary_components(id),

    amount DECIMAL(12,2) NOT NULL,
    is_overridden BOOLEAN DEFAULT false, -- True if different from structure default

    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(staff_salary_id, component_id)
);

-- RLS Policies
ALTER TABLE salary_components ENABLE ROW LEVEL SECURITY;
ALTER TABLE salary_structures ENABLE ROW LEVEL SECURITY;
ALTER TABLE staff_salaries ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_salary_components ON salary_components
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_salary_structures ON salary_structures
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

CREATE POLICY tenant_isolation_staff_salaries ON staff_salaries
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);
```

### API Endpoints

```
# Salary Components
GET    /api/v1/salary-components           - List salary components
POST   /api/v1/salary-components           - Create component
GET    /api/v1/salary-components/:id       - Get component
PUT    /api/v1/salary-components/:id       - Update component
DELETE /api/v1/salary-components/:id       - Delete component
GET    /api/v1/salary-components/dropdown  - Get for dropdown

# Salary Structures
GET    /api/v1/salary-structures           - List structures
POST   /api/v1/salary-structures           - Create structure
GET    /api/v1/salary-structures/:id       - Get structure with components
PUT    /api/v1/salary-structures/:id       - Update structure
DELETE /api/v1/salary-structures/:id       - Delete structure
GET    /api/v1/salary-structures/dropdown  - Get for dropdown

# Staff Salary
GET    /api/v1/staff/:id/salary            - Get staff current salary
POST   /api/v1/staff/:id/salary            - Assign/revise salary
GET    /api/v1/staff/:id/salary/history    - Get salary revision history
```

### Frontend Components

```
src/app/features/admin/salary/
├── salary.routes.ts
├── models/
│   ├── salary-component.model.ts
│   └── salary-structure.model.ts
├── services/
│   ├── salary-component.service.ts
│   └── salary-structure.service.ts
└── pages/
    ├── salary-components/
    │   ├── salary-components.component.ts
    │   └── salary-component-form.component.ts
    └── salary-structures/
        ├── salary-structures.component.ts
        └── salary-structure-form.component.ts

# Staff salary assignment in staff module
src/app/features/staff/
├── components/
│   └── staff-salary/
│       ├── staff-salary.component.ts
│       └── salary-assignment-form.component.ts
```

---

## Tasks

### Backend Tasks

- [x] BE-5.5.1: Create salary tables migration
- [x] BE-5.5.2: Create salary component module (model, repository, service, handler)
- [x] BE-5.5.3: Create salary structure module (model, repository, service, handler)
- [x] BE-5.5.4: Create staff salary module (assignment, revision, history)
- [x] BE-5.5.5: Add salary permissions and seed default components

### Frontend Tasks

- [x] FE-5.5.1: Create salary component models and service
- [x] FE-5.5.2: Create salary components list and form pages
- [x] FE-5.5.3: Create salary structure models and service
- [x] FE-5.5.4: Create salary structures list and form pages
- [x] FE-5.5.5: Create staff salary assignment component
- [x] FE-5.5.6: Add navigation links and routes

### Testing Tasks

- [ ] BE-5.5.6: Write unit tests for salary modules

---

## Dependencies

- Story 5.1: Staff Profile Management (completed)
- Story 5.3: Department & Designation Hierarchy (completed)

## Notes

- Default Indian salary components: Basic, HRA (40-50% of Basic), DA, TA, Special Allowance
- Default deductions: PF (12% of Basic), ESI (0.75% if applicable), Professional Tax
- Gross = Sum of all earnings
- Net = Gross - Sum of all deductions
- CTC = Gross + Employer contributions (PF employer share, etc.)
