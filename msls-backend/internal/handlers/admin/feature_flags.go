// Package admin provides HTTP handlers for administrative endpoints.
package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
	"msls-backend/internal/services/featureflag"
)

// FeatureFlagHandler handles feature flag HTTP requests.
type FeatureFlagHandler struct {
	service *featureflag.Service
}

// NewFeatureFlagHandler creates a new FeatureFlagHandler.
func NewFeatureFlagHandler(service *featureflag.Service) *FeatureFlagHandler {
	return &FeatureFlagHandler{
		service: service,
	}
}

// ListFlags returns all feature flags.
// @Summary List feature flags
// @Description Get all feature flags (admin only)
// @Tags Admin - Feature Flags
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=ListFeatureFlagsResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admin/feature-flags [get]
func (h *FeatureFlagHandler) ListFlags(c *gin.Context) {
	flags, err := h.service.ListFlags(c.Request.Context())
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to list feature flags"))
		return
	}

	dtos := make([]FeatureFlagDTO, len(flags))
	for i, flag := range flags {
		dtos[i] = flagToDTO(&flag)
	}

	response.OK(c, ListFeatureFlagsResponse{Flags: dtos})
}

// GetFlag returns a specific feature flag.
// @Summary Get feature flag
// @Description Get a feature flag by ID (admin only)
// @Tags Admin - Feature Flags
// @Produce json
// @Security BearerAuth
// @Param id path string true "Feature Flag ID"
// @Success 200 {object} response.Success{data=FeatureFlagDTO}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admin/feature-flags/{id} [get]
func (h *FeatureFlagHandler) GetFlag(c *gin.Context) {
	flagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid flag ID"))
		return
	}

	flag, err := h.service.GetFlag(c.Request.Context(), flagID)
	if err != nil {
		if err == featureflag.ErrFlagNotFound {
			apperrors.Abort(c, apperrors.NotFound("Feature flag not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to get feature flag"))
		return
	}

	response.OK(c, flagToDTO(flag))
}

// CreateFlag creates a new feature flag.
// @Summary Create feature flag
// @Description Create a new feature flag (admin only)
// @Tags Admin - Feature Flags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateFeatureFlagRequest true "Feature flag details"
// @Success 201 {object} response.Success{data=FeatureFlagDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/admin/feature-flags [post]
func (h *FeatureFlagHandler) CreateFlag(c *gin.Context) {
	var req CreateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	createReq := featureflag.CreateFlagRequest{
		Key:          req.Key,
		Name:         req.Name,
		Description:  req.Description,
		DefaultValue: req.DefaultValue,
	}

	if req.Metadata != nil {
		createReq.Metadata = models.FeatureFlagMetadata{
			Category:          req.Metadata.Category,
			RequiresSetup:     req.Metadata.RequiresSetup,
			Beta:              req.Metadata.Beta,
			RolloutPercentage: req.Metadata.RolloutPercentage,
		}
	}

	flag, err := h.service.CreateFlag(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case featureflag.ErrFlagKeyExists:
			apperrors.Abort(c, apperrors.Conflict("Feature flag key already exists"))
		case featureflag.ErrInvalidFlagKey:
			apperrors.Abort(c, apperrors.BadRequest("Invalid flag key format. Use lowercase letters, numbers, and underscores only."))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create feature flag"))
		}
		return
	}

	response.Created(c, flagToDTO(flag))
}

// UpdateFlag updates an existing feature flag.
// @Summary Update feature flag
// @Description Update a feature flag (admin only)
// @Tags Admin - Feature Flags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Feature Flag ID"
// @Param request body UpdateFeatureFlagRequest true "Feature flag updates"
// @Success 200 {object} response.Success{data=FeatureFlagDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admin/feature-flags/{id} [put]
func (h *FeatureFlagHandler) UpdateFlag(c *gin.Context) {
	flagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid flag ID"))
		return
	}

	var req UpdateFeatureFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := featureflag.UpdateFlagRequest{
		Name:         req.Name,
		Description:  req.Description,
		DefaultValue: req.DefaultValue,
	}

	if req.Metadata != nil {
		updateReq.Metadata = &models.FeatureFlagMetadata{
			Category:          req.Metadata.Category,
			RequiresSetup:     req.Metadata.RequiresSetup,
			Beta:              req.Metadata.Beta,
			RolloutPercentage: req.Metadata.RolloutPercentage,
		}
	}

	flag, err := h.service.UpdateFlag(c.Request.Context(), flagID, updateReq)
	if err != nil {
		if err == featureflag.ErrFlagNotFound {
			apperrors.Abort(c, apperrors.NotFound("Feature flag not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update feature flag"))
		return
	}

	response.OK(c, flagToDTO(flag))
}

// DeleteFlag deletes a feature flag.
// @Summary Delete feature flag
// @Description Delete a feature flag (admin only)
// @Tags Admin - Feature Flags
// @Produce json
// @Security BearerAuth
// @Param id path string true "Feature Flag ID"
// @Success 204 "No Content"
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admin/feature-flags/{id} [delete]
func (h *FeatureFlagHandler) DeleteFlag(c *gin.Context) {
	flagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid flag ID"))
		return
	}

	if err := h.service.DeleteFlag(c.Request.Context(), flagID); err != nil {
		if err == featureflag.ErrFlagNotFound {
			apperrors.Abort(c, apperrors.NotFound("Feature flag not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete feature flag"))
		return
	}

	response.NoContent(c)
}

// SetTenantFlags sets feature flag overrides for a tenant.
// @Summary Set tenant feature flags
// @Description Set feature flag overrides for a specific tenant (admin only)
// @Tags Admin - Feature Flags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Param request body SetTenantFlagsRequest true "Flag overrides"
// @Success 200 {object} response.Success{data=[]TenantFeatureFlagDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admin/tenants/{id}/feature-flags [put]
func (h *FeatureFlagHandler) SetTenantFlags(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid tenant ID"))
		return
	}

	var req SetTenantFlagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Set each flag override
	for _, override := range req.Flags {
		setReq := featureflag.SetForTenantRequest{
			TenantID:    tenantID,
			FlagKey:     override.Key,
			Enabled:     override.Enabled,
			CustomValue: override.CustomValue,
		}
		if err := h.service.SetForTenant(c.Request.Context(), setReq); err != nil {
			if err == featureflag.ErrFlagNotFound {
				apperrors.Abort(c, apperrors.NotFound("Feature flag not found: "+override.Key))
				return
			}
			apperrors.Abort(c, apperrors.InternalError("Failed to set tenant feature flag"))
			return
		}
	}

	// Return updated overrides
	overrides, err := h.service.GetTenantOverrides(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get tenant feature flags"))
		return
	}

	dtos := make([]TenantFeatureFlagDTO, len(overrides))
	for i, override := range overrides {
		dtos[i] = TenantFeatureFlagDTO{
			ID:          override.ID,
			TenantID:    override.TenantID,
			FlagKey:     override.FeatureFlag.Key,
			FlagName:    override.FeatureFlag.Name,
			Enabled:     override.Enabled,
			CustomValue: override.CustomValue,
			CreatedAt:   override.CreatedAt,
			UpdatedAt:   override.UpdatedAt,
		}
	}

	response.OK(c, dtos)
}

// GetTenantFlags returns feature flag overrides for a tenant.
// @Summary Get tenant feature flags
// @Description Get all feature flag overrides for a specific tenant (admin only)
// @Tags Admin - Feature Flags
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Success 200 {object} response.Success{data=[]TenantFeatureFlagDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admin/tenants/{id}/feature-flags [get]
func (h *FeatureFlagHandler) GetTenantFlags(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid tenant ID"))
		return
	}

	overrides, err := h.service.GetTenantOverrides(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get tenant feature flags"))
		return
	}

	dtos := make([]TenantFeatureFlagDTO, len(overrides))
	for i, override := range overrides {
		dtos[i] = TenantFeatureFlagDTO{
			ID:          override.ID,
			TenantID:    override.TenantID,
			FlagKey:     override.FeatureFlag.Key,
			FlagName:    override.FeatureFlag.Name,
			Enabled:     override.Enabled,
			CustomValue: override.CustomValue,
			CreatedAt:   override.CreatedAt,
			UpdatedAt:   override.UpdatedAt,
		}
	}

	response.OK(c, dtos)
}

// SetUserFlags sets feature flag overrides for a user.
// @Summary Set user feature flags
// @Description Set feature flag overrides for a specific user for beta testing (admin only)
// @Tags Admin - Feature Flags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body SetUserFlagsRequest true "Flag overrides"
// @Success 200 {object} response.Success{data=[]UserFeatureFlagDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/admin/users/{id}/feature-flags [put]
func (h *FeatureFlagHandler) SetUserFlags(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid user ID"))
		return
	}

	var req SetUserFlagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Set each flag override
	for _, override := range req.Flags {
		setReq := featureflag.SetForUserRequest{
			UserID:  userID,
			FlagKey: override.Key,
			Enabled: override.Enabled,
		}
		if err := h.service.SetForUser(c.Request.Context(), setReq); err != nil {
			if err == featureflag.ErrFlagNotFound {
				apperrors.Abort(c, apperrors.NotFound("Feature flag not found: "+override.Key))
				return
			}
			apperrors.Abort(c, apperrors.InternalError("Failed to set user feature flag"))
			return
		}
	}

	// Return updated overrides
	overrides, err := h.service.GetUserOverrides(c.Request.Context(), userID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get user feature flags"))
		return
	}

	dtos := make([]UserFeatureFlagDTO, len(overrides))
	for i, override := range overrides {
		dtos[i] = UserFeatureFlagDTO{
			ID:        override.ID,
			UserID:    override.UserID,
			FlagKey:   override.FeatureFlag.Key,
			FlagName:  override.FeatureFlag.Name,
			Enabled:   override.Enabled,
			CreatedAt: override.CreatedAt,
			UpdatedAt: override.UpdatedAt,
		}
	}

	response.OK(c, dtos)
}

// GetUserFlags returns feature flag overrides for a user.
// @Summary Get user feature flags
// @Description Get all feature flag overrides for a specific user (admin only)
// @Tags Admin - Feature Flags
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Success{data=[]UserFeatureFlagDTO}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/admin/users/{id}/feature-flags [get]
func (h *FeatureFlagHandler) GetUserFlags(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid user ID"))
		return
	}

	overrides, err := h.service.GetUserOverrides(c.Request.Context(), userID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to get user feature flags"))
		return
	}

	dtos := make([]UserFeatureFlagDTO, len(overrides))
	for i, override := range overrides {
		dtos[i] = UserFeatureFlagDTO{
			ID:        override.ID,
			UserID:    override.UserID,
			FlagKey:   override.FeatureFlag.Key,
			FlagName:  override.FeatureFlag.Name,
			Enabled:   override.Enabled,
			CreatedAt: override.CreatedAt,
			UpdatedAt: override.UpdatedAt,
		}
	}

	response.OK(c, dtos)
}

// GetCurrentFlags returns the current user's active feature flags.
// @Summary Get current user's feature flags
// @Description Get all feature flags with their current state for the authenticated user
// @Tags Feature Flags
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Success{data=CurrentFlagsResponse}
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/feature-flags [get]
func (h *FeatureFlagHandler) GetCurrentFlags(c *gin.Context) {
	// Get user and tenant IDs from context
	userID, _ := middleware.GetCurrentUserID(c)
	tenantID, _ := middleware.GetCurrentTenantID(c)

	states := h.service.GetAllFlagsForContext(c.Request.Context(), tenantID, userID)

	dtos := make([]FeatureFlagStateDTO, len(states))
	for i, state := range states {
		dtos[i] = FeatureFlagStateDTO{
			Key:         state.Key,
			Name:        state.Name,
			Description: state.Description,
			Enabled:     state.Enabled,
			CustomValue: state.CustomValue,
			Source:      state.Source,
		}
	}

	response.OK(c, CurrentFlagsResponse{Flags: dtos})
}

// IsEnabled checks if a specific feature flag is enabled for the current user.
// @Summary Check if feature flag is enabled
// @Description Check if a specific feature flag is enabled for the authenticated user
// @Tags Feature Flags
// @Produce json
// @Security BearerAuth
// @Param key path string true "Feature Flag Key"
// @Success 200 {object} response.Success{data=FeatureFlagStateDTO}
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/feature-flags/{key} [get]
func (h *FeatureFlagHandler) IsEnabled(c *gin.Context) {
	flagKey := c.Param("key")
	if flagKey == "" {
		apperrors.Abort(c, apperrors.BadRequest("Flag key is required"))
		return
	}

	// Get user and tenant IDs from context
	userID, _ := middleware.GetCurrentUserID(c)
	tenantID, _ := middleware.GetCurrentTenantID(c)

	state := h.service.GetFlagState(c.Request.Context(), flagKey, tenantID, userID)
	if state == nil {
		apperrors.Abort(c, apperrors.NotFound("Feature flag not found"))
		return
	}

	dto := FeatureFlagStateDTO{
		Key:         state.Key,
		Name:        state.Name,
		Description: state.Description,
		Enabled:     state.Enabled,
		CustomValue: state.CustomValue,
		Source:      state.Source,
	}

	response.OK(c, dto)
}

// Helper function to convert model to DTO.
func flagToDTO(flag *models.FeatureFlag) FeatureFlagDTO {
	return FeatureFlagDTO{
		ID:           flag.ID,
		Key:          flag.Key,
		Name:         flag.Name,
		Description:  flag.Description,
		DefaultValue: flag.DefaultValue,
		Metadata: FeatureFlagMetadataDTO{
			Category:          flag.Metadata.Category,
			RequiresSetup:     flag.Metadata.RequiresSetup,
			Beta:              flag.Metadata.Beta,
			RolloutPercentage: flag.Metadata.RolloutPercentage,
		},
		CreatedAt: flag.CreatedAt,
		UpdatedAt: flag.UpdatedAt,
	}
}
