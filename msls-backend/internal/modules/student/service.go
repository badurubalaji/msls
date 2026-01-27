// Package student provides student management functionality.
package student

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
	"msls-backend/internal/services/branch"
)

// Service handles student business logic.
type Service struct {
	repo          *Repository
	branchService *branch.Service
	db            *gorm.DB
}

// NewService creates a new student service.
func NewService(db *gorm.DB, branchService *branch.Service) *Service {
	return &Service{
		repo:          NewRepository(db),
		branchService: branchService,
		db:            db,
	}
}

// Create creates a new student with auto-generated admission number.
func (s *Service) Create(ctx context.Context, dto CreateStudentDTO) (*models.Student, error) {
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

	// Verify branch exists
	branchEntity, err := s.branchService.GetByID(ctx, dto.TenantID, dto.BranchID)
	if err != nil {
		return nil, ErrBranchNotFound
	}

	var student *models.Student

	// Use transaction to ensure atomic creation
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Generate admission number
		year := time.Now().Year()
		sequence, err := s.repo.GetNextSequence(ctx, tx, dto.TenantID, dto.BranchID, year)
		if err != nil {
			return fmt.Errorf("generate admission number: %w", err)
		}

		admissionNumber := fmt.Sprintf("%s-%d-%05d", branchEntity.Code, year, sequence)

		// Set admission date
		admissionDate := time.Now()
		if dto.AdmissionDate != nil {
			admissionDate = *dto.AdmissionDate
		}

		// Create student
		student = &models.Student{
			TenantID:        dto.TenantID,
			BranchID:        dto.BranchID,
			AdmissionNumber: admissionNumber,
			FirstName:       dto.FirstName,
			MiddleName:      dto.MiddleName,
			LastName:        dto.LastName,
			DateOfBirth:     dto.DateOfBirth,
			Gender:          dto.Gender,
			BloodGroup:      dto.BloodGroup,
			AadhaarNumber:   dto.AadhaarNumber,
			Status:          models.StudentStatusActive,
			AdmissionDate:   admissionDate,
			CreatedBy:       dto.CreatedBy,
			UpdatedBy:       dto.CreatedBy,
		}

		if err := s.repo.CreateWithTx(ctx, tx, student); err != nil {
			return err
		}

		// Create addresses
		if dto.CurrentAddress != nil {
			currentAddr := &models.StudentAddress{
				TenantID:     dto.TenantID,
				StudentID:    student.ID,
				AddressType:  models.AddressTypeCurrent,
				AddressLine1: dto.CurrentAddress.AddressLine1,
				AddressLine2: dto.CurrentAddress.AddressLine2,
				City:         dto.CurrentAddress.City,
				State:        dto.CurrentAddress.State,
				PostalCode:   dto.CurrentAddress.PostalCode,
				Country:      dto.CurrentAddress.Country,
			}
			if currentAddr.Country == "" {
				currentAddr.Country = "India"
			}
			if err := s.repo.CreateAddress(ctx, tx, currentAddr); err != nil {
				return err
			}
			student.Addresses = append(student.Addresses, *currentAddr)
		}

		if dto.PermanentAddress != nil && !dto.SameAsCurrentAddress {
			permAddr := &models.StudentAddress{
				TenantID:     dto.TenantID,
				StudentID:    student.ID,
				AddressType:  models.AddressTypePermanent,
				AddressLine1: dto.PermanentAddress.AddressLine1,
				AddressLine2: dto.PermanentAddress.AddressLine2,
				City:         dto.PermanentAddress.City,
				State:        dto.PermanentAddress.State,
				PostalCode:   dto.PermanentAddress.PostalCode,
				Country:      dto.PermanentAddress.Country,
			}
			if permAddr.Country == "" {
				permAddr.Country = "India"
			}
			if err := s.repo.CreateAddress(ctx, tx, permAddr); err != nil {
				return err
			}
			student.Addresses = append(student.Addresses, *permAddr)
		} else if dto.SameAsCurrentAddress && dto.CurrentAddress != nil {
			// Copy current address as permanent
			permAddr := &models.StudentAddress{
				TenantID:     dto.TenantID,
				StudentID:    student.ID,
				AddressType:  models.AddressTypePermanent,
				AddressLine1: dto.CurrentAddress.AddressLine1,
				AddressLine2: dto.CurrentAddress.AddressLine2,
				City:         dto.CurrentAddress.City,
				State:        dto.CurrentAddress.State,
				PostalCode:   dto.CurrentAddress.PostalCode,
				Country:      dto.CurrentAddress.Country,
			}
			if permAddr.Country == "" {
				permAddr.Country = "India"
			}
			if err := s.repo.CreateAddress(ctx, tx, permAddr); err != nil {
				return err
			}
			student.Addresses = append(student.Addresses, *permAddr)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch the complete student with all relations
	return s.GetByID(ctx, dto.TenantID, student.ID)
}

// GetByID retrieves a student by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Student, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetByAdmissionNumber retrieves a student by admission number.
func (s *Service) GetByAdmissionNumber(ctx context.Context, tenantID uuid.UUID, admissionNumber string) (*models.Student, error) {
	return s.repo.GetByAdmissionNumber(ctx, tenantID, admissionNumber)
}

// Update updates a student.
func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, dto UpdateStudentDTO) (*models.Student, error) {
	// Get existing student
	student, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check version for optimistic locking
	if dto.Version != 0 && student.Version != dto.Version {
		return nil, ErrOptimisticLockConflict
	}

	// Use transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update fields
		if dto.FirstName != nil {
			student.FirstName = *dto.FirstName
		}
		if dto.MiddleName != nil {
			student.MiddleName = *dto.MiddleName
		}
		if dto.LastName != nil {
			student.LastName = *dto.LastName
		}
		if dto.DateOfBirth != nil {
			if dto.DateOfBirth.After(time.Now()) {
				return ErrInvalidDateOfBirth
			}
			student.DateOfBirth = *dto.DateOfBirth
		}
		if dto.Gender != nil {
			if !dto.Gender.IsValid() {
				return ErrInvalidGender
			}
			student.Gender = *dto.Gender
		}
		if dto.BloodGroup != nil {
			student.BloodGroup = *dto.BloodGroup
		}
		if dto.AadhaarNumber != nil {
			student.AadhaarNumber = *dto.AadhaarNumber
		}
		if dto.Status != nil {
			if !dto.Status.IsValid() {
				return ErrInvalidStatus
			}
			student.Status = *dto.Status
		}
		if dto.PhotoURL != nil {
			student.PhotoURL = *dto.PhotoURL
		}
		if dto.BirthCertificateURL != nil {
			student.BirthCertificateURL = *dto.BirthCertificateURL
		}
		student.UpdatedBy = dto.UpdatedBy

		if err := s.repo.UpdateWithTx(ctx, tx, student); err != nil {
			return err
		}

		// Update addresses if provided
		if dto.CurrentAddress != nil {
			existingAddr, err := s.repo.GetAddressByType(ctx, student.ID, models.AddressTypeCurrent)
			if err == ErrAddressNotFound {
				// Create new address
				currentAddr := &models.StudentAddress{
					TenantID:     tenantID,
					StudentID:    student.ID,
					AddressType:  models.AddressTypeCurrent,
					AddressLine1: dto.CurrentAddress.AddressLine1,
					AddressLine2: dto.CurrentAddress.AddressLine2,
					City:         dto.CurrentAddress.City,
					State:        dto.CurrentAddress.State,
					PostalCode:   dto.CurrentAddress.PostalCode,
					Country:      dto.CurrentAddress.Country,
				}
				if currentAddr.Country == "" {
					currentAddr.Country = "India"
				}
				if err := s.repo.CreateAddress(ctx, tx, currentAddr); err != nil {
					return err
				}
			} else if err == nil {
				// Update existing address
				existingAddr.AddressLine1 = dto.CurrentAddress.AddressLine1
				existingAddr.AddressLine2 = dto.CurrentAddress.AddressLine2
				existingAddr.City = dto.CurrentAddress.City
				existingAddr.State = dto.CurrentAddress.State
				existingAddr.PostalCode = dto.CurrentAddress.PostalCode
				if dto.CurrentAddress.Country != "" {
					existingAddr.Country = dto.CurrentAddress.Country
				}
				if err := s.repo.UpdateAddress(ctx, tx, existingAddr); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		if dto.PermanentAddress != nil && (dto.SameAsCurrentAddress == nil || !*dto.SameAsCurrentAddress) {
			existingAddr, err := s.repo.GetAddressByType(ctx, student.ID, models.AddressTypePermanent)
			if err == ErrAddressNotFound {
				// Create new address
				permAddr := &models.StudentAddress{
					TenantID:     tenantID,
					StudentID:    student.ID,
					AddressType:  models.AddressTypePermanent,
					AddressLine1: dto.PermanentAddress.AddressLine1,
					AddressLine2: dto.PermanentAddress.AddressLine2,
					City:         dto.PermanentAddress.City,
					State:        dto.PermanentAddress.State,
					PostalCode:   dto.PermanentAddress.PostalCode,
					Country:      dto.PermanentAddress.Country,
				}
				if permAddr.Country == "" {
					permAddr.Country = "India"
				}
				if err := s.repo.CreateAddress(ctx, tx, permAddr); err != nil {
					return err
				}
			} else if err == nil {
				// Update existing address
				existingAddr.AddressLine1 = dto.PermanentAddress.AddressLine1
				existingAddr.AddressLine2 = dto.PermanentAddress.AddressLine2
				existingAddr.City = dto.PermanentAddress.City
				existingAddr.State = dto.PermanentAddress.State
				existingAddr.PostalCode = dto.PermanentAddress.PostalCode
				if dto.PermanentAddress.Country != "" {
					existingAddr.Country = dto.PermanentAddress.Country
				}
				if err := s.repo.UpdateAddress(ctx, tx, existingAddr); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch updated student
	return s.GetByID(ctx, tenantID, id)
}

// Delete soft deletes a student.
func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List retrieves students with optional filtering.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.Student, string, int64, error) {
	return s.repo.List(ctx, filter)
}

// Count returns the total number of students for a tenant.
func (s *Service) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, tenantID)
}

// UpdateStatus updates the status of a student.
func (s *Service) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status models.StudentStatus, updatedBy *uuid.UUID) (*models.Student, error) {
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if err := s.repo.UpdateStatus(ctx, tenantID, id, status, updatedBy); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, id)
}

// UpdatePhoto updates the photo URL for a student.
func (s *Service) UpdatePhoto(ctx context.Context, tenantID, id uuid.UUID, photoURL string, updatedBy *uuid.UUID) (*models.Student, error) {
	if err := s.repo.UpdatePhoto(ctx, tenantID, id, photoURL, updatedBy); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, tenantID, id)
}

// GenerateAdmissionNumber generates an admission number without creating a student.
// Useful for previewing the next admission number.
func (s *Service) GenerateAdmissionNumber(ctx context.Context, tenantID, branchID uuid.UUID) (string, error) {
	branchEntity, err := s.branchService.GetByID(ctx, tenantID, branchID)
	if err != nil {
		return "", ErrBranchNotFound
	}

	year := time.Now().Year()

	// Get current sequence without incrementing
	var seq models.StudentAdmissionSequence
	err = s.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ? AND year = ?", tenantID, branchID, year).
		First(&seq).Error

	nextSeq := 1
	if err == nil {
		nextSeq = seq.LastSequence + 1
	}

	return fmt.Sprintf("%s-%d-%05d", branchEntity.Code, year, nextSeq), nil
}
