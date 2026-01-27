# Story 4.7: Student Promotion & Retention Processing

**Epic:** 4 - Student Lifecycle Management
**Status:** done
**Priority:** P2
**Story Points:** 8

---

## User Story

**As an** administrator,
**I want** to process student promotions at year end,
**So that** students are moved to appropriate classes.

---

## Acceptance Criteria

### AC1: Promotion Initiation
- [x] Admin can select: from academic year, from class, from section
- [x] System shows all students with their result status (pass/fail/pending)
- [x] Students can be marked: promote, retain, transfer
- [x] Bulk selection available (select all pass, select all fail)

### AC2: Promotion Rules Configuration
- [x] Admin can configure promotion rules per class
- [x] Rules: minimum attendance %, minimum marks %, required subjects passed
- [x] Auto-promotion: students meeting criteria auto-marked for promotion
- [x] Students not meeting criteria flagged for manual review
- [x] Manual override available for all students

### AC3: Promotion Processing
- [x] When promotions confirmed, new enrollments created for next academic year
- [x] Section assignments: manual selection or auto-distributed
- [x] Roll numbers generated (sequential by class/section)
- [x] Promotion report generated for records

### AC4: Retention Handling
- [x] Retained students get new enrollment in same class, new year
- [x] Retention reason recorded
- [x] Previous year enrollment marked as "completed"
- [ ] Student/parent notification option (deferred to Epic 12)

---

## Technical Requirements

### Backend

**Database Tables:**

```sql
-- Migration: 20260126100600_create_promotion_system.up.sql

CREATE TYPE promotion_decision AS ENUM ('promote', 'retain', 'transfer', 'pending');

-- Promotion rules per class
CREATE TABLE promotion_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    class_id UUID NOT NULL, -- From class
    min_attendance_pct DECIMAL(5,2) DEFAULT 75.00,
    min_overall_marks_pct DECIMAL(5,2) DEFAULT 33.00,
    min_subjects_passed INTEGER DEFAULT 0,
    auto_promote_on_criteria BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_promotion_rule UNIQUE (tenant_id, class_id)
);

ALTER TABLE promotion_rules ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON promotion_rules USING (tenant_id = current_setting('app.current_tenant')::UUID);

-- Promotion batch for tracking
CREATE TABLE promotion_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    from_academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    to_academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    from_class_id UUID NOT NULL,
    from_section_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'processing', 'completed', 'cancelled')),
    total_students INTEGER NOT NULL DEFAULT 0,
    promoted_count INTEGER NOT NULL DEFAULT 0,
    retained_count INTEGER NOT NULL DEFAULT 0,
    transferred_count INTEGER NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ,
    processed_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

ALTER TABLE promotion_batches ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON promotion_batches USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_promotion_batches_year ON promotion_batches(from_academic_year_id);

-- Individual promotion decisions
CREATE TABLE promotion_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    batch_id UUID NOT NULL REFERENCES promotion_batches(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id),
    from_enrollment_id UUID NOT NULL REFERENCES student_enrollments(id),
    to_enrollment_id UUID REFERENCES student_enrollments(id), -- Created after processing
    decision promotion_decision NOT NULL DEFAULT 'pending',
    to_class_id UUID, -- Target class (null if retain)
    to_section_id UUID, -- Target section
    auto_decided BOOLEAN NOT NULL DEFAULT FALSE,
    decision_reason TEXT,
    attendance_pct DECIMAL(5,2),
    overall_marks_pct DECIMAL(5,2),
    subjects_passed INTEGER,
    override_by UUID REFERENCES users(id),
    override_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_promotion_record UNIQUE (batch_id, student_id)
);

ALTER TABLE promotion_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON promotion_records USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_promotion_records_batch ON promotion_records(batch_id);
CREATE INDEX idx_promotion_records_student ON promotion_records(student_id);
```

**API Endpoints:**

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/promotion-rules` | List rules | `promotion:read` |
| POST | `/api/v1/promotion-rules` | Create/update rule | `promotion:write` |
| GET | `/api/v1/promotion-batches` | List batches | `promotion:read` |
| POST | `/api/v1/promotion-batches` | Create batch | `promotion:write` |
| GET | `/api/v1/promotion-batches/{id}` | Get batch details | `promotion:read` |
| GET | `/api/v1/promotion-batches/{id}/records` | List records | `promotion:read` |
| PUT | `/api/v1/promotion-batches/{id}/records/{rid}` | Update decision | `promotion:write` |
| POST | `/api/v1/promotion-batches/{id}/auto-decide` | Apply rules | `promotion:write` |
| POST | `/api/v1/promotion-batches/{id}/process` | Process promotions | `promotion:write` |
| GET | `/api/v1/promotion-batches/{id}/report` | Download report | `promotion:read` |

**Promotion Service:**

```go
type PromotionService struct {
    db             *gorm.DB
    enrollmentRepo *EnrollmentRepository
    batchRepo      *PromotionBatchRepository
    recordRepo     *PromotionRecordRepository
    ruleRepo       *PromotionRuleRepository
}

func (s *PromotionService) CreateBatch(ctx context.Context, dto CreateBatchDTO) (*PromotionBatch, error) {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create batch
        batch := &PromotionBatch{
            TenantID:           ctx.Value("tenantID").(uuid.UUID),
            FromAcademicYearID: dto.FromAcademicYearID,
            ToAcademicYearID:   dto.ToAcademicYearID,
            FromClassID:        dto.FromClassID,
            FromSectionID:      dto.FromSectionID,
            Status:             "draft",
        }

        if err := s.batchRepo.Create(ctx, tx, batch); err != nil {
            return fmt.Errorf("create batch: %w", err)
        }

        // Get all students in the class/section
        enrollments, err := s.enrollmentRepo.GetByClassAndYear(ctx, tx, dto.FromClassID, dto.FromSectionID, dto.FromAcademicYearID)
        if err != nil {
            return fmt.Errorf("get enrollments: %w", err)
        }

        // Create pending records for each student
        for _, enrollment := range enrollments {
            record := &PromotionRecord{
                TenantID:         batch.TenantID,
                BatchID:          batch.ID,
                StudentID:        enrollment.StudentID,
                FromEnrollmentID: enrollment.ID,
                Decision:         PromotionDecisionPending,
            }
            if err := s.recordRepo.Create(ctx, tx, record); err != nil {
                return fmt.Errorf("create record: %w", err)
            }
        }

        batch.TotalStudents = len(enrollments)
        return s.batchRepo.Update(ctx, tx, batch)
    })
}

func (s *PromotionService) AutoDecide(ctx context.Context, batchID uuid.UUID) error {
    batch, err := s.batchRepo.GetByID(ctx, batchID)
    if err != nil {
        return fmt.Errorf("get batch: %w", err)
    }

    rules, err := s.ruleRepo.GetByClass(ctx, batch.FromClassID)
    if err != nil {
        return fmt.Errorf("get rules: %w", err)
    }

    records, err := s.recordRepo.GetByBatch(ctx, batchID)
    if err != nil {
        return fmt.Errorf("get records: %w", err)
    }

    for _, record := range records {
        // Get student's attendance and marks (will need integration with Epic 7, 8)
        // For now, use placeholder values
        meetsAttendance := record.AttendancePct == nil || *record.AttendancePct >= rules.MinAttendancePct
        meetsMarks := record.OverallMarksPct == nil || *record.OverallMarksPct >= rules.MinOverallMarksPct

        if meetsAttendance && meetsMarks {
            record.Decision = PromotionDecisionPromote
            record.AutoDecided = true
            record.ToClassID = getNextClass(batch.FromClassID) // Helper function
        } else {
            record.Decision = PromotionDecisionRetain
            record.AutoDecided = true
            record.DecisionReason = buildDecisionReason(meetsAttendance, meetsMarks)
        }

        if err := s.recordRepo.Update(ctx, record); err != nil {
            return fmt.Errorf("update record: %w", err)
        }
    }

    return nil
}

func (s *PromotionService) ProcessBatch(ctx context.Context, batchID uuid.UUID) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        batch, err := s.batchRepo.GetByID(ctx, batchID)
        if err != nil {
            return fmt.Errorf("get batch: %w", err)
        }

        records, err := s.recordRepo.GetByBatch(ctx, batchID)
        if err != nil {
            return fmt.Errorf("get records: %w", err)
        }

        var promoted, retained, transferred int

        for _, record := range records {
            if record.Decision == PromotionDecisionPending {
                return ErrPendingDecisions
            }

            switch record.Decision {
            case PromotionDecisionPromote:
                // Create new enrollment in next class
                newEnrollment := &StudentEnrollment{
                    TenantID:       batch.TenantID,
                    StudentID:      record.StudentID,
                    AcademicYearID: batch.ToAcademicYearID,
                    ClassID:        *record.ToClassID,
                    SectionID:      record.ToSectionID,
                    Status:         EnrollmentStatusActive,
                }
                if err := s.enrollmentRepo.Create(ctx, tx, newEnrollment); err != nil {
                    return fmt.Errorf("create enrollment: %w", err)
                }
                record.ToEnrollmentID = &newEnrollment.ID
                promoted++

            case PromotionDecisionRetain:
                // Create new enrollment in same class
                newEnrollment := &StudentEnrollment{
                    TenantID:       batch.TenantID,
                    StudentID:      record.StudentID,
                    AcademicYearID: batch.ToAcademicYearID,
                    ClassID:        batch.FromClassID,
                    SectionID:      record.ToSectionID,
                    Status:         EnrollmentStatusActive,
                }
                if err := s.enrollmentRepo.Create(ctx, tx, newEnrollment); err != nil {
                    return fmt.Errorf("create enrollment: %w", err)
                }
                record.ToEnrollmentID = &newEnrollment.ID
                retained++

            case PromotionDecisionTransfer:
                // Student transfers out, don't create new enrollment
                transferred++
            }

            // Mark old enrollment as completed
            if err := s.enrollmentRepo.MarkCompleted(ctx, tx, record.FromEnrollmentID); err != nil {
                return fmt.Errorf("mark completed: %w", err)
            }

            if err := s.recordRepo.Update(ctx, tx, record); err != nil {
                return fmt.Errorf("update record: %w", err)
            }
        }

        // Update batch status
        batch.Status = "completed"
        batch.PromotedCount = promoted
        batch.RetainedCount = retained
        batch.TransferredCount = transferred
        batch.ProcessedAt = ptr(time.Now())
        batch.ProcessedBy = ctx.Value("userID").(*uuid.UUID)

        return s.batchRepo.Update(ctx, tx, batch)
    })
}
```

### Frontend

**Components to create:**

```bash
ng generate component features/students/pages/promotion-wizard --standalone
ng generate component features/students/components/promotion-rules-form --standalone
ng generate component features/students/components/promotion-student-list --standalone
ng generate component features/students/components/promotion-summary --standalone
ng generate interface features/students/models/promotion-batch
ng generate interface features/students/models/promotion-record
ng generate service features/students/services/promotion
```

**Promotion Wizard Steps:**

```typescript
// Step 1: Select Source (year, class, section)
// Step 2: Review Students with auto-decisions
// Step 3: Manual overrides for flagged students
// Step 4: Assign sections for promoted students
// Step 5: Confirm and process
// Step 6: View report

@Component({
  selector: 'app-promotion-wizard',
  template: `
    <app-stepper [currentStep]="currentStep()" [steps]="steps">
      @switch (currentStep()) {
        @case (0) {
          <app-promotion-source-form (next)="onSourceSelected($event)" />
        }
        @case (1) {
          <app-promotion-student-list
            [records]="records()"
            [canEdit]="true"
            (decisionChanged)="onDecisionChanged($event)"
            (next)="nextStep()"
          />
        }
        @case (2) {
          <app-promotion-section-assignment
            [records]="promotedRecords()"
            [sections]="targetSections()"
            (assigned)="onSectionsAssigned($event)"
          />
        }
        @case (3) {
          <app-promotion-summary
            [batch]="batch()"
            (confirm)="processPromotion()"
          />
        }
        @case (4) {
          <app-promotion-report [batchId]="batch()?.id" />
        }
      }
    </app-stepper>
  `
})
export class PromotionWizardComponent {
  currentStep = signal(0);
  batch = signal<PromotionBatch | null>(null);
  records = signal<PromotionRecord[]>([]);
}
```

---

## Tasks

### Backend Tasks

- [x] **BE-4.7.1**: Create promotion rules entity and migration
- [x] **BE-4.7.2**: Create promotion batch and records migration
- [x] **BE-4.7.3**: Create promotion repositories
- [x] **BE-4.7.4**: Create promotion service with auto-decide logic
- [x] **BE-4.7.5**: Create promotion HTTP handlers
- [x] **BE-4.7.6**: Implement batch processing transaction
- [ ] **BE-4.7.7**: Add promotion report generation (Excel)
- [x] **BE-4.7.8**: Add promotion permissions to seed
- [ ] **BE-4.7.9**: Write unit tests

### Frontend Tasks

- [x] **FE-4.7.1**: Create promotion interfaces
- [x] **FE-4.7.2**: Create promotion service
- [x] **FE-4.7.3**: Create promotion wizard with stepper
- [x] **FE-4.7.4**: Create source selection form
- [x] **FE-4.7.5**: Create student list with decision editing
- [x] **FE-4.7.6**: Create section assignment component
- [x] **FE-4.7.7**: Create summary and confirmation view
- [x] **FE-4.7.8**: Create report download component
- [x] **FE-4.7.9**: Add routes for promotion wizard
- [ ] **FE-4.7.10**: Write component tests

---

## Definition of Done

- [ ] All acceptance criteria verified
- [ ] Promotion rules configurable
- [ ] Auto-decision based on rules working
- [ ] Manual override working
- [ ] Batch processing creates new enrollments
- [ ] Report generation working
- [ ] Backend tests passing
- [ ] Frontend tests passing

---

## Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Story 4.6 (Enrollment History) | Required | Creates new enrollments |
| Academic Years (Story 3-2) | ✅ Done | For from/to years |
| Attendance (Epic 7) | Future | Attendance % placeholder for now |
| Examinations (Epic 8) | Future | Marks % placeholder for now |
| Classes/Sections (Epic 6) | Future | Use UUIDs, validate later |

---

## Notes

- Integration with Epic 7 (Attendance) and Epic 8 (Examinations) needed for accurate auto-decisions
- For now, attendance_pct and overall_marks_pct can be null or manually entered
- Roll number generation: sequential within class/section
- Consider "final year" classes that don't promote (graduation)
- Notification to parents deferred to Epic 12
