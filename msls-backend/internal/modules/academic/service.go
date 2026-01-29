// Package academic provides academic structure management functionality.
package academic

import (
	"context"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// Service handles business logic for academic entities.
type Service struct {
	repo *Repository
}

// NewService creates a new academic service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ========================================
// Class Service Methods
// ========================================

// ListClasses returns all classes for a tenant with filters.
func (s *Service) ListClasses(ctx context.Context, filter ClassFilter) ([]models.Class, int64, error) {
	return s.repo.ListClasses(ctx, filter)
}

// GetClassByID returns a class by ID.
func (s *Service) GetClassByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Class, error) {
	return s.repo.GetClassByID(ctx, tenantID, id)
}

// CreateClass creates a new class.
func (s *Service) CreateClass(ctx context.Context, tenantID uuid.UUID, req CreateClassRequest, userID uuid.UUID) (*models.Class, error) {
	// Check if code already exists
	existing, err := s.repo.GetClassByCode(ctx, tenantID, req.BranchID, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrClassCodeExists
	}

	class := &models.Class{
		TenantID:     tenantID,
		BranchID:     req.BranchID,
		Name:         req.Name,
		Code:         req.Code,
		DisplayOrder: req.DisplayOrder,
		Description:  req.Description,
		HasStreams:   req.HasStreams,
		IsActive:     true,
		CreatedBy:    &userID,
	}

	// Set level if provided
	if req.Level != "" {
		level := models.ClassLevel(req.Level)
		class.Level = &level
	}

	if err := s.repo.CreateClass(ctx, class); err != nil {
		return nil, err
	}

	// Set streams if provided
	if len(req.StreamIDs) > 0 && req.HasStreams {
		if err := s.repo.SetClassStreams(ctx, tenantID, class.ID, req.StreamIDs); err != nil {
			return nil, err
		}
	}

	// Reload with relationships
	return s.repo.GetClassByID(ctx, tenantID, class.ID)
}

// UpdateClass updates an existing class.
func (s *Service) UpdateClass(ctx context.Context, tenantID, id uuid.UUID, req UpdateClassRequest) (*models.Class, error) {
	class, err := s.repo.GetClassByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check code uniqueness if changing
	if req.Code != nil && *req.Code != class.Code {
		existing, err := s.repo.GetClassByCode(ctx, tenantID, class.BranchID, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrClassCodeExists
		}
		class.Code = *req.Code
	}

	if req.Name != nil {
		class.Name = *req.Name
	}
	if req.DisplayOrder != nil {
		class.DisplayOrder = *req.DisplayOrder
	}
	if req.Description != nil {
		class.Description = *req.Description
	}
	if req.HasStreams != nil {
		class.HasStreams = *req.HasStreams
	}
	if req.IsActive != nil {
		class.IsActive = *req.IsActive
	}
	if req.Level != nil {
		level := models.ClassLevel(*req.Level)
		class.Level = &level
	}

	if err := s.repo.UpdateClass(ctx, class); err != nil {
		return nil, err
	}

	// Update streams if provided
	if req.StreamIDs != nil {
		if err := s.repo.SetClassStreams(ctx, tenantID, id, req.StreamIDs); err != nil {
			return nil, err
		}
	}

	// Reload with relationships
	return s.repo.GetClassByID(ctx, tenantID, id)
}

// DeleteClass deletes a class.
func (s *Service) DeleteClass(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if class exists
	_, err := s.repo.GetClassByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Check if class has sections
	count, err := s.repo.CountSectionsByClassID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrClassHasSections
	}

	return s.repo.DeleteClass(ctx, tenantID, id)
}

// ========================================
// Section Service Methods
// ========================================

// ListSections returns all sections for a tenant with filters.
func (s *Service) ListSections(ctx context.Context, filter SectionFilter) ([]models.Section, int64, error) {
	sections, total, err := s.repo.ListSections(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get student counts
	if len(sections) > 0 {
		sectionIDs := make([]uuid.UUID, len(sections))
		for i, sec := range sections {
			sectionIDs[i] = sec.ID
		}

		counts, err := s.repo.GetSectionStudentCounts(ctx, sectionIDs)
		if err != nil {
			return nil, 0, err
		}

		for i := range sections {
			sections[i].StudentCount = counts[sections[i].ID]
		}
	}

	return sections, total, nil
}

// GetSectionByID returns a section by ID.
func (s *Service) GetSectionByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Section, error) {
	section, err := s.repo.GetSectionByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Get student count
	count, err := s.repo.CountStudentsBySectionID(ctx, id)
	if err != nil {
		return nil, err
	}
	section.StudentCount = int(count)

	return section, nil
}

// CreateSection creates a new section.
func (s *Service) CreateSection(ctx context.Context, tenantID uuid.UUID, req CreateSectionRequest, userID uuid.UUID) (*models.Section, error) {
	// Check if code already exists for this class
	existing, err := s.repo.GetSectionByCode(ctx, tenantID, req.ClassID, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrSectionCodeExists
	}

	capacity := 40
	if req.Capacity > 0 {
		capacity = req.Capacity
	}

	section := &models.Section{
		TenantID:       tenantID,
		ClassID:        req.ClassID,
		AcademicYearID: req.AcademicYearID,
		StreamID:       req.StreamID,
		ClassTeacherID: req.ClassTeacherID,
		Name:           req.Name,
		Code:           req.Code,
		Capacity:       &capacity,
		RoomNumber:     req.RoomNumber,
		IsActive:       true,
		CreatedBy:      &userID,
	}

	if err := s.repo.CreateSection(ctx, section); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.repo.GetSectionByID(ctx, tenantID, section.ID)
}

// UpdateSection updates an existing section.
func (s *Service) UpdateSection(ctx context.Context, tenantID, id uuid.UUID, req UpdateSectionRequest) (*models.Section, error) {
	section, err := s.repo.GetSectionByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check code uniqueness if changing
	if req.Code != nil && *req.Code != section.Code {
		existing, err := s.repo.GetSectionByCode(ctx, tenantID, section.ClassID, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrSectionCodeExists
		}
		section.Code = *req.Code
	}

	if req.Name != nil {
		section.Name = *req.Name
	}
	if req.Capacity != nil {
		section.Capacity = req.Capacity
	}
	if req.DisplayOrder != nil {
		section.DisplayOrder = *req.DisplayOrder
	}
	if req.RoomNumber != nil {
		section.RoomNumber = *req.RoomNumber
	}
	if req.IsActive != nil {
		section.IsActive = *req.IsActive
	}

	// Handle nullable fields
	if req.AcademicYearID != nil {
		section.AcademicYearID = req.AcademicYearID
	}
	if req.StreamID != nil {
		section.StreamID = req.StreamID
	}
	if req.ClassTeacherID != nil {
		section.ClassTeacherID = req.ClassTeacherID
	}

	if err := s.repo.UpdateSection(ctx, section); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.repo.GetSectionByID(ctx, tenantID, id)
}

// DeleteSection deletes a section.
func (s *Service) DeleteSection(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if section exists
	_, err := s.repo.GetSectionByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Check if section has students
	count, err := s.repo.CountStudentsBySectionID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSectionHasStudents
	}

	return s.repo.DeleteSection(ctx, tenantID, id)
}

// AssignClassTeacher assigns a class teacher to a section.
func (s *Service) AssignClassTeacher(ctx context.Context, tenantID, sectionID uuid.UUID, teacherID *uuid.UUID) (*models.Section, error) {
	section, err := s.repo.GetSectionByID(ctx, tenantID, sectionID)
	if err != nil {
		return nil, err
	}

	section.ClassTeacherID = teacherID

	if err := s.repo.UpdateSection(ctx, section); err != nil {
		return nil, err
	}

	return s.repo.GetSectionByID(ctx, tenantID, sectionID)
}

// ========================================
// Stream Service Methods
// ========================================

// ListStreams returns all streams for a tenant with filters.
func (s *Service) ListStreams(ctx context.Context, filter StreamFilter) ([]models.Stream, int64, error) {
	return s.repo.ListStreams(ctx, filter)
}

// GetStreamByID returns a stream by ID.
func (s *Service) GetStreamByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Stream, error) {
	return s.repo.GetStreamByID(ctx, tenantID, id)
}

// CreateStream creates a new stream.
func (s *Service) CreateStream(ctx context.Context, tenantID uuid.UUID, req CreateStreamRequest, userID uuid.UUID) (*models.Stream, error) {
	// Check if code already exists
	existing, err := s.repo.GetStreamByCode(ctx, tenantID, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrStreamCodeExists
	}

	stream := &models.Stream{
		TenantID:     tenantID,
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		DisplayOrder: req.DisplayOrder,
		IsActive:     true,
		CreatedBy:    &userID,
	}

	if err := s.repo.CreateStream(ctx, stream); err != nil {
		return nil, err
	}

	return stream, nil
}

// UpdateStream updates an existing stream.
func (s *Service) UpdateStream(ctx context.Context, tenantID, id uuid.UUID, req UpdateStreamRequest) (*models.Stream, error) {
	stream, err := s.repo.GetStreamByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check code uniqueness if changing
	if req.Code != nil && *req.Code != stream.Code {
		existing, err := s.repo.GetStreamByCode(ctx, tenantID, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrStreamCodeExists
		}
		stream.Code = *req.Code
	}

	if req.Name != nil {
		stream.Name = *req.Name
	}
	if req.Description != nil {
		stream.Description = *req.Description
	}
	if req.DisplayOrder != nil {
		stream.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		stream.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateStream(ctx, stream); err != nil {
		return nil, err
	}

	return stream, nil
}

// DeleteStream deletes a stream.
func (s *Service) DeleteStream(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if stream exists
	_, err := s.repo.GetStreamByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Check if stream is in use
	inUse, err := s.repo.IsStreamInUse(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return ErrStreamInUse
	}

	return s.repo.DeleteStream(ctx, tenantID, id)
}

// ========================================
// Structure Service Methods
// ========================================

// GetClassStructure returns the hierarchical class-section structure.
func (s *Service) GetClassStructure(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID, academicYearID *uuid.UUID) ([]ClassWithSectionsResponse, error) {
	classes, err := s.repo.GetClassStructure(ctx, tenantID, branchID, academicYearID)
	if err != nil {
		return nil, err
	}

	// Collect all section IDs
	var sectionIDs []uuid.UUID
	for _, class := range classes {
		for _, section := range class.Sections {
			sectionIDs = append(sectionIDs, section.ID)
		}
	}

	// Get student counts
	var studentCounts map[uuid.UUID]int
	if len(sectionIDs) > 0 {
		studentCounts, err = s.repo.GetSectionStudentCounts(ctx, sectionIDs)
		if err != nil {
			return nil, err
		}
	} else {
		studentCounts = make(map[uuid.UUID]int)
	}

	// Build response
	result := make([]ClassWithSectionsResponse, len(classes))
	for i, class := range classes {
		classResp := ClassWithSectionsResponse{
			ID:           class.ID,
			Name:         class.Name,
			Code:         class.Code,
			DisplayOrder: class.DisplayOrder,
			HasStreams:   class.HasStreams,
			IsActive:     class.IsActive,
			Sections:     make([]SectionStructureResponse, len(class.Sections)),
		}

		var totalStudents, totalCapacity int
		for j, section := range class.Sections {
			studentCount := studentCounts[section.ID]
			capacity := 40
			if section.Capacity != nil {
				capacity = *section.Capacity
			}

			capacityUsage := 0.0
			if capacity > 0 {
				capacityUsage = float64(studentCount) / float64(capacity) * 100
			}

			secResp := SectionStructureResponse{
				ID:            section.ID,
				Name:          section.Name,
				Code:          section.Code,
				Capacity:      capacity,
				StudentCount:  studentCount,
				RoomNumber:    section.RoomNumber,
				CapacityUsage: capacityUsage,
			}

			if section.ClassTeacherID != nil {
				idStr := section.ClassTeacherID.String()
				secResp.ClassTeacherID = &idStr
				if section.ClassTeacher != nil {
					secResp.ClassTeacherName = section.ClassTeacher.FirstName
					if section.ClassTeacher.LastName != "" {
						secResp.ClassTeacherName += " " + section.ClassTeacher.LastName
					}
				}
			}

			if section.Stream != nil {
				secResp.StreamName = section.Stream.Name
			}

			classResp.Sections[j] = secResp
			totalStudents += studentCount
			totalCapacity += capacity
		}

		classResp.TotalStudents = totalStudents
		classResp.TotalCapacity = totalCapacity
		result[i] = classResp
	}

	return result, nil
}

// ========================================
// Subject Service Methods
// ========================================

// ListSubjects returns all subjects for a tenant with filters.
func (s *Service) ListSubjects(ctx context.Context, filter SubjectFilter) ([]models.Subject, int64, error) {
	return s.repo.ListSubjects(ctx, filter)
}

// GetSubjectByID returns a subject by ID.
func (s *Service) GetSubjectByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Subject, error) {
	return s.repo.GetSubjectByID(ctx, tenantID, id)
}

// CreateSubject creates a new subject.
func (s *Service) CreateSubject(ctx context.Context, tenantID uuid.UUID, req CreateSubjectRequest, userID uuid.UUID) (*models.Subject, error) {
	// Check if code already exists
	existing, err := s.repo.GetSubjectByCode(ctx, tenantID, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrSubjectCodeExists
	}

	// Set defaults
	maxMarks := 100
	if req.MaxMarks > 0 {
		maxMarks = req.MaxMarks
	}

	passingMarks := 35
	if req.PassingMarks > 0 {
		passingMarks = req.PassingMarks
	}

	subject := &models.Subject{
		TenantID:     tenantID,
		Name:         req.Name,
		Code:         req.Code,
		ShortName:    req.ShortName,
		Description:  req.Description,
		SubjectType:  models.SubjectType(req.SubjectType),
		MaxMarks:     &maxMarks,
		PassingMarks: &passingMarks,
		CreditHours:  &req.CreditHours,
		DisplayOrder: req.DisplayOrder,
		IsActive:     true,
		CreatedBy:    &userID,
	}

	if err := s.repo.CreateSubject(ctx, subject); err != nil {
		return nil, err
	}

	return subject, nil
}

// UpdateSubject updates an existing subject.
func (s *Service) UpdateSubject(ctx context.Context, tenantID, id uuid.UUID, req UpdateSubjectRequest) (*models.Subject, error) {
	subject, err := s.repo.GetSubjectByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Check code uniqueness if changing
	if req.Code != nil && *req.Code != subject.Code {
		existing, err := s.repo.GetSubjectByCode(ctx, tenantID, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrSubjectCodeExists
		}
		subject.Code = *req.Code
	}

	if req.Name != nil {
		subject.Name = *req.Name
	}
	if req.ShortName != nil {
		subject.ShortName = *req.ShortName
	}
	if req.Description != nil {
		subject.Description = *req.Description
	}
	if req.SubjectType != nil {
		subject.SubjectType = models.SubjectType(*req.SubjectType)
	}
	if req.MaxMarks != nil {
		subject.MaxMarks = req.MaxMarks
	}
	if req.PassingMarks != nil {
		subject.PassingMarks = req.PassingMarks
	}
	if req.CreditHours != nil {
		subject.CreditHours = req.CreditHours
	}
	if req.DisplayOrder != nil {
		subject.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		subject.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateSubject(ctx, subject); err != nil {
		return nil, err
	}

	return subject, nil
}

// DeleteSubject deletes a subject.
func (s *Service) DeleteSubject(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if subject exists
	_, err := s.repo.GetSubjectByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Check if subject is in use
	inUse, err := s.repo.IsSubjectInUse(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return ErrSubjectInUse
	}

	return s.repo.DeleteSubject(ctx, tenantID, id)
}

// ========================================
// Class-Subject Service Methods
// ========================================

// ListClassSubjects returns all class-subject mappings for a tenant with filters.
func (s *Service) ListClassSubjects(ctx context.Context, filter ClassSubjectFilter) ([]models.ClassSubject, int64, error) {
	return s.repo.ListClassSubjects(ctx, filter)
}

// GetClassSubjectByID returns a class-subject mapping by ID.
func (s *Service) GetClassSubjectByID(ctx context.Context, tenantID, id uuid.UUID) (*models.ClassSubject, error) {
	return s.repo.GetClassSubjectByID(ctx, tenantID, id)
}

// CreateClassSubject creates a new class-subject mapping.
func (s *Service) CreateClassSubject(ctx context.Context, tenantID, classID uuid.UUID, req CreateClassSubjectRequest) (*models.ClassSubject, error) {
	// Check if mapping already exists
	existing, err := s.repo.GetClassSubjectByClassAndSubject(ctx, tenantID, classID, req.SubjectID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrClassSubjectExists
	}

	// Set default periods per week
	periodsPerWeek := 5
	if req.PeriodsPerWeek > 0 {
		periodsPerWeek = req.PeriodsPerWeek
	}

	classSubject := &models.ClassSubject{
		TenantID:       tenantID,
		ClassID:        classID,
		SubjectID:      req.SubjectID,
		IsMandatory:    req.IsMandatory,
		PeriodsPerWeek: periodsPerWeek,
		IsActive:       true,
	}

	if err := s.repo.CreateClassSubject(ctx, classSubject); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.repo.GetClassSubjectByID(ctx, tenantID, classSubject.ID)
}

// UpdateClassSubject updates an existing class-subject mapping.
func (s *Service) UpdateClassSubject(ctx context.Context, tenantID, id uuid.UUID, req UpdateClassSubjectRequest) (*models.ClassSubject, error) {
	classSubject, err := s.repo.GetClassSubjectByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if req.IsMandatory != nil {
		classSubject.IsMandatory = *req.IsMandatory
	}
	if req.PeriodsPerWeek != nil {
		classSubject.PeriodsPerWeek = *req.PeriodsPerWeek
	}
	if req.IsActive != nil {
		classSubject.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateClassSubject(ctx, classSubject); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.repo.GetClassSubjectByID(ctx, tenantID, id)
}

// DeleteClassSubject deletes a class-subject mapping.
func (s *Service) DeleteClassSubject(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if mapping exists
	_, err := s.repo.GetClassSubjectByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteClassSubject(ctx, tenantID, id)
}
