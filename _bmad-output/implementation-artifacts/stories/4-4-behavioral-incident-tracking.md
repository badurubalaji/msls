# Story 4.4: Student Behavioral Tracking

**Epic:** 4 - Student Lifecycle Management
**Status:** done
**Priority:** P2
**Story Points:** 5

---

## User Story

**As a** teacher or administrator,
**I want** to record and track student behavioral incidents,
**So that** patterns can be identified and addressed.

---

## Acceptance Criteria

### AC1: Incident Recording
- [x] Can record incident: date, time, location, incident type
- [x] Incident types: positive recognition, minor infraction, major violation
- [x] Can describe: what happened, witnesses, student response
- [x] Can record: action taken, parent meeting required
- [x] Recorded by (user) automatically captured

### AC2: Follow-up Actions
- [x] Can schedule follow-up: meeting date, participants, expected outcomes
- [ ] Notification sent to relevant parties when follow-up scheduled - deferred to notifications epic
- [x] Meeting outcome can be recorded later
- [x] Follow-up status: pending, completed, cancelled

### AC3: Behavioral History View
- [x] Chronological list of all incidents for a student
- [x] Filters: type, date range, severity
- [x] Pattern analysis showing trends (improving/declining based on positive vs negative)
- [x] Summary statistics: total incidents by type, this month vs last month

### AC4: Behavior Tab in Student Profile
- [x] Tab shows behavior summary with trend indicator
- [x] Quick add incident button
- [x] Recent incidents list (last 10)
- [x] Link to full behavioral history

---

## Technical Requirements

### Backend

**Database Tables:**

```sql
-- Migration: 20260126100300_create_behavioral_tracking.up.sql

CREATE TYPE incident_type AS ENUM ('positive_recognition', 'minor_infraction', 'major_violation');
CREATE TYPE incident_severity AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE follow_up_status AS ENUM ('pending', 'completed', 'cancelled');

CREATE TABLE behavioral_incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    incident_type incident_type NOT NULL,
    severity incident_severity NOT NULL DEFAULT 'medium',
    incident_date DATE NOT NULL,
    incident_time TIME NOT NULL,
    location VARCHAR(200),
    description TEXT NOT NULL,
    witnesses TEXT, -- JSON array of names
    student_response TEXT,
    action_taken TEXT NOT NULL,
    parent_meeting_required BOOLEAN NOT NULL DEFAULT FALSE,
    reported_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE behavioral_incidents ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON behavioral_incidents USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_behavioral_student ON behavioral_incidents(student_id);
CREATE INDEX idx_behavioral_date ON behavioral_incidents(incident_date DESC);
CREATE INDEX idx_behavioral_type ON behavioral_incidents(incident_type);

CREATE TABLE incident_follow_ups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    incident_id UUID NOT NULL REFERENCES behavioral_incidents(id) ON DELETE CASCADE,
    scheduled_date DATE NOT NULL,
    scheduled_time TIME,
    participants TEXT, -- JSON array: [{name, role}]
    expected_outcomes TEXT,
    meeting_notes TEXT,
    actual_outcomes TEXT,
    status follow_up_status NOT NULL DEFAULT 'pending',
    completed_at TIMESTAMPTZ,
    completed_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

ALTER TABLE incident_follow_ups ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON incident_follow_ups USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_follow_ups_incident ON incident_follow_ups(incident_id);
CREATE INDEX idx_follow_ups_date ON incident_follow_ups(scheduled_date) WHERE status = 'pending';
```

**API Endpoints:**

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/students/{id}/behavioral-incidents` | List incidents | `student:read` |
| POST | `/api/v1/students/{id}/behavioral-incidents` | Record incident | `behavior:write` |
| GET | `/api/v1/students/{id}/behavioral-incidents/{bid}` | Get incident details | `student:read` |
| PUT | `/api/v1/students/{id}/behavioral-incidents/{bid}` | Update incident | `behavior:write` |
| GET | `/api/v1/students/{id}/behavioral-summary` | Get behavior statistics | `student:read` |
| POST | `/api/v1/behavioral-incidents/{bid}/follow-ups` | Schedule follow-up | `behavior:write` |
| PUT | `/api/v1/behavioral-incidents/{bid}/follow-ups/{fid}` | Update follow-up | `behavior:write` |
| GET | `/api/v1/follow-ups/pending` | List pending follow-ups | `behavior:read` |

**Behavior Summary Response:**

```go
type BehaviorSummary struct {
    TotalIncidents       int     `json:"totalIncidents"`
    PositiveCount        int     `json:"positiveCount"`
    MinorInfractionCount int     `json:"minorInfractionCount"`
    MajorViolationCount  int     `json:"majorViolationCount"`
    ThisMonthCount       int     `json:"thisMonthCount"`
    LastMonthCount       int     `json:"lastMonthCount"`
    Trend                string  `json:"trend"` // improving, declining, stable
    PendingFollowUps     int     `json:"pendingFollowUps"`
}

func (s *Service) CalculateTrend(summary BehaviorSummary) string {
    positiveRatio := float64(summary.PositiveCount) / float64(summary.TotalIncidents)
    if positiveRatio > 0.7 {
        return "improving"
    } else if positiveRatio < 0.3 {
        return "declining"
    }
    return "stable"
}
```

### Frontend

**Components to create:**

```bash
ng generate component features/students/components/behavior-summary --standalone
ng generate component features/students/components/incident-form --standalone
ng generate component features/students/components/incident-list --standalone
ng generate component features/students/components/follow-up-form --standalone
ng generate interface features/students/models/behavioral-incident
ng generate interface features/students/models/follow-up
```

**Behavior Summary Component:**

```typescript
@Component({
  selector: 'app-behavior-summary',
  template: `
    <div class="grid grid-cols-4 gap-4">
      <app-stat-card
        title="Total Incidents"
        [value]="summary().totalIncidents"
        icon="clipboard-list"
      />
      <app-stat-card
        title="Positive Recognition"
        [value]="summary().positiveCount"
        icon="star"
        variant="success"
      />
      <app-stat-card
        title="Infractions"
        [value]="summary().minorInfractionCount + summary().majorViolationCount"
        icon="exclamation-triangle"
        variant="warning"
      />
      <app-stat-card
        title="Trend"
        [value]="summary().trend"
        [icon]="trendIcon()"
        [variant]="trendVariant()"
      />
    </div>
  `
})
export class BehaviorSummaryComponent {
  summary = input.required<BehaviorSummary>();

  trendIcon = computed(() => {
    switch (this.summary().trend) {
      case 'improving': return 'arrow-up';
      case 'declining': return 'arrow-down';
      default: return 'minus';
    }
  });
}
```

---

## Tasks

### Backend Tasks

- [ ] **BE-4.4.1**: Create behavioral incident entity and migration
- [ ] **BE-4.4.2**: Create follow-up entity and migration
- [ ] **BE-4.4.3**: Create behavioral incident repository
- [ ] **BE-4.4.4**: Create behavioral service with summary calculation
- [ ] **BE-4.4.5**: Create behavioral HTTP handlers
- [ ] **BE-4.4.6**: Add behavior permissions to seed
- [ ] **BE-4.4.7**: Write unit tests

### Frontend Tasks

- [ ] **FE-4.4.1**: Create behavioral interfaces
- [ ] **FE-4.4.2**: Create behavioral service
- [ ] **FE-4.4.3**: Create behavior summary component
- [ ] **FE-4.4.4**: Create incident form component
- [ ] **FE-4.4.5**: Create incident list with filters
- [ ] **FE-4.4.6**: Create follow-up form component
- [ ] **FE-4.4.7**: Add Behavior tab to student profile
- [ ] **FE-4.4.8**: Write component tests

---

## Definition of Done

- [ ] All acceptance criteria verified
- [ ] Incident types working correctly
- [ ] Follow-up scheduling working
- [ ] Trend calculation accurate
- [ ] Backend tests passing
- [ ] Frontend tests passing

---

## Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Story 4.1 (Student Profile) | Required | Incidents link to student |
| Notification system (Epic 12) | Future | Follow-up notifications deferred |

---

## Notes

- Positive recognition encourages good behavior tracking
- Trend analysis helps identify students needing intervention
- Follow-up notifications will be implemented in Epic 12
- Consider privacy: who can see behavioral records?
