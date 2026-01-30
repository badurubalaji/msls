package examination

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for examinations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new examination repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ========================================
// Examination Methods
// ========================================

// List returns examinations based on filters
func (r *Repository) List(tenantID uuid.UUID, filter ExaminationFilter) ([]models.Examination, error) {
	var exams []models.Examination

	query := r.db.Where("tenant_id = ?", tenantID).
		Preload("ExamType").
		Preload("AcademicYear").
		Preload("Classes").
		Preload("Schedules").
		Preload("Schedules.Subject")

	if filter.AcademicYearID != nil {
		query = query.Where("academic_year_id = ?", *filter.AcademicYearID)
	}

	if filter.ExamTypeID != nil {
		query = query.Where("exam_type_id = ?", *filter.ExamTypeID)
	}

	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ?", searchTerm)
	}

	if filter.ClassID != nil {
		query = query.Joins("JOIN examination_classes ec ON ec.examination_id = examinations.id").
			Where("ec.class_id = ?", *filter.ClassID)
	}

	err := query.Order("start_date DESC, name ASC").Find(&exams).Error
	return exams, err
}

// GetByID returns an examination by ID
func (r *Repository) GetByID(tenantID, id uuid.UUID) (*models.Examination, error) {
	var exam models.Examination
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("ExamType").
		Preload("AcademicYear").
		Preload("Classes").
		Preload("Schedules", func(db *gorm.DB) *gorm.DB {
			return db.Order("exam_date ASC, start_time ASC")
		}).
		Preload("Schedules.Subject").
		First(&exam).Error

	if err != nil {
		return nil, err
	}
	return &exam, nil
}

// Create creates a new examination
func (r *Repository) Create(exam *models.Examination) error {
	return r.db.Create(exam).Error
}

// Update updates an examination
func (r *Repository) Update(exam *models.Examination) error {
	return r.db.Save(exam).Error
}

// Delete deletes an examination
func (r *Repository) Delete(tenantID, id uuid.UUID) error {
	return r.db.Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&models.Examination{}).Error
}

// UpdateStatus updates the examination status
func (r *Repository) UpdateStatus(tenantID, id uuid.UUID, status models.ExamStatus) error {
	return r.db.Model(&models.Examination{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Update("status", status).Error
}

// ========================================
// Examination Classes Methods
// ========================================

// SetClasses sets the classes for an examination (replaces existing)
func (r *Repository) SetClasses(examID uuid.UUID, classIDs []uuid.UUID) error {
	// Start transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing classes
		if err := tx.Where("examination_id = ?", examID).Delete(&models.ExaminationClass{}).Error; err != nil {
			return err
		}

		// Insert new classes
		for _, classID := range classIDs {
			ec := models.ExaminationClass{
				ExaminationID: examID,
				ClassID:       classID,
			}
			if err := tx.Create(&ec).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ========================================
// Schedule Methods
// ========================================

// CreateSchedule creates an exam schedule
func (r *Repository) CreateSchedule(schedule *models.ExamSchedule) error {
	return r.db.Create(schedule).Error
}

// GetScheduleByID returns a schedule by ID
func (r *Repository) GetScheduleByID(scheduleID uuid.UUID) (*models.ExamSchedule, error) {
	var schedule models.ExamSchedule
	err := r.db.Where("id = ?", scheduleID).
		Preload("Subject").
		First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// UpdateSchedule updates an exam schedule
func (r *Repository) UpdateSchedule(schedule *models.ExamSchedule) error {
	return r.db.Save(schedule).Error
}

// DeleteSchedule deletes an exam schedule
func (r *Repository) DeleteSchedule(scheduleID uuid.UUID) error {
	return r.db.Where("id = ?", scheduleID).Delete(&models.ExamSchedule{}).Error
}

// CheckSubjectExists checks if a subject is already scheduled for an examination
func (r *Repository) CheckSubjectExists(examID, subjectID uuid.UUID, excludeScheduleID *uuid.UUID) (bool, error) {
	query := r.db.Model(&models.ExamSchedule{}).
		Where("examination_id = ? AND subject_id = ?", examID, subjectID)

	if excludeScheduleID != nil {
		query = query.Where("id != ?", *excludeScheduleID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// GetSchedulesByExamID returns all schedules for an examination
func (r *Repository) GetSchedulesByExamID(examID uuid.UUID) ([]models.ExamSchedule, error) {
	var schedules []models.ExamSchedule
	err := r.db.Where("examination_id = ?", examID).
		Preload("Subject").
		Order("exam_date ASC, start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

// CountSchedules counts schedules for an examination
func (r *Repository) CountSchedules(examID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.ExamSchedule{}).Where("examination_id = ?", examID).Count(&count).Error
	return count, err
}

// ========================================
// Validation Helpers
// ========================================

// GetExamTypeByID returns an exam type by ID
func (r *Repository) GetExamTypeByID(tenantID, examTypeID uuid.UUID) (*models.ExamType, error) {
	var examType models.ExamType
	err := r.db.Where("tenant_id = ? AND id = ?", tenantID, examTypeID).First(&examType).Error
	if err != nil {
		return nil, err
	}
	return &examType, nil
}

// ClassExists checks if a class exists
func (r *Repository) ClassExists(tenantID, classID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Class{}).Where("tenant_id = ? AND id = ?", tenantID, classID).Count(&count).Error
	return count > 0, err
}

// SubjectExists checks if a subject exists
func (r *Repository) SubjectExists(tenantID, subjectID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Subject{}).Where("tenant_id = ? AND id = ?", tenantID, subjectID).Count(&count).Error
	return count > 0, err
}
