// Package payroll provides PDF generation for payslips.
package payroll

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// PDFGenerator generates PDF documents for payroll.
type PDFGenerator struct {
	companyName    string
	companyAddress string
}

// NewPDFGenerator creates a new PDF generator.
func NewPDFGenerator(companyName, companyAddress string) *PDFGenerator {
	return &PDFGenerator{
		companyName:    companyName,
		companyAddress: companyAddress,
	}
}

// GeneratePayslipPDF generates a PDF for a payslip.
func (g *PDFGenerator) GeneratePayslipPDF(payslip *models.Payslip, payRun *models.PayRun) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pageWidth := 170.0 // 210 - 40 (margins)

	// Colors
	primaryColor := []int{31, 41, 55}     // Dark gray
	accentColor := []int{79, 70, 229}     // Indigo
	lightBg := []int{249, 250, 251}       // Very light gray
	borderColor := []int{229, 231, 235}   // Light border
	successColor := []int{22, 163, 74}    // Green
	mutedText := []int{107, 114, 128}     // Muted gray

	// ========================================
	// HEADER SECTION
	// ========================================

	// Company name
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 8, g.companyName)
	pdf.Ln(10)

	// Title
	pdf.SetTextColor(accentColor[0], accentColor[1], accentColor[2])
	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(0, 12, "Salary Slip")
	pdf.Ln(8)

	// Pay period subtitle
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 6, formatPayPeriodPDF(payRun.PayPeriodMonth, payRun.PayPeriodYear))
	pdf.Ln(12)

	// Divider line
	pdf.SetDrawColor(borderColor[0], borderColor[1], borderColor[2])
	pdf.SetLineWidth(0.5)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(8)

	// ========================================
	// EMPLOYEE DETAILS SECTION
	// ========================================

	y := pdf.GetY()

	// Left column - Employee info
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.SetXY(20, y)
	pdf.Cell(40, 5, "EMPLOYEE NAME")

	pdf.SetXY(20, y+12)
	pdf.Cell(40, 5, "EMPLOYEE ID")

	pdf.SetXY(20, y+24)
	pdf.Cell(40, 5, "DEPARTMENT")

	// Left column - Values
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])

	staffName := getStaffName(payslip)
	pdf.SetXY(20, y+5)
	pdf.Cell(60, 6, staffName)

	pdf.SetXY(20, y+17)
	pdf.Cell(60, 6, payslip.Staff.EmployeeID)

	deptName := "-"
	if payslip.Staff.Department != nil {
		deptName = payslip.Staff.Department.Name
	}
	pdf.SetXY(20, y+29)
	pdf.Cell(60, 6, deptName)

	// Right column - Pay info
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.SetXY(120, y)
	pdf.Cell(40, 5, "PAY PERIOD")

	pdf.SetXY(120, y+12)
	pdf.Cell(40, 5, "WORKING DAYS")

	pdf.SetXY(120, y+24)
	pdf.Cell(40, 5, "DAYS WORKED")

	// Right column - Values
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])

	pdf.SetXY(120, y+5)
	pdf.Cell(60, 6, formatPayPeriodPDF(payRun.PayPeriodMonth, payRun.PayPeriodYear))

	pdf.SetXY(120, y+17)
	pdf.Cell(60, 6, fmt.Sprintf("%d", payslip.WorkingDays))

	pdf.SetXY(120, y+29)
	pdf.Cell(60, 6, fmt.Sprintf("%d", payslip.PresentDays))

	pdf.SetY(y + 42)
	pdf.Ln(5)

	// ========================================
	// EARNINGS & DEDUCTIONS TABLE
	// ========================================

	tableY := pdf.GetY()
	colWidth := pageWidth / 2

	// Table headers
	pdf.SetFillColor(lightBg[0], lightBg[1], lightBg[2])
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.SetFont("Arial", "B", 10)

	// Earnings header
	pdf.SetXY(20, tableY)
	pdf.CellFormat(colWidth-5, 10, "  Earnings", "1", 0, "L", true, 0, "")

	// Deductions header
	pdf.SetXY(20+colWidth, tableY)
	pdf.CellFormat(colWidth-5, 10, "  Deductions", "1", 0, "L", true, 0, "")

	tableY += 12

	// Earnings items
	earningsY := tableY
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])

	for _, comp := range payslip.Components {
		if comp.ComponentType == "earning" {
			pdf.SetXY(22, earningsY)
			pdf.Cell(50, 6, comp.ComponentName)
			pdf.SetXY(72, earningsY)
			pdf.CellFormat(colWidth-57, 6, formatCurrencyPDF(comp.Amount), "", 0, "R", false, 0, "")
			earningsY += 7
		}
	}

	// Deductions items
	deductionsY := tableY
	for _, comp := range payslip.Components {
		if comp.ComponentType == "deduction" {
			pdf.SetXY(20+colWidth+2, deductionsY)
			pdf.Cell(50, 6, comp.ComponentName)
			pdf.SetXY(20+colWidth+52, deductionsY)
			pdf.CellFormat(colWidth-57, 6, formatCurrencyPDF(comp.Amount), "", 0, "R", false, 0, "")
			deductionsY += 7
		}
	}

	// LOP Deduction if applicable
	if payslip.LOPDeduction.GreaterThan(decimal.Zero) {
		pdf.SetTextColor(200, 50, 50)
		pdf.SetXY(20+colWidth+2, deductionsY)
		pdf.Cell(50, 6, fmt.Sprintf("LOP (%d days)", payslip.LOPDays))
		pdf.SetXY(20+colWidth+52, deductionsY)
		pdf.CellFormat(colWidth-57, 6, formatCurrencyPDF(payslip.LOPDeduction), "", 0, "R", false, 0, "")
		deductionsY += 7
		pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	}

	// Find max Y for totals line
	maxY := earningsY
	if deductionsY > earningsY {
		maxY = deductionsY
	}
	maxY += 3

	// Divider lines
	pdf.SetDrawColor(borderColor[0], borderColor[1], borderColor[2])
	pdf.Line(22, maxY, 20+colWidth-7, maxY)
	pdf.Line(20+colWidth+2, maxY, 188, maxY)

	// Totals
	maxY += 5
	pdf.SetFont("Arial", "B", 10)

	// Total Earnings
	pdf.SetTextColor(successColor[0], successColor[1], successColor[2])
	pdf.SetXY(22, maxY)
	pdf.Cell(50, 6, "Total Earnings")
	pdf.SetXY(72, maxY)
	pdf.CellFormat(colWidth-57, 6, formatCurrencyPDF(payslip.TotalEarnings), "", 0, "R", false, 0, "")

	// Total Deductions
	pdf.SetTextColor(200, 50, 50)
	pdf.SetXY(20+colWidth+2, maxY)
	pdf.Cell(50, 6, "Total Deductions")
	pdf.SetXY(20+colWidth+52, maxY)
	pdf.CellFormat(colWidth-57, 6, formatCurrencyPDF(payslip.TotalDeductions), "", 0, "R", false, 0, "")

	// ========================================
	// NET PAY SECTION
	// ========================================

	netY := maxY + 20

	// Net pay box
	pdf.SetFillColor(accentColor[0], accentColor[1], accentColor[2])
	pdf.RoundedRect(20, netY, pageWidth, 30, 4, "1234", "F")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "", 11)
	pdf.SetXY(30, netY+8)
	pdf.Cell(40, 6, "Net Pay")

	pdf.SetFont("Arial", "B", 20)
	pdf.SetXY(80, netY+6)
	pdf.CellFormat(100, 10, formatCurrencyPDF(payslip.NetSalary), "", 0, "R", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(30, netY+18)
	pdf.Cell(0, 5, fmt.Sprintf("Gross: %s  -  Deductions: %s",
		formatCurrencyPDF(payslip.GrossSalary),
		formatCurrencyPDF(payslip.TotalDeductions)))

	// ========================================
	// ATTENDANCE SUMMARY (if there are leaves/absents)
	// ========================================

	if payslip.LeaveDays > 0 || payslip.AbsentDays > 0 || payslip.LOPDays > 0 {
		summaryY := netY + 40

		pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
		pdf.SetFont("Arial", "", 9)
		pdf.SetXY(20, summaryY)

		parts := []string{}
		if payslip.LeaveDays > 0 {
			parts = append(parts, fmt.Sprintf("Leave: %d days", payslip.LeaveDays))
		}
		if payslip.AbsentDays > 0 {
			parts = append(parts, fmt.Sprintf("Absent: %d days", payslip.AbsentDays))
		}
		if payslip.LOPDays > 0 {
			parts = append(parts, fmt.Sprintf("Loss of Pay: %d days", payslip.LOPDays))
		}

		pdf.Cell(0, 5, "Attendance Note: "+strings.Join(parts, " | "))
	}

	// ========================================
	// FOOTER
	// ========================================

	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.SetFont("Arial", "I", 8)
	pdf.SetXY(20, 265)
	pdf.Cell(0, 4, "This is a computer-generated document and does not require a signature.")
	pdf.SetXY(20, 270)
	pdf.Cell(0, 4, fmt.Sprintf("Generated on %s", time.Now().Format("02 Jan 2006")))

	// Generate PDF bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// GetPayslipFilename returns the filename for a payslip PDF.
func GetPayslipFilename(payRun *models.PayRun) string {
	months := []string{
		"", "jan", "feb", "mar", "apr", "may", "jun",
		"jul", "aug", "sep", "oct", "nov", "dec",
	}
	monthStr := "unknown"
	if payRun.PayPeriodMonth >= 1 && payRun.PayPeriodMonth <= 12 {
		monthStr = months[payRun.PayPeriodMonth]
	}
	return fmt.Sprintf("payslip_%s%d.pdf", monthStr, payRun.PayPeriodYear)
}

// Helper functions

func getStaffName(payslip *models.Payslip) string {
	if payslip.Staff.FirstName == "" {
		return "Employee"
	}
	name := payslip.Staff.FirstName
	if payslip.Staff.LastName != "" {
		name += " " + payslip.Staff.LastName
	}
	return name
}

func formatPayPeriodPDF(month, year int) string {
	months := []string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	if month >= 1 && month <= 12 {
		return fmt.Sprintf("%s %d", months[month], year)
	}
	return fmt.Sprintf("%d/%d", month, year)
}

func formatCurrencyPDF(amount decimal.Decimal) string {
	// Use "Rs." instead of Unicode rupee symbol for PDF compatibility
	return fmt.Sprintf("Rs. %s", formatNumberWithCommas(amount.StringFixed(2)))
}

func formatNumberWithCommas(s string) string {
	// Split into integer and decimal parts
	parts := strings.Split(s, ".")
	intPart := parts[0]

	// Indian number format: last 3 digits, then groups of 2
	n := len(intPart)
	if n <= 3 {
		if len(parts) > 1 {
			return intPart + "." + parts[1]
		}
		return intPart
	}

	// Start from the right
	result := intPart[n-3:]
	remaining := intPart[:n-3]

	for len(remaining) > 0 {
		if len(remaining) >= 2 {
			result = remaining[len(remaining)-2:] + "," + result
			remaining = remaining[:len(remaining)-2]
		} else {
			result = remaining + "," + result
			remaining = ""
		}
	}

	if len(parts) > 1 {
		return result + "." + parts[1]
	}
	return result
}
