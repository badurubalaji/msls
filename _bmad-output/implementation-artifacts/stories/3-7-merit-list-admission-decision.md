# Story 3.7: Merit List & Admission Decision

**Epic:** 3 - School Setup & Admissions
**Status:** done
**Priority:** High
**Estimated Effort:** Medium

---

## User Story

As an **admission committee**,
I want **to generate merit lists and make admission decisions**,
So that **students can be selected fairly based on criteria**.

---

## Acceptance Criteria

### AC1: Generate Merit List
**Given** test results are entered
**When** generating merit list
**Then** students are ranked by total score
**And** merit list shows: rank, name, score, status
**And** cutoff score can be applied to filter

### AC2: Make Admission Decisions
**Given** a merit list is generated
**When** making admission decisions
**Then** admin can select students for admission
**And** admin can mark: selected, waitlisted, rejected
**And** selected students' status changes to "offer_sent"
**And** offer letter is generated with fee details

### AC3: Accept Offer and Enroll
**Given** an offer is sent
**When** parent accepts and pays admission fee
**Then** payment is recorded
**And** application status changes to "enrolled"
**And** student record is created from application data
**And** admission process is complete

---

## Technical Requirements

### Backend (Go)

#### Database Schema

```sql
-- Admission Decisions
CREATE TABLE IF NOT EXISTS admission_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    decision VARCHAR(20) NOT NULL,
    decision_date DATE NOT NULL,
    decided_by UUID REFERENCES users(id),
    section_assigned VARCHAR(50),
    waitlist_position INT,
    rejection_reason TEXT,
    offer_letter_url VARCHAR(500),
    offer_valid_until DATE,
    offer_accepted BOOLEAN,
    offer_accepted_at TIMESTAMPTZ,
    remarks TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Merit Lists (snapshot for audit)
CREATE TABLE IF NOT EXISTS merit_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id) ON DELETE CASCADE,
    class_name VARCHAR(50) NOT NULL,
    test_id UUID REFERENCES entrance_tests(id),
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    generated_by UUID REFERENCES users(id),
    cutoff_score DECIMAL(5,2),
    entries JSONB NOT NULL DEFAULT '[]',
    is_final BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| POST | /api/v1/admission-sessions/:id/merit-list | Generate merit list | admissions:update |
| GET | /api/v1/admission-sessions/:id/merit-list | Get merit list | admissions:read |
| POST | /api/v1/applications/:id/decision | Make decision | admissions:update |
| POST | /api/v1/applications/:id/offer-letter | Generate offer letter | admissions:update |
| POST | /api/v1/applications/:id/accept-offer | Accept offer | applications:update |
| POST | /api/v1/applications/:id/enroll | Complete enrollment | admissions:update |

### Frontend (Angular)

#### Feature Structure

```
src/app/features/admissions/
├── merit/
│   ├── merit-list.component.ts
│   └── decision-form.component.ts
├── enrollment/
│   └── enrollment.component.ts
```

---

## Definition of Done

- [x] Backend: Merit list generation
- [x] Backend: Decision endpoints
- [x] Backend: Offer letter PDF generation (placeholder URL generation)
- [x] Frontend: Merit list view with actions
- [x] Frontend: Decision making UI
- [x] Frontend: Enrollment completion
