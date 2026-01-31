// Package hallticket provides hall ticket generation and management.
package hallticket

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// QRCodeGenerator generates and verifies QR codes for hall tickets.
type QRCodeGenerator struct {
	secret string
}

// NewQRCodeGenerator creates a new QR code generator with the given secret.
func NewQRCodeGenerator(secret string) *QRCodeGenerator {
	return &QRCodeGenerator{secret: secret}
}

// QRCodePayload represents the data stored in a hall ticket QR code.
type QRCodePayload struct {
	TicketID  string `json:"t"` // Hall ticket ID (short form)
	StudentID string `json:"s"` // Student ID (short form)
	ExamID    string `json:"e"` // Exam ID (short form)
	Hash      string `json:"v"` // Verification hash
}

// GenerateQRCodeData generates the QR code data string for a hall ticket.
func (g *QRCodeGenerator) GenerateQRCodeData(ticketID, studentID, examID uuid.UUID) (string, error) {
	// Create short form of IDs (first 8 chars of UUID)
	shortTicketID := ticketID.String()[:8]
	shortStudentID := studentID.String()[:8]
	shortExamID := examID.String()[:8]

	// Generate verification hash
	hashInput := fmt.Sprintf("%s:%s:%s:%s", ticketID.String(), studentID.String(), examID.String(), g.secret)
	hash := sha256.Sum256([]byte(hashInput))
	shortHash := hex.EncodeToString(hash[:])[:16] // First 16 chars of hash

	payload := QRCodePayload{
		TicketID:  shortTicketID,
		StudentID: shortStudentID,
		ExamID:    shortExamID,
		Hash:      shortHash,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal qr payload: %w", err)
	}

	return string(data), nil
}

// GenerateQRCodePNG generates a PNG image of the QR code.
func (g *QRCodeGenerator) GenerateQRCodePNG(data string, size int) ([]byte, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("encode qr code: %w", err)
	}
	return png, nil
}

// GenerateQRCodeBase64 generates a base64-encoded PNG of the QR code.
func (g *QRCodeGenerator) GenerateQRCodeBase64(data string, size int) (string, error) {
	png, err := g.GenerateQRCodePNG(data, size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}

// VerifyQRCodeData verifies that a QR code data string is valid for a given hall ticket.
func (g *QRCodeGenerator) VerifyQRCodeData(qrData string, ticketID, studentID, examID uuid.UUID) bool {
	var payload QRCodePayload
	if err := json.Unmarshal([]byte(qrData), &payload); err != nil {
		return false
	}

	// Verify short IDs match
	if payload.TicketID != ticketID.String()[:8] {
		return false
	}
	if payload.StudentID != studentID.String()[:8] {
		return false
	}
	if payload.ExamID != examID.String()[:8] {
		return false
	}

	// Verify hash
	hashInput := fmt.Sprintf("%s:%s:%s:%s", ticketID.String(), studentID.String(), examID.String(), g.secret)
	hash := sha256.Sum256([]byte(hashInput))
	expectedHash := hex.EncodeToString(hash[:])[:16]

	return payload.Hash == expectedHash
}

// ParseQRCodeData parses the QR code data and returns the payload.
func ParseQRCodeData(qrData string) (*QRCodePayload, error) {
	var payload QRCodePayload
	if err := json.Unmarshal([]byte(qrData), &payload); err != nil {
		return nil, fmt.Errorf("parse qr data: %w", err)
	}
	return &payload, nil
}
