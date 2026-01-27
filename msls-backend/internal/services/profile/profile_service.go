// Package profile provides profile management services for the MSLS application.
package profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
	"msls-backend/internal/services/auth"
)

// Import image decoders
var _ = image.Decode

// Max avatar size: 2MB.
const MaxAvatarSize = 2 * 1024 * 1024

// Allowed avatar MIME types.
var AllowedAvatarTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// Valid timezones (a subset - in production you'd use a complete list).
var ValidTimezones = map[string]bool{
	"UTC":                 true,
	"America/New_York":    true,
	"America/Chicago":     true,
	"America/Denver":      true,
	"America/Los_Angeles": true,
	"Europe/London":       true,
	"Europe/Paris":        true,
	"Europe/Berlin":       true,
	"Asia/Tokyo":          true,
	"Asia/Shanghai":       true,
	"Asia/Kolkata":        true,
	"Asia/Dubai":          true,
	"Australia/Sydney":    true,
	"Pacific/Auckland":    true,
}

// Valid locales.
var ValidLocales = map[string]bool{
	"en":    true,
	"en-US": true,
	"en-GB": true,
	"es":    true,
	"fr":    true,
	"de":    true,
	"ja":    true,
	"zh":    true,
	"hi":    true,
	"ar":    true,
}

// NotificationPreferences represents user notification settings.
type NotificationPreferences struct {
	Email bool `json:"email"`
	Push  bool `json:"push"`
	SMS   bool `json:"sms"`
}

// UpdateProfileRequest represents a profile update request.
type UpdateProfileRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Timezone  *string `json:"timezone,omitempty"`
	Locale    *string `json:"locale,omitempty"`
}

// ChangePasswordRequest represents a password change request.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// UpdatePreferencesRequest represents a notification preferences update request.
type UpdatePreferencesRequest struct {
	Email *bool `json:"email,omitempty"`
	Push  *bool `json:"push,omitempty"`
	SMS   *bool `json:"sms,omitempty"`
}

// UserPreference represents a user preference stored in the database.
type UserPreference struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v7()"`
	UserID    uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	Category  string          `json:"category" gorm:"type:varchar(50);not null"`
	Key       string          `json:"key" gorm:"type:varchar(100);not null"`
	Value     json.RawMessage `json:"value" gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// TableName returns the table name for UserPreference.
func (UserPreference) TableName() string {
	return "user_preferences"
}

// ProfileService handles profile management operations.
type ProfileService struct {
	db              *gorm.DB
	passwordService *auth.PasswordService
	uploadDir       string
}

// Config holds configuration for the profile service.
type Config struct {
	UploadDir string
}

// NewProfileService creates a new ProfileService instance.
func NewProfileService(db *gorm.DB, cfg Config) *ProfileService {
	uploadDir := cfg.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads/avatars"
	}

	// Ensure upload directory exists
	_ = os.MkdirAll(uploadDir, 0755)

	return &ProfileService{
		db:              db,
		passwordService: auth.NewPasswordService(),
		uploadDir:       uploadDir,
	}
}

// GetProfile retrieves a user's full profile by ID.
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
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

// UpdateProfile updates a user's profile information.
func (s *ProfileService) UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Validate timezone if provided
	if req.Timezone != nil && *req.Timezone != "" {
		if !ValidTimezones[*req.Timezone] {
			return nil, ErrInvalidTimezone
		}
	}

	// Validate locale if provided
	if req.Locale != nil && *req.Locale != "" {
		if !ValidLocales[*req.Locale] {
			return nil, ErrInvalidLocale
		}
	}

	// Build updates
	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Timezone != nil {
		updates["timezone"] = *req.Timezone
	}
	if req.Locale != nil {
		updates["locale"] = *req.Locale
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := s.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// Reload user with associations
	return s.GetProfile(ctx, userID)
}

// ChangePassword changes a user's password.
func (s *ProfileService) ChangePassword(ctx context.Context, userID uuid.UUID, req ChangePasswordRequest, ipAddress net.IP, userAgent string) error {
	// Validate new password matches confirmation
	if req.NewPassword != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	// Validate new password strength
	if err := s.passwordService.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Verify current password
	if user.PasswordHash == nil {
		return ErrInvalidCurrentPassword
	}
	if err := s.passwordService.VerifyPassword(req.CurrentPassword, *user.PasswordHash); err != nil {
		return ErrInvalidCurrentPassword
	}

	// Hash new password
	passwordHash, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	if err := s.db.WithContext(ctx).Model(&user).Updates(map[string]interface{}{
		"password_hash": passwordHash,
		"updated_at":    time.Now(),
	}).Error; err != nil {
		return err
	}

	// Create audit log
	s.createAuditLog(ctx, &user, models.AuditActionPasswordChange, ipAddress, userAgent)

	return nil
}

// UploadAvatar handles avatar file upload.
func (s *ProfileService) UploadAvatar(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Check file size
	if header.Size > MaxAvatarSize {
		return "", ErrAvatarTooLarge
	}

	// Check content type
	contentType := header.Header.Get("Content-Type")
	if !AllowedAvatarTypes[contentType] {
		return "", ErrInvalidAvatarFormat
	}

	// Verify user exists
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		if contentType == "image/jpeg" {
			ext = ".jpg"
		} else {
			ext = ".png"
		}
	}
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(s.uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Reset file position
	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	// Copy file
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Generate thumbnail (optional - resize for consistency)
	if err := s.generateThumbnail(filePath, ext); err != nil {
		// Log error but don't fail the upload
		fmt.Printf("Warning: failed to generate thumbnail: %v\n", err)
	}

	// Update user avatar URL
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := s.db.WithContext(ctx).Model(&user).Update("avatar_url", avatarURL).Error; err != nil {
		// Clean up the uploaded file
		os.Remove(filePath)
		return "", err
	}

	// Delete old avatar if exists
	if user.AvatarURL != nil && *user.AvatarURL != "" {
		oldPath := filepath.Join(".", *user.AvatarURL)
		os.Remove(oldPath)
	}

	return avatarURL, nil
}

// generateThumbnail creates a resized version of the avatar image.
func (s *ProfileService) generateThumbnail(filePath, ext string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize to max 256x256 maintaining aspect ratio
	thumbnail := resize.Thumbnail(256, 256, img, resize.Lanczos3)

	// Create thumbnail file
	thumbPath := strings.TrimSuffix(filePath, ext) + "_thumb" + ext
	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		return err
	}
	defer thumbFile.Close()

	// Encode and save
	ext = strings.ToLower(ext)
	if ext == ".jpg" || ext == ".jpeg" {
		return jpeg.Encode(thumbFile, thumbnail, &jpeg.Options{Quality: 85})
	}
	return png.Encode(thumbFile, thumbnail)
}

// UpdateNotificationPreferences updates a user's notification preferences.
func (s *ProfileService) UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, req UpdatePreferencesRequest) (*NotificationPreferences, error) {
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Get current preferences
	var currentPrefs NotificationPreferences
	if user.NotificationPreferences != nil {
		if err := json.Unmarshal(user.NotificationPreferences, &currentPrefs); err != nil {
			// Use defaults if parsing fails
			currentPrefs = NotificationPreferences{Email: true, Push: true, SMS: false}
		}
	}

	// Apply updates
	if req.Email != nil {
		currentPrefs.Email = *req.Email
	}
	if req.Push != nil {
		currentPrefs.Push = *req.Push
	}
	if req.SMS != nil {
		currentPrefs.SMS = *req.SMS
	}

	// Marshal back to JSON
	prefsJSON, err := json.Marshal(currentPrefs)
	if err != nil {
		return nil, err
	}

	// Update in database
	if err := s.db.WithContext(ctx).Model(&user).Updates(map[string]interface{}{
		"notification_preferences": prefsJSON,
		"updated_at":               time.Now(),
	}).Error; err != nil {
		return nil, err
	}

	return &currentPrefs, nil
}

// GetNotificationPreferences retrieves a user's notification preferences.
func (s *ProfileService) GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreferences, error) {
	var user models.User
	err := s.db.WithContext(ctx).Select("notification_preferences").First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	var prefs NotificationPreferences
	if user.NotificationPreferences != nil {
		if err := json.Unmarshal(user.NotificationPreferences, &prefs); err != nil {
			// Return defaults if parsing fails
			return &NotificationPreferences{Email: true, Push: true, SMS: false}, nil
		}
	} else {
		prefs = NotificationPreferences{Email: true, Push: true, SMS: false}
	}

	return &prefs, nil
}

// RequestAccountDeletion initiates an account deletion request.
func (s *ProfileService) RequestAccountDeletion(ctx context.Context, userID uuid.UUID, ipAddress net.IP, userAgent string) error {
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check if deletion already requested
	if user.AccountDeletionRequestedAt != nil {
		return ErrAccountDeletionAlreadyRequested
	}

	// Set deletion request timestamp
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&user).Updates(map[string]interface{}{
		"account_deletion_requested_at": now,
		"status":                        models.StatusInactive,
		"updated_at":                    now,
	}).Error; err != nil {
		return err
	}

	// Create audit log
	s.createAuditLog(ctx, &user, models.AuditActionAccountDeletionRequested, ipAddress, userAgent)

	return nil
}

// GetUserPreferences retrieves all preferences for a user in a category.
func (s *ProfileService) GetUserPreferences(ctx context.Context, userID uuid.UUID, category string) ([]UserPreference, error) {
	var prefs []UserPreference
	query := s.db.WithContext(ctx).Where("user_id = ?", userID)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if err := query.Find(&prefs).Error; err != nil {
		return nil, err
	}
	return prefs, nil
}

// SetUserPreference sets a user preference.
func (s *ProfileService) SetUserPreference(ctx context.Context, userID uuid.UUID, category, key string, value interface{}) error {
	// Verify user exists
	var user models.User
	err := s.db.WithContext(ctx).Select("id").First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Marshal value to JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Upsert preference
	pref := UserPreference{
		UserID:    userID,
		Category:  category,
		Key:       key,
		Value:     valueJSON,
		UpdatedAt: time.Now(),
	}

	return s.db.WithContext(ctx).
		Where("user_id = ? AND category = ? AND key = ?", userID, category, key).
		Assign(map[string]interface{}{
			"value":      valueJSON,
			"updated_at": time.Now(),
		}).
		FirstOrCreate(&pref).Error
}

// DeleteUserPreference deletes a user preference.
func (s *ProfileService) DeleteUserPreference(ctx context.Context, userID uuid.UUID, category, key string) error {
	result := s.db.WithContext(ctx).
		Where("user_id = ? AND category = ? AND key = ?", userID, category, key).
		Delete(&UserPreference{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPreferenceNotFound
	}
	return nil
}

// RecordLogin updates the last login information.
func (s *ProfileService) RecordLogin(ctx context.Context, userID uuid.UUID, ipAddress net.IP) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ipAddress,
	}).Error
}

// createAuditLog creates an audit log entry.
func (s *ProfileService) createAuditLog(ctx context.Context, user *models.User, action models.AuditAction, ipAddress net.IP, userAgent string) {
	log := models.NewAuditLog(action, "user").
		WithTenant(user.TenantID).
		WithUser(user.ID).
		WithEntity(user.ID).
		WithIPAddress(ipAddress).
		WithUserAgent(userAgent).
		Build()
	s.db.WithContext(ctx).Create(log)
}
