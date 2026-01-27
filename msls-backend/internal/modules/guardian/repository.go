// Package guardian provides guardian and emergency contact management functionality.
package guardian

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for guardians and emergency contacts.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new guardian repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// =========================================================================
// Guardian Operations
// =========================================================================

// CreateGuardian creates a new guardian in the database.
func (r *Repository) CreateGuardian(ctx context.Context, guardian *models.StudentGuardian) error {
	return r.db.WithContext(ctx).Create(guardian).Error
}

// GetGuardianByID retrieves a guardian by ID.
func (r *Repository) GetGuardianByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StudentGuardian, error) {
	var guardian models.StudentGuardian
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&guardian).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGuardianNotFound
		}
		return nil, err
	}
	return &guardian, nil
}

// GetGuardiansByStudentID retrieves all guardians for a student.
func (r *Repository) GetGuardiansByStudentID(ctx context.Context, tenantID, studentID uuid.UUID) ([]models.StudentGuardian, error) {
	var guardians []models.StudentGuardian
	err := r.db.WithContext(ctx).
		Where("student_id = ? AND tenant_id = ?", studentID, tenantID).
		Order("is_primary DESC, relation ASC, created_at ASC").
		Find(&guardians).Error
	return guardians, err
}

// UpdateGuardian updates a guardian in the database.
func (r *Repository) UpdateGuardian(ctx context.Context, guardian *models.StudentGuardian) error {
	return r.db.WithContext(ctx).Save(guardian).Error
}

// DeleteGuardian deletes a guardian from the database.
func (r *Repository) DeleteGuardian(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.StudentGuardian{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrGuardianNotFound
	}
	return nil
}

// HasPrimaryGuardian checks if a student already has a primary guardian.
func (r *Repository) HasPrimaryGuardian(ctx context.Context, tenantID, studentID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.StudentGuardian{}).
		Where("student_id = ? AND tenant_id = ? AND is_primary = true", studentID, tenantID)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// ClearPrimaryGuardian clears the primary flag from all guardians of a student.
func (r *Repository) ClearPrimaryGuardian(ctx context.Context, tenantID, studentID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.StudentGuardian{}).
		Where("student_id = ? AND tenant_id = ?", studentID, tenantID).
		Update("is_primary", false).Error
}

// =========================================================================
// Emergency Contact Operations
// =========================================================================

// CreateEmergencyContact creates a new emergency contact in the database.
func (r *Repository) CreateEmergencyContact(ctx context.Context, contact *models.StudentEmergencyContact) error {
	return r.db.WithContext(ctx).Create(contact).Error
}

// GetEmergencyContactByID retrieves an emergency contact by ID.
func (r *Repository) GetEmergencyContactByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StudentEmergencyContact, error) {
	var contact models.StudentEmergencyContact
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmergencyContactNotFound
		}
		return nil, err
	}
	return &contact, nil
}

// GetEmergencyContactsByStudentID retrieves all emergency contacts for a student.
func (r *Repository) GetEmergencyContactsByStudentID(ctx context.Context, tenantID, studentID uuid.UUID) ([]models.StudentEmergencyContact, error) {
	var contacts []models.StudentEmergencyContact
	err := r.db.WithContext(ctx).
		Where("student_id = ? AND tenant_id = ?", studentID, tenantID).
		Order("priority ASC").
		Find(&contacts).Error
	return contacts, err
}

// UpdateEmergencyContact updates an emergency contact in the database.
func (r *Repository) UpdateEmergencyContact(ctx context.Context, contact *models.StudentEmergencyContact) error {
	return r.db.WithContext(ctx).Save(contact).Error
}

// DeleteEmergencyContact deletes an emergency contact from the database.
func (r *Repository) DeleteEmergencyContact(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.StudentEmergencyContact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrEmergencyContactNotFound
	}
	return nil
}

// IsPriorityTaken checks if a priority is already used for a student's emergency contacts.
func (r *Repository) IsPriorityTaken(ctx context.Context, tenantID, studentID uuid.UUID, priority int, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.StudentEmergencyContact{}).
		Where("student_id = ? AND tenant_id = ? AND priority = ?", studentID, tenantID, priority)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetNextPriority returns the next available priority for a student's emergency contacts.
func (r *Repository) GetNextPriority(ctx context.Context, tenantID, studentID uuid.UUID) (int, error) {
	var maxPriority int
	err := r.db.WithContext(ctx).
		Model(&models.StudentEmergencyContact{}).
		Where("student_id = ? AND tenant_id = ?", studentID, tenantID).
		Select("COALESCE(MAX(priority), 0)").
		Scan(&maxPriority).Error
	if err != nil {
		return 0, err
	}
	return maxPriority + 1, nil
}

// =========================================================================
// Student Validation
// =========================================================================

// StudentExists checks if a student exists.
func (r *Repository) StudentExists(ctx context.Context, tenantID, studentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ? AND tenant_id = ?", studentID, tenantID).
		Count(&count).Error
	return count > 0, err
}
