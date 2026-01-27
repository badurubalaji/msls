// Package designation provides designation management functionality.
package designation

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// CreateDesignationDTO represents a request to create a designation.
type CreateDesignationDTO struct {
	TenantID     uuid.UUID
	Name         string
	Level        int
	DepartmentID *uuid.UUID
	IsActive     bool
}

// UpdateDesignationDTO represents a request to update a designation.
type UpdateDesignationDTO struct {
	Name         *string
	Level        *int
	DepartmentID *uuid.UUID
	IsActive     *bool
}

// DesignationResponse represents a designation in API responses.
type DesignationResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Level          int    `json:"level"`
	DepartmentID   string `json:"departmentId,omitempty"`
	DepartmentName string `json:"departmentName,omitempty"`
	IsActive       bool   `json:"isActive"`
	StaffCount     int    `json:"staffCount"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// DesignationListResponse represents a paginated list of designations.
type DesignationListResponse struct {
	Designations []DesignationResponse `json:"designations"`
	Total        int64                 `json:"total"`
}

// DropdownItem represents a simple item for dropdowns.
type DropdownItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

// ListFilter contains filter options for listing designations.
type ListFilter struct {
	TenantID     uuid.UUID
	DepartmentID *uuid.UUID
	IsActive     *bool
	Search       string
}

// ToDesignationResponse converts a Designation model to a DesignationResponse.
func ToDesignationResponse(d *models.Designation, staffCount int) DesignationResponse {
	resp := DesignationResponse{
		ID:         d.ID.String(),
		Name:       d.Name,
		Level:      d.Level,
		IsActive:   d.IsActive,
		StaffCount: staffCount,
		CreatedAt:  d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  d.UpdatedAt.Format(time.RFC3339),
	}

	if d.DepartmentID != nil {
		resp.DepartmentID = d.DepartmentID.String()
	}

	if d.Department != nil {
		resp.DepartmentName = d.Department.Name
	}

	return resp
}

// ToDesignationResponses converts a slice of Designation models to DesignationResponses.
func ToDesignationResponses(designations []models.Designation, staffCounts map[uuid.UUID]int) []DesignationResponse {
	responses := make([]DesignationResponse, len(designations))
	for i, desig := range designations {
		count := 0
		if staffCounts != nil {
			count = staffCounts[desig.ID]
		}
		responses[i] = ToDesignationResponse(&desig, count)
	}
	return responses
}

// ToDropdownItem converts a Designation to a DropdownItem.
func ToDropdownItem(d *models.Designation) DropdownItem {
	return DropdownItem{
		ID:    d.ID.String(),
		Name:  d.Name,
		Level: d.Level,
	}
}

// ToDropdownItems converts designations to dropdown items.
func ToDropdownItems(designations []models.Designation) []DropdownItem {
	items := make([]DropdownItem, len(designations))
	for i, desig := range designations {
		items[i] = ToDropdownItem(&desig)
	}
	return items
}
