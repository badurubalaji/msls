// Package database provides database connection and management utilities.
package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds the database connection configuration.
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	LogLevel        logger.LogLevel
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "",
		DBName:          "msls",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
		LogLevel:        logger.Warn,
	}
}

// DSN returns the PostgreSQL connection string.
func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// DatabaseURL returns the PostgreSQL connection URL format.
func (c Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

// Connection wraps a GORM database connection with additional functionality.
type Connection struct {
	db *gorm.DB
}

// New creates a new database connection using the provided configuration.
func New(cfg Config) (*Connection, error) {
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(cfg.LogLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return &Connection{db: db}, nil
}

// NewFromDSN creates a new database connection using a DSN string.
func NewFromDSN(dsn string, logLevel logger.LogLevel) (*Connection, error) {
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Connection{db: db}, nil
}

// DB returns the underlying GORM database instance.
func (c *Connection) DB() *gorm.DB {
	return c.db
}

// Close closes the database connection.
func (c *Connection) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Ping verifies the database connection is alive.
func (c *Connection) Ping(ctx context.Context) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

// SetPoolConfig updates the connection pool configuration.
func (c *Connection) SetPoolConfig(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	return nil
}

// Stats returns database connection pool statistics.
func (c *Connection) Stats() (Stats, error) {
	sqlDB, err := c.db.DB()
	if err != nil {
		return Stats{}, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	dbStats := sqlDB.Stats()
	return Stats{
		MaxOpenConnections: dbStats.MaxOpenConnections,
		OpenConnections:    dbStats.OpenConnections,
		InUse:              dbStats.InUse,
		Idle:               dbStats.Idle,
		WaitCount:          dbStats.WaitCount,
		WaitDuration:       dbStats.WaitDuration,
		MaxIdleClosed:      dbStats.MaxIdleClosed,
		MaxIdleTimeClosed:  dbStats.MaxIdleTimeClosed,
		MaxLifetimeClosed:  dbStats.MaxLifetimeClosed,
	}, nil
}

// Stats represents database connection pool statistics.
type Stats struct {
	MaxOpenConnections int
	OpenConnections    int
	InUse              int
	Idle               int
	WaitCount          int64
	WaitDuration       time.Duration
	MaxIdleClosed      int64
	MaxIdleTimeClosed  int64
	MaxLifetimeClosed  int64
}

// HealthCheck performs a comprehensive health check on the database connection.
func (c *Connection) HealthCheck(ctx context.Context) error {
	// Check if we can ping the database
	if err := c.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Check if we can execute a simple query
	var result int
	if err := c.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	return nil
}
