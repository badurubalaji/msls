// Package admission provides admission management services including
// application reviews and entrance tests.
package admission

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ApplicationStatus represents the status of an admission application.
type ApplicationStatus string

// Application status constants.
const (
	ApplicationStatusDraft            ApplicationStatus = "draft"
	ApplicationStatusSubmitted        ApplicationStatus = "submitted"
	ApplicationStatusUnderReview      ApplicationStatus = "under_review"
	ApplicationStatusDocumentsPending ApplicationStatus = "documents_pending"
	ApplicationStatusTestScheduled    ApplicationStatus = "test_scheduled"
	ApplicationStatusTestCompleted    ApplicationStatus = "test_completed"
	ApplicationStatusShortlisted      ApplicationStatus = "shortlisted"
	ApplicationStatusApproved         ApplicationStatus = "approved"
	ApplicationStatusRejected         ApplicationStatus = "rejected"
	ApplicationStatusWaitlisted       ApplicationStatus = "waitlisted"
	ApplicationStatusEnrolled         ApplicationStatus = "enrolled"
)

// IsValid checks if the application status is valid.
func (s ApplicationStatus) IsValid() bool {
	switch s {
	case ApplicationStatusDraft, ApplicationStatusSubmitted, ApplicationStatusUnderReview,
		ApplicationStatusDocumentsPending, ApplicationStatusTestScheduled, ApplicationStatusTestCompleted,
		ApplicationStatusShortlisted, ApplicationStatusApproved, ApplicationStatusRejected,
		ApplicationStatusWaitlisted, ApplicationStatusEnrolled:
		return true
	}
	return false
}

// String returns the string representation of the application status.
func (s ApplicationStatus) String() string {
	return string(s)
}

// VerificationStatus represents the verification status of a document.
type VerificationStatus string

// Verification status constants.
const (
	VerificationStatusPending         VerificationStatus = "pending"
	VerificationStatusVerified        VerificationStatus = "verified"
	VerificationStatusRejected        VerificationStatus = "rejected"
	VerificationStatusResubmitRequired VerificationStatus = "resubmit_required"
)

// IsValid checks if the verification status is valid.
func (s VerificationStatus) IsValid() bool {
	switch s {
	case VerificationStatusPending, VerificationStatusVerified,
		VerificationStatusRejected, VerificationStatusResubmitRequired:
		return true
	}
	return false
}

// TestRegistrationStatus represents the status of a test registration.
type TestRegistrationStatus string

// Test registration status constants.
const (
	TestRegStatusRegistered         TestRegistrationStatus = "registered"
	TestRegStatusHallTicketGenerated TestRegistrationStatus = "hall_ticket_generated"
	TestRegStatusAppeared           TestRegistrationStatus = "appeared"
	TestRegStatusAbsent             TestRegistrationStatus = "absent"
	TestRegStatusCancelled          TestRegistrationStatus = "cancelled"
)

// IsValid checks if the test registration status is valid.
func (s TestRegistrationStatus) IsValid() bool {
	switch s {
	case TestRegStatusRegistered, TestRegStatusHallTicketGenerated,
		TestRegStatusAppeared, TestRegStatusAbsent, TestRegStatusCancelled:
		return true
	}
	return false
}

// TestResult represents the result of a test.
type TestResult string

// Test result constants.
const (
	TestResultPass        TestResult = "pass"
	TestResultFail        TestResult = "fail"
	TestResultMerit       TestResult = "merit"
	TestResultDistinction TestResult = "distinction"
)

// IsValid checks if the test result is valid.
func (r TestResult) IsValid() bool {
	switch r {
	case TestResultPass, TestResultFail, TestResultMerit, TestResultDistinction:
		return true
	}
	return false
}

// ReviewType represents the type of application review.
type ReviewType string

// Review type constants.
const (
	ReviewTypeInitialScreening    ReviewType = "initial_screening"
	ReviewTypeDocumentVerification ReviewType = "document_verification"
	ReviewTypeAcademicReview      ReviewType = "academic_review"
	ReviewTypeInterview           ReviewType = "interview"
	ReviewTypeFinalDecision       ReviewType = "final_decision"
)

// IsValid checks if the review type is valid.
func (r ReviewType) IsValid() bool {
	switch r {
	case ReviewTypeInitialScreening, ReviewTypeDocumentVerification,
		ReviewTypeAcademicReview, ReviewTypeInterview, ReviewTypeFinalDecision:
		return true
	}
	return false
}

// ReviewStatus represents the status of a review.
type ReviewStatus string

// Review status constants.
const (
	ReviewStatusApproved    ReviewStatus = "approved"
	ReviewStatusRejected    ReviewStatus = "rejected"
	ReviewStatusPendingInfo ReviewStatus = "pending_info"
	ReviewStatusEscalated   ReviewStatus = "escalated"
)

// IsValid checks if the review status is valid.
func (s ReviewStatus) IsValid() bool {
	switch s {
	case ReviewStatusApproved, ReviewStatusRejected,
		ReviewStatusPendingInfo, ReviewStatusEscalated:
		return true
	}
	return false
}

// EntranceTestStatus represents the status of an entrance test.
type EntranceTestStatus string

// Entrance test status constants.
const (
	TestStatusScheduled  EntranceTestStatus = "scheduled"
	TestStatusInProgress EntranceTestStatus = "in_progress"
	TestStatusCompleted  EntranceTestStatus = "completed"
	TestStatusCancelled  EntranceTestStatus = "cancelled"
)

// IsValid checks if the test status is valid.
func (s EntranceTestStatus) IsValid() bool {
	switch s {
	case TestStatusScheduled, TestStatusInProgress, TestStatusCompleted, TestStatusCancelled:
		return true
	}
	return false
}

// StringArray represents a JSON array of strings.
type StringArray []string

// Value implements the driver.Valuer interface.
func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return "[]", nil
	}
	return json.Marshal(sa)
}

// Scan implements the sql.Scanner interface.
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = StringArray{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, sa)
}

// SubjectMarks represents marks configuration for a subject.
type SubjectMarks struct {
	Subject  string          `json:"subject"`
	MaxMarks decimal.Decimal `json:"maxMarks"`
}

// TestSubjects represents a list of subjects for a test.
type TestSubjects []SubjectMarks

// Value implements the driver.Valuer interface.
func (ts TestSubjects) Value() (driver.Value, error) {
	if ts == nil {
		return "[]", nil
	}
	return json.Marshal(ts)
}

// Scan implements the sql.Scanner interface.
func (ts *TestSubjects) Scan(value interface{}) error {
	if value == nil {
		*ts = TestSubjects{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, ts)
}

// MarksMap represents marks obtained per subject.
type MarksMap map[string]decimal.Decimal

// Value implements the driver.Valuer interface.
func (mm MarksMap) Value() (driver.Value, error) {
	if mm == nil {
		return "{}", nil
	}
	return json.Marshal(mm)
}

// Scan implements the sql.Scanner interface.
func (mm *MarksMap) Scan(value interface{}) error {
	if value == nil {
		*mm = MarksMap{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, mm)
}

// ExtraData represents additional custom fields.
type ExtraData map[string]interface{}

// Value implements the driver.Valuer interface.
func (ed ExtraData) Value() (driver.Value, error) {
	if ed == nil {
		return "{}", nil
	}
	return json.Marshal(ed)
}

// Scan implements the sql.Scanner interface.
func (ed *ExtraData) Scan(value interface{}) error {
	if value == nil {
		*ed = ExtraData{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, ed)
}

// AdmissionApplication represents an admission application in the database.
type AdmissionApplication struct {
	ID                uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID          uuid.UUID         `gorm:"type:uuid;not null;index" json:"tenant_id"`
	SessionID         uuid.UUID         `gorm:"type:uuid;not null;index" json:"session_id"`
	BranchID          *uuid.UUID        `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	EnquiryID         *uuid.UUID        `gorm:"type:uuid;index" json:"enquiry_id,omitempty"`
	ApplicationNumber string            `gorm:"size:50;not null" json:"application_number"`

	// Student Information
	StudentName         string     `gorm:"size:200;not null" json:"student_name"`
	DateOfBirth         time.Time  `gorm:"type:date;not null" json:"date_of_birth"`
	Gender              string     `gorm:"size:20;not null" json:"gender"`
	BloodGroup          string     `gorm:"size:10" json:"blood_group,omitempty"`
	Nationality         string     `gorm:"size:100;default:'Indian'" json:"nationality"`
	Religion            string     `gorm:"size:100" json:"religion,omitempty"`
	Caste               string     `gorm:"size:100" json:"caste,omitempty"`
	Category            string     `gorm:"size:50" json:"category,omitempty"`
	MotherTongue        string     `gorm:"size:100" json:"mother_tongue,omitempty"`
	AadharNumber        string     `gorm:"size:12" json:"aadhar_number,omitempty"`

	// Academic Information
	ClassApplying       string          `gorm:"size:50;not null" json:"class_applying"`
	PreviousSchool      string          `gorm:"size:200" json:"previous_school,omitempty"`
	PreviousClass       string          `gorm:"size:50" json:"previous_class,omitempty"`
	PreviousPercentage  *decimal.Decimal `gorm:"type:decimal(5,2)" json:"previous_percentage,omitempty"`
	MediumOfInstruction string          `gorm:"size:50" json:"medium_of_instruction,omitempty"`

	// Contact Information
	AddressLine1 string `gorm:"size:255" json:"address_line1,omitempty"`
	AddressLine2 string `gorm:"size:255" json:"address_line2,omitempty"`
	City         string `gorm:"size:100" json:"city,omitempty"`
	State        string `gorm:"size:100" json:"state,omitempty"`
	PostalCode   string `gorm:"size:20" json:"postal_code,omitempty"`
	Country      string `gorm:"size:100;default:'India'" json:"country"`

	// Parent/Guardian Information
	FatherName          string `gorm:"size:200" json:"father_name,omitempty"`
	FatherPhone         string `gorm:"size:20" json:"father_phone,omitempty"`
	FatherEmail         string `gorm:"size:255" json:"father_email,omitempty"`
	FatherOccupation    string `gorm:"size:200" json:"father_occupation,omitempty"`
	FatherQualification string `gorm:"size:100" json:"father_qualification,omitempty"`
	MotherName          string `gorm:"size:200" json:"mother_name,omitempty"`
	MotherPhone         string `gorm:"size:20" json:"mother_phone,omitempty"`
	MotherEmail         string `gorm:"size:255" json:"mother_email,omitempty"`
	MotherOccupation    string `gorm:"size:200" json:"mother_occupation,omitempty"`
	MotherQualification string `gorm:"size:100" json:"mother_qualification,omitempty"`
	GuardianName        string `gorm:"size:200" json:"guardian_name,omitempty"`
	GuardianPhone       string `gorm:"size:20" json:"guardian_phone,omitempty"`
	GuardianEmail       string `gorm:"size:255" json:"guardian_email,omitempty"`
	GuardianRelation    string `gorm:"size:50" json:"guardian_relation,omitempty"`

	// Application Details
	Status        ApplicationStatus `gorm:"size:20;not null;default:'draft'" json:"status"`
	SubmittedAt   *time.Time        `gorm:"type:timestamptz" json:"submitted_at,omitempty"`
	Remarks       string            `gorm:"type:text" json:"remarks,omitempty"`
	InternalNotes string            `gorm:"type:text" json:"internal_notes,omitempty"`
	Priority      int               `gorm:"default:0" json:"priority"`

	// Payment Information
	FeePaid          bool       `gorm:"default:false" json:"fee_paid"`
	PaymentReference string     `gorm:"size:100" json:"payment_reference,omitempty"`
	PaymentDate      *time.Time `gorm:"type:timestamptz" json:"payment_date,omitempty"`

	// Additional Data
	ExtraData ExtraData `gorm:"type:jsonb;default:'{}'" json:"extra_data"`

	// Audit Fields
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
	DeletedAt *time.Time `gorm:"type:timestamptz" json:"deleted_at,omitempty"`

	// Relationships
	Documents []ApplicationDocument `gorm:"foreignKey:ApplicationID" json:"documents,omitempty"`
	Reviews   []ApplicationReview   `gorm:"foreignKey:ApplicationID" json:"reviews,omitempty"`
}

// TableName specifies the table name for AdmissionApplication.
func (AdmissionApplication) TableName() string {
	return "admission_applications"
}

// ApplicationDocument represents a document submitted with an application.
type ApplicationDocument struct {
	ID                 uuid.UUID          `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID           uuid.UUID          `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ApplicationID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"application_id"`
	DocumentType       string             `gorm:"size:100;not null" json:"document_type"`
	DocumentName       string             `gorm:"size:255;not null" json:"document_name"`
	FilePath           string             `gorm:"size:500;not null" json:"file_path"`
	FileSize           int                `gorm:"" json:"file_size,omitempty"`
	MimeType           string             `gorm:"size:100" json:"mime_type,omitempty"`
	VerificationStatus VerificationStatus `gorm:"size:20;not null;default:'pending'" json:"verification_status"`
	VerifiedBy         *uuid.UUID         `gorm:"type:uuid" json:"verified_by,omitempty"`
	VerifiedAt         *time.Time         `gorm:"type:timestamptz" json:"verified_at,omitempty"`
	VerificationRemarks string            `gorm:"type:text" json:"verification_remarks,omitempty"`
	CreatedAt          time.Time          `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt          time.Time          `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName specifies the table name for ApplicationDocument.
func (ApplicationDocument) TableName() string {
	return "application_documents"
}

// EntranceTest represents an entrance test for admissions.
type EntranceTest struct {
	ID                uuid.UUID          `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID          uuid.UUID          `gorm:"type:uuid;not null;index" json:"tenant_id"`
	SessionID         uuid.UUID          `gorm:"type:uuid;not null;index" json:"session_id"`
	TestName          string             `gorm:"size:200;not null" json:"test_name"`
	TestDate          time.Time          `gorm:"type:date;not null" json:"test_date"`
	StartTime         string             `gorm:"type:time;not null" json:"start_time"`
	DurationMinutes   int                `gorm:"not null;default:60" json:"duration_minutes"`
	Venue             string             `gorm:"size:200" json:"venue,omitempty"`
	ClassNames        StringArray        `gorm:"type:jsonb;default:'[]'" json:"class_names"`
	MaxCandidates     int                `gorm:"default:100" json:"max_candidates"`
	Status            EntranceTestStatus `gorm:"size:20;default:'scheduled'" json:"status"`
	Subjects          TestSubjects       `gorm:"type:jsonb;default:'[]'" json:"subjects"`
	Instructions      string             `gorm:"type:text" json:"instructions,omitempty"`
	PassingPercentage decimal.Decimal    `gorm:"type:decimal(5,2);default:33.00" json:"passing_percentage"`
	CreatedAt         time.Time          `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time          `gorm:"not null;default:now()" json:"updated_at"`
	CreatedBy         *uuid.UUID         `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID         `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Registrations []TestRegistration `gorm:"foreignKey:TestID" json:"registrations,omitempty"`
}

// TableName specifies the table name for EntranceTest.
func (EntranceTest) TableName() string {
	return "entrance_tests"
}

// TestRegistration represents a student's registration for an entrance test.
type TestRegistration struct {
	ID                    uuid.UUID              `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID              uuid.UUID              `gorm:"type:uuid;not null;index" json:"tenant_id"`
	TestID                uuid.UUID              `gorm:"type:uuid;not null;index" json:"test_id"`
	ApplicationID         uuid.UUID              `gorm:"type:uuid;not null;index" json:"application_id"`
	RollNumber            string                 `gorm:"size:20" json:"roll_number,omitempty"`
	Status                TestRegistrationStatus `gorm:"size:20;not null;default:'registered'" json:"status"`
	Marks                 MarksMap               `gorm:"type:jsonb;default:'{}'" json:"marks"`
	TotalMarks            *decimal.Decimal       `gorm:"type:decimal(6,2)" json:"total_marks,omitempty"`
	MaxMarks              *decimal.Decimal       `gorm:"type:decimal(6,2)" json:"max_marks,omitempty"`
	Percentage            *decimal.Decimal       `gorm:"type:decimal(5,2)" json:"percentage,omitempty"`
	Result                *TestResult            `gorm:"size:20" json:"result,omitempty"`
	Remarks               string                 `gorm:"type:text" json:"remarks,omitempty"`
	HallTicketGeneratedAt *time.Time             `gorm:"type:timestamptz" json:"hall_ticket_generated_at,omitempty"`
	CreatedAt             time.Time              `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt             time.Time              `gorm:"not null;default:now()" json:"updated_at"`

	// Relationships
	Test        *EntranceTest         `gorm:"foreignKey:TestID" json:"test,omitempty"`
	Application *AdmissionApplication `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
}

// TableName specifies the table name for TestRegistration.
func (TestRegistration) TableName() string {
	return "test_registrations"
}

// ApplicationReview represents a review of an admission application.
type ApplicationReview struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"`
	TenantID      uuid.UUID    `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ApplicationID uuid.UUID    `gorm:"type:uuid;not null;index" json:"application_id"`
	ReviewerID    uuid.UUID    `gorm:"type:uuid;not null" json:"reviewer_id"`
	ReviewType    ReviewType   `gorm:"size:50;not null" json:"review_type"`
	Status        ReviewStatus `gorm:"size:20;not null" json:"status"`
	Comments      string       `gorm:"type:text" json:"comments,omitempty"`
	CreatedAt     time.Time    `gorm:"not null;default:now()" json:"created_at"`
}

// TableName specifies the table name for ApplicationReview.
func (ApplicationReview) TableName() string {
	return "application_reviews"
}
