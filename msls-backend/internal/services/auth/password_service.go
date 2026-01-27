// Package auth provides authentication services for the MSLS application.
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"unicode"

	"golang.org/x/crypto/argon2"
)

// Password hashing parameters using Argon2id.
const (
	// Argon2id parameters (OWASP recommended)
	argon2Time    = 1
	argon2Memory  = 64 * 1024 // 64 MB
	argon2Threads = 4
	argon2KeyLen  = 32
	argon2SaltLen = 16
)

// Password validation requirements.
const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
)

// Note: Password validation errors are defined in errors.go

// PasswordService handles password hashing and validation.
type PasswordService struct {
	// Configuration can be added here for customization
}

// NewPasswordService creates a new PasswordService instance.
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword hashes a password using Argon2id.
// Returns the hash in the format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
func (s *PasswordService) HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Encode salt and hash to base64
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	// Return in PHC string format
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Time, argon2Threads, saltB64, hashB64), nil
}

// VerifyPassword verifies a password against a hash.
func (s *PasswordService) VerifyPassword(password, encodedHash string) error {
	// Parse the encoded hash
	var version int
	var memory, time uint32
	var threads uint8
	var saltB64, hashB64 string

	// Parse the PHC string format
	_, err := fmt.Sscanf(encodedHash, "$argon2id$v=%d$m=%d,t=%d,p=%d$%s",
		&version, &memory, &time, &threads, &saltB64)
	if err != nil {
		return ErrInvalidPasswordHash
	}

	// Split saltB64 and hashB64 (they're separated by $)
	parts := regexp.MustCompile(`\$`).Split(encodedHash, -1)
	if len(parts) != 6 {
		return ErrInvalidPasswordHash
	}
	saltB64 = parts[4]
	hashB64 = parts[5]

	// Decode the salt
	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return ErrInvalidPasswordHash
	}

	// Decode the expected hash
	expectedHash, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return ErrInvalidPasswordHash
	}

	// Compute the hash of the provided password
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(expectedHash)))

	// Compare hashes using constant-time comparison
	if !constantTimeCompare(computedHash, expectedHash) {
		return ErrPasswordHashMismatch
	}

	return nil
}

// ValidatePassword validates a password against the security requirements.
func (s *PasswordService) ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}
	if len(password) > MaxPasswordLength {
		return ErrPasswordTooLong
	}

	var hasUppercase, hasLowercase, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUppercase {
		return ErrPasswordNoUppercase
	}
	if !hasLowercase {
		return ErrPasswordNoLowercase
	}
	if !hasDigit {
		return ErrPasswordNoDigit
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

// GenerateRandomToken generates a cryptographically secure random token.
func (s *PasswordService) GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// constantTimeCompare compares two byte slices in constant time.
func constantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
