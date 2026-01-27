// Package bulk provides bulk operation functionality.
package bulk

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"msls-backend/internal/modules/enrollment"
	"msls-backend/internal/pkg/database/models"
)

// ExportService handles student export functionality.
type ExportService struct {
	db        *gorm.DB
	uploadDir string
}

// NewExportService creates a new export service.
func NewExportService(db *gorm.DB, uploadDir string) *ExportService {
	return &ExportService{
		db:        db,
		uploadDir: uploadDir,
	}
}

// ExportStudents exports students to a file and returns the URL.
func (s *ExportService) ExportStudents(ctx context.Context, tenantID uuid.UUID, studentIDs []uuid.UUID, format string, columns []string) (string, error) {
	// Fetch students with related data
	students, err := s.fetchStudentsForExport(ctx, tenantID, studentIDs)
	if err != nil {
		return "", fmt.Errorf("fetch students: %w", err)
	}

	// Create export directory
	exportDir := filepath.Join(s.uploadDir, "exports", tenantID.String())
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("create export directory: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("students-%s.%s", timestamp, format)
	filePath := filepath.Join(exportDir, filename)

	// Generate file
	switch format {
	case "xlsx":
		if err := s.createExcelFile(filePath, students, columns); err != nil {
			return "", fmt.Errorf("create excel: %w", err)
		}
	case "csv":
		if err := s.createCSVFile(filePath, students, columns); err != nil {
			return "", fmt.Errorf("create csv: %w", err)
		}
	default:
		return "", ErrInvalidExportFormat
	}

	// Return relative URL
	fileURL := fmt.Sprintf("/uploads/exports/%s/%s", tenantID.String(), filename)
	return fileURL, nil
}

// StudentExportData holds flattened student data for export.
type StudentExportData struct {
	AdmissionNumber  string
	FirstName        string
	LastName         string
	FullName         string
	Gender           string
	DateOfBirth      string
	BloodGroup       string
	AadhaarNumber    string
	AdmissionDate    string
	Status           string
	Branch           string
	Class            string
	Section          string
	RollNumber       string
	Phone            string
	Email            string
	GuardianName     string
	GuardianPhone    string
	GuardianEmail    string
	GuardianRelation string
	Address          string
	City             string
	State            string
}

// fetchStudentsForExport fetches students with all related data.
func (s *ExportService) fetchStudentsForExport(ctx context.Context, tenantID uuid.UUID, studentIDs []uuid.UUID) ([]StudentExportData, error) {
	var students []models.Student

	query := s.db.WithContext(ctx).
		Preload("Branch").
		Preload("Addresses").
		Where("tenant_id = ? AND id IN ?", tenantID, studentIDs)

	if err := query.Find(&students).Error; err != nil {
		return nil, err
	}

	// Fetch guardians for all students
	var guardians []models.StudentGuardian
	s.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id IN ? AND is_primary = true", tenantID, studentIDs).
		Find(&guardians)

	// Create guardian lookup
	guardianMap := make(map[uuid.UUID]models.StudentGuardian)
	for _, g := range guardians {
		guardianMap[g.StudentID] = g
	}

	// Fetch active enrollments
	var enrollments []enrollment.StudentEnrollment
	s.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id IN ? AND status = 'active'", tenantID, studentIDs).
		Find(&enrollments)

	// Create enrollment lookup
	enrollmentMap := make(map[uuid.UUID]enrollment.StudentEnrollment)
	for _, e := range enrollments {
		enrollmentMap[e.StudentID] = e
	}

	// Build export data
	result := make([]StudentExportData, len(students))
	for i, student := range students {
		data := StudentExportData{
			AdmissionNumber: student.AdmissionNumber,
			FirstName:       student.FirstName,
			LastName:        student.LastName,
			FullName:        student.FullName(),
			Gender:          string(student.Gender),
			DateOfBirth:     student.DateOfBirth.Format("2006-01-02"),
			BloodGroup:      student.BloodGroup,
			AadhaarNumber:   student.AadhaarNumber,
			AdmissionDate:   student.AdmissionDate.Format("2006-01-02"),
			Status:          string(student.Status),
		}

		// Add branch
		if student.Branch.ID != uuid.Nil {
			data.Branch = student.Branch.Name
		}

		// Add address (current)
		for _, addr := range student.Addresses {
			if addr.AddressType == models.AddressTypeCurrent {
				data.Address = addr.AddressLine1
				if addr.AddressLine2 != "" {
					data.Address += ", " + addr.AddressLine2
				}
				data.City = addr.City
				data.State = addr.State
				break
			}
		}

		// Add guardian info
		if guardian, ok := guardianMap[student.ID]; ok {
			data.GuardianName = guardian.FullName()
			data.GuardianPhone = guardian.Phone
			data.GuardianEmail = guardian.Email
			data.GuardianRelation = string(guardian.Relation)
		}

		// Add enrollment info
		if enrollment, ok := enrollmentMap[student.ID]; ok {
			data.RollNumber = enrollment.RollNumber
			// Note: class and section names would require additional joins
			// For now, we'll leave them blank as they depend on Epic 6
		}

		result[i] = data
	}

	return result, nil
}

// createExcelFile creates an Excel file with student data.
func (s *ExportService) createExcelFile(filePath string, students []StudentExportData, columns []string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Students"
	f.SetSheetName("Sheet1", sheet)

	// Write headers
	for i, col := range columns {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		label := ExportColumnLabels[col]
		if label == "" {
			label = col
		}
		f.SetCellValue(sheet, cell, label)
	}

	// Style header row
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	// Write data rows
	for rowIdx, student := range students {
		for colIdx, col := range columns {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			value := s.getColumnValue(student, col)
			f.SetCellValue(sheet, cell, value)
		}
	}

	// Auto-fit columns
	for i := range columns {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, colName, colName, 15)
	}

	return f.SaveAs(filePath)
}

// createCSVFile creates a CSV file with student data.
func (s *ExportService) createCSVFile(filePath string, students []StudentExportData, columns []string) error {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := make([]string, len(columns))
	for i, col := range columns {
		label := ExportColumnLabels[col]
		if label == "" {
			label = col
		}
		headers[i] = label
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data rows
	for _, student := range students {
		row := make([]string, len(columns))
		for i, col := range columns {
			row[i] = s.getColumnValue(student, col)
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// getColumnValue returns the value for a column from student data.
func (s *ExportService) getColumnValue(student StudentExportData, column string) string {
	switch column {
	case "admission_number":
		return student.AdmissionNumber
	case "first_name":
		return student.FirstName
	case "last_name":
		return student.LastName
	case "full_name":
		return student.FullName
	case "gender":
		return student.Gender
	case "date_of_birth":
		return student.DateOfBirth
	case "blood_group":
		return student.BloodGroup
	case "aadhaar_number":
		return student.AadhaarNumber
	case "admission_date":
		return student.AdmissionDate
	case "status":
		return student.Status
	case "branch":
		return student.Branch
	case "class":
		return student.Class
	case "section":
		return student.Section
	case "roll_number":
		return student.RollNumber
	case "phone":
		return student.Phone
	case "email":
		return student.Email
	case "guardian_name":
		return student.GuardianName
	case "guardian_phone":
		return student.GuardianPhone
	case "guardian_email":
		return student.GuardianEmail
	case "guardian_relation":
		return student.GuardianRelation
	case "address":
		return student.Address
	case "city":
		return student.City
	case "state":
		return student.State
	default:
		return ""
	}
}
