# Story 5.8: Staff Document Management

**Epic:** 5 - Staff Management
**Status:** review
**Priority:** P2 (Required for HR compliance)
**Story Points:** 8

---

## User Story

**As an** HR administrator,
**I want** to manage staff documents with expiry tracking,
**So that** compliance requirements are met and documents are always up-to-date.

---

## Acceptance Criteria

### AC1: Document Upload

**Given** a staff member's documents section
**When** uploading a document
**Then** they can select type: Aadhaar, PAN, offer letter, contract, ID proof, education certificate
**And** they can enter: document number, issue date, expiry date (if applicable)
**And** they can upload file (PDF, images - JPG, PNG)
**And** file size limit is enforced (max 5MB per file)
**And** document is stored securely with encryption

### AC2: Document Expiry Tracking

**Given** documents have expiry dates
**When** expiry approaches
**Then** notification is sent 30 days before expiry
**And** notification is sent 7 days before expiry
**And** dashboard shows expiring documents count
**And** expired documents are flagged with visual indicator

### AC3: Document Verification

**Given** a document is uploaded
**When** admin verifies the document
**Then** they can mark status: pending, verified, rejected
**And** verification details are recorded (verified by, date, notes)
**And** rejection reason is captured if rejected
**And** staff is notified of verification status

### AC4: Document Listing & Search

**Given** HR is viewing staff documents
**When** searching or filtering
**Then** they can filter by: document type, status, expiry date range
**And** they can search by: staff name, employee ID, document number
**And** bulk download of documents is available
**And** document compliance report can be generated

### AC5: Document Categories & Types

**Given** admin is configuring document types
**When** managing document types
**Then** they can create custom document types
**And** they can set: name, category, is_mandatory, has_expiry, validity_period
**And** they can define required documents per staff type (teaching/non-teaching)

### AC6: Compliance Dashboard

**Given** HR is on compliance dashboard
**When** viewing document compliance
**Then** they see: total staff, documents submitted, pending verification, expired
**And** they see: staff with missing mandatory documents
**And** they see: upcoming expirations (next 30/60/90 days)
**And** compliance percentage per document type is shown

---

## Technical Implementation

### Database Schema

```sql
-- Migration: 000046_staff_documents.up.sql

-- Document Types Configuration
CREATE TABLE staff_document_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),

    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL, -- identity, education, employment, other
    description TEXT,

    is_mandatory BOOLEAN NOT NULL DEFAULT false,
    has_expiry BOOLEAN NOT NULL DEFAULT false,
    default_validity_months INTEGER, -- Default validity period if has_expiry

    applicable_to VARCHAR(20)[] DEFAULT ARRAY['teaching', 'non_teaching'], -- Staff types
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_document_type_code UNIQUE (tenant_id, code)
);

-- Staff Documents
CREATE TABLE staff_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    document_type_id UUID NOT NULL REFERENCES staff_document_types(id),

    document_number VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,

    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL, -- S3/MinIO path
    file_size INTEGER NOT NULL, -- in bytes
    mime_type VARCHAR(100) NOT NULL,

    -- Verification
    verification_status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (verification_status IN ('pending', 'verified', 'rejected')),
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    verification_notes TEXT,
    rejection_reason TEXT,

    -- Metadata
    remarks TEXT,
    is_current BOOLEAN NOT NULL DEFAULT true, -- For document versions

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Enable RLS
ALTER TABLE staff_document_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE staff_documents ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_document_types ON staff_document_types
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

CREATE POLICY tenant_isolation_staff_documents ON staff_documents
    USING (
        current_setting('app.tenant_id', true) IS NULL
        OR current_setting('app.tenant_id', true) = ''
        OR tenant_id = current_setting('app.tenant_id', true)::UUID
    );

-- Indexes
CREATE INDEX idx_staff_documents_tenant ON staff_documents(tenant_id);
CREATE INDEX idx_staff_documents_staff ON staff_documents(staff_id);
CREATE INDEX idx_staff_documents_type ON staff_documents(document_type_id);
CREATE INDEX idx_staff_documents_expiry ON staff_documents(expiry_date) WHERE expiry_date IS NOT NULL;
CREATE INDEX idx_staff_documents_status ON staff_documents(verification_status);
CREATE INDEX idx_document_types_tenant ON staff_document_types(tenant_id);

-- Document Expiry Notifications (audit trail)
CREATE TABLE staff_document_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    document_id UUID NOT NULL REFERENCES staff_documents(id) ON DELETE CASCADE,

    notification_type VARCHAR(50) NOT NULL, -- expiry_30_days, expiry_7_days, expired
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_to UUID[] NOT NULL, -- User IDs notified

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed default document types
INSERT INTO staff_document_types (id, tenant_id, name, code, category, is_mandatory, has_expiry, applicable_to)
SELECT
    uuid_generate_v7(),
    t.id,
    dt.name,
    dt.code,
    dt.category,
    dt.is_mandatory,
    dt.has_expiry,
    dt.applicable_to
FROM tenants t
CROSS JOIN (
    VALUES
        ('Aadhaar Card', 'aadhaar', 'identity', true, false, ARRAY['teaching', 'non_teaching']),
        ('PAN Card', 'pan', 'identity', true, false, ARRAY['teaching', 'non_teaching']),
        ('Passport', 'passport', 'identity', false, true, ARRAY['teaching', 'non_teaching']),
        ('Driving License', 'driving_license', 'identity', false, true, ARRAY['teaching', 'non_teaching']),
        ('10th Marksheet', 'ssc_marksheet', 'education', true, false, ARRAY['teaching', 'non_teaching']),
        ('12th Marksheet', 'hsc_marksheet', 'education', true, false, ARRAY['teaching', 'non_teaching']),
        ('Degree Certificate', 'degree_certificate', 'education', false, false, ARRAY['teaching']),
        ('B.Ed Certificate', 'bed_certificate', 'education', false, false, ARRAY['teaching']),
        ('Offer Letter', 'offer_letter', 'employment', true, false, ARRAY['teaching', 'non_teaching']),
        ('Appointment Letter', 'appointment_letter', 'employment', true, false, ARRAY['teaching', 'non_teaching']),
        ('Experience Letter', 'experience_letter', 'employment', false, false, ARRAY['teaching', 'non_teaching']),
        ('Police Verification', 'police_verification', 'compliance', true, true, ARRAY['teaching', 'non_teaching']),
        ('Medical Fitness', 'medical_fitness', 'compliance', false, true, ARRAY['teaching', 'non_teaching'])
) AS dt(name, code, category, is_mandatory, has_expiry, applicable_to)
ON CONFLICT DO NOTHING;

-- Add permissions
INSERT INTO permissions (code, name, description, module, created_at, updated_at)
VALUES
    ('staff_document:view', 'View Staff Documents', 'Permission to view staff documents', 'staff_document', NOW(), NOW()),
    ('staff_document:upload', 'Upload Staff Documents', 'Permission to upload documents', 'staff_document', NOW(), NOW()),
    ('staff_document:verify', 'Verify Staff Documents', 'Permission to verify/reject documents', 'staff_document', NOW(), NOW()),
    ('staff_document:delete', 'Delete Staff Documents', 'Permission to delete documents', 'staff_document', NOW(), NOW()),
    ('staff_document:download', 'Download Staff Documents', 'Permission to download documents', 'staff_document', NOW(), NOW()),
    ('staff_document:manage_types', 'Manage Document Types', 'Permission to manage document type configuration', 'staff_document', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;
```

### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/staff/:id/documents` | List staff documents | `staff_document:view` |
| POST | `/api/v1/staff/:id/documents` | Upload document | `staff_document:upload` |
| GET | `/api/v1/staff/:id/documents/:docId` | Get document details | `staff_document:view` |
| PUT | `/api/v1/staff/:id/documents/:docId` | Update document metadata | `staff_document:upload` |
| DELETE | `/api/v1/staff/:id/documents/:docId` | Delete document | `staff_document:delete` |
| GET | `/api/v1/staff/:id/documents/:docId/download` | Download document file | `staff_document:download` |
| PUT | `/api/v1/staff/:id/documents/:docId/verify` | Verify document | `staff_document:verify` |
| PUT | `/api/v1/staff/:id/documents/:docId/reject` | Reject document | `staff_document:verify` |
| GET | `/api/v1/documents/expiring` | List expiring documents | `staff_document:view` |
| GET | `/api/v1/documents/compliance` | Get compliance report | `staff_document:view` |
| GET | `/api/v1/document-types` | List document types | `staff_document:view` |
| POST | `/api/v1/document-types` | Create document type | `staff_document:manage_types` |
| PUT | `/api/v1/document-types/:id` | Update document type | `staff_document:manage_types` |
| DELETE | `/api/v1/document-types/:id` | Delete document type | `staff_document:manage_types` |

---

## Tasks

### Backend Tasks

- [x] **BE-5.8.1**: Create database migration for staff documents tables
- [x] **BE-5.8.2**: Create document type module (entity, dto, repository, service, handler)
- [x] **BE-5.8.3**: Create staff document module with file upload handling
- [x] **BE-5.8.4**: Implement document verification workflow
- [x] **BE-5.8.5**: Create expiry tracking and notification service
- [x] **BE-5.8.6**: Create compliance report generation
- [x] **BE-5.8.7**: Register routes and add permissions

### Frontend Tasks

- [x] **FE-5.8.1**: Create document models and services
- [x] **FE-5.8.2**: Create document upload component with drag-and-drop
- [x] **FE-5.8.3**: Create document list component with filters
- [x] **FE-5.8.4**: Create document verification component
- [x] **FE-5.8.5**: Create document type configuration page
- [x] **FE-5.8.6**: Create expiring documents dashboard widget
- [x] **FE-5.8.7**: Create compliance report page
- [x] **FE-5.8.8**: Add routes and navigation

---

## Definition of Done

- [x] Documents can be uploaded with metadata
- [x] Documents support multiple file types (PDF, images)
- [x] Expiry dates are tracked with notifications
- [x] HR can verify or reject documents
- [x] Compliance dashboard shows document status
- [x] Missing mandatory documents are highlighted
- [x] Documents can be downloaded securely
- [ ] Unit tests pass

---

## Dependencies

- Story 5.1: Staff Profile Management (completed)
- MinIO/S3 storage configuration (infrastructure)

---

## Dev Notes

### Backend Structure

```
internal/modules/staffdocument/
├── entity.go
├── dto.go
├── repository.go
├── service.go
├── handler.go
├── errors.go
└── routes.go

internal/services/notification/
└── document_expiry.go  # Expiry notification service
```

### Frontend Structure

```
src/app/features/staff/documents/
├── documents.routes.ts
├── models/
│   ├── document.model.ts
│   └── document-type.model.ts
├── services/
│   └── document.service.ts
├── pages/
│   ├── document-list/
│   ├── document-types/
│   └── compliance-report/
└── components/
    ├── document-upload/
    ├── document-card/
    ├── document-verification/
    └── expiring-documents-widget/
```

### File Storage

- Use MinIO for document storage
- Path format: `/{tenant_id}/staff/{staff_id}/documents/{document_id}/{filename}`
- Generate presigned URLs for secure downloads
- Max file size: 5MB
- Allowed types: PDF, JPG, JPEG, PNG

### Notification Schedule

- 30 days before expiry: First reminder
- 7 days before expiry: Urgent reminder
- On expiry day: Expired notification
- Use background job/cron for checking expirations daily
