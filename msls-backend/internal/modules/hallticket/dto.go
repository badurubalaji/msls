// Package hallticket provides hall ticket generation and management.
package hallticket

import (
	"time"

	"github.com/google/uuid"
)

// HallTicketStatus represents the status of a hall ticket.
type HallTicketStatus string

const (
	// StatusGenerated indicates the hall ticket has been generated.
	StatusGenerated HallTicketStatus = "generated"
	// StatusPrinted indicates the hall ticket has been printed.
	StatusPrinted HallTicketStatus = "printed"
	// StatusDownloaded indicates the hall ticket has been downloaded.
	StatusDownloaded HallTicketStatus = "downloaded"
)

// HallTicket represents a hall ticket for a student taking an examination.
type HallTicket struct {
	ID            uuid.UUID        `json:"id"`
	TenantID      uuid.UUID        `json:"tenantId"`
	ExaminationID uuid.UUID        `json:"examinationId"`
	StudentID     uuid.UUID        `json:"studentId"`
	RollNumber    string           `json:"rollNumber"`
	QRCodeData    string           `json:"qrCodeData"`
	Status        HallTicketStatus `json:"status"`
	GeneratedAt   time.Time        `json:"generatedAt"`
	PrintedAt     *time.Time       `json:"printedAt,omitempty"`
	DownloadedAt  *time.Time       `json:"downloadedAt,omitempty"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`

	// Joined data
	StudentName       string `json:"studentName,omitempty"`
	StudentPhoto      string `json:"studentPhoto,omitempty"`
	AdmissionNumber   string `json:"admissionNumber,omitempty"`
	ClassName         string `json:"className,omitempty"`
	SectionName       string `json:"sectionName,omitempty"`
	ExaminationName   string `json:"examinationName,omitempty"`
}

// HallTicketTemplate represents a template for generating hall tickets.
type HallTicketTemplate struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenantId"`
	Name          string     `json:"name"`
	HeaderLogoURL string     `json:"headerLogoUrl,omitempty"`
	SchoolName    string     `json:"schoolName,omitempty"`
	SchoolAddress string     `json:"schoolAddress,omitempty"`
	Instructions  string     `json:"instructions,omitempty"`
	IsDefault     bool       `json:"isDefault"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	CreatedBy     *uuid.UUID `json:"createdBy,omitempty"`
	UpdatedBy     *uuid.UUID `json:"updatedBy,omitempty"`
}

// CreateTemplateRequest is the request for creating a hall ticket template.
type CreateTemplateRequest struct {
	TenantID      uuid.UUID  `json:"-"`
	Name          string     `json:"name" binding:"required"`
	HeaderLogoURL string     `json:"headerLogoUrl"`
	SchoolName    string     `json:"schoolName"`
	SchoolAddress string     `json:"schoolAddress"`
	Instructions  string     `json:"instructions"`
	IsDefault     bool       `json:"isDefault"`
	CreatedBy     *uuid.UUID `json:"-"`
}

// UpdateTemplateRequest is the request for updating a hall ticket template.
type UpdateTemplateRequest struct {
	Name          *string `json:"name"`
	HeaderLogoURL *string `json:"headerLogoUrl"`
	SchoolName    *string `json:"schoolName"`
	SchoolAddress *string `json:"schoolAddress"`
	Instructions  *string `json:"instructions"`
	IsDefault     *bool   `json:"isDefault"`
	UpdatedBy     *uuid.UUID `json:"-"`
}

// GenerateRequest is the request for generating hall tickets.
type GenerateRequest struct {
	TenantID         uuid.UUID  `json:"-"`
	ExaminationID    uuid.UUID  `json:"-"` // From URL path
	ClassID          *uuid.UUID `json:"classId"`          // Optional - if nil, generate for all classes
	SectionID        *uuid.UUID `json:"sectionId"`        // Optional - filter by section
	RollNumberPrefix string     `json:"rollNumberPrefix"` // Optional - custom prefix
	CreatedBy        *uuid.UUID `json:"-"`
}

// GenerateResponse is the response after generating hall tickets.
type GenerateResponse struct {
	TotalStudents int      `json:"totalStudents"`
	Generated     int      `json:"generated"`
	Skipped       int      `json:"skipped"` // Already had hall tickets
	Failed        int      `json:"failed"`
	Errors        []string `json:"errors,omitempty"`
}

// ListFilter is the filter for listing hall tickets.
type ListFilter struct {
	TenantID      uuid.UUID
	ExaminationID uuid.UUID
	ClassID       *uuid.UUID
	SectionID     *uuid.UUID
	Status        *HallTicketStatus
	Search        string // Search by student name or roll number
	Limit         int
	Offset        int
}

// VerifyResponse is the response for verifying a hall ticket QR code.
type VerifyResponse struct {
	Valid           bool      `json:"valid"`
	HallTicketID    uuid.UUID `json:"hallTicketId,omitempty"`
	StudentName     string    `json:"studentName,omitempty"`
	AdmissionNumber string    `json:"admissionNumber,omitempty"`
	RollNumber      string    `json:"rollNumber,omitempty"`
	ExaminationName string    `json:"examinationName,omitempty"`
	ClassName       string    `json:"className,omitempty"`
	Message         string    `json:"message,omitempty"`
}

// ExamScheduleItem represents a single exam schedule entry for hall ticket display.
type ExamScheduleItem struct {
	SubjectName  string    `json:"subjectName"`
	SubjectCode  string    `json:"subjectCode"`
	ExamDate     time.Time `json:"examDate"`
	StartTime    string    `json:"startTime"`
	EndTime      string    `json:"endTime"`
	MaxMarks     int       `json:"maxMarks"`
	Venue        string    `json:"venue,omitempty"`
}

// HallTicketPDFData contains all data needed to generate a hall ticket PDF.
type HallTicketPDFData struct {
	HallTicket      *HallTicket
	Template        *HallTicketTemplate
	ExamSchedules   []ExamScheduleItem
	ExamName        string
	ExamStartDate   time.Time
	ExamEndDate     time.Time
	StudentPhotoURL string
}
