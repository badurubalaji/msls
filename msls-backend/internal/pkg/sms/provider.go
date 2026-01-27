// Package sms provides SMS sending capabilities with support for multiple providers.
package sms

import (
	"context"
	"errors"
)

// Common errors for SMS operations.
var (
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrSendFailed         = errors.New("failed to send SMS")
	ErrProviderNotReady   = errors.New("SMS provider not ready")
	ErrRateLimited        = errors.New("SMS rate limit exceeded")
)

// Message represents an SMS message to be sent.
type Message struct {
	To      string // Phone number in E.164 format (e.g., +919876543210)
	Body    string // Message content
	From    string // Optional sender ID
}

// SendResult represents the result of sending an SMS.
type SendResult struct {
	MessageID string // Provider-specific message ID
	Status    string // Status of the send operation
}

// Provider defines the interface for SMS providers.
// Implementations can include Twilio, AWS SNS, Vonage, etc.
type Provider interface {
	// Send sends an SMS message and returns the result.
	Send(ctx context.Context, msg Message) (*SendResult, error)

	// Name returns the provider name for logging purposes.
	Name() string

	// IsReady returns true if the provider is configured and ready to send.
	IsReady() bool
}

// ProviderConfig holds configuration for SMS providers.
type ProviderConfig struct {
	// Twilio configuration
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string

	// AWS SNS configuration
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	// Mock configuration for development
	MockEnabled bool
	MockLogPath string
}
