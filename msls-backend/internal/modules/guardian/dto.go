// Package guardian provides guardian and emergency contact management functionality.
package guardian

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// CreateGuardianDTO represents a request to create a guardian.
type CreateGuardianDTO struct {
	TenantID        uuid.UUID
	StudentID       uuid.UUID
	Relation        models.GuardianRelation
	FirstName       string
	LastName        string
	Phone           string
	Email           string
	Occupation      string
	AnnualIncome    decimal.Decimal
	Education       string
	IsPrimary       bool
	HasPortalAccess bool
	AddressLine1    string
	AddressLine2    string
	City            string
	State           string
	PostalCode      string
	Country         string
	CreatedBy       *uuid.UUID
}

// UpdateGuardianDTO represents a request to update a guardian.
type UpdateGuardianDTO struct {
	Relation        *models.GuardianRelation
	FirstName       *string
	LastName        *string
	Phone           *string
	Email           *string
	Occupation      *string
	AnnualIncome    *decimal.Decimal
	Education       *string
	IsPrimary       *bool
	HasPortalAccess *bool
	AddressLine1    *string
	AddressLine2    *string
	City            *string
	State           *string
	PostalCode      *string
	Country         *string
	UpdatedBy       *uuid.UUID
}

// CreateEmergencyContactDTO represents a request to create an emergency contact.
type CreateEmergencyContactDTO struct {
	TenantID       uuid.UUID
	StudentID      uuid.UUID
	Name           string
	Relation       string
	Phone          string
	AlternatePhone string
	Priority       int
	Notes          string
	CreatedBy      *uuid.UUID
}

// UpdateEmergencyContactDTO represents a request to update an emergency contact.
type UpdateEmergencyContactDTO struct {
	Name           *string
	Relation       *string
	Phone          *string
	AlternatePhone *string
	Priority       *int
	Notes          *string
	UpdatedBy      *uuid.UUID
}

// GuardianResponse represents a guardian in API responses.
type GuardianResponse struct {
	ID              string  `json:"id"`
	StudentID       string  `json:"studentId"`
	Relation        string  `json:"relation"`
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	FullName        string  `json:"fullName"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email,omitempty"`
	Occupation      string  `json:"occupation,omitempty"`
	AnnualIncome    string  `json:"annualIncome,omitempty"`
	Education       string  `json:"education,omitempty"`
	IsPrimary       bool    `json:"isPrimary"`
	HasPortalAccess bool    `json:"hasPortalAccess"`
	UserID          *string `json:"userId,omitempty"`
	AddressLine1    string  `json:"addressLine1,omitempty"`
	AddressLine2    string  `json:"addressLine2,omitempty"`
	City            string  `json:"city,omitempty"`
	State           string  `json:"state,omitempty"`
	PostalCode      string  `json:"postalCode,omitempty"`
	Country         string  `json:"country,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// EmergencyContactResponse represents an emergency contact in API responses.
type EmergencyContactResponse struct {
	ID             string `json:"id"`
	StudentID      string `json:"studentId"`
	Name           string `json:"name"`
	Relation       string `json:"relation"`
	Phone          string `json:"phone"`
	AlternatePhone string `json:"alternatePhone,omitempty"`
	Priority       int    `json:"priority"`
	Notes          string `json:"notes,omitempty"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// GuardianListResponse represents a list of guardians.
type GuardianListResponse struct {
	Guardians []GuardianResponse `json:"guardians"`
	Total     int                `json:"total"`
}

// EmergencyContactListResponse represents a list of emergency contacts.
type EmergencyContactListResponse struct {
	Contacts []EmergencyContactResponse `json:"contacts"`
	Total    int                        `json:"total"`
}

// ToGuardianResponse converts a StudentGuardian model to a GuardianResponse.
func ToGuardianResponse(guardian *models.StudentGuardian) GuardianResponse {
	resp := GuardianResponse{
		ID:              guardian.ID.String(),
		StudentID:       guardian.StudentID.String(),
		Relation:        string(guardian.Relation),
		FirstName:       guardian.FirstName,
		LastName:        guardian.LastName,
		FullName:        guardian.FullName(),
		Phone:           guardian.Phone,
		Email:           guardian.Email,
		Occupation:      guardian.Occupation,
		Education:       guardian.Education,
		IsPrimary:       guardian.IsPrimary,
		HasPortalAccess: guardian.HasPortalAccess,
		AddressLine1:    guardian.AddressLine1,
		AddressLine2:    guardian.AddressLine2,
		City:            guardian.City,
		State:           guardian.State,
		PostalCode:      guardian.PostalCode,
		Country:         guardian.Country,
		CreatedAt:       guardian.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       guardian.UpdatedAt.Format(time.RFC3339),
	}

	if !guardian.AnnualIncome.IsZero() {
		resp.AnnualIncome = guardian.AnnualIncome.String()
	}

	if guardian.UserID != nil {
		userIDStr := guardian.UserID.String()
		resp.UserID = &userIDStr
	}

	return resp
}

// ToGuardianResponses converts a slice of StudentGuardian models to GuardianResponses.
func ToGuardianResponses(guardians []models.StudentGuardian) []GuardianResponse {
	responses := make([]GuardianResponse, len(guardians))
	for i, guardian := range guardians {
		responses[i] = ToGuardianResponse(&guardian)
	}
	return responses
}

// ToEmergencyContactResponse converts a StudentEmergencyContact model to an EmergencyContactResponse.
func ToEmergencyContactResponse(contact *models.StudentEmergencyContact) EmergencyContactResponse {
	return EmergencyContactResponse{
		ID:             contact.ID.String(),
		StudentID:      contact.StudentID.String(),
		Name:           contact.Name,
		Relation:       contact.Relation,
		Phone:          contact.Phone,
		AlternatePhone: contact.AlternatePhone,
		Priority:       contact.Priority,
		Notes:          contact.Notes,
		CreatedAt:      contact.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      contact.UpdatedAt.Format(time.RFC3339),
	}
}

// ToEmergencyContactResponses converts a slice of StudentEmergencyContact models to EmergencyContactResponses.
func ToEmergencyContactResponses(contacts []models.StudentEmergencyContact) []EmergencyContactResponse {
	responses := make([]EmergencyContactResponse, len(contacts))
	for i, contact := range contacts {
		responses[i] = ToEmergencyContactResponse(&contact)
	}
	return responses
}
