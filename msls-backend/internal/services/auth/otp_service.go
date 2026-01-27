// Package auth provides authentication services for the MSLS application.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
	"msls-backend/internal/pkg/sms"
)

// OTPService handles OTP-based authentication operations.
type OTPService struct {
	db          *gorm.DB
	jwtService  *JWTService
	smsProvider sms.Provider
}

// OTPConfig holds OTP service configuration.
type OTPConfig struct {
	SMSProvider sms.Provider
}

// NewOTPService creates a new OTPService instance.
func NewOTPService(db *gorm.DB, jwtService *JWTService, config OTPConfig) *OTPService {
	return &OTPService{
		db:          db,
		jwtService:  jwtService,
		smsProvider: config.SMSProvider,
	}
}

// RequestOTPRequest represents a request to send an OTP.
type RequestOTPRequest struct {
	Identifier string           // Phone number or email
	Type       models.OTPType   // Type of OTP (login, verify)
	Channel    models.OTPChannel // Delivery channel (sms, email)
	TenantID   uuid.UUID        // Optional tenant ID for login OTPs
}

// RequestOTPResponse represents the response after sending an OTP.
type RequestOTPResponse struct {
	Message          string `json:"message"`
	ExpiresIn        int    `json:"expires_in"` // seconds
	MaskedIdentifier string `json:"masked_identifier"`
}

// VerifyOTPRequest represents a request to verify an OTP.
type VerifyOTPRequest struct {
	Identifier string
	Code       string
	TenantID   uuid.UUID
	IPAddress  net.IP
	UserAgent  string
}

// RequestOTP generates and sends an OTP to the specified identifier.
func (s *OTPService) RequestOTP(ctx context.Context, req RequestOTPRequest) (*RequestOTPResponse, error) {
	// Validate identifier
	if err := s.validateIdentifier(req.Identifier, req.Channel); err != nil {
		return nil, err
	}

	// Normalize identifier
	identifier := s.normalizeIdentifier(req.Identifier, req.Channel)

	// Check rate limit
	if err := s.checkRateLimit(ctx, identifier, req.Channel); err != nil {
		return nil, err
	}

	// Find user if this is a login OTP
	var user *models.User
	if req.Type == models.OTPTypeLogin {
		var err error
		user, err = s.findUserByIdentifier(ctx, identifier, req.Channel, req.TenantID)
		if err != nil {
			return nil, err
		}
	}

	// Generate OTP code
	code, err := s.generateOTPCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Hash the OTP code
	codeHash := s.hashOTPCode(code)

	// Create OTP record
	otp := &models.OTPCode{
		Identifier: identifier,
		CodeHash:   codeHash,
		Type:       req.Type,
		Channel:    req.Channel,
		ExpiresAt:  time.Now().Add(models.OTPExpiryDuration),
	}

	if user != nil {
		otp.UserID = &user.ID
	}

	if err := s.db.WithContext(ctx).Create(otp).Error; err != nil {
		return nil, fmt.Errorf("failed to create OTP record: %w", err)
	}

	// Update rate limit
	if err := s.updateRateLimit(ctx, identifier, req.Channel); err != nil {
		// Log error but don't fail the request
		// Rate limit tracking is not critical for sending
	}

	// Send OTP
	if err := s.sendOTP(ctx, identifier, code, req.Channel); err != nil {
		return nil, err
	}

	return &RequestOTPResponse{
		Message:          "OTP sent",
		ExpiresIn:        int(models.OTPExpiryDuration.Seconds()),
		MaskedIdentifier: s.maskIdentifier(identifier, req.Channel),
	}, nil
}

// VerifyOTP verifies an OTP code and returns a token pair if successful.
func (s *OTPService) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*TokenPair, *models.User, error) {
	// Validate identifier
	channel := s.detectChannel(req.Identifier)
	if err := s.validateIdentifier(req.Identifier, channel); err != nil {
		return nil, nil, err
	}

	// Normalize identifier
	identifier := s.normalizeIdentifier(req.Identifier, channel)

	// Hash the provided code
	codeHash := s.hashOTPCode(req.Code)

	// Find the most recent valid OTP for this identifier
	var otp models.OTPCode
	err := s.db.WithContext(ctx).
		Where("identifier = ? AND type = ? AND verified_at IS NULL AND expires_at > ?",
			identifier, models.OTPTypeLogin, time.Now()).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrOTPExpired
		}
		return nil, nil, fmt.Errorf("failed to find OTP: %w", err)
	}

	// Check if max attempts exceeded
	if otp.HasExceededMaxAttempts() {
		return nil, nil, ErrOTPMaxAttempts
	}

	// Increment attempts
	otp.IncrementAttempts()
	if err := s.db.WithContext(ctx).Save(&otp).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to update OTP attempts: %w", err)
	}

	// Verify the code
	if otp.CodeHash != codeHash {
		if otp.HasExceededMaxAttempts() {
			return nil, nil, ErrOTPMaxAttempts
		}
		return nil, nil, ErrOTPInvalid
	}

	// Mark OTP as verified
	otp.MarkVerified()
	if err := s.db.WithContext(ctx).Save(&otp).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to mark OTP as verified: %w", err)
	}

	// Find or create user
	user, err := s.findUserByIdentifier(ctx, identifier, channel, req.TenantID)
	if err != nil {
		return nil, nil, err
	}

	// Check if account is locked
	if user.IsLocked() {
		s.recordLoginAttempt(ctx, &user.ID, identifier, req.IPAddress, req.UserAgent, false, models.LoginFailureAccountLocked)
		return nil, nil, ErrAccountLocked
	}

	// Check if account is active
	if !user.IsActive() {
		s.recordLoginAttempt(ctx, &user.ID, identifier, req.IPAddress, req.UserAgent, false, models.LoginFailureAccountInactive)
		return nil, nil, ErrAccountInactive
	}

	// Verify phone if this was an SMS OTP
	if channel == models.OTPChannelSMS && !user.IsPhoneVerified() {
		user.VerifyPhone()
		if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to verify phone: %w", err)
		}
	}

	// Record login
	user.RecordLogin()
	user.ResetFailedAttempts()
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Load user with roles and permissions
	if err := s.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(user, "id = ?", user.ID).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to load user roles: %w", err)
	}

	// Generate token pair
	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	// Record successful login
	s.recordLoginAttempt(ctx, &user.ID, identifier, req.IPAddress, req.UserAgent, true, "")
	s.createAuditLog(ctx, user, models.AuditActionLogin, req.IPAddress, req.UserAgent)

	return tokenPair, user, nil
}

// ResendOTP resends an OTP to the specified identifier.
func (s *OTPService) ResendOTP(ctx context.Context, req RequestOTPRequest) (*RequestOTPResponse, error) {
	// Validate identifier
	if err := s.validateIdentifier(req.Identifier, req.Channel); err != nil {
		return nil, err
	}

	// Normalize identifier
	identifier := s.normalizeIdentifier(req.Identifier, req.Channel)

	// Check cooldown
	var rateLimit models.OTPRateLimit
	err := s.db.WithContext(ctx).
		Where("identifier = ? AND channel = ?", identifier, req.Channel).
		First(&rateLimit).Error

	if err == nil && !rateLimit.CanRequestNewOTP(models.OTPCooldownDuration) {
		return nil, ErrOTPCooldown
	}

	// Invalidate any existing unused OTPs
	s.db.WithContext(ctx).
		Model(&models.OTPCode{}).
		Where("identifier = ? AND channel = ? AND verified_at IS NULL", identifier, req.Channel).
		Update("expires_at", time.Now())

	// Request a new OTP
	return s.RequestOTP(ctx, req)
}

// CleanupExpiredOTPs removes expired OTP codes from the database.
func (s *OTPService) CleanupExpiredOTPs(ctx context.Context) (int64, error) {
	result := s.db.WithContext(ctx).
		Where("expires_at < ? OR verified_at IS NOT NULL", time.Now().Add(-time.Hour)).
		Delete(&models.OTPCode{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup expired OTPs: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// CleanupOldRateLimits removes old rate limit records.
func (s *OTPService) CleanupOldRateLimits(ctx context.Context) (int64, error) {
	result := s.db.WithContext(ctx).
		Where("window_start < ?", time.Now().Add(-2*time.Hour)).
		Delete(&models.OTPRateLimit{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup old rate limits: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// generateOTPCode generates a 6-digit OTP code.
func (s *OTPService) generateOTPCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	// Pad to 6 digits
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// hashOTPCode hashes an OTP code using SHA-256.
func (s *OTPService) hashOTPCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

// validateIdentifier validates a phone number or email address.
func (s *OTPService) validateIdentifier(identifier string, channel models.OTPChannel) error {
	switch channel {
	case models.OTPChannelSMS:
		if !s.isValidPhoneNumber(identifier) {
			return ErrInvalidIdentifier
		}
	case models.OTPChannelEmail:
		if !s.isValidEmail(identifier) {
			return ErrInvalidIdentifier
		}
	default:
		return ErrInvalidOTPChannel
	}
	return nil
}

// isValidPhoneNumber validates a phone number (E.164 format).
func (s *OTPService) isValidPhoneNumber(phone string) bool {
	// E.164 format: +[country code][subscriber number]
	// Examples: +919876543210, +14155552671
	pattern := `^\+[1-9]\d{6,14}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// isValidEmail validates an email address.
func (s *OTPService) isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// normalizeIdentifier normalizes a phone number or email.
func (s *OTPService) normalizeIdentifier(identifier string, channel models.OTPChannel) string {
	switch channel {
	case models.OTPChannelSMS:
		// Remove spaces and dashes
		return strings.ReplaceAll(strings.ReplaceAll(identifier, " ", ""), "-", "")
	case models.OTPChannelEmail:
		return strings.ToLower(strings.TrimSpace(identifier))
	}
	return identifier
}

// detectChannel detects the OTP channel from the identifier.
func (s *OTPService) detectChannel(identifier string) models.OTPChannel {
	if strings.HasPrefix(identifier, "+") {
		return models.OTPChannelSMS
	}
	return models.OTPChannelEmail
}

// maskIdentifier masks a phone number or email for display.
func (s *OTPService) maskIdentifier(identifier string, channel models.OTPChannel) string {
	switch channel {
	case models.OTPChannelSMS:
		if len(identifier) > 4 {
			return "****" + identifier[len(identifier)-4:]
		}
		return "****"
	case models.OTPChannelEmail:
		parts := strings.Split(identifier, "@")
		if len(parts) != 2 {
			return "****"
		}
		local := parts[0]
		if len(local) > 2 {
			return local[:2] + "****@" + parts[1]
		}
		return "****@" + parts[1]
	}
	return "****"
}

// checkRateLimit checks if the identifier has exceeded the rate limit.
func (s *OTPService) checkRateLimit(ctx context.Context, identifier string, channel models.OTPChannel) error {
	var rateLimit models.OTPRateLimit
	err := s.db.WithContext(ctx).
		Where("identifier = ? AND channel = ?", identifier, channel).
		First(&rateLimit).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No rate limit record, allow the request
			return nil
		}
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Check if rate limit exceeded
	if rateLimit.HasExceededLimit(models.MaxOTPRequestsPerHr, models.OTPRateLimitWindow) {
		return ErrOTPRateLimited
	}

	// Check cooldown
	if !rateLimit.CanRequestNewOTP(models.OTPCooldownDuration) {
		return ErrOTPCooldown
	}

	return nil
}

// updateRateLimit updates the rate limit record for an identifier.
func (s *OTPService) updateRateLimit(ctx context.Context, identifier string, channel models.OTPChannel) error {
	var rateLimit models.OTPRateLimit
	err := s.db.WithContext(ctx).
		Where("identifier = ? AND channel = ?", identifier, channel).
		First(&rateLimit).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new rate limit record
			rateLimit = models.OTPRateLimit{
				Identifier:    identifier,
				Channel:       channel,
				RequestCount:  1,
				WindowStart:   time.Now(),
				LastRequestAt: time.Now(),
			}
			return s.db.WithContext(ctx).Create(&rateLimit).Error
		}
		return err
	}

	// Check if window has expired
	if !rateLimit.IsWithinWindow(models.OTPRateLimitWindow) {
		rateLimit.ResetWindow()
	} else {
		rateLimit.IncrementCount()
	}

	return s.db.WithContext(ctx).Save(&rateLimit).Error
}

// findUserByIdentifier finds a user by phone number or email.
func (s *OTPService) findUserByIdentifier(ctx context.Context, identifier string, channel models.OTPChannel, tenantID uuid.UUID) (*models.User, error) {
	var user models.User
	var query *gorm.DB

	switch channel {
	case models.OTPChannelSMS:
		query = s.db.WithContext(ctx).Where("phone = ?", identifier)
	case models.OTPChannelEmail:
		query = s.db.WithContext(ctx).Where("email = ?", identifier)
	default:
		return nil, ErrInvalidOTPChannel
	}

	// Add tenant filter if provided
	if tenantID != uuid.Nil {
		query = query.Where("tenant_id = ?", tenantID)
	}

	err := query.First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIdentifierNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// sendOTP sends the OTP via the appropriate channel.
func (s *OTPService) sendOTP(ctx context.Context, identifier, code string, channel models.OTPChannel) error {
	switch channel {
	case models.OTPChannelSMS:
		return s.sendSMS(ctx, identifier, code)
	case models.OTPChannelEmail:
		return s.sendEmail(ctx, identifier, code)
	default:
		return ErrInvalidOTPChannel
	}
}

// sendSMS sends an OTP via SMS.
func (s *OTPService) sendSMS(ctx context.Context, phone, code string) error {
	if s.smsProvider == nil || !s.smsProvider.IsReady() {
		return ErrSMSSendFailed
	}

	msg := sms.Message{
		To:   phone,
		Body: fmt.Sprintf("Your MSLS verification code is: %s. This code expires in 5 minutes.", code),
	}

	_, err := s.smsProvider.Send(ctx, msg)
	if err != nil {
		return ErrSMSSendFailed
	}

	return nil
}

// sendEmail sends an OTP via email.
// TODO: Integrate with actual email service
func (s *OTPService) sendEmail(ctx context.Context, email, code string) error {
	// For now, just log the email (mock implementation)
	// In production, integrate with email service like SendGrid, SES, etc.
	fmt.Printf("[EMAIL MOCK] OTP to %s: %s\n", email, code)
	return nil
}

// generateTokenPair generates a new access and refresh token pair.
func (s *OTPService) generateTokenPair(ctx context.Context, user *models.User) (*TokenPair, error) {
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
func (s *OTPService) recordLoginAttempt(ctx context.Context, userID *uuid.UUID, identifier string, ipAddress net.IP, userAgent string, success bool, failureReason string) {
	attempt := &models.LoginAttempt{
		UserID:        userID,
		Email:         identifier,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Success:       success,
		FailureReason: failureReason,
	}
	s.db.WithContext(ctx).Create(attempt)
}

// createAuditLog creates an audit log entry.
func (s *OTPService) createAuditLog(ctx context.Context, user *models.User, action models.AuditAction, ipAddress net.IP, userAgent string) {
	log := models.NewAuditLog(action, "user").
		WithTenant(user.TenantID).
		WithUser(user.ID).
		WithEntity(user.ID).
		WithIPAddress(ipAddress).
		WithUserAgent(userAgent).
		Build()
	s.db.WithContext(ctx).Create(log)
}
