# Story 3.5: Online Admission Application

**Epic:** 3 - School Setup & Admissions
**Status:** done
**Priority:** High
**Estimated Effort:** Large

---

## User Story

As a **parent**,
I want **to submit an admission application online**,
So that **I don't need to visit the school for initial application**.

---

## Acceptance Criteria

### AC1: Submit Application
**Given** a parent accesses the admission portal
**When** they select class and fill the application form
**Then** they can enter: student details (name, DOB, gender, Aadhaar)
**And** they can enter: parent/guardian details
**And** they can enter: previous school details (if applicable)
**And** they can upload required documents (birth certificate, photos, etc.)
**And** application is saved with unique application number

### AC2: Application Confirmation
**Given** application is submitted
**When** all required fields and documents are provided
**Then** application status is "submitted"
**And** confirmation email/SMS is sent to parent
**And** application appears in admin dashboard for review

### AC3: Check Application Status
**Given** a parent wants to check application status
**When** they login with application number and phone
**Then** they see current status and any remarks
**And** they can upload additional documents if requested

---

## Technical Requirements

### Backend (Go)

#### Database Schema

```sql
-- Admission Applications
CREATE TABLE IF NOT EXISTS admission_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES admission_sessions(id),
    enquiry_id UUID REFERENCES admission_enquiries(id),
    application_number VARCHAR(50) NOT NULL,
    current_stage VARCHAR(50) DEFAULT 'submitted',
    stage_history JSONB DEFAULT '[]',

    -- Student Details
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(20) NOT NULL,
    blood_group VARCHAR(10),
    nationality VARCHAR(50) DEFAULT 'Indian',
    religion VARCHAR(50),
    category VARCHAR(50),
    aadhaar_number VARCHAR(12),
    photo_url VARCHAR(500),

    -- Previous School
    previous_school VARCHAR(255),
    previous_class VARCHAR(50),
    previous_percentage DECIMAL(5,2),
    transfer_certificate_url VARCHAR(500),

    -- Contact
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),

    -- Metadata
    submitted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(tenant_id, application_number)
);

-- Application Parents
CREATE TABLE IF NOT EXISTS application_parents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    relation VARCHAR(20) NOT NULL,
    name VARCHAR(200) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    occupation VARCHAR(100),
    education VARCHAR(100),
    annual_income VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Application Documents
CREATE TABLE IF NOT EXISTS application_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES admission_applications(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | /api/v1/applications | List applications | applications:read |
| GET | /api/v1/applications/:id | Get by ID | applications:read |
| POST | /api/v1/applications | Create application | applications:create |
| PUT | /api/v1/applications/:id | Update application | applications:update |
| PATCH | /api/v1/applications/:id/submit | Submit application | applications:update |
| PATCH | /api/v1/applications/:id/stage | Update stage | applications:update |
| POST | /api/v1/applications/:id/documents | Upload document | applications:update |
| GET | /api/v1/applications/:id/documents | List documents | applications:read |
| DELETE | /api/v1/applications/:id/documents/:docId | Delete document | applications:update |
| POST | /api/v1/applications/:id/parents | Add parent | applications:update |
| PUT | /api/v1/applications/:id/parents/:parentId | Update parent | applications:update |
| GET | /api/v1/applications/status | Check status (public) | - |

### Frontend (Angular)

#### Feature Structure

```
src/app/features/admissions/applications/
├── application.model.ts
├── application.service.ts
├── applications.component.ts
├── application-form.component.ts
├── application-review.component.ts
└── public/
    └── application-status.component.ts
```

---

## Definition of Done

- [x] Backend: All tables and endpoints
  - Created admission_applications, application_parents, application_documents tables (migration exists)
  - Implemented ApplicationService with CRUD operations
  - Implemented ApplicationHandler with all HTTP endpoints
  - Created DTOs for request/response serialization
- [x] Backend: Document upload to storage
  - AddDocument endpoint created with file metadata support
  - Document verification/rejection workflow implemented
- [x] Backend: Application number generation
  - Auto-generated format: APP-YYYYMMDD-XXXX (e.g., APP-20260123-0001)
  - ApplicationNumberSequence table tracks daily sequences
- [x] Frontend: Application list for admin
  - ApplicationsComponent with filtering, search, and pagination
  - Status badges, stage history display
- [x] Frontend: Multi-step application form
  - ApplicationFormComponent with tabs for student, parents, documents
  - Draft save and submit functionality
- [x] Frontend: Document upload UI
  - Document type selection and file upload
  - Verification status display
- [x] Frontend: Public status check page
  - PublicStatusCheck component at /check-status
  - No authentication required
  - Validates application number + phone number
