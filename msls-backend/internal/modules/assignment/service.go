// Package assignment provides teacher subject assignment functionality.
package assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Default workload settings if not configured.
const (
	defaultMinPeriodsPerWeek = 20
	defaultMaxPeriodsPerWeek = 35
)

// Service handles assignment business logic.
type Service struct {
	repo *Repository
	db   *gorm.DB
}

// NewService creates a new assignment service.
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: NewRepository(db),
		db:   db,
	}
}

// Create creates a new teacher assignment with validation.
func (s *Service) Create(ctx context.Context, dto CreateAssignmentDTO) (*models.TeacherSubjectAssignment, error) {
	// Validate effective dates
	if dto.EffectiveTo != nil && dto.EffectiveFrom.After(*dto.EffectiveTo) {
		return nil, ErrInvalidEffectiveDate
	}

	// Verify staff exists and is teaching type
	var staff models.Staff
	if err := s.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", dto.TenantID, dto.StaffID).First(&staff).Error; err != nil {
		return nil, ErrStaffNotFound
	}
	if staff.StaffType != models.StaffTypeTeaching {
		return nil, fmt.Errorf("staff member is not of teaching type")
	}

	// Verify subject exists
	var subject models.Subject
	if err := s.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", dto.TenantID, dto.SubjectID).First(&subject).Error; err != nil {
		return nil, ErrSubjectNotFound
	}

	// Verify class exists
	var class models.Class
	if err := s.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", dto.TenantID, dto.ClassID).First(&class).Error; err != nil {
		return nil, ErrClassNotFound
	}

	// Verify section exists if provided
	if dto.SectionID != nil {
		var section models.Section
		if err := s.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", dto.TenantID, *dto.SectionID).First(&section).Error; err != nil {
			return nil, ErrSectionNotFound
		}
	}

	// Verify academic year exists
	var academicYear models.AcademicYear
	if err := s.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", dto.TenantID, dto.AcademicYearID).First(&academicYear).Error; err != nil {
		return nil, ErrAcademicYearNotFound
	}

	var assignment *models.TeacherSubjectAssignment
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check workload if periods are being added
		if dto.PeriodsPerWeek > 0 {
			currentPeriods, err := s.repo.GetTotalPeriodsForStaff(ctx, dto.TenantID, dto.StaffID, dto.AcademicYearID)
			if err != nil {
				return err
			}

			maxPeriods, err := s.getMaxPeriodsForBranch(ctx, dto.TenantID, staff.BranchID)
			if err != nil {
				return err
			}

			if currentPeriods+dto.PeriodsPerWeek > maxPeriods {
				return ErrTeacherOverAssigned
			}
		}

		// Handle class teacher assignment
		if dto.IsClassTeacher {
			// Check if class teacher already exists
			existingCT, err := s.repo.GetClassTeacher(ctx, dto.TenantID, dto.ClassID, dto.SectionID, dto.AcademicYearID)
			if err != nil {
				return err
			}
			if existingCT != nil && existingCT.StaffID != dto.StaffID {
				// Clear existing class teacher
				if err := s.repo.ClearClassTeacher(ctx, tx, dto.TenantID, dto.ClassID, dto.SectionID, dto.AcademicYearID); err != nil {
					return err
				}
			}
		}

		assignment = &models.TeacherSubjectAssignment{
			TenantID:       dto.TenantID,
			StaffID:        dto.StaffID,
			SubjectID:      dto.SubjectID,
			ClassID:        dto.ClassID,
			SectionID:      dto.SectionID,
			AcademicYearID: dto.AcademicYearID,
			PeriodsPerWeek: dto.PeriodsPerWeek,
			IsClassTeacher: dto.IsClassTeacher,
			EffectiveFrom:  dto.EffectiveFrom,
			EffectiveTo:    dto.EffectiveTo,
			Remarks:        dto.Remarks,
			Status:         models.AssignmentStatusActive,
			CreatedBy:      dto.CreatedBy,
		}

		txRepo := NewRepository(tx)
		if err := txRepo.CreateWithTx(ctx, tx, assignment); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch complete assignment with relations
	return s.repo.GetByID(ctx, dto.TenantID, assignment.ID)
}

// GetByID retrieves an assignment by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.TeacherSubjectAssignment, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// Update updates an assignment.
func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, dto UpdateAssignmentDTO) (*models.TeacherSubjectAssignment, error) {
	assignment, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if assignment.Status == models.AssignmentStatusInactive {
		return nil, ErrCannotModifyInactiveAssignment
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Handle periods change with workload validation
		if dto.PeriodsPerWeek != nil && *dto.PeriodsPerWeek != assignment.PeriodsPerWeek {
			periodsDiff := *dto.PeriodsPerWeek - assignment.PeriodsPerWeek
			if periodsDiff > 0 {
				currentPeriods, err := s.repo.GetTotalPeriodsForStaff(ctx, tenantID, assignment.StaffID, assignment.AcademicYearID)
				if err != nil {
					return err
				}

				var staff models.Staff
				if err := s.db.WithContext(ctx).Where("id = ?", assignment.StaffID).First(&staff).Error; err != nil {
					return err
				}

				maxPeriods, err := s.getMaxPeriodsForBranch(ctx, tenantID, staff.BranchID)
				if err != nil {
					return err
				}

				// Subtract current assignment's periods and add new periods
				newTotal := currentPeriods - assignment.PeriodsPerWeek + *dto.PeriodsPerWeek
				if newTotal > maxPeriods {
					return ErrTeacherOverAssigned
				}
			}
			assignment.PeriodsPerWeek = *dto.PeriodsPerWeek
		}

		// Handle class teacher change
		if dto.IsClassTeacher != nil && *dto.IsClassTeacher != assignment.IsClassTeacher {
			if *dto.IsClassTeacher {
				// Setting as class teacher - clear existing
				if err := s.repo.ClearClassTeacher(ctx, tx, tenantID, assignment.ClassID, assignment.SectionID, assignment.AcademicYearID); err != nil {
					return err
				}
			}
			assignment.IsClassTeacher = *dto.IsClassTeacher
		}

		if dto.EffectiveFrom != nil {
			assignment.EffectiveFrom = *dto.EffectiveFrom
		}
		if dto.EffectiveTo != nil {
			assignment.EffectiveTo = dto.EffectiveTo
		}
		if dto.Remarks != nil {
			assignment.Remarks = *dto.Remarks
		}
		if dto.Status != nil {
			if !dto.Status.IsValid() {
				return ErrInvalidStatus
			}
			assignment.Status = *dto.Status
		}

		// Validate effective dates
		if assignment.EffectiveTo != nil && assignment.EffectiveFrom.After(*assignment.EffectiveTo) {
			return ErrInvalidEffectiveDate
		}

		return tx.WithContext(ctx).Save(assignment).Error
	})

	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, tenantID, id)
}

// Delete deletes an assignment.
func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List retrieves assignments with optional filtering.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.TeacherSubjectAssignment, string, int64, error) {
	return s.repo.List(ctx, filter)
}

// GetStaffAssignments retrieves all assignments for a staff member.
func (s *Service) GetStaffAssignments(ctx context.Context, tenantID, staffID uuid.UUID, academicYearID *uuid.UUID) ([]models.TeacherSubjectAssignment, error) {
	return s.repo.GetStaffAssignments(ctx, tenantID, staffID, academicYearID)
}

// GetClassTeacher retrieves the class teacher for a class-section.
func (s *Service) GetClassTeacher(ctx context.Context, tenantID, classID uuid.UUID, sectionID *uuid.UUID, academicYearID uuid.UUID) (*models.TeacherSubjectAssignment, error) {
	return s.repo.GetClassTeacher(ctx, tenantID, classID, sectionID, academicYearID)
}

// SetClassTeacher sets a teacher as the class teacher for a class-section.
func (s *Service) SetClassTeacher(ctx context.Context, tenantID, classID uuid.UUID, sectionID *uuid.UUID, academicYearID, staffID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Clear existing class teacher
		if err := s.repo.ClearClassTeacher(ctx, tx, tenantID, classID, sectionID, academicYearID); err != nil {
			return err
		}

		// Find an existing assignment for this staff-class-section or create a dummy one
		var assignment models.TeacherSubjectAssignment
		query := tx.WithContext(ctx).
			Where("tenant_id = ? AND staff_id = ? AND class_id = ? AND academic_year_id = ? AND status = ?",
				tenantID, staffID, classID, academicYearID, models.AssignmentStatusActive)

		if sectionID != nil {
			query = query.Where("section_id = ?", *sectionID)
		} else {
			query = query.Where("section_id IS NULL")
		}

		err := query.First(&assignment).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// No existing assignment - this teacher must have at least one assignment to be class teacher
				return fmt.Errorf("teacher must have at least one subject assignment in this class to be class teacher")
			}
			return err
		}

		// Update the assignment to be class teacher
		return tx.WithContext(ctx).Model(&assignment).Update("is_class_teacher", true).Error
	})
}

// GetWorkloadReport generates a workload report for all teachers.
func (s *Service) GetWorkloadReport(ctx context.Context, tenantID, academicYearID uuid.UUID, branchID *uuid.UUID) (*WorkloadReportResponse, error) {
	teachers, err := s.repo.GetTeachersWithWorkload(ctx, tenantID, academicYearID, branchID)
	if err != nil {
		return nil, err
	}

	// Get workload settings
	var minPeriods, maxPeriods int
	if branchID != nil {
		settings, err := s.repo.GetWorkloadSettings(ctx, tenantID, *branchID)
		if err != nil {
			return nil, err
		}
		if settings != nil {
			minPeriods = settings.MinPeriodsPerWeek
			maxPeriods = settings.MaxPeriodsPerWeek
		}
	}
	if minPeriods == 0 {
		minPeriods = defaultMinPeriodsPerWeek
	}
	if maxPeriods == 0 {
		maxPeriods = defaultMaxPeriodsPerWeek
	}

	// Calculate status for each teacher
	var overAssigned, underAssigned, normalAssigned int
	for i := range teachers {
		teachers[i].MinPeriods = minPeriods
		teachers[i].MaxPeriods = maxPeriods

		if teachers[i].TotalPeriods > maxPeriods {
			teachers[i].WorkloadStatus = "over"
			overAssigned++
		} else if teachers[i].TotalPeriods < minPeriods {
			teachers[i].WorkloadStatus = "under"
			underAssigned++
		} else {
			teachers[i].WorkloadStatus = "normal"
			normalAssigned++
		}
	}

	return &WorkloadReportResponse{
		Teachers:       teachers,
		TotalTeachers:  len(teachers),
		OverAssigned:   overAssigned,
		UnderAssigned:  underAssigned,
		NormalAssigned: normalAssigned,
	}, nil
}

// GetUnassignedSubjects retrieves subjects without teachers for a given academic year.
func (s *Service) GetUnassignedSubjects(ctx context.Context, tenantID, academicYearID uuid.UUID) (*UnassignedSubjectsResponse, error) {
	subjects, err := s.repo.GetUnassignedSubjects(ctx, tenantID, academicYearID)
	if err != nil {
		return nil, err
	}

	return &UnassignedSubjectsResponse{
		Subjects: subjects,
		Total:    len(subjects),
	}, nil
}

// GetWorkloadSettings retrieves workload settings for a branch.
func (s *Service) GetWorkloadSettings(ctx context.Context, tenantID, branchID uuid.UUID) (*models.TeacherWorkloadSettings, error) {
	settings, err := s.repo.GetWorkloadSettings(ctx, tenantID, branchID)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		// Return default settings
		return &models.TeacherWorkloadSettings{
			TenantID:          tenantID,
			BranchID:          branchID,
			MinPeriodsPerWeek: defaultMinPeriodsPerWeek,
			MaxPeriodsPerWeek: defaultMaxPeriodsPerWeek,
		}, nil
	}
	return settings, nil
}

// UpdateWorkloadSettings creates or updates workload settings.
func (s *Service) UpdateWorkloadSettings(ctx context.Context, tenantID, branchID uuid.UUID, dto WorkloadSettingsDTO) (*models.TeacherWorkloadSettings, error) {
	settings := &models.TeacherWorkloadSettings{
		TenantID:              tenantID,
		BranchID:              branchID,
		MinPeriodsPerWeek:     dto.MinPeriodsPerWeek,
		MaxPeriodsPerWeek:     dto.MaxPeriodsPerWeek,
		MaxSubjectsPerTeacher: dto.MaxSubjectsPerTeacher,
		MaxClassesPerTeacher:  dto.MaxClassesPerTeacher,
	}

	if err := s.repo.CreateOrUpdateWorkloadSettings(ctx, settings); err != nil {
		return nil, err
	}

	return s.repo.GetWorkloadSettings(ctx, tenantID, branchID)
}

// BulkCreate creates multiple assignments at once.
func (s *Service) BulkCreate(ctx context.Context, dto BulkCreateAssignmentDTO) ([]models.TeacherSubjectAssignment, []error) {
	var created []models.TeacherSubjectAssignment
	var errs []error

	for _, item := range dto.Assignments {
		assignment, err := s.Create(ctx, CreateAssignmentDTO{
			TenantID:       dto.TenantID,
			StaffID:        item.StaffID,
			SubjectID:      item.SubjectID,
			ClassID:        item.ClassID,
			SectionID:      item.SectionID,
			AcademicYearID: item.AcademicYearID,
			PeriodsPerWeek: item.PeriodsPerWeek,
			IsClassTeacher: item.IsClassTeacher,
			EffectiveFrom:  item.EffectiveFrom,
			CreatedBy:      dto.CreatedBy,
		})
		if err != nil {
			errs = append(errs, err)
		} else {
			created = append(created, *assignment)
		}
	}

	return created, errs
}

// DeactivateAssignment marks an assignment as inactive.
func (s *Service) DeactivateAssignment(ctx context.Context, tenantID, id uuid.UUID) (*models.TeacherSubjectAssignment, error) {
	assignment, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if assignment.Status == models.AssignmentStatusInactive {
		return nil, ErrAssignmentAlreadyInactive
	}

	assignment.Status = models.AssignmentStatusInactive
	now := time.Now()
	assignment.EffectiveTo = &now

	if err := s.repo.Update(ctx, assignment); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, tenantID, id)
}

// getMaxPeriodsForBranch retrieves the max periods per week for a branch.
func (s *Service) getMaxPeriodsForBranch(ctx context.Context, tenantID, branchID uuid.UUID) (int, error) {
	settings, err := s.repo.GetWorkloadSettings(ctx, tenantID, branchID)
	if err != nil {
		return 0, err
	}
	if settings != nil {
		return settings.MaxPeriodsPerWeek, nil
	}
	return defaultMaxPeriodsPerWeek, nil
}
