package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// EvaluationType represents the type of evaluation for an exam
type EvaluationType string

const (
	EvaluationTypeMarks EvaluationType = "marks"
	EvaluationTypeGrade EvaluationType = "grade"
)

// ExamType represents a type of examination (e.g., Unit Test, Half Yearly, Annual)
type ExamType struct {
	ID                  uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID            uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name                string          `gorm:"type:varchar(100);not null" json:"name"`
	Code                string          `gorm:"type:varchar(20);not null" json:"code"`
	Description         *string         `gorm:"type:text" json:"description,omitempty"`
	Weightage           decimal.Decimal `gorm:"type:decimal(5,2);not null;default:0" json:"weightage"`
	EvaluationType      EvaluationType  `gorm:"type:varchar(20);not null;default:'marks'" json:"evaluation_type"`
	DefaultMaxMarks     int             `gorm:"not null;default:100" json:"default_max_marks"`
	DefaultPassingMarks *int            `gorm:"" json:"default_passing_marks,omitempty"`
	DisplayOrder        int             `gorm:"not null;default:0" json:"display_order"`
	IsActive            bool            `gorm:"not null;default:true" json:"is_active"`
	CreatedAt           time.Time       `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt           time.Time       `gorm:"not null;default:now()" json:"updated_at"`
	CreatedBy           *uuid.UUID      `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID      `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relations
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"-"`
}

// TableName returns the table name for ExamType
func (ExamType) TableName() string {
	return "exam_types"
}
