// Package enrollment provides student enrollment management functionality.
package enrollment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles database operations for enrollments.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new enrollment repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new enrollment in the database.
func (r *Repository) Create(ctx context.Context, enrollment *StudentEnrollment) error {
	if err := r.db.WithContext(ctx).Create(enrollment).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "uniq_student_year") {
				return ErrDuplicateEnrollment
			}
			if strings.Contains(err.Error(), "uniq_active_enrollment") {
				return ErrActiveEnrollmentExists
			}
		}
		return fmt.Errorf("create enrollment: %w", err)
	}
	return nil
}

// CreateWithTx creates a new enrollment within a transaction.
func (r *Repository) CreateWithTx(ctx context.Context, tx *gorm.DB, enrollment *StudentEnrollment) error {
	if err := tx.WithContext(ctx).Create(enrollment).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "uniq_student_year") {
				return ErrDuplicateEnrollment
			}
			if strings.Contains(err.Error(), "uniq_active_enrollment") {
				return ErrActiveEnrollmentExists
			}
		}
		return fmt.Errorf("create enrollment: %w", err)
	}
	return nil
}

// GetByID retrieves an enrollment by ID.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*StudentEnrollment, error) {
	var enrollment StudentEnrollment
	err := r.db.WithContext(ctx).
		Preload("AcademicYear").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEnrollmentNotFound
		}
		return nil, fmt.Errorf("get enrollment by id: %w", err)
	}
	return &enrollment, nil
}

// GetActiveByStudent retrieves the active enrollment for a student.
func (r *Repository) GetActiveByStudent(ctx context.Context, tenantID, studentID uuid.UUID) (*StudentEnrollment, error) {
	var enrollment StudentEnrollment
	err := r.db.WithContext(ctx).
		Preload("AcademicYear").
		Where("tenant_id = ? AND student_id = ? AND status = ?", tenantID, studentID, EnrollmentStatusActive).
		First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActiveEnrollmentNotFound
		}
		return nil, fmt.Errorf("get active enrollment: %w", err)
	}
	return &enrollment, nil
}

// GetActiveByStudentWithTx retrieves the active enrollment for a student within a transaction.
func (r *Repository) GetActiveByStudentWithTx(ctx context.Context, tx *gorm.DB, tenantID, studentID uuid.UUID) (*StudentEnrollment, error) {
	var enrollment StudentEnrollment
	err := tx.WithContext(ctx).
		Preload("AcademicYear").
		Where("tenant_id = ? AND student_id = ? AND status = ?", tenantID, studentID, EnrollmentStatusActive).
		First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActiveEnrollmentNotFound
		}
		return nil, fmt.Errorf("get active enrollment: %w", err)
	}
	return &enrollment, nil
}

// GetByStudentAndYear retrieves an enrollment for a specific student and academic year.
func (r *Repository) GetByStudentAndYear(ctx context.Context, tenantID, studentID, academicYearID uuid.UUID) (*StudentEnrollment, error) {
	var enrollment StudentEnrollment
	err := r.db.WithContext(ctx).
		Preload("AcademicYear").
		Where("tenant_id = ? AND student_id = ? AND academic_year_id = ?", tenantID, studentID, academicYearID).
		First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEnrollmentNotFound
		}
		return nil, fmt.Errorf("get enrollment by student and year: %w", err)
	}
	return &enrollment, nil
}

// ListByStudent retrieves all enrollments for a student ordered by enrollment date descending.
func (r *Repository) ListByStudent(ctx context.Context, tenantID, studentID uuid.UUID) ([]StudentEnrollment, error) {
	var enrollments []StudentEnrollment
	err := r.db.WithContext(ctx).
		Preload("AcademicYear").
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("enrollment_date DESC").
		Find(&enrollments).Error
	if err != nil {
		return nil, fmt.Errorf("list enrollments by student: %w", err)
	}
	return enrollments, nil
}

// ListByClass retrieves all enrollments for a class.
func (r *Repository) ListByClass(ctx context.Context, tenantID, classID uuid.UUID) ([]StudentEnrollment, error) {
	var enrollments []StudentEnrollment
	err := r.db.WithContext(ctx).
		Preload("AcademicYear").
		Where("tenant_id = ? AND class_id = ? AND status = ?", tenantID, classID, EnrollmentStatusActive).
		Order("roll_number ASC").
		Find(&enrollments).Error
	if err != nil {
		return nil, fmt.Errorf("list enrollments by class: %w", err)
	}
	return enrollments, nil
}

// ListBySection retrieves all enrollments for a section.
func (r *Repository) ListBySection(ctx context.Context, tenantID, sectionID uuid.UUID) ([]StudentEnrollment, error) {
	var enrollments []StudentEnrollment
	err := r.db.WithContext(ctx).
		Preload("AcademicYear").
		Where("tenant_id = ? AND section_id = ? AND status = ?", tenantID, sectionID, EnrollmentStatusActive).
		Order("roll_number ASC").
		Find(&enrollments).Error
	if err != nil {
		return nil, fmt.Errorf("list enrollments by section: %w", err)
	}
	return enrollments, nil
}

// Update updates an enrollment in the database.
func (r *Repository) Update(ctx context.Context, enrollment *StudentEnrollment) error {
	result := r.db.WithContext(ctx).Save(enrollment)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			if strings.Contains(result.Error.Error(), "uniq_active_enrollment") {
				return ErrActiveEnrollmentExists
			}
		}
		return fmt.Errorf("update enrollment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrEnrollmentNotFound
	}
	return nil
}

// UpdateWithTx updates an enrollment within a transaction.
func (r *Repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, enrollment *StudentEnrollment) error {
	result := tx.WithContext(ctx).Save(enrollment)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			if strings.Contains(result.Error.Error(), "uniq_active_enrollment") {
				return ErrActiveEnrollmentExists
			}
		}
		return fmt.Errorf("update enrollment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrEnrollmentNotFound
	}
	return nil
}

// Delete deletes an enrollment.
func (r *Repository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&StudentEnrollment{})
	if result.Error != nil {
		return fmt.Errorf("delete enrollment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrEnrollmentNotFound
	}
	return nil
}

// CreateStatusChange creates a status change log entry.
func (r *Repository) CreateStatusChange(ctx context.Context, change *EnrollmentStatusChange) error {
	if err := r.db.WithContext(ctx).Create(change).Error; err != nil {
		return fmt.Errorf("create status change: %w", err)
	}
	return nil
}

// CreateStatusChangeWithTx creates a status change log entry within a transaction.
func (r *Repository) CreateStatusChangeWithTx(ctx context.Context, tx *gorm.DB, change *EnrollmentStatusChange) error {
	if err := tx.WithContext(ctx).Create(change).Error; err != nil {
		return fmt.Errorf("create status change: %w", err)
	}
	return nil
}

// ListStatusChanges retrieves status change history for an enrollment.
func (r *Repository) ListStatusChanges(ctx context.Context, tenantID, enrollmentID uuid.UUID) ([]EnrollmentStatusChange, error) {
	var changes []EnrollmentStatusChange
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND enrollment_id = ?", tenantID, enrollmentID).
		Order("changed_at DESC").
		Find(&changes).Error
	if err != nil {
		return nil, fmt.Errorf("list status changes: %w", err)
	}
	return changes, nil
}

// CountByAcademicYear counts enrollments by academic year and status.
func (r *Repository) CountByAcademicYear(ctx context.Context, tenantID, academicYearID uuid.UUID, status *EnrollmentStatus) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&StudentEnrollment{}).
		Where("tenant_id = ? AND academic_year_id = ?", tenantID, academicYearID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count enrollments: %w", err)
	}
	return count, nil
}

// DB returns the underlying database connection for transactions.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
