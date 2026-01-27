// Package admission provides admission management services.
package admission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// SessionService handles admission session operations.
type SessionService struct {
	db *gorm.DB
}

// NewSessionService creates a new SessionService instance.
func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

// CreateSessionRequest represents a request to create an admission session.
type CreateSessionRequest struct {
	TenantID          uuid.UUID
	BranchID          *uuid.UUID
	AcademicYearID    *uuid.UUID
	Name              string
	Description       string
	StartDate         time.Time
	EndDate           time.Time
	ApplicationFee    decimal.Decimal
	RequiredDocuments []string
	Settings          models.SessionSettings
	CreatedBy         *uuid.UUID
}

// UpdateSessionRequest represents a request to update an admission session.
type UpdateSessionRequest struct {
	Name              *string
	Description       *string
	StartDate         *time.Time
	EndDate           *time.Time
	ApplicationFee    *decimal.Decimal
	RequiredDocuments *[]string
	Settings          *models.SessionSettings
	UpdatedBy         *uuid.UUID
}

// ListSessionFilter contains filters for listing admission sessions.
type ListSessionFilter struct {
	TenantID       uuid.UUID
	BranchID       *uuid.UUID
	AcademicYearID *uuid.UUID
	Status         *models.AdmissionSessionStatus
	Search         string
	IncludeSeats   bool
}

// SessionStats contains statistics for an admission session.
type SessionStats struct {
	TotalApplications int
	ApprovedCount     int
	PendingCount      int
	RejectedCount     int
	TotalSeats        int
	FilledSeats       int
	AvailableSeats    int
}

// Create creates a new admission session.
func (s *SessionService) Create(ctx context.Context, req CreateSessionRequest) (*models.AdmissionSession, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.Name == "" {
		return nil, ErrSessionNameRequired
	}
	if req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDateRange
	}

	// Check if session name already exists for tenant and academic year
	var existing models.AdmissionSession
	query := s.db.WithContext(ctx).
		Where("tenant_id = ? AND name = ?", req.TenantID, req.Name)
	if req.AcademicYearID != nil {
		query = query.Where("academic_year_id = ?", *req.AcademicYearID)
	} else {
		query = query.Where("academic_year_id IS NULL")
	}
	err := query.First(&existing).Error
	if err == nil {
		return nil, ErrSessionNameExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing session: %w", err)
	}

	// Set default required documents if not provided
	requiredDocs := models.RequiredDocuments(req.RequiredDocuments)
	if requiredDocs == nil {
		requiredDocs = models.RequiredDocuments{}
	}

	session := &models.AdmissionSession{
		TenantID:          req.TenantID,
		BranchID:          req.BranchID,
		AcademicYearID:    req.AcademicYearID,
		Name:              req.Name,
		Description:       req.Description,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		Status:            models.SessionStatusUpcoming,
		ApplicationFee:    req.ApplicationFee,
		RequiredDocuments: requiredDocs,
		Settings:          req.Settings,
		CreatedBy:         req.CreatedBy,
		UpdatedBy:         req.CreatedBy,
	}

	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create admission session: %w", err)
	}

	return s.GetByID(ctx, req.TenantID, session.ID, true)
}

// GetByID retrieves an admission session by ID.
func (s *SessionService) GetByID(ctx context.Context, tenantID, id uuid.UUID, includeSeats bool) (*models.AdmissionSession, error) {
	var session models.AdmissionSession
	query := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id)

	if includeSeats {
		query = query.Preload("Seats")
	}

	err := query.First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get admission session: %w", err)
	}
	return &session, nil
}

// List retrieves admission sessions with optional filtering.
func (s *SessionService) List(ctx context.Context, filter ListSessionFilter) ([]models.AdmissionSession, error) {
	query := s.db.WithContext(ctx).
		Model(&models.AdmissionSession{}).
		Where("tenant_id = ?", filter.TenantID).
		Order("start_date DESC, name")

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.AcademicYearID != nil {
		query = query.Where("academic_year_id = ?", *filter.AcademicYearID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	if filter.IncludeSeats {
		query = query.Preload("Seats")
	}

	var sessions []models.AdmissionSession
	if err := query.Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to list admission sessions: %w", err)
	}

	return sessions, nil
}

// Update updates an admission session.
func (s *SessionService) Update(ctx context.Context, tenantID, id uuid.UUID, req UpdateSessionRequest) (*models.AdmissionSession, error) {
	session, err := s.GetByID(ctx, tenantID, id, false)
	if err != nil {
		return nil, err
	}

	// Cannot modify closed sessions
	if session.Status == models.SessionStatusClosed {
		return nil, ErrCannotModifyClosedSession
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		if *req.Name == "" {
			return nil, ErrSessionNameRequired
		}
		// Check if name already exists for another session
		var existing models.AdmissionSession
		query := s.db.WithContext(ctx).
			Where("tenant_id = ? AND name = ? AND id != ?", tenantID, *req.Name, id)
		if session.AcademicYearID != nil {
			query = query.Where("academic_year_id = ?", *session.AcademicYearID)
		} else {
			query = query.Where("academic_year_id IS NULL")
		}
		err := query.First(&existing).Error
		if err == nil {
			return nil, ErrSessionNameExists
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check existing session: %w", err)
		}
		updates["name"] = *req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	startDate := session.StartDate
	endDate := session.EndDate

	if req.StartDate != nil {
		startDate = *req.StartDate
		updates["start_date"] = *req.StartDate
	}

	if req.EndDate != nil {
		endDate = *req.EndDate
		updates["end_date"] = *req.EndDate
	}

	// Validate date range
	if endDate.Before(startDate) {
		return nil, ErrInvalidDateRange
	}

	if req.ApplicationFee != nil {
		updates["application_fee"] = *req.ApplicationFee
	}

	if req.RequiredDocuments != nil {
		updates["required_documents"] = models.RequiredDocuments(*req.RequiredDocuments)
	}

	if req.Settings != nil {
		updates["settings"] = *req.Settings
	}

	if req.UpdatedBy != nil {
		updates["updated_by"] = req.UpdatedBy
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(session).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update admission session: %w", err)
		}
	}

	return s.GetByID(ctx, tenantID, id, true)
}

// ChangeStatus changes the status of an admission session.
func (s *SessionService) ChangeStatus(ctx context.Context, tenantID, id uuid.UUID, newStatus models.AdmissionSessionStatus, updatedBy *uuid.UUID) (*models.AdmissionSession, error) {
	if !newStatus.IsValid() {
		return nil, ErrInvalidStatus
	}

	session, err := s.GetByID(ctx, tenantID, id, false)
	if err != nil {
		return nil, err
	}

	// Validate status transitions
	if !isValidStatusTransition(session.Status, newStatus) {
		return nil, ErrInvalidStatusTransition
	}

	updates := map[string]interface{}{
		"status":     newStatus,
		"updated_by": updatedBy,
	}

	if err := s.db.WithContext(ctx).Model(session).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update session status: %w", err)
	}

	return s.GetByID(ctx, tenantID, id, true)
}

// Delete deletes an admission session.
func (s *SessionService) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	session, err := s.GetByID(ctx, tenantID, id, false)
	if err != nil {
		return err
	}

	// Cannot delete open sessions
	if session.Status == models.SessionStatusOpen {
		return ErrCannotDeleteOpenSession
	}

	// TODO: Check for associated applications when application table is available

	// Delete in a transaction (seats will be cascade deleted)
	if err := s.db.WithContext(ctx).Delete(session).Error; err != nil {
		return fmt.Errorf("failed to delete admission session: %w", err)
	}

	return nil
}

// Count returns the total number of admission sessions for a tenant.
func (s *SessionService) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.AdmissionSession{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count admission sessions: %w", err)
	}
	return count, nil
}

// GetStats retrieves statistics for an admission session.
func (s *SessionService) GetStats(ctx context.Context, tenantID, id uuid.UUID) (*SessionStats, error) {
	session, err := s.GetByID(ctx, tenantID, id, true)
	if err != nil {
		return nil, err
	}

	stats := &SessionStats{}

	// Calculate seat statistics from seats
	for _, seat := range session.Seats {
		stats.TotalSeats += seat.TotalSeats
		stats.FilledSeats += seat.FilledSeats
	}
	stats.AvailableSeats = stats.TotalSeats - stats.FilledSeats

	// TODO: Add application statistics when applications table is available

	return stats, nil
}

// ExtendDeadline extends the end date of an admission session.
func (s *SessionService) ExtendDeadline(ctx context.Context, tenantID, id uuid.UUID, newEndDate time.Time, updatedBy *uuid.UUID) (*models.AdmissionSession, error) {
	session, err := s.GetByID(ctx, tenantID, id, false)
	if err != nil {
		return nil, err
	}

	// Cannot modify closed sessions
	if session.Status == models.SessionStatusClosed {
		return nil, ErrCannotModifyClosedSession
	}

	// New end date must be after start date
	if newEndDate.Before(session.StartDate) {
		return nil, ErrInvalidDateRange
	}

	updates := map[string]interface{}{
		"end_date":   newEndDate,
		"updated_by": updatedBy,
	}

	if err := s.db.WithContext(ctx).Model(session).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to extend deadline: %w", err)
	}

	return s.GetByID(ctx, tenantID, id, true)
}

// isValidStatusTransition checks if a status transition is allowed.
func isValidStatusTransition(from, to models.AdmissionSessionStatus) bool {
	// Valid transitions:
	// upcoming -> open
	// upcoming -> closed
	// open -> closed
	// closed -> open (reopen)
	switch from {
	case models.SessionStatusUpcoming:
		return to == models.SessionStatusOpen || to == models.SessionStatusClosed
	case models.SessionStatusOpen:
		return to == models.SessionStatusClosed
	case models.SessionStatusClosed:
		return to == models.SessionStatusOpen // Allow reopening
	}
	return false
}

// =============================================================================
// Seat Management
// =============================================================================

// CreateSeatRequest represents a request to create a seat configuration.
type CreateSeatRequest struct {
	TenantID      uuid.UUID
	SessionID     uuid.UUID
	ClassName     string
	TotalSeats    int
	WaitlistLimit int
	ReservedSeats models.ReservedSeats
}

// UpdateSeatRequest represents a request to update a seat configuration.
type UpdateSeatRequest struct {
	TotalSeats    *int
	WaitlistLimit *int
	ReservedSeats *models.ReservedSeats
}

// CreateSeat creates a new seat configuration for a session.
func (s *SessionService) CreateSeat(ctx context.Context, req CreateSeatRequest) (*models.AdmissionSeat, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.SessionID == uuid.Nil {
		return nil, ErrSessionIDRequired
	}
	if req.ClassName == "" {
		return nil, ErrClassNameRequired
	}
	if req.TotalSeats < 0 {
		return nil, ErrInvalidTotalSeats
	}

	// Verify session exists and belongs to tenant
	session, err := s.GetByID(ctx, req.TenantID, req.SessionID, false)
	if err != nil {
		return nil, err
	}

	// Cannot modify closed sessions
	if session.Status == models.SessionStatusClosed {
		return nil, ErrCannotModifyClosedSession
	}

	// Check if class already exists for this session
	var existing models.AdmissionSeat
	err = s.db.WithContext(ctx).
		Where("session_id = ? AND class_name = ?", req.SessionID, req.ClassName).
		First(&existing).Error
	if err == nil {
		return nil, ErrClassAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing seat configuration: %w", err)
	}

	// Set defaults
	waitlistLimit := req.WaitlistLimit
	if waitlistLimit == 0 {
		waitlistLimit = 10
	}

	reservedSeats := req.ReservedSeats
	if reservedSeats == nil {
		reservedSeats = models.ReservedSeats{}
	}

	seat := &models.AdmissionSeat{
		TenantID:      req.TenantID,
		SessionID:     req.SessionID,
		ClassName:     req.ClassName,
		TotalSeats:    req.TotalSeats,
		FilledSeats:   0,
		WaitlistLimit: waitlistLimit,
		ReservedSeats: reservedSeats,
	}

	if err := s.db.WithContext(ctx).Create(seat).Error; err != nil {
		return nil, fmt.Errorf("failed to create seat configuration: %w", err)
	}

	return seat, nil
}

// GetSeatByID retrieves a seat configuration by ID.
func (s *SessionService) GetSeatByID(ctx context.Context, tenantID, seatID uuid.UUID) (*models.AdmissionSeat, error) {
	var seat models.AdmissionSeat
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, seatID).
		First(&seat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSeatNotFound
		}
		return nil, fmt.Errorf("failed to get seat configuration: %w", err)
	}
	return &seat, nil
}

// ListSeats retrieves all seat configurations for a session.
func (s *SessionService) ListSeats(ctx context.Context, tenantID, sessionID uuid.UUID) ([]models.AdmissionSeat, error) {
	// Verify session exists
	_, err := s.GetByID(ctx, tenantID, sessionID, false)
	if err != nil {
		return nil, err
	}

	var seats []models.AdmissionSeat
	err = s.db.WithContext(ctx).
		Where("tenant_id = ? AND session_id = ?", tenantID, sessionID).
		Order("class_name").
		Find(&seats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list seat configurations: %w", err)
	}

	return seats, nil
}

// UpdateSeat updates a seat configuration.
func (s *SessionService) UpdateSeat(ctx context.Context, tenantID, seatID uuid.UUID, req UpdateSeatRequest) (*models.AdmissionSeat, error) {
	seat, err := s.GetSeatByID(ctx, tenantID, seatID)
	if err != nil {
		return nil, err
	}

	// Check if session is closed
	session, err := s.GetByID(ctx, tenantID, seat.SessionID, false)
	if err != nil {
		return nil, err
	}
	if session.Status == models.SessionStatusClosed {
		return nil, ErrCannotModifyClosedSession
	}

	updates := make(map[string]interface{})

	if req.TotalSeats != nil {
		if *req.TotalSeats < 0 {
			return nil, ErrInvalidTotalSeats
		}
		if seat.FilledSeats > *req.TotalSeats {
			return nil, ErrFilledExceedsTotal
		}
		updates["total_seats"] = *req.TotalSeats
	}

	if req.WaitlistLimit != nil {
		updates["waitlist_limit"] = *req.WaitlistLimit
	}

	if req.ReservedSeats != nil {
		updates["reserved_seats"] = *req.ReservedSeats
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(seat).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update seat configuration: %w", err)
		}
	}

	return s.GetSeatByID(ctx, tenantID, seatID)
}

// DeleteSeat deletes a seat configuration.
func (s *SessionService) DeleteSeat(ctx context.Context, tenantID, seatID uuid.UUID) error {
	seat, err := s.GetSeatByID(ctx, tenantID, seatID)
	if err != nil {
		return err
	}

	// Check if session is closed
	session, err := s.GetByID(ctx, tenantID, seat.SessionID, false)
	if err != nil {
		return err
	}
	if session.Status == models.SessionStatusClosed {
		return ErrCannotModifyClosedSession
	}

	// TODO: Check for associated applications when available

	if err := s.db.WithContext(ctx).Delete(seat).Error; err != nil {
		return fmt.Errorf("failed to delete seat configuration: %w", err)
	}

	return nil
}

// IncrementFilledSeats increments the filled seats count for a seat configuration.
func (s *SessionService) IncrementFilledSeats(ctx context.Context, tenantID, seatID uuid.UUID, count int) error {
	seat, err := s.GetSeatByID(ctx, tenantID, seatID)
	if err != nil {
		return err
	}

	newFilled := seat.FilledSeats + count
	if newFilled > seat.TotalSeats {
		return ErrFilledExceedsTotal
	}
	if newFilled < 0 {
		newFilled = 0
	}

	if err := s.db.WithContext(ctx).Model(seat).Update("filled_seats", newFilled).Error; err != nil {
		return fmt.Errorf("failed to update filled seats: %w", err)
	}

	return nil
}
