// Package staffdocument provides staff document management functionality.
package staffdocument

import "errors"

// Module errors.
var (
	// Document errors
	ErrDocumentNotFound      = errors.New("document not found")
	ErrDocumentTypeNotFound  = errors.New("document type not found")
	ErrInvalidDocumentType   = errors.New("invalid document type")
	ErrInvalidCategory       = errors.New("invalid document category")
	ErrInvalidStatus         = errors.New("invalid verification status")
	ErrFileTooLarge          = errors.New("file exceeds maximum size limit")
	ErrInvalidFileType       = errors.New("invalid file type, only PDF and images are allowed")
	ErrDuplicateDocument     = errors.New("document of this type already exists for this staff member")
	ErrDuplicateDocTypeCode  = errors.New("document type with this code already exists")
	ErrDocumentAlreadyVerified = errors.New("document is already verified")
	ErrDocumentAlreadyRejected = errors.New("document is already rejected")
	ErrMissingRejectionReason = errors.New("rejection reason is required")
	ErrDocumentTypeInUse     = errors.New("document type is in use and cannot be deleted")

	// File storage errors
	ErrFileUploadFailed   = errors.New("file upload failed")
	ErrFileDownloadFailed = errors.New("file download failed")
	ErrFileDeleteFailed   = errors.New("file delete failed")
)
