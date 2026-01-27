// Package attendance provides staff attendance management functionality.
package attendance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// StaffService interface for staff operations.
type StaffService interface {
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Staff, error)
}

// Service handles attendance business logic.
type Service struct {
	repo         *Repository
	staffService StaffService
	db           *gorm.DB
}

// NewService creates a new attendance service.
func NewService(db *gorm.DB, staffService StaffService) *Service {
	return &Service{
		repo:         NewRepository(db),
		staffService: staffService,
		db:           db,
	}
}

// CheckIn marks check-in for a staff member.
func (s *Service) CheckIn(ctx context.Context, dto CheckInDTO) (*models.StaffAttendance, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.StaffID == uuid.Nil {
		return nil, ErrStaffIDRequired
	}

	// Verify staff exists
	staff, err := s.staffService.GetByID(ctx, dto.TenantID, dto.StaffID)
	if err != nil {
		return nil, ErrStaffNotFound
	}

	today := time.Now().Truncate(24 * time.Hour)

	// Check if already checked in today
	existing, err := s.repo.GetAttendanceByStaffAndDate(ctx, dto.TenantID, dto.StaffID, today)
	if err == nil && existing != nil {
		if existing.CheckInTime != nil {
			return nil, ErrAlreadyCheckedIn
		}
	}

	// Get branch settings for late calculation
	settings, _ := s.repo.GetSettings(ctx, dto.TenantID, staff.BranchID)

	now := time.Now()
	isLate := false
	lateMinutes := 0

	if settings != nil {
		// Calculate if late
		workStart := time.Date(now.Year(), now.Month(), now.Day(),
			settings.WorkStartTime.Hour(), settings.WorkStartTime.Minute(), 0, 0, now.Location())
		threshold := workStart.Add(time.Duration(settings.LateThresholdMinutes) * time.Minute)

		if now.After(threshold) {
			isLate = true
			lateMinutes = int(now.Sub(workStart).Minutes())
		}
	}

	status := models.AttendanceStatusPresent
	var halfDayType models.HalfDayType
	if dto.HalfDayType != "" {
		if !HalfDayType(dto.HalfDayType).IsValid() {
			return nil, ErrInvalidHalfDayType
		}
		status = models.AttendanceStatusHalfDay
		halfDayType = models.HalfDayType(dto.HalfDayType)
	}

	if existing != nil {
		// Update existing record
		existing.CheckInTime = &now
		existing.Status = status
		existing.HalfDayType = halfDayType
		existing.IsLate = isLate
		existing.LateMinutes = lateMinutes
		existing.Remarks = dto.Remarks
		existing.MarkedBy = dto.MarkedBy
		existing.MarkedAt = now
		existing.UpdatedAt = now

		if err := s.repo.UpdateAttendance(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Create new attendance record
	attendance := &models.StaffAttendance{
		TenantID:       dto.TenantID,
		StaffID:        dto.StaffID,
		AttendanceDate: today,
		Status:         status,
		CheckInTime:    &now,
		IsLate:         isLate,
		LateMinutes:    lateMinutes,
		HalfDayType:    halfDayType,
		Remarks:        dto.Remarks,
		MarkedBy:       dto.MarkedBy,
		MarkedAt:       now,
	}

	if err := s.repo.CreateAttendance(ctx, attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// CheckOut marks check-out for a staff member.
func (s *Service) CheckOut(ctx context.Context, dto CheckOutDTO) (*models.StaffAttendance, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.StaffID == uuid.Nil {
		return nil, ErrStaffIDRequired
	}

	today := time.Now().Truncate(24 * time.Hour)

	// Get today's attendance
	attendance, err := s.repo.GetAttendanceByStaffAndDate(ctx, dto.TenantID, dto.StaffID, today)
	if err != nil {
		if errors.Is(err, ErrAttendanceNotFound) {
			return nil, ErrNotCheckedIn
		}
		return nil, err
	}

	if attendance.CheckInTime == nil {
		return nil, ErrNotCheckedIn
	}

	if attendance.CheckOutTime != nil {
		return nil, ErrAlreadyCheckedOut
	}

	now := time.Now()
	attendance.CheckOutTime = &now
	if dto.Remarks != "" {
		if attendance.Remarks != "" {
			attendance.Remarks = attendance.Remarks + "; " + dto.Remarks
		} else {
			attendance.Remarks = dto.Remarks
		}
	}
	attendance.UpdatedAt = now

	if err := s.repo.UpdateAttendance(ctx, attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// MarkAttendance marks attendance for a staff member (used by HR).
func (s *Service) MarkAttendance(ctx context.Context, dto MarkAttendanceDTO) (*models.StaffAttendance, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.StaffID == uuid.Nil {
		return nil, ErrStaffIDRequired
	}
	if dto.AttendanceDate.IsZero() {
		return nil, ErrDateRequired
	}
	if dto.AttendanceDate.After(time.Now()) {
		return nil, ErrFutureDate
	}
	if !AttendanceStatus(dto.Status).IsValid() {
		return nil, ErrInvalidStatus
	}
	if dto.Status == StatusHalfDay && dto.HalfDayType != "" {
		if !HalfDayType(dto.HalfDayType).IsValid() {
			return nil, ErrInvalidHalfDayType
		}
	}

	// Verify staff exists
	staff, err := s.staffService.GetByID(ctx, dto.TenantID, dto.StaffID)
	if err != nil {
		return nil, ErrStaffNotFound
	}

	date := dto.AttendanceDate.Truncate(24 * time.Hour)

	// Check if attendance already exists
	existing, err := s.repo.GetAttendanceByStaffAndDate(ctx, dto.TenantID, dto.StaffID, date)
	if err != nil && !errors.Is(err, ErrAttendanceNotFound) {
		return nil, err
	}

	// Calculate late status if present
	isLate := false
	lateMinutes := 0
	if dto.Status == StatusPresent && dto.CheckInTime != nil {
		settings, _ := s.repo.GetSettings(ctx, dto.TenantID, staff.BranchID)
		if settings != nil {
			workStart := time.Date(dto.CheckInTime.Year(), dto.CheckInTime.Month(), dto.CheckInTime.Day(),
				settings.WorkStartTime.Hour(), settings.WorkStartTime.Minute(), 0, 0, dto.CheckInTime.Location())
			threshold := workStart.Add(time.Duration(settings.LateThresholdMinutes) * time.Minute)

			if dto.CheckInTime.After(threshold) {
				isLate = true
				lateMinutes = int(dto.CheckInTime.Sub(workStart).Minutes())
			}
		}
	}

	now := time.Now()

	if existing != nil {
		// Update existing record
		existing.Status = models.AttendanceStatus(dto.Status)
		existing.CheckInTime = dto.CheckInTime
		existing.CheckOutTime = dto.CheckOutTime
		existing.IsLate = isLate
		existing.LateMinutes = lateMinutes
		existing.HalfDayType = models.HalfDayType(dto.HalfDayType)
		existing.Remarks = dto.Remarks
		existing.MarkedBy = dto.MarkedBy
		existing.MarkedAt = now
		existing.UpdatedAt = now

		if err := s.repo.UpdateAttendance(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Create new attendance record
	attendance := &models.StaffAttendance{
		TenantID:       dto.TenantID,
		StaffID:        dto.StaffID,
		AttendanceDate: date,
		Status:         models.AttendanceStatus(dto.Status),
		CheckInTime:    dto.CheckInTime,
		CheckOutTime:   dto.CheckOutTime,
		IsLate:         isLate,
		LateMinutes:    lateMinutes,
		HalfDayType:    models.HalfDayType(dto.HalfDayType),
		Remarks:        dto.Remarks,
		MarkedBy:       dto.MarkedBy,
		MarkedAt:       now,
	}

	if err := s.repo.CreateAttendance(ctx, attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// GetTodayAttendance retrieves today's attendance for a staff member.
func (s *Service) GetTodayAttendance(ctx context.Context, tenantID, staffID uuid.UUID) (*models.StaffAttendance, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.repo.GetAttendanceByStaffAndDate(ctx, tenantID, staffID, today)
}

// GetAttendanceByID retrieves an attendance record by ID.
func (s *Service) GetAttendanceByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StaffAttendance, error) {
	return s.repo.GetAttendanceByID(ctx, tenantID, id)
}

// ListAttendance retrieves attendance records with filtering.
func (s *Service) ListAttendance(ctx context.Context, filter ListFilter) ([]models.StaffAttendance, string, int64, error) {
	return s.repo.ListAttendance(ctx, filter)
}

// GetMonthlySummary retrieves monthly attendance summary for a staff member.
func (s *Service) GetMonthlySummary(ctx context.Context, tenantID, staffID uuid.UUID, year, month int) (*AttendanceSummaryResponse, error) {
	attendanceList, err := s.repo.GetMonthlyAttendance(ctx, tenantID, staffID, year, month)
	if err != nil {
		return nil, err
	}

	summary := &AttendanceSummaryResponse{
		Month:     time.Month(month).String(),
		Year:      year,
		TotalDays: len(attendanceList),
	}

	for _, a := range attendanceList {
		switch a.Status {
		case models.AttendanceStatusPresent:
			summary.PresentDays++
		case models.AttendanceStatusAbsent:
			summary.AbsentDays++
		case models.AttendanceStatusHalfDay:
			summary.HalfDays++
		case models.AttendanceStatusOnLeave:
			summary.LeaveDays++
		case models.AttendanceStatusHoliday:
			summary.HolidayDays++
		}
		if a.IsLate {
			summary.LateDays++
			summary.TotalLateMinutes += a.LateMinutes
		}
	}

	return summary, nil
}

// RequestRegularization submits a regularization request.
func (s *Service) RequestRegularization(ctx context.Context, dto RegularizationRequestDTO) (*models.StaffAttendanceRegularization, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.StaffID == uuid.Nil {
		return nil, ErrStaffIDRequired
	}
	if dto.RequestDate.IsZero() {
		return nil, ErrDateRequired
	}
	if dto.RequestDate.After(time.Now()) {
		return nil, ErrFutureDate
	}
	if !AttendanceStatus(dto.RequestedStatus).IsValid() {
		return nil, ErrInvalidStatus
	}
	if dto.Reason == "" {
		return nil, ErrReasonRequired
	}

	// Verify staff exists
	_, err := s.staffService.GetByID(ctx, dto.TenantID, dto.StaffID)
	if err != nil {
		return nil, ErrStaffNotFound
	}

	// Check if there's already a pending regularization for this date
	existing, err := s.repo.GetPendingRegularizationByStaffAndDate(ctx, dto.TenantID, dto.StaffID, dto.RequestDate)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCannotRegularizePendingRequest
	}

	// Get the attendance record if it exists
	date := dto.RequestDate.Truncate(24 * time.Hour)
	attendance, _ := s.repo.GetAttendanceByStaffAndDate(ctx, dto.TenantID, dto.StaffID, date)

	regularization := &models.StaffAttendanceRegularization{
		TenantID:              dto.TenantID,
		StaffID:               dto.StaffID,
		RequestDate:           date,
		RequestedStatus:       models.AttendanceStatus(dto.RequestedStatus),
		Reason:                dto.Reason,
		SupportingDocumentURL: dto.SupportingDocumentURL,
		Status:                models.RegularizationStatusPending,
	}

	if attendance != nil {
		regularization.AttendanceID = &attendance.ID
	}

	if err := s.repo.CreateRegularization(ctx, regularization); err != nil {
		return nil, err
	}

	return regularization, nil
}

// ApproveRegularization approves a regularization request.
func (s *Service) ApproveRegularization(ctx context.Context, dto RegularizationReviewDTO) (*models.StaffAttendanceRegularization, error) {
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.RegularizationID == uuid.Nil {
		return nil, ErrRegularizationNotFound
	}

	regularization, err := s.repo.GetRegularizationByID(ctx, dto.TenantID, dto.RegularizationID)
	if err != nil {
		return nil, err
	}

	if regularization.Status != models.RegularizationStatusPending {
		return nil, ErrRegularizationAlreadyProcessed
	}

	now := time.Now()

	// Use transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update regularization status
		regularization.Status = models.RegularizationStatusApproved
		regularization.ReviewedBy = dto.ReviewedBy
		regularization.ReviewedAt = &now
		regularization.UpdatedAt = now

		if err := s.repo.UpdateRegularization(ctx, regularization); err != nil {
			return err
		}

		// Update or create attendance record
		date := regularization.RequestDate
		attendance, err := s.repo.GetAttendanceByStaffAndDate(ctx, dto.TenantID, regularization.StaffID, date)
		if err != nil && !errors.Is(err, ErrAttendanceNotFound) {
			return err
		}

		if attendance != nil {
			// Update existing attendance
			attendance.Status = regularization.RequestedStatus
			attendance.Remarks = fmt.Sprintf("Regularized: %s", regularization.Reason)
			attendance.UpdatedAt = now
			return s.repo.UpdateAttendance(ctx, attendance)
		}

		// Create new attendance record
		newAttendance := &models.StaffAttendance{
			TenantID:       dto.TenantID,
			StaffID:        regularization.StaffID,
			AttendanceDate: date,
			Status:         regularization.RequestedStatus,
			Remarks:        fmt.Sprintf("Regularized: %s", regularization.Reason),
			MarkedBy:       dto.ReviewedBy,
			MarkedAt:       now,
		}
		return s.repo.CreateAttendance(ctx, newAttendance)
	})

	if err != nil {
		return nil, err
	}

	return regularization, nil
}

// RejectRegularization rejects a regularization request.
func (s *Service) RejectRegularization(ctx context.Context, dto RegularizationReviewDTO) (*models.StaffAttendanceRegularization, error) {
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.RegularizationID == uuid.Nil {
		return nil, ErrRegularizationNotFound
	}
	if dto.RejectionReason == "" {
		return nil, ErrReasonRequired
	}

	regularization, err := s.repo.GetRegularizationByID(ctx, dto.TenantID, dto.RegularizationID)
	if err != nil {
		return nil, err
	}

	if regularization.Status != models.RegularizationStatusPending {
		return nil, ErrRegularizationAlreadyProcessed
	}

	now := time.Now()
	regularization.Status = models.RegularizationStatusRejected
	regularization.ReviewedBy = dto.ReviewedBy
	regularization.ReviewedAt = &now
	regularization.RejectionReason = dto.RejectionReason
	regularization.UpdatedAt = now

	if err := s.repo.UpdateRegularization(ctx, regularization); err != nil {
		return nil, err
	}

	return regularization, nil
}

// ListRegularizations retrieves regularization requests with filtering.
func (s *Service) ListRegularizations(ctx context.Context, filter RegularizationFilter) ([]models.StaffAttendanceRegularization, string, int64, error) {
	return s.repo.ListRegularizations(ctx, filter)
}

// GetRegularizationByID retrieves a regularization request by ID.
func (s *Service) GetRegularizationByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StaffAttendanceRegularization, error) {
	return s.repo.GetRegularizationByID(ctx, tenantID, id)
}

// GetSettings retrieves attendance settings for a branch.
func (s *Service) GetSettings(ctx context.Context, tenantID, branchID uuid.UUID) (*models.StaffAttendanceSettings, error) {
	return s.repo.GetSettings(ctx, tenantID, branchID)
}

// UpdateSettings updates attendance settings for a branch.
func (s *Service) UpdateSettings(ctx context.Context, dto SettingsDTO) (*models.StaffAttendanceSettings, error) {
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.BranchID == uuid.Nil {
		return nil, ErrBranchIDRequired
	}

	settings := &models.StaffAttendanceSettings{
		TenantID:                      dto.TenantID,
		BranchID:                      dto.BranchID,
		WorkStartTime:                 dto.WorkStartTime,
		WorkEndTime:                   dto.WorkEndTime,
		LateThresholdMinutes:          dto.LateThresholdMinutes,
		HalfDayThresholdHours:         dto.HalfDayThresholdHours,
		AllowSelfCheckout:             dto.AllowSelfCheckout,
		RequireRegularizationApproval: dto.RequireRegularizationApproval,
	}

	if err := s.repo.CreateOrUpdateSettings(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}
