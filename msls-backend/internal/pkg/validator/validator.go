// Package validator provides custom validation functions and utilities.
package validator

import (
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	// validate is the singleton validator instance.
	validate *validator.Validate

	// once ensures validator is initialized only once.
	once sync.Once

	// phoneRegex validates phone numbers in various formats.
	// Supports: +1234567890, 123-456-7890, (123) 456-7890, etc.
	phoneRegex = regexp.MustCompile(`^[\+]?[(]?[0-9]{1,4}[)]?[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,9}$`)

	// slugRegex validates URL-friendly slugs.
	slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

	// alphanumericSpaceRegex validates alphanumeric strings with spaces.
	alphanumericSpaceRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
)

// Get returns the singleton validator instance with all custom validators registered.
func Get() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		registerCustomValidators(validate)
	})
	return validate
}

// registerCustomValidators registers all custom validation functions.
func registerCustomValidators(v *validator.Validate) {
	// Register custom validators
	_ = v.RegisterValidation("phone", validatePhone)
	_ = v.RegisterValidation("slug", validateSlug)
	_ = v.RegisterValidation("alphanumeric_space", validateAlphanumericSpace)
	_ = v.RegisterValidation("strong_password", validateStrongPassword)
	_ = v.RegisterValidation("tenant_id", validateTenantID)
	_ = v.RegisterValidation("not_blank", validateNotBlank)
}

// validatePhone validates phone numbers.
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // Allow empty, use 'required' tag to enforce non-empty
	}
	// Remove all non-digit characters except + for country code
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '+' {
			return r
		}
		return -1
	}, phone)

	// Phone should have between 7 and 15 digits
	digitCount := len(strings.ReplaceAll(cleaned, "+", ""))
	if digitCount < 7 || digitCount > 15 {
		return false
	}

	return phoneRegex.MatchString(phone)
}

// validateSlug validates URL-friendly slugs.
func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	if slug == "" {
		return true
	}
	return slugRegex.MatchString(slug)
}

// validateAlphanumericSpace validates alphanumeric strings with spaces.
func validateAlphanumericSpace(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "" {
		return true
	}
	return alphanumericSpaceRegex.MatchString(str)
}

// validateStrongPassword validates password strength.
// Requirements: 8+ chars, at least one uppercase, one lowercase, one digit, one special char.
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return true
	}

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// validateTenantID validates a tenant ID format (UUID or custom format).
func validateTenantID(fl validator.FieldLevel) bool {
	tenantID := fl.Field().String()
	if tenantID == "" {
		return true
	}

	// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	uuidRegex := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	return uuidRegex.MatchString(tenantID)
}

// validateNotBlank validates that a string is not just whitespace.
func validateNotBlank(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	return strings.TrimSpace(str) != ""
}

// ValidateStruct validates a struct and returns validation errors.
func ValidateStruct(s interface{}) error {
	return Get().Struct(s)
}

// ValidateVar validates a single variable against a tag.
func ValidateVar(field interface{}, tag string) error {
	return Get().Var(field, tag)
}

// FieldErrors extracts field errors from a validator error.
func FieldErrors(err error) []FieldErrorDetail {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make([]FieldErrorDetail, 0, len(validationErrors))
		for _, e := range validationErrors {
			errors = append(errors, FieldErrorDetail{
				Field:   toSnakeCase(e.Field()),
				Tag:     e.Tag(),
				Value:   e.Value(),
				Message: getErrorMessage(e),
			})
		}
		return errors
	}
	return nil
}

// FieldErrorDetail represents a single field validation error.
type FieldErrorDetail struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value,omitempty"`
	Message string      `json:"message"`
}

// getErrorMessage returns a human-readable error message for a validation error.
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "phone":
		return "Invalid phone number format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "gte":
		return "Value must be greater than or equal to " + e.Param()
	case "lte":
		return "Value must be less than or equal to " + e.Param()
	case "uuid":
		return "Invalid UUID format"
	case "slug":
		return "Invalid slug format (use lowercase letters, numbers, and hyphens)"
	case "strong_password":
		return "Password must be at least 8 characters with uppercase, lowercase, number, and special character"
	case "not_blank":
		return "This field cannot be blank"
	case "oneof":
		return "Value must be one of: " + e.Param()
	case "alphanum":
		return "Value must be alphanumeric"
	case "alphanumeric_space":
		return "Value must contain only letters, numbers, and spaces"
	case "url":
		return "Invalid URL format"
	case "datetime":
		return "Invalid datetime format"
	case "tenant_id":
		return "Invalid tenant ID format"
	default:
		return "Validation failed on '" + e.Tag() + "' constraint"
	}
}

// toSnakeCase converts a string to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// IsValid checks if a struct is valid without returning detailed errors.
func IsValid(s interface{}) bool {
	return ValidateStruct(s) == nil
}

// IsValidPhone checks if a phone number is valid.
func IsValidPhone(phone string) bool {
	return ValidateVar(phone, "phone") == nil
}

// IsValidEmail checks if an email is valid.
func IsValidEmail(email string) bool {
	return ValidateVar(email, "email") == nil
}

// IsValidUUID checks if a string is a valid UUID.
func IsValidUUID(uuid string) bool {
	return ValidateVar(uuid, "uuid") == nil
}

// IsValidSlug checks if a string is a valid slug.
func IsValidSlug(slug string) bool {
	return ValidateVar(slug, "slug") == nil
}
