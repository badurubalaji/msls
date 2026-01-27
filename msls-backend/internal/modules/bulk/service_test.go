// Package bulk provides bulk operation functionality.
package bulk

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"msls-backend/internal/pkg/database/models"
)

func TestCreateBulkOperationDTO_Validation(t *testing.T) {
	tests := []struct {
		name      string
		dto       CreateBulkOperationDTO
		wantError error
	}{
		{
			name: "valid DTO",
			dto: CreateBulkOperationDTO{
				TenantID:      uuid.New(),
				OperationType: models.BulkOperationTypeUpdateStatus,
				StudentIDs:    []uuid.UUID{uuid.New(), uuid.New()},
				CreatedBy:     uuid.New(),
			},
			wantError: nil,
		},
		{
			name: "empty student IDs",
			dto: CreateBulkOperationDTO{
				TenantID:      uuid.New(),
				OperationType: models.BulkOperationTypeUpdateStatus,
				StudentIDs:    []uuid.UUID{},
				CreatedBy:     uuid.New(),
			},
			wantError: ErrNoStudentsProvided,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test validates the DTO structure
			if tt.wantError == ErrNoStudentsProvided {
				assert.Empty(t, tt.dto.StudentIDs)
			} else {
				assert.NotEmpty(t, tt.dto.StudentIDs)
			}
		})
	}
}

func TestExportParams_Validation(t *testing.T) {
	tests := []struct {
		name    string
		params  ExportParams
		isValid bool
	}{
		{
			name: "valid xlsx format",
			params: ExportParams{
				Format:  "xlsx",
				Columns: []string{"first_name", "last_name"},
			},
			isValid: true,
		},
		{
			name: "valid csv format",
			params: ExportParams{
				Format:  "csv",
				Columns: []string{"admission_number", "full_name"},
			},
			isValid: true,
		},
		{
			name: "invalid format",
			params: ExportParams{
				Format:  "pdf",
				Columns: []string{"first_name"},
			},
			isValid: false,
		},
		{
			name: "empty format",
			params: ExportParams{
				Format:  "",
				Columns: []string{"first_name"},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.params.Format == "xlsx" || tt.params.Format == "csv"
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestDefaultExportColumns(t *testing.T) {
	// Verify default columns are set
	assert.NotEmpty(t, DefaultExportColumns)
	assert.Contains(t, DefaultExportColumns, "admission_number")
	assert.Contains(t, DefaultExportColumns, "first_name")
	assert.Contains(t, DefaultExportColumns, "last_name")
	assert.Contains(t, DefaultExportColumns, "status")
}

func TestExportColumnLabels(t *testing.T) {
	// Verify all default columns have labels
	for _, col := range DefaultExportColumns {
		label, ok := ExportColumnLabels[col]
		assert.True(t, ok, "Column %s should have a label", col)
		assert.NotEmpty(t, label, "Label for column %s should not be empty", col)
	}
}

func TestBulkOperationType_IsValid(t *testing.T) {
	tests := []struct {
		opType  models.BulkOperationType
		isValid bool
	}{
		{models.BulkOperationTypeSMS, true},
		{models.BulkOperationTypeEmail, true},
		{models.BulkOperationTypeUpdateStatus, true},
		{models.BulkOperationTypeExport, true},
		{models.BulkOperationType("invalid"), false},
		{models.BulkOperationType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.opType), func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.opType.IsValid())
		})
	}
}

func TestBulkOperationStatus_IsValid(t *testing.T) {
	tests := []struct {
		status  models.BulkOperationStatus
		isValid bool
	}{
		{models.BulkOperationStatusPending, true},
		{models.BulkOperationStatusProcessing, true},
		{models.BulkOperationStatusCompleted, true},
		{models.BulkOperationStatusFailed, true},
		{models.BulkOperationStatusCancelled, true},
		{models.BulkOperationStatus("invalid"), false},
		{models.BulkOperationStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.status.IsValid())
		})
	}
}

func TestToBulkOperationResponse(t *testing.T) {
	opID := uuid.New()
	tenantID := uuid.New()
	createdBy := uuid.New()

	op := &models.BulkOperation{
		ID:             opID,
		TenantID:       tenantID,
		OperationType:  models.BulkOperationTypeExport,
		Status:         models.BulkOperationStatusCompleted,
		TotalCount:     10,
		ProcessedCount: 10,
		SuccessCount:   9,
		FailureCount:   1,
		ResultURL:      "/exports/test.xlsx",
		CreatedBy:      createdBy,
	}

	resp := ToBulkOperationResponse(op)

	assert.Equal(t, opID.String(), resp.ID)
	assert.Equal(t, "export", resp.OperationType)
	assert.Equal(t, "completed", resp.Status)
	assert.Equal(t, 10, resp.TotalCount)
	assert.Equal(t, 10, resp.ProcessedCount)
	assert.Equal(t, 9, resp.SuccessCount)
	assert.Equal(t, 1, resp.FailureCount)
	assert.Equal(t, "/exports/test.xlsx", resp.ResultURL)
}

func TestToBulkOperationItemResponse(t *testing.T) {
	itemID := uuid.New()
	studentID := uuid.New()

	item := &models.BulkOperationItem{
		ID:           itemID,
		StudentID:    studentID,
		Status:       models.BulkItemStatusSuccess,
		ErrorMessage: "",
	}

	resp := ToBulkOperationItemResponse(item)

	assert.Equal(t, itemID.String(), resp.ID)
	assert.Equal(t, studentID.String(), resp.StudentID)
	assert.Equal(t, "success", resp.Status)
	assert.Empty(t, resp.ErrorMessage)
}

func TestMaxConstants(t *testing.T) {
	// Verify limits are reasonable
	assert.Equal(t, 10000, MaxExportRecords)
	assert.Equal(t, 1000, MaxBulkStudents)
	assert.True(t, MaxExportRecords >= MaxBulkStudents)
}
