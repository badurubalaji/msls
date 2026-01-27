// Package profile provides HTTP handlers for profile management endpoints.
package profile

import (
	"net"

	"github.com/gin-gonic/gin"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/services/auth"
	profileservice "msls-backend/internal/services/profile"
)

// Handler handles profile HTTP requests.
type Handler struct {
	profileService *profileservice.ProfileService
}

// NewHandler creates a new profile Handler.
func NewHandler(profileService *profileservice.ProfileService) *Handler {
	return &Handler{
		profileService: profileService,
	}
}

// GetProfile returns the current user's profile.
// @Summary Get current user profile
// @Description Get the authenticated user's full profile information
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=ProfileResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	user, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case profileservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to get profile"))
		}
		return
	}

	response.OK(c, UserToProfileResponse(user))
}

// UpdateProfile updates the current user's profile.
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update data"
// @Success 200 {object} response.Success{data=ProfileResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	user, err := h.profileService.UpdateProfile(c.Request.Context(), userID, req.ToServiceRequest())
	if err != nil {
		switch err {
		case profileservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		case profileservice.ErrInvalidTimezone:
			apperrors.Abort(c, apperrors.BadRequest("Invalid timezone"))
		case profileservice.ErrInvalidLocale:
			apperrors.Abort(c, apperrors.BadRequest("Invalid locale"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update profile"))
		}
		return
	}

	response.OK(c, UserToProfileResponse(user))
}

// ChangePassword changes the current user's password.
// @Summary Change password
// @Description Change the authenticated user's password
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change data"
// @Success 200 {object} response.Success{data=MessageResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/profile/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")

	err := h.profileService.ChangePassword(c.Request.Context(), userID, req.ToServiceRequest(), ipAddress, userAgent)
	if err != nil {
		switch err {
		case profileservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		case profileservice.ErrInvalidCurrentPassword:
			apperrors.Abort(c, apperrors.BadRequest("Current password is incorrect"))
		case profileservice.ErrPasswordMismatch:
			apperrors.Abort(c, apperrors.BadRequest("New password and confirmation do not match"))
		case auth.ErrPasswordTooShort,
			auth.ErrPasswordTooLong,
			auth.ErrPasswordNoUppercase,
			auth.ErrPasswordNoLowercase,
			auth.ErrPasswordNoDigit,
			auth.ErrPasswordNoSpecial:
			apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to change password"))
		}
		return
	}

	response.OK(c, MessageResponse{Message: "Password changed successfully"})
}

// UploadAvatar handles avatar file upload.
// @Summary Upload avatar
// @Description Upload a new avatar image for the authenticated user
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "Avatar image (JPEG or PNG, max 2MB)"
// @Success 200 {object} response.Success{data=AvatarUploadResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/profile/avatar [post]
func (h *Handler) UploadAvatar(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Avatar file is required"))
		return
	}
	defer file.Close()

	avatarURL, err := h.profileService.UploadAvatar(c.Request.Context(), userID, file, header)
	if err != nil {
		switch err {
		case profileservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		case profileservice.ErrInvalidAvatarFormat:
			apperrors.Abort(c, apperrors.BadRequest("Invalid avatar format. Only JPEG and PNG are allowed."))
		case profileservice.ErrAvatarTooLarge:
			apperrors.Abort(c, apperrors.BadRequest("Avatar file too large. Maximum size is 2MB."))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to upload avatar"))
		}
		return
	}

	response.OK(c, AvatarUploadResponse{AvatarURL: avatarURL})
}

// GetPreferences returns the current user's notification preferences.
// @Summary Get notification preferences
// @Description Get the authenticated user's notification preferences
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=NotificationPreferencesResponse}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/profile/preferences [get]
func (h *Handler) GetPreferences(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	prefs, err := h.profileService.GetNotificationPreferences(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case profileservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to get preferences"))
		}
		return
	}

	response.OK(c, NotificationPreferencesResponse{
		Email: prefs.Email,
		Push:  prefs.Push,
		SMS:   prefs.SMS,
	})
}

// UpdatePreferences updates the current user's notification preferences.
// @Summary Update notification preferences
// @Description Update the authenticated user's notification preferences
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdatePreferencesRequest true "Preferences update data"
// @Success 200 {object} response.Success{data=NotificationPreferencesResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/profile/preferences [put]
func (h *Handler) UpdatePreferences(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	var req UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	prefs, err := h.profileService.UpdateNotificationPreferences(c.Request.Context(), userID, req.ToServiceRequest())
	if err != nil {
		switch err {
		case profileservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update preferences"))
		}
		return
	}

	response.OK(c, NotificationPreferencesResponse{
		Email: prefs.Email,
		Push:  prefs.Push,
		SMS:   prefs.SMS,
	})
}

// RequestAccountDeletion initiates an account deletion request.
// @Summary Request account deletion
// @Description Request deletion of the authenticated user's account
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=MessageResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/profile [delete]
func (h *Handler) RequestAccountDeletion(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	ipAddress := net.ParseIP(c.ClientIP())
	userAgent := c.GetHeader("User-Agent")

	err := h.profileService.RequestAccountDeletion(c.Request.Context(), userID, ipAddress, userAgent)
	if err != nil {
		switch err {
		case profileservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		case profileservice.ErrAccountDeletionAlreadyRequested:
			apperrors.Abort(c, apperrors.Conflict("Account deletion already requested"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to request account deletion"))
		}
		return
	}

	response.OK(c, MessageResponse{
		Message: "Account deletion requested. Your account will be deactivated and scheduled for deletion.",
	})
}

// GetUserPreferences returns the current user's extended preferences.
// @Summary Get extended preferences
// @Description Get the authenticated user's extended preferences by category
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Param category query string false "Preference category filter"
// @Success 200 {object} response.Success{data=UserPreferencesResponse}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/profile/preferences/extended [get]
func (h *Handler) GetUserPreferences(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	category := c.Query("category")
	prefs, err := h.profileService.GetUserPreferences(c.Request.Context(), userID, category)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get preferences"))
		return
	}

	// Convert to DTOs
	dtos := make([]UserPreferenceDTO, len(prefs))
	for i, p := range prefs {
		var value interface{}
		if err := p.Value.UnmarshalJSON(p.Value); err == nil {
			value = p.Value
		}
		dtos[i] = UserPreferenceDTO{
			Category: p.Category,
			Key:      p.Key,
			Value:    value,
		}
	}

	response.OK(c, UserPreferencesResponse{Preferences: dtos})
}

// SetUserPreference sets a user preference.
// @Summary Set extended preference
// @Description Set a preference for the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SetPreferenceRequest true "Preference data"
// @Success 200 {object} response.Success{data=MessageResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/profile/preferences/extended [post]
func (h *Handler) SetUserPreference(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	var req SetPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	err := h.profileService.SetUserPreference(c.Request.Context(), userID, req.Category, req.Key, req.Value)
	if err != nil {
		switch err {
		case profileservice.ErrUserNotFound:
			apperrors.Abort(c, apperrors.NotFound("User not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to set preference"))
		}
		return
	}

	response.OK(c, MessageResponse{Message: "Preference saved successfully"})
}

// DeleteUserPreference deletes a user preference.
// @Summary Delete extended preference
// @Description Delete a preference for the authenticated user
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Param category query string true "Preference category"
// @Param key query string true "Preference key"
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/profile/preferences/extended [delete]
func (h *Handler) DeleteUserPreference(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		apperrors.Abort(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	category := c.Query("category")
	key := c.Query("key")

	if category == "" || key == "" {
		apperrors.Abort(c, apperrors.BadRequest("Category and key are required"))
		return
	}

	err := h.profileService.DeleteUserPreference(c.Request.Context(), userID, category, key)
	if err != nil {
		switch err {
		case profileservice.ErrPreferenceNotFound:
			apperrors.Abort(c, apperrors.NotFound("Preference not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete preference"))
		}
		return
	}

	response.NoContent(c)
}
