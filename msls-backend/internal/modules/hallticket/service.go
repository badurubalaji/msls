// Package hallticket provides hall ticket generation and management.
package hallticket

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service provides business logic for hall ticket operations.
type Service struct {
	repo         *Repository
	db           *gorm.DB
	qrGenerator  *QRCodeGenerator
	pdfGenerator *PDFGenerator
}

// NewService creates a new hall ticket service.
func NewService(db *gorm.DB, qrSecret string) *Service {
	repo := NewRepository(db)
	qrGen := NewQRCodeGenerator(qrSecret)
	pdfGen := NewPDFGenerator(qrGen)

	return &Service{
		repo:         repo,
		db:           db,
		qrGenerator:  qrGen,
		pdfGenerator: pdfGen,
	}
}

// GenerateHallTickets generates hall tickets for students in an examination.
func (s *Service) GenerateHallTickets(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.ExaminationID == uuid.Nil {
		return nil, ErrExaminationIDRequired
	}

	// Verify examination exists and is scheduled
	var examStatus string
	var examClasses []uuid.UUID
	err := s.db.WithContext(ctx).
		Table("examinations").
		Select("status").
		Where("tenant_id = ? AND id = ?", req.TenantID, req.ExaminationID).
		Scan(&examStatus).Error
	if err != nil {
		return nil, fmt.Errorf("get examination: %w", err)
	}
	if examStatus != "scheduled" {
		return nil, ErrExaminationNotScheduled
	}

	// Get classes for the examination
	err = s.db.WithContext(ctx).
		Table("examination_classes").
		Select("class_id").
		Where("examination_id = ?", req.ExaminationID).
		Pluck("class_id", &examClasses).Error
	if err != nil {
		return nil, fmt.Errorf("get examination classes: %w", err)
	}
	if len(examClasses) == 0 {
		return nil, ErrNoClassesForExam
	}

	// Filter by class if specified
	if req.ClassID != nil {
		found := false
		for _, c := range examClasses {
			if c == *req.ClassID {
				found = true
				break
			}
		}
		if !found {
			return nil, ErrNoClassesForExam
		}
		examClasses = []uuid.UUID{*req.ClassID}
	}

	// Get students enrolled in these classes
	query := s.db.WithContext(ctx).
		Table("student_enrollments se").
		Select("se.student_id, s.first_name, s.last_name, s.admission_number, c.name as class_name, c.code as class_code").
		Joins("JOIN students s ON se.student_id = s.id").
		Joins("JOIN classes c ON se.class_id = c.id").
		Where("se.tenant_id = ? AND se.status = 'active'", req.TenantID).
		Where("se.class_id IN ?", examClasses)

	if req.SectionID != nil {
		query = query.Where("se.section_id = ?", *req.SectionID)
	}

	type studentInfo struct {
		StudentID       uuid.UUID `gorm:"column:student_id"`
		FirstName       string    `gorm:"column:first_name"`
		LastName        string    `gorm:"column:last_name"`
		AdmissionNumber string    `gorm:"column:admission_number"`
		ClassName       string    `gorm:"column:class_name"`
		ClassCode       string    `gorm:"column:class_code"`
	}

	var students []studentInfo
	if err := query.Scan(&students).Error; err != nil {
		return nil, fmt.Errorf("get students: %w", err)
	}

	if len(students) == 0 {
		return nil, ErrNoStudentsInClass
	}

	// Get existing hall tickets
	existing, err := s.repo.GetExistingStudentIDs(ctx, req.TenantID, req.ExaminationID)
	if err != nil {
		return nil, err
	}

	// Determine roll number prefix
	year := time.Now().Year()
	rollNumberPrefix := fmt.Sprintf("%d-", year)
	if req.RollNumberPrefix != "" {
		rollNumberPrefix = req.RollNumberPrefix + "-"
	}

	// Get current sequence per class
	classSequences := make(map[string]int)

	response := &GenerateResponse{
		TotalStudents: len(students),
	}

	var ticketsToCreate []*HallTicket

	for _, student := range students {
		// Skip if already has hall ticket
		if existing[student.StudentID] {
			response.Skipped++
			continue
		}

		// Get or initialize sequence for this class
		classPrefix := rollNumberPrefix + student.ClassCode + "-"
		if _, ok := classSequences[classPrefix]; !ok {
			seq, _ := s.repo.GetLastRollNumber(ctx, req.TenantID, req.ExaminationID, classPrefix)
			classSequences[classPrefix] = seq
		}
		classSequences[classPrefix]++
		seq := classSequences[classPrefix]

		// Generate roll number
		rollNumber := fmt.Sprintf("%s%03d", classPrefix, seq)

		// Generate ticket ID and QR code
		ticketID := uuid.New()
		qrData, err := s.qrGenerator.GenerateQRCodeData(ticketID, student.StudentID, req.ExaminationID)
		if err != nil {
			response.Failed++
			response.Errors = append(response.Errors, fmt.Sprintf("QR code for %s: %v", student.AdmissionNumber, err))
			continue
		}

		ticket := &HallTicket{
			ID:            ticketID,
			TenantID:      req.TenantID,
			ExaminationID: req.ExaminationID,
			StudentID:     student.StudentID,
			RollNumber:    rollNumber,
			QRCodeData:    qrData,
			Status:        StatusGenerated,
			GeneratedAt:   time.Now(),
		}

		ticketsToCreate = append(ticketsToCreate, ticket)
	}

	// Batch create tickets
	if len(ticketsToCreate) > 0 {
		if err := s.repo.CreateHallTicketBatch(ctx, ticketsToCreate); err != nil {
			return nil, fmt.Errorf("create hall tickets: %w", err)
		}
		response.Generated = len(ticketsToCreate)
	}

	return response, nil
}

// ListHallTickets lists hall tickets with filters.
func (s *Service) ListHallTickets(ctx context.Context, filter ListFilter) ([]*HallTicket, int64, error) {
	return s.repo.ListHallTickets(ctx, filter)
}

// GetHallTicket retrieves a single hall ticket.
func (s *Service) GetHallTicket(ctx context.Context, tenantID, id uuid.UUID) (*HallTicket, error) {
	return s.repo.GetHallTicketByID(ctx, tenantID, id)
}

// GetHallTicketPDF generates a PDF for a single hall ticket.
func (s *Service) GetHallTicketPDF(ctx context.Context, tenantID, examID, ticketID uuid.UUID) ([]byte, string, error) {
	ticket, err := s.repo.GetHallTicketByID(ctx, tenantID, ticketID)
	if err != nil {
		return nil, "", err
	}

	// Get full ticket data with joins
	tickets, _, err := s.repo.ListHallTickets(ctx, ListFilter{
		TenantID:      tenantID,
		ExaminationID: examID,
		Limit:         1,
	})
	if err != nil {
		return nil, "", err
	}

	// Find our ticket in the results
	var fullTicket *HallTicket
	for _, t := range tickets {
		if t.ID == ticketID {
			fullTicket = t
			break
		}
	}
	if fullTicket == nil {
		fullTicket = ticket
	}

	pdfData, err := s.buildPDFData(ctx, tenantID, examID, fullTicket)
	if err != nil {
		return nil, "", err
	}

	pdf, err := s.pdfGenerator.GenerateHallTicketPDF(pdfData)
	if err != nil {
		return nil, "", err
	}

	// Update status to downloaded
	_ = s.repo.UpdateHallTicketStatus(ctx, tenantID, ticketID, StatusDownloaded)

	return pdf, GetFilename(ticket.RollNumber), nil
}

// GetBatchPDF generates a PDF containing all hall tickets for an examination.
func (s *Service) GetBatchPDF(ctx context.Context, tenantID, examID uuid.UUID, classID *uuid.UUID) ([]byte, string, error) {
	filter := ListFilter{
		TenantID:      tenantID,
		ExaminationID: examID,
		ClassID:       classID,
		Limit:         1000, // Reasonable max for PDF
	}

	tickets, _, err := s.repo.ListHallTickets(ctx, filter)
	if err != nil {
		return nil, "", err
	}

	if len(tickets) == 0 {
		return nil, "", ErrHallTicketNotFound
	}

	var pdfDataList []*HallTicketPDFData
	for _, ticket := range tickets {
		pdfData, err := s.buildPDFData(ctx, tenantID, examID, ticket)
		if err != nil {
			continue // Skip problematic tickets
		}
		pdfDataList = append(pdfDataList, pdfData)
	}

	if len(pdfDataList) == 0 {
		return nil, "", ErrHallTicketNotFound
	}

	pdf, err := s.pdfGenerator.GenerateBatchPDF(pdfDataList)
	if err != nil {
		return nil, "", err
	}

	// Get exam name for filename
	var examName string
	s.db.WithContext(ctx).
		Table("examinations").
		Select("name").
		Where("id = ?", examID).
		Scan(&examName)

	return pdf, GetBatchFilename(examName), nil
}

func (s *Service) buildPDFData(ctx context.Context, tenantID, examID uuid.UUID, ticket *HallTicket) (*HallTicketPDFData, error) {
	// Get template
	template, err := s.repo.GetDefaultTemplate(ctx, tenantID)
	if err != nil && err != ErrTemplateNotFound {
		return nil, err
	}

	// Get exam details
	var exam struct {
		Name      string    `gorm:"column:name"`
		StartDate time.Time `gorm:"column:start_date"`
		EndDate   time.Time `gorm:"column:end_date"`
	}
	s.db.WithContext(ctx).
		Table("examinations").
		Select("name, start_date, end_date").
		Where("id = ?", examID).
		Scan(&exam)

	// Get exam schedules
	var schedules []ExamScheduleItem
	s.db.WithContext(ctx).
		Table("exam_schedules es").
		Select("sub.name as subject_name, sub.code as subject_code, es.exam_date, es.start_time::text, es.end_time::text, es.max_marks, COALESCE(es.venue, '') as venue").
		Joins("JOIN subjects sub ON es.subject_id = sub.id").
		Where("es.examination_id = ?", examID).
		Order("es.exam_date ASC, es.start_time ASC").
		Scan(&schedules)

	// Get student photo
	var photoURL string
	s.db.WithContext(ctx).
		Table("students").
		Select("photo_url").
		Where("id = ?", ticket.StudentID).
		Scan(&photoURL)

	return &HallTicketPDFData{
		HallTicket:      ticket,
		Template:        template,
		ExamSchedules:   schedules,
		ExamName:        exam.Name,
		ExamStartDate:   exam.StartDate,
		ExamEndDate:     exam.EndDate,
		StudentPhotoURL: photoURL,
	}, nil
}

// VerifyHallTicket verifies a hall ticket QR code.
func (s *Service) VerifyHallTicket(ctx context.Context, qrData string) (*VerifyResponse, error) {
	payload, err := ParseQRCodeData(qrData)
	if err != nil {
		return &VerifyResponse{
			Valid:   false,
			Message: "Invalid QR code format",
		}, nil
	}

	// Find hall ticket by short ID prefix
	var ticket struct {
		ID              uuid.UUID
		TenantID        uuid.UUID
		ExaminationID   uuid.UUID
		StudentID       uuid.UUID
		RollNumber      string
		QRCodeData      string
		StudentName     string
		AdmissionNumber string
		ExamName        string
		ClassName       string
	}

	err = s.db.WithContext(ctx).
		Table("hall_tickets ht").
		Select(`ht.id, ht.tenant_id, ht.examination_id, ht.student_id, ht.roll_number, ht.qr_code_data,
			COALESCE(s.first_name || ' ' || s.last_name, '') as student_name,
			s.admission_number,
			e.name as exam_name,
			COALESCE(c.name, '') as class_name`).
		Joins("LEFT JOIN students s ON ht.student_id = s.id").
		Joins("LEFT JOIN examinations e ON ht.examination_id = e.id").
		Joins("LEFT JOIN student_enrollments se ON se.student_id = s.id AND se.status = 'active'").
		Joins("LEFT JOIN classes c ON se.class_id = c.id").
		Where("ht.id::text LIKE ?", payload.TicketID+"%").
		First(&ticket).Error

	if err != nil {
		return &VerifyResponse{
			Valid:   false,
			Message: "Hall ticket not found",
		}, nil
	}

	// Verify the QR code data matches
	if !s.qrGenerator.VerifyQRCodeData(ticket.QRCodeData, ticket.ID, ticket.StudentID, ticket.ExaminationID) {
		return &VerifyResponse{
			Valid:   false,
			Message: "QR code verification failed",
		}, nil
	}

	return &VerifyResponse{
		Valid:           true,
		HallTicketID:    ticket.ID,
		StudentName:     ticket.StudentName,
		AdmissionNumber: ticket.AdmissionNumber,
		RollNumber:      ticket.RollNumber,
		ExaminationName: ticket.ExamName,
		ClassName:       ticket.ClassName,
		Message:         "Valid hall ticket",
	}, nil
}

// DeleteHallTicket deletes a hall ticket.
func (s *Service) DeleteHallTicket(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteHallTicket(ctx, tenantID, id)
}

// DeleteHallTicketsByExam deletes all hall tickets for an examination.
func (s *Service) DeleteHallTicketsByExam(ctx context.Context, tenantID, examID uuid.UUID) error {
	return s.repo.DeleteHallTicketsByExam(ctx, tenantID, examID)
}

// Template methods

// CreateTemplate creates a new hall ticket template.
func (s *Service) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*HallTicketTemplate, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.Name == "" {
		return nil, ErrTemplateNameRequired
	}

	tmpl := &HallTicketTemplate{
		ID:            uuid.New(),
		TenantID:      req.TenantID,
		Name:          req.Name,
		HeaderLogoURL: req.HeaderLogoURL,
		SchoolName:    req.SchoolName,
		SchoolAddress: req.SchoolAddress,
		Instructions:  req.Instructions,
		IsDefault:     req.IsDefault,
		CreatedBy:     req.CreatedBy,
		UpdatedBy:     req.CreatedBy,
	}

	if err := s.repo.CreateTemplate(ctx, tmpl); err != nil {
		return nil, err
	}

	return tmpl, nil
}

// GetTemplate retrieves a template by ID.
func (s *Service) GetTemplate(ctx context.Context, tenantID, id uuid.UUID) (*HallTicketTemplate, error) {
	return s.repo.GetTemplateByID(ctx, tenantID, id)
}

// ListTemplates lists all templates for a tenant.
func (s *Service) ListTemplates(ctx context.Context, tenantID uuid.UUID) ([]*HallTicketTemplate, error) {
	return s.repo.ListTemplates(ctx, tenantID)
}

// UpdateTemplate updates a template.
func (s *Service) UpdateTemplate(ctx context.Context, tenantID, id uuid.UUID, req *UpdateTemplateRequest) (*HallTicketTemplate, error) {
	return s.repo.UpdateTemplate(ctx, tenantID, id, req)
}

// DeleteTemplate deletes a template.
func (s *Service) DeleteTemplate(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, tenantID, id)
}
