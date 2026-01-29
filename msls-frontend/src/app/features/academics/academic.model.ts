/**
 * Academic module models for Class, Section, and Stream management.
 */

// Class level enum matching backend
export type ClassLevel = 'nursery' | 'primary' | 'middle' | 'secondary' | 'senior_secondary';

export const CLASS_LEVELS: { value: ClassLevel; label: string }[] = [
  { value: 'nursery', label: 'Nursery (LKG/UKG)' },
  { value: 'primary', label: 'Primary (Class 1-5)' },
  { value: 'middle', label: 'Middle (Class 6-8)' },
  { value: 'secondary', label: 'Secondary (Class 9-10)' },
  { value: 'senior_secondary', label: 'Senior Secondary (Class 11-12)' },
];

// ========================================
// Stream Models
// ========================================

export interface Stream {
  id: string;
  name: string;
  code: string;
  description?: string;
  displayOrder: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateStreamRequest {
  name: string;
  code: string;
  description?: string;
  displayOrder?: number;
}

export interface UpdateStreamRequest {
  name?: string;
  code?: string;
  description?: string;
  displayOrder?: number;
  isActive?: boolean;
}

export interface StreamListResponse {
  streams: Stream[];
  total: number;
}

// ========================================
// Class Models
// ========================================

export interface Class {
  id: string;
  branchId: string;
  branchName?: string;
  name: string;
  code: string;
  level?: ClassLevel;
  displayOrder: number;
  description?: string;
  hasStreams: boolean;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
  sections?: Section[];
  streams?: Stream[];
}

export interface CreateClassRequest {
  branchId: string;
  name: string;
  code: string;
  level?: ClassLevel;
  displayOrder?: number;
  description?: string;
  hasStreams?: boolean;
  streamIds?: string[];
}

export interface UpdateClassRequest {
  name?: string;
  code?: string;
  level?: ClassLevel;
  displayOrder?: number;
  description?: string;
  hasStreams?: boolean;
  isActive?: boolean;
  streamIds?: string[];
}

export interface ClassListResponse {
  classes: Class[];
  total: number;
}

export interface ClassFilter {
  branchId?: string;
  level?: ClassLevel;
  isActive?: boolean;
  hasStreams?: boolean;
  search?: string;
}

// ========================================
// Section Models
// ========================================

export interface Section {
  id: string;
  classId: string;
  className?: string;
  academicYearId?: string;
  academicYearName?: string;
  streamId?: string;
  streamName?: string;
  classTeacherId?: string;
  classTeacherName?: string;
  name: string;
  code: string;
  capacity: number;
  roomNumber?: string;
  displayOrder: number;
  isActive: boolean;
  studentCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateSectionRequest {
  classId: string;
  academicYearId?: string;
  streamId?: string;
  classTeacherId?: string;
  name: string;
  code: string;
  capacity?: number;
  roomNumber?: string;
}

export interface UpdateSectionRequest {
  academicYearId?: string;
  streamId?: string;
  classTeacherId?: string;
  name?: string;
  code?: string;
  capacity?: number;
  roomNumber?: string;
  displayOrder?: number;
  isActive?: boolean;
}

export interface SectionListResponse {
  sections: Section[];
  total: number;
}

export interface SectionFilter {
  classId?: string;
  academicYearId?: string;
  streamId?: string;
  isActive?: boolean;
  search?: string;
}

// ========================================
// Class Structure Models (Hierarchical View)
// ========================================

export interface SectionStructure {
  id: string;
  name: string;
  code: string;
  capacity: number;
  studentCount: number;
  classTeacherId?: string;
  classTeacherName?: string;
  streamName?: string;
  roomNumber?: string;
  capacityUsage: number; // Percentage
}

export interface ClassWithSections {
  id: string;
  name: string;
  code: string;
  displayOrder: number;
  hasStreams: boolean;
  isActive: boolean;
  sections: SectionStructure[];
  totalStudents: number;
  totalCapacity: number;
}

export interface ClassStructureResponse {
  classes: ClassWithSections[];
}

// ========================================
// Subject Models
// ========================================

export type SubjectType = 'core' | 'elective' | 'language' | 'co_curricular' | 'vocational';

export const SUBJECT_TYPES: { value: SubjectType; label: string }[] = [
  { value: 'core', label: 'Core' },
  { value: 'elective', label: 'Elective' },
  { value: 'language', label: 'Language' },
  { value: 'co_curricular', label: 'Co-Curricular' },
  { value: 'vocational', label: 'Vocational' },
];

export interface Subject {
  id: string;
  name: string;
  code: string;
  shortName?: string;
  description?: string;
  subjectType: SubjectType;
  maxMarks: number;
  passingMarks: number;
  creditHours: number;
  displayOrder: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateSubjectRequest {
  name: string;
  code: string;
  shortName?: string;
  description?: string;
  subjectType: SubjectType;
  maxMarks?: number;
  passingMarks?: number;
  creditHours?: number;
  displayOrder?: number;
}

export interface UpdateSubjectRequest {
  name?: string;
  code?: string;
  shortName?: string;
  description?: string;
  subjectType?: SubjectType;
  maxMarks?: number;
  passingMarks?: number;
  creditHours?: number;
  displayOrder?: number;
  isActive?: boolean;
}

export interface SubjectListResponse {
  subjects: Subject[];
  total: number;
}

export interface SubjectFilter {
  subjectType?: SubjectType;
  isActive?: boolean;
  search?: string;
}

// ========================================
// Class-Subject Models
// ========================================

export interface ClassSubject {
  id: string;
  classId: string;
  className?: string;
  subjectId: string;
  subjectName?: string;
  subjectCode?: string;
  isMandatory: boolean;
  periodsPerWeek: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateClassSubjectRequest {
  subjectId: string;
  isMandatory?: boolean;
  periodsPerWeek?: number;
}

export interface UpdateClassSubjectRequest {
  isMandatory?: boolean;
  periodsPerWeek?: number;
  isActive?: boolean;
}

export interface ClassSubjectListResponse {
  classSubjects: ClassSubject[];
  total: number;
}
