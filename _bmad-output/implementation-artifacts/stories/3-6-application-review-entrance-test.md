# Story 3.6: Application Review & Entrance Test

**Epic:** 3 - School Setup & Admissions
**Status:** ready-for-dev
**Priority:** High
**Estimated Effort:** Medium

---

## User Story

As an **admission committee member**,
I want **to review applications and schedule entrance tests**,
So that **qualified candidates can be evaluated**.

---

## Acceptance Criteria

### AC1: Review Application
**Given** applications are submitted
**When** reviewing an application
**Then** admin sees all submitted information and documents
**And** admin can verify document authenticity
**And** admin can add review comments
**And** admin can update status: under_review, documents_pending, test_scheduled, rejected

### AC2: Schedule Entrance Test
**Given** entrance test is configured for a class
**When** scheduling tests
**Then** admin can set test date, time, and venue
**And** admin can assign students to test slots
**And** hall ticket is generated for each student
**And** SMS/email notification is sent with test details

### AC3: Record Test Results
**Given** entrance test is conducted
**When** entering test results
**Then** admin can enter marks for each subject
**And** total score is calculated automatically
**And** application moves to "test_completed" status

---

## Technical Requirements

### Backend (Go)

#### Database Schema

```sql
-- Entrance Tests
CREATE TABLE IF NOT EXISTS entrance_tests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id) ON DELETE CASCADE,
    test_name VARCHAR(200) NOT NULL,
    test_date DATE NOT NULL,
    start_time TIME NOT NULL,
    duration_minutes INT NOT NULL DEFAULT 60,
    venue VARCHAR(200),
    class_names JSONB DEFAULT '[]',
    max_candidates INT DEFAULT 100,
    status VARCHAR(20) DEFAULT 'scheduled',
    subjects JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Test Registrations
CREATE TABLE IF NOT EXISTS test_registrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES entrance_tests(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    roll_number VARCHAR(20),
    status VARCHAR(20) DEFAULT 'registered',
    marks JSONB DEFAULT '{}',
    total_marks DECIMAL(6,2),
    max_marks DECIMAL(6,2),
    percentage DECIMAL(5,2),
    result VARCHAR(20),
    remarks TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Application Reviews
CREATE TABLE IF NOT EXISTS application_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(id),
    review_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    comments TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| POST | /api/v1/applications/:id/review | Add review | applications:update |
| GET | /api/v1/applications/:id/reviews | Get reviews | applications:read |
| PATCH | /api/v1/applications/:id/documents/:docId/verify | Verify document | applications:update |
| GET | /api/v1/entrance-tests | List tests | admissions:read |
| POST | /api/v1/entrance-tests | Create test | admissions:create |
| PUT | /api/v1/entrance-tests/:id | Update test | admissions:update |
| DELETE | /api/v1/entrance-tests/:id | Delete test | admissions:delete |
| POST | /api/v1/entrance-tests/:id/register | Register candidate | admissions:update |
| POST | /api/v1/entrance-tests/:id/results | Submit results | admissions:update |
| GET | /api/v1/entrance-tests/:id/hall-tickets | Generate hall tickets | admissions:read |

### Frontend (Angular)

#### Feature Structure

```
src/app/features/admissions/
├── review/
│   ├── application-review.component.ts
│   └── document-verification.component.ts
├── tests/
│   ├── entrance-test.model.ts
│   ├── entrance-test.service.ts
│   ├── tests.component.ts
│   ├── test-form.component.ts
│   └── results-entry.component.ts
```

---

## Definition of Done

- [ ] Backend: Review and test endpoints
- [ ] Backend: Hall ticket generation
- [ ] Frontend: Application review page
- [ ] Frontend: Test scheduling UI
- [ ] Frontend: Results entry form
