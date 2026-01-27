// Package config provides application configuration management using Viper.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration values.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MinIO    MinIOConfig
	Log      LogConfig
	App      AppConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds PostgreSQL database configuration.
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DSN returns the PostgreSQL connection string.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Addr returns the Redis address in host:port format.
func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig holds JWT authentication configuration.
type JWTConfig struct {
	Secret           string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
	Issuer           string
}

// MinIOConfig holds MinIO object storage configuration.
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level  string
	Format string
}

// AppConfig holds general application configuration.
type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
}

// IsDevelopment returns true if the application is running in development mode.
func (c AppConfig) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction returns true if the application is running in production mode.
func (c AppConfig) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// Load reads configuration from environment variables and returns a Config struct.
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Read from environment variables
	v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind environment variables explicitly
	bindEnvVars(v)

	cfg := &Config{
		App: AppConfig{
			Name:        v.GetString("APP_NAME"),
			Environment: v.GetString("APP_ENV"),
			Debug:       v.GetBool("APP_DEBUG"),
		},
		Server: ServerConfig{
			Host:         v.GetString("SERVER_HOST"),
			Port:         v.GetInt("SERVER_PORT"),
			ReadTimeout:  v.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  v.GetDuration("SERVER_IDLE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetInt("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			Name:            v.GetString("DB_NAME"),
			SSLMode:         v.GetString("DB_SSLMODE"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetInt("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:           v.GetString("JWT_SECRET"),
			AccessExpiresIn:  v.GetDuration("JWT_ACCESS_EXPIRES_IN"),
			RefreshExpiresIn: v.GetDuration("JWT_REFRESH_EXPIRES_IN"),
			Issuer:           v.GetString("JWT_ISSUER"),
		},
		MinIO: MinIOConfig{
			Endpoint:        v.GetString("MINIO_ENDPOINT"),
			AccessKeyID:     v.GetString("MINIO_ACCESS_KEY_ID"),
			SecretAccessKey: v.GetString("MINIO_SECRET_ACCESS_KEY"),
			UseSSL:          v.GetBool("MINIO_USE_SSL"),
			BucketName:      v.GetString("MINIO_BUCKET_NAME"),
		},
		Log: LogConfig{
			Level:  v.GetString("LOG_LEVEL"),
			Format: v.GetString("LOG_FORMAT"),
		},
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("APP_NAME", "msls-backend")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_DEBUG", true)

	// Server defaults
	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("SERVER_READ_TIMEOUT", "15s")
	v.SetDefault("SERVER_WRITE_TIMEOUT", "15s")
	v.SetDefault("SERVER_IDLE_TIMEOUT", "60s")

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "msls")
	v.SetDefault("DB_PASSWORD", "msls_password")
	v.SetDefault("DB_NAME", "msls")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	// Redis defaults
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	// JWT defaults
	v.SetDefault("JWT_SECRET", "change-me-in-production")
	v.SetDefault("JWT_ACCESS_EXPIRES_IN", "15m")
	v.SetDefault("JWT_REFRESH_EXPIRES_IN", "168h") // 7 days
	v.SetDefault("JWT_ISSUER", "msls-backend")

	// MinIO defaults
	v.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	v.SetDefault("MINIO_ACCESS_KEY_ID", "minioadmin")
	v.SetDefault("MINIO_SECRET_ACCESS_KEY", "minioadmin")
	v.SetDefault("MINIO_USE_SSL", false)
	v.SetDefault("MINIO_BUCKET_NAME", "msls")

	// Log defaults
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")
}

func bindEnvVars(v *viper.Viper) {
	envVars := []string{
		"APP_NAME", "APP_ENV", "APP_DEBUG",
		"SERVER_HOST", "SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "DB_CONN_MAX_LIFETIME",
		"REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD", "REDIS_DB",
		"JWT_SECRET", "JWT_ACCESS_EXPIRES_IN", "JWT_REFRESH_EXPIRES_IN", "JWT_ISSUER",
		"MINIO_ENDPOINT", "MINIO_ACCESS_KEY_ID", "MINIO_SECRET_ACCESS_KEY", "MINIO_USE_SSL", "MINIO_BUCKET_NAME",
		"LOG_LEVEL", "LOG_FORMAT",
	}

	for _, env := range envVars {
		_ = v.BindEnv(env)
	}
}
