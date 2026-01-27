# Story 4.5: Student Document Management

**Epic:** 4 - Student Lifecycle Management
**Status:** done
**Priority:** P1
**Story Points:** 5

---

## User Story

**As an** administrator,
**I want** to manage and verify student documents,
**So that** required documents are maintained and verified.

---

## Acceptance Criteria

### AC1: Document Upload
- [x] Can upload documents: birth certificate, Aadhaar, transfer certificate, caste certificate, income certificate, photo ID, medical certificate
- [x] Supported formats: PDF, JPG, PNG (max 5MB per file)
- [x] Can enter: document number, issue date, expiry date (if applicable)
- [x] Document status defaults to "pending_verification"
- [x] Document metadata stored (uploaded by, uploaded at)

### AC2: Document Verification
- [x] Admin can verify documents: mark as verified or rejected
- [x] If rejected, reason must be provided
- [x] Verification date and verifier automatically recorded
- [x] Document status updates: pending_verification → verified/rejected

### AC3: Document Checklist
- [x] Required documents configurable per class/branch
- [x] Student profile shows document checklist with status indicators
- [x] Missing documents highlighted in red
- [ ] Bulk reminder can be sent for missing documents (future: Epic 12)

### AC4: Documents Tab in Student Profile
- [x] Tab shows all uploaded documents by category
- [x] Each document shows: type, status badge, uploaded date, actions
- [x] Actions: View, Download, Re-upload, Verify/Reject (admin only)
- [x] Document preview in modal (PDF/image viewer)

---

## Technical Requirements

### Backend

**Database Tables:**

```sql
-- Migration: 20260126100400_create_student_documents.up.sql

CREATE TYPE document_status AS ENUM ('pending_verification', 'verified', 'rejected');

CREATE TABLE document_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_mandatory BOOLEAN NOT NULL DEFAULT FALSE,
    has_expiry BOOLEAN NOT NULL DEFAULT FALSE,
    allowed_extensions VARCHAR(100) NOT NULL DEFAULT 'pdf,jpg,png',
    max_size_mb INTEGER NOT NULL DEFAULT 5,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_doc_type_code UNIQUE (tenant_id, code)
);

ALTER TABLE document_types ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON document_types USING (tenant_id = current_setting('app.current_tenant')::UUID);

-- Seed default document types
INSERT INTO document_types (tenant_id, code, name, is_mandatory, sort_order) VALUES
    ('system', 'birth_certificate', 'Birth Certificate', TRUE, 1),
    ('system', 'aadhaar', 'Aadhaar Card', TRUE, 2),
    ('system', 'transfer_certificate', 'Transfer Certificate', FALSE, 3),
    ('system', 'caste_certificate', 'Caste Certificate', FALSE, 4),
    ('system', 'income_certificate', 'Income Certificate', FALSE, 5),
    ('system', 'photo_id', 'Photo ID', FALSE, 6),
    ('system', 'medical_certificate', 'Medical Certificate', FALSE, 7);

CREATE TABLE student_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    document_type_id UUID NOT NULL REFERENCES document_types(id),
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size_bytes INTEGER NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    document_number VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,
    status document_status NOT NULL DEFAULT 'pending_verification',
    rejection_reason TEXT,
    verified_at TIMESTAMPTZ,
    verified_by UUID REFERENCES users(id),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    uploaded_by UUID NOT NULL REFERENCES users(id),
    CONSTRAINT uniq_student_doc_type UNIQUE (student_id, document_type_id)
);

ALTER TABLE student_documents ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON student_documents USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_student_documents_student ON student_documents(student_id);
CREATE INDEX idx_student_documents_status ON student_documents(status);

-- Required documents per class (optional configuration)
CREATE TABLE class_required_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    class_id UUID, -- NULL means applies to all classes
    document_type_id UUID NOT NULL REFERENCES document_types(id),
    is_required BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_class_doc UNIQUE (tenant_id, class_id, document_type_id)
);

ALTER TABLE class_required_documents ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON class_required_documents USING (tenant_id = current_setting('app.current_tenant')::UUID);
```

**API Endpoints:**

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/document-types` | List document types | `student:read` |
| POST | `/api/v1/document-types` | Create document type | `admin:settings` |
| GET | `/api/v1/students/{id}/documents` | List student documents | `student:read` |
| POST | `/api/v1/students/{id}/documents` | Upload document | `student:update` |
| GET | `/api/v1/students/{id}/documents/{did}` | Get document | `student:read` |
| PUT | `/api/v1/students/{id}/documents/{did}` | Update document metadata | `student:update` |
| DELETE | `/api/v1/students/{id}/documents/{did}` | Delete document | `student:update` |
| POST | `/api/v1/students/{id}/documents/{did}/verify` | Verify document | `document:verify` |
| POST | `/api/v1/students/{id}/documents/{did}/reject` | Reject document | `document:verify` |
| GET | `/api/v1/students/{id}/document-checklist` | Get checklist with status | `student:read` |

**File Storage:**

```go
// Storage path pattern
func (s *Service) GetDocumentPath(tenantID, studentID uuid.UUID, docType string, fileName string) string {
    return fmt.Sprintf("%s/students/%s/documents/%s/%s",
        tenantID.String(),
        studentID.String(),
        docType,
        fileName,
    )
}

// Upload handler
func (h *Handler) UploadDocument(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.Error(ErrNoFileUploaded)
        return
    }
    defer file.Close()

    // Validate file size
    if header.Size > int64(docType.MaxSizeMB)*1024*1024 {
        c.Error(ErrFileTooLarge)
        return
    }

    // Validate extension
    ext := filepath.Ext(header.Filename)
    if !isAllowedExtension(ext, docType.AllowedExtensions) {
        c.Error(ErrInvalidFileType)
        return
    }

    // Upload to storage
    url, err := h.storage.Upload(ctx, path, file)
    if err != nil {
        c.Error(fmt.Errorf("upload file: %w", err))
        return
    }

    // Create document record
    // ...
}
```

### Frontend

**Components to create:**

```bash
ng generate component features/students/components/document-list --standalone
ng generate component features/students/components/document-upload --standalone
ng generate component features/students/components/document-preview --standalone
ng generate component features/students/components/document-checklist --standalone
ng generate interface features/students/models/document
ng generate interface features/students/models/document-type
ng generate service features/students/services/document
```

**Document Checklist Component:**

```typescript
@Component({
  selector: 'app-document-checklist',
  template: `
    <div class="space-y-2">
      <h4 class="font-medium text-gray-900">Document Checklist</h4>
      @for (item of checklist(); track item.documentType.id) {
        <div class="flex items-center justify-between p-2 border rounded">
          <div class="flex items-center gap-2">
            @if (item.document) {
              @switch (item.document.status) {
                @case ('verified') {
                  <span class="text-green-500">✓</span>
                }
                @case ('rejected') {
                  <span class="text-red-500">✗</span>
                }
                @case ('pending_verification') {
                  <span class="text-yellow-500">⏳</span>
                }
              }
            } @else {
              <span class="text-gray-400">○</span>
            }
            <span [class.text-red-600]="item.documentType.isMandatory && !item.document">
              {{ item.documentType.name }}
              @if (item.documentType.isMandatory) {
                <span class="text-red-500">*</span>
              }
            </span>
          </div>
          <app-button
            size="sm"
            [variant]="item.document ? 'ghost' : 'primary'"
            (click)="uploadDocument(item.documentType)"
          >
            {{ item.document ? 'Re-upload' : 'Upload' }}
          </app-button>
        </div>
      }
    </div>
  `
})
export class DocumentChecklistComponent {
  studentId = input.required<string>();
  checklist = signal<ChecklistItem[]>([]);
}
```

---

## Tasks

### Backend Tasks

- [x] **BE-4.5.1**: Create document types entity and migration
- [x] **BE-4.5.2**: Create student documents entity and migration
- [x] **BE-4.5.3**: Implement file storage service (local/MinIO)
- [x] **BE-4.5.4**: Create document repository
- [x] **BE-4.5.5**: Create document service with upload/verify logic
- [x] **BE-4.5.6**: Create document HTTP handlers
- [x] **BE-4.5.7**: Add document permissions to seed
- [ ] **BE-4.5.8**: Write unit tests

### Frontend Tasks

- [x] **FE-4.5.1**: Create document interfaces
- [x] **FE-4.5.2**: Create document service with upload
- [x] **FE-4.5.3**: Create document list component
- [x] **FE-4.5.4**: Create document upload component (drag-drop)
- [x] **FE-4.5.5**: Create document preview modal (PDF viewer)
- [x] **FE-4.5.6**: Create document checklist component
- [x] **FE-4.5.7**: Add Documents tab to student profile
- [ ] **FE-4.5.8**: Write component tests

---

## Definition of Done

- [x] All acceptance criteria verified
- [x] File upload working (PDF, images)
- [x] Verification workflow working
- [x] Document checklist accurate
- [x] File size/type validation working
- [ ] Backend tests passing
- [ ] Frontend tests passing

---

## Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Story 4.1 (Student Profile) | Required | Documents link to student |
| File storage service | Required | Need storage abstraction |
| Notification (Epic 12) | Future | Bulk reminders deferred |

---

## Notes

- Consider virus scanning for uploads (future enhancement)
- Expiry date reminders for documents like medical certificates
- Document types are tenant-configurable
- Storage path ensures tenant isolation
