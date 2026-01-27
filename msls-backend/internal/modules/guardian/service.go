// Package guardian provides guardian and emergency contact management functionality.
package guardian

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"msls-backend/internal/pkg/database/models"
	"msls-backend/internal/pkg/logger"
)

// Service handles business logic for guardians and emergency contacts.
type Service struct {
	repo *Repository
}

// NewService creates a new guardian service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// =========================================================================
// Guardian Operations
// =========================================================================

// CreateGuardian creates a new guardian for a student.
func (s *Service) CreateGuardian(ctx context.Context, dto CreateGuardianDTO) (*models.StudentGuardian, error) {
	// Validate required fields
	if dto.FirstName == "" {
		return nil, ErrFirstNameRequired
	}
	if dto.LastName == "" {
		return nil, ErrLastNameRequired
	}
	if dto.Phone == "" {
		return nil, ErrPhoneRequired
	}
	if dto.Relation == "" || !dto.Relation.IsValid() {
		return nil, ErrInvalidRelation
	}

	// Check if student exists
	exists, err := s.repo.StudentExists(ctx, dto.TenantID, dto.StudentID)
	if err != nil {
		logger.Error("Failed to check student existence", zap.Error(err))
		return nil, err
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	// If setting as primary, clear existing primary
	if dto.IsPrimary {
		if err := s.repo.ClearPrimaryGuardian(ctx, dto.TenantID, dto.StudentID); err != nil {
			logger.Error("Failed to clear primary guardian", zap.Error(err))
			return nil, err
		}
	}

	guardian := &models.StudentGuardian{
		TenantID:        dto.TenantID,
		StudentID:       dto.StudentID,
		Relation:        dto.Relation,
		FirstName:       dto.FirstName,
		LastName:        dto.LastName,
		Phone:           dto.Phone,
		Email:           dto.Email,
		Occupation:      dto.Occupation,
		AnnualIncome:    dto.AnnualIncome,
		Education:       dto.Education,
		IsPrimary:       dto.IsPrimary,
		HasPortalAccess: dto.HasPortalAccess,
		AddressLine1:    dto.AddressLine1,
		AddressLine2:    dto.AddressLine2,
		City:            dto.City,
		State:           dto.State,
		PostalCode:      dto.PostalCode,
		Country:         dto.Country,
		CreatedBy:       dto.CreatedBy,
		UpdatedBy:       dto.CreatedBy,
	}

	if err := s.repo.CreateGuardian(ctx, guardian); err != nil {
		logger.Error("Failed to create guardian",
			zap.String("student_id", dto.StudentID.String()),
			zap.Error(err))
		return nil, err
	}

	return guardian, nil
}

// GetGuardianByID retrieves a guardian by ID.
func (s *Service) GetGuardianByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StudentGuardian, error) {
	return s.repo.GetGuardianByID(ctx, tenantID, id)
}

// GetGuardiansByStudentID retrieves all guardians for a student.
func (s *Service) GetGuardiansByStudentID(ctx context.Context, tenantID, studentID uuid.UUID) ([]models.StudentGuardian, error) {
	// Check if student exists
	exists, err := s.repo.StudentExists(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	return s.repo.GetGuardiansByStudentID(ctx, tenantID, studentID)
}

// UpdateGuardian updates an existing guardian.
func (s *Service) UpdateGuardian(ctx context.Context, tenantID, id uuid.UUID, dto UpdateGuardianDTO) (*models.StudentGuardian, error) {
	guardian, err := s.repo.GetGuardianByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Validate relation if provided
	if dto.Relation != nil && !dto.Relation.IsValid() {
		return nil, ErrInvalidRelation
	}

	// If setting as primary, clear existing primary
	if dto.IsPrimary != nil && *dto.IsPrimary && !guardian.IsPrimary {
		if err := s.repo.ClearPrimaryGuardian(ctx, tenantID, guardian.StudentID); err != nil {
			logger.Error("Failed to clear primary guardian", zap.Error(err))
			return nil, err
		}
	}

	// Apply updates
	if dto.Relation != nil {
		guardian.Relation = *dto.Relation
	}
	if dto.FirstName != nil {
		guardian.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		guardian.LastName = *dto.LastName
	}
	if dto.Phone != nil {
		guardian.Phone = *dto.Phone
	}
	if dto.Email != nil {
		guardian.Email = *dto.Email
	}
	if dto.Occupation != nil {
		guardian.Occupation = *dto.Occupation
	}
	if dto.AnnualIncome != nil {
		guardian.AnnualIncome = *dto.AnnualIncome
	}
	if dto.Education != nil {
		guardian.Education = *dto.Education
	}
	if dto.IsPrimary != nil {
		guardian.IsPrimary = *dto.IsPrimary
	}
	if dto.HasPortalAccess != nil {
		guardian.HasPortalAccess = *dto.HasPortalAccess
	}
	if dto.AddressLine1 != nil {
		guardian.AddressLine1 = *dto.AddressLine1
	}
	if dto.AddressLine2 != nil {
		guardian.AddressLine2 = *dto.AddressLine2
	}
	if dto.City != nil {
		guardian.City = *dto.City
	}
	if dto.State != nil {
		guardian.State = *dto.State
	}
	if dto.PostalCode != nil {
		guardian.PostalCode = *dto.PostalCode
	}
	if dto.Country != nil {
		guardian.Country = *dto.Country
	}
	if dto.UpdatedBy != nil {
		guardian.UpdatedBy = dto.UpdatedBy
	}

	if err := s.repo.UpdateGuardian(ctx, guardian); err != nil {
		logger.Error("Failed to update guardian",
			zap.String("guardian_id", id.String()),
			zap.Error(err))
		return nil, err
	}

	return guardian, nil
}

// DeleteGuardian deletes a guardian.
func (s *Service) DeleteGuardian(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteGuardian(ctx, tenantID, id)
}

// SetPrimaryGuardian sets a guardian as the primary guardian for a student.
func (s *Service) SetPrimaryGuardian(ctx context.Context, tenantID, id uuid.UUID) (*models.StudentGuardian, error) {
	guardian, err := s.repo.GetGuardianByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Clear existing primary
	if err := s.repo.ClearPrimaryGuardian(ctx, tenantID, guardian.StudentID); err != nil {
		return nil, err
	}

	// Set this guardian as primary
	guardian.IsPrimary = true
	if err := s.repo.UpdateGuardian(ctx, guardian); err != nil {
		return nil, err
	}

	return guardian, nil
}

// =========================================================================
// Emergency Contact Operations
// =========================================================================

// CreateEmergencyContact creates a new emergency contact for a student.
func (s *Service) CreateEmergencyContact(ctx context.Context, dto CreateEmergencyContactDTO) (*models.StudentEmergencyContact, error) {
	// Validate required fields
	if dto.Name == "" {
		return nil, ErrNameRequired
	}
	if dto.Relation == "" {
		return nil, ErrRelationRequired
	}
	if dto.Phone == "" {
		return nil, ErrPhoneRequired
	}

	// Check if student exists
	exists, err := s.repo.StudentExists(ctx, dto.TenantID, dto.StudentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	// Handle priority
	priority := dto.Priority
	if priority == 0 {
		// Get next available priority
		nextPriority, err := s.repo.GetNextPriority(ctx, dto.TenantID, dto.StudentID)
		if err != nil {
			return nil, err
		}
		priority = nextPriority
	} else {
		// Validate priority range
		if priority < 1 || priority > 5 {
			return nil, ErrInvalidPriority
		}
		// Check if priority is taken
		taken, err := s.repo.IsPriorityTaken(ctx, dto.TenantID, dto.StudentID, priority, nil)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, ErrPriorityConflict
		}
	}

	contact := &models.StudentEmergencyContact{
		TenantID:       dto.TenantID,
		StudentID:      dto.StudentID,
		Name:           dto.Name,
		Relation:       dto.Relation,
		Phone:          dto.Phone,
		AlternatePhone: dto.AlternatePhone,
		Priority:       priority,
		Notes:          dto.Notes,
		CreatedBy:      dto.CreatedBy,
		UpdatedBy:      dto.CreatedBy,
	}

	if err := s.repo.CreateEmergencyContact(ctx, contact); err != nil {
		logger.Error("Failed to create emergency contact",
			zap.String("student_id", dto.StudentID.String()),
			zap.Error(err))
		return nil, err
	}

	return contact, nil
}

// GetEmergencyContactByID retrieves an emergency contact by ID.
func (s *Service) GetEmergencyContactByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StudentEmergencyContact, error) {
	return s.repo.GetEmergencyContactByID(ctx, tenantID, id)
}

// GetEmergencyContactsByStudentID retrieves all emergency contacts for a student.
func (s *Service) GetEmergencyContactsByStudentID(ctx context.Context, tenantID, studentID uuid.UUID) ([]models.StudentEmergencyContact, error) {
	// Check if student exists
	exists, err := s.repo.StudentExists(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	return s.repo.GetEmergencyContactsByStudentID(ctx, tenantID, studentID)
}

// UpdateEmergencyContact updates an existing emergency contact.
func (s *Service) UpdateEmergencyContact(ctx context.Context, tenantID, id uuid.UUID, dto UpdateEmergencyContactDTO) (*models.StudentEmergencyContact, error) {
	contact, err := s.repo.GetEmergencyContactByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check priority if being updated
	if dto.Priority != nil {
		if *dto.Priority < 1 || *dto.Priority > 5 {
			return nil, ErrInvalidPriority
		}
		if *dto.Priority != contact.Priority {
			taken, err := s.repo.IsPriorityTaken(ctx, tenantID, contact.StudentID, *dto.Priority, &id)
			if err != nil {
				return nil, err
			}
			if taken {
				return nil, ErrPriorityConflict
			}
		}
	}

	// Apply updates
	if dto.Name != nil {
		contact.Name = *dto.Name
	}
	if dto.Relation != nil {
		contact.Relation = *dto.Relation
	}
	if dto.Phone != nil {
		contact.Phone = *dto.Phone
	}
	if dto.AlternatePhone != nil {
		contact.AlternatePhone = *dto.AlternatePhone
	}
	if dto.Priority != nil {
		contact.Priority = *dto.Priority
	}
	if dto.Notes != nil {
		contact.Notes = *dto.Notes
	}
	if dto.UpdatedBy != nil {
		contact.UpdatedBy = dto.UpdatedBy
	}

	if err := s.repo.UpdateEmergencyContact(ctx, contact); err != nil {
		logger.Error("Failed to update emergency contact",
			zap.String("contact_id", id.String()),
			zap.Error(err))
		return nil, err
	}

	return contact, nil
}

// DeleteEmergencyContact deletes an emergency contact.
func (s *Service) DeleteEmergencyContact(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteEmergencyContact(ctx, tenantID, id)
}
