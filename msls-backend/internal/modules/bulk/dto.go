// Package bulk provides bulk operation functionality.
package bulk

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// CreateBulkOperationDTO represents a request to create a bulk operation.
type CreateBulkOperationDTO struct {
	TenantID      uuid.UUID
	OperationType models.BulkOperationType
	StudentIDs    []uuid.UUID
	Parameters    models.BulkOperationParams
	CreatedBy     uuid.UUID
}

// BulkStatusUpdateParams contains parameters for bulk status update.
type BulkStatusUpdateParams struct {
	NewStatus models.StudentStatus `json:"newStatus"`
}

// ExportParams contains parameters for export operation.
type ExportParams struct {
	Format  string   `json:"format"`  // xlsx or csv
	Columns []string `json:"columns"` // columns to export
}

// BulkOperationResponse represents a bulk operation in API responses.
type BulkOperationResponse struct {
	ID             string                      `json:"id"`
	OperationType  string                      `json:"operationType"`
	Status         string                      `json:"status"`
	TotalCount     int                         `json:"totalCount"`
	ProcessedCount int                         `json:"processedCount"`
	SuccessCount   int                         `json:"successCount"`
	FailureCount   int                         `json:"failureCount"`
	ResultURL      string                      `json:"resultUrl,omitempty"`
	ErrorMessage   string                      `json:"errorMessage,omitempty"`
	StartedAt      string                      `json:"startedAt,omitempty"`
	CompletedAt    string                      `json:"completedAt,omitempty"`
	CreatedAt      string                      `json:"createdAt"`
	Items          []BulkOperationItemResponse `json:"items,omitempty"`
}

// BulkOperationItemResponse represents a bulk operation item in API responses.
type BulkOperationItemResponse struct {
	ID           string `json:"id"`
	StudentID    string `json:"studentId"`
	StudentName  string `json:"studentName,omitempty"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ProcessedAt  string `json:"processedAt,omitempty"`
}

// ToBulkOperationResponse converts a BulkOperation model to a response.
func ToBulkOperationResponse(op *models.BulkOperation) BulkOperationResponse {
	resp := BulkOperationResponse{
		ID:             op.ID.String(),
		OperationType:  string(op.OperationType),
		Status:         string(op.Status),
		TotalCount:     op.TotalCount,
		ProcessedCount: op.ProcessedCount,
		SuccessCount:   op.SuccessCount,
		FailureCount:   op.FailureCount,
		ResultURL:      op.ResultURL,
		ErrorMessage:   op.ErrorMessage,
		CreatedAt:      op.CreatedAt.Format(time.RFC3339),
	}

	if op.StartedAt != nil {
		resp.StartedAt = op.StartedAt.Format(time.RFC3339)
	}
	if op.CompletedAt != nil {
		resp.CompletedAt = op.CompletedAt.Format(time.RFC3339)
	}

	// Convert items
	if len(op.Items) > 0 {
		resp.Items = make([]BulkOperationItemResponse, len(op.Items))
		for i, item := range op.Items {
			resp.Items[i] = ToBulkOperationItemResponse(&item)
		}
	}

	return resp
}

// ToBulkOperationItemResponse converts a BulkOperationItem model to a response.
func ToBulkOperationItemResponse(item *models.BulkOperationItem) BulkOperationItemResponse {
	resp := BulkOperationItemResponse{
		ID:           item.ID.String(),
		StudentID:    item.StudentID.String(),
		Status:       string(item.Status),
		ErrorMessage: item.ErrorMessage,
	}

	if item.ProcessedAt != nil {
		resp.ProcessedAt = item.ProcessedAt.Format(time.RFC3339)
	}

	// Add student name if loaded
	if item.Student.ID != uuid.Nil {
		resp.StudentName = item.Student.FullName()
	}

	return resp
}

// DefaultExportColumns are the default columns for student export.
var DefaultExportColumns = []string{
	"admission_number",
	"first_name",
	"last_name",
	"class",
	"section",
	"phone",
	"guardian_name",
	"guardian_phone",
	"status",
}

// ExportColumnLabels maps column keys to display labels.
var ExportColumnLabels = map[string]string{
	"admission_number": "Admission Number",
	"first_name":       "First Name",
	"last_name":        "Last Name",
	"full_name":        "Full Name",
	"class":            "Class",
	"section":          "Section",
	"roll_number":      "Roll Number",
	"gender":           "Gender",
	"date_of_birth":    "Date of Birth",
	"blood_group":      "Blood Group",
	"aadhaar_number":   "Aadhaar Number",
	"admission_date":   "Admission Date",
	"phone":            "Phone",
	"email":            "Email",
	"guardian_name":    "Guardian Name",
	"guardian_phone":   "Guardian Phone",
	"guardian_email":   "Guardian Email",
	"guardian_relation": "Guardian Relation",
	"address":          "Address",
	"city":             "City",
	"state":            "State",
	"status":           "Status",
	"branch":           "Branch",
}
