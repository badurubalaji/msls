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

// ============================================================================
// Period-wise Attendance Service Methods (Story 7.2)
// ============================================================================

// GetTeacherPeriods returns the periods available for a section on a given date.
func (s *Service) GetTeacherPeriods(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) (*TeacherPeriodsResponse, error) {
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

	// Get timetable entries for this day
	dayOfWeek := int(date.Weekday())
	entries, err := s.repo.GetTimetableEntriesForSection(ctx, tenantID, sectionID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, ErrNoTimetableForSection
	}

	// Get student count for section
	studentCount, _ := s.repo.CountStudentsInSection(ctx, tenantID, sectionID)

	// Build period responses
	periods := make([]PeriodInfoResponse, 0, len(entries))
	for _, entry := range entries {
		if entry.PeriodSlot == nil || !entry.PeriodSlot.IsTeachingPeriod() {
			continue // Skip breaks and non-teaching periods
		}

		// Check if attendance is marked for this period
		markedCount, _ := s.repo.GetPeriodAttendanceCount(ctx, tenantID, sectionID, entry.PeriodSlotID, date)

		periodResp := PeriodInfoResponse{
			PeriodSlotID:     entry.PeriodSlotID.String(),
			TimetableEntryID: entry.ID.String(),
			PeriodName:       entry.PeriodSlot.Name,
			PeriodNumber:     entry.PeriodSlot.PeriodNumber,
			StartTime:        entry.PeriodSlot.StartTime,
			EndTime:          entry.PeriodSlot.EndTime,
			IsMarked:         markedCount > 0,
			MarkedCount:      int(markedCount),
			TotalStudents:    int(studentCount),
		}

		if entry.Subject != nil {
			periodResp.SubjectID = entry.SubjectID.String()
			periodResp.SubjectName = entry.Subject.Name
			periodResp.SubjectCode = entry.Subject.Code
		}

		if entry.Staff != nil {
			periodResp.StaffID = entry.StaffID.String()
			periodResp.StaffName = entry.Staff.FirstName + " " + entry.Staff.LastName
		}

		periods = append(periods, periodResp)
	}

	className := ""
	if section.Class.ID != uuid.Nil {
		className = section.Class.Name
	}

	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	return &TeacherPeriodsResponse{
		SectionID:   sectionID.String(),
		SectionName: section.Name,
		ClassName:   className,
		Date:        date.Format("2006-01-02"),
		DayOfWeek:   dayOfWeek,
		DayName:     dayNames[dayOfWeek],
		Periods:     periods,
	}, nil
}

// GetPeriodAttendance returns students with their attendance status for a specific period.
func (s *Service) GetPeriodAttendance(ctx context.Context, tenantID, sectionID, periodID uuid.UUID, date time.Time) (*PeriodAttendanceResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if sectionID == uuid.Nil {
		return nil, ErrSectionIDRequired
	}
	if periodID == uuid.Nil {
		return nil, ErrPeriodIDRequired
	}

	// Get section details
	section, err := s.repo.GetSectionByID(ctx, tenantID, sectionID)
	if err != nil {
		return nil, err
	}

	// Get period slot details
	periodSlot, err := s.repo.GetPeriodSlotByID(ctx, tenantID, periodID)
	if err != nil {
		return nil, err
	}

	// Get timetable entry for this period on this day
	dayOfWeek := int(date.Weekday())
	entries, err := s.repo.GetTimetableEntriesForSection(ctx, tenantID, sectionID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	// Find the matching timetable entry
	var timetableEntry *models.TimetableEntry
	for _, entry := range entries {
		if entry.PeriodSlotID == periodID {
			timetableEntry = &entry
			break
		}
	}

	if timetableEntry == nil {
		return nil, ErrInvalidPeriodForSection
	}

	// Get students in section
	students, err := s.repo.GetStudentsInSection(ctx, tenantID, sectionID)
	if err != nil {
		return nil, err
	}

	if len(students) == 0 {
		return nil, ErrNoStudentsInSection
	}

	// Get existing attendance records for this period
	existingAttendance, err := s.repo.GetPeriodAttendance(ctx, tenantID, sectionID, periodID, date)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	attendanceMap := make(map[uuid.UUID]models.StudentAttendance)
	for _, a := range existingAttendance {
		attendanceMap[a.StudentID] = a
	}

	// Build student responses
	studentResponses := make([]StudentForPeriodAttendance, len(students))
	for i, student := range students {
		resp := StudentForPeriodAttendance{
			StudentID:       student.ID.String(),
			AdmissionNumber: student.AdmissionNumber,
			FirstName:       student.FirstName,
			LastName:        student.LastName,
			FullName:        student.FullName(),
			PhotoURL:        student.PhotoURL,
		}

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

	// Check if can edit
	canEdit := true
	isMarked := len(existingAttendance) > 0
	var markedAt, markedByName string

	if isMarked {
		firstRecord, _ := s.repo.GetFirstPeriodAttendanceRecord(ctx, tenantID, sectionID, periodID, date)
		if firstRecord != nil {
			markedAt = firstRecord.MarkedAt.Format(time.RFC3339)
			if firstRecord.MarkedByUser.ID != uuid.Nil {
				markedByName = firstRecord.MarkedByUser.FirstName + " " + firstRecord.MarkedByUser.LastName
			}
			canEdit, _ = s.CanEditPeriodAttendance(ctx, tenantID, sectionID, periodID, date)
		}
	}

	className := ""
	if section.Class.ID != uuid.Nil {
		className = section.Class.Name
	}

	response := &PeriodAttendanceResponse{
		SectionID:    sectionID.String(),
		SectionName:  section.Name,
		ClassName:    className,
		Date:         date.Format("2006-01-02"),
		PeriodSlotID: periodID.String(),
		PeriodName:   periodSlot.Name,
		PeriodNumber: periodSlot.PeriodNumber,
		StartTime:    periodSlot.StartTime,
		EndTime:      periodSlot.EndTime,
		Students:     studentResponses,
		IsMarked:     isMarked,
		CanEdit:      canEdit,
		MarkedAt:     markedAt,
		MarkedByName: markedByName,
		Summary:      summary,
	}

	if timetableEntry.Subject != nil {
		response.SubjectID = timetableEntry.SubjectID.String()
		response.SubjectName = timetableEntry.Subject.Name
		response.SubjectCode = timetableEntry.Subject.Code
	}
	response.TimetableEntryID = timetableEntry.ID.String()

	return response, nil
}

// MarkPeriodAttendance marks attendance for a specific period.
func (s *Service) MarkPeriodAttendance(ctx context.Context, dto MarkPeriodAttendanceDTO) (*MarkPeriodAttendanceResult, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.SectionID == uuid.Nil {
		return nil, ErrSectionIDRequired
	}
	if dto.PeriodID == uuid.Nil {
		return nil, ErrPeriodIDRequired
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
	_, err := s.repo.GetSectionByID(ctx, dto.TenantID, dto.SectionID)
	if err != nil {
		return nil, err
	}

	// Verify period slot exists
	_, err = s.repo.GetPeriodSlotByID(ctx, dto.TenantID, dto.PeriodID)
	if err != nil {
		return nil, err
	}

	// Check if can edit (if already marked)
	existingCount, _ := s.repo.GetPeriodAttendanceCount(ctx, dto.TenantID, dto.SectionID, dto.PeriodID, dto.Date)
	if existingCount > 0 {
		canEdit, err := s.CanEditPeriodAttendance(ctx, dto.TenantID, dto.SectionID, dto.PeriodID, dto.Date)
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
			TenantID:         dto.TenantID,
			StudentID:        record.StudentID,
			SectionID:        dto.SectionID,
			AttendanceDate:   dto.Date,
			PeriodID:         &dto.PeriodID,
			TimetableEntryID: &dto.TimetableEntryID,
			Status:           models.StudentAttendanceStatus(status),
			LateArrivalTime:  record.LateArrivalTime,
			Remarks:          record.Remarks,
			MarkedBy:         dto.MarkedBy,
			MarkedAt:         now,
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
	if err := s.repo.SavePeriodAttendance(ctx, dto.TenantID, attendanceRecords); err != nil {
		return nil, err
	}

	return &MarkPeriodAttendanceResult{
		SectionID: dto.SectionID.String(),
		PeriodID:  dto.PeriodID.String(),
		Date:      dto.Date.Format("2006-01-02"),
		Summary:   summary,
		MarkedAt:  now.Format(time.RFC3339),
		Message:   fmt.Sprintf("Period attendance marked: %d present, %d absent, %d late, %d half-day", summary.Present, summary.Absent, summary.Late, summary.HalfDay),
	}, nil
}

// CanEditPeriodAttendance checks if period attendance can still be edited.
func (s *Service) CanEditPeriodAttendance(ctx context.Context, tenantID, sectionID, periodID uuid.UUID, date time.Time) (bool, error) {
	firstRecord, err := s.repo.GetFirstPeriodAttendanceRecord(ctx, tenantID, sectionID, periodID, date)
	if err != nil {
		return false, err
	}
	if firstRecord == nil {
		return true, nil
	}

	// Default edit window
	editWindowMinutes := 120

	section, err := s.repo.GetSectionByID(ctx, tenantID, sectionID)
	if err == nil && section.Class.ID != uuid.Nil {
		settings, err := s.repo.GetSettings(ctx, tenantID, tenantID)
		if err == nil {
			editWindowMinutes = settings.EditWindowMinutes
		}
	}

	windowEnd := firstRecord.MarkedAt.Add(time.Duration(editWindowMinutes) * time.Minute)
	return time.Now().Before(windowEnd), nil
}

// GetDailySummary returns aggregated attendance summary for all periods in a day.
func (s *Service) GetDailySummary(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) (*DailySummaryResponse, error) {
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

	// Get timetable entries for this day
	dayOfWeek := int(date.Weekday())
	entries, err := s.repo.GetTimetableEntriesForSection(ctx, tenantID, sectionID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	// Filter teaching periods only
	teachingPeriods := make([]models.TimetableEntry, 0)
	for _, entry := range entries {
		if entry.PeriodSlot != nil && entry.PeriodSlot.IsTeachingPeriod() {
			teachingPeriods = append(teachingPeriods, entry)
		}
	}

	if len(teachingPeriods) == 0 {
		return nil, ErrNoTimetableForSection
	}

	// Get students in section
	students, err := s.repo.GetStudentsInSection(ctx, tenantID, sectionID)
	if err != nil {
		return nil, err
	}

	// Get all period attendance for the day
	allAttendance, err := s.repo.GetDailyPeriodAttendance(ctx, tenantID, sectionID, date)
	if err != nil {
		return nil, err
	}

	// Build attendance map: studentID -> periodID -> status
	attendanceByStudent := make(map[uuid.UUID]map[uuid.UUID]models.StudentAttendanceStatus)
	for _, a := range allAttendance {
		if a.PeriodID == nil {
			continue
		}
		if _, ok := attendanceByStudent[a.StudentID]; !ok {
			attendanceByStudent[a.StudentID] = make(map[uuid.UUID]models.StudentAttendanceStatus)
		}
		attendanceByStudent[a.StudentID][*a.PeriodID] = a.Status
	}

	// Build period info
	periods := make([]PeriodInfoResponse, len(teachingPeriods))
	for i, entry := range teachingPeriods {
		markedCount, _ := s.repo.GetPeriodAttendanceCount(ctx, tenantID, sectionID, entry.PeriodSlotID, date)
		periods[i] = PeriodInfoResponse{
			PeriodSlotID:     entry.PeriodSlotID.String(),
			TimetableEntryID: entry.ID.String(),
			PeriodName:       entry.PeriodSlot.Name,
			PeriodNumber:     entry.PeriodSlot.PeriodNumber,
			StartTime:        entry.PeriodSlot.StartTime,
			EndTime:          entry.PeriodSlot.EndTime,
			IsMarked:         markedCount > 0,
			MarkedCount:      int(markedCount),
			TotalStudents:    len(students),
		}
		if entry.Subject != nil {
			periods[i].SubjectID = entry.SubjectID.String()
			periods[i].SubjectName = entry.Subject.Name
			periods[i].SubjectCode = entry.Subject.Code
		}
	}

	// Build student summaries
	studentSummaries := make([]DailySummaryStudent, len(students))
	totalPresent := 0
	fullPresentCount := 0
	absentCount := 0

	for i, student := range students {
		periodStatuses := make(map[string]string)
		periodsPresent := 0
		periodsAbsent := 0
		periodsLate := 0

		studentAttendance := attendanceByStudent[student.ID]
		for _, entry := range teachingPeriods {
			if status, ok := studentAttendance[entry.PeriodSlotID]; ok {
				periodStatuses[entry.PeriodSlotID.String()] = string(status)
				switch status {
				case models.StudentAttendancePresent:
					periodsPresent++
				case models.StudentAttendanceAbsent:
					periodsAbsent++
				case models.StudentAttendanceLate:
					periodsLate++
					periodsPresent++ // Late counts as present for percentage
				case models.StudentAttendanceHalfDay:
					periodsPresent++ // Half-day counts as present for percentage
				}
			}
		}

		totalPeriods := len(teachingPeriods)
		attendancePercentage := 0.0
		if totalPeriods > 0 {
			attendancePercentage = float64(periodsPresent) / float64(totalPeriods) * 100
		}

		// Derive overall status (present if >50% periods attended)
		overallStatus := string(models.StudentAttendanceAbsent)
		if attendancePercentage > 50 {
			overallStatus = string(models.StudentAttendancePresent)
		}

		if periodsPresent == totalPeriods && totalPeriods > 0 {
			fullPresentCount++
		}
		if attendancePercentage <= 50 && totalPeriods > 0 {
			absentCount++
		}
		totalPresent += periodsPresent

		studentSummaries[i] = DailySummaryStudent{
			StudentID:            student.ID.String(),
			AdmissionNumber:      student.AdmissionNumber,
			FullName:             student.FullName(),
			PhotoURL:             student.PhotoURL,
			PeriodStatuses:       periodStatuses,
			TotalPeriods:         totalPeriods,
			PeriodsPresent:       periodsPresent,
			PeriodsAbsent:        periodsAbsent,
			PeriodsLate:          periodsLate,
			AttendancePercentage: attendancePercentage,
			OverallStatus:        overallStatus,
		}
	}

	// Calculate overall summary
	avgAttendance := 0.0
	totalPossible := len(students) * len(teachingPeriods)
	if totalPossible > 0 {
		avgAttendance = float64(totalPresent) / float64(totalPossible) * 100
	}

	className := ""
	if section.Class.ID != uuid.Nil {
		className = section.Class.Name
	}

	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	return &DailySummaryResponse{
		SectionID:   sectionID.String(),
		SectionName: section.Name,
		ClassName:   className,
		Date:        date.Format("2006-01-02"),
		DayName:     dayNames[dayOfWeek],
		Periods:     periods,
		Students:    studentSummaries,
		Summary: DailySummarySummary{
			TotalStudents:     len(students),
			TotalPeriods:      len(teachingPeriods),
			AverageAttendance: avgAttendance,
			FullPresentCount:  fullPresentCount,
			AbsentCount:       absentCount,
		},
	}, nil
}

// GetSubjectAttendance returns attendance statistics for a student by subject.
func (s *Service) GetSubjectAttendance(ctx context.Context, tenantID, studentID, subjectID uuid.UUID, dateFrom, dateTo *time.Time) (*SubjectAttendanceResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if studentID == uuid.Nil {
		return nil, ErrStudentIDRequired
	}
	if subjectID == uuid.Nil {
		return nil, ErrSubjectIDRequired
	}

	// Get attendance records for the subject
	attendance, err := s.repo.GetSubjectAttendanceForStudent(ctx, tenantID, studentID, subjectID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	totalPeriods := len(attendance)
	periodsPresent := 0
	periodsAbsent := 0
	periodsLate := 0

	var subjectName, subjectCode string

	for _, a := range attendance {
		switch a.Status {
		case models.StudentAttendancePresent:
			periodsPresent++
		case models.StudentAttendanceAbsent:
			periodsAbsent++
		case models.StudentAttendanceLate:
			periodsLate++
			periodsPresent++ // Late counts as present for percentage
		case models.StudentAttendanceHalfDay:
			periodsPresent++ // Half-day counts as present
		}

		// Extract subject info from first record with subject
		if a.TimetableEntry != nil && a.TimetableEntry.Subject != nil && subjectName == "" {
			subjectName = a.TimetableEntry.Subject.Name
			subjectCode = a.TimetableEntry.Subject.Code
		}
	}

	attendancePercentage := 0.0
	if totalPeriods > 0 {
		attendancePercentage = float64(periodsPresent) / float64(totalPeriods) * 100
	}

	// Minimum required is 75% for exam eligibility
	minimumRequired := 75.0
	isEligible := attendancePercentage >= minimumRequired

	return &SubjectAttendanceResponse{
		StudentID:            studentID.String(),
		SubjectID:            subjectID.String(),
		SubjectName:          subjectName,
		SubjectCode:          subjectCode,
		TotalPeriods:         totalPeriods,
		PeriodsPresent:       periodsPresent,
		PeriodsAbsent:        periodsAbsent,
		PeriodsLate:          periodsLate,
		AttendancePercentage: attendancePercentage,
		MinimumRequired:      minimumRequired,
		IsEligible:           isEligible,
	}, nil
}

// ============================================================================
// Attendance Edit & Audit Service Methods (Story 7.3)
// ============================================================================

// EditAttendance edits an attendance record with audit trail.
func (s *Service) EditAttendance(ctx context.Context, dto EditAttendanceDTO) (*EditAttendanceResult, error) {
	// Validate required fields
	if dto.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if dto.AttendanceID == uuid.Nil {
		return nil, ErrAttendanceIDRequired
	}
	if dto.Reason == "" {
		return nil, ErrEditReasonRequired
	}

	// Validate status if provided
	if dto.Status != "" {
		status := AttendanceStatus(dto.Status)
		if !status.IsValid() {
			return nil, ErrInvalidStatus
		}
	}

	// Get existing attendance record
	attendance, err := s.repo.GetAttendanceByID(ctx, dto.TenantID, dto.AttendanceID)
	if err != nil {
		return nil, err
	}

	// Check if user can edit this attendance
	canEdit, reason := s.CanUserEditAttendance(ctx, dto.TenantID, attendance, dto.EditedBy, dto.IsAdmin)
	if !canEdit {
		if reason == "window_expired" {
			return nil, ErrAdminOnlyEdit
		}
		return nil, ErrNotOriginalMarker
	}

	// Store previous values for audit
	previousStatus := attendance.Status
	previousRemarks := attendance.Remarks
	previousLateArrivalTime := attendance.LateArrivalTime

	// Update attendance fields
	now := time.Now()
	if dto.Status != "" {
		attendance.Status = models.StudentAttendanceStatus(dto.Status)
	}
	if dto.Remarks != nil {
		attendance.Remarks = *dto.Remarks
	}
	if dto.LateArrivalTime != nil {
		attendance.LateArrivalTime = dto.LateArrivalTime
	}
	attendance.UpdatedAt = now

	// Create audit record
	audit := &models.StudentAttendanceAudit{
		TenantID:                dto.TenantID,
		AttendanceID:            dto.AttendanceID,
		PreviousStatus:          &previousStatus,
		NewStatus:               attendance.Status,
		PreviousRemarks:         &previousRemarks,
		NewRemarks:              &attendance.Remarks,
		PreviousLateArrivalTime: previousLateArrivalTime,
		NewLateArrivalTime:      attendance.LateArrivalTime,
		ChangeType:              models.AttendanceChangeEdit,
		ChangeReason:            dto.Reason,
		ChangedBy:               dto.EditedBy,
		ChangedAt:               now,
	}

	// Save attendance and audit in transaction
	if err := s.repo.UpdateAttendanceWithAudit(ctx, attendance, audit); err != nil {
		return nil, err
	}

	return &EditAttendanceResult{
		AttendanceID: dto.AttendanceID.String(),
		StudentID:    attendance.StudentID.String(),
		Date:         attendance.AttendanceDate.Format("2006-01-02"),
		Status:       string(attendance.Status),
		EditedAt:     now.Format(time.RFC3339),
		EditedBy:     dto.EditedBy.String(),
		Message:      "Attendance updated successfully",
	}, nil
}

// CanUserEditAttendance checks if a specific user can edit an attendance record.
// Returns (canEdit, reason) where reason explains why editing is not allowed.
func (s *Service) CanUserEditAttendance(ctx context.Context, tenantID uuid.UUID, attendance *models.StudentAttendance, userID uuid.UUID, isAdmin bool) (bool, string) {
	// Admin can always edit
	if isAdmin {
		return true, ""
	}

	// Get settings for edit window
	editWindowMinutes := 120 // Default 2 hours

	section, err := s.repo.GetSectionByID(ctx, tenantID, attendance.SectionID)
	if err == nil && section.Class.ID != uuid.Nil {
		settings, err := s.repo.GetSettings(ctx, tenantID, tenantID)
		if err == nil {
			editWindowMinutes = settings.EditWindowMinutes
		}
	}

	// Check if within edit window
	windowEnd := attendance.MarkedAt.Add(time.Duration(editWindowMinutes) * time.Minute)
	isWithinWindow := time.Now().Before(windowEnd)

	if isWithinWindow {
		// Within window - original marker can edit
		if attendance.MarkedBy == userID {
			return true, ""
		}
		return false, "not_original_marker"
	}

	// Outside window - only admin can edit
	return false, "window_expired"
}

// GetAttendanceAuditTrail returns the audit history for an attendance record.
func (s *Service) GetAttendanceAuditTrail(ctx context.Context, tenantID, attendanceID uuid.UUID) (*AttendanceAuditTrailResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if attendanceID == uuid.Nil {
		return nil, ErrAttendanceIDRequired
	}

	// Get attendance record
	attendance, err := s.repo.GetAttendanceByID(ctx, tenantID, attendanceID)
	if err != nil {
		return nil, err
	}

	// Get audit trail
	audits, err := s.repo.GetAuditTrail(ctx, tenantID, attendanceID)
	if err != nil {
		return nil, err
	}

	// Build response
	auditResponses := make([]AttendanceAuditEntry, len(audits))
	for i, audit := range audits {
		entry := AttendanceAuditEntry{
			ID:           audit.ID.String(),
			ChangeType:   string(audit.ChangeType),
			NewStatus:    string(audit.NewStatus),
			ChangeReason: audit.ChangeReason,
			ChangedAt:    audit.ChangedAt.Format(time.RFC3339),
		}

		if audit.PreviousStatus != nil {
			entry.PreviousStatus = string(*audit.PreviousStatus)
		}
		if audit.PreviousRemarks != nil {
			entry.PreviousRemarks = *audit.PreviousRemarks
		}
		if audit.NewRemarks != nil {
			entry.NewRemarks = *audit.NewRemarks
		}

		if audit.ChangedByUser.ID != uuid.Nil {
			entry.ChangedByID = audit.ChangedByUser.ID.String()
			entry.ChangedByName = audit.ChangedByUser.FirstName + " " + audit.ChangedByUser.LastName
		}

		auditResponses[i] = entry
	}

	studentName := ""
	if attendance.Student.ID != uuid.Nil {
		studentName = attendance.Student.FullName()
	}

	return &AttendanceAuditTrailResponse{
		AttendanceID: attendanceID.String(),
		StudentID:    attendance.StudentID.String(),
		StudentName:  studentName,
		Date:         attendance.AttendanceDate.Format("2006-01-02"),
		AuditEntries: auditResponses,
		TotalChanges: len(audits),
	}, nil
}

// GetEditWindowStatus returns the edit window status for an attendance record.
func (s *Service) GetEditWindowStatus(ctx context.Context, tenantID, attendanceID, userID uuid.UUID, isAdmin bool) (*EditWindowStatusResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if attendanceID == uuid.Nil {
		return nil, ErrAttendanceIDRequired
	}

	// Get attendance record
	attendance, err := s.repo.GetAttendanceByID(ctx, tenantID, attendanceID)
	if err != nil {
		return nil, err
	}

	// Get edit window settings
	editWindowMinutes := 120

	section, err := s.repo.GetSectionByID(ctx, tenantID, attendance.SectionID)
	if err == nil && section.Class.ID != uuid.Nil {
		settings, err := s.repo.GetSettings(ctx, tenantID, tenantID)
		if err == nil {
			editWindowMinutes = settings.EditWindowMinutes
		}
	}

	// Calculate window end
	windowEnd := attendance.MarkedAt.Add(time.Duration(editWindowMinutes) * time.Minute)
	isWithinWindow := time.Now().Before(windowEnd)

	remainingMinutes := 0
	if isWithinWindow {
		remaining := time.Until(windowEnd)
		remainingMinutes = int(remaining.Minutes())
	}

	// Check if user can edit
	canEdit, reason := s.CanUserEditAttendance(ctx, tenantID, attendance, userID, isAdmin)

	return &EditWindowStatusResponse{
		AttendanceID:      attendanceID.String(),
		MarkedAt:          attendance.MarkedAt.Format(time.RFC3339),
		WindowEndAt:       windowEnd.Format(time.RFC3339),
		WindowMinutes:     editWindowMinutes,
		RemainingMinutes:  remainingMinutes,
		IsWithinWindow:    isWithinWindow,
		CanEdit:           canEdit,
		EditDeniedReason:  reason,
		IsOriginalMarker:  attendance.MarkedBy == userID,
		RequiresAdminEdit: !isWithinWindow && !isAdmin,
	}, nil
}

// ============================================================================
// Calendar & Reports Service Methods (Stories 7.4-7.8)
// ============================================================================

// GetStudentCalendar returns a student's monthly attendance calendar.
func (s *Service) GetStudentCalendar(ctx context.Context, tenantID, studentID uuid.UUID, year, month int) (*MonthlyCalendarResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if studentID == uuid.Nil {
		return nil, ErrStudentIDRequired
	}

	// Get student details
	student, err := s.repo.GetStudentByID(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}

	// Get attendance for the month
	attendance, err := s.repo.GetStudentMonthlyAttendance(ctx, tenantID, studentID, year, month)
	if err != nil {
		return nil, err
	}

	// Build attendance map
	attendanceMap := make(map[string]models.StudentAttendance)
	for _, a := range attendance {
		attendanceMap[a.AttendanceDate.Format("2006-01-02")] = a
	}

	// Generate calendar days
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	days := make([]CalendarDayResponse, 0)
	summary := MonthlySummary{}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		day := CalendarDayResponse{
			Date:      d.Format("2006-01-02"),
			DayOfWeek: int(d.Weekday()),
			IsWeekend: d.Weekday() == time.Sunday || d.Weekday() == time.Saturday,
		}

		if day.IsWeekend {
			summary.Holidays++
		} else {
			summary.WorkingDays++
			if a, ok := attendanceMap[d.Format("2006-01-02")]; ok {
				day.Status = string(a.Status)
				day.Remarks = a.Remarks
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
		}
		days = append(days, day)
	}

	// Calculate percentage
	if summary.WorkingDays > 0 {
		presentCount := summary.Present + summary.Late + summary.HalfDay
		summary.Percentage = float64(presentCount) / float64(summary.WorkingDays) * 100
	}

	// Get student's current section
	section, _ := s.repo.GetStudentCurrentSection(ctx, tenantID, studentID)

	// Get class average
	classAverage := 0.0
	if section != nil {
		classAverage, _ = s.repo.GetSectionAttendanceStats(ctx, tenantID, section.ID, startDate, endDate)
	}

	// Determine trend (compare with previous month)
	trend := "stable"
	prevMonth := startDate.AddDate(0, -1, 0)
	prevAttendance, _ := s.repo.GetStudentMonthlyAttendance(ctx, tenantID, studentID, prevMonth.Year(), int(prevMonth.Month()))
	if len(prevAttendance) > 0 {
		prevPresent := 0
		for _, a := range prevAttendance {
			if a.Status == models.StudentAttendancePresent || a.Status == models.StudentAttendanceLate || a.Status == models.StudentAttendanceHalfDay {
				prevPresent++
			}
		}
		prevPercentage := float64(prevPresent) / float64(len(prevAttendance)) * 100
		if summary.Percentage > prevPercentage+5 {
			trend = "improving"
		} else if summary.Percentage < prevPercentage-5 {
			trend = "declining"
		}
	}

	monthNames := []string{"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	return &MonthlyCalendarResponse{
		StudentID:    studentID.String(),
		StudentName:  student.FullName(),
		Year:         year,
		Month:        month,
		MonthName:    monthNames[month],
		Days:         days,
		Summary:      summary,
		ClassAverage: classAverage,
		Trend:        trend,
	}, nil
}

// GetStudentSummary returns a student's attendance summary with trend.
func (s *Service) GetStudentSummary(ctx context.Context, tenantID, studentID uuid.UUID, dateFrom, dateTo time.Time) (*StudentSummaryResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if studentID == uuid.Nil {
		return nil, ErrStudentIDRequired
	}

	// Get student details
	student, err := s.repo.GetStudentByID(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}

	// Get attendance for date range
	attendance, err := s.repo.GetStudentAttendanceRange(ctx, tenantID, studentID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	// Calculate summary
	summary := MonthlySummary{WorkingDays: len(attendance)}
	for _, a := range attendance {
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

	if summary.WorkingDays > 0 {
		presentCount := summary.Present + summary.Late + summary.HalfDay
		summary.Percentage = float64(presentCount) / float64(summary.WorkingDays) * 100
	}

	// Get student's current section
	studentSection, _ := s.repo.GetStudentCurrentSection(ctx, tenantID, studentID)

	// Get class average
	classAverage := 0.0
	if studentSection != nil {
		classAverage, _ = s.repo.GetSectionAttendanceStats(ctx, tenantID, studentSection.ID, dateFrom, dateTo)
	}

	// Calculate trend
	trend := "stable"
	midPoint := dateFrom.Add(dateTo.Sub(dateFrom) / 2)
	firstHalf, _ := s.repo.GetStudentAttendanceRange(ctx, tenantID, studentID, dateFrom, midPoint)
	secondHalf, _ := s.repo.GetStudentAttendanceRange(ctx, tenantID, studentID, midPoint.AddDate(0, 0, 1), dateTo)

	if len(firstHalf) > 0 && len(secondHalf) > 0 {
		firstPresent := 0
		for _, a := range firstHalf {
			if a.Status != models.StudentAttendanceAbsent {
				firstPresent++
			}
		}
		secondPresent := 0
		for _, a := range secondHalf {
			if a.Status != models.StudentAttendanceAbsent {
				secondPresent++
			}
		}
		firstPct := float64(firstPresent) / float64(len(firstHalf)) * 100
		secondPct := float64(secondPresent) / float64(len(secondHalf)) * 100
		if secondPct > firstPct+5 {
			trend = "improving"
		} else if secondPct < firstPct-5 {
			trend = "declining"
		}
	}

	sectionName := ""
	className := ""
	sectionID := ""
	if studentSection != nil {
		sectionID = studentSection.ID.String()
		sectionName = studentSection.Name
		if studentSection.Class.ID != uuid.Nil {
			className = studentSection.Class.Name
		}
	}

	return &StudentSummaryResponse{
		StudentID:    studentID.String(),
		StudentName:  student.FullName(),
		SectionID:    sectionID,
		SectionName:  sectionName + " - " + className,
		DateFrom:     dateFrom.Format("2006-01-02"),
		DateTo:       dateTo.Format("2006-01-02"),
		Summary:      summary,
		ClassAverage: classAverage,
		Trend:        trend,
	}, nil
}

// GetClassReport returns class-level attendance report for a date.
func (s *Service) GetClassReport(ctx context.Context, tenantID, sectionID uuid.UUID, date time.Time) (*ClassReportResponse, error) {
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

	// Get attendance for the date
	attendance, err := s.repo.GetClassAttendance(ctx, tenantID, sectionID, date)
	if err != nil {
		return nil, err
	}

	// Build attendance map
	attendanceMap := make(map[uuid.UUID]models.StudentAttendance)
	for _, a := range attendance {
		attendanceMap[a.StudentID] = a
	}

	// Build student entries
	entries := make([]StudentReportEntry, len(students))
	summary := AttendanceSummary{Total: len(students)}

	for i, student := range students {
		entry := StudentReportEntry{
			StudentID:       student.ID.String(),
			AdmissionNumber: student.AdmissionNumber,
			FullName:        student.FullName(),
		}

		if a, ok := attendanceMap[student.ID]; ok {
			entry.Status = string(a.Status)
			entry.StatusLabel = AttendanceStatus(a.Status).Label()
			entry.Remarks = a.Remarks

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
		entries[i] = entry
	}

	attendanceRate := 0.0
	if summary.Total > 0 {
		attendanceRate = float64(summary.Present+summary.Late+summary.HalfDay) / float64(summary.Total) * 100
	}

	className := ""
	if section.Class.ID != uuid.Nil {
		className = section.Class.Name
	}

	return &ClassReportResponse{
		SectionID:      sectionID.String(),
		SectionName:    section.Name,
		ClassName:      className,
		Date:           date.Format("2006-01-02"),
		Students:       entries,
		Summary:        summary,
		AttendanceRate: attendanceRate,
	}, nil
}

// GetMonthlyClassReport returns monthly class attendance report in grid format.
func (s *Service) GetMonthlyClassReport(ctx context.Context, tenantID, sectionID uuid.UUID, year, month int) (*MonthlyClassReportResponse, error) {
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

	// Get attendance for the month
	attendance, err := s.repo.GetClassMonthlyAttendance(ctx, tenantID, sectionID, year, month)
	if err != nil {
		return nil, err
	}

	// Build attendance map: studentID -> date -> attendance
	attendanceMap := make(map[uuid.UUID]map[string]models.StudentAttendance)
	for _, a := range attendance {
		if _, ok := attendanceMap[a.StudentID]; !ok {
			attendanceMap[a.StudentID] = make(map[string]models.StudentAttendance)
		}
		attendanceMap[a.StudentID][a.AttendanceDate.Format("2006-01-02")] = a
	}

	// Generate working dates (excluding weekends)
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	var dates []string
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Sunday && d.Weekday() != time.Saturday {
			dates = append(dates, d.Format("2006-01-02"))
		}
	}

	// Build student reports
	studentReports := make([]MonthlyStudentReport, len(students))
	summaryStats := ClassMonthlySummary{TotalStudents: len(students)}
	totalAttendance := 0.0

	for i, student := range students {
		report := MonthlyStudentReport{
			StudentID:       student.ID.String(),
			AdmissionNumber: student.AdmissionNumber,
			FullName:        student.FullName(),
			DailyStatus:     make(map[string]string),
		}

		studentAttendance := attendanceMap[student.ID]
		for _, date := range dates {
			if a, ok := studentAttendance[date]; ok {
				report.DailyStatus[date] = string(a.Status)
				switch a.Status {
				case models.StudentAttendancePresent:
					report.Present++
				case models.StudentAttendanceAbsent:
					report.Absent++
				case models.StudentAttendanceLate:
					report.Late++
					report.Present++ // Late counts as present
				case models.StudentAttendanceHalfDay:
					report.Present++ // Half-day counts as present
				}
			}
		}

		if len(dates) > 0 {
			report.Percentage = float64(report.Present) / float64(len(dates)) * 100
		}
		totalAttendance += report.Percentage

		if report.Percentage >= 90 {
			summaryStats.StudentsAbove90++
		}
		if report.Percentage < 75 {
			summaryStats.StudentsBelow75++
		}
		if report.Percentage < 60 {
			summaryStats.StudentsBelow60++
		}

		studentReports[i] = report
	}

	if summaryStats.TotalStudents > 0 {
		summaryStats.AverageAttendance = totalAttendance / float64(summaryStats.TotalStudents)
	}

	className := ""
	if section.Class.ID != uuid.Nil {
		className = section.Class.Name
	}

	monthNames := []string{"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	return &MonthlyClassReportResponse{
		SectionID:   sectionID.String(),
		SectionName: section.Name,
		ClassName:   className,
		Year:        year,
		Month:       month,
		MonthName:   monthNames[month],
		WorkingDays: len(dates),
		Dates:       dates,
		Students:    studentReports,
		Summary:     summaryStats,
	}, nil
}

// GetLowAttendanceDashboard returns dashboard data for low attendance students.
func (s *Service) GetLowAttendanceDashboard(ctx context.Context, tenantID uuid.UUID, dateFrom, dateTo time.Time, threshold, criticalThreshold float64) (*LowAttendanceDashboardResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	// Get all students with their attendance
	lowAttendanceData, err := s.repo.GetLowAttendanceStudents(ctx, tenantID, dateFrom, dateTo, threshold)
	if err != nil {
		return nil, err
	}

	// Get all sections for total student count
	sections, err := s.repo.GetAllActiveSections(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	totalStudents := 0
	for _, section := range sections {
		count, _ := s.repo.CountStudentsInSection(ctx, tenantID, section.ID)
		totalStudents += int(count)
	}

	// Build low attendance students list
	students := make([]LowAttendanceStudent, 0)
	chronicCount := 0

	for _, data := range lowAttendanceData {
		percentage := 0.0
		if data.TotalRecords > 0 {
			percentage = float64(data.TotalPresent) / float64(data.TotalRecords) * 100
		}

		student, err := s.repo.GetStudentByID(ctx, tenantID, data.StudentID)
		if err != nil {
			continue
		}

		// Get student's section
		studentSection, _ := s.repo.GetStudentCurrentSection(ctx, tenantID, data.StudentID)

		className := ""
		sectionName := ""
		if studentSection != nil {
			sectionName = studentSection.Name
			if studentSection.Class.ID != uuid.Nil {
				className = studentSection.Class.Name
			}
		}

		entry := LowAttendanceStudent{
			StudentID:       data.StudentID.String(),
			AdmissionNumber: student.AdmissionNumber,
			FullName:        student.FullName(),
			ClassName:       className,
			SectionName:     sectionName,
			AttendanceRate:  percentage,
			DaysAbsent:      int(data.TotalRecords - data.TotalPresent),
		}
		students = append(students, entry)

		if percentage < criticalThreshold {
			chronicCount++
		}
	}

	// Get overall attendance rate
	present, _, total, _ := s.repo.GetDailyAttendanceSummary(ctx, tenantID, time.Now())
	overallRate := 0.0
	if total > 0 {
		overallRate = float64(present) / float64(total) * 100
	}

	return &LowAttendanceDashboardResponse{
		DateFrom:              dateFrom.Format("2006-01-02"),
		DateTo:                dateTo.Format("2006-01-02"),
		Threshold:             threshold,
		CriticalThreshold:     criticalThreshold,
		TotalStudents:         totalStudents,
		BelowThreshold:        len(students),
		ChronicAbsentees:      chronicCount,
		OverallAttendanceRate: overallRate,
		Students:              students,
	}, nil
}

// GetUnmarkedAttendance returns classes with unmarked attendance for a date.
func (s *Service) GetUnmarkedAttendance(ctx context.Context, tenantID uuid.UUID, date time.Time, deadlineTime string) (*UnmarkedAttendanceResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	// Get unmarked sections
	unmarkedIDs, err := s.repo.GetUnmarkedSections(ctx, tenantID, date)
	if err != nil {
		return nil, err
	}

	// Get all sections for total count
	allSections, err := s.repo.GetAllActiveSections(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Check if past deadline
	now := time.Now()
	deadline, _ := time.Parse("15:04", deadlineTime)
	deadlineToday := time.Date(now.Year(), now.Month(), now.Day(), deadline.Hour(), deadline.Minute(), 0, 0, now.Location())
	isPastDeadline := now.After(deadlineToday)

	// Build unmarked class info
	unmarkedClasses := make([]UnmarkedClassInfo, 0)
	for _, sectionID := range unmarkedIDs {
		section, err := s.repo.GetSectionByID(ctx, tenantID, sectionID)
		if err != nil {
			continue
		}

		studentCount, _ := s.repo.CountStudentsInSection(ctx, tenantID, sectionID)

		className := ""
		if section.Class.ID != uuid.Nil {
			className = section.Class.Name
		}

		unmarkedClasses = append(unmarkedClasses, UnmarkedClassInfo{
			SectionID:    sectionID.String(),
			SectionName:  section.Name,
			ClassName:    className,
			StudentCount: int(studentCount),
			IsEscalated:  isPastDeadline,
		})
	}

	return &UnmarkedAttendanceResponse{
		Date:            date.Format("2006-01-02"),
		Deadline:        deadlineTime,
		IsPostDeadline:  isPastDeadline,
		UnmarkedClasses: unmarkedClasses,
		TotalClasses:    len(allSections),
		MarkedClasses:   len(allSections) - len(unmarkedIDs),
	}, nil
}

// GetDailyReport returns end-of-day attendance summary report.
func (s *Service) GetDailyReport(ctx context.Context, tenantID uuid.UUID, date time.Time) (*DailyReportSummaryResponse, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	// Get overall attendance stats
	present, absent, total, err := s.repo.GetDailyAttendanceSummary(ctx, tenantID, date)
	if err != nil {
		return nil, err
	}

	overallRate := 0.0
	if total > 0 {
		overallRate = float64(present) / float64(total) * 100
	}

	// Get all sections
	sections, err := s.repo.GetAllActiveSections(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Get unmarked sections
	unmarkedIDs, _ := s.repo.GetUnmarkedSections(ctx, tenantID, date)
	unmarkedMap := make(map[uuid.UUID]bool)
	for _, id := range unmarkedIDs {
		unmarkedMap[id] = true
	}

	// Get low attendance classes (below 75%)
	lowClasses := make([]ClassAttendanceBreakdown, 0)
	for _, section := range sections {
		if unmarkedMap[section.ID] {
			continue
		}

		rate, _ := s.repo.GetSectionAttendanceStats(ctx, tenantID, section.ID, date, date)
		studentCount, _ := s.repo.CountStudentsInSection(ctx, tenantID, section.ID)

		className := ""
		if section.Class.ID != uuid.Nil {
			className = section.Class.Name
		}

		if rate < 75 {
			lowClasses = append(lowClasses, ClassAttendanceBreakdown{
				ClassName:      className,
				SectionName:    section.Name,
				SectionID:      section.ID.String(),
				TotalStudents:  int(studentCount),
				AttendanceRate: rate,
			})
		}
	}

	totalStudents := 0
	for _, section := range sections {
		count, _ := s.repo.CountStudentsInSection(ctx, tenantID, section.ID)
		totalStudents += int(count)
	}

	return &DailyReportSummaryResponse{
		Date:                  date.Format("2006-01-02"),
		OverallAttendanceRate: overallRate,
		TotalStudents:         totalStudents,
		TotalPresent:          int(present),
		TotalAbsent:           int(absent),
		ClassesMarked:         len(sections) - len(unmarkedIDs),
		ClassesTotal:          len(sections),
		LowAttendanceClasses:  lowClasses,
		GeneratedAt:           time.Now().Format(time.RFC3339),
	}, nil
}
