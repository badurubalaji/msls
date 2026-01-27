// Package admission provides admission management services.
package admission

import "errors"

// Service-level errors for admission operations.
var (
	// ErrSessionNotFound is returned when an admission session is not found.
	ErrSessionNotFound = errors.New("admission session not found")

	// ErrSeatNotFound is returned when an admission seat configuration is not found.
	ErrSeatNotFound = errors.New("admission seat configuration not found")

	// ErrSessionNameRequired is returned when session name is missing.
	ErrSessionNameRequired = errors.New("session name is required")

	// ErrSessionNameExists is returned when session name already exists for tenant and academic year.
	ErrSessionNameExists = errors.New("session with this name already exists for the academic year")

	// ErrInvalidDateRange is returned when end date is before start date.
	ErrInvalidDateRange = errors.New("end date must be after or equal to start date")

	// ErrInvalidStatus is returned when an invalid status is provided.
	ErrInvalidStatus = errors.New("invalid session status")

	// ErrInvalidStatusTransition is returned when status transition is not allowed.
	ErrInvalidStatusTransition = errors.New("invalid status transition")

	// ErrCannotDeleteOpenSession is returned when trying to delete an open session.
	ErrCannotDeleteOpenSession = errors.New("cannot delete an open admission session")

	// ErrSessionHasApplications is returned when session has applications.
	ErrSessionHasApplications = errors.New("session has applications and cannot be deleted")

	// ErrTenantIDRequired is returned when tenant ID is missing.
	ErrTenantIDRequired = errors.New("tenant ID is required")

	// ErrSessionIDRequired is returned when session ID is missing.
	ErrSessionIDRequired = errors.New("session ID is required")

	// ErrClassNameRequired is returned when class name is missing.
	ErrClassNameRequired = errors.New("class name is required")

	// ErrClassAlreadyExists is returned when a class seat configuration already exists for the session.
	ErrClassAlreadyExists = errors.New("seat configuration for this class already exists")

	// ErrInvalidTotalSeats is returned when total seats is invalid.
	ErrInvalidTotalSeats = errors.New("total seats must be greater than or equal to zero")

	// ErrFilledExceedsTotal is returned when filled seats would exceed total seats.
	ErrFilledExceedsTotal = errors.New("filled seats cannot exceed total seats")

	// ErrCannotModifyClosedSession is returned when trying to modify a closed session.
	ErrCannotModifyClosedSession = errors.New("cannot modify a closed admission session")

	// Enquiry-related errors

	// ErrEnquiryNotFound is returned when an enquiry is not found.
	ErrEnquiryNotFound = errors.New("enquiry not found")

	// ErrFollowUpNotFound is returned when a follow-up is not found.
	ErrFollowUpNotFound = errors.New("follow-up not found")

	// ErrStudentNameRequired is returned when student name is missing.
	ErrStudentNameRequired = errors.New("student name is required")

	// ErrClassApplyingRequired is returned when class applying is missing.
	ErrClassApplyingRequired = errors.New("class applying is required")

	// ErrParentNameRequired is returned when parent name is missing.
	ErrParentNameRequired = errors.New("parent name is required")

	// ErrParentPhoneRequired is returned when parent phone is missing.
	ErrParentPhoneRequired = errors.New("parent phone is required")

	// ErrInvalidEnquiryStatus is returned when an invalid status is provided.
	ErrInvalidEnquiryStatus = errors.New("invalid enquiry status")

	// ErrInvalidEnquirySource is returned when an invalid source is provided.
	ErrInvalidEnquirySource = errors.New("invalid enquiry source")

	// ErrInvalidContactMode is returned when an invalid contact mode is provided.
	ErrInvalidContactMode = errors.New("invalid contact mode")

	// ErrInvalidFollowUpOutcome is returned when an invalid follow-up outcome is provided.
	ErrInvalidFollowUpOutcome = errors.New("invalid follow-up outcome")

	// ErrEnquiryAlreadyConverted is returned when trying to convert an already converted enquiry.
	ErrEnquiryAlreadyConverted = errors.New("enquiry has already been converted to an application")

	// ErrEnquiryClosed is returned when trying to update a closed enquiry.
	ErrEnquiryClosed = errors.New("enquiry is closed and cannot be modified")

	// ErrSessionRequired is returned when session ID is required but not provided.
	ErrSessionRequired = errors.New("session ID is required for conversion")

	// ErrFollowUpDateRequired is returned when follow-up date is missing.
	ErrFollowUpDateRequired = errors.New("follow-up date is required")

	// ErrInvalidGender is returned when an invalid gender is provided.
	ErrInvalidGender = errors.New("invalid gender value")

	// Merit List and Decision errors

	// ErrApplicationNotFound is returned when an application is not found.
	ErrApplicationNotFound = errors.New("application not found")

	// ErrDecisionNotFound is returned when a decision is not found.
	ErrDecisionNotFound = errors.New("admission decision not found")

	// ErrMeritListNotFound is returned when a merit list is not found.
	ErrMeritListNotFound = errors.New("merit list not found")

	// ErrMeritListExists is returned when a merit list already exists.
	ErrMeritListExists = errors.New("merit list already exists for this session, class, and test")

	// ErrMeritListFinalized is returned when trying to modify a finalized merit list.
	ErrMeritListFinalized = errors.New("merit list is already finalized")

	// ErrDecisionExists is returned when a decision already exists for an application.
	ErrDecisionExists = errors.New("decision already exists for this application")

	// ErrInvalidDecisionType is returned when an invalid decision type is provided.
	ErrInvalidDecisionType = errors.New("invalid decision type")

	// ErrWaitlistPositionRequired is returned when waitlist position is missing for waitlisted decision.
	ErrWaitlistPositionRequired = errors.New("waitlist position is required for waitlisted decision")

	// ErrRejectionReasonRequired is returned when rejection reason is missing for rejected decision.
	ErrRejectionReasonRequired = errors.New("rejection reason is required for rejected decision")

	// ErrOfferNotFound is returned when no offer exists for an application.
	ErrOfferNotFound = errors.New("no offer found for this application")

	// ErrOfferExpired is returned when the offer has expired.
	ErrOfferExpired = errors.New("offer has expired")

	// ErrOfferAlreadyAccepted is returned when the offer was already accepted.
	ErrOfferAlreadyAccepted = errors.New("offer has already been accepted")

	// ErrOfferNotAccepted is returned when trying to enroll without accepting the offer.
	ErrOfferNotAccepted = errors.New("offer must be accepted before enrollment")

	// ErrAlreadyEnrolled is returned when the application is already enrolled.
	ErrAlreadyEnrolled = errors.New("application is already enrolled")

	// ErrNoApplicantsForMeritList is returned when there are no applicants to generate a merit list.
	ErrNoApplicantsForMeritList = errors.New("no applicants found to generate merit list")

	// ErrApplicationIDRequired is returned when application ID is missing.
	ErrApplicationIDRequired = errors.New("application ID is required")

	// ErrDecisionRequired is returned when decision is missing.
	ErrDecisionRequired = errors.New("decision is required")

	// ErrDecisionDateRequired is returned when decision date is missing.
	ErrDecisionDateRequired = errors.New("decision date is required")

	// ErrInvalidApplicationStatus is returned when application status is invalid for the operation.
	ErrInvalidApplicationStatus = errors.New("invalid application status for this operation")

	// Document-related errors

	// ErrDocumentNotFound is returned when a document is not found.
	ErrDocumentNotFound = errors.New("document not found")

	// ErrDocumentAlreadyVerified is returned when trying to verify an already verified document.
	ErrDocumentAlreadyVerified = errors.New("document already verified")

	// ErrInvalidVerificationStatus is returned when an invalid verification status is provided.
	ErrInvalidVerificationStatus = errors.New("invalid verification status")

	// Review-related errors

	// ErrReviewNotFound is returned when a review is not found.
	ErrReviewNotFound = errors.New("review not found")

	// ErrInvalidReviewType is returned when an invalid review type is provided.
	ErrInvalidReviewType = errors.New("invalid review type")

	// ErrInvalidReviewStatus is returned when an invalid review status is provided.
	ErrInvalidReviewStatus = errors.New("invalid review status")

	// Entrance test-related errors

	// ErrTestNotFound is returned when an entrance test is not found.
	ErrTestNotFound = errors.New("entrance test not found")

	// ErrTestAlreadyExists is returned when test with same name exists.
	ErrTestAlreadyExists = errors.New("entrance test with same name already exists")

	// ErrInvalidTestStatus is returned when an invalid test status is provided.
	ErrInvalidTestStatus = errors.New("invalid test status")

	// ErrTestFull is returned when test has reached maximum candidates.
	ErrTestFull = errors.New("test has reached maximum candidates")

	// ErrTestNotScheduled is returned when test is not in scheduled status.
	ErrTestNotScheduled = errors.New("test is not in scheduled status")

	// ErrTestDatePassed is returned when test date has already passed.
	ErrTestDatePassed = errors.New("test date has already passed")

	// ErrCannotModifyCompletedTest is returned when trying to modify a completed test.
	ErrCannotModifyCompletedTest = errors.New("cannot modify a completed test")

	// ErrTestNameRequired is returned when test name is missing.
	ErrTestNameRequired = errors.New("test name is required")

	// ErrTestDateRequired is returned when test date is missing.
	ErrTestDateRequired = errors.New("test date is required")

	// ErrTestStartTimeRequired is returned when test start time is missing.
	ErrTestStartTimeRequired = errors.New("test start time is required")

	// Registration-related errors

	// ErrRegistrationNotFound is returned when a test registration is not found.
	ErrRegistrationNotFound = errors.New("test registration not found")

	// ErrAlreadyRegistered is returned when application is already registered for the test.
	ErrAlreadyRegistered = errors.New("application already registered for this test")

	// ErrInvalidRegistrationStatus is returned when an invalid registration status is provided.
	ErrInvalidRegistrationStatus = errors.New("invalid registration status")

	// ErrCannotSubmitResults is returned when results cannot be submitted for registration.
	ErrCannotSubmitResults = errors.New("cannot submit results for this registration")

	// ErrResultsAlreadySubmitted is returned when results are already submitted.
	ErrResultsAlreadySubmitted = errors.New("results already submitted for this registration")

	// ErrCannotUpdateApplication is returned when application cannot be updated.
	ErrCannotUpdateApplication = errors.New("cannot update application in current status")

	// Story 3.5: Online Admission Application - Additional errors

	// ErrSessionClosed is returned when trying to create an application for a closed session.
	ErrSessionClosed = errors.New("admission session is closed")

	// ErrApplicationAlreadySubmitted is returned when trying to modify a submitted application.
	ErrApplicationAlreadySubmitted = errors.New("application has already been submitted")

	// ErrInvalidStageTransition is returned when the stage transition is not allowed.
	ErrInvalidStageTransition = errors.New("invalid stage transition")

	// ErrMissingRequiredFields is returned when required fields are not provided.
	ErrMissingRequiredFields = errors.New("missing required fields")

	// ErrMissingRequiredDocuments is returned when required documents are not uploaded.
	ErrMissingRequiredDocuments = errors.New("missing required documents")

	// ErrParentNotFound is returned when a parent record is not found.
	ErrParentNotFound = errors.New("parent not found")

	// ErrInvalidParentRelation is returned when an invalid parent relation is provided.
	ErrInvalidParentRelation = errors.New("invalid parent relation")

	// ErrInvalidDocumentType is returned when an invalid document type is provided.
	ErrInvalidDocumentType = errors.New("invalid document type")

	// ErrApplicationNumberConflict is returned when application number already exists.
	ErrApplicationNumberConflict = errors.New("application number already exists")

	// ErrCannotDeleteSubmittedApplication is returned when trying to delete a submitted application.
	ErrCannotDeleteSubmittedApplication = errors.New("cannot delete a submitted application")

	// ErrInvalidPhoneNumber is returned when the phone number format is invalid.
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")

	// ErrApplicationNotInDraft is returned when trying to submit a non-draft application.
	ErrApplicationNotInDraft = errors.New("application is not in draft status")
)

// StageTransitionError provides detailed information about invalid stage transitions.
type StageTransitionError struct {
	CurrentStatus    string
	RequestedStatus  string
	ValidTransitions []string
}

func (e *StageTransitionError) Error() string {
	return "invalid stage transition"
}

// NewStageTransitionError creates a new StageTransitionError with the given details.
func NewStageTransitionError(current, requested string, validTransitions []string) *StageTransitionError {
	return &StageTransitionError{
		CurrentStatus:    current,
		RequestedStatus:  requested,
		ValidTransitions: validTransitions,
	}
}
