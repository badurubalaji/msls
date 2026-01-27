// Package student provides student management functionality.
package student

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// CreateStudentDTO represents a request to create a student.
type CreateStudentDTO struct {
	TenantID            uuid.UUID
	BranchID            uuid.UUID
	FirstName           string
	MiddleName          string
	LastName            string
	DateOfBirth         time.Time
	Gender              models.Gender
	BloodGroup          string
	AadhaarNumber       string
	AdmissionDate       *time.Time
	CurrentAddress      *AddressDTO
	PermanentAddress    *AddressDTO
	SameAsCurrentAddress bool
	CreatedBy           *uuid.UUID
}

// UpdateStudentDTO represents a request to update a student.
type UpdateStudentDTO struct {
	FirstName           *string
	MiddleName          *string
	LastName            *string
	DateOfBirth         *time.Time
	Gender              *models.Gender
	BloodGroup          *string
	AadhaarNumber       *string
	Status              *models.StudentStatus
	PhotoURL            *string
	BirthCertificateURL *string
	CurrentAddress      *AddressDTO
	PermanentAddress    *AddressDTO
	SameAsCurrentAddress *bool
	Version             int // For optimistic locking
	UpdatedBy           *uuid.UUID
}

// AddressDTO represents an address in requests.
type AddressDTO struct {
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
	Country      string
}

// StudentResponse represents a student in API responses.
type StudentResponse struct {
	ID                  string            `json:"id"`
	AdmissionNumber     string            `json:"admissionNumber"`
	FirstName           string            `json:"firstName"`
	MiddleName          string            `json:"middleName,omitempty"`
	LastName            string            `json:"lastName"`
	FullName            string            `json:"fullName"`
	DateOfBirth         string            `json:"dateOfBirth"`
	Gender              string            `json:"gender"`
	BloodGroup          string            `json:"bloodGroup,omitempty"`
	AadhaarNumber       string            `json:"aadhaarNumber,omitempty"`
	PhotoURL            string            `json:"photoUrl,omitempty"`
	BirthCertificateURL string            `json:"birthCertificateUrl,omitempty"`
	Status              string            `json:"status"`
	AdmissionDate       string            `json:"admissionDate"`
	BranchID            string            `json:"branchId"`
	BranchName          string            `json:"branchName,omitempty"`
	Initials            string            `json:"initials"`
	CurrentAddress      *AddressResponse  `json:"currentAddress,omitempty"`
	PermanentAddress    *AddressResponse  `json:"permanentAddress,omitempty"`
	CreatedAt           string            `json:"createdAt"`
	UpdatedAt           string            `json:"updatedAt"`
	Version             int               `json:"version"`
}

// AddressResponse represents an address in API responses.
type AddressResponse struct {
	ID           string `json:"id"`
	AddressType  string `json:"addressType"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
	Country      string `json:"country"`
}

// StudentListResponse represents a paginated list of students.
type StudentListResponse struct {
	Students   []StudentResponse `json:"students"`
	NextCursor string            `json:"nextCursor,omitempty"`
	HasMore    bool              `json:"hasMore"`
	Total      int64             `json:"total,omitempty"`
}

// ListFilter contains filters for listing students.
type ListFilter struct {
	TenantID      uuid.UUID
	BranchID      *uuid.UUID
	ClassID       *uuid.UUID
	SectionID     *uuid.UUID
	Status        *models.StudentStatus
	Gender        *models.Gender
	AdmissionFrom *time.Time
	AdmissionTo   *time.Time
	Search        string
	Cursor        string
	Limit         int
	SortBy        string
	SortOrder     string
}

// ToStudentResponse converts a Student model to a StudentResponse.
func ToStudentResponse(student *models.Student) StudentResponse {
	resp := StudentResponse{
		ID:                  student.ID.String(),
		AdmissionNumber:     student.AdmissionNumber,
		FirstName:           student.FirstName,
		MiddleName:          student.MiddleName,
		LastName:            student.LastName,
		FullName:            student.FullName(),
		DateOfBirth:         student.DateOfBirth.Format("2006-01-02"),
		Gender:              string(student.Gender),
		BloodGroup:          student.BloodGroup,
		AadhaarNumber:       student.AadhaarNumber,
		PhotoURL:            student.PhotoURL,
		BirthCertificateURL: student.BirthCertificateURL,
		Status:              string(student.Status),
		AdmissionDate:       student.AdmissionDate.Format("2006-01-02"),
		BranchID:            student.BranchID.String(),
		Initials:            student.GetInitials(),
		CreatedAt:           student.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           student.UpdatedAt.Format(time.RFC3339),
		Version:             student.Version,
	}

	// Add branch name if loaded
	if student.Branch.ID != uuid.Nil {
		resp.BranchName = student.Branch.Name
	}

	// Add addresses
	for _, addr := range student.Addresses {
		addrResp := &AddressResponse{
			ID:           addr.ID.String(),
			AddressType:  string(addr.AddressType),
			AddressLine1: addr.AddressLine1,
			AddressLine2: addr.AddressLine2,
			City:         addr.City,
			State:        addr.State,
			PostalCode:   addr.PostalCode,
			Country:      addr.Country,
		}
		if addr.AddressType == models.AddressTypeCurrent {
			resp.CurrentAddress = addrResp
		} else if addr.AddressType == models.AddressTypePermanent {
			resp.PermanentAddress = addrResp
		}
	}

	return resp
}

// ToStudentResponses converts a slice of Student models to StudentResponses.
func ToStudentResponses(students []models.Student) []StudentResponse {
	responses := make([]StudentResponse, len(students))
	for i, student := range students {
		responses[i] = ToStudentResponse(&student)
	}
	return responses
}
