// Package promotion provides student promotion and retention processing functionality.
package promotion

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/modules/enrollment"
)

// Service handles business logic for promotions.
type Service struct {
	repo           *Repository
	enrollmentRepo *enrollment.Repository
}

// NewService creates a new promotion service.
func NewService(repo *Repository, enrollmentRepo *enrollment.Repository) *Service {
	return &Service{
		repo:           repo,
		enrollmentRepo: enrollmentRepo,
	}
}

// ========================================================================
// Promotion Rules
// ========================================================================

// CreateOrUpdateRule creates or updates a promotion rule for a class.
func (s *Service) CreateOrUpdateRule(ctx context.Context, tenantID uuid.UUID, req CreateRuleRequest, userID *uuid.UUID) (*PromotionRule, error) {
	rule := &PromotionRule{
		TenantID:              tenantID,
		ClassID:               req.ClassID,
		MinAttendancePct:      req.MinAttendancePct,
		MinOverallMarksPct:    req.MinOverallMarksPct,
		AutoPromoteOnCriteria: true,
		IsActive:              true,
		CreatedBy:             userID,
		UpdatedBy:             userID,
	}

	if req.MinSubjectsPassed != nil {
		rule.MinSubjectsPassed = *req.MinSubjectsPassed
	}
	if req.AutoPromoteOnCriteria != nil {
		rule.AutoPromoteOnCriteria = *req.AutoPromoteOnCriteria
	}

	if err := s.repo.UpsertRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("upsert rule: %w", err)
	}

	// Fetch the updated rule
	return s.repo.GetRuleByClass(ctx, tenantID, req.ClassID)
}

// GetRule retrieves a promotion rule by ID.
func (s *Service) GetRule(ctx context.Context, tenantID, id uuid.UUID) (*PromotionRule, error) {
	return s.repo.GetRuleByID(ctx, tenantID, id)
}

// GetRuleByClass retrieves a promotion rule by class ID.
func (s *Service) GetRuleByClass(ctx context.Context, tenantID, classID uuid.UUID) (*PromotionRule, error) {
	return s.repo.GetRuleByClass(ctx, tenantID, classID)
}

// ListRules retrieves all promotion rules for a tenant.
func (s *Service) ListRules(ctx context.Context, tenantID uuid.UUID) ([]PromotionRule, error) {
	return s.repo.ListRules(ctx, tenantID)
}

// DeleteRule deletes a promotion rule.
func (s *Service) DeleteRule(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteRule(ctx, tenantID, id)
}

// ========================================================================
// Promotion Batches
// ========================================================================

// CreateBatch creates a new promotion batch with records for all students in the class/section.
func (s *Service) CreateBatch(ctx context.Context, tenantID uuid.UUID, req CreateBatchRequest, userID *uuid.UUID) (*PromotionBatch, error) {
	// Validate academic years
	if req.FromAcademicYearID == req.ToAcademicYearID {
		return nil, ErrSameAcademicYear
	}

	// Get enrollments for the source class/section
	enrollments, err := s.getEnrollmentsForBatch(ctx, tenantID, req.FromAcademicYearID, req.FromClassID, req.FromSectionID)
	if err != nil {
		return nil, err
	}

	if len(enrollments) == 0 {
		return nil, ErrNoStudentsInClass
	}

	var batch *PromotionBatch

	// Use transaction for batch creation
	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		// Create the batch
		batch = &PromotionBatch{
			TenantID:           tenantID,
			FromAcademicYearID: req.FromAcademicYearID,
			ToAcademicYearID:   req.ToAcademicYearID,
			FromClassID:        req.FromClassID,
			FromSectionID:      req.FromSectionID,
			ToClassID:          req.ToClassID,
			Status:             BatchStatusDraft,
			TotalStudents:      len(enrollments),
			Notes:              req.Notes,
			CreatedBy:          userID,
		}

		if err := s.repo.CreateBatchWithTx(ctx, tx, batch); err != nil {
			return fmt.Errorf("create batch: %w", err)
		}

		// Create promotion records for each student
		records := make([]PromotionRecord, 0, len(enrollments))
		for _, e := range enrollments {
			record := PromotionRecord{
				TenantID:         tenantID,
				BatchID:          batch.ID,
				StudentID:        e.StudentID,
				FromEnrollmentID: e.ID,
				Decision:         DecisionPending,
				ToClassID:        req.ToClassID, // Pre-populate target class
			}
			records = append(records, record)
		}

		if err := s.repo.CreateRecordsBatchWithTx(ctx, tx, records); err != nil {
			return fmt.Errorf("create records: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch the created batch with preloads
	return s.repo.GetBatchByID(ctx, tenantID, batch.ID)
}

// getEnrollmentsForBatch retrieves enrollments for batch creation.
func (s *Service) getEnrollmentsForBatch(ctx context.Context, tenantID, academicYearID, classID uuid.UUID, sectionID *uuid.UUID) ([]enrollment.StudentEnrollment, error) {
	db := s.repo.DB().WithContext(ctx)

	query := db.Model(&enrollment.StudentEnrollment{}).
		Where("tenant_id = ? AND academic_year_id = ? AND class_id = ? AND status = ?",
			tenantID, academicYearID, classID, enrollment.EnrollmentStatusActive)

	if sectionID != nil {
		query = query.Where("section_id = ?", *sectionID)
	}

	var enrollments []enrollment.StudentEnrollment
	if err := query.Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("get enrollments: %w", err)
	}

	return enrollments, nil
}

// GetBatch retrieves a promotion batch by ID.
func (s *Service) GetBatch(ctx context.Context, tenantID, id uuid.UUID) (*PromotionBatch, error) {
	return s.repo.GetBatchByID(ctx, tenantID, id)
}

// ListBatches retrieves all promotion batches for a tenant.
func (s *Service) ListBatches(ctx context.Context, tenantID uuid.UUID, status *BatchStatus) ([]PromotionBatch, error) {
	return s.repo.ListBatches(ctx, tenantID, status)
}

// CancelBatch cancels a promotion batch.
func (s *Service) CancelBatch(ctx context.Context, tenantID, id uuid.UUID, reason string, userID *uuid.UUID) error {
	batch, err := s.repo.GetBatchByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if batch.Status != BatchStatusDraft {
		return ErrBatchNotDraft
	}

	now := time.Now()
	batch.Status = BatchStatusCancelled
	batch.CancelledAt = &now
	batch.CancelledBy = userID
	batch.CancellationReason = reason

	return s.repo.UpdateBatch(ctx, batch)
}

// DeleteBatch deletes a draft promotion batch.
func (s *Service) DeleteBatch(ctx context.Context, tenantID, id uuid.UUID) error {
	batch, err := s.repo.GetBatchByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if batch.Status != BatchStatusDraft {
		return ErrBatchNotDraft
	}

	return s.repo.DeleteBatch(ctx, tenantID, id)
}

// ========================================================================
// Promotion Records
// ========================================================================

// GetRecordsByBatch retrieves all promotion records for a batch.
func (s *Service) GetRecordsByBatch(ctx context.Context, tenantID, batchID uuid.UUID) ([]PromotionRecord, *RecordsSummary, error) {
	records, err := s.repo.ListRecordsByBatch(ctx, tenantID, batchID)
	if err != nil {
		return nil, nil, err
	}

	summary, err := s.repo.GetRecordsSummary(ctx, tenantID, batchID)
	if err != nil {
		return nil, nil, err
	}

	return records, summary, nil
}

// UpdateRecord updates a single promotion record.
func (s *Service) UpdateRecord(ctx context.Context, tenantID, batchID, recordID uuid.UUID, req UpdateRecordRequest, userID *uuid.UUID) (*PromotionRecord, error) {
	// Verify batch is in draft status
	batch, err := s.repo.GetBatchByID(ctx, tenantID, batchID)
	if err != nil {
		return nil, err
	}

	if batch.Status != BatchStatusDraft {
		return nil, ErrBatchNotDraft
	}

	// Get the record
	record, err := s.repo.GetRecordByID(ctx, tenantID, recordID)
	if err != nil {
		return nil, err
	}

	if record.BatchID != batchID {
		return nil, ErrRecordNotFound
	}

	// Update fields
	if req.Decision != nil {
		record.Decision = *req.Decision
		record.AutoDecided = false
		record.OverrideBy = userID
		now := time.Now()
		record.OverrideAt = &now
	}
	if req.ToClassID != nil {
		record.ToClassID = req.ToClassID
	}
	if req.ToSectionID != nil {
		record.ToSectionID = req.ToSectionID
	}
	if req.RollNumber != nil {
		record.RollNumber = *req.RollNumber
	}
	if req.OverrideReason != nil {
		record.OverrideReason = *req.OverrideReason
	}
	if req.RetentionReason != nil {
		record.RetentionReason = *req.RetentionReason
	}
	if req.TransferDestination != nil {
		record.TransferDestination = *req.TransferDestination
	}

	if err := s.repo.UpdateRecord(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

// BulkUpdateRecords updates multiple promotion records with the same decision.
func (s *Service) BulkUpdateRecords(ctx context.Context, tenantID, batchID uuid.UUID, req BulkUpdateRecordsRequest, userID *uuid.UUID) error {
	// Verify batch is in draft status
	batch, err := s.repo.GetBatchByID(ctx, tenantID, batchID)
	if err != nil {
		return err
	}

	if batch.Status != BatchStatusDraft {
		return ErrBatchNotDraft
	}

	now := time.Now()

	for _, recordID := range req.RecordIDs {
		record, err := s.repo.GetRecordByID(ctx, tenantID, recordID)
		if err != nil {
			continue // Skip not found records
		}

		if record.BatchID != batchID {
			continue
		}

		record.Decision = req.Decision
		record.AutoDecided = false
		record.OverrideBy = userID
		record.OverrideAt = &now
		record.OverrideReason = req.Reason

		if req.ToClassID != nil {
			record.ToClassID = req.ToClassID
		}
		if req.ToSectionID != nil {
			record.ToSectionID = req.ToSectionID
		}

		if err := s.repo.UpdateRecord(ctx, record); err != nil {
			return fmt.Errorf("update record %s: %w", recordID, err)
		}
	}

	return nil
}

// ========================================================================
// Auto-Decision
// ========================================================================

// AutoDecide applies promotion rules to auto-decide pending records.
func (s *Service) AutoDecide(ctx context.Context, tenantID, batchID uuid.UUID) error {
	batch, err := s.repo.GetBatchByID(ctx, tenantID, batchID)
	if err != nil {
		return err
	}

	if batch.Status != BatchStatusDraft {
		return ErrBatchNotDraft
	}

	// Get promotion rules for the class
	rules, err := s.repo.GetRuleByClass(ctx, tenantID, batch.FromClassID)
	if err != nil && err != ErrRuleNotFound {
		return err
	}

	// Get pending records
	records, err := s.repo.ListRecordsByDecision(ctx, tenantID, batchID, DecisionPending)
	if err != nil {
		return err
	}

	for _, record := range records {
		decision, reason := s.evaluateRecord(&record, rules, batch.ToClassID)
		record.Decision = decision
		record.DecisionReason = reason
		record.AutoDecided = true

		if decision == DecisionPromote && batch.ToClassID != nil {
			record.ToClassID = batch.ToClassID
		} else if decision == DecisionRetain {
			record.ToClassID = &batch.FromClassID
		}

		if err := s.repo.UpdateRecord(ctx, &record); err != nil {
			return fmt.Errorf("update record %s: %w", record.ID, err)
		}
	}

	return nil
}

// evaluateRecord evaluates a student against promotion rules.
func (s *Service) evaluateRecord(record *PromotionRecord, rules *PromotionRule, targetClass *uuid.UUID) (PromotionDecision, string) {
	// If no rules defined, default to promote
	if rules == nil || !rules.AutoPromoteOnCriteria {
		if targetClass != nil {
			return DecisionPromote, "No rules defined - default promotion"
		}
		return DecisionPending, "No rules defined and no target class specified"
	}

	reasons := []string{}
	meetsAllCriteria := true

	// Check attendance (if available and rules specify)
	if rules.MinAttendancePct != nil && record.AttendancePct != nil {
		if *record.AttendancePct < *rules.MinAttendancePct {
			meetsAllCriteria = false
			reasons = append(reasons, fmt.Sprintf("Attendance %.1f%% < required %.1f%%", *record.AttendancePct, *rules.MinAttendancePct))
		}
	}

	// Check marks (if available and rules specify)
	if rules.MinOverallMarksPct != nil && record.OverallMarksPct != nil {
		if *record.OverallMarksPct < *rules.MinOverallMarksPct {
			meetsAllCriteria = false
			reasons = append(reasons, fmt.Sprintf("Marks %.1f%% < required %.1f%%", *record.OverallMarksPct, *rules.MinOverallMarksPct))
		}
	}

	// Check subjects passed (if available and rules specify)
	if rules.MinSubjectsPassed > 0 && record.SubjectsPassed != nil {
		if *record.SubjectsPassed < rules.MinSubjectsPassed {
			meetsAllCriteria = false
			reasons = append(reasons, fmt.Sprintf("Subjects passed %d < required %d", *record.SubjectsPassed, rules.MinSubjectsPassed))
		}
	}

	if meetsAllCriteria {
		return DecisionPromote, "Meets all promotion criteria"
	}

	reason := "Does not meet criteria: "
	for i, r := range reasons {
		if i > 0 {
			reason += "; "
		}
		reason += r
	}

	return DecisionRetain, reason
}

// ========================================================================
// Batch Processing
// ========================================================================

// ProcessBatch processes a promotion batch, creating new enrollments.
func (s *Service) ProcessBatch(ctx context.Context, tenantID, batchID uuid.UUID, generateRollNumbers bool, userID *uuid.UUID) error {
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		// Get batch
		batch, err := s.repo.GetBatchByIDWithTx(ctx, tx, tenantID, batchID)
		if err != nil {
			return err
		}

		if batch.Status == BatchStatusCompleted {
			return ErrBatchAlreadyProcessed
		}
		if batch.Status == BatchStatusCancelled {
			return ErrBatchCancelled
		}
		if batch.Status != BatchStatusDraft {
			return ErrBatchNotProcessable
		}

		// Update batch status to processing
		batch.Status = BatchStatusProcessing
		if err := s.repo.UpdateBatchWithTx(ctx, tx, batch); err != nil {
			return err
		}

		// Get all records
		records, err := s.repo.ListRecordsByBatchWithTx(ctx, tx, tenantID, batchID)
		if err != nil {
			return err
		}

		// Check for pending decisions
		for _, r := range records {
			if r.Decision == DecisionPending {
				batch.Status = BatchStatusDraft
				s.repo.UpdateBatchWithTx(ctx, tx, batch)
				return ErrPendingDecisions
			}
		}

		var promotedCount, retainedCount, transferredCount int
		rollNumberCounter := make(map[string]int) // section -> counter

		for i := range records {
			record := &records[i]

			switch record.Decision {
			case DecisionPromote:
				if record.ToClassID == nil {
					batch.Status = BatchStatusDraft
					s.repo.UpdateBatchWithTx(ctx, tx, batch)
					return ErrMissingTargetClass
				}

				// Generate roll number if requested
				rollNumber := record.RollNumber
				if generateRollNumbers && rollNumber == "" {
					sectionKey := ""
					if record.ToSectionID != nil {
						sectionKey = record.ToSectionID.String()
					}
					rollNumberCounter[sectionKey]++
					rollNumber = fmt.Sprintf("%d", rollNumberCounter[sectionKey])
				}

				// Create new enrollment
				newEnrollment := &enrollment.StudentEnrollment{
					TenantID:       tenantID,
					StudentID:      record.StudentID,
					AcademicYearID: batch.ToAcademicYearID,
					ClassID:        record.ToClassID,
					SectionID:      record.ToSectionID,
					RollNumber:     rollNumber,
					Status:         enrollment.EnrollmentStatusActive,
					EnrollmentDate: time.Now(),
					CreatedBy:      userID,
				}

				if err := s.enrollmentRepo.CreateWithTx(ctx, tx, newEnrollment); err != nil {
					return fmt.Errorf("create enrollment for student %s: %w", record.StudentID, err)
				}

				record.ToEnrollmentID = &newEnrollment.ID
				record.RollNumber = rollNumber
				promotedCount++

			case DecisionRetain:
				// Create enrollment in same class
				rollNumber := record.RollNumber
				if generateRollNumbers && rollNumber == "" {
					sectionKey := ""
					if record.ToSectionID != nil {
						sectionKey = record.ToSectionID.String()
					}
					rollNumberCounter[sectionKey]++
					rollNumber = fmt.Sprintf("%d", rollNumberCounter[sectionKey])
				}

				newEnrollment := &enrollment.StudentEnrollment{
					TenantID:       tenantID,
					StudentID:      record.StudentID,
					AcademicYearID: batch.ToAcademicYearID,
					ClassID:        &batch.FromClassID, // Same class
					SectionID:      record.ToSectionID,
					RollNumber:     rollNumber,
					Status:         enrollment.EnrollmentStatusActive,
					EnrollmentDate: time.Now(),
					CreatedBy:      userID,
				}

				if err := s.enrollmentRepo.CreateWithTx(ctx, tx, newEnrollment); err != nil {
					return fmt.Errorf("create retention enrollment for student %s: %w", record.StudentID, err)
				}

				record.ToEnrollmentID = &newEnrollment.ID
				record.RollNumber = rollNumber
				retainedCount++

			case DecisionTransfer:
				// No new enrollment created - student is leaving
				transferredCount++
			}

			// Mark old enrollment as completed
			oldEnrollment, err := s.enrollmentRepo.GetActiveByStudentWithTx(ctx, tx, tenantID, record.StudentID)
			if err == nil && oldEnrollment != nil {
				now := time.Now()
				oldEnrollment.Status = enrollment.EnrollmentStatusCompleted
				oldEnrollment.CompletionDate = &now
				oldEnrollment.UpdatedBy = userID
				if err := s.enrollmentRepo.UpdateWithTx(ctx, tx, oldEnrollment); err != nil {
					return fmt.Errorf("complete old enrollment for student %s: %w", record.StudentID, err)
				}
			}

			// Update record
			if err := s.repo.UpdateRecordWithTx(ctx, tx, record); err != nil {
				return fmt.Errorf("update record %s: %w", record.ID, err)
			}
		}

		// Update batch counts and status
		now := time.Now()
		batch.Status = BatchStatusCompleted
		batch.PromotedCount = promotedCount
		batch.RetainedCount = retainedCount
		batch.TransferredCount = transferredCount
		batch.ProcessedAt = &now
		batch.ProcessedBy = userID

		return s.repo.UpdateBatchWithTx(ctx, tx, batch)
	})
}

// ========================================================================
// Reporting
// ========================================================================

// GetPromotionReport generates a report for a completed batch.
func (s *Service) GetPromotionReport(ctx context.Context, tenantID, batchID uuid.UUID) ([]PromotionReportRow, error) {
	batch, err := s.repo.GetBatchByID(ctx, tenantID, batchID)
	if err != nil {
		return nil, err
	}

	records, err := s.repo.ListRecordsByBatch(ctx, tenantID, batchID)
	if err != nil {
		return nil, err
	}

	rows := make([]PromotionReportRow, 0, len(records))
	for _, r := range records {
		row := PromotionReportRow{
			Decision:   string(r.Decision),
			RollNumber: r.RollNumber,
		}

		if r.Student != nil {
			row.StudentAdmissionNo = r.Student.AdmissionNumber
			row.StudentName = r.Student.FullName()
		}

		// Format class IDs (would be names if we had class lookup)
		row.FromClass = batch.FromClassID.String()
		if batch.FromSectionID != nil {
			row.FromSection = batch.FromSectionID.String()
		}

		if r.ToClassID != nil {
			row.ToClass = r.ToClassID.String()
		}
		if r.ToSectionID != nil {
			row.ToSection = r.ToSectionID.String()
		}

		if r.AttendancePct != nil {
			row.AttendancePct = fmt.Sprintf("%.1f%%", *r.AttendancePct)
		}
		if r.OverallMarksPct != nil {
			row.MarksPct = fmt.Sprintf("%.1f%%", *r.OverallMarksPct)
		}

		switch r.Decision {
		case DecisionRetain:
			row.Reason = r.RetentionReason
			if row.Reason == "" {
				row.Reason = r.DecisionReason
			}
		case DecisionTransfer:
			row.Reason = r.TransferDestination
		default:
			row.Reason = r.DecisionReason
		}

		rows = append(rows, row)
	}

	return rows, nil
}
