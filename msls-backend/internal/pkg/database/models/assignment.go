// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AssignmentStatus represents the status of a teacher assignment.
type AssignmentStatus string

// AssignmentStatus constants.
const (
	AssignmentStatusActive   AssignmentStatus = "active"
	AssignmentStatusInactive AssignmentStatus = "inactive"
)

// IsValid checks if the assignment status is a valid value.
func (s AssignmentStatus) IsValid() bool {
	switch s {
	case AssignmentStatusActive, AssignmentStatusInactive:
		return true
	}
	return false
}

// String returns the string representation of the status.
func (s AssignmentStatus) String() string {
	return string(s)
}

// TeacherSubjectAssignment represents a teacher's assignment to a subject and class.
type TeacherSubjectAssignment struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"tenantId"`
	StaffID        uuid.UUID        `gorm:"type:uuid;not null;index" json:"staffId"`
	SubjectID      uuid.UUID        `gorm:"type:uuid;not null;index" json:"subjectId"`
	ClassID        uuid.UUID        `gorm:"type:uuid;not null;index" json:"classId"`
	SectionID      *uuid.UUID       `gorm:"type:uuid;index" json:"sectionId,omitempty"`
	AcademicYearID uuid.UUID        `gorm:"type:uuid;not null;index" json:"academicYearId"`

	PeriodsPerWeek int              `gorm:"not null;default:0" json:"periodsPerWeek"`
	IsClassTeacher bool             `gorm:"not null;default:false" json:"isClassTeacher"`

	EffectiveFrom time.Time         `gorm:"type:date;not null" json:"effectiveFrom"`
	EffectiveTo   *time.Time        `gorm:"type:date" json:"effectiveTo,omitempty"`

	Status  AssignmentStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	Remarks string           `gorm:"type:text" json:"remarks,omitempty"`

	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"createdBy,omitempty"`

	// Relationships
	Staff        Staff        `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	Subject      Subject      `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Class        Class        `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Section      *Section     `gorm:"foreignKey:SectionID" json:"section,omitempty"`
	AcademicYear AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academicYear,omitempty"`
}

// TableName returns the table name for the TeacherSubjectAssignment model.
func (TeacherSubjectAssignment) TableName() string {
	return "teacher_subject_assignments"
}

// BeforeCreate hook for TeacherSubjectAssignment.
func (a *TeacherSubjectAssignment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Status == "" {
		a.Status = AssignmentStatusActive
	}
	return nil
}

// TeacherWorkloadSettings represents workload configuration per branch.
type TeacherWorkloadSettings struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenantId"`
	BranchID uuid.UUID `gorm:"type:uuid;not null;index" json:"branchId"`

	MinPeriodsPerWeek    int `gorm:"not null;default:20" json:"minPeriodsPerWeek"`
	MaxPeriodsPerWeek    int `gorm:"not null;default:35" json:"maxPeriodsPerWeek"`
	MaxSubjectsPerTeacher *int `gorm:"default:5" json:"maxSubjectsPerTeacher,omitempty"`
	MaxClassesPerTeacher  *int `gorm:"default:8" json:"maxClassesPerTeacher,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Branch Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

// TableName returns the table name for the TeacherWorkloadSettings model.
func (TeacherWorkloadSettings) TableName() string {
	return "teacher_workload_settings"
}

// BeforeCreate hook for TeacherWorkloadSettings.
func (w *TeacherWorkloadSettings) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// ClassLevel represents the level of a class.
type ClassLevel string

// ClassLevel constants.
const (
	ClassLevelNursery         ClassLevel = "nursery"
	ClassLevelPrimary         ClassLevel = "primary"
	ClassLevelMiddle          ClassLevel = "middle"
	ClassLevelSecondary       ClassLevel = "secondary"
	ClassLevelSeniorSecondary ClassLevel = "senior_secondary"
)

// IsValid checks if the class level is a valid value.
func (l ClassLevel) IsValid() bool {
	switch l {
	case ClassLevelNursery, ClassLevelPrimary, ClassLevelMiddle, ClassLevelSecondary, ClassLevelSeniorSecondary:
		return true
	}
	return false
}

// String returns the string representation of the class level.
func (l ClassLevel) String() string {
	return string(l)
}

// Class represents a class in the academic structure.
type Class struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"tenantId"`
	BranchID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"branchId"`
	Name         string      `gorm:"type:varchar(50);not null" json:"name"`
	Code         string      `gorm:"type:varchar(20);not null" json:"code"`
	Level        *ClassLevel `gorm:"type:varchar(30)" json:"level,omitempty"`
	DisplayOrder int         `gorm:"not null;default:0" json:"displayOrder"`
	Description  string      `gorm:"type:text" json:"description,omitempty"`
	HasStreams   bool        `gorm:"not null;default:false" json:"hasStreams"`
	IsActive     bool        `gorm:"not null;default:true" json:"isActive"`
	CreatedAt    time.Time   `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy    *uuid.UUID  `gorm:"type:uuid" json:"createdBy,omitempty"`

	// Relationships
	Branch   Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Sections []Section `gorm:"foreignKey:ClassID" json:"sections,omitempty"`
	Streams  []Stream  `gorm:"many2many:class_streams;foreignKey:ID;joinForeignKey:ClassID;References:ID;joinReferences:StreamID" json:"streams,omitempty"`
}

// TableName returns the table name for the Class model.
func (Class) TableName() string {
	return "classes"
}

// Section represents a section within a class.
type Section struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenantId"`
	ClassID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"classId"`
	AcademicYearID  *uuid.UUID `gorm:"type:uuid;index" json:"academicYearId,omitempty"`
	StreamID        *uuid.UUID `gorm:"type:uuid;index" json:"streamId,omitempty"`
	ClassTeacherID  *uuid.UUID `gorm:"type:uuid;index" json:"classTeacherId,omitempty"`
	Name            string     `gorm:"type:varchar(20);not null" json:"name"`
	Code            string     `gorm:"type:varchar(20);not null" json:"code"`
	Capacity        *int       `gorm:"default:40" json:"capacity,omitempty"`
	RoomNumber      string     `gorm:"type:varchar(50)" json:"roomNumber,omitempty"`
	DisplayOrder    int        `gorm:"not null;default:0" json:"displayOrder"`
	IsActive        bool       `gorm:"not null;default:true" json:"isActive"`
	CreatedAt       time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy       *uuid.UUID `gorm:"type:uuid" json:"createdBy,omitempty"`

	// Relationships
	Class        Class         `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	AcademicYear *AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academicYear,omitempty"`
	Stream       *Stream       `gorm:"foreignKey:StreamID" json:"stream,omitempty"`
	ClassTeacher *Staff        `gorm:"foreignKey:ClassTeacherID" json:"classTeacher,omitempty"`

	// Computed fields (not stored)
	StudentCount int `gorm:"-" json:"studentCount,omitempty"`
}

// TableName returns the table name for the Section model.
func (Section) TableName() string {
	return "sections"
}

// SubjectType represents the type of subject.
type SubjectType string

// SubjectType constants.
const (
	SubjectTypeCore        SubjectType = "core"
	SubjectTypeElective    SubjectType = "elective"
	SubjectTypeLanguage    SubjectType = "language"
	SubjectTypeCoCurricular SubjectType = "co_curricular"
	SubjectTypeVocational  SubjectType = "vocational"
)

// IsValid checks if the subject type is a valid value.
func (t SubjectType) IsValid() bool {
	switch t {
	case SubjectTypeCore, SubjectTypeElective, SubjectTypeLanguage, SubjectTypeCoCurricular, SubjectTypeVocational:
		return true
	}
	return false
}

// Subject represents a subject in the academic structure.
type Subject struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"tenantId"`
	Name         string      `gorm:"type:varchar(100);not null" json:"name"`
	Code         string      `gorm:"type:varchar(20);not null" json:"code"`
	ShortName    string      `gorm:"type:varchar(20)" json:"shortName,omitempty"`
	Description  string      `gorm:"type:text" json:"description,omitempty"`
	SubjectType  SubjectType `gorm:"type:varchar(30);not null;default:'core'" json:"subjectType"`
	MaxMarks     *int        `gorm:"default:100" json:"maxMarks,omitempty"`
	PassingMarks *int        `gorm:"default:35" json:"passingMarks,omitempty"`
	CreditHours  *float64    `gorm:"type:decimal(4,2);default:0" json:"creditHours,omitempty"`
	IsActive     bool        `gorm:"not null;default:true" json:"isActive"`
	DisplayOrder int         `gorm:"not null;default:0" json:"displayOrder"`
	CreatedAt    time.Time   `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy    *uuid.UUID  `gorm:"type:uuid" json:"createdBy,omitempty"`
}

// TableName returns the table name for the Subject model.
func (Subject) TableName() string {
	return "subjects"
}

// ClassSubject represents the mapping between classes and subjects.
type ClassSubject struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID      uuid.UUID `gorm:"type:uuid;not null;index" json:"tenantId"`
	ClassID       uuid.UUID `gorm:"type:uuid;not null;index" json:"classId"`
	SubjectID     uuid.UUID `gorm:"type:uuid;not null;index" json:"subjectId"`
	IsMandatory   bool      `gorm:"not null;default:true" json:"isMandatory"`
	PeriodsPerWeek int      `gorm:"not null;default:5" json:"periodsPerWeek"`
	IsActive      bool      `gorm:"not null;default:true" json:"isActive"`
	CreatedAt     time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Class   Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Subject Subject `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
}

// TableName returns the table name for the ClassSubject model.
func (ClassSubject) TableName() string {
	return "class_subjects"
}

// Stream represents an academic stream (Science, Commerce, Arts).
type Stream struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenantId"`
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`
	Code         string     `gorm:"type:varchar(20);not null" json:"code"`
	Description  string     `gorm:"type:text" json:"description,omitempty"`
	DisplayOrder int        `gorm:"not null;default:0" json:"displayOrder"`
	IsActive     bool       `gorm:"not null;default:true" json:"isActive"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy    *uuid.UUID `gorm:"type:uuid" json:"createdBy,omitempty"`
}

// TableName returns the table name for the Stream model.
func (Stream) TableName() string {
	return "streams"
}

// ClassStream represents the mapping between classes and streams.
type ClassStream struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenantId"`
	ClassID   uuid.UUID `gorm:"type:uuid;not null;index" json:"classId"`
	StreamID  uuid.UUID `gorm:"type:uuid;not null;index" json:"streamId"`
	IsActive  bool      `gorm:"not null;default:true" json:"isActive"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`

	// Relationships
	Class  Class  `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Stream Stream `gorm:"foreignKey:StreamID" json:"stream,omitempty"`
}

// TableName returns the table name for the ClassStream model.
func (ClassStream) TableName() string {
	return "class_streams"
}
