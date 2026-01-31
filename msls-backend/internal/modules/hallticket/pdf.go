// Package hallticket provides hall ticket generation and management.
package hallticket

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
)

// PDFGenerator generates PDF documents for hall tickets.
type PDFGenerator struct {
	qrGenerator *QRCodeGenerator
}

// NewPDFGenerator creates a new PDF generator.
func NewPDFGenerator(qrGenerator *QRCodeGenerator) *PDFGenerator {
	return &PDFGenerator{
		qrGenerator: qrGenerator,
	}
}

// GenerateHallTicketPDF generates a PDF for a single hall ticket.
func (g *PDFGenerator) GenerateHallTicketPDF(data *HallTicketPDFData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	g.addHallTicketPage(pdf, data)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generate pdf: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateBatchPDF generates a PDF containing multiple hall tickets.
func (g *PDFGenerator) GenerateBatchPDF(tickets []*HallTicketPDFData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")

	for _, data := range tickets {
		g.addHallTicketPage(pdf, data)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generate batch pdf: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *PDFGenerator) addHallTicketPage(pdf *fpdf.Fpdf, data *HallTicketPDFData) {
	pdf.AddPage()
	pdf.SetMargins(15, 15, 15)

	pageWidth := 180.0 // 210 - 30 (margins)

	// Colors
	primaryColor := []int{31, 41, 55}     // Dark gray
	accentColor := []int{16, 185, 129}    // Green (Tailwind emerald-500)
	lightBg := []int{249, 250, 251}       // Very light gray
	borderColor := []int{229, 231, 235}   // Light border
	mutedText := []int{107, 114, 128}     // Muted gray

	// ========================================
	// HEADER SECTION
	// ========================================

	// School name from template
	schoolName := "School Name"
	if data.Template != nil && data.Template.SchoolName != "" {
		schoolName = data.Template.SchoolName
	}

	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(pageWidth, 10, schoolName, "", 1, "C", false, 0, "")

	// School address
	if data.Template != nil && data.Template.SchoolAddress != "" {
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
		pdf.CellFormat(pageWidth, 5, data.Template.SchoolAddress, "", 1, "C", false, 0, "")
	}

	pdf.Ln(5)

	// Title Box
	pdf.SetFillColor(accentColor[0], accentColor[1], accentColor[2])
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(pageWidth, 10, "HALL TICKET / ADMIT CARD", "0", 1, "C", true, 0, "")

	pdf.Ln(3)

	// Exam name
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(pageWidth, 8, data.ExamName, "", 1, "C", false, 0, "")

	// Exam period
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	examPeriod := fmt.Sprintf("%s to %s", data.ExamStartDate.Format("02 Jan 2006"), data.ExamEndDate.Format("02 Jan 2006"))
	pdf.CellFormat(pageWidth, 6, examPeriod, "", 1, "C", false, 0, "")

	pdf.Ln(5)

	// Divider
	pdf.SetDrawColor(borderColor[0], borderColor[1], borderColor[2])
	pdf.SetLineWidth(0.5)
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(5)

	// ========================================
	// STUDENT DETAILS SECTION
	// ========================================

	startY := pdf.GetY()

	// Left side - Student details
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])

	// Row 1: Roll Number
	pdf.SetXY(15, startY)
	pdf.Cell(35, 6, "Roll Number:")
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(accentColor[0], accentColor[1], accentColor[2])
	pdf.Cell(70, 6, data.HallTicket.RollNumber)

	// Row 2: Student Name
	pdf.SetXY(15, startY+10)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.Cell(35, 6, "Student Name:")
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.Cell(70, 6, data.HallTicket.StudentName)

	// Row 3: Admission Number
	pdf.SetXY(15, startY+20)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.Cell(35, 6, "Admission No:")
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.Cell(70, 6, data.HallTicket.AdmissionNumber)

	// Row 4: Class/Section
	pdf.SetXY(15, startY+30)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.Cell(35, 6, "Class / Section:")
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	classSection := data.HallTicket.ClassName
	if data.HallTicket.SectionName != "" {
		classSection += " - " + data.HallTicket.SectionName
	}
	pdf.Cell(70, 6, classSection)

	// Right side - QR Code
	if g.qrGenerator != nil {
		qrData := data.HallTicket.QRCodeData
		qrBase64, err := g.qrGenerator.GenerateQRCodeBase64(qrData, 150)
		if err == nil && qrBase64 != "" {
			// Register the image
			imgName := "qr_" + data.HallTicket.ID.String()
			pdf.RegisterImageOptionsReader(imgName, fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(mustDecodeBase64(qrBase64)))
			pdf.ImageOptions(imgName, 150, startY, 40, 40, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		}
	}

	pdf.SetY(startY + 45)
	pdf.Ln(5)

	// ========================================
	// EXAM SCHEDULE TABLE
	// ========================================

	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.Cell(pageWidth, 8, "Examination Schedule")
	pdf.Ln(8)

	// Table header
	pdf.SetFillColor(lightBg[0], lightBg[1], lightBg[2])
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.SetFont("Arial", "B", 9)

	colWidths := []float64{60, 35, 30, 30, 25}
	headers := []string{"Subject", "Date", "Time", "Max Marks", "Venue"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])

	for _, schedule := range data.ExamSchedules {
		pdf.CellFormat(colWidths[0], 7, schedule.SubjectName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidths[1], 7, schedule.ExamDate.Format("02 Jan 2006"), "1", 0, "C", false, 0, "")
		timeRange := fmt.Sprintf("%s-%s", schedule.StartTime, schedule.EndTime)
		pdf.CellFormat(colWidths[2], 7, timeRange, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[3], 7, fmt.Sprintf("%d", schedule.MaxMarks), "1", 0, "C", false, 0, "")
		venue := schedule.Venue
		if venue == "" {
			venue = "-"
		}
		pdf.CellFormat(colWidths[4], 7, venue, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(8)

	// ========================================
	// INSTRUCTIONS SECTION
	// ========================================

	if data.Template != nil && data.Template.Instructions != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
		pdf.Cell(pageWidth, 6, "Important Instructions:")
		pdf.Ln(6)

		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
		pdf.MultiCell(pageWidth, 5, data.Template.Instructions, "", "L", false)
	} else {
		// Default instructions
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
		pdf.Cell(pageWidth, 6, "Important Instructions:")
		pdf.Ln(6)

		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
		defaultInstructions := `1. Students must carry this hall ticket to the examination hall.
2. Students should report 30 minutes before the exam time.
3. Electronic devices are not allowed in the examination hall.
4. Students must follow all examination rules and regulations.
5. Any misconduct will result in disqualification.`
		pdf.MultiCell(pageWidth, 5, defaultInstructions, "", "L", false)
	}

	pdf.Ln(10)

	// ========================================
	// SIGNATURE SECTION
	// ========================================

	sigY := pdf.GetY()

	// Student signature
	pdf.SetXY(15, sigY)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.Cell(60, 20, "")
	pdf.Line(15, sigY+15, 75, sigY+15)
	pdf.SetXY(15, sigY+16)
	pdf.Cell(60, 5, "Student's Signature")

	// Principal signature
	pdf.SetXY(130, sigY)
	pdf.Cell(60, 20, "")
	pdf.Line(130, sigY+15, 190, sigY+15)
	pdf.SetXY(130, sigY+16)
	pdf.Cell(60, 5, "Principal's Signature & Seal")

	// ========================================
	// FOOTER
	// ========================================

	pdf.SetY(270)
	pdf.SetTextColor(mutedText[0], mutedText[1], mutedText[2])
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(pageWidth/2, 4, fmt.Sprintf("Generated on %s", time.Now().Format("02 Jan 2006 15:04")), "", 0, "L", false, 0, "")
	pdf.CellFormat(pageWidth/2, 4, "This is a computer-generated document", "", 0, "R", false, 0, "")
}

// GetFilename returns the filename for a hall ticket PDF.
func GetFilename(rollNumber string) string {
	return fmt.Sprintf("hall_ticket_%s.pdf", rollNumber)
}

// GetBatchFilename returns the filename for a batch hall ticket PDF.
func GetBatchFilename(examName string) string {
	// Sanitize exam name for filename
	safe := ""
	for _, r := range examName {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			safe += string(r)
		} else if r == ' ' {
			safe += "_"
		}
	}
	return fmt.Sprintf("hall_tickets_%s.pdf", safe)
}

func mustDecodeBase64(s string) []byte {
	data, _ := base64.StdEncoding.DecodeString(s)
	return data
}
