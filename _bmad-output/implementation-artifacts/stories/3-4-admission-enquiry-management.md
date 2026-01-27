# Story 3.4: Admission Enquiry Management

**Epic:** 3 - School Setup & Admissions
**Status:** done
**Priority:** High
**Estimated Effort:** Medium

---

## User Story

As a **front office staff**,
I want **to capture and track admission enquiries**,
So that **interested families can be followed up systematically**.

---

## Acceptance Criteria

### AC1: Create Enquiry
**Given** a staff member is handling an enquiry
**When** they create a new enquiry
**Then** they can enter: parent name, phone, email, student name, class interested
**And** they can add notes from the conversation
**And** enquiry is assigned a unique enquiry number
**And** enquiry status is "new"

### AC2: List and Filter Enquiries
**Given** enquiries exist in the system
**When** viewing the enquiry list
**Then** they see all enquiries with status (new, contacted, interested, converted, closed)
**And** they can filter by status, class, date range
**And** they can add follow-up notes with date

### AC3: Convert to Application
**Given** an enquiry is converted
**When** creating an application from enquiry
**Then** enquiry data pre-fills the application form
**And** enquiry status changes to "converted"
**And** enquiry links to the application

---

## Technical Requirements

### Backend (Go)

#### Database Schema

```sql
-- Admission Enquiries
CREATE TABLE IF NOT EXISTS admission_enquiries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    session_id UUID REFERENCES admission_sessions(id),
    enquiry_number VARCHAR(50) NOT NULL,
    student_name VARCHAR(200) NOT NULL,
    date_of_birth DATE,
    gender VARCHAR(20),
    class_applying VARCHAR(50) NOT NULL,
    parent_name VARCHAR(200) NOT NULL,
    parent_phone VARCHAR(20) NOT NULL,
    parent_email VARCHAR(255),
    source VARCHAR(50) DEFAULT 'walk_in',
    referral_details TEXT,
    remarks TEXT,
    status VARCHAR(20) DEFAULT 'new',
    follow_up_date DATE,
    assigned_to UUID REFERENCES users(id),
    converted_application_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    UNIQUE(tenant_id, enquiry_number)
);

-- Enquiry Follow-ups
CREATE TABLE IF NOT EXISTS enquiry_follow_ups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    enquiry_id UUID NOT NULL REFERENCES admission_enquiries(id) ON DELETE CASCADE,
    follow_up_date DATE NOT NULL,
    contact_mode VARCHAR(20) DEFAULT 'phone',
    notes TEXT,
    outcome VARCHAR(20),
    next_follow_up DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);
```

#### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | /api/v1/enquiries | List enquiries | enquiries:read |
| GET | /api/v1/enquiries/:id | Get by ID | enquiries:read |
| POST | /api/v1/enquiries | Create enquiry | enquiries:create |
| PUT | /api/v1/enquiries/:id | Update enquiry | enquiries:update |
| DELETE | /api/v1/enquiries/:id | Delete enquiry | enquiries:delete |
| POST | /api/v1/enquiries/:id/follow-ups | Add follow-up | enquiries:update |
| GET | /api/v1/enquiries/:id/follow-ups | List follow-ups | enquiries:read |
| POST | /api/v1/enquiries/:id/convert | Convert to application | enquiries:update |

### Frontend (Angular)

#### Feature Structure

```
src/app/features/admissions/enquiries/
├── enquiry.model.ts
├── enquiry.service.ts
├── enquiries.component.ts
├── enquiry-form.component.ts
└── follow-up-form.component.ts
```

---

## Definition of Done

- [x] Backend: Migration and endpoints
- [x] Backend: Enquiry number generation
- [x] Frontend: Enquiry list with filters
- [x] Frontend: Enquiry form modal
- [x] Frontend: Follow-up tracking
- [x] Frontend: Convert to application action
