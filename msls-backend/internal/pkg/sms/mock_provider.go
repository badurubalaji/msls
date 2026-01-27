// Package sms provides SMS sending capabilities with support for multiple providers.
package sms

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MockProvider is a mock SMS provider for development and testing.
// It logs SMS messages instead of actually sending them.
type MockProvider struct {
	mu       sync.Mutex
	messages []SentMessage
	logger   *log.Logger
	logFile  *os.File
}

// SentMessage represents a message that was "sent" by the mock provider.
type SentMessage struct {
	ID        string
	To        string
	Body      string
	From      string
	SentAt    time.Time
}

// NewMockProvider creates a new mock SMS provider.
// If logPath is provided, messages will also be written to a file.
func NewMockProvider(logPath string) (*MockProvider, error) {
	provider := &MockProvider{
		messages: make([]SentMessage, 0),
	}

	if logPath != "" {
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		provider.logFile = file
		provider.logger = log.New(file, "[SMS] ", log.LstdFlags)
	} else {
		provider.logger = log.New(os.Stdout, "[SMS MOCK] ", log.LstdFlags)
	}

	return provider, nil
}

// Send logs the SMS message instead of sending it.
func (p *MockProvider) Send(ctx context.Context, msg Message) (*SendResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Generate a mock message ID
	messageID := uuid.New().String()

	sentMsg := SentMessage{
		ID:     messageID,
		To:     msg.To,
		Body:   msg.Body,
		From:   msg.From,
		SentAt: time.Now(),
	}

	p.messages = append(p.messages, sentMsg)

	// Log the message
	p.logger.Printf("SMS to %s: %s (ID: %s)", msg.To, msg.Body, messageID)

	return &SendResult{
		MessageID: messageID,
		Status:    "mock_sent",
	}, nil
}

// Name returns the provider name.
func (p *MockProvider) Name() string {
	return "mock"
}

// IsReady always returns true for the mock provider.
func (p *MockProvider) IsReady() bool {
	return true
}

// GetSentMessages returns all messages sent through this mock provider.
// Useful for testing.
func (p *MockProvider) GetSentMessages() []SentMessage {
	p.mu.Lock()
	defer p.mu.Unlock()

	result := make([]SentMessage, len(p.messages))
	copy(result, p.messages)
	return result
}

// GetLastMessage returns the most recently sent message.
// Useful for testing.
func (p *MockProvider) GetLastMessage() *SentMessage {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.messages) == 0 {
		return nil
	}
	return &p.messages[len(p.messages)-1]
}

// ClearMessages clears all stored messages.
// Useful for testing.
func (p *MockProvider) ClearMessages() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.messages = make([]SentMessage, 0)
}

// Close closes the log file if one was opened.
func (p *MockProvider) Close() error {
	if p.logFile != nil {
		return p.logFile.Close()
	}
	return nil
}

// Ensure MockProvider implements Provider interface.
var _ Provider = (*MockProvider)(nil)
