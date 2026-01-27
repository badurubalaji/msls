// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentStatus represents the status of a student.
type StudentStatus string

// StudentStatus constants.
const (
	StudentStatusActive      StudentStatus = "active"
	StudentStatusInactive    StudentStatus = "inactive"
	StudentStatusTransferred StudentStatus = "transferred"
	StudentStatusGraduated   StudentStatus = "graduated"
)

// IsValid checks if the student status is a valid value.
func (s StudentStatus) IsValid() bool {
	switch s {
	case StudentStatusActive, StudentStatusInactive, StudentStatusTransferred, StudentStatusGraduated:
		return true
	}
	return false
}

// String returns the string representation of the status.
func (s StudentStatus) String() string {
	return string(s)
}

// Gender represents the gender of a student.
type Gender string

// Gender constants.
const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// IsValid checks if the gender is a valid value.
func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderOther:
		return true
	}
	return false
}

// AddressType represents the type of address.
type AddressType string

// AddressType constants.
const (
	AddressTypeCurrent   AddressType = "current"
	AddressTypePermanent AddressType = "permanent"
)

// IsValid checks if the address type is a valid value.
func (a AddressType) IsValid() bool {
	switch a {
	case AddressTypeCurrent, AddressTypePermanent:
		return true
	}
	return false
}

// Student validation errors.
var (
	ErrStudentTenantIDRequired    = errors.New("tenant ID is required")
	ErrStudentBranchIDRequired    = errors.New("branch ID is required")
	ErrStudentFirstNameRequired   = errors.New("first name is required")
	ErrStudentLastNameRequired    = errors.New("last name is required")
	ErrStudentDOBRequired         = errors.New("date of birth is required")
	ErrStudentGenderRequired      = errors.New("gender is required")
	ErrStudentInvalidGender       = errors.New("invalid gender value")
	ErrStudentInvalidStatus       = errors.New("invalid student status")
	ErrStudentInvalidAddressType  = errors.New("invalid address type")
	ErrStudentAddressLine1Required = errors.New("address line 1 is required")
	ErrStudentCityRequired        = errors.New("city is required")
	ErrStudentStateRequired       = errors.New("state is required")
	ErrStudentPostalCodeRequired  = errors.New("postal code is required")
)

// Student represents a student in the system.
type Student struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenantId"`
	BranchID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"branchId"`
	AdmissionNumber     string         `gorm:"type:varchar(20);not null" json:"admissionNumber"`
	FirstName           string         `gorm:"type:varchar(100);not null" json:"firstName"`
	MiddleName          string         `gorm:"type:varchar(100)" json:"middleName,omitempty"`
	LastName            string         `gorm:"type:varchar(100);not null" json:"lastName"`
	DateOfBirth         time.Time      `gorm:"type:date;not null" json:"dateOfBirth"`
	Gender              Gender         `gorm:"type:varchar(10);not null" json:"gender"`
	BloodGroup          string         `gorm:"type:varchar(5)" json:"bloodGroup,omitempty"`
	AadhaarNumber       string         `gorm:"type:varchar(12)" json:"aadhaarNumber,omitempty"`
	PhotoURL            string         `gorm:"type:varchar(500)" json:"photoUrl,omitempty"`
	BirthCertificateURL string         `gorm:"type:varchar(500)" json:"birthCertificateUrl,omitempty"`
	Status              StudentStatus  `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	AdmissionDate       time.Time      `gorm:"type:date;not null;default:CURRENT_DATE" json:"admissionDate"`
	CreatedAt           time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt           time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy           *uuid.UUID     `gorm:"type:uuid" json:"createdBy,omitempty"`
	UpdatedBy           *uuid.UUID     `gorm:"type:uuid" json:"updatedBy,omitempty"`
	Version             int            `gorm:"not null;default:1" json:"version"`

	// Relationships
	Tenant    Tenant            `gorm:"foreignKey:TenantID" json:"-"`
	Branch    Branch            `gorm:"foreignKey:BranchID" json:"-"`
	Addresses []StudentAddress  `gorm:"foreignKey:StudentID" json:"addresses,omitempty"`
}

// TableName returns the table name for the Student model.
func (Student) TableName() string {
	return "students"
}

// BeforeCreate hook for Student.
func (s *Student) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.Status == "" {
		s.Status = StudentStatusActive
	}
	if s.AdmissionDate.IsZero() {
		s.AdmissionDate = time.Now()
	}
	return s.Validate()
}

// BeforeUpdate hook for Student to increment version for optimistic locking.
func (s *Student) BeforeUpdate(tx *gorm.DB) error {
	s.Version++
	return nil
}

// Validate performs validation on the Student model.
func (s *Student) Validate() error {
	if s.TenantID == uuid.Nil {
		return ErrStudentTenantIDRequired
	}
	if s.BranchID == uuid.Nil {
		return ErrStudentBranchIDRequired
	}
	if s.FirstName == "" {
		return ErrStudentFirstNameRequired
	}
	if s.LastName == "" {
		return ErrStudentLastNameRequired
	}
	if s.DateOfBirth.IsZero() {
		return ErrStudentDOBRequired
	}
	if s.Gender == "" {
		return ErrStudentGenderRequired
	}
	if !s.Gender.IsValid() {
		return ErrStudentInvalidGender
	}
	if s.Status != "" && !s.Status.IsValid() {
		return ErrStudentInvalidStatus
	}
	return nil
}

// FullName returns the full name of the student.
func (s *Student) FullName() string {
	if s.MiddleName != "" {
		return s.FirstName + " " + s.MiddleName + " " + s.LastName
	}
	return s.FirstName + " " + s.LastName
}

// GetInitials returns the initials of the student's name.
func (s *Student) GetInitials() string {
	initials := ""
	if len(s.FirstName) > 0 {
		initials += string(s.FirstName[0])
	}
	if len(s.LastName) > 0 {
		initials += string(s.LastName[0])
	}
	return initials
}

// StudentAddress represents an address associated with a student.
type StudentAddress struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"tenantId"`
	StudentID    uuid.UUID   `gorm:"type:uuid;not null;index" json:"studentId"`
	AddressType  AddressType `gorm:"type:varchar(20);not null" json:"addressType"`
	AddressLine1 string      `gorm:"type:varchar(255);not null" json:"addressLine1"`
	AddressLine2 string      `gorm:"type:varchar(255)" json:"addressLine2,omitempty"`
	City         string      `gorm:"type:varchar(100);not null" json:"city"`
	State        string      `gorm:"type:varchar(100);not null" json:"state"`
	PostalCode   string      `gorm:"type:varchar(10);not null" json:"postalCode"`
	Country      string      `gorm:"type:varchar(100);not null;default:'India'" json:"country"`
	CreatedAt    time.Time   `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Student Student `gorm:"foreignKey:StudentID" json:"-"`
}

// TableName returns the table name for the StudentAddress model.
func (StudentAddress) TableName() string {
	return "student_addresses"
}

// BeforeCreate hook for StudentAddress.
func (a *StudentAddress) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Country == "" {
		a.Country = "India"
	}
	return a.Validate()
}

// Validate performs validation on the StudentAddress model.
func (a *StudentAddress) Validate() error {
	if a.TenantID == uuid.Nil {
		return ErrStudentTenantIDRequired
	}
	if a.StudentID == uuid.Nil {
		return errors.New("student ID is required")
	}
	if !a.AddressType.IsValid() {
		return ErrStudentInvalidAddressType
	}
	if a.AddressLine1 == "" {
		return ErrStudentAddressLine1Required
	}
	if a.City == "" {
		return ErrStudentCityRequired
	}
	if a.State == "" {
		return ErrStudentStateRequired
	}
	if a.PostalCode == "" {
		return ErrStudentPostalCodeRequired
	}
	return nil
}

// StudentAdmissionSequence tracks the admission number sequence per branch/year.
type StudentAdmissionSequence struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index" json:"tenantId"`
	BranchID     uuid.UUID `gorm:"type:uuid;not null" json:"branchId"`
	Year         int       `gorm:"not null" json:"year"`
	LastSequence int       `gorm:"not null;default:0" json:"lastSequence"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updatedAt"`
}

// TableName returns the table name for the StudentAdmissionSequence model.
func (StudentAdmissionSequence) TableName() string {
	return "student_admission_sequences"
}
