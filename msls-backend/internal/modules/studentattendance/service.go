// Package studentattendance provides student attendance management functionality.
package studentattendance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Service handles student attendance business logic.
type Service struct {
	repo *Repository
	db   *gorm.DB
}

// NewService creates a new student attendance service.
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: NewRepository(db),
		db:   db,
	}
}

// GetTeacherSections returns the sections available for a teacher to mark attendance.
// For now, returns all sections. In the future, this should be filtered by teacher assignment.
func (s *Service) GetTeacherSections(ctx context.Context, tenantID, branchID uuid.UUID, date time.Time) ([]TeacherClassResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	// Get all active sections for the branch
	var sections []models.Section
	err := s.db.WithContext(ctx).
		Preload("Class").
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Order("display_order ASC").
		Find(&sections).Error
	if err != nil {
		return nil, fmt.Errorf("get sections: %w", err)
	}

	dateStr := date.Format("2006-01-02")
	responses := make([]TeacherClassResponse, 0, len(sections))

	for _, section := range sections {
		// Count students in section
		studentCount, err := s.repo.CountStudentsInSection(ctx, tenantID, section.ID)
		if err != nil {
			continue
		}

		// Skip sections with no students
		if studentCount == 0 {
			continue
		}

		// Check if attendance is already marked
		markedCount, _ := s.repo.GetMarkedAttendanceCount(ctx, tenantID, section.ID, date)

		className := ""
		classCode := ""
		if section.Class.ID != uuid.Nil {
			className = section.Class.Name
			classCode = section.Class.Code
		}

		responses = append(responses, TeacherClassResponse{
			SectionID:     section.ID.String(),
			SectionName:   section.Name,
			SectionCode:   section.Code,
			ClassName:     className,
			ClassCode:     classCode,
			StudentCount:  int(studentCount),
			IsMarkedToday: markedCount > 0,
			MarkedCount:   int(markedCount),
		})

		_ = dateStr // Used for logging if needed
	}

	return responses, nil
}

// GetClassStudentsForAttendance returns students in a section with their attendance status for a date.
func (s *Service) GetClassStudentsForAttendance(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) (*ClassAttendanceResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if sectionID == uuid.Nil {
		return nil, ErrSectionIDRequired
	}

	// Get section details
	section, err := s.repo.GetSectionByID(ctx, tenantID, sectionID)
	if err != nil {
		return nil, err
	}

	// Get students in section
	students, err := s.repo.GetStudentsInSection(ctx, tenantID, sectionID)
	if err != nil {
		return nil, err
	}

	if len(students) == 0 {
		return nil, ErrNoStudentsInSection
	}

	// Get existing attendance records for this date
	existingAttendance, err := s.repo.GetClassAttendance(ctx, tenantID, sectionID, date)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	attendanceMap := make(map[uuid.UUID]models.StudentAttendance)
	for _, a := range existingAttendance {
		attendanceMap[a.StudentID] = a
	}

	// Get attendance history for each student (last 5 days)
	studentResponses := make([]StudentForAttendance, len(students))
	for i, student := range students {
		// Get last 5 days history
		history, _ := s.repo.GetStudentAttendanceHistory(ctx, tenantID, student.ID, 5)
		last5Days := make([]string, len(history))
		for j, h := range history {
			last5Days[j] = h.Status.ShortLabel()
		}

		resp := StudentForAttendance{
			StudentID:       student.ID.String(),
			AdmissionNumber: student.AdmissionNumber,
			FirstName:       student.FirstName,
			LastName:        student.LastName,
			FullName:        student.FullName(),
			PhotoURL:        student.PhotoURL,
			Last5Days:       last5Days,
		}

		// Check if attendance already exists
		if attendance, ok := attendanceMap[student.ID]; ok {
			resp.Status = string(attendance.Status)
			resp.Remarks = attendance.Remarks
			if attendance.LateArrivalTime != nil {
				resp.LateArrivalTime = attendance.LateArrivalTime.Format("15:04")
			}
		}

		studentResponses[i] = resp
	}

	// Calculate summary
	summary := AttendanceSummary{Total: len(students)}
	for _, a := range existingAttendance {
		switch a.Status {
		case models.StudentAttendancePresent:
			summary.Present++
		case models.StudentAttendanceAbsent:
			summary.Absent++
		case models.StudentAttendanceLate:
			summary.Late++
		case models.StudentAttendanceHalfDay:
			summary.HalfDay++
		}
	}

	// Check if can edit (based on edit window)
	canEdit := true
	isMarked := len(existingAttendance) > 0
	var markedAt, markedByName string

	if isMarked {
		// Get first attendance record to check marked time
		firstRecord, _ := s.repo.GetFirstAttendanceRecord(ctx, tenantID, sectionID, date)
		if firstRecord != nil {
			markedAt = firstRecord.MarkedAt.Format(time.RFC3339)
			if firstRecord.MarkedByUser.ID != uuid.Nil {
				markedByName = firstRecord.MarkedByUser.FirstName + " " + firstRecord.MarkedByUser.LastName
			}

			// Check edit window
			canEdit, _ = s.CanEditAttendance(ctx, tenantID, sectionID, date)
		}
	}

	className := ""
	if section.Class.ID != uuid.Nil {
		className = section.Class.Name
	}

	return &ClassAttendanceResponse{
		SectionID:    sectionID.String(),
		SectionName:  section.Name,
		ClassName:    className,
		Date:         date.Format("2006-01-02"),
		Students:     studentResponses,
		IsMarked:     isMarked,
		CanEdit:      canEdit,
		MarkedAt:     markedAt,
		MarkedByName: markedByName,
		Summary:      summary,
	}, nil
}

// MarkClassAttendance marks attendance for all students in a class section.
func (s *Service) MarkClassAttendance(ctx context.Context, dto MarkClassAttendanceDTO) (*MarkAttendanceResult, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.SectionID == uuid.Nil {
		return nil, ErrSectionIDRequired
	}
	if dto.Date.IsZero() {
		return nil, ErrDateRequired
	}
	if dto.Date.After(time.Now()) {
		return nil, ErrFutureDate
	}
	if len(dto.Records) == 0 {
		return nil, ErrEmptyAttendanceRecords
	}

	// Verify section exists
	section, err := s.repo.GetSectionByID(ctx, dto.TenantID, dto.SectionID)
	if err != nil {
		return nil, err
	}

	// Check if can edit (if already marked)
	existingCount, _ := s.repo.GetMarkedAttendanceCount(ctx, dto.TenantID, dto.SectionID, dto.Date)
	if existingCount > 0 {
		canEdit, err := s.CanEditAttendance(ctx, dto.TenantID, dto.SectionID, dto.Date)
		if err != nil {
			return nil, err
		}
		if !canEdit {
			return nil, ErrEditWindowExpired
		}
	}

	// Validate and prepare attendance records
	now := time.Now()
	attendanceRecords := make([]models.StudentAttendance, len(dto.Records))
	summary := AttendanceSummary{Total: len(dto.Records)}

	for i, record := range dto.Records {
		// Validate status
		status := AttendanceStatus(record.Status)
		if !status.IsValid() {
			return nil, ErrInvalidStatus
		}

		attendanceRecords[i] = models.StudentAttendance{
			TenantID:        dto.TenantID,
			StudentID:       record.StudentID,
			SectionID:       dto.SectionID,
			AttendanceDate:  dto.Date,
			Status:          models.StudentAttendanceStatus(status),
			LateArrivalTime: record.LateArrivalTime,
			Remarks:         record.Remarks,
			MarkedBy:        dto.MarkedBy,
			MarkedAt:        now,
		}

		// Update summary
		switch status {
		case StatusPresent:
			summary.Present++
		case StatusAbsent:
			summary.Absent++
		case StatusLate:
			summary.Late++
		case StatusHalfDay:
			summary.HalfDay++
		}
	}

	// Save attendance records
	if err := s.repo.SaveClassAttendance(ctx, dto.TenantID, attendanceRecords); err != nil {
		return nil, err
	}

	_ = section // Used for validation

	return &MarkAttendanceResult{
		SectionID: dto.SectionID.String(),
		Date:      dto.Date.Format("2006-01-02"),
		Summary:   summary,
		MarkedAt:  now.Format(time.RFC3339),
		Message:   fmt.Sprintf("Attendance marked: %d present, %d absent, %d late, %d half-day", summary.Present, summary.Absent, summary.Late, summary.HalfDay),
	}, nil
}

// CanEditAttendance checks if attendance can still be edited based on edit window.
func (s *Service) CanEditAttendance(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) (bool, error) {
	// Get first attendance record to check marked time
	firstRecord, err := s.repo.GetFirstAttendanceRecord(ctx, tenantID, sectionID, date)
	if err != nil {
		return false, err
	}
	if firstRecord == nil {
		// No attendance marked yet, can edit
		return true, nil
	}

	// Get section to find branch
	section, err := s.repo.GetSectionByID(ctx, tenantID, sectionID)
	if err != nil {
		return false, err
	}

	// Get settings for the branch (if section has branch info)
	// Default to 120 minutes if settings not found
	editWindowMinutes := 120

	// Try to get branch-specific settings
	// Note: Section might not have direct branch reference, we use class's branch
	if section.Class.ID != uuid.Nil {
		// For now, use default. In production, we'd look up branch from class
		settings, err := s.repo.GetSettings(ctx, tenantID, tenantID) // Using tenantID as placeholder
		if err == nil {
			editWindowMinutes = settings.EditWindowMinutes
		}
	}

	// Check if within edit window
	windowEnd := firstRecord.MarkedAt.Add(time.Duration(editWindowMinutes) * time.Minute)
	return time.Now().Before(windowEnd), nil
}

// GetAttendanceByID retrieves an attendance record by ID.
func (s *Service) GetAttendanceByID(ctx context.Context, tenantID, id uuid.UUID) (*models.StudentAttendance, error) {
	return s.repo.GetAttendanceByID(ctx, tenantID, id)
}

// ListAttendance retrieves attendance records with filtering.
func (s *Service) ListAttendance(ctx context.Context, filter ListFilter) ([]models.StudentAttendance, string, int64, error) {
	return s.repo.ListAttendance(ctx, filter)
}

// GetSettings retrieves student attendance settings for a branch.
func (s *Service) GetSettings(ctx context.Context, tenantID, branchID uuid.UUID) (*models.StudentAttendanceSettings, error) {
	return s.repo.GetSettings(ctx, tenantID, branchID)
}

// UpdateSettings updates student attendance settings for a branch.
func (s *Service) UpdateSettings(ctx context.Context, dto SettingsDTO) (*models.StudentAttendanceSettings, error) {
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.BranchID == uuid.Nil {
		return nil, ErrBranchIDRequired
	}

	settings := &models.StudentAttendanceSettings{
		TenantID:             dto.TenantID,
		BranchID:             dto.BranchID,
		EditWindowMinutes:    dto.EditWindowMinutes,
		LateThresholdMinutes: dto.LateThresholdMinutes,
		SMSOnAbsent:          dto.SMSOnAbsent,
	}

	if err := s.repo.CreateOrUpdateSettings(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}
