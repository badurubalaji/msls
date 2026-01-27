package models

import (
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Enums
// =============================================================================

// AllergyType represents the type of allergy
type AllergyType string

const (
	AllergyTypeFood          AllergyType = "food"
	AllergyTypeMedication    AllergyType = "medication"
	AllergyTypeEnvironmental AllergyType = "environmental"
	AllergyTypeInsect        AllergyType = "insect"
	AllergyTypeOther         AllergyType = "other"
)

func (a AllergyType) IsValid() bool {
	switch a {
	case AllergyTypeFood, AllergyTypeMedication, AllergyTypeEnvironmental, AllergyTypeInsect, AllergyTypeOther:
		return true
	}
	return false
}

// AllergySeverity represents the severity of an allergy
type AllergySeverity string

const (
	AllergySeverityMild            AllergySeverity = "mild"
	AllergySeverityModerate        AllergySeverity = "moderate"
	AllergySeveritySevere          AllergySeverity = "severe"
	AllergySeverityLifeThreatening AllergySeverity = "life_threatening"
)

func (s AllergySeverity) IsValid() bool {
	switch s {
	case AllergySeverityMild, AllergySeverityModerate, AllergySeveritySevere, AllergySeverityLifeThreatening:
		return true
	}
	return false
}

// ConditionType represents the type of chronic condition
type ConditionType string

const (
	ConditionTypeRespiratory  ConditionType = "respiratory"
	ConditionTypeCardiac      ConditionType = "cardiac"
	ConditionTypeNeurological ConditionType = "neurological"
	ConditionTypeEndocrine    ConditionType = "endocrine"
	ConditionTypeOther        ConditionType = "other"
)

func (c ConditionType) IsValid() bool {
	switch c {
	case ConditionTypeRespiratory, ConditionTypeCardiac, ConditionTypeNeurological, ConditionTypeEndocrine, ConditionTypeOther:
		return true
	}
	return false
}

// ConditionSeverity represents the severity of a chronic condition
type ConditionSeverity string

const (
	ConditionSeverityMild     ConditionSeverity = "mild"
	ConditionSeverityModerate ConditionSeverity = "moderate"
	ConditionSeveritySevere   ConditionSeverity = "severe"
)

func (s ConditionSeverity) IsValid() bool {
	switch s {
	case ConditionSeverityMild, ConditionSeverityModerate, ConditionSeveritySevere:
		return true
	}
	return false
}

// MedicationFrequency represents how often medication is taken
type MedicationFrequency string

const (
	MedicationFrequencyDaily      MedicationFrequency = "daily"
	MedicationFrequencyTwiceDaily MedicationFrequency = "twice_daily"
	MedicationFrequencyThreeDaily MedicationFrequency = "three_daily"
	MedicationFrequencyWeekly     MedicationFrequency = "weekly"
	MedicationFrequencyAsNeeded   MedicationFrequency = "as_needed"
	MedicationFrequencyOther      MedicationFrequency = "other"
)

func (f MedicationFrequency) IsValid() bool {
	switch f {
	case MedicationFrequencyDaily, MedicationFrequencyTwiceDaily, MedicationFrequencyThreeDaily,
		MedicationFrequencyWeekly, MedicationFrequencyAsNeeded, MedicationFrequencyOther:
		return true
	}
	return false
}

// MedicationRoute represents how medication is administered
type MedicationRoute string

const (
	MedicationRouteOral      MedicationRoute = "oral"
	MedicationRouteInjection MedicationRoute = "injection"
	MedicationRouteInhaler   MedicationRoute = "inhaler"
	MedicationRouteTopical   MedicationRoute = "topical"
	MedicationRouteDrops     MedicationRoute = "drops"
	MedicationRouteOther     MedicationRoute = "other"
)

func (r MedicationRoute) IsValid() bool {
	switch r {
	case MedicationRouteOral, MedicationRouteInjection, MedicationRouteInhaler,
		MedicationRouteTopical, MedicationRouteDrops, MedicationRouteOther:
		return true
	}
	return false
}

// VaccineType represents the type of vaccine
type VaccineType string

const (
	VaccineTypeRequired VaccineType = "required"
	VaccineTypeOptional VaccineType = "optional"
	VaccineTypeBooster  VaccineType = "booster"
)

func (v VaccineType) IsValid() bool {
	switch v {
	case VaccineTypeRequired, VaccineTypeOptional, VaccineTypeBooster:
		return true
	}
	return false
}

// IncidentType represents the type of medical incident
type IncidentType string

const (
	IncidentTypeIllness   IncidentType = "illness"
	IncidentTypeInjury    IncidentType = "injury"
	IncidentTypeEmergency IncidentType = "emergency"
	IncidentTypeOther     IncidentType = "other"
)

func (i IncidentType) IsValid() bool {
	switch i {
	case IncidentTypeIllness, IncidentTypeInjury, IncidentTypeEmergency, IncidentTypeOther:
		return true
	}
	return false
}

// =============================================================================
// Models
// =============================================================================

// StudentHealthProfile represents the main health profile for a student
type StudentHealthProfile struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:uq_student_health_profile"`

	// Basic health info
	BloodGroup    *string  `gorm:"type:varchar(5)"`
	HeightCm      *float64 `gorm:"type:decimal(5,2)"`
	WeightKg      *float64 `gorm:"type:decimal(5,2)"`
	VisionLeft    *string  `gorm:"type:varchar(10)"`
	VisionRight   *string  `gorm:"type:varchar(10)"`
	HearingStatus *string  `gorm:"type:varchar(20);default:'normal'"`

	// Medical history
	MedicalNotes *string `gorm:"type:text"`

	// Insurance
	InsuranceProvider     *string    `gorm:"type:varchar(100)"`
	InsurancePolicyNumber *string    `gorm:"type:varchar(50)"`
	InsuranceExpiry       *time.Time `gorm:"type:date"`

	// Emergency medical info
	PreferredHospital *string `gorm:"type:varchar(200)"`
	FamilyDoctorName  *string `gorm:"type:varchar(100)"`
	FamilyDoctorPhone *string `gorm:"type:varchar(15)"`

	// Metadata
	LastCheckupDate *time.Time `gorm:"type:date"`
	CreatedAt       time.Time  `gorm:"not null;default:now()"`
	UpdatedAt       time.Time  `gorm:"not null;default:now()"`
	CreatedBy       *uuid.UUID `gorm:"type:uuid"`
	UpdatedBy       *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID"`
}

func (StudentHealthProfile) TableName() string {
	return "student_health_profiles"
}

// StudentAllergy represents an allergy record for a student
type StudentAllergy struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Allergy details
	Allergen    string          `gorm:"type:varchar(100);not null"`
	AllergyType AllergyType     `gorm:"type:varchar(30);not null"`
	Severity    AllergySeverity `gorm:"type:varchar(20);not null"`

	// Reaction and treatment
	ReactionDescription   *string `gorm:"type:text"`
	TreatmentInstructions *string `gorm:"type:text"`

	// Medication for emergency
	EmergencyMedication *string `gorm:"type:varchar(100)"`

	// Status
	DiagnosedDate *time.Time `gorm:"type:date"`
	IsActive      bool       `gorm:"not null;default:true"`
	Notes         *string    `gorm:"type:text"`

	// Metadata
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID"`
}

func (StudentAllergy) TableName() string {
	return "student_allergies"
}

// StudentChronicCondition represents a chronic condition record for a student
type StudentChronicCondition struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Condition details
	ConditionName string            `gorm:"type:varchar(100);not null"`
	ConditionType ConditionType     `gorm:"type:varchar(30);not null"`
	Severity      ConditionSeverity `gorm:"type:varchar(20);not null"`

	// Management
	ManagementPlan *string `gorm:"type:text"`
	Restrictions   *string `gorm:"type:text"`
	Triggers       *string `gorm:"type:text"`

	// Status
	DiagnosedDate *time.Time `gorm:"type:date"`
	IsActive      bool       `gorm:"not null;default:true"`
	Notes         *string    `gorm:"type:text"`

	// Metadata
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID"`
}

func (StudentChronicCondition) TableName() string {
	return "student_chronic_conditions"
}

// StudentMedication represents a medication record for a student
type StudentMedication struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Medication details
	MedicationName string              `gorm:"type:varchar(100);not null"`
	Dosage         string              `gorm:"type:varchar(50);not null"`
	Frequency      MedicationFrequency `gorm:"type:varchar(50);not null"`
	Route          MedicationRoute     `gorm:"type:varchar(30);not null"`

	// Purpose and instructions
	Purpose             *string `gorm:"type:varchar(200)"`
	SpecialInstructions *string `gorm:"type:text"`

	// Timing
	StartDate time.Time  `gorm:"type:date;not null"`
	EndDate   *time.Time `gorm:"type:date"`

	// Administration at school
	AdministeredAtSchool     bool    `gorm:"not null;default:false"`
	SchoolAdministrationTime *string `gorm:"type:varchar(50)"`

	// Prescribing doctor
	PrescribingDoctor *string    `gorm:"type:varchar(100)"`
	PrescriptionDate  *time.Time `gorm:"type:date"`

	// Status
	IsActive bool    `gorm:"not null;default:true"`
	Notes    *string `gorm:"type:text"`

	// Metadata
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID"`
}

func (StudentMedication) TableName() string {
	return "student_medications"
}

// StudentVaccination represents a vaccination record for a student
type StudentVaccination struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Vaccination details
	VaccineName string       `gorm:"type:varchar(100);not null"`
	VaccineType *VaccineType `gorm:"type:varchar(50)"`
	DoseNumber  int          `gorm:"not null;default:1"`

	// Administration
	AdministeredDate time.Time `gorm:"type:date;not null"`
	AdministeredBy   *string   `gorm:"type:varchar(100)"`
	AdministrationSite *string `gorm:"type:varchar(50)"`
	BatchNumber      *string   `gorm:"type:varchar(50)"`

	// Next dose
	NextDueDate *time.Time `gorm:"type:date"`

	// Reaction
	HadReaction         bool    `gorm:"not null;default:false"`
	ReactionDescription *string `gorm:"type:text"`

	// Certificate
	CertificateUrl *string `gorm:"type:text"`

	// Status
	IsVerified bool       `gorm:"not null;default:false"`
	VerifiedBy *uuid.UUID `gorm:"type:uuid"`
	VerifiedAt *time.Time `gorm:"type:timestamptz"`
	Notes      *string    `gorm:"type:text"`

	// Metadata
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	CreatedBy *uuid.UUID `gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID"`
}

func (StudentVaccination) TableName() string {
	return "student_vaccinations"
}

// StudentMedicalIncident represents a medical incident record for a student
type StudentMedicalIncident struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Incident details
	IncidentDate time.Time    `gorm:"type:date;not null"`
	IncidentTime string       `gorm:"type:time;not null"`
	Location     *string      `gorm:"type:varchar(100)"`
	IncidentType IncidentType `gorm:"type:varchar(30);not null"`

	// Description
	Description string  `gorm:"type:text;not null"`
	Symptoms    *string `gorm:"type:text"`

	// Response
	FirstAidGiven       bool    `gorm:"not null;default:false"`
	FirstAidDescription *string `gorm:"type:text"`
	ActionTaken         string  `gorm:"type:text;not null"`

	// Follow-up
	ParentNotified   bool       `gorm:"not null;default:false"`
	ParentNotifiedAt *time.Time `gorm:"type:timestamptz"`
	ParentNotifiedBy *uuid.UUID `gorm:"type:uuid"`

	HospitalVisitRequired bool       `gorm:"not null;default:false"`
	HospitalName          *string    `gorm:"type:varchar(200)"`
	HospitalVisitDate     *time.Time `gorm:"type:date"`

	// Recovery
	StudentSentHome   bool       `gorm:"not null;default:false"`
	ReturnToClassTime *string    `gorm:"type:time"`
	FollowUpRequired  bool       `gorm:"not null;default:false"`
	FollowUpDate      *time.Time `gorm:"type:date"`
	FollowUpNotes     *string    `gorm:"type:text"`

	// Outcome
	Outcome *string `gorm:"type:text"`

	// Reported by
	ReportedBy uuid.UUID `gorm:"type:uuid;not null"`
	Notes      *string   `gorm:"type:text"`

	// Metadata
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Student  Student `gorm:"foreignKey:StudentID"`
	Reporter User    `gorm:"foreignKey:ReportedBy"`
}

func (StudentMedicalIncident) TableName() string {
	return "student_medical_incidents"
}
