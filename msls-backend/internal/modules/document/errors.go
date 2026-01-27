// Package document provides student document management functionality.
package document

import "errors"

// Domain errors for document module.
var (
	ErrDocumentTypeNotFound     = errors.New("document type not found")
	ErrDocumentNotFound         = errors.New("document not found")
	ErrStudentNotFound          = errors.New("student not found")
	ErrNoFileUploaded           = errors.New("no file uploaded")
	ErrFileTooLarge             = errors.New("file size exceeds maximum allowed")
	ErrInvalidFileType          = errors.New("invalid file type")
	ErrDocumentAlreadyExists    = errors.New("document of this type already exists for the student")
	ErrDocumentAlreadyVerified  = errors.New("document is already verified")
	ErrInvalidDocumentStatus    = errors.New("invalid document status")
	ErrRejectionReasonRequired  = errors.New("rejection reason is required")
	ErrOptimisticLockConflict   = errors.New("document was modified by another user")
	ErrInvalidExpiryDate        = errors.New("expiry date must be in the future")
	ErrInvalidIssueDate         = errors.New("issue date cannot be in the future")
)
