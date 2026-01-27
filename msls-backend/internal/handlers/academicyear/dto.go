// Package academicyear provides HTTP handlers for academic year management endpoints.
package academicyear

import (
	"time"

	"msls-backend/internal/pkg/database/models"
)

// ============================================================================
// Request DTOs
// ============================================================================

// CreateAcademicYearRequest represents the request body for creating an academic year.
type CreateAcademicYearRequest struct {
	Name      string  `json:"name" binding:"required,max=50"`
	StartDate string  `json:"startDate" binding:"required"`
	EndDate   string  `json:"endDate" binding:"required"`
	BranchID  *string `json:"branchId" binding:"omitempty,uuid"`
	IsCurrent bool    `json:"isCurrent"`
}

// UpdateAcademicYearRequest represents the request body for updating an academic year.
type UpdateAcademicYearRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=50"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
	IsActive  *bool   `json:"isActive"`
}

// CreateTermRequest represents the request body for creating a term.
type CreateTermRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	Sequence  int    `json:"sequence"`
}

// UpdateTermRequest represents the request body for updating a term.
type UpdateTermRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=100"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
	Sequence  *int    `json:"sequence"`
}

// CreateHolidayRequest represents the request body for creating a holiday.
type CreateHolidayRequest struct {
	Name       string  `json:"name" binding:"required,max=200"`
	Date       string  `json:"date" binding:"required"`
	Type       string  `json:"type" binding:"omitempty,oneof=public religious national school other"`
	BranchID   *string `json:"branchId" binding:"omitempty,uuid"`
	IsOptional bool    `json:"isOptional"`
}

// UpdateHolidayRequest represents the request body for updating a holiday.
type UpdateHolidayRequest struct {
	Name       *string `json:"name" binding:"omitempty,max=200"`
	Date       *string `json:"date"`
	Type       *string `json:"type" binding:"omitempty,oneof=public religious national school other"`
	IsOptional *bool   `json:"isOptional"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// AcademicYearResponse represents an academic year in API responses.
type AcademicYearResponse struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	StartDate string             `json:"startDate"`
	EndDate   string             `json:"endDate"`
	BranchID  *string            `json:"branchId,omitempty"`
	IsCurrent bool               `json:"isCurrent"`
	IsActive  bool               `json:"isActive"`
	Terms     []TermResponse     `json:"terms,omitempty"`
	Holidays  []HolidayResponse  `json:"holidays,omitempty"`
	CreatedAt string             `json:"createdAt"`
	UpdatedAt string             `json:"updatedAt"`
}

// AcademicYearListResponse represents a list of academic years.
type AcademicYearListResponse struct {
	AcademicYears []AcademicYearResponse `json:"academicYears"`
	Total         int64                  `json:"total"`
}

// TermResponse represents a term in API responses.
type TermResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Sequence  int    `json:"sequence"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// TermListResponse represents a list of terms.
type TermListResponse struct {
	Terms []TermResponse `json:"terms"`
	Total int            `json:"total"`
}

// HolidayResponse represents a holiday in API responses.
type HolidayResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Date       string  `json:"date"`
	Type       string  `json:"type"`
	BranchID   *string `json:"branchId,omitempty"`
	IsOptional bool    `json:"isOptional"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// HolidayListResponse represents a list of holidays.
type HolidayListResponse struct {
	Holidays []HolidayResponse `json:"holidays"`
	Total    int               `json:"total"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// academicYearToResponse converts an AcademicYear model to an AcademicYearResponse.
func academicYearToResponse(ay *models.AcademicYear) AcademicYearResponse {
	response := AcademicYearResponse{
		ID:        ay.ID.String(),
		Name:      ay.Name,
		StartDate: ay.StartDate.Format("2006-01-02"),
		EndDate:   ay.EndDate.Format("2006-01-02"),
		IsCurrent: ay.IsCurrent,
		IsActive:  ay.IsActive,
		Terms:     make([]TermResponse, 0),
		Holidays:  make([]HolidayResponse, 0),
		CreatedAt: ay.CreatedAt.Format(time.RFC3339),
		UpdatedAt: ay.UpdatedAt.Format(time.RFC3339),
	}

	if ay.BranchID != nil {
		branchID := ay.BranchID.String()
		response.BranchID = &branchID
	}

	// Convert terms
	for _, term := range ay.Terms {
		response.Terms = append(response.Terms, termToResponse(&term))
	}

	// Convert holidays
	for _, holiday := range ay.Holidays {
		response.Holidays = append(response.Holidays, holidayToResponse(&holiday))
	}

	return response
}

// academicYearsToResponses converts a slice of AcademicYear models to AcademicYearResponses.
func academicYearsToResponses(academicYears []models.AcademicYear) []AcademicYearResponse {
	responses := make([]AcademicYearResponse, len(academicYears))
	for i, ay := range academicYears {
		responses[i] = academicYearToResponse(&ay)
	}
	return responses
}

// termToResponse converts an AcademicTerm model to a TermResponse.
func termToResponse(term *models.AcademicTerm) TermResponse {
	return TermResponse{
		ID:        term.ID.String(),
		Name:      term.Name,
		StartDate: term.StartDate.Format("2006-01-02"),
		EndDate:   term.EndDate.Format("2006-01-02"),
		Sequence:  term.Sequence,
		CreatedAt: term.CreatedAt.Format(time.RFC3339),
		UpdatedAt: term.UpdatedAt.Format(time.RFC3339),
	}
}

// termsToResponses converts a slice of AcademicTerm models to TermResponses.
func termsToResponses(terms []models.AcademicTerm) []TermResponse {
	responses := make([]TermResponse, len(terms))
	for i, term := range terms {
		responses[i] = termToResponse(&term)
	}
	return responses
}

// holidayToResponse converts a Holiday model to a HolidayResponse.
func holidayToResponse(holiday *models.Holiday) HolidayResponse {
	response := HolidayResponse{
		ID:         holiday.ID.String(),
		Name:       holiday.Name,
		Date:       holiday.Date.Format("2006-01-02"),
		Type:       string(holiday.Type),
		IsOptional: holiday.IsOptional,
		CreatedAt:  holiday.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  holiday.UpdatedAt.Format(time.RFC3339),
	}

	if holiday.BranchID != nil {
		branchID := holiday.BranchID.String()
		response.BranchID = &branchID
	}

	return response
}

// holidaysToResponses converts a slice of Holiday models to HolidayResponses.
func holidaysToResponses(holidays []models.Holiday) []HolidayResponse {
	responses := make([]HolidayResponse, len(holidays))
	for i, holiday := range holidays {
		responses[i] = holidayToResponse(&holiday)
	}
	return responses
}
