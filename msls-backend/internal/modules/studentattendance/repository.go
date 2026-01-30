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

// ============================================================================
// Period-wise Attendance Repository Methods (Story 7.2)
// ============================================================================

// GetPeriodAttendance retrieves attendance records for a section and period on a specific date.
func (r *Repository) GetPeriodAttendance(ctx context.Context, tenantID, sectionID, periodID uuid.UUID, date time.Time) ([]models.StudentAttendance, error) {
	var attendance []models.StudentAttendance
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("MarkedByUser").
		Preload("PeriodSlot").
		Preload("TimetableEntry").
		Preload("TimetableEntry.Subject").
		Where("tenant_id = ? AND section_id = ? AND period_id = ? AND attendance_date = ?",
			tenantID, sectionID, periodID, date.Format("2006-01-02")).
		Order("created_at ASC").
		Find(&attendance).Error
	if err != nil {
		return nil, fmt.Errorf("get period attendance: %w", err)
	}
	return attendance, nil
}

// SavePeriodAttendance saves or updates period-specific attendance records.
func (r *Repository) SavePeriodAttendance(ctx context.Context, tenantID uuid.UUID, records []models.StudentAttendance) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			// Use upsert with period_id in the match criteria
			err := tx.Where("tenant_id = ? AND student_id = ? AND attendance_date = ? AND period_id = ?",
				tenantID, record.StudentID, record.AttendanceDate.Format("2006-01-02"), record.PeriodID).
				Assign(models.StudentAttendance{
					SectionID:        record.SectionID,
					PeriodID:         record.PeriodID,
					TimetableEntryID: record.TimetableEntryID,
					Status:           record.Status,
					LateArrivalTime:  record.LateArrivalTime,
					Remarks:          record.Remarks,
					MarkedBy:         record.MarkedBy,
					MarkedAt:         record.MarkedAt,
					UpdatedAt:        time.Now(),
				}).
				FirstOrCreate(&record).Error
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key") {
					return ErrDuplicateAttendance
				}
				return fmt.Errorf("save period attendance record: %w", err)
			}
		}
		return nil
	})
}

// GetDailyPeriodAttendance retrieves all period attendance records for a section on a date.
func (r *Repository) GetDailyPeriodAttendance(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) ([]models.StudentAttendance, error) {
	var attendance []models.StudentAttendance
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("PeriodSlot").
		Preload("TimetableEntry").
		Preload("TimetableEntry.Subject").
		Where("tenant_id = ? AND section_id = ? AND attendance_date = ? AND period_id IS NOT NULL",
			tenantID, sectionID, date.Format("2006-01-02")).
		Order("period_id ASC, student_id ASC").
		Find(&attendance).Error
	if err != nil {
		return nil, fmt.Errorf("get daily period attendance: %w", err)
	}
	return attendance, nil
}

// GetPeriodAttendanceCount returns the count of marked attendance for a period.
func (r *Repository) GetPeriodAttendanceCount(ctx context.Context, tenantID, sectionID, periodID uuid.UUID, date time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Where("tenant_id = ? AND section_id = ? AND period_id = ? AND attendance_date = ?",
			tenantID, sectionID, periodID, date.Format("2006-01-02")).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count period attendance: %w", err)
	}
	return count, nil
}

// GetFirstPeriodAttendanceRecord returns the first attendance record for a period.
func (r *Repository) GetFirstPeriodAttendanceRecord(ctx context.Context, tenantID, sectionID, periodID uuid.UUID, date time.Time) (*models.StudentAttendance, error) {
	var attendance models.StudentAttendance
	err := r.db.WithContext(ctx).
		Preload("MarkedByUser").
		Where("tenant_id = ? AND section_id = ? AND period_id = ? AND attendance_date = ?",
			tenantID, sectionID, periodID, date.Format("2006-01-02")).
		Order("marked_at ASC").
		First(&attendance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get first period attendance record: %w", err)
	}
	return &attendance, nil
}

// GetTimetableEntriesForSection retrieves timetable entries for a section on a specific day.
func (r *Repository) GetTimetableEntriesForSection(ctx context.Context, tenantID, sectionID uuid.UUID, dayOfWeek int) ([]models.TimetableEntry, error) {
	var entries []models.TimetableEntry

	// Find published timetable for section
	var timetable models.Timetable
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND section_id = ? AND status = ? AND deleted_at IS NULL",
			tenantID, sectionID, models.TimetableStatusPublished).
		First(&timetable).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No published timetable
		}
		return nil, fmt.Errorf("get published timetable: %w", err)
	}

	// Get entries for this day
	err = r.db.WithContext(ctx).
		Preload("PeriodSlot").
		Preload("Subject").
		Preload("Staff").
		Where("tenant_id = ? AND timetable_id = ? AND day_of_week = ?",
			tenantID, timetable.ID, dayOfWeek).
		Order("period_slot_id ASC").
		Find(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("get timetable entries: %w", err)
	}
	return entries, nil
}

// GetPeriodSlotByID retrieves a period slot by ID.
func (r *Repository) GetPeriodSlotByID(ctx context.Context, tenantID, periodID uuid.UUID) (*models.PeriodSlot, error) {
	var slot models.PeriodSlot
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, periodID).
		First(&slot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPeriodNotFound
		}
		return nil, fmt.Errorf("get period slot: %w", err)
	}
	return &slot, nil
}

// GetTimetableEntryByID retrieves a timetable entry by ID.
func (r *Repository) GetTimetableEntryByID(ctx context.Context, tenantID, entryID uuid.UUID) (*models.TimetableEntry, error) {
	var entry models.TimetableEntry
	err := r.db.WithContext(ctx).
		Preload("PeriodSlot").
		Preload("Subject").
		Preload("Staff").
		Preload("Timetable").
		Preload("Timetable.Section").
		Where("tenant_id = ? AND id = ?", tenantID, entryID).
		First(&entry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTimetableEntryNotFound
		}
		return nil, fmt.Errorf("get timetable entry: %w", err)
	}
	return &entry, nil
}

// GetSubjectAttendanceForStudent retrieves period attendance for a student by subject.
func (r *Repository) GetSubjectAttendanceForStudent(ctx context.Context, tenantID, studentID, subjectID uuid.UUID, dateFrom, dateTo *time.Time) ([]models.StudentAttendance, error) {
	query := r.db.WithContext(ctx).
		Preload("TimetableEntry").
		Preload("TimetableEntry.Subject").
		Preload("PeriodSlot").
		Joins("JOIN timetable_entries ON timetable_entries.id = student_attendance.timetable_entry_id").
		Where("student_attendance.tenant_id = ? AND student_attendance.student_id = ? AND timetable_entries.subject_id = ? AND student_attendance.period_id IS NOT NULL",
			tenantID, studentID, subjectID)

	if dateFrom != nil {
		query = query.Where("student_attendance.attendance_date >= ?", dateFrom.Format("2006-01-02"))
	}
	if dateTo != nil {
		query = query.Where("student_attendance.attendance_date <= ?", dateTo.Format("2006-01-02"))
	}

	var attendance []models.StudentAttendance
	err := query.Order("student_attendance.attendance_date DESC").Find(&attendance).Error
	if err != nil {
		return nil, fmt.Errorf("get subject attendance: %w", err)
	}
	return attendance, nil
}

// ============================================================================
// Attendance Audit Repository Methods (Story 7.3)
// ============================================================================

// CreateAuditRecord creates a new audit record for an attendance change.
func (r *Repository) CreateAuditRecord(ctx context.Context, audit *models.StudentAttendanceAudit) error {
	if err := r.db.WithContext(ctx).Create(audit).Error; err != nil {
		return fmt.Errorf("create audit record: %w", err)
	}
	return nil
}

// GetAuditTrail retrieves all audit records for an attendance entry.
func (r *Repository) GetAuditTrail(ctx context.Context, tenantID, attendanceID uuid.UUID) ([]models.StudentAttendanceAudit, error) {
	var audits []models.StudentAttendanceAudit
	err := r.db.WithContext(ctx).
		Preload("ChangedByUser").
		Where("tenant_id = ? AND attendance_id = ?", tenantID, attendanceID).
		Order("changed_at DESC").
		Find(&audits).Error
	if err != nil {
		return nil, fmt.Errorf("get audit trail: %w", err)
	}
	return audits, nil
}

// UpdateAttendanceWithAudit updates an attendance record and creates an audit entry.
func (r *Repository) UpdateAttendanceWithAudit(ctx context.Context, attendance *models.StudentAttendance, audit *models.StudentAttendanceAudit) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update attendance record
		if err := tx.Save(attendance).Error; err != nil {
			return fmt.Errorf("update attendance: %w", err)
		}

		// Create audit record
		if err := tx.Create(audit).Error; err != nil {
			return fmt.Errorf("create audit: %w", err)
		}

		return nil
	})
}

// GetUserByID retrieves a user by ID.
func (r *Repository) GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, userID).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

// ============================================================================
// Calendar & Reports Repository Methods (Stories 7.4-7.8)
// ============================================================================

// GetStudentMonthlyAttendance retrieves a student's attendance for a specific month.
func (r *Repository) GetStudentMonthlyAttendance(ctx context.Context, tenantID, studentID uuid.UUID, year, month int) ([]models.StudentAttendance, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	var attendance []models.StudentAttendance
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ? AND attendance_date >= ? AND attendance_date <= ? AND period_id IS NULL",
			tenantID, studentID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("attendance_date ASC").
		Find(&attendance).Error
	if err != nil {
		return nil, fmt.Errorf("get student monthly attendance: %w", err)
	}
	return attendance, nil
}

// GetStudentAttendanceRange retrieves attendance for a date range.
func (r *Repository) GetStudentAttendanceRange(ctx context.Context, tenantID, studentID uuid.UUID, dateFrom, dateTo time.Time) ([]models.StudentAttendance, error) {
	var attendance []models.StudentAttendance
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ? AND attendance_date >= ? AND attendance_date <= ? AND period_id IS NULL",
			tenantID, studentID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")).
		Order("attendance_date ASC").
		Find(&attendance).Error
	if err != nil {
		return nil, fmt.Errorf("get student attendance range: %w", err)
	}
	return attendance, nil
}

// GetSectionAttendanceStats retrieves attendance statistics for a section.
func (r *Repository) GetSectionAttendanceStats(ctx context.Context, tenantID, sectionID uuid.UUID, dateFrom, dateTo time.Time) (float64, error) {
	var result struct {
		TotalPresent int64
		TotalRecords int64
	}

	err := r.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Select("COUNT(CASE WHEN status IN ('present', 'late', 'half_day') THEN 1 END) as total_present, COUNT(*) as total_records").
		Where("tenant_id = ? AND section_id = ? AND attendance_date >= ? AND attendance_date <= ? AND period_id IS NULL",
			tenantID, sectionID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")).
		Scan(&result).Error
	if err != nil {
		return 0, fmt.Errorf("get section attendance stats: %w", err)
	}

	if result.TotalRecords == 0 {
		return 0, nil
	}
	return float64(result.TotalPresent) / float64(result.TotalRecords) * 100, nil
}

// GetClassMonthlyAttendance retrieves all attendance for a section for a month.
func (r *Repository) GetClassMonthlyAttendance(ctx context.Context, tenantID, sectionID uuid.UUID, year, month int) ([]models.StudentAttendance, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	var attendance []models.StudentAttendance
	err := r.db.WithContext(ctx).
		Preload("Student").
		Where("tenant_id = ? AND section_id = ? AND attendance_date >= ? AND attendance_date <= ? AND period_id IS NULL",
			tenantID, sectionID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("student_id, attendance_date").
		Find(&attendance).Error
	if err != nil {
		return nil, fmt.Errorf("get class monthly attendance: %w", err)
	}
	return attendance, nil
}

// GetLowAttendanceStudents retrieves students below attendance threshold.
func (r *Repository) GetLowAttendanceStudents(ctx context.Context, tenantID uuid.UUID, dateFrom, dateTo time.Time, threshold float64) ([]struct {
	StudentID    uuid.UUID
	SectionID    uuid.UUID
	TotalPresent int64
	TotalRecords int64
}, error) {
	var results []struct {
		StudentID    uuid.UUID
		SectionID    uuid.UUID
		TotalPresent int64
		TotalRecords int64
	}

	subquery := r.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Select("student_id, section_id, COUNT(CASE WHEN status IN ('present', 'late', 'half_day') THEN 1 END) as total_present, COUNT(*) as total_records").
		Where("tenant_id = ? AND attendance_date >= ? AND attendance_date <= ? AND period_id IS NULL",
			tenantID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")).
		Group("student_id, section_id").
		Having("CAST(COUNT(CASE WHEN status IN ('present', 'late', 'half_day') THEN 1 END) AS FLOAT) / NULLIF(COUNT(*), 0) * 100 < ?", threshold)

	err := subquery.Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("get low attendance students: %w", err)
	}
	return results, nil
}

// GetUnmarkedSections retrieves sections without attendance for a date.
func (r *Repository) GetUnmarkedSections(ctx context.Context, tenantID uuid.UUID, date time.Time) ([]uuid.UUID, error) {
	var allSections []models.Section
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Find(&allSections).Error
	if err != nil {
		return nil, fmt.Errorf("get all sections: %w", err)
	}

	var markedSections []uuid.UUID
	err = r.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Distinct("section_id").
		Where("tenant_id = ? AND attendance_date = ? AND period_id IS NULL", tenantID, date.Format("2006-01-02")).
		Pluck("section_id", &markedSections).Error
	if err != nil {
		return nil, fmt.Errorf("get marked sections: %w", err)
	}

	markedMap := make(map[uuid.UUID]bool)
	for _, id := range markedSections {
		markedMap[id] = true
	}

	var unmarked []uuid.UUID
	for _, section := range allSections {
		if !markedMap[section.ID] {
			unmarked = append(unmarked, section.ID)
		}
	}

	return unmarked, nil
}

// GetDailyAttendanceSummary retrieves overall attendance summary for a date.
func (r *Repository) GetDailyAttendanceSummary(ctx context.Context, tenantID uuid.UUID, date time.Time) (present, absent, total int64, err error) {
	var result struct {
		Present int64
		Absent  int64
		Total   int64
	}

	err = r.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Select("COUNT(CASE WHEN status = 'present' THEN 1 END) as present, COUNT(CASE WHEN status = 'absent' THEN 1 END) as absent, COUNT(*) as total").
		Where("tenant_id = ? AND attendance_date = ? AND period_id IS NULL", tenantID, date.Format("2006-01-02")).
		Scan(&result).Error
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get daily summary: %w", err)
	}

	return result.Present, result.Absent, result.Total, nil
}

// GetStudentByID retrieves a student by ID.
func (r *Repository) GetStudentByID(ctx context.Context, tenantID, studentID uuid.UUID) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, studentID).
		First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("get student: %w", err)
	}
	return &student, nil
}

// GetStudentCurrentSection retrieves a student's current section from student_sections.
func (r *Repository) GetStudentCurrentSection(ctx context.Context, tenantID, studentID uuid.UUID) (*models.Section, error) {
	var section models.Section
	err := r.db.WithContext(ctx).
		Preload("Class").
		Joins("JOIN student_sections ON student_sections.section_id = sections.id").
		Where("student_sections.student_id = ? AND sections.tenant_id = ?", studentID, tenantID).
		First(&section).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get student section: %w", err)
	}
	return &section, nil
}

// GetAllActiveSections retrieves all active sections.
func (r *Repository) GetAllActiveSections(ctx context.Context, tenantID uuid.UUID) ([]models.Section, error) {
	var sections []models.Section
	err := r.db.WithContext(ctx).
		Preload("Class").
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Order("display_order").
		Find(&sections).Error
	if err != nil {
		return nil, fmt.Errorf("get all sections: %w", err)
	}
	return sections, nil
}
