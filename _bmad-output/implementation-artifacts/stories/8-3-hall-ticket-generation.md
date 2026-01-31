# Story 8.3: Hall Ticket Generation

Status: review

## Story

As an **administrator**,
I want **to generate hall tickets for exams**,
So that **students have proper exam identification**.

## Acceptance Criteria

### AC1: Generate Hall Tickets
**Given** an exam is scheduled (status = 'scheduled')
**When** generating hall tickets
**Then** they can select: exam, class, section (or all)
**And** hall tickets are generated in batch
**And** each ticket has: unique roll number, student photo, exam schedule

### AC2: Hall Ticket Template
**Given** hall ticket template exists
**When** generating tickets
**Then** template is used with school branding
**And** exam-wise schedule is printed
**And** important instructions are included
**And** QR code for verification is included

### AC3: Distribute Hall Tickets
**Given** hall tickets are generated
**When** distributing
**Then** bulk print option is available
**And** PDF download for individual or batch
**And** student portal shows downloadable hall ticket

## Tasks / Subtasks

### Backend Tasks

- [x] **Task 1: Database Schema** (AC: 1,2,3)
  - [x] Create `hall_tickets` table: id, tenant_id, examination_id, student_id, roll_number, qr_code_data, status (generated/printed/downloaded), generated_at
  - [x] Create `hall_ticket_templates` table: id, tenant_id, name, content (JSON), is_default, header_logo_url, instructions
  - [x] Add migration file with RLS policies and indexes
  - [x] Add unique constraint on (examination_id, student_id)

- [x] **Task 2: Model & Repository** (AC: 1,2,3)
  - [x] Create `HallTicket` model with relationships
  - [x] Create `HallTicketTemplate` model
  - [x] Implement repository for hall ticket CRUD
  - [x] Implement batch generation queries
  - [x] Add roll number generation logic (per exam, sequential or custom pattern)

- [x] **Task 3: PDF Generation Service** (AC: 1,2)
  - [x] Create `HallTicketPDFGenerator` using `github.com/go-pdf/fpdf` (existing pattern)
  - [x] Include student photo from `Student.PhotoURL`
  - [x] Include exam schedule from `exam_schedules` table
  - [x] Generate QR code using `github.com/skip2/go-qrcode`
  - [x] Include school branding (logo, name, address)
  - [x] Add instructions section from template

- [x] **Task 4: Service Layer** (AC: 1,2,3)
  - [x] Create hall ticket service with business logic
  - [x] Validate exam is in 'scheduled' status before generation
  - [x] Generate roll numbers (configurable pattern: YYYY-CLASS-SEQ or custom)
  - [x] Batch generate hall tickets for class/section
  - [x] Generate QR code with verification data (ticket_id, student_id, exam_id hash)

- [x] **Task 5: API Endpoints** (AC: 1,2,3)
  - [x] `GET /api/v1/examinations/:id/hall-tickets` - List hall tickets for exam
  - [x] `POST /api/v1/examinations/:id/hall-tickets/generate` - Generate hall tickets (class_id optional)
  - [x] `GET /api/v1/examinations/:id/hall-tickets/:ticketId/pdf` - Download single PDF
  - [x] `GET /api/v1/examinations/:id/hall-tickets/pdf` - Download batch PDF
  - [x] `GET /api/v1/hall-tickets/verify/:qrCode` - Verify hall ticket (public endpoint)
  - [x] `GET /api/v1/hall-ticket-templates` - List templates
  - [x] `POST /api/v1/hall-ticket-templates` - Create template
  - [x] `PUT /api/v1/hall-ticket-templates/:id` - Update template

- [x] **Task 6: Permissions** (AC: 1,2,3)
  - [x] Add permissions: `hall-ticket:view`, `hall-ticket:generate`, `hall-ticket:download`, `hall-ticket:template-manage`

### Frontend Tasks

- [x] **Task 7: Model & Service** (AC: 1,2,3)
  - [x] Create `HallTicket` and `HallTicketTemplate` interfaces
  - [x] Create `HallTicketService` with API methods (added to ExamService)
  - [x] Add download helpers for PDF blob handling

- [x] **Task 8: Hall Ticket List Component** (AC: 1,3)
  - [x] Create list page showing hall tickets for an exam
  - [x] Show student name, roll number, class, status
  - [x] Add filters: class, section, generation status
  - [x] Add bulk download button
  - [x] Add individual download buttons

- [x] **Task 9: Hall Ticket Generation Component** (AC: 1,2)
  - [x] Create generation dialog/page (modal in list component)
  - [x] Select exam (only 'scheduled' exams) - navigated from scheduled exams
  - [x] Select class/section (or all)
  - [x] Show preview of roll number pattern
  - [x] Generate button with progress indicator
  - [x] Show generation results (success/failed count)

- [x] **Task 10: Hall Ticket Template Component** (AC: 2)
  - [x] Create template management page
  - [x] Form for template details (name, instructions, logo upload)
  - [ ] Preview of template layout (not implemented - optional enhancement)
  - [x] Set default template

- [x] **Task 11: Routing & Navigation** (AC: 1,2,3)
  - [x] Add hall ticket routes under /exams/:id/hall-tickets
  - [x] Add template routes under /exams/hall-ticket-templates
  - [x] Add navigation links from examination detail page
  - [ ] Add permission guards (uses existing auth middleware)

## Dev Notes

### Database Design

```sql
-- Hall Ticket Templates
CREATE TABLE hall_ticket_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    header_logo_url VARCHAR(500),
    school_name VARCHAR(200),
    school_address TEXT,
    instructions TEXT,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Hall Tickets
CREATE TABLE hall_tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    examination_id UUID NOT NULL REFERENCES examinations(id),
    student_id UUID NOT NULL REFERENCES students(id),
    roll_number VARCHAR(50) NOT NULL,
    qr_code_data VARCHAR(500) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'generated',
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    printed_at TIMESTAMPTZ,
    downloaded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_hall_ticket_status CHECK (status IN ('generated', 'printed', 'downloaded')),
    CONSTRAINT uq_hall_ticket_exam_student UNIQUE (examination_id, student_id)
);

-- RLS Policies
ALTER TABLE hall_ticket_templates ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_hall_ticket_templates ON hall_ticket_templates
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

ALTER TABLE hall_tickets ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation_hall_tickets ON hall_tickets
    USING (tenant_id = current_setting('app.tenant_id', true)::UUID);

-- Indexes
CREATE INDEX idx_hall_tickets_tenant ON hall_tickets(tenant_id);
CREATE INDEX idx_hall_tickets_examination ON hall_tickets(examination_id);
CREATE INDEX idx_hall_tickets_student ON hall_tickets(student_id);
CREATE INDEX idx_hall_ticket_templates_tenant ON hall_ticket_templates(tenant_id);
```

### PDF Generation Pattern

Follow existing pattern from `msls-backend/internal/modules/payroll/pdf.go`:
- Use `github.com/go-pdf/fpdf` for PDF generation
- Create `HallTicketPDFGenerator` struct
- Include school branding, student info, exam schedule, QR code
- Support both single and batch PDF generation

### QR Code Data Format

```json
{
  "t": "hall_ticket_id",
  "s": "student_id_short",
  "e": "exam_id_short",
  "v": "verification_hash"
}
```

### Roll Number Pattern

Default: `{YEAR}-{CLASS_CODE}-{SEQUENCE}`
Example: `2026-10A-001`, `2026-10A-002`

### Related Entities

- `examinations` (8.2) - Get exam details and schedules
- `exam_schedules` (8.2) - Subject-wise exam schedule for hall ticket
- `students` - Student name, photo, admission number
- `classes` - Class name for display
- `examination_classes` - Link exam to applicable classes

### Project Structure Notes

**Backend files to create:**
```
msls-backend/internal/modules/hallticket/
├── handler.go
├── service.go
├── repository.go
├── dto.go
├── errors.go
├── pdf.go  (PDF generation)
└── qrcode.go  (QR code generation)
```

**Frontend files to create:**
```
msls-frontend/src/app/features/exams/
├── hall-ticket.model.ts
├── hall-ticket.service.ts
├── hall-ticket-list/
│   ├── hall-ticket-list.ts
│   ├── hall-ticket-list.html
│   └── hall-ticket-list.scss
├── hall-ticket-generate/
│   ├── hall-ticket-generate.ts
│   ├── hall-ticket-generate.html
│   └── hall-ticket-generate.scss
└── hall-ticket-template/
    ├── hall-ticket-template.ts
    ├── hall-ticket-template.html
    └── hall-ticket-template.scss
```

### Dependencies

**Backend (go.mod):**
```
github.com/go-pdf/fpdf (already exists - used by payroll)
github.com/skip2/go-qrcode (new - for QR code generation)
```

**Frontend:**
- No new dependencies needed

### Previous Story Learnings (from 8.2)

1. Examination module exists with full CRUD
2. `exam_schedules` table has subject-wise schedule
3. `examination_classes` links exams to classes
4. Permissions pattern: `exam:view`, `exam:create`, etc.
5. Frontend uses signals for state management
6. ExaminationService already exists for API calls

### References

- [Source: architecture.md#PDF-Generation] - Use fpdf pattern from payroll module
- [Source: 8-2-examination-scheduling.md] - Examination and schedule models
- [Source: models/student.go] - Student model with PhotoURL field
- [Source: payroll/pdf.go] - PDF generation pattern to follow
- [Source: project-context.md] - API design and naming conventions

## How It Works (Manual Testing)

### Prerequisites
1. Backend running: `cd msls-backend && go run cmd/api/main.go`
2. Frontend running: `cd msls-frontend && npm start`
3. Login as admin at http://localhost:4200
4. An examination in "scheduled" status (see Story 8.2)
5. Students enrolled in the exam's classes

### Seed Data
Run seed data to create test examination with schedules:
```bash
cd msls-backend
PGPASSWORD=postgres psql -h localhost -U postgres -d msls -f scripts/seed_exam_data.sql
```
This creates:
- Hall ticket template: "Default Template"
- Examination: "Mid Term Examination - LKG" (scheduled status)
- 3 exam schedules: English, Hindi, Mathematics

### Test Scenarios

#### 1. Manage Hall Ticket Templates
- Navigate to **Exams → Hall Ticket Templates** in sidebar
- See list of templates
- Click **+ Add Template** to create new
  - Enter: Name, School Name, Address, Instructions
  - Set as default template
- Edit or delete existing templates

#### 2. Generate Hall Tickets
- Navigate to **Exams → Examinations**
- Find a "Scheduled" or "Ongoing" exam
- Click **Hall Tickets** button on the exam row
- Click **Generate Hall Tickets** button
- Select class filter (optional)
- Click **Generate**
- See progress and results (success/failed count)

#### 3. View Hall Tickets
- Hall ticket list shows:
  - Student name
  - Roll number (auto-generated: YYYY-CLASS-SEQ)
  - Class and section
  - Status (Generated/Downloaded)
- Filter by class or status

#### 4. Download Hall Tickets
- **Individual**: Click download icon on a row → PDF downloaded
- **Bulk**: Click **Download All** button → Combined PDF

#### 5. Verify Hall Ticket (Public)
- Each hall ticket has a QR code
- Scan QR code or visit `/api/v1/hall-tickets/verify/{qr_code}`
- Returns verification status (valid/invalid)

### PDF Contents
The generated hall ticket PDF includes:
- School name and logo
- Student photo (if available)
- Student details: Name, Roll Number, Class
- Examination name and dates
- Subject-wise schedule (date, time, venue)
- Instructions from template
- QR code for verification
- Footer with generation date

### API Testing (curl)
```bash
# List hall tickets for exam
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/examinations/{exam_id}/hall-tickets

# Generate hall tickets
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"class_id":"..."}' \
  http://localhost:8080/api/v1/examinations/{exam_id}/hall-tickets/generate

# Download single hall ticket PDF
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/examinations/{exam_id}/hall-tickets/{ticket_id}/pdf -o ticket.pdf

# Download batch PDF
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/examinations/{exam_id}/hall-tickets/pdf?class_id=..." -o batch.pdf

# Verify hall ticket (public endpoint)
curl http://localhost:8080/api/v1/hall-tickets/verify/{qr_code_data}

# Hall ticket templates
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/hall-ticket-templates
```

### Roll Number Format
Default pattern: `{YEAR}-{CLASS_CODE}-{SEQUENCE}`
Example: `2026-LKG-001`, `2026-LKG-002`

### QR Code Verification
QR code contains encrypted data with:
- Hall ticket ID
- Student ID (shortened)
- Examination ID (shortened)
- Verification hash (SHA-256)

The verification endpoint validates the hash to prevent tampering.

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Implementation Plan

Followed the task list in order:
1. Created database migration with hall_tickets and hall_ticket_templates tables
2. Implemented backend module (hallticket) with full CRUD operations
3. Integrated PDF generation with QR codes
4. Added frontend models, service methods, and components

### Completion Notes

Story 8.3 Hall Ticket Generation has been implemented with the following features:

**Backend:**
- Migration 000060 with hall_tickets and hall_ticket_templates tables
- Full hallticket module with handler, service, repository, dto, errors
- PDF generation with school branding, exam schedule, and QR codes
- QR code generation and verification for authentication
- Permissions: hall-ticket:view, hall-ticket:generate, hall-ticket:download, hall-ticket:template-manage

**Frontend:**
- HallTicket and HallTicketTemplate interfaces added to exam.model.ts
- ExamService extended with hall ticket API methods
- HallTicketListComponent for viewing/managing hall tickets per exam
- HallTicketTemplateComponent for managing templates
- Routes added for /exams/:id/hall-tickets and /exams/hall-ticket-templates
- Navigation link added to examination list (for scheduled/ongoing exams)

**Not Implemented (Optional Enhancements):**
- Template layout preview in frontend
- Dedicated permission guards (uses existing auth middleware)

### Change Log
| Date | Change | Author |
|------|--------|--------|
| 2026-01-31 | Story created with comprehensive context | create-story workflow |
| 2026-01-31 | Backend implementation complete | dev agent |
| 2026-01-31 | Frontend implementation complete | dev agent |
| 2026-01-31 | Story marked for review | dev agent |

### File List

**Backend Files Created:**
- msls-backend/migrations/000060_hall_tickets.up.sql
- msls-backend/migrations/000060_hall_tickets.down.sql
- msls-backend/internal/modules/hallticket/errors.go
- msls-backend/internal/modules/hallticket/dto.go
- msls-backend/internal/modules/hallticket/repository.go
- msls-backend/internal/modules/hallticket/qrcode.go
- msls-backend/internal/modules/hallticket/pdf.go
- msls-backend/internal/modules/hallticket/service.go
- msls-backend/internal/modules/hallticket/handler.go

**Backend Files Modified:**
- msls-backend/cmd/api/main.go (added hallticket handler registration)

**Frontend Files Created:**
- msls-frontend/src/app/features/exams/hall-ticket-list/hall-ticket-list.component.ts
- msls-frontend/src/app/features/exams/hall-ticket-template/hall-ticket-template.component.ts

**Frontend Files Modified:**
- msls-frontend/src/app/features/exams/exam.model.ts (added hall ticket models)
- msls-frontend/src/app/features/exams/exam.service.ts (added hall ticket API methods)
- msls-frontend/src/app/features/exams/exams.routes.ts (added hall ticket routes)
- msls-frontend/src/app/features/exams/examination-list/examination-list.component.ts (added hall ticket nav link)
