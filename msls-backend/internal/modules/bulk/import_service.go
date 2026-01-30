// Package bulk provides bulk operation functionality.
package bulk

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// ImportService handles bulk import functionality.
type ImportService struct {
	db *gorm.DB
}

// NewImportService creates a new import service.
func NewImportService(db *gorm.DB) *ImportService {
	return &ImportService{db: db}
}

// StudentImportRow represents a row in the student import template.
type StudentImportRow struct {
	RowNum          int
	AdmissionNumber string
	FirstName       string
	MiddleName      string
	LastName        string
	DateOfBirth     string
	Gender          string
	BloodGroup      string
	AadhaarNumber   string
	ClassName       string
	SectionName     string
	RollNumber      string
	AdmissionDate   string
	// Guardian info
	GuardianName     string
	GuardianRelation string
	GuardianPhone    string
	GuardianEmail    string
	// Address
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
}

// ImportResult contains the result of a bulk import operation.
type ImportResult struct {
	TotalRows    int            `json:"totalRows"`
	SuccessCount int            `json:"successCount"`
	FailedCount  int            `json:"failedCount"`
	Errors       []ImportError  `json:"errors,omitempty"`
	CreatedIDs   []string       `json:"createdIds,omitempty"`
}

// ImportError represents an error during import.
type ImportError struct {
	Row     int    `json:"row"`
	Column  string `json:"column,omitempty"`
	Message string `json:"message"`
}

// Template column headers
var studentImportHeaders = []string{
	"Admission Number*",
	"First Name*",
	"Middle Name",
	"Last Name*",
	"Date of Birth* (YYYY-MM-DD)",
	"Gender* (male/female/other)",
	"Blood Group",
	"Aadhaar Number",
	"Class Name*",
	"Section Name*",
	"Roll Number",
	"Admission Date (YYYY-MM-DD)",
	"Guardian Name*",
	"Guardian Relation* (father/mother/guardian)",
	"Guardian Phone*",
	"Guardian Email",
	"Address Line 1",
	"Address Line 2",
	"City",
	"State",
	"Postal Code",
}

// GenerateTemplate generates an Excel template for student import.
func (s *ImportService) GenerateTemplate(ctx context.Context, tenantID uuid.UUID) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Students"
	f.SetSheetName("Sheet1", sheetName)

	// Style for headers
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// Write headers
	for i, header := range studentImportHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
		// Set column width
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, colName, colName, 20)
	}

	// Add sample row
	sampleData := []string{
		"ADM2024001",
		"John",
		"",
		"Doe",
		"2015-06-15",
		"male",
		"O+",
		"123456789012",
		"Class 5",
		"Section A",
		"1",
		"2024-04-01",
		"Robert Doe",
		"father",
		"9876543210",
		"robert@email.com",
		"123 Main Street",
		"Apt 4B",
		"Mumbai",
		"Maharashtra",
		"400001",
	}
	for i, val := range sampleData {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, val)
	}

	// Add instructions sheet
	f.NewSheet("Instructions")
	instructions := []string{
		"Student Bulk Import Instructions",
		"",
		"1. Fill in the 'Students' sheet with student data",
		"2. Fields marked with * are mandatory",
		"3. Date format: YYYY-MM-DD (e.g., 2015-06-15)",
		"4. Gender values: male, female, other",
		"5. Guardian relation: father, mother, guardian",
		"6. Class Name and Section Name must exist in the system",
		"7. Remove the sample row before uploading",
		"8. Maximum 500 students per import",
		"",
		"Notes:",
		"- Admission Number must be unique",
		"- If admission date is empty, current date will be used",
		"- Students will be enrolled in the current academic year",
	}
	for i, line := range instructions {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue("Instructions", cell, line)
	}

	// Get classes and sections for reference
	var classes []struct {
		ID   uuid.UUID
		Name string
	}
	s.db.WithContext(ctx).
		Model(&models.Class{}).
		Select("id, name").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("display_order, name").
		Find(&classes)

	// Add reference sheet with valid classes and sections
	f.NewSheet("Reference")
	f.SetCellValue("Reference", "A1", "Available Classes")
	f.SetCellValue("Reference", "B1", "Available Sections")

	row := 2
	for _, class := range classes {
		f.SetCellValue("Reference", fmt.Sprintf("A%d", row), class.Name)

		var sections []struct {
			Name string
		}
		s.db.WithContext(ctx).
			Model(&models.Section{}).
			Select("name").
			Where("class_id = ?", class.ID).
			Where("deleted_at IS NULL").
			Order("name").
			Find(&sections)

		sectionNames := make([]string, len(sections))
		for i, sec := range sections {
			sectionNames[i] = sec.Name
		}
		f.SetCellValue("Reference", fmt.Sprintf("B%d", row), strings.Join(sectionNames, ", "))
		row++
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("write excel: %w", err)
	}

	return buf.Bytes(), nil
}

// ImportStudents imports students from an Excel or CSV file.
func (s *ImportService) ImportStudents(ctx context.Context, tenantID, branchID, academicYearID uuid.UUID, userID uuid.UUID, fileContent []byte, fileType string) (*ImportResult, error) {
	var rows []StudentImportRow
	var err error

	if fileType == "csv" {
		rows, err = s.parseCSV(fileContent)
	} else {
		rows, err = s.parseExcel(fileContent)
	}

	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no data rows found in file")
	}

	if len(rows) > 500 {
		return nil, fmt.Errorf("too many rows (max 500)")
	}

	result := &ImportResult{
		TotalRows:  len(rows),
		CreatedIDs: make([]string, 0),
		Errors:     make([]ImportError, 0),
	}

	// Build lookup maps for classes and sections
	classMap, err := s.buildClassMap(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("build class map: %w", err)
	}

	sectionMap, err := s.buildSectionMap(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("build section map: %w", err)
	}

	// Process each row
	for _, row := range rows {
		// Validate row
		validationErrors := s.validateRow(row, classMap, sectionMap)
		if len(validationErrors) > 0 {
			result.Errors = append(result.Errors, validationErrors...)
			result.FailedCount++
			continue
		}

		// Create student with enrollment in transaction
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Parse dates
			dob, _ := time.Parse("2006-01-02", row.DateOfBirth)
			admissionDate := time.Now()
			if row.AdmissionDate != "" {
				admissionDate, _ = time.Parse("2006-01-02", row.AdmissionDate)
			}

			// Create student
			student := models.Student{
				TenantID:        tenantID,
				BranchID:        branchID,
				AdmissionNumber: row.AdmissionNumber,
				FirstName:       row.FirstName,
				MiddleName:      row.MiddleName,
				LastName:        row.LastName,
				DateOfBirth:     dob,
				Gender:          models.Gender(strings.ToLower(row.Gender)),
				BloodGroup:      row.BloodGroup,
				AadhaarNumber:   row.AadhaarNumber,
				Status:          models.StudentStatusActive,
				AdmissionDate:   admissionDate,
				CreatedBy:       &userID,
			}

			if err := tx.Create(&student).Error; err != nil {
				return fmt.Errorf("create student: %w", err)
			}

			// Create enrollment in student_enrollments table
			classID := classMap[strings.ToLower(row.ClassName)]
			sectionKey := fmt.Sprintf("%s|%s", strings.ToLower(row.ClassName), strings.ToLower(row.SectionName))
			sectionID := sectionMap[sectionKey]

			enrollmentData := map[string]interface{}{
				"tenant_id":       tenantID,
				"student_id":      student.ID,
				"academic_year_id": academicYearID,
				"class_id":        classID,
				"section_id":      sectionID,
				"roll_number":     row.RollNumber,
				"status":          "active",
				"enrollment_date": admissionDate,
				"created_by":      userID,
				"created_at":      time.Now(),
				"updated_at":      time.Now(),
			}

			if err := tx.Table("student_enrollments").Create(enrollmentData).Error; err != nil {
				return fmt.Errorf("create enrollment: %w", err)
			}

			// Create guardian if provided
			if row.GuardianName != "" && row.GuardianPhone != "" {
				relation := models.GuardianRelation(strings.ToLower(row.GuardianRelation))
				if !relation.IsValid() {
					relation = models.GuardianRelationGuardian
				}

				// Parse guardian name into first/last
				nameParts := strings.SplitN(row.GuardianName, " ", 2)
				firstName := nameParts[0]
				lastName := ""
				if len(nameParts) > 1 {
					lastName = nameParts[1]
				}

				guardian := models.StudentGuardian{
					TenantID:  tenantID,
					StudentID: student.ID,
					FirstName: firstName,
					LastName:  lastName,
					Relation:  relation,
					Phone:     row.GuardianPhone,
					Email:     row.GuardianEmail,
					IsPrimary: true,
					CreatedBy: &userID,
				}

				if err := tx.Create(&guardian).Error; err != nil {
					return fmt.Errorf("create guardian: %w", err)
				}
			}

			// Create address if provided
			if row.AddressLine1 != "" {
				address := models.StudentAddress{
					TenantID:     tenantID,
					StudentID:    student.ID,
					AddressType:  models.AddressTypeCurrent,
					AddressLine1: row.AddressLine1,
					AddressLine2: row.AddressLine2,
					City:         row.City,
					State:        row.State,
					PostalCode:   row.PostalCode,
					Country:      "India",
				}

				if err := tx.Create(&address).Error; err != nil {
					return fmt.Errorf("create address: %w", err)
				}
			}

			result.CreatedIDs = append(result.CreatedIDs, student.ID.String())
			return nil
		})

		if err != nil {
			result.Errors = append(result.Errors, ImportError{
				Row:     row.RowNum,
				Message: err.Error(),
			})
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// parseExcel parses an Excel file into StudentImportRow slice.
func (s *ImportService) parseExcel(content []byte) ([]StudentImportRow, error) {
	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, nil
	}

	result := make([]StudentImportRow, 0, len(rows)-1)
	for i, row := range rows[1:] { // Skip header
		if len(row) == 0 || (len(row) > 0 && strings.TrimSpace(row[0]) == "") {
			continue
		}

		importRow := StudentImportRow{RowNum: i + 2}
		for j, cell := range row {
			s.setCellValue(&importRow, j, strings.TrimSpace(cell))
		}
		result = append(result, importRow)
	}

	return result, nil
}

// parseCSV parses a CSV file into StudentImportRow slice.
func (s *ImportService) parseCSV(content []byte) ([]StudentImportRow, error) {
	reader := csv.NewReader(bytes.NewReader(content))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, nil
	}

	result := make([]StudentImportRow, 0, len(rows)-1)
	for i, row := range rows[1:] { // Skip header
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}

		importRow := StudentImportRow{RowNum: i + 2}
		for j, cell := range row {
			s.setCellValue(&importRow, j, strings.TrimSpace(cell))
		}
		result = append(result, importRow)
	}

	return result, nil
}

// setCellValue sets a cell value on the import row based on column index.
func (s *ImportService) setCellValue(row *StudentImportRow, colIndex int, value string) {
	switch colIndex {
	case 0:
		row.AdmissionNumber = value
	case 1:
		row.FirstName = value
	case 2:
		row.MiddleName = value
	case 3:
		row.LastName = value
	case 4:
		row.DateOfBirth = value
	case 5:
		row.Gender = value
	case 6:
		row.BloodGroup = value
	case 7:
		row.AadhaarNumber = value
	case 8:
		row.ClassName = value
	case 9:
		row.SectionName = value
	case 10:
		row.RollNumber = value
	case 11:
		row.AdmissionDate = value
	case 12:
		row.GuardianName = value
	case 13:
		row.GuardianRelation = value
	case 14:
		row.GuardianPhone = value
	case 15:
		row.GuardianEmail = value
	case 16:
		row.AddressLine1 = value
	case 17:
		row.AddressLine2 = value
	case 18:
		row.City = value
	case 19:
		row.State = value
	case 20:
		row.PostalCode = value
	}
}

// validateRow validates a single import row.
func (s *ImportService) validateRow(row StudentImportRow, classMap map[string]uuid.UUID, sectionMap map[string]uuid.UUID) []ImportError {
	var errors []ImportError

	// Required fields
	if row.AdmissionNumber == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Admission Number", Message: "Admission number is required"})
	}
	if row.FirstName == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "First Name", Message: "First name is required"})
	}
	if row.LastName == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Last Name", Message: "Last name is required"})
	}
	if row.DateOfBirth == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Date of Birth", Message: "Date of birth is required"})
	} else if _, err := time.Parse("2006-01-02", row.DateOfBirth); err != nil {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Date of Birth", Message: "Invalid date format (use YYYY-MM-DD)"})
	}
	if row.Gender == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Gender", Message: "Gender is required"})
	} else {
		gender := models.Gender(strings.ToLower(row.Gender))
		if !gender.IsValid() {
			errors = append(errors, ImportError{Row: row.RowNum, Column: "Gender", Message: "Invalid gender (use male/female/other)"})
		}
	}
	if row.ClassName == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Class Name", Message: "Class name is required"})
	} else if _, ok := classMap[strings.ToLower(row.ClassName)]; !ok {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Class Name", Message: "Class not found in system"})
	}
	if row.SectionName == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Section Name", Message: "Section name is required"})
	} else {
		sectionKey := fmt.Sprintf("%s|%s", strings.ToLower(row.ClassName), strings.ToLower(row.SectionName))
		if _, ok := sectionMap[sectionKey]; !ok {
			errors = append(errors, ImportError{Row: row.RowNum, Column: "Section Name", Message: "Section not found for the specified class"})
		}
	}
	if row.GuardianName == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Guardian Name", Message: "Guardian name is required"})
	}
	if row.GuardianPhone == "" {
		errors = append(errors, ImportError{Row: row.RowNum, Column: "Guardian Phone", Message: "Guardian phone is required"})
	}

	// Optional validation
	if row.AdmissionDate != "" {
		if _, err := time.Parse("2006-01-02", row.AdmissionDate); err != nil {
			errors = append(errors, ImportError{Row: row.RowNum, Column: "Admission Date", Message: "Invalid date format (use YYYY-MM-DD)"})
		}
	}

	return errors
}

// buildClassMap builds a map of class names to IDs.
func (s *ImportService) buildClassMap(ctx context.Context, tenantID uuid.UUID) (map[string]uuid.UUID, error) {
	var classes []struct {
		ID   uuid.UUID
		Name string
	}

	err := s.db.WithContext(ctx).
		Model(&models.Class{}).
		Select("id, name").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Find(&classes).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]uuid.UUID)
	for _, c := range classes {
		result[strings.ToLower(c.Name)] = c.ID
	}

	return result, nil
}

// buildSectionMap builds a map of "className|sectionName" to section IDs.
func (s *ImportService) buildSectionMap(ctx context.Context, tenantID uuid.UUID) (map[string]uuid.UUID, error) {
	var sections []struct {
		ID        uuid.UUID
		Name      string
		ClassName string
	}

	err := s.db.WithContext(ctx).
		Table("sections").
		Select("sections.id, sections.name, classes.name as class_name").
		Joins("JOIN classes ON classes.id = sections.class_id").
		Where("sections.tenant_id = ?", tenantID).
		Where("sections.deleted_at IS NULL").
		Find(&sections).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]uuid.UUID)
	for _, sec := range sections {
		key := fmt.Sprintf("%s|%s", strings.ToLower(sec.ClassName), strings.ToLower(sec.Name))
		result[key] = sec.ID
	}

	return result, nil
}
