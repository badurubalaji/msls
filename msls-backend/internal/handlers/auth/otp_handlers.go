// Package auth provides HTTP handlers for authentication endpoints.
package auth

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	authservice "msls-backend/internal/services/auth"
)

// OTPHandler handles OTP-related HTTP requests.
type OTPHandler struct {
	otpService *authservice.OTPService
}

// NewOTPHandler creates a new OTP Handler.
func NewOTPHandler(otpService *authservice.OTPService) *OTPHandler {
	return &OTPHandler{
		otpService: otpService,
	}
}

// RequestOTP handles OTP request.
// @Summary Request OTP
// @Description Request an OTP for passwordless login
// @Tags Auth/OTP
// @Accept json
// @Produce json
// @Param request body OTPRequestRequest true "OTP request details"
// @Success 200 {object} response.Success{data=OTPRequestResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 429 {object} apperrors.AppError
// @Router /api/v1/auth/otp/request [post]
func (h *OTPHandler) RequestOTP(c *gin.Context) {
	var req OTPRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Determine OTP channel from type
	var channel models.OTPChannel
	switch req.Type {
	case "sms":
		channel = models.OTPChannelSMS
	case "email":
		channel = models.OTPChannelEmail
	default:
		apperrors.Abort(c, apperrors.BadRequest("Invalid type: must be 'sms' or 'email'"))
		return
	}

	// Parse tenant ID if provided
	var tenantID uuid.UUID
	if req.TenantID != "" {
		var err error
		tenantID, err = uuid.Parse(req.TenantID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid tenant ID"))
			return
		}
	}

	// Create service request
	serviceReq := authservice.RequestOTPRequest{
		Identifier: req.Identifier,
		Type:       models.OTPTypeLogin,
		Channel:    channel,
		TenantID:   tenantID,
	}

	result, err := h.otpService.RequestOTP(c.Request.Context(), serviceReq)
	if err != nil {
		handleOTPError(c, err)
		return
	}

	resp := OTPRequestResponse{
		Message:          result.Message,
		ExpiresIn:        result.ExpiresIn,
		MaskedIdentifier: result.MaskedIdentifier,
	}

	response.OK(c, resp)
}

// VerifyOTP handles OTP verification.
// @Summary Verify OTP
// @Description Verify an OTP and login
// @Tags Auth/OTP
// @Accept json
// @Produce json
// @Param request body OTPVerifyRequest true "OTP verification details"
// @Success 200 {object} response.Success{data=OTPVerifyResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 423 {object} apperrors.AppError
// @Router /api/v1/auth/otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse tenant ID if provided
	var tenantID uuid.UUID
	if req.TenantID != "" {
		var err error
		tenantID, err = uuid.Parse(req.TenantID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid tenant ID"))
			return
		}
	}

	// Get client IP and user agent
	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")

	// Create service request
	serviceReq := authservice.VerifyOTPRequest{
		Identifier: req.Identifier,
		Code:       req.Code,
		TenantID:   tenantID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	tokenPair, user, err := h.otpService.VerifyOTP(c.Request.Context(), serviceReq)
	if err != nil {
		handleOTPError(c, err)
		return
	}

	resp := OTPVerifyResponse{
		User:         userToDTO(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}

	response.OK(c, resp)
}

// ResendOTP handles OTP resend requests.
// @Summary Resend OTP
// @Description Resend an OTP (with cooldown)
// @Tags Auth/OTP
// @Accept json
// @Produce json
// @Param request body OTPResendRequest true "OTP resend details"
// @Success 200 {object} response.Success{data=OTPRequestResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 429 {object} apperrors.AppError
// @Router /api/v1/auth/otp/resend [post]
func (h *OTPHandler) ResendOTP(c *gin.Context) {
	var req OTPResendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Determine OTP channel from type
	var channel models.OTPChannel
	switch req.Type {
	case "sms":
		channel = models.OTPChannelSMS
	case "email":
		channel = models.OTPChannelEmail
	default:
		apperrors.Abort(c, apperrors.BadRequest("Invalid type: must be 'sms' or 'email'"))
		return
	}

	// Parse tenant ID if provided
	var tenantID uuid.UUID
	if req.TenantID != "" {
		var err error
		tenantID, err = uuid.Parse(req.TenantID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid tenant ID"))
			return
		}
	}

	// Create service request
	serviceReq := authservice.RequestOTPRequest{
		Identifier: req.Identifier,
		Type:       models.OTPTypeLogin,
		Channel:    channel,
		TenantID:   tenantID,
	}

	result, err := h.otpService.ResendOTP(c.Request.Context(), serviceReq)
	if err != nil {
		handleOTPError(c, err)
		return
	}

	resp := OTPRequestResponse{
		Message:          result.Message,
		ExpiresIn:        result.ExpiresIn,
		MaskedIdentifier: result.MaskedIdentifier,
	}

	response.OK(c, resp)
}

// handleOTPError handles OTP-related errors and returns appropriate HTTP responses.
func handleOTPError(c *gin.Context, err error) {
	switch err {
	case authservice.ErrOTPExpired:
		apperrors.Abort(c, apperrors.BadRequest("OTP has expired. Please request a new one."))
	case authservice.ErrOTPInvalid:
		apperrors.Abort(c, apperrors.Unauthorized("Invalid OTP code"))
	case authservice.ErrOTPAlreadyUsed:
		apperrors.Abort(c, apperrors.BadRequest("OTP has already been used"))
	case authservice.ErrOTPMaxAttempts:
		apperrors.Abort(c, apperrors.BadRequest("Maximum verification attempts exceeded. Please request a new OTP."))
	case authservice.ErrOTPRateLimited:
		apperrors.Abort(c, apperrors.TooManyRequests(3600)) // 1 hour
	case authservice.ErrOTPCooldown:
		apperrors.Abort(c, apperrors.TooManyRequests(60)) // 60 seconds
	case authservice.ErrInvalidIdentifier:
		apperrors.Abort(c, apperrors.BadRequest("Invalid phone number or email address"))
	case authservice.ErrIdentifierNotFound:
		// Don't reveal whether the identifier exists for security
		apperrors.Abort(c, apperrors.BadRequest("Unable to send OTP. Please check the identifier."))
	case authservice.ErrSMSSendFailed:
		apperrors.Abort(c, apperrors.InternalError("Failed to send SMS. Please try again."))
	case authservice.ErrEmailSendFailed:
		apperrors.Abort(c, apperrors.InternalError("Failed to send email. Please try again."))
	case authservice.ErrAccountLocked:
		apperrors.Abort(c, &apperrors.AppError{
			Type:   "https://httpstatuses.com/423",
			Title:  "Locked",
			Status: http.StatusLocked,
			Detail: "Account is locked. Please try again later.",
		})
	case authservice.ErrAccountInactive:
		apperrors.Abort(c, apperrors.Unauthorized("Account is not active"))
	case authservice.ErrInvalidOTPChannel:
		apperrors.Abort(c, apperrors.BadRequest("Invalid OTP channel"))
	default:
		apperrors.Abort(c, apperrors.InternalError("An error occurred. Please try again."))
	}
}
