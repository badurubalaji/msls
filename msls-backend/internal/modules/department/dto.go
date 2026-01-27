// Package department provides department management functionality.
package department

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// CreateDepartmentDTO represents a request to create a department.
type CreateDepartmentDTO struct {
	TenantID    uuid.UUID
	BranchID    uuid.UUID
	Name        string
	Code        string
	Description *string
	HeadID      *uuid.UUID
	IsActive    bool
}

// UpdateDepartmentDTO represents a request to update a department.
type UpdateDepartmentDTO struct {
	Name        *string
	Code        *string
	Description *string
	HeadID      *uuid.UUID
	IsActive    *bool
}

// DepartmentResponse represents a department in API responses.
type DepartmentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	HeadID      string `json:"headId,omitempty"`
	HeadName    string `json:"headName,omitempty"`
	BranchID    string `json:"branchId"`
	BranchName  string `json:"branchName,omitempty"`
	IsActive    bool   `json:"isActive"`
	StaffCount  int    `json:"staffCount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// DepartmentListResponse represents a paginated list of departments.
type DepartmentListResponse struct {
	Departments []DepartmentResponse `json:"departments"`
	Total       int64                `json:"total"`
}

// DropdownItem represents a simple item for dropdowns.
type DropdownItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListFilter contains filter options for listing departments.
type ListFilter struct {
	TenantID uuid.UUID
	BranchID *uuid.UUID
	IsActive *bool
	Search   string
}

// ToDepartmentResponse converts a Department model to a DepartmentResponse.
func ToDepartmentResponse(d *models.Department, staffCount int) DepartmentResponse {
	resp := DepartmentResponse{
		ID:         d.ID.String(),
		Name:       d.Name,
		Code:       d.Code,
		BranchID:   d.BranchID.String(),
		IsActive:   d.IsActive,
		StaffCount: staffCount,
		CreatedAt:  d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  d.UpdatedAt.Format(time.RFC3339),
	}

	if d.Description != nil {
		resp.Description = *d.Description
	}

	if d.HeadID != nil {
		resp.HeadID = d.HeadID.String()
	}

	if d.Head != nil {
		resp.HeadName = d.Head.FullName()
	}

	if d.Branch != nil {
		resp.BranchName = d.Branch.Name
	}

	return resp
}

// ToDepartmentResponses converts a slice of Department models to DepartmentResponses.
func ToDepartmentResponses(departments []models.Department, staffCounts map[uuid.UUID]int) []DepartmentResponse {
	responses := make([]DepartmentResponse, len(departments))
	for i, dept := range departments {
		count := 0
		if staffCounts != nil {
			count = staffCounts[dept.ID]
		}
		responses[i] = ToDepartmentResponse(&dept, count)
	}
	return responses
}

// ToDropdownItem converts a Department to a DropdownItem.
func ToDropdownItem(d *models.Department) DropdownItem {
	return DropdownItem{
		ID:   d.ID.String(),
		Name: d.Name,
	}
}

// ToDropdownItems converts departments to dropdown items.
func ToDropdownItems(departments []models.Department) []DropdownItem {
	items := make([]DropdownItem, len(departments))
	for i, dept := range departments {
		items[i] = ToDropdownItem(&dept)
	}
	return items
}
