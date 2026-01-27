// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StaffStatus represents the status of a staff member.
type StaffStatus string

// StaffStatus constants.
const (
	StaffStatusActive     StaffStatus = "active"
	StaffStatusInactive   StaffStatus = "inactive"
	StaffStatusTerminated StaffStatus = "terminated"
	StaffStatusOnLeave    StaffStatus = "on_leave"
)

// IsValid checks if the staff status is a valid value.
func (s StaffStatus) IsValid() bool {
	switch s {
	case StaffStatusActive, StaffStatusInactive, StaffStatusTerminated, StaffStatusOnLeave:
		return true
	}
	return false
}

// String returns the string representation of the status.
func (s StaffStatus) String() string {
	return string(s)
}

// StaffType represents the type of staff member.
type StaffType string

// StaffType constants.
const (
	StaffTypeTeaching    StaffType = "teaching"
	StaffTypeNonTeaching StaffType = "non_teaching"
)

// IsValid checks if the staff type is a valid value.
func (t StaffType) IsValid() bool {
	switch t {
	case StaffTypeTeaching, StaffTypeNonTeaching:
		return true
	}
	return false
}

// Staff validation errors.
var (
	ErrStaffTenantIDRequired    = errors.New("tenant ID is required")
	ErrStaffBranchIDRequired    = errors.New("branch ID is required")
	ErrStaffFirstNameRequired   = errors.New("first name is required")
	ErrStaffLastNameRequired    = errors.New("last name is required")
	ErrStaffDOBRequired         = errors.New("date of birth is required")
	ErrStaffGenderRequired      = errors.New("gender is required")
	ErrStaffInvalidGender       = errors.New("invalid gender value")
	ErrStaffInvalidStatus       = errors.New("invalid staff status")
	ErrStaffInvalidStaffType    = errors.New("invalid staff type")
	ErrStaffWorkEmailRequired   = errors.New("work email is required")
	ErrStaffWorkPhoneRequired   = errors.New("work phone is required")
	ErrStaffJoinDateRequired    = errors.New("join date is required")
)

// Staff represents a staff member in the system.
type Staff struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenantId"`
	BranchID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"branchId"`
	
	// Employee identification
	EmployeeID       string         `gorm:"type:varchar(50);not null" json:"employeeId"`
	EmployeeIDPrefix string         `gorm:"type:varchar(10);not null;default:'EMP'" json:"employeeIdPrefix"`
	
	// Personal details
	FirstName        string         `gorm:"type:varchar(100);not null" json:"firstName"`
	MiddleName       string         `gorm:"type:varchar(100)" json:"middleName,omitempty"`
	LastName         string         `gorm:"type:varchar(100);not null" json:"lastName"`
	DateOfBirth      time.Time      `gorm:"type:date;not null" json:"dateOfBirth"`
	Gender           Gender         `gorm:"type:varchar(20);not null" json:"gender"`
	BloodGroup       string         `gorm:"type:varchar(10)" json:"bloodGroup,omitempty"`
	Nationality      string         `gorm:"type:varchar(50);default:'Indian'" json:"nationality,omitempty"`
	Religion         string         `gorm:"type:varchar(50)" json:"religion,omitempty"`
	MaritalStatus    string         `gorm:"type:varchar(20)" json:"maritalStatus,omitempty"`
	
	// Contact details
	PersonalEmail           string `gorm:"type:varchar(255)" json:"personalEmail,omitempty"`
	WorkEmail               string `gorm:"type:varchar(255);not null" json:"workEmail"`
	PersonalPhone           string `gorm:"type:varchar(20)" json:"personalPhone,omitempty"`
	WorkPhone               string `gorm:"type:varchar(20);not null" json:"workPhone"`
	EmergencyContactName    string `gorm:"type:varchar(200)" json:"emergencyContactName,omitempty"`
	EmergencyContactPhone   string `gorm:"type:varchar(20)" json:"emergencyContactPhone,omitempty"`
	EmergencyContactRelation string `gorm:"type:varchar(50)" json:"emergencyContactRelation,omitempty"`
	
	// Current Address
	CurrentAddressLine1 string `gorm:"type:varchar(255)" json:"currentAddressLine1,omitempty"`
	CurrentAddressLine2 string `gorm:"type:varchar(255)" json:"currentAddressLine2,omitempty"`
	CurrentCity         string `gorm:"type:varchar(100)" json:"currentCity,omitempty"`
	CurrentState        string `gorm:"type:varchar(100)" json:"currentState,omitempty"`
	CurrentPincode      string `gorm:"type:varchar(10)" json:"currentPincode,omitempty"`
	CurrentCountry      string `gorm:"type:varchar(100);default:'India'" json:"currentCountry,omitempty"`
	
	// Permanent Address
	PermanentAddressLine1 string `gorm:"type:varchar(255)" json:"permanentAddressLine1,omitempty"`
	PermanentAddressLine2 string `gorm:"type:varchar(255)" json:"permanentAddressLine2,omitempty"`
	PermanentCity         string `gorm:"type:varchar(100)" json:"permanentCity,omitempty"`
	PermanentState        string `gorm:"type:varchar(100)" json:"permanentState,omitempty"`
	PermanentPincode      string `gorm:"type:varchar(10)" json:"permanentPincode,omitempty"`
	PermanentCountry      string `gorm:"type:varchar(100);default:'India'" json:"permanentCountry,omitempty"`
	SameAsCurrent         bool   `gorm:"default:false" json:"sameAsCurrent"`
	
	// Employment details
	StaffType          StaffType  `gorm:"type:varchar(20);not null" json:"staffType"`
	DepartmentID       *uuid.UUID `gorm:"type:uuid;index" json:"departmentId,omitempty"`
	DesignationID      *uuid.UUID `gorm:"type:uuid;index" json:"designationId,omitempty"`
	ReportingManagerID *uuid.UUID `gorm:"type:uuid" json:"reportingManagerId,omitempty"`
	JoinDate           time.Time  `gorm:"type:date;not null" json:"joinDate"`
	ConfirmationDate   *time.Time `gorm:"type:date" json:"confirmationDate,omitempty"`
	ProbationEndDate   *time.Time `gorm:"type:date" json:"probationEndDate,omitempty"`
	
	// Status
	Status          StaffStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	StatusReason    string      `gorm:"type:text" json:"statusReason,omitempty"`
	TerminationDate *time.Time  `gorm:"type:date" json:"terminationDate,omitempty"`
	
	// Profile
	PhotoURL string `gorm:"type:varchar(500)" json:"photoUrl,omitempty"`
	Bio      string `gorm:"type:text" json:"bio,omitempty"`
	
	// Audit
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy *uuid.UUID     `gorm:"type:uuid" json:"createdBy,omitempty"`
	UpdatedBy *uuid.UUID     `gorm:"type:uuid" json:"updatedBy,omitempty"`
	Version   int            `gorm:"not null;default:1" json:"version"`
	
	// Relationships
	Tenant           Tenant       `gorm:"foreignKey:TenantID" json:"-"`
	Branch           Branch       `gorm:"foreignKey:BranchID" json:"-"`
	Department       *Department  `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Designation      *Designation `gorm:"foreignKey:DesignationID" json:"designation,omitempty"`
	ReportingManager *Staff       `gorm:"foreignKey:ReportingManagerID" json:"reportingManager,omitempty"`
}

// TableName returns the table name for the Staff model.
func (Staff) TableName() string {
	return "staff"
}

// BeforeCreate hook for Staff.
func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.Status == "" {
		s.Status = StaffStatusActive
	}
	if s.EmployeeIDPrefix == "" {
		s.EmployeeIDPrefix = "EMP"
	}
	return s.Validate()
}

// BeforeUpdate hook for Staff to increment version for optimistic locking.
func (s *Staff) BeforeUpdate(tx *gorm.DB) error {
	s.Version++
	return nil
}

// Validate performs validation on the Staff model.
func (s *Staff) Validate() error {
	if s.TenantID == uuid.Nil {
		return ErrStaffTenantIDRequired
	}
	if s.BranchID == uuid.Nil {
		return ErrStaffBranchIDRequired
	}
	if s.FirstName == "" {
		return ErrStaffFirstNameRequired
	}
	if s.LastName == "" {
		return ErrStaffLastNameRequired
	}
	if s.DateOfBirth.IsZero() {
		return ErrStaffDOBRequired
	}
	if s.Gender == "" {
		return ErrStaffGenderRequired
	}
	if !s.Gender.IsValid() {
		return ErrStaffInvalidGender
	}
	if s.WorkEmail == "" {
		return ErrStaffWorkEmailRequired
	}
	if s.WorkPhone == "" {
		return ErrStaffWorkPhoneRequired
	}
	if !s.StaffType.IsValid() {
		return ErrStaffInvalidStaffType
	}
	if s.JoinDate.IsZero() {
		return ErrStaffJoinDateRequired
	}
	if s.Status != "" && !s.Status.IsValid() {
		return ErrStaffInvalidStatus
	}
	return nil
}

// FullName returns the full name of the staff member.
func (s *Staff) FullName() string {
	if s.MiddleName != "" {
		return s.FirstName + " " + s.MiddleName + " " + s.LastName
	}
	return s.FirstName + " " + s.LastName
}

// GetInitials returns the initials of the staff member's name.
func (s *Staff) GetInitials() string {
	initials := ""
	if len(s.FirstName) > 0 {
		initials += string(s.FirstName[0])
	}
	if len(s.LastName) > 0 {
		initials += string(s.LastName[0])
	}
	return initials
}

// StaffEmployeeSequence tracks the employee ID sequence per tenant/prefix.
type StaffEmployeeSequence struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index" json:"tenantId"`
	Prefix       string    `gorm:"type:varchar(10);not null;default:'EMP'" json:"prefix"`
	LastSequence int       `gorm:"not null;default:0" json:"lastSequence"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updatedAt"`
}

// TableName returns the table name for the StaffEmployeeSequence model.
func (StaffEmployeeSequence) TableName() string {
	return "staff_employee_sequences"
}

// StaffStatusHistory tracks status changes for staff members.
type StaffStatusHistory struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenantId"`
	StaffID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"staffId"`
	OldStatus     string     `gorm:"type:varchar(20)" json:"oldStatus,omitempty"`
	NewStatus     string     `gorm:"type:varchar(20);not null" json:"newStatus"`
	Reason        string     `gorm:"type:text" json:"reason,omitempty"`
	EffectiveDate time.Time  `gorm:"type:date;not null" json:"effectiveDate"`
	ChangedBy     *uuid.UUID `gorm:"type:uuid" json:"changedBy,omitempty"`
	ChangedAt     time.Time  `gorm:"not null;default:now()" json:"changedAt"`
	
	// Relationships
	Staff Staff `gorm:"foreignKey:StaffID" json:"-"`
}

// TableName returns the table name for the StaffStatusHistory model.
func (StaffStatusHistory) TableName() string {
	return "staff_status_history"
}
