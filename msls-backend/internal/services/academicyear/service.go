// Package academicyear provides academic year management services.
package academicyear

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Service handles academic year-related operations.
type Service struct {
	db *gorm.DB
}

// NewService creates a new AcademicYear Service instance.
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ============================================================================
// Academic Year Request/Response Types
// ============================================================================

// CreateAcademicYearRequest represents a request to create an academic year.
type CreateAcademicYearRequest struct {
	TenantID  uuid.UUID
	BranchID  *uuid.UUID
	Name      string
	StartDate time.Time
	EndDate   time.Time
	IsCurrent bool
	CreatedBy *uuid.UUID
}

// UpdateAcademicYearRequest represents a request to update an academic year.
type UpdateAcademicYearRequest struct {
	Name      *string
	StartDate *time.Time
	EndDate   *time.Time
	IsActive  *bool
	UpdatedBy *uuid.UUID
}

// ListAcademicYearFilter contains filters for listing academic years.
type ListAcademicYearFilter struct {
	TenantID  uuid.UUID
	BranchID  *uuid.UUID
	IsCurrent *bool
	IsActive  *bool
	Search    string
}

// ============================================================================
// Academic Year Operations
// ============================================================================

// Create creates a new academic year.
func (s *Service) Create(ctx context.Context, req CreateAcademicYearRequest) (*models.AcademicYear, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.Name == "" {
		return nil, ErrAcademicYearNameRequired
	}
	if req.StartDate.IsZero() {
		return nil, ErrAcademicYearStartDateRequired
	}
	if req.EndDate.IsZero() {
		return nil, ErrAcademicYearEndDateRequired
	}
	if !req.EndDate.After(req.StartDate) {
		return nil, ErrAcademicYearInvalidDates
	}

	// Check if name already exists for this tenant/branch
	var existing models.AcademicYear
	query := s.db.WithContext(ctx).
		Where("tenant_id = ? AND name = ?", req.TenantID, req.Name)
	if req.BranchID != nil {
		query = query.Where("branch_id = ?", *req.BranchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}
	err := query.First(&existing).Error
	if err == nil {
		return nil, ErrAcademicYearNameExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing academic year: %w", err)
	}

	academicYear := &models.AcademicYear{
		TenantID:  req.TenantID,
		BranchID:  req.BranchID,
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsCurrent: req.IsCurrent,
		IsActive:  true,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
	}

	if err := s.db.WithContext(ctx).Create(academicYear).Error; err != nil {
		return nil, fmt.Errorf("failed to create academic year: %w", err)
	}

	return s.GetByID(ctx, req.TenantID, academicYear.ID)
}

// GetByID retrieves an academic year by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	err := s.db.WithContext(ctx).
		Preload("Terms", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Preload("Holidays", func(db *gorm.DB) *gorm.DB {
			return db.Order("date ASC")
		}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&academicYear).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		return nil, fmt.Errorf("failed to get academic year: %w", err)
	}
	return &academicYear, nil
}

// GetCurrent retrieves the current academic year for the tenant.
func (s *Service) GetCurrent(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	query := s.db.WithContext(ctx).
		Preload("Terms", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Preload("Holidays", func(db *gorm.DB) *gorm.DB {
			return db.Order("date ASC")
		}).
		Where("tenant_id = ? AND is_current = ?", tenantID, true)

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}

	err := query.First(&academicYear).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAcademicYearNotFound
		}
		return nil, fmt.Errorf("failed to get current academic year: %w", err)
	}
	return &academicYear, nil
}

// List retrieves academic years with optional filtering.
func (s *Service) List(ctx context.Context, filter ListAcademicYearFilter) ([]models.AcademicYear, error) {
	query := s.db.WithContext(ctx).
		Model(&models.AcademicYear{}).
		Preload("Terms", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Preload("Holidays", func(db *gorm.DB) *gorm.DB {
			return db.Order("date ASC")
		}).
		Where("tenant_id = ?", filter.TenantID).
		Order("start_date DESC")

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.IsCurrent != nil {
		query = query.Where("is_current = ?", *filter.IsCurrent)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ?", search)
	}

	var academicYears []models.AcademicYear
	if err := query.Find(&academicYears).Error; err != nil {
		return nil, fmt.Errorf("failed to list academic years: %w", err)
	}

	return academicYears, nil
}

// Update updates an academic year.
func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, req UpdateAcademicYearRequest) (*models.AcademicYear, error) {
	academicYear, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil && *req.Name != "" {
		// Check if new name already exists
		var existing models.AcademicYear
		nameQuery := s.db.WithContext(ctx).
			Where("tenant_id = ? AND name = ? AND id != ?", tenantID, *req.Name, id)
		if academicYear.BranchID != nil {
			nameQuery = nameQuery.Where("branch_id = ?", *academicYear.BranchID)
		} else {
			nameQuery = nameQuery.Where("branch_id IS NULL")
		}
		if err := nameQuery.First(&existing).Error; err == nil {
			return nil, ErrAcademicYearNameExists
		}
		updates["name"] = *req.Name
	}

	startDate := academicYear.StartDate
	endDate := academicYear.EndDate

	if req.StartDate != nil {
		startDate = *req.StartDate
		updates["start_date"] = startDate
	}
	if req.EndDate != nil {
		endDate = *req.EndDate
		updates["end_date"] = endDate
	}

	// Validate dates if either was changed
	if req.StartDate != nil || req.EndDate != nil {
		if !endDate.After(startDate) {
			return nil, ErrAcademicYearInvalidDates
		}
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if req.UpdatedBy != nil {
		updates["updated_by"] = req.UpdatedBy
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(academicYear).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update academic year: %w", err)
		}
	}

	return s.GetByID(ctx, tenantID, id)
}

// SetCurrent sets an academic year as the current year.
func (s *Service) SetCurrent(ctx context.Context, tenantID, id uuid.UUID, updatedBy *uuid.UUID) (*models.AcademicYear, error) {
	academicYear, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if academicYear.IsCurrent {
		// Already current, return as-is
		return academicYear, nil
	}

	if !academicYear.IsActive {
		return nil, ErrCannotSetInactiveAsCurrent
	}

	// The database trigger will handle unsetting other current years
	updates := map[string]interface{}{
		"is_current": true,
		"updated_by": updatedBy,
	}

	if err := s.db.WithContext(ctx).Model(academicYear).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to set current academic year: %w", err)
	}

	return s.GetByID(ctx, tenantID, id)
}

// Delete deletes an academic year.
func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	academicYear, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Check for dependent records
	var termCount int64
	if err := s.db.WithContext(ctx).Model(&models.AcademicTerm{}).
		Where("academic_year_id = ?", id).Count(&termCount).Error; err != nil {
		return fmt.Errorf("failed to check term dependencies: %w", err)
	}

	var holidayCount int64
	if err := s.db.WithContext(ctx).Model(&models.Holiday{}).
		Where("academic_year_id = ?", id).Count(&holidayCount).Error; err != nil {
		return fmt.Errorf("failed to check holiday dependencies: %w", err)
	}

	// Delete in transaction to cascade
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete holidays first
		if err := tx.Where("academic_year_id = ?", id).Delete(&models.Holiday{}).Error; err != nil {
			return fmt.Errorf("failed to delete holidays: %w", err)
		}

		// Delete terms
		if err := tx.Where("academic_year_id = ?", id).Delete(&models.AcademicTerm{}).Error; err != nil {
			return fmt.Errorf("failed to delete terms: %w", err)
		}

		// Delete academic year
		if err := tx.Delete(academicYear).Error; err != nil {
			return fmt.Errorf("failed to delete academic year: %w", err)
		}

		return nil
	})

	return err
}

// Count returns the total number of academic years for a tenant.
func (s *Service) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.AcademicYear{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count academic years: %w", err)
	}
	return count, nil
}

// ============================================================================
// Academic Term Request Types
// ============================================================================

// CreateTermRequest represents a request to create an academic term.
type CreateTermRequest struct {
	TenantID       uuid.UUID
	AcademicYearID uuid.UUID
	Name           string
	StartDate      time.Time
	EndDate        time.Time
	Sequence       int
}

// UpdateTermRequest represents a request to update an academic term.
type UpdateTermRequest struct {
	Name      *string
	StartDate *time.Time
	EndDate   *time.Time
	Sequence  *int
}

// ============================================================================
// Academic Term Operations
// ============================================================================

// CreateTerm creates a new term for an academic year.
func (s *Service) CreateTerm(ctx context.Context, req CreateTermRequest) (*models.AcademicTerm, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.AcademicYearID == uuid.Nil {
		return nil, ErrAcademicYearNotFound
	}
	if req.Name == "" {
		return nil, ErrTermNameRequired
	}
	if req.StartDate.IsZero() {
		return nil, ErrTermStartDateRequired
	}
	if req.EndDate.IsZero() {
		return nil, ErrTermEndDateRequired
	}
	if !req.EndDate.After(req.StartDate) {
		return nil, ErrTermInvalidDates
	}

	// Verify academic year exists
	academicYear, err := s.GetByID(ctx, req.TenantID, req.AcademicYearID)
	if err != nil {
		return nil, err
	}

	// Validate term dates are within academic year
	if req.StartDate.Before(academicYear.StartDate) || req.EndDate.After(academicYear.EndDate) {
		return nil, ErrTermOutsideAcademicYear
	}

	// Check if term name already exists for this academic year
	var existing models.AcademicTerm
	err = s.db.WithContext(ctx).
		Where("academic_year_id = ? AND name = ?", req.AcademicYearID, req.Name).
		First(&existing).Error
	if err == nil {
		return nil, ErrTermNameExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing term: %w", err)
	}

	// Auto-assign sequence if not provided
	sequence := req.Sequence
	if sequence <= 0 {
		var maxSequence int
		s.db.WithContext(ctx).
			Model(&models.AcademicTerm{}).
			Where("academic_year_id = ?", req.AcademicYearID).
			Select("COALESCE(MAX(sequence), 0)").
			Scan(&maxSequence)
		sequence = maxSequence + 1
	}

	term := &models.AcademicTerm{
		TenantID:       req.TenantID,
		AcademicYearID: req.AcademicYearID,
		Name:           req.Name,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Sequence:       sequence,
	}

	if err := s.db.WithContext(ctx).Create(term).Error; err != nil {
		return nil, fmt.Errorf("failed to create term: %w", err)
	}

	return s.GetTermByID(ctx, req.TenantID, term.ID)
}

// GetTermByID retrieves a term by ID.
func (s *Service) GetTermByID(ctx context.Context, tenantID, id uuid.UUID) (*models.AcademicTerm, error) {
	var term models.AcademicTerm
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&term).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTermNotFound
		}
		return nil, fmt.Errorf("failed to get term: %w", err)
	}
	return &term, nil
}

// ListTerms retrieves terms for an academic year.
func (s *Service) ListTerms(ctx context.Context, tenantID, academicYearID uuid.UUID) ([]models.AcademicTerm, error) {
	var terms []models.AcademicTerm
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND academic_year_id = ?", tenantID, academicYearID).
		Order("sequence ASC").
		Find(&terms).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list terms: %w", err)
	}
	return terms, nil
}

// UpdateTerm updates an academic term.
func (s *Service) UpdateTerm(ctx context.Context, tenantID, academicYearID, termID uuid.UUID, req UpdateTermRequest) (*models.AcademicTerm, error) {
	term, err := s.GetTermByID(ctx, tenantID, termID)
	if err != nil {
		return nil, err
	}

	// Verify term belongs to the academic year
	if term.AcademicYearID != academicYearID {
		return nil, ErrTermNotFound
	}

	// Get academic year for date validation
	academicYear, err := s.GetByID(ctx, tenantID, academicYearID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil && *req.Name != "" {
		// Check if new name already exists
		var existing models.AcademicTerm
		err = s.db.WithContext(ctx).
			Where("academic_year_id = ? AND name = ? AND id != ?", academicYearID, *req.Name, termID).
			First(&existing).Error
		if err == nil {
			return nil, ErrTermNameExists
		}
		updates["name"] = *req.Name
	}

	startDate := term.StartDate
	endDate := term.EndDate

	if req.StartDate != nil {
		startDate = *req.StartDate
		updates["start_date"] = startDate
	}
	if req.EndDate != nil {
		endDate = *req.EndDate
		updates["end_date"] = endDate
	}

	// Validate dates if either was changed
	if req.StartDate != nil || req.EndDate != nil {
		if !endDate.After(startDate) {
			return nil, ErrTermInvalidDates
		}
		if startDate.Before(academicYear.StartDate) || endDate.After(academicYear.EndDate) {
			return nil, ErrTermOutsideAcademicYear
		}
	}

	if req.Sequence != nil && *req.Sequence > 0 {
		updates["sequence"] = *req.Sequence
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(term).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update term: %w", err)
		}
	}

	return s.GetTermByID(ctx, tenantID, termID)
}

// DeleteTerm deletes an academic term.
func (s *Service) DeleteTerm(ctx context.Context, tenantID, academicYearID, termID uuid.UUID) error {
	term, err := s.GetTermByID(ctx, tenantID, termID)
	if err != nil {
		return err
	}

	// Verify term belongs to the academic year
	if term.AcademicYearID != academicYearID {
		return ErrTermNotFound
	}

	if err := s.db.WithContext(ctx).Delete(term).Error; err != nil {
		return fmt.Errorf("failed to delete term: %w", err)
	}

	return nil
}

// ============================================================================
// Holiday Request Types
// ============================================================================

// CreateHolidayRequest represents a request to create a holiday.
type CreateHolidayRequest struct {
	TenantID       uuid.UUID
	AcademicYearID uuid.UUID
	BranchID       *uuid.UUID
	Name           string
	Date           time.Time
	Type           models.HolidayType
	IsOptional     bool
}

// UpdateHolidayRequest represents a request to update a holiday.
type UpdateHolidayRequest struct {
	Name       *string
	Date       *time.Time
	Type       *models.HolidayType
	IsOptional *bool
}

// ============================================================================
// Holiday Operations
// ============================================================================

// CreateHoliday creates a new holiday for an academic year.
func (s *Service) CreateHoliday(ctx context.Context, req CreateHolidayRequest) (*models.Holiday, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.AcademicYearID == uuid.Nil {
		return nil, ErrAcademicYearNotFound
	}
	if req.Name == "" {
		return nil, ErrHolidayNameRequired
	}
	if req.Date.IsZero() {
		return nil, ErrHolidayDateRequired
	}

	// Set default type if not provided
	holidayType := req.Type
	if holidayType == "" {
		holidayType = models.HolidayTypePublic
	}
	if !holidayType.IsValid() {
		return nil, ErrHolidayInvalidType
	}

	// Verify academic year exists
	academicYear, err := s.GetByID(ctx, req.TenantID, req.AcademicYearID)
	if err != nil {
		return nil, err
	}

	// Validate holiday date is within academic year
	if req.Date.Before(academicYear.StartDate) || req.Date.After(academicYear.EndDate) {
		return nil, ErrHolidayOutsideAcademicYear
	}

	holiday := &models.Holiday{
		TenantID:       req.TenantID,
		AcademicYearID: req.AcademicYearID,
		BranchID:       req.BranchID,
		Name:           req.Name,
		Date:           req.Date,
		Type:           holidayType,
		IsOptional:     req.IsOptional,
	}

	if err := s.db.WithContext(ctx).Create(holiday).Error; err != nil {
		return nil, fmt.Errorf("failed to create holiday: %w", err)
	}

	return s.GetHolidayByID(ctx, req.TenantID, holiday.ID)
}

// GetHolidayByID retrieves a holiday by ID.
func (s *Service) GetHolidayByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Holiday, error) {
	var holiday models.Holiday
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&holiday).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHolidayNotFound
		}
		return nil, fmt.Errorf("failed to get holiday: %w", err)
	}
	return &holiday, nil
}

// ListHolidays retrieves holidays for an academic year.
func (s *Service) ListHolidays(ctx context.Context, tenantID, academicYearID uuid.UUID) ([]models.Holiday, error) {
	var holidays []models.Holiday
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND academic_year_id = ?", tenantID, academicYearID).
		Order("date ASC").
		Find(&holidays).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list holidays: %w", err)
	}
	return holidays, nil
}

// UpdateHoliday updates a holiday.
func (s *Service) UpdateHoliday(ctx context.Context, tenantID, academicYearID, holidayID uuid.UUID, req UpdateHolidayRequest) (*models.Holiday, error) {
	holiday, err := s.GetHolidayByID(ctx, tenantID, holidayID)
	if err != nil {
		return nil, err
	}

	// Verify holiday belongs to the academic year
	if holiday.AcademicYearID != academicYearID {
		return nil, ErrHolidayNotFound
	}

	// Get academic year for date validation
	academicYear, err := s.GetByID(ctx, tenantID, academicYearID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil && *req.Name != "" {
		updates["name"] = *req.Name
	}

	if req.Date != nil {
		// Validate date is within academic year
		if req.Date.Before(academicYear.StartDate) || req.Date.After(academicYear.EndDate) {
			return nil, ErrHolidayOutsideAcademicYear
		}
		updates["date"] = *req.Date
	}

	if req.Type != nil {
		if !req.Type.IsValid() {
			return nil, ErrHolidayInvalidType
		}
		updates["type"] = *req.Type
	}

	if req.IsOptional != nil {
		updates["is_optional"] = *req.IsOptional
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(holiday).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update holiday: %w", err)
		}
	}

	return s.GetHolidayByID(ctx, tenantID, holidayID)
}

// DeleteHoliday deletes a holiday.
func (s *Service) DeleteHoliday(ctx context.Context, tenantID, academicYearID, holidayID uuid.UUID) error {
	holiday, err := s.GetHolidayByID(ctx, tenantID, holidayID)
	if err != nil {
		return err
	}

	// Verify holiday belongs to the academic year
	if holiday.AcademicYearID != academicYearID {
		return ErrHolidayNotFound
	}

	if err := s.db.WithContext(ctx).Delete(holiday).Error; err != nil {
		return fmt.Errorf("failed to delete holiday: %w", err)
	}

	return nil
}

// IsHoliday checks if a given date is a holiday.
func (s *Service) IsHoliday(ctx context.Context, tenantID uuid.UUID, date time.Time, branchID *uuid.UUID) (bool, *models.Holiday, error) {
	var holiday models.Holiday
	query := s.db.WithContext(ctx).
		Where("tenant_id = ? AND date = ?", tenantID, date)

	if branchID != nil {
		// Check for branch-specific or tenant-wide holidays
		query = query.Where("branch_id = ? OR branch_id IS NULL", *branchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}

	err := query.First(&holiday).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, fmt.Errorf("failed to check holiday: %w", err)
	}

	return true, &holiday, nil
}
