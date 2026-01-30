package examination

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Service handles business logic for examinations
type Service struct {
	repo *Repository
}

// NewService creates a new examination service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ========================================
// Examination Methods
// ========================================

// List returns examinations based on filters
func (s *Service) List(tenantID uuid.UUID, filter ExaminationFilter) ([]models.Examination, error) {
	return s.repo.List(tenantID, filter)
}

// GetByID returns an examination by ID
func (s *Service) GetByID(tenantID, id uuid.UUID) (*models.Examination, error) {
	exam, err := s.repo.GetByID(tenantID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrExaminationNotFound
		}
		return nil, err
	}
	return exam, nil
}

// Create creates a new examination
func (s *Service) Create(tenantID uuid.UUID, req CreateExaminationRequest, userID uuid.UUID) (*models.Examination, error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, ErrInvalidDateRange
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, ErrInvalidDateRange
	}

	// Validate date range
	if endDate.Before(startDate) {
		return nil, ErrInvalidDateRange
	}

	// Validate exam type exists and is active
	examType, err := s.repo.GetExamTypeByID(tenantID, req.ExamTypeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrExamTypeNotActive
		}
		return nil, err
	}
	if !examType.IsActive {
		return nil, ErrExamTypeNotActive
	}

	// Validate classes
	if len(req.ClassIDs) == 0 {
		return nil, ErrNoClassesSpecified
	}

	for _, classID := range req.ClassIDs {
		exists, err := s.repo.ClassExists(tenantID, classID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrNoClassesSpecified
		}
	}

	// Create examination
	exam := &models.Examination{
		TenantID:       tenantID,
		Name:           req.Name,
		ExamTypeID:     req.ExamTypeID,
		AcademicYearID: req.AcademicYearID,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         models.ExamStatusDraft,
		Description:    req.Description,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := s.repo.Create(exam); err != nil {
		return nil, err
	}

	// Set classes
	if err := s.repo.SetClasses(exam.ID, req.ClassIDs); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.repo.GetByID(tenantID, exam.ID)
}

// Update updates an examination
func (s *Service) Update(tenantID, id uuid.UUID, req UpdateExaminationRequest, userID uuid.UUID) (*models.Examination, error) {
	exam, err := s.GetByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	// Can only update draft examinations
	if exam.Status != models.ExamStatusDraft {
		return nil, ErrCannotUpdatePublished
	}

	// Update fields
	if req.Name != nil {
		exam.Name = *req.Name
	}

	if req.ExamTypeID != nil {
		examType, err := s.repo.GetExamTypeByID(tenantID, *req.ExamTypeID)
		if err != nil {
			return nil, err
		}
		if !examType.IsActive {
			return nil, ErrExamTypeNotActive
		}
		exam.ExamTypeID = *req.ExamTypeID
	}

	if req.AcademicYearID != nil {
		exam.AcademicYearID = *req.AcademicYearID
	}

	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, ErrInvalidDateRange
		}
		exam.StartDate = startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, ErrInvalidDateRange
		}
		exam.EndDate = endDate
	}

	// Validate date range
	if exam.EndDate.Before(exam.StartDate) {
		return nil, ErrInvalidDateRange
	}

	if req.Description != nil {
		exam.Description = req.Description
	}

	exam.UpdatedBy = &userID

	if err := s.repo.Update(exam); err != nil {
		return nil, err
	}

	// Update classes if provided
	if req.ClassIDs != nil {
		if len(req.ClassIDs) == 0 {
			return nil, ErrNoClassesSpecified
		}

		for _, classID := range req.ClassIDs {
			exists, err := s.repo.ClassExists(tenantID, classID)
			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, ErrNoClassesSpecified
			}
		}

		if err := s.repo.SetClasses(exam.ID, req.ClassIDs); err != nil {
			return nil, err
		}
	}

	return s.repo.GetByID(tenantID, id)
}

// Delete deletes an examination
func (s *Service) Delete(tenantID, id uuid.UUID) error {
	exam, err := s.GetByID(tenantID, id)
	if err != nil {
		return err
	}

	// Can only delete draft examinations
	if exam.Status != models.ExamStatusDraft {
		return ErrCannotDeletePublished
	}

	return s.repo.Delete(tenantID, id)
}

// Publish publishes an examination (changes status to scheduled)
func (s *Service) Publish(tenantID, id uuid.UUID, userID uuid.UUID) (*models.Examination, error) {
	exam, err := s.GetByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check status transition
	if !exam.Status.CanTransitionTo(models.ExamStatusScheduled) {
		return nil, ErrInvalidStatusTransition
	}

	// Must have at least one schedule
	count, err := s.repo.CountSchedules(id)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrNoSchedulesForPublish
	}

	// Update status
	if err := s.repo.UpdateStatus(tenantID, id, models.ExamStatusScheduled); err != nil {
		return nil, err
	}

	return s.repo.GetByID(tenantID, id)
}

// UnPublish reverts examination to draft status
func (s *Service) UnPublish(tenantID, id uuid.UUID, userID uuid.UUID) (*models.Examination, error) {
	exam, err := s.GetByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check status transition
	if !exam.Status.CanTransitionTo(models.ExamStatusDraft) {
		return nil, ErrInvalidStatusTransition
	}

	// Update status
	if err := s.repo.UpdateStatus(tenantID, id, models.ExamStatusDraft); err != nil {
		return nil, err
	}

	return s.repo.GetByID(tenantID, id)
}

// ========================================
// Schedule Methods
// ========================================

// CreateSchedule creates an exam schedule
func (s *Service) CreateSchedule(tenantID, examID uuid.UUID, req CreateScheduleRequest) (*models.ExamSchedule, error) {
	// Get examination
	exam, err := s.GetByID(tenantID, examID)
	if err != nil {
		return nil, err
	}

	// Can only add schedules to draft examinations
	if exam.Status != models.ExamStatusDraft {
		return nil, ErrCannotUpdatePublished
	}

	// Parse exam date
	examDate, err := time.Parse("2006-01-02", req.ExamDate)
	if err != nil {
		return nil, ErrScheduleOutsideExamDates
	}

	// Validate date is within examination range
	if examDate.Before(exam.StartDate) || examDate.After(exam.EndDate) {
		return nil, ErrScheduleOutsideExamDates
	}

	// Validate times
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, ErrInvalidTimeRange
	}
	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, ErrInvalidTimeRange
	}
	if !endTime.After(startTime) {
		return nil, ErrInvalidTimeRange
	}

	// Validate passing marks
	if req.PassingMarks != nil && *req.PassingMarks > req.MaxMarks {
		return nil, ErrInvalidPassingMarks
	}

	// Check if subject already scheduled
	exists, err := s.repo.CheckSubjectExists(examID, req.SubjectID, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSubjectAlreadyScheduled
	}

	// Create schedule
	schedule := &models.ExamSchedule{
		ExaminationID: examID,
		SubjectID:     req.SubjectID,
		ExamDate:      examDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		MaxMarks:      req.MaxMarks,
		PassingMarks:  req.PassingMarks,
		Venue:         req.Venue,
		Notes:         req.Notes,
	}

	if err := s.repo.CreateSchedule(schedule); err != nil {
		return nil, err
	}

	return s.repo.GetScheduleByID(schedule.ID)
}

// UpdateSchedule updates an exam schedule
func (s *Service) UpdateSchedule(tenantID, examID, scheduleID uuid.UUID, req UpdateScheduleRequest) (*models.ExamSchedule, error) {
	// Get examination
	exam, err := s.GetByID(tenantID, examID)
	if err != nil {
		return nil, err
	}

	// Can only update schedules for draft examinations
	if exam.Status != models.ExamStatusDraft {
		return nil, ErrCannotUpdatePublished
	}

	// Get schedule
	schedule, err := s.repo.GetScheduleByID(scheduleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}

	// Verify schedule belongs to this examination
	if schedule.ExaminationID != examID {
		return nil, ErrScheduleNotFound
	}

	// Update fields
	if req.SubjectID != nil {
		// Check if new subject already scheduled
		exists, err := s.repo.CheckSubjectExists(examID, *req.SubjectID, &scheduleID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrSubjectAlreadyScheduled
		}
		schedule.SubjectID = *req.SubjectID
	}

	if req.ExamDate != nil {
		examDate, err := time.Parse("2006-01-02", *req.ExamDate)
		if err != nil {
			return nil, ErrScheduleOutsideExamDates
		}
		if examDate.Before(exam.StartDate) || examDate.After(exam.EndDate) {
			return nil, ErrScheduleOutsideExamDates
		}
		schedule.ExamDate = examDate
	}

	if req.StartTime != nil {
		schedule.StartTime = *req.StartTime
	}

	if req.EndTime != nil {
		schedule.EndTime = *req.EndTime
	}

	// Validate times
	startTime, _ := time.Parse("15:04", schedule.StartTime)
	endTime, _ := time.Parse("15:04", schedule.EndTime)
	if !endTime.After(startTime) {
		return nil, ErrInvalidTimeRange
	}

	if req.MaxMarks != nil {
		schedule.MaxMarks = *req.MaxMarks
	}

	if req.PassingMarks != nil {
		if *req.PassingMarks > schedule.MaxMarks {
			return nil, ErrInvalidPassingMarks
		}
		schedule.PassingMarks = req.PassingMarks
	}

	if req.Venue != nil {
		schedule.Venue = req.Venue
	}

	if req.Notes != nil {
		schedule.Notes = req.Notes
	}

	if err := s.repo.UpdateSchedule(schedule); err != nil {
		return nil, err
	}

	return s.repo.GetScheduleByID(scheduleID)
}

// DeleteSchedule deletes an exam schedule
func (s *Service) DeleteSchedule(tenantID, examID, scheduleID uuid.UUID) error {
	// Get examination
	exam, err := s.GetByID(tenantID, examID)
	if err != nil {
		return err
	}

	// Can only delete schedules from draft examinations
	if exam.Status != models.ExamStatusDraft {
		return ErrCannotUpdatePublished
	}

	// Get schedule to verify it belongs to this examination
	schedule, err := s.repo.GetScheduleByID(scheduleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrScheduleNotFound
		}
		return err
	}

	if schedule.ExaminationID != examID {
		return ErrScheduleNotFound
	}

	return s.repo.DeleteSchedule(scheduleID)
}

// GetSchedules returns all schedules for an examination
func (s *Service) GetSchedules(tenantID, examID uuid.UUID) ([]models.ExamSchedule, error) {
	// Verify examination exists
	if _, err := s.GetByID(tenantID, examID); err != nil {
		return nil, err
	}

	return s.repo.GetSchedulesByExamID(examID)
}
