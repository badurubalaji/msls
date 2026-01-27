// Package student provides student management functionality.
package student

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for students.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new student repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new student in the database.
func (r *Repository) Create(ctx context.Context, student *models.Student) error {
	if err := r.db.WithContext(ctx).Create(student).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "admission_number") {
				return ErrDuplicateAdmissionNumber
			}
		}
		return fmt.Errorf("create student: %w", err)
	}
	return nil
}

// CreateWithTx creates a new student in the database within a transaction.
func (r *Repository) CreateWithTx(ctx context.Context, tx *gorm.DB, student *models.Student) error {
	if err := tx.WithContext(ctx).Create(student).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "admission_number") {
				return ErrDuplicateAdmissionNumber
			}
		}
		return fmt.Errorf("create student: %w", err)
	}
	return nil
}

// GetByID retrieves a student by ID.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Addresses").
		Preload("Branch").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("get student by id: %w", err)
	}
	return &student, nil
}

// GetByAdmissionNumber retrieves a student by admission number.
func (r *Repository) GetByAdmissionNumber(ctx context.Context, tenantID uuid.UUID, admissionNumber string) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Addresses").
		Preload("Branch").
		Where("tenant_id = ? AND admission_number = ?", tenantID, admissionNumber).
		First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("get student by admission number: %w", err)
	}
	return &student, nil
}

// Update updates a student in the database.
func (r *Repository) Update(ctx context.Context, student *models.Student) error {
	// Use current version for optimistic locking check
	// BeforeUpdate hook will increment version during the update
	result := r.db.WithContext(ctx).
		Model(student).
		Where("version = ?", student.Version). // Optimistic locking
		Updates(student)

	if result.Error != nil {
		return fmt.Errorf("update student: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrOptimisticLockConflict
	}

	return nil
}

// UpdateWithTx updates a student within a transaction.
func (r *Repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, student *models.Student) error {
	// Use current version for optimistic locking check
	// BeforeUpdate hook will increment version during the update
	result := tx.WithContext(ctx).
		Model(student).
		Where("version = ?", student.Version). // Optimistic locking
		Updates(student)

	if result.Error != nil {
		return fmt.Errorf("update student: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrOptimisticLockConflict
	}

	return nil
}

// Delete soft deletes a student.
func (r *Repository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Student{})

	if result.Error != nil {
		return fmt.Errorf("delete student: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}

	return nil
}

// List retrieves students with optional filtering and pagination.
func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.Student, string, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Preload("Branch").
		Where("students.tenant_id = ?", filter.TenantID)

	// Apply filters
	if filter.BranchID != nil {
		query = query.Where("students.branch_id = ?", *filter.BranchID)
	}

	if filter.Status != nil {
		query = query.Where("students.status = ?", *filter.Status)
	}

	if filter.Gender != nil {
		query = query.Where("students.gender = ?", *filter.Gender)
	}

	// Filter by class/section via enrollment join
	if filter.ClassID != nil || filter.SectionID != nil {
		query = query.Joins("JOIN student_enrollments ON students.id = student_enrollments.student_id AND student_enrollments.status = 'active'")
		if filter.ClassID != nil {
			query = query.Where("student_enrollments.class_id = ?", *filter.ClassID)
		}
		if filter.SectionID != nil {
			query = query.Where("student_enrollments.section_id = ?", *filter.SectionID)
		}
	}

	// Filter by admission date range
	if filter.AdmissionFrom != nil {
		query = query.Where("students.admission_date >= ?", *filter.AdmissionFrom)
	}
	if filter.AdmissionTo != nil {
		query = query.Where("students.admission_date <= ?", *filter.AdmissionTo)
	}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"LOWER(students.first_name) LIKE ? OR LOWER(students.last_name) LIKE ? OR LOWER(students.first_name || ' ' || students.last_name) LIKE ? OR LOWER(students.admission_number) LIKE ?",
			search, search, search, search,
		)
	}

	// Get total count (use distinct to avoid counting duplicates from join)
	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Distinct("students.id").Count(&total).Error; err != nil {
		return nil, "", 0, fmt.Errorf("count students: %w", err)
	}

	// Apply cursor pagination
	if filter.Cursor != "" {
		cursorID, err := uuid.Parse(filter.Cursor)
		if err == nil {
			query = query.Where("students.id > ?", cursorID)
		}
	}

	// Apply sorting
	sortBy := "students.last_name, students.first_name"
	if filter.SortBy != "" {
		sortOrder := "ASC"
		if strings.ToLower(filter.SortOrder) == "desc" {
			sortOrder = "DESC"
		}
		sortBy = fmt.Sprintf("students.%s %s", filter.SortBy, sortOrder)
	}
	query = query.Order(sortBy)

	// Apply limit
	limit := 20
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Distinct("students.*").Limit(limit + 1)

	var students []models.Student
	if err := query.Find(&students).Error; err != nil {
		return nil, "", 0, fmt.Errorf("list students: %w", err)
	}

	// Calculate next cursor
	var nextCursor string
	if len(students) > limit {
		students = students[:limit]
		nextCursor = students[len(students)-1].ID.String()
	}

	return students, nextCursor, total, nil
}

// Count returns the total number of students for a tenant.
func (r *Repository) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count students: %w", err)
	}
	return count, nil
}

// GetNextSequence gets and increments the admission number sequence.
func (r *Repository) GetNextSequence(ctx context.Context, tx *gorm.DB, tenantID, branchID uuid.UUID, year int) (int, error) {
	var seq models.StudentAdmissionSequence

	// Try to find existing sequence
	err := tx.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ? AND year = ?", tenantID, branchID, year).
		First(&seq).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new sequence
		seq = models.StudentAdmissionSequence{
			TenantID:     tenantID,
			BranchID:     branchID,
			Year:         year,
			LastSequence: 1,
		}
		if err := tx.WithContext(ctx).Create(&seq).Error; err != nil {
			return 0, fmt.Errorf("create sequence: %w", err)
		}
		return seq.LastSequence, nil
	}

	if err != nil {
		return 0, fmt.Errorf("get sequence: %w", err)
	}

	// Increment sequence
	seq.LastSequence++
	seq.UpdatedAt = time.Now()
	if err := tx.WithContext(ctx).Save(&seq).Error; err != nil {
		return 0, fmt.Errorf("update sequence: %w", err)
	}

	return seq.LastSequence, nil
}

// CreateAddress creates a new address for a student.
func (r *Repository) CreateAddress(ctx context.Context, tx *gorm.DB, address *models.StudentAddress) error {
	if err := tx.WithContext(ctx).Create(address).Error; err != nil {
		return fmt.Errorf("create address: %w", err)
	}
	return nil
}

// UpdateAddress updates an existing address.
func (r *Repository) UpdateAddress(ctx context.Context, tx *gorm.DB, address *models.StudentAddress) error {
	if err := tx.WithContext(ctx).Save(address).Error; err != nil {
		return fmt.Errorf("update address: %w", err)
	}
	return nil
}

// DeleteAddresses deletes all addresses for a student.
func (r *Repository) DeleteAddresses(ctx context.Context, tx *gorm.DB, studentID uuid.UUID) error {
	if err := tx.WithContext(ctx).
		Where("student_id = ?", studentID).
		Delete(&models.StudentAddress{}).Error; err != nil {
		return fmt.Errorf("delete addresses: %w", err)
	}
	return nil
}

// GetAddressByType gets an address by student ID and type.
func (r *Repository) GetAddressByType(ctx context.Context, studentID uuid.UUID, addressType models.AddressType) (*models.StudentAddress, error) {
	var address models.StudentAddress
	err := r.db.WithContext(ctx).
		Where("student_id = ? AND address_type = ?", studentID, addressType).
		First(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("get address: %w", err)
	}
	return &address, nil
}

// UpdateStatus updates the status of a student.
func (r *Repository) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status models.StudentStatus, updatedBy *uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("update status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}

	return nil
}

// UpdatePhoto updates the photo URL for a student.
func (r *Repository) UpdatePhoto(ctx context.Context, tenantID, id uuid.UUID, photoURL string, updatedBy *uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(map[string]interface{}{
			"photo_url":  photoURL,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("update photo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}

	return nil
}

// DB returns the underlying database connection for transactions.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
