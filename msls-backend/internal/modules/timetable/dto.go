// Package timetable provides timetable management functionality.
package timetable

import (
	"msls-backend/internal/pkg/database/models"

	"github.com/google/uuid"
)

// ========================================
// Shift DTOs
// ========================================

// ShiftResponse represents a shift in API responses.
type ShiftResponse struct {
	ID           uuid.UUID `json:"id"`
	BranchID     uuid.UUID `json:"branchId"`
	BranchName   string    `json:"branchName,omitempty"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	StartTime    string    `json:"startTime"`
	EndTime      string    `json:"endTime"`
	Description  string    `json:"description,omitempty"`
	DisplayOrder int       `json:"displayOrder"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

// ShiftListResponse represents the response for listing shifts.
type ShiftListResponse struct {
	Shifts []ShiftResponse `json:"shifts"`
	Total  int64           `json:"total"`
}

// CreateShiftRequest represents the request body for creating a shift.
type CreateShiftRequest struct {
	BranchID     uuid.UUID `json:"branchId" binding:"required"`
	Name         string    `json:"name" binding:"required,max=50"`
	Code         string    `json:"code" binding:"required,max=20"`
	StartTime    string    `json:"startTime" binding:"required"`
	EndTime      string    `json:"endTime" binding:"required"`
	Description  string    `json:"description"`
	DisplayOrder int       `json:"displayOrder"`
}

// UpdateShiftRequest represents the request body for updating a shift.
type UpdateShiftRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=50"`
	Code         *string `json:"code" binding:"omitempty,max=20"`
	StartTime    *string `json:"startTime"`
	EndTime      *string `json:"endTime"`
	Description  *string `json:"description"`
	DisplayOrder *int    `json:"displayOrder"`
	IsActive     *bool   `json:"isActive"`
}

// ShiftFilter represents filters for listing shifts.
type ShiftFilter struct {
	TenantID uuid.UUID
	BranchID *uuid.UUID
	IsActive *bool
}

// ShiftToResponse converts a Shift model to ShiftResponse.
func ShiftToResponse(s *models.Shift) ShiftResponse {
	resp := ShiftResponse{
		ID:           s.ID,
		BranchID:     s.BranchID,
		Name:         s.Name,
		Code:         s.Code,
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
		Description:  s.Description,
		DisplayOrder: s.DisplayOrder,
		IsActive:     s.IsActive,
		CreatedAt:    s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if s.Branch != nil {
		resp.BranchName = s.Branch.Name
	}

	return resp
}

// ========================================
// Day Pattern DTOs
// ========================================

// DayPatternResponse represents a day pattern in API responses.
type DayPatternResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description,omitempty"`
	TotalPeriods int       `json:"totalPeriods"`
	DisplayOrder int       `json:"displayOrder"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

// DayPatternListResponse represents the response for listing day patterns.
type DayPatternListResponse struct {
	DayPatterns []DayPatternResponse `json:"dayPatterns"`
	Total       int64                `json:"total"`
}

// CreateDayPatternRequest represents the request body for creating a day pattern.
type CreateDayPatternRequest struct {
	Name         string `json:"name" binding:"required,max=50"`
	Code         string `json:"code" binding:"required,max=20"`
	Description  string `json:"description"`
	TotalPeriods int    `json:"totalPeriods"`
	DisplayOrder int    `json:"displayOrder"`
}

// UpdateDayPatternRequest represents the request body for updating a day pattern.
type UpdateDayPatternRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=50"`
	Code         *string `json:"code" binding:"omitempty,max=20"`
	Description  *string `json:"description"`
	TotalPeriods *int    `json:"totalPeriods"`
	DisplayOrder *int    `json:"displayOrder"`
	IsActive     *bool   `json:"isActive"`
}

// DayPatternFilter represents filters for listing day patterns.
type DayPatternFilter struct {
	TenantID uuid.UUID
	IsActive *bool
}

// DayPatternToResponse converts a DayPattern model to DayPatternResponse.
func DayPatternToResponse(dp *models.DayPattern) DayPatternResponse {
	return DayPatternResponse{
		ID:           dp.ID,
		Name:         dp.Name,
		Code:         dp.Code,
		Description:  dp.Description,
		TotalPeriods: dp.TotalPeriods,
		DisplayOrder: dp.DisplayOrder,
		IsActive:     dp.IsActive,
		CreatedAt:    dp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    dp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ========================================
// Day Pattern Assignment DTOs
// ========================================

// DayPatternAssignmentResponse represents a day pattern assignment in API responses.
type DayPatternAssignmentResponse struct {
	ID             uuid.UUID `json:"id"`
	BranchID       uuid.UUID `json:"branchId"`
	DayOfWeek      int       `json:"dayOfWeek"`
	DayName        string    `json:"dayName"`
	DayPatternID   *string   `json:"dayPatternId,omitempty"`
	DayPatternName string    `json:"dayPatternName,omitempty"`
	DayPatternCode string    `json:"dayPatternCode,omitempty"`
	IsWorkingDay   bool      `json:"isWorkingDay"`
}

// DayPatternAssignmentListResponse represents the response for listing day pattern assignments.
type DayPatternAssignmentListResponse struct {
	Assignments []DayPatternAssignmentResponse `json:"assignments"`
}

// UpdateDayPatternAssignmentRequest represents the request to update a day assignment.
type UpdateDayPatternAssignmentRequest struct {
	DayPatternID *uuid.UUID `json:"dayPatternId"`
	IsWorkingDay *bool      `json:"isWorkingDay"`
}

// DayPatternAssignmentToResponse converts a DayPatternAssignment model to response.
func DayPatternAssignmentToResponse(a *models.DayPatternAssignment) DayPatternAssignmentResponse {
	resp := DayPatternAssignmentResponse{
		ID:           a.ID,
		BranchID:     a.BranchID,
		DayOfWeek:    a.DayOfWeek,
		DayName:      a.GetDayName(),
		IsWorkingDay: a.IsWorkingDay,
	}

	if a.DayPatternID != nil {
		idStr := a.DayPatternID.String()
		resp.DayPatternID = &idStr
		if a.DayPattern != nil {
			resp.DayPatternName = a.DayPattern.Name
			resp.DayPatternCode = a.DayPattern.Code
		}
	}

	return resp
}

// ========================================
// Period Slot DTOs
// ========================================

// PeriodSlotResponse represents a period slot in API responses.
type PeriodSlotResponse struct {
	ID              uuid.UUID `json:"id"`
	BranchID        uuid.UUID `json:"branchId"`
	Name            string    `json:"name"`
	PeriodNumber    *int      `json:"periodNumber,omitempty"`
	SlotType        string    `json:"slotType"`
	StartTime       string    `json:"startTime"`
	EndTime         string    `json:"endTime"`
	DurationMinutes int       `json:"durationMinutes"`
	DayPatternID    *string   `json:"dayPatternId,omitempty"`
	DayPatternName  string    `json:"dayPatternName,omitempty"`
	ShiftID         *string   `json:"shiftId,omitempty"`
	ShiftName       string    `json:"shiftName,omitempty"`
	DisplayOrder    int       `json:"displayOrder"`
	IsActive        bool      `json:"isActive"`
	IsTeaching      bool      `json:"isTeaching"`
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string    `json:"updatedAt"`
}

// PeriodSlotListResponse represents the response for listing period slots.
type PeriodSlotListResponse struct {
	PeriodSlots []PeriodSlotResponse `json:"periodSlots"`
	Total       int64                `json:"total"`
}

// CreatePeriodSlotRequest represents the request body for creating a period slot.
type CreatePeriodSlotRequest struct {
	BranchID        uuid.UUID  `json:"branchId" binding:"required"`
	Name            string     `json:"name" binding:"required,max=50"`
	PeriodNumber    *int       `json:"periodNumber"`
	SlotType        string     `json:"slotType" binding:"required,oneof=regular short assembly break lunch activity zero_period"`
	StartTime       string     `json:"startTime" binding:"required"`
	EndTime         string     `json:"endTime" binding:"required"`
	DurationMinutes int        `json:"durationMinutes" binding:"required,min=1"`
	DayPatternID    *uuid.UUID `json:"dayPatternId"`
	ShiftID         *uuid.UUID `json:"shiftId"`
	DisplayOrder    int        `json:"displayOrder"`
}

// UpdatePeriodSlotRequest represents the request body for updating a period slot.
type UpdatePeriodSlotRequest struct {
	Name            *string    `json:"name" binding:"omitempty,max=50"`
	PeriodNumber    *int       `json:"periodNumber"`
	SlotType        *string    `json:"slotType" binding:"omitempty,oneof=regular short assembly break lunch activity zero_period"`
	StartTime       *string    `json:"startTime"`
	EndTime         *string    `json:"endTime"`
	DurationMinutes *int       `json:"durationMinutes" binding:"omitempty,min=1"`
	DayPatternID    *uuid.UUID `json:"dayPatternId"`
	ShiftID         *uuid.UUID `json:"shiftId"`
	DisplayOrder    *int       `json:"displayOrder"`
	IsActive        *bool      `json:"isActive"`
}

// PeriodSlotFilter represents filters for listing period slots.
type PeriodSlotFilter struct {
	TenantID     uuid.UUID
	BranchID     *uuid.UUID
	DayPatternID *uuid.UUID
	ShiftID      *uuid.UUID
	SlotType     *string
	IsActive     *bool
}

// PeriodSlotToResponse converts a PeriodSlot model to PeriodSlotResponse.
func PeriodSlotToResponse(ps *models.PeriodSlot) PeriodSlotResponse {
	resp := PeriodSlotResponse{
		ID:              ps.ID,
		BranchID:        ps.BranchID,
		Name:            ps.Name,
		PeriodNumber:    ps.PeriodNumber,
		SlotType:        string(ps.SlotType),
		StartTime:       ps.StartTime,
		EndTime:         ps.EndTime,
		DurationMinutes: ps.DurationMinutes,
		DisplayOrder:    ps.DisplayOrder,
		IsActive:        ps.IsActive,
		IsTeaching:      ps.IsTeachingPeriod(),
		CreatedAt:       ps.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       ps.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if ps.DayPatternID != nil {
		idStr := ps.DayPatternID.String()
		resp.DayPatternID = &idStr
		if ps.DayPattern != nil {
			resp.DayPatternName = ps.DayPattern.Name
		}
	}

	if ps.ShiftID != nil {
		idStr := ps.ShiftID.String()
		resp.ShiftID = &idStr
		if ps.Shift != nil {
			resp.ShiftName = ps.Shift.Name
		}
	}

	return resp
}
