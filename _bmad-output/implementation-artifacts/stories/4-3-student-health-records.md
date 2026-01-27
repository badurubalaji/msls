# Story 4.3: Student Health Records

**Epic:** 4 - Student Lifecycle Management
**Status:** done
**Priority:** P1
**Story Points:** 5

---

## User Story

**As a** school nurse/administrator,
**I want** to maintain student health records,
**So that** medical emergencies can be handled with proper information.

---

## Acceptance Criteria

### AC1: Health Information Entry
- [x] Can enter basic health info: blood group, height (cm), weight (kg), vision status
- [x] Can record allergies with severity level (mild, moderate, severe)
- [x] Can record chronic conditions (asthma, diabetes, epilepsy, etc.)
- [x] Can record regular medications (name, dosage, frequency)
- [x] Health data marked as confidential (restricted access via permission)

### AC2: Vaccination Tracking
- [x] Can enter vaccination records: vaccine name, date administered, next due date
- [ ] Can set reminders for upcoming vaccinations (30 days before due) - deferred to notifications epic
- [x] Can upload vaccination certificate (PDF/image)
- [x] Vaccination history displayed chronologically

### AC3: Medical Incident Recording
- [x] Can record medical incidents: date, time, description, action taken
- [x] Can record: parent notified (yes/no), hospital visit required (yes/no)
- [x] Incident history maintained and searchable
- [x] Incidents linked to student timeline

### AC4: Health Tab in Student Profile
- [x] Health section accessible only with `health:read` permission
- [x] Tab shows summary cards: allergies, conditions, last checkup
- [x] Quick access to vaccination schedule
- [x] Medical incident history with filters

---

## Technical Requirements

### Backend

**Database Tables:**

```sql
-- Migration: 20260126100200_create_health_records_table.up.sql

CREATE TABLE student_health_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    blood_group VARCHAR(5),
    height_cm DECIMAL(5,2),
    weight_kg DECIMAL(5,2),
    vision_status VARCHAR(50), -- normal, corrected, impaired
    last_checkup_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    CONSTRAINT uniq_health_record_student UNIQUE (student_id)
);

ALTER TABLE student_health_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON student_health_records USING (tenant_id = current_setting('app.current_tenant')::UUID);

CREATE TABLE student_allergies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    allergen VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('mild', 'moderate', 'severe')),
    reaction_description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE student_allergies ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON student_allergies USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_allergies_student ON student_allergies(student_id);

CREATE TABLE student_conditions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    condition_name VARCHAR(100) NOT NULL,
    diagnosed_date DATE,
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE student_conditions ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON student_conditions USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_conditions_student ON student_conditions(student_id);

CREATE TABLE student_medications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    medication_name VARCHAR(100) NOT NULL,
    dosage VARCHAR(50) NOT NULL,
    frequency VARCHAR(50) NOT NULL, -- daily, twice daily, as needed
    start_date DATE,
    end_date DATE,
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE student_medications ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON student_medications USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_medications_student ON student_medications(student_id);

CREATE TABLE student_vaccinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    vaccine_name VARCHAR(100) NOT NULL,
    date_administered DATE NOT NULL,
    next_due_date DATE,
    administered_by VARCHAR(200),
    certificate_url VARCHAR(500),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE student_vaccinations ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON student_vaccinations USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_vaccinations_student ON student_vaccinations(student_id);
CREATE INDEX idx_vaccinations_due ON student_vaccinations(next_due_date) WHERE next_due_date IS NOT NULL;

CREATE TABLE medical_incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    incident_date DATE NOT NULL,
    incident_time TIME NOT NULL,
    location VARCHAR(200),
    description TEXT NOT NULL,
    symptoms TEXT,
    action_taken TEXT NOT NULL,
    parent_notified BOOLEAN NOT NULL DEFAULT FALSE,
    parent_notified_at TIMESTAMPTZ,
    hospital_visit_required BOOLEAN NOT NULL DEFAULT FALSE,
    hospital_name VARCHAR(200),
    follow_up_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

ALTER TABLE medical_incidents ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON medical_incidents USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_medical_incidents_student ON medical_incidents(student_id);
CREATE INDEX idx_medical_incidents_date ON medical_incidents(incident_date DESC);
```

**API Endpoints:**

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/students/{id}/health` | Get health record | `health:read` |
| PUT | `/api/v1/students/{id}/health` | Update health record | `health:write` |
| GET | `/api/v1/students/{id}/allergies` | List allergies | `health:read` |
| POST | `/api/v1/students/{id}/allergies` | Add allergy | `health:write` |
| DELETE | `/api/v1/students/{id}/allergies/{aid}` | Remove allergy | `health:write` |
| GET | `/api/v1/students/{id}/conditions` | List conditions | `health:read` |
| POST | `/api/v1/students/{id}/conditions` | Add condition | `health:write` |
| PUT | `/api/v1/students/{id}/conditions/{cid}` | Update condition | `health:write` |
| GET | `/api/v1/students/{id}/medications` | List medications | `health:read` |
| POST | `/api/v1/students/{id}/medications` | Add medication | `health:write` |
| PUT | `/api/v1/students/{id}/medications/{mid}` | Update medication | `health:write` |
| GET | `/api/v1/students/{id}/vaccinations` | List vaccinations | `health:read` |
| POST | `/api/v1/students/{id}/vaccinations` | Add vaccination | `health:write` |
| GET | `/api/v1/students/{id}/medical-incidents` | List incidents | `health:read` |
| POST | `/api/v1/students/{id}/medical-incidents` | Record incident | `health:write` |
| GET | `/api/v1/vaccinations/due` | List students with due vaccinations | `health:read` |

**New Permissions to Add:**

```go
var HealthPermissions = []Permission{
    {Code: "health:read", Name: "View Health Records", Module: "health"},
    {Code: "health:write", Name: "Edit Health Records", Module: "health"},
}
```

### Frontend

**Components to create:**

```bash
ng generate component features/students/components/health-overview --standalone
ng generate component features/students/components/allergy-form --standalone
ng generate component features/students/components/vaccination-form --standalone
ng generate component features/students/components/medical-incident-form --standalone
ng generate interface features/students/models/health-record
ng generate interface features/students/models/vaccination
ng generate interface features/students/models/medical-incident
```

---

## Tasks

### Backend Tasks

- [ ] **BE-4.3.1**: Create health-related entity migrations
- [ ] **BE-4.3.2**: Create health record repository
- [ ] **BE-4.3.3**: Create health service with CRUD operations
- [ ] **BE-4.3.4**: Create health HTTP handlers
- [ ] **BE-4.3.5**: Add vaccination reminder query (due in 30 days)
- [ ] **BE-4.3.6**: Add health permissions to seed
- [ ] **BE-4.3.7**: Write unit tests

### Frontend Tasks

- [ ] **FE-4.3.1**: Create health-related interfaces
- [ ] **FE-4.3.2**: Create health service
- [ ] **FE-4.3.3**: Create health overview component
- [ ] **FE-4.3.4**: Create allergy/condition/medication forms
- [ ] **FE-4.3.5**: Create vaccination form with certificate upload
- [ ] **FE-4.3.6**: Create medical incident form
- [ ] **FE-4.3.7**: Add Health tab to student profile
- [ ] **FE-4.3.8**: Implement permission-based visibility
- [ ] **FE-4.3.9**: Write component tests

---

## Definition of Done

- [ ] All acceptance criteria verified
- [ ] Health records only visible with health:read permission
- [ ] Vaccination reminders working
- [ ] Certificate upload working
- [ ] Medical incident history displays correctly
- [ ] Backend tests passing
- [ ] Frontend tests passing

---

## Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Story 4.1 (Student Profile) | Required | Health links to student |
| RBAC (Epic 2) | ✅ Done | For health permissions |
| File storage | Required | For vaccination certificates |

---

## Notes

- Health data is sensitive - ensure proper access control
- Consider HIPAA-like privacy requirements
- Vaccination reminders could trigger notifications (Epic 12)
- BMI calculation: weight_kg / (height_cm/100)^2
