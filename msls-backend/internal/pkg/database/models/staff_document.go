// Package models contains database models for the MSLS application.
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// DocumentCategory represents the category of a document type.
type DocumentCategory string

const (
	DocumentCategoryIdentity   DocumentCategory = "identity"
	DocumentCategoryEducation  DocumentCategory = "education"
	DocumentCategoryEmployment DocumentCategory = "employment"
	DocumentCategoryCompliance DocumentCategory = "compliance"
	DocumentCategoryOther      DocumentCategory = "other"
)

// ValidDocumentCategories contains all valid document categories.
var ValidDocumentCategories = []DocumentCategory{
	DocumentCategoryIdentity,
	DocumentCategoryEducation,
	DocumentCategoryEmployment,
	DocumentCategoryCompliance,
	DocumentCategoryOther,
}

// IsValid checks if the category is valid.
func (c DocumentCategory) IsValid() bool {
	for _, valid := range ValidDocumentCategories {
		if c == valid {
			return true
		}
	}
	return false
}

// Note: VerificationStatus type is defined in admission.go and reused here.

// StaffDocumentType represents a configurable document type.
type StaffDocumentType struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index"`

	Name        string           `gorm:"type:varchar(100);not null"`
	Code        string           `gorm:"type:varchar(50);not null"`
	Category    DocumentCategory `gorm:"type:varchar(50);not null"`
	Description *string          `gorm:"type:text"`

	IsMandatory           bool           `gorm:"not null;default:false"`
	HasExpiry             bool           `gorm:"not null;default:false"`
	DefaultValidityMonths *int           `gorm:"type:integer"`
	ApplicableTo          pq.StringArray `gorm:"type:varchar(20)[]"` // teaching, non_teaching
	IsActive              bool           `gorm:"not null;default:true"`
	DisplayOrder          int            `gorm:"default:0"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`

	// Relations
	Tenant Tenant `gorm:"foreignKey:TenantID"`
}

// TableName specifies the table name for StaffDocumentType.
func (StaffDocumentType) TableName() string {
	return "staff_document_types"
}

// StaffDocument represents a document uploaded for a staff member.
type StaffDocument struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index"`
	StaffID        uuid.UUID `gorm:"type:uuid;not null;index"`
	DocumentTypeID uuid.UUID `gorm:"type:uuid;not null;index"`

	DocumentNumber *string    `gorm:"type:varchar(100)"`
	IssueDate      *time.Time `gorm:"type:date"`
	ExpiryDate     *time.Time `gorm:"type:date;index"`

	FileName string `gorm:"type:varchar(255);not null"`
	FilePath string `gorm:"type:text;not null"`
	FileSize int    `gorm:"not null"` // in bytes
	MimeType string `gorm:"type:varchar(100);not null"`

	// Verification
	VerificationStatus VerificationStatus `gorm:"type:varchar(20);not null;default:'pending';index"`
	VerifiedBy         *uuid.UUID         `gorm:"type:uuid"`
	VerifiedAt         *time.Time         `gorm:"type:timestamptz"`
	VerificationNotes  *string            `gorm:"type:text"`
	RejectionReason    *string            `gorm:"type:text"`

	// Metadata
	Remarks   *string `gorm:"type:text"`
	IsCurrent bool    `gorm:"not null;default:true"`

	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Tenant       Tenant             `gorm:"foreignKey:TenantID"`
	Staff        Staff              `gorm:"foreignKey:StaffID"`
	DocumentType *StaffDocumentType `gorm:"foreignKey:DocumentTypeID"`
	Verifier     *User              `gorm:"foreignKey:VerifiedBy"`
	Creator      *User              `gorm:"foreignKey:CreatedBy"`
}

// TableName specifies the table name for StaffDocument.
func (StaffDocument) TableName() string {
	return "staff_documents"
}

// IsExpired checks if the document has expired.
func (d *StaffDocument) IsExpired() bool {
	if d.ExpiryDate == nil {
		return false
	}
	return d.ExpiryDate.Before(time.Now())
}

// IsExpiringSoon checks if the document expires within the given number of days.
func (d *StaffDocument) IsExpiringSoon(days int) bool {
	if d.ExpiryDate == nil {
		return false
	}
	threshold := time.Now().AddDate(0, 0, days)
	return d.ExpiryDate.Before(threshold) && !d.IsExpired()
}

// NotificationType represents the type of document notification.
type NotificationType string

const (
	NotificationTypeExpiry30Days NotificationType = "expiry_30_days"
	NotificationTypeExpiry7Days  NotificationType = "expiry_7_days"
	NotificationTypeExpired      NotificationType = "expired"
)

// StaffDocumentNotification records notifications sent for document expirations.
type StaffDocumentNotification struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID   uuid.UUID `gorm:"type:uuid;not null;index"`
	DocumentID uuid.UUID `gorm:"type:uuid;not null;index"`

	NotificationType NotificationType `gorm:"type:varchar(50);not null"`
	SentAt           time.Time        `gorm:"not null;default:now()"`
	SentTo           pq.StringArray   `gorm:"type:uuid[]"` // User IDs notified

	CreatedAt time.Time `gorm:"not null;default:now()"`

	// Relations
	Tenant   Tenant        `gorm:"foreignKey:TenantID"`
	Document StaffDocument `gorm:"foreignKey:DocumentID"`
}

// TableName specifies the table name for StaffDocumentNotification.
func (StaffDocumentNotification) TableName() string {
	return "staff_document_notifications"
}
