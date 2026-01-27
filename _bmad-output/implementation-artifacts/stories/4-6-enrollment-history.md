# Story 4.6: Student Enrollment History

**Epic:** 4 - Student Lifecycle Management
**Status:** done
**Priority:** P1
**Story Points:** 5

---

## User Story

**As an** administrator,
**I want** to track student enrollment across academic years,
**So that** complete academic history is maintained.

---

## Acceptance Criteria

### AC1: Enrollment History View
- [x] Shows all enrollment records: academic year, class, section, roll number, status
- [x] History displayed chronologically from admission to current
- [x] Each record shows: class teacher (if assigned), attendance summary, exam results link

### AC2: New Academic Year Enrollment
- [x] When promotions processed, new enrollment record created automatically
- [x] Previous year enrollment marked as "completed"
- [x] Class teacher assignment recorded (links to Staff module - nullable until Epic 5)

### AC3: Transfer Processing
- [x] When student transfers out, enrollment status changes to "transferred"
- [x] Transfer date and reason recorded
- [x] Student status changes to "inactive"
- [ ] Transfer certificate generation link available (Epic 15)

### AC4: Enrollment Status Management
- [x] Status options: active, completed, transferred, dropout
- [x] Status change logged with reason and date
- [x] Only active enrollment per student per academic year

---

## Technical Requirements

### Backend

**Database Tables:**

```sql
-- Migration: 20260126100500_create_student_enrollments.up.sql

CREATE TYPE enrollment_status AS ENUM ('active', 'completed', 'transferred', 'dropout');

CREATE TABLE student_enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    class_id UUID NOT NULL, -- References classes table (Epic 6)
    section_id UUID, -- References sections table (Epic 6), nullable
    roll_number VARCHAR(20),
    class_teacher_id UUID, -- References staff (Epic 5), nullable
    status enrollment_status NOT NULL DEFAULT 'active',
    enrollment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    completion_date DATE,
    transfer_date DATE,
    transfer_reason TEXT,
    dropout_date DATE,
    dropout_reason TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    CONSTRAINT uniq_student_year UNIQUE (student_id, academic_year_id)
);

ALTER TABLE student_enrollments ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON student_enrollments USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_enrollments_student ON student_enrollments(student_id);
CREATE INDEX idx_enrollments_year ON student_enrollments(academic_year_id);
CREATE INDEX idx_enrollments_class ON student_enrollments(class_id);
CREATE INDEX idx_enrollments_status ON student_enrollments(status) WHERE status = 'active';

-- Ensure only one active enrollment per student
CREATE UNIQUE INDEX uniq_active_enrollment ON student_enrollments(student_id)
    WHERE status = 'active';

-- Enrollment status change log
CREATE TABLE enrollment_status_changes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    enrollment_id UUID NOT NULL REFERENCES student_enrollments(id) ON DELETE CASCADE,
    from_status enrollment_status,
    to_status enrollment_status NOT NULL,
    change_reason TEXT,
    change_date DATE NOT NULL DEFAULT CURRENT_DATE,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    changed_by UUID NOT NULL REFERENCES users(id)
);

ALTER TABLE enrollment_status_changes ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON enrollment_status_changes USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_status_changes_enrollment ON enrollment_status_changes(enrollment_id);
```

**API Endpoints:**

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/students/{id}/enrollments` | List enrollment history | `student:read` |
| GET | `/api/v1/students/{id}/enrollments/current` | Get current enrollment | `student:read` |
| POST | `/api/v1/students/{id}/enrollments` | Create enrollment | `student:update` |
| PUT | `/api/v1/students/{id}/enrollments/{eid}` | Update enrollment | `student:update` |
| POST | `/api/v1/students/{id}/enrollments/{eid}/transfer` | Process transfer | `student:update` |
| POST | `/api/v1/students/{id}/enrollments/{eid}/dropout` | Process dropout | `student:update` |
| GET | `/api/v1/enrollments/by-class/{classId}` | List enrollments by class | `student:read` |
| GET | `/api/v1/enrollments/by-section/{sectionId}` | List enrollments by section | `student:read` |

**Transfer Service:**

```go
type TransferDTO struct {
    TransferDate   time.Time `json:"transferDate" binding:"required"`
    TransferReason string    `json:"transferReason" binding:"required"`
}

func (s *Service) ProcessTransfer(ctx context.Context, studentID uuid.UUID, dto TransferDTO) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Get current enrollment
        enrollment, err := s.enrollmentRepo.GetActiveByStudent(ctx, tx, studentID)
        if err != nil {
            return fmt.Errorf("get active enrollment: %w", err)
        }

        // Update enrollment status
        enrollment.Status = EnrollmentStatusTransferred
        enrollment.TransferDate = &dto.TransferDate
        enrollment.TransferReason = dto.TransferReason
        if err := s.enrollmentRepo.Update(ctx, tx, enrollment); err != nil {
            return fmt.Errorf("update enrollment: %w", err)
        }

        // Log status change
        if err := s.logStatusChange(ctx, tx, enrollment.ID, EnrollmentStatusActive, EnrollmentStatusTransferred, dto.TransferReason); err != nil {
            return fmt.Errorf("log status change: %w", err)
        }

        // Update student status to inactive
        if err := s.studentRepo.UpdateStatus(ctx, tx, studentID, StudentStatusInactive); err != nil {
            return fmt.Errorf("update student status: %w", err)
        }

        return nil
    })
}
```

**Enrollment History Response:**

```go
type EnrollmentHistoryItem struct {
    ID              uuid.UUID        `json:"id"`
    AcademicYear    AcademicYearRef  `json:"academicYear"`
    Class           ClassRef         `json:"class"`
    Section         *SectionRef      `json:"section,omitempty"`
    RollNumber      string           `json:"rollNumber,omitempty"`
    ClassTeacher    *StaffRef        `json:"classTeacher,omitempty"`
    Status          string           `json:"status"`
    EnrollmentDate  time.Time        `json:"enrollmentDate"`
    CompletionDate  *time.Time       `json:"completionDate,omitempty"`
    TransferDate    *time.Time       `json:"transferDate,omitempty"`
    AttendancePct   *float64         `json:"attendancePercentage,omitempty"` // Calculated, optional
}
```

### Frontend

**Components to create:**

```bash
ng generate component features/students/components/enrollment-history --standalone
ng generate component features/students/components/enrollment-form --standalone
ng generate component features/students/components/transfer-form --standalone
ng generate interface features/students/models/enrollment
ng generate interface features/students/models/enrollment-history
```

**Enrollment History Timeline:**

```typescript
@Component({
  selector: 'app-enrollment-history',
  template: `
    <div class="relative">
      <!-- Timeline line -->
      <div class="absolute left-4 top-0 bottom-0 w-0.5 bg-gray-200"></div>

      @for (enrollment of enrollments(); track enrollment.id; let i = $index) {
        <div class="relative pl-10 pb-8">
          <!-- Timeline dot -->
          <div class="absolute left-2.5 w-3 h-3 rounded-full"
               [class.bg-green-500]="enrollment.status === 'active'"
               [class.bg-gray-400]="enrollment.status === 'completed'"
               [class.bg-orange-500]="enrollment.status === 'transferred'"
               [class.bg-red-500]="enrollment.status === 'dropout'">
          </div>

          <div class="bg-white border rounded-lg p-4">
            <div class="flex justify-between items-start">
              <div>
                <h4 class="font-medium">{{ enrollment.academicYear.name }}</h4>
                <p class="text-sm text-gray-600">
                  {{ enrollment.class.name }}
                  @if (enrollment.section) {
                    - {{ enrollment.section.name }}
                  }
                </p>
                @if (enrollment.rollNumber) {
                  <p class="text-sm text-gray-500">Roll No: {{ enrollment.rollNumber }}</p>
                }
              </div>
              <app-badge [variant]="getStatusVariant(enrollment.status)">
                {{ enrollment.status | titlecase }}
              </app-badge>
            </div>

            <div class="mt-2 text-sm text-gray-500">
              <p>Enrolled: {{ enrollment.enrollmentDate | date:'mediumDate' }}</p>
              @if (enrollment.classTeacher) {
                <p>Class Teacher: {{ enrollment.classTeacher.name }}</p>
              }
            </div>
          </div>
        </div>
      }
    </div>
  `
})
export class EnrollmentHistoryComponent {
  studentId = input.required<string>();
  enrollments = signal<EnrollmentHistoryItem[]>([]);
}
```

---

## Tasks

### Backend Tasks

- [x] **BE-4.6.1**: Create enrollment entity and migration
- [x] **BE-4.6.2**: Create enrollment status change log migration
- [x] **BE-4.6.3**: Create enrollment repository
- [x] **BE-4.6.4**: Create enrollment service with transfer/dropout logic
- [x] **BE-4.6.5**: Create enrollment HTTP handlers
- [x] **BE-4.6.6**: Add status change logging
- [x] **BE-4.6.7**: Write unit tests

### Frontend Tasks

- [x] **FE-4.6.1**: Create enrollment interfaces
- [x] **FE-4.6.2**: Create enrollment service
- [x] **FE-4.6.3**: Create enrollment history timeline component
- [x] **FE-4.6.4**: Create enrollment form (for manual entry)
- [x] **FE-4.6.5**: Create transfer form modal
- [x] **FE-4.6.6**: Add Enrollment History tab to student profile
- [x] **FE-4.6.7**: Write model tests (component tests require TestBed setup)

---

## Definition of Done

- [x] All acceptance criteria verified
- [x] Enrollment history displays correctly
- [x] Transfer process updates student status
- [x] Only one active enrollment per student enforced
- [x] Status changes logged
- [x] Backend tests passing (10/10)
- [x] Frontend model tests passing (11/11)

---

## Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Story 4.1 (Student Profile) | Required | Enrollment links to student |
| Academic Years (Story 3-2) | ✅ Done | For academic_year_id |
| Classes/Sections (Epic 6) | Future | Use UUID placeholders, nullable |
| Staff (Epic 5) | Future | Class teacher nullable |

---

## Notes

- class_id and section_id reference tables from Epic 6 (Academic Structure)
- Until Epic 6 is done, store as UUIDs without FK constraint
- class_teacher_id references staff from Epic 5 - nullable until then
- Enrollment history is key for promotion processing (Story 4.7)
- Transfer certificate generation will be in Epic 15
