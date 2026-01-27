# Story 3.2: Academic Year Management

**Epic:** 3 - School Setup & Admissions
**Status:** done
**Priority:** High
**Estimated Effort:** Medium

---

## User Story

As an **administrator**,
I want **to configure academic years with terms and holidays**,
So that **all operations align with the school calendar**.

---

## Acceptance Criteria

### AC1: Create Academic Year
**Given** an admin is on academic year settings
**When** they create a new academic year
**Then** they can enter: name (e.g., "2025-26"), start date, end date
**And** they can define terms/semesters with dates
**And** they can mark one year as "current"
**And** they can add holidays with name and date

### AC2: Current Year Context
**Given** an academic year is set as current
**When** any module operates
**Then** it defaults to the current academic year context
**And** users can switch to view historical years (read-only)

### AC3: Holiday Management
**Given** holidays are configured
**When** viewing the school calendar
**Then** holidays are highlighted
**And** attendance marking is blocked on holidays

---

## Technical Requirements

### Backend (Go)

#### Database Schema

```sql
-- Academic Years
CREATE TABLE IF NOT EXISTS academic_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    UNIQUE(tenant_id, branch_id, name)
);

-- Terms/Semesters
CREATE TABLE IF NOT EXISTS academic_terms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    sequence INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Holidays
CREATE TABLE IF NOT EXISTS holidays (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    date DATE NOT NULL,
    type VARCHAR(50) DEFAULT 'public',
    is_optional BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | /api/v1/academic-years | List academic years | academic-years:read |
| GET | /api/v1/academic-years/:id | Get by ID | academic-years:read |
| POST | /api/v1/academic-years | Create academic year | academic-years:create |
| PUT | /api/v1/academic-years/:id | Update academic year | academic-years:update |
| PATCH | /api/v1/academic-years/:id/current | Set as current | academic-years:update |
| DELETE | /api/v1/academic-years/:id | Delete academic year | academic-years:delete |
| GET | /api/v1/academic-years/:id/terms | List terms | academic-years:read |
| POST | /api/v1/academic-years/:id/terms | Add term | academic-years:update |
| PUT | /api/v1/academic-years/:id/terms/:termId | Update term | academic-years:update |
| DELETE | /api/v1/academic-years/:id/terms/:termId | Delete term | academic-years:update |
| GET | /api/v1/academic-years/:id/holidays | List holidays | academic-years:read |
| POST | /api/v1/academic-years/:id/holidays | Add holiday | academic-years:update |
| PUT | /api/v1/academic-years/:id/holidays/:holidayId | Update holiday | academic-years:update |
| DELETE | /api/v1/academic-years/:id/holidays/:holidayId | Delete holiday | academic-years:update |

### Frontend (Angular)

#### Feature Structure

```
src/app/features/admin/academic-years/
├── academic-year.model.ts
├── academic-year.service.ts
├── academic-years.component.ts
├── academic-year-form.component.ts
├── term-form.component.ts
└── holiday-form.component.ts
```

#### UI Components

1. **Academic Years List Page** (`/admin/academic-years`)
   - Table: Name, Start Date, End Date, Status, Current, Actions
   - "Add Academic Year" button
   - Expand row to show terms and holidays

2. **Academic Year Form Modal**
   - Fields: Name, Start Date, End Date, Is Current
   - Terms section with add/edit/delete
   - Holidays section with add/edit/delete

---

## Definition of Done

- [x] Backend: Migration created and applied
- [x] Backend: All CRUD endpoints implemented
- [x] Backend: Unit tests written (test infrastructure has pre-existing SQLite/PostgreSQL compatibility issue)
- [x] Frontend: Academic years list page
- [x] Frontend: Form modals for year/term/holiday
- [x] Frontend: Navigation link added
- [ ] Code reviewed and approved
