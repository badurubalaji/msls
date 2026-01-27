// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// GuardianRelation represents the relation of a guardian to a student.
type GuardianRelation string

// GuardianRelation constants.
const (
	GuardianRelationFather       GuardianRelation = "father"
	GuardianRelationMother       GuardianRelation = "mother"
	GuardianRelationGrandfather  GuardianRelation = "grandfather"
	GuardianRelationGrandmother  GuardianRelation = "grandmother"
	GuardianRelationUncle        GuardianRelation = "uncle"
	GuardianRelationAunt         GuardianRelation = "aunt"
	GuardianRelationSibling      GuardianRelation = "sibling"
	GuardianRelationGuardian     GuardianRelation = "guardian"
	GuardianRelationOther        GuardianRelation = "other"
)

// IsValid checks if the guardian relation is a valid value.
func (r GuardianRelation) IsValid() bool {
	switch r {
	case GuardianRelationFather, GuardianRelationMother,
		GuardianRelationGrandfather, GuardianRelationGrandmother,
		GuardianRelationUncle, GuardianRelationAunt,
		GuardianRelationSibling, GuardianRelationGuardian, GuardianRelationOther:
		return true
	}
	return false
}

// String returns the string representation of the relation.
func (r GuardianRelation) String() string {
	return string(r)
}

// Guardian validation errors.
var (
	ErrGuardianTenantIDRequired   = errors.New("tenant ID is required")
	ErrGuardianStudentIDRequired  = errors.New("student ID is required")
	ErrGuardianRelationRequired   = errors.New("relation is required")
	ErrGuardianInvalidRelation    = errors.New("invalid guardian relation")
	ErrGuardianFirstNameRequired  = errors.New("first name is required")
	ErrGuardianLastNameRequired   = errors.New("last name is required")
	ErrGuardianPhoneRequired      = errors.New("phone number is required")
	ErrEmergencyNameRequired      = errors.New("emergency contact name is required")
	ErrEmergencyRelationRequired  = errors.New("emergency contact relation is required")
	ErrEmergencyPhoneRequired     = errors.New("emergency contact phone is required")
	ErrEmergencyInvalidPriority   = errors.New("priority must be between 1 and 5")
)

// StudentGuardian represents a guardian associated with a student.
type StudentGuardian struct {
	ID              uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID        uuid.UUID        `gorm:"type:uuid;not null;index" json:"tenantId"`
	StudentID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"studentId"`
	Relation        GuardianRelation `gorm:"type:varchar(20);not null" json:"relation"`
	FirstName       string           `gorm:"type:varchar(100);not null" json:"firstName"`
	LastName        string           `gorm:"type:varchar(100);not null" json:"lastName"`
	Phone           string           `gorm:"type:varchar(15);not null" json:"phone"`
	Email           string           `gorm:"type:varchar(255)" json:"email,omitempty"`
	Occupation      string           `gorm:"type:varchar(100)" json:"occupation,omitempty"`
	AnnualIncome    decimal.Decimal  `gorm:"type:decimal(15,2)" json:"annualIncome,omitempty"`
	Education       string           `gorm:"type:varchar(100)" json:"education,omitempty"`
	IsPrimary       bool             `gorm:"not null;default:false" json:"isPrimary"`
	HasPortalAccess bool             `gorm:"not null;default:false" json:"hasPortalAccess"`
	UserID          *uuid.UUID       `gorm:"type:uuid" json:"userId,omitempty"`
	AddressLine1    string           `gorm:"type:varchar(255)" json:"addressLine1,omitempty"`
	AddressLine2    string           `gorm:"type:varchar(255)" json:"addressLine2,omitempty"`
	City            string           `gorm:"type:varchar(100)" json:"city,omitempty"`
	State           string           `gorm:"type:varchar(100)" json:"state,omitempty"`
	PostalCode      string           `gorm:"type:varchar(10)" json:"postalCode,omitempty"`
	Country         string           `gorm:"type:varchar(100);default:'India'" json:"country,omitempty"`
	CreatedAt       time.Time        `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt       time.Time        `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy       *uuid.UUID       `gorm:"type:uuid" json:"createdBy,omitempty"`
	UpdatedBy       *uuid.UUID       `gorm:"type:uuid" json:"updatedBy,omitempty"`

	// Relationships
	Student Student `gorm:"foreignKey:StudentID" json:"-"`
	User    *User   `gorm:"foreignKey:UserID" json:"-"`
}

// TableName returns the table name for the StudentGuardian model.
func (StudentGuardian) TableName() string {
	return "student_guardians"
}

// BeforeCreate hook for StudentGuardian.
func (g *StudentGuardian) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	if g.Country == "" {
		g.Country = "India"
	}
	return g.Validate()
}

// Validate performs validation on the StudentGuardian model.
func (g *StudentGuardian) Validate() error {
	if g.TenantID == uuid.Nil {
		return ErrGuardianTenantIDRequired
	}
	if g.StudentID == uuid.Nil {
		return ErrGuardianStudentIDRequired
	}
	if g.Relation == "" {
		return ErrGuardianRelationRequired
	}
	if !g.Relation.IsValid() {
		return ErrGuardianInvalidRelation
	}
	if g.FirstName == "" {
		return ErrGuardianFirstNameRequired
	}
	if g.LastName == "" {
		return ErrGuardianLastNameRequired
	}
	if g.Phone == "" {
		return ErrGuardianPhoneRequired
	}
	return nil
}

// FullName returns the full name of the guardian.
func (g *StudentGuardian) FullName() string {
	return g.FirstName + " " + g.LastName
}

// StudentEmergencyContact represents an emergency contact for a student.
type StudentEmergencyContact struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenantId"`
	StudentID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"studentId"`
	Name           string     `gorm:"type:varchar(200);not null" json:"name"`
	Relation       string     `gorm:"type:varchar(50);not null" json:"relation"`
	Phone          string     `gorm:"type:varchar(15);not null" json:"phone"`
	AlternatePhone string     `gorm:"type:varchar(15)" json:"alternatePhone,omitempty"`
	Priority       int        `gorm:"not null;default:1" json:"priority"`
	Notes          string     `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	CreatedBy      *uuid.UUID `gorm:"type:uuid" json:"createdBy,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:uuid" json:"updatedBy,omitempty"`

	// Relationships
	Student Student `gorm:"foreignKey:StudentID" json:"-"`
}

// TableName returns the table name for the StudentEmergencyContact model.
func (StudentEmergencyContact) TableName() string {
	return "student_emergency_contacts"
}

// BeforeCreate hook for StudentEmergencyContact.
func (e *StudentEmergencyContact) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.Priority == 0 {
		e.Priority = 1
	}
	return e.Validate()
}

// Validate performs validation on the StudentEmergencyContact model.
func (e *StudentEmergencyContact) Validate() error {
	if e.TenantID == uuid.Nil {
		return ErrGuardianTenantIDRequired
	}
	if e.StudentID == uuid.Nil {
		return ErrGuardianStudentIDRequired
	}
	if e.Name == "" {
		return ErrEmergencyNameRequired
	}
	if e.Relation == "" {
		return ErrEmergencyRelationRequired
	}
	if e.Phone == "" {
		return ErrEmergencyPhoneRequired
	}
	if e.Priority < 1 || e.Priority > 5 {
		return ErrEmergencyInvalidPriority
	}
	return nil
}
