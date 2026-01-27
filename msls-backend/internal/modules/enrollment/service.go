// Package enrollment provides student enrollment management functionality.
package enrollment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// StudentRepository defines the interface for student operations.
type StudentRepository interface {
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Student, error)
	UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status models.StudentStatus, updatedBy *uuid.UUID) error
}

// AcademicYearRepository defines the interface for academic year operations.
type AcademicYearRepository interface {
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.AcademicYear, error)
}

// Service handles enrollment business logic.
type Service struct {
	repo            *Repository
	studentRepo     StudentRepository
	academicYearRepo AcademicYearRepository
	db              *gorm.DB
}

// NewService creates a new enrollment service.
func NewService(db *gorm.DB, studentRepo StudentRepository, academicYearRepo AcademicYearRepository) *Service {
	return &Service{
		repo:            NewRepository(db),
		studentRepo:     studentRepo,
		academicYearRepo: academicYearRepo,
		db:              db,
	}
}

// Create creates a new enrollment for a student.
func (s *Service) Create(ctx context.Context, dto CreateEnrollmentDTO) (*StudentEnrollment, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.StudentID == uuid.Nil {
		return nil, ErrStudentIDRequired
	}
	if dto.AcademicYearID == uuid.Nil {
		return nil, ErrAcademicYearIDRequired
	}

	// Verify student exists
	if s.studentRepo != nil {
		_, err := s.studentRepo.GetByID(ctx, dto.TenantID, dto.StudentID)
		if err != nil {
			return nil, ErrStudentNotFound
		}
	}

	// Verify academic year exists
	if s.academicYearRepo != nil {
		_, err := s.academicYearRepo.GetByID(ctx, dto.TenantID, dto.AcademicYearID)
		if err != nil {
			return nil, ErrAcademicYearNotFound
		}
	}

	// Set enrollment date
	enrollmentDate := time.Now()
	if dto.EnrollmentDate != nil {
		enrollmentDate = *dto.EnrollmentDate
	}

	enrollment := &StudentEnrollment{
		TenantID:       dto.TenantID,
		StudentID:      dto.StudentID,
		AcademicYearID: dto.AcademicYearID,
		ClassID:        dto.ClassID,
		SectionID:      dto.SectionID,
		RollNumber:     dto.RollNumber,
		ClassTeacherID: dto.ClassTeacherID,
		Status:         EnrollmentStatusActive,
		EnrollmentDate: enrollmentDate,
		Notes:          dto.Notes,
		CreatedBy:      dto.CreatedBy,
		UpdatedBy:      dto.CreatedBy,
	}

	var createdEnrollment *StudentEnrollment

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the enrollment
		if err := s.repo.CreateWithTx(ctx, tx, enrollment); err != nil {
			return err
		}

		// Log the initial status
		if dto.CreatedBy != nil {
			statusChange := &EnrollmentStatusChange{
				TenantID:     dto.TenantID,
				EnrollmentID: enrollment.ID,
				FromStatus:   nil, // Initial creation
				ToStatus:     EnrollmentStatusActive,
				ChangeReason: "Initial enrollment",
				ChangeDate:   enrollmentDate,
				ChangedBy:    *dto.CreatedBy,
			}
			if err := s.repo.CreateStatusChangeWithTx(ctx, tx, statusChange); err != nil {
				return fmt.Errorf("log initial status: %w", err)
			}
		}

		createdEnrollment = enrollment
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch the complete enrollment with relations
	return s.GetByID(ctx, dto.TenantID, createdEnrollment.ID)
}

// GetByID retrieves an enrollment by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*StudentEnrollment, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetCurrentByStudent retrieves the active enrollment for a student.
func (s *Service) GetCurrentByStudent(ctx context.Context, tenantID, studentID uuid.UUID) (*StudentEnrollment, error) {
	return s.repo.GetActiveByStudent(ctx, tenantID, studentID)
}

// ListByStudent retrieves all enrollments for a student (enrollment history).
func (s *Service) ListByStudent(ctx context.Context, tenantID, studentID uuid.UUID) ([]StudentEnrollment, error) {
	return s.repo.ListByStudent(ctx, tenantID, studentID)
}

// ListByClass retrieves all active enrollments for a class.
func (s *Service) ListByClass(ctx context.Context, tenantID, classID uuid.UUID) ([]StudentEnrollment, error) {
	return s.repo.ListByClass(ctx, tenantID, classID)
}

// ListBySection retrieves all active enrollments for a section.
func (s *Service) ListBySection(ctx context.Context, tenantID, sectionID uuid.UUID) ([]StudentEnrollment, error) {
	return s.repo.ListBySection(ctx, tenantID, sectionID)
}

// Update updates an enrollment.
func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, dto UpdateEnrollmentDTO) (*StudentEnrollment, error) {
	enrollment, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Can only update active enrollments
	if enrollment.Status != EnrollmentStatusActive {
		return nil, ErrCannotModifyInactiveEnrollment
	}

	// Apply updates
	if dto.ClassID != nil {
		enrollment.ClassID = dto.ClassID
	}
	if dto.SectionID != nil {
		enrollment.SectionID = dto.SectionID
	}
	if dto.RollNumber != nil {
		enrollment.RollNumber = *dto.RollNumber
	}
	if dto.ClassTeacherID != nil {
		enrollment.ClassTeacherID = dto.ClassTeacherID
	}
	if dto.Notes != nil {
		enrollment.Notes = *dto.Notes
	}
	enrollment.UpdatedBy = dto.UpdatedBy

	if err := s.repo.Update(ctx, enrollment); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, id)
}

// ProcessTransfer processes a student transfer.
func (s *Service) ProcessTransfer(ctx context.Context, tenantID, studentID uuid.UUID, dto TransferDTO) (*StudentEnrollment, error) {
	// Validate required fields
	if dto.TransferDate.IsZero() {
		return nil, ErrTransferDateRequired
	}
	if dto.TransferReason == "" {
		return nil, ErrTransferReasonRequired
	}
	if dto.UpdatedBy == nil {
		return nil, ErrChangedByRequired
	}

	var enrollment *StudentEnrollment

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the active enrollment
		var err error
		enrollment, err = s.repo.GetActiveByStudentWithTx(ctx, tx, tenantID, studentID)
		if err != nil {
			return err
		}

		// Record previous status for logging
		previousStatus := enrollment.Status

		// Update enrollment status
		enrollment.Status = EnrollmentStatusTransferred
		enrollment.TransferDate = &dto.TransferDate
		enrollment.TransferReason = dto.TransferReason
		enrollment.UpdatedBy = dto.UpdatedBy

		if err := s.repo.UpdateWithTx(ctx, tx, enrollment); err != nil {
			return fmt.Errorf("update enrollment: %w", err)
		}

		// Log status change
		statusChange := &EnrollmentStatusChange{
			TenantID:     tenantID,
			EnrollmentID: enrollment.ID,
			FromStatus:   &previousStatus,
			ToStatus:     EnrollmentStatusTransferred,
			ChangeReason: dto.TransferReason,
			ChangeDate:   dto.TransferDate,
			ChangedBy:    *dto.UpdatedBy,
		}
		if err := s.repo.CreateStatusChangeWithTx(ctx, tx, statusChange); err != nil {
			return fmt.Errorf("log status change: %w", err)
		}

		// Update student status to inactive
		if s.studentRepo != nil {
			if err := s.studentRepo.UpdateStatus(ctx, tenantID, studentID, models.StudentStatusInactive, dto.UpdatedBy); err != nil {
				return fmt.Errorf("update student status: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, enrollment.ID)
}

// ProcessDropout processes a student dropout.
func (s *Service) ProcessDropout(ctx context.Context, tenantID, studentID uuid.UUID, dto DropoutDTO) (*StudentEnrollment, error) {
	// Validate required fields
	if dto.DropoutDate.IsZero() {
		return nil, ErrDropoutDateRequired
	}
	if dto.DropoutReason == "" {
		return nil, ErrDropoutReasonRequired
	}
	if dto.UpdatedBy == nil {
		return nil, ErrChangedByRequired
	}

	var enrollment *StudentEnrollment

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the active enrollment
		var err error
		enrollment, err = s.repo.GetActiveByStudentWithTx(ctx, tx, tenantID, studentID)
		if err != nil {
			return err
		}

		// Record previous status for logging
		previousStatus := enrollment.Status

		// Update enrollment status
		enrollment.Status = EnrollmentStatusDropout
		enrollment.DropoutDate = &dto.DropoutDate
		enrollment.DropoutReason = dto.DropoutReason
		enrollment.UpdatedBy = dto.UpdatedBy

		if err := s.repo.UpdateWithTx(ctx, tx, enrollment); err != nil {
			return fmt.Errorf("update enrollment: %w", err)
		}

		// Log status change
		statusChange := &EnrollmentStatusChange{
			TenantID:     tenantID,
			EnrollmentID: enrollment.ID,
			FromStatus:   &previousStatus,
			ToStatus:     EnrollmentStatusDropout,
			ChangeReason: dto.DropoutReason,
			ChangeDate:   dto.DropoutDate,
			ChangedBy:    *dto.UpdatedBy,
		}
		if err := s.repo.CreateStatusChangeWithTx(ctx, tx, statusChange); err != nil {
			return fmt.Errorf("log status change: %w", err)
		}

		// Update student status to inactive
		if s.studentRepo != nil {
			if err := s.studentRepo.UpdateStatus(ctx, tenantID, studentID, models.StudentStatusInactive, dto.UpdatedBy); err != nil {
				return fmt.Errorf("update student status: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, enrollment.ID)
}

// CompleteEnrollment marks an enrollment as completed (end of academic year).
func (s *Service) CompleteEnrollment(ctx context.Context, tenantID, enrollmentID uuid.UUID, dto CompleteEnrollmentDTO) (*StudentEnrollment, error) {
	if dto.UpdatedBy == nil {
		return nil, ErrChangedByRequired
	}

	enrollment, err := s.repo.GetByID(ctx, tenantID, enrollmentID)
	if err != nil {
		return nil, err
	}

	// Can only complete active enrollments
	if enrollment.Status != EnrollmentStatusActive {
		return nil, ErrInvalidStatusTransition
	}

	// Record previous status for logging
	previousStatus := enrollment.Status

	// Set completion date
	completionDate := time.Now()
	if dto.CompletionDate != nil {
		completionDate = *dto.CompletionDate
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update enrollment status
		enrollment.Status = EnrollmentStatusCompleted
		enrollment.CompletionDate = &completionDate
		enrollment.UpdatedBy = dto.UpdatedBy

		if err := s.repo.UpdateWithTx(ctx, tx, enrollment); err != nil {
			return fmt.Errorf("update enrollment: %w", err)
		}

		// Log status change
		statusChange := &EnrollmentStatusChange{
			TenantID:     tenantID,
			EnrollmentID: enrollment.ID,
			FromStatus:   &previousStatus,
			ToStatus:     EnrollmentStatusCompleted,
			ChangeReason: "Academic year completed",
			ChangeDate:   completionDate,
			ChangedBy:    *dto.UpdatedBy,
		}
		if err := s.repo.CreateStatusChangeWithTx(ctx, tx, statusChange); err != nil {
			return fmt.Errorf("log status change: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, enrollmentID)
}

// GetStatusHistory retrieves the status change history for an enrollment.
func (s *Service) GetStatusHistory(ctx context.Context, tenantID, enrollmentID uuid.UUID) ([]EnrollmentStatusChange, error) {
	// Verify enrollment exists
	_, err := s.repo.GetByID(ctx, tenantID, enrollmentID)
	if err != nil {
		return nil, err
	}

	return s.repo.ListStatusChanges(ctx, tenantID, enrollmentID)
}

// Delete deletes an enrollment.
func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}
