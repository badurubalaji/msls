// Package auth provides authentication services for the MSLS application.
package auth

import "errors"

// Authentication-related errors.
var (
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidCredentials        = errors.New("invalid email or password")
	ErrAccountLocked             = errors.New("account is locked")
	ErrAccountInactive           = errors.New("account is not active")
	ErrEmailNotVerified          = errors.New("email is not verified")
	ErrTenantNotFound            = errors.New("tenant not found")
	ErrTenantInactive            = errors.New("tenant is not active")
	ErrRefreshTokenNotFound      = errors.New("refresh token not found")
	ErrRefreshTokenRevoked       = errors.New("refresh token has been revoked")
	ErrRefreshTokenExpired       = errors.New("refresh token has expired")
	ErrVerificationTokenNotFound = errors.New("verification token not found")
	ErrVerificationTokenUsed     = errors.New("verification token has already been used")
	ErrVerificationTokenExpired  = errors.New("verification token has expired")
	ErrEmailAlreadyExists        = errors.New("email already exists")
	ErrRoleNotFound              = errors.New("role not found")
)

// Password validation errors.
var (
	ErrPasswordTooShort     = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong      = errors.New("password must be at most 128 characters")
	ErrPasswordNoUppercase  = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase  = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoDigit      = errors.New("password must contain at least one digit")
	ErrPasswordNoSpecial    = errors.New("password must contain at least one special character")
	ErrPasswordHashMismatch = errors.New("password hash does not match")
	ErrInvalidPasswordHash  = errors.New("invalid password hash format")
)

// JWT-related errors.
var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidClaims    = errors.New("invalid token claims")
	ErrTokenNotYetValid = errors.New("token is not yet valid")
)

// OTP-related errors.
var (
	ErrOTPExpired          = errors.New("OTP has expired")
	ErrOTPInvalid          = errors.New("invalid OTP code")
	ErrOTPAlreadyUsed      = errors.New("OTP has already been used")
	ErrOTPMaxAttempts      = errors.New("maximum OTP verification attempts exceeded")
	ErrOTPRateLimited      = errors.New("too many OTP requests, please try again later")
	ErrOTPCooldown         = errors.New("please wait before requesting another OTP")
	ErrInvalidIdentifier   = errors.New("invalid phone number or email address")
	ErrIdentifierNotFound  = errors.New("no account found with this identifier")
	ErrSMSSendFailed       = errors.New("failed to send SMS")
	ErrEmailSendFailed     = errors.New("failed to send email")
	ErrInvalidOTPType      = errors.New("invalid OTP type")
	ErrInvalidOTPChannel   = errors.New("invalid OTP channel")
)

// TOTP/2FA-related errors.
var (
	ErrTOTPNotSetup          = errors.New("2FA setup not initiated")
	ErrTOTPAlreadyEnabled    = errors.New("2FA is already enabled")
	ErrTOTPNotEnabled        = errors.New("2FA is not enabled for this account")
	ErrTOTPInvalidCode       = errors.New("invalid 2FA code")
	ErrTOTPRateLimitExceeded = errors.New("too many 2FA attempts, please try again later")
	ErrTOTPRequired          = errors.New("2FA verification required")
)
