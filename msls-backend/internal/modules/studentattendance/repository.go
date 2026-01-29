// Package studentattendance provides student attendance management functionality.
package studentattendance

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

// Repository handles database operations for student attendance.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new student attendance repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetClassAttendance retrieves attendance records for a section on a specific date.
func (r *Repository) GetClassAttendance(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) ([]models.StudentAttendance, error) {
	var attendance []models.StudentAttendance
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("MarkedByUser").
		Where("tenant_id = ? AND section_id = ? AND attendance_date = ?",
			tenantID, sectionID, date.Format("2006-01-02")).
		Order("created_at ASC").
		Find(&attendance).Error
	if err != nil {
		return nil, fmt.Errorf("get class attendance: %w", err)
	}
	return attendance, nil
}

// GetStudentAttendance retrieves attendance for a specific student on a date.
func (r *Repository) GetStudentAttendance(ctx context.Context, tenantID, studentID uuid.UUID, date time.Time) (*models.StudentAttendance, error) {
	var attendance models.StudentAttendance
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Section").
		Preload("MarkedByUser").
		Where("tenant_id = ? AND student_id = ? AND attendance_date = ?",
			tenantID, studentID, date.Format("2006-01-02")).
		First(&attendance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, fmt.Errorf("get student attendance: %w", err)
	}
	return &attendance, nil
}

// GetStudentAttendanceHistory retrieves the last N days of attendance for a student.
func (r *Repository) GetStudentAttendanceHistory(ctx context.Context, tenantID, studentID uuid.UUID, days int) ([]StudentAttendanceHistoryEntry, error) {
	endDate := time.Now().Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -days)

	var attendance []models.StudentAttendance
	err := r.db.WithContext(ctx).
		Select("attendance_date, status").
		Where("tenant_id = ? AND student_id = ? AND attendance_date >= ? AND attendance_date <= ?",
			tenantID, studentID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("attendance_date DESC").
		Find(&attendance).Error
	if err != nil {
		return nil, fmt.Errorf("get student attendance history: %w", err)
	}

	history := make([]StudentAttendanceHistoryEntry, len(attendance))
	for i, a := range attendance {
		history[i] = StudentAttendanceHistoryEntry{
			Date:   a.AttendanceDate,
			Status: AttendanceStatus(a.Status),
		}
	}
	return history, nil
}

// SaveClassAttendance saves or updates attendance records for a class section.
func (r *Repository) SaveClassAttendance(ctx context.Context, tenantID uuid.UUID, records []models.StudentAttendance) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			// Use upsert - update if exists, create if not
			err := tx.Where("tenant_id = ? AND student_id = ? AND attendance_date = ?",
				tenantID, record.StudentID, record.AttendanceDate.Format("2006-01-02")).
				Assign(models.StudentAttendance{
					SectionID:       record.SectionID,
					Status:          record.Status,
					LateArrivalTime: record.LateArrivalTime,
					Remarks:         record.Remarks,
					MarkedBy:        record.MarkedBy,
					MarkedAt:        record.MarkedAt,
					UpdatedAt:       time.Now(),
				}).
				FirstOrCreate(&record).Error
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key") {
					return ErrDuplicateAttendance
				}
				return fmt.Errorf("save attendance record: %w", err)
			}
		}
		return nil
	})
}

// GetAttendanceByID retrieves an attendance record by ID.
func (r *Repository) GetAttendanceByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StudentAttendance, error) {
	var attendance models.StudentAttendance
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Section").
		Preload("Section.Class").
		Preload("MarkedByUser").
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

// ListAttendance retrieves attendance records with optional filtering.
func (r *Repository) ListAttendance(ctx context.Context, filter ListFilter) ([]models.StudentAttendance, string, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Preload("Student").
		Preload("Section").
		Preload("Section.Class").
		Preload("MarkedByUser").
		Where("student_attendance.tenant_id = ?", filter.TenantID)

	// Apply filters
	if filter.SectionID != nil {
		query = query.Where("student_attendance.section_id = ?", *filter.SectionID)
	}
	if filter.StudentID != nil {
		query = query.Where("student_attendance.student_id = ?", *filter.StudentID)
	}
	if filter.Status != nil {
		query = query.Where("student_attendance.status = ?", *filter.Status)
	}
	if filter.DateFrom != nil {
		query = query.Where("student_attendance.attendance_date >= ?", filter.DateFrom.Format("2006-01-02"))
	}
	if filter.DateTo != nil {
		query = query.Where("student_attendance.attendance_date <= ?", filter.DateTo.Format("2006-01-02"))
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
			query = query.Where("student_attendance.id > ?", cursorID)
		}
	}

	// Apply sorting
	sortBy := "student_attendance.attendance_date DESC, student_attendance.created_at DESC"
	if filter.SortBy != "" {
		sortOrder := "DESC"
		if strings.ToLower(filter.SortOrder) == "asc" {
			sortOrder = "ASC"
		}
		sortBy = fmt.Sprintf("student_attendance.%s %s", filter.SortBy, sortOrder)
	}
	query = query.Order(sortBy)

	// Apply limit
	limit := 20
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Limit(limit + 1)

	var attendanceList []models.StudentAttendance
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

// GetSettings retrieves student attendance settings for a branch.
func (r *Repository) GetSettings(ctx context.Context, tenantID, branchID uuid.UUID) (*models.StudentAttendanceSettings, error) {
	var settings models.StudentAttendanceSettings
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

// CreateOrUpdateSettings creates or updates student attendance settings for a branch.
func (r *Repository) CreateOrUpdateSettings(ctx context.Context, settings *models.StudentAttendanceSettings) error {
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ?", settings.TenantID, settings.BranchID).
		Assign(models.StudentAttendanceSettings{
			EditWindowMinutes:    settings.EditWindowMinutes,
			LateThresholdMinutes: settings.LateThresholdMinutes,
			SMSOnAbsent:          settings.SMSOnAbsent,
			UpdatedAt:            time.Now(),
		}).
		FirstOrCreate(settings).Error
	if err != nil {
		return fmt.Errorf("create or update settings: %w", err)
	}
	return nil
}

// GetSectionByID retrieves a section by ID.
func (r *Repository) GetSectionByID(ctx context.Context, tenantID, sectionID uuid.UUID) (*models.Section, error) {
	var section models.Section
	err := r.db.WithContext(ctx).
		Preload("Class").
		Where("tenant_id = ? AND id = ?", tenantID, sectionID).
		First(&section).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSectionNotFound
		}
		return nil, fmt.Errorf("get section: %w", err)
	}
	return &section, nil
}

// GetStudentsInSection retrieves all active students in a section.
func (r *Repository) GetStudentsInSection(ctx context.Context, tenantID, sectionID uuid.UUID) ([]models.Student, error) {
	var students []models.Student

	// Get students enrolled in this section via student_sections table
	err := r.db.WithContext(ctx).
		Joins("JOIN student_sections ON student_sections.student_id = students.id").
		Where("students.tenant_id = ? AND student_sections.section_id = ? AND students.status = ?",
			tenantID, sectionID, models.StudentStatusActive).
		Order("students.first_name ASC, students.last_name ASC").
		Find(&students).Error
	if err != nil {
		return nil, fmt.Errorf("get students in section: %w", err)
	}
	return students, nil
}

// CountStudentsInSection returns the number of students in a section.
func (r *Repository) CountStudentsInSection(ctx context.Context, tenantID, sectionID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Joins("JOIN student_sections ON student_sections.student_id = students.id").
		Where("students.tenant_id = ? AND student_sections.section_id = ? AND students.status = ?",
			tenantID, sectionID, models.StudentStatusActive).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count students in section: %w", err)
	}
	return count, nil
}

// GetMarkedAttendanceCount returns the number of marked attendance records for a section on a date.
func (r *Repository) GetMarkedAttendanceCount(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Where("tenant_id = ? AND section_id = ? AND attendance_date = ?",
			tenantID, sectionID, date.Format("2006-01-02")).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count marked attendance: %w", err)
	}
	return count, nil
}

// GetFirstAttendanceRecord returns the first attendance record for a section on a date.
func (r *Repository) GetFirstAttendanceRecord(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) (*models.StudentAttendance, error) {
	var attendance models.StudentAttendance
	err := r.db.WithContext(ctx).
		Preload("MarkedByUser").
		Where("tenant_id = ? AND section_id = ? AND attendance_date = ?",
			tenantID, sectionID, date.Format("2006-01-02")).
		Order("marked_at ASC").
		First(&attendance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get first attendance record: %w", err)
	}
	return &attendance, nil
}

// DB returns the underlying database connection.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
