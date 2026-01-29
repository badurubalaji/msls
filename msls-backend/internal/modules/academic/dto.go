// Package academic provides academic structure management functionality.
package academic

import (
	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// ========================================
// Class DTOs
// ========================================

// ClassResponse represents a class in API responses.
type ClassResponse struct {
	ID           uuid.UUID         `json:"id"`
	BranchID     uuid.UUID         `json:"branchId"`
	BranchName   string            `json:"branchName,omitempty"`
	Name         string            `json:"name"`
	Code         string            `json:"code"`
	Level        string            `json:"level,omitempty"`
	DisplayOrder int               `json:"displayOrder"`
	Description  string            `json:"description,omitempty"`
	HasStreams   bool              `json:"hasStreams"`
	IsActive     bool              `json:"isActive"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
	Sections     []SectionResponse `json:"sections,omitempty"`
	Streams      []StreamResponse  `json:"streams,omitempty"`
}

// ClassListResponse represents the response for listing classes.
type ClassListResponse struct {
	Classes []ClassResponse `json:"classes"`
	Total   int64           `json:"total"`
}

// CreateClassRequest represents the request body for creating a class.
type CreateClassRequest struct {
	BranchID     uuid.UUID   `json:"branchId" binding:"required"`
	Name         string      `json:"name" binding:"required,max=50"`
	Code         string      `json:"code" binding:"required,max=20"`
	Level        string      `json:"level" binding:"omitempty,oneof=nursery primary middle secondary senior_secondary"`
	DisplayOrder int         `json:"displayOrder"`
	Description  string      `json:"description"`
	HasStreams   bool        `json:"hasStreams"`
	StreamIDs    []uuid.UUID `json:"streamIds"`
}

// UpdateClassRequest represents the request body for updating a class.
type UpdateClassRequest struct {
	Name         *string     `json:"name" binding:"omitempty,max=50"`
	Code         *string     `json:"code" binding:"omitempty,max=20"`
	Level        *string     `json:"level" binding:"omitempty,oneof=nursery primary middle secondary senior_secondary"`
	DisplayOrder *int        `json:"displayOrder"`
	Description  *string     `json:"description"`
	HasStreams   *bool       `json:"hasStreams"`
	IsActive     *bool       `json:"isActive"`
	StreamIDs    []uuid.UUID `json:"streamIds"`
}

// ClassFilter represents filters for listing classes.
type ClassFilter struct {
	TenantID   uuid.UUID
	BranchID   *uuid.UUID
	Level      *string
	IsActive   *bool
	HasStreams *bool
	Search     string
}

// ToResponse converts a Class model to ClassResponse.
func ClassToResponse(c *models.Class) ClassResponse {
	resp := ClassResponse{
		ID:           c.ID,
		BranchID:     c.BranchID,
		Name:         c.Name,
		Code:         c.Code,
		DisplayOrder: c.DisplayOrder,
		Description:  c.Description,
		HasStreams:   c.HasStreams,
		IsActive:     c.IsActive,
		CreatedAt:    c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if c.Level != nil {
		resp.Level = c.Level.String()
	}

	if c.Branch.ID != uuid.Nil {
		resp.BranchName = c.Branch.Name
	}

	if len(c.Sections) > 0 {
		resp.Sections = make([]SectionResponse, len(c.Sections))
		for i, s := range c.Sections {
			resp.Sections[i] = SectionToResponse(&s)
		}
	}

	if len(c.Streams) > 0 {
		resp.Streams = make([]StreamResponse, len(c.Streams))
		for i, s := range c.Streams {
			resp.Streams[i] = StreamToResponse(&s)
		}
	}

	return resp
}

// ========================================
// Section DTOs
// ========================================

// SectionResponse represents a section in API responses.
type SectionResponse struct {
	ID               uuid.UUID `json:"id"`
	ClassID          uuid.UUID `json:"classId"`
	ClassName        string    `json:"className,omitempty"`
	AcademicYearID   *string   `json:"academicYearId,omitempty"`
	AcademicYearName string    `json:"academicYearName,omitempty"`
	StreamID         *string   `json:"streamId,omitempty"`
	StreamName       string    `json:"streamName,omitempty"`
	ClassTeacherID   *string   `json:"classTeacherId,omitempty"`
	ClassTeacherName string    `json:"classTeacherName,omitempty"`
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	Capacity         int       `json:"capacity"`
	RoomNumber       string    `json:"roomNumber,omitempty"`
	DisplayOrder     int       `json:"displayOrder"`
	IsActive         bool      `json:"isActive"`
	StudentCount     int       `json:"studentCount"`
	CreatedAt        string    `json:"createdAt"`
	UpdatedAt        string    `json:"updatedAt"`
}

// SectionListResponse represents the response for listing sections.
type SectionListResponse struct {
	Sections []SectionResponse `json:"sections"`
	Total    int64             `json:"total"`
}

// CreateSectionRequest represents the request body for creating a section.
type CreateSectionRequest struct {
	ClassID        uuid.UUID  `json:"classId" binding:"required"`
	AcademicYearID *uuid.UUID `json:"academicYearId"`
	StreamID       *uuid.UUID `json:"streamId"`
	ClassTeacherID *uuid.UUID `json:"classTeacherId"`
	Name           string     `json:"name" binding:"required,max=20"`
	Code           string     `json:"code" binding:"required,max=20"`
	Capacity       int        `json:"capacity"`
	RoomNumber     string     `json:"roomNumber"`
}

// UpdateSectionRequest represents the request body for updating a section.
type UpdateSectionRequest struct {
	AcademicYearID *uuid.UUID `json:"academicYearId"`
	StreamID       *uuid.UUID `json:"streamId"`
	ClassTeacherID *uuid.UUID `json:"classTeacherId"`
	Name           *string    `json:"name" binding:"omitempty,max=20"`
	Code           *string    `json:"code" binding:"omitempty,max=20"`
	Capacity       *int       `json:"capacity"`
	RoomNumber     *string    `json:"roomNumber"`
	DisplayOrder   *int       `json:"displayOrder"`
	IsActive       *bool      `json:"isActive"`
}

// SectionFilter represents filters for listing sections.
type SectionFilter struct {
	TenantID       uuid.UUID
	ClassID        *uuid.UUID
	AcademicYearID *uuid.UUID
	StreamID       *uuid.UUID
	IsActive       *bool
	Search         string
}

// SectionToResponse converts a Section model to SectionResponse.
func SectionToResponse(s *models.Section) SectionResponse {
	resp := SectionResponse{
		ID:           s.ID,
		ClassID:      s.ClassID,
		Name:         s.Name,
		Code:         s.Code,
		DisplayOrder: s.DisplayOrder,
		IsActive:     s.IsActive,
		StudentCount: s.StudentCount,
		CreatedAt:    s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if s.Capacity != nil {
		resp.Capacity = *s.Capacity
	} else {
		resp.Capacity = 40
	}

	resp.RoomNumber = s.RoomNumber

	if s.Class.ID != uuid.Nil {
		resp.ClassName = s.Class.Name
	}

	if s.AcademicYearID != nil {
		idStr := s.AcademicYearID.String()
		resp.AcademicYearID = &idStr
		if s.AcademicYear != nil && s.AcademicYear.ID != uuid.Nil {
			resp.AcademicYearName = s.AcademicYear.Name
		}
	}

	if s.StreamID != nil {
		idStr := s.StreamID.String()
		resp.StreamID = &idStr
		if s.Stream != nil && s.Stream.ID != uuid.Nil {
			resp.StreamName = s.Stream.Name
		}
	}

	if s.ClassTeacherID != nil {
		idStr := s.ClassTeacherID.String()
		resp.ClassTeacherID = &idStr
		if s.ClassTeacher != nil {
			resp.ClassTeacherName = s.ClassTeacher.FirstName
			if s.ClassTeacher.LastName != "" {
				resp.ClassTeacherName += " " + s.ClassTeacher.LastName
			}
		}
	}

	return resp
}

// ========================================
// Stream DTOs
// ========================================

// StreamResponse represents a stream in API responses.
type StreamResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description,omitempty"`
	DisplayOrder int       `json:"displayOrder"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

// StreamListResponse represents the response for listing streams.
type StreamListResponse struct {
	Streams []StreamResponse `json:"streams"`
	Total   int64            `json:"total"`
}

// CreateStreamRequest represents the request body for creating a stream.
type CreateStreamRequest struct {
	Name         string `json:"name" binding:"required,max=100"`
	Code         string `json:"code" binding:"required,max=20"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"displayOrder"`
}

// UpdateStreamRequest represents the request body for updating a stream.
type UpdateStreamRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=100"`
	Code         *string `json:"code" binding:"omitempty,max=20"`
	Description  *string `json:"description"`
	DisplayOrder *int    `json:"displayOrder"`
	IsActive     *bool   `json:"isActive"`
}

// StreamFilter represents filters for listing streams.
type StreamFilter struct {
	TenantID uuid.UUID
	IsActive *bool
	Search   string
}

// StreamToResponse converts a Stream model to StreamResponse.
func StreamToResponse(s *models.Stream) StreamResponse {
	return StreamResponse{
		ID:           s.ID,
		Name:         s.Name,
		Code:         s.Code,
		Description:  s.Description,
		DisplayOrder: s.DisplayOrder,
		IsActive:     s.IsActive,
		CreatedAt:    s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ========================================
// Subject DTOs
// ========================================

// SubjectResponse represents a subject in API responses.
type SubjectResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	ShortName    string    `json:"shortName,omitempty"`
	Description  string    `json:"description,omitempty"`
	SubjectType  string    `json:"subjectType"`
	MaxMarks     int       `json:"maxMarks"`
	PassingMarks int       `json:"passingMarks"`
	CreditHours  float64   `json:"creditHours"`
	DisplayOrder int       `json:"displayOrder"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

// SubjectListResponse represents the response for listing subjects.
type SubjectListResponse struct {
	Subjects []SubjectResponse `json:"subjects"`
	Total    int64             `json:"total"`
}

// CreateSubjectRequest represents the request body for creating a subject.
type CreateSubjectRequest struct {
	Name         string  `json:"name" binding:"required,max=100"`
	Code         string  `json:"code" binding:"required,max=20"`
	ShortName    string  `json:"shortName" binding:"max=20"`
	Description  string  `json:"description"`
	SubjectType  string  `json:"subjectType" binding:"required,oneof=core elective language co_curricular vocational"`
	MaxMarks     int     `json:"maxMarks"`
	PassingMarks int     `json:"passingMarks"`
	CreditHours  float64 `json:"creditHours"`
	DisplayOrder int     `json:"displayOrder"`
}

// UpdateSubjectRequest represents the request body for updating a subject.
type UpdateSubjectRequest struct {
	Name         *string  `json:"name" binding:"omitempty,max=100"`
	Code         *string  `json:"code" binding:"omitempty,max=20"`
	ShortName    *string  `json:"shortName" binding:"omitempty,max=20"`
	Description  *string  `json:"description"`
	SubjectType  *string  `json:"subjectType" binding:"omitempty,oneof=core elective language co_curricular vocational"`
	MaxMarks     *int     `json:"maxMarks"`
	PassingMarks *int     `json:"passingMarks"`
	CreditHours  *float64 `json:"creditHours"`
	DisplayOrder *int     `json:"displayOrder"`
	IsActive     *bool    `json:"isActive"`
}

// SubjectFilter represents filters for listing subjects.
type SubjectFilter struct {
	TenantID    uuid.UUID
	SubjectType *string
	IsActive    *bool
	Search      string
}

// SubjectToResponse converts a Subject model to SubjectResponse.
func SubjectToResponse(s *models.Subject) SubjectResponse {
	resp := SubjectResponse{
		ID:           s.ID,
		Name:         s.Name,
		Code:         s.Code,
		ShortName:    s.ShortName,
		Description:  s.Description,
		SubjectType:  string(s.SubjectType),
		DisplayOrder: s.DisplayOrder,
		IsActive:     s.IsActive,
		CreatedAt:    s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if s.MaxMarks != nil {
		resp.MaxMarks = *s.MaxMarks
	} else {
		resp.MaxMarks = 100
	}

	if s.PassingMarks != nil {
		resp.PassingMarks = *s.PassingMarks
	} else {
		resp.PassingMarks = 35
	}

	if s.CreditHours != nil {
		resp.CreditHours = *s.CreditHours
	}

	return resp
}

// ========================================
// Class-Subject DTOs
// ========================================

// ClassSubjectResponse represents a class-subject mapping in API responses.
type ClassSubjectResponse struct {
	ID             uuid.UUID `json:"id"`
	ClassID        uuid.UUID `json:"classId"`
	ClassName      string    `json:"className,omitempty"`
	SubjectID      uuid.UUID `json:"subjectId"`
	SubjectName    string    `json:"subjectName,omitempty"`
	SubjectCode    string    `json:"subjectCode,omitempty"`
	IsMandatory    bool      `json:"isMandatory"`
	PeriodsPerWeek int       `json:"periodsPerWeek"`
	IsActive       bool      `json:"isActive"`
	CreatedAt      string    `json:"createdAt"`
	UpdatedAt      string    `json:"updatedAt"`
}

// ClassSubjectListResponse represents the response for listing class-subject mappings.
type ClassSubjectListResponse struct {
	ClassSubjects []ClassSubjectResponse `json:"classSubjects"`
	Total         int64                  `json:"total"`
}

// CreateClassSubjectRequest represents the request body for creating a class-subject mapping.
type CreateClassSubjectRequest struct {
	SubjectID      uuid.UUID `json:"subjectId" binding:"required"`
	IsMandatory    bool      `json:"isMandatory"`
	PeriodsPerWeek int       `json:"periodsPerWeek"`
}

// UpdateClassSubjectRequest represents the request body for updating a class-subject mapping.
type UpdateClassSubjectRequest struct {
	IsMandatory    *bool `json:"isMandatory"`
	PeriodsPerWeek *int  `json:"periodsPerWeek"`
	IsActive       *bool `json:"isActive"`
}

// ClassSubjectFilter represents filters for listing class-subject mappings.
type ClassSubjectFilter struct {
	TenantID  uuid.UUID
	ClassID   *uuid.UUID
	SubjectID *uuid.UUID
	IsActive  *bool
}

// ClassSubjectToResponse converts a ClassSubject model to ClassSubjectResponse.
func ClassSubjectToResponse(cs *models.ClassSubject) ClassSubjectResponse {
	resp := ClassSubjectResponse{
		ID:             cs.ID,
		ClassID:        cs.ClassID,
		SubjectID:      cs.SubjectID,
		IsMandatory:    cs.IsMandatory,
		PeriodsPerWeek: cs.PeriodsPerWeek,
		IsActive:       cs.IsActive,
		CreatedAt:      cs.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      cs.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if cs.Class.ID != uuid.Nil {
		resp.ClassName = cs.Class.Name
	}

	if cs.Subject.ID != uuid.Nil {
		resp.SubjectName = cs.Subject.Name
		resp.SubjectCode = cs.Subject.Code
	}

	return resp
}

// ========================================
// Structure DTOs
// ========================================

// ClassStructureResponse represents the hierarchical class-section structure.
type ClassStructureResponse struct {
	Classes []ClassWithSectionsResponse `json:"classes"`
}

// ClassWithSectionsResponse represents a class with its sections for structure view.
type ClassWithSectionsResponse struct {
	ID           uuid.UUID                  `json:"id"`
	Name         string                     `json:"name"`
	Code         string                     `json:"code"`
	DisplayOrder int                        `json:"displayOrder"`
	HasStreams   bool                       `json:"hasStreams"`
	IsActive     bool                       `json:"isActive"`
	Sections     []SectionStructureResponse `json:"sections"`
	TotalStudents int                       `json:"totalStudents"`
	TotalCapacity int                       `json:"totalCapacity"`
}

// SectionStructureResponse represents a section in the structure view.
type SectionStructureResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	Capacity         int       `json:"capacity"`
	StudentCount     int       `json:"studentCount"`
	ClassTeacherID   *string   `json:"classTeacherId,omitempty"`
	ClassTeacherName string    `json:"classTeacherName,omitempty"`
	StreamName       string    `json:"streamName,omitempty"`
	RoomNumber       string    `json:"roomNumber,omitempty"`
	CapacityUsage    float64   `json:"capacityUsage"` // Percentage
}
