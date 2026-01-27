// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"time"

	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// ============================================================================
// Request DTOs
// ============================================================================

// CreateSessionRequest represents the request body for creating an admission session.
type CreateSessionRequest struct {
	BranchID          *string         `json:"branchId"`
	AcademicYearID    *string         `json:"academicYearId"`
	Name              string          `json:"name" binding:"required,max=200"`
	Description       string          `json:"description" binding:"max=1000"`
	StartDate         string          `json:"startDate" binding:"required"`
	EndDate           string          `json:"endDate" binding:"required"`
	ApplicationFee    decimal.Decimal `json:"applicationFee"`
	RequiredDocuments []string        `json:"requiredDocuments"`
	Settings          *SessionSettings `json:"settings"`
}

// UpdateSessionRequest represents the request body for updating an admission session.
type UpdateSessionRequest struct {
	Name              *string          `json:"name" binding:"omitempty,max=200"`
	Description       *string          `json:"description" binding:"omitempty,max=1000"`
	StartDate         *string          `json:"startDate"`
	EndDate           *string          `json:"endDate"`
	ApplicationFee    *decimal.Decimal `json:"applicationFee"`
	RequiredDocuments *[]string        `json:"requiredDocuments"`
	Settings          *SessionSettings `json:"settings"`
}

// ChangeStatusRequest represents the request body for changing session status.
type ChangeStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=upcoming open closed"`
}

// ExtendDeadlineRequest represents the request body for extending session deadline.
type ExtendDeadlineRequest struct {
	EndDate string `json:"endDate" binding:"required"`
}

// CreateSeatRequest represents the request body for creating a seat configuration.
type CreateSeatRequest struct {
	ClassName     string         `json:"className" binding:"required,max=100"`
	TotalSeats    int            `json:"totalSeats" binding:"min=0"`
	WaitlistLimit int            `json:"waitlistLimit" binding:"min=0"`
	ReservedSeats map[string]int `json:"reservedSeats"`
}

// UpdateSeatRequest represents the request body for updating a seat configuration.
type UpdateSeatRequest struct {
	TotalSeats    *int            `json:"totalSeats" binding:"omitempty,min=0"`
	WaitlistLimit *int            `json:"waitlistLimit" binding:"omitempty,min=0"`
	ReservedSeats *map[string]int `json:"reservedSeats"`
}

// SessionSettings represents session settings in API requests/responses.
type SessionSettings struct {
	AllowOnlineApplication bool   `json:"allowOnlineApplication"`
	NotifyOnApplication    bool   `json:"notifyOnApplication"`
	AutoConfirmPayment     bool   `json:"autoConfirmPayment"`
	MaxApplicationsPerDay  int    `json:"maxApplicationsPerDay,omitempty"`
	Instructions           string `json:"instructions,omitempty"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// SessionResponse represents an admission session in API responses.
type SessionResponse struct {
	ID                string           `json:"id"`
	BranchID          *string          `json:"branchId,omitempty"`
	AcademicYearID    *string          `json:"academicYearId,omitempty"`
	Name              string           `json:"name"`
	Description       string           `json:"description,omitempty"`
	StartDate         string           `json:"startDate"`
	EndDate           string           `json:"endDate"`
	Status            string           `json:"status"`
	ApplicationFee    decimal.Decimal  `json:"applicationFee"`
	RequiredDocuments []string         `json:"requiredDocuments"`
	Settings          SessionSettings  `json:"settings"`
	CreatedAt         string           `json:"createdAt"`
	UpdatedAt         string           `json:"updatedAt"`
	Seats             []SeatResponse   `json:"seats,omitempty"`
	Stats             *SessionStatsDTO `json:"stats,omitempty"`
}

// SeatResponse represents a seat configuration in API responses.
type SeatResponse struct {
	ID             string         `json:"id"`
	SessionID      string         `json:"sessionId"`
	ClassName      string         `json:"className"`
	TotalSeats     int            `json:"totalSeats"`
	FilledSeats    int            `json:"filledSeats"`
	AvailableSeats int            `json:"availableSeats"`
	WaitlistLimit  int            `json:"waitlistLimit"`
	ReservedSeats  map[string]int `json:"reservedSeats"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      string         `json:"updatedAt"`
}

// SessionStatsDTO represents session statistics in API responses.
type SessionStatsDTO struct {
	TotalApplications int `json:"totalApplications"`
	ApprovedCount     int `json:"approvedCount"`
	PendingCount      int `json:"pendingCount"`
	RejectedCount     int `json:"rejectedCount"`
	TotalSeats        int `json:"totalSeats"`
	FilledSeats       int `json:"filledSeats"`
	AvailableSeats    int `json:"availableSeats"`
}

// SessionListResponse represents a list of admission sessions.
type SessionListResponse struct {
	Sessions []SessionResponse `json:"sessions"`
	Total    int64             `json:"total"`
}

// SeatListResponse represents a list of seat configurations.
type SeatListResponse struct {
	Seats []SeatResponse `json:"seats"`
	Total int            `json:"total"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// sessionToResponse converts a Session model to a SessionResponse.
func sessionToResponse(session *models.AdmissionSession) SessionResponse {
	resp := SessionResponse{
		ID:             session.ID.String(),
		Name:           session.Name,
		Description:    session.Description,
		StartDate:      session.StartDate.Format("2006-01-02"),
		EndDate:        session.EndDate.Format("2006-01-02"),
		Status:         string(session.Status),
		ApplicationFee: session.ApplicationFee,
		RequiredDocuments: func() []string {
			if session.RequiredDocuments == nil {
				return []string{}
			}
			return session.RequiredDocuments
		}(),
		Settings: SessionSettings{
			AllowOnlineApplication: session.Settings.AllowOnlineApplication,
			NotifyOnApplication:    session.Settings.NotifyOnApplication,
			AutoConfirmPayment:     session.Settings.AutoConfirmPayment,
			MaxApplicationsPerDay:  session.Settings.MaxApplicationsPerDay,
			Instructions:           session.Settings.Instructions,
		},
		CreatedAt: session.CreatedAt.Format(time.RFC3339),
		UpdatedAt: session.UpdatedAt.Format(time.RFC3339),
	}

	if session.BranchID != nil {
		branchIDStr := session.BranchID.String()
		resp.BranchID = &branchIDStr
	}

	if session.AcademicYearID != nil {
		academicYearIDStr := session.AcademicYearID.String()
		resp.AcademicYearID = &academicYearIDStr
	}

	// Convert seats if present
	if len(session.Seats) > 0 {
		resp.Seats = make([]SeatResponse, len(session.Seats))
		for i, seat := range session.Seats {
			resp.Seats[i] = seatToResponse(&seat)
		}
	}

	return resp
}

// sessionsToResponses converts a slice of Session models to SessionResponses.
func sessionsToResponses(sessions []models.AdmissionSession) []SessionResponse {
	responses := make([]SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = sessionToResponse(&session)
	}
	return responses
}

// seatToResponse converts a Seat model to a SeatResponse.
func seatToResponse(seat *models.AdmissionSeat) SeatResponse {
	reservedSeats := map[string]int{}
	if seat.ReservedSeats != nil {
		reservedSeats = seat.ReservedSeats
	}

	return SeatResponse{
		ID:             seat.ID.String(),
		SessionID:      seat.SessionID.String(),
		ClassName:      seat.ClassName,
		TotalSeats:     seat.TotalSeats,
		FilledSeats:    seat.FilledSeats,
		AvailableSeats: seat.AvailableSeats(),
		WaitlistLimit:  seat.WaitlistLimit,
		ReservedSeats:  reservedSeats,
		CreatedAt:      seat.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      seat.UpdatedAt.Format(time.RFC3339),
	}
}

// seatsToResponses converts a slice of Seat models to SeatResponses.
func seatsToResponses(seats []models.AdmissionSeat) []SeatResponse {
	responses := make([]SeatResponse, len(seats))
	for i, seat := range seats {
		responses[i] = seatToResponse(&seat)
	}
	return responses
}

// settingsToModel converts a SessionSettings DTO to a model.
func settingsToModel(s *SessionSettings) models.SessionSettings {
	if s == nil {
		return models.SessionSettings{}
	}
	return models.SessionSettings{
		AllowOnlineApplication: s.AllowOnlineApplication,
		NotifyOnApplication:    s.NotifyOnApplication,
		AutoConfirmPayment:     s.AutoConfirmPayment,
		MaxApplicationsPerDay:  s.MaxApplicationsPerDay,
		Instructions:           s.Instructions,
	}
}
