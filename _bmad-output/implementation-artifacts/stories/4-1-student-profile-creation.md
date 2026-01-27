# Story 4.1: Student Profile Creation

**Epic:** 4 - Student Lifecycle Management
**Status:** done
**Priority:** P0 (Foundation for all Epic 4 stories)
**Story Points:** 8

---

## User Story

**As an** administrator,
**I want** to create and manage student profiles with complete information,
**So that** all student data is centrally managed.

---

## Acceptance Criteria

### AC1: Student Creation Form
- [ ] Admin can enter personal details: first name, last name, DOB, gender, blood group, Aadhaar number
- [ ] Admin can enter academic details: class, section (roll number assigned after enrollment)
- [ ] Admin can enter address details: current address, permanent address (with "same as current" option)
- [ ] Admin can upload: photo (max 2MB, jpg/png), birth certificate (PDF/image)
- [ ] Admission number is auto-generated: `{BRANCH_CODE}-{YEAR}-{SEQUENCE}` (e.g., `MUM-2026-00001`)
- [ ] Student status defaults to "active"
- [ ] Form validates all required fields before submission

### AC2: Student Profile View
- [ ] Profile displays all information in organized sections (Personal, Academic, Address, Documents)
- [ ] Quick actions available: Edit, View Attendance, View Fees, View Documents
- [ ] Edit history available via audit log link
- [ ] Photo displayed with fallback avatar (initials)
- [ ] Status badge shows current status (active, inactive, transferred, graduated)

### AC3: Student List View
- [ ] Paginated list (20 per page default, cursor-based)
- [ ] Columns: Photo, Name, Admission No, Class/Section, Status, Actions
- [ ] Quick filters: Class, Section, Status
- [ ] Search by: Name, Admission Number
- [ ] Click row to view profile

### AC4: Student Edit
- [ ] All fields editable except admission number
- [ ] Changes logged to audit trail (old value, new value, user, timestamp)
- [ ] Optimistic locking to prevent concurrent edit conflicts

---

## Technical Requirements

### Backend

**New Module:** `internal/modules/student/`

```
internal/modules/student/
├── entity.go          # Student, StudentAddress entities
├── dto.go             # CreateStudentDTO, UpdateStudentDTO, StudentResponse
├── repository.go      # GORM repository
├── service.go         # Business logic
├── handler.go         # HTTP handlers
├── errors.go          # ErrStudentNotFound, ErrDuplicateAdmissionNumber
└── student_test.go    # Unit tests
```

**Database Tables:**

```sql
-- Migration: 20260126100000_create_students_table.up.sql

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID NOT NULL REFERENCES branches(id),
    admission_number VARCHAR(20) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    blood_group VARCHAR(5),
    aadhaar_number VARCHAR(12),
    photo_url VARCHAR(500),
    birth_certificate_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'transferred', 'graduated')),
    admission_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    CONSTRAINT uniq_students_tenant_admission UNIQUE (tenant_id, admission_number)
);

ALTER TABLE students ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON students USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_students_tenant ON students(tenant_id);
CREATE INDEX idx_students_branch ON students(branch_id);
CREATE INDEX idx_students_status ON students(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_name ON students(tenant_id, last_name, first_name);

CREATE TABLE student_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    address_type VARCHAR(20) NOT NULL CHECK (address_type IN ('current', 'permanent')),
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'India',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_student_address_type UNIQUE (student_id, address_type)
);

ALTER TABLE student_addresses ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON student_addresses USING (tenant_id = current_setting('app.current_tenant')::UUID);
```

**API Endpoints:**

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/students` | List students (paginated) | `student:read` |
| GET | `/api/v1/students/{id}` | Get student by ID | `student:read` |
| POST | `/api/v1/students` | Create student | `student:create` |
| PUT | `/api/v1/students/{id}` | Update student | `student:update` |
| DELETE | `/api/v1/students/{id}` | Soft delete student | `student:delete` |
| POST | `/api/v1/students/{id}/photo` | Upload photo | `student:update` |
| POST | `/api/v1/students/{id}/documents` | Upload document | `student:update` |

**Admission Number Generation:**

```go
func (s *Service) generateAdmissionNumber(ctx context.Context, branchID uuid.UUID) (string, error) {
    branch, err := s.branchRepo.GetByID(ctx, branchID)
    if err != nil {
        return "", fmt.Errorf("get branch: %w", err)
    }

    year := time.Now().Year()
    sequence, err := s.repo.GetNextSequence(ctx, branchID, year)
    if err != nil {
        return "", fmt.Errorf("get sequence: %w", err)
    }

    return fmt.Sprintf("%s-%d-%05d", branch.Code, year, sequence), nil
}
```

### Frontend

**Feature:** `src/app/features/students/`

**Components to create (using ng generate):**

```bash
# Pages
ng generate component features/students/pages/student-list --standalone
ng generate component features/students/pages/student-detail --standalone
ng generate component features/students/pages/student-form --standalone

# Components
ng generate component features/students/components/student-card --standalone
ng generate component features/students/components/address-form --standalone
ng generate component features/students/components/photo-upload --standalone

# Service
ng generate service features/students/services/student

# Models
ng generate interface features/students/models/student
ng generate interface features/students/models/student-address
```

**Routes:**

```typescript
// students.routes.ts
export const STUDENT_ROUTES: Routes = [
  { path: '', component: StudentListComponent },
  { path: 'new', component: StudentFormComponent },
  { path: ':id', component: StudentDetailComponent },
  { path: ':id/edit', component: StudentFormComponent },
];
```

**State Management (Signals):**

```typescript
@Injectable({ providedIn: 'root' })
export class StudentService {
  private _students = signal<Student[]>([]);
  private _loading = signal<boolean>(false);
  private _selectedStudent = signal<Student | null>(null);

  readonly students = this._students.asReadonly();
  readonly loading = this._loading.asReadonly();
  readonly selectedStudent = this._selectedStudent.asReadonly();
}
```

---

## Tasks

### Backend Tasks

- [x] **BE-4.1.1**: Create student entity and migration
- [x] **BE-4.1.2**: Create student repository with GORM
- [x] **BE-4.1.3**: Create student service with admission number generation
- [x] **BE-4.1.4**: Create student HTTP handlers
- [x] **BE-4.1.5**: Add student routes to main router
- [x] **BE-4.1.6**: Create file upload endpoint for photos
- [x] **BE-4.1.7**: Add audit logging hooks
- [x] **BE-4.1.8**: Write unit tests (target: 80% coverage)

### Frontend Tasks

- [x] **FE-4.1.1**: Create student model interfaces
- [x] **FE-4.1.2**: Create student service with HTTP calls
- [x] **FE-4.1.3**: Create student list page with data table
- [x] **FE-4.1.4**: Create student form page (create/edit)
- [x] **FE-4.1.5**: Create student detail page
- [x] **FE-4.1.6**: Create photo upload component
- [x] **FE-4.1.7**: Create address form component
- [x] **FE-4.1.8**: Add routes and navigation
- [x] **FE-4.1.9**: Write component tests

---

## Definition of Done

- [x] All acceptance criteria verified
- [x] Backend unit tests passing (80%+ coverage)
- [x] Frontend component tests passing
- [ ] API documented in OpenAPI spec
- [x] No lint/type errors
- [ ] Code reviewed
- [ ] Audit logging verified
- [ ] RLS policies verified with different tenants

---

## Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Branch module (Epic 3) | ✅ Done | For branch code in admission number |
| User/Auth module (Epic 2) | ✅ Done | For created_by, updated_by |
| File storage infrastructure | ⚠️ | May need to implement storage service |
| Admission module (Epic 3) | ✅ Done | Students can be linked from admitted applications |

---

## Notes

- Students created here are typically converted from admitted applications (Story 3-7)
- Consider adding a "Convert to Student" action in admission decision workflow
- Photo storage path: `/{tenant_id}/students/{student_id}/photo.{ext}`
- Document storage path: `/{tenant_id}/students/{student_id}/documents/{doc_type}.{ext}`
