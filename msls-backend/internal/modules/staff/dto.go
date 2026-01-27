// Package staff provides staff management functionality.
package staff

import (
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// CreateStaffDTO represents a request to create a staff member.
type CreateStaffDTO struct {
	TenantID  uuid.UUID
	BranchID  uuid.UUID
	
	// Personal details
	FirstName     string
	MiddleName    string
	LastName      string
	DateOfBirth   time.Time
	Gender        models.Gender
	BloodGroup    string
	Nationality   string
	Religion      string
	MaritalStatus string
	
	// Contact details
	PersonalEmail            string
	WorkEmail                string
	PersonalPhone            string
	WorkPhone                string
	EmergencyContactName     string
	EmergencyContactPhone    string
	EmergencyContactRelation string
	
	// Address
	CurrentAddress   *AddressDTO
	PermanentAddress *AddressDTO
	SameAsCurrent    bool
	
	// Employment details
	StaffType          models.StaffType
	DepartmentID       *uuid.UUID
	DesignationID      *uuid.UUID
	ReportingManagerID *uuid.UUID
	JoinDate           time.Time
	ConfirmationDate   *time.Time
	ProbationEndDate   *time.Time
	
	// Profile
	Bio string
	
	CreatedBy *uuid.UUID
}

// UpdateStaffDTO represents a request to update a staff member.
type UpdateStaffDTO struct {
	// Personal details
	FirstName     *string
	MiddleName    *string
	LastName      *string
	DateOfBirth   *time.Time
	Gender        *models.Gender
	BloodGroup    *string
	Nationality   *string
	Religion      *string
	MaritalStatus *string
	
	// Contact details
	PersonalEmail            *string
	WorkEmail                *string
	PersonalPhone            *string
	WorkPhone                *string
	EmergencyContactName     *string
	EmergencyContactPhone    *string
	EmergencyContactRelation *string
	
	// Address
	CurrentAddress   *AddressDTO
	PermanentAddress *AddressDTO
	SameAsCurrent    *bool
	
	// Employment details
	StaffType          *models.StaffType
	DepartmentID       *uuid.UUID
	DesignationID      *uuid.UUID
	ReportingManagerID *uuid.UUID
	ConfirmationDate   *time.Time
	ProbationEndDate   *time.Time
	
	// Profile
	Bio      *string
	PhotoURL *string
	
	Version   int // For optimistic locking
	UpdatedBy *uuid.UUID
}

// StatusUpdateDTO represents a request to update staff status.
type StatusUpdateDTO struct {
	Status        models.StaffStatus
	Reason        string
	EffectiveDate time.Time
	UpdatedBy     *uuid.UUID
}

// AddressDTO represents an address in requests.
type AddressDTO struct {
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	Pincode      string
	Country      string
}

// StaffResponse represents a staff member in API responses.
type StaffResponse struct {
	ID               string `json:"id"`
	EmployeeID       string `json:"employeeId"`
	EmployeeIDPrefix string `json:"employeeIdPrefix"`
	
	// Personal details
	FirstName     string `json:"firstName"`
	MiddleName    string `json:"middleName,omitempty"`
	LastName      string `json:"lastName"`
	FullName      string `json:"fullName"`
	Initials      string `json:"initials"`
	DateOfBirth   string `json:"dateOfBirth"`
	Gender        string `json:"gender"`
	BloodGroup    string `json:"bloodGroup,omitempty"`
	Nationality   string `json:"nationality,omitempty"`
	Religion      string `json:"religion,omitempty"`
	MaritalStatus string `json:"maritalStatus,omitempty"`
	
	// Contact details
	PersonalEmail            string `json:"personalEmail,omitempty"`
	WorkEmail                string `json:"workEmail"`
	PersonalPhone            string `json:"personalPhone,omitempty"`
	WorkPhone                string `json:"workPhone"`
	EmergencyContactName     string `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone    string `json:"emergencyContactPhone,omitempty"`
	EmergencyContactRelation string `json:"emergencyContactRelation,omitempty"`
	
	// Address
	CurrentAddress   *AddressResponse `json:"currentAddress,omitempty"`
	PermanentAddress *AddressResponse `json:"permanentAddress,omitempty"`
	SameAsCurrent    bool             `json:"sameAsCurrent"`
	
	// Employment details
	StaffType          string              `json:"staffType"`
	DepartmentID       string              `json:"departmentId,omitempty"`
	DepartmentName     string              `json:"departmentName,omitempty"`
	DesignationID      string              `json:"designationId,omitempty"`
	DesignationName    string              `json:"designationName,omitempty"`
	ReportingManagerID string              `json:"reportingManagerId,omitempty"`
	ReportingManager   *StaffRefResponse   `json:"reportingManager,omitempty"`
	JoinDate           string              `json:"joinDate"`
	ConfirmationDate   string              `json:"confirmationDate,omitempty"`
	ProbationEndDate   string              `json:"probationEndDate,omitempty"`
	
	// Status
	Status          string `json:"status"`
	StatusReason    string `json:"statusReason,omitempty"`
	TerminationDate string `json:"terminationDate,omitempty"`
	
	// Profile
	PhotoURL string `json:"photoUrl,omitempty"`
	Bio      string `json:"bio,omitempty"`
	
	// Branch
	BranchID   string `json:"branchId"`
	BranchName string `json:"branchName,omitempty"`
	
	// Audit
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Version   int    `json:"version"`
}

// StaffRefResponse is a minimal staff reference for nested responses.
type StaffRefResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl,omitempty"`
}

// AddressResponse represents an address in API responses.
type AddressResponse struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	Pincode      string `json:"pincode"`
	Country      string `json:"country"`
}

// StaffListResponse represents a paginated list of staff.
type StaffListResponse struct {
	Staff      []StaffResponse `json:"staff"`
	NextCursor string          `json:"nextCursor,omitempty"`
	HasMore    bool            `json:"hasMore"`
	Total      int64           `json:"total,omitempty"`
}

// ListFilter contains filters for listing staff.
type ListFilter struct {
	TenantID      uuid.UUID
	BranchID      *uuid.UUID
	DepartmentID  *uuid.UUID
	DesignationID *uuid.UUID
	StaffType     *models.StaffType
	Status        *models.StaffStatus
	Gender        *models.Gender
	JoinDateFrom  *time.Time
	JoinDateTo    *time.Time
	Search        string
	Cursor        string
	Limit         int
	SortBy        string
	SortOrder     string
}

// StatusHistoryResponse represents a status history entry.
type StatusHistoryResponse struct {
	ID            string `json:"id"`
	OldStatus     string `json:"oldStatus,omitempty"`
	NewStatus     string `json:"newStatus"`
	Reason        string `json:"reason,omitempty"`
	EffectiveDate string `json:"effectiveDate"`
	ChangedBy     string `json:"changedBy,omitempty"`
	ChangedAt     string `json:"changedAt"`
}

// ToStaffResponse converts a Staff model to a StaffResponse.
func ToStaffResponse(staff *models.Staff) StaffResponse {
	resp := StaffResponse{
		ID:               staff.ID.String(),
		EmployeeID:       staff.EmployeeID,
		EmployeeIDPrefix: staff.EmployeeIDPrefix,
		FirstName:        staff.FirstName,
		MiddleName:       staff.MiddleName,
		LastName:         staff.LastName,
		FullName:         staff.FullName(),
		Initials:         staff.GetInitials(),
		DateOfBirth:      staff.DateOfBirth.Format("2006-01-02"),
		Gender:           string(staff.Gender),
		BloodGroup:       staff.BloodGroup,
		Nationality:      staff.Nationality,
		Religion:         staff.Religion,
		MaritalStatus:    staff.MaritalStatus,
		PersonalEmail:    staff.PersonalEmail,
		WorkEmail:        staff.WorkEmail,
		PersonalPhone:    staff.PersonalPhone,
		WorkPhone:        staff.WorkPhone,
		EmergencyContactName:     staff.EmergencyContactName,
		EmergencyContactPhone:    staff.EmergencyContactPhone,
		EmergencyContactRelation: staff.EmergencyContactRelation,
		SameAsCurrent:    staff.SameAsCurrent,
		StaffType:        string(staff.StaffType),
		JoinDate:         staff.JoinDate.Format("2006-01-02"),
		Status:           string(staff.Status),
		StatusReason:     staff.StatusReason,
		PhotoURL:         staff.PhotoURL,
		Bio:              staff.Bio,
		BranchID:         staff.BranchID.String(),
		CreatedAt:        staff.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        staff.UpdatedAt.Format(time.RFC3339),
		Version:          staff.Version,
	}
	
	// Add current address
	if staff.CurrentAddressLine1 != "" {
		resp.CurrentAddress = &AddressResponse{
			AddressLine1: staff.CurrentAddressLine1,
			AddressLine2: staff.CurrentAddressLine2,
			City:         staff.CurrentCity,
			State:        staff.CurrentState,
			Pincode:      staff.CurrentPincode,
			Country:      staff.CurrentCountry,
		}
	}
	
	// Add permanent address
	if staff.PermanentAddressLine1 != "" {
		resp.PermanentAddress = &AddressResponse{
			AddressLine1: staff.PermanentAddressLine1,
			AddressLine2: staff.PermanentAddressLine2,
			City:         staff.PermanentCity,
			State:        staff.PermanentState,
			Pincode:      staff.PermanentPincode,
			Country:      staff.PermanentCountry,
		}
	}
	
	// Add optional dates
	if staff.ConfirmationDate != nil {
		resp.ConfirmationDate = staff.ConfirmationDate.Format("2006-01-02")
	}
	if staff.ProbationEndDate != nil {
		resp.ProbationEndDate = staff.ProbationEndDate.Format("2006-01-02")
	}
	if staff.TerminationDate != nil {
		resp.TerminationDate = staff.TerminationDate.Format("2006-01-02")
	}
	
	// Add department
	if staff.DepartmentID != nil {
		resp.DepartmentID = staff.DepartmentID.String()
		if staff.Department != nil {
			resp.DepartmentName = staff.Department.Name
		}
	}
	
	// Add designation
	if staff.DesignationID != nil {
		resp.DesignationID = staff.DesignationID.String()
		if staff.Designation != nil {
			resp.DesignationName = staff.Designation.Name
		}
	}
	
	// Add reporting manager
	if staff.ReportingManagerID != nil {
		resp.ReportingManagerID = staff.ReportingManagerID.String()
		if staff.ReportingManager != nil {
			resp.ReportingManager = &StaffRefResponse{
				ID:       staff.ReportingManager.ID.String(),
				Name:     staff.ReportingManager.FullName(),
				PhotoURL: staff.ReportingManager.PhotoURL,
			}
		}
	}
	
	// Add branch name if loaded
	if staff.Branch.ID != uuid.Nil {
		resp.BranchName = staff.Branch.Name
	}
	
	return resp
}

// ToStaffResponses converts a slice of Staff models to StaffResponses.
func ToStaffResponses(staffList []models.Staff) []StaffResponse {
	responses := make([]StaffResponse, len(staffList))
	for i, staff := range staffList {
		responses[i] = ToStaffResponse(&staff)
	}
	return responses
}

// ToStatusHistoryResponse converts a StaffStatusHistory model to a response.
func ToStatusHistoryResponse(history *models.StaffStatusHistory) StatusHistoryResponse {
	resp := StatusHistoryResponse{
		ID:            history.ID.String(),
		OldStatus:     history.OldStatus,
		NewStatus:     history.NewStatus,
		Reason:        history.Reason,
		EffectiveDate: history.EffectiveDate.Format("2006-01-02"),
		ChangedAt:     history.ChangedAt.Format(time.RFC3339),
	}
	if history.ChangedBy != nil {
		resp.ChangedBy = history.ChangedBy.String()
	}
	return resp
}

// ToStatusHistoryResponses converts a slice of status history models to responses.
func ToStatusHistoryResponses(histories []models.StaffStatusHistory) []StatusHistoryResponse {
	responses := make([]StatusHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = ToStatusHistoryResponse(&history)
	}
	return responses
}
