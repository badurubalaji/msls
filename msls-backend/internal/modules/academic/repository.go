// Package academic provides academic structure management functionality.
package academic

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for academic entities.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new academic repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ========================================
// Class Repository Methods
// ========================================

// ListClasses returns all classes for a tenant with filters.
func (r *Repository) ListClasses(ctx context.Context, filter ClassFilter) ([]models.Class, int64, error) {
	var classes []models.Class
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Class{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.HasStreams != nil {
		query = query.Where("has_streams = ?", *filter.HasStreams)
	}

	if filter.Level != nil {
		query = query.Where("level = ?", *filter.Level)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Branch").
		Preload("Streams").
		Order("display_order ASC, name ASC").
		Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// GetClassByID returns a class by ID.
func (r *Repository) GetClassByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Class, error) {
	var class models.Class
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC, name ASC")
		}).
		Preload("Sections.ClassTeacher").
		Preload("Sections.Stream").
		Preload("Sections.AcademicYear").
		Preload("Streams").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&class).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrClassNotFound
	}
	return &class, err
}

// GetClassByCode returns a class by code.
func (r *Repository) GetClassByCode(ctx context.Context, tenantID, branchID uuid.UUID, code string) (*models.Class, error) {
	var class models.Class
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ? AND code = ?", tenantID, branchID, code).
		First(&class).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &class, err
}

// CreateClass creates a new class.
func (r *Repository) CreateClass(ctx context.Context, class *models.Class) error {
	return r.db.WithContext(ctx).Create(class).Error
}

// UpdateClass updates an existing class.
func (r *Repository) UpdateClass(ctx context.Context, class *models.Class) error {
	return r.db.WithContext(ctx).Save(class).Error
}

// DeleteClass deletes a class.
func (r *Repository) DeleteClass(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Class{}).Error
}

// CountSectionsByClassID returns the number of sections for a class.
func (r *Repository) CountSectionsByClassID(ctx context.Context, classID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Section{}).
		Where("class_id = ?", classID).
		Count(&count).Error
	return count, err
}

// SetClassStreams sets the streams for a class.
func (r *Repository) SetClassStreams(ctx context.Context, tenantID, classID uuid.UUID, streamIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing class-stream mappings
		if err := tx.Where("class_id = ?", classID).Delete(&models.ClassStream{}).Error; err != nil {
			return err
		}

		// Create new mappings
		for _, streamID := range streamIDs {
			mapping := &models.ClassStream{
				TenantID: tenantID,
				ClassID:  classID,
				StreamID: streamID,
				IsActive: true,
			}
			if err := tx.Create(mapping).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ========================================
// Section Repository Methods
// ========================================

// ListSections returns all sections for a tenant with filters.
func (r *Repository) ListSections(ctx context.Context, filter SectionFilter) ([]models.Section, int64, error) {
	var sections []models.Section
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Section{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.ClassID != nil {
		query = query.Where("class_id = ?", *filter.ClassID)
	}

	if filter.AcademicYearID != nil {
		query = query.Where("academic_year_id = ?", *filter.AcademicYearID)
	}

	if filter.StreamID != nil {
		query = query.Where("stream_id = ?", *filter.StreamID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Class").
		Preload("AcademicYear").
		Preload("Stream").
		Preload("ClassTeacher").
		Order("display_order ASC, name ASC").
		Find(&sections).Error; err != nil {
		return nil, 0, err
	}

	return sections, total, nil
}

// GetSectionByID returns a section by ID.
func (r *Repository) GetSectionByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Section, error) {
	var section models.Section
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("AcademicYear").
		Preload("Stream").
		Preload("ClassTeacher").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&section).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSectionNotFound
	}
	return &section, err
}

// GetSectionByCode returns a section by code within a class.
func (r *Repository) GetSectionByCode(ctx context.Context, tenantID, classID uuid.UUID, code string) (*models.Section, error) {
	var section models.Section
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND class_id = ? AND code = ?", tenantID, classID, code).
		First(&section).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &section, err
}

// CreateSection creates a new section.
func (r *Repository) CreateSection(ctx context.Context, section *models.Section) error {
	return r.db.WithContext(ctx).Create(section).Error
}

// UpdateSection updates an existing section.
func (r *Repository) UpdateSection(ctx context.Context, section *models.Section) error {
	return r.db.WithContext(ctx).Save(section).Error
}

// DeleteSection deletes a section.
func (r *Repository) DeleteSection(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Section{}).Error
}

// CountStudentsBySectionID returns the number of students in a section.
func (r *Repository) CountStudentsBySectionID(ctx context.Context, sectionID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Student{}).
		Where("section_id = ?", sectionID).
		Count(&count).Error
	return count, err
}

// GetSectionStudentCounts returns student counts for multiple sections.
func (r *Repository) GetSectionStudentCounts(ctx context.Context, sectionIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	type result struct {
		SectionID uuid.UUID
		Count     int
	}

	var results []result
	err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Select("section_id, COUNT(*) as count").
		Where("section_id IN ?", sectionIDs).
		Group("section_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[uuid.UUID]int)
	for _, r := range results {
		counts[r.SectionID] = r.Count
	}
	return counts, nil
}

// ========================================
// Stream Repository Methods
// ========================================

// ListStreams returns all streams for a tenant with filters.
func (r *Repository) ListStreams(ctx context.Context, filter StreamFilter) ([]models.Stream, int64, error) {
	var streams []models.Stream
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Stream{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("display_order ASC, name ASC").
		Find(&streams).Error; err != nil {
		return nil, 0, err
	}

	return streams, total, nil
}

// GetStreamByID returns a stream by ID.
func (r *Repository) GetStreamByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Stream, error) {
	var stream models.Stream
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&stream).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStreamNotFound
	}
	return &stream, err
}

// GetStreamByCode returns a stream by code.
func (r *Repository) GetStreamByCode(ctx context.Context, tenantID uuid.UUID, code string) (*models.Stream, error) {
	var stream models.Stream
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&stream).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &stream, err
}

// CreateStream creates a new stream.
func (r *Repository) CreateStream(ctx context.Context, stream *models.Stream) error {
	return r.db.WithContext(ctx).Create(stream).Error
}

// UpdateStream updates an existing stream.
func (r *Repository) UpdateStream(ctx context.Context, stream *models.Stream) error {
	return r.db.WithContext(ctx).Save(stream).Error
}

// DeleteStream deletes a stream.
func (r *Repository) DeleteStream(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Stream{}).Error
}

// IsStreamInUse checks if a stream is used by any class or section.
func (r *Repository) IsStreamInUse(ctx context.Context, streamID uuid.UUID) (bool, error) {
	var classCount int64
	if err := r.db.WithContext(ctx).Model(&models.ClassStream{}).
		Where("stream_id = ?", streamID).
		Count(&classCount).Error; err != nil {
		return false, err
	}

	if classCount > 0 {
		return true, nil
	}

	var sectionCount int64
	if err := r.db.WithContext(ctx).Model(&models.Section{}).
		Where("stream_id = ?", streamID).
		Count(&sectionCount).Error; err != nil {
		return false, err
	}

	return sectionCount > 0, nil
}

// ========================================
// Structure Repository Methods
// ========================================

// GetClassStructure returns the hierarchical class-section structure.
func (r *Repository) GetClassStructure(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID, academicYearID *uuid.UUID) ([]models.Class, error) {
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = true", tenantID)

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}

	sectionQuery := func(db *gorm.DB) *gorm.DB {
		q := db.Where("is_active = true").Order("display_order ASC, name ASC")
		if academicYearID != nil {
			q = q.Where("academic_year_id = ?", *academicYearID)
		}
		return q
	}

	var classes []models.Class
	err := query.
		Preload("Sections", sectionQuery).
		Preload("Sections.ClassTeacher").
		Preload("Sections.Stream").
		Order("display_order ASC, name ASC").
		Find(&classes).Error

	return classes, err
}

// ========================================
// Subject Repository Methods
// ========================================

// ListSubjects returns all subjects for a tenant with filters.
func (r *Repository) ListSubjects(ctx context.Context, filter SubjectFilter) ([]models.Subject, int64, error) {
	var subjects []models.Subject
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Subject{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.SubjectType != nil {
		query = query.Where("subject_type = ?", *filter.SubjectType)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR short_name ILIKE ?", search, search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("display_order ASC, name ASC").
		Find(&subjects).Error; err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}

// GetSubjectByID returns a subject by ID.
func (r *Repository) GetSubjectByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&subject).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSubjectNotFound
	}
	return &subject, err
}

// GetSubjectByCode returns a subject by code.
func (r *Repository) GetSubjectByCode(ctx context.Context, tenantID uuid.UUID, code string) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&subject).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subject, err
}

// CreateSubject creates a new subject.
func (r *Repository) CreateSubject(ctx context.Context, subject *models.Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

// UpdateSubject updates an existing subject.
func (r *Repository) UpdateSubject(ctx context.Context, subject *models.Subject) error {
	return r.db.WithContext(ctx).Save(subject).Error
}

// DeleteSubject deletes a subject.
func (r *Repository) DeleteSubject(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Subject{}).Error
}

// IsSubjectInUse checks if a subject is used in any class-subject mapping or teacher assignment.
func (r *Repository) IsSubjectInUse(ctx context.Context, subjectID uuid.UUID) (bool, error) {
	var classSubjectCount int64
	if err := r.db.WithContext(ctx).Model(&models.ClassSubject{}).
		Where("subject_id = ?", subjectID).
		Count(&classSubjectCount).Error; err != nil {
		return false, err
	}

	if classSubjectCount > 0 {
		return true, nil
	}

	var assignmentCount int64
	if err := r.db.WithContext(ctx).Model(&models.TeacherSubjectAssignment{}).
		Where("subject_id = ?", subjectID).
		Count(&assignmentCount).Error; err != nil {
		return false, err
	}

	return assignmentCount > 0, nil
}

// ========================================
// Class-Subject Repository Methods
// ========================================

// ListClassSubjects returns all class-subject mappings for a tenant with filters.
func (r *Repository) ListClassSubjects(ctx context.Context, filter ClassSubjectFilter) ([]models.ClassSubject, int64, error) {
	var classSubjects []models.ClassSubject
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ClassSubject{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.ClassID != nil {
		query = query.Where("class_id = ?", *filter.ClassID)
	}

	if filter.SubjectID != nil {
		query = query.Where("subject_id = ?", *filter.SubjectID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Class").
		Preload("Subject").
		Order("created_at ASC").
		Find(&classSubjects).Error; err != nil {
		return nil, 0, err
	}

	return classSubjects, total, nil
}

// GetClassSubjectByID returns a class-subject mapping by ID.
func (r *Repository) GetClassSubjectByID(ctx context.Context, tenantID, id uuid.UUID) (*models.ClassSubject, error) {
	var classSubject models.ClassSubject
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("Subject").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&classSubject).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrClassSubjectNotFound
	}
	return &classSubject, err
}

// GetClassSubjectByClassAndSubject returns a class-subject mapping by class and subject IDs.
func (r *Repository) GetClassSubjectByClassAndSubject(ctx context.Context, tenantID, classID, subjectID uuid.UUID) (*models.ClassSubject, error) {
	var classSubject models.ClassSubject
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND class_id = ? AND subject_id = ?", tenantID, classID, subjectID).
		First(&classSubject).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &classSubject, err
}

// CreateClassSubject creates a new class-subject mapping.
func (r *Repository) CreateClassSubject(ctx context.Context, classSubject *models.ClassSubject) error {
	return r.db.WithContext(ctx).Create(classSubject).Error
}

// UpdateClassSubject updates an existing class-subject mapping.
func (r *Repository) UpdateClassSubject(ctx context.Context, classSubject *models.ClassSubject) error {
	return r.db.WithContext(ctx).Save(classSubject).Error
}

// DeleteClassSubject deletes a class-subject mapping.
func (r *Repository) DeleteClassSubject(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.ClassSubject{}).Error
}
