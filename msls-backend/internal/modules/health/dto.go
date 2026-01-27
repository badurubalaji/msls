package health

import (
	"time"
)

// =============================================================================
// Health Profile DTOs
// =============================================================================

type HealthProfileResponse struct {
	ID        string `json:"id"`
	StudentID string `json:"studentId"`

	// Basic health info
	BloodGroup    *string  `json:"bloodGroup,omitempty"`
	HeightCm      *float64 `json:"heightCm,omitempty"`
	WeightKg      *float64 `json:"weightKg,omitempty"`
	VisionLeft    *string  `json:"visionLeft,omitempty"`
	VisionRight   *string  `json:"visionRight,omitempty"`
	HearingStatus *string  `json:"hearingStatus,omitempty"`

	// Medical history
	MedicalNotes *string `json:"medicalNotes,omitempty"`

	// Insurance
	InsuranceProvider     *string `json:"insuranceProvider,omitempty"`
	InsurancePolicyNumber *string `json:"insurancePolicyNumber,omitempty"`
	InsuranceExpiry       *string `json:"insuranceExpiry,omitempty"`

	// Emergency medical info
	PreferredHospital *string `json:"preferredHospital,omitempty"`
	FamilyDoctorName  *string `json:"familyDoctorName,omitempty"`
	FamilyDoctorPhone *string `json:"familyDoctorPhone,omitempty"`

	// Metadata
	LastCheckupDate *string `json:"lastCheckupDate,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type CreateHealthProfileRequest struct {
	BloodGroup            *string  `json:"bloodGroup"`
	HeightCm              *float64 `json:"heightCm"`
	WeightKg              *float64 `json:"weightKg"`
	VisionLeft            *string  `json:"visionLeft"`
	VisionRight           *string  `json:"visionRight"`
	HearingStatus         *string  `json:"hearingStatus"`
	MedicalNotes          *string  `json:"medicalNotes"`
	InsuranceProvider     *string  `json:"insuranceProvider"`
	InsurancePolicyNumber *string  `json:"insurancePolicyNumber"`
	InsuranceExpiry       *string  `json:"insuranceExpiry"`
	PreferredHospital     *string  `json:"preferredHospital"`
	FamilyDoctorName      *string  `json:"familyDoctorName"`
	FamilyDoctorPhone     *string  `json:"familyDoctorPhone"`
	LastCheckupDate       *string  `json:"lastCheckupDate"`
}

type UpdateHealthProfileRequest struct {
	BloodGroup            *string  `json:"bloodGroup"`
	HeightCm              *float64 `json:"heightCm"`
	WeightKg              *float64 `json:"weightKg"`
	VisionLeft            *string  `json:"visionLeft"`
	VisionRight           *string  `json:"visionRight"`
	HearingStatus         *string  `json:"hearingStatus"`
	MedicalNotes          *string  `json:"medicalNotes"`
	InsuranceProvider     *string  `json:"insuranceProvider"`
	InsurancePolicyNumber *string  `json:"insurancePolicyNumber"`
	InsuranceExpiry       *string  `json:"insuranceExpiry"`
	PreferredHospital     *string  `json:"preferredHospital"`
	FamilyDoctorName      *string  `json:"familyDoctorName"`
	FamilyDoctorPhone     *string  `json:"familyDoctorPhone"`
	LastCheckupDate       *string  `json:"lastCheckupDate"`
}

// =============================================================================
// Allergy DTOs
// =============================================================================

type AllergyResponse struct {
	ID        string `json:"id"`
	StudentID string `json:"studentId"`

	Allergen              string  `json:"allergen"`
	AllergyType           string  `json:"allergyType"`
	Severity              string  `json:"severity"`
	ReactionDescription   *string `json:"reactionDescription,omitempty"`
	TreatmentInstructions *string `json:"treatmentInstructions,omitempty"`
	EmergencyMedication   *string `json:"emergencyMedication,omitempty"`
	DiagnosedDate         *string `json:"diagnosedDate,omitempty"`
	IsActive              bool    `json:"isActive"`
	Notes                 *string `json:"notes,omitempty"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AllergyListResponse struct {
	Allergies []AllergyResponse `json:"allergies"`
	Total     int               `json:"total"`
}

type CreateAllergyRequest struct {
	Allergen              string  `json:"allergen" binding:"required"`
	AllergyType           string  `json:"allergyType" binding:"required"`
	Severity              string  `json:"severity" binding:"required"`
	ReactionDescription   *string `json:"reactionDescription"`
	TreatmentInstructions *string `json:"treatmentInstructions"`
	EmergencyMedication   *string `json:"emergencyMedication"`
	DiagnosedDate         *string `json:"diagnosedDate"`
	Notes                 *string `json:"notes"`
}

type UpdateAllergyRequest struct {
	Allergen              *string `json:"allergen"`
	AllergyType           *string `json:"allergyType"`
	Severity              *string `json:"severity"`
	ReactionDescription   *string `json:"reactionDescription"`
	TreatmentInstructions *string `json:"treatmentInstructions"`
	EmergencyMedication   *string `json:"emergencyMedication"`
	DiagnosedDate         *string `json:"diagnosedDate"`
	IsActive              *bool   `json:"isActive"`
	Notes                 *string `json:"notes"`
}

// =============================================================================
// Chronic Condition DTOs
// =============================================================================

type ConditionResponse struct {
	ID        string `json:"id"`
	StudentID string `json:"studentId"`

	ConditionName  string  `json:"conditionName"`
	ConditionType  string  `json:"conditionType"`
	Severity       string  `json:"severity"`
	ManagementPlan *string `json:"managementPlan,omitempty"`
	Restrictions   *string `json:"restrictions,omitempty"`
	Triggers       *string `json:"triggers,omitempty"`
	DiagnosedDate  *string `json:"diagnosedDate,omitempty"`
	IsActive       bool    `json:"isActive"`
	Notes          *string `json:"notes,omitempty"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ConditionListResponse struct {
	Conditions []ConditionResponse `json:"conditions"`
	Total      int                 `json:"total"`
}

type CreateConditionRequest struct {
	ConditionName  string  `json:"conditionName" binding:"required"`
	ConditionType  string  `json:"conditionType" binding:"required"`
	Severity       string  `json:"severity" binding:"required"`
	ManagementPlan *string `json:"managementPlan"`
	Restrictions   *string `json:"restrictions"`
	Triggers       *string `json:"triggers"`
	DiagnosedDate  *string `json:"diagnosedDate"`
	Notes          *string `json:"notes"`
}

type UpdateConditionRequest struct {
	ConditionName  *string `json:"conditionName"`
	ConditionType  *string `json:"conditionType"`
	Severity       *string `json:"severity"`
	ManagementPlan *string `json:"managementPlan"`
	Restrictions   *string `json:"restrictions"`
	Triggers       *string `json:"triggers"`
	DiagnosedDate  *string `json:"diagnosedDate"`
	IsActive       *bool   `json:"isActive"`
	Notes          *string `json:"notes"`
}

// =============================================================================
// Medication DTOs
// =============================================================================

type MedicationResponse struct {
	ID        string `json:"id"`
	StudentID string `json:"studentId"`

	MedicationName           string  `json:"medicationName"`
	Dosage                   string  `json:"dosage"`
	Frequency                string  `json:"frequency"`
	Route                    string  `json:"route"`
	Purpose                  *string `json:"purpose,omitempty"`
	SpecialInstructions      *string `json:"specialInstructions,omitempty"`
	StartDate                string  `json:"startDate"`
	EndDate                  *string `json:"endDate,omitempty"`
	AdministeredAtSchool     bool    `json:"administeredAtSchool"`
	SchoolAdministrationTime *string `json:"schoolAdministrationTime,omitempty"`
	PrescribingDoctor        *string `json:"prescribingDoctor,omitempty"`
	PrescriptionDate         *string `json:"prescriptionDate,omitempty"`
	IsActive                 bool    `json:"isActive"`
	Notes                    *string `json:"notes,omitempty"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type MedicationListResponse struct {
	Medications []MedicationResponse `json:"medications"`
	Total       int                  `json:"total"`
}

type CreateMedicationRequest struct {
	MedicationName           string  `json:"medicationName" binding:"required"`
	Dosage                   string  `json:"dosage" binding:"required"`
	Frequency                string  `json:"frequency" binding:"required"`
	Route                    string  `json:"route" binding:"required"`
	Purpose                  *string `json:"purpose"`
	SpecialInstructions      *string `json:"specialInstructions"`
	StartDate                string  `json:"startDate" binding:"required"`
	EndDate                  *string `json:"endDate"`
	AdministeredAtSchool     bool    `json:"administeredAtSchool"`
	SchoolAdministrationTime *string `json:"schoolAdministrationTime"`
	PrescribingDoctor        *string `json:"prescribingDoctor"`
	PrescriptionDate         *string `json:"prescriptionDate"`
	Notes                    *string `json:"notes"`
}

type UpdateMedicationRequest struct {
	MedicationName           *string `json:"medicationName"`
	Dosage                   *string `json:"dosage"`
	Frequency                *string `json:"frequency"`
	Route                    *string `json:"route"`
	Purpose                  *string `json:"purpose"`
	SpecialInstructions      *string `json:"specialInstructions"`
	StartDate                *string `json:"startDate"`
	EndDate                  *string `json:"endDate"`
	AdministeredAtSchool     *bool   `json:"administeredAtSchool"`
	SchoolAdministrationTime *string `json:"schoolAdministrationTime"`
	PrescribingDoctor        *string `json:"prescribingDoctor"`
	PrescriptionDate         *string `json:"prescriptionDate"`
	IsActive                 *bool   `json:"isActive"`
	Notes                    *string `json:"notes"`
}

// =============================================================================
// Vaccination DTOs
// =============================================================================

type VaccinationResponse struct {
	ID        string `json:"id"`
	StudentID string `json:"studentId"`

	VaccineName        string  `json:"vaccineName"`
	VaccineType        *string `json:"vaccineType,omitempty"`
	DoseNumber         int     `json:"doseNumber"`
	AdministeredDate   string  `json:"administeredDate"`
	AdministeredBy     *string `json:"administeredBy,omitempty"`
	AdministrationSite *string `json:"administrationSite,omitempty"`
	BatchNumber        *string `json:"batchNumber,omitempty"`
	NextDueDate        *string `json:"nextDueDate,omitempty"`
	HadReaction        bool    `json:"hadReaction"`
	ReactionDescription *string `json:"reactionDescription,omitempty"`
	CertificateUrl     *string `json:"certificateUrl,omitempty"`
	IsVerified         bool    `json:"isVerified"`
	VerifiedAt         *string `json:"verifiedAt,omitempty"`
	Notes              *string `json:"notes,omitempty"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type VaccinationListResponse struct {
	Vaccinations []VaccinationResponse `json:"vaccinations"`
	Total        int                   `json:"total"`
}

type CreateVaccinationRequest struct {
	VaccineName         string  `json:"vaccineName" binding:"required"`
	VaccineType         *string `json:"vaccineType"`
	DoseNumber          int     `json:"doseNumber"`
	AdministeredDate    string  `json:"administeredDate" binding:"required"`
	AdministeredBy      *string `json:"administeredBy"`
	AdministrationSite  *string `json:"administrationSite"`
	BatchNumber         *string `json:"batchNumber"`
	NextDueDate         *string `json:"nextDueDate"`
	HadReaction         bool    `json:"hadReaction"`
	ReactionDescription *string `json:"reactionDescription"`
	CertificateUrl      *string `json:"certificateUrl"`
	Notes               *string `json:"notes"`
}

type UpdateVaccinationRequest struct {
	VaccineName         *string `json:"vaccineName"`
	VaccineType         *string `json:"vaccineType"`
	DoseNumber          *int    `json:"doseNumber"`
	AdministeredDate    *string `json:"administeredDate"`
	AdministeredBy      *string `json:"administeredBy"`
	AdministrationSite  *string `json:"administrationSite"`
	BatchNumber         *string `json:"batchNumber"`
	NextDueDate         *string `json:"nextDueDate"`
	HadReaction         *bool   `json:"hadReaction"`
	ReactionDescription *string `json:"reactionDescription"`
	CertificateUrl      *string `json:"certificateUrl"`
	Notes               *string `json:"notes"`
}

// =============================================================================
// Medical Incident DTOs
// =============================================================================

type IncidentResponse struct {
	ID        string `json:"id"`
	StudentID string `json:"studentId"`

	IncidentDate        string  `json:"incidentDate"`
	IncidentTime        string  `json:"incidentTime"`
	Location            *string `json:"location,omitempty"`
	IncidentType        string  `json:"incidentType"`
	Description         string  `json:"description"`
	Symptoms            *string `json:"symptoms,omitempty"`
	FirstAidGiven       bool    `json:"firstAidGiven"`
	FirstAidDescription *string `json:"firstAidDescription,omitempty"`
	ActionTaken         string  `json:"actionTaken"`
	ParentNotified      bool    `json:"parentNotified"`
	ParentNotifiedAt    *string `json:"parentNotifiedAt,omitempty"`
	HospitalVisitRequired bool  `json:"hospitalVisitRequired"`
	HospitalName        *string `json:"hospitalName,omitempty"`
	HospitalVisitDate   *string `json:"hospitalVisitDate,omitempty"`
	StudentSentHome     bool    `json:"studentSentHome"`
	ReturnToClassTime   *string `json:"returnToClassTime,omitempty"`
	FollowUpRequired    bool    `json:"followUpRequired"`
	FollowUpDate        *string `json:"followUpDate,omitempty"`
	FollowUpNotes       *string `json:"followUpNotes,omitempty"`
	Outcome             *string `json:"outcome,omitempty"`
	ReportedBy          string  `json:"reportedBy"`
	Notes               *string `json:"notes,omitempty"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type IncidentListResponse struct {
	Incidents []IncidentResponse `json:"incidents"`
	Total     int                `json:"total"`
}

type CreateIncidentRequest struct {
	IncidentDate          string  `json:"incidentDate" binding:"required"`
	IncidentTime          string  `json:"incidentTime" binding:"required"`
	Location              *string `json:"location"`
	IncidentType          string  `json:"incidentType" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	Symptoms              *string `json:"symptoms"`
	FirstAidGiven         bool    `json:"firstAidGiven"`
	FirstAidDescription   *string `json:"firstAidDescription"`
	ActionTaken           string  `json:"actionTaken" binding:"required"`
	ParentNotified        bool    `json:"parentNotified"`
	ParentNotifiedAt      *string `json:"parentNotifiedAt"`
	HospitalVisitRequired bool    `json:"hospitalVisitRequired"`
	HospitalName          *string `json:"hospitalName"`
	HospitalVisitDate     *string `json:"hospitalVisitDate"`
	StudentSentHome       bool    `json:"studentSentHome"`
	ReturnToClassTime     *string `json:"returnToClassTime"`
	FollowUpRequired      bool    `json:"followUpRequired"`
	FollowUpDate          *string `json:"followUpDate"`
	FollowUpNotes         *string `json:"followUpNotes"`
	Outcome               *string `json:"outcome"`
	Notes                 *string `json:"notes"`
}

type UpdateIncidentRequest struct {
	IncidentDate          *string `json:"incidentDate"`
	IncidentTime          *string `json:"incidentTime"`
	Location              *string `json:"location"`
	IncidentType          *string `json:"incidentType"`
	Description           *string `json:"description"`
	Symptoms              *string `json:"symptoms"`
	FirstAidGiven         *bool   `json:"firstAidGiven"`
	FirstAidDescription   *string `json:"firstAidDescription"`
	ActionTaken           *string `json:"actionTaken"`
	ParentNotified        *bool   `json:"parentNotified"`
	ParentNotifiedAt      *string `json:"parentNotifiedAt"`
	HospitalVisitRequired *bool   `json:"hospitalVisitRequired"`
	HospitalName          *string `json:"hospitalName"`
	HospitalVisitDate     *string `json:"hospitalVisitDate"`
	StudentSentHome       *bool   `json:"studentSentHome"`
	ReturnToClassTime     *string `json:"returnToClassTime"`
	FollowUpRequired      *bool   `json:"followUpRequired"`
	FollowUpDate          *string `json:"followUpDate"`
	FollowUpNotes         *string `json:"followUpNotes"`
	Outcome               *string `json:"outcome"`
	Notes                 *string `json:"notes"`
}

// =============================================================================
// Health Summary DTO
// =============================================================================

type HealthSummaryResponse struct {
	Profile    *HealthProfileResponse `json:"profile,omitempty"`
	Allergies  []AllergyResponse      `json:"allergies"`
	Conditions []ConditionResponse    `json:"conditions"`
	Medications []MedicationResponse  `json:"medications"`
	Vaccinations []VaccinationResponse `json:"vaccinations"`
	RecentIncidents []IncidentResponse `json:"recentIncidents"`
}

// =============================================================================
// Helper Functions
// =============================================================================

func formatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func formatDateTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
