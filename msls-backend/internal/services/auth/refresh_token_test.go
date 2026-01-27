// Package auth provides authentication services for the MSLS application.
package auth

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// setupRefreshTokenTestDB creates an in-memory SQLite database for testing.
func setupRefreshTokenTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create simplified tables compatible with SQLite (not using PostgreSQL-specific features)
	err = db.Exec(`
		CREATE TABLE tenants (
			id TEXT PRIMARY KEY,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			name TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			settings TEXT DEFAULT '{}',
			status TEXT NOT NULL DEFAULT 'active'
		)
	`).Error
	require.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			tenant_id TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT,
			updated_by TEXT,
			email TEXT,
			phone TEXT,
			password_hash TEXT,
			first_name TEXT DEFAULT '',
			last_name TEXT DEFAULT '',
			avatar_url TEXT,
			bio TEXT,
			timezone TEXT DEFAULT 'UTC',
			locale TEXT DEFAULT 'en',
			notification_preferences TEXT DEFAULT '{}',
			status TEXT NOT NULL DEFAULT 'active',
			two_factor_enabled INTEGER DEFAULT 0,
			two_factor_secret TEXT,
			totp_verified_at DATETIME,
			email_verified_at DATETIME,
			phone_verified_at DATETIME,
			last_login_at DATETIME,
			last_login_ip TEXT,
			locked_until DATETIME,
			failed_login_attempts INTEGER DEFAULT 0,
			account_deletion_requested_at DATETIME,
			FOREIGN KEY (tenant_id) REFERENCES tenants(id)
		)
	`).Error
	require.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE refresh_tokens (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			token_hash TEXT NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			revoked_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`).Error
	require.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE roles (
			id TEXT PRIMARY KEY,
			tenant_id TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			name TEXT NOT NULL,
			display_name TEXT,
			description TEXT,
			is_system INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active'
		)
	`).Error
	require.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE permissions (
			id TEXT PRIMARY KEY,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			description TEXT,
			module TEXT NOT NULL
		)
	`).Error
	require.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE user_roles (
			user_id TEXT NOT NULL,
			role_id TEXT NOT NULL,
			PRIMARY KEY (user_id, role_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (role_id) REFERENCES roles(id)
		)
	`).Error
	require.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE role_permissions (
			role_id TEXT NOT NULL,
			permission_id TEXT NOT NULL,
			PRIMARY KEY (role_id, permission_id),
			FOREIGN KEY (role_id) REFERENCES roles(id),
			FOREIGN KEY (permission_id) REFERENCES permissions(id)
		)
	`).Error
	require.NoError(t, err)

	err = db.Exec(`
		CREATE TABLE audit_logs (
			id TEXT PRIMARY KEY,
			tenant_id TEXT,
			user_id TEXT,
			action TEXT NOT NULL,
			entity_type TEXT NOT NULL,
			entity_id TEXT,
			old_data TEXT,
			new_data TEXT,
			ip_address TEXT,
			user_agent TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	require.NoError(t, err)

	return db
}

// createTestTenantForRefresh creates a test tenant in the database.
func createTestTenantForRefresh(t *testing.T, db *gorm.DB) *models.Tenant {
	t.Helper()

	tenant := &models.Tenant{
		Name:   "Test Tenant",
		Slug:   "test-tenant",
		Status: models.StatusActive,
	}
	tenant.ID = uuid.New()
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	err := db.Create(tenant).Error
	require.NoError(t, err)

	return tenant
}

// createTestUserForRefresh creates a test user in the database.
func createTestUserForRefresh(t *testing.T, db *gorm.DB, tenant *models.Tenant) *models.User {
	t.Helper()

	email := "test@example.com"
	passwordHash := "$argon2id$v=19$m=65536,t=3,p=4$hash" // Dummy hash for testing
	user := &models.User{
		TenantID:     tenant.ID,
		Email:        &email,
		PasswordHash: &passwordHash,
		FirstName:    "Test",
		LastName:     "User",
		Status:       models.StatusActive,
	}
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := db.Create(user).Error
	require.NoError(t, err)

	return user
}

// createTestRefreshToken creates a test refresh token in the database.
func createTestRefreshToken(t *testing.T, db *gorm.DB, user *models.User, jwtService *JWTService, expired bool, revoked bool) (string, *models.RefreshToken) {
	t.Helper()

	// Generate a raw refresh token
	rawToken, expiresAt, err := jwtService.GenerateRefreshToken()
	require.NoError(t, err)

	// Adjust expiry if testing expired token
	if expired {
		expiresAt = time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	}

	// Create the stored token
	storedToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: jwtService.HashRefreshToken(rawToken),
		ExpiresAt: expiresAt,
	}
	storedToken.ID = uuid.New()
	storedToken.CreatedAt = time.Now()
	storedToken.UpdatedAt = time.Now()

	if revoked {
		now := time.Now()
		storedToken.RevokedAt = &now
	}

	err = db.Create(storedToken).Error
	require.NoError(t, err)

	// Reload with user relationship
	err = db.Preload("User").First(storedToken, "id = ?", storedToken.ID).Error
	require.NoError(t, err)

	return rawToken, storedToken
}

// createTestAuthService creates an AuthService for testing.
func createTestAuthService(t *testing.T, db *gorm.DB) (*AuthService, *JWTService) {
	t.Helper()

	jwtConfig := JWTConfig{
		Secret:     "test-secret-key-at-least-32-bytes-long",
		Issuer:     "msls-test",
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 7 * 24 * time.Hour,
	}

	jwtService := NewJWTService(jwtConfig)
	authService := NewAuthService(db, jwtService)

	return authService, jwtService
}

// TestRefreshToken_Success tests successful token refresh.
func TestRefreshToken_Success(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, jwtService := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data
	tenant := createTestTenantForRefresh(t, db)
	user := createTestUserForRefresh(t, db, tenant)
	rawToken, storedToken := createTestRefreshToken(t, db, user, jwtService, false, false)

	// Execute refresh
	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	tokenPair, err := authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.NotEqual(t, rawToken, tokenPair.RefreshToken, "New refresh token should be different")
	assert.Greater(t, tokenPair.ExpiresIn, 0)

	// Verify old token is revoked
	var oldToken models.RefreshToken
	err = db.First(&oldToken, "id = ?", storedToken.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, oldToken.RevokedAt, "Old token should be revoked")

	// Verify new token is stored
	var tokenCount int64
	db.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked_at IS NULL", user.ID).Count(&tokenCount)
	assert.Equal(t, int64(1), tokenCount, "Should have exactly one active refresh token")

	// Verify audit logs were created
	var auditLogs []models.AuditLog
	db.Where("user_id = ?", user.ID).Order("created_at desc").Find(&auditLogs)
	assert.GreaterOrEqual(t, len(auditLogs), 2, "Should have audit logs for revocation and refresh")

	// Check for token_refresh audit log
	hasRefreshLog := false
	for _, log := range auditLogs {
		if log.Action == models.AuditActionTokenRefresh {
			hasRefreshLog = true
			assert.Equal(t, ipAddress.String(), log.IPAddress.String())
			assert.Equal(t, userAgent, log.UserAgent)
			break
		}
	}
	assert.True(t, hasRefreshLog, "Should have token_refresh audit log")
}

// TestRefreshToken_NotFound tests refresh with invalid token.
func TestRefreshToken_NotFound(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, _ := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data but don't create a token
	_ = createTestTenantForRefresh(t, db)

	// Try to refresh with invalid token
	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	tokenPair, err := authService.RefreshToken(ctx, "invalid-token", ipAddress, userAgent)

	// Verify
	assert.ErrorIs(t, err, ErrRefreshTokenNotFound)
	assert.Nil(t, tokenPair)

	// Verify audit log for failed refresh
	var auditLog models.AuditLog
	err = db.Where("action = ?", models.AuditActionTokenRefreshFailed).First(&auditLog).Error
	require.NoError(t, err)
	assert.Equal(t, "refresh_token", auditLog.EntityType)
}

// TestRefreshToken_Expired tests refresh with expired token.
func TestRefreshToken_Expired(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, jwtService := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data with expired token
	tenant := createTestTenantForRefresh(t, db)
	user := createTestUserForRefresh(t, db, tenant)
	rawToken, _ := createTestRefreshToken(t, db, user, jwtService, true, false) // expired=true

	// Try to refresh with expired token
	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	tokenPair, err := authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)

	// Verify
	assert.ErrorIs(t, err, ErrRefreshTokenExpired)
	assert.Nil(t, tokenPair)

	// Verify audit log for failed refresh
	var auditLog models.AuditLog
	err = db.Where("action = ? AND user_id = ?", models.AuditActionTokenRefreshFailed, user.ID).First(&auditLog).Error
	require.NoError(t, err)
}

// TestRefreshToken_Revoked tests refresh with revoked token.
func TestRefreshToken_Revoked(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, jwtService := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data with revoked token
	tenant := createTestTenantForRefresh(t, db)
	user := createTestUserForRefresh(t, db, tenant)
	rawToken, _ := createTestRefreshToken(t, db, user, jwtService, false, true) // revoked=true

	// Try to refresh with revoked token
	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	tokenPair, err := authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)

	// Verify
	assert.ErrorIs(t, err, ErrRefreshTokenRevoked)
	assert.Nil(t, tokenPair)

	// Verify audit log for failed refresh
	var auditLog models.AuditLog
	err = db.Where("action = ? AND user_id = ?", models.AuditActionTokenRefreshFailed, user.ID).First(&auditLog).Error
	require.NoError(t, err)
}

// TestRefreshToken_Rotation tests that old token is revoked after refresh.
func TestRefreshToken_Rotation(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, jwtService := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data
	tenant := createTestTenantForRefresh(t, db)
	user := createTestUserForRefresh(t, db, tenant)
	rawToken, storedToken := createTestRefreshToken(t, db, user, jwtService, false, false)

	// First refresh
	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	tokenPair1, err := authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)
	require.NoError(t, err)

	// Verify old token is revoked
	var oldToken models.RefreshToken
	db.First(&oldToken, "id = ?", storedToken.ID)
	assert.NotNil(t, oldToken.RevokedAt, "Old token should be revoked after first refresh")

	// Try to use old token again (should fail)
	_, err = authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)
	assert.ErrorIs(t, err, ErrRefreshTokenRevoked, "Using old token after rotation should fail")

	// Use new token (should succeed)
	tokenPair2, err := authService.RefreshToken(ctx, tokenPair1.RefreshToken, ipAddress, userAgent)
	require.NoError(t, err)
	assert.NotNil(t, tokenPair2)
	assert.NotEqual(t, tokenPair1.RefreshToken, tokenPair2.RefreshToken)
}

// TestRefreshToken_ConcurrentRefresh tests concurrent refresh attempts.
func TestRefreshToken_ConcurrentRefresh(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, jwtService := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data
	tenant := createTestTenantForRefresh(t, db)
	user := createTestUserForRefresh(t, db, tenant)
	rawToken, _ := createTestRefreshToken(t, db, user, jwtService, false, false)

	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	// Execute concurrent refresh attempts
	results := make(chan error, 2)

	go func() {
		_, err := authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)
		results <- err
	}()

	go func() {
		time.Sleep(10 * time.Millisecond) // Slight delay to ensure ordering
		_, err := authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)
		results <- err
	}()

	// Collect results
	err1 := <-results
	err2 := <-results

	// At least one should succeed, the other should fail with revoked
	successCount := 0
	revokedCount := 0
	for _, err := range []error{err1, err2} {
		if err == nil {
			successCount++
		} else if err == ErrRefreshTokenRevoked {
			revokedCount++
		}
	}

	assert.Equal(t, 1, successCount, "Exactly one refresh should succeed")
	assert.Equal(t, 1, revokedCount, "Exactly one refresh should fail with revoked")
}

// TestRefreshToken_AuditLogContainsIP tests that audit logs contain IP address.
func TestRefreshToken_AuditLogContainsIP(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, jwtService := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data
	tenant := createTestTenantForRefresh(t, db)
	user := createTestUserForRefresh(t, db, tenant)
	rawToken, _ := createTestRefreshToken(t, db, user, jwtService, false, false)

	// Execute refresh with specific IP
	ipAddress := net.ParseIP("10.0.0.42")
	userAgent := "CustomApp/1.0"

	_, err := authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)
	require.NoError(t, err)

	// Verify audit log contains correct IP and user agent
	var auditLog models.AuditLog
	err = db.Where("action = ? AND user_id = ?", models.AuditActionTokenRefresh, user.ID).First(&auditLog).Error
	require.NoError(t, err)

	assert.Equal(t, "10.0.0.42", auditLog.IPAddress.String())
	assert.Equal(t, "CustomApp/1.0", auditLog.UserAgent)
}

// TestRefreshToken_MultipleUsersIsolation tests that tokens are isolated between users.
func TestRefreshToken_MultipleUsersIsolation(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, jwtService := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data for two users
	tenant := createTestTenantForRefresh(t, db)
	user1 := createTestUserForRefresh(t, db, tenant)

	// Create second user with different email
	email2 := "test2@example.com"
	passwordHash := "$argon2id$v=19$m=65536,t=3,p=4$hash"
	user2 := &models.User{
		TenantID:     tenant.ID,
		Email:        &email2,
		PasswordHash: &passwordHash,
		FirstName:    "Test2",
		LastName:     "User2",
		Status:       models.StatusActive,
	}
	user2.ID = uuid.New()
	user2.CreatedAt = time.Now()
	user2.UpdatedAt = time.Now()
	db.Create(user2)

	// Create tokens for both users
	rawToken1, _ := createTestRefreshToken(t, db, user1, jwtService, false, false)
	rawToken2, _ := createTestRefreshToken(t, db, user2, jwtService, false, false)

	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	// Refresh user1's token
	tokenPair1, err := authService.RefreshToken(ctx, rawToken1, ipAddress, userAgent)
	require.NoError(t, err)
	assert.NotNil(t, tokenPair1)

	// User2's token should still work
	tokenPair2, err := authService.RefreshToken(ctx, rawToken2, ipAddress, userAgent)
	require.NoError(t, err)
	assert.NotNil(t, tokenPair2)

	// Tokens should be different
	assert.NotEqual(t, tokenPair1.AccessToken, tokenPair2.AccessToken)
	assert.NotEqual(t, tokenPair1.RefreshToken, tokenPair2.RefreshToken)
}

// TestRefreshToken_GeneratesValidAccessToken tests that the generated access token is valid.
func TestRefreshToken_GeneratesValidAccessToken(t *testing.T) {
	t.Parallel()

	db := setupRefreshTokenTestDB(t)
	authService, jwtService := createTestAuthService(t, db)
	ctx := context.Background()

	// Create test data
	tenant := createTestTenantForRefresh(t, db)
	user := createTestUserForRefresh(t, db, tenant)
	rawToken, _ := createTestRefreshToken(t, db, user, jwtService, false, false)

	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	// Execute refresh
	tokenPair, err := authService.RefreshToken(ctx, rawToken, ipAddress, userAgent)
	require.NoError(t, err)

	// Validate the access token
	claims, err := jwtService.ValidateAccessToken(tokenPair.AccessToken)
	require.NoError(t, err)

	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, tenant.ID, claims.TenantID)
	assert.Equal(t, *user.Email, claims.Email)
}
