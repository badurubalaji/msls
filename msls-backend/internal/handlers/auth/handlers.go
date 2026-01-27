// Package auth provides HTTP handlers for authentication endpoints.
package auth

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
	authservice "msls-backend/internal/services/auth"
)

// Handler handles authentication HTTP requests.
type Handler struct {
	authService *authservice.AuthService
}

// NewHandler creates a new auth Handler.
func NewHandler(authService *authservice.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

// Register handles user registration.
// @Summary Register a new user
// @Description Create a new user account (admin only)
// @Tags Auth
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} response.Success{data=UserDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Get tenant ID from middleware context
	tenantID := middleware.GetTenantID(c)
	if tenantID == "" {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid tenant ID"))
		return
	}

	// Convert role IDs to UUIDs
	roleIDs := make([]uuid.UUID, len(req.RoleIDs))
	for i, id := range req.RoleIDs {
		roleUUID, err := uuid.Parse(id)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid role ID: "+id))
			return
		}
		roleIDs[i] = roleUUID
	}

	// Create registration request
	registerReq := authservice.RegisterRequest{
		TenantID:  tenantUUID,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleIDs:   roleIDs,
	}

	user, err := h.authService.Register(c.Request.Context(), registerReq)
	if err != nil {
		switch err {
		case authservice.ErrEmailAlreadyExists:
			apperrors.Abort(c, apperrors.Conflict("Email already exists"))
		case authservice.ErrTenantNotFound:
			apperrors.Abort(c, apperrors.NotFound("Tenant not found"))
		case authservice.ErrTenantInactive:
			apperrors.Abort(c, apperrors.BadRequest("Tenant is not active"))
		case authservice.ErrRoleNotFound:
			apperrors.Abort(c, apperrors.NotFound("One or more roles not found"))
		case authservice.ErrPasswordTooShort,
			authservice.ErrPasswordTooLong,
			authservice.ErrPasswordNoUppercase,
			authservice.ErrPasswordNoLowercase,
			authservice.ErrPasswordNoDigit,
			authservice.ErrPasswordNoSpecial:
			apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to register user"))
		}
		return
	}

	response.Created(c, userToDTO(user))
}

// Login handles user login.
// @Summary Login
// @Description Authenticate with email and password. If 2FA is enabled, returns partial_token and requires_two_factor=true.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.Success{data=LoginResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 423 {object} apperrors.AppError
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	tenantUUID, err := uuid.Parse(req.TenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid tenant ID"))
		return
	}

	// Get client IP and user agent
	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")

	loginReq := authservice.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		TenantID:  tenantUUID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	result, err := h.authService.Login(c.Request.Context(), loginReq)
	if err != nil {
		switch err {
		case authservice.ErrInvalidCredentials:
			apperrors.Abort(c, apperrors.Unauthorized("Invalid email or password"))
		case authservice.ErrAccountLocked:
			apperrors.Abort(c, &apperrors.AppError{
				Type:   "https://httpstatuses.com/423",
				Title:  "Locked",
				Status: http.StatusLocked,
				Detail: "Account is locked due to too many failed login attempts. Please try again later.",
			})
		case authservice.ErrAccountInactive:
			apperrors.Abort(c, apperrors.Unauthorized("Account is not active"))
		case authservice.ErrEmailNotVerified:
			apperrors.Abort(c, apperrors.Unauthorized("Email is not verified"))
		case authservice.ErrTenantNotFound:
			apperrors.Abort(c, apperrors.Unauthorized("Invalid credentials"))
		case authservice.ErrTenantInactive:
			apperrors.Abort(c, apperrors.Unauthorized("Tenant is not active"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Login failed"))
		}
		return
	}

	// Check if 2FA is required
	if result.RequiresTwoFactor {
		resp := LoginResponse{
			User:              userToDTO(result.User),
			TwoFactorEnabled:  true,
			RequiresTwoFactor: true,
			PartialToken:      result.PartialToken,
		}
		response.OK(c, resp)
		return
	}

	// Normal login response
	resp := LoginResponse{
		User:              userToDTO(result.User),
		AccessToken:       result.TokenPair.AccessToken,
		RefreshToken:      result.TokenPair.RefreshToken,
		ExpiresIn:         result.TokenPair.ExpiresIn,
		TwoFactorEnabled:  result.User.TwoFactorEnabled,
		RequiresTwoFactor: false,
	}

	response.OK(c, resp)
}

// RefreshToken handles token refresh.
// @Summary Refresh token
// @Description Get a new access token using a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Success{data=RefreshTokenResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Get client IP and user agent for audit logging
	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")

	tokenPair, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken, ipAddress, userAgent)
	if err != nil {
		switch err {
		case authservice.ErrRefreshTokenNotFound:
			apperrors.Abort(c, apperrors.Unauthorized("Invalid refresh token"))
		case authservice.ErrRefreshTokenRevoked:
			apperrors.Abort(c, apperrors.Unauthorized("Refresh token has been revoked"))
		case authservice.ErrRefreshTokenExpired:
			apperrors.Abort(c, apperrors.Unauthorized("Refresh token has expired"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Token refresh failed"))
		}
		return
	}

	resp := RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}

	response.OK(c, resp)
}

// Logout handles user logout.
// @Summary Logout
// @Description Invalidate the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LogoutRequest true "Refresh token to invalidate"
// @Success 200 {object} response.Success{data=MessageResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Get current user ID from context
	userID, _ := middleware.GetCurrentUserID(c)
	var user *models.User
	if userID != uuid.Nil {
		user, _ = h.authService.GetUserByID(c.Request.Context(), userID)
	}

	// Get client IP and user agent
	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")

	err := h.authService.Logout(c.Request.Context(), req.RefreshToken, user, ipAddress, userAgent)
	if err != nil {
		switch err {
		case authservice.ErrRefreshTokenNotFound:
			apperrors.Abort(c, apperrors.BadRequest("Invalid refresh token"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Logout failed"))
		}
		return
	}

	response.OK(c, MessageResponse{Message: "Logged out successfully"})
}

// VerifyEmail handles email verification.
// @Summary Verify email
// @Description Verify user email with token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Verification token"
// @Success 200 {object} response.Success{data=MessageResponse}
// @Failure 400 {object} apperrors.AppError
// @Router /api/v1/auth/verify-email [post]
func (h *Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	err := h.authService.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		switch err {
		case authservice.ErrVerificationTokenNotFound:
			apperrors.Abort(c, apperrors.BadRequest("Invalid verification token"))
		case authservice.ErrVerificationTokenUsed:
			apperrors.Abort(c, apperrors.BadRequest("Verification token has already been used"))
		case authservice.ErrVerificationTokenExpired:
			apperrors.Abort(c, apperrors.BadRequest("Verification token has expired"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Email verification failed"))
		}
		return
	}

	response.OK(c, MessageResponse{Message: "Email verified successfully"})
}

// ForgotPassword handles password reset request.
// @Summary Forgot password
// @Description Request a password reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email and tenant ID"
// @Success 200 {object} response.Success{data=MessageResponse}
// @Failure 400 {object} apperrors.AppError
// @Router /api/v1/auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	tenantUUID, err := uuid.Parse(req.TenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid tenant ID"))
		return
	}

	// Request password reset (don't reveal if email exists)
	_, err = h.authService.RequestPasswordReset(c.Request.Context(), req.Email, tenantUUID)
	if err != nil {
		// Log error but don't reveal to client
		// In production, you would send an email with the reset token
	}

	// Always return success to prevent email enumeration
	response.OK(c, MessageResponse{
		Message: "If an account with that email exists, a password reset link has been sent.",
	})
}

// ResetPassword handles password reset.
// @Summary Reset password
// @Description Reset password with token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} response.Success{data=MessageResponse}
// @Failure 400 {object} apperrors.AppError
// @Router /api/v1/auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Get client IP and user agent
	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")

	err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword, ipAddress, userAgent)
	if err != nil {
		switch err {
		case authservice.ErrVerificationTokenNotFound:
			apperrors.Abort(c, apperrors.BadRequest("Invalid reset token"))
		case authservice.ErrVerificationTokenUsed:
			apperrors.Abort(c, apperrors.BadRequest("Reset token has already been used"))
		case authservice.ErrVerificationTokenExpired:
			apperrors.Abort(c, apperrors.BadRequest("Reset token has expired"))
		case authservice.ErrPasswordTooShort,
			authservice.ErrPasswordTooLong,
			authservice.ErrPasswordNoUppercase,
			authservice.ErrPasswordNoLowercase,
			authservice.ErrPasswordNoDigit,
			authservice.ErrPasswordNoSpecial:
			apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		default:
			apperrors.Abort(c, apperrors.InternalError("Password reset failed"))
		}
		return
	}

	response.OK(c, MessageResponse{Message: "Password reset successfully"})
}

// Me returns the current user's information.
// @Summary Get current user
// @Description Get the authenticated user's information
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=MeResponse}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/auth/me [get]
func (h *Handler) Me(c *gin.Context) {
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

	permissions := user.GetPermissions()

	resp := MeResponse{
		User:        userToDTO(user),
		Permissions: permissions,
	}

	response.OK(c, resp)
}

// userToDTO converts a User model to a UserDTO.
func userToDTO(user *models.User) UserDTO {
	dto := UserDTO{
		ID:               user.ID,
		TenantID:         user.TenantID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		FullName:         user.FullName(),
		Status:           string(user.Status),
		TwoFactorEnabled: user.TwoFactorEnabled,
		EmailVerifiedAt:  user.EmailVerifiedAt,
		PhoneVerifiedAt:  user.PhoneVerifiedAt,
		LastLoginAt:      user.LastLoginAt,
		CreatedAt:        user.CreatedAt,
		Permissions:      user.GetPermissions(),
	}

	if user.Email != nil {
		dto.Email = *user.Email
	}
	if user.Phone != nil {
		dto.Phone = *user.Phone
	}

	// Convert roles
	dto.Roles = make([]RoleDTO, len(user.Roles))
	for i, role := range user.Roles {
		dto.Roles[i] = RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsSystem:    role.IsSystem,
		}
	}

	return dto
}
