// Package hallticket provides hall ticket generation and management.
package hallticket

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles database operations for hall tickets.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new hall ticket repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// hallTicketModel is the database model for hall tickets.
type hallTicketModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	TenantID      uuid.UUID  `gorm:"type:uuid;not null"`
	ExaminationID uuid.UUID  `gorm:"type:uuid;not null"`
	StudentID     uuid.UUID  `gorm:"type:uuid;not null"`
	RollNumber    string     `gorm:"type:varchar(50);not null"`
	QRCodeData    string     `gorm:"type:varchar(500);not null"`
	Status        string     `gorm:"type:varchar(20);not null;default:'generated'"`
	GeneratedAt   time.Time  `gorm:"not null;default:now()"`
	PrintedAt     *time.Time
	DownloadedAt  *time.Time
	CreatedAt     time.Time  `gorm:"not null;default:now()"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()"`
}

func (hallTicketModel) TableName() string {
	return "hall_tickets"
}

// templateModel is the database model for hall ticket templates.
type templateModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	TenantID      uuid.UUID  `gorm:"type:uuid;not null"`
	Name          string     `gorm:"type:varchar(100);not null"`
	HeaderLogoURL string     `gorm:"type:varchar(500)"`
	SchoolName    string     `gorm:"type:varchar(200)"`
	SchoolAddress string     `gorm:"type:text"`
	Instructions  string     `gorm:"type:text"`
	IsDefault     bool       `gorm:"not null;default:false"`
	CreatedAt     time.Time  `gorm:"not null;default:now()"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()"`
	CreatedBy     *uuid.UUID `gorm:"type:uuid"`
	UpdatedBy     *uuid.UUID `gorm:"type:uuid"`
}

func (templateModel) TableName() string {
	return "hall_ticket_templates"
}

// CreateHallTicket creates a new hall ticket.
func (r *Repository) CreateHallTicket(ctx context.Context, ht *HallTicket) error {
	model := &hallTicketModel{
		ID:            ht.ID,
		TenantID:      ht.TenantID,
		ExaminationID: ht.ExaminationID,
		StudentID:     ht.StudentID,
		RollNumber:    ht.RollNumber,
		QRCodeData:    ht.QRCodeData,
		Status:        string(ht.Status),
		GeneratedAt:   ht.GeneratedAt,
		PrintedAt:     ht.PrintedAt,
		DownloadedAt:  ht.DownloadedAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("create hall ticket: %w", err)
	}

	ht.CreatedAt = model.CreatedAt
	ht.UpdatedAt = model.UpdatedAt
	return nil
}

// CreateHallTicketBatch creates multiple hall tickets in a batch.
func (r *Repository) CreateHallTicketBatch(ctx context.Context, tickets []*HallTicket) error {
	if len(tickets) == 0 {
		return nil
	}

	models := make([]hallTicketModel, len(tickets))
	now := time.Now()
	for i, ht := range tickets {
		models[i] = hallTicketModel{
			ID:            ht.ID,
			TenantID:      ht.TenantID,
			ExaminationID: ht.ExaminationID,
			StudentID:     ht.StudentID,
			RollNumber:    ht.RollNumber,
			QRCodeData:    ht.QRCodeData,
			Status:        string(ht.Status),
			GeneratedAt:   ht.GeneratedAt,
			PrintedAt:     ht.PrintedAt,
			DownloadedAt:  ht.DownloadedAt,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
	}

	if err := r.db.WithContext(ctx).CreateInBatches(models, 100).Error; err != nil {
		return fmt.Errorf("create hall tickets batch: %w", err)
	}

	return nil
}

// GetHallTicketByID retrieves a hall ticket by ID.
func (r *Repository) GetHallTicketByID(ctx context.Context, tenantID, id uuid.UUID) (*HallTicket, error) {
	var model hallTicketModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrHallTicketNotFound
		}
		return nil, fmt.Errorf("get hall ticket by id: %w", err)
	}

	return r.modelToHallTicket(&model), nil
}

// GetHallTicketByExamAndStudent retrieves a hall ticket by examination and student.
func (r *Repository) GetHallTicketByExamAndStudent(ctx context.Context, tenantID, examID, studentID uuid.UUID) (*HallTicket, error) {
	var model hallTicketModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND examination_id = ? AND student_id = ?", tenantID, examID, studentID).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrHallTicketNotFound
		}
		return nil, fmt.Errorf("get hall ticket by exam and student: %w", err)
	}

	return r.modelToHallTicket(&model), nil
}

// ListHallTickets lists hall tickets with filters.
func (r *Repository) ListHallTickets(ctx context.Context, filter ListFilter) ([]*HallTicket, int64, error) {
	query := r.db.WithContext(ctx).
		Table("hall_tickets ht").
		Select(`ht.*,
			COALESCE(s.first_name || ' ' || s.last_name, '') as student_name,
			s.photo_url as student_photo,
			s.admission_number,
			COALESCE(c.name, '') as class_name,
			COALESCE(sec.name, '') as section_name,
			COALESCE(e.name, '') as examination_name`).
		Joins("LEFT JOIN students s ON ht.student_id = s.id").
		Joins("LEFT JOIN student_enrollments se ON se.student_id = s.id AND se.status = 'active'").
		Joins("LEFT JOIN classes c ON se.class_id = c.id").
		Joins("LEFT JOIN sections sec ON se.section_id = sec.id").
		Joins("LEFT JOIN examinations e ON ht.examination_id = e.id").
		Where("ht.tenant_id = ?", filter.TenantID).
		Where("ht.examination_id = ?", filter.ExaminationID)

	if filter.ClassID != nil {
		query = query.Where("se.class_id = ?", *filter.ClassID)
	}

	if filter.SectionID != nil {
		query = query.Where("se.section_id = ?", *filter.SectionID)
	}

	if filter.Status != nil {
		query = query.Where("ht.status = ?", string(*filter.Status))
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("(LOWER(s.first_name || ' ' || s.last_name) LIKE LOWER(?) OR ht.roll_number LIKE ?)", searchPattern, searchPattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count hall tickets: %w", err)
	}

	query = query.Order("ht.roll_number ASC")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	type hallTicketWithJoins struct {
		hallTicketModel
		StudentName     string `gorm:"column:student_name"`
		StudentPhoto    string `gorm:"column:student_photo"`
		AdmissionNumber string `gorm:"column:admission_number"`
		ClassName       string `gorm:"column:class_name"`
		SectionName     string `gorm:"column:section_name"`
		ExaminationName string `gorm:"column:examination_name"`
	}

	var models []hallTicketWithJoins
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("list hall tickets: %w", err)
	}

	tickets := make([]*HallTicket, len(models))
	for i, m := range models {
		tickets[i] = &HallTicket{
			ID:              m.ID,
			TenantID:        m.TenantID,
			ExaminationID:   m.ExaminationID,
			StudentID:       m.StudentID,
			RollNumber:      m.RollNumber,
			QRCodeData:      m.QRCodeData,
			Status:          HallTicketStatus(m.Status),
			GeneratedAt:     m.GeneratedAt,
			PrintedAt:       m.PrintedAt,
			DownloadedAt:    m.DownloadedAt,
			CreatedAt:       m.CreatedAt,
			UpdatedAt:       m.UpdatedAt,
			StudentName:     m.StudentName,
			StudentPhoto:    m.StudentPhoto,
			AdmissionNumber: m.AdmissionNumber,
			ClassName:       m.ClassName,
			SectionName:     m.SectionName,
			ExaminationName: m.ExaminationName,
		}
	}

	return tickets, total, nil
}

// UpdateHallTicketStatus updates the status of a hall ticket.
func (r *Repository) UpdateHallTicketStatus(ctx context.Context, tenantID, id uuid.UUID, status HallTicketStatus) error {
	updates := map[string]interface{}{
		"status":     string(status),
		"updated_at": time.Now(),
	}

	now := time.Now()
	switch status {
	case StatusPrinted:
		updates["printed_at"] = now
	case StatusDownloaded:
		updates["downloaded_at"] = now
	}

	result := r.db.WithContext(ctx).
		Model(&hallTicketModel{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("update hall ticket status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrHallTicketNotFound
	}

	return nil
}

// GetExistingStudentIDs returns student IDs that already have hall tickets for an exam.
func (r *Repository) GetExistingStudentIDs(ctx context.Context, tenantID, examID uuid.UUID) (map[uuid.UUID]bool, error) {
	var studentIDs []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&hallTicketModel{}).
		Where("tenant_id = ? AND examination_id = ?", tenantID, examID).
		Pluck("student_id", &studentIDs).Error

	if err != nil {
		return nil, fmt.Errorf("get existing student ids: %w", err)
	}

	result := make(map[uuid.UUID]bool, len(studentIDs))
	for _, id := range studentIDs {
		result[id] = true
	}
	return result, nil
}

// GetLastRollNumber gets the last roll number sequence for an exam.
func (r *Repository) GetLastRollNumber(ctx context.Context, tenantID, examID uuid.UUID, prefix string) (int, error) {
	var maxRollNumber string
	err := r.db.WithContext(ctx).
		Model(&hallTicketModel{}).
		Where("tenant_id = ? AND examination_id = ? AND roll_number LIKE ?", tenantID, examID, prefix+"%").
		Order("roll_number DESC").
		Limit(1).
		Pluck("roll_number", &maxRollNumber).Error

	if err != nil || maxRollNumber == "" {
		return 0, nil
	}

	// Extract sequence number from roll number (e.g., "2026-10A-015" -> 15)
	var seq int
	_, err = fmt.Sscanf(maxRollNumber, prefix+"%d", &seq)
	if err != nil {
		return 0, nil
	}
	return seq, nil
}

// DeleteHallTicket deletes a hall ticket.
func (r *Repository) DeleteHallTicket(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&hallTicketModel{})

	if result.Error != nil {
		return fmt.Errorf("delete hall ticket: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrHallTicketNotFound
	}

	return nil
}

// DeleteHallTicketsByExam deletes all hall tickets for an examination.
func (r *Repository) DeleteHallTicketsByExam(ctx context.Context, tenantID, examID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND examination_id = ?", tenantID, examID).
		Delete(&hallTicketModel{})

	if result.Error != nil {
		return fmt.Errorf("delete hall tickets by exam: %w", result.Error)
	}

	return nil
}

// Template Repository Methods

// CreateTemplate creates a new hall ticket template.
func (r *Repository) CreateTemplate(ctx context.Context, tmpl *HallTicketTemplate) error {
	// If this is set as default, unset other defaults
	if tmpl.IsDefault {
		if err := r.unsetDefaultTemplate(ctx, tmpl.TenantID); err != nil {
			return err
		}
	}

	model := &templateModel{
		ID:            tmpl.ID,
		TenantID:      tmpl.TenantID,
		Name:          tmpl.Name,
		HeaderLogoURL: tmpl.HeaderLogoURL,
		SchoolName:    tmpl.SchoolName,
		SchoolAddress: tmpl.SchoolAddress,
		Instructions:  tmpl.Instructions,
		IsDefault:     tmpl.IsDefault,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CreatedBy:     tmpl.CreatedBy,
		UpdatedBy:     tmpl.UpdatedBy,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("create template: %w", err)
	}

	tmpl.CreatedAt = model.CreatedAt
	tmpl.UpdatedAt = model.UpdatedAt
	return nil
}

// GetTemplateByID retrieves a template by ID.
func (r *Repository) GetTemplateByID(ctx context.Context, tenantID, id uuid.UUID) (*HallTicketTemplate, error) {
	var model templateModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("get template by id: %w", err)
	}

	return r.modelToTemplate(&model), nil
}

// GetDefaultTemplate retrieves the default template for a tenant.
func (r *Repository) GetDefaultTemplate(ctx context.Context, tenantID uuid.UUID) (*HallTicketTemplate, error) {
	var model templateModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_default = true", tenantID).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return first template if no default set
			err = r.db.WithContext(ctx).
				Where("tenant_id = ?", tenantID).
				Order("created_at ASC").
				First(&model).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, ErrTemplateNotFound
				}
				return nil, fmt.Errorf("get first template: %w", err)
			}
		} else {
			return nil, fmt.Errorf("get default template: %w", err)
		}
	}

	return r.modelToTemplate(&model), nil
}

// ListTemplates lists all templates for a tenant.
func (r *Repository) ListTemplates(ctx context.Context, tenantID uuid.UUID) ([]*HallTicketTemplate, error) {
	var models []templateModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("is_default DESC, name ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}

	templates := make([]*HallTicketTemplate, len(models))
	for i, m := range models {
		templates[i] = r.modelToTemplate(&m)
	}
	return templates, nil
}

// UpdateTemplate updates a template.
func (r *Repository) UpdateTemplate(ctx context.Context, tenantID, id uuid.UUID, req *UpdateTemplateRequest) (*HallTicketTemplate, error) {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.HeaderLogoURL != nil {
		updates["header_logo_url"] = *req.HeaderLogoURL
	}
	if req.SchoolName != nil {
		updates["school_name"] = *req.SchoolName
	}
	if req.SchoolAddress != nil {
		updates["school_address"] = *req.SchoolAddress
	}
	if req.Instructions != nil {
		updates["instructions"] = *req.Instructions
	}
	if req.UpdatedBy != nil {
		updates["updated_by"] = *req.UpdatedBy
	}

	if req.IsDefault != nil && *req.IsDefault {
		if err := r.unsetDefaultTemplate(ctx, tenantID); err != nil {
			return nil, err
		}
		updates["is_default"] = true
	} else if req.IsDefault != nil {
		updates["is_default"] = false
	}

	result := r.db.WithContext(ctx).
		Model(&templateModel{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("update template: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, ErrTemplateNotFound
	}

	return r.GetTemplateByID(ctx, tenantID, id)
}

// DeleteTemplate deletes a template.
func (r *Repository) DeleteTemplate(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&templateModel{})

	if result.Error != nil {
		return fmt.Errorf("delete template: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrTemplateNotFound
	}

	return nil
}

func (r *Repository) unsetDefaultTemplate(ctx context.Context, tenantID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&templateModel{}).
		Where("tenant_id = ? AND is_default = true", tenantID).
		Update("is_default", false).Error
}

func (r *Repository) modelToHallTicket(m *hallTicketModel) *HallTicket {
	return &HallTicket{
		ID:            m.ID,
		TenantID:      m.TenantID,
		ExaminationID: m.ExaminationID,
		StudentID:     m.StudentID,
		RollNumber:    m.RollNumber,
		QRCodeData:    m.QRCodeData,
		Status:        HallTicketStatus(m.Status),
		GeneratedAt:   m.GeneratedAt,
		PrintedAt:     m.PrintedAt,
		DownloadedAt:  m.DownloadedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (r *Repository) modelToTemplate(m *templateModel) *HallTicketTemplate {
	return &HallTicketTemplate{
		ID:            m.ID,
		TenantID:      m.TenantID,
		Name:          m.Name,
		HeaderLogoURL: m.HeaderLogoURL,
		SchoolName:    m.SchoolName,
		SchoolAddress: m.SchoolAddress,
		Instructions:  m.Instructions,
		IsDefault:     m.IsDefault,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		CreatedBy:     m.CreatedBy,
		UpdatedBy:     m.UpdatedBy,
	}
}
