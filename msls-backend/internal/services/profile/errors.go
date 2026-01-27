// Package profile provides profile management services for the MSLS application.
package profile

import "errors"

// Profile service errors.
var (
	// ErrUserNotFound indicates the user was not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCurrentPassword indicates the current password is incorrect.
	ErrInvalidCurrentPassword = errors.New("current password is incorrect")

	// ErrPasswordMismatch indicates the new password and confirmation don't match.
	ErrPasswordMismatch = errors.New("new password and confirmation do not match")

	// ErrInvalidAvatarFormat indicates the avatar file format is not supported.
	ErrInvalidAvatarFormat = errors.New("invalid avatar format, only JPEG and PNG are allowed")

	// ErrAvatarTooLarge indicates the avatar file is too large.
	ErrAvatarTooLarge = errors.New("avatar file too large, maximum size is 2MB")

	// ErrPreferenceNotFound indicates the preference was not found.
	ErrPreferenceNotFound = errors.New("preference not found")

	// ErrAccountDeletionAlreadyRequested indicates account deletion was already requested.
	ErrAccountDeletionAlreadyRequested = errors.New("account deletion already requested")

	// ErrInvalidTimezone indicates an invalid timezone was provided.
	ErrInvalidTimezone = errors.New("invalid timezone")

	// ErrInvalidLocale indicates an invalid locale was provided.
	ErrInvalidLocale = errors.New("invalid locale")
)
