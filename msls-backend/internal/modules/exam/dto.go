package exam

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// ========================================
// Filter DTOs
// ========================================

// ExamTypeFilter contains filter options for listing exam types.
type ExamTypeFilter struct {
	TenantID uuid.UUID
	IsActive *bool
	Search   string
}

// ========================================
// Request DTOs
// ========================================

// CreateExamTypeRequest contains the data for creating an exam type.
type CreateExamTypeRequest struct {
	Name                string                  `json:"name" binding:"required,max=100"`
	Code                string                  `json:"code" binding:"required,max=20"`
	Description         *string                 `json:"description"`
	Weightage           decimal.Decimal         `json:"weightage" binding:"gte=0,lte=100"`
	EvaluationType      models.EvaluationType   `json:"evaluation_type" binding:"required,oneof=marks grade"`
	DefaultMaxMarks     int                     `json:"default_max_marks" binding:"required,gt=0"`
	DefaultPassingMarks *int                    `json:"default_passing_marks" binding:"omitempty,gt=0"`
}

// UpdateExamTypeRequest contains the data for updating an exam type.
type UpdateExamTypeRequest struct {
	Name                *string                 `json:"name" binding:"omitempty,max=100"`
	Code                *string                 `json:"code" binding:"omitempty,max=20"`
	Description         *string                 `json:"description"`
	Weightage           *decimal.Decimal        `json:"weightage" binding:"omitempty,gte=0,lte=100"`
	EvaluationType      *models.EvaluationType  `json:"evaluation_type" binding:"omitempty,oneof=marks grade"`
	DefaultMaxMarks     *int                    `json:"default_max_marks" binding:"omitempty,gt=0"`
	DefaultPassingMarks *int                    `json:"default_passing_marks" binding:"omitempty,gte=0"`
}

// UpdateDisplayOrderRequest contains the data for updating display order.
type UpdateDisplayOrderRequest struct {
	Items []DisplayOrderItem `json:"items" binding:"required,min=1,dive"`
}

// DisplayOrderItem represents a single item in display order update.
type DisplayOrderItem struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	DisplayOrder int       `json:"display_order" binding:"gte=0"`
}

// ToggleActiveRequest contains the data for toggling active status.
type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}

// ========================================
// Response DTOs
// ========================================

// ExamTypeResponse represents an exam type in API responses.
type ExamTypeResponse struct {
	ID                  uuid.UUID               `json:"id"`
	Name                string                  `json:"name"`
	Code                string                  `json:"code"`
	Description         *string                 `json:"description,omitempty"`
	Weightage           decimal.Decimal         `json:"weightage"`
	EvaluationType      models.EvaluationType   `json:"evaluation_type"`
	DefaultMaxMarks     int                     `json:"default_max_marks"`
	DefaultPassingMarks *int                    `json:"default_passing_marks,omitempty"`
	DisplayOrder        int                     `json:"display_order"`
	IsActive            bool                    `json:"is_active"`
	CreatedAt           string                  `json:"created_at"`
	UpdatedAt           string                  `json:"updated_at"`
}

// ExamTypeListResponse represents a list of exam types.
type ExamTypeListResponse struct {
	Items []ExamTypeResponse `json:"items"`
	Total int64              `json:"total"`
}

// ========================================
// Conversion Functions
// ========================================

// ToExamTypeResponse converts a model to response DTO.
func ToExamTypeResponse(et *models.ExamType) ExamTypeResponse {
	return ExamTypeResponse{
		ID:                  et.ID,
		Name:                et.Name,
		Code:                et.Code,
		Description:         et.Description,
		Weightage:           et.Weightage,
		EvaluationType:      et.EvaluationType,
		DefaultMaxMarks:     et.DefaultMaxMarks,
		DefaultPassingMarks: et.DefaultPassingMarks,
		DisplayOrder:        et.DisplayOrder,
		IsActive:            et.IsActive,
		CreatedAt:           et.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:           et.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToExamTypeListResponse converts a slice of models to list response.
func ToExamTypeListResponse(examTypes []models.ExamType, total int64) ExamTypeListResponse {
	items := make([]ExamTypeResponse, len(examTypes))
	for i, et := range examTypes {
		items[i] = ToExamTypeResponse(&et)
	}
	return ExamTypeListResponse{
		Items: items,
		Total: total,
	}
}
