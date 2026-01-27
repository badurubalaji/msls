// Package auth provides authentication services for the MSLS application.
package auth

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Note: All errors are defined in errors.go

// Account lockout configuration.
const (
	MaxFailedAttempts = 5
	LockoutDuration   = 30 * time.Minute
)

// Partial token configuration for 2FA.
const (
	PartialTokenTTL = 5 * time.Minute
)

// TokenPair represents an access and refresh token pair.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	ExpiresIn    int       `json:"expires_in"` // seconds
}

// RegisterRequest represents a user registration request.
type RegisterRequest struct {
	TenantID  uuid.UUID
	Email     string
	Password  string
	FirstName string
	LastName  string
	RoleIDs   []uuid.UUID
}

// LoginRequest represents a login request.
type LoginRequest struct {
	Email     string
	Password  string
	TenantID  uuid.UUID
	IPAddress net.IP
	UserAgent string
}

// LoginResult represents the result of a login attempt.
type LoginResult struct {
	TokenPair         *TokenPair
	User              *models.User
	RequiresTwoFactor bool
	PartialToken      string
}

// AuthService handles authentication operations.
type AuthService struct {
	db              *gorm.DB
	jwtService      *JWTService
	passwordService *PasswordService
	totpService     *TOTPService
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(db *gorm.DB, jwtService *JWTService) *AuthService {
	return &AuthService{
		db:              db,
		jwtService:      jwtService,
		passwordService: NewPasswordService(),
		totpService:     nil, // Set via SetTOTPService
	}
}

// SetTOTPService sets the TOTP service for 2FA operations.
func (s *AuthService) SetTOTPService(totpService *TOTPService) {
	s.totpService = totpService
}

// GetJWTService returns the JWT service.
func (s *AuthService) GetJWTService() *JWTService {
	return s.jwtService
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// Validate password
	if err := s.passwordService.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if tenant exists and is active
	var tenant models.Tenant
	if err := s.db.WithContext(ctx).First(&tenant, "id = ?", req.TenantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	if tenant.Status != models.StatusActive {
		return nil, ErrTenantInactive
	}

	// Check if email already exists for this tenant
	var existingUser models.User
	err := s.db.WithContext(ctx).Where("tenant_id = ? AND email = ?", req.TenantID, req.Email).First(&existingUser).Error
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Validate roles exist
	var roles []models.Role
	if err := s.db.WithContext(ctx).Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	if len(roles) != len(req.RoleIDs) {
		return nil, ErrRoleNotFound
	}

	// Hash password
	passwordHash, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		TenantID:     req.TenantID,
		Email:        &req.Email,
		PasswordHash: &passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Status:       models.StatusActive,
		Roles:        roles,
	}

	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a token pair or indicates 2FA is required.
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResult, error) {
	// Find tenant (tenants table has no RLS, so no context needed)
	var tenant models.Tenant
	if err := s.db.WithContext(ctx).First(&tenant, "id = ?", req.TenantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.recordLoginAttempt(ctx, nil, req.Email, req.IPAddress, req.UserAgent, false, models.LoginFailureTenantInactive)
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	if tenant.Status != models.StatusActive {
		s.recordLoginAttempt(ctx, nil, req.Email, req.IPAddress, req.UserAgent, false, models.LoginFailureTenantInactive)
		return nil, ErrTenantInactive
	}

	// Find user by email within tenant
	// Use a transaction with RLS bypass since user is not authenticated yet
	var user models.User
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Bypass RLS for login - user is not authenticated yet
		if err := tx.Exec("SET LOCAL app.bypass_rls = 'true'").Error; err != nil {
			return err
		}
		return tx.Preload("Roles.Permissions").
			Where("tenant_id = ? AND email = ?", req.TenantID, req.Email).
			First(&user).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.recordLoginAttempt(ctx, nil, req.Email, req.IPAddress, req.UserAgent, false, models.LoginFailureUserNotFound)
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if account is locked
	if user.IsLocked() {
		s.recordLoginAttempt(ctx, &user.ID, req.Email, req.IPAddress, req.UserAgent, false, models.LoginFailureAccountLocked)
		return nil, ErrAccountLocked
	}

	// Check if account is active
	if !user.IsActive() {
		s.recordLoginAttempt(ctx, &user.ID, req.Email, req.IPAddress, req.UserAgent, false, models.LoginFailureAccountInactive)
		return nil, ErrAccountInactive
	}

	// Verify password
	if user.PasswordHash == nil {
		s.recordLoginAttempt(ctx, &user.ID, req.Email, req.IPAddress, req.UserAgent, false, models.LoginFailureInvalidCredentials)
		return nil, ErrInvalidCredentials
	}
	if err := s.passwordService.VerifyPassword(req.Password, *user.PasswordHash); err != nil {
		// Increment failed attempts
		user.IncrementFailedAttempts()
		if user.FailedLoginAttempts >= MaxFailedAttempts {
			user.Lock(LockoutDuration)
			s.createAuditLog(ctx, &user, models.AuditActionAccountLocked, req.IPAddress, req.UserAgent)
		}
		// Save with RLS bypass
		_ = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			tx.Exec("SET LOCAL app.bypass_rls = 'true'")
			return tx.Save(&user).Error
		})
		s.recordLoginAttempt(ctx, &user.ID, req.Email, req.IPAddress, req.UserAgent, false, models.LoginFailureInvalidCredentials)
		return nil, ErrInvalidCredentials
	}

	// Reset failed attempts on successful password verification
	user.ResetFailedAttempts()
	_ = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Exec("SET LOCAL app.bypass_rls = 'true'")
		return tx.Save(&user).Error
	})

	// Check if 2FA is enabled
	if user.TwoFactorEnabled && user.TwoFactorSecret != nil && user.TOTPVerifiedAt != nil {
		// Generate a partial token for 2FA verification
		partialToken, err := s.generatePartialToken(&user)
		if err != nil {
			return nil, err
		}

		return &LoginResult{
			User:              &user,
			RequiresTwoFactor: true,
			PartialToken:      partialToken,
		}, nil
	}

	// No 2FA required - complete login
	user.RecordLogin()
	_ = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Exec("SET LOCAL app.bypass_rls = 'true'")
		return tx.Save(&user).Error
	})

	// Generate token pair
	tokenPair, err := s.generateTokenPair(ctx, &user)
	if err != nil {
		return nil, err
	}

	// Record successful login
	s.recordLoginAttempt(ctx, &user.ID, req.Email, req.IPAddress, req.UserAgent, true, "")
	s.createAuditLog(ctx, &user, models.AuditActionLogin, req.IPAddress, req.UserAgent)

	return &LoginResult{
		TokenPair:         tokenPair,
		User:              &user,
		RequiresTwoFactor: false,
	}, nil
}

// generatePartialToken generates a temporary token for 2FA verification.
func (s *AuthService) generatePartialToken(user *models.User) (string, error) {
	// Use JWT with short expiry for the partial token
	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	// Generate a token with limited claims (no permissions until 2FA is verified)
	token, _, err := s.jwtService.GenerateAccessToken(user.ID, user.TenantID, email, []string{"2fa:pending"})
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateTwoFactorLogin validates the 2FA code and completes login.
func (s *AuthService) ValidateTwoFactorLogin(ctx context.Context, partialToken, code string, ipAddress net.IP, userAgent string) (*TokenPair, *models.User, error) {
	// Validate partial token
	claims, err := s.jwtService.ValidateAccessToken(partialToken)
	if err != nil {
		return nil, nil, err
	}

	// Check if this is a 2FA pending token
	is2FAPending := false
	for _, perm := range claims.Permissions {
		if perm == "2fa:pending" {
			is2FAPending = true
			break
		}
	}
	if !is2FAPending {
		return nil, nil, ErrInvalidToken
	}

	// Get user
	user, err := s.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, nil, err
	}

	// Validate 2FA code using TOTP service
	if s.totpService == nil {
		return nil, nil, errors.New("TOTP service not configured")
	}

	if err := s.totpService.ValidateTOTPForLogin(ctx, user, code, ipAddress); err != nil {
		s.createAuditLog(ctx, user, models.AuditAction2FAFailed, ipAddress, userAgent)
		return nil, nil, err
	}

	// Complete login
	user.RecordLogin()
	s.db.WithContext(ctx).Save(user)

	// Generate full token pair
	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	// Record successful login
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	s.recordLoginAttempt(ctx, &user.ID, email, ipAddress, userAgent, true, "")
	s.createAuditLog(ctx, user, models.AuditActionLogin, ipAddress, userAgent)

	return tokenPair, user, nil
}

// RefreshToken refreshes an access token using a refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string, ipAddress net.IP, userAgent string) (*TokenPair, error) {
	// Hash the refresh token
	tokenHash := s.jwtService.HashRefreshToken(refreshToken)

	// Find the refresh token in the database
	var storedToken models.RefreshToken
	err := s.db.WithContext(ctx).
		Preload("User.Roles.Permissions").
		Where("token_hash = ?", tokenHash).
		First(&storedToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.createTokenRefreshAuditLog(ctx, nil, models.AuditActionTokenRefreshFailed, ipAddress, userAgent, "token_not_found")
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	// Check if token is revoked
	if storedToken.IsRevoked() {
		s.createTokenRefreshAuditLog(ctx, storedToken.User, models.AuditActionTokenRefreshFailed, ipAddress, userAgent, "token_revoked")
		return nil, ErrRefreshTokenRevoked
	}

	// Check if token is expired
	if storedToken.IsExpired() {
		s.createTokenRefreshAuditLog(ctx, storedToken.User, models.AuditActionTokenRefreshFailed, ipAddress, userAgent, "token_expired")
		return nil, ErrRefreshTokenExpired
	}

	// Revoke the old refresh token (rotation)
	storedToken.Revoke()
	s.db.WithContext(ctx).Save(&storedToken)

	// Create audit log for token revocation due to rotation
	s.createTokenRefreshAuditLog(ctx, storedToken.User, models.AuditActionTokenRevoked, ipAddress, userAgent, "rotation")

	// Generate new token pair
	tokenPair, err := s.generateTokenPair(ctx, storedToken.User)
	if err != nil {
		return nil, err
	}

	// Create audit log for successful refresh
	s.createTokenRefreshAuditLog(ctx, storedToken.User, models.AuditActionTokenRefresh, ipAddress, userAgent, "")

	return tokenPair, nil
}

// createTokenRefreshAuditLog creates an audit log entry for token refresh events.
func (s *AuthService) createTokenRefreshAuditLog(ctx context.Context, user *models.User, action models.AuditAction, ipAddress net.IP, userAgent string, reason string) {
	builder := models.NewAuditLog(action, "refresh_token").
		WithIPAddress(ipAddress).
		WithUserAgent(userAgent)

	if user != nil {
		builder.WithTenant(user.TenantID).WithUser(user.ID)
	}

	if reason != "" {
		builder.WithNewData(map[string]string{"reason": reason})
	}

	log := builder.Build()
	s.db.WithContext(ctx).Create(log)
}

// Logout invalidates a refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshToken string, user *models.User, ipAddress net.IP, userAgent string) error {
	// Hash the refresh token
	tokenHash := s.jwtService.HashRefreshToken(refreshToken)

	// Find and revoke the refresh token
	var storedToken models.RefreshToken
	err := s.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&storedToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRefreshTokenNotFound
		}
		return err
	}

	storedToken.Revoke()
	if err := s.db.WithContext(ctx).Save(&storedToken).Error; err != nil {
		return err
	}

	// Create audit log
	if user != nil {
		s.createAuditLog(ctx, user, models.AuditActionLogout, ipAddress, userAgent)
	}

	return nil
}

// VerifyEmail verifies a user's email using a verification token.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	// Hash the token
	tokenHash := s.jwtService.HashRefreshToken(token)

	// Find the verification token
	var verificationToken models.VerificationToken
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("token_hash = ? AND type = ?", tokenHash, models.VerificationTokenTypeEmailVerify).
		First(&verificationToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVerificationTokenNotFound
		}
		return err
	}

	// Check if token is already used
	if verificationToken.IsUsed() {
		return ErrVerificationTokenUsed
	}

	// Check if token is expired
	if verificationToken.IsExpired() {
		return ErrVerificationTokenExpired
	}

	// Mark token as used
	verificationToken.MarkUsed()
	if err := s.db.WithContext(ctx).Save(&verificationToken).Error; err != nil {
		return err
	}

	// Verify user email
	verificationToken.User.VerifyEmail()
	if err := s.db.WithContext(ctx).Save(verificationToken.User).Error; err != nil {
		return err
	}

	// Create audit log
	s.createAuditLog(ctx, verificationToken.User, models.AuditActionEmailVerified, nil, "")

	return nil
}

// RequestPasswordReset creates a password reset token.
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string, tenantID uuid.UUID) (string, error) {
	// Find user by email
	var user models.User
	err := s.db.WithContext(ctx).Where("tenant_id = ? AND email = ?", tenantID, email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal whether the email exists
			return "", nil
		}
		return "", err
	}

	// Generate reset token
	token, _, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	// Store hashed token
	verificationToken := &models.VerificationToken{
		UserID:    user.ID,
		TokenHash: s.jwtService.HashRefreshToken(token),
		Type:      models.VerificationTokenTypePasswordReset,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiry for password reset
	}

	if err := s.db.WithContext(ctx).Create(verificationToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword resets a user's password using a reset token.
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string, ipAddress net.IP, userAgent string) error {
	// Validate new password
	if err := s.passwordService.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash the token
	tokenHash := s.jwtService.HashRefreshToken(token)

	// Find the verification token
	var verificationToken models.VerificationToken
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("token_hash = ? AND type = ?", tokenHash, models.VerificationTokenTypePasswordReset).
		First(&verificationToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVerificationTokenNotFound
		}
		return err
	}

	// Check if token is already used
	if verificationToken.IsUsed() {
		return ErrVerificationTokenUsed
	}

	// Check if token is expired
	if verificationToken.IsExpired() {
		return ErrVerificationTokenExpired
	}

	// Hash new password
	passwordHash, err := s.passwordService.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user password
	verificationToken.User.PasswordHash = &passwordHash
	verificationToken.User.Unlock() // Unlock account if locked
	if err := s.db.WithContext(ctx).Save(verificationToken.User).Error; err != nil {
		return err
	}

	// Mark token as used
	verificationToken.MarkUsed()
	if err := s.db.WithContext(ctx).Save(&verificationToken).Error; err != nil {
		return err
	}

	// Revoke all existing refresh tokens for this user
	s.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", verificationToken.User.ID).
		Update("revoked_at", time.Now())

	// Create audit log
	s.createAuditLog(ctx, verificationToken.User, models.AuditActionPasswordReset, ipAddress, userAgent)

	return nil
}

// CreateEmailVerificationToken creates a new email verification token for a user.
func (s *AuthService) CreateEmailVerificationToken(ctx context.Context, user *models.User) (string, error) {
	// Generate token
	token, _, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	// Store hashed token
	verificationToken := &models.VerificationToken{
		UserID:    user.ID,
		TokenHash: s.jwtService.HashRefreshToken(token),
		Type:      models.VerificationTokenTypeEmailVerify,
		ExpiresAt: time.Now().Add(72 * time.Hour), // 72 hour expiry for email verification
	}

	if err := s.db.WithContext(ctx).Create(verificationToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID retrieves a user by ID with roles and permissions.
func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Preload("Tenant").
		First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// generateTokenPair generates a new access and refresh token pair.
func (s *AuthService) generateTokenPair(ctx context.Context, user *models.User) (*TokenPair, error) {
	// Get user permissions
	permissions := user.GetPermissions()

	// Get email (handle nil pointer)
	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	// Generate access token
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(user.ID, user.TenantID, email, permissions)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, refreshExpiresAt, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store refresh token hash in database
	storedToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: s.jwtService.HashRefreshToken(refreshToken),
		ExpiresAt: refreshExpiresAt,
	}
	if err := s.db.WithContext(ctx).Create(storedToken).Error; err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		ExpiresIn:    int(s.jwtService.GetAccessTTL().Seconds()),
	}, nil
}

// recordLoginAttempt records a login attempt.
func (s *AuthService) recordLoginAttempt(ctx context.Context, userID *uuid.UUID, email string, ipAddress net.IP, userAgent string, success bool, failureReason string) {
	attempt := &models.LoginAttempt{
		UserID:        userID,
		Email:         email,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Success:       success,
		FailureReason: failureReason,
	}
	s.db.WithContext(ctx).Create(attempt)
}

// createAuditLog creates an audit log entry (internal use).
func (s *AuthService) createAuditLog(ctx context.Context, user *models.User, action models.AuditAction, ipAddress net.IP, userAgent string) {
	log := models.NewAuditLog(action, "user").
		WithTenant(user.TenantID).
		WithUser(user.ID).
		WithEntity(user.ID).
		WithIPAddress(ipAddress).
		WithUserAgent(userAgent).
		Build()
	s.db.WithContext(ctx).Create(log)
}

// CreateAuditLog creates an audit log entry (exported for handlers).
func (s *AuthService) CreateAuditLog(ctx context.Context, user *models.User, action models.AuditAction, ipAddress net.IP, userAgent string) {
	s.createAuditLog(ctx, user, action, ipAddress, userAgent)
}
