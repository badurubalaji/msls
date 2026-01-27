// Package auth provides authentication services for the MSLS application.
package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// TOTP configuration constants.
const (
	// TOTPIssuer is the issuer name shown in authenticator apps.
	TOTPIssuer = "MSLS"

	// TOTPPeriod is the time step in seconds (standard is 30 seconds).
	TOTPPeriod = 30

	// TOTPDigits is the number of digits in the TOTP code.
	TOTPDigits = 6

	// TOTPAlgorithm is the hash algorithm used for TOTP.
	TOTPAlgorithm = otp.AlgorithmSHA1

	// BackupCodeCount is the number of backup codes generated.
	BackupCodeCount = 8

	// BackupCodeLength is the length of each backup code.
	BackupCodeLength = 8

	// TOTPRateLimitAttempts is the maximum number of 2FA attempts per minute.
	TOTPRateLimitAttempts = 5

	// TOTPRateLimitWindow is the time window for rate limiting.
	TOTPRateLimitWindow = time.Minute
)

// TOTPService handles TOTP-based two-factor authentication.
type TOTPService struct {
	db            *gorm.DB
	encryptionKey []byte
}

// TOTPSetupResult contains the data needed for 2FA setup.
type TOTPSetupResult struct {
	Secret        string `json:"secret"`
	QRCodeDataURL string `json:"qr_code_data_url"`
	ManualEntry   string `json:"manual_entry"`
}

// NewTOTPService creates a new TOTPService instance.
// encryptionKey must be 32 bytes for AES-256 encryption.
func NewTOTPService(db *gorm.DB, encryptionKey string) (*TOTPService, error) {
	key := deriveKey(encryptionKey)
	if len(key) != 32 {
		return nil, errors.New("encryption key must derive to 32 bytes")
	}
	return &TOTPService{
		db:            db,
		encryptionKey: key,
	}, nil
}

// deriveKey derives a 32-byte key from the input using SHA-256.
func deriveKey(input string) []byte {
	hash := sha256.Sum256([]byte(input))
	return hash[:]
}

// GenerateSecret generates a new TOTP secret for a user.
func (s *TOTPService) GenerateSecret(accountName string) (*TOTPSetupResult, error) {
	// Generate TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      TOTPIssuer,
		AccountName: accountName,
		Period:      TOTPPeriod,
		Digits:      otp.DigitsSix,
		Algorithm:   TOTPAlgorithm,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	return &TOTPSetupResult{
		Secret:        key.Secret(),
		QRCodeDataURL: key.URL(),
		ManualEntry:   key.Secret(),
	}, nil
}

// ValidateCode validates a TOTP code against a secret.
// It allows for a 30-second time window before and after the current time.
func (s *TOTPService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// ValidateCodeWithWindow validates a TOTP code with a specified time window.
func (s *TOTPService) ValidateCodeWithWindow(secret, code string, skew uint) bool {
	valid, _ := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    TOTPPeriod,
		Skew:      skew,
		Digits:    otp.DigitsSix,
		Algorithm: TOTPAlgorithm,
	})
	return valid
}

// EncryptSecret encrypts a TOTP secret for storage.
func (s *TOTPService) EncryptSecret(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSecret decrypts a stored TOTP secret.
func (s *TOTPService) DecryptSecret(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// GenerateBackupCodes generates a set of backup codes for account recovery.
func (s *TOTPService) GenerateBackupCodes() ([]string, error) {
	codes := make([]string, BackupCodeCount)
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for i := 0; i < BackupCodeCount; i++ {
		code := make([]byte, BackupCodeLength)
		randomBytes := make([]byte, BackupCodeLength)
		if _, err := rand.Read(randomBytes); err != nil {
			return nil, fmt.Errorf("failed to generate random bytes: %w", err)
		}

		for j := 0; j < BackupCodeLength; j++ {
			code[j] = charset[randomBytes[j]%byte(len(charset))]
		}
		codes[i] = string(code)
	}

	return codes, nil
}

// HashBackupCode hashes a backup code for storage using SHA-256.
func (s *TOTPService) HashBackupCode(code string) string {
	// Normalize the code (uppercase, remove spaces/dashes)
	normalized := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(code, " ", ""), "-", ""))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

// SetupTOTP initializes 2FA setup for a user without enabling it.
// The secret is stored encrypted but 2FA is not yet enabled.
func (s *TOTPService) SetupTOTP(ctx context.Context, user *models.User) (*TOTPSetupResult, error) {
	// Get account name from email
	accountName := ""
	if user.Email != nil {
		accountName = *user.Email
	} else {
		accountName = user.ID.String()
	}

	// Generate new secret
	result, err := s.GenerateSecret(accountName)
	if err != nil {
		return nil, err
	}

	// Encrypt the secret
	encryptedSecret, err := s.EncryptSecret(result.Secret)
	if err != nil {
		return nil, err
	}

	// Store the encrypted secret (but don't enable 2FA yet)
	user.TwoFactorSecret = &encryptedSecret
	user.TwoFactorEnabled = false
	user.TOTPVerifiedAt = nil

	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, fmt.Errorf("failed to save TOTP secret: %w", err)
	}

	return result, nil
}

// VerifyAndEnableTOTP verifies the TOTP code and enables 2FA for the user.
func (s *TOTPService) VerifyAndEnableTOTP(ctx context.Context, user *models.User, code string, ipAddress net.IP) ([]string, error) {
	// Check rate limit
	if exceeded, err := s.checkRateLimit(ctx, user.ID, ipAddress); err != nil {
		return nil, err
	} else if exceeded {
		return nil, ErrTOTPRateLimitExceeded
	}

	// Check if user has a pending TOTP secret
	if user.TwoFactorSecret == nil || *user.TwoFactorSecret == "" {
		s.recordTOTPAttempt(ctx, user.ID, ipAddress, false)
		return nil, ErrTOTPNotSetup
	}

	// Check if 2FA is already enabled
	if user.TwoFactorEnabled {
		return nil, ErrTOTPAlreadyEnabled
	}

	// Decrypt the secret
	secret, err := s.DecryptSecret(*user.TwoFactorSecret)
	if err != nil {
		s.recordTOTPAttempt(ctx, user.ID, ipAddress, false)
		return nil, fmt.Errorf("failed to decrypt TOTP secret: %w", err)
	}

	// Validate the code
	if !s.ValidateCode(secret, code) {
		s.recordTOTPAttempt(ctx, user.ID, ipAddress, false)
		return nil, ErrTOTPInvalidCode
	}

	// Generate backup codes
	backupCodes, err := s.GenerateBackupCodes()
	if err != nil {
		return nil, err
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete any existing backup codes
	if err := tx.Where("user_id = ?", user.ID).Delete(&models.BackupCode{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Store hashed backup codes
	for _, code := range backupCodes {
		backupCode := &models.BackupCode{
			UserID:   user.ID,
			CodeHash: s.HashBackupCode(code),
		}
		if err := tx.Create(backupCode).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Enable 2FA
	now := time.Now()
	user.TwoFactorEnabled = true
	user.TOTPVerifiedAt = &now

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Record successful attempt
	s.recordTOTPAttempt(ctx, user.ID, ipAddress, true)

	return backupCodes, nil
}

// ValidateTOTPForLogin validates a TOTP code during login.
func (s *TOTPService) ValidateTOTPForLogin(ctx context.Context, user *models.User, code string, ipAddress net.IP) error {
	// Check rate limit
	if exceeded, err := s.checkRateLimit(ctx, user.ID, ipAddress); err != nil {
		return err
	} else if exceeded {
		return ErrTOTPRateLimitExceeded
	}

	// Check if 2FA is enabled
	if !user.TwoFactorEnabled || user.TwoFactorSecret == nil {
		return ErrTOTPNotEnabled
	}

	// Decrypt the secret
	secret, err := s.DecryptSecret(*user.TwoFactorSecret)
	if err != nil {
		s.recordTOTPAttempt(ctx, user.ID, ipAddress, false)
		return fmt.Errorf("failed to decrypt TOTP secret: %w", err)
	}

	// Validate the code
	if !s.ValidateCode(secret, code) {
		// Check if it's a backup code
		if s.validateBackupCode(ctx, user.ID, code, ipAddress) {
			s.recordTOTPAttempt(ctx, user.ID, ipAddress, true)
			return nil
		}
		s.recordTOTPAttempt(ctx, user.ID, ipAddress, false)
		return ErrTOTPInvalidCode
	}

	// Record successful attempt
	s.recordTOTPAttempt(ctx, user.ID, ipAddress, true)

	return nil
}

// validateBackupCode validates and consumes a backup code.
func (s *TOTPService) validateBackupCode(ctx context.Context, userID uuid.UUID, code string, ipAddress net.IP) bool {
	codeHash := s.HashBackupCode(code)

	var backupCode models.BackupCode
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND code_hash = ? AND used_at IS NULL", userID, codeHash).
		First(&backupCode).Error
	if err != nil {
		return false
	}

	// Mark the code as used
	now := time.Now()
	backupCode.UsedAt = &now
	if err := s.db.WithContext(ctx).Save(&backupCode).Error; err != nil {
		return false
	}

	return true
}

// DisableTOTP disables 2FA for a user.
func (s *TOTPService) DisableTOTP(ctx context.Context, user *models.User, password string, passwordService *PasswordService, ipAddress net.IP) error {
	// Verify password
	if user.PasswordHash == nil {
		return ErrInvalidCredentials
	}
	if err := passwordService.VerifyPassword(password, *user.PasswordHash); err != nil {
		return ErrInvalidCredentials
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete backup codes
	if err := tx.Where("user_id = ?", user.ID).Delete(&models.BackupCode{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Disable 2FA
	user.TwoFactorEnabled = false
	user.TwoFactorSecret = nil
	user.TOTPVerifiedAt = nil

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// GetBackupCodes retrieves the remaining (unused) backup codes count for a user.
func (s *TOTPService) GetBackupCodesCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&models.BackupCode{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// RegenerateBackupCodes generates new backup codes for a user.
func (s *TOTPService) RegenerateBackupCodes(ctx context.Context, user *models.User, code string, ipAddress net.IP) ([]string, error) {
	// Check rate limit
	if exceeded, err := s.checkRateLimit(ctx, user.ID, ipAddress); err != nil {
		return nil, err
	} else if exceeded {
		return nil, ErrTOTPRateLimitExceeded
	}

	// Check if 2FA is enabled
	if !user.TwoFactorEnabled || user.TwoFactorSecret == nil {
		return nil, ErrTOTPNotEnabled
	}

	// Decrypt the secret
	secret, err := s.DecryptSecret(*user.TwoFactorSecret)
	if err != nil {
		s.recordTOTPAttempt(ctx, user.ID, ipAddress, false)
		return nil, fmt.Errorf("failed to decrypt TOTP secret: %w", err)
	}

	// Validate the current code
	if !s.ValidateCode(secret, code) {
		s.recordTOTPAttempt(ctx, user.ID, ipAddress, false)
		return nil, ErrTOTPInvalidCode
	}

	// Generate new backup codes
	backupCodes, err := s.GenerateBackupCodes()
	if err != nil {
		return nil, err
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing backup codes
	if err := tx.Where("user_id = ?", user.ID).Delete(&models.BackupCode{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Store new hashed backup codes
	for _, code := range backupCodes {
		backupCode := &models.BackupCode{
			UserID:   user.ID,
			CodeHash: s.HashBackupCode(code),
		}
		if err := tx.Create(backupCode).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Record successful attempt
	s.recordTOTPAttempt(ctx, user.ID, ipAddress, true)

	return backupCodes, nil
}

// checkRateLimit checks if the user has exceeded the TOTP attempt rate limit.
func (s *TOTPService) checkRateLimit(ctx context.Context, userID uuid.UUID, ipAddress net.IP) (bool, error) {
	var count int64
	windowStart := time.Now().Add(-TOTPRateLimitWindow)

	// Count attempts from either the user ID or IP address
	err := s.db.WithContext(ctx).Model(&models.TOTPAttempt{}).
		Where("(user_id = ? OR ip_address = ?) AND created_at > ? AND success = false", userID, ipAddress, windowStart).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count >= TOTPRateLimitAttempts, nil
}

// recordTOTPAttempt records a TOTP validation attempt.
func (s *TOTPService) recordTOTPAttempt(ctx context.Context, userID uuid.UUID, ipAddress net.IP, success bool) {
	attempt := &models.TOTPAttempt{
		UserID:    userID,
		IPAddress: ipAddress,
		Success:   success,
	}
	s.db.WithContext(ctx).Create(attempt)
}

// IsTOTPEnabled checks if a user has 2FA enabled.
func (s *TOTPService) IsTOTPEnabled(user *models.User) bool {
	return user.TwoFactorEnabled && user.TwoFactorSecret != nil && user.TOTPVerifiedAt != nil
}
