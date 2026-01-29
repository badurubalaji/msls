// Package timetable provides timetable management functionality.
package timetable

import (
	"context"

	"msls-backend/internal/pkg/database/models"

	"github.com/google/uuid"
)

// Service handles business logic for timetable entities.
type Service struct {
	repo *Repository
}

// NewService creates a new timetable service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ========================================
// Shift Service Methods
// ========================================

// ListShifts returns all shifts for a tenant with filters.
func (s *Service) ListShifts(ctx context.Context, filter ShiftFilter) ([]models.Shift, int64, error) {
	return s.repo.ListShifts(ctx, filter)
}

// GetShiftByID returns a shift by ID.
func (s *Service) GetShiftByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Shift, error) {
	return s.repo.GetShiftByID(ctx, tenantID, id)
}

// CreateShift creates a new shift.
func (s *Service) CreateShift(ctx context.Context, tenantID uuid.UUID, req CreateShiftRequest, userID uuid.UUID) (*models.Shift, error) {
	// Check if code already exists
	existing, err := s.repo.GetShiftByCode(ctx, tenantID, req.BranchID, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrShiftCodeExists
	}

	shift := &models.Shift{
		TenantID:     tenantID,
		BranchID:     req.BranchID,
		Name:         req.Name,
		Code:         req.Code,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Description:  req.Description,
		DisplayOrder: req.DisplayOrder,
		IsActive:     true,
		CreatedBy:    &userID,
	}

	if err := s.repo.CreateShift(ctx, shift); err != nil {
		return nil, err
	}

	return s.repo.GetShiftByID(ctx, tenantID, shift.ID)
}

// UpdateShift updates an existing shift.
func (s *Service) UpdateShift(ctx context.Context, tenantID, id uuid.UUID, req UpdateShiftRequest) (*models.Shift, error) {
	shift, err := s.repo.GetShiftByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check code uniqueness if changing
	if req.Code != nil && *req.Code != shift.Code {
		existing, err := s.repo.GetShiftByCode(ctx, tenantID, shift.BranchID, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrShiftCodeExists
		}
		shift.Code = *req.Code
	}

	if req.Name != nil {
		shift.Name = *req.Name
	}
	if req.StartTime != nil {
		shift.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		shift.EndTime = *req.EndTime
	}
	if req.Description != nil {
		shift.Description = *req.Description
	}
	if req.DisplayOrder != nil {
		shift.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		shift.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateShift(ctx, shift); err != nil {
		return nil, err
	}

	return s.repo.GetShiftByID(ctx, tenantID, id)
}

// DeleteShift deletes a shift.
func (s *Service) DeleteShift(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := s.repo.GetShiftByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteShift(ctx, tenantID, id)
}

// ========================================
// Day Pattern Service Methods
// ========================================

// ListDayPatterns returns all day patterns for a tenant with filters.
func (s *Service) ListDayPatterns(ctx context.Context, filter DayPatternFilter) ([]models.DayPattern, int64, error) {
	return s.repo.ListDayPatterns(ctx, filter)
}

// GetDayPatternByID returns a day pattern by ID.
func (s *Service) GetDayPatternByID(ctx context.Context, tenantID, id uuid.UUID) (*models.DayPattern, error) {
	return s.repo.GetDayPatternByID(ctx, tenantID, id)
}

// CreateDayPattern creates a new day pattern.
func (s *Service) CreateDayPattern(ctx context.Context, tenantID uuid.UUID, req CreateDayPatternRequest, userID uuid.UUID) (*models.DayPattern, error) {
	// Check if code already exists
	existing, err := s.repo.GetDayPatternByCode(ctx, tenantID, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDayPatternCodeExists
	}

	totalPeriods := 8
	if req.TotalPeriods > 0 {
		totalPeriods = req.TotalPeriods
	}

	pattern := &models.DayPattern{
		TenantID:     tenantID,
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		TotalPeriods: totalPeriods,
		DisplayOrder: req.DisplayOrder,
		IsActive:     true,
		CreatedBy:    &userID,
	}

	if err := s.repo.CreateDayPattern(ctx, pattern); err != nil {
		return nil, err
	}

	return pattern, nil
}

// UpdateDayPattern updates an existing day pattern.
func (s *Service) UpdateDayPattern(ctx context.Context, tenantID, id uuid.UUID, req UpdateDayPatternRequest) (*models.DayPattern, error) {
	pattern, err := s.repo.GetDayPatternByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check code uniqueness if changing
	if req.Code != nil && *req.Code != pattern.Code {
		existing, err := s.repo.GetDayPatternByCode(ctx, tenantID, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrDayPatternCodeExists
		}
		pattern.Code = *req.Code
	}

	if req.Name != nil {
		pattern.Name = *req.Name
	}
	if req.Description != nil {
		pattern.Description = *req.Description
	}
	if req.TotalPeriods != nil {
		pattern.TotalPeriods = *req.TotalPeriods
	}
	if req.DisplayOrder != nil {
		pattern.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		pattern.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateDayPattern(ctx, pattern); err != nil {
		return nil, err
	}

	return pattern, nil
}

// DeleteDayPattern deletes a day pattern.
func (s *Service) DeleteDayPattern(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := s.repo.GetDayPatternByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	inUse, err := s.repo.IsDayPatternInUse(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return ErrDayPatternInUse
	}

	return s.repo.DeleteDayPattern(ctx, tenantID, id)
}

// ========================================
// Day Pattern Assignment Service Methods
// ========================================

// ListDayPatternAssignments returns all day pattern assignments for a branch.
func (s *Service) ListDayPatternAssignments(ctx context.Context, tenantID, branchID uuid.UUID) ([]models.DayPatternAssignment, error) {
	return s.repo.ListDayPatternAssignments(ctx, tenantID, branchID)
}

// UpdateDayPatternAssignment updates a day pattern assignment.
func (s *Service) UpdateDayPatternAssignment(ctx context.Context, tenantID, branchID uuid.UUID, dayOfWeek int, req UpdateDayPatternAssignmentRequest) (*models.DayPatternAssignment, error) {
	assignment, err := s.repo.GetDayPatternAssignment(ctx, tenantID, branchID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	if assignment == nil {
		assignment = &models.DayPatternAssignment{
			TenantID:     tenantID,
			BranchID:     branchID,
			DayOfWeek:    dayOfWeek,
			IsWorkingDay: true,
		}
	}

	if req.DayPatternID != nil {
		assignment.DayPatternID = req.DayPatternID
	}
	if req.IsWorkingDay != nil {
		assignment.IsWorkingDay = *req.IsWorkingDay
	}

	if err := s.repo.UpsertDayPatternAssignment(ctx, assignment); err != nil {
		return nil, err
	}

	return s.repo.GetDayPatternAssignment(ctx, tenantID, branchID, dayOfWeek)
}

// ========================================
// Period Slot Service Methods
// ========================================

// ListPeriodSlots returns all period slots for a tenant with filters.
func (s *Service) ListPeriodSlots(ctx context.Context, filter PeriodSlotFilter) ([]models.PeriodSlot, int64, error) {
	return s.repo.ListPeriodSlots(ctx, filter)
}

// GetPeriodSlotByID returns a period slot by ID.
func (s *Service) GetPeriodSlotByID(ctx context.Context, tenantID, id uuid.UUID) (*models.PeriodSlot, error) {
	return s.repo.GetPeriodSlotByID(ctx, tenantID, id)
}

// CreatePeriodSlot creates a new period slot.
func (s *Service) CreatePeriodSlot(ctx context.Context, tenantID uuid.UUID, req CreatePeriodSlotRequest, userID uuid.UUID) (*models.PeriodSlot, error) {
	slot := &models.PeriodSlot{
		TenantID:        tenantID,
		BranchID:        req.BranchID,
		Name:            req.Name,
		PeriodNumber:    req.PeriodNumber,
		SlotType:        models.PeriodSlotType(req.SlotType),
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationMinutes: req.DurationMinutes,
		DayPatternID:    req.DayPatternID,
		ShiftID:         req.ShiftID,
		DisplayOrder:    req.DisplayOrder,
		IsActive:        true,
		CreatedBy:       &userID,
	}

	if err := s.repo.CreatePeriodSlot(ctx, slot); err != nil {
		return nil, err
	}

	return s.repo.GetPeriodSlotByID(ctx, tenantID, slot.ID)
}

// UpdatePeriodSlot updates an existing period slot.
func (s *Service) UpdatePeriodSlot(ctx context.Context, tenantID, id uuid.UUID, req UpdatePeriodSlotRequest) (*models.PeriodSlot, error) {
	slot, err := s.repo.GetPeriodSlotByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		slot.Name = *req.Name
	}
	if req.PeriodNumber != nil {
		slot.PeriodNumber = req.PeriodNumber
	}
	if req.SlotType != nil {
		slot.SlotType = models.PeriodSlotType(*req.SlotType)
	}
	if req.StartTime != nil {
		slot.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		slot.EndTime = *req.EndTime
	}
	if req.DurationMinutes != nil {
		slot.DurationMinutes = *req.DurationMinutes
	}
	if req.DayPatternID != nil {
		slot.DayPatternID = req.DayPatternID
	}
	if req.ShiftID != nil {
		slot.ShiftID = req.ShiftID
	}
	if req.DisplayOrder != nil {
		slot.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		slot.IsActive = *req.IsActive
	}

	if err := s.repo.UpdatePeriodSlot(ctx, slot); err != nil {
		return nil, err
	}

	return s.repo.GetPeriodSlotByID(ctx, tenantID, id)
}

// DeletePeriodSlot deletes a period slot.
func (s *Service) DeletePeriodSlot(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := s.repo.GetPeriodSlotByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	return s.repo.DeletePeriodSlot(ctx, tenantID, id)
}
