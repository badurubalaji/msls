// Package auth provides authentication services for the MSLS application.
package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Note: JWT errors are defined in errors.go

// Claims represents the JWT claims structure.
type Claims struct {
	jwt.RegisteredClaims
	UserID      uuid.UUID `json:"user_id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Email       string    `json:"email,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
}

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// JWTConfig holds JWT service configuration.
type JWTConfig struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// NewJWTService creates a new JWTService instance.
func NewJWTService(config JWTConfig) *JWTService {
	return &JWTService{
		secret:     []byte(config.Secret),
		issuer:     config.Issuer,
		accessTTL:  config.AccessTTL,
		refreshTTL: config.RefreshTTL,
	}
}

// GenerateAccessToken generates a new JWT access token for a user.
func (s *JWTService) GenerateAccessToken(userID, tenantID uuid.UUID, email string, permissions []string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessTTL)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		UserID:      userID,
		TenantID:    tenantID,
		Email:       email,
		Permissions: permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a cryptographically secure refresh token.
func (s *JWTService) GenerateRefreshToken() (string, time.Time, error) {
	ps := NewPasswordService()
	token, err := ps.GenerateRandomToken(32)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Now().Add(s.refreshTTL)
	return token, expiresAt, nil
}

// ValidateAccessToken validates an access token and returns its claims.
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotYetValid
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// HashRefreshToken hashes a refresh token using SHA-256.
func (s *JWTService) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GetAccessTTL returns the access token TTL.
func (s *JWTService) GetAccessTTL() time.Duration {
	return s.accessTTL
}

// GetRefreshTTL returns the refresh token TTL.
func (s *JWTService) GetRefreshTTL() time.Duration {
	return s.refreshTTL
}
