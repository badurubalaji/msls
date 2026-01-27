// Package logger provides structured logging using Zap with request ID support.
package logger

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// RequestIDKey is the context key for request IDs.
	RequestIDKey contextKey = "request_id"
)

// Logger wraps zap.Logger to provide additional functionality.
type Logger struct {
	*zap.Logger
}

var (
	// globalLogger is the default logger instance.
	globalLogger *Logger
)

// Config holds logger configuration options.
type Config struct {
	Level  string
	Format string
}

// New creates a new Logger instance with the given configuration.
func New(cfg Config) (*Logger, error) {
	level := parseLevel(cfg.Level)

	var encoder zapcore.Encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if strings.ToLower(cfg.Format) == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	logger := &Logger{Logger: zapLogger}
	globalLogger = logger

	return logger, nil
}

// parseLevel converts a string log level to zapcore.Level.
func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Global returns the global logger instance.
func Global() *Logger {
	if globalLogger == nil {
		// Create a default logger if none exists
		logger, _ := New(Config{Level: "info", Format: "json"})
		globalLogger = logger
	}
	return globalLogger
}

// WithRequestID returns a new logger with the request ID field.
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{Logger: l.With(zap.String("request_id", requestID))}
}

// WithContext returns a new logger with fields extracted from the context.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return l.WithRequestID(requestID)
	}
	return l
}

// WithField returns a new logger with an additional field.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{Logger: l.With(zap.Any(key, value))}
}

// WithFields returns a new logger with additional fields.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{Logger: l.With(zapFields...)}
}

// WithError returns a new logger with an error field.
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.With(zap.Error(err))}
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Debug logs a debug message.
func Debug(msg string, fields ...zap.Field) {
	Global().Debug(msg, fields...)
}

// Info logs an info message.
func Info(msg string, fields ...zap.Field) {
	Global().Info(msg, fields...)
}

// Warn logs a warning message.
func Warn(msg string, fields ...zap.Field) {
	Global().Warn(msg, fields...)
}

// Error logs an error message.
func Error(msg string, fields ...zap.Field) {
	Global().Error(msg, fields...)
}

// Fatal logs a fatal message and exits.
func Fatal(msg string, fields ...zap.Field) {
	Global().Fatal(msg, fields...)
}

// ContextWithRequestID returns a new context with the request ID.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// RequestIDFromContext extracts the request ID from the context.
func RequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
