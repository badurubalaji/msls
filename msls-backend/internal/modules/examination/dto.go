package examination

import (
	"time"

	"github.com/google/uuid"
	"msls-backend/internal/pkg/database/models"
)

// ========================================
// Request DTOs
// ========================================

// CreateExaminationRequest represents the request to create an examination
type CreateExaminationRequest struct {
	Name           string      `json:"name" binding:"required,max=200"`
	ExamTypeID     uuid.UUID   `json:"examTypeId" binding:"required"`
	AcademicYearID uuid.UUID   `json:"academicYearId" binding:"required"`
	StartDate      string      `json:"startDate" binding:"required"` // YYYY-MM-DD
	EndDate        string      `json:"endDate" binding:"required"`   // YYYY-MM-DD
	Description    *string     `json:"description"`
	ClassIDs       []uuid.UUID `json:"classIds" binding:"required,min=1"`
}

// UpdateExaminationRequest represents the request to update an examination
type UpdateExaminationRequest struct {
	Name           *string     `json:"name,omitempty" binding:"omitempty,max=200"`
	ExamTypeID     *uuid.UUID  `json:"examTypeId,omitempty"`
	AcademicYearID *uuid.UUID  `json:"academicYearId,omitempty"`
	StartDate      *string     `json:"startDate,omitempty"` // YYYY-MM-DD
	EndDate        *string     `json:"endDate,omitempty"`   // YYYY-MM-DD
	Description    *string     `json:"description,omitempty"`
	ClassIDs       []uuid.UUID `json:"classIds,omitempty"`
}

// CreateScheduleRequest represents the request to create an exam schedule
type CreateScheduleRequest struct {
	SubjectID    uuid.UUID `json:"subjectId" binding:"required"`
	ExamDate     string    `json:"examDate" binding:"required"` // YYYY-MM-DD
	StartTime    string    `json:"startTime" binding:"required"` // HH:MM
	EndTime      string    `json:"endTime" binding:"required"`   // HH:MM
	MaxMarks     int       `json:"maxMarks" binding:"required,min=1"`
	PassingMarks *int      `json:"passingMarks,omitempty"`
	Venue        *string   `json:"venue,omitempty" binding:"omitempty,max=100"`
	Notes        *string   `json:"notes,omitempty"`
}

// UpdateScheduleRequest represents the request to update an exam schedule
type UpdateScheduleRequest struct {
	SubjectID    *uuid.UUID `json:"subjectId,omitempty"`
	ExamDate     *string    `json:"examDate,omitempty"` // YYYY-MM-DD
	StartTime    *string    `json:"startTime,omitempty"` // HH:MM
	EndTime      *string    `json:"endTime,omitempty"`   // HH:MM
	MaxMarks     *int       `json:"maxMarks,omitempty" binding:"omitempty,min=1"`
	PassingMarks *int       `json:"passingMarks,omitempty"`
	Venue        *string    `json:"venue,omitempty" binding:"omitempty,max=100"`
	Notes        *string    `json:"notes,omitempty"`
}

// ExaminationFilter represents filters for listing examinations
type ExaminationFilter struct {
	AcademicYearID *uuid.UUID `form:"academicYearId"`
	ExamTypeID     *uuid.UUID `form:"examTypeId"`
	ClassID        *uuid.UUID `form:"classId"`
	Status         *string    `form:"status"`
	Search         *string    `form:"search"`
}

// ========================================
// Response DTOs
// ========================================

// ExaminationResponse represents an examination in responses
type ExaminationResponse struct {
	ID             uuid.UUID                `json:"id"`
	Name           string                   `json:"name"`
	ExamTypeID     uuid.UUID                `json:"examTypeId"`
	ExamTypeName   string                   `json:"examTypeName"`
	AcademicYearID uuid.UUID                `json:"academicYearId"`
	AcademicYear   string                   `json:"academicYear"`
	StartDate      string                   `json:"startDate"`
	EndDate        string                   `json:"endDate"`
	Status         models.ExamStatus        `json:"status"`
	Description    *string                  `json:"description,omitempty"`
	Classes        []ClassSummary           `json:"classes"`
	Schedules      []ExamScheduleResponse   `json:"schedules,omitempty"`
	ScheduleCount  int                      `json:"scheduleCount"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
}

// ClassSummary represents a class in examination response
type ClassSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ExamScheduleResponse represents an exam schedule in responses
type ExamScheduleResponse struct {
	ID           uuid.UUID `json:"id"`
	SubjectID    uuid.UUID `json:"subjectId"`
	SubjectName  string    `json:"subjectName"`
	SubjectCode  string    `json:"subjectCode"`
	ExamDate     string    `json:"examDate"`
	StartTime    string    `json:"startTime"`
	EndTime      string    `json:"endTime"`
	MaxMarks     int       `json:"maxMarks"`
	PassingMarks *int      `json:"passingMarks,omitempty"`
	Venue        *string   `json:"venue,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
}

// ========================================
// Mappers
// ========================================

// ToResponse converts an Examination model to ExaminationResponse
func ToResponse(exam *models.Examination) ExaminationResponse {
	resp := ExaminationResponse{
		ID:             exam.ID,
		Name:           exam.Name,
		ExamTypeID:     exam.ExamTypeID,
		AcademicYearID: exam.AcademicYearID,
		StartDate:      exam.StartDate.Format("2006-01-02"),
		EndDate:        exam.EndDate.Format("2006-01-02"),
		Status:         exam.Status,
		Description:    exam.Description,
		Classes:        make([]ClassSummary, 0),
		Schedules:      make([]ExamScheduleResponse, 0),
		ScheduleCount:  len(exam.Schedules),
		CreatedAt:      exam.CreatedAt,
		UpdatedAt:      exam.UpdatedAt,
	}

	// Map exam type name
	if exam.ExamType != nil {
		resp.ExamTypeName = exam.ExamType.Name
	}

	// Map academic year
	if exam.AcademicYear != nil {
		resp.AcademicYear = exam.AcademicYear.Name
	}

	// Map classes
	for _, class := range exam.Classes {
		resp.Classes = append(resp.Classes, ClassSummary{
			ID:   class.ID,
			Name: class.Name,
		})
	}

	// Map schedules
	for _, schedule := range exam.Schedules {
		scheduleResp := ExamScheduleResponse{
			ID:           schedule.ID,
			SubjectID:    schedule.SubjectID,
			ExamDate:     schedule.ExamDate.Format("2006-01-02"),
			StartTime:    schedule.StartTime,
			EndTime:      schedule.EndTime,
			MaxMarks:     schedule.MaxMarks,
			PassingMarks: schedule.PassingMarks,
			Venue:        schedule.Venue,
			Notes:        schedule.Notes,
		}
		if schedule.Subject != nil {
			scheduleResp.SubjectName = schedule.Subject.Name
			scheduleResp.SubjectCode = schedule.Subject.Code
		}
		resp.Schedules = append(resp.Schedules, scheduleResp)
	}

	return resp
}

// ToResponseList converts a slice of Examination models to responses
func ToResponseList(exams []models.Examination) []ExaminationResponse {
	responses := make([]ExaminationResponse, len(exams))
	for i, exam := range exams {
		responses[i] = ToResponse(&exam)
	}
	return responses
}
