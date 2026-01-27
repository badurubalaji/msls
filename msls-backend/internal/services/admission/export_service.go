// Package admission provides admission management services.
package admission

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ExportFormat represents the export file format.
type ExportFormat string

// Export format constants.
const (
	ExportFormatExcel ExportFormat = "excel"
	ExportFormatPDF   ExportFormat = "pdf"
)

// ExportReportType represents the type of report to export.
type ExportReportType string

// Export report type constants.
const (
	ExportReportTypeDashboard      ExportReportType = "dashboard"
	ExportReportTypeClassWise      ExportReportType = "class-wise"
	ExportReportTypeFunnel         ExportReportType = "funnel"
	ExportReportTypeSourceAnalysis ExportReportType = "source-analysis"
	ExportReportTypeDailyTrend     ExportReportType = "daily-trend"
)

// ExportRequest represents a request to export a report.
type ExportRequest struct {
	ReportType ExportReportType
	Format     ExportFormat
	Filter     DashboardFilter
}

// ExportResult contains the exported file data.
type ExportResult struct {
	Data        []byte
	Filename    string
	ContentType string
}

// ExportService handles report export operations.
type ExportService struct {
	db            *gorm.DB
	reportService *ReportService
}

// NewExportService creates a new ExportService instance.
func NewExportService(db *gorm.DB, reportService *ReportService) *ExportService {
	return &ExportService{
		db:            db,
		reportService: reportService,
	}
}

// ExportReport generates an export of the specified report.
func (s *ExportService) ExportReport(ctx context.Context, req ExportRequest) (*ExportResult, error) {
	if req.Filter.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	switch req.Format {
	case ExportFormatExcel:
		return s.exportToExcel(ctx, req)
	case ExportFormatPDF:
		// PDF export can be implemented with a library like gofpdf
		// For now, we'll return an error indicating not implemented
		return nil, fmt.Errorf("PDF export is not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported export format: %s", req.Format)
	}
}

// exportToExcel generates an Excel file for the specified report.
func (s *ExportService) exportToExcel(ctx context.Context, req ExportRequest) (*ExportResult, error) {
	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	var sheetName string
	var err error

	switch req.ReportType {
	case ExportReportTypeDashboard:
		sheetName = "Dashboard"
		err = s.exportDashboardToExcel(ctx, f, sheetName, req.Filter)
	case ExportReportTypeClassWise:
		sheetName = "Class-wise Report"
		err = s.exportClassWiseToExcel(ctx, f, sheetName, req.Filter)
	case ExportReportTypeFunnel:
		sheetName = "Conversion Funnel"
		err = s.exportFunnelToExcel(ctx, f, sheetName, req.Filter)
	case ExportReportTypeSourceAnalysis:
		sheetName = "Source Analysis"
		err = s.exportSourceAnalysisToExcel(ctx, f, sheetName, req.Filter)
	case ExportReportTypeDailyTrend:
		sheetName = "Daily Trend"
		err = s.exportDailyTrendToExcel(ctx, f, sheetName, req.Filter)
	default:
		return nil, fmt.Errorf("unsupported report type: %s", req.ReportType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export %s: %w", req.ReportType, err)
	}

	// Delete the default Sheet1
	index, err := f.GetSheetIndex("Sheet1")
	if err == nil && index != -1 {
		_ = f.DeleteSheet("Sheet1")
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel file: %w", err)
	}

	filename := fmt.Sprintf("admission-%s-%s.xlsx", req.ReportType, time.Now().Format("2006-01-02"))

	return &ExportResult{
		Data:        buf.Bytes(),
		Filename:    filename,
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}, nil
}

// exportDashboardToExcel exports dashboard stats to Excel.
func (s *ExportService) exportDashboardToExcel(ctx context.Context, f *excelize.File, sheetName string, filter DashboardFilter) error {
	stats, err := s.reportService.GetDashboardStats(ctx, filter)
	if err != nil {
		return err
	}

	// Create sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Set column widths
	_ = f.SetColWidth(sheetName, "A", "A", 30)
	_ = f.SetColWidth(sheetName, "B", "B", 15)

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"3B82F6"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Title row
	_ = f.SetCellValue(sheetName, "A1", "Admission Dashboard Report")
	_ = f.MergeCell(sheetName, "A1", "B1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	_ = f.SetCellStyle(sheetName, "A1", "B1", titleStyle)

	// Generated date
	_ = f.SetCellValue(sheetName, "A2", fmt.Sprintf("Generated: %s", time.Now().Format("02-Jan-2006 15:04")))
	_ = f.MergeCell(sheetName, "A2", "B2")

	// Headers
	_ = f.SetCellValue(sheetName, "A4", "Metric")
	_ = f.SetCellValue(sheetName, "B4", "Value")
	_ = f.SetCellStyle(sheetName, "A4", "B4", headerStyle)

	// Data rows
	data := [][]interface{}{
		{"Total Enquiries", stats.TotalEnquiries},
		{"Total Applications", stats.TotalApplications},
		{"Approved", stats.Approved},
		{"Enrolled", stats.Enrolled},
		{"Pending", stats.Pending},
		{"Rejected", stats.Rejected},
		{"Waitlisted", stats.Waitlisted},
		{"", ""},
		{"Conversion Rates", ""},
		{"Enquiry to Application", fmt.Sprintf("%.1f%%", stats.ConversionRates.EnquiryToApplication)},
		{"Application to Approved", fmt.Sprintf("%.1f%%", stats.ConversionRates.ApplicationToApproved)},
		{"Approved to Enrolled", fmt.Sprintf("%.1f%%", stats.ConversionRates.ApprovedToEnrolled)},
	}

	for i, row := range data {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+5), row[0])
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+5), row[1])
	}

	return nil
}

// exportClassWiseToExcel exports class-wise report to Excel.
func (s *ExportService) exportClassWiseToExcel(ctx context.Context, f *excelize.File, sheetName string, filter DashboardFilter) error {
	report, err := s.reportService.GetClassWiseReport(ctx, filter)
	if err != nil {
		return err
	}

	// Create sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Set column widths
	_ = f.SetColWidth(sheetName, "A", "A", 15)
	_ = f.SetColWidth(sheetName, "B", "G", 12)

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"3B82F6"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Title row
	_ = f.SetCellValue(sheetName, "A1", "Class-wise Admission Report")
	_ = f.MergeCell(sheetName, "A1", "G1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	_ = f.SetCellStyle(sheetName, "A1", "G1", titleStyle)

	// Generated date
	_ = f.SetCellValue(sheetName, "A2", fmt.Sprintf("Generated: %s", time.Now().Format("02-Jan-2006 15:04")))
	_ = f.MergeCell(sheetName, "A2", "G2")

	// Headers
	headers := []string{"Class", "Total Seats", "Applications", "Approved", "Enrolled", "Waitlisted", "Vacant"}
	for i, h := range headers {
		col := string(rune('A' + i))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s4", col), h)
	}
	_ = f.SetCellStyle(sheetName, "A4", "G4", headerStyle)

	// Data rows
	for i, c := range report.Classes {
		row := i + 5
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), c.ClassName)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), c.TotalSeats)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), c.Applications)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), c.Approved)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), c.Enrolled)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), c.Waitlisted)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), c.Vacant)
	}

	return nil
}

// exportFunnelToExcel exports funnel report to Excel.
func (s *ExportService) exportFunnelToExcel(ctx context.Context, f *excelize.File, sheetName string, filter DashboardFilter) error {
	report, err := s.reportService.GetFunnelReport(ctx, filter)
	if err != nil {
		return err
	}

	// Create sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Set column widths
	_ = f.SetColWidth(sheetName, "A", "A", 20)
	_ = f.SetColWidth(sheetName, "B", "C", 15)

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"3B82F6"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Title row
	_ = f.SetCellValue(sheetName, "A1", "Admission Conversion Funnel")
	_ = f.MergeCell(sheetName, "A1", "C1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	_ = f.SetCellStyle(sheetName, "A1", "C1", titleStyle)

	// Generated date
	_ = f.SetCellValue(sheetName, "A2", fmt.Sprintf("Generated: %s", time.Now().Format("02-Jan-2006 15:04")))
	_ = f.MergeCell(sheetName, "A2", "C2")

	// Headers
	_ = f.SetCellValue(sheetName, "A4", "Stage")
	_ = f.SetCellValue(sheetName, "B4", "Count")
	_ = f.SetCellValue(sheetName, "C4", "Percentage")
	_ = f.SetCellStyle(sheetName, "A4", "C4", headerStyle)

	// Data rows
	for i, stage := range report.Stages {
		row := i + 5
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), stage.Stage)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), stage.Count)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("%.1f%%", stage.Percentage))
	}

	return nil
}

// exportSourceAnalysisToExcel exports source analysis to Excel.
func (s *ExportService) exportSourceAnalysisToExcel(ctx context.Context, f *excelize.File, sheetName string, filter DashboardFilter) error {
	report, err := s.reportService.GetSourceAnalysis(ctx, filter)
	if err != nil {
		return err
	}

	// Create sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Set column widths
	_ = f.SetColWidth(sheetName, "A", "A", 20)
	_ = f.SetColWidth(sheetName, "B", "D", 15)

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"3B82F6"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Title row
	_ = f.SetCellValue(sheetName, "A1", "Enquiry Source Analysis")
	_ = f.MergeCell(sheetName, "A1", "D1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	_ = f.SetCellStyle(sheetName, "A1", "D1", titleStyle)

	// Generated date
	_ = f.SetCellValue(sheetName, "A2", fmt.Sprintf("Generated: %s", time.Now().Format("02-Jan-2006 15:04")))
	_ = f.MergeCell(sheetName, "A2", "D2")

	// Total count
	_ = f.SetCellValue(sheetName, "A3", fmt.Sprintf("Total Enquiries: %d", report.TotalCount))
	_ = f.MergeCell(sheetName, "A3", "D3")

	// Headers
	_ = f.SetCellValue(sheetName, "A5", "Source")
	_ = f.SetCellValue(sheetName, "B5", "Count")
	_ = f.SetCellValue(sheetName, "C5", "Percentage")
	_ = f.SetCellValue(sheetName, "D5", "Converted")
	_ = f.SetCellStyle(sheetName, "A5", "D5", headerStyle)

	// Data rows
	for i, source := range report.Sources {
		row := i + 6
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), source.Source)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), source.Count)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("%.1f%%", source.Percentage))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), source.Converted)
	}

	return nil
}

// exportDailyTrendToExcel exports daily trend to Excel.
func (s *ExportService) exportDailyTrendToExcel(ctx context.Context, f *excelize.File, sheetName string, filter DashboardFilter) error {
	report, err := s.reportService.GetDailyTrend(ctx, filter, 30)
	if err != nil {
		return err
	}

	// Create sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Set column widths
	_ = f.SetColWidth(sheetName, "A", "A", 15)
	_ = f.SetColWidth(sheetName, "B", "C", 15)

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"3B82F6"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Title row
	_ = f.SetCellValue(sheetName, "A1", "Daily Application Trend")
	_ = f.MergeCell(sheetName, "A1", "C1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	_ = f.SetCellStyle(sheetName, "A1", "C1", titleStyle)

	// Generated date
	_ = f.SetCellValue(sheetName, "A2", fmt.Sprintf("Generated: %s", time.Now().Format("02-Jan-2006 15:04")))
	_ = f.MergeCell(sheetName, "A2", "C2")

	// Headers
	_ = f.SetCellValue(sheetName, "A4", "Date")
	_ = f.SetCellValue(sheetName, "B4", "Enquiries")
	_ = f.SetCellValue(sheetName, "C4", "Applications")
	_ = f.SetCellStyle(sheetName, "A4", "C4", headerStyle)

	// Data rows
	for i, trend := range report.Trends {
		row := i + 5
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), trend.Date)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), trend.Enquiries)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), trend.Applications)
	}

	return nil
}
