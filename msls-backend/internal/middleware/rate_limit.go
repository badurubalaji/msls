// Package middleware provides HTTP middleware components for the Gin framework.
package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig holds configuration for rate limiting middleware.
type RateLimitConfig struct {
	// Rate is the number of requests allowed per interval.
	Rate int

	// Interval is the time window for rate limiting.
	Interval time.Duration

	// BurstSize is the maximum burst size (bucket capacity).
	BurstSize int

	// KeyFunc extracts the rate limiting key from the request.
	// Default: client IP address.
	KeyFunc func(c *gin.Context) string

	// ExcludedPaths lists paths to exclude from rate limiting.
	ExcludedPaths []string

	// OnLimitReached is an optional callback when rate limit is exceeded.
	OnLimitReached func(c *gin.Context)
}

// DefaultRateLimitConfig returns the default rate limit configuration.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Rate:           200,
		Interval:       time.Minute,
		BurstSize:      50,
		KeyFunc:        nil,
		ExcludedPaths:  []string{"/health", "/ready", "/metrics"},
		OnLimitReached: nil,
	}
}

// tokenBucket implements the token bucket algorithm.
type tokenBucket struct {
	tokens     float64
	capacity   float64
	rate       float64 // tokens per second
	lastUpdate time.Time
	mu         sync.Mutex
}

// newTokenBucket creates a new token bucket.
func newTokenBucket(capacity int, rate float64) *tokenBucket {
	return &tokenBucket{
		tokens:     float64(capacity),
		capacity:   float64(capacity),
		rate:       rate,
		lastUpdate: time.Now(),
	}
}

// take attempts to take n tokens from the bucket.
// Returns true if successful, false if rate limited.
func (b *tokenBucket) take(n float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastUpdate).Seconds()
	b.lastUpdate = now

	// Add tokens based on elapsed time
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}

	// Check if we have enough tokens
	if b.tokens < n {
		return false
	}

	b.tokens -= n
	return true
}

// remaining returns the current number of tokens.
func (b *tokenBucket) remaining() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return int(b.tokens)
}

// rateLimiter manages rate limiting for multiple clients.
type rateLimiter struct {
	buckets   map[string]*tokenBucket
	capacity  int
	rate      float64
	mu        sync.RWMutex
	cleanupMu sync.Mutex
	lastClean time.Time
}

// newRateLimiter creates a new rate limiter.
func newRateLimiter(capacity int, rate float64) *rateLimiter {
	return &rateLimiter{
		buckets:   make(map[string]*tokenBucket),
		capacity:  capacity,
		rate:      rate,
		lastClean: time.Now(),
	}
}

// getBucket returns the token bucket for a key, creating one if necessary.
func (r *rateLimiter) getBucket(key string) *tokenBucket {
	r.mu.RLock()
	bucket, exists := r.buckets[key]
	r.mu.RUnlock()

	if exists {
		return bucket
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock
	if bucket, exists = r.buckets[key]; exists {
		return bucket
	}

	bucket = newTokenBucket(r.capacity, r.rate)
	r.buckets[key] = bucket

	// Periodic cleanup of old buckets
	r.maybeCleanup()

	return bucket
}

// maybeCleanup removes stale buckets periodically.
func (r *rateLimiter) maybeCleanup() {
	// Only cleanup every 5 minutes
	if time.Since(r.lastClean) < 5*time.Minute {
		return
	}

	r.cleanupMu.Lock()
	defer r.cleanupMu.Unlock()

	// Double-check after acquiring lock
	if time.Since(r.lastClean) < 5*time.Minute {
		return
	}

	r.lastClean = time.Now()
	threshold := time.Now().Add(-10 * time.Minute)

	for key, bucket := range r.buckets {
		bucket.mu.Lock()
		if bucket.lastUpdate.Before(threshold) {
			delete(r.buckets, key)
		}
		bucket.mu.Unlock()
	}
}

// take attempts to take a token for the given key.
func (r *rateLimiter) take(key string) bool {
	return r.getBucket(key).take(1)
}

// remaining returns the remaining tokens for a key.
func (r *rateLimiter) remaining(key string) int {
	return r.getBucket(key).remaining()
}

// RateLimit returns a middleware that implements token bucket rate limiting.
// It returns 429 Too Many Requests when the rate limit is exceeded.
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	// Calculate rate (tokens per second)
	rate := float64(config.Rate) / config.Interval.Seconds()

	// Create rate limiter
	limiter := newRateLimiter(config.BurstSize, rate)

	// Build excluded paths map
	excludedPaths := make(map[string]bool)
	for _, path := range config.ExcludedPaths {
		excludedPaths[path] = true
	}

	// Default key function: client IP
	keyFunc := config.KeyFunc
	if keyFunc == nil {
		keyFunc = func(c *gin.Context) string {
			return c.ClientIP()
		}
	}

	return func(c *gin.Context) {
		// Skip rate limiting for excluded paths
		if excludedPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		key := keyFunc(c)

		// Attempt to take a token
		if !limiter.take(key) {
			// Rate limit exceeded
			retryAfter := int(config.Interval.Seconds())

			// Set headers
			c.Header("X-RateLimit-Limit", strconv.Itoa(config.Rate))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Interval).Unix(), 10))
			c.Header("Retry-After", strconv.Itoa(retryAfter))

			// Call optional callback
			if config.OnLimitReached != nil {
				config.OnLimitReached(c)
			}

			// Return 429 with RFC 7807 format
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"type":     "https://httpstatuses.com/429",
				"title":    "Too Many Requests",
				"status":   http.StatusTooManyRequests,
				"detail":   "Rate limit exceeded. Please try again later.",
				"instance": c.Request.URL.Path,
				"extensions": gin.H{
					"retry_after": retryAfter,
				},
			})
			return
		}

		// Set rate limit headers for successful requests
		remaining := limiter.remaining(key)
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Rate))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Interval).Unix(), 10))

		c.Next()
	}
}

// RateLimitDefault returns a middleware with default configuration.
func RateLimitDefault() gin.HandlerFunc {
	return RateLimit(DefaultRateLimitConfig())
}

// RateLimitPerMinute returns a middleware with the specified requests per minute.
func RateLimitPerMinute(rate int) gin.HandlerFunc {
	config := DefaultRateLimitConfig()
	config.Rate = rate
	config.Interval = time.Minute
	return RateLimit(config)
}

// RateLimitPerSecond returns a middleware with the specified requests per second.
func RateLimitPerSecond(rate int) gin.HandlerFunc {
	config := DefaultRateLimitConfig()
	config.Rate = rate
	config.Interval = time.Second
	return RateLimit(config)
}

// RateLimitByTenant returns a middleware that rate limits by tenant ID.
func RateLimitByTenant(rate int, interval time.Duration) gin.HandlerFunc {
	config := DefaultRateLimitConfig()
	config.Rate = rate
	config.Interval = interval
	config.KeyFunc = func(c *gin.Context) string {
		tenantID := GetTenantID(c)
		if tenantID == "" {
			return c.ClientIP()
		}
		return tenantID
	}
	return RateLimit(config)
}

// RateLimitByUser returns a middleware that rate limits by user ID (from context).
func RateLimitByUser(rate int, interval time.Duration, userIDKey string) gin.HandlerFunc {
	config := DefaultRateLimitConfig()
	config.Rate = rate
	config.Interval = interval
	config.KeyFunc = func(c *gin.Context) string {
		if userID, exists := c.Get(userIDKey); exists {
			if id, ok := userID.(string); ok && id != "" {
				return id
			}
		}
		return c.ClientIP()
	}
	return RateLimit(config)
}
