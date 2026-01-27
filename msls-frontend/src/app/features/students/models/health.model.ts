/**
 * Student Health Models
 */

// =============================================================================
// Enums and Types
// =============================================================================

export type AllergyType = 'food' | 'medication' | 'environmental' | 'insect' | 'other';
export type AllergySeverity = 'mild' | 'moderate' | 'severe' | 'life_threatening';
export type ConditionType = 'respiratory' | 'cardiac' | 'neurological' | 'endocrine' | 'other';
export type ConditionSeverity = 'mild' | 'moderate' | 'severe';
export type MedicationFrequency = 'daily' | 'twice_daily' | 'three_daily' | 'weekly' | 'as_needed' | 'other';
export type MedicationRoute = 'oral' | 'injection' | 'inhaler' | 'topical' | 'drops' | 'other';
export type VaccineType = 'required' | 'optional' | 'booster';
export type IncidentType = 'illness' | 'injury' | 'emergency' | 'other';

// =============================================================================
// Option Constants
// =============================================================================

export const ALLERGY_TYPE_OPTIONS: { value: AllergyType; label: string }[] = [
  { value: 'food', label: 'Food' },
  { value: 'medication', label: 'Medication' },
  { value: 'environmental', label: 'Environmental' },
  { value: 'insect', label: 'Insect' },
  { value: 'other', label: 'Other' },
];

export const ALLERGY_SEVERITY_OPTIONS: { value: AllergySeverity; label: string }[] = [
  { value: 'mild', label: 'Mild' },
  { value: 'moderate', label: 'Moderate' },
  { value: 'severe', label: 'Severe' },
  { value: 'life_threatening', label: 'Life Threatening' },
];

export const CONDITION_TYPE_OPTIONS: { value: ConditionType; label: string }[] = [
  { value: 'respiratory', label: 'Respiratory' },
  { value: 'cardiac', label: 'Cardiac' },
  { value: 'neurological', label: 'Neurological' },
  { value: 'endocrine', label: 'Endocrine' },
  { value: 'other', label: 'Other' },
];

export const CONDITION_SEVERITY_OPTIONS: { value: ConditionSeverity; label: string }[] = [
  { value: 'mild', label: 'Mild' },
  { value: 'moderate', label: 'Moderate' },
  { value: 'severe', label: 'Severe' },
];

export const MEDICATION_FREQUENCY_OPTIONS: { value: MedicationFrequency; label: string }[] = [
  { value: 'daily', label: 'Once Daily' },
  { value: 'twice_daily', label: 'Twice Daily' },
  { value: 'three_daily', label: 'Three Times Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'as_needed', label: 'As Needed' },
  { value: 'other', label: 'Other' },
];

export const MEDICATION_ROUTE_OPTIONS: { value: MedicationRoute; label: string }[] = [
  { value: 'oral', label: 'Oral' },
  { value: 'injection', label: 'Injection' },
  { value: 'inhaler', label: 'Inhaler' },
  { value: 'topical', label: 'Topical' },
  { value: 'drops', label: 'Drops' },
  { value: 'other', label: 'Other' },
];

export const VACCINE_TYPE_OPTIONS: { value: VaccineType; label: string }[] = [
  { value: 'required', label: 'Required' },
  { value: 'optional', label: 'Optional' },
  { value: 'booster', label: 'Booster' },
];

export const INCIDENT_TYPE_OPTIONS: { value: IncidentType; label: string }[] = [
  { value: 'illness', label: 'Illness' },
  { value: 'injury', label: 'Injury' },
  { value: 'emergency', label: 'Emergency' },
  { value: 'other', label: 'Other' },
];

export const BLOOD_GROUP_OPTIONS = ['A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-'];

// =============================================================================
// Interfaces
// =============================================================================

export interface HealthProfile {
  id: string;
  studentId: string;
  bloodGroup?: string;
  heightCm?: number;
  weightKg?: number;
  visionLeft?: string;
  visionRight?: string;
  hearingStatus?: string;
  medicalNotes?: string;
  insuranceProvider?: string;
  insurancePolicyNumber?: string;
  insuranceExpiry?: string;
  preferredHospital?: string;
  familyDoctorName?: string;
  familyDoctorPhone?: string;
  lastCheckupDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Allergy {
  id: string;
  studentId: string;
  allergen: string;
  allergyType: AllergyType;
  severity: AllergySeverity;
  reactionDescription?: string;
  treatmentInstructions?: string;
  emergencyMedication?: string;
  diagnosedDate?: string;
  isActive: boolean;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ChronicCondition {
  id: string;
  studentId: string;
  conditionName: string;
  conditionType: ConditionType;
  severity: ConditionSeverity;
  managementPlan?: string;
  restrictions?: string;
  triggers?: string;
  diagnosedDate?: string;
  isActive: boolean;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Medication {
  id: string;
  studentId: string;
  medicationName: string;
  dosage: string;
  frequency: MedicationFrequency;
  route: MedicationRoute;
  purpose?: string;
  specialInstructions?: string;
  startDate: string;
  endDate?: string;
  administeredAtSchool: boolean;
  schoolAdministrationTime?: string;
  prescribingDoctor?: string;
  prescriptionDate?: string;
  isActive: boolean;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Vaccination {
  id: string;
  studentId: string;
  vaccineName: string;
  vaccineType?: VaccineType;
  doseNumber: number;
  administeredDate: string;
  administeredBy?: string;
  administrationSite?: string;
  batchNumber?: string;
  nextDueDate?: string;
  hadReaction: boolean;
  reactionDescription?: string;
  certificateUrl?: string;
  isVerified: boolean;
  verifiedAt?: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface MedicalIncident {
  id: string;
  studentId: string;
  incidentDate: string;
  incidentTime: string;
  location?: string;
  incidentType: IncidentType;
  description: string;
  symptoms?: string;
  firstAidGiven: boolean;
  firstAidDescription?: string;
  actionTaken: string;
  parentNotified: boolean;
  parentNotifiedAt?: string;
  hospitalVisitRequired: boolean;
  hospitalName?: string;
  hospitalVisitDate?: string;
  studentSentHome: boolean;
  returnToClassTime?: string;
  followUpRequired: boolean;
  followUpDate?: string;
  followUpNotes?: string;
  outcome?: string;
  reportedBy: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface HealthSummary {
  profile?: HealthProfile;
  allergies: Allergy[];
  conditions: ChronicCondition[];
  medications: Medication[];
  vaccinations: Vaccination[];
  recentIncidents: MedicalIncident[];
}

// =============================================================================
// Request/Response Types
// =============================================================================

export interface AllergyListResponse {
  allergies: Allergy[];
  total: number;
}

export interface ConditionListResponse {
  conditions: ChronicCondition[];
  total: number;
}

export interface MedicationListResponse {
  medications: Medication[];
  total: number;
}

export interface VaccinationListResponse {
  vaccinations: Vaccination[];
  total: number;
}

export interface IncidentListResponse {
  incidents: MedicalIncident[];
  total: number;
}

export interface CreateHealthProfileRequest {
  bloodGroup?: string;
  heightCm?: number;
  weightKg?: number;
  visionLeft?: string;
  visionRight?: string;
  hearingStatus?: string;
  medicalNotes?: string;
  insuranceProvider?: string;
  insurancePolicyNumber?: string;
  insuranceExpiry?: string;
  preferredHospital?: string;
  familyDoctorName?: string;
  familyDoctorPhone?: string;
  lastCheckupDate?: string;
}

export interface CreateAllergyRequest {
  allergen: string;
  allergyType: AllergyType;
  severity: AllergySeverity;
  reactionDescription?: string;
  treatmentInstructions?: string;
  emergencyMedication?: string;
  diagnosedDate?: string;
  notes?: string;
}

export interface CreateConditionRequest {
  conditionName: string;
  conditionType: ConditionType;
  severity: ConditionSeverity;
  managementPlan?: string;
  restrictions?: string;
  triggers?: string;
  diagnosedDate?: string;
  notes?: string;
}

export interface CreateMedicationRequest {
  medicationName: string;
  dosage: string;
  frequency: MedicationFrequency;
  route: MedicationRoute;
  purpose?: string;
  specialInstructions?: string;
  startDate: string;
  endDate?: string;
  administeredAtSchool: boolean;
  schoolAdministrationTime?: string;
  prescribingDoctor?: string;
  prescriptionDate?: string;
  notes?: string;
}

export interface CreateVaccinationRequest {
  vaccineName: string;
  vaccineType?: VaccineType;
  doseNumber: number;
  administeredDate: string;
  administeredBy?: string;
  administrationSite?: string;
  batchNumber?: string;
  nextDueDate?: string;
  hadReaction: boolean;
  reactionDescription?: string;
  certificateUrl?: string;
  notes?: string;
}

export interface CreateIncidentRequest {
  incidentDate: string;
  incidentTime: string;
  location?: string;
  incidentType: IncidentType;
  description: string;
  symptoms?: string;
  firstAidGiven: boolean;
  firstAidDescription?: string;
  actionTaken: string;
  parentNotified: boolean;
  parentNotifiedAt?: string;
  hospitalVisitRequired: boolean;
  hospitalName?: string;
  hospitalVisitDate?: string;
  studentSentHome: boolean;
  returnToClassTime?: string;
  followUpRequired: boolean;
  followUpDate?: string;
  followUpNotes?: string;
  outcome?: string;
  notes?: string;
}

// =============================================================================
// Helper Functions
// =============================================================================

export function getAllergyTypeLabel(type: AllergyType): string {
  return ALLERGY_TYPE_OPTIONS.find(opt => opt.value === type)?.label || type;
}

export function getAllergySeverityLabel(severity: AllergySeverity): string {
  return ALLERGY_SEVERITY_OPTIONS.find(opt => opt.value === severity)?.label || severity;
}

export function getConditionTypeLabel(type: ConditionType): string {
  return CONDITION_TYPE_OPTIONS.find(opt => opt.value === type)?.label || type;
}

export function getMedicationFrequencyLabel(frequency: MedicationFrequency): string {
  return MEDICATION_FREQUENCY_OPTIONS.find(opt => opt.value === frequency)?.label || frequency;
}

export function getMedicationRouteLabel(route: MedicationRoute): string {
  return MEDICATION_ROUTE_OPTIONS.find(opt => opt.value === route)?.label || route;
}

export function getVaccineTypeLabel(type: VaccineType): string {
  return VACCINE_TYPE_OPTIONS.find(opt => opt.value === type)?.label || type;
}

export function getIncidentTypeLabel(type: IncidentType): string {
  return INCIDENT_TYPE_OPTIONS.find(opt => opt.value === type)?.label || type;
}

export function getSeverityColor(severity: AllergySeverity | ConditionSeverity): string {
  switch (severity) {
    case 'mild':
      return 'success';
    case 'moderate':
      return 'warning';
    case 'severe':
    case 'life_threatening':
      return 'danger';
    default:
      return 'neutral';
  }
}
