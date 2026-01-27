// Package auth provides HTTP handlers for authentication endpoints.
package auth

import (
	"net"

	"github.com/gin-gonic/gin"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
	"msls-backend/internal/pkg/response"
	authservice "msls-backend/internal/services/auth"
)

// TwoFactorHandler handles 2FA HTTP requests.
type TwoFactorHandler struct {
	authService *authservice.AuthService
	totpService *authservice.TOTPService
}

// NewTwoFactorHandler creates a new TwoFactorHandler.
func NewTwoFactorHandler(authService *authservice.AuthService, totpService *authservice.TOTPService) *TwoFactorHandler {
	return &TwoFactorHandler{
		authService: authService,
		totpService: totpService,
	}
}

// Setup2FA handles 2FA setup (generates TOTP secret and QR code).
// @Summary Setup 2FA
// @Description Generate TOTP secret and QR code for 2FA setup
// @Tags Auth - 2FA
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=TwoFactorSetupResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/auth/2fa/setup [post]
func (h *TwoFactorHandler) Setup2FA(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == authservice.ErrUserNotFound {
			apperrors.Abort(c, apperrors.NotFound("User not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get user"))
		return
	}

	// Check if 2FA is already enabled
	if user.TwoFactorEnabled && user.TOTPVerifiedAt != nil {
		apperrors.Abort(c, apperrors.Conflict("2FA is already enabled"))
		return
	}

	// Generate TOTP secret
	result, err := h.totpService.SetupTOTP(c.Request.Context(), user)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to setup 2FA"))
		return
	}

	// Create audit log
	h.createAuditLog(c, user, models.AuditAction2FASetup)

	response.OK(c, TwoFactorSetupResponse{
		Secret:        result.Secret,
		QRCodeDataURL: result.QRCodeDataURL,
		ManualEntry:   result.ManualEntry,
	})
}

// Verify2FA handles verification and enabling of 2FA.
// @Summary Verify and enable 2FA
// @Description Verify TOTP code and enable 2FA for the user
// @Tags Auth - 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TwoFactorVerifyRequest true "TOTP code"
// @Success 200 {object} response.Success{data=TwoFactorVerifyResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Failure 429 {object} apperrors.AppError
// @Router /api/v1/auth/2fa/verify [post]
func (h *TwoFactorHandler) Verify2FA(c *gin.Context) {
	var req TwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == authservice.ErrUserNotFound {
			apperrors.Abort(c, apperrors.NotFound("User not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get user"))
		return
	}

	ipAddress := net.ParseIP(c.ClientIP())

	// Verify and enable 2FA
	backupCodes, err := h.totpService.VerifyAndEnableTOTP(c.Request.Context(), user, req.Code, ipAddress)
	if err != nil {
		switch err {
		case authservice.ErrTOTPNotSetup:
			apperrors.Abort(c, apperrors.BadRequest("2FA setup not initiated. Please call /2fa/setup first."))
		case authservice.ErrTOTPAlreadyEnabled:
			apperrors.Abort(c, apperrors.Conflict("2FA is already enabled"))
		case authservice.ErrTOTPInvalidCode:
			apperrors.Abort(c, apperrors.BadRequest("Invalid 2FA code"))
		case authservice.ErrTOTPRateLimitExceeded:
			apperrors.Abort(c, apperrors.TooManyRequests(60))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to verify 2FA"))
		}
		return
	}

	// Create audit log
	h.createAuditLog(c, user, models.AuditAction2FAEnabled)

	response.OK(c, TwoFactorVerifyResponse{
		Message:     "2FA enabled successfully. Please save your backup codes securely.",
		BackupCodes: backupCodes,
	})
}

// Validate2FA handles validation of 2FA code during login.
// @Summary Validate 2FA code
// @Description Validate TOTP code during login to complete authentication
// @Tags Auth - 2FA
// @Accept json
// @Produce json
// @Param request body TwoFactorValidateRequest true "Partial token and TOTP code"
// @Success 200 {object} response.Success{data=LoginResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 429 {object} apperrors.AppError
// @Router /api/v1/auth/2fa/validate [post]
func (h *TwoFactorHandler) Validate2FA(c *gin.Context) {
	var req TwoFactorValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")

	// Validate partial token and get token pair
	tokenPair, user, err := h.authService.ValidateTwoFactorLogin(c.Request.Context(), req.PartialToken, req.Code, ipAddress, userAgent)
	if err != nil {
		switch err {
		case authservice.ErrInvalidToken, authservice.ErrExpiredToken:
			apperrors.Abort(c, apperrors.Unauthorized("Invalid or expired partial token. Please login again."))
		case authservice.ErrTOTPInvalidCode:
			apperrors.Abort(c, apperrors.BadRequest("Invalid 2FA code"))
		case authservice.ErrTOTPRateLimitExceeded:
			apperrors.Abort(c, apperrors.TooManyRequests(60))
		case authservice.ErrTOTPNotEnabled:
			apperrors.Abort(c, apperrors.BadRequest("2FA is not enabled for this account"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to validate 2FA"))
		}
		return
	}

	// Create audit log
	h.createAuditLog(c, user, models.AuditAction2FAValidated)

	resp := LoginResponse{
		User:              userToDTO(user),
		AccessToken:       tokenPair.AccessToken,
		RefreshToken:      tokenPair.RefreshToken,
		ExpiresIn:         tokenPair.ExpiresIn,
		TwoFactorEnabled:  true,
		RequiresTwoFactor: false,
	}

	response.OK(c, resp)
}

// Disable2FA handles disabling of 2FA.
// @Summary Disable 2FA
// @Description Disable 2FA for the authenticated user (requires password)
// @Tags Auth - 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TwoFactorDisableRequest true "Password"
// @Success 200 {object} response.Success{data=MessageResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/auth/2fa/disable [post]
func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {
	var req TwoFactorDisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == authservice.ErrUserNotFound {
			apperrors.Abort(c, apperrors.NotFound("User not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get user"))
		return
	}

	ipAddress := net.ParseIP(c.ClientIP())

	// Disable 2FA
	err = h.totpService.DisableTOTP(c.Request.Context(), user, req.Password, authservice.NewPasswordService(), ipAddress)
	if err != nil {
		switch err {
		case authservice.ErrInvalidCredentials:
			apperrors.Abort(c, apperrors.Unauthorized("Invalid password"))
		case authservice.ErrTOTPNotEnabled:
			apperrors.Abort(c, apperrors.BadRequest("2FA is not enabled"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to disable 2FA"))
		}
		return
	}

	// Create audit log
	h.createAuditLog(c, user, models.AuditAction2FADisabled)

	response.OK(c, MessageResponse{Message: "2FA disabled successfully"})
}

// GetBackupCodes returns the backup codes status for the user.
// @Summary Get backup codes status
// @Description Get the count of remaining backup codes
// @Tags Auth - 2FA
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=BackupCodesResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/auth/2fa/backup-codes [get]
func (h *TwoFactorHandler) GetBackupCodes(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == authservice.ErrUserNotFound {
			apperrors.Abort(c, apperrors.NotFound("User not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get user"))
		return
	}

	// Check if 2FA is enabled
	if !user.TwoFactorEnabled {
		apperrors.Abort(c, apperrors.BadRequest("2FA is not enabled"))
		return
	}

	// Get backup codes count
	count, err := h.totpService.GetBackupCodesCount(c.Request.Context(), userID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get backup codes"))
		return
	}

	response.OK(c, BackupCodesResponse{
		BackupCodesCount: count,
	})
}

// RegenerateBackupCodes regenerates backup codes for the user.
// @Summary Regenerate backup codes
// @Description Generate new backup codes (requires TOTP code)
// @Tags Auth - 2FA
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RegenerateBackupCodesRequest true "TOTP code"
// @Success 200 {object} response.Success{data=BackupCodesResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 429 {object} apperrors.AppError
// @Router /api/v1/auth/2fa/regenerate-backup [post]
func (h *TwoFactorHandler) RegenerateBackupCodes(c *gin.Context) {
	var req RegenerateBackupCodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == authservice.ErrUserNotFound {
			apperrors.Abort(c, apperrors.NotFound("User not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get user"))
		return
	}

	ipAddress := net.ParseIP(c.ClientIP())

	// Regenerate backup codes
	backupCodes, err := h.totpService.RegenerateBackupCodes(c.Request.Context(), user, req.Code, ipAddress)
	if err != nil {
		switch err {
		case authservice.ErrTOTPNotEnabled:
			apperrors.Abort(c, apperrors.BadRequest("2FA is not enabled"))
		case authservice.ErrTOTPInvalidCode:
			apperrors.Abort(c, apperrors.BadRequest("Invalid 2FA code"))
		case authservice.ErrTOTPRateLimitExceeded:
			apperrors.Abort(c, apperrors.TooManyRequests(60))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to regenerate backup codes"))
		}
		return
	}

	// Create audit log
	h.createAuditLog(c, user, models.AuditActionBackupCodesRegen)

	response.OK(c, BackupCodesResponse{
		BackupCodes:      backupCodes,
		BackupCodesCount: int64(len(backupCodes)),
	})
}

// GetStatus returns the 2FA status for the user.
// @Summary Get 2FA status
// @Description Get the current 2FA status for the authenticated user
// @Tags Auth - 2FA
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=TwoFactorStatusResponse}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/auth/2fa/status [get]
func (h *TwoFactorHandler) GetStatus(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == authservice.ErrUserNotFound {
			apperrors.Abort(c, apperrors.NotFound("User not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get user"))
		return
	}

	var backupCodesCount int64
	if user.TwoFactorEnabled {
		count, err := h.totpService.GetBackupCodesCount(c.Request.Context(), userID)
		if err == nil {
			backupCodesCount = count
		}
	}

	response.OK(c, TwoFactorStatusResponse{
		Enabled:          user.TwoFactorEnabled,
		BackupCodesCount: backupCodesCount,
	})
}

// createAuditLog creates an audit log entry for 2FA actions.
func (h *TwoFactorHandler) createAuditLog(c *gin.Context, user *models.User, action models.AuditAction) {
	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")
	h.authService.CreateAuditLog(c.Request.Context(), user, action, ipAddress, userAgent)
}
