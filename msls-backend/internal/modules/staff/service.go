// Package staff provides staff management functionality.
package staff

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
	"msls-backend/internal/services/branch"
)

// Service handles staff business logic.
type Service struct {
	repo          *Repository
	branchService *branch.Service
	db            *gorm.DB
}

// NewService creates a new staff service.
func NewService(db *gorm.DB, branchService *branch.Service) *Service {
	return &Service{
		repo:          NewRepository(db),
		branchService: branchService,
		db:            db,
	}
}

// Create creates a new staff member with auto-generated employee ID.
func (s *Service) Create(ctx context.Context, dto CreateStaffDTO) (*models.Staff, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.BranchID == uuid.Nil {
		return nil, ErrBranchIDRequired
	}
	if dto.FirstName == "" {
		return nil, ErrFirstNameRequired
	}
	if dto.LastName == "" {
		return nil, ErrLastNameRequired
	}
	if dto.DateOfBirth.IsZero() {
		return nil, ErrDateOfBirthRequired
	}
	if dto.DateOfBirth.After(time.Now()) {
		return nil, ErrInvalidDateOfBirth
	}
	if dto.Gender == "" {
		return nil, ErrGenderRequired
	}
	if !dto.Gender.IsValid() {
		return nil, ErrInvalidGender
	}
	if dto.WorkEmail == "" {
		return nil, ErrWorkEmailRequired
	}
	if dto.WorkPhone == "" {
		return nil, ErrWorkPhoneRequired
	}
	if !dto.StaffType.IsValid() {
		return nil, ErrInvalidStaffType
	}
	if dto.JoinDate.IsZero() {
		return nil, ErrJoinDateRequired
	}

	// Verify branch exists
	_, err := s.branchService.GetByID(ctx, dto.TenantID, dto.BranchID)
	if err != nil {
		return nil, ErrBranchNotFound
	}

	// Verify reporting manager exists if provided
	if dto.ReportingManagerID != nil {
		_, err := s.repo.GetByID(ctx, dto.TenantID, *dto.ReportingManagerID)
		if err != nil {
			return nil, ErrReportingManagerNotFound
		}
	}

	var staff *models.Staff
	prefix := "EMP" // Default prefix

	// Use transaction to ensure atomic creation
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Generate employee ID
		sequence, err := s.repo.GetNextSequence(ctx, tx, dto.TenantID, prefix)
		if err != nil {
			return fmt.Errorf("generate employee id: %w", err)
		}

		employeeID := fmt.Sprintf("%s%05d", prefix, sequence)

		// Create staff
		staff = &models.Staff{
			TenantID:         dto.TenantID,
			BranchID:         dto.BranchID,
			EmployeeID:       employeeID,
			EmployeeIDPrefix: prefix,
			FirstName:        dto.FirstName,
			MiddleName:       dto.MiddleName,
			LastName:         dto.LastName,
			DateOfBirth:      dto.DateOfBirth,
			Gender:           dto.Gender,
			BloodGroup:       dto.BloodGroup,
			Nationality:      dto.Nationality,
			Religion:         dto.Religion,
			MaritalStatus:    dto.MaritalStatus,
			PersonalEmail:    dto.PersonalEmail,
			WorkEmail:        dto.WorkEmail,
			PersonalPhone:    dto.PersonalPhone,
			WorkPhone:        dto.WorkPhone,
			EmergencyContactName:     dto.EmergencyContactName,
			EmergencyContactPhone:    dto.EmergencyContactPhone,
			EmergencyContactRelation: dto.EmergencyContactRelation,
			StaffType:          dto.StaffType,
			DepartmentID:       dto.DepartmentID,
			DesignationID:      dto.DesignationID,
			ReportingManagerID: dto.ReportingManagerID,
			JoinDate:           dto.JoinDate,
			ConfirmationDate:   dto.ConfirmationDate,
			ProbationEndDate:   dto.ProbationEndDate,
			Bio:                dto.Bio,
			Status:             models.StaffStatusActive,
			CreatedBy:          dto.CreatedBy,
			UpdatedBy:          dto.CreatedBy,
		}

		// Set nationality default
		if staff.Nationality == "" {
			staff.Nationality = "Indian"
		}

		// Handle addresses
		if dto.CurrentAddress != nil {
			staff.CurrentAddressLine1 = dto.CurrentAddress.AddressLine1
			staff.CurrentAddressLine2 = dto.CurrentAddress.AddressLine2
			staff.CurrentCity = dto.CurrentAddress.City
			staff.CurrentState = dto.CurrentAddress.State
			staff.CurrentPincode = dto.CurrentAddress.Pincode
			staff.CurrentCountry = dto.CurrentAddress.Country
			if staff.CurrentCountry == "" {
				staff.CurrentCountry = "India"
			}
		}

		if dto.SameAsCurrent && dto.CurrentAddress != nil {
			staff.SameAsCurrent = true
			staff.PermanentAddressLine1 = dto.CurrentAddress.AddressLine1
			staff.PermanentAddressLine2 = dto.CurrentAddress.AddressLine2
			staff.PermanentCity = dto.CurrentAddress.City
			staff.PermanentState = dto.CurrentAddress.State
			staff.PermanentPincode = dto.CurrentAddress.Pincode
			staff.PermanentCountry = dto.CurrentAddress.Country
			if staff.PermanentCountry == "" {
				staff.PermanentCountry = "India"
			}
		} else if dto.PermanentAddress != nil {
			staff.PermanentAddressLine1 = dto.PermanentAddress.AddressLine1
			staff.PermanentAddressLine2 = dto.PermanentAddress.AddressLine2
			staff.PermanentCity = dto.PermanentAddress.City
			staff.PermanentState = dto.PermanentAddress.State
			staff.PermanentPincode = dto.PermanentAddress.Pincode
			staff.PermanentCountry = dto.PermanentAddress.Country
			if staff.PermanentCountry == "" {
				staff.PermanentCountry = "India"
			}
		}

		if err := s.repo.CreateWithTx(ctx, tx, staff); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch the complete staff with all relations
	return s.GetByID(ctx, dto.TenantID, staff.ID)
}

// GetByID retrieves a staff member by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Staff, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetByEmployeeID retrieves a staff member by employee ID.
func (s *Service) GetByEmployeeID(ctx context.Context, tenantID uuid.UUID, employeeID string) (*models.Staff, error) {
	return s.repo.GetByEmployeeID(ctx, tenantID, employeeID)
}

// Update updates a staff member.
func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, dto UpdateStaffDTO) (*models.Staff, error) {
	// Get existing staff
	staff, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check version for optimistic locking
	if dto.Version != 0 && staff.Version != dto.Version {
		return nil, ErrOptimisticLockConflict
	}

	// Verify reporting manager if changed
	if dto.ReportingManagerID != nil && *dto.ReportingManagerID != uuid.Nil {
		if staff.ReportingManagerID == nil || *dto.ReportingManagerID != *staff.ReportingManagerID {
			_, err := s.repo.GetByID(ctx, tenantID, *dto.ReportingManagerID)
			if err != nil {
				return nil, ErrReportingManagerNotFound
			}
		}
	}

	// Use transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update fields
		if dto.FirstName != nil {
			staff.FirstName = *dto.FirstName
		}
		if dto.MiddleName != nil {
			staff.MiddleName = *dto.MiddleName
		}
		if dto.LastName != nil {
			staff.LastName = *dto.LastName
		}
		if dto.DateOfBirth != nil {
			if dto.DateOfBirth.After(time.Now()) {
				return ErrInvalidDateOfBirth
			}
			staff.DateOfBirth = *dto.DateOfBirth
		}
		if dto.Gender != nil {
			if !dto.Gender.IsValid() {
				return ErrInvalidGender
			}
			staff.Gender = *dto.Gender
		}
		if dto.BloodGroup != nil {
			staff.BloodGroup = *dto.BloodGroup
		}
		if dto.Nationality != nil {
			staff.Nationality = *dto.Nationality
		}
		if dto.Religion != nil {
			staff.Religion = *dto.Religion
		}
		if dto.MaritalStatus != nil {
			staff.MaritalStatus = *dto.MaritalStatus
		}
		if dto.PersonalEmail != nil {
			staff.PersonalEmail = *dto.PersonalEmail
		}
		if dto.WorkEmail != nil {
			staff.WorkEmail = *dto.WorkEmail
		}
		if dto.PersonalPhone != nil {
			staff.PersonalPhone = *dto.PersonalPhone
		}
		if dto.WorkPhone != nil {
			staff.WorkPhone = *dto.WorkPhone
		}
		if dto.EmergencyContactName != nil {
			staff.EmergencyContactName = *dto.EmergencyContactName
		}
		if dto.EmergencyContactPhone != nil {
			staff.EmergencyContactPhone = *dto.EmergencyContactPhone
		}
		if dto.EmergencyContactRelation != nil {
			staff.EmergencyContactRelation = *dto.EmergencyContactRelation
		}
		if dto.StaffType != nil {
			if !dto.StaffType.IsValid() {
				return ErrInvalidStaffType
			}
			staff.StaffType = *dto.StaffType
		}
		if dto.DepartmentID != nil {
			staff.DepartmentID = dto.DepartmentID
		}
		if dto.DesignationID != nil {
			staff.DesignationID = dto.DesignationID
		}
		if dto.ReportingManagerID != nil {
			staff.ReportingManagerID = dto.ReportingManagerID
		}
		if dto.ConfirmationDate != nil {
			staff.ConfirmationDate = dto.ConfirmationDate
		}
		if dto.ProbationEndDate != nil {
			staff.ProbationEndDate = dto.ProbationEndDate
		}
		if dto.Bio != nil {
			staff.Bio = *dto.Bio
		}
		if dto.PhotoURL != nil {
			staff.PhotoURL = *dto.PhotoURL
		}

		// Handle addresses
		if dto.CurrentAddress != nil {
			staff.CurrentAddressLine1 = dto.CurrentAddress.AddressLine1
			staff.CurrentAddressLine2 = dto.CurrentAddress.AddressLine2
			staff.CurrentCity = dto.CurrentAddress.City
			staff.CurrentState = dto.CurrentAddress.State
			staff.CurrentPincode = dto.CurrentAddress.Pincode
			if dto.CurrentAddress.Country != "" {
				staff.CurrentCountry = dto.CurrentAddress.Country
			}
		}

		if dto.SameAsCurrent != nil && *dto.SameAsCurrent && dto.CurrentAddress != nil {
			staff.SameAsCurrent = true
			staff.PermanentAddressLine1 = dto.CurrentAddress.AddressLine1
			staff.PermanentAddressLine2 = dto.CurrentAddress.AddressLine2
			staff.PermanentCity = dto.CurrentAddress.City
			staff.PermanentState = dto.CurrentAddress.State
			staff.PermanentPincode = dto.CurrentAddress.Pincode
			staff.PermanentCountry = dto.CurrentAddress.Country
		} else if dto.PermanentAddress != nil {
			staff.SameAsCurrent = false
			staff.PermanentAddressLine1 = dto.PermanentAddress.AddressLine1
			staff.PermanentAddressLine2 = dto.PermanentAddress.AddressLine2
			staff.PermanentCity = dto.PermanentAddress.City
			staff.PermanentState = dto.PermanentAddress.State
			staff.PermanentPincode = dto.PermanentAddress.Pincode
			if dto.PermanentAddress.Country != "" {
				staff.PermanentCountry = dto.PermanentAddress.Country
			}
		}

		staff.UpdatedBy = dto.UpdatedBy

		if err := s.repo.UpdateWithTx(ctx, tx, staff); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch updated staff
	return s.GetByID(ctx, tenantID, id)
}

// Delete soft deletes a staff member.
func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List retrieves staff with optional filtering.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.Staff, string, int64, error) {
	return s.repo.List(ctx, filter)
}

// Count returns the total number of staff for a tenant.
func (s *Service) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, tenantID)
}

// UpdateStatus updates the status of a staff member with history tracking.
func (s *Service) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, dto StatusUpdateDTO) (*models.Staff, error) {
	if !dto.Status.IsValid() {
		return nil, ErrInvalidStatus
	}
	if dto.Reason == "" {
		return nil, ErrStatusReasonRequired
	}
	if dto.EffectiveDate.IsZero() {
		return nil, ErrEffectiveDateRequired
	}

	// Get current staff to record old status
	staff, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Use transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create status history record
		history := &models.StaffStatusHistory{
			TenantID:      tenantID,
			StaffID:       id,
			OldStatus:     string(staff.Status),
			NewStatus:     string(dto.Status),
			Reason:        dto.Reason,
			EffectiveDate: dto.EffectiveDate,
			ChangedBy:     dto.UpdatedBy,
		}
		if err := s.repo.CreateStatusHistory(ctx, tx, history); err != nil {
			return err
		}

		// Update staff status
		updates := map[string]interface{}{
			"status":        dto.Status,
			"status_reason": dto.Reason,
			"updated_by":    dto.UpdatedBy,
			"updated_at":    time.Now(),
		}

		if dto.Status == models.StaffStatusTerminated {
			updates["termination_date"] = dto.EffectiveDate
		}

		result := tx.WithContext(ctx).
			Model(&models.Staff{}).
			Where("tenant_id = ? AND id = ?", tenantID, id).
			Updates(updates)

		if result.Error != nil {
			return fmt.Errorf("update status: %w", result.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, id)
}

// GetStatusHistory retrieves the status history for a staff member.
func (s *Service) GetStatusHistory(ctx context.Context, tenantID, staffID uuid.UUID) ([]models.StaffStatusHistory, error) {
	return s.repo.GetStatusHistory(ctx, tenantID, staffID)
}

// UpdatePhoto updates the photo URL for a staff member.
func (s *Service) UpdatePhoto(ctx context.Context, tenantID, id uuid.UUID, photoURL string, updatedBy *uuid.UUID) (*models.Staff, error) {
	if err := s.repo.UpdatePhoto(ctx, tenantID, id, photoURL, updatedBy); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, id)
}

// GenerateEmployeeID generates an employee ID without creating a staff member.
// Useful for previewing the next employee ID.
func (s *Service) GenerateEmployeeID(ctx context.Context, tenantID uuid.UUID) (string, error) {
	prefix := "EMP"

	// Get current sequence without incrementing
	var seq models.StaffEmployeeSequence
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND prefix = ?", tenantID, prefix).
		First(&seq).Error

	nextSeq := 1
	if err == nil {
		nextSeq = seq.LastSequence + 1
	}

	return fmt.Sprintf("%s%05d", prefix, nextSeq), nil
}
