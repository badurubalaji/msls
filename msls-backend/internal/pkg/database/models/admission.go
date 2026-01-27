// Package models provides GORM model definitions for the MSLS database.
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AdmissionSessionStatus represents the status of an admission session.
type AdmissionSessionStatus string

// Admission session status constants.
const (
	SessionStatusUpcoming AdmissionSessionStatus = "upcoming"
	SessionStatusOpen     AdmissionSessionStatus = "open"
	SessionStatusClosed   AdmissionSessionStatus = "closed"
)

// IsValid checks if the session status is valid.
func (s AdmissionSessionStatus) IsValid() bool {
	switch s {
	case SessionStatusUpcoming, SessionStatusOpen, SessionStatusClosed:
		return true
	}
	return false
}

// String returns the string representation of the session status.
func (s AdmissionSessionStatus) String() string {
	return string(s)
}

// RequiredDocuments represents a list of required document types.
type RequiredDocuments []string

// Value implements the driver.Valuer interface for database serialization.
func (rd RequiredDocuments) Value() (driver.Value, error) {
	if rd == nil {
		return "[]", nil
	}
	return json.Marshal(rd)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (rd *RequiredDocuments) Scan(value interface{}) error {
	if value == nil {
		*rd = RequiredDocuments{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, rd)
}

// SessionSettings represents additional settings for an admission session.
type SessionSettings struct {
	AllowOnlineApplication bool   `json:"allowOnlineApplication,omitempty"`
	NotifyOnApplication    bool   `json:"notifyOnApplication,omitempty"`
	AutoConfirmPayment     bool   `json:"autoConfirmPayment,omitempty"`
	MaxApplicationsPerDay  int    `json:"maxApplicationsPerDay,omitempty"`
	Instructions           string `json:"instructions,omitempty"`
}

// Value implements the driver.Valuer interface for database serialization.
func (ss SessionSettings) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (ss *SessionSettings) Scan(value interface{}) error {
	if value == nil {
		*ss = SessionSettings{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, ss)
}

// ReservedSeats represents seat reservations by category.
type ReservedSeats map[string]int

// Value implements the driver.Valuer interface for database serialization.
func (rs ReservedSeats) Value() (driver.Value, error) {
	if rs == nil {
		return "{}", nil
	}
	return json.Marshal(rs)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (rs *ReservedSeats) Scan(value interface{}) error {
	if value == nil {
		*rs = ReservedSeats{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, rs)
}

// AdmissionSession represents an admission session in the database.
type AdmissionSession struct {
	ID                uuid.UUID              `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID          uuid.UUID              `gorm:"type:uuid;not null;index" json:"tenant_id"`
	BranchID          *uuid.UUID             `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	AcademicYearID    *uuid.UUID             `gorm:"type:uuid;index" json:"academic_year_id,omitempty"`
	Name              string                 `gorm:"size:200;not null" json:"name"`
	Description       string                 `gorm:"type:text" json:"description,omitempty"`
	StartDate         time.Time              `gorm:"type:date;not null" json:"start_date"`
	EndDate           time.Time              `gorm:"type:date;not null" json:"end_date"`
	Status            AdmissionSessionStatus `gorm:"size:20;not null;default:'upcoming'" json:"status"`
	ApplicationFee    decimal.Decimal        `gorm:"type:decimal(10,2);default:0" json:"application_fee"`
	RequiredDocuments RequiredDocuments      `gorm:"type:jsonb;default:'[]'" json:"required_documents"`
	Settings          SessionSettings        `gorm:"type:jsonb;default:'{}'" json:"settings"`
	CreatedAt         time.Time              `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time              `gorm:"not null;default:now()" json:"updated_at"`
	CreatedBy         *uuid.UUID             `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID             `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Seats []AdmissionSeat `gorm:"foreignKey:SessionID" json:"seats,omitempty"`
}

// TableName specifies the table name for AdmissionSession.
func (AdmissionSession) TableName() string {
	return "admission_sessions"
}

// AdmissionSeat represents seat configuration for a class in an admission session.
type AdmissionSeat struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID      uuid.UUID     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	SessionID     uuid.UUID     `gorm:"type:uuid;not null;index" json:"session_id"`
	ClassName     string        `gorm:"size:100;not null" json:"class_name"`
	TotalSeats    int           `gorm:"not null;default:0" json:"total_seats"`
	FilledSeats   int           `gorm:"not null;default:0" json:"filled_seats"`
	WaitlistLimit int           `gorm:"default:10" json:"waitlist_limit"`
	ReservedSeats ReservedSeats `gorm:"type:jsonb;default:'{}'" json:"reserved_seats"`
	CreatedAt     time.Time     `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"not null;default:now()" json:"updated_at"`

	// Relationships
	Session *AdmissionSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// TableName specifies the table name for AdmissionSeat.
func (AdmissionSeat) TableName() string {
	return "admission_seats"
}

// AvailableSeats returns the number of available seats.
func (as *AdmissionSeat) AvailableSeats() int {
	available := as.TotalSeats - as.FilledSeats
	if available < 0 {
		return 0
	}
	return available
}

// IsWaitlistAvailable checks if the waitlist has space.
func (as *AdmissionSeat) IsWaitlistAvailable(currentWaitlist int) bool {
	return currentWaitlist < as.WaitlistLimit
}

// TotalReservedSeats returns the total number of reserved seats across all categories.
func (as *AdmissionSeat) TotalReservedSeats() int {
	total := 0
	for _, count := range as.ReservedSeats {
		total += count
	}
	return total
}

// GeneralSeats returns the number of general (unreserved) seats available.
func (as *AdmissionSeat) GeneralSeats() int {
	general := as.TotalSeats - as.TotalReservedSeats()
	if general < 0 {
		return 0
	}
	return general
}

// EnquirySource represents the source of an admission enquiry.
type EnquirySource string

// Enquiry source constants.
const (
	EnquirySourceWalkIn   EnquirySource = "walk_in"
	EnquirySourceWebsite  EnquirySource = "website"
	EnquirySourceReferral EnquirySource = "referral"
	EnquirySourcePhone    EnquirySource = "phone"
	EnquirySourceEmail    EnquirySource = "email"
	EnquirySourceOther    EnquirySource = "other"
)

// IsValid checks if the enquiry source is valid.
func (es EnquirySource) IsValid() bool {
	switch es {
	case EnquirySourceWalkIn, EnquirySourceWebsite, EnquirySourceReferral,
		EnquirySourcePhone, EnquirySourceEmail, EnquirySourceOther:
		return true
	}
	return false
}

// String returns the string representation of the enquiry source.
func (es EnquirySource) String() string {
	return string(es)
}

// EnquiryStatus represents the status of an admission enquiry.
type EnquiryStatus string

// Enquiry status constants.
const (
	EnquiryStatusNew        EnquiryStatus = "new"
	EnquiryStatusContacted  EnquiryStatus = "contacted"
	EnquiryStatusConverted  EnquiryStatus = "converted"
	EnquiryStatusNotInterested EnquiryStatus = "not_interested"
	EnquiryStatusClosed     EnquiryStatus = "closed"
)

// IsValid checks if the enquiry status is valid.
func (es EnquiryStatus) IsValid() bool {
	switch es {
	case EnquiryStatusNew, EnquiryStatusContacted, EnquiryStatusConverted,
		EnquiryStatusNotInterested, EnquiryStatusClosed:
		return true
	}
	return false
}

// String returns the string representation of the enquiry status.
func (es EnquiryStatus) String() string {
	return string(es)
}

// AdmissionEnquiry represents an admission enquiry in the database.
type AdmissionEnquiry struct {
	ID              uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID        uuid.UUID     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	BranchID        *uuid.UUID    `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	SessionID       *uuid.UUID    `gorm:"type:uuid;index" json:"session_id,omitempty"`
	EnquiryNumber   string        `gorm:"size:50;not null;uniqueIndex:idx_enquiry_number_tenant" json:"enquiry_number"`
	StudentName     string        `gorm:"size:200;not null" json:"student_name"`
	ParentName      string        `gorm:"size:200" json:"parent_name,omitempty"`
	Email           string        `gorm:"size:255" json:"email,omitempty"`
	Phone           string        `gorm:"size:20;not null" json:"phone"`
	ClassName       string        `gorm:"size:100" json:"class_name,omitempty"`
	Source          EnquirySource `gorm:"size:50;not null;default:'walk_in'" json:"source"`
	Status          EnquiryStatus `gorm:"size:50;not null;default:'new'" json:"status"`
	Notes           string        `gorm:"type:text" json:"notes,omitempty"`
	ReferralName    string        `gorm:"size:200" json:"referral_name,omitempty"`
	EnquiryDate     time.Time     `gorm:"type:date;not null;default:CURRENT_DATE" json:"enquiry_date"`
	FollowUpDate    *time.Time    `gorm:"type:date" json:"follow_up_date,omitempty"`
	CreatedAt       time.Time     `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"not null;default:now()" json:"updated_at"`
	CreatedBy       *uuid.UUID    `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy       *uuid.UUID    `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Session *AdmissionSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// TableName specifies the table name for AdmissionEnquiry.
func (AdmissionEnquiry) TableName() string {
	return "admission_enquiries"
}

// ApplicationStatus represents the status of an admission application.
type ApplicationStatus string

// Application status constants.
const (
	ApplicationStatusDraft              ApplicationStatus = "draft"
	ApplicationStatusSubmitted          ApplicationStatus = "submitted"
	ApplicationStatusUnderReview        ApplicationStatus = "under_review"
	ApplicationStatusDocumentsPending   ApplicationStatus = "documents_pending"
	ApplicationStatusDocumentsVerified  ApplicationStatus = "documents_verified"
	ApplicationStatusTestScheduled      ApplicationStatus = "test_scheduled"
	ApplicationStatusTestCompleted      ApplicationStatus = "test_completed"
	ApplicationStatusInterviewScheduled ApplicationStatus = "interview_scheduled"
	ApplicationStatusInterviewCompleted ApplicationStatus = "interview_completed"
	ApplicationStatusShortlisted        ApplicationStatus = "shortlisted"
	ApplicationStatusApproved           ApplicationStatus = "approved"
	ApplicationStatusRejected           ApplicationStatus = "rejected"
	ApplicationStatusWaitlisted         ApplicationStatus = "waitlisted"
	ApplicationStatusEnrolled           ApplicationStatus = "enrolled"
	ApplicationStatusWithdrawn          ApplicationStatus = "withdrawn"
)

// IsValid checks if the application status is valid.
func (as ApplicationStatus) IsValid() bool {
	switch as {
	case ApplicationStatusDraft, ApplicationStatusSubmitted, ApplicationStatusUnderReview,
		ApplicationStatusDocumentsPending, ApplicationStatusDocumentsVerified,
		ApplicationStatusTestScheduled, ApplicationStatusTestCompleted,
		ApplicationStatusInterviewScheduled, ApplicationStatusInterviewCompleted,
		ApplicationStatusShortlisted, ApplicationStatusApproved, ApplicationStatusRejected,
		ApplicationStatusWaitlisted, ApplicationStatusEnrolled, ApplicationStatusWithdrawn:
		return true
	}
	return false
}

// String returns the string representation of the application status.
func (as ApplicationStatus) String() string {
	return string(as)
}

// ApplicantDetails represents additional details about the applicant.
type ApplicantDetails struct {
	Gender          string `json:"gender,omitempty"`
	DateOfBirth     string `json:"dateOfBirth,omitempty"`
	Nationality     string `json:"nationality,omitempty"`
	Religion        string `json:"religion,omitempty"`
	Category        string `json:"category,omitempty"`
	BloodGroup      string `json:"bloodGroup,omitempty"`
	Address         string `json:"address,omitempty"`
	City            string `json:"city,omitempty"`
	State           string `json:"state,omitempty"`
	PinCode         string `json:"pinCode,omitempty"`
	PreviousSchool  string `json:"previousSchool,omitempty"`
	TransferReason  string `json:"transferReason,omitempty"`
}

// Value implements the driver.Valuer interface for database serialization.
func (ad ApplicantDetails) Value() (driver.Value, error) {
	return json.Marshal(ad)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (ad *ApplicantDetails) Scan(value interface{}) error {
	if value == nil {
		*ad = ApplicantDetails{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, ad)
}

// ParentGuardianInfo represents parent/guardian information.
type ParentGuardianInfo struct {
	FatherName       string `json:"fatherName,omitempty"`
	FatherOccupation string `json:"fatherOccupation,omitempty"`
	FatherPhone      string `json:"fatherPhone,omitempty"`
	FatherEmail      string `json:"fatherEmail,omitempty"`
	MotherName       string `json:"motherName,omitempty"`
	MotherOccupation string `json:"motherOccupation,omitempty"`
	MotherPhone      string `json:"motherPhone,omitempty"`
	MotherEmail      string `json:"motherEmail,omitempty"`
	GuardianName     string `json:"guardianName,omitempty"`
	GuardianRelation string `json:"guardianRelation,omitempty"`
	GuardianPhone    string `json:"guardianPhone,omitempty"`
	GuardianEmail    string `json:"guardianEmail,omitempty"`
}

// Value implements the driver.Valuer interface for database serialization.
func (pg ParentGuardianInfo) Value() (driver.Value, error) {
	return json.Marshal(pg)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (pg *ParentGuardianInfo) Scan(value interface{}) error {
	if value == nil {
		*pg = ParentGuardianInfo{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, pg)
}

// AdmissionApplication represents an admission application in the database.
type AdmissionApplication struct {
	ID                uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID          uuid.UUID         `gorm:"type:uuid;not null;index" json:"tenant_id"`
	BranchID          *uuid.UUID        `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	SessionID         uuid.UUID         `gorm:"type:uuid;not null;index" json:"session_id"`
	EnquiryID         *uuid.UUID        `gorm:"type:uuid;index" json:"enquiry_id,omitempty"`
	ApplicationNumber string            `gorm:"column:application_number;size:50;not null" json:"application_number"`

	// Student Information (matches DB columns)
	StudentName         string     `gorm:"column:student_name;size:200;not null" json:"student_name"`
	DateOfBirth         *time.Time `gorm:"column:date_of_birth;type:date" json:"date_of_birth,omitempty"`
	Gender              string     `gorm:"column:gender;size:20" json:"gender,omitempty"`
	BloodGroup          string     `gorm:"column:blood_group;size:10" json:"blood_group,omitempty"`
	Nationality         string     `gorm:"column:nationality;size:100;default:'Indian'" json:"nationality,omitempty"`
	Religion            string     `gorm:"column:religion;size:100" json:"religion,omitempty"`
	Caste               string     `gorm:"column:caste;size:100" json:"caste,omitempty"`
	Category            string     `gorm:"column:category;size:50" json:"category,omitempty"`
	MotherTongue        string     `gorm:"column:mother_tongue;size:100" json:"mother_tongue,omitempty"`
	AadharNumber        string     `gorm:"column:aadhar_number;size:12" json:"aadhar_number,omitempty"`

	// Academic Information
	ClassApplying       string           `gorm:"column:class_applying;size:50;not null" json:"class_applying"`
	PreviousSchool      string           `gorm:"column:previous_school;size:200" json:"previous_school,omitempty"`
	PreviousClass       string           `gorm:"column:previous_class;size:50" json:"previous_class,omitempty"`
	PreviousPercentage  *decimal.Decimal `gorm:"column:previous_percentage;type:decimal(5,2)" json:"previous_percentage,omitempty"`
	MediumOfInstruction string           `gorm:"column:medium_of_instruction;size:50" json:"medium_of_instruction,omitempty"`

	// Contact Information
	AddressLine1 string `gorm:"column:address_line1;size:255" json:"address_line1,omitempty"`
	AddressLine2 string `gorm:"column:address_line2;size:255" json:"address_line2,omitempty"`
	City         string `gorm:"column:city;size:100" json:"city,omitempty"`
	State        string `gorm:"column:state;size:100" json:"state,omitempty"`
	PostalCode   string `gorm:"column:postal_code;size:20" json:"postal_code,omitempty"`
	Country      string `gorm:"column:country;size:100;default:'India'" json:"country,omitempty"`

	// Parent/Guardian Information
	FatherName          string `gorm:"column:father_name;size:200" json:"father_name,omitempty"`
	FatherPhone         string `gorm:"column:father_phone;size:20" json:"father_phone,omitempty"`
	FatherEmail         string `gorm:"column:father_email;size:255" json:"father_email,omitempty"`
	FatherOccupation    string `gorm:"column:father_occupation;size:200" json:"father_occupation,omitempty"`
	FatherQualification string `gorm:"column:father_qualification;size:100" json:"father_qualification,omitempty"`
	MotherName          string `gorm:"column:mother_name;size:200" json:"mother_name,omitempty"`
	MotherPhone         string `gorm:"column:mother_phone;size:20" json:"mother_phone,omitempty"`
	MotherEmail         string `gorm:"column:mother_email;size:255" json:"mother_email,omitempty"`
	MotherOccupation    string `gorm:"column:mother_occupation;size:200" json:"mother_occupation,omitempty"`
	MotherQualification string `gorm:"column:mother_qualification;size:100" json:"mother_qualification,omitempty"`
	GuardianName        string `gorm:"column:guardian_name;size:200" json:"guardian_name,omitempty"`
	GuardianPhone       string `gorm:"column:guardian_phone;size:20" json:"guardian_phone,omitempty"`
	GuardianEmail       string `gorm:"column:guardian_email;size:255" json:"guardian_email,omitempty"`
	GuardianRelation    string `gorm:"column:guardian_relation;size:50" json:"guardian_relation,omitempty"`

	// Application Details
	Status        ApplicationStatus `gorm:"column:status;size:50;not null;default:'draft'" json:"status"`
	SubmittedAt   *time.Time        `gorm:"column:submitted_at;type:timestamptz" json:"submitted_at,omitempty"`
	Remarks       string            `gorm:"column:remarks;type:text" json:"remarks,omitempty"`
	InternalNotes string            `gorm:"column:internal_notes;type:text" json:"internal_notes,omitempty"`
	Priority      int               `gorm:"column:priority;default:0" json:"priority,omitempty"`

	// Payment Information
	FeePaid          bool       `gorm:"column:fee_paid;not null;default:false" json:"fee_paid"`
	PaymentReference string     `gorm:"column:payment_reference;size:100" json:"payment_reference,omitempty"`
	PaymentDate      *time.Time `gorm:"column:payment_date;type:timestamptz" json:"payment_date,omitempty"`

	// Extra Data (for custom fields)
	ExtraData map[string]interface{} `gorm:"column:extra_data;type:jsonb;default:'{}';serializer:json" json:"extra_data,omitempty"`

	// Audit Fields
	CreatedAt time.Time  `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	CreatedBy *uuid.UUID `gorm:"column:created_by;type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"column:updated_by;type:uuid" json:"updated_by,omitempty"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at,omitempty"`

	// Relationships
	Session   *AdmissionSession      `gorm:"foreignKey:SessionID" json:"session,omitempty"`
	Enquiry   *AdmissionEnquiry      `gorm:"foreignKey:EnquiryID" json:"enquiry,omitempty"`
	Parents   []ApplicationParent    `gorm:"foreignKey:ApplicationID" json:"parents,omitempty"`
	Documents []ApplicationDocument  `gorm:"foreignKey:ApplicationID" json:"documents,omitempty"`
}

// TableName specifies the table name for AdmissionApplication.
func (AdmissionApplication) TableName() string {
	return "admission_applications"
}

// DecisionType represents the type of admission decision.
type DecisionType string

// DecisionType constants.
const (
	DecisionApproved   DecisionType = "approved"
	DecisionWaitlisted DecisionType = "waitlisted"
	DecisionRejected   DecisionType = "rejected"
)

// IsValid checks if the decision type is valid.
func (d DecisionType) IsValid() bool {
	switch d {
	case DecisionApproved, DecisionWaitlisted, DecisionRejected:
		return true
	}
	return false
}

// String returns the string representation of the decision type.
func (d DecisionType) String() string {
	return string(d)
}

// AdmissionDecision represents an admission decision for an application.
type AdmissionDecision struct {
	ID               uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID         uuid.UUID    `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ApplicationID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"application_id"`
	Decision         DecisionType `gorm:"type:varchar(20);not null" json:"decision"`
	DecisionDate     time.Time    `gorm:"type:date;not null" json:"decision_date"`
	DecidedBy        *uuid.UUID   `gorm:"type:uuid" json:"decided_by,omitempty"`
	SectionAssigned  *string      `gorm:"type:varchar(50)" json:"section_assigned,omitempty"`
	WaitlistPosition *int         `gorm:"type:int" json:"waitlist_position,omitempty"`
	RejectionReason  *string      `gorm:"type:text" json:"rejection_reason,omitempty"`
	OfferLetterURL   *string      `gorm:"type:varchar(500)" json:"offer_letter_url,omitempty"`
	OfferValidUntil  *time.Time   `gorm:"type:date" json:"offer_valid_until,omitempty"`
	OfferAccepted    *bool        `gorm:"type:boolean" json:"offer_accepted,omitempty"`
	OfferAcceptedAt  *time.Time   `gorm:"type:timestamptz" json:"offer_accepted_at,omitempty"`
	Remarks          *string      `gorm:"type:text" json:"remarks,omitempty"`
	CreatedAt        time.Time    `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time    `gorm:"not null;default:now()" json:"updated_at"`
	CreatedBy        *uuid.UUID   `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy        *uuid.UUID   `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Application *AdmissionApplication `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
}

// TableName returns the table name for AdmissionDecision.
func (AdmissionDecision) TableName() string {
	return "admission_decisions"
}

// MeritListEntry represents a single entry in the merit list.
type MeritListEntry struct {
	Rank           int        `json:"rank"`
	ApplicationID  uuid.UUID  `json:"applicationId"`
	StudentName    string     `json:"studentName"`
	Score          float64    `json:"score"`
	TestScore      *float64   `json:"testScore,omitempty"`
	InterviewScore *float64   `json:"interviewScore,omitempty"`
	PreviousMarks  *float64   `json:"previousMarks,omitempty"`
	Status         string     `json:"status"` // pending, approved, waitlisted, rejected
	ParentPhone    string     `json:"parentPhone,omitempty"`
	ParentEmail    *string    `json:"parentEmail,omitempty"`
}

// MeritListEntries is a slice of MeritListEntry that implements Scanner/Valuer for JSONB.
type MeritListEntries []MeritListEntry

// Value implements the driver.Valuer interface for database storage.
func (m MeritListEntries) Value() (driver.Value, error) {
	if m == nil {
		return "[]", nil
	}
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface for database retrieval.
func (m *MeritListEntries) Scan(value interface{}) error {
	if value == nil {
		*m = MeritListEntries{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed for MeritListEntries")
	}

	return json.Unmarshal(bytes, m)
}

// MeritList represents a merit list snapshot for an admission session.
type MeritList struct {
	ID          uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"tenant_id"`
	SessionID   uuid.UUID        `gorm:"type:uuid;not null;index" json:"session_id"`
	ClassName   string           `gorm:"type:varchar(50);not null" json:"class_name"`
	TestID      *uuid.UUID       `gorm:"type:uuid" json:"test_id,omitempty"`
	GeneratedAt time.Time        `gorm:"type:timestamptz;not null;default:now()" json:"generated_at"`
	GeneratedBy *uuid.UUID       `gorm:"type:uuid" json:"generated_by,omitempty"`
	CutoffScore *float64         `gorm:"type:decimal(5,2)" json:"cutoff_score,omitempty"`
	Entries     MeritListEntries `gorm:"type:jsonb;not null;default:'[]'" json:"entries"`
	IsFinal     bool             `gorm:"type:boolean;default:false" json:"is_final"`
	CreatedAt   time.Time        `gorm:"not null;default:now()" json:"created_at"`

	// Relationships
	Session *AdmissionSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// TableName returns the table name for MeritList.
func (MeritList) TableName() string {
	return "merit_lists"
}

// TotalEntries returns the total number of entries in the merit list.
func (m *MeritList) TotalEntries() int {
	return len(m.Entries)
}

// EntriesAboveCutoff returns entries with score above the cutoff.
func (m *MeritList) EntriesAboveCutoff() []MeritListEntry {
	if m.CutoffScore == nil {
		return m.Entries
	}

	var above []MeritListEntry
	for _, entry := range m.Entries {
		if entry.Score >= *m.CutoffScore {
			above = append(above, entry)
		}
	}
	return above
}

// ============================================================================
// Story 3.5: Online Admission Application - Additional Types
// ============================================================================

// StageHistoryEntry represents a single entry in the application stage history.
type StageHistoryEntry struct {
	Stage     ApplicationStatus `json:"stage"`
	Timestamp time.Time         `json:"timestamp"`
	ChangedBy *uuid.UUID        `json:"changedBy,omitempty"`
	Remarks   string            `json:"remarks,omitempty"`
}

// StageHistory represents the history of stage transitions.
type StageHistory []StageHistoryEntry

// Value implements the driver.Valuer interface for database storage.
func (h StageHistory) Value() (driver.Value, error) {
	return json.Marshal(h)
}

// Scan implements the sql.Scanner interface for database retrieval.
func (h *StageHistory) Scan(value interface{}) error {
	if value == nil {
		*h = StageHistory{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed for StageHistory")
	}
	return json.Unmarshal(bytes, h)
}

// CanTransitionTo checks if the status can transition to the target status.
func (s ApplicationStatus) CanTransitionTo(target ApplicationStatus) bool {
	// Define valid transitions for each status
	validTransitions := map[ApplicationStatus][]ApplicationStatus{
		ApplicationStatusDraft: {
			ApplicationStatusSubmitted,
		},
		ApplicationStatusSubmitted: {
			ApplicationStatusUnderReview,
			ApplicationStatusDocumentsPending,
			ApplicationStatusDocumentsVerified,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusUnderReview: {
			ApplicationStatusDocumentsPending,
			ApplicationStatusDocumentsVerified,
			ApplicationStatusTestScheduled,
			ApplicationStatusInterviewScheduled,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusDocumentsPending: {
			ApplicationStatusDocumentsVerified,
			ApplicationStatusUnderReview,
			ApplicationStatusRejected,
		},
		ApplicationStatusDocumentsVerified: {
			ApplicationStatusTestScheduled,
			ApplicationStatusInterviewScheduled,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusTestScheduled: {
			ApplicationStatusTestCompleted,
			ApplicationStatusRejected,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusTestCompleted: {
			ApplicationStatusInterviewScheduled,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusInterviewScheduled: {
			ApplicationStatusInterviewCompleted,
			ApplicationStatusRejected,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusInterviewCompleted: {
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusShortlisted: {
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusApproved: {
			ApplicationStatusEnrolled,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusWaitlisted: {
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusRejected: {
			// Terminal state - no transitions allowed (except admin override)
		},
		ApplicationStatusEnrolled: {
			// Terminal state - no transitions allowed (except admin override)
		},
		ApplicationStatusWithdrawn: {
			// Terminal state - no transitions allowed (except admin override)
		},
	}

	allowed, exists := validTransitions[s]
	if !exists {
		return false
	}

	for _, valid := range allowed {
		if target == valid {
			return true
		}
	}
	return false
}

// GetValidTransitions returns the list of valid status transitions from the current status.
func (s ApplicationStatus) GetValidTransitions() []ApplicationStatus {
	validTransitions := map[ApplicationStatus][]ApplicationStatus{
		ApplicationStatusDraft: {
			ApplicationStatusSubmitted,
		},
		ApplicationStatusSubmitted: {
			ApplicationStatusUnderReview,
			ApplicationStatusDocumentsPending,
			ApplicationStatusDocumentsVerified,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusUnderReview: {
			ApplicationStatusDocumentsPending,
			ApplicationStatusDocumentsVerified,
			ApplicationStatusTestScheduled,
			ApplicationStatusInterviewScheduled,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusDocumentsPending: {
			ApplicationStatusDocumentsVerified,
			ApplicationStatusUnderReview,
			ApplicationStatusRejected,
		},
		ApplicationStatusDocumentsVerified: {
			ApplicationStatusTestScheduled,
			ApplicationStatusInterviewScheduled,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusTestScheduled: {
			ApplicationStatusTestCompleted,
			ApplicationStatusRejected,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusTestCompleted: {
			ApplicationStatusInterviewScheduled,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusInterviewScheduled: {
			ApplicationStatusInterviewCompleted,
			ApplicationStatusRejected,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusInterviewCompleted: {
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusShortlisted: {
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusApproved: {
			ApplicationStatusEnrolled,
			ApplicationStatusWithdrawn,
		},
		ApplicationStatusWaitlisted: {
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWithdrawn,
		},
	}

	return validTransitions[s]
}

// GetFullName returns the full name of the student from applicant details.
func (a *AdmissionApplication) GetFullName() string {
	return a.StudentName
}

// ParentRelation represents the relation type of a parent/guardian.
type ParentRelation string

// ParentRelation constants.
const (
	RelationFather   ParentRelation = "father"
	RelationMother   ParentRelation = "mother"
	RelationGuardian ParentRelation = "guardian"
)

// IsValid checks if the relation is a valid value.
func (r ParentRelation) IsValid() bool {
	switch r {
	case RelationFather, RelationMother, RelationGuardian:
		return true
	}
	return false
}

// DocumentType represents the type of document uploaded.
type DocumentType string

// DocumentType constants.
const (
	DocBirthCertificate    DocumentType = "birth_certificate"
	DocPhoto               DocumentType = "photo"
	DocAadhaarCard         DocumentType = "aadhaar_card"
	DocTransferCertificate DocumentType = "transfer_certificate"
	DocMarksheet           DocumentType = "marksheet"
	DocMedicalCertificate  DocumentType = "medical_certificate"
	DocAddressProof        DocumentType = "address_proof"
	DocCasteCertificate    DocumentType = "caste_certificate"
	DocIncomeCertificate   DocumentType = "income_certificate"
	DocOther               DocumentType = "other"
)

// IsValid checks if the document type is a valid value.
func (d DocumentType) IsValid() bool {
	switch d {
	case DocBirthCertificate, DocPhoto, DocAadhaarCard, DocTransferCertificate,
		DocMarksheet, DocMedicalCertificate, DocAddressProof, DocCasteCertificate,
		DocIncomeCertificate, DocOther:
		return true
	}
	return false
}

// ApplicationParent represents a parent/guardian associated with an application.
type ApplicationParent struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ApplicationID uuid.UUID      `gorm:"type:uuid;not null;index" json:"application_id"`
	Relation      ParentRelation `gorm:"size:20;not null" json:"relation"`
	Name          string         `gorm:"size:200;not null" json:"name"`
	Phone         string         `gorm:"size:20" json:"phone"`
	Email         string         `gorm:"size:255" json:"email"`
	Occupation    string         `gorm:"size:100" json:"occupation"`
	Education     string         `gorm:"size:100" json:"education"`
	AnnualIncome  string         `gorm:"size:50" json:"annual_income"`
	CreatedAt     time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName returns the table name for ApplicationParent.
func (ApplicationParent) TableName() string {
	return "application_parents"
}

// VerificationStatus represents the verification status of a document.
type VerificationStatus string

const (
	VerificationStatusPending         VerificationStatus = "pending"
	VerificationStatusVerified        VerificationStatus = "verified"
	VerificationStatusRejected        VerificationStatus = "rejected"
	VerificationStatusResubmitRequired VerificationStatus = "resubmit_required"
)

// ApplicationDocument represents a document uploaded with an application.
type ApplicationDocument struct {
	ID                  uuid.UUID          `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID            uuid.UUID          `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ApplicationID       uuid.UUID          `gorm:"type:uuid;not null;index" json:"application_id"`
	DocumentType        DocumentType       `gorm:"column:document_type;size:100;not null" json:"document_type"`
	FileName            string             `gorm:"column:document_name;size:255;not null" json:"file_name"`
	FileURL             string             `gorm:"column:file_path;size:500;not null" json:"file_url"`
	FileSize            int64              `gorm:"column:file_size" json:"file_size"`
	MimeType            string             `gorm:"column:mime_type;size:100" json:"mime_type"`
	VerificationStatus  VerificationStatus `gorm:"column:verification_status;size:20;not null;default:'pending'" json:"verification_status"`
	VerifiedBy          *uuid.UUID         `gorm:"type:uuid" json:"verified_by,omitempty"`
	VerifiedAt          *time.Time         `json:"verified_at,omitempty"`
	VerificationRemarks string             `gorm:"column:verification_remarks;type:text" json:"verification_remarks,omitempty"`
	CreatedAt           time.Time          `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt           time.Time          `gorm:"not null;default:now()" json:"updated_at"`
}

// IsVerified returns true if the document is verified.
func (d *ApplicationDocument) IsVerified() bool {
	return d.VerificationStatus == VerificationStatusVerified
}

// RejectionReason returns the verification remarks (alias for backwards compatibility).
func (d *ApplicationDocument) RejectionReason() string {
	return d.VerificationRemarks
}

// TableName returns the table name for ApplicationDocument.
func (ApplicationDocument) TableName() string {
	return "application_documents"
}

// ApplicationNumberSequence tracks sequence numbers for generating unique application numbers.
type ApplicationNumberSequence struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	DatePrefix   string    `gorm:"size:8;not null" json:"date_prefix"`
	LastSequence int       `gorm:"not null;default:0" json:"last_sequence"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName returns the table name for ApplicationNumberSequence.
func (ApplicationNumberSequence) TableName() string {
	return "application_number_sequences"
}
