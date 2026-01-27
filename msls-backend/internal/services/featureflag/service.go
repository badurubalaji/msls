// Package featureflag provides feature flag management services.
package featureflag

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Flag key format validation.
var flagKeyRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Cache TTL for feature flags.
const (
	defaultCacheTTL = 5 * time.Minute
)

// cachedFlags holds the in-memory cache of feature flags.
type cachedFlags struct {
	flags      map[string]*models.FeatureFlag
	tenantMap  map[uuid.UUID]map[string]bool           // tenantID -> flagKey -> enabled
	tenantVals map[uuid.UUID]map[string]json.RawMessage // tenantID -> flagKey -> customValue
	userMap    map[uuid.UUID]map[string]bool           // userID -> flagKey -> enabled
	expiresAt  time.Time
	mu         sync.RWMutex
}

// Service provides feature flag management functionality.
type Service struct {
	db       *gorm.DB
	cache    *cachedFlags
	cacheTTL time.Duration
}

// Config holds configuration for the feature flag service.
type Config struct {
	// CacheTTL is the duration for which cached flags are valid.
	CacheTTL time.Duration
}

// DefaultConfig returns the default service configuration.
func DefaultConfig() Config {
	return Config{
		CacheTTL: defaultCacheTTL,
	}
}

// NewService creates a new feature flag service.
func NewService(db *gorm.DB, cfg Config) *Service {
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = defaultCacheTTL
	}

	return &Service{
		db:       db,
		cacheTTL: cfg.CacheTTL,
		cache: &cachedFlags{
			flags:      make(map[string]*models.FeatureFlag),
			tenantMap:  make(map[uuid.UUID]map[string]bool),
			tenantVals: make(map[uuid.UUID]map[string]json.RawMessage),
			userMap:    make(map[uuid.UUID]map[string]bool),
		},
	}
}

// IsEnabled checks if a feature flag is enabled for the given context.
// Priority: User override > Tenant override > Default value
func (s *Service) IsEnabled(ctx context.Context, flagKey string, tenantID, userID uuid.UUID) bool {
	s.ensureCacheValid(ctx)

	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	// Check user override first (highest priority)
	if userID != uuid.Nil {
		if userFlags, ok := s.cache.userMap[userID]; ok {
			if enabled, ok := userFlags[flagKey]; ok {
				return enabled
			}
		}
	}

	// Check tenant override
	if tenantID != uuid.Nil {
		if tenantFlags, ok := s.cache.tenantMap[tenantID]; ok {
			if enabled, ok := tenantFlags[flagKey]; ok {
				return enabled
			}
		}
	}

	// Fall back to default value
	if flag, ok := s.cache.flags[flagKey]; ok {
		return flag.DefaultValue
	}

	return false
}

// GetValue returns the custom value for a feature flag if set.
// Returns the tenant's custom value if set, otherwise nil.
func (s *Service) GetValue(ctx context.Context, flagKey string, tenantID uuid.UUID) json.RawMessage {
	s.ensureCacheValid(ctx)

	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	if tenantID != uuid.Nil {
		if tenantVals, ok := s.cache.tenantVals[tenantID]; ok {
			if val, ok := tenantVals[flagKey]; ok {
				return val
			}
		}
	}

	return nil
}

// GetFlagState returns the complete state of a feature flag for a given context.
func (s *Service) GetFlagState(ctx context.Context, flagKey string, tenantID, userID uuid.UUID) *models.FeatureFlagState {
	s.ensureCacheValid(ctx)

	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	flag, ok := s.cache.flags[flagKey]
	if !ok {
		return nil
	}

	state := &models.FeatureFlagState{
		Key:         flag.Key,
		Name:        flag.Name,
		Description: flag.Description,
		Enabled:     flag.DefaultValue,
		Source:      "default",
	}

	// Check tenant override
	if tenantID != uuid.Nil {
		if tenantFlags, ok := s.cache.tenantMap[tenantID]; ok {
			if enabled, ok := tenantFlags[flagKey]; ok {
				state.Enabled = enabled
				state.Source = "tenant"
			}
		}
		if tenantVals, ok := s.cache.tenantVals[tenantID]; ok {
			if val, ok := tenantVals[flagKey]; ok {
				state.CustomValue = val
			}
		}
	}

	// Check user override (highest priority)
	if userID != uuid.Nil {
		if userFlags, ok := s.cache.userMap[userID]; ok {
			if enabled, ok := userFlags[flagKey]; ok {
				state.Enabled = enabled
				state.Source = "user"
			}
		}
	}

	return state
}

// GetAllFlagsForContext returns all feature flags with their current state for a given context.
func (s *Service) GetAllFlagsForContext(ctx context.Context, tenantID, userID uuid.UUID) []models.FeatureFlagState {
	s.ensureCacheValid(ctx)

	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	states := make([]models.FeatureFlagState, 0, len(s.cache.flags))

	for _, flag := range s.cache.flags {
		state := models.FeatureFlagState{
			Key:         flag.Key,
			Name:        flag.Name,
			Description: flag.Description,
			Enabled:     flag.DefaultValue,
			Source:      "default",
		}

		// Check tenant override
		if tenantID != uuid.Nil {
			if tenantFlags, ok := s.cache.tenantMap[tenantID]; ok {
				if enabled, ok := tenantFlags[flag.Key]; ok {
					state.Enabled = enabled
					state.Source = "tenant"
				}
			}
			if tenantVals, ok := s.cache.tenantVals[tenantID]; ok {
				if val, ok := tenantVals[flag.Key]; ok {
					state.CustomValue = val
				}
			}
		}

		// Check user override (highest priority)
		if userID != uuid.Nil {
			if userFlags, ok := s.cache.userMap[userID]; ok {
				if enabled, ok := userFlags[flag.Key]; ok {
					state.Enabled = enabled
					state.Source = "user"
				}
			}
		}

		states = append(states, state)
	}

	return states
}

// ListFlags returns all feature flags (admin function).
func (s *Service) ListFlags(ctx context.Context) ([]models.FeatureFlag, error) {
	var flags []models.FeatureFlag
	if err := s.db.WithContext(ctx).Order("key").Find(&flags).Error; err != nil {
		return nil, err
	}
	return flags, nil
}

// GetFlag returns a feature flag by ID.
func (s *Service) GetFlag(ctx context.Context, flagID uuid.UUID) (*models.FeatureFlag, error) {
	var flag models.FeatureFlag
	if err := s.db.WithContext(ctx).First(&flag, "id = ?", flagID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFlagNotFound
		}
		return nil, err
	}
	return &flag, nil
}

// GetFlagByKey returns a feature flag by key.
func (s *Service) GetFlagByKey(ctx context.Context, key string) (*models.FeatureFlag, error) {
	var flag models.FeatureFlag
	if err := s.db.WithContext(ctx).First(&flag, "key = ?", key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFlagNotFound
		}
		return nil, err
	}
	return &flag, nil
}

// CreateFlagRequest represents a request to create a feature flag.
type CreateFlagRequest struct {
	Key          string
	Name         string
	Description  string
	DefaultValue bool
	Metadata     models.FeatureFlagMetadata
}

// CreateFlag creates a new feature flag.
func (s *Service) CreateFlag(ctx context.Context, req CreateFlagRequest) (*models.FeatureFlag, error) {
	// Validate key format
	if !flagKeyRegex.MatchString(req.Key) {
		return nil, ErrInvalidFlagKey
	}

	// Check if key already exists
	var existing models.FeatureFlag
	if err := s.db.WithContext(ctx).Where("key = ?", req.Key).First(&existing).Error; err == nil {
		return nil, ErrFlagKeyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	flag := &models.FeatureFlag{
		Key:          req.Key,
		Name:         req.Name,
		Description:  req.Description,
		DefaultValue: req.DefaultValue,
		Metadata:     req.Metadata,
	}

	if err := flag.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Create(flag).Error; err != nil {
		return nil, err
	}

	// Invalidate cache
	s.invalidateCache()

	return flag, nil
}

// UpdateFlagRequest represents a request to update a feature flag.
type UpdateFlagRequest struct {
	Name         *string
	Description  *string
	DefaultValue *bool
	Metadata     *models.FeatureFlagMetadata
}

// UpdateFlag updates a feature flag.
func (s *Service) UpdateFlag(ctx context.Context, flagID uuid.UUID, req UpdateFlagRequest) (*models.FeatureFlag, error) {
	flag, err := s.GetFlag(ctx, flagID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		flag.Name = *req.Name
	}
	if req.Description != nil {
		flag.Description = *req.Description
	}
	if req.DefaultValue != nil {
		flag.DefaultValue = *req.DefaultValue
	}
	if req.Metadata != nil {
		flag.Metadata = *req.Metadata
	}

	if err := flag.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Save(flag).Error; err != nil {
		return nil, err
	}

	// Invalidate cache
	s.invalidateCache()

	return flag, nil
}

// DeleteFlag deletes a feature flag.
func (s *Service) DeleteFlag(ctx context.Context, flagID uuid.UUID) error {
	flag, err := s.GetFlag(ctx, flagID)
	if err != nil {
		return err
	}

	// Delete the flag (cascade will remove tenant and user overrides)
	if err := s.db.WithContext(ctx).Delete(flag).Error; err != nil {
		return err
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// SetForTenantRequest represents a request to set a tenant-specific flag override.
type SetForTenantRequest struct {
	TenantID    uuid.UUID
	FlagKey     string
	Enabled     bool
	CustomValue json.RawMessage
}

// SetForTenant sets a feature flag override for a specific tenant.
func (s *Service) SetForTenant(ctx context.Context, req SetForTenantRequest) error {
	// Get the flag by key
	flag, err := s.GetFlagByKey(ctx, req.FlagKey)
	if err != nil {
		return err
	}

	// Check if override already exists
	var existing models.TenantFeatureFlag
	err = s.db.WithContext(ctx).Where("tenant_id = ? AND flag_id = ?", req.TenantID, flag.ID).First(&existing).Error

	if err == nil {
		// Update existing override
		existing.Enabled = req.Enabled
		existing.CustomValue = req.CustomValue
		if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
			return err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new override
		override := &models.TenantFeatureFlag{
			TenantID:    req.TenantID,
			FlagID:      flag.ID,
			Enabled:     req.Enabled,
			CustomValue: req.CustomValue,
		}
		if err := s.db.WithContext(ctx).Create(override).Error; err != nil {
			return err
		}
	} else {
		return err
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// RemoveForTenant removes a tenant-specific flag override.
func (s *Service) RemoveForTenant(ctx context.Context, tenantID uuid.UUID, flagKey string) error {
	flag, err := s.GetFlagByKey(ctx, flagKey)
	if err != nil {
		return err
	}

	result := s.db.WithContext(ctx).Where("tenant_id = ? AND flag_id = ?", tenantID, flag.ID).Delete(&models.TenantFeatureFlag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTenantOverrideNotFound
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// GetTenantOverrides returns all feature flag overrides for a tenant.
func (s *Service) GetTenantOverrides(ctx context.Context, tenantID uuid.UUID) ([]models.TenantFeatureFlag, error) {
	var overrides []models.TenantFeatureFlag
	if err := s.db.WithContext(ctx).
		Preload("FeatureFlag").
		Where("tenant_id = ?", tenantID).
		Find(&overrides).Error; err != nil {
		return nil, err
	}
	return overrides, nil
}

// SetForUserRequest represents a request to set a user-specific flag override.
type SetForUserRequest struct {
	UserID  uuid.UUID
	FlagKey string
	Enabled bool
}

// SetForUser sets a feature flag override for a specific user (beta testing).
func (s *Service) SetForUser(ctx context.Context, req SetForUserRequest) error {
	// Get the flag by key
	flag, err := s.GetFlagByKey(ctx, req.FlagKey)
	if err != nil {
		return err
	}

	// Check if override already exists
	var existing models.UserFeatureFlag
	err = s.db.WithContext(ctx).Where("user_id = ? AND flag_id = ?", req.UserID, flag.ID).First(&existing).Error

	if err == nil {
		// Update existing override
		existing.Enabled = req.Enabled
		if err := s.db.WithContext(ctx).Save(&existing).Error; err != nil {
			return err
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new override
		override := &models.UserFeatureFlag{
			UserID:  req.UserID,
			FlagID:  flag.ID,
			Enabled: req.Enabled,
		}
		if err := s.db.WithContext(ctx).Create(override).Error; err != nil {
			return err
		}
	} else {
		return err
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// RemoveForUser removes a user-specific flag override.
func (s *Service) RemoveForUser(ctx context.Context, userID uuid.UUID, flagKey string) error {
	flag, err := s.GetFlagByKey(ctx, flagKey)
	if err != nil {
		return err
	}

	result := s.db.WithContext(ctx).Where("user_id = ? AND flag_id = ?", userID, flag.ID).Delete(&models.UserFeatureFlag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserOverrideNotFound
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// GetUserOverrides returns all feature flag overrides for a user.
func (s *Service) GetUserOverrides(ctx context.Context, userID uuid.UUID) ([]models.UserFeatureFlag, error) {
	var overrides []models.UserFeatureFlag
	if err := s.db.WithContext(ctx).
		Preload("FeatureFlag").
		Where("user_id = ?", userID).
		Find(&overrides).Error; err != nil {
		return nil, err
	}
	return overrides, nil
}

// RefreshCache forces a refresh of the feature flag cache.
func (s *Service) RefreshCache(ctx context.Context) error {
	return s.loadCache(ctx)
}

// ensureCacheValid checks if the cache is valid and refreshes if needed.
func (s *Service) ensureCacheValid(ctx context.Context) {
	s.cache.mu.RLock()
	expired := time.Now().After(s.cache.expiresAt)
	s.cache.mu.RUnlock()

	if expired {
		s.loadCache(ctx)
	}
}

// invalidateCache marks the cache as expired.
func (s *Service) invalidateCache() {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()
	s.cache.expiresAt = time.Time{}
}

// loadCache loads all feature flags and overrides into memory.
func (s *Service) loadCache(ctx context.Context) error {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	// Load all flags
	var flags []models.FeatureFlag
	if err := s.db.WithContext(ctx).Find(&flags).Error; err != nil {
		return err
	}

	flagMap := make(map[string]*models.FeatureFlag)
	for i := range flags {
		flagMap[flags[i].Key] = &flags[i]
	}

	// Load all tenant overrides
	var tenantOverrides []models.TenantFeatureFlag
	if err := s.db.WithContext(ctx).Preload("FeatureFlag").Find(&tenantOverrides).Error; err != nil {
		return err
	}

	tenantMap := make(map[uuid.UUID]map[string]bool)
	tenantVals := make(map[uuid.UUID]map[string]json.RawMessage)
	for _, override := range tenantOverrides {
		if _, ok := tenantMap[override.TenantID]; !ok {
			tenantMap[override.TenantID] = make(map[string]bool)
			tenantVals[override.TenantID] = make(map[string]json.RawMessage)
		}
		tenantMap[override.TenantID][override.FeatureFlag.Key] = override.Enabled
		if override.CustomValue != nil {
			tenantVals[override.TenantID][override.FeatureFlag.Key] = override.CustomValue
		}
	}

	// Load all user overrides
	var userOverrides []models.UserFeatureFlag
	if err := s.db.WithContext(ctx).Preload("FeatureFlag").Find(&userOverrides).Error; err != nil {
		return err
	}

	userMap := make(map[uuid.UUID]map[string]bool)
	for _, override := range userOverrides {
		if _, ok := userMap[override.UserID]; !ok {
			userMap[override.UserID] = make(map[string]bool)
		}
		userMap[override.UserID][override.FeatureFlag.Key] = override.Enabled
	}

	// Update cache
	s.cache.flags = flagMap
	s.cache.tenantMap = tenantMap
	s.cache.tenantVals = tenantVals
	s.cache.userMap = userMap
	s.cache.expiresAt = time.Now().Add(s.cacheTTL)

	return nil
}
