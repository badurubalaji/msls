# Story 3.3: Admission Session Configuration

**Epic:** 3 - School Setup & Admissions
**Status:** done
**Priority:** High
**Estimated Effort:** Medium

---

## User Story

As an **administrator**,
I want **to configure admission sessions for different classes**,
So that **admissions can be processed in organized cycles**.

---

## Acceptance Criteria

### AC1: Create Admission Session
**Given** an admin is on admission settings
**When** they create an admission session
**Then** they can enter: name, academic year, start/end dates
**And** they can select applicable classes for this session
**And** they can set maximum seats per class
**And** they can configure required documents list
**And** they can set admission fee amount

### AC2: Manage Session
**Given** an admission session exists
**When** viewing the session
**Then** they see: applications count, seats filled, available seats
**And** they can open/close the session
**And** they can extend deadline if needed

---

## Technical Requirements

### Backend (Go)

#### Database Schema

```sql
-- Admission Sessions
CREATE TABLE IF NOT EXISTS admission_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'upcoming',
    application_fee DECIMAL(10,2) DEFAULT 0,
    required_documents JSONB DEFAULT '[]',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Admission Seats (per class)
CREATE TABLE IF NOT EXISTS admission_seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id) ON DELETE CASCADE,
    class_name VARCHAR(50) NOT NULL,
    total_seats INT NOT NULL DEFAULT 0,
    filled_seats INT NOT NULL DEFAULT 0,
    waitlist_limit INT DEFAULT 10,
    reserved_seats JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | /api/v1/admission-sessions | List sessions | admissions:read |
| GET | /api/v1/admission-sessions/:id | Get by ID | admissions:read |
| POST | /api/v1/admission-sessions | Create session | admissions:create |
| PUT | /api/v1/admission-sessions/:id | Update session | admissions:update |
| PATCH | /api/v1/admission-sessions/:id/status | Change status | admissions:update |
| DELETE | /api/v1/admission-sessions/:id | Delete session | admissions:delete |
| GET | /api/v1/admission-sessions/:id/seats | List seats config | admissions:read |
| POST | /api/v1/admission-sessions/:id/seats | Add seat config | admissions:update |
| PUT | /api/v1/admission-sessions/:id/seats/:seatId | Update seats | admissions:update |

### Frontend (Angular)

#### Feature Structure

```
src/app/features/admissions/
├── admissions.routes.ts
├── sessions/
│   ├── admission-session.model.ts
│   ├── admission-session.service.ts
│   ├── sessions.component.ts
│   ├── session-form.component.ts
│   └── seat-config.component.ts
```

---

## Definition of Done

- [x] Backend: Migration and endpoints
- [x] Backend: Unit tests
- [x] Frontend: Session list and form
- [x] Frontend: Seat configuration UI
- [x] Navigation link added

## Implementation Notes

### Backend
- **Session Service**: `/home/ashulabs/workspace/msls/msls-backend/internal/services/admission/session_service.go`
- **Session Handler**: `/home/ashulabs/workspace/msls/msls-backend/internal/handlers/admission/session_handler.go`
- **Session DTOs**: `/home/ashulabs/workspace/msls/msls-backend/internal/handlers/admission/session_dto.go`
- **Unit Tests**: `/home/ashulabs/workspace/msls/msls-backend/internal/services/admission/session_service_test.go`
- **Migration**: `/home/ashulabs/workspace/msls/msls-backend/migrations/000015_admission_sessions.up.sql`

### Frontend
- **Models**: `/home/ashulabs/workspace/msls/msls-frontend/src/app/features/admissions/sessions/admission-session.model.ts`
- **Service**: `/home/ashulabs/workspace/msls/msls-frontend/src/app/features/admissions/sessions/admission-session.service.ts`
- **Sessions Page**: `/home/ashulabs/workspace/msls/msls-frontend/src/app/features/admissions/sessions/sessions.component.ts`
- **Session Form**: `/home/ashulabs/workspace/msls/msls-frontend/src/app/features/admissions/sessions/session-form.component.ts`
- **Seat Config**: `/home/ashulabs/workspace/msls/msls-frontend/src/app/features/admissions/sessions/seat-config.component.ts`
- **Routes**: `/home/ashulabs/workspace/msls/msls-frontend/src/app/features/admissions/admissions.routes.ts`
- **Navigation**: `/home/ashulabs/workspace/msls/msls-frontend/src/app/layouts/nav-config.ts`

### Key Features
- Full CRUD for admission sessions with status management (upcoming/open/closed)
- Per-class seat configuration with total seats, filled seats, and waitlist limits
- Required documents configuration
- Application fee setting
- Session statistics (applications count, seats filled, available seats)
- Deadline extension capability
- Status change workflow validation
