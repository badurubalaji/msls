// Package branch provides HTTP handlers for branch management endpoints.
package branch

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// ============================================================================
// Request DTOs
// ============================================================================

// CreateBranchRequest represents the request body for creating a branch.
type CreateBranchRequest struct {
	Code         string                 `json:"code" binding:"required,max=20"`
	Name         string                 `json:"name" binding:"required,max=200"`
	AddressLine1 string                 `json:"addressLine1" binding:"max=255"`
	AddressLine2 string                 `json:"addressLine2" binding:"max=255"`
	City         string                 `json:"city" binding:"max=100"`
	State        string                 `json:"state" binding:"max=100"`
	PostalCode   string                 `json:"postalCode" binding:"max=20"`
	Country      string                 `json:"country" binding:"max=100"`
	Phone        string                 `json:"phone" binding:"max=20"`
	Email        string                 `json:"email" binding:"omitempty,email,max=255"`
	LogoURL      string                 `json:"logoUrl" binding:"omitempty,url,max=500"`
	Timezone     string                 `json:"timezone" binding:"max=50"`
	IsPrimary    bool                   `json:"isPrimary"`
	Settings     map[string]interface{} `json:"settings"`
}

// UpdateBranchRequest represents the request body for updating a branch.
type UpdateBranchRequest struct {
	Name         *string                `json:"name" binding:"omitempty,max=200"`
	AddressLine1 *string                `json:"addressLine1" binding:"omitempty,max=255"`
	AddressLine2 *string                `json:"addressLine2" binding:"omitempty,max=255"`
	City         *string                `json:"city" binding:"omitempty,max=100"`
	State        *string                `json:"state" binding:"omitempty,max=100"`
	PostalCode   *string                `json:"postalCode" binding:"omitempty,max=20"`
	Country      *string                `json:"country" binding:"omitempty,max=100"`
	Phone        *string                `json:"phone" binding:"omitempty,max=20"`
	Email        *string                `json:"email" binding:"omitempty,email,max=255"`
	LogoURL      *string                `json:"logoUrl" binding:"omitempty,url,max=500"`
	Timezone     *string                `json:"timezone" binding:"omitempty,max=50"`
	IsPrimary    *bool                  `json:"isPrimary"`
	Settings     map[string]interface{} `json:"settings"`
}

// SetStatusRequest represents the request body for setting branch status.
type SetStatusRequest struct {
	IsActive bool `json:"isActive"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// BranchResponse represents a branch in API responses.
type BranchResponse struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	LogoURL      string `json:"logoUrl,omitempty"`
	Timezone     string `json:"timezone"`
	IsPrimary    bool   `json:"isPrimary"`
	IsActive     bool   `json:"isActive"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	// Placeholder fields for future use
	StudentCount int64 `json:"studentCount,omitempty"`
	StaffCount   int64 `json:"staffCount,omitempty"`
}

// BranchListResponse represents a list of branches.
type BranchListResponse struct {
	Branches []BranchResponse `json:"branches"`
	Total    int64            `json:"total"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// branchToResponse converts a Branch model to a BranchResponse.
func branchToResponse(branch *models.Branch) BranchResponse {
	// Parse address lines from the combined street field
	addressLine1 := branch.Address.Street
	addressLine2 := ""
	if parts := strings.SplitN(branch.Address.Street, "\n", 2); len(parts) == 2 {
		addressLine1 = parts[0]
		addressLine2 = parts[1]
	}

	// Get phone and email from settings
	phone := branch.Settings.ContactPhone
	email := branch.Settings.ContactEmail

	// Default timezone
	timezone := "Asia/Kolkata"

	return BranchResponse{
		ID:           branch.ID.String(),
		Code:         branch.Code,
		Name:         branch.Name,
		AddressLine1: addressLine1,
		AddressLine2: addressLine2,
		City:         branch.Address.City,
		State:        branch.Address.State,
		PostalCode:   branch.Address.PostalCode,
		Country:      branch.Address.Country,
		Phone:        phone,
		Email:        email,
		LogoURL:      "", // Logo URL would be stored in settings if needed
		Timezone:     timezone,
		IsPrimary:    branch.IsPrimary,
		IsActive:     branch.Status == models.StatusActive,
		Status:       string(branch.Status),
		CreatedAt:    branch.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    branch.UpdatedAt.Format(time.RFC3339),
		StudentCount: 0, // Placeholder - will be populated when student module is available
		StaffCount:   0, // Placeholder - will be populated when staff module is available
	}
}

// branchesToResponses converts a slice of Branch models to BranchResponses.
func branchesToResponses(branches []models.Branch) []BranchResponse {
	responses := make([]BranchResponse, len(branches))
	for i, branch := range branches {
		responses[i] = branchToResponse(&branch)
	}
	return responses
}

// parseUUID parses a string to a UUID.
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
