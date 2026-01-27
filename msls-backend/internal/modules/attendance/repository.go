// Package attendance provides staff attendance management functionality.
package attendance

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for attendance.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new attendance repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateAttendance creates a new attendance record.
func (r *Repository) CreateAttendance(ctx context.Context, attendance *models.StaffAttendance) error {
	if err := r.db.WithContext(ctx).Create(attendance).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "uniq_staff_attendance") {
			return ErrDuplicateAttendance
		}
		return fmt.Errorf("create attendance: %w", err)
	}
	return nil
}

// UpdateAttendance updates an attendance record.
func (r *Repository) UpdateAttendance(ctx context.Context, attendance *models.StaffAttendance) error {
	result := r.db.WithContext(ctx).Save(attendance)
	if result.Error != nil {
		return fmt.Errorf("update attendance: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrAttendanceNotFound
	}
	return nil
}

// GetAttendanceByID retrieves an attendance record by ID.
func (r *Repository) GetAttendanceByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StaffAttendance, error) {
	var attendance models.StaffAttendance
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&attendance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, fmt.Errorf("get attendance by id: %w", err)
	}
	return &attendance, nil
}

// GetAttendanceByStaffAndDate retrieves an attendance record for a staff member on a specific date.
func (r *Repository) GetAttendanceByStaffAndDate(ctx context.Context, tenantID, staffID uuid.UUID, date time.Time) (*models.StaffAttendance, error) {
	var attendance models.StaffAttendance
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Where("tenant_id = ? AND staff_id = ? AND attendance_date = ?", tenantID, staffID, date.Format("2006-01-02")).
		First(&attendance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, fmt.Errorf("get attendance by staff and date: %w", err)
	}
	return &attendance, nil
}

// ListAttendance retrieves attendance records with optional filtering.
func (r *Repository) ListAttendance(ctx context.Context, filter ListFilter) ([]models.StaffAttendance, string, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.StaffAttendance{}).
		Preload("Staff").
		Where("staff_attendance.tenant_id = ?", filter.TenantID)

	// Apply filters
	if filter.StaffID != nil {
		query = query.Where("staff_attendance.staff_id = ?", *filter.StaffID)
	}

	if filter.BranchID != nil {
		query = query.Joins("JOIN staff ON staff.id = staff_attendance.staff_id").
			Where("staff.branch_id = ?", *filter.BranchID)
	}

	if filter.DepartmentID != nil {
		query = query.Joins("JOIN staff s ON s.id = staff_attendance.staff_id").
			Where("s.department_id = ?", *filter.DepartmentID)
	}

	if filter.Status != nil {
		query = query.Where("staff_attendance.status = ?", *filter.Status)
	}

	if filter.DateFrom != nil {
		query = query.Where("staff_attendance.attendance_date >= ?", filter.DateFrom.Format("2006-01-02"))
	}

	if filter.DateTo != nil {
		query = query.Where("staff_attendance.attendance_date <= ?", filter.DateTo.Format("2006-01-02"))
	}

	// Get total count
	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, "", 0, fmt.Errorf("count attendance: %w", err)
	}

	// Apply cursor pagination
	if filter.Cursor != "" {
		cursorID, err := uuid.Parse(filter.Cursor)
		if err == nil {
			query = query.Where("staff_attendance.id > ?", cursorID)
		}
	}

	// Apply sorting
	sortBy := "staff_attendance.attendance_date DESC, staff_attendance.created_at DESC"
	if filter.SortBy != "" {
		sortOrder := "DESC"
		if strings.ToLower(filter.SortOrder) == "asc" {
			sortOrder = "ASC"
		}
		sortBy = fmt.Sprintf("staff_attendance.%s %s", filter.SortBy, sortOrder)
	}
	query = query.Order(sortBy)

	// Apply limit
	limit := 20
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Limit(limit + 1)

	var attendanceList []models.StaffAttendance
	if err := query.Find(&attendanceList).Error; err != nil {
		return nil, "", 0, fmt.Errorf("list attendance: %w", err)
	}

	// Calculate next cursor
	var nextCursor string
	if len(attendanceList) > limit {
		attendanceList = attendanceList[:limit]
		nextCursor = attendanceList[len(attendanceList)-1].ID.String()
	}

	return attendanceList, nextCursor, total, nil
}

// GetMonthlyAttendance retrieves attendance records for a staff member in a specific month.
func (r *Repository) GetMonthlyAttendance(ctx context.Context, tenantID, staffID uuid.UUID, year int, month int) ([]models.StaffAttendance, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	var attendanceList []models.StaffAttendance
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND staff_id = ? AND attendance_date >= ? AND attendance_date <= ?",
			tenantID, staffID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("attendance_date ASC").
		Find(&attendanceList).Error
	if err != nil {
		return nil, fmt.Errorf("get monthly attendance: %w", err)
	}
	return attendanceList, nil
}

// CreateRegularization creates a new regularization request.
func (r *Repository) CreateRegularization(ctx context.Context, regularization *models.StaffAttendanceRegularization) error {
	if err := r.db.WithContext(ctx).Create(regularization).Error; err != nil {
		return fmt.Errorf("create regularization: %w", err)
	}
	return nil
}

// UpdateRegularization updates a regularization request.
func (r *Repository) UpdateRegularization(ctx context.Context, regularization *models.StaffAttendanceRegularization) error {
	result := r.db.WithContext(ctx).Save(regularization)
	if result.Error != nil {
		return fmt.Errorf("update regularization: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrRegularizationNotFound
	}
	return nil
}

// GetRegularizationByID retrieves a regularization request by ID.
func (r *Repository) GetRegularizationByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StaffAttendanceRegularization, error) {
	var regularization models.StaffAttendanceRegularization
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Attendance").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&regularization).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRegularizationNotFound
		}
		return nil, fmt.Errorf("get regularization by id: %w", err)
	}
	return &regularization, nil
}

// GetPendingRegularizationByStaffAndDate checks if there's a pending regularization request for a staff member on a specific date.
func (r *Repository) GetPendingRegularizationByStaffAndDate(ctx context.Context, tenantID, staffID uuid.UUID, date time.Time) (*models.StaffAttendanceRegularization, error) {
	var regularization models.StaffAttendanceRegularization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND staff_id = ? AND request_date = ? AND status = ?",
			tenantID, staffID, date.Format("2006-01-02"), models.RegularizationStatusPending).
		First(&regularization).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get pending regularization: %w", err)
	}
	return &regularization, nil
}

// ListRegularizations retrieves regularization requests with optional filtering.
func (r *Repository) ListRegularizations(ctx context.Context, filter RegularizationFilter) ([]models.StaffAttendanceRegularization, string, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.StaffAttendanceRegularization{}).
		Preload("Staff").
		Preload("Attendance").
		Where("staff_attendance_regularization.tenant_id = ?", filter.TenantID)

	// Apply filters
	if filter.StaffID != nil {
		query = query.Where("staff_attendance_regularization.staff_id = ?", *filter.StaffID)
	}

	if filter.Status != nil {
		query = query.Where("staff_attendance_regularization.status = ?", *filter.Status)
	}

	if filter.DateFrom != nil {
		query = query.Where("staff_attendance_regularization.request_date >= ?", filter.DateFrom.Format("2006-01-02"))
	}

	if filter.DateTo != nil {
		query = query.Where("staff_attendance_regularization.request_date <= ?", filter.DateTo.Format("2006-01-02"))
	}

	// Get total count
	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, "", 0, fmt.Errorf("count regularizations: %w", err)
	}

	// Apply cursor pagination
	if filter.Cursor != "" {
		cursorID, err := uuid.Parse(filter.Cursor)
		if err == nil {
			query = query.Where("staff_attendance_regularization.id > ?", cursorID)
		}
	}

	// Sort by created_at DESC
	query = query.Order("staff_attendance_regularization.created_at DESC")

	// Apply limit
	limit := 20
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Limit(limit + 1)

	var regularizations []models.StaffAttendanceRegularization
	if err := query.Find(&regularizations).Error; err != nil {
		return nil, "", 0, fmt.Errorf("list regularizations: %w", err)
	}

	// Calculate next cursor
	var nextCursor string
	if len(regularizations) > limit {
		regularizations = regularizations[:limit]
		nextCursor = regularizations[len(regularizations)-1].ID.String()
	}

	return regularizations, nextCursor, total, nil
}

// GetSettings retrieves attendance settings for a branch.
func (r *Repository) GetSettings(ctx context.Context, tenantID, branchID uuid.UUID) (*models.StaffAttendanceSettings, error) {
	var settings models.StaffAttendanceSettings
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Where("tenant_id = ? AND branch_id = ?", tenantID, branchID).
		First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSettingsNotFound
		}
		return nil, fmt.Errorf("get settings: %w", err)
	}
	return &settings, nil
}

// CreateOrUpdateSettings creates or updates attendance settings for a branch.
func (r *Repository) CreateOrUpdateSettings(ctx context.Context, settings *models.StaffAttendanceSettings) error {
	// Use upsert
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ?", settings.TenantID, settings.BranchID).
		Assign(models.StaffAttendanceSettings{
			WorkStartTime:                 settings.WorkStartTime,
			WorkEndTime:                   settings.WorkEndTime,
			LateThresholdMinutes:          settings.LateThresholdMinutes,
			HalfDayThresholdHours:         settings.HalfDayThresholdHours,
			AllowSelfCheckout:             settings.AllowSelfCheckout,
			RequireRegularizationApproval: settings.RequireRegularizationApproval,
			UpdatedAt:                     time.Now(),
		}).
		FirstOrCreate(settings).Error
	if err != nil {
		return fmt.Errorf("create or update settings: %w", err)
	}
	return nil
}

// DB returns the underlying database connection.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
