// Package assignment provides teacher subject assignment functionality.
package assignment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for teacher assignments.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new assignment repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new teacher assignment.
func (r *Repository) Create(ctx context.Context, assignment *models.TeacherSubjectAssignment) error {
	if err := r.db.WithContext(ctx).Create(assignment).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "uniq_teacher_assignment") {
			return ErrDuplicateAssignment
		}
		return fmt.Errorf("create assignment: %w", err)
	}
	return nil
}

// CreateWithTx creates a new assignment within a transaction.
func (r *Repository) CreateWithTx(ctx context.Context, tx *gorm.DB, assignment *models.TeacherSubjectAssignment) error {
	if err := tx.WithContext(ctx).Create(assignment).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "uniq_teacher_assignment") {
			return ErrDuplicateAssignment
		}
		return fmt.Errorf("create assignment: %w", err)
	}
	return nil
}

// GetByID retrieves an assignment by ID.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.TeacherSubjectAssignment, error) {
	var assignment models.TeacherSubjectAssignment
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Subject").
		Preload("Class").
		Preload("Section").
		Preload("AcademicYear").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&assignment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAssignmentNotFound
		}
		return nil, fmt.Errorf("get assignment by id: %w", err)
	}
	return &assignment, nil
}

// Update updates an assignment.
func (r *Repository) Update(ctx context.Context, assignment *models.TeacherSubjectAssignment) error {
	result := r.db.WithContext(ctx).Save(assignment)
	if result.Error != nil {
		return fmt.Errorf("update assignment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrAssignmentNotFound
	}
	return nil
}

// Delete deletes an assignment.
func (r *Repository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.TeacherSubjectAssignment{})
	if result.Error != nil {
		return fmt.Errorf("delete assignment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrAssignmentNotFound
	}
	return nil
}

// List retrieves assignments with optional filtering and pagination.
func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.TeacherSubjectAssignment, string, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.TeacherSubjectAssignment{}).
		Preload("Staff").
		Preload("Subject").
		Preload("Class").
		Preload("Section").
		Preload("AcademicYear").
		Where("teacher_subject_assignments.tenant_id = ?", filter.TenantID)

	// Apply filters
	if filter.StaffID != nil {
		query = query.Where("teacher_subject_assignments.staff_id = ?", *filter.StaffID)
	}
	if filter.SubjectID != nil {
		query = query.Where("teacher_subject_assignments.subject_id = ?", *filter.SubjectID)
	}
	if filter.ClassID != nil {
		query = query.Where("teacher_subject_assignments.class_id = ?", *filter.ClassID)
	}
	if filter.SectionID != nil {
		query = query.Where("teacher_subject_assignments.section_id = ?", *filter.SectionID)
	}
	if filter.AcademicYearID != nil {
		query = query.Where("teacher_subject_assignments.academic_year_id = ?", *filter.AcademicYearID)
	}
	if filter.IsClassTeacher != nil {
		query = query.Where("teacher_subject_assignments.is_class_teacher = ?", *filter.IsClassTeacher)
	}
	if filter.Status != nil {
		query = query.Where("teacher_subject_assignments.status = ?", *filter.Status)
	}

	// Get total count
	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, "", 0, fmt.Errorf("count assignments: %w", err)
	}

	// Apply cursor pagination
	if filter.Cursor != "" {
		cursorID, err := uuid.Parse(filter.Cursor)
		if err == nil {
			query = query.Where("teacher_subject_assignments.id > ?", cursorID)
		}
	}

	// Apply sorting and limit
	query = query.Order("teacher_subject_assignments.created_at DESC")
	limit := 20
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Limit(limit + 1)

	var assignments []models.TeacherSubjectAssignment
	if err := query.Find(&assignments).Error; err != nil {
		return nil, "", 0, fmt.Errorf("list assignments: %w", err)
	}

	// Calculate next cursor
	var nextCursor string
	if len(assignments) > limit {
		assignments = assignments[:limit]
		nextCursor = assignments[len(assignments)-1].ID.String()
	}

	return assignments, nextCursor, total, nil
}

// GetStaffAssignments retrieves all assignments for a staff member.
func (r *Repository) GetStaffAssignments(ctx context.Context, tenantID, staffID uuid.UUID, academicYearID *uuid.UUID) ([]models.TeacherSubjectAssignment, error) {
	query := r.db.WithContext(ctx).
		Preload("Subject").
		Preload("Class").
		Preload("Section").
		Preload("AcademicYear").
		Where("tenant_id = ? AND staff_id = ? AND status = ?", tenantID, staffID, models.AssignmentStatusActive)

	if academicYearID != nil {
		query = query.Where("academic_year_id = ?", *academicYearID)
	}

	var assignments []models.TeacherSubjectAssignment
	if err := query.Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("get staff assignments: %w", err)
	}
	return assignments, nil
}

// GetClassTeacher retrieves the class teacher for a class-section.
func (r *Repository) GetClassTeacher(ctx context.Context, tenantID, classID uuid.UUID, sectionID *uuid.UUID, academicYearID uuid.UUID) (*models.TeacherSubjectAssignment, error) {
	query := r.db.WithContext(ctx).
		Preload("Staff").
		Where("tenant_id = ? AND class_id = ? AND academic_year_id = ? AND is_class_teacher = ? AND status = ?",
			tenantID, classID, academicYearID, true, models.AssignmentStatusActive)

	if sectionID != nil {
		query = query.Where("section_id = ?", *sectionID)
	} else {
		query = query.Where("section_id IS NULL")
	}

	var assignment models.TeacherSubjectAssignment
	err := query.First(&assignment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No class teacher assigned
		}
		return nil, fmt.Errorf("get class teacher: %w", err)
	}
	return &assignment, nil
}

// GetTotalPeriodsForStaff calculates the total periods per week for a staff member.
func (r *Repository) GetTotalPeriodsForStaff(ctx context.Context, tenantID, staffID, academicYearID uuid.UUID) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.TeacherSubjectAssignment{}).
		Where("tenant_id = ? AND staff_id = ? AND academic_year_id = ? AND status = ?",
			tenantID, staffID, academicYearID, models.AssignmentStatusActive).
		Select("COALESCE(SUM(periods_per_week), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("get total periods: %w", err)
	}
	return int(total), nil
}

// GetWorkloadSettings retrieves workload settings for a branch.
func (r *Repository) GetWorkloadSettings(ctx context.Context, tenantID, branchID uuid.UUID) (*models.TeacherWorkloadSettings, error) {
	var settings models.TeacherWorkloadSettings
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Where("tenant_id = ? AND branch_id = ?", tenantID, branchID).
		First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No settings configured, use defaults
		}
		return nil, fmt.Errorf("get workload settings: %w", err)
	}
	return &settings, nil
}

// CreateOrUpdateWorkloadSettings creates or updates workload settings.
func (r *Repository) CreateOrUpdateWorkloadSettings(ctx context.Context, settings *models.TeacherWorkloadSettings) error {
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ?", settings.TenantID, settings.BranchID).
		Assign(map[string]interface{}{
			"min_periods_per_week":     settings.MinPeriodsPerWeek,
			"max_periods_per_week":     settings.MaxPeriodsPerWeek,
			"max_subjects_per_teacher": settings.MaxSubjectsPerTeacher,
			"max_classes_per_teacher":  settings.MaxClassesPerTeacher,
		}).
		FirstOrCreate(settings).Error
	if err != nil {
		return fmt.Errorf("create or update workload settings: %w", err)
	}
	return nil
}

// GetTeachersWithWorkload retrieves all teaching staff with their workload summaries.
func (r *Repository) GetTeachersWithWorkload(ctx context.Context, tenantID, academicYearID uuid.UUID, branchID *uuid.UUID) ([]WorkloadSummary, error) {
	// Build query to get all teaching staff with their assignment counts
	query := r.db.WithContext(ctx).
		Table("staff").
		Select(`
			staff.id as staff_id,
			CONCAT(staff.first_name, ' ', staff.last_name) as staff_name,
			staff.employee_id as staff_employee_id,
			COALESCE(departments.name, '') as department_name,
			COALESCE(SUM(tsa.periods_per_week), 0) as total_periods,
			COUNT(DISTINCT tsa.subject_id) as total_subjects,
			COUNT(DISTINCT CONCAT(tsa.class_id, COALESCE(tsa.section_id::text, ''))) as total_classes,
			BOOL_OR(tsa.is_class_teacher) as is_class_teacher
		`).
		Joins("LEFT JOIN teacher_subject_assignments tsa ON tsa.staff_id = staff.id AND tsa.academic_year_id = ? AND tsa.status = 'active' AND tsa.tenant_id = ?", academicYearID, tenantID).
		Joins("LEFT JOIN departments ON departments.id = staff.department_id").
		Where("staff.tenant_id = ? AND staff.staff_type = ? AND staff.status = ?", tenantID, models.StaffTypeTeaching, models.StaffStatusActive).
		Group("staff.id, staff.first_name, staff.last_name, staff.employee_id, departments.name")

	if branchID != nil {
		query = query.Where("staff.branch_id = ?", *branchID)
	}

	type workloadRow struct {
		StaffID        uuid.UUID `gorm:"column:staff_id"`
		StaffName      string    `gorm:"column:staff_name"`
		StaffEmployeeID string   `gorm:"column:staff_employee_id"`
		DepartmentName string    `gorm:"column:department_name"`
		TotalPeriods   int       `gorm:"column:total_periods"`
		TotalSubjects  int       `gorm:"column:total_subjects"`
		TotalClasses   int       `gorm:"column:total_classes"`
		IsClassTeacher bool      `gorm:"column:is_class_teacher"`
	}

	var rows []workloadRow
	if err := query.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("get teachers with workload: %w", err)
	}

	summaries := make([]WorkloadSummary, len(rows))
	for i, row := range rows {
		summaries[i] = WorkloadSummary{
			StaffID:         row.StaffID.String(),
			StaffName:       row.StaffName,
			StaffEmployeeID: row.StaffEmployeeID,
			DepartmentName:  row.DepartmentName,
			TotalPeriods:    row.TotalPeriods,
			TotalSubjects:   row.TotalSubjects,
			TotalClasses:    row.TotalClasses,
			IsClassTeacher:  row.IsClassTeacher,
		}
	}

	return summaries, nil
}

// GetUnassignedSubjects retrieves subjects without teachers for a given academic year.
func (r *Repository) GetUnassignedSubjects(ctx context.Context, tenantID, academicYearID uuid.UUID) ([]UnassignedSubject, error) {
	// Get all class-subjects that don't have an active assignment
	query := r.db.WithContext(ctx).
		Table("class_subjects cs").
		Select(`
			cs.subject_id,
			subjects.name as subject_name,
			subjects.code as subject_code,
			cs.class_id,
			classes.name as class_name,
			sections.id as section_id,
			sections.name as section_name
		`).
		Joins("JOIN subjects ON subjects.id = cs.subject_id").
		Joins("JOIN classes ON classes.id = cs.class_id").
		Joins("LEFT JOIN sections ON sections.class_id = cs.class_id AND sections.is_active = true").
		Joins(`LEFT JOIN teacher_subject_assignments tsa ON
			tsa.subject_id = cs.subject_id AND
			tsa.class_id = cs.class_id AND
			(tsa.section_id = sections.id OR (tsa.section_id IS NULL AND sections.id IS NULL)) AND
			tsa.academic_year_id = ? AND
			tsa.status = 'active' AND
			tsa.tenant_id = ?`, academicYearID, tenantID).
		Where("cs.tenant_id = ? AND cs.is_active = true AND tsa.id IS NULL", tenantID)

	type unassignedRow struct {
		SubjectID   uuid.UUID  `gorm:"column:subject_id"`
		SubjectName string     `gorm:"column:subject_name"`
		SubjectCode string     `gorm:"column:subject_code"`
		ClassID     uuid.UUID  `gorm:"column:class_id"`
		ClassName   string     `gorm:"column:class_name"`
		SectionID   *uuid.UUID `gorm:"column:section_id"`
		SectionName *string    `gorm:"column:section_name"`
	}

	var rows []unassignedRow
	if err := query.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("get unassigned subjects: %w", err)
	}

	subjects := make([]UnassignedSubject, len(rows))
	for i, row := range rows {
		subjects[i] = UnassignedSubject{
			SubjectID:   row.SubjectID.String(),
			SubjectName: row.SubjectName,
			SubjectCode: row.SubjectCode,
			ClassID:     row.ClassID.String(),
			ClassName:   row.ClassName,
		}
		if row.SectionID != nil {
			subjects[i].SectionID = row.SectionID.String()
		}
		if row.SectionName != nil {
			subjects[i].SectionName = *row.SectionName
		}
	}

	return subjects, nil
}

// ClearClassTeacher removes the class teacher flag from existing assignments for a class-section.
func (r *Repository) ClearClassTeacher(ctx context.Context, tx *gorm.DB, tenantID, classID uuid.UUID, sectionID *uuid.UUID, academicYearID uuid.UUID) error {
	query := tx.WithContext(ctx).
		Model(&models.TeacherSubjectAssignment{}).
		Where("tenant_id = ? AND class_id = ? AND academic_year_id = ? AND is_class_teacher = ?",
			tenantID, classID, academicYearID, true)

	if sectionID != nil {
		query = query.Where("section_id = ?", *sectionID)
	} else {
		query = query.Where("section_id IS NULL")
	}

	if err := query.Update("is_class_teacher", false).Error; err != nil {
		return fmt.Errorf("clear class teacher: %w", err)
	}
	return nil
}

// DB returns the underlying database connection for transactions.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
