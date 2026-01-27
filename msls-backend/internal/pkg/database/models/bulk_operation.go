// Package models contains database models.
package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BulkOperationType represents the type of bulk operation.
type BulkOperationType string

const (
	BulkOperationTypeSMS          BulkOperationType = "send_sms"
	BulkOperationTypeEmail        BulkOperationType = "send_email"
	BulkOperationTypeUpdateStatus BulkOperationType = "update_status"
	BulkOperationTypeExport       BulkOperationType = "export"
)

// IsValid checks if the bulk operation type is valid.
func (t BulkOperationType) IsValid() bool {
	switch t {
	case BulkOperationTypeSMS, BulkOperationTypeEmail, BulkOperationTypeUpdateStatus, BulkOperationTypeExport:
		return true
	}
	return false
}

// BulkOperationStatus represents the status of a bulk operation.
type BulkOperationStatus string

const (
	BulkOperationStatusPending    BulkOperationStatus = "pending"
	BulkOperationStatusProcessing BulkOperationStatus = "processing"
	BulkOperationStatusCompleted  BulkOperationStatus = "completed"
	BulkOperationStatusFailed     BulkOperationStatus = "failed"
	BulkOperationStatusCancelled  BulkOperationStatus = "cancelled"
)

// IsValid checks if the bulk operation status is valid.
func (s BulkOperationStatus) IsValid() bool {
	switch s {
	case BulkOperationStatusPending, BulkOperationStatusProcessing, BulkOperationStatusCompleted, BulkOperationStatusFailed, BulkOperationStatusCancelled:
		return true
	}
	return false
}

// BulkOperationItemStatus represents the status of a bulk operation item.
type BulkOperationItemStatus string

const (
	BulkItemStatusPending BulkOperationItemStatus = "pending"
	BulkItemStatusSuccess BulkOperationItemStatus = "success"
	BulkItemStatusFailed  BulkOperationItemStatus = "failed"
	BulkItemStatusSkipped BulkOperationItemStatus = "skipped"
)

// BulkOperationParams is a JSON type for operation parameters.
type BulkOperationParams map[string]interface{}

// Value implements driver.Valuer for database storage.
func (p BulkOperationParams) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

// Scan implements sql.Scanner for database retrieval.
func (p *BulkOperationParams) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// BulkOperation represents a bulk operation record.
type BulkOperation struct {
	ID             uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID       uuid.UUID           `gorm:"type:uuid;not null;index"`
	OperationType  BulkOperationType   `gorm:"type:varchar(50);not null"`
	Status         BulkOperationStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	TotalCount     int                 `gorm:"not null;default:0"`
	ProcessedCount int                 `gorm:"not null;default:0"`
	SuccessCount   int                 `gorm:"not null;default:0"`
	FailureCount   int                 `gorm:"not null;default:0"`
	Parameters     BulkOperationParams `gorm:"type:jsonb"`
	ResultURL      string              `gorm:"type:varchar(500)"`
	ErrorMessage   string              `gorm:"type:text"`
	StartedAt      *time.Time          `gorm:"type:timestamptz"`
	CompletedAt    *time.Time          `gorm:"type:timestamptz"`
	CreatedAt      time.Time           `gorm:"type:timestamptz;not null;default:now()"`
	CreatedBy      uuid.UUID           `gorm:"type:uuid;not null"`

	// Associations
	Items []BulkOperationItem `gorm:"foreignKey:OperationID;references:ID"`
}

// TableName returns the table name for BulkOperation.
func (BulkOperation) TableName() string {
	return "bulk_operations"
}

// BulkOperationItem represents an individual item in a bulk operation.
type BulkOperationItem struct {
	ID           uuid.UUID               `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID     uuid.UUID               `gorm:"type:uuid;not null;index"`
	OperationID  uuid.UUID               `gorm:"type:uuid;not null;index"`
	StudentID    uuid.UUID               `gorm:"type:uuid;not null;index"`
	Status       BulkOperationItemStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	ErrorMessage string                  `gorm:"type:text"`
	ProcessedAt  *time.Time              `gorm:"type:timestamptz"`
	CreatedAt    time.Time               `gorm:"type:timestamptz;not null;default:now()"`

	// Associations
	Operation BulkOperation `gorm:"foreignKey:OperationID;references:ID"`
	Student   Student       `gorm:"foreignKey:StudentID;references:ID"`
}

// TableName returns the table name for BulkOperationItem.
func (BulkOperationItem) TableName() string {
	return "bulk_operation_items"
}
