// Package salary provides salary management functionality.
package salary

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// ========================================
// Salary Component DTOs
// ========================================

// CreateComponentDTO represents a request to create a salary component.
type CreateComponentDTO struct {
	TenantID        uuid.UUID
	Name            string
	Code            string
	Description     *string
	ComponentType   models.ComponentType
	CalculationType models.CalculationType
	PercentageOfID  *uuid.UUID
	IsTaxable       bool
	DisplayOrder    int
}

// UpdateComponentDTO represents a request to update a salary component.
type UpdateComponentDTO struct {
	Name            *string
	Code            *string
	Description     *string
	ComponentType   *models.ComponentType
	CalculationType *models.CalculationType
	PercentageOfID  *uuid.UUID
	IsTaxable       *bool
	IsActive        *bool
	DisplayOrder    *int
}

// ComponentResponse represents a salary component in API responses.
type ComponentResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Code            string `json:"code"`
	Description     string `json:"description,omitempty"`
	ComponentType   string `json:"componentType"`
	CalculationType string `json:"calculationType"`
	PercentageOfID  string `json:"percentageOfId,omitempty"`
	PercentageOfName string `json:"percentageOfName,omitempty"`
	IsTaxable       bool   `json:"isTaxable"`
	IsActive        bool   `json:"isActive"`
	DisplayOrder    int    `json:"displayOrder"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// ComponentListResponse represents a list of salary components.
type ComponentListResponse struct {
	Components []ComponentResponse `json:"components"`
	Total      int64               `json:"total"`
}

// ComponentFilter contains filter options for listing components.
type ComponentFilter struct {
	TenantID      uuid.UUID
	ComponentType *models.ComponentType
	IsActive      *bool
	Search        string
}

// ToComponentResponse converts a SalaryComponent model to a response.
func ToComponentResponse(c *models.SalaryComponent) ComponentResponse {
	resp := ComponentResponse{
		ID:              c.ID.String(),
		Name:            c.Name,
		Code:            c.Code,
		ComponentType:   string(c.ComponentType),
		CalculationType: string(c.CalculationType),
		IsTaxable:       c.IsTaxable,
		IsActive:        c.IsActive,
		DisplayOrder:    c.DisplayOrder,
		CreatedAt:       c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       c.UpdatedAt.Format(time.RFC3339),
	}

	if c.Description != nil {
		resp.Description = *c.Description
	}

	if c.PercentageOfID != nil {
		resp.PercentageOfID = c.PercentageOfID.String()
	}

	if c.PercentageOf != nil {
		resp.PercentageOfName = c.PercentageOf.Name
	}

	return resp
}

// ToComponentResponses converts a slice of SalaryComponent models to responses.
func ToComponentResponses(components []models.SalaryComponent) []ComponentResponse {
	responses := make([]ComponentResponse, len(components))
	for i, comp := range components {
		responses[i] = ToComponentResponse(&comp)
	}
	return responses
}

// ========================================
// Salary Structure DTOs
// ========================================

// CreateStructureDTO represents a request to create a salary structure.
type CreateStructureDTO struct {
	TenantID      uuid.UUID
	Name          string
	Code          string
	Description   *string
	DesignationID *uuid.UUID
	Components    []StructureComponentDTO
}

// StructureComponentDTO represents a component in a structure.
type StructureComponentDTO struct {
	ComponentID uuid.UUID
	Amount      *decimal.Decimal
	Percentage  *decimal.Decimal
}

// UpdateStructureDTO represents a request to update a salary structure.
type UpdateStructureDTO struct {
	Name          *string
	Code          *string
	Description   *string
	DesignationID *uuid.UUID
	IsActive      *bool
	Components    []StructureComponentDTO
}

// StructureResponse represents a salary structure in API responses.
type StructureResponse struct {
	ID              string                       `json:"id"`
	Name            string                       `json:"name"`
	Code            string                       `json:"code"`
	Description     string                       `json:"description,omitempty"`
	DesignationID   string                       `json:"designationId,omitempty"`
	DesignationName string                       `json:"designationName,omitempty"`
	IsActive        bool                         `json:"isActive"`
	ComponentCount  int                          `json:"componentCount"`
	StaffCount      int                          `json:"staffCount"`
	Components      []StructureComponentResponse `json:"components,omitempty"`
	CreatedAt       string                       `json:"createdAt"`
	UpdatedAt       string                       `json:"updatedAt"`
}

// StructureComponentResponse represents a component in a structure response.
type StructureComponentResponse struct {
	ID            string  `json:"id"`
	ComponentID   string  `json:"componentId"`
	ComponentName string  `json:"componentName"`
	ComponentCode string  `json:"componentCode"`
	ComponentType string  `json:"componentType"`
	Amount        *string `json:"amount,omitempty"`
	Percentage    *string `json:"percentage,omitempty"`
}

// StructureListResponse represents a list of salary structures.
type StructureListResponse struct {
	Structures []StructureResponse `json:"structures"`
	Total      int64               `json:"total"`
}

// StructureFilter contains filter options for listing structures.
type StructureFilter struct {
	TenantID      uuid.UUID
	DesignationID *uuid.UUID
	IsActive      *bool
	Search        string
}

// ToStructureResponse converts a SalaryStructure model to a response.
func ToStructureResponse(s *models.SalaryStructure, staffCount int, includeComponents bool) StructureResponse {
	resp := StructureResponse{
		ID:             s.ID.String(),
		Name:           s.Name,
		Code:           s.Code,
		IsActive:       s.IsActive,
		ComponentCount: len(s.Components),
		StaffCount:     staffCount,
		CreatedAt:      s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      s.UpdatedAt.Format(time.RFC3339),
	}

	if s.Description != nil {
		resp.Description = *s.Description
	}

	if s.DesignationID != nil {
		resp.DesignationID = s.DesignationID.String()
	}

	if s.Designation != nil {
		resp.DesignationName = s.Designation.Name
	}

	if includeComponents && len(s.Components) > 0 {
		resp.Components = make([]StructureComponentResponse, len(s.Components))
		for i, comp := range s.Components {
			resp.Components[i] = ToStructureComponentResponse(&comp)
		}
	}

	return resp
}

// ToStructureComponentResponse converts a SalaryStructureComponent to a response.
func ToStructureComponentResponse(sc *models.SalaryStructureComponent) StructureComponentResponse {
	resp := StructureComponentResponse{
		ID:          sc.ID.String(),
		ComponentID: sc.ComponentID.String(),
	}

	if sc.Component != nil {
		resp.ComponentName = sc.Component.Name
		resp.ComponentCode = sc.Component.Code
		resp.ComponentType = string(sc.Component.ComponentType)
	}

	if sc.Amount != nil {
		amt := sc.Amount.StringFixed(2)
		resp.Amount = &amt
	}

	if sc.Percentage != nil {
		pct := sc.Percentage.StringFixed(2)
		resp.Percentage = &pct
	}

	return resp
}

// ToStructureResponses converts a slice of SalaryStructure models to responses.
func ToStructureResponses(structures []models.SalaryStructure, staffCounts map[uuid.UUID]int) []StructureResponse {
	responses := make([]StructureResponse, len(structures))
	for i, s := range structures {
		count := 0
		if staffCounts != nil {
			count = staffCounts[s.ID]
		}
		responses[i] = ToStructureResponse(&s, count, false)
	}
	return responses
}

// ========================================
// Staff Salary DTOs
// ========================================

// AssignSalaryDTO represents a request to assign salary to staff.
type AssignSalaryDTO struct {
	TenantID      uuid.UUID
	StaffID       uuid.UUID
	StructureID   *uuid.UUID
	EffectiveFrom time.Time
	Components    []StaffComponentDTO
	RevisionReason *string
	CreatedBy     *uuid.UUID
}

// StaffComponentDTO represents a component value for staff salary.
type StaffComponentDTO struct {
	ComponentID  uuid.UUID
	Amount       decimal.Decimal
	IsOverridden bool
}

// StaffSalaryResponse represents a staff salary in API responses.
type StaffSalaryResponse struct {
	ID             string                       `json:"id"`
	StaffID        string                       `json:"staffId"`
	StaffName      string                       `json:"staffName,omitempty"`
	StructureID    string                       `json:"structureId,omitempty"`
	StructureName  string                       `json:"structureName,omitempty"`
	EffectiveFrom  string                       `json:"effectiveFrom"`
	EffectiveTo    string                       `json:"effectiveTo,omitempty"`
	GrossSalary    string                       `json:"grossSalary"`
	NetSalary      string                       `json:"netSalary"`
	CTC            string                       `json:"ctc,omitempty"`
	RevisionReason string                       `json:"revisionReason,omitempty"`
	IsCurrent      bool                         `json:"isCurrent"`
	Components     []StaffSalaryComponentResponse `json:"components,omitempty"`
	CreatedAt      string                       `json:"createdAt"`
	UpdatedAt      string                       `json:"updatedAt"`
}

// StaffSalaryComponentResponse represents a component in staff salary response.
type StaffSalaryComponentResponse struct {
	ID            string `json:"id"`
	ComponentID   string `json:"componentId"`
	ComponentName string `json:"componentName"`
	ComponentCode string `json:"componentCode"`
	ComponentType string `json:"componentType"`
	Amount        string `json:"amount"`
	IsOverridden  bool   `json:"isOverridden"`
}

// StaffSalaryHistoryResponse represents salary history.
type StaffSalaryHistoryResponse struct {
	History []StaffSalaryResponse `json:"history"`
	Total   int64                 `json:"total"`
}

// ToStaffSalaryResponse converts a StaffSalary model to a response.
func ToStaffSalaryResponse(s *models.StaffSalary, includeComponents bool) StaffSalaryResponse {
	resp := StaffSalaryResponse{
		ID:          s.ID.String(),
		StaffID:     s.StaffID.String(),
		EffectiveFrom: s.EffectiveFrom.Format("2006-01-02"),
		GrossSalary: s.GrossSalary.StringFixed(2),
		NetSalary:   s.NetSalary.StringFixed(2),
		IsCurrent:   s.IsCurrent,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}

	if s.Staff != nil {
		resp.StaffName = s.Staff.FullName()
	}

	if s.StructureID != nil {
		resp.StructureID = s.StructureID.String()
	}

	if s.Structure != nil {
		resp.StructureName = s.Structure.Name
	}

	if s.EffectiveTo != nil {
		resp.EffectiveTo = s.EffectiveTo.Format("2006-01-02")
	}

	if s.CTC != nil {
		resp.CTC = s.CTC.StringFixed(2)
	}

	if s.RevisionReason != nil {
		resp.RevisionReason = *s.RevisionReason
	}

	if includeComponents && len(s.Components) > 0 {
		resp.Components = make([]StaffSalaryComponentResponse, len(s.Components))
		for i, comp := range s.Components {
			resp.Components[i] = ToStaffSalaryComponentResponse(&comp)
		}
	}

	return resp
}

// ToStaffSalaryComponentResponse converts a StaffSalaryComponent to a response.
func ToStaffSalaryComponentResponse(sc *models.StaffSalaryComponent) StaffSalaryComponentResponse {
	resp := StaffSalaryComponentResponse{
		ID:           sc.ID.String(),
		ComponentID:  sc.ComponentID.String(),
		Amount:       sc.Amount.StringFixed(2),
		IsOverridden: sc.IsOverridden,
	}

	if sc.Component != nil {
		resp.ComponentName = sc.Component.Name
		resp.ComponentCode = sc.Component.Code
		resp.ComponentType = string(sc.Component.ComponentType)
	}

	return resp
}

// ToStaffSalaryResponses converts a slice of StaffSalary models to responses.
func ToStaffSalaryResponses(salaries []models.StaffSalary) []StaffSalaryResponse {
	responses := make([]StaffSalaryResponse, len(salaries))
	for i, s := range salaries {
		responses[i] = ToStaffSalaryResponse(&s, false)
	}
	return responses
}

// ========================================
// Dropdown DTOs
// ========================================

// ComponentDropdownItem represents a component for dropdown.
type ComponentDropdownItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	ComponentType string `json:"componentType"`
}

// StructureDropdownItem represents a structure for dropdown.
type StructureDropdownItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ToComponentDropdownItem converts a SalaryComponent to a dropdown item.
func ToComponentDropdownItem(c *models.SalaryComponent) ComponentDropdownItem {
	return ComponentDropdownItem{
		ID:            c.ID.String(),
		Name:          c.Name,
		Code:          c.Code,
		ComponentType: string(c.ComponentType),
	}
}

// ToComponentDropdownItems converts components to dropdown items.
func ToComponentDropdownItems(components []models.SalaryComponent) []ComponentDropdownItem {
	items := make([]ComponentDropdownItem, len(components))
	for i, c := range components {
		items[i] = ToComponentDropdownItem(&c)
	}
	return items
}

// ToStructureDropdownItem converts a SalaryStructure to a dropdown item.
func ToStructureDropdownItem(s *models.SalaryStructure) StructureDropdownItem {
	return StructureDropdownItem{
		ID:   s.ID.String(),
		Name: s.Name,
		Code: s.Code,
	}
}

// ToStructureDropdownItems converts structures to dropdown items.
func ToStructureDropdownItems(structures []models.SalaryStructure) []StructureDropdownItem {
	items := make([]StructureDropdownItem, len(structures))
	for i, s := range structures {
		items[i] = ToStructureDropdownItem(&s)
	}
	return items
}
